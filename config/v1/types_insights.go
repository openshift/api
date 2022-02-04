package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// Insights holds cluster-wide information about Insights.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type Insights struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec InsightsSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status InsightsStatus `json:"status"`
}

type InsightsSpec struct {
	// GatheringConfig spec attribute includes all the configuration options related to
	// gathering of the Insights archive and its uploading to the ingress.
	// +optional
	GatheringConfig *GatheringConfig `json:"gatheringConfig,omitempty"`
}

type InsightsStatus struct {
}

type GatheringConfig struct {
	// DataPolicy allows user to enable additional global obfuscation of the IP addresses and base domain
	// in the Insights archive data.
	// +kubebuilder:default=NoPolicy
	DataPolicy DataPolicy `json:"dataPolicy"`
	// ForceGatherReason enables user to force Insights data gathering by setting a new reason.
	// When there is some gathering in the progress then it is interrupted.
	// When all the gatherers are deactivated by the `DisabledGatherers`, nothing happens.
	// When the forced gathering is finished then the value is cleared.
	ForceGatherReason string `json:"forceGatherReason"`
	// List of gatherers to be excluded from the gathering. All the gatherers can be disabled by providing "all" value.
	// If all the gatherers are disabled, the Insights operator does not gather any data.
	// The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// An example of disabling gatherers looks like this: `disabledGatherers: ["clusterconfig/machine_configs", "workloads/workload_info"]`
	DisabledGatherers []string `json:"disabledGatherers"`
}

const (
	// No data obfuscation
	NoPolicy DataPolicy = "NoPolicy"
	// IP addresses and cluster domain name is obfuscated
	IPsAndClusterDomainPolicy DataPolicy = "IPsAndClusterDomainPolicy"
)

// +kubebuilder:validation:Enum="";NoPolicy;IPsAndClusterDomainPolicy
type DataPolicy string

// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type InsightsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Insights `json:"items"`
}
