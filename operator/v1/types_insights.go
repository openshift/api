package v1

import (
	corev1 "k8s.io/api/core/v1"
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
	// GatheringStatus provides basic information about the last Insights gathering.
	GatheringStatus *GatheringStatus `json:"gatheringStatus,omitempty"`
	// ReportStatus provides general Insights analysis results.
	ReportStatus *ReportStatus `json:"reportStatus,omitempty"`
}

type GatheringStatus struct {
	// LastGatherTime is the last time when Insights gathering finished.
	LastGatherTime metav1.Time `json:"lastGatherTime,omitempty"`
	// LastGatherReason provides last known reason of gathering. This is helpful
	// especially when gathering was forced by user
	LastGatherReason string `json:"lastGatherReason,omitempty"`
	// StartGatherTime is the time when gathering started. The value is 0
	// when there is no gathering in progress.
	StartGatherTime metav1.Time `json:"startGatherTime,omitempty"`
	// List of active gatherers (and their statuses) in the last gathering.
	GathererStatuses []GathererStatus `json:"gathererStatuses,omitempty"`
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
	GathererConditions []GathererCondition `json:"conditions"`
	// DurationMillisecond represents the time spent gathering.
	DurationMillisecond int64 `json:"durationMillisecond"`
}

type GathererCondition struct {
	// Type of the gatherer condition
	Type GathererConditionType `json:"type"`
	// Status is last known status of the particular gatherer condition
	Status corev1.ConditionStatus `json:"status"`
	// Messages is an optional attribute that provides error and warning messages
	// from the gatherer
	Messages []string `json:"messages,omitempty"`
	// Last time the condition transit from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
}

// GathererConditionType is a valid value for GathererCondition.Type
// +kubebuilder:validation:Enum="";Successful;Disabled;Failed;Warning
type GathererConditionType string

const (
	// Gatherer was successful withnout any error reported
	Successful GathererConditionType = "Successful"
	// Gatherer was disabled by user
	Disabled GathererConditionType = "Disabled"
	// Gatherer failed to run
	Failed GathererConditionType = "Failed"
	// Gatherer was running, but there were some errors
	Warning GathererConditionType = "Warning"
)

// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type InsightsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Insights `json:"items"`
}
