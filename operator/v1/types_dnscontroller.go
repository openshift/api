package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSController manages the CoreDNS component to provide a name resolution service
// for pods and services in the cluster.
//
// This supports the DNS-based service discovery specification:
// https://github.com/kubernetes/dns/blob/master/docs/specification.md
//
// More details: https://kubernetes.io/docs/tasks/administer-cluster/coredns
type DNSController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the DNSController.
	Spec DNSControllerSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the DNSController.
	Status DNSControllerStatus `json:"status,omitempty"`
}

// DNSControllerSpec is the specification of the desired behavior of the DNSController.
type DNSControllerSpec struct {
}

const (
	// Available indicates the DNS controller daemonset is available.
	DNSControllerAvailableConditionType = "Available"
)

// DNSControllerStatus defines the observed status of the DNSController.
type DNSControllerStatus struct {
	// clusterIP is the service IP through which this DNS is made available.
	//
	// In the case of the default DNSController, this will be a well known IP that is used
	// as the default nameserver for pods that are using the default ClusterFirst DNS policy.
	//
	// In general, this IP can be specified in a pod's spec.dnsConfig.nameservers list
	// or used explicitly when performing name resolution from within the cluster.
	// Example: dig foo.com @<service IP>
	//
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	ClusterIP string `json:"clusterIP"`

	// clusterDomain is the local cluster DNS domain suffix for DNS controller services.
	// This will be a subdomain as defined in RFC 1034,
	// section 3.5: https://tools.ietf.org/html/rfc1034#section-3.5
	// Example: "cluster.local"
	//
	// More info: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service
	ClusterDomain string `json:"clusterDomain"`

	// conditions provide information about the state of the DNS controller on the cluster.
	//
	// These are the supported DNS controller conditions:
	//
	//   * Available
	//   - True if the following conditions are met:
	//     * DNS controller daemonset is available.
	//   - False if any of those conditions are unsatisfied.
	//
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []OperatorCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSControllerList contains a list of DNSControllers.
type DNSControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DNSController `json:"items"`
}
