package manifestmerge

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	kyaml "sigs.k8s.io/yaml"
	"strings"
)

func FilterForFeatureSet(featureSetName string) (ManifestFilter, error) {
	if featureSetName == "CustomNoUpgrade" {
		return &CustomNoUpgrade{}, nil
	}

	// TODO find pointers to serialized featuregates.  I'm sure Joel has some sneaky thing that prevents me from doing the "simple" relative path.
	featureGateFilename := ""
	switch {
	case featureSetName == "TechPreviewNoUpgrade":
		featureGateFilename = "payload-manifests/featuregates/featureGate-TechPreviewNoUpgrade.yaml"
	case featureSetName == "Default":
		featureGateFilename = "payload-manifests/featuregates/featureGate-Default.yaml"
	default:
		return nil, fmt.Errorf("unrecognized featureset name %q", featureSetName)
	}

	enabledFeatureGatesSet := sets.NewString()

	featureGateBytes, err := os.ReadFile(featureGateFilename)
	if err != nil {
		return nil, err
	}

	// use unstructured to pull this information to avoid vendoring openshift/api
	uncastFeatureGate := map[string]interface{}{}
	if err := kyaml.Unmarshal(featureGateBytes, &uncastFeatureGate); err != nil {
		return nil, fmt.Errorf("unable to parse featuregate %q: %w", featureGateFilename, err)
	}

	uncastFeatureGateSlice, _, err := unstructured.NestedSlice(uncastFeatureGate, "status", "featureGates")
	if err != nil {
		return nil, fmt.Errorf("no slice found %w", err)
	}
	enabledFeatureGates, _, err := unstructured.NestedSlice(uncastFeatureGateSlice[0].(map[string]interface{}), "enabled")
	if err != nil {
		return nil, fmt.Errorf("no enabled found %w", err)
	}
	for _, currGate := range enabledFeatureGates {
		featureGateName, _, err := unstructured.NestedString(currGate.(map[string]interface{}), "name")
		if err != nil {
			return nil, fmt.Errorf("no gate name found %w", err)
		}
		enabledFeatureGatesSet.Insert(featureGateName)
	}

	return &ForFeatureGates{
		allowedFeatureGates: enabledFeatureGatesSet,
	}, nil
}

type ManifestFilter interface {
	UseManifest([]byte) (bool, error)
}

type AllFeatureGates struct{}

func (*AllFeatureGates) UseManifest([]byte) (bool, error) {
	return true, nil
}

type CustomNoUpgrade struct{}

func (*CustomNoUpgrade) UseManifest([]byte) (bool, error) {
	return true, nil
}

type ForFeatureGates struct {
	allowedFeatureGates sets.String
}

func (f *ForFeatureGates) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}

	manifestFeatureGates := featureGatesFromManifest(partialObject)
	if len(manifestFeatureGates) == 0 || manifestFeatureGates.Has("") {
		// always include ungated manifests
		return true, nil
	}

	return manifestFeatureGates.HasAny(f.allowedFeatureGates.UnsortedList()...), nil
}

func featureGatesFromManifest(manifest metav1.Object) sets.String {
	ret := sets.String{}
	for existingAnnotation := range manifest.GetAnnotations() {
		if strings.HasPrefix(existingAnnotation, "feature-gate.release.openshift.io/") {
			featureGateName := strings.TrimPrefix(existingAnnotation, "feature-gate.release.openshift.io/")
			ret.Insert(featureGateName)
		}
	}
	return ret
}
