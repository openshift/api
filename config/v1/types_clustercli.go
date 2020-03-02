package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCLI is the Schema for the cluster cli API
type ClusterCLI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	// +required
	Spec ClusterCLISpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ClusterCLIStatus `json:"status,omitempty"`
}

// ClusterCLISpec defines the desired state of ClusterCLI
type ClusterCLISpec struct {
	// Description of ClusterCLI
	Description string `json:"description"`
	// DisplayName for CLI
	DisplayName string `json:"displayName"`
	// Image is the cli image that contains the cli artifacts
	Image string `json:"image"`
	// Mapping defines extract targets for ClusterCLIs
	Mapping []ClusterCLIMapping `json:"mapping,omitempty"`
}

// ClusterCLIStatus defines the observed state of ClusterCLI
type ClusterCLIStatus struct {
	Version string `json:"version,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCLIList contains a list of ClusterCLI
type ClusterCLIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCLI `json:"items"`
}

// ClusterCLIMapping holds mapping information from cli images for extracting clis
type ClusterCLIMapping struct {
	// OS is GOOS
	OS string `json:"os,omitempty"`
	// Arch is GOARCH
	Arch string `json:"arch,omitempty"`
	// From is the directory or file in the image to extract
	From string `json:"from,omitempty"`
}
