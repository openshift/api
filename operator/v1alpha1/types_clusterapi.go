package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterapis,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2564
// +openshift:file-pattern=cvoRunLevel=0000_30,operatorName=cluster-api,operatorOrdering=01
// +openshift:enable:FeatureGate=ClusterAPIMachineManagement

// ClusterAPI provides configuration for the capi-operator.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.spec) || !has(oldSelf.spec.unmanagedCustomResourceDefinitions) || has(self.spec.unmanagedCustomResourceDefinitions)",message="unmanagedCustomResourceDefinitions cannot be unset once set"
type ClusterAPI struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the capi-operator.
	// +optional
	Spec *ClusterAPISpec `json:"spec,omitempty"`

	// status defines the observed status of the capi-operator.
	// +optional
	Status ClusterAPIStatus `json:"status,omitzero"`
}

// ClusterAPISpec defines the desired configuration of the capi-operator.
type ClusterAPISpec struct {
	// unmanagedCustomResourceDefinitions is a list of ClusterResourceDefinition (CRD)
	// names that should not be managed by the capi-operator installer
	// controller. This allows external actors to own specific CRDs while
	// capi-operator manages others.
	//
	// Each CRD name must be a valid DNS-1123 subdomain consisting of lowercase
	// alphanumeric characters, '-' or '.', and must start and end with an
	// alphanumeric character, with a maximum length of 253 characters.
	// Example: "clusters.cluster.x-k8s.io"
	//
	// Items cannot be removed from this list once added.
	//
	// The maximum number of unmanagedCustomResourceDefinitions is 128.
	//
	// +optional
	// +listType=set
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:XValidation:rule="oldSelf.all(item, item in self)",message="items cannot be removed from unmanagedCustomResourceDefinitions list"
	// +kubebuilder:validation:items:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:items:MaxLength=253
	UnmanagedCustomResourceDefinitions []string `json:"unmanagedCustomResourceDefinitions,omitempty"`
}

// ClusterAPIStatus describes the current state of the capi-operator.
// +kubebuilder:validation:MinProperties=1
type ClusterAPIStatus struct {
	// targetConfigMaps is a list of ConfigMap names that the staging controller
	// has validated and approved for reconciliation. The installer controller
	// will reconcile these ConfigMaps.
	//
	// Each ConfigMap name must be a valid DNS-1123 label consisting of lowercase
	// alphanumeric characters or hyphens, starting and ending with an alphanumeric
	// character, with a maximum length of 63 characters.
	//
	// This field is owned by the staging controller and is updated atomically to a
	// consistent set of transport ConfigMaps that have passed validation checks.
	//
	// The maximum number of targetConfigMaps is 128.
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MaxLength=63
	// +kubebuilder:validation:items:XValidation:rule="!format.dns1123Label().validate(self).hasValue()",message="each ConfigMap name must be a valid DNS1123 label consisting of lowercase alphanumeric characters or hyphens, starting and ending with an alphanumeric character"
	TargetConfigMaps []string `json:"targetConfigMaps,omitempty"`

	// activeConfigMaps is a list of ConfigMap names that the installer
	// controller has successfully reconciled. This represents the currently
	// deployed CAPI provider components.
	//
	// Each ConfigMap name must be a valid DNS-1123 label consisting of lowercase
	// alphanumeric characters or hyphens, starting and ending with an alphanumeric
	// character, with a maximum length of 63 characters.
	//
	// This field is owned by the installer controller and is updated atomically after
	// a successful reconciliation.
	//
	// The maximum number of activeConfigMaps is 128.
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MaxLength=63
	// +kubebuilder:validation:items:XValidation:rule="!format.dns1123Label().validate(self).hasValue()",message="each ConfigMap name must be a valid DNS1123 label consisting of lowercase alphanumeric characters or hyphens, starting and ending with an alphanumeric character"
	ActiveConfigMaps []string `json:"activeConfigMaps,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAPIList contains a list of ClusterAPI configurations
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ClusterAPIList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	// items contains the items
	Items []ClusterAPI `json:"items"`
}
