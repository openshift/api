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
	for _, capability := range KnownClusterVersionCapabilities {
		known[capability] = exists
	}

	for set, caps := range ClusterVersionCapabilitySets {
		for _, capability := range caps {
			if _, ok := known[capability]; !ok {
				t.Errorf("Capability set %s contains %s, which needs to be added to KnownClusterVersionCapabilities", set, capability)
			}
		}
	}
}

// TestSubsetsVersionCapabilities verifies that all preceding sets
// (order given by capabilitySetOrder) are subsets of following sets
func TestSubsetsVersionCapabilities(t *testing.T) {
	// defines the order by which subsets are tested
	var capabilitySetOrder = []ClusterVersionCapabilitySet{
		ClusterVersionCapabilitySetNone,
		ClusterVersionCapabilitySet4_11,
		ClusterVersionCapabilitySet4_12,
		ClusterVersionCapabilitySetCurrent,
	}

	for i := 0; i < len(capabilitySetOrder)-1; i++ {
		var setName = capabilitySetOrder[i]
		var setName1 = capabilitySetOrder[i+1]

		for _, capability := range ClusterVersionCapabilitySets[setName] {
			found := false
			for _, capability1 := range ClusterVersionCapabilitySets[setName1] {
				if capability == capability1 {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Capability set %s is not a subset of %s (contains %s)", setName, setName1, capability)
			}
		}
	}
}
