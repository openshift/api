package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kueue is the Schema for the kueue API
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +k8s:openapi-gen=true
// +genclient
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
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

type ManageJobsWithoutQueueNameOption string

const (
	// NoQueueName means that all jobs will be gated by Kueue
	NoQueueName ManageJobsWithoutQueueNameOption = "NoQueueName"
	// QueueName means that the jobs require a queue label.
	QueueName ManageJobsWithoutQueueNameOption = "QueueName"
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
	FeatureGates map[string]bool `json:"featureGates,omitempty"`
	// resources provides additional configuration options for handling the resources.
	// Supports https://github.com/kubernetes-sigs/kueue/blob/release-0.10/keps/2937-resource-transformer/README.md
	// +optional
	Resources Resources `json:"resources,omitempty"`
	// ManageJobsWithoutQueueName controls whether or not Kueue reconciles
	// jobs that don't set the annotation kueue.x-k8s.io/queue-name.
	// Allowed values are NoQueueName and QueueName
	// Default will be QueueName
	// +optional
	ManageJobsWithoutQueueName *ManageJobsWithoutQueueNameOption `json:"manageJobsWithoutQueueName,omitempty"`
	// ManagedJobsNamespaceSelector can be used to omit some namespaces from ManagedJobsWithoutQueueName
	// +optional
	ManagedJobsNamespaceSelector *metav1.LabelSelector `json:"managedJobsNamespaceSelector,omitempty"`
	// FairSharing controls the fair sharing semantics across the cluster.
	FairSharing FairSharing `json:"fairSharing,omitempty"`
	// Disable Metrics
	// Microshift does not enable metrics by default
	// Default will assume metrics are enabled.
	// +optional
	DisableMetrics *bool `json:"disableMetrics,omitempty"`
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
	// +required
	Items []Kueue `json:"items"`
}

// These structs come directly from Kueue.

type Resources struct {
	// ExcludedResourcePrefixes defines which resources should be ignored by Kueue
	ExcludeResourcePrefixes []string `json:"excludeResourcePrefixes,omitempty"`

	// Transformations defines how to transform PodSpec resources into Workload resource requests.
	// This is intended to be a map with Input as the key (enforced by validation code)
	Transformations []ResourceTransformation `json:"transformations,omitempty"`
}

type ResourceTransformationStrategy string

const Retain ResourceTransformationStrategy = "Retain"
const Replace ResourceTransformationStrategy = "Replace"

type ResourceTransformation struct {
	// Input is the name of the input resource.
	Input corev1.ResourceName `json:"input"`

	// Strategy specifies if the input resource should be replaced or retained.
	// Defaults to Retain
	Strategy *ResourceTransformationStrategy `json:"strategy,omitempty"`

	// Outputs specifies the output resources and quantities per unit of input resource.
	// An empty Outputs combined with a `Replace` Strategy causes the Input resource to be ignored by Kueue.
	Outputs corev1.ResourceList `json:"outputs,omitempty"`
}

type PreemptionStrategy string

const (
	LessThanOrEqualToFinalShare PreemptionStrategy = "LessThanOrEqualToFinalShare"
	LessThanInitialShare        PreemptionStrategy = "LessThanInitialShare"
)

type FairSharing struct {
	// enable indicates whether to enable fair sharing for all cohorts.
	// Defaults to false.
	Enable bool `json:"enable"`

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
	PreemptionStrategies []PreemptionStrategy `json:"preemptionStrategies,omitempty"`
}

// WaitForPodsReady defines configuration for the Wait For Pods Ready feature,
// which is used to ensure that all Pods are ready within the specified time.
type WaitForPodsReady struct {
	// Enable indicates whether to enable wait for pods ready feature.
	// Defaults to false.
	Enable bool `json:"enable,omitempty"`

	// Timeout defines the time for an admitted workload to reach the
	// PodsReady=true condition. When the timeout is exceeded, the workload
	// evicted and requeued in the same cluster queue.
	// Defaults to 5min.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// BlockAdmission when true, cluster queue will block admissions for all
	// subsequent jobs until the jobs reach the PodsReady=true condition.
	// This setting is only honored when `Enable` is set to true.
	BlockAdmission *bool `json:"blockAdmission,omitempty"`

	// RequeuingStrategy defines the strategy for requeuing a Workload.
	// +optional
	RequeuingStrategy *RequeuingStrategy `json:"requeuingStrategy,omitempty"`

	// RecoveryTimeout defines an opt-in timeout, measured since the
	// last transition to the PodsReady=false condition after a Workload is Admitted and running.
	// Such a transition may happen when a Pod failed and the replacement Pod
	// is awaited to be scheduled.
	// After exceeding the timeout the corresponding job gets suspended again
	// and requeued after the backoff delay. The timeout is enforced only if waitForPodsReady.enable=true.
	// If not set, there is no timeout.
	// +optional
	RecoveryTimeout *metav1.Duration `json:"recoveryTimeout,omitempty"`
}

type MultiKueue struct {
	// GCInterval defines the time interval between two consecutive garbage collection runs.
	// Defaults to 1min. If 0, the garbage collection is disabled.
	// +optional
	GCInterval *metav1.Duration `json:"gcInterval"`

	// Origin defines a label value used to track the creator of workloads in the worker
	// clusters.
	// This is used by multikueue in components like its garbage collector to identify
	// remote objects that ware created by this multikueue manager cluster and delete
	// them if their local counterpart no longer exists.
	// +optional
	Origin *string `json:"origin,omitempty"`

	// WorkerLostTimeout defines the time a local workload's multikueue admission check state is kept Ready
	// if the connection with its reserving worker cluster is lost.
	//
	// Defaults to 15 minutes.
	// +optional
	WorkerLostTimeout *metav1.Duration `json:"workerLostTimeout,omitempty"`
}

type RequeuingStrategy struct {
	// Timestamp defines the timestamp used for re-queuing a Workload
	// that was evicted due to Pod readiness. The possible values are:
	//
	// - `Eviction` (default) indicates from Workload `Evicted` condition with `PodsReadyTimeout` reason.
	// - `Creation` indicates from Workload .metadata.creationTimestamp.
	//
	// +optional
	Timestamp *RequeuingTimestamp `json:"timestamp,omitempty"`

	// BackoffLimitCount defines the maximum number of re-queuing retries.
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

	// BackoffBaseSeconds defines the base for the exponential backoff for
	// re-queuing an evicted workload.
	//
	// Defaults to 60.
	// +optional
	BackoffBaseSeconds *int32 `json:"backoffBaseSeconds,omitempty"`

	// BackoffMaxSeconds defines the maximum backoff time to re-queue an evicted workload.
	//
	// Defaults to 3600.
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
	// List of framework names to be enabled.
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
	Frameworks []string `json:"frameworks,omitempty"`
	// List of GroupVersionKinds that are managed for Kueue by external controllers;
	// the expected format is `Kind.version.group.com`.
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
	LabelKeysToCopy []string `json:"labelKeysToCopy,omitempty"`
}
