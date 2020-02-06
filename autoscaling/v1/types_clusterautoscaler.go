package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAutoscaler provides a means to configure the cluster-autoscaler.
//
// +kubebuilder:subresource:status
type ClusterAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Desired state of ClusterAutoscaler resource
	Spec ClusterAutoscalerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of ClusterAutoscaler resource
	Status ClusterAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ClusterAutoscalerSpec represents the desired state of a ClusterAutoscaler.
type ClusterAutoscalerSpec struct {
	// Constraints of autoscaling resources.
	ResourceLimits *ResourceLimits `json:"resourceLimits,omitempty" protobuf:"bytes,1,opt,name=resourceLimits"`

	// Configuration of scale down operation.
	ScaleDown *ScaleDownConfig `json:"scaleDown,omitempty" protobuf:"bytes,2,opt,name=scaleDown"`

	// Gives pods graceful termination time before scaling down.
	MaxPodGracePeriod *int32 `json:"maxPodGracePeriod,omitempty" protobuf:"varint,3,opt,name=maxPodGracePeriod"`

	// To allow users to schedule "best-effort" pods, which shouldn't trigger
	// Cluster Autoscaler actions, but only run when there are spare resources available,
	// More info: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-cluster-autoscaler-work-with-pod-priority-and-preemption
	PodPriorityThreshold *int32 `json:"podPriorityThreshold,omitempty" protobuf:"varint,4,opt,name=podPriorityThreshold"`

	// BalanceSimilarNodeGroups enables/disables the
	// `--balance-similar-node-groups` cluster-autocaler feature.
	// This feature will automatically identify node groups with
	// the same instance type and the same set of labels and try
	// to keep the respective sizes of those node groups balanced.
	BalanceSimilarNodeGroups *bool `json:"balanceSimilarNodeGroups,omitempty" protobuf:"varint,5,opt,name=balanceSimilarNodeGroups"`

	// Enables/Disables `--ignore-daemonsets-utilization` CA feature flag. Should CA ignore DaemonSet pods when calculating resource utilization for scaling down. false by default
	IgnoreDaemonsetsUtilization *bool `json:"ignoreDaemonsetsUtilization,omitempty" protobuf:"varint,6,opt,name=ignoreDaemonsetsUtilization"`

	// Enables/Disables `--skip-nodes-with-local-storage` CA feature flag. If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath. true by default at autoscaler
	SkipNodesWithLocalStorage *bool `json:"skipNodesWithLocalStorage,omitempty" protobuf:"varint,7,opt,name=skipNodesWithLocalStorage"`
}

// ResourceLimits represents global resource limits the cluster-autoscaler
// should not allow to be exceeded.
type ResourceLimits struct {
	// Maximum number of nodes in all node groups.
	// Cluster autoscaler will not grow the cluster beyond this number.
	// +kubebuilder:validation:Minimum=0
	MaxNodesTotal *int32 `json:"maxNodesTotal,omitempty" protobuf:"varint,1,opt,name=maxNodesTotal"`

	// Minimum and maximum number of cores in cluster, in the format <min>:<max>.
	// Cluster autoscaler will not scale the cluster beyond these numbers.
	Cores *ResourceRange `json:"cores,omitempty" protobuf:"bytes,2,opt,name=cores"`

	// Minimum and maximum number of gigabytes of memory in cluster, in the format <min>:<max>.
	// Cluster autoscaler will not scale the cluster beyond these numbers.
	Memory *ResourceRange `json:"memory,omitempty" protobuf:"bytes,3,opt,name=memory"`

	// Minimum and maximum number of different GPUs in cluster, in the format <gpu_type>:<min>:<max>.
	// Cluster autoscaler will not scale the cluster beyond these numbers. Can be passed multiple times.
	GPUS []GPULimit `json:"gpus,omitempty" protobuf:"bytes,4,rep,name=gpus"`
}

// GPULimit represents limits on the number of GPU devices in the cluster.
type GPULimit struct {
	// +kubebuilder:validation:MinLength=1
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`

	// +kubebuilder:validation:Minimum=0
	Min int32 `json:"min" protobuf:"varint,2,opt,name=min"`
	// +kubebuilder:validation:Minimum=1
	Max int32 `json:"max" protobuf:"varint,3,opt,name=max"`
}

// ResourceRange represents a min and max range for a resource.
type ResourceRange struct {
	// +kubebuilder:validation:Minimum=0
	Min int32 `json:"min" protobuf:"varint,1,opt,name=min"`
	Max int32 `json:"max" protobuf:"varint,2,opt,name=max"`
}

// ScaleDownConfig represents the cluster-autoscaler configuration related to
// scaling down operations.
type ScaleDownConfig struct {
	// Should CA scale down the cluster
	Enabled bool `json:"enabled" protobuf:"varint,1,opt,name=enabled"`

	// How long after scale up that scale down evaluation resumes
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterAdd *string `json:"delayAfterAdd,omitempty" protobuf:"bytes,2,opt,name=delayAfterAdd"`

	// How long after node deletion that scale down evaluation resumes, defaults to scan-interval
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterDelete *string `json:"delayAfterDelete,omitempty" protobuf:"bytes,3,opt,name=delayAfterDelete"`

	// How long after scale down failure that scale down evaluation resumes
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterFailure *string `json:"delayAfterFailure,omitempty" protobuf:"bytes,4,opt,name=delayAfterFailure"`

	// How long a node should be unneeded before it is eligible for scale down
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	UnneededTime *string `json:"unneededTime,omitempty" protobuf:"bytes,5,opt,name=unneededTime"`
}

// ClusterAutoscalerStatus defines the observed state of the cluster-autoscaler.
type ClusterAutoscalerStatus struct {
	// TODO: Add status fields.
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAutoscalerList contains a list of ClusterAutoscaler resources.
type ClusterAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ClusterAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
