package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +openshift:compatibility-gen:level=4

// DNSNameResolver stores the DNS name resolution information of a DNS name. It can be enabled by the TechPreviewNoUpgrade feature set.
// It can also be enabled by the feature gate DNSNameResolver when using CustomNoUpgrade feature set.
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
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9]([-A-Za-z0-9]*[A-Za-z0-9])?\.)*[A-Za-z0-9]([-A-Za-z0-9]*[A-Za-z0-9])?\.$
	// +kubebuilder:validation:MaxLength=254
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec.name is immutable"
	Name string `json:"name"`
}

// DNSNameResolverStatus defines the observed status of DNSNameResolver.
type DNSNameResolverStatus struct {
	// resolvedNames contains a list of matching DNS names and their corresponding IP addresses
	// along with their TTL and last DNS lookup times.
	// +listType=map
	// +listMapKey=dnsName
	// +patchMergeKey=dnsName
	// +patchStrategy=merge
	// +optional
	ResolvedNames []DNSNameResolverStatusItem `json:"resolvedNames,omitempty" patchStrategy:"merge" patchMergeKey:"dnsName"`
}

// DNSNameResolverStatusItem describes the details of a resolved DNS name.
type DNSNameResolverStatusItem struct {
	// conditions provide information about the state of the DNS name.
	// Known .status.conditions.type is: "Degraded"
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// dnsName is the resolved DNS name matching the name field of DNSNameResolverSpec. This field can
	// store both regular and wildcard DNS names which match the spec.name field. When the spec.name
	// field contains a regular DNS name, this field will store the same regular DNS name after it is
	// successfully resolved. When the spec.name field contains a wildcard DNS name, each resolvedName.dnsName
	// will store the regular DNS names which match the wildcard DNS name and have been successfully resolved.
	// If the wildcard DNS name can also be successfully resolved, then this field will store the wildcard
	// DNS name as well.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^(\*\.)?([A-Za-z0-9]([-A-Za-z0-9]*[A-Za-z0-9])?\.)*[A-Za-z0-9]([-A-Za-z0-9]*[A-Za-z0-9])?\.$
	// +kubebuilder:validation:MaxLength=254
	DNSName string `json:"dnsName"`
	// resolvedAddresses gives the list of associated IP addresses and their corresponding TTLs and last
	// lookup times for the dnsName.
	// +kubebuilder:validation:Required
	// +listType=map
	// +listMapKey=ip
	ResolvedAddresses []DNSNameResolverInfo `json:"resolvedAddresses"`
	// resolutionFailures keeps the count of how many consecutive times the DNS resolution failed
	// for the dnsName. If the DNS resolution succeeds then the field will be set to zero. Upon
	// every failure, the value of the field will be incremented by one. The details about the DNS
	// name will be removed, if the value of resolutionFailures reaches 5 and the TTL of all the
	// associated IP addresses have expired.
	ResolutionFailures int32 `json:"resolutionFailures,omitempty"`
}

type DNSNameResolverInfo struct {
	// ip is an IP address associated with the dnsName. The validity of the IP address expires after
	// lastLookupTime + ttlSeconds. To refresh the information a DNS lookup will be performed on the
	// expiration of the IP address's validity. If the information is not refreshed then it will be
	// removed with a grace period after the expiration of the IP address's validity.
	// +kubebuilder:validation:Required
	IP string `json:"ip"`
	// ttlSeconds is the time-to-live value of the IP address. The validity of the IP address expires after
	// lastLookupTime + ttlSeconds. On a successful DNS lookup the value of this field will be updated with
	// the current time-to-live value. If the information is not refreshed then it will be removed with a
	// grace period after the expiration of the IP address's validity.
	// +kubebuilder:validation:Required
	TTLSeconds int32 `json:"ttlSeconds"`
	// lastLookupTime is the timestamp when the last DNS lookup was completed successfully. The validity of
	// the IP address expires after lastLookupTime + ttlSeconds. The value of this field will be updated to
	// the current time on a successful DNS lookup. If the information is not refreshed then it will be
	// removed with a grace period after the expiration of the IP address's validity.
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
