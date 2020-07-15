package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Etcd provides information to configure an operator to manage kube-apiserver.
type Etcd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// +kubebuilder:validation:Required
	// +required
	Spec EtcdSpec `json:"spec"`
	// +optional
	Status EtcdStatus `json:"status"`
}

type EtcdSpec struct {
	StaticPodOperatorSpec `json:",inline"`
}

type EtcdStatus struct {
	StaticPodOperatorStatus `json:",inline"`

	// leader describes the etcd member which is currently the leader. This field
	// is a hint and shouldn't be considered authoritative.
	Leader LeaderStatus `json:"leader"`
}

// leaderStatus describes the etcd leader member details.
type LeaderStatus struct {
	// name is the etcd leader member name, if available.
	Name string `json:"name,omitempty"`
	// node is the etcd leader member node, if available.
	Node string `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeAPISOperatorConfigList is a collection of items
type EtcdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items contains the items
	Items []Etcd `json:"items"`
}
