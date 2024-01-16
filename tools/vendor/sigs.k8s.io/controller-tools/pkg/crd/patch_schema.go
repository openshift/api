package crd

import (
	"fmt"

	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// mayHandleField returns true if the field should be considered by this invocation of the generator.
// Right now, the only skip is based on the featureset marker.
func mayHandleField(field markers.FieldInfo) bool {
	uncastFeatureSet := field.Markers.Get(crdmarkers.OpenShiftFeatureSetMarkerName)
	if uncastFeatureSet != nil {
		featureSetsForField, ok := uncastFeatureSet.([]string)
		if !ok {
				panic(fmt.Sprintf("actually got %t", uncastFeatureSet))
		}
		//  if any of the field's declared featureSets match any of the manifest's declared featuresets, include the field.
		for _, currFeatureSetForField := range featureSetsForField {
			if crdmarkers.RequiredFeatureSets.Has(currFeatureSetForField) {
				return true
			}
		}
		return false
	}

	uncastFeatureSet = field.Markers.Get(crdmarkers.OpenShiftClusterProfileAwareFeatureSets)
	if uncastFeatureSet == nil{
		return true
	}
	clusterProfileAndFeatureSet, ok := uncastFeatureSet.(crdmarkers.FeatureSetClusterProfileTuple)
	if !ok {
		panic(fmt.Sprintf("actually got %t", uncastFeatureSet))
	}
	foundFeatureSet := false
	//  if any of the field's declared featureSets match any of the manifest's declared featuresets, include the field.
	for _, currFeatureSetForField := range clusterProfileAndFeatureSet.FeatureSetNames {
		if crdmarkers.RequiredFeatureSets.Has(currFeatureSetForField) {
			foundFeatureSet = true
			break
		}
	}
	if !foundFeatureSet{
		return false
	}
	if len(clusterProfileAndFeatureSet.ClusterProfiles) == 0{
		return true
	}

	for _, currClusterProfileForFeature := range clusterProfileAndFeatureSet.ClusterProfiles{
		if crdmarkers.RequiredClusterProfiles.Has(currClusterProfileForFeature) {
			return true
		}
	}

	return false
}
