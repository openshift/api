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
	// integrations is a required field that configures the Kueue's workload integrations.
	// Kueue has both standard integrations, known as job frameworks, and external integrations known as external frameworks.
	// Kueue will only manage workloads that correspond to the specified integrations.
	// +required
	Integrations Integrations `json:"integrations"`
	// queueLabelPolicy controls how kueue manages workloads
	// The default behavior of Kueue will manage workloads that have a queue-name label.
	// +optional
	QueueLabelPolicy QueueLabelPolicy `json:"queueLabelPolicy,omitempty"`
	// kueueGangSchedulingPolicy controls how Kueue admits workloads
	// +optional
	KueueGangSchedulingPolicy KueueGangSchedulingPolicy `json:"kueueGangSchedulingPolicy,omitempty"`
	// fairSharing TODO not done yet
	// +optional
	FairSharing FairSharing `json:"fairSharing,omitempty"`
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

// +kubebuilder:validation:Enum=BatchJob;RayJob;RayCluster;JobSet;MPIJob;PaddleJob;PytorchJob;TFJob;XGBoostJob;AppWrapper;Pod;Deployment;StatefulSet;LeaderWorkerSet
type KueueIntegration string

const (
	KueueIntegrationBatchJob        KueueIntegration = "BatchJob"
	KueueIntegrationRayJob          KueueIntegration = "RayJob"
	KueueIntegrationRayCluster      KueueIntegration = "RayCluster"
	KueueIntegrationJobSet          KueueIntegration = "JobSet"
	KueueIntegrationMPIJob          KueueIntegration = "MPIJob"
	KueueIntegrationPaddeJob        KueueIntegration = "PaddeJob"
	KueueIntegrationPyTorchJob      KueueIntegration = "PyTorchJob"
	KueueIntegrationTFJob           KueueIntegration = "TFJob"
	KueueIntegrationXGBoostJob      KueueIntegration = "XGBoostJob"
	KueueIntegrationAppWrapper      KueueIntegration = "AppWrapper"
	KueueIntegrationPod             KueueIntegration = "Pod"
	KueueIntegrationDeployment      KueueIntegration = "Deployment"
	KueueIntegrationStatefulSet     KueueIntegration = "StatefulSet"
	KueueIntegrationLeaderWorkerSet KueueIntegration = "LeaderWorkerSet"
)

// This is the GVR for an external framework.
// Controller runtime requires this in this format
// for api discoverability.
type ExternalFramework struct {
	// group is the API group of the externalFramework.
	// Must be a valid DNS 1123 subdomain consisting of of lower-case alphanumeric characters,
	// hyphens and periods, of at most 253 characters in length.
	// Each period separated segment within the subdomain must start and end with an alphanumeric character.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.dns1123Label().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Group string `json:"group"`
	// resource is the Resource type of the external framework.
	// Resource types are lowercase and plural (e.g. pods, deployments).
	// Must be a valid DNS 1123 label consisting of a lower-case alphanumeric string
	// and hyphens of at most 63 characters in length.
	// The value must start and end with an alphanumeric character.
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.dns1123Label().validate(self).hasValue()",message="a lowercase RFC 1123 label must consist of lower case alphanumeric characters and '-', and must start and end with an alphanumeric character."
	// +required
	Resource string `json:"resource"`
	// version is the version of the api (e.g. v1alpha1, v1beta1, v1).
	// Must be a valid DNS 1035 label consisting of a lower-case alphanumeric string
	// and hyphens of at most 63 characters in length.
	// The value must start with an alphabetic character and end with an alphanumeric character.
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.dns1035Label().validate(self).hasValue()",message="a lowercase RFC 1035 label must consist of lower case alphanumeric characters, '-' or '.', and must start with an alphabetic character and end with an alphanumeric character."
	// +required
	Version string `json:"version"`
}

// This is the integrations for Kueue.
// Kueue uses these apis to determine
// which jobs will be managed by Kueue.
type Integrations struct {
	// frameworks are a unique list of names to be enabled.
	// This is required and must have at least one element.
	// Each framework represents a type of job that Kueue will manage.
	// Frameworks are a list of frameworks that Kueue has support for.
	// The allowed values are BatchJob, RayJob, RayCluster, JobSet, MPIJob, PaddleJob, PytorchJob, TFJob, XGBoostJob, AppWrapper, Pod, Deployment, StatefulSet and LeaderWorkerSet.
	// +kubebuilder:validation:MaxItems=14
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:XValidation:rule="self.all(x, self.exists_one(y, x == y))",message="each item in frameworks must be unique"
	// +listType=atomic
	// +required
	Frameworks []KueueIntegration `json:"frameworks"`
	// externalFrameworks are a list of GroupVersionResources
	// that are managed for Kueue by external controllers.
	// These are optional and should only be used if you have an external controller
	// that integrates with Kueue.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=32
	// +optional
	ExternalFrameworks []ExternalFramework `json:"externalFrameworks,omitempty"`
	// labelKeysToCopy are a list of label keys that are copied once a workload is created.
	// These keys are persisted to the internal Kueue workload object.
	// If not specified, only the Kueue labels will be copied.
	// +kubebuilder:validation:MaxItems=64
	// +listType=atomic
	// +optional
	LabelKeysToCopy []LabelKeys `json:"labelKeysToCopy,omitempty"`
}

type LabelKeys struct {
	// key is the label key
	// A label key must be a valid qualified name consisting of a lower-case alphanumeric string,
	// and hyphens of at most 63 characters in length.
	// The name must start and end with an alphanumeric character.
	// The name may be optionally prefixed with a subdomain consisting of lower-case alphanumeric characters,
	// hyphens and periods, of at most 253 characters in length.
	// Each period separated segment within the subdomain must start and end with an alphanumeric character.
	// The optional prefix and the name are separate by a forward slash (/).
	// +kubebuilder:validation:MaxLength=317
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="!format.qualifiedName().validate(self).hasValue()",message="a qualified name must consist of a lower-case alphanumeric and hyphenated string of at most 63 characters in length, starting and ending with alphanumeric chracters. The name may be optionally prefixed with a subdomain consisting of lower-case alphanumeric characters, hyphens and periods, of at most 253 characters in length. Each period separated segment within the subdomain must start and end with an alphanumeric character."
	// +optional
	Key string `json:"key,omitempty"`
}

// +kubebuilder:validation:Enum=ByWorkload;Disabled
type KueueGangSchedulingPolicyOptions string

const (
	KueueGangSchedulingPolicyEvictNotReadyWorkloads KueueGangSchedulingPolicyOptions = "ByWorkload"
	KueueGangSchedulingPolicyDisabled               KueueGangSchedulingPolicyOptions = "Disabled"
)

// +kubebuilder:validation:Enum=Parallel;Sequential
type KueueGangSchedulingAdmissionOptions string

const (
	KueueGangSchedulingAdmissionOptionsSequential KueueGangSchedulingAdmissionOptions = "Sequential"
	KueueGangSchedulingAdmissionOptionsParallel   KueueGangSchedulingAdmissionOptions = "Parallel"
)

// Kueue provides the ability to admit workloads all in one (gang admission)
// and evicts workloads if they are not ready within a specific time.
type KueueGangSchedulingPolicy struct {
	// policy allows for changing the kinds of gang scheduling Kueue does.
	// This is an optional field.
	// The allowed values are ByWorkload and Disabled.
	// The default value will be Disabled.
	// ByWorkload allows for configuration how admission is performed
	// for Kueue.
	// +optional
	Policy KueueGangSchedulingPolicyOptions `json:"policy"`
	// byWorkload controls how admission is done.
	// The options are Sequential and Parallel.
	// Sequential means workloads are admitted in sequential order and Kueue waits
	// for those workloads to be ready.
	// Parallel admits workloads in parallel and does not wait on these workloads to be ready.
	// Workloads can be evicted in the background if they are not ready within default settings.
	// +kubebuilder:validation:XValidation:rule="self.policy==ByWorkload",message="byWorkload is only valid if policy equals ByWorkload"
	// +optional
	ByWorkload KueueGangSchedulingAdmissionOptions `json:"byWorkload"`
}

// +kubebuilder:validation:Enum=QueueNameRequired;QueueNameOptional
type QueueLabelNamePolicy string

const (
	QueueLabelNamePolicyRequired QueueLabelNamePolicy = "QueueNameRequired"
	QueueLabelNamePolicyOptional QueueLabelNamePolicy = "QueueNameOptional"
)

type QueueLabelPolicy struct {
	// queueLabelPolicy controls whether or not Kueue reconciles
	// jobs that don't set the label kueue.x-k8s.io/queue-name.
	// The allowed values are QueueNameRequired and QueueNameOptional.
	// If set to QueueNameRequired, then those jobs will be suspended and never started unless
	// they are assigned a queue and eventually admitted. This also applies to
	// jobs created before starting the kueue controller.
	// Defaults to QueueNameOptional; therefore, those jobs are not managed and if they are created
	// unsuspended, they will start immediately.
	// +optional
	QueueLabelPolicy QueueLabelNamePolicy `json:"queueLabelPolicy"`
}

// +kubebuilder:validation:Enum=LessThanOrEqualToFinalShare;LessThanInitialShare
type PreemptionStrategy string

const (
	LessThanOrEqualToFinalShare PreemptionStrategy = "LessThanOrEqualToFinalShare"
	LessThanInitialShare        PreemptionStrategy = "LessThanInitialShare"
)

type FairSharing struct {
	// preemptionStrategies indicates which constraints should a preemption satisfy.
	// The preemption algorithm will only use the next strategy in the list if the
	// incoming workload (preemptor) doesn't fit after using the previous strategies.
	// Possible values are:
	// - LessThanOrEqualToFinalShare: Only preempt a workload if the share of the preemptor CQ
	//   with the preemptor workload is less than or equal to the share of the preemptee CQ
	//   without the workload to be preempted.
	//   This strategy might favor preemption of smaller workloads in the preemptee CQ,
	//   regardless of priority or start time, in an effort to keep the share of the CQ
	//   as high as possible.
	// - LessThanInitialShare: Only preempt a workload if the share of the preemptor CQ
	//   with the incoming workload is strictly less than the share of the preemptee CQ.
	//   This strategy doesn't depend on the share usage of the workload being preempted.
	//   As a result, the strategy chooses to preempt workloads with the lowest priority and
	//   newest start time first.
	// The default strategy is ["LessThanOrEqualToFinalShare", "LessThanInitialShare"].
	// +kubebuilder:validation:MaxItems=2
	// +optional
	PreemptionStrategies []PreemptionStrategy `json:"preemptionStrategies,omitempty"`
}
