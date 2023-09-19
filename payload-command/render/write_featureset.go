package render

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	for featureSetName := range configv1.FeatureSets {
		featureGates := &configv1.FeatureGate{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster",
			},
			Spec: configv1.FeatureGateSpec{
				FeatureGateSelection: configv1.FeatureGateSelection{
					FeatureSet: featureSetName,
				},
			},
		}

		currentDetails, err := FeaturesGateDetailsFromFeatureSets(configv1.FeatureSets, featureGates, o.PayloadVersion)
		if err != nil {
			return fmt.Errorf("error determining FeatureGates: %w", err)
		}
		featureGates.Status.FeatureGates = []configv1.FeatureGateDetails{*currentDetails}

		featureGateOutBytes := writeFeatureGateV1OrDie(featureGates)
		featureSetFileName := fmt.Sprintf("featureGate-%s.yaml", featureSetName)
		if len(featureSetName) == 0 {
			featureSetFileName = fmt.Sprintf("featureGate-%s.yaml", "Default")
		}

		destFile := filepath.Join(o.AssetOutputDir, featureSetFileName)
		if err := os.WriteFile(destFile, []byte(featureGateOutBytes), 0644); err != nil {
			return fmt.Errorf("error writing FeatureGate manifest: %w", err)
		}
	}

	return nil
}
