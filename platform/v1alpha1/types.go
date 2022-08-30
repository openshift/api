package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Package contains fields to configure which OLM package this PlatformOperator will install
type Package struct {
	// name contains the desired OLM-based Operator package name
	// that is defined in an existing CatalogSource resource in the cluster.
	//
	// This configured package will be managed with the cluster's lifecycle. In
	// the current implementation, it will be retrieving this name from a list of
	// supported operators out of the catalogs included with OpenShift.
	// +kubebuilder:validation:Required
	//
	// +kubebuilder:validation:Pattern:=[a-z0-9]([-a-z0-9]*[a-z0-9])?
	// +kubebuilder:validation:MaxLength:=56
	// ---
	// + More restrictions to package names supported is an intentional design
	// + decision that, while limiting to user options, allows code built on these
	// + API's to make more confident assumptions on data structure.
	Name string `json:"name"`
}

// PlatformOperatorSpec defines the desired state of PlatformOperator.
type PlatformOperatorSpec struct {
	// package contains the desired package and its configuration for this
	// PlatformOperator.
	// +kubebuilder:validation:Required
	Package Package `json:"package"`
}

// ActiveBundleDeployment references a BundleDeployment resource.
type ActiveBundleDeployment struct {
	// name is the metadata.name of the referenced BundleDeployment object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// PlatformOperatorStatus defines the observed state of PlatformOperator
type PlatformOperatorStatus struct {
	// conditions represent the latest available observations of a platform operator's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// activeBundleDeployment is the reference to the BundleDeployment resource that's
	// being managed by this PO resource. If this field is not populated in the status
	// then it means the PlatformOperator has either not been installed yet or is
	// failing to install.
	// +optional
	ActiveBundleDeployment ActiveBundleDeployment `json:"activeBundleDeployment,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlatformOperator is the Schema for the PlatformOperators API.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type PlatformOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformOperatorSpec   `json:"spec"`
	Status PlatformOperatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlatformOperatorList contains a list of PlatformOperators
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type PlatformOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlatformOperator `json:"items"`
}
