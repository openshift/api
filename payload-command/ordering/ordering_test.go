package ordering

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// This test enforces a CVO apply-ordering invariant for the manifests shipped by
// this repository: every empty custom resource (CR) in
// payload-command/empty-resources whose CustomResourceDefinition (CRD) is also
// shipped in payload-manifests/crds must have that CRD applied before it.
//
// The Cluster Version Operator applies manifests ordered by run-level (the
// 0000_NN_ filename prefix) and, within a run-level, by byte-order of the
// filename; it will not advance to a higher run-level until lower ones complete.
// A CR that sorts before its CRD therefore deadlocks an in-progress update once
// the CR's feature gate is enabled (see OCPBUGS-99266 / OCPBUGS-113639).
//
// A CRD counts as "applied before" a CR when any of the following hold:
//   - the CRD is release.openshift.io/bootstrap-required (created at bootstrap,
//     before the CVO's ordered apply runs);
//   - the CRD is at a lower run-level than the CR;
//   - the CRD is at the same run-level and its filename sorts before the CR's.

type manifest struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name        string            `json:"name"`
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Group string `json:"group"`
		Names struct {
			Kind string `json:"kind"`
		} `json:"names"`
	} `json:"spec"`

	// filled in by the loader
	filename string
}

type groupKind struct{ group, kind string }

// runLevel returns the "0000_NN" prefix that the CVO uses as a run-level barrier.
func runLevel(filename string) string {
	parts := strings.SplitN(filename, "_", 3)
	if len(parts) >= 2 {
		return parts[0] + "_" + parts[1]
	}
	return parts[0]
}

// groupOf extracts the API group from an apiVersion ("" for core).
func groupOf(apiVersion string) string {
	if i := strings.Index(apiVersion, "/"); i >= 0 {
		return apiVersion[:i]
	}
	return ""
}

func (m manifest) isBootstrapRequired() bool {
	return strings.EqualFold(m.Metadata.Annotations["release.openshift.io/bootstrap-required"], "true")
}

// appliesBefore reports whether crd is guaranteed to be applied (and, for a CRD,
// established) before cr under CVO ordering.
func appliesBefore(crd, cr manifest) bool {
	if crd.isBootstrapRequired() {
		return true
	}
	crdRL, crRL := runLevel(crd.filename), runLevel(cr.filename)
	if crdRL != crRL {
		return crdRL < crRL
	}
	// Same run-level: byte-order of the filename decides. Go string comparison
	// is byte-wise, matching the CVO's C-locale sort.
	return crd.filename < cr.filename
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate repo root (go.mod not found)")
		}
		dir = parent
	}
}

func loadManifests(t *testing.T, dir string) []manifest {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	var out []manifest
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if ext := filepath.Ext(e.Name()); ext != ".yaml" && ext != ".yml" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		reader := utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(raw)))
		for {
			doc, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("parse %s: %v", e.Name(), err)
			}
			if strings.TrimSpace(string(doc)) == "" {
				continue
			}
			var m manifest
			if err := yaml.Unmarshal(doc, &m); err != nil {
				t.Fatalf("parse %s: %v", e.Name(), err)
			}
			if m.Kind == "" {
				continue
			}
			m.filename = e.Name()
			out = append(out, m)
		}
	}
	return out
}

func TestEmptyResourceCRDOrdering(t *testing.T) {
	root := repoRoot(t)

	// Build the registry of CRDs shipped in the payload from this repo.
	crds := map[groupKind][]manifest{}
	for _, m := range loadManifests(t, filepath.Join(root, "payload-manifests", "crds")) {
		if m.Kind != "CustomResourceDefinition" {
			continue
		}
		if m.Spec.Group == "" || m.Spec.Names.Kind == "" {
			continue
		}
		gk := groupKind{m.Spec.Group, m.Spec.Names.Kind}
		crds[gk] = append(crds[gk], m)
	}

	// Every empty CR whose CRD ships here must be ordered after that CRD.
	crs := loadManifests(t, filepath.Join(root, "payload-command", "empty-resources"))
	if len(crs) == 0 {
		t.Fatal("no empty-resource manifests found")
	}

	checked := 0
	for _, cr := range crs {
		gk := groupKind{groupOf(cr.APIVersion), cr.Kind}
		variants, ok := crds[gk]
		if !ok {
			// CRD is provided elsewhere (e.g. a bootstrap-only CRD not in this
			// repo, or another component). Nothing to enforce here.
			continue
		}
		checked++
		for _, crd := range variants {
			if !appliesBefore(crd, cr) {
				t.Errorf("CR %q (%s.%s) is not ordered after its CRD:\n"+
					"  CR : %s\n"+
					"  CRD: %s\n"+
					"The CR sorts before its CRD under CVO ordering, which can deadlock an\n"+
					"in-progress update. Place the CR at a run-level/filename that sorts after\n"+
					"the CRD (or make the CRD bootstrap-required).",
					cr.Metadata.Name, gk.kind, gk.group, cr.filename, crd.filename)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no CR/CRD pairs were checked; test wiring is likely broken")
	}
	t.Logf("verified ordering for %d empty CR(s) against in-repo CRDs", checked)
}
