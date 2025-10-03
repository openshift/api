package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	machineosconfig "github.com/openshift/api/machineconfiguration/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=internalreleaseimages,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:file-pattern=cvoRunLevel=0000_80,operatorName=machine-config,operatorOrdering=01
// +openshift:enable:FeatureGate=NoRegistryClusterOperations
// +kubebuilder:metadata:labels=openshift.io/operator-managed=

// InternalReleaseImage is used to keep track and manage a set
// of release bundles (OCP and OLM operators images) that are stored
// into the control planes nodes.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type InternalReleaseImage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the configuration of this internal release image.
	// +required
	Spec InternalReleaseImageSpec `json:"spec,omitzero"`

	// status describes the last observed state of this internal release image.
	// +optional
	Status InternalReleaseImageStatus `json:"status"`
}

// InternalReleaseImageStatus describes the current state of a InternalReleaseImage.
type InternalReleaseImageStatus struct {
	// conditions represent the observations of an internal release image current state.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=20
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// availableReleases is a list of release bundle identifiers currently detected
	// from the attached ISO.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	// +optional
	AvailableReleases []InternalReleaseImageRef `json:"availableReleases,omitempty"`

	// releases is a list of the currently managed release bundles.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	// +optional
	Releases []InternalReleaseImageRef `json:"releases,omitempty"`
}

// InternalReleaseImageSpec defines the desired state of a InternalReleaseImage.
type InternalReleaseImageSpec struct {
	// releases is a list of release bundle identifiers that the user wants to
	// add/remove to/from the control plane nodes.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=20
	// +listType=map
	// +listMapKey=name
	Releases []InternalReleaseImageRef `json:"releases,omitempty"`
}

type InternalReleaseImageRef struct {
	// name indicates the desired release bundle identifier.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	Name string `json:"name,omitempty"`

	// image is an OCP release imaged referenced by digest.
	// The format of the image pull spec is: host[:port][/namespace]/name@sha256:<digest>,
	// where the digest must be 64 characters long, and consist only of lowercase hexadecimal characters, a-f and 0-9.
	// The length of the whole spec must be between 1 to 447 characters.
	// +optional
	Image machineosconfig.ImageDigestFormat `json:"image,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InternalReleaseImageList is a list of InternalReleaseImage resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type InternalReleaseImageList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	Items []InternalReleaseImage `json:"items"`
}
