package features

import (
	"slices"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes/scheme"
)

// GroupVersionResourcesForFeatureGate is intended to be used to fetch the appropriate group-version-resource
// pairings for a given feature gate name for upstream Kubernetes feature gates that have an OpenShift equivalent gate.
// This is so that it can be used to configure the kube-apiserver with the correct --runtime-config options
// to only ever enable the most up-to-date versions of an API that a gate depends on.
// This will only ever return GVRs for alpha or beta APIs. Stable APIs are served by default and therefore do not require
// an explicit --runtime-config entry.
// If the provided gate does not exist or is a stable API version, the return value will be nil.
func GroupVersionResourcesForFeatureGate(featureGate configv1.FeatureGateName) []schema.GroupVersionResource {
	gateStatuses, ok := allFeatureGates[featureGate]
	if !ok {
		return nil
	}

	// aggregate GRs across all gate statuses. Realistically, this should only ever be a single
	// gate status for upstream equivalent gates.
	groupKindResources := sets.New[groupKindResource]()
	for _, status := range gateStatuses {
		groupKindResources.Insert(status.groupKindResources.UnsortedList()...)
	}

	gvrs := []schema.GroupVersionResource{}

	for _, gkr := range groupKindResources.UnsortedList() {
		versions := scheme.Scheme.VersionsForGroupKind(schema.GroupKind{Group: gkr.Group, Kind: gkr.Kind})

		// ensure that we are always sorting in descending order of version priority
		slices.SortFunc(versions, func(a, b schema.GroupVersion) int {
			return version.CompareKubeAwareVersionStrings(b.Version, a.Version)
		})

		if len(versions) == 0 {
			continue
		}

		gvrs = append(gvrs, schema.GroupVersionResource{
			Group:    versions[0].Group,
			Version:  versions[0].Version,
			Resource: gkr.Resource,
		})
	}

	return gvrs
}
