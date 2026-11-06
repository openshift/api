package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestClusterAPICRDFeatureGate(t *testing.T) {
	for _, profile := range []string{"Hypershift", "SelfManagedHA"} {
		t.Run(profile, func(t *testing.T) {
			filename := "0000_30_cluster-api_01_clusterapis-" + profile + ".crd.yaml"
			data, err := os.ReadFile(filepath.Join("zz_generated.crd-manifests", filename))
			if err != nil {
				t.Fatal(err)
			}
			var crd metav1.PartialObjectMetadata
			if err := yaml.Unmarshal(data, &crd); err != nil {
				t.Fatal(err)
			}
			if got := crd.Annotations["release.openshift.io/feature-gate"]; got != "ClusterAPIOperator" {
				t.Errorf("release feature gate = %q, want ClusterAPIOperator", got)
			}
		})
	}
}
