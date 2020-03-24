package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeControllerManager provides information to configure an operator to manage kube-controller-manager.
type KubeControllerManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec is the specification of the desired behavior of the Kubernetes Controller Manager
	// +kubebuilder:validation:Required
	// +required
	Spec KubeControllerManagerSpec `json:"spec"`

	// status is the most recently observed status of the Kubernetes Controller Manager
	// +optional
	Status KubeControllerManagerStatus `json:"status"`
}

type KubeControllerManagerSpec struct {
	StaticPodOperatorSpec `json:",inline"`

	// enableDeprecatedAndRemovedServiceCAKeyUntilNextRelease_ThisMakesClusterImpossibleToUpgrade
	// enables service ca injection into all legacy service account token secrets. Defaults to
	// false. If set to true, will make it impossible to upgrade the cluster.
	// +optional
	EnableDeprecatedAndRemovedServiceCAKeyUntilNextRelease_ThisMakesClusterImpossibleToUpgrade bool `json:"enableDeprecatedAndRemovedServiceCAKeyUntilNextRelease_ThisMakesClusterImpossibleToUpgrade"`
}

type KubeControllerManagerStatus struct {
	StaticPodOperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeControllerManagerList is a collection of items
type KubeControllerManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items contains the items
	Items []KubeControllerManager `json:"items"`
}
