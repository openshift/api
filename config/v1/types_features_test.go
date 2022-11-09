package v1

import (
	"reflect"
	"testing"
)

func TestFeatureBuilder(t *testing.T) {
	tests := []struct {
		name     string
		actual   *FeatureGateEnabledDisabled
		expected *FeatureGateEnabledDisabled
	}{
		{
			name:     "nothing",
			actual:   newDefaultFeatures().toFeatures(),
			expected: defaultFeatures,
		},
		{
			name:   "disable-existing",
			actual: newDefaultFeatures().without("APIPriorityAndFairness").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
				},
				Disabled: []string{
					"RetroactiveDefaultStorageClass",
					"APIPriorityAndFairness",
				},
			},
		},
		{
			name:   "enable-existing",
			actual: newDefaultFeatures().with("CSIMigrationAzureFile").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
					"CSIMigrationAzureFile",
				},
				Disabled: []string{
					"RetroactiveDefaultStorageClass",
				},
			},
		},
		{
			name:   "disable-more",
			actual: newDefaultFeatures().without("APIPriorityAndFairness", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
				},
				Disabled: []string{
					"RetroactiveDefaultStorageClass",
					"APIPriorityAndFairness",
					"other",
				},
			},
		},
		{
			name:   "enable-more",
			actual: newDefaultFeatures().with("CSIMigrationAzureFile", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
					"CSIMigrationAzureFile",
					"other",
				},
				Disabled: []string{
					"RetroactiveDefaultStorageClass",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.expected, tc.actual) {
				t.Errorf("\nExpected feature gates: \n Enabled: %s \n Disabled: %s \nBut got:\n Enabled: %v \n Disabled: %s\n", tc.expected.Enabled, tc.expected.Disabled, tc.actual.Enabled, tc.actual.Disabled)
			}
		})
	}
}
