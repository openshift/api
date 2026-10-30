package v1

import (
	"testing"
)

// TestKnownClusterVersionCapabilities verifies that all capabilities
// referenced by capability sets are contained in
// KnownClusterVersionCapabilities.
func TestKnownClusterVersionCapabilities(t *testing.T) {
	exists := struct{}{}
	known := make(map[ClusterVersionCapability]struct{}, len(KnownClusterVersionCapabilities))
	for _, cap := range KnownClusterVersionCapabilities {
		known[cap] = exists
	}

	for set, caps := range ClusterVersionCapabilitySets {
		for _, cap := range caps {
			if _, ok := known[cap]; !ok {
				t.Errorf("Capability set %s contains %s, which needs to be added to KnownClusterVersionCapabilities", set, cap)
			}
		}
	}
}

// TestMachineAPIDependencyBaselineAssumptions checks the baseline membership
// encoded in the CEL rule on ClusterVersionCapabilitiesSpec. Update that rule
// and its admission fixtures if these assumptions change.
func TestMachineAPIDependencyBaselineAssumptions(t *testing.T) {
	for baseline, capabilities := range ClusterVersionCapabilitySets {
		t.Run(string(baseline), func(t *testing.T) {
			var machineAPI, clusterAPI bool
			for _, capability := range capabilities {
				machineAPI = machineAPI || capability == ClusterVersionCapabilityMachineAPI
				clusterAPI = clusterAPI || capability == ClusterVersionCapabilityClusterAPI
			}
			if expected := baseline != ClusterVersionCapabilitySetNone; machineAPI != expected {
				t.Errorf("baseline %s: MachineAPI membership is %t, CEL assumes %t", baseline, machineAPI, expected)
			}
			if expected := baseline == ClusterVersionCapabilitySetCurrent; clusterAPI != expected {
				t.Errorf("baseline %s: ClusterAPI membership is %t, CEL assumes %t", baseline, clusterAPI, expected)
			}
		})
	}
}
