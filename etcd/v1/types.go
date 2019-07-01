package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterMemberRequest describes a request to join the cluster by a potential etcd cluster member.
type ClusterMemberRequest struct {
	// Standard object's metadata.
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// The membership request itself.
	// +optional
	Spec ClusterMemberRequestSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Derived information about the request.
	// +optional
	Status ClusterMemberRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ClusterMemberRequestSpec defines an individual etcd cluster member request.
type ClusterMemberRequestSpec struct {
	// Name of the etcd peer member requesting membership.
	Name string `json:"Name" protobuf:"bytes,1,opt,name=name"`
	// PeerlURLs is comma seperated string including the members PeerURLs.
	// +optional
	PeerURLs string `json:"PeerURLs" protobuf:"bytes,2,opt,name=peerUrls"`
}

type ClusterMemberRequestStatus struct {
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []ClusterMemberRequestCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`

	// If request was approved, the controller will place the populated config here.
	// +optional
	Config ClusterMemberConfig `json:"config,omitempty" protobuf:"bytes,2,opt,name=config"`
}

type RequestConditionType string

// These are the possible conditions for the etcd membership request.
const (
	ClusterMemberApproved  RequestConditionType = "Approved"
	ClusterMemberDelivered RequestConditionType = "Delivered"
	ClusterMemberDenied    RequestConditionType = "Denied"
	ClusterMemberPending   RequestConditionType = "Pending"
)

type ClusterMemberRequestCondition struct {
	// request approval state, currently Approved, Delivered, Denied or Pending.
	Type RequestConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=RequestConditionType"`
	// brief reason for the request state
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// human readable message with details about the request state
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// timestamp for the last update to this condition
	// +optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,4,opt,name=lastUpdateTime"`
}

type ClusterMemberConfig struct {
	// Initial cluster configuration for bootstrapping. The value will be a comma delimted string
	// containing a list of etcd peers in the format <name>=<address>.
	InitialCluster string `json:"initialCluster,omitempty" protobuf:"bytes,2,opt,name=initialCluster"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterMemberRequestList is a collection of ClusterMemberRequests
type ClusterMemberRequestList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is the list of ClusterMemberRequests
	Items []ClusterMemberRequest `json:"items" protobuf:"bytes,2,rep,name=items"`
}
