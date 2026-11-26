package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/yaml"
)

// TestBGPVIPPeersJSONCapacity checks that every gated ControllerConfig schema
// accepts the largest transport MCO can derive from the bounded BGPVIPConfig.
// MCO's bgpVIPPeersJSONFromCR adds VIPs, maps hostnames to peer lists, and
// serializes passwordSecret references (including the namespace), not passwords.
func TestBGPVIPPeersJSONCapacity(t *testing.T) {
	bgpSchemas := readBGPSchemas(t, "../v1alpha1/zz_generated.crd-manifests/0000_80_machine-config_01_bgpvipconfigs.crd.yaml")
	bgpSpec := bgpSchemas[0].Properties["spec"]
	defaultPeers := bgpSpec.Properties["defaultPeers"]
	hostOverrides := bgpSpec.Properties["hostOverrides"]
	peerSchema := defaultPeers.Items.Schema
	hostnameSchema := hostOverrides.Items.Schema.Properties["hostname"]
	secretSchema := peerSchema.Properties["passwordSecret"].Properties["name"]
	communitySchema := bgpSpec.Properties["communities"]

	// Use the API length bounds, even where canonical IPs are shorter, so this
	// is a conservative upper bound. All serialized strings are ASCII without
	// JSON escapes. The field names match MCO's BGPVIPPeersConfig/BGPVIPPeer.
	peer := map[string]interface{}{
		"peerAddress": strings.Repeat("f", int(*peerSchema.Properties["peerAddress"].MaxLength)),
		"peerASN":     int64(4294967295),
		"passwordSecret": map[string]interface{}{
			"name":      strings.Repeat("s", int(*secretSchema.MaxLength)),
			"namespace": "openshift-config",
		},
		"port":          65535,
		"bfdEnabled":    "true",
		"ebgpMultiHop":  "true",
		"holdTime":      "65535",
		"keepaliveTime": "65535",
	}
	peers := make([]interface{}, int(*defaultPeers.MaxItems))
	for i := range peers {
		peers[i] = peer
	}
	overridePeers := make([]interface{}, int(*hostOverrides.Items.Schema.Properties["peers"].MaxItems))
	for i := range overridePeers {
		overridePeers[i] = peer
	}
	overrides := make(map[string]interface{}, int(*hostOverrides.MaxItems))
	for i := range int(*hostOverrides.MaxItems) {
		prefix := fmt.Sprintf("%d-", i)
		hostname := prefix + strings.Repeat("h", int(*hostnameSchema.MaxLength)-len(prefix))
		overrides[hostname] = overridePeers
	}
	communities := make([]string, int(*communitySchema.MaxItems))
	for i := range communities {
		communities[i] = strings.Repeat("4", int(*communitySchema.Items.Schema.MaxLength))
	}
	payload, err := json.Marshal(map[string]interface{}{
		"localASN":      int64(4294967295),
		"defaultPeers":  peers,
		"communities":   communities,
		"apiVIPs":       []string{strings.Repeat("f", 45), strings.Repeat("f", 45)},
		"ingressVIPs":   []string{strings.Repeat("f", 45), strings.Repeat("f", 45)},
		"hostOverrides": overrides,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("conservative maximum MCO transport: %d characters", len(payload))

	files, err := filepath.Glob("zz_generated.crd-manifests/*_controllerconfigs-*.crd.yaml")
	if err != nil {
		t.Fatal(err)
	}
	gated := 0
	for _, file := range files {
		for _, schema := range readBGPSchemas(t, file) {
			field, present := schema.Properties["spec"].Properties["bgpVIPPeersJSON"]
			if !present {
				continue
			}
			gated++
			if field.MaxLength == nil {
				t.Errorf("%s: transport must have a bounded length", file)
			} else if int64(len(payload)) > *field.MaxLength {
				t.Errorf("%s: maximum transport needs %d characters, schema allows %d", file, len(payload), *field.MaxLength)
			}
		}
	}
	if gated == 0 {
		t.Fatal("no BGP-enabled ControllerConfig schemas tested")
	}
}

func readBGPSchemas(t *testing.T, file string) []spec.Schema {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var crd struct {
		Spec struct {
			Versions []struct {
				Schema struct {
					OpenAPIV3Schema spec.Schema `json:"openAPIV3Schema"`
				} `json:"schema"`
			} `json:"versions"`
		} `json:"spec"`
	}
	if err := yaml.Unmarshal(data, &crd); err != nil {
		t.Fatal(err)
	}
	if len(crd.Spec.Versions) == 0 {
		t.Fatalf("%s: no CRD versions found", file)
	}
	schemas := make([]spec.Schema, 0, len(crd.Spec.Versions))
	for _, version := range crd.Spec.Versions {
		schemas = append(schemas, version.Schema.OpenAPIV3Schema)
	}
	return schemas
}
