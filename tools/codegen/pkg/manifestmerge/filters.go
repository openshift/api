package manifestmerge

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/openshift/api/tools/codegen/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "sigs.k8s.io/yaml"
)

func AllKnownFeatureSets(payloadFeatureGatePath string) (sets.Set[string], error) {
	allFeatureSets := sets.Set[string]{}
	allFeatureSets.Insert("CustomNoUpgrade") // this one won't have a rendered version since we don't know the gates

	featureSetManifestFiles, err := os.ReadDir(payloadFeatureGatePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read FeatureSetManifestDir: %w", err)
	}
	for _, currFeatureSetManifestFile := range featureSetManifestFiles {
		featureGateFilename := filepath.Join(payloadFeatureGatePath, currFeatureSetManifestFile.Name())
		featureGateBytes, err := os.ReadFile(featureGateFilename)
		if err != nil {
			return nil, fmt.Errorf("unable to read %q: %w", featureGateFilename, err)
		}

		// use unstructured to pull this information to avoid vendoring openshift/api
		featureGateMap := map[string]interface{}{}
		if err := kyaml.Unmarshal(featureGateBytes, &featureGateMap); err != nil {
			return nil, fmt.Errorf("unable to parse featuregate %q: %w", featureGateFilename, err)
		}
		uncastFeatureGate := unstructured.Unstructured{
			Object: featureGateMap,
		}

		currFeatureSet, _, _ := unstructured.NestedString(uncastFeatureGate.Object, "spec", "featureSet")
		if len(currFeatureSet) == 0 {
			currFeatureSet = "Default"
		}
		allFeatureSets.Insert(currFeatureSet)
	}

	return allFeatureSets, nil
}

// VersionedFeatureSet represents a feature set with its version information
type VersionedFeatureSet struct {
	FeatureSet      string
	ClusterProfile  string
	VersionRange    VersionRange
	FeatureGateFile string
}

// VersionRange represents a list of major versions (may not be consecutive)
type VersionRange []uint64

// String returns a string representation of the version range
func (vr VersionRange) String() string {
	if len(vr) == 0 {
		return ""
	}
	if len(vr) == 1 {
		return fmt.Sprintf("%d", vr[0])
	}

	// Check if versions are consecutive to display as range
	allConsecutive := true
	for i := 1; i < len(vr); i++ {
		if vr[i] != vr[i-1]+1 {
			allConsecutive = false
			break
		}
	}

	if allConsecutive {
		return fmt.Sprintf("%d-%d", vr[0], vr[len(vr)-1])
	}

	// Display as comma-separated list
	versionStrs := make([]string, len(vr))
	for i, v := range vr {
		versionStrs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(versionStrs, ",")
}

// Contains checks if a version is within this range
func (vr VersionRange) Contains(version uint64) bool {
	for _, v := range vr {
		if v == version {
			return true
		}
	}
	return false
}

// StartVersion returns the lowest version in the range
func (vr VersionRange) StartVersion() uint64 {
	if len(vr) == 0 {
		return 0
	}

	return vr[0]
}

// EndVersion returns the highest version in the range
func (vr VersionRange) EndVersion() uint64 {
	if len(vr) == 0 {
		return 0
	}

	return vr[len(vr)-1]
}

// parseFeatureGateFile reads and parses a featureGate file to extract version information from annotations
// Returns VersionedFeatureSet with parsed information, or error if format doesn't match
func parseFeatureGateFile(filePath string) (*VersionedFeatureSet, error) {
	filename := filepath.Base(filePath)

	// Read the file content
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse YAML to extract annotations
	featureGateMap := map[string]interface{}{}
	if err := kyaml.Unmarshal(fileBytes, &featureGateMap); err != nil {
		return nil, fmt.Errorf("unable to parse featuregate %s: %w", filename, err)
	}

	uncastFeatureGate := unstructured.Unstructured{
		Object: featureGateMap,
	}

	// Extract cluster profile from annotations
	annotations := uncastFeatureGate.GetAnnotations()
	if annotations == nil {
		return nil, fmt.Errorf("no annotations found in %s", filename)
	}

	// Extract version information from annotation
	versionStr, found := annotations["release.openshift.io/major-version"]
	if !found {
		return nil, fmt.Errorf("no version annotation found in %s", filename)
	}

	// Parse comma-separated version list
	versionRange, err := parseVersionsFromAnnotation(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version annotation in %s: %w", filename, err)
	}

	// Determine cluster profile from the annotations
	var clusterProfile string
	for annotationKey, annotationValue := range annotations {
		if annotationValue == "false-except-for-the-config-operator" {
			// Convert long cluster profile names to short names
			switch annotationKey {
			case "include.release.openshift.io/ibm-cloud-managed":
				clusterProfile = "Hypershift"
			case "include.release.openshift.io/self-managed-high-availability":
				clusterProfile = "SelfManagedHA"
			default:
				clusterProfile = annotationKey // fallback to the key itself
			}
			break
		}
	}

	if clusterProfile == "" {
		return nil, fmt.Errorf("no cluster profile annotation found in %s", filename)
	}

	// Extract feature set from spec
	featureSet, _, _ := unstructured.NestedString(uncastFeatureGate.Object, "spec", "featureSet")
	if featureSet == "" {
		featureSet = "Default"
	}

	return &VersionedFeatureSet{
		FeatureSet:      featureSet,
		ClusterProfile:  clusterProfile,
		VersionRange:    versionRange,
		FeatureGateFile: filename,
	}, nil
}

// parseVersionsFromAnnotation parses a comma-separated list of versions like "4,5,6,7,8,9,10"
func parseVersionsFromAnnotation(versionStr string) (VersionRange, error) {
	versionStr = strings.TrimSpace(versionStr)
	if versionStr == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Split by comma and parse each version
	versionParts := strings.Split(versionStr, ",")
	var versions VersionRange

	for _, part := range versionParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		version, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid version '%s': %w", part, err)
		}
		versions = append(versions, version)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no valid versions found")
	}

	// Sort versions
	slices.Sort(versions)

	return versions, nil
}

// AllVersionedFeatureSets returns all versioned feature sets found in the payload directory
func AllVersionedFeatureSets(payloadFeatureGatePath string) ([]VersionedFeatureSet, error) {
	var versionedFeatureSets []VersionedFeatureSet

	featureSetManifestFiles, err := os.ReadDir(payloadFeatureGatePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read FeatureSetManifestDir: %w", err)
	}

	for _, currFeatureSetManifestFile := range featureSetManifestFiles {
		if !strings.HasSuffix(currFeatureSetManifestFile.Name(), ".yaml") {
			continue
		}

		filePath := filepath.Join(payloadFeatureGatePath, currFeatureSetManifestFile.Name())
		versionedFeatureSet, err := parseFeatureGateFile(filePath)
		if err != nil {
			// Skip files that don't contain version annotations or can't be parsed
			// This allows for backward compatibility with other files in the directory
			continue
		}

		versionedFeatureSets = append(versionedFeatureSets, *versionedFeatureSet)
	}

	// Sort by version range start, then by cluster profile, then by feature set
	sort.Slice(versionedFeatureSets, func(i, j int) bool {
		a, b := versionedFeatureSets[i], versionedFeatureSets[j]
		if a.VersionRange.StartVersion() != b.VersionRange.StartVersion() {
			return a.VersionRange.StartVersion() < b.VersionRange.StartVersion()
		}
		if a.ClusterProfile != b.ClusterProfile {
			return a.ClusterProfile < b.ClusterProfile
		}
		return a.FeatureSet < b.FeatureSet
	})

	return versionedFeatureSets, nil
}

// GetVersionedFeatureSets returns versioned feature sets for a specific version
func GetVersionedFeatureSets(payloadFeatureGatePath string, targetVersion uint64) ([]VersionedFeatureSet, error) {
	allVersioned, err := AllVersionedFeatureSets(payloadFeatureGatePath)
	if err != nil {
		return nil, err
	}

	var applicable []VersionedFeatureSet
	for _, vfs := range allVersioned {
		if vfs.VersionRange.Contains(targetVersion) {
			applicable = append(applicable, vfs)
		}
	}

	return applicable, nil
}

// AllKnownVersions returns all versions found in versioned feature sets
func AllKnownVersions(payloadFeatureGatePath string) ([]uint64, error) {
	allVersioned, err := AllVersionedFeatureSets(payloadFeatureGatePath)
	if err != nil {
		return nil, err
	}

	versionSet := sets.New[uint64]()
	for _, vfs := range allVersioned {
		for _, version := range vfs.VersionRange {
			versionSet.Insert(version)
		}
	}

	versions := make([]uint64, 0, versionSet.Len())
	for _, v := range versionSet.UnsortedList() {
		versions = append(versions, v)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	return versions, nil
}

func FilterForFeatureSet(payloadFeatureGatePath, clusterProfile, featureSetName string) (ManifestFilter, error) {
	if featureSetName == "CustomNoUpgrade" {
		return &AndManifestFilter{
			filters: []ManifestFilter{
				&CustomNoUpgrade{},
				&ClusterProfileFilter{
					clusterProfile: sets.New[string](clusterProfile),
				},
			},
		}, nil
	}

	allKnownFeatureSets, err := AllKnownFeatureSets(payloadFeatureGatePath)
	if err != nil {
		return nil, fmt.Errorf("failed reading featuresets from %q", payloadFeatureGatePath)
	}
	if !allKnownFeatureSets.Has(featureSetName) {
		return nil, fmt.Errorf("unrecognized featureset name %q", featureSetName)
	}
	// First, try to find versioned feature gate files for this combination
	versionedFeatureSets, err := AllVersionedFeatureSets(payloadFeatureGatePath)
	if err != nil {
		return nil, fmt.Errorf("failed reading versioned featuresets from %q: %w", payloadFeatureGatePath, err)
	}

	clusterProfileShortName, err := utils.ClusterProfileToShortName(clusterProfile)
	if err != nil {
		return nil, fmt.Errorf("unrecognized clusterprofile name %q: %w", clusterProfile, err)
	}

	// Look for versioned files that match this cluster profile and feature set
	var matchingVersionedFiles []string
	for _, vfs := range versionedFeatureSets {
		if vfs.ClusterProfile == clusterProfileShortName && vfs.FeatureSet == featureSetName {
			matchingVersionedFiles = append(matchingVersionedFiles, vfs.FeatureGateFile)
		}
	}

	var featureGateFilename string
	if len(matchingVersionedFiles) > 0 {
		// Use the first matching versioned file (they should have identical content if properly consolidated)
		featureGateFilename = path.Join(payloadFeatureGatePath, matchingVersionedFiles[0])
	} else {
		// Fallback to legacy filename pattern for backward compatibility
		featureGateFilename = path.Join(payloadFeatureGatePath, fmt.Sprintf("featureGate-%s-%s.yaml", clusterProfileShortName, featureSetName))
	}

	enabledFeatureGatesSet := sets.NewString()

	featureGateBytes, err := os.ReadFile(featureGateFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read feature gate file %s: %w", featureGateFilename, err)
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
				clusterProfile: sets.New[string](clusterProfile),
			},
		},
	}, nil
}

// FilterForVersionedFeatureSet returns a filter for a specific version, cluster profile, and feature set
func FilterForVersionedFeatureSet(payloadFeatureGatePath, clusterProfile, featureSetName string, targetVersion uint64) (ManifestFilter, error) {
	if featureSetName == "CustomNoUpgrade" {
		return &AndManifestFilter{
			filters: []ManifestFilter{
				&CustomNoUpgrade{},
				&ClusterProfileFilter{
					clusterProfile: sets.New[string](clusterProfile),
				},
			},
		}, nil
	}

	// Get versioned feature sets for the target version
	versionedFeatureSets, err := GetVersionedFeatureSets(payloadFeatureGatePath, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed reading versioned featuresets for version %d from %q: %w", targetVersion, payloadFeatureGatePath, err)
	}

	clusterProfileShortName, err := utils.ClusterProfileToShortName(clusterProfile)
	if err != nil {
		return nil, fmt.Errorf("unrecognized clusterprofile name %q: %w", clusterProfile, err)
	}

	// Find the specific versioned feature set for this combination
	var matchingFeatureGateFile string
	for _, vfs := range versionedFeatureSets {
		if vfs.ClusterProfile == clusterProfileShortName && vfs.FeatureSet == featureSetName {
			matchingFeatureGateFile = vfs.FeatureGateFile
			break
		}
	}

	if matchingFeatureGateFile == "" {
		return nil, fmt.Errorf("no feature gate file found for version %d, cluster profile %s, feature set %s", targetVersion, clusterProfile, featureSetName)
	}

	featureGateFilename := path.Join(payloadFeatureGatePath, matchingFeatureGateFile)
	enabledFeatureGatesSet := sets.NewString()

	featureGateBytes, err := os.ReadFile(featureGateFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read feature gate file %s: %w", featureGateFilename, err)
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
				clusterProfile: sets.New[string](clusterProfile),
			},
			&VersionFilter{
				targetVersion: sets.New[uint64](targetVersion),
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

	return f.allowedFeatureGates.HasAll(manifestFeatureGates.UnsortedList()...), nil
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
	clusterProfile sets.Set[string]
}

func (f *ClusterProfileFilter) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}
	// if there's no preferenceinclude everywhere
	if !utils.HasClusterProfilePreference(partialObject.GetAnnotations()) {
		return true, nil
	}

	for _, currClusterProfile := range f.clusterProfile.UnsortedList() {
		if partialObject.GetAnnotations()[currClusterProfile] == "true" {
			return true, nil
		}
	}
	return false, nil
}

func (f *ClusterProfileFilter) UseCRD(metadata crdForFeatureSet) bool {
	return metadata.clusterProfile.Equal(f.clusterProfile)
}

func (f *ClusterProfileFilter) String() string {
	return fmt.Sprintf("clusterProfile=%v", f.clusterProfile)
}

type VersionFilter struct {
	targetVersion sets.Set[uint64]
}

func (f *VersionFilter) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}

	versionStr, found := partialObject.GetAnnotations()["release.openshift.io/major-version"]
	if !found {
		// If the manifest doesn't restrict itself to a specific version, we must include it.
		return true, nil
	}

	versionRange, err := parseVersionsFromAnnotation(versionStr)
	if err != nil {
		return false, fmt.Errorf("invalid version annotation in %s: %w", partialObject.GetName(), err)
	}
	return f.targetVersion.Equal(sets.New[uint64](versionRange...)), nil
}

func (f *VersionFilter) UseCRD(metadata crdForFeatureSet) bool {
	return f.targetVersion.Equal(metadata.version)
}

func (f *VersionFilter) String() string {
	return fmt.Sprintf("version=%v", f.targetVersion.UnsortedList())
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
	featureSetName sets.Set[string]
}

func (f *FeatureSetFilter) UseManifest(data []byte) (bool, error) {
	partialObject := &metav1.PartialObjectMetadata{}
	if err := kyaml.Unmarshal(data, partialObject); err != nil {
		return false, err
	}

	manifestFeatureSet := partialObject.GetAnnotations()["release.openshift.io/feature-set"]
	return f.featureSetName.Has(manifestFeatureSet), nil
}

func (f *FeatureSetFilter) UseCRD(metadata crdForFeatureSet) bool {
	return metadata.featureSet.Equal(f.featureSetName)
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

type EqualData struct {
	data *unstructured.Unstructured
}

func (f *EqualData) UseCRD(metadata crdForFeatureSet) bool {
	return reflect.DeepEqual(metadata.crd, f.data)
}

func (f *EqualData) String() string {
	return fmt.Sprintf("EqualData=%v", f.data)
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
