package v1

import (
	"fmt"
	"testing"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func validateFeatureGateDescription(in FeatureGateDescription) error {
	errs := []error{}

	if len(in.FeatureGateAttributes.Name) == 0 {
		errs = append(errs, fmt.Errorf("must have name"))
	}
	if len(in.OwningJiraComponent) == 0 {
		errs = append(errs, fmt.Errorf("featureGate/%v must have owningJiraComponent", in.FeatureGateAttributes.Name))
	}
	if len(in.ResponsiblePerson) == 0 {
		errs = append(errs, fmt.Errorf("featureGate/%v must have responsiblePerson", in.FeatureGateAttributes.Name))
	}
	if len(in.OwningProduct) == 0 {
		errs = append(errs, fmt.Errorf("featureGate/%v must have owningProduct", in.FeatureGateAttributes.Name))
	}
	if in.OwningProduct != kubernetes && in.OwningProduct != ocpSpecific {
		errs = append(errs, fmt.Errorf("featureGate/%v owningProduct must be either %q or %q", in.FeatureGateAttributes.Name, kubernetes, ocpSpecific))
	}

	return utilerrors.NewAggregate(errs)
}

func TestAllFeatureGates(t *testing.T) {
	for featureSet, currFeatures := range FeatureSets {
		for _, enabled := range currFeatures.Enabled {
			for _, disabled := range currFeatures.Disabled {
				if enabled.FeatureGateAttributes.Name == disabled.FeatureGateAttributes.Name {
					t.Errorf("featureSet/%v has featureGate/%v both enabled and disabled", featureSet, enabled.FeatureGateAttributes.Name)
				}
			}

			if err := validateFeatureGateDescription(enabled); err != nil {
				t.Error(err)
			}
		}

		for _, disabled := range currFeatures.Disabled {
			if err := validateFeatureGateDescription(disabled); err != nil {
				t.Error(err)
			}
		}
	}
}
