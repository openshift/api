package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPoolProgressInsight provides summary information about an ongoing node pool update in Standalone clusters
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=machineconfigpoolprogressinsights,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Provides summary information about an ongoing node pool update in Standalone clusters"
// +kubebuilder:metadata:annotations="displayName=MachineConfigPoolProgressInsights"
// +kubebuilder:validation:XValidation:rule="self.metadata.name == self.status.name",message="Progress Insight .metadata.name must match .status.name"
// MachineConfigPoolProgressInsight reports the state of a MachineConfigPool resource (which represents a pool of nodes
// update in standalone clusters), during a cluster update.
type MachineConfigPoolProgressInsight struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, MachineConfigPoolProgressInsight is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +required
	Spec MachineConfigPoolProgressInsightSpec `json:"spec"`
	// status exposes the health and status of the ongoing cluster update
	// +optional
	Status MachineConfigPoolProgressInsightStatus `json:"status"`
}

// MachineConfigPoolProgressInsightSpec is empty for now, MachineConfigPoolProgressInsightSpec is purely status-reporting API. In the future spec may be used
// to hold configuration to drive what information is surfaced and how
type MachineConfigPoolProgressInsightSpec struct {
}

// MachineConfigPoolProgressInsightStatus reports the state of a MachineConfigPool resource (which represents a pool of nodes
// update in standalone clusters), during a cluster update.
// +kubebuilder:validation:XValidation:rule="self.name == self.resource.name",message=".name must match .resource.name"
type MachineConfigPoolProgressInsightStatus struct {
	// conditions provide details about the machine config pool update. It contains at most 10 items. Known conditions are:
	// - Updating: whether the pool is updating; When Updating=False, the reason field can be Pending, Updated or Excluded
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// name is the name of the machine config pool
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Name string `json:"name"`

	// resource is the MachineConfigPool resource that represents the pool
	//
	// +Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// +resource name (because the rest is implied by status insight type). However, because we use resource references in
	// +many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// +than type safety for producers.
	// +required
	// +kubebuilder:validation:XValidation:rule="self.group == 'machineconfiguration.openshift.io' && self.resource == 'machineconfigpools'",message="resource must be a machineconfigpools.machineconfiguration.openshift.io resource"
	Resource ResourceRef `json:"resource"`

	// scopeType describes whether the pool is a control plane or a worker pool
	// +required
	Scope ScopeType `json:"scopeType"`

	// assessment is the assessment of the machine config pool update process. Valid values are: Pending, Completed, Degraded, Excluded, Progressing
	// +required
	Assessment PoolAssessment `json:"assessment"`

	// completion is a percentage of the update completion (0-100)
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Completion int32 `json:"completion"`

	// summaries is a list of counts of nodes matching certain criteria (e.g. updated, degraded, etc.). Maximum 16 items can be listed.
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Summaries []NodeSummary `json:"summaries,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// PoolAssessment is the assessment of the node pool update process
// +kubebuilder:validation:Enum=Pending;Completed;Degraded;Excluded;Progressing
type PoolAssessment string

const (
	// Pending means the nodes in the pool will be updated but none have even started yet
	PoolPending PoolAssessment = "Pending"
	// Completed means all nodes in the pool have been updated
	PoolCompleted PoolAssessment = "Completed"
	// Degraded means the process of updating the pool suffers from an observed problem
	PoolDegraded PoolAssessment = "Degraded"
	// Excluded means some (or all) nodes in the pool would be normally updated but a configuration (such as paused MCP)
	// prevents that from happening
	PoolExcluded PoolAssessment = "Excluded"
	// Progressing means the nodes in the pool are being updated and no problems or slowness are detected
	PoolProgressing PoolAssessment = "Progressing"
)

// NodeSummary is a count of nodes matching certain criteria (e.g. updated, degraded, etc.)
type NodeSummary struct {
	// type is the type of the summary. Valid values are: Total, Available, Progressing, Outdated, Draining, Excluded, Degraded
	// The summaries are not exclusive, a single node may be counted in multiple summaries.
	// +required
	Type NodeSummaryType `json:"type"`

	// count is the number of nodes matching the criteria, between 0 and 1024
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1024
	Count int32 `json:"count"`
}

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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPoolProgressInsightList is a list of MachineConfigPoolProgressInsightList resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type MachineConfigPoolProgressInsightList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a list of MachineConfigPoolProgressInsight resources
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Items []MachineConfigPoolProgressInsight `json:"items"`
}
