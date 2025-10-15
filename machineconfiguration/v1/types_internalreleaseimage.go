package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// InternalReleaseImageStatus describes the current state of a InternalReleaseImage.
// +openshift:enable:FeatureGate=NoRegistryClusterOperations
type InternalReleaseImageStatus struct {
	// conditions represent the observations of an internal release image current state.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=256
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

// +openshift:enable:FeatureGate=NoRegistryClusterOperations
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
	Image ImageDigestFormat `json:"image,omitempty"`
}
