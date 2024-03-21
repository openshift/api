package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineOSImage describes an image build for a MachineConfigPool
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type MachineOSImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the machineosimage spec
	// +kubebuilder:validation:Required
	Spec MachineOSImageSpec `json:"spec"`

	// status describes the machineosimage status
	// +optional
	Status MachineOSImageStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineOSImageList describes all of the OS Images on the system
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type MachineOSImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MachineOSImage `json:"items"`
}

// MachineOSImageSpec contains user-configurable options and other information about an OS image.
type MachineOSImageSpec struct {
	// baseImage is the base OS image that this object was built from
	// +kubebuilder:validation:Required
	BaseImage string `json:"baseImage"`
	// imagePullSpec contains the final pull spec of the built image
	// +kubebuilder:validation:Required
	ImagePullSpec string `json:"imagePullSpec"`
	// machineConfigPool contains the MCP that the image has been rolled out to
	// +kubebuilder:validation:Required
	MachineConfigPool MachineOSObjectReference `json:"machineConfigPool"`
	// machineConfig is a reference to the MC that was used to build this image.
	// +optional
	RenderedMachineConfig MachineOSObjectReference `json:"machineConfig"`
}

// MachineOSImageStatus contains state ralted information about an image.
type MachineOSImageStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +kubebuilder:validation:Required
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// conditions contains status conditions related to the image.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	// age is how long the image has existed
	// +kubebuilder:validation:Required
	Age metav1.Duration `json:"age"`
	// rolloutStatus represents how the image is progressing in its rollout after being built
	// Valid types are: RolledOut, RolloutPending, RolloutFailed
	// +kubebuilder:validation:Required
	RolloutStatus ImageRolloutStatus `json:"rolloutStatus"`
	// imageUsage describes how recently this image has been used. Once this gets to Rotten, the image can be garbage collected
	// Valid types are: InUse, Stale, Rotten
	// +kubebuilder:validation:Required
	ImageUsage MachineOSImageState `json:"imageUsage"`
}

type ImageRolloutStatus string

const (
	ImageRolledOut      ImageRolloutStatus = "RolledOut"
	ImageRolloutPending ImageRolloutStatus = "RolloutPending"
	ImageRolloutFailed  ImageRolloutStatus = "RolloutFailed"
)

// MachineOSImageState describes the lifecycle of an image
type MachineOSImageState struct {
	// usage describes which point in its lifecycle an image is at
	// Valid types are: InUse, Stale, Rotten
	// +kubebuilder:validation:Required
	Usage MachineOSImageUsage `json:"usage"`
}

// MachineOSImageUsage describes the condition types used to describe the lifecycle of the image
type MachineOSImageUsage string

const (
	MachineOSImageInUse  MachineOSImageUsage = "InUse"
	MachineOSImageStale  MachineOSImageUsage = "Stale"
	MachineOSImageRotten MachineOSImageUsage = "Rotten"
)

type MachineOSObjectReference struct {
	// name is the name of the referenced object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}
