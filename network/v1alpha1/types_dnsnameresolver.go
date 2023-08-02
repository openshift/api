package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +openshift:compatibility-gen:level=4

// DNSNameResolver stores the DNS name resolution information of a DNS name. It is TechPreviewNoUpgrade only.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type DNSNameResolver struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the DNSNameResolver.
	// +kubebuilder:validation:Required
	Spec DNSNameResolverSpec `json:"spec"`
	// status is the most recently observed status of the DNSNameResolver.
	// +optional
	Status DNSNameResolverStatus `json:"status,omitempty"`
}

// DNSNameResolverSpec is a desired state description of DNSNameResolver.
type DNSNameResolverSpec struct {
	// name is the DNS name for which the DNS name resolution information will be stored.
	// For a regular DNS name, only the DNS name resolution information of the regular DNS
	// name will be stored. For a wildcard DNS name, the DNS name resolution information
	// of all the DNS names, that matches the wildcard DNS name, will be stored.
	// For a wildcard DNS name, the '*' will match only one label. Additionally, only a single
	// '*' can be used at the beginning of the wildcard DNS name. For example, '*.example.com.'
	// will match 'sub1.example.com.' but won't match 'sub2.sub1.example.com.'
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9-]+\.)*[A-Za-z0-9-]+\.$
	// +kubebuilder:validation:MaxLength=254
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec.name is immutable"
	Name string `json:"name"`
}

// DNSNameResolverStatus defines the observed status of DNSNameResolver.
type DNSNameResolverStatus struct {
	// resolvedNames contains a list of matching DNS names and their corresponding IP addresses
	// along with TTL and last DNS lookup time.
	// +listType=map
	// +listMapKey=dnsName
	// +patchMergeKey=dnsName
	// +patchStrategy=merge
	// +optional
	ResolvedNames []DNSNameResolverStatusItem `json:"resolvedNames,omitempty" patchStrategy:"merge" patchMergeKey:"dnsName"`
}

// +kubebuilder:validation:Pattern=`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/([0-9]|[1-2][0-9]|3[0-2])$|^s*((([0-9A-Fa-f]{1,4}:){7}(:|([0-9A-Fa-f]{1,4})))|(([0-9A-Fa-f]{1,4}:){6}:([0-9A-Fa-f]{1,4})?)|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){0,1}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){0,2}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){0,3}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){0,4}):([0-9A-Fa-f]{1,4})?))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){0,5}):([0-9A-Fa-f]{1,4})?))|(:(:|((:[0-9A-Fa-f]{1,4}){1,7}))))(%.+)?s*/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])$`
// IPAddressStr is used for validation of an IP address.
type IPAddressStr string

// DNSNameResolverStatusItem describes the details of a resolved DNS name.
type DNSNameResolverStatusItem struct {
	// conditions provide information about the state of the DNS name.
	// Known .status.conditions.type is: "Degraded"
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// dnsName is the resolved DNS name matching the name field of DNSNameResolverSpec.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9-]+\.)*[A-Za-z0-9-]+\.$
	// +kubebuilder:validation:MaxLength=254
	DNSName string `json:"dnsName"`
	// info gives the list of associated IP addresses and the corresponding TTL and last
	// lookup time for the dnsName.
	// +kubebuilder:validation:Required
	// +listType=map
	// +listMapKey=ip
	Info []DNSNameResolverInfo `json:"info"`
	// resolutionFailures keeps the count of how many consecutive times the DNS resolution failed
	// for the dnsName. If the DNS resolution succeeds then the field will be set to zero. Upon
	// every failure, the value of the field will be incremented by one. Upon reaching the value
	// of 5, the details about the DNS name will be removed.
	ResolutionFailures int `json:"resolutionFailures,omitempty"`
}

type DNSNameResolverInfo struct {
	// ip is an IP address associated with the dnsName. The validity of the IP address expires after
	// lastLookupTime + ttlSeconds. To refresh the information a DNS lookup will be performed on the
	// expiration of the IP address's validity. If the information is not refreshed then it will be
	// removed after a grace period of 1 second after the expiration of the IP address's validity.
	// +kubebuilder:validation:Required
	IP IPAddressStr `json:"ip"`
	// ttlSeconds is the minimum time-to-live value among all the IP addresses.
	// +kubebuilder:validation:Required
	TTLSeconds int32 `json:"ttlSeconds"`
	// lastLookupTime is the timestamp when the last DNS lookup was completed.
	// +kubebuilder:validation:Required
	LastLookupTime *metav1.Time `json:"lastLookupTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +openshift:compatibility-gen:level=4

// DNSNameResolverList contains a list of DNSNameResolvers.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type DNSNameResolverList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DNSNameResolver `json:"items"`
}
