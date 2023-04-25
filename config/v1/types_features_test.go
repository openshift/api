package v1

import (
	"reflect"
	"testing"
)

var (
	featureGateFoo = FeatureGateName("Foo")
	foo            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: featureGateFoo,
		},
		OwningJiraComponent: "auth",
		ResponsiblePerson:   "stlaz",
		OwningProduct:       ocpSpecific,
	}

	featureGateBar = FeatureGateName("Bar")
	bar            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: featureGateBar,
		},
		OwningJiraComponent: "storage",
		ResponsiblePerson:   "RomanBednar",
		OwningProduct:       kubernetes,
	}

	featureGateBaz = FeatureGateName("Baz")
	baz            = FeatureGateDescription{
		FeatureGateAttributes: FeatureGateAttributes{
			Name: featureGateBaz,
		},
		OwningJiraComponent: "cloud-provider",
		ResponsiblePerson:   "jspeed",
		OwningProduct:       ocpSpecific,
	}
)

var testFeatures = &FeatureGateEnabledDisabled{
	Enabled: []FeatureGateDescription{
		foo,
	},
	Disabled: []FeatureGateDescription{
		bar,
	},
}

func TestFeatureBuilder(t *testing.T) {
	tests := []struct {
		name     string
		actual   *FeatureGateEnabledDisabled
		expected *FeatureGateEnabledDisabled
	}{
		{
			name:     "nothing",
			actual:   newDefaultFeatures().toFeatures(testFeatures),
			expected: testFeatures,
		},
		{
			name:   "disable-existing",
			actual: newDefaultFeatures().without(foo).toFeatures(testFeatures),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []FeatureGateDescription{},
				Disabled: []FeatureGateDescription{
					bar,
					foo,
				},
			},
		},
		{
			name:   "enable-existing",
			actual: newDefaultFeatures().with(foo).toFeatures(testFeatures),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []FeatureGateDescription{
					foo,
				},
				Disabled: []FeatureGateDescription{
					bar,
				},
			},
		},
		{
			name:   "disable-more",
			actual: newDefaultFeatures().without(baz).toFeatures(testFeatures),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []FeatureGateDescription{
					foo,
				},
				Disabled: []FeatureGateDescription{
					bar,
					baz,
				},
			},
		},
		{
			name:   "enable-more",
			actual: newDefaultFeatures().with(baz).toFeatures(testFeatures),
			expected: &FeatureGateEnabledDisabled{
				Enabled: []FeatureGateDescription{
					foo,
					baz,
				},
				Disabled: []FeatureGateDescription{
					bar,
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
