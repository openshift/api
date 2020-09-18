package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeDescheduler provides information to configure an operator to manage kube-descheduler.
type KubeDescheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec KubeDeschedulerSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status KubeDeschedulerStatus `json:"status"`
}

// KubeDeschedulerSpec defines the desired state of KubeDescheduler
type KubeDeschedulerSpec struct {
	OperatorSpec `json:",inline"`
}

// KubeDeschedulerStatus defines the observed state of KubeDescheduler
type KubeDeschedulerStatus struct {
	OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeDeschedulerList contains a list of KubeDescheduler
type KubeDeschedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeDescheduler `json:"items"`
}
