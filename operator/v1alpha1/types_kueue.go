package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
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
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="olm is a singleton, .metadata.name must be 'cluster'"
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
	// config that is persisted to a config map
	// +required
	Config KueueConfiguration `json:"config"`
}

// +kubebuilder:validation:Enum="NoQueueName;QueueName"
type ManageJobsWithoutQueueNameOption string

const (
	// NoQueueName means that all jobs will be gated by Kueue
	NoQueueName ManageJobsWithoutQueueNameOption = "NoQueueName"
	// QueueName means that the jobs require a queue label.
	QueueName ManageJobsWithoutQueueNameOption = "QueueName"
)

// +kubebuilder:validation:Enum="Enabled;Disabled"
type EnabledOrDisabled string

const (
	Enabled  EnabledOrDisabled = "Enabled"
	Disabled EnabledOrDisabled = "Disabled"
)

type KueueConfiguration struct {
	// waitForPodsReady configures gang admission
	// +optional
	WaitForPodsReady WaitForPodsReady `json:"waitForPodsReady,omitempty"`
	// integrations are the types of integrations Kueue will manager
	// +required
	Integrations Integrations `json:"integrations"`
	// featureGates are advanced features for Kueue
	// +optional
	FeatureGates map[string]EnabledOrDisabled `json:"featureGates,omitempty"`
	// resources provides additional configuration options for handling the resources.
	// Supports https://github.com/kubernetes-sigs/kueue/blob/release-0.10/keps/2937-resource-transformer/README.md
	// +optional
	Resources Resources `json:"resources,omitempty"`
	// manageJobsWithoutQueueName controls whether or not Kueue reconciles
	// jobs that don't set the annotation kueue.x-k8s.io/queue-name.
	// Allowed values are NoQueueName and QueueName
	// Default will be QueueName
	// +kubebuilder:default=QueueName
	// +optional
	ManageJobsWithoutQueueName *ManageJobsWithoutQueueNameOption `json:"manageJobsWithoutQueueName,omitempty"`
	// managedJobsNamespaceSelector can be used to omit some namespaces from ManagedJobsWithoutQueueName
	// Only valid if ManagedJobsWithoutQueueName is NoQueueName
	// +optional
	ManagedJobsNamespaceSelector *metav1.LabelSelector `json:"managedJobsNamespaceSelector,omitempty"`
	// fairSharing controls the fair sharing semantics across the cluster.
	// +optional
	FairSharing FairSharing `json:"fairSharing,omitempty"`
	// metrics allows one to change if metrics
	// are enabled or disabled.
	// Microshift does not enable metrics by default
	// Default will assume metrics are enabled.
	// +kubebuilder:default=Enabled
	// +optional
	Metrics *EnabledOrDisabled `json:"metrics,omitempty"`
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

// These structs come directly from Kueue.
type Resources struct {
	// excludeResourcePrefixes defines which resources should be ignored by Kueue
	// +optional
	// +kubebuilder:validation:items:MaxLength=64
	// +kubebuilder:validation:MaxItems=64
	ExcludeResourcePrefixes []string `json:"excludeResourcePrefixes,omitempty"`

	// transformations defines how to transform PodSpec resources into Workload resource requests.
	// This is intended to be a map with Input as the key (enforced by validation code)
	// +optional
	// +kubebuilder:validation:items:MaxLength=64
	// +kubebuilder:validation:MaxItems=64
	Transformations []ResourceTransformation `json:"transformations,omitempty"`
}

// +kubebuilder:validation:Enum=Retain;Replace
type ResourceTransformationStrategy string

const Retain ResourceTransformationStrategy = "Retain"
const Replace ResourceTransformationStrategy = "Replace"

type ResourceTransformation struct {
	// input is the name of the input resource.
	// +required
	Input corev1.ResourceName `json:"input"`

	// strategy specifies if the input resource should be replaced or retained.
	// +kubebuilder:default=Retain
	// +optional
	Strategy *ResourceTransformationStrategy `json:"strategy,omitempty"`

	// outputs specifies the output resources and quantities per unit of input resource.
	// An empty Outputs combined with a `Replace` Strategy causes the Input resource to be ignored by Kueue.
	// +optional
	Outputs corev1.ResourceList `json:"outputs,omitempty"`
}

// +kubebuilder:validation:Enum=LessThanOrEqualToFinalShare;LessThanInitialShare
type PreemptionStrategy string

const (
	LessThanOrEqualToFinalShare PreemptionStrategy = "LessThanOrEqualToFinalShare"
	LessThanInitialShare        PreemptionStrategy = "LessThanInitialShare"
)

type FairSharing struct {
	// enable indicates whether to enable fair sharing for all cohorts.
	// Defaults to false.
	// +kubebuilder:default=Disabled
	// +optional
	Enable EnabledOrDisabled `json:"enable"`

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
	// +optional
	// +kubebuilder:validation:MaxItems=2
	PreemptionStrategies []PreemptionStrategy `json:"preemptionStrategies,omitempty"`
}

// WaitForPodsReady defines configuration for the Wait For Pods Ready feature,
// which is used to ensure that all Pods are ready within the specified time.
type WaitForPodsReady struct {
	// timeout defines the time for an admitted workload to reach the
	// PodsReady=true condition. When the timeout is exceeded, the workload
	// evicted and requeued in the same cluster queue.
	// Defaults to 5min.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// blockAdmission when true, cluster queue will block admissions for all
	// subsequent jobs until the jobs reach the PodsReady=true condition.
	// This setting is only honored when `Enable` is set to true.
	// +optional
	BlockAdmission *EnabledOrDisabled `json:"blockAdmission,omitempty"`

	// requeuingStrategy defines the strategy for requeuing a Workload.
	// +optional
	RequeuingStrategy *RequeuingStrategy `json:"requeuingStrategy,omitempty"`

	// recoveryTimeout defines an opt-in timeout, measured since the
	// last transition to the PodsReady=false condition after a Workload is Admitted and running.
	// Such a transition may happen when a Pod failed and the replacement Pod
	// is awaited to be scheduled.
	// After exceeding the timeout the corresponding job gets suspended again
	// and requeued after the backoff delay. The timeout is enforced only if waitForPodsReady.enable=true.
	// If not set, there is no timeout.
	// +optional
	RecoveryTimeout *metav1.Duration `json:"recoveryTimeout,omitempty"`
}

type RequeuingStrategy struct {
	// timestamp defines the timestamp used for re-queuing a Workload
	// that was evicted due to Pod readiness. The possible values are:
	//
	// - `Eviction` (default) indicates from Workload `Evicted` condition with `PodsReadyTimeout` reason.
	// - `Creation` indicates from Workload .metadata.creationTimestamp.
	//
	// +kubebuilder:default=Eviction
	// +kubebuilder:validation:MaxLength=8
	// +optional
	Timestamp *RequeuingTimestamp `json:"timestamp,omitempty"`

	// backoffLimitCount defines the maximum number of re-queuing retries.
	// Once the number is reached, the workload is deactivated (`.spec.activate`=`false`).
	// When it is null, the workloads will repeatedly and endless re-queueing.
	//
	// Every backoff duration is about "b*2^(n-1)+Rand" where:
	// - "b" represents the base set by "BackoffBaseSeconds" parameter,
	// - "n" represents the "workloadStatus.requeueState.count",
	// - "Rand" represents the random jitter.
	// During this time, the workload is taken as an inadmissible and
	// other workloads will have a chance to be admitted.
	// By default, the consecutive requeue delays are around: (60s, 120s, 240s, ...).
	//
	// Defaults to null.
	// +optional
	BackoffLimitCount *int32 `json:"backoffLimitCount,omitempty"`

	// backoffBaseSeconds defines the base for the exponential backoff for
	// re-queuing an evicted workload.
	//
	// Defaults to 60.
	// +kubebuilder:default=60
	// +optional
	BackoffBaseSeconds *int32 `json:"backoffBaseSeconds,omitempty"`

	// backoffMaxSeconds defines the maximum backoff time to re-queue an evicted workload.
	//
	// Defaults to 3600.
	// +kubebuilder:default=3600
	// +optional
	BackoffMaxSeconds *int32 `json:"backoffMaxSeconds,omitempty"`
}

type RequeuingTimestamp string

const (
	// CreationTimestamp timestamp (from Workload .metadata.creationTimestamp).
	CreationTimestamp RequeuingTimestamp = "Creation"

	// EvictionTimestamp timestamp (from Workload .status.conditions).
	EvictionTimestamp RequeuingTimestamp = "Eviction"
)

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
	// +kubebuilder:validation:items:MaxLength=64
	// +required
	Frameworks []string `json:"frameworks"`
	// externalFrameworks are a list of GroupVersionKinds
	// that are managed for Kueue by external controllers;
	// the expected format is `Kind.version.group.com`.
	// As far as
	// +optional
	// +kubebuilder:validation:items:MaxLength=64
	// +kubebuilder:validation:MaxItems=4
	ExternalFrameworks []string `json:"externalFrameworks,omitempty"`

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
