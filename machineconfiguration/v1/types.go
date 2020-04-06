package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// MachineConfigRoleLabelKey is metadata key in the MachineConfig. Specifies the node role that config should be applied to.
// For example: `master` or `worker`
const MachineConfigRoleLabelKey = "machineconfiguration.openshift.io/role"

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfig defines the configuration for a machine
type MachineConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec MachineConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// MachineConfigSpec is the spec for MachineConfig
type MachineConfigSpec struct {
	// OSImageURL specifies the remote location that will be used to
	// fetch the OS.
	OSImageURL string `json:"osImageURL" protobuf:"bytes,1,opt,name=osImageURL"`
	// Config is a Ignition Config object.
	Config runtime.RawExtension `json:"config" protobuf:"bytes,2,opt,name=config"`

	// +nullable
	KernelArguments []string `json:"kernelArguments" protobuf:"bytes,3,rep,name=kernelArguments"`

	FIPS       bool   `json:"fips" protobuf:"varint,4,opt,name=fips"`
	KernelType string `json:"kernelType" protobuf:"bytes,5,opt,name=kernelType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigList is a list of MachineConfig resources
type MachineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []MachineConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPool describes a pool of MachineConfigs.
type MachineConfigPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +required
	Spec MachineConfigPoolSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status MachineConfigPoolStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// MachineConfigPoolSpec is the spec for MachineConfigPool resource.
type MachineConfigPoolSpec struct {
	// machineConfigSelector specifies a label selector for MachineConfigs.
	// Refer https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ on how label and selectors work.
	MachineConfigSelector *metav1.LabelSelector `json:"machineConfigSelector,omitempty" protobuf:"bytes,1,opt,name=machineConfigSelector"`

	// nodeSelector specifies a label selector for Machines
	NodeSelector *metav1.LabelSelector `json:"nodeSelector,omitempty" protobuf:"bytes,2,opt,name=nodeSelector"`

	// paused specifies whether or not changes to this machine config pool should be stopped.
	// This includes generating new desiredMachineConfig and update of machines.
	Paused bool `json:"paused" protobuf:"varint,3,opt,name=paused"`

	// maxUnavailable specifies the percentage or constant number of machines that can be updating at any given time.
	// default is 1.
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,4,opt,name=maxUnavailable"`

	// The targeted MachineConfig object for the machine config pool.
	Configuration MachineConfigPoolStatusConfiguration `json:"configuration" protobuf:"bytes,5,opt,name=configuration"`
}

// MachineConfigPoolStatus is the status for MachineConfigPool resource.
type MachineConfigPoolStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// configuration represents the current MachineConfig object for the machine config pool.
	Configuration MachineConfigPoolStatusConfiguration `json:"configuration" protobuf:"bytes,2,opt,name=configuration"`

	// machineCount represents the total number of machines in the machine config pool.
	MachineCount int32 `json:"machineCount" protobuf:"varint,3,opt,name=machineCount"`

	// updatedMachineCount represents the total number of machines targeted by the pool that have the CurrentMachineConfig as their config.
	UpdatedMachineCount int32 `json:"updatedMachineCount" protobuf:"varint,4,opt,name=updatedMachineCount"`

	// readyMachineCount represents the total number of ready machines targeted by the pool.
	ReadyMachineCount int32 `json:"readyMachineCount" protobuf:"varint,5,opt,name=readyMachineCount"`

	// unavailableMachineCount represents the total number of unavailable (non-ready) machines targeted by the pool.
	// A node is marked unavailable if it is in updating state or NodeReady condition is false.
	UnavailableMachineCount int32 `json:"unavailableMachineCount" protobuf:"varint,6,opt,name=unavailableMachineCount"`

	// degradedMachineCount represents the total number of machines marked degraded (or unreconcilable).
	// A node is marked degraded if applying a configuration failed..
	DegradedMachineCount int32 `json:"degradedMachineCount" protobuf:"varint,7,opt,name=degradedMachineCount"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []MachineConfigPoolCondition `json:"conditions" protobuf:"bytes,8,rep,name=conditions"`
}

// MachineConfigPoolStatusConfiguration stores the current configuration for the pool, and
// optionally also stores the list of MachineConfig objects used to generate the configuration.
type MachineConfigPoolStatusConfiguration struct {
	corev1.ObjectReference `json:",inline" protobuf:"bytes,1,opt,name=objectReference"`

	// source is the list of MachineConfig objects that were used to generate the single MachineConfig object specified in `content`.
	// +optional
	Source []corev1.ObjectReference `json:"source,omitempty" protobuf:"bytes,2,rep,name=source"`
}

// MachineConfigPoolCondition contains condition information for an MachineConfigPool.
type MachineConfigPoolCondition struct {
	// type of the condition, currently ('Done', 'Updating', 'Failed').
	Type MachineConfigPoolConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=MachineConfigPoolConditionType"`

	// status of the condition, one of ('True', 'False', 'Unknown').
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is a brief machine readable explanation for the condition's last
	// transition.
	Reason string `json:"reason" protobuf:"bytes,4,opt,name=reason"`

	// message is a human readable description of the details of the last
	// transition, complementing reason.
	Message string `json:"message" protobuf:"bytes,5,opt,name=message"`
}

// MachineConfigPoolConditionType valid conditions of a MachineConfigPool
type MachineConfigPoolConditionType string

const (
	// MachineConfigPoolUpdated means MachineConfigPool is updated completely.
	// When the all the machines in the pool are updated to the correct machine config.
	MachineConfigPoolUpdated MachineConfigPoolConditionType = "Updated"

	// MachineConfigPoolUpdating means MachineConfigPool is updating.
	// When at least one of machine is not either not updated or is in the process of updating
	// to the desired machine config.
	MachineConfigPoolUpdating MachineConfigPoolConditionType = "Updating"

	// MachineConfigPoolNodeDegraded means the update for one of the machine is not progressing
	MachineConfigPoolNodeDegraded MachineConfigPoolConditionType = "NodeDegraded"

	// MachineConfigPoolRenderDegraded means the rendered configuration for the pool cannot be generated because of an error
	MachineConfigPoolRenderDegraded MachineConfigPoolConditionType = "RenderDegraded"

	// MachineConfigPoolDegraded is the overall status of the pool based, today, on whether we fail with NodeDegraded or RenderDegraded
	MachineConfigPoolDegraded MachineConfigPoolConditionType = "Degraded"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPoolList is a list of MachineConfigPool resources
type MachineConfigPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []MachineConfigPool `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeletConfig describes a customized Kubelet configuration.
type KubeletConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +required
	Spec KubeletConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status KubeletConfigStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// KubeletConfigSpec defines the desired state of KubeletConfig
type KubeletConfigSpec struct {
	MachineConfigPoolSelector *metav1.LabelSelector `json:"machineConfigPoolSelector,omitempty" protobuf:"bytes,1,opt,name=machineConfigPoolSelector"`
	KubeletConfig             *runtime.RawExtension `json:"kubeletConfig,omitempty" protobuf:"bytes,2,opt,name=kubeletConfig"`
}

// KubeletConfigStatus defines the observed state of a KubeletConfig
type KubeletConfigStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []KubeletConfigCondition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
}

// KubeletConfigCondition defines the state of the KubeletConfig
type KubeletConfigCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	Type KubeletConfigStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=KubeletConfigStatusConditionType"`

	// status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the time of the last update to the current status object.
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are PascalCase
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// KubeletConfigStatusConditionType is the state of the operator's reconciliation functionality.
type KubeletConfigStatusConditionType string

const (
	// KubeletConfigSuccess designates a successful application of a KubeletConfig CR.
	KubeletConfigSuccess KubeletConfigStatusConditionType = "Success"

	// KubeletConfigFailure designates a failure applying a KubeletConfig CR.
	KubeletConfigFailure KubeletConfigStatusConditionType = "Failure"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeletConfigList is a list of KubeletConfig resources
type KubeletConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []KubeletConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContainerRuntimeConfig describes a customized Container Runtime configuration.
type ContainerRuntimeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +required
	Spec ContainerRuntimeConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ContainerRuntimeConfigStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// ContainerRuntimeConfigSpec defines the desired state of ContainerRuntimeConfig
type ContainerRuntimeConfigSpec struct {
	MachineConfigPoolSelector *metav1.LabelSelector          `json:"machineConfigPoolSelector,omitempty" protobuf:"bytes,1,opt,name=machineConfigPoolSelector"`
	ContainerRuntimeConfig    *ContainerRuntimeConfiguration `json:"containerRuntimeConfig,omitempty" protobuf:"bytes,2,opt,name=containerRuntimeConfig"`
}

// ContainerRuntimeConfiguration defines the tuneables of the container runtime
type ContainerRuntimeConfiguration struct {
	// pidsLimit specifies the maximum number of processes allowed in a container
	PidsLimit int64 `json:"pidsLimit,omitempty" protobuf:"varint,1,opt,name=pidsLimit"`

	// logLevel specifies the verbosity of the logs based on the level it is set to.
	// Options are fatal, panic, error, warn, info, and debug.
	LogLevel string `json:"logLevel,omitempty" protobuf:"bytes,2,opt,name=logLevel"`

	// logSizeMax specifies the Maximum size allowed for the container log file.
	// Negative numbers indicate that no size limit is imposed.
	// If it is positive, it must be >= 8192 to match/exceed conmon's read buffer.
	LogSizeMax resource.Quantity `json:"logSizeMax" protobuf:"bytes,3,opt,name=logSizeMax"`

	// overlaySize specifies the maximum size of a container image.
	// This flag can be used to set quota on the size of container images. (default: 10GB)
	OverlaySize resource.Quantity `json:"overlaySize" protobuf:"bytes,4,opt,name=overlaySize"`
}

// ContainerRuntimeConfigStatus defines the observed state of a ContainerRuntimeConfig
type ContainerRuntimeConfigStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []ContainerRuntimeConfigCondition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
}

// ContainerRuntimeConfigCondition defines the state of the ContainerRuntimeConfig
type ContainerRuntimeConfigCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	Type ContainerRuntimeConfigStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ContainerRuntimeConfigStatusConditionType"`

	// status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the time of the last update to the current status object.
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are PascalCase
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// ContainerRuntimeConfigStatusConditionType is the state of the operator's reconciliation functionality.
type ContainerRuntimeConfigStatusConditionType string

const (
	// ContainerRuntimeConfigSuccess designates a successful application of a ContainerRuntimeConfig CR.
	ContainerRuntimeConfigSuccess ContainerRuntimeConfigStatusConditionType = "Success"

	// ContainerRuntimeConfigFailure designates a failure applying a ContainerRuntimeConfig CR.
	ContainerRuntimeConfigFailure ContainerRuntimeConfigStatusConditionType = "Failure"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContainerRuntimeConfigList is a list of ContainerRuntimeConfig resources
type ContainerRuntimeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []ContainerRuntimeConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}
