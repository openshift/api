package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PacemakerDaemonStateType represents the state of the pacemaker daemon
// +kubebuilder:validation:Enum=Running;KnownNotRunning
type PacemakerDaemonStateType string

const (
	// PacemakerDaemonStateRunning indicates the pacemaker daemon is in the 'running' state.
	PacemakerDaemonStateRunning PacemakerDaemonStateType = "Running"

	// PacemakerDaemonStateNotRunning indicates the pacemaker daemon is not in the 'running' state.
	// This is left as a blanket state to cover states like init, wait_for_ping, starting_daemons, shutting_down, shutdown_complete, etc.
	PacemakerDaemonStateNotRunning PacemakerDaemonStateType = "KnownNotRunning"
)

// QuorumStatusType represents the quorum status of a Pacemaker cluster
// +kubebuilder:validation:Enum=Quorate;NoQuorum
type QuorumStatusType string

const (
	// QuorumStatusQuorate indicates the cluster has quorum
	QuorumStatusQuorate QuorumStatusType = "Quorate"

	// QuorumStatusNoQuorum indicates the cluster does not have quorum
	QuorumStatusNoQuorum QuorumStatusType = "NoQuorum"
)

// NodeOnlineStatusType represents whether a node is online or offline
// +kubebuilder:validation:Enum=Online;Offline
type NodeOnlineStatusType string

const (
	// NodeOnlineStatusOnline indicates the node is online
	NodeOnlineStatusOnline NodeOnlineStatusType = "Online"

	// NodeOnlineStatusOffline indicates the node is offline
	NodeOnlineStatusOffline NodeOnlineStatusType = "Offline"
)

// NodeModeType represents whether a node is in active or standby mode
// +kubebuilder:validation:Enum=Active;Standby
type NodeModeType string

const (
	// NodeModeActive indicates the node is in active mode
	NodeModeActive NodeModeType = "Active"

	// NodeModeStandby indicates the node is in standby mode
	NodeModeStandby NodeModeType = "Standby"
)

// ResourceRoleType represents the role of a resource in the Pacemaker cluster
// We don't use promoted and unpromoted, so resources in those roles would omit the role field.
// +kubebuilder:validation:Enum=Started;Stopped
type ResourceRoleType string

const (
	// ResourceRoleStarted indicates the resource is started
	ResourceRoleStarted ResourceRoleType = "Started"

	// ResourceRoleStopped indicates the resource is stopped
	ResourceRoleStopped ResourceRoleType = "Stopped"
)

// ResourceActiveStatusType represents whether a resource is active or inactive
// +kubebuilder:validation:Enum=Active;Inactive
type ResourceActiveStatusType string

const (
	// ResourceActiveStatusActive indicates the resource is active
	ResourceActiveStatusActive ResourceActiveStatusType = "Active"

	// ResourceActiveStatusInactive indicates the resource is inactive
	ResourceActiveStatusInactive ResourceActiveStatusType = "Inactive"
)

// FencingActionType represents the action taken during a fencing event
// +kubebuilder:validation:Enum=reboot;off;on
type FencingActionType string

const (
	// FencingActionReboot indicates the node was rebooted
	FencingActionReboot FencingActionType = "reboot"

	// FencingActionOff indicates the node was turned off
	FencingActionOff FencingActionType = "off"

	// FencingActionOn indicates the node was turned on
	FencingActionOn FencingActionType = "on"
)

// FencingStatusType represents the status of a fencing event
// +kubebuilder:validation:Enum=success;failed
type FencingStatusType string

const (
	// FencingStatusSuccess indicates the fencing event was successful
	FencingStatusSuccess FencingStatusType = "success"

	// FencingStatusFailed indicates the fencing event failed
	FencingStatusFailed FencingStatusType = "failed"

	// FencingStatusPending indicates the fencing event is pending
	FencingStatusPending FencingStatusType = "pending"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// # PacemakerCluster represents the current state of the Pacemaker cluster as reported by the pcs status command
//
// This resource provides a view into the health and status of a Pacemaker-managed cluster in dual-replica (two-node)
// deployments. The status is periodically collected by a privileged controller and made available for monitoring
// and health checking purposes.
//
// Design Principle: Act on Deterministic Information
// Almost all fields are optional and will be populated if the data is available. If a field is not populated, it means that
// no actions can be taken by the cluster-etcd-operator with regard to whether pacemaker is healthy. This means that the
// operator will only transition between the PacemakerHealthy and PacemakerDegraded states based on deterministic information. Otherwise,
// the last known PacemakerHealthy or PacemakerDegraded status will be preserved.
//
// Some examples actions taken on deterministic information would be:
// - If the cluster is known to have quorum and critical resources are started, the operator report PacemakerHealthy state.
// - If the cluster is known to have lost quorum, or one or more critical resources are not started, the operator report PacemakerDegraded state.
// - If the cluster is trying to replace a failed node, it will check to see if the replacement node matches the expected node configuration in
//   pacemaker before proceeding with node replacement operations.

// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pacemakerclusters,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:file-pattern=cvoRunLevel=0000_25,operatorName=etcd,operatorOrdering=01,operatorComponent=two-node-fencing
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2544
type PacemakerCluster struct {
	metav1.TypeMeta   `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is an empty spec to satisfy Kubernetes API conventions.
	// PacemakerCluster is a status-only resource and does not use spec for configuration.
	// +optional
	Spec *PacemakerClusterSpec `json:"spec,omitempty"`

	// status contains the actual pacemaker cluster status information collected from the cluster.
	// This follows the Design Principle: Act on Deterministic Information.
	// When not present, pacemaker status is treated as unknown and no actions are taken by the cluster-etcd-operator.
	// +optional
	Status PacemakerClusterStatus `json:"status,omitempty,omitzero"`
}

// PacemakerClusterSpec is an empty spec as PacemakerCluster is a status-only resource
type PacemakerClusterSpec struct {
}

// PacemakerClusterStatus contains the actual pacemaker cluster status information
type PacemakerClusterStatus struct {
	// lastUpdated is the timestamp when this status was last updated
	// This is the only required field in the status object because we can use it to warn if the status collection has gone stale.
	// +required
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`

	// rawXML contains the raw XML output from pcs status xml command.
	// Kept for debugging purposes only; healthcheck should not need to parse this.
	// When present, it must be between 1 and 262144 characters long (max 256KB).
	// When not present, the raw XML output is not available.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=262144
	// +optional
	RawXML string `json:"rawXML,omitempty"`

	// collectionError contains any error encountered while collecting status
	// When present, it must be between 1 and 2048 characters long (max 2KB).
	// When not present, no collection errors are available and the status collection is assumed to be successful.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +optional
	CollectionError string `json:"collectionError,omitempty"`

	// summary provides high-level counts and flags for the cluster state
	// When present, it must be a valid PacemakerSummary object.
	// When not present, the summary is not available. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	Summary *PacemakerSummary `json:"summary,omitempty"`

	// nodes provides detailed information about each node in the cluster
	// When present, it must be a list of 1 or 2 PacemakerNodeStatus objects. Two is expected in a healthy cluster.
	// When not present, the nodes are not available. This likely indicates that there is an error parsing the raw XML output.
  // If only one node is present, this indicates that the cluster is in the process of replacing a failed node.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +optional
	Nodes []PacemakerNodeStatus `json:"nodes,omitempty"`

	// resources provides detailed information about each resource in the cluster
	// When present, it must be a list of 1 or more PacemakerResourceStatus objects.
	// When not present, the resources are not available. This likely indicates that there is an error parsing the raw XML output.
	// The number of resources is expected to be between 1 and 16, but is most likely to be exactly 6.
	// The critical resources that expect to run on both nodes are: kubelet, etcd, and a fencing resource (i.e. redfish) for each node.
	// This could drift over time as Two Node Fencing matures, so this is left flexible.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +optional
	Resources []PacemakerResourceStatus `json:"resources,omitempty"`

	// nodeHistory provides recent operation history for troubleshooting
	// When present, it must be a list of 1 or more PacemakerNodeHistoryEntry objects.
	// When not present, the node history is not available. This is the expected status for a healthy cluster.
	// Node history being capped at 16 is a reasonable limit to prevent abuse of the API, since the action history reported by the cluster
	// needs to be reported in a presentable format.
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +optional
	NodeHistory []PacemakerNodeHistoryEntry `json:"nodeHistory,omitempty"`

	// fencingHistory provides recent fencing events
	// When present, it must be a list of 1 or more PacemakerFencingEvent objects.
	// When not present, the fencing history is not available. This is the expected status for a healthy cluster.
	// Fencing history being capped at 16 is a reasonable limit to prevent abuse of the API, since the fencing history reported by the cluster
	// needs to be reported in a presentable format.
	// +listType=atomic
 	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +optional
	FencingHistory []PacemakerFencingEvent `json:"fencingHistory,omitempty"`
}

// PacemakerSummary provides a high-level summary of cluster state
type PacemakerSummary struct {
	// pacemakerDaemonState indicates the state of the pacemaker daemon
	// PacemakerDaemonStateType can be one of the following values:
	// - Running - the pacemaker daemon is in the 'running' state
	// - KnownNotRunning - the pacemaker daemon is not in the 'running' state. This is left as a blanket state
	//   to cover states like init, wait_for_ping, starting_daemons, shutting_down, shutdown_complete, etc.
	// +optional
	PacemakerDaemonState PacemakerDaemonStateType `json:"pacemakerDaemonState,omitempty"`

	// quorumStatus indicates if the cluster has quorum
	// QuorumStatusType can be one of the following values:
	// - Quorate - the cluster has quorum
	// - NoQuorum - the cluster does not have quorum
	// +optional
	QuorumStatus QuorumStatusType `json:"quorumStatus,omitempty"`

	// nodesOnline is the count of online nodes
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +optional
	NodesOnline *int32 `json:"nodesOnline,omitempty"`

	// nodesTotal is the total count of configured nodes
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +optional
	NodesTotal *int32 `json:"nodesTotal,omitempty"`

	// resourcesStarted is the count of started resources
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=16
	// +optional
	ResourcesStarted *int32 `json:"resourcesStarted,omitempty"`

	// resourcesTotal is the total count of configured resources
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=16
	// +optional
	ResourcesTotal *int32 `json:"resourcesTotal,omitempty"`
}

// NodeStatus represents the status of a single node in the Pacemaker cluster
type PacemakerNodeStatus struct {
	// name is the name of the node
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Name string `json:"name,omitempty"`

	// ipv4Address is the IPv4 address of the node, if registered via IPv4
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=15
  // +kubebuilder:validation:Pattern="^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	// +optional
	IPv4Address string `json:"ipv4Address,omitempty"`

	// ipv6Address is the IPv6 address of the node, if registered via IPv6
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=39
	// +kubebuilder:validation:Format=ipv6
	// +kubebuilder:validation:Pattern=`^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	// +optional
	IPv6Address string `json:"ipv6Address,omitempty"`

	// onlineStatus indicates if the node is online or offline
	// NodeOnlineStatusType can be one of the following values:
	// - Online - the node is online
	// - Offline - the node is offline
	// +optional
	OnlineStatus NodeOnlineStatusType `json:"onlineStatus,omitempty"`

	// mode indicates if the node is in active or standby mode
	// NodeModeType can be one of the following values:
	// - Active - the node is in active mode
	// - Standby - the node is in standby mode
	// +optional
	Mode NodeModeType `json:"mode,omitempty"`
}

// PacemakerResourceStatus represents the status of a single resource in the Pacemaker cluster
type PacemakerResourceStatus struct {
	// name is the name of the resource
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Name string `json:"name,omitempty"`

	// resourceAgent is the resource agent type (e.g., "ocf:heartbeat:IPaddr2", "systemd:kubelet")
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	ResourceAgent string `json:"resourceAgent,omitempty"`

	// role is the current role of the resource
	// ResourceRoleType can be one of the following values:
	// - Started - the resource is started
	// - Stopped - the resource is stopped
	// We don't use promoted and unpromoted, so resources in those roles would omit the role field.
	// +optional
	Role ResourceRoleType `json:"role,omitempty"`

	// activeStatus indicates if the resource is active or inactive
	// ResourceActiveStatusType can be one of the following values:
	// - Active - the resource is active
	// - Inactive - the resource is inactive
	// +optional
	ActiveStatus ResourceActiveStatusType `json:"activeStatus,omitempty"`

	// node is the node where the resource is running
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Node string `json:"node,omitempty"`
}

// PacemakerNodeHistoryEntry represents a single operation history entry from node_history
type PacemakerNodeHistoryEntry struct {
	// node is the node where the operation occurred
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Node string `json:"node,omitempty"`

	// resource is the resource that was operated on
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Resource string `json:"resource,omitempty"`

	// operation is the operation that was performed (e.g., "monitor", "start", "stop")
	// Unlike other fields, this is not an enum because while "monitor", "start" and "stop"
	// are the most common, resource agents can define their own operations.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=16
	// +optional
	Operation string `json:"operation,omitempty"`

	// rc is the return code from the operation
	// +optional
	RC *int32 `json:"rc,omitempty"`

	// rcText is the human-readable return code text (e.g., "ok", "error", "not running")
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=16
	// +optional
	RCText string `json:"rcText,omitempty"`

	// lastRCChange is the timestamp when the RC last changed
	// +optional
	LastRCChange metav1.Time `json:"lastRCChange,omitempty"`
}

// PacemakerFencingEvent represents a single fencing event from fence history
type PacemakerFencingEvent struct {
	// target is the node that was fenced
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Target string `json:"target,omitempty"`

	// action is the fencing action performed
	// FencingActionType can be one of the following values:
	// - reboot - the node was rebooted
	// - off - the node was turned off
	// - on - the node was turned on
	// +optional
	Action FencingActionType `json:"action,omitempty"`

	// status is the status of the fencing operation
	// FencingStatusType can be one of the following values:
	// - success - the fencing event was successful
	// - failed - the fencing event failed
	// - pending - the fencing event is pending
	// +optional
	Status FencingStatusType `json:"status,omitempty"`

	// completed is the timestamp when the fencing event was completed
	// +optional
	Completed metav1.Time `json:"completed,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4

// PacemakerClusterList contains a list of PacemakerCluster objects.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type PacemakerClusterList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of PacemakerCluster objects.
	Items []PacemakerCluster `json:"items"`
}
