package manifestmerge

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	"path"
	kyaml "sigs.k8s.io/yaml"
	"strings"
)

var clusterProfileToShortName = map[string]string{
	"include.release.openshift.io/ibm-cloud-managed":              "Hypershift",
	"include.release.openshift.io/self-managed-high-availability": "SelfManagedHA",
	"include.release.openshift.io/single-node-developer":          "SingleNode",
}

func FilterForFeatureSet(payloadFeatureGatePath, clusterProfile, featureSetName string) (ManifestFilter, error) {
	if featureSetName == "CustomNoUpgrade" {
		return &AndManifestFilter{
			filters: []ManifestFilter{
				&CustomNoUpgrade{},
				&ClusterProfileFilter{
					clusterProfile: clusterProfile,
				},
			},
		}, nil
	}

	switch {
	case featureSetName == "TechPreviewNoUpgrade":
	case featureSetName == "Default":
	default:
		return nil, fmt.Errorf("unrecognized featureset name %q", featureSetName)
	}
	featureGateFilename := path.Join(payloadFeatureGatePath, fmt.Sprintf("featureGate-%s-%s.yaml", clusterProfileToShortName[clusterProfile], featureSetName))

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

	return &AndManifestFilter{
		filters: []ManifestFilter{
			&ForFeatureGates{
				allowedFeatureGates: enabledFeatureGatesSet,
			},
			&ClusterProfileFilter{
				clusterProfile: clusterProfile,
			},
		},
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

func (f *CustomNoUpgrade) String() string {
	return fmt.Sprintf("CustomNoUpgrade")
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

func (f *ForFeatureGates) String() string {
	return fmt.Sprintf("featureGates/%d", len(f.allowedFeatureGates))
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

type ClusterProfileFilter struct {
	clusterProfile string
}

func (f *ClusterProfileFilter) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}
	// if there's no preferenceinclude everywhere
	if !hasClusterProfilePreference(partialObject.GetAnnotations()) {
		return true, nil
	}

	forThisProfile := partialObject.GetAnnotations()[f.clusterProfile] == "true"
	return forThisProfile, nil
}

func (f *ClusterProfileFilter) UseCRD(metadata crdForFeatureSet) bool {
	return metadata.clusterProfile == f.clusterProfile
}

func (f *ClusterProfileFilter) String() string {
	return fmt.Sprintf("clusterProfile=%v", f.clusterProfile)
}

type AndManifestFilter struct {
	filters []ManifestFilter
}

func (f *AndManifestFilter) UseManifest(data []byte) (bool, error) {
	for _, curr := range f.filters {
		ret, err := curr.UseManifest(data)
		if err != nil {
			return false, err
		}
		if !ret {
			return false, nil
		}
	}

	return true, nil
}

func (f *AndManifestFilter) String() string {
	str := []string{}
	for _, curr := range f.filters {
		str = append(str, fmt.Sprintf("%v", curr))
	}
	return strings.Join(str, " AND ")
}

type CRDFilter interface {
	UseCRD(metadata crdForFeatureSet) bool
}

type AndCRDFilter struct {
	filters []CRDFilter
}

func (f *AndCRDFilter) UseCRD(metadata crdForFeatureSet) bool {
	for _, curr := range f.filters {
		ret := curr.UseCRD(metadata)
		if !ret {
			return false
		}
	}

	return true
}

func (f *AndCRDFilter) String() string {
	str := []string{}
	for _, curr := range f.filters {
		str = append(str, fmt.Sprintf("%v", curr))
	}
	return strings.Join(str, " AND ")
}

type FeatureSetFilter struct {
	featureSetName string
}

func (f *FeatureSetFilter) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}

	forThisFeatureSet := partialObject.GetAnnotations()["release.openshift.io/feature-set"] == f.featureSetName
	return forThisFeatureSet, nil
}

func (f *FeatureSetFilter) UseCRD(metadata crdForFeatureSet) bool {
	return metadata.featureSet == f.featureSetName
}

func (f *FeatureSetFilter) String() string {
	return fmt.Sprintf("featureSetName=%v", f.featureSetName)
}

type HasData struct {
}

func (f *HasData) UseCRD(metadata crdForFeatureSet) bool {
	return metadata.noData == false
}

func (f *HasData) String() string {
	return "HasData"
}

type Everything struct {
}

func (f *Everything) UseManifest(data []byte) (bool, error) {
	return true, nil
}

func (f *Everything) UseCRD(metadata crdForFeatureSet) bool {
	return true
}

func (f *Everything) String() string {
	return "Everything"
}
