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
					"OpenShiftPodSecurityAdmission",
				},
				Disabled: []string{
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
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
					"OpenShiftPodSecurityAdmission",
					"CSIMigrationAzureFile",
				},
				Disabled: []string{
					"CSIMigrationvSphere",
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
					"OpenShiftPodSecurityAdmission",
				},
				Disabled: []string{
					"CSIMigrationAzureFile",
					"CSIMigrationvSphere",
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
					"OpenShiftPodSecurityAdmission",
					"CSIMigrationAzureFile",
					"other",
				},
				Disabled: []string{
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
