package v1alpha2

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
	// +optional
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
	// in the Insights archive data. Valid values are "ClearText" and "ObfuscateNetworking".
	// When set to ClearText the data is not obfuscated.
	// When set to ObfuscateNetworking the IP addresses and the cluster domain name are obfuscated.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The current default is ClearText.
	// +optional
	DataPolicy DataPolicy `json:"dataPolicy,omitempty"`
	// gatherers is a list of gatherers configurations.
	// The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// Run the following command to get the names of last active gatherers:
	// "oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'"
	// +kubebuilder:validation:MaxItems=100
	// +listType=atomic
	// +optional
	Gatherers []GathererConfig `json:"gatherers,omitempty"`
}

// state declares valid gatherer state types.
// +kubebuilder:validation:Enum="";Enabled;Disabled
type GathererState string

// gathererConfig allows to configure specific gatherers
type GathererConfig struct {
	// name is the name of specific gatherer
	// +kubebuilder:validation:MaxLength=256
	// +required
	Name string `json:"name"`
	// state allows you to configure specific gatherer. Valid values are "Enabled", "Disabled" and omitted.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default.
	// The current default is Enabled.
	// +optional
	State GathererState `json:"state"`
}

const (
	// No data obfuscation
	ClearText DataPolicyOption = "ClearText"
	// IP addresses and cluster domain name are obfuscated
	ObfuscateNetworking DataPolicyOption = "ObfuscateNetworking"
	// Gatherer state marked as disabled, which means that the gatherer will not run.
	Disabled GathererState = "Disabled"
	// Gatherer state marked as enabled, which means that the gatherer will run.
	Enabled GathererState = "Enabled"
)

// dataPolicyOption declares valid data policy options
// +kubebuilder:validation:Enum="";ClearText;ObfuscateNetworking
type DataPolicyOption string

// DataPolicy is a list of data policy options
// +listType=set
// +kubebuilder:validation:MaxItems=5
type DataPolicy []DataPolicyOption

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InsightsDataGatherList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type InsightsDataGatherList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ListMeta `json:"metadata"`
	// items is the list of InsightsDataGather objects
	// +kubebuilder:validation:MaxItems=100
	// +required
	Items []InsightsDataGather `json:"items"`
}
