package v1

import (
	"reflect"
	"testing"
)

func TestFeatureBuilder(t *testing.T) {

	var testStart = &FeatureGateEnabledDisabled{
		Enabled:  []string{"alpha", "bravo"},
		Disabled: []string{"charlie", "delta"},
	}

	tests := []struct {
		name     string
		actual   *FeatureGateEnabledDisabled
		expected *FeatureGateEnabledDisabled
	}{
		{
			name:     "nothing",
			actual:   newDefaultFeatures(defaultFeatures).toFeatures(),
			expected: defaultFeatures,
		},
		{
			name:   "disable-existing",
			actual: newDefaultFeatures(testStart).without("alpha").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{"bravo"},
				Disabled: []string{
					"charlie",
					"delta",
					"alpha",
				},
			},
		},
		{
			name:   "enable-existing",
			actual: newDefaultFeatures(testStart).with("charlie").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"alpha",
					"bravo",
					"charlie",
				},
				Disabled: []string{
					"delta",
				},
			},
		},
		{
			name:   "disable-more",
			actual: newDefaultFeatures(testStart).without("alpha", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{"bravo"},
				Disabled: []string{
					"charlie",
					"delta",
					"alpha",
					"other",
				},
			},
		},
		{
			name:   "enable-more",
			actual: newDefaultFeatures(testStart).with("charlie", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"alpha",
					"bravo",
					"charlie",
					"other",
				},
				Disabled: []string{
					"delta",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.expected, tc.actual) {
				t.Error(tc.actual)
			}
		})
	}
}
