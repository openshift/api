package features

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
)

func TestIngressControllerResourceConfigurationFeatureSets(t *testing.T) {
	t.Parallel()
	for version, profiles := range AllFeatureSets() {
		for profile, featureSets := range profiles {
			for featureSet, gates := range featureSets {
				if featureSet != configv1.Default && featureSet != configv1.OKD && featureSet != configv1.TechPreviewNoUpgrade && featureSet != configv1.DevPreviewNoUpgrade {
					continue
				}
				enabled := false
				for _, gate := range gates.Enabled {
					if gate.FeatureGateAttributes.Name == FeatureGateIngressControllerResourceConfiguration {
						enabled = true
					}
				}
				wantEnabled := featureSet == configv1.TechPreviewNoUpgrade || featureSet == configv1.DevPreviewNoUpgrade
				if enabled != wantEnabled {
					t.Errorf("version=%d profile=%s featureSet=%s: enabled=%t, want %t", version, profile, featureSet, enabled, wantEnabled)
				}
			}
		}
	}
}
