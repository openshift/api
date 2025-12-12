package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PacemakerCluster is used in Two Node OpenShift with Fencing (TNF) deployments to monitor the health
// of etcd running under pacemaker.

// PacemakerCluster condition types (cluster-level)
const (
	// PacemakerClusterHealthyConditionType tracks the overall health of the pacemaker cluster.
	// This is an aggregate condition that reflects the health of all cluster-level conditions and node health.
	// Specifically, it aggregates the following conditions:
	// - PacemakerClusterNotInMaintenanceModeConditionType
	// - PacemakerClusterExpectedNodeCountConditionType
	// - PacemakerClusterNodeHealthyConditionType (for each node)
	// When True, the cluster is healthy with reason "ClusterHealthy".
	// When False, the cluster is unhealthy with reason "ClusterUnhealthy".
	PacemakerClusterHealthyConditionType = "PacemakerClusterHealthy"

	// PacemakerClusterNotInMaintenanceModeConditionType tracks whether the cluster is not in maintenance mode.
	// Maintenance mode is a cluster-wide setting that prevents pacemaker from starting or stopping resources.
	// When True, the cluster is not in maintenance mode with reason "NotInMaintenance". This is the normal operating state.
	// When False, the cluster is in maintenance mode with reason "InMaintenance". This is an unexpected state.
	PacemakerClusterNotInMaintenanceModeConditionType = "PacemakerClusterNotInMaintenanceMode"

	// PacemakerClusterExpectedNodeCountConditionType tracks whether the cluster has the expected number of nodes.
	// For TNF, we are expecting exactly 2 nodes.
	// When True, the expected number of nodes are present with reason "ExpectedNodeCount".
	// When False, the node count is incorrect with reason "InsufficientNodes" or "ExcessiveNodes".
	PacemakerClusterExpectedNodeCountConditionType = "PacemakerClusterExpectedNodeCount"
)

// PacemakerClusterHealthy condition reasons
const (
	// PacemakerClusterHealthyReasonClusterHealthy means the pacemaker cluster is healthy and operating normally.
	PacemakerClusterHealthyReasonClusterHealthy = "ClusterHealthy"

	// PacemakerClusterHealthyReasonClusterUnhealthy means the pacemaker cluster has issues that need investigation.
	PacemakerClusterHealthyReasonClusterUnhealthy = "ClusterUnhealthy"
)

// PacemakerClusterNotInMaintenanceMode condition reasons
const (
	// PacemakerClusterNotInMaintenanceModeReasonNotInMaintenance means the cluster is not in maintenance mode.
	// This is the normal operating state.
	PacemakerClusterNotInMaintenanceModeReasonNotInMaintenance = "NotInMaintenance"

	// PacemakerClusterNotInMaintenanceModeReasonInMaintenance means the cluster is in maintenance mode.
	// In maintenance mode, pacemaker will not start or stop any resources. Entering and exiting this state requires
	// manual user intervention, and is unexpected during normal cluster operation.
	PacemakerClusterNotInMaintenanceModeReasonInMaintenance = "InMaintenance"
)

// PacemakerClusterExpectedNodeCount condition reasons
const (
	// PacemakerClusterExpectedNodeCountReasonExpectedNodeCount means the expected number of nodes are present.
	// For TNF, we are expecting exactly 2 nodes. This is the expected healthy state.
	PacemakerClusterExpectedNodeCountReasonExpectedNodeCount = "ExpectedNodeCount"

	// PacemakerClusterExpectedNodeCountReasonInsufficientNodes means fewer nodes than expected are present.
	// For TNF, this means that less than 2 nodes are present. Under normal operation, this will only happen during
	// a node replacement operation. It's also possible to enter this state with manual user intervention, but
	// will also require user intervention to restore normal functionality.
	PacemakerClusterExpectedNodeCountReasonInsufficientNodes = "InsufficientNodes"

	// PacemakerClusterExpectedNodeCountReasonExcessiveNodes means more nodes than expected are present.
	// For TNF, this means more than 2 nodes are present. This should be investigated as it is unexpected and should
	// never happen during normal cluster operation. It is possible to enter this state with manual user intervention,
	// but will also require user intervention to restore normal functionality.
	PacemakerClusterExpectedNodeCountReasonExcessiveNodes = "ExcessiveNodes"
)

// PacemakerClusterNode condition types (node-level)
const (
	// PacemakerClusterNodeHealthyConditionType tracks the overall health of a node in the pacemaker cluster.
	// This is an aggregate condition that reflects the health of all node-level conditions and resource health.
	// Specifically, it aggregates the following conditions:
	// - PacemakerClusterNodeOnlineConditionType
	// - PacemakerClusterNodeNotInMaintenanceConditionType
	// - PacemakerClusterNodeActiveConditionType
	// - PacemakerClusterNodeReadyConditionType
	// - PacemakerClusterNodeCleanConditionType
	// - PacemakerClusterNodeMemberConditionType
	// - PacemakerClusterResourceHealthyConditionType (for the node's kubelet resource)
	// - PacemakerClusterResourceHealthyConditionType (for the node's etcd resource)
	// - PacemakerClusterResourceHealthyConditionType (for the node's fencing agent resource)
	// When True, the node is healthy with reason "NodeHealthy".
	// When False, the node is unhealthy with reason "NodeUnhealthy".
	PacemakerClusterNodeHealthyConditionType = "PacemakerClusterNodeHealthy"

	// PacemakerClusterNodeOnlineConditionType tracks whether a node is online.
	// When True, the node is online with reason "Online". This is the normal operating state.
	// When False, the node is offline with reason "Offline". This is an expected state.
	PacemakerClusterNodeOnlineConditionType = "PacemakerClusterNodeOnline"

	// PacemakerClusterNodeNotInMaintenanceConditionType tracks whether a node is not in maintenance mode.
	// A node in maintenance mode is ignored by pacemaker while maintenance mode is active.
	// When True, the node is not in maintenance mode with reason "NotInMaintenance". This is the normal operating state.
	// When False, the node is in maintenance mode with reason "InMaintenance". This is an unexpected state.
	PacemakerClusterNodeNotInMaintenanceConditionType = "PacemakerClusterNodeNotInMaintenance"

	// PacemakerClusterNodeActiveConditionType tracks whether a node is active (not in standby mode).
	// When a node enters standby mode, pacemaker moves its resources to other nodes in the cluster.
	// In TNF, we do not use standby mode during normal operation.
	// When True, the node is active with reason "NodeActive". This is the normal operating state.
	// When False, the node is in standby mode with reason "NodeStandby". This is an unexpected state.
	PacemakerClusterNodeActiveConditionType = "PacemakerClusterNodeActive"

	// PacemakerClusterNodeReadyConditionType tracks whether a node is ready (not in a pending state).
	// A node in a pending state is in the process of joining or leaving the cluster.
	// When True, the node is ready with reason "NodeReady". This is the normal operating state.
	// When False, the node is pending with reason "NodePending". This is expected to be temporary.
	PacemakerClusterNodeReadyConditionType = "PacemakerClusterNodeReady"

	// PacemakerClusterNodeCleanConditionType tracks whether a node is in a clean state.
	// An unclean state means that pacemaker was unable to confirm the node's state, which signifies issues
	// in fencing, communication, or configuration.
	// When True, the node is clean with reason "NodeClean". This is the normal operating state.
	// When False, the node is unclean with reason "NodeUnclean". This is an unexpected state.
	PacemakerClusterNodeCleanConditionType = "PacemakerClusterNodeClean"

	// PacemakerClusterNodeMemberConditionType tracks whether a node is a member of the cluster.
	// Some configurations may use remote nodes or ping nodes, which are nodes that are not members.
	// For TNF, we expect both nodes to be members.
	// When True, the node is a member with reason "Member". This is the normal operating state.
	// When False, the node is not a member with reason "NotMember". This is an unexpected state.
	PacemakerClusterNodeMemberConditionType = "PacemakerClusterNodeMember"
)

// PacemakerClusterNodeHealthy condition reasons
const (
	// PacemakerClusterNodeHealthyReasonNodeHealthy means the node is healthy and operating normally.
	PacemakerClusterNodeHealthyReasonNodeHealthy = "NodeHealthy"

	// PacemakerClusterNodeHealthyReasonNodeUnhealthy means the node has issues that need investigation.
	PacemakerClusterNodeHealthyReasonNodeUnhealthy = "NodeUnhealthy"
)

// PacemakerClusterNodeOnline condition reasons
const (
	// PacemakerClusterNodeOnlineReasonOnline means the node is online. This is the normal operating state.
	PacemakerClusterNodeOnlineReasonOnline = "Online"

	// PacemakerClusterNodeOnlineReasonOffline means the node is offline.
	PacemakerClusterNodeOnlineReasonOffline = "Offline"
)

// PacemakerClusterNodeNotInMaintenance condition reasons
const (
	// PacemakerClusterNodeNotInMaintenanceReasonNotInMaintenance means the node is not in maintenance mode.
	// This is the normal operating state.
	PacemakerClusterNodeNotInMaintenanceReasonNotInMaintenance = "NotInMaintenance"

	// PacemakerClusterNodeNotInMaintenanceReasonInMaintenance means the node is in maintenance mode.
	// This is an unexpected state.
	PacemakerClusterNodeNotInMaintenanceReasonInMaintenance = "InMaintenance"
)

// PacemakerClusterNodeActive condition reasons
const (
	// PacemakerClusterNodeActiveReasonNodeActive means the node is active (not in standby mode).
	// This is the normal operating state.
	PacemakerClusterNodeActiveReasonNodeActive = "NodeActive"

	// PacemakerClusterNodeActiveReasonNodeStandby means the node is in standby mode.
	// This is an unexpected state.
	PacemakerClusterNodeActiveReasonNodeStandby = "NodeStandby"
)

// PacemakerClusterNodeReady condition reasons
const (
	// PacemakerClusterNodeReadyReasonNodeReady means the node is ready (not in a pending state).
	// This is the normal operating state.
	PacemakerClusterNodeReadyReasonNodeReady = "NodeReady"

	// PacemakerClusterNodeReadyReasonNodePending means the node is joining or leaving the cluster.
	// This state is expected to be temporary.
	PacemakerClusterNodeReadyReasonNodePending = "NodePending"
)

// PacemakerClusterNodeClean condition reasons
const (
	// PacemakerClusterNodeCleanReasonNodeClean means the node is in a clean state.
	// This is the normal operating state.
	PacemakerClusterNodeCleanReasonNodeClean = "NodeClean"

	// PacemakerClusterNodeCleanReasonNodeUnclean means the node is in an unclean state.
	// Pacemaker was unable to confirm the node's state, which signifies issues in fencing, communication, or configuration.
	// This is an unexpected state.
	PacemakerClusterNodeCleanReasonNodeUnclean = "NodeUnclean"
)

// PacemakerClusterNodeMember condition reasons
const (
	// PacemakerClusterNodeMemberReasonMember means the node is a member of the cluster.
	// For TNF, we expect both nodes to be members. This is the normal operating state.
	PacemakerClusterNodeMemberReasonMember = "Member"

	// PacemakerClusterNodeMemberReasonNotMember means the node is not a member of the cluster.
	// Some configurations may use remote nodes or ping nodes, which are nodes that are not members.
	// This is an unexpected state.
	PacemakerClusterNodeMemberReasonNotMember = "NotMember"
)

// PacemakerClusterResource condition types (resource-level)
const (
	// PacemakerClusterResourceHealthyConditionType tracks the overall health of a pacemaker resource.
	// This is an aggregate condition that reflects the health of all resource-level conditions.
	// Specifically, it aggregates the following conditions:
	// - PacemakerClusterResourceNotInMaintenanceConditionType
	// - PacemakerClusterResourceManagedConditionType
	// - PacemakerClusterResourceEnabledConditionType
	// - PacemakerClusterResourceOperationalConditionType
	// - PacemakerClusterResourceActiveConditionType
	// - PacemakerClusterResourceStartedConditionType
	// - PacemakerClusterResourceUnblockedConditionType
	// When True, the resource is healthy with reason "ResourceHealthy".
	// When False, the resource is unhealthy with reason "ResourceUnhealthy".
	PacemakerClusterResourceHealthyConditionType = "PacemakerClusterResourceHealthy"

	// PacemakerClusterResourceNotInMaintenanceConditionType tracks whether a resource is not in maintenance mode.
	// Resources in maintenance mode are not monitored or moved by pacemaker.
	// In TNF, we do not expect any resources to be in maintenance mode.
	// When True, the resource is not in maintenance mode with reason "NotInMaintenance". This is the normal operating state.
	// When False, the resource is in maintenance mode with reason "InMaintenance". This is an unexpected state.
	PacemakerClusterResourceNotInMaintenanceConditionType = "PacemakerClusterResourceNotInMaintenance"

	// PacemakerClusterResourceManagedConditionType tracks whether a resource is managed by pacemaker.
	// Resources that are not managed by pacemaker are effectively invisible to the pacemaker HA logic.
	// For TNF, all resources are expected to be managed.
	// When True, the resource is managed with reason "Managed". This is the normal operating state.
	// When False, the resource is not managed with reason "NotManaged". This is an unexpected state.
	PacemakerClusterResourceManagedConditionType = "PacemakerClusterResourceManaged"

	// PacemakerClusterResourceEnabledConditionType tracks whether a resource is enabled.
	// Resources that are disabled are stopped and not automatically managed or started by the cluster.
	// In TNF, we do not expect any resources to be disabled.
	// When True, the resource is enabled with reason "Enabled". This is the normal operating state.
	// When False, the resource is disabled with reason "Disabled". This is an unexpected state.
	PacemakerClusterResourceEnabledConditionType = "PacemakerClusterResourceEnabled"

	// PacemakerClusterResourceOperationalConditionType tracks whether a resource is operational (not failed).
	// A failed resource is one that is not able to start or is in an error state.
	// When True, the resource is operational with reason "Operational". This is the normal operating state.
	// When False, the resource has failed with reason "Failed". This is an unexpected state.
	PacemakerClusterResourceOperationalConditionType = "PacemakerClusterResourceOperational"

	// PacemakerClusterResourceActiveConditionType tracks whether a resource is active.
	// An active resource is running on a cluster node.
	// In TNF, all resources are expected to be active.
	// When True, the resource is active with reason "Active". This is the normal operating state.
	// When False, the resource is not active with reason "NotActive". This is an unexpected state.
	PacemakerClusterResourceActiveConditionType = "PacemakerClusterResourceActive"

	// PacemakerClusterResourceStartedConditionType tracks whether a resource is started.
	// It's normal for a resource like etcd to become stopped in the event of a quorum loss event because
	// the pacemaker recovery logic will fence a node and restore etcd quorum on the surviving node as a cluster-of-one.
	// A resource that stays stopped for an extended period of time is an unexpected state and should be investigated.
	// When True, the resource is started with reason "Started". This is the normal operating state.
	// When False, the resource is not started with reason "Stopped". This is expected to be temporary.
	PacemakerClusterResourceStartedConditionType = "PacemakerClusterResourceStarted"

	// PacemakerClusterResourceUnblockedConditionType tracks whether a resource is unblocked.
	// A resource that is blocked is unable to start or move to a different node.
	// In TNF, we do not expect any resources to be blocked.
	// When True, the resource is unblocked with reason "Unblocked". This is the normal operating state.
	// When False, the resource is blocked with reason "Blocked". This is an unexpected state.
	PacemakerClusterResourceUnblockedConditionType = "PacemakerClusterResourceUnblocked"
)

// PacemakerClusterResourceHealthy condition reasons
const (
	// PacemakerClusterResourceHealthyReasonResourceHealthy means the resource is healthy and operating normally.
	PacemakerClusterResourceHealthyReasonResourceHealthy = "ResourceHealthy"

	// PacemakerClusterResourceHealthyReasonResourceUnhealthy means the resource has issues that need investigation.
	PacemakerClusterResourceHealthyReasonResourceUnhealthy = "ResourceUnhealthy"
)

// PacemakerClusterResourceNotInMaintenance condition reasons
const (
	// PacemakerClusterResourceNotInMaintenanceReasonNotInMaintenance means the resource is not in maintenance mode.
	// This is the normal operating state.
	PacemakerClusterResourceNotInMaintenanceReasonNotInMaintenance = "NotInMaintenance"

	// PacemakerClusterResourceNotInMaintenanceReasonInMaintenance means the resource is in maintenance mode.
	// Resources in maintenance mode are not monitored or moved by pacemaker. This is an unexpected state.
	PacemakerClusterResourceNotInMaintenanceReasonInMaintenance = "InMaintenance"
)

// PacemakerClusterResourceManaged condition reasons
const (
	// PacemakerClusterResourceManagedReasonManaged means the resource is managed by pacemaker.
	// This is the normal operating state.
	PacemakerClusterResourceManagedReasonManaged = "Managed"

	// PacemakerClusterResourceManagedReasonNotManaged means the resource is not managed by pacemaker.
	// Resources that are not managed by pacemaker are effectively invisible to the pacemaker HA logic.
	// This is an unexpected state.
	PacemakerClusterResourceManagedReasonNotManaged = "NotManaged"
)

// PacemakerClusterResourceEnabled condition reasons
const (
	// PacemakerClusterResourceEnabledReasonEnabled means the resource is enabled.
	// This is the normal operating state.
	PacemakerClusterResourceEnabledReasonEnabled = "Enabled"

	// PacemakerClusterResourceEnabledReasonDisabled means the resource is disabled.
	// Resources that are disabled are stopped and not automatically managed or started by the cluster.
	// This is an unexpected state.
	PacemakerClusterResourceEnabledReasonDisabled = "Disabled"
)

// PacemakerClusterResourceOperational condition reasons
const (
	// PacemakerClusterResourceOperationalReasonOperational means the resource is operational (not failed).
	// This is the normal operating state.
	PacemakerClusterResourceOperationalReasonOperational = "Operational"

	// PacemakerClusterResourceOperationalReasonFailed means the resource has failed.
	// A failed resource is one that is not able to start or is in an error state. This is an unexpected state.
	PacemakerClusterResourceOperationalReasonFailed = "Failed"
)

// PacemakerClusterResourceActive condition reasons
const (
	// PacemakerClusterResourceActiveReasonActive means the resource is active.
	// An active resource is running on a cluster node. This is the normal operating state.
	PacemakerClusterResourceActiveReasonActive = "Active"

	// PacemakerClusterResourceActiveReasonNotActive means the resource is not active.
	// This is an unexpected state.
	PacemakerClusterResourceActiveReasonNotActive = "NotActive"
)

// PacemakerClusterResourceStarted condition reasons
const (
	// PacemakerClusterResourceStartedReasonStarted means the resource is started.
	// This is the normal operating state.
	PacemakerClusterResourceStartedReasonStarted = "Started"

	// PacemakerClusterResourceStartedReasonStopped means the resource is stopped.
	// It's normal for a resource like etcd to become stopped in the event of a quorum loss event because
	// the pacemaker recovery logic will fence a node and restore etcd quorum on the surviving node as a cluster-of-one.
	// A resource that stays stopped for an extended period of time is an unexpected state and should be investigated.
	PacemakerClusterResourceStartedReasonStopped = "Stopped"
)

// PacemakerClusterResourceUnblocked condition reasons
const (
	// PacemakerClusterResourceUnblockedReasonUnblocked means the resource is unblocked.
	// This is the normal operating state.
	PacemakerClusterResourceUnblockedReasonUnblocked = "Unblocked"

	// PacemakerClusterResourceUnblockedReasonBlocked means the resource is blocked.
	// A resource that is blocked is unable to start or move to a different node. This is an unexpected state.
	PacemakerClusterResourceUnblockedReasonBlocked = "Blocked"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// PacemakerCluster represents the current state of the pacemaker cluster as reported by the pcs status command.
// PacemakerCluster is a cluster-scoped singleton resource. The name of this instance is "cluster". This
// resource provides a view into the health and status of a pacemaker-managed cluster in TNF deployments.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pacemakerclusters,scope=Cluster,singular=pacemakercluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2544
// +openshift:file-pattern=cvoRunLevel=0000_25,operatorName=etcd,operatorOrdering=01,operatorComponent=two-node-fencing
// +openshift:enable:FeatureGate=DualReplica
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="PacemakerCluster must be named 'cluster'"
type PacemakerCluster struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// status contains the actual pacemaker cluster status information collected from the cluster.
	// The goal of this status is to be able to quickly identify if pacemaker is in a healthy state.
	// In TNF, a healthy pacemaker cluster has 2 nodes, both of which have healthy kubelet, etcd, and fencing resources.
	// +optional
	Status *PacemakerClusterStatus `json:"status,omitempty"`
}

// PacemakerClusterStatus contains the actual pacemaker cluster status information. As part of validating the status
// object, we need to ensure that the lastUpdated timestamp is always newer than the current value. We allow the
// lastUpdated timestamp to be empty on initial creation, but it is required once it has been set.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.lastUpdated) || (has(self.lastUpdated) && self.lastUpdated > oldSelf.lastUpdated)",message="lastUpdated must be a newer timestamp than the current value"
type PacemakerClusterStatus struct {
	// conditions represent the observations of the pacemaker cluster's current state.
	// Known condition types are: "PacemakerClusterHealthy", "PacemakerClusterNotInMaintenanceMode",
	// "PacemakerClusterExpectedNodeCount".
	// The "PacemakerClusterHealthy" condition is an aggregate that tracks the overall health of the cluster.
	// The "PacemakerClusterNotInMaintenanceMode" condition tracks whether the cluster is not in maintenance mode.
	// The "PacemakerClusterExpectedNodeCount" condition tracks whether the expected number of nodes are present.
	// Each of these conditions is required.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterHealthy')",message="conditions must contain a condition of type PacemakerClusterHealthy"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNotInMaintenanceMode')",message="conditions must contain a condition of type PacemakerClusterNotInMaintenanceMode"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterExpectedNodeCount')",message="conditions must contain a condition of type PacemakerClusterExpectedNodeCount"
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// lastUpdated is the timestamp when this status was last updated. This is useful for identifying
	// stale status reports.
	// When present, it must be a valid timestamp in RFC3339 format and cannot be set to an earlier timestamp than the
	// current value. It is optional upon initial creation, but required once it has been set.
	// +kubebuilder:validation:Format=date-time
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`

	// nodes provides detailed information about each node in the cluster including per-node resource health.
	// Each node entry includes the node's name, IP address, conditions, and resource status.
	// The list can contain between 1 and 32 nodes. The upper limit is imposed by pacemaker.
	// For TNF, exactly 2 nodes are expected in a healthy cluster.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +optional
	Nodes []PacemakerClusterNodeStatus `json:"nodes,omitempty"`
}

// PacemakerClusterNodeStatus represents the status of a single node in the pacemaker cluster including
// the node's conditions and the health of critical resources running on that node.
type PacemakerClusterNodeStatus struct {
	// conditions represent the observations of the node's current state.
	// Known condition types are: "PacemakerClusterNodeHealthy", "PacemakerClusterNodeOnline",
	// "PacemakerClusterNodeNotInMaintenance", "PacemakerClusterNodeActive", "PacemakerClusterNodeReady",
	// "PacemakerClusterNodeClean", "PacemakerClusterNodeMember".
	// The "PacemakerClusterNodeHealthy" condition is an aggregate that tracks the overall health of the node.
	// The "PacemakerClusterNodeOnline" condition tracks whether the node is online.
	// The "PacemakerClusterNodeNotInMaintenance" condition tracks whether the node is not in maintenance mode.
	// The "PacemakerClusterNodeActive" condition tracks whether the node is active (not in standby mode).
	// The "PacemakerClusterNodeReady" condition tracks whether the node is ready (not in a pending state).
	// The "PacemakerClusterNodeClean" condition tracks whether the node is in a clean (status known) state.
	// The "PacemakerClusterNodeMember" condition tracks whether the node is a member of the cluster.
	// Each of these conditions is required.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeHealthy')",message="conditions must contain a condition of type PacemakerClusterNodeHealthy"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeOnline')",message="conditions must contain a condition of type PacemakerClusterNodeOnline"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeNotInMaintenance')",message="conditions must contain a condition of type PacemakerClusterNodeNotInMaintenance"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeActive')",message="conditions must contain a condition of type PacemakerClusterNodeActive"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeReady')",message="conditions must contain a condition of type PacemakerClusterNodeReady"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeClean')",message="conditions must contain a condition of type PacemakerClusterNodeClean"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterNodeMember')",message="conditions must contain a condition of type PacemakerClusterNodeMember"
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// name is the name of the node. This is expected to match the Kubernetes node's name, which must be a lowercase
	// RFC 1123 subdomain consisting of lowercase alphanumeric characters, '-' or '.', starting and ending with
	// an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="name must be a lowercase RFC 1123 subdomain consisting of lowercase alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character"
	// +required
	Name string `json:"name,omitempty"`

	// ipAddress is the canonical IPv4 or IPv6 address of the node. It must be a valid canonical global unicast IPv4
	// or IPv6 address (including private/RFC1918 addresses). This excludes special addresses like unspecified,
	// loopback, link-local, and multicast.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=39
	// +kubebuilder:validation:XValidation:rule="isIP(self) && ip.isCanonical(self) && ip(self).isGlobalUnicast()",message="ipAddress must be a valid canonical global unicast IPv4 or IPv6 address"
	// +required
	IPAddress string `json:"ipAddress,omitempty"`

	// kubelet contains the status of the kubelet pacemaker resource on this node.
	// The kubelet resource is a prerequisite for etcd in TNF deployments. It is expected to be running on all nodes.
	// +optional
	Kubelet *PacemakerClusterResourceStatus `json:"kubelet,omitempty"`

	// etcd contains the status of the etcd pacemaker resource on this node.
	// The etcd resource may temporarily transition to stopped during pacemaker quorum-recovery operations.
	// +optional
	Etcd *PacemakerClusterResourceStatus `json:"etcd,omitempty"`

	// fencingAgent contains the status of the fencing agent pacemaker resource on this node.
	// The fencing agent is used to fence the other node during a quorum loss event.
	// +optional
	FencingAgent *PacemakerClusterResourceStatus `json:"fencingAgent,omitempty"`
}

// PacemakerClusterResourceStatus represents the status of a pacemaker resource on a node.
// A pacemaker resource is a unit of work managed by pacemaker. In pacemaker terminology, resources are services or
// applications that pacemaker monitors, starts, stops, and moves between nodes to maintain high availability.
// For TNF, we manage three resources:
//  - kubelet (the Kubernetes node agent and a prerequisite for etcd),
//  - etcd (the distributed key-value store)
//  - fencing agent (used to isolate failed nodes during a quorum loss event)
type PacemakerClusterResourceStatus struct {
	// conditions represent the observations of the resource's current state.
	// Known condition types are: "PacemakerClusterResourceHealthy", "PacemakerClusterResourceNotInMaintenance",
	// "PacemakerClusterResourceManaged", "PacemakerClusterResourceEnabled", "PacemakerClusterResourceOperational",
	// "PacemakerClusterResourceActive", "PacemakerClusterResourceStarted", "PacemakerClusterResourceUnblocked".
	// The "PacemakerClusterResourceHealthy" condition is an aggregate that tracks the overall health of the resource.
	// The "PacemakerClusterResourceNotInMaintenance" condition tracks whether the resource is not in maintenance mode.
	// The "PacemakerClusterResourceManaged" condition tracks whether the resource is managed by pacemaker.
	// The "PacemakerClusterResourceEnabled" condition tracks whether the resource is enabled.
	// The "PacemakerClusterResourceOperational" condition tracks whether the resource is operational (not failed).
	// The "PacemakerClusterResourceActive" condition tracks whether the resource is active (available to be used).
	// The "PacemakerClusterResourceStarted" condition tracks whether the resource is started.
	// The "PacemakerClusterResourceUnblocked" condition tracks whether the resource is unblocked.
	// Each of these conditions is required.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceHealthy')",message="conditions must contain a condition of type PacemakerClusterResourceHealthy"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceNotInMaintenance')",message="conditions must contain a condition of type PacemakerClusterResourceNotInMaintenance"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceManaged')",message="conditions must contain a condition of type PacemakerClusterResourceManaged"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceEnabled')",message="conditions must contain a condition of type PacemakerClusterResourceEnabled"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceOperational')",message="conditions must contain a condition of type PacemakerClusterResourceOperational"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceActive')",message="conditions must contain a condition of type PacemakerClusterResourceActive"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceStarted')",message="conditions must contain a condition of type PacemakerClusterResourceStarted"
	// +kubebuilder:validation:XValidation:rule="self.exists(c, c.type == 'PacemakerClusterResourceUnblocked')",message="conditions must contain a condition of type PacemakerClusterResourceUnblocked"
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// PacemakerClusterList contains a list of PacemakerCluster objects. PacemakerCluster is a cluster-scoped singleton
// resource; only one instance named "cluster" may exist. This list type exists only to satisfy Kubernetes API
// conventions.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type PacemakerClusterList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of PacemakerCluster objects.
	Items []PacemakerCluster `json:"items"`
}
