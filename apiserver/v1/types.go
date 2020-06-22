package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CriticalResource prevents the .spec.provider from being deleted until the .spec.criteria are satisfied.
// Because the ability to block deletion of a namespace scoped resource implies the ability to prevent namespaces
// from being cleaned up, these are namespace scoped, but they have a higher degree of privilege required for mutation
// than namespace admin.
// TODO in a future release we can consider making the permissions broader, but because criteria span namespaces,
//  it seems unlikely that we can broaden the permissions.
type CriticalResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CriticalResourceSpec   `json:"spec"`
	Status CriticalResourceStatus `json:"status"`
}

type CriticalResourceSpec struct {
	// provider indicates the resource which provides critical functionality.  It must be in the same namespace as the
	// critical resource.
	// TODO if we need to protect cluster scoped resources, we will need to create a resource type.  We cannot reasonably
	//  block deletion of a cluster scoped resource from within a namespace.
	Provider CriticalResourceProvider   `json:"provider"`
	Criteria []CriticalResourceCriteria `json:"criteria"`
}

type GroupResource struct {
	// TODO provide validation
	Group string `json:"group"`
	// TODO provide validation
	Resource string `json:"resource"`
}

type CriticalResourceProvider struct {
	// only allow deployments.apps to start
	GroupResource `json:",inline"`
	// TODO provide validation
	Name string `json:"name"`
}

type CriticalResourceCriteriaType string

var (
	FinalizerType        CriticalResourceCriteriaType = "Finalizer"
	SpecificResourceType CriticalResourceCriteriaType = "SpecificResource"
)

type CriticalResourceCriteria struct {
	Type             CriticalResourceCriteriaType              `json:"type"`
	Finalizer        *FinalizerCriticalResourceCriteria        `json:"finalizer"`
	SpecificResource *SpecificResourceCriticalResourceCriteria `json:"specificResource"`
}

type FinalizerCriticalResourceCriteria struct {
	GroupResource `json:",inline"`
	FinalizerName string `json:"finalizerName"`
}

type SpecificResourceCriticalResourceCriteria struct {
	GroupResource `json:",inline"`
	// namespace can be empty for cluster scoped resources.
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type CriticalResourceStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthenticationList is a collection of items
type CriticalResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CriticalResource `json:"items"`
}
