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
			actual: newDefaultFeatures().without("PodSecurity").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
				},
				Disabled: []string{
					"CSIMigrationAWS",
					"CSIMigrationGCE",
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
					"PodSecurity",
				},
			},
		},
		{
			name:   "enable-existing",
			actual: newDefaultFeatures().with("CSIMigrationAWS").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
					"PodSecurity",
					"CSIMigrationAWS",
				},
				Disabled: []string{
					"CSIMigrationGCE",
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
				},
			},
		},
		{
			name:   "disable-more",
			actual: newDefaultFeatures().without("PodSecurity", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
				},
				Disabled: []string{
					"CSIMigrationAWS",
					"CSIMigrationGCE",
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
					"PodSecurity",
					"other",
				},
			},
		},
		{
			name:   "enable-more",
			actual: newDefaultFeatures().with("CSIMigrationAWS", "other").toFeatures(),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []string{
					"APIPriorityAndFairness",
					"RotateKubeletServerCertificate",
					"DownwardAPIHugePages",
					"PodSecurity",
					"CSIMigrationAWS",
					"other",
				},
				Disabled: []string{
					"CSIMigrationGCE",
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
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
