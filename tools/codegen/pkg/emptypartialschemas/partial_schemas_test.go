package emptypartialschemas

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func testCRDInfo(crdName string) *CRDInfo {
	return &CRDInfo{
		CRDName:              crdName,
		TopLevelFeatureGates: []string{"ExampleGate"},
		GroupName:            "example.openshift.io",
		Version:              "v1",
		PluralName:           "examples",
		KindName:             "Example",
		Scope:                "Cluster",
		HasStatus:            true,
		FeatureGates:         []string{"ExampleGate", "OtherGate"},
		ShortNames:           []string{"ex"},
		ApprovedPRNumber:     "https://github.com/openshift/api/pull/1",
		FilenameRunLevel:     "0000_10",
		FilenameOperatorName: "config-operator",
		PrinterColumns: []apiextensionsv1.CustomResourceColumnDefinition{
			{Name: "Age", Type: "date", JSONPath: ".metadata.creationTimestamp"},
		},
		Annotations: map[string]string{"example.openshift.io/annotation": "value"},
	}
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

// TestWriteCRDInfoEntry covers the empty, single and multiple entry cases of the
// zz_generated.featuregated-crd-manifests.yaml report: the output ends with exactly
// one newline, contains no blank lines and decodes to the complete entries written.
func TestWriteCRDInfoEntry(t *testing.T) {
	tests := []struct {
		name     string
		crdNames []string
	}{
		{
			name:     "empty",
			crdNames: nil,
		},
		{
			name:     "single entry",
			crdNames: []string{"examples.example.openshift.io"},
		},
		{
			name:     "multiple entries",
			crdNames: []string{"bars.example.openshift.io", "examples.example.openshift.io", "foos.example.openshift.io"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			want := map[string]*CRDInfo{}
			for _, crdName := range tt.crdNames {
				want[crdName] = testCRDInfo(crdName)
				if err := writeCRDInfoEntry(&buf, crdName, testCRDInfo(crdName)); err != nil {
					t.Fatalf("writeCRDInfoEntry(%q) returned an unexpected error: %v", crdName, err)
				}
			}

			out := buf.Bytes()
			if len(tt.crdNames) == 0 {
				if len(out) != 0 {
					t.Fatalf("expected no output, got %q", out)
				}
				return
			}

			if !bytes.HasSuffix(out, []byte("\n")) || bytes.HasSuffix(out, []byte("\n\n")) {
				t.Errorf("expected output to end with exactly one newline, got %q", out)
			}
			if bytes.Contains(out, []byte("\n\n")) {
				t.Errorf("expected output without blank lines, got %q", out)
			}

			got := map[string]*CRDInfo{}
			if err := yaml.Unmarshal(out, &got); err != nil {
				t.Fatalf("output is not valid YAML: %v\n%s", err, out)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("decoded entries do not match the written entries:\ngot:  %#v\nwant: %#v", got, want)
			}
		})
	}
}

func TestWriteCRDInfoEntryWriterError(t *testing.T) {
	writeErr := errors.New("write failed")
	err := writeCRDInfoEntry(failingWriter{err: writeErr}, "examples.example.openshift.io", testCRDInfo("examples.example.openshift.io"))
	if !errors.Is(err, writeErr) {
		t.Fatalf("expected the writer error %v, got %v", writeErr, err)
	}
}
