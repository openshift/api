package schemapatcher

import (
	"strings"

	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "sigs.k8s.io/yaml"
)

// mayHandleFile returns true if this manifest should progress past the file collection stage.
// Currently, the only check is the feature-set annotation.
func mayHandleFile(filename string, rawContent []byte) bool {
	manifest := &unstructured.Unstructured{}
	if err := kyaml.Unmarshal(rawContent, &manifest); err != nil {
		return true
	}

	manifestFeatureSets := sets.String{}
	if manifestFeatureSetString := manifest.GetAnnotations()["release.openshift.io/feature-set"]; len(manifestFeatureSetString) > 0 {
		for _, curr := range strings.Split(manifestFeatureSetString, ",") {
			manifestFeatureSets.Insert(curr)
		}
	}
	return manifestFeatureSets.Equal(crdmarkers.RequiredFeatureSets)
}
