package registry

import (
	"fmt"

	"k8s.io/component-base/featuregate"

	configv1 "github.com/openshift/api/config/v1"
)

type OpenShiftGateFields struct {
	EnhancementPullRequest string
	ContactPerson          string
	JiraComponent          string
	EnabledFeatureSets     []configv1.FeatureSet
}

type UpstreamGateFieldOverridesRegistry map[featuregate.Feature]OpenShiftGateFields

func (u UpstreamGateFieldOverridesRegistry) Register(gate featuregate.Feature, overrides OpenShiftGateFields) {
	if _, ok := u[gate]; ok {
		panic(fmt.Sprintf("gate %q has already been registered with an override", gate))
	}

	u[gate] = overrides
}

const (
	defaultEnhancement   = "https://github.com/openshift/enhancements/pull/2084"
	defaultContactPerson = "bpalmer"
	defaultJiraComponent = "kube-apiserver"
)

func (u UpstreamGateFieldOverridesRegistry) FieldsForGate(gate featuregate.Feature) OpenShiftGateFields {
	if _, ok := upstreamGateOverrideRegistry[gate]; !ok {
		return OpenShiftGateFields{
			EnhancementPullRequest: defaultEnhancement,
			ContactPerson:          defaultContactPerson,
			JiraComponent:          defaultJiraComponent,
		}
	}

	return upstreamGateOverrideRegistry[gate]
}

var upstreamGateOverrideRegistry = UpstreamGateFieldOverridesRegistry{}

func RegisterGateOverride(gate featuregate.Feature, overrides OpenShiftGateFields) {
	upstreamGateOverrideRegistry.Register(gate, overrides)
}

func FieldsForGate(gate featuregate.Feature) OpenShiftGateFields {
	return upstreamGateOverrideRegistry.FieldsForGate(gate)
}
