package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=http01challengeproxies,scope=Namespaced
// +openshift:api-approved.openshift.io=true
// +openshift:file-pattern=cvoRunLevel=0000_70,operatorName=network,operatorOrdering=00
// +openshift:compatibility-gen:level=4
// +openshift:enable:FeatureGate=HTTP01ChallengeProxy

// HTTP01ChallengeProxy is the schema for the HTTP01ChallengeProxy API.
type HTTP01ChallengeProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the HTTP01ChallengeProxy.
	// +required
	Spec HTTP01ChallengeProxySpec `json:"spec"`
	// Status is the most recently observed status of the HTTP01ChallengeProxy.
	// +optional
	Status HTTP01ChallengeProxyStatus `json:"status,omitempty"`
}

// HTTP01ChallengeProxySpec is a desired state description of HTTP01ChallengeProxy.
type HTTP01ChallengeProxySpec struct {
}

// HTTP01ChallengeProxyStatus defines the observed status of HTTP01ChallengeProxy.
type HTTP01ChallengeProxyStatus struct {
	// Conditions is a list of conditions and their status.
	// +listType=map
	// +listMapKey=type
	Conditions []HTTP01ChallengeProxyCondition `json:"conditions"`
}

// HTTP01ChallengeProxyCondition describes the state of the HTTP01ChallengeProxy.
type HTTP01ChallengeProxyCondition struct {
	// Type is the type of the condition.
	// +required
	Type string `json:"type"`
	// Status is the status of the condition.
	// +required
	Status metav1.ConditionStatus `json:"status"`
	// LastTransitionTime is the time of the last transition of the condition.
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is the reason for the condition's last transition.
	// +required
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about the last transition.
	// +required
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HTTP01ChallengeProxyList contains a list of HTTP01ChallengeProxy.
type HTTP01ChallengeProxyList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of HTTP01ChallengeProxy.
	// +listType=map
	Items []HTTP01ChallengeProxy `json:"items"`
}
