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
	metav1.ObjectMeta `json:"metadata"`

	// spec is the specification of the desired behavior of the Insights
	// +kubebuilder:validation:Required
	// +required
	Spec InsightsSpec `json:"spec"`

	// status is the most recently observed status of the Insights operator
	// +optional
	Status InsightsStatus `json:"status"`
}

type InsightsSpec struct {
	OperatorSpec `json:",inline"`
}

type InsightsStatus struct {
	OperatorStatus `json:",inline"`
	// GatherStatus provides basic information about the last Insights data gathering.
	GatherStatus *GatherStatus `json:"gatherStatus,omitempty"`
	// ReportStatus provides general Insights analysis results.
	ReportStatus *ReportStatus `json:"reportStatus,omitempty"`
}

type GatherStatus struct {
	// LastGatherTime is the last time when Insights data gathering finished.
	LastGatherTime metav1.Time `json:"lastGatherTime,omitempty"`
	// LastGatherReason provides last known reason of data gathering. This is helpful
	// especially when data gathering was forced by user
	LastGatherReason string `json:"lastGatherReason,omitempty"`
	// StartGatherTime is the time when data gathering started. The value is 0
	// when there is no data gathering in progress.
	StartGatherTime metav1.Time `json:"startGatherTime,omitempty"`
	// List of active gatherers (and their statuses) in the last gathering.
	Gatherers []GathererStatus `json:"gatherers,omitempty"`
}

type ReportStatus struct {
	// Number of active Insights healthchecks with low severity
	LowHealthChecksCount int `json:"low"`
	// Number of active Insights healthchecks with moderate severity
	ModerateHealthChecksCount int `json:"moderate"`
	// Number of active Insights healthchecks with important severity
	ImportantHealthChecksCount int `json:"important"`
	// Number of active Insights healthchecks with critical severity
	CriticalHealthChecksCount int `json:"critical"`
	// TotalCount is the count of all active Insights healthchecks
	TotalHealthChecksCount int `json:"total"`
}

type GathererStatus struct {
	// Name is the name of the gatherer.
	Name string `json:"name"`
	// GathererConditions provide details on the status of each gatherer.
	GathererConditions []metav1.Condition `json:"conditions"`
	// Duration represents the time spent gathering.
	Duration metav1.Duration `json:"duration"`
}

// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type InsightsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Insights `json:"items"`
}
