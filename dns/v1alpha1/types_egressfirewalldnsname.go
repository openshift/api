package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +openshift:compatibility-gen:level=4

// EgressFirewallDNSName describes a DNS name used in a EgressFirewall rule. It is TechPreviewNoUpgrade only.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type EgressFirewallDNSName struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the EgressFirewallDNSName.
	// +kubebuilder:validation:Required
	Spec EgressFirewallDNSNameSpec `json:"spec"`
	// status is the most recently observed status of the EgressFirewallDNSName.
	// +optional
	Status EgressFirewallDNSNameStatus `json:"status,omitempty"`
}

// EgressFirewallDNSNameSpec is a desired state description of EgressFirewallDNSName.
type EgressFirewallDNSNameSpec struct {
	// name is the DNS name used in a EgressFirewall rule.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9-]+\.)*[A-Za-z0-9-]+\.$
	Name string `json:"name"`
}

// EgressFirewallDNSNameStatus defines the observed status of EgressFirewallDNSName.
type EgressFirewallDNSNameStatus struct {
	// resolvedNames contains a list of matching DNS names and their corresponding IP addresses along with TTL and last
	// DNS lookup time.
	// +optional
	ResolvedNames []EgressFirewallDNSNameStatusItem `json:"resolvedNames,omitempty"`
}

// EgressFirewallDNSNameStatusItem describes the details of a resolved DNS name.
type EgressFirewallDNSNameStatusItem struct {
	// dnsName is the resolved DNS name matching the name field of EgressFirewallDNSNameSpec.
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9-]+\.)*[A-Za-z0-9-]+\.$
	DNSName string `json:"dnsName"`
	// The IP addresses associated with the DNS name used in a EgressFirewall rule.
	// +listType=set
	IPs []string `json:"ips"`
	// Minimum time-to-live value among all the IP addresses.
	TTL int64 `json:"ttl"`
	// Timestamp when the last DNS lookup was successfully completed.
	LastLookupTime metav1.Time `json:"lastLookupTime"`
	// retryCounter keeps the count of how many times the DNS lookup failed for the dnsName field.
	RetryCounter int `json:"retryCounter"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +openshift:compatibility-gen:level=4

// EgressFirewallDNSNameList contains a list of EgressFirewallDNSNames.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type EgressFirewallDNSNameList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EgressFirewallDNSName `json:"items"`
}
