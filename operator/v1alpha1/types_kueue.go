package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kueue is the CRD to represent the kueue operator
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
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="kueue is a singleton, .metadata.name must be 'cluster'"
type Kueue struct {
	metav1.TypeMeta `json:",inline"`
	// metadata for kueue
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// kueue configuration must not be changed once the object exists
	// to change the configuration, one can delete the object and create a new object.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="values are immutable once set"
	// +required
	Spec KueueOperandSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status KueueStatus `json:"status,omitempty"`
}

type KueueOperandSpec struct {
	operatorv1.OperatorSpec `json:",inline"`
	// config is the desired configuration
	// for the kueue operator.
	// +required
	Config KueueConfiguration `json:"config"`
}

type KueueConfiguration struct {
	// integrations are the workloads Kueue will manage
	// kueue has integrations in the codebase and it also allows
	// for external frameworks
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
	// items is a slice of kueue
	// this is a cluster scoped resource and there can only be 1 kueue
	// +kubebuilder:validation:MaxItems=1
	// +required
	Items []Kueue `json:"items"`
}

// +kubebuilder:validation:Enum=batch/job;ray.io/rayjob;ray.io/raycluster;jobset.x-k8s.io/jobset;kubeflow.org/mpijob;kubeflow.org/paddlejob;kubeflow.org/pytorchjob;kubeflow.org/tfjob;kubeflow.org/xgboostjob;workload.codeflare.dev/appwrapper;pod;deployment;statefulset;leaderworkerset.x-k8s.io/leaderworkerset
type KueueIntegrations string

const (
	BatchJob        KueueIntegrations = "batch/job"
	RayJob          KueueIntegrations = "ray.io/rayjob"
	RayCluster      KueueIntegrations = "ray.io/raycluster"
	JobSet          KueueIntegrations = "jobset.x-k8s.io/jobset"
	MPIJob          KueueIntegrations = "kubeflow.org/mpijob"
	PaddeJob        KueueIntegrations = "kubeflow.org/paddlejob"
	PyTorchJob      KueueIntegrations = "kubeflow.org/pytorchjob"
	TfJob           KueueIntegrations = "kubeflow.org/tfjob"
	XGBoostJob      KueueIntegrations = "kubeflow.org/xgboostjob"
	AppWrappers     KueueIntegrations = "workload.codeflare.dev/appwrapper"
	Pod             KueueIntegrations = "pod"
	Deployment      KueueIntegrations = "deployment"
	Statefulset     KueueIntegrations = "statefulset"
	LeaderWorkerSet KueueIntegrations = "leaderworkerset.x-k8s.io/leaderworkerset"
)

// This is the GVK for an external framework.
// Controller runtime requires this in this format
// for api discoverability.
type ExternalFramework struct {
	// group of externalFramework
	// +kubebuilder:validation:MaxLength=256
	// +required
	Group string `json:"group"`
	// resourceType of external framework
	// this is the same as Kind in the GVK settings
	// +required
	// +kubebuilder:validation:MaxLength=256
	ResourceType string `json:"resourceType"`
	// version is the version of the api
	// +required
	// +kubebuilder:validation:MaxLength=256
	Version string `json:"version"`
}

type Integrations struct {
	// frameworks are a list of names to be enabled.
	// Possible options:
	//  - "batch/job"
	//  - "kubeflow.org/mpijob"
	//  - "ray.io/rayjob"
	//  - "ray.io/raycluster"
	//  - "jobset.x-k8s.io/jobset"
	//  - "kubeflow.org/paddlejob"
	//  - "kubeflow.org/pytorchjob"
	//  - "kubeflow.org/tfjob"
	//  - "kubeflow.org/xgboostjob"
	//  - "workload.codeflare.dev/appwrapper"
	//  - "pod"
	//  - "deployment" (requires enabling pod integration)
	//  - "statefulset" (requires enabling pod integration)
	//  - "leaderworkerset.x-k8s.io/leaderworkerset" (requires enabling pod integration)
	// +kubebuilder:validation:MaxItems=14
	// +kubebuilder:validation:MinItems=1
	// kubebuilder:validation:UniqueItems=true
	// This is required and must have at least one element.
	// The frameworks are jobs that Kueue will manage.
	// +required
	Frameworks []KueueIntegrations `json:"frameworks"`
	// externalFrameworks are a list of GroupVersionKinds
	// that are managed for Kueue by external controllers;
	// the expected format is `Kind.version.group.com`.
	// These are optional and should only be used if you have an external controller
	// that integrations with kueue.
	// +optional
	// +kubebuilder:validation:MaxItems=4
	ExternalFrameworks []ExternalFramework `json:"externalFrameworks,omitempty"`

	// labelKeysToCopy is a list of label keys that should be copied from the job into the
	// workload object. It is not required for the job to have all the labels from this
	// list. If a job does not have some label with the given key from this list, the
	// constructed workload object will be created without this label. In the case
	// of creating a workload from a composable job (pod group), if multiple objects
	// have labels with some key from the list, the values of these labels must
	// match or otherwise the workload creation would fail. The labels are copied only
	// during the workload creation and are not updated even if the labels of the
	// underlying job are changed.
	// +kubebuilder:validation:items:MaxLength=64
	// +kubebuilder:validation:MaxItems=64
	// +optional
	LabelKeysToCopy []string `json:"labelKeysToCopy,omitempty"`
}
