package render

import (
	"flag"
	"fmt"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
)

var (
	clusterProfileToShortName = map[configv1.ClusterProfileName]string{
		configv1.Hypershift:  "Hypershift",
		configv1.SelfManaged: "SelfManagedHA",
		configv1.SingleNode:  "SingleNode",
	}
)

// WriteFeatureSets holds values to drive the render command.
type WriteFeatureSets struct {
	PayloadVersion string
	AssetOutputDir string
}

func (o *WriteFeatureSets) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.PayloadVersion, "payload-version", o.PayloadVersion, "Version that will eventually be placed into ClusterOperator.status.  This normally comes from the CVO set via env var: OPERATOR_IMAGE_VERSION.")
	fs.StringVar(&o.AssetOutputDir, "asset-output-dir", o.AssetOutputDir, "Output path for rendered manifests.")
}

// Validate verifies the inputs.
func (o *WriteFeatureSets) Validate() error {
	return nil
}

// Complete fills in missing values before command execution.
func (o *WriteFeatureSets) Complete() error {
	return nil
}

// Run contains the logic of the render command.
func (o *WriteFeatureSets) Run() error {
	err := os.MkdirAll(o.AssetOutputDir, 0755)
	if err != nil {
		return err
	}

	statusByClusterProfileByFeatureSet := configv1.AllFeatureSets()
	for clusterProfile, byFeatureSet := range statusByClusterProfileByFeatureSet {
		for featureSetName, featureGateStatuses := range byFeatureSet {
			currentDetails := FeaturesGateDetailsFromFeatureSets(featureGateStatuses, o.PayloadVersion)

			featureGateInstance := &configv1.FeatureGate{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cluster",
					Annotations: map[string]string{
						string(clusterProfile): "true",
					},
				},
				Spec: configv1.FeatureGateSpec{
					FeatureGateSelection: configv1.FeatureGateSelection{
						FeatureSet: featureSetName,
					},
				},
				Status: configv1.FeatureGateStatus{
					FeatureGates: []configv1.FeatureGateDetails{
						*currentDetails,
					},
				},
			}

			featureGateOutBytes := writeFeatureGateV1OrDie(featureGateInstance)
			featureSetFileName := fmt.Sprintf("featureGate-%s-%s.yaml", featureSetName, clusterProfileToShortName[clusterProfile])
			if len(featureSetName) == 0 {
				featureSetFileName = fmt.Sprintf("featureGate-%s-%s.yaml", "Default", clusterProfileToShortName[clusterProfile])
			}

			destFile := filepath.Join(o.AssetOutputDir, featureSetFileName)
			if err := os.WriteFile(destFile, []byte(featureGateOutBytes), 0644); err != nil {
				return fmt.Errorf("error writing FeatureGate manifest: %w", err)
			}

			// for compatibility during the transition, we'll copy the old, invalid featuregates
			legacyFilename := ""
			switch {
			case len(featureSetName) == 0 && clusterProfile == configv1.SelfManaged:
				legacyFilename = "featureGate-Default.yaml"
			case featureSetName == configv1.TechPreviewNoUpgrade && clusterProfile == configv1.SelfManaged:
				legacyFilename = "featureGate-TechPreviewNoUpgrade.yaml"
			}
			if len(legacyFilename) > 0 {
				legacyFeatureGateInstance := featureGateInstance.DeepCopy()
				delete(legacyFeatureGateInstance.Annotations, "include.release.openshift.io/self-managed-high-availability")
				legacyFeatureGateBytes := writeFeatureGateV1OrDie(legacyFeatureGateInstance)
				legacyFile := filepath.Join(o.AssetOutputDir, legacyFilename)
				if err := os.WriteFile(legacyFile, []byte(legacyFeatureGateBytes), 0644); err != nil {
					return fmt.Errorf("error writing FeatureGate manifest: %w", err)
				}
			}
		}
	}

	return nil
}
