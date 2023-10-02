package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:resource:scope=Cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// DataGather provides data gather configuration options and status for the particular Insights data gathering.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type DataGather struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	Spec DataGatherSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status DataGatherStatus `json:"status"`
}

// DataGatherSpec contains the configuration for the DataGather.
type DataGatherSpec struct {
	// dataPolicy allows user to enable additional global obfuscation of the IP addresses and base domain
	// in the Insights archive data. Valid values are "ClearText" and "ObfuscateNetworking".
	// When set to ClearText the data is not obfuscated.
	// When set to ObfuscateNetworking the IP addresses and the cluster domain name are obfuscated.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The current default is ClearText.
	// +optional
	DataPolicy DataPolicy `json:"dataPolicy"`
	// gatherers is a list of gatherers configurations.
	// The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// Run the following command to get the names of last active gatherers:
	// "oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'"
	// +optional
	Gatherers []GathererConfig `json:"gatherers"`
}

const (
	// No data obfuscation
	NoPolicy DataPolicy = "ClearText"
	// IP addresses and cluster domain name are obfuscated
	ObfuscateNetworking DataPolicy = "ObfuscateNetworking"
	// Data gathering is running
	Running DataGatherState = "Running"
	// Data gathering is completed
	Completed DataGatherState = "Completed"
	// Data gathering failed
	Failed DataGatherState = "Failed"
	// Data gathering is pending
	Pending DataGatherState = "Pending"
	// Gatherer state marked as disabled, which means that the gatherer will not run.
	Disabled GathererState = "Disabled"
	// Gatherer state marked as enabled, which means that the gatherer will run.
	Enabled GathererState = "Enabled"
)

// dataPolicy declares valid data policy types
// +kubebuilder:validation:Enum="";ClearText;ObfuscateNetworking
type DataPolicy string

// state declares valid gatherer state types.
// +kubebuilder:validation:Enum="";Enabled;Disabled
type GathererState string

// gathererConfig allows to configure specific gatherers
type GathererConfig struct {
	// name is the name of specific gatherer
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// state allows you to configure specific gatherer. Valid values are "Enabled", "Disabled" and omitted.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default.
	// The current default is Enabled.
	// +optional
	State GathererState `json:"state"`
}

// dataGatherState declares valid gathering state types
// +kubebuilder:validation:Optional
// +kubebuilder:validation:Enum=Running;Completed;Failed;Pending
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Running' && self == 'Pending')", message="dataGatherState cannot transition from Running to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Completed' && self == 'Pending')", message="dataGatherState cannot transition from Completed to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Failed' && self == 'Pending')", message="dataGatherState cannot transition from Failed to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Completed' && self == 'Running')", message="dataGatherState cannot transition from Completed to Running"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Failed' && self == 'Running')", message="dataGatherState cannot transition from Failed to Running"
type DataGatherState string

// DataGatherStatus contains information relating to the DataGather state.
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.insightsRequestID) || has(self.insightsRequestID))",message="cannot remove insightsRequestID attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.startTime) || has(self.startTime))",message="cannot remove startTime attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.finishTime) || has(self.finishTime))",message="cannot remove finishTime attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.dataGatherState) || has(self.dataGatherState))",message="cannot remove dataGatherState attribute from status"
// +kubebuilder:validation:Optional
type DataGatherStatus struct {
	// conditions provide details on the status of the gatherer job.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions" patchStrategy:"merge" patchMergeKey:"type"`
	// dataGatherState reflects the current state of the data gathering process.
	// +optional
	State DataGatherState `json:"dataGatherState,omitempty"`
	// gatherers is a list of active gatherers (and their statuses) in the last gathering.
	// +listType=map
	// +listMapKey=name
	// +optional
	Gatherers []GathererStatus `json:"gatherers,omitempty"`
	// startTime is the time when Insights data gathering started.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="startTime is immutable once set"
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// finishTime is the time when Insights data gathering finished.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="finishTime is immutable once set"
	// +optional
	FinishTime metav1.Time `json:"finishTime,omitempty"`
	// relatedObjects is a list of resources which are useful when debugging or inspecting the data
	// gathering Pod
	// +optional
	RelatedObjects []ObjectReference `json:"relatedObjects,omitempty"`
	// insightsRequestID is an Insights request ID to track the status of the
	// Insights analysis (in console.redhat.com processing pipeline) for the corresponding Insights data archive.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="insightsRequestID is immutable once set"
	// +kubebuilder:validation:Optional
	// +optional
	InsightsRequestID string `json:"insightsRequestID,omitempty"`
	// insightsReport provides general Insights analysis results.
	// When omitted, this means no data gathering has taken place yet or the
	// corresponding Insights analysis (identified by "insightsRequestID") is not available.
	// +optional
	InsightsReport InsightsReport `json:"insightsReport,omitempty"`
}

// gathererStatus represents information about a particular
// data gatherer.
type GathererStatus struct {
	// conditions provide details on the status of each gatherer.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Conditions []metav1.Condition `json:"conditions" patchStrategy:"merge" patchMergeKey:"type"`
	// name is the name of the gatherer.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:MinLength=5
	Name string `json:"name"`
	// lastGatherDuration represents the time spent gathering.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([1-9][0-9]*(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	LastGatherDuration metav1.Duration `json:"lastGatherDuration"`
}

// insightsReport provides Insights health check report based on the most
// recently sent Insights data.
type InsightsReport struct {
	// downloadedAt is the time when the last Insights report was downloaded.
	// An empty value means that there has not been any Insights report downloaded yet and
	// it usually appears in disconnected clusters (or clusters when the Insights data gathering is disabled).
	// +optional
	DownloadedAt metav1.Time `json:"downloadedAt,omitempty"`
	// healthChecks provides basic information about active Insights health checks
	// in a cluster.
	// +listType=atomic
	// +optional
	HealthChecks []HealthCheck `json:"healthChecks,omitempty"`
	// uri provides the URL link from which the report was downloaded.
	// +kubebuilder:validation:Pattern=`^https:\/\/\S+`
	// +optional
	URI string `json:"uri,omitempty"`
}

// healthCheck represents an Insights health check attributes.
type HealthCheck struct {
	// description provides basic description of the healtcheck.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:MinLength=10
	Description string `json:"description"`
	// totalRisk of the healthcheck. Indicator of the total risk posed
	// by the detected issue; combination of impact and likelihood. The values can be from 1 to 4,
	// and the higher the number, the more important the issue.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=4
	TotalRisk int32 `json:"totalRisk"`
	// advisorURI provides the URL link to the Insights Advisor.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https:\/\/\S+`
	AdvisorURI string `json:"advisorURI"`
	// state determines what the current state of the health check is.
	// Health check is enabled by default and can be disabled
	// by the user in the Insights advisor user interface.
	// +kubebuilder:validation:Required
	State HealthCheckState `json:"state"`
}

// healthCheckState provides information about the status of the
// health check (for example, the health check may be marked as disabled by the user).
// +kubebuilder:validation:Enum:=Enabled;Disabled
type HealthCheckState string

const (
	// enabled marks the health check as enabled
	HealthCheckEnabled HealthCheckState = "Enabled"
	// disabled marks the health check as disabled
	HealthCheckDisabled HealthCheckState = "Disabled"
)

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// group is the API Group of the Resource.
	// Enter empty string for the core group.
	// This value should consist of only lowercase alphanumeric characters, hyphens and periods.
	// Example: "", "apps", "build.openshift.io", etc.
	// +kubebuilder:validation:Pattern:="^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
	// +kubebuilder:validation:Required
	Group string `json:"group"`
	// resource is the type that is being referenced.
	// It is normally the plural form of the resource kind in lowercase.
	// This value should consist of only lowercase alphanumeric characters and hyphens.
	// Example: "deployments", "deploymentconfigs", "pods", etc.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:="^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	Resource string `json:"resource"`
	// name of the referent.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// namespace of the referent.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataGatherList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type DataGatherList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	// items contains a list of DataGather resources.
	Items []DataGather `json:"items"`
}
