package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpdateStatus is the API about in-progress updates, kept populated by Update Status Controller by
// aggregating and summarizing UpdateInformers
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=updatestatuses,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=TODO
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +openshift:compatibility-gen:level=4
type UpdateStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec UpdateStatusSpec `json:"spec"`
	// +optional
	Status UpdateStatusStatus `json:"status,omitempty"`
}

// UpdateStatusSpec is empty for now, can possibly hold configuration for Update Status Controller in the future
type UpdateStatusSpec struct {
}

// +k8s:deepcopy-gen=true

// UpdateStatusStatus is the API about in-progress updates, kept populated by Update Status Controller by
// aggregating and summarizing UpdateInformers
type UpdateStatusStatus struct {
	// ControlPlaneUpdateStatus contains a summary and insights related to the control plane update
	ControlPlane ControlPlaneUpdateStatus `json:"controlPlane"`

	// WorkerPoolsUpdateStatus contains summaries and insights related to the worker pools update
	WorkerPools []PoolUpdateStatus `json:"workerPools"`

	// Conditions provide details about Update Status Controller operational matters
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type ControlPlaneConditionType string

const (
	ControlPlaneConditionTypeUpdating ControlPlaneConditionType = "Updating"
)

type ControlPlaneConditionUpdatingReason string

const (
	ControlPlaneConditionUpdatingReasonClusterVersionProgressing        ControlPlaneConditionUpdatingReason = "ClusterVersionProgressing"
	ControlPlaneConditionUpdatingReasonClusterVersionNotProgressing     ControlPlaneConditionUpdatingReason = "ClusterVersionNotProgressing"
	ControlPlaneConditionUpdatingReasonClusterVersionProgressingUnknown ControlPlaneConditionUpdatingReason = "ClusterVersionProgressingUnknown"
	ControlPlaneConditionUpdatingReasonClusterVersionWithoutProgressing ControlPlaneConditionUpdatingReason = "ClusterVersionWithoutProgressing"
)

// ControlPlaneUpdateStatus contains a summary and insights related to the control plane update
type ControlPlaneUpdateStatus struct {
	// Informers is a list of insight producers, each carries a list of insights
	// +listType=map
	// +listMapKey=name
	Informers []UpdateInformer `json:"informers,omitempty"`

	// Conditions provides details about the control plane update
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// UpdateInformer is an insight producer identified by a name, carrying a list of insights it produced
type UpdateInformer struct {
	// Name is the name of the insight producer
	// +required
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Insights is a list of insights produced by this producer
	Insights []UpdateInsight `json:"insights,omitempty"`
}

type ControlPlaneUpdateAssessment string

const (
	ControlPlaneUpdateAssessmentProgressing ControlPlaneUpdateAssessment = "Progressing"
	ControlPlaneUpdateAssessmentCompleted   ControlPlaneUpdateAssessment = "Completed"
	ControlPlaneUpdateAssessmentDegraded    ControlPlaneUpdateAssessment = "Degraded"
)

type ClusterVersionStatusInsightConditionType string

const (
	ClusterVersionStatusInsightConditionTypeUpdating ClusterVersionStatusInsightConditionType = "Updating"
)

type ClusterVersionStatusInsightUpdatingReason string

const (
	ClusterVersionStatusInsightUpdatingReasonNoProgressing ClusterVersionStatusInsightUpdatingReason = "MissingProgressingCondition"
)

// ControlPlaneUpdateVersions contains the original and target versions of the upgrade
type ControlPlaneUpdateVersions struct {
	// Previous is the version of the control plane before the update
	Previous string `json:"previous,omitempty"`

	// IsPreviousPartial is true if the update was initiated in a state where the previous upgrade (to the original version)
	// was not fully completed
	IsPreviousPartial bool `json:"previousPartial,omitempty"`

	// Target is the version of the control plane after the update
	Target string `json:"target"`

	// IsTargetInstall is true if the current (or last completed) work is an installation, not an upgrade
	IsTargetInstall bool `json:"targetInstall,omitempty"`
}

type ClusterVersionStatusInsight struct {
	// Resource is the ClusterVersion resource that represents the control plane
	Resource ResourceRef `json:"resource"`

	// Assessment is the assessment of the control plane update process
	Assessment ControlPlaneUpdateAssessment `json:"assessment"`

	// Versions contains the original and target versions of the upgrade
	Versions ControlPlaneUpdateVersions `json:"versions"`

	// Completion is a percentage of the update completion (0-100)
	Completion int32 `json:"completion"`

	// StartedAt is the time when the update started
	StartedAt metav1.Time `json:"startedAt"`

	// CompletedAt is the time when the update completed
	CompletedAt metav1.Time `json:"completedAt"`

	// EstimatedCompletedAt is the estimated time when the update will complete
	EstimatedCompletedAt metav1.Time `json:"estimatedCompletedAt"`

	// Conditions provides details about the control plane update
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
type ClusterOperatorStatusInsightConditionType string

const (
	ClusterOperatorStatusInsightConditionTypeUpdating ClusterOperatorStatusInsightConditionType = "Updating"
	ClusterOperatorStatusInsightConditionTypeHealthy  ClusterOperatorStatusInsightConditionType = "Healthy"
)

type ClusterOperatorStatusInsightUpdatingReason string

const (
	ClusterOperatorStatusInsightUpdatingReasonUpdated        ClusterOperatorStatusInsightUpdatingReason = "Updated"
	ClusterOperatorStatusInsightUpdatingReasonPending        ClusterOperatorStatusInsightUpdatingReason = "Pending"
	ClusterOperatorStatusInsightUpdatingReasonProgressing    ClusterOperatorStatusInsightUpdatingReason = "Progressing"
	ClusterOperatorStatusInsightUpdatingReasonUnknownUpdate  ClusterOperatorStatusInsightUpdatingReason = "UnclearClusterState"
	ClusterOperatorStatusInsightUpdatingReasonUnknownVersion ClusterOperatorStatusInsightUpdatingReason = "UnknownVersion"
)

type ClusterOperatorStatusInsightHealthyReason string

const (
	ClusterOperatorUpdateStatusInsightHealthyReasonAllIsWell        ClusterOperatorStatusInsightHealthyReason = "AllIsWell"
	ClusterOperatorUpdateStatusInsightHealthyReasonUnavailable      ClusterOperatorStatusInsightHealthyReason = "Unavailable"
	ClusterOperatorUpdateStatusInsightHealthyReasonDegraded         ClusterOperatorStatusInsightHealthyReason = "Degraded"
	ClusterOperatorUpdateStatusInsightHealthyReasonMissingAvailable ClusterOperatorStatusInsightHealthyReason = "MissingAvailable"
	ClusterOperatorUpdateStatusInsightHealthyReasonMissingDegraded  ClusterOperatorStatusInsightHealthyReason = "MissingDegraded"
)

type ClusterOperatorStatusInsight struct {
	// Name is the name of the operator
	Name string `json:"name"`

	// Resource is the ClusterOperator resource that represents the operator
	Resource ResourceRef `json:"resource"`

	// Conditions provide details about the operator
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// PoolUpdateStatus contains a summary and insights related to a node pool update
type PoolUpdateStatus struct {
	// Name is the name of the pool
	Name string `json:"name"`

	// Resource is the resource that represents the pool
	Resource PoolResourceRef `json:"resource"`

	// Informers is a list of insight producers, each carries a list of insights
	// +listType=map
	// +listMapKey=name
	Informers []UpdateInformer `json:"informers,omitempty"`

	// Conditions provide details about the pool
}

type PoolUpdateAssessment string

const (
	PoolUpdateAssessmentPending     PoolUpdateAssessment = "Pending"
	PoolUpdateAssessmentCompleted   PoolUpdateAssessment = "Completed"
	PoolUpdateAssessmentDegraded    PoolUpdateAssessment = "Degraded"
	PoolUpdateAssessmentExcluded    PoolUpdateAssessment = "Excluded"
	PoolUpdateAssessmentProgressing PoolUpdateAssessment = "Progressing"
)

type PoolNodesSummaryType string

const (
	PoolNodesSummaryTypeTotal       PoolNodesSummaryType = "Total"
	PoolNodesSummaryTypeAvailable   PoolNodesSummaryType = "Available"
	PoolNodesSummaryTypeProgressing PoolNodesSummaryType = "Progressing"
	PoolNodesSummaryTypeOutdated    PoolNodesSummaryType = "Outdated"
	PoolNodesSummaryTypeDraining    PoolNodesSummaryType = "Draining"
	PoolNodesSummaryTypeExcluded    PoolNodesSummaryType = "Excluded"
	PoolNodesSummaryTypeDegraded    PoolNodesSummaryType = "Degraded"
)

type PoolNodesUpdateSummary struct {
	// Type is the type of the summary
	// +required
	// +kubebuilder:validation:Required
	Type PoolNodesSummaryType `json:"type"`

	// Count is the number of nodes matching the criteria
	Count int32 `json:"count"`
}

type MachineConfigPoolStatusInsight struct {
	// Name is the name of the machine config pool
	Name string `json:"name"`

	// Resource is the MachineConfigPool resource that represents the pool
	Resource ResourceRef `json:"resource"`

	// Scope describes whether the pool is a control plane or a worker pool
	Scope ScopeType `json:"scopeType"`

	// Assessment is the assessment of the machine config pool update process
	Assessment PoolUpdateAssessment `json:"assessment"`

	// Completion is a percentage of the update completion (0-100)
	Completion int32 `json:"completion"`

	// Summaries is a list of counts of nodes matching certain criteria (e.g. updated, degraded, etc.)
	// +listType=map
	// +listMapKey=type
	Summaries []PoolNodesUpdateSummary `json:"summaries,omitempty"`

	// Conditions provide details about the machine config pool
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type NodeStatusInsightConditionType string

const (
	NodeStatusInsightConditionTypeUpdating  NodeStatusInsightConditionType = "Updating"
	NodeStatusInsightConditionTypeDegraded  NodeStatusInsightConditionType = "Degraded"
	NodeStatusInsightConditionTypeAvailable NodeStatusInsightConditionType = "Available"
)

type NodeStatusInsightUpdatingReason string

const (
	// Updating=True reasons

	NodeStatusInsightUpdatingReasonDraining  NodeStatusInsightUpdatingReason = "Draining"
	NodeStatusInsightUpdatingReasonUpdating  NodeStatusInsightUpdatingReason = "Updating"
	NodeStatusInsightUpdatingReasonRebooting NodeStatusInsightUpdatingReason = "Rebooting"

	// Updating=False reasons

	NodeStatusInsightUpdatingReasonPaused    NodeStatusInsightUpdatingReason = "Paused"
	NodeStatusInsightUpdatingReasonPending   NodeStatusInsightUpdatingReason = "Pending"
	NodeStatusInsightUpdatingReasonCompleted NodeStatusInsightUpdatingReason = "Completed"
)

type NodeStatusInsight struct {
	// Name is the name of the node
	Name string `json:"name"`

	// Resource is the Node resource that represents the node
	Resource ResourceRef `json:"resource"`

	// PoolResource is the resource that represents the pool the node is a member of
	PoolResource PoolResourceRef `json:"poolResource"`

	// Version is the version of the node, when known
	Version string `json:"version,omitempty"`

	// EstToComplete is the estimated time to complete the update, when known
	EstToComplete metav1.Duration `json:"estToComplete,omitempty"`

	// Message is a human-readable message about the node update status
	Message string `json:"message,omitempty"`

	// Conditions provides details about the control plane update
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type UpdateInsightType string

const (
	UpdateInsightTypeClusterVersionStatusInsight    UpdateInsightType = "ClusterVersion"
	UpdateInsightTypeClusterOperatorStatusInsight   UpdateInsightType = "ClusterOperator"
	UpdateInsightTypeMachineConfigPoolStatusInsight UpdateInsightType = "MachineConfigPool"
	UpdateInsightTypeNodeStatusInsight              UpdateInsightType = "Node"
	UpdateInsightTypeUpdateHealthInsight            UpdateInsightType = "UpdateHealth"
)

type UpdateInsight struct {
	// +unionDiscriminator
	Type UpdateInsightType `json:"type"`

	// UID identifies an insight over time
	UID string `json:"uid"`

	// AcquiredAt is the time when the data was acquired by the producer
	AcquiredAt metav1.Time `json:"acquisitionTime"`

	// ClusterVersionStatusInsight is a status insight about the state of a control plane update, where
	// the control plane is represented by a ClusterVersion resource usually managed by CVO
	// +optional
	ClusterVersionStatusInsight *ClusterVersionStatusInsight `json:"cv,omitempty"`

	// ClusterOperatorStatusInsight is a status insight about the state of a control plane cluster operator update
	// represented by a ClusterOperator resource
	// +optional
	ClusterOperatorStatusInsight *ClusterOperatorStatusInsight `json:"co,omitempty"`

	// MachineConfigPoolStatusInsight is a status insight about the state of a worker pool update, where the worker pool
	// is represented by a MachineConfigPool resource
	// +optional
	MachineConfigPoolStatusInsight *MachineConfigPoolStatusInsight `json:"mcp,omitempty"`

	// NodeStatusInsight is a status insight about the state of a worker node update, where the worker node is represented
	// by a Node resource
	// +optional
	NodeStatusInsight *NodeStatusInsight `json:"node,omitempty"`

	// UpdateHealthInsight is a generic health insight about the update. It does not represent a status of any specific
	// resource but surfaces actionable information about the health of the cluster or an update
	// +optional
	UpdateHealthInsight *UpdateHealthInsight `json:"health,omitempty"`
}

// UpdateHealthInsight is a piece of actionable information produced by an insight producer about the health
// of the cluster or an update
type UpdateHealthInsight struct {
	// StartedAt is the time when the condition reported by the insight started
	StartedAt metav1.Time `json:"startedAt"`

	// Scope is list of objects involved in the insight
	// +optional
	Scope UpdateInsightScope `json:"scope,omitempty"`

	// Impact describes the impact the reported condition has on the cluster or update
	Impact UpdateInsightImpact `json:"impact"`

	// Remediation contains ... TODO
	Remediation UpdateInsightRemediation `json:"remediation"`
}

// ScopeType is one of ControlPlane or WorkerPool
// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
type ScopeType string

const (
	ScopeTypeControlPlane ScopeType = "ControlPlane"
	ScopeTypeWorkerPool   ScopeType = "WorkerPool"
)

// UpdateInsightScope is a list of objects involved in the insight
type UpdateInsightScope struct {
	// Type is either ControlPlane or WorkerPool
	// +kubebuilder:validation:Required
	Type ScopeType `json:"type"`

	// Resources is a list of resources involved in the insight
	// +optional
	Resources []ResourceRef `json:"resources,omitempty"`
}

// ResourceRef is a reference to a kubernetes resource, typically involved in an
// insight
type ResourceRef struct {
	// Kind of object being referenced
	Kind string `json:"kind"`

	// APIGroup of the object being referenced
	// +optional
	APIGroup string `json:"apiGroup,omitempty"`

	// Name of the object being referenced
	Name string `json:"name"`

	// Namespace of the object being referenced, if any
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// InsightImpactLevel describes the severity of the impact the reported condition has on the cluster or update
// +kubebuilder:validation:Enum=info;warning;error;critical
type InsightImpactLevel string

const (
	// InfoImpactLevel should be used for insights that are strictly informational or even positive (things go well or
	// something recently healed)
	InfoImpactLevel InsightImpactLevel = "info"
	// WarningImpactLevel should be used for insights that explain a minor or transient problem. Anything that requires
	// admin attention or manual action should not be a warning but at least an error.
	WarningImpactLevel InsightImpactLevel = "warning"
	// ErrorImpactLevel should be used for insights that inform about a problem that requires admin attention. Insights of
	// level error and higher should be as actionable as possible, and should be accompanied by links to documentation,
	// KB articles or other resources that help the admin to resolve the problem.
	ErrorImpactLevel InsightImpactLevel = "error"
	// CriticalInfoLevel should be used rarely, for insights that inform about a severe problem, threatening with data
	// loss, destroyed cluster or other catastrophic consequences. Insights of this level should be accompanied by
	// links to documentation, KB articles or other resources that help the admin to resolve the problem, or at least
	// prevent the severe consequences from happening.
	CriticalInfoLevel InsightImpactLevel = "critical"
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

// UpdateInsightImpact describes the impact the reported condition has on the cluster or update
type UpdateInsightImpact struct {
	// Level is the severity of the impact
	Level InsightImpactLevel `json:"level"`

	// Type is the type of the impact
	Type InsightImpactType `json:"type"`

	// Summary is a short summary of the impact
	Summary string `json:"summary"`

	// Description is a human-oriented description of the condition reported by the insight
	Description string `json:"description"`
}

// UpdateInsightRemediation contains ... TODO
type UpdateInsightRemediation struct {
	// Reference is a URL where administrators can find information to resolve or prevent the reported condition
	Reference string `json:"reference"`

	// EstimatedFinish is the estimated time when the informer expects the condition to be resolved, if applicable.
	// This should normally only be provided by system level insights (impact level=status)
	EstimatedFinish metav1.Time `json:"estimatedFinish"`
}

// PoolResourceRef is a reference to a kubernetes resource that represents a worker pool
type PoolResourceRef struct {
	ResourceRef `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpdateStatusList is a list of UpdateStatus resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type UpdateStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []UpdateStatus `json:"items"`
}
