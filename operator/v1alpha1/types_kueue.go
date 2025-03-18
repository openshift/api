package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kueue is the CRD to represent the Kueue operator
// This CRD defines the configuration that the Kueue
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kueue,scope=Cluster
// +k8s:openapi-gen=true
// +genclient
// +genclient:nonNamespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="Kueue is a singleton, .metadata.name must be 'cluster'"
type Kueue struct {
	metav1.TypeMeta `json:",inline"`
	// metadata for kueue
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec KueueOperandSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status KueueStatus `json:"status,omitempty"`
}

type KueueOperandSpec struct {
	operatorv1.OperatorSpec `json:",inline"`
	// config is the desired configuration
	// for the Kueue operator.
	// +required
	Config KueueConfiguration `json:"config"`
}

type KueueConfiguration struct {
	// integrations are the workloads Kueue will manage
	// Kueue has integrations in the codebase and it also allows
	// for external frameworks
	// Kueue are an important part to specify for the API as Kueue will
	// only manage the workloads that are specfied in this list.
	// This is a required field.
	// +required
	Integrations Integrations `json:"integrations"`
}

// KueueStatus defines the observed state of Kueue
type KueueStatus struct {
	operatorv1.OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KueueList contains a list of Kueue
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type KueueList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata for the list
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// items is a slice of Kueue
	// this is a cluster scoped resource and there can only be 1 Kueue
	// +kubebuilder:validation:MaxItems=1
	// +required
	Items []Kueue `json:"items"`
}

// +kubebuilder:validation:Enum=BatchJob;RayJob;RayCluster;JobSet;MPIJob;PaddleJob;PytorchJob;TFJob;XGBoostJob;AppWrappers;Pod;Deployment;StatefulSet;LeaderWorkerSet
type KueueIntegrations string

const (
	BatchJob        KueueIntegrations = "BatchJob"
	RayJob          KueueIntegrations = "RayJob"
	RayCluster      KueueIntegrations = "RayCluster"
	JobSet          KueueIntegrations = "JobSet"
	MPIJob          KueueIntegrations = "MPIJob"
	PaddeJob        KueueIntegrations = "PaddeJob"
	PyTorchJob      KueueIntegrations = "PyTorchJob"
	TFJob           KueueIntegrations = "TFJob"
	XGBoostJob      KueueIntegrations = "XGBoostJob"
	AppWrappers     KueueIntegrations = "AppWrappers"
	Pod             KueueIntegrations = "Pod"
	Deployment      KueueIntegrations = "Deployment"
	StatefulSet     KueueIntegrations = "StatefulSet"
	LeaderWorkerSet KueueIntegrations = "LeaderWorkerSet"
)

// This is the GVK for an external framework.
// Controller runtime requires this in this format
// for api discoverability.
type ExternalFramework struct {
	// group of externalFramework
	// must be a valid qualified name consisting of a lower-case alphanumeric string,
	// and hyphens of at most 63 characters in length. The name must start and end with an alphanumeric character.
	// The name may be optionally prefixed with a subdomain consisting of lower-case alphanumeric characters,
	// hyphens and periods, of at most 253 characters in length.
	// Each period separated segment within the subdomain must start and end with an alphanumeric character.
	// The optional prefix and the name are separate by a forward slash (/).
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.qualifiedName().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Group string `json:"group"`
	// resource of external framework
	// must be a valid qualified name consisting of a lower-case alphanumeric string,
	// and hyphens of at most 63 characters in length. The name must start and end with an alphanumeric character.
	// The name may be optionally prefixed with a subdomain consisting of lower-case alphanumeric characters,
	// hyphens and periods, of at most 253 characters in length.
	// Each period separated segment within the subdomain must start and end with an alphanumeric character.
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.qualifiedName().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MinLength=1
	// +required
	Resource string `json:"resource"`
	// version is the version of the api
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:MinLength=1
	// +required
	Version string `json:"version"`
}

type Integrations struct {
	// frameworks are a list of names to be enabled.
	// This is required and must have at least one element.
	// The frameworks are jobs that Kueue will manage.
	// KueueIntegrations is a list of frameworks that Kueue has support for.
	// The allowed values are BatchJob;RayJob;RayCluster;JobSet;MPIJob;PaddleJob;PytorchJob;TFJob;XGBoostJob;AppWrappers;Pod;Deployment;StatefulSet;LeaderWorkerSet.
	// +kubebuilder:validation:MaxItems=14
	// +kubebuilder:validation:MinItems=1
	// kubebuilder:validation:UniqueItems=true
	// +listType=atomic
	// +required
	Frameworks []KueueIntegrations `json:"frameworks"`
	// externalFrameworks are a list of GroupVersionResources
	// that are managed for Kueue by external controllers;
	// These are optional and should only be used if you have an external controller
	// that integrates with Kueue.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=32
	// +optional
	ExternalFrameworks []ExternalFramework `json:"externalFrameworks,omitempty"`

	// labelKeysToCopy are a list of label keys that are copied once a workload is created
	// these keys are persisted to the internal Kueue workload object.
	// otherwise only the Kueue labels will be copied.
	// +kubebuilder:validation:MaxItems=64
	// +listType=atomic
	// +optional
	LabelKeysToCopy []LabelKeys `json:"labelKeysToCopy,omitempty"`
}

type LabelKeys struct {
	// key is the label key
	// A label key must be a valid qualified name consisting of a lower-case alphanumeric string,
	// and hyphens of at most 63 characters in length. The name must start and end with an alphanumeric character.
	// The name may be optionally prefixed with a subdomain consisting of lower-case alphanumeric characters,
	// hyphens and periods, of at most 253 characters in length.
	// Each period separated segment within the subdomain must start and end with an alphanumeric character.
	// The optional prefix and the name are separate by a forward slash (/).
	// +kubebuilder:validation:MaxLength=317
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.qualifiedName().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +optional
	Key string `json:"key,omitempty"`
}
