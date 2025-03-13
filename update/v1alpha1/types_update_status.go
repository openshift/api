package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpdateStatus reports status for in-progress cluster version updates
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=updatestatuses,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Provides health and status information about OpenShift cluster updates."
// +kubebuilder:metadata:annotations="displayName=UpdateStatuses"
type UpdateStatus struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, UpdateStatus is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +required
	Spec UpdateStatusSpec `json:"spec"`
	// status exposes the health and status of the ongoing cluster update
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
	// conditions provide details about the controller operational matters, exposing whether the controller managing this
	// UpdateStatus is functioning well, receives insights from individual informers, and is able to interpret them and
	// relay them through this UpdateStatus. These condition do not communicate anything about the state of the update
	// itself but may indicate whether the UpdateStatus content is reliable or not.
	// +TODO(UpdateStatus API GA): Update the list of conditions expected to be present
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// controlPlane contains a summary and insights related to the control plane update
	// +optional
	ControlPlane *ControlPlane `json:"controlPlane,omitempty"`

	// workerPools contains summaries and insights related to the worker pools update. Each item in the list represents
	// a single worker pool and carries all insights reported for it by informers. It has at most 32 items.
	// +TODO(UpdateStatus API GA): Determine a proper limit for MCPs / NodePool
	// +TODO(UpdateStatus): How to handle degenerate clusters with many pools? Worst case clusters can have per-node pools
	//                     so hundreds, and hypothetically more empty ones.
	// +listType=map
	// +listMapKey=name
	// +patchStrategy=merge
	// +patchMergeKey=name
	// +optional
	// +kubebuilder:validation:MaxItems=32
	WorkerPools []Pool `json:"workerPools,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// ControlPlane contains a summary and insights related to the control plane update.
type ControlPlane struct {
	// conditions provides details about the control plane update. This is a high-level status of an abstract control plane
	// concept, and will typically be the controller's interpretation / summarization of the insights it received (that
	// will be placed in .informers[].insights for clients that want to perform further analysis of the data).
	// Known condition types are:
	// * "Updating": Whether the cluster control plane is currently updating or not
	//
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// resource is the resource that represents the control plane. It will typically be a ClusterVersion resource
	// in standalone OpenShift and HostedCluster in Hosted Control Planes. This field is optional because the information
	// may be unknown temporarily.
	//
	// +Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// +only the "correct" resource types to be referenced (here, ClusterVersion and HostedCluster). However, because we
	// +use resource references in many places and this API is intended to be consumed by clients, not produced, consistency
	// +seems to be more valuable than type safety for producers.
	// +optional
	// +kubebuilder:validation:XValidation:rule="(self.group == 'config.openshift.io' && self.resource == 'clusterversions') || (self.group == 'hypershift.openshift.io' && self.resource == 'hostedclusters')",message="controlPlane.resource must be either a clusterversions.config.openshift.io or a hostedclusters.hypershift.openshift.io resource"
	Resource *ResourceRef `json:"resource"`

	// poolResource is the resource that represents control plane node pool, typically a MachineConfigPool. This field
	// is optional because some form factors (like Hosted Control Planes) do not have dedicated control plane node pools,
	// and also because the information may be unknown temporarily.
	//
	// +Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// +only the "correct" resource types to be referenced (here, MachineConfigPool). However, because we use resource
	// +references in many places and this API is intended to be consumed by clients, not produced, consistency seems to be
	// +more valuable than type safety for producers.
	// +optional
	// +kubebuilder:validation:XValidation:rule="(self.group == 'machineconfiguration.openshift.io' && self.resource == 'machineconfigpools')",message="controlPlane.poolResource must be a machineconfigpools.machineconfiguration.openshift.io resource"
	PoolResource *PoolResourceRef `json:"poolResource,omitempty"`

	// informers is a list of insight producers. An informer is a system, internal or external to the cluster, that
	// produces units of information relevant to the update process, either about its progress or its health. Each
	// informer in the list is identified by a name, and contains a list of insights it contributed to the Update Status API,
	// relevant to the control plane update. Contains at most 16 items.
	//
	// +listType=map
	// +listMapKey=name
	// +patchStrategy=merge
	// +patchMergeKey=name
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Informers []ControlPlaneInformer `json:"informers,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
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
	ControlPlaneClusterVersionProgressing ControlPlaneUpdatingReason = "ClusterVersionProgressing"
	// ClusterVersionNotProgressing is used for Updating=False set because we observed a ClusterVersion resource to
	// have Progressing=False condition
	ControlPlaneClusterVersionNotProgressing ControlPlaneUpdatingReason = "ClusterVersionNotProgressing"
	// CannotDetermineUpdating is used with Updating=Unknown. This covers many different actual reasons such as
	// missing or Unknown Progressing condition on ClusterVersion, but it does not seem useful to track the individual
	// reasons to that granularity for Updating=Unknown
	ControlPlaneCannotDetermineUpdating ControlPlaneUpdatingReason = "CannotDetermineUpdating"
)

// Pool contains a summary and insights related to a node pool update
// +kubebuilder:validation:XValidation:rule="self.name == self.resource.name",message="workerPools .name must match .resource.name"
type Pool struct {
	// conditions provides details about the node pool update. This is a high-level status of an abstract "pool of nodes"
	// concept, and will typically be the controller's interpretation / summarization of the insights it received (that
	// will be placed in .informers[].insights for clients that want to perform further analysis of the data).
	// Known condition types are:
	// * "Updating": Whether the pool of nodes is currently updating or not
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// name is the name of the pool, follows the same rules as a Kubernetes resource name (RFC-1123 subdomain)
	// +required
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	Name string `json:"name"`

	// resource is the resource that represents the pool, either a MachineConfigPool in Standalone OpenShift or a NodePool
	// in Hosted Control Planes.
	//
	// +Note: By OpenShift API conventions, in isolation this should probably be a specialized reference type that allows
	// +only the "correct" resource types to be referenced (here, MachineConfigPool or NodePool). However, because we use
	// +resource references in many places and this API is intended to be consumed by clients, not produced, consistency
	// +seems to be more valuable than type safety for producers.
	// +required
	// +kubebuilder:validation:XValidation:rule="(self.group == 'machineconfiguration.openshift.io' && self.resource == 'machineconfigpools') || (self.group == 'hypershift.openshift.io' && self.resource == 'nodepools')",message="workerPools[].poolResource must be a machineconfigpools.machineconfiguration.openshift.io or hostedclusters.hypershift.openshift.io resource"
	Resource PoolResourceRef `json:"resource"`

	// informers is a list of insight producers. An informer is a system, internal or external to the cluster, that
	// produces units of information relevant to the update process, either about its progress or its health. Each
	// informer in the list is identified by a name, and contains a list of insights it contributed to the Update Status API,
	// relevant to the process of updating this pool of nodes.
	// +listType=map
	// +listMapKey=name
	// +patchStrategy=merge
	// +patchMergeKey=name
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Informers []WorkerPoolInformer `json:"informers,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// ControlPlaneInformer represents a system, internal or external to the cluster, that  produces units of information
// relevant to the update process, either about its progress or its health. Each informer is identified by a name, and
// contains a list of insights it contributed to the Update Status API, relevant to the control plane update.
type ControlPlaneInformer struct {
	// name is the name of the insight producer
	// +required
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	Name string `json:"name"`

	// insights is a list of insights produced by this producer. Insights are units of information relevant to an update
	// progress or health information. There are two types of update insights: status insights and health insights. The
	// first type are directly tied to the update process, regardless of whether it is proceeding smoothly or not.
	//
	// Status Insights expose the state of a single resource that is directly involved in the update process, usually a resource
	// that either has a notion of "being updated," (such as a Node or ClusterOperator) or represents a higher-level
	// abstraction (such as a ClusterVersion resource tahat represents the control plane or MachineConfigPool that represents
	// a pool of nodes).
	//
	// Health Insights report a state or condition in the cluster that is abnormal or negative and either affects or is
	// affected by the update. Ideally, none would be generated in a standard healthy update. Health insights communicate
	// a condition that warrants attention by the cluster administrator.
	// +listType=map
	// +listMapKey=uid
	// +patchStrategy=merge
	// +patchMergeKey=uid
	// +optional
	// +kubebuilder:validation:MaxItems=128
	Insights []ControlPlaneInsight `json:"insights,omitempty" patchStrategy:"merge" patchMergeKey:"uid"`
}

// WorkerPoolInformer represents a system, internal or external to the cluster, that  produces units of information
// relevant to the update process, either about its progress or its health. Each informer is identified by a name, and
// contains a list of insights it contributed to the Update Status API, relevant to a specific worker pool.
type WorkerPoolInformer struct {
	// name is the name of the insight producer
	// +required
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	Name string `json:"name"`

	// insights is a list of insights produced by this producer. Insights are units of information relevant to an update
	// progress or health information. There are two types of update insights: status insights and health insights. The
	// first type are directly tied to the update process, regardless of whether it is proceeding smoothly or not.
	//
	// Status Insights expose the state of a single resource that is directly involved in the update process, usually a resource
	// that either has a notion of "being updated," (such as a Node or ClusterOperator) or represents a higher-level
	// abstraction (such as a ClusterVersion resource tahat represents the control plane or MachineConfigPool that represents
	// a pool of nodes).
	//
	// Health Insights report a state or condition in the cluster that is abnormal or negative and either affects or is
	// affected by the update. Ideally, none would be generated in a standard healthy update. Health insights communicate
	// a condition that warrants attention by the cluster administrator.
	// +listType=map
	// +listMapKey=uid
	// +patchStrategy=merge
	// +patchMergeKey=uid
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Insights []WorkerPoolInsight `json:"insights,omitempty" patchStrategy:"merge" patchMergeKey:"uid"`
}

// ControlPlaneInsight is a unique piece of either status/progress or update health information produced by update informer
type ControlPlaneInsight struct {
	// uid identifies the insight over time
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	UID string `json:"uid"`

	// acquiredAt is the time when the data was acquired by the producer
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	AcquiredAt metav1.Time `json:"acquiredAt"`

	// insight is a discriminated union of all insights types that can be reported for the control plane
	// +required
	Insight ControlPlaneInsightUnion `json:"insight"`
}

// WorkerPoolInsight is a unique piece of either status/progress or update health information produced by update informer
type WorkerPoolInsight struct {
	// uid identifies the insight over time
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	UID string `json:"uid"`

	// acquiredAt is the time when the data was acquired by the producer
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	AcquiredAt metav1.Time `json:"acquiredAt"`

	// insight is a discriminated union of all insights types that can be reported for a worker pool
	// +required
	Insight WorkerPoolInsightUnion `json:"insight"`
}

// ControlPlaneInsightUnion is the discriminated union of all insights types that can be reported for the control plane,
// identified by type field
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'ClusterVersion' ?  has(self.clusterVersion) : !has(self.clusterVersion)",message="clusterVersion is required when type is ClusterVersion, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'ClusterOperator' ?  has(self.clusterOperator) : !has(self.clusterOperator)",message="clusterOperator is required when type is ClusterOperator, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'MachineConfigPool' ?  has(self.machineConfigPool) : !has(self.machineConfigPool)",message="machineConfigPool is required when type is MachineConfigPool, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Node' ?  has(self.node) : !has(self.node)",message="node is required when type is Node, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Health' ?  has(self.health) : !has(self.health)",message="health is required when type is Health, and forbidden otherwise"
// +union
type ControlPlaneInsightUnion struct {
	// type identifies the type of the update insight, one of: ClusterVersion, ClusterOperator, MachineConfigPool, Node, Health
	// ClusterVersion, ClusterOperator, MachineConfigPool, Node types are progress insights about a resource directly
	// involved in the update process
	// Health insights report a state or condition in the cluster that is abnormal or negative and either affects or is
	// affected by the update.
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum=ClusterVersion;ClusterOperator;MachineConfigPool;Node;Health
	Type InsightType `json:"type"`

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
	HealthInsight *HealthInsight `json:"health,omitempty"`
}

// WorkerPoolInsightUnion is the discriminated union of insights types that can be reported for a worker pool,
// identified by type field
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'MachineConfigPool' ?  has(self.machineConfigPool) : !has(self.machineConfigPool)",message="machineConfigPool is required when type is MachineConfigPool, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Node' ?  has(self.node) : !has(self.node)",message="node is required when type is Node, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Health' ?  has(self.health) : !has(self.health)",message="health is required when type is Health, and forbidden otherwise"
// +union
type WorkerPoolInsightUnion struct {
	// type identifies the type of the update insight, one of: MachineConfigPool, Node, Health
	// MachineConfigPool, Node types are progress insights about a resource directly involved in the update process
	// Health insights report a state or condition in the cluster that is abnormal or negative and either affects or is
	// affected by the update.
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum=MachineConfigPool;Node;Health
	Type InsightType `json:"type"`

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
	HealthInsight *HealthInsight `json:"health,omitempty"`
}

// InsightType identifies the type of the update insight as either one of the resource-specific status insight,
// or a generic health insight
type InsightType string

const (
	// Resource-specific status insights should be reported continuously during the update process and mostly communicate
	// progress and high-level state

	// ClusterVersion status insight reports progress and high-level state of a ClusterVersion resource, representing
	// control plane in standalone clusters
	ClusterVersionStatusInsightType InsightType = "ClusterVersion"
	// ClusterOperator status insight reports progress and high-level state of a ClusterOperator, representing a control
	// plane component
	ClusterOperatorStatusInsightType InsightType = "ClusterOperator"
	// MachineConfigPool status insight reports progress and high-level state of a MachineConfigPool resource, representing
	// a pool of nodes in clusters using Machine API
	MachineConfigPoolStatusInsightType InsightType = "MachineConfigPool"
	// Node status insight reports progress and high-level state of a Node resource, representing a node (both control
	// plane and worker) in a cluster
	NodeStatusInsightType InsightType = "Node"

	// Health insights are reported only when an informer observes a condition that requires admin attention
	HealthInsightType InsightType = "Health"
)

// ResourceRef is a reference to a kubernetes resource, typically involved in an insight
type ResourceRef struct {
	// group of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Group string `json:"group,omitempty"`

	// resource of object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Resource string `json:"resource"`

	// name of the object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Name string `json:"name"`

	// namespace of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Namespace string `json:"namespace,omitempty"`
}

// PoolResourceRef is a reference to a kubernetes resource that represents a node pool
// +kubebuilder:validation:XValidation:rule="has(self.resource) && (self.resource == 'machineconfigpools' && self.group == 'machineconfiguration.openshift.io'",message="a poolResource must be a machineconfigpools.machineconfiguration.openshift.io resource"
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
	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a  list of UpdateStatus resources
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Items []UpdateStatus `json:"items"`
}
