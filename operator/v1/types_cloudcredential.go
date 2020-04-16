package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudCredential provides a means to configure an operator to manage the cloud-credential-operator. `cluster` is the canonical name.
type CloudCredential struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec CloudCredentialSpec `json:"spec"`

	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status CloudCredentialStatus `json:"status"`
}

// CloudCredentialSpec is the specification of the desired behavior of the CloudCredential operator.
type CloudCredentialSpec struct {
	OperatorSpec `json:",inline"`
}

// CloudCredentialStatus defines the observed status of the CloudCredential operator.
type CloudCredentialStatus struct {
	OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// CloudCredentialList contains a list of CloudCredentials.
type CloudCredentialList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudCredential `json:"items"`
}
