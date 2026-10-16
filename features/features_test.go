package features

import (
	"fmt"
	"testing"
)

func TestNewOLMOwnSingleNamespaceNotRegistered(t *testing.T) {
	t.Parallel()

	const removedFeatureGate = "NewOLMOwnSingleNamespace"
	if _, registered := allFeatureGates[removedFeatureGate]; registered {
		t.Errorf("%s must not be registered", removedFeatureGate)
	}

	for version, byClusterProfile := range AllFeatureSets() {
		for clusterProfile, byFeatureSet := range byClusterProfile {
			for featureSet, gates := range byFeatureSet {
				t.Run(fmt.Sprintf("%d/%s/%s", version, clusterProfile, featureSet), func(t *testing.T) {
					t.Parallel()

					for state, descriptions := range map[string][]FeatureGateDescription{
						"enabled":  gates.Enabled,
						"disabled": gates.Disabled,
					} {
						for _, gate := range descriptions {
							if gate.FeatureGateAttributes.Name == removedFeatureGate {
								t.Errorf("%s must not appear in %s gates", removedFeatureGate, state)
							}
						}
					}
				})
			}
		}
	}
}
