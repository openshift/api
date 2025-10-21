package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PacemakerDaemonStateType represents the state of the pacemaker daemon
// PacemakerDaemonStateType can be one of the following values:
// - Running - the pacemaker daemon is in the 'running' state
// - KnownNotRunning - the pacemaker daemon is not in the 'running' state. This is left as a blanket state
//   to cover states like init, wait_for_ping, starting_daemons, shutting_down, shutdown_complete, etc.
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
// QuorumStatusType can be one of the following values:
// - Quorate - the cluster has quorum
// - NoQuorum - the cluster does not have quorum
// +kubebuilder:validation:Enum=Quorate;NoQuorum
type QuorumStatusType string

const (
	// QuorumStatusQuorate indicates the cluster has quorum
	QuorumStatusQuorate QuorumStatusType = "Quorate"

	// QuorumStatusNoQuorum indicates the cluster does not have quorum
	QuorumStatusNoQuorum QuorumStatusType = "NoQuorum"
)

// NodeOnlineStatusType represents whether a node is online or offline
// NodeOnlineStatusType can be one of the following values:
// - Online - the node is online
// - Offline - the node is offline
// +kubebuilder:validation:Enum=Online;Offline
type NodeOnlineStatusType string

const (
	// NodeOnlineStatusOnline indicates the node is online
	NodeOnlineStatusOnline NodeOnlineStatusType = "Online"

	// NodeOnlineStatusOffline indicates the node is offline
	NodeOnlineStatusOffline NodeOnlineStatusType = "Offline"
)

// NodeModeType represents whether a node is in active or standby mode
// NodeModeType can be one of the following values:
// - Active - the node is in active mode
// - Standby - the node is in standby mode
// +kubebuilder:validation:Enum=Active;Standby
type NodeModeType string

const (
	// NodeModeActive indicates the node is in active mode
	NodeModeActive NodeModeType = "Active"

	// NodeModeStandby indicates the node is in standby mode
	NodeModeStandby NodeModeType = "Standby"
)

// ResourceRoleType represents the role of a resource in the Pacemaker cluster
// ResourceRoleType can be one of the following values:
// - Started - the resource is started
// - Stopped - the resource is stopped
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
// ResourceActiveStatusType can be one of the following values:
// - Active - the resource is active
// - Inactive - the resource is inactive
// +kubebuilder:validation:Enum=Active;Inactive
type ResourceActiveStatusType string

const (
	// ResourceActiveStatusActive indicates the resource is active
	ResourceActiveStatusActive ResourceActiveStatusType = "Active"

	// ResourceActiveStatusInactive indicates the resource is inactive
	ResourceActiveStatusInactive ResourceActiveStatusType = "Inactive"
)

// FencingActionType represents the action taken during a fencing event
// FencingActionType can be one of the following values:
// - Reboot - the node was rebooted
// - Off - the node was turned off
// - On - the node was turned on
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
// FencingStatusType can be one of the following values:
// - Success - the fencing event was successful
// - Failed - the fencing event failed
// - Pending - the fencing event is pending
// +kubebuilder:validation:Enum=success;failed;pending
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
	metav1.TypeMeta `json:",inline"`

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
	// When fields are not present, they are treated as unknown and no actions are taken by the cluster-etcd-operator.
	// +optional
	Status *PacemakerClusterStatus `json:"status,omitempty,omitzero"`
}

// PacemakerClusterSpec is an empty spec as PacemakerCluster is a status-only resource
type PacemakerClusterSpec struct {
}

// PacemakerClusterStatus contains the actual pacemaker cluster status information
// As part of validating the status object, we need to ensure that the lastUpdated timestamp is always newer than the current value.
// We allow the lastUpdated timestamp to be empty on initial creation, but it is required once it has been set.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.lastUpdated) || (has(self.lastUpdated) && self.lastUpdated >= oldSelf.lastUpdated)",message="lastUpdated must be a newer timestamp than the current value"
type PacemakerClusterStatus struct {
	// lastUpdated is the timestamp when this status was last updated
	// When present, it must be a valid timestamp in RFC3339 format and cannot be set to an earlier timestamp than the current value
	// This is the most critical field in the status object because we can use it to warn if the status collection has gone stale.
	// It is optional upon initial creation, but required once it has been set.
	// +kubebuilder:validation:Format=RFC3339
	// +optional
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`

	// rawXML contains the raw XML output from pcs status xml command.
	// Kept for debugging purposes only; healthcheck should not need to parse this.
	// When present, it must be between 1 and 262144 characters long (max 256KB). This limit protects the API from XML bombs and excessive memory consumption.
	// When not present, the raw XML output is not available.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=262144
	// +optional
	RawXML string `json:"rawXML,omitempty"`

	// collectionError contains any error encountered while collecting status
	// When present, it must be between 1 and 2048 characters long (max 2KB). This limit ensures that the error message can be displayed in a human-readable format.
	// When not present, no collection errors are available and the status collection is assumed to be successful.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +optional
	CollectionError string `json:"collectionError,omitempty"`

	// summary provides high-level counts and flags for the cluster state
	// TNF clusters don't use advanced pacemaker features like maintenance mode or pacemaker remote nodes, so these are omitted from the summary.
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
	// is presented as a recorded event.
	// The healthchecker runs every 30 seconds and creates events for failed operations that occured within the last 5 minutes (unless they have already been reported).
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +optional
	NodeHistory []PacemakerNodeHistoryEntry `json:"nodeHistory,omitempty"`

	// fencingHistory provides recent fencing events
	// When present, it must be a list of 1 or more PacemakerFencingEvent objects.
	// When not present, the fencing history is not available. This is the expected status for a healthy cluster.
	// Fencing history being capped at 16 is a reasonable limit to prevent abuse of the API, since the fencing history reported by the cluster
	// is presented as a recorded event.
	// The healthchecker runs every 30 seconds and creates events for fencing events that occured within the last 24 hours (unless they have already been reported).
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
	// When present, it must be a valid PacemakerDaemonStateType.
	// When not present, the pacemaker daemon state is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	PacemakerDaemonState PacemakerDaemonStateType `json:"pacemakerDaemonState,omitempty"`

	// quorumStatus indicates if the cluster has quorum
	// QuorumStatusType can be one of the following values:
	// - Quorate - the cluster has quorum
	// - NoQuorum - the cluster does not have quorum
	// When present, it must be a valid QuorumStatusType.
	// When not present, the quorum status is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	QuorumStatus QuorumStatusType `json:"quorumStatus,omitempty"`

	// nodesOnline is the count of online nodes
	// When present, it must be a valid integer between 0 and 2 inclusive.
	// When not present, the nodes online count is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +optional
	NodesOnline *int32 `json:"nodesOnline,omitempty"`

	// nodesTotal is the total count of configured nodes
	// When present, it must be a valid integer between 0 and 2 inclusive.
	// When not present, the nodes total count is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +optional
	NodesTotal *int32 `json:"nodesTotal,omitempty"`

	// resourcesStarted is the count of started resources
	// When present, it must be a valid integer between 0 and 16 inclusive. For a healthy Two Node Fencing (TNF) cluster, this is expected to be 6.
	// The expected resources are kubelet, etcd, and a fencing resource (i.e. redfish) for each node.
	// The number could be less than 6 if the cluster is starting up or not healthy.
	// The total number of resources managed by the cluster could drift over time as Two Node Fencing matures, so this is left flexible.
	// When not present, the resources started count is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=16
	// +optional
	ResourcesStarted *int32 `json:"resourcesStarted,omitempty"`

	// resourcesTotal is the total count of configured resources
	// When present, it must be a valid integer between 0 and 16 inclusive. For a healthy Two Node Fencing (TNF) cluster, this is expected to be 6.
	// The expected resources are kubelet, etcd, and a fencing resource (i.e. redfish) for each node.
	// The total number of resources managed by the cluster could drift over time as Two Node Fencing matures, so this is left flexible.
	// When not present, the resources total count is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=16
	// +optional
	ResourcesTotal *int32 `json:"resourcesTotal,omitempty"`
}

// NodeStatus represents the status of a single node in the Pacemaker cluster
type PacemakerNodeStatus struct {
	// name is the name of the node
	// It must be a valid string between 1 and 256 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +required
	Name string `json:"name,omitempty"`

	// ipAddress is the canonical IPv4 or IPv6 address of the node
	// When present, it must be a valid canonical global unicast IPv4 or IPv6 address (including private/RFC1918 addresses).
	// This excludes special addresses like unspecified, loopback, link-local, and multicast.
	// When not present, the IP address is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=39
	// +kubebuilder:validation:XValidation:rule="isIP(self) && ip.isCanonical(self) && ip(self).isGlobalUnicast()",message="ipAddress must be a valid canonical global unicast IPv4 or IPv6 address"
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// onlineStatus indicates if the node is online or offline
	// NodeOnlineStatusType can be one of the following values:
	// - Online - the node is online
	// - Offline - the node is offline
	// When present, it must be a valid NodeOnlineStatusType.
	// When not present, the online status is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	OnlineStatus NodeOnlineStatusType `json:"onlineStatus,omitempty"`

	// mode indicates if the node is in active or standby mode
	// NodeModeType can be one of the following values:
	// - Active - the node is in active mode
	// - Standby - the node is in standby mode
	// When present, it must be a valid NodeModeType.
	// When not present, the node mode is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	Mode NodeModeType `json:"mode,omitempty"`
}

// PacemakerResourceStatus represents the status of a single resource in the Pacemaker cluster
type PacemakerResourceStatus struct {
	// name is the name of the resource
	// It must be a valid string between 1 and 256 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +required
	Name string `json:"name,omitempty"`

	// resourceAgent is the resource agent type (e.g., "ocf:heartbeat:IPaddr2", "systemd:kubelet")
	// When present, it must be a valid string between 1 and 256 characters long.
	// When not present, the resource agent type is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	ResourceAgent string `json:"resourceAgent,omitempty"`

	// role is the current role of the resource
	// ResourceRoleType can be one of the following values:
	// - Started - the resource is started
	// - Stopped - the resource is stopped
	// We don't use promoted and unpromoted, so resources in those roles would omit the role field.
	// When present, it must be a valid ResourceRoleType.
	// When not present, the resource role is unknown (or an unsupported type like promoted or unpromoted). This likely indicates that there is an error parsing the raw XML output.
	// +optional
	Role ResourceRoleType `json:"role,omitempty"`

	// activeStatus indicates if the resource is active or inactive
	// ResourceActiveStatusType can be one of the following values:
	// - Active - the resource is active
	// - Inactive - the resource is inactive
	// When present, it must be a valid ResourceActiveStatusType.
	// When not present, the active status is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	ActiveStatus ResourceActiveStatusType `json:"activeStatus,omitempty"`

	// node is the node where the resource is running
	// When present, it must be a valid string between 1 and 256 characters long.
	// When not present, the resource is not assigned to a node. This typically indicates a stopped or unscheduled resource. It could also imply an error parsing the raw XML output.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +optional
	Node string `json:"node,omitempty"`
}

// PacemakerNodeHistoryEntry represents a single operation history entry from node_history
type PacemakerNodeHistoryEntry struct {
	// node is the node where the operation occurred
	// It must be a valid string between 1 and 256 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +required
	Node string `json:"node,omitempty"`

	// resource is the resource that was operated on
	// It must be a valid string between 1 and 256 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +required
	Resource string `json:"resource,omitempty"`

	// operation is the operation that was performed (e.g., "monitor", "start", "stop")
	// Unlike other fields, this is not an enum because while "monitor", "start" and "stop"
	// are the most common, resource agents can define their own operations.
	// It must be a valid string between 1 and 32 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +required
	Operation string `json:"operation,omitempty"`

	// rc is the return code from the operation
	// When present, it must be a valid integer between 0 and 2147483647 (max 32-bit int) inclusive.
	// When not present, the return code is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2147483647
	// +optional
	RC *int32 `json:"rc,omitempty"`

	// rcText is the human-readable return code text (e.g., "ok", "error", "not running")
	// When present, it must be a valid string between 1 and 32 characters long. This is a human-readable string and is not validated against any specific format.
	// When not present, the return code text is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +optional
	RCText string `json:"rcText,omitempty"`

	// lastRCChange is the timestamp when the RC last changed
	// It must be a valid timestamp in RFC3339 format and cannot be empty.
	// +kubebuilder:validation:Format=RFC3339
	// +required
	LastRCChange metav1.Time `json:"lastRCChange,omitempty"`
}

// PacemakerFencingEvent represents a single fencing event from fence history
type PacemakerFencingEvent struct {
	// target is the node that was fenced
	// It must be a valid string between 1 and 256 characters long and cannot be empty.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +required
	Target string `json:"target,omitempty"`

	// action is the fencing action performed
	// FencingActionType can be one of the following values:
	// - reboot - the node was rebooted
	// - off - the node was turned off
	// - on - the node was turned on
	// When present, it must be a valid FencingActionType.
	// When not present, the fencing action is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	Action FencingActionType `json:"action,omitempty"`

	// status is the status of the fencing operation
	// FencingStatusType can be one of the following values:
	// - success - the fencing event was successful
	// - failed - the fencing event failed
	// - pending - the fencing event is pending
	// When present, it must be a valid FencingStatusType.
	// When not present, the fencing status is unknown. This likely indicates that there is an error parsing the raw XML output.
	// +optional
	Status FencingStatusType `json:"status,omitempty"`

	// completed is the timestamp when the fencing event was completed
	// It must be a valid timestamp in RFC3339 format and cannot be empty.
	// +kubebuilder:validation:Format=RFC3339
	// +required
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
