package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// InsightsDataGather provides data gather configuration options for the the Insights Operator.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=insightsdatagathers,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/1245
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=InsightsConfig
// +openshift:compatibility-gen:level=4
type InsightsDataGather struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec InsightsDataGatherSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status InsightsDataGatherStatus `json:"status"`
}

type InsightsDataGatherSpec struct {
	// gatherConfig spec attribute includes all the configuration options related to
	// gathering of the Insights data and its uploading to the ingress.
	// +optional
	GatherConfig GatherConfig `json:"gatherConfig,omitempty"`
}

type InsightsDataGatherStatus struct{}

// gatherConfig provides data gathering configuration options.
type GatherConfig struct {
	// dataPolicy allows user to enable additional global obfuscation of the IP addresses and base domain
	// in the Insights archive data. Valid values are "None" and "ObfuscateNetworking".
	// When set to None the data is not obfuscated.
	// When set to ObfuscateNetworking the IP addresses and the cluster domain name are obfuscated.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	DataPolicy DataPolicy `json:"dataPolicy,omitempty"`
	// disabledGatherers is a list of gatherers to be excluded from the gathering. All the gatherers can be disabled by providing "all" value.
	// If all the gatherers are disabled, the Insights operator does not gather any data.
	// The format for the disabledGatherer should be: {gatherer}/{function} where the function is optional.
	// Gatherer consists of a lowercase string that may include underscores (_).
	// Function consists of a lowercase string that may include underscores (_) and is separated from the gatherer by a forward slash (/).
	// The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// Run the following command to get the names of last active gatherers:
	// "oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'"
	// An example of disabling gatherers looks like this: `disabledGatherers: ["clusterconfig/machine_configs", "workloads/workload_info"]`
	// +kubebuilder:validation:MaxItems=100
	// +optional
	DisabledGatherers []DisabledGatherer `json:"disabledGatherers"`
	// storageSpec is an optional field that allows user to define persistent storage for on-demand gathering
	// jobs to store the Insights data archive.
	// If omitted, the gathering job will use ephemeral storage.
	// +optional
	StorageSpec *StorageSpec `json:"storageSpec,omitempty"`
}

// disabledGatherer is a string that represents a gatherer that should be disabled
// +kubebuilder:validation:MaxLength=256
// +kubebuilder:validation:XValidation:rule=`self.matches("^[a-z]+[_a-z]*[a-z]([/a-z][_a-z]*)?[a-z]$")`,message=`disabledGatherer must be in the format of {gatherer}/{function} where the gatherer and function are lowercase strings that may include underscores (_) and are separated by a forward slash (/) if the function is provided`
type DisabledGatherer string

// storageSpec provides persistent storage configuration options for on-demand gathering jobs.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'PersistentVolumeClaim' ?  has(self.persistentVolume) : !has(self.persistentVolume)",message="persistentVolume is required when type is PersistentVolumeClaim, and forbidden otherwise"
type StorageSpec struct {
	// type is a required field that specifies the type of storage that will be used to store the Insights data archive.
	// Valid values are "PersistentVolumeClaim" and "Ephemeral".
	// +required
	Type StorageType `json:"type"`
	// persistentVolume is an optional field that specifies the PersistentVolume that will be used to store the Insights data archive.
	// The PersistentVolume must be created in the openshift-insights namespace.
	// +optional
	PersistentVolume *PersistentVolumeConfig `json:"persistentVolume,omitempty"`
}

// storageType declares valid storage types
// +kubebuilder:validation:Enum=PersistentVolumeClaim;Ephemeral
type StorageType string

const (
	// StorageTypePersistentVolumeClaim storage type
	StorageTypePersistentVolumeClaim StorageType = "PersistentVolumeClaim"
	// StorageTypeEphemeral storage type
	StorageTypeEphemeral StorageType = "Ephemeral"
)

// persistentVolumeConfig provides configuration options for PersistentVolume storage.
type PersistentVolumeConfig struct {
	// persistentVolumeClaim is a required field that specifies the configuration of the PersistentVolumeClaim that will
	// be used to store the Insights data archive. The PersistentVolumeClaim must be created in the openshift-insights namespace.
	// +required
	PersistentVolumeClaim PersistentVolumeClaimReference `json:"persistentVolumeClaim"`
	// mountPath is an optional field specifying the directory where the PVC will be mounted inside the
	// Insights data gathering Pod. If omitted, the path that is used to store the Insights data archive by Insights
	// operator will be used instead. The path cannot exceed 1024 characters and must not contain a colon.
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="!self.contains(':')",message="mountPath must not contain a colon"
	// +optional
	MountPath string `json:"mountPath,omitempty"`
}

// persistentVolumeClaimReference is a reference to a PersistentVolumeClaim.
type PersistentVolumeClaimReference struct {
	// name is a string that follows the DNS1123 subdomain format.
	// It must be at most 253 characters in length, and must consist only of lower case alphanumeric characters,
	//  '-' and '.', and must start and end with an alphanumeric character.
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	// +required
	Name string `json:"name"`
}

const (
	// No data obfuscation
	NoPolicy DataPolicy = "None"
	// IP addresses and cluster domain name are obfuscated
	ObfuscateNetworking DataPolicy = "ObfuscateNetworking"
)

// dataPolicy declares valid data policy types
// +kubebuilder:validation:Enum="";None;ObfuscateNetworking
type DataPolicy string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InsightsDataGatherList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type InsightsDataGatherList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`
	Items           []InsightsDataGather `json:"items"`
}
