package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftBootstrapConfig is repsonsible for providing ignition configuration to bootstrp nodes
// on an OpenShift cluster managed by Cluster API.
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type OpenShiftBootstrapConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// TODO: Since the configuration for bootstrapping is owned by the cluster,
	// and the bootstrap provider gets the Machine for bootstrapping via the owner reference,
	// it doesn't currently need a spec.
	// Should cofirm this is true during initial POC.

	// status is the observed state of the OpenShiftBootstrapProvider.
	// +optional
	Status OpenShiftBootstrapConfigStatus `json:"status,omitempty"`
}

// OpenShiftBootstrapConfigStatus contains status related to the OpenShiftBootstrapConfig state.
type OpenShiftBootstrapConfigStatus struct {
	// conditions represents the observations of the OpenShiftBootstrapConfig's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ready denotes whether or not the bootstrap data has been created.
	// +optional
	// + ---
	// + This field is required as part of the Cluster API bootstrap provider API contract.
	Ready bool `json:"ready,omitempty"`

	// dataSecretName is the name of the user data secret that should be used to bootstrap
	// the owning Machine's Node when it first boots.
	// + ---
	// + This field is required as part of the Cluster API bootstrap provider API contract.
	// +kubebuilder:validation:Pattern="[a-z0-9]([-.a-z0-9]{,251}[a-z0-9])?"
	// +kubebuilder:validation:MaxLength=253
	// +optional
	DataSecretName string `json:"dataSecretName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftBootstrapConfigList contains a list of OpenShiftBootstrapConfig
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type OpenShiftBootstrapConfigList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items contains a list of OpenShiftBootstrapConfigs.
	Items []OpenShiftBootstrapConfig `json:"items"`
}
