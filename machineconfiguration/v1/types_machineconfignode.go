package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=machineconfignodes,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2171
// +openshift:file-pattern=cvoRunLevel=0000_80,operatorName=machine-config,operatorOrdering=01
// +openshift:enable:FeatureGate=MachineConfigNodes
// +kubebuilder:printcolumn:name="PoolName",type="string",JSONPath=.spec.pool.name,priority=0
// +kubebuilder:printcolumn:name="DesiredConfig",type="string",JSONPath=.spec.configVersion.desired,priority=0
// +kubebuilder:printcolumn:name="CurrentConfig",type="string",JSONPath=.status.configVersion.current,priority=0
// +kubebuilder:printcolumn:name="Updated",type="string",JSONPath=.status.conditions[?(@.type=="Updated")].status,priority=0
// +kubebuilder:printcolumn:name="UpdatePrepared",type="string",JSONPath=.status.conditions[?(@.type=="UpdatePrepared")].status,priority=1
// +kubebuilder:printcolumn:name="UpdateExecuted",type="string",JSONPath=.status.conditions[?(@.type=="UpdateExecuted")].status,priority=1
// +kubebuilder:printcolumn:name="UpdatePostActionComplete",type="string",JSONPath=.status.conditions[?(@.type=="UpdatePostActionComplete")].status,priority=1
// +kubebuilder:printcolumn:name="UpdateComplete",type="string",JSONPath=.status.conditions[?(@.type=="UpdateComplete")].status,priority=1
// +kubebuilder:printcolumn:name="Resumed",type="string",JSONPath=.status.conditions[?(@.type=="Resumed")].status,priority=1
// +kubebuilder:printcolumn:name="UpdateCompatible",type="string",JSONPath=.status.conditions[?(@.type=="UpdateCompatible")].status,priority=1
// +kubebuilder:printcolumn:name="UpdatedFilesAndOS",type="string",JSONPath=.status.conditions[?(@.type=="AppliedFilesAndOS")].status,priority=1
// +kubebuilder:printcolumn:name="CordonedNode",type="string",JSONPath=.status.conditions[?(@.type=="Cordoned")].status,priority=1
// +kubebuilder:printcolumn:name="DrainedNode",type="string",JSONPath=.status.conditions[?(@.type=="Drained")].status,priority=1
// +kubebuilder:printcolumn:name="RebootedNode",type="string",JSONPath=.status.conditions[?(@.type=="RebootedNode")].status,priority=1
// +kubebuilder:printcolumn:name="ReloadedCRIO",type="string",JSONPath=.status.conditions[?(@.type=="ReloadedCRIO")].status,priority=1
// +kubebuilder:printcolumn:name="UncordonedNode",type="string",JSONPath=.status.conditions[?(@.type=="Uncordoned")].status,priority=1
// +kubebuilder:metadata:labels=openshift.io/operator-managed=

// MachineConfigNode describes the health of the Machines on the system
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +kubebuilder:validation:XValidation:rule="self.metadata.name == self.spec.node.name",message="spec.node.name should match metadata.name"
type MachineConfigNode struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the configuration of the machine config node.
	// +required
	Spec MachineConfigNodeSpec `json:"spec"`

	// status describes the last observed state of this machine config node.
	// +optional
	Status MachineConfigNodeStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigNodeList describes all of the MachinesStates on the system
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfigNodeList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items contains a collection of MachineConfigNode resources.
	// +kubebuilder:validation:MaxItems=100000
	// +optional
	// TODO: Decide on an appropriate maxItems value for this list of MachineConfigNodes. There is a 1:1 ratio between # of nodes & MCN objects.
	Items []MachineConfigNode `json:"items"`
}

// MachineConfigNodeSpec describes the MachineConfigNode we are managing.
type MachineConfigNodeSpec struct {
	// node contains a reference to the node for this machine config node.
	// +required
	Node MCOObjectReference `json:"node"`

	// pool contains a reference to the machine config pool that this machine config node's
	// referenced node belongs to.
	// +required
	Pool MCOObjectReference `json:"pool"`

	// configVersion holds the desired config version for the node targeted by this machine config node resource.
	// The desired version represents the machine config the node will attempt to update to. This gets set before the machine config operator validates
	// the new machine config against the current machine config.
	// +required
	// TODO: Potentially remove since it might be redundant. There may need to be some changes to how ObservedGeneration is updated if this field
	// is consolidated with MachineConfigNodeStatus.ConfigVersion.Current.
	ConfigVersion MachineConfigNodeSpecMachineConfigVersion `json:"configVersion"`

	// pinnedImageSets is a user defined value that holds the names of the desired image sets that the node should pull and pin.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=100
	// +optional
	PinnedImageSets []MachineConfigNodeSpecPinnedImageSet `json:"pinnedImageSets,omitempty"`
}

// MCOObjectReference holds information about an object the MCO either owns
// or modifies in some way
type MCOObjectReference struct {
	// name is the name of the object being referenced. For example, this can represent a pool
	// or node name.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Name string `json:"name"`
}

// MachineConfigNodeStatus holds the reported information on a particular machine config node.
type MachineConfigNodeStatus struct {
	// conditions represent the observations of a machine config node's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=15
	// +optional
	// TODO: Reevaluate maxItems value as status are trimmed down/API structure is updated.
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// observedGeneration represents the generation of the MachineConfigNode object observed by the Machine Config Operator's controller.
	// This field is updated when the controller observes a change to the desiredConfig in the configVersion of the machine config node spec.
	// +kubebuilder:validation:XValidation:rule="self >= oldSelf", message="observedGeneration must not move backwards"
	// +kubebuilder:validation:Minimum=0
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// configVersion describes the current and desired machine config version for this node.
	// +required
	ConfigVersion MachineConfigNodeStatusMachineConfigVersion `json:"configVersion"`
	// pinnedImageSets describes the current and desired pinned image sets for this node.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=100
	// +optional
	PinnedImageSets []MachineConfigNodeStatusPinnedImageSet `json:"pinnedImageSets,omitempty"`
}

// MachineConfigNodeStatusPinnedImageSet holds information about the current and desired pinned image sets for the current observed machine config node.
// +kubebuilder:validation:XValidation:rule="has(self.desiredGeneration) && has(self.currentGeneration) ? self.desiredGeneration >= self.currentGeneration : true",message="desired generation must be greater than or equal to the current generation"
// +kubebuilder:validation:XValidation:rule="has(self.lastFailedGeneration) && has(self.desiredGeneration) ? self.desiredGeneration >= self.lastFailedGeneration : true",message="desired generation must be greater than or equal to last failed generation"
// +kubebuilder:validation:XValidation:rule="has(self.lastFailedGeneration) ? has(self.desiredGeneration): true",message="desired generation must be defined if last failed generation is defined"
type MachineConfigNodeStatusPinnedImageSet struct {
	// name is the name of the pinned image set.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Name string `json:"name"`
	// currentGeneration is the generation of the pinned image set that has most recently been successfully pulled and pinned on this node.
	// +optional
	CurrentGeneration int32 `json:"currentGeneration,omitempty"`
	// desiredGeneration is the generation of the pinned image set that is targeted to be pulled and pinned on this node.
	// +kubebuilder:validation:Minimum=0
	// +optional
	DesiredGeneration int32 `json:"desiredGeneration,omitempty"`
	// lastFailedGeneration is the generation of the most recent pinned image set that failed to be pulled and pinned on this node.
	// +kubebuilder:validation:Minimum=0
	// +optional
	LastFailedGeneration int32 `json:"lastFailedGeneration,omitempty"`
	// lastFailedGenerationErrors is a list of errors explaining why the lastFailed generation failed to be pulled and pinned.
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=message
	// +optional
	LastFailedGenerationErrors []MachineConfigNodeStatusPinnedImageSetError `json:"lastFailedGenerationErrors,omitempty"`
}

// MachineConfigNodeStatusMachineConfigVersion holds the current and desired config versions as last updated in the MCN status.
// When the current and desired versions are not matched, the machine config pool is processing an upgrade and the machine config node will
// monitor the upgrade process.
// When the current and desired versions do not match, the machine config node will ignore these events given that certain operations
// happen both during the MCO's upgrade mode and the daily operations mode. //TODO: Check if this should instead say *do* match to be consistent
// with comment on lines 204-205.
type MachineConfigNodeStatusMachineConfigVersion struct {
	// current is the name of the machine config currently in use on the node.
	// This value is updated once the machine config daemon has completed the update of the configuration for the node.
	// This value should match the desired version unless an upgrade is in progress.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +optional
	Current string `json:"current"`
	// desired is the MachineConfig the node wants to upgrade to.
	// This value gets set in the machine config node status once the machine config has been validated
	// against the current machine config.
	// This value should match the desired version unless an upgrade is in progress.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Desired string `json:"desired"`
}

// MachineConfigNodeSpecMachineConfigVersion holds the desired config version for the current observed machine config node.
// When Current is not equal to Desired; the MachineConfigOperator is in an upgrade phase and the machine config node will
// take account of upgrade related events. Otherwise they will be ignored given that certain operations
// happen both during the MCO's upgrade mode and the daily operations mode.
type MachineConfigNodeSpecMachineConfigVersion struct {
	// desired is the name of the machine config that the the node should be upgraded to.
	// This value is set when the machine config pool generates a new version of its rendered configuration.
	// When this value is changed, the machine config daemon starts the node upgrade process.
	// This value gets set in the machine config node spec once the machine config has been targeted for upgrade and before it is validated.
	// This value should match the desired version unless an upgrade is in progress.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Desired string `json:"desired"`
}

// MachineConfigNodeSpecPinnedImageSet holds information on the desired pinned image sets that the current observed machine config node
// should pin and pull.
type MachineConfigNodeSpecPinnedImageSet struct {
	// name is the name of the pinned image set.
	// This value should match the desired version unless an upgrade is in progress.
	// Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) consisting
	// of only lowercase alphanumeric characters, hyphens (-), and periods (.), start and end
	// with an alphanumeric character, and be at most 253 characters in length.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Name string `json:"name"`
}

// MachineConfigNodeStatusPinnedImageSetError holds information on the the errors associated with failed attempts at pulling and pinning images
type MachineConfigNodeStatusPinnedImageSetError struct {
	// message is the message of the pinned image error.
	// +kubebuilder:validation:MaxLength=2048
	// +required
	// TODO: Update MaxLength value once this API design has been reevaluated.
	Message string `json:"message"`
}

// StateProgress is each possible state for each possible MachineConfigNodeType
// UpgradeProgression Kind will only use the "MachinConfigPoolUpdate..." types for example //TODO: figure out how to make comment more relevant.
// +enum
type StateProgress string

// TODO: Trim down to only helpful statues, prioritizing OCL status needs. Relevant PR: https://github.com/openshift/api/pull/1596
const (
	// MachineConfigNodeUpdatePrepared describes a machine that is preparing in the daemon to trigger an update
	MachineConfigNodeUpdatePrepared StateProgress = "UpdatePrepared"
	// MachineConfigNodeUpdateExecuted describes a machine that has executed the body of the upgrade
	MachineConfigNodeUpdateExecuted StateProgress = "UpdateExecuted"
	// MachineConfigNodeUpdatePostActionComplete describes a machine that has executed its post update action
	MachineConfigNodeUpdatePostActionComplete StateProgress = "UpdatePostActionComplete"
	// MachineConfigNodeUpdateComplete describes a machine that has completed the core parts of an upgrade
	MachineConfigNodeUpdateComplete StateProgress = "UpdateComplete"
	// MachineConfigNodeUpdated describes a machine that has a matching desired and current config after executing an update
	MachineConfigNodeUpdated StateProgress = "Updated"
	// MachineConfigNodeUpdateResumed describes a machine that has resumed normal processes
	MachineConfigNodeUpdateResumed StateProgress = "Resumed"
	// MachineConfigNodeUpdateCompatible the part of the preparing phase where the mco decides whether it can update
	MachineConfigNodeUpdateCompatible StateProgress = "UpdateCompatible"
	// MachineConfigNodeUpdateDrained describes the part of the inprogress phase where the node drains
	MachineConfigNodeUpdateDrained StateProgress = "Drained"
	// MachineConfigNodeUpdateFilesAndOS describes the part of the inprogress phase where the nodes file and OS config change
	MachineConfigNodeUpdateFilesAndOS StateProgress = "AppliedFilesAndOS"
	// MachineConfigNodeUpdateCordoned describes the part of the completing phase where the node cordons
	MachineConfigNodeUpdateCordoned StateProgress = "Cordoned"
	// MachineConfigNodeUpdateUncordoned describes the part of the completing phase where the node uncordons
	MachineConfigNodeUpdateUncordoned StateProgress = "Uncordoned"
	// MachineConfigNodeUpdateRebooted describes the part of the post action phase where the node reboots itself
	MachineConfigNodeUpdateRebooted StateProgress = "RebootedNode"
	// MachineConfigNodeUpdateReloaded describes the part of the post action phase where the node reloads its CRIO service
	MachineConfigNodeUpdateReloaded StateProgress = "ReloadedCRIO"
	// MachineConfigNodePinnedImageSetsProgressing describes a machine currently progressing to the desired pinned image sets
	MachineConfigNodePinnedImageSetsProgressing StateProgress = "PinnedImageSetsProgressing"
	// MachineConfigNodePinnedImageSetsDegraded describes a machine that has failed to progress to the desired pinned image sets
	MachineConfigNodePinnedImageSetsDegraded StateProgress = "PinnedImageSetsDegraded"
)
