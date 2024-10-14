package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpdateStatus reports status for in-progress cluster version updates
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=updatestatuses,scope=Namespaced
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Provides health and status information about OpenShift cluster updates."
// +kubebuilder:metadata:annotations="displayName=UpdateStatuses"
type UpdateStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, UpdateStatus is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +kubebuilder:validation:Required
	Spec UpdateStatusSpec `json:"spec"`
	// +optional
	Status UpdateStatusStatus `json:"status"`
}

// UpdateStatusSpec is empty for now, UpdateStatus is purely status-reporting API. In the future spec may be used
// to hold configuration to drive what information is surfaced and how
type UpdateStatusSpec struct {
}

// +k8s:deepcopy-gen=true

// UpdateStatusStatus is the API about in-progress updates. It aggregates and summarizes UpdateInsights produced by
// update informers
type UpdateStatusStatus struct {
	// controlPlane contains a summary and insights related to the control plane update
	// +kubebuilder:validation:Required
	ControlPlane ControlPlaneUpdateStatus `json:"controlPlane"`

	// workerPools contains summaries and insights related to the worker pools update
	// +listType=map
	// +listMapKey=name
	// +optional
	WorkerPools []PoolUpdateStatus `json:"workerPools,omitempty"`

	// conditions provide details about the controller operational matters
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ControlPlaneConditionType are types of conditions that can be reported on control plane level
type ControlPlaneConditionType string

const (
	// Updating is the condition type that communicate whether the whole control plane is updating or not
	ControlPlaneUpdating ControlPlaneConditionType = "Updating"
)

// ControlPlaneUpdatingReason are well-known reasons for the Updating condition
// +kubebuilder:validation:Enum=ClusterVersionProgressing;ClusterVersionNotProgressing;CannotDetermineUpdating
type ControlPlaneUpdatingReason string

const (
	// ClusterVersionProgressing is used for Updating=True set because we observed a ClusterVersion resource to
	// have Progressing=True condition
	ReasonClusterVersionProgressing ControlPlaneUpdatingReason = "ClusterVersionProgressing"
	// ClusterVersionNotProgressing is used for Updating=False set because we observed a ClusterVersion resource to
	// have Progressing=False condition
	ReasonClusterVersionNotProgressing ControlPlaneUpdatingReason = "ClusterVersionNotProgressing"
	// CannotDetermineUpdating is used with Updating=Unknown. This covers many different actual reasons such as
	// missing or Unknown Progressing condition on ClusterVersion, but it does not seem useful to track the individual
	// reasons to that granularity for Updating=Unknown
	ReasonClusterVersionCannotDetermine ControlPlaneUpdatingReason = "CannotDetermineUpdating"
)

// ControlPlaneUpdateStatus contains a summary and insights related to the control plane update
type ControlPlaneUpdateStatus struct {
	// resource is the resource that represents the control plane. It will typically be a ClusterVersion resource
	// in standalone OpenShift and HostedCluster in Hosted Control Planes.
	//
	// Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// only the "correct" resource types to be referenced (here, ClusterVersion and HostedCluster). However, because we
	// use resource references in many places and this API is intended to be consumed by clients, not produced, consistency
	// seems to be more valuable than type safety for producers.
	// +kubebuilder:validation:Required
	Resource ResourceRef `json:"resource"`

	// poolResource is the resource that represents control plane node pool, typically a MachineConfigPool. This field
	// is optional because some form factors (like Hosted Control Planes) do not have dedicated control plane node pools.
	//
	// Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// only the "correct" resource types to be referenced (here, MachineConfigPool). However, because we use resource
	// references in many places and this API is intended to be consumed by clients, not produced, consistency seems to be
	// more valuable than type safety for producers.
	// +optional
	PoolResource *PoolResourceRef `json:"poolResource,omitempty"`

	// informers is a list of insight producers, each carries a list of insights relevant for control plane
	// +listType=map
	// +listMapKey=name
	// +optional
	Informers []UpdateInformer `json:"informers,omitempty"`

	// conditions provides details about the control plane update
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// UpdateInformer is an insight producer identified by a name, carrying a list of insights it produced
type UpdateInformer struct {
	// name is the name of the insight producer
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// insights is a list of insights produced by this producer
	// +optional
	// +listType=map
	// +listMapKey=uid
	Insights []UpdateInsight `json:"insights,omitempty"`
}

// ControlPlaneAssessment is the assessment of the control plane update process
type ControlPlaneAssessment string

const (
	// Unknown means the update status and health cannot be determined
	ControlPlaneAssessmentUnknown ControlPlaneAssessment = "Unknown"
	// Progressing means the control plane is updating and no problems or slowness are detected
	ControlPlaneAssessmentProgressing ControlPlaneAssessment = "Progressing"
	// Completed means the control plane successfully completed updating and no problems are detected
	ControlPlaneAssessmentCompleted ControlPlaneAssessment = "Completed"
	// Degraded means the process of updating the control plane suffers from an observed problem
	ControlPlaneAssessmentDegraded ControlPlaneAssessment = "Degraded"
)

// ClusterVersionStatusInsightConditionType are types of conditions that can be reported on ClusterVersion status insight
type ClusterVersionStatusInsightConditionType string

const (
	// Updating condition communicates whether the ClusterVersion is updating
	ClusterVersionStatusInsightUpdating ClusterVersionStatusInsightConditionType = "Updating"
)

// ClusterVersionStatusInsightUpdatingReason are well-known reasons for the Updating condition on ClusterVersion status insights
type ClusterVersionStatusInsightUpdatingReason string

const (
	// CannotDetermineUpdating is used with Updating=Unknown
	ClusterVersionCannotDetermineUpdating ClusterVersionStatusInsightUpdatingReason = "CannotDetermineUpdating"
	// ClusterVersionProgressing means that ClusterVersion is considered to be Updating=True because it has a Progressing=True condition
	ClusterVersionProgressing ClusterVersionStatusInsightUpdatingReason = "ClusterVersionProgressing"
	// ClusterVersionNotProgressing means that ClusterVersion is considered to be Updating=False because it has a Progressing=False condition
	ClusterVersionNotProgressing ClusterVersionStatusInsightUpdatingReason = "ClusterVersionNotProgressing"
)

// VersionMetadataKey is a key for a metadata value associated with a version
// +kubebuilder:validation:Enum=Installation;Partial;Architecture
type VersionMetadataKey string

const (
	// Installation denotes a boolean that indicates the update was initiated as an installation
	InstallationMetadata VersionMetadataKey = "Installation"
	// Partial denotes a boolean that indicates the update was initiated in a state where the previous upgrade
	// (to the original version) was not fully completed
	PartialMetadata VersionMetadataKey = "Partial"
	// Architecture denotes a string that indicates the architecture of the payload image of the version,
	// when relevant
	ArchitectureMetadata VersionMetadataKey = "Architecture"
)

type VersionMetadata struct {
	// key is the name of this metadata value
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=installation;partial;architecture
	Key VersionMetadataKey `json:"key"`

	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=32
	Value string `json:"value,omitempty"`
}

// Version describes a version involved in an update, typically on one side of an update edge
type Version struct {
	// version is a semantic version string, or a placeholder '<none>' for the special case where this
	// is a "previous" version in a new installation
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=64
	Version string `json:"version,omitempty"`

	// metadata is a list of metadata associated with the version. It is a list of key-value pairs. The value is optional
	// and when not provided, the metadata item has boolean semantics (presence indicates true)
	// +listType=map
	// +listMapKey=key
	// +optional
	Metadata []VersionMetadata `json:"metadata,omitempty"`
}

// ControlPlaneUpdateVersions contains the original and target versions of the upgrade
type ControlPlaneUpdateVersions struct {
	// previous is the version of the control plane before the update. When the cluster is being installed
	// for the first time, the version will have a placeholder value like '<none>' and the target version
	// will have a boolean installation=true metadata
	// +kubebuilder:validation:Required
	Previous Version `json:"previous"`

	// target is the version of the control plane after the update
	// +kubebuilder:validation:Required
	Target Version `json:"target"`
}

// ClusterVersionStatusInsight reports the state of a ClusterVersion resource (which represents a control plane
// update in standalone clusters), during the update.
type ClusterVersionStatusInsight struct {
	// resource is the ClusterVersion resource that represents the control plane
	//
	// Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// resource name (because the rest is implied by status insight type). However, because we use resource references in
	// many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// than type safety for producers.
	// +kubebuilder:validation:Required
	Resource ResourceRef `json:"resource"`

	// assessment is the assessment of the control plane update process
	// +kubebuilder:validation:Required
	Assessment ControlPlaneAssessment `json:"assessment"`

	// versions contains the original and target versions of the upgrade
	// +kubebuilder:validation:Required
	Versions ControlPlaneUpdateVersions `json:"versions"`

	// completion is a percentage of the update completion (0-100)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Completion int32 `json:"completion"`

	// startedAt is the time when the update started
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	StartedAt metav1.Time `json:"startedAt"`

	// completedAt is the time when the update completed
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`

	// estimatedCompletedAt is the estimated time when the update will complete
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	EstimatedCompletedAt *metav1.Time `json:"estimatedCompletedAt,omitempty"`

	// conditions provides detailed observed conditions about ClusterVersion
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ClusterOperatorStatusInsightConditionType are types of conditions that can be reported on ClusterOperator status insights
type ClusterOperatorStatusInsightConditionType string

const (
	// Updating condition communicates whether the ClusterOperator is updating
	ClusterOperatorStatusInsightUpdating ClusterOperatorStatusInsightConditionType = "Updating"
	// Healthy condition communicates whether the ClusterOperator is considered healthy
	ClusterOperatorStatusInsightHealthy ClusterOperatorStatusInsightConditionType = "Healthy"
)

// ClusterOperatorUpdatingReason are well-known reasons for the Updating condition on ClusterOperator status insights
type ClusterOperatorUpdatingReason string

const (
	// Updated is used with Updating=False when the ClusterOperator finished updating
	ClusterOperatorUpdatingReasonUpdated ClusterOperatorUpdatingReason = "Updated"
	// Pending is used with Updating=False when the ClusterOperator is not updating and is still running previous version
	ClusterOperatorUpdatingReasonPending ClusterOperatorUpdatingReason = "Pending"
	// Progressing is used with Updating=True when the ClusterOperator is updating
	ClusterOperatorUpdatingReasonProgressing ClusterOperatorUpdatingReason = "Progressing"
	// CannotDetermine is used with Updating=Unknown
	ClusterOperatorUpdatingCannotDetermine ClusterOperatorUpdatingReason = "CannotDetermine"
)

// ClusterOperatorHealthyReason are well-known reasons for the Healthy condition on ClusterOperator status insights
type ClusterOperatorHealthyReason string

const (
	// AsExpected is used with Healthy=True when no issues are observed
	ClusterOperatorHealthyReasonAsExpected ClusterOperatorHealthyReason = "AsExpected"
	// Unavailable is used with Healthy=False when the ClusterOperator has Available=False condition
	ClusterOperatorHealthyReasonUnavailable ClusterOperatorHealthyReason = "Unavailable"
	// Degraded is used with Healthy=False when the ClusterOperator has Degraded=True condition
	ClusterOperatorHealthyReasonDegraded ClusterOperatorHealthyReason = "Degraded"
	// CannotDetermine is used with Healthy=Unknown
	ClusterOperatorHealthyReasonCannotDetermine ClusterOperatorHealthyReason = "CannotDetermine"
)

// ClusterOperatorStatusInsight reports the state of a ClusterOperator resource (which represents a control plane
// component update in standalone clusters), during the update
type ClusterOperatorStatusInsight struct {
	// name is the name of the operator
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// resource is the ClusterOperator resource that represents the operator
	//
	// Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// resource name (because the rest is implied by status insight type). However, because we use resource references in
	// many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// than type safety for producers.
	// +kubebuilder:validation:Required
	Resource ResourceRef `json:"resource"`

	// conditions provide details about the operator
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// PoolUpdateStatus contains a summary and insights related to a node pool update
type PoolUpdateStatus struct {
	// name is the name of the pool
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// resource is the resource that represents the pool
	//
	// Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// only the "correct" resource types to be referenced (here, MachineConfigPool or NodePool). However, because we use
	// resource references in many places and this API is intended to be consumed by clients, not produced, consistency
	// seems to be more valuable than type safety for producers.
	// +kubebuilder:validation:Required
	Resource PoolResourceRef `json:"resource"`

	// informers is a list of insight producers, each carries a list of insights
	// +listType=map
	// +listMapKey=name
	// +optional
	Informers []UpdateInformer `json:"informers,omitempty"`

	// conditions provide details about the pool
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// PoolUpdateAssessment is the assessment of the node pool update process
type PoolUpdateAssessment string

const (
	// Pending means the nodes in the pool will be updated but none have even started yet
	PoolUpdatePending PoolUpdateAssessment = "Pending"
	// Completed means all nodes in the pool have been updated
	PoolUpdateCompleted PoolUpdateAssessment = "Completed"
	// Degraded means the process of updating the pool suffers from an observed problem
	PoolUpdateDegraded PoolUpdateAssessment = "Degraded"
	// Excluded means some (or all) nodes in the pool would be normally updated but a configuration (such as paused MCP)
	// prevents that from happening
	PoolUpdateExcluded PoolUpdateAssessment = "Excluded"
	// Progressing means the nodes in the pool are being updated and no problems or slowness are detected
	PoolUpdateProgressing PoolUpdateAssessment = "Progressing"
)

// NodeSummaryType are types of summaries (how many nodes match certain criteria, such as updated, degraded, etc.)
// reported for a node pool
// +kubebuilder:validation:Enum=Total;Available;Progressing;Outdated;Draining;Excluded;Degraded
type NodeSummaryType string

const (
	// Total is the total number of nodes in the pool
	NodesTotal NodeSummaryType = "Total"
	// Available is the number of nodes in the pool that are available (accepting workloads)
	NodesAvailable NodeSummaryType = "Available"
	// Progressing is the number of nodes in the pool that are updating
	NodesProgressing NodeSummaryType = "Progressing"
	// Outdated is the number of nodes in the pool that are running an outdated version
	NodesOutdated NodeSummaryType = "Outdated"
	// Draining is the number of nodes in the pool that are being drained
	NodesDraining NodeSummaryType = "Draining"
	// Excluded is the number of nodes in the pool that would normally be updated but configuration (such as paused MCP)
	// prevents that from happening
	NodesExcluded NodeSummaryType = "Excluded"
	// Degraded is the number of nodes in the pool that are degraded
	NodesDegraded NodeSummaryType = "Degraded"
)

// NodeSummary is a count of nodes matching certain criteria (e.g. updated, degraded, etc.)
type NodeSummary struct {
	// type is the type of the summary
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Total;Available;Progressing;Outdated;Draining;Excluded;Degraded
	Type NodeSummaryType `json:"type"`

	// count is the number of nodes matching the criteria
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	Count int32 `json:"count"`
}

// ClusterVersionStatusInsight reports the state of a MachineConfigPool resource during the update
type MachineConfigPoolStatusInsight struct {
	// name is the name of the machine config pool
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// resource is the MachineConfigPool resource that represents the pool
	//
	// Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// resource name (because the rest is implied by status insight type). However, because we use resource references in
	// many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// than type safety for producers.
	// +kubebuilder:validation:Required
	Resource PoolResourceRef `json:"resource"`

	// scopeType describes whether the pool is a control plane or a worker pool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
	Scope ScopeType `json:"scopeType"`

	// assessment is the assessment of the machine config pool update process
	// +kubebuilder:validation:Required
	Assessment PoolUpdateAssessment `json:"assessment"`

	// completion is a percentage of the update completion (0-100)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Completion int32 `json:"completion"`

	// summaries is a list of counts of nodes matching certain criteria (e.g. updated, degraded, etc.)
	// +listType=map
	// +listMapKey=type
	// +optional
	Summaries []NodeSummary `json:"summaries,omitempty"`

	// conditions provide details about the machine config pool update
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NodeStatusInsightConditionType are types of conditions that can be reported on Node status insights
type NodeStatusInsightConditionType string

const (
	// Updating condition communicates whether the Node is updating
	NodeStatusInsightUpdating NodeStatusInsightConditionType = "Updating"
	// Degraded condition communicates whether the Node is degraded (problem observed)
	NodeStatusInsightDegraded NodeStatusInsightConditionType = "Degraded"
	// Available condition communicates whether the Node is available (accepting workloads)
	NodeStatusInsightAvailable NodeStatusInsightConditionType = "Available"
)

// NodeUpdatingReason are well-known reasons for the Updating condition on Node status insights
type NodeUpdatingReason string

const (
	// Draining is used with Updating=True when the Node is being drained
	NodeDraining NodeUpdatingReason = "Draining"
	// Updating is used with Updating=True when new node configuration is being applied
	NodeUpdating NodeUpdatingReason = "Updating"
	// Rebooting is used with Updating=True when the Node is rebooting into the new version
	NodeRebooting NodeUpdatingReason = "Rebooting"

	// Updated is used with Updating=False when the Node is prevented by configuration from updating
	NodePaused NodeUpdatingReason = "Paused"
	// Updated is used with Updating=False when the Node is waiting to be eventually updated
	NodeUpdatePending NodeUpdatingReason = "Pending"
	// Updated is used with Updating=False when the Node has been updated
	NodeCompleted NodeUpdatingReason = "Completed"

	// CannotDetermine is used with Updating=Unknown
	NodeCannotDetermine NodeUpdatingReason = "CannotDetermine"
)

// NodeStatusInsight reports the state of a Node during the update
type NodeStatusInsight struct {
	// name is the name of the node
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// resource is the Node resource that represents the node
	//
	// Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// resource name (because the rest is implied by status insight type). However, because we use resource references in
	// many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// than type safety for producers.
	// +kubebuilder:validation:Required
	Resource ResourceRef `json:"resource"`

	// poolResource is the resource that represents the pool the node is a member of
	//
	// Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// only the "correct" resource types to be referenced (here, MachineConfigPool or NodePool). However, because we use
	// resource references in many places and this API is intended to be consumed by clients, not produced, consistency
	// seems to be more valuable than type safety for producers.
	// +kubebuilder:validation:Required
	PoolResource PoolResourceRef `json:"poolResource"`

	// version is the version of the node, when known
	// +optional
	// +kubebuilder:validation:Type=string
	Version string `json:"version,omitempty"`

	// estToComplete is the estimated time to complete the update, when known
	// +optional
	// +kubebuilder:validation:Type=string
	EstToComplete *metav1.Duration `json:"estToComplete,omitempty"`

	// message is a short human-readable message about the node update status
	// +optional
	Message string `json:"message,omitempty"`

	// conditions provides details about the control plane update
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// UpdateInsightType identifies the type of the update insight as either one of the resource-specific status insight,
// or a generic health insight
// +kubebuilder:validation:Enum=ClusterVersion;ClusterOperator;MachineConfigPool;Node;UpdateHealth
type UpdateInsightType string

const (
	// Resource-specific status insights should be reported continuously during the update process and mostly communicate
	// progress and high-level state

	// ClusterVersion status insight reports progress and high-level state of a ClusterVersion resource, representing
	// control plane in standalone clusters
	ClusterVersionStatusInsightType UpdateInsightType = "ClusterVersion"
	// ClusterOperator status insight reports progress and high-level state of a ClusterOperator, representing a control
	// plane component
	ClusterOperatorStatusInsightType UpdateInsightType = "ClusterOperator"
	// MachineConfigPool status insight reports progress and high-level state of a MachineConfigPool resource, representing
	// a pool of nodes in clusters using Machine API
	MachineConfigPoolStatusInsightType UpdateInsightType = "MachineConfigPool"
	// Node status insight reports progress and high-level state of a Node resource, representing a node (both control
	// plane and worker) in a cluster
	NodeStatusInsightType UpdateInsightType = "Node"

	// Health insights are reported only when an informer observes a condition that requires admin attention
	UpdateHealthInsightType UpdateInsightType = "UpdateHealth"
)

type UpdateInsight struct {
	// uid identifies the insight over time
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	UID string `json:"uid"`

	// acquiredAt is the time when the data was acquired by the producer
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	AcquiredAt metav1.Time `json:"acquiredAt"`

	UpdateInsightUnion `json:",inline"`
}

type UpdateInsightUnion struct {
	// type identifies the type of the update insight
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=ClusterVersion;ClusterOperator;MachineConfigPool;Node;UpdateHealth
	Type UpdateInsightType `json:"type"`

	// clusterVersion is a status insight about the state of a control plane update, where
	// the control plane is represented by a ClusterVersion resource usually managed by CVO
	// +optional
	// +unionMember
	ClusterVersionStatusInsight *ClusterVersionStatusInsight `json:"clusterVersion,omitempty"`

	// clusterOperator is a status insight about the state of a control plane cluster operator update
	// represented by a ClusterOperator resource
	// +optional
	// +unionMember
	ClusterOperatorStatusInsight *ClusterOperatorStatusInsight `json:"clusterOperator,omitempty"`

	// machineConfigPool is a status insight about the state of a worker pool update, where the worker pool
	// is represented by a MachineConfigPool resource
	// +optional
	// +unionMember
	MachineConfigPoolStatusInsight *MachineConfigPoolStatusInsight `json:"machineConfigPool,omitempty"`

	// node is a status insight about the state of a worker node update, where the worker node is represented
	// by a Node resource
	// +optional
	// +unionMember
	NodeStatusInsight *NodeStatusInsight `json:"node,omitempty"`

	// health is a generic health insight about the update. It does not represent a status of any specific
	// resource but surfaces actionable information about the health of the cluster or an update
	// +optional
	// +unionMember
	UpdateHealthInsight *UpdateHealthInsight `json:"health,omitempty"`
}

// UpdateHealthInsight is a piece of actionable information produced by an insight producer about the health
// of the cluster or an update
type UpdateHealthInsight struct {
	// startedAt is the time when the condition reported by the insight started
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	StartedAt metav1.Time `json:"startedAt"`

	// scope is list of objects involved in the insight
	// +kubebuilder:validation:Required
	Scope UpdateInsightScope `json:"scope"`

	// impact describes the impact the reported condition has on the cluster or update
	// +kubebuilder:validation:Required
	Impact UpdateInsightImpact `json:"impact"`

	// remediation contains information about how to resolve or prevent the reported condition
	Remediation UpdateInsightRemediation `json:"remediation"`
}

// ScopeType is one of ControlPlane or WorkerPool
// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
type ScopeType string

const (
	// ControlPlane is used for insights that are related to the control plane (including control plane pool or nodes)
	ControlPlaneScope ScopeType = "ControlPlane"
	// WorkerPool is used for insights that are related to a worker pools and nodes (excluding control plane)
	WorkerPoolScope ScopeType = "WorkerPool"
)

// UpdateInsightScope is a list of resources involved in the insight
type UpdateInsightScope struct {
	// type is either ControlPlane or WorkerPool
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
	Type ScopeType `json:"type"`

	// resources is a list of resources involved in the insight, of any group/kind
	// +optional
	// +listType=atomic
	Resources []ResourceRef `json:"resources,omitempty"`
}

// ResourceRef is a reference to a kubernetes resource, typically involved in an insight
type ResourceRef struct {
	// group of the object being referenced, if any
	// +optional
	Group string `json:"group,omitempty"`

	// resource of object being referenced
	// +kubebuilder:validation:Required
	Resource string `json:"resource"`

	// name of the object being referenced
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// namespace of the object being referenced, if any
	// +optional
	Namespace string `json:"namespace,omitempty"`
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

// UpdateInsightImpact describes the impact the reported condition has on the cluster or update
type UpdateInsightImpact struct {
	// level is the severity of the impact
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Unknown;Info;Warning;Error;Critical
	Level InsightImpactLevel `json:"level"`

	// type is the type of the impact
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=None;Unknown;API Availability;Cluster Capacity;Application Availability;Application Outage;Data Loss;Update Speed;Update Stalled
	Type InsightImpactType `json:"type"`

	// summary is a short summary of the impact
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	Summary string `json:"summary"`

	// description is a human-oriented, possibly longer-form description of the condition reported by the insight
	// +optional
	// +kubebuilder:validation:Type=string
	Description string `json:"description,omitempty"`
}

// UpdateInsightRemediation contains information about how to resolve or prevent the reported condition
type UpdateInsightRemediation struct {
	// reference is a URL where administrators can find information to resolve or prevent the reported condition
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=uri
	Reference string `json:"reference"`

	// estimatedFinish is the estimated time when the informer expects the condition to be resolved, if applicable.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	EstimatedFinish *metav1.Time `json:"estimatedFinish,omitempty"`
}

// PoolResourceRef is a reference to a kubernetes resource that represents a node pool
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
