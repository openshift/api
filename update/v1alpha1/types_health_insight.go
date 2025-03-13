package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HealthInsight is a piece of actionable information produced by an insight producer about the health
// of the cluster in the context of an update
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=healthinsights,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Reports a piece of actionable information about the health of the cluster in the context of an update"
// +kubebuilder:metadata:annotations="displayName=HealthInsights"
// HealthInsight is a piece of actionable information produced by an insight producer about the health
// of the cluster in the context of an update
type HealthInsight struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, HealthInsight is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +required
	Spec HealthInsightSpec `json:"spec"`
	// status reports a piece of actionable information produced by an insight producer about the health
	// of the cluster in the context of an update
	// +optional
	Status HealthInsightStatus `json:"status"`
}

// HealthInsightSpec is empty for now, HealthInsightSpec is purely status-reporting API. In the future spec may be used
// to hold configuration to drive what information is surfaced and how
type HealthInsightSpec struct {
}

// HealthInsightStatus reports a piece of actionable information produced by an insight producer about the health
// of the cluster in the context of an update
type HealthInsightStatus struct {
	// startedAt is the time when the condition reported by the insight started
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	StartedAt metav1.Time `json:"startedAt"`

	// scope is list of objects involved in the insight
	// +required
	Scope InsightScope `json:"scope"`

	// impact describes the impact the reported condition has on the cluster or update
	// +required
	Impact InsightImpact `json:"impact"`

	// remediation contains information about how to resolve or prevent the reported condition
	// +required
	Remediation InsightRemediation `json:"remediation"`
}

// InsightScope is a list of resources involved in the insight
type InsightScope struct {
	// type is either ControlPlane or WorkerPool
	// +required
	Type ScopeType `json:"type"`

	// resources is a list of resources involved in the insight, of any group/kind. Maximum 16 resources can be listed.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=16
	Resources []ResourceRef `json:"resources,omitempty"`
}

// InsightImpact describes the impact the reported condition has on the cluster or update
type InsightImpact struct {
	// level is the severity of the impact. Valid values are Unknown, Info, Warning, Error, Critical.
	// +required
	Level InsightImpactLevel `json:"level"`

	// type is the type of the impact. Valid values are None, Unknown, API Availability, Cluster Capacity,
	// Application Availability, Application Outage, Data Loss, Update Speed, Update Stalled.
	// +required
	Type InsightImpactType `json:"type"`

	// summary is a short summary of the impact. It must not be empty and must be shorter than 256 characters.
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:MinLength=1
	Summary string `json:"summary"`

	// description is a human-oriented, possibly longer-form description of the condition reported by the insight It must
	// be shorter than 4096 characters.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=4096
	Description string `json:"description,omitempty"`
}

// InsightImpactLevel describes the severity of the impact the reported condition has on the cluster or update
// +kubebuilder:validation:Enum=Unknown;Info;Warning;Error;Critical
type InsightImpactLevel string

const (
	// UnknownImpactLevel is used when the impact level is not known
	UnknownImpactLevel InsightImpactLevel = "Unknown"
	// info should be used for insights that are strictly informational or even positive (things go well or
	// something recently healed)
	InfoImpactLevel InsightImpactLevel = "Info"
	// warning should be used for insights that explain a minor or transient problem. Anything that requires
	// admin attention or manual action should not be a warning but at least an error.
	WarningImpactLevel InsightImpactLevel = "Warning"
	// error should be used for insights that inform about a problem that requires admin attention. Insights of
	// level error and higher should be as actionable as possible, and should be accompanied by links to documentation,
	// KB articles or other resources that help the admin to resolve the problem.
	ErrorImpactLevel InsightImpactLevel = "Error"
	// critical should be used rarely, for insights that inform about a severe problem, threatening with data
	// loss, destroyed cluster or other catastrophic consequences. Insights of this level should be accompanied by
	// links to documentation, KB articles or other resources that help the admin to resolve the problem, or at least
	// prevent the severe consequences from happening.
	CriticalInfoLevel InsightImpactLevel = "Critical"
)

// InsightImpactType describes the type of the impact the reported condition has on the cluster or update
// +kubebuilder:validation:Enum=None;Unknown;API Availability;Cluster Capacity;Application Availability;Application Outage;Data Loss;Update Speed;Update Stalled
type InsightImpactType string

const (
	NoneImpactType                    InsightImpactType = "None"
	UnknownImpactType                 InsightImpactType = "Unknown"
	ApiAvailabilityImpactType         InsightImpactType = "API Availability"
	ClusterCapacityImpactType         InsightImpactType = "Cluster Capacity"
	ApplicationAvailabilityImpactType InsightImpactType = "Application Availability"
	ApplicationOutageImpactType       InsightImpactType = "Application Outage"
	DataLossImpactType                InsightImpactType = "Data Loss"
	UpdateSpeedImpactType             InsightImpactType = "Update Speed"
	UpdateStalledImpactType           InsightImpactType = "Update Stalled"
)

// InsightRemediation contains information about how to resolve or prevent the reported condition
type InsightRemediation struct {
	// reference is a URL where administrators can find information to resolve or prevent the reported condition
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=512
	// +kubebuilder:validation:XValidation:rule="isURL(self)",message="reference must a valid URL"
	Reference string `json:"reference"`

	// estimatedFinish is the estimated time when the informer expects the condition to be resolved, if applicable.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	EstimatedFinish *metav1.Time `json:"estimatedFinish,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HealthInsightList is a list of HealthInsightList resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type HealthInsightList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a list of HealthInsight resources
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Items []HealthInsight `json:"items"`
}
