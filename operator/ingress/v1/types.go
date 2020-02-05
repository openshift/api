package v1

import (
	configv1 "github.com/openshift/api/config/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DNSRecord is a DNS record managed in the zones defined by
// dns.config.openshift.io/cluster .spec.publicZone and .spec.privateZone.
//
// Cluster admin manipulation of this resource is not supported. This resource
// is only for internal communication of OpenShift operators.
type DNSRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the dnsecord.
	Spec DNSRecordSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the dnsRecord.
	Status DNSRecordStatus `json:"status,omitempty"`
}

// DNSRecordSpec contains the details of a DNS record.
type DNSRecordSpec struct {
	// dnsName is the hostname of the DNS record
	//
	// +kubebuilder:validation:Required
	// +required
	DNSName string `json:"dnsName,omitempty"`
	// targets are record targets.
	//
	// +kubebuilder:validation:Required
	// +required
	Targets []string `json:"targets,omitempty"`
	// recordType is the DNS record type. For example, "A" or "CNAME".
	//
	// +kubebuilder:validation:Required
	// +required
	RecordType DNSRecordType `json:"recordType,omitempty"`
	// recordTTL is the record TTL in seconds. The default is 30.
	//
	// +kubebuilder:validation:Optional
	// +optional
	RecordTTL int64 `json:"recordTTL,omitempty"`
}

// DNSRecordStatus is the most recently observed status of each record.
type DNSRecordStatus struct {
	// zones are the status of the record in each zone.
	//
	// +kubebuilder:validation:Optional
	// +optional
	Zones []DNSZoneStatus `json:"zones,omitempty"`
}

// DNSZoneStatus is the status of a record within a specific zone.
type DNSZoneStatus struct {
	// dnsZone is the zone where the record is published.
	//
	// +kubebuilder:validation:Required
	// +required
	DNSZone configv1.DNSZone `json:"dnsZone"`
	// conditions are any conditions associated with the record in the zone.
	//
	// If publishing the record fails, the "Failed" condition will be set with a
	// reason and message describing the cause of the failure.
	//
	// +kubebuilder:validation:Optional
	// +optional
	Conditions []DNSZoneCondition `json:"conditions,omitempty"`
}

var (
	// Failed means the record is not available within a zone.
	DNSRecordFailedConditionType = "Failed"
)

// DNSZoneCondition is just the standard condition fields.
type DNSZoneCondition struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

type DNSRecordType string

const (
	// CNAME is an RFC 1035 CNAME record.
	CNAMERecordType string = "CNAME"

	// CNAME is an RFC 1035 A record.
	ARecordType string = "A"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// DNSRecordList contains a list of dnsrecords.
type DNSRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSRecord `json:"items"`
}
