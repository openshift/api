package render

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
)

func TestRenderCustomNoUpgradeFeatureGate(t *testing.T) {
	const (
		payloadVersion      = "4.22.0"
		payloadMajorVersion = uint64(4)
	)

	defaultFeatureGates := features.FeatureSets(payloadMajorVersion, features.SelfManaged, configv1.Default)
	if len(defaultFeatureGates.Enabled) == 0 || len(defaultFeatureGates.Disabled) == 0 {
		t.Fatal("default feature gates must include enabled and disabled gates for override testing")
	}
	forceEnabled := defaultFeatureGates.Disabled[0].FeatureGateAttributes.Name
	forceDisabled := defaultFeatureGates.Enabled[0].FeatureGateAttributes.Name

	tests := []struct {
		name             string
		custom           *configv1.CustomFeatureGates
		wantFeatureGates configv1.FeatureGateDetails
		wantErr          string
	}{
		{
			name:             "omitted customNoUpgrade uses defaults",
			wantFeatureGates: *FeaturesGateDetailsFromFeatureSets(defaultFeatureGates, payloadVersion),
		},
		{
			name:             "empty customNoUpgrade uses defaults",
			custom:           &configv1.CustomFeatureGates{},
			wantFeatureGates: *FeaturesGateDetailsFromFeatureSets(defaultFeatureGates, payloadVersion),
		},
		{
			name: "explicit overrides replace defaults",
			custom: &configv1.CustomFeatureGates{
				Enabled:  []configv1.FeatureGateName{forceEnabled},
				Disabled: []configv1.FeatureGateName{forceDisabled},
			},
			wantFeatureGates: featureGateDetailsWithOverrides(defaultFeatureGates, payloadVersion, forceEnabled, forceDisabled),
		},
		{
			name: "conflicting overrides are rejected",
			custom: &configv1.CustomFeatureGates{
				Enabled:  []configv1.FeatureGateName{"Example"},
				Disabled: []configv1.FeatureGateName{"Example"},
			},
			wantErr: `trying to enable and disable "Example"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &configv1.FeatureGate{
				Spec: configv1.FeatureGateSpec{
					FeatureGateSelection: configv1.FeatureGateSelection{
						FeatureSet:      configv1.CustomNoUpgrade,
						CustomNoUpgrade: tt.custom,
					},
				},
			}

			got, err := renderCustomNoUpgradeFeatureGate(in, features.SelfManaged, payloadVersion, payloadMajorVersion)
			if len(tt.wantErr) > 0 {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("renderCustomNoUpgradeFeatureGate() error = %v, want error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("renderCustomNoUpgradeFeatureGate() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tt.wantFeatureGates, got.Status.FeatureGates[0]) {
				t.Errorf("renderCustomNoUpgradeFeatureGate() feature gates = %#v, want %#v", got.Status.FeatureGates[0], tt.wantFeatureGates)
			}
		})
	}
}

func featureGateDetailsWithOverrides(defaults *features.FeatureGateEnabledDisabled, payloadVersion string, forceEnabled, forceDisabled configv1.FeatureGateName) configv1.FeatureGateDetails {
	details := *FeaturesGateDetailsFromFeatureSets(defaults, payloadVersion)
	details.Enabled = replaceFeatureGate(details.Enabled, forceDisabled, forceEnabled)
	details.Disabled = replaceFeatureGate(details.Disabled, forceEnabled, forceDisabled)
	return details
}

func replaceFeatureGate(in []configv1.FeatureGateAttributes, remove, add configv1.FeatureGateName) []configv1.FeatureGateAttributes {
	out := make([]configv1.FeatureGateAttributes, 0, len(in))
	for _, gate := range in {
		if gate.Name != remove {
			out = append(out, gate)
		}
	}
	out = append(out, configv1.FeatureGateAttributes{Name: add})
	sort.Sort(byName(out))
	return out
}
