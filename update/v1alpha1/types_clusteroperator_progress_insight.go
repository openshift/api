package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterOperatorProgressInsight reports the state of a Cluster Operator (an individual control plane component) during an update
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=clusteroperatorprogressinsights,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Provides information about a Cluster Operator update"
// +kubebuilder:metadata:annotations="displayName=ClusterOperatorProgressInsights"
// ClusterOperatorProgressInsight reports the state of a Cluster Operator (an individual control plane component) during an update
type ClusterOperatorProgressInsight struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, ClusterOperatorProgressInsight is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +required
	Spec ClusterOperatorProgressInsightSpec `json:"spec"`
	// status exposes the health and status of the ongoing cluster operator update
	// +optional
	Status ClusterOperatorProgressInsightStatus `json:"status"`
}

// ClusterOperatorProgressInsightSpec is empty for now, ClusterOperatorProgressInsightSpec is purely status-reporting API. In the future spec may be used
// to hold configuration to drive what information is surfaced and how
type ClusterOperatorProgressInsightSpec struct {
}

// ClusterOperatorProgressInsight reports the state of a ClusterOperator resource (which represents a control plane
// component update in standalone clusters), during the update
// +kubebuilder:validation:XValidation:rule="self.name == self.resource.name",message=".name must match .resource.name"
type ClusterOperatorProgressInsightStatus struct {
	// conditions provide details about the operator. It contains at most 10 items. Known conditions are:
	// - Updating: whether the operator is updating; When Updating=False, the reason field can be Pending or Updated
	// - Healthy: whether the operator is considered healthy; When Healthy=False, the reason field can be Unavailable or Degraded, and Unavailable is "stronger" than Degraded
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// name is the name of the operator
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Name string `json:"name"`

	// resource is the ClusterOperator resource that represents the operator
	//
	// +Note: By OpenShift API conventions, in isolation this should be a specialized reference that refers just to
	// +resource name (because the rest is implied by status insight type). However, because we use resource references in
	// +many places and this API is intended to be consumed by clients, not produced, consistency seems to be more valuable
	// +than type safety for producers.
	// +required
	// +kubebuilder:validation:XValidation:rule="self.group == 'config.openshift.io' && self.resource == 'clusteroperators'",message="resource must be a clusteroperators.config.openshift.io resource"
	Resource ResourceRef `json:"resource"`
}

// ClusterOperatorProgressInsightConditionType are types of conditions that can be reported on ClusterOperator status insights
type ClusterOperatorProgressInsightConditionType string

const (
	// Updating condition communicates whether the ClusterOperator is updating
	ClusterOperatorProgressInsightUpdating ClusterOperatorProgressInsightConditionType = "Updating"
	// Healthy condition communicates whether the ClusterOperator is considered healthy
	ClusterOperatorProgressInsightHealthy ClusterOperatorProgressInsightConditionType = "Healthy"
)

// ClusterOperatorUpdatingReason are well-known reasons for the Updating condition on ClusterOperator status insights
type ClusterOperatorUpdatingReason string

const (
	// Updated is used with Updating=False when the ClusterOperator finished updating
	ClusterOperatorUpdatingReasonUpdated ClusterOperatorUpdatingReason = "Updated"
	// Pending is used with Updating=False when the ClusterOperator is not updating and is still running previous version
	ClusterOperatorUpdatingReasonPending ClusterOperatorUpdatingReason = "Pending"
	// Progressing is used with Updating=True when the ClusterOperator is updating
	ClusterOperatorUpdatingReasonProgressing ClusterOperatorUpdatingReason = "Progressing"
	// CannotDetermine is used with Updating=Unknown
	ClusterOperatorUpdatingCannotDetermine ClusterOperatorUpdatingReason = "CannotDetermine"
)

// ClusterOperatorHealthyReason are well-known reasons for the Healthy condition on ClusterOperator status insights
type ClusterOperatorHealthyReason string

const (
	// AsExpected is used with Healthy=True when no issues are observed
	ClusterOperatorHealthyReasonAsExpected ClusterOperatorHealthyReason = "AsExpected"
	// Unavailable is used with Healthy=False when the ClusterOperator has Available=False condition
	ClusterOperatorHealthyReasonUnavailable ClusterOperatorHealthyReason = "Unavailable"
	// Degraded is used with Healthy=False when the ClusterOperator has Degraded=True condition
	ClusterOperatorHealthyReasonDegraded ClusterOperatorHealthyReason = "Degraded"
	// CannotDetermine is used with Healthy=Unknown
	ClusterOperatorHealthyReasonCannotDetermine ClusterOperatorHealthyReason = "CannotDetermine"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterOperatorProgressInsightList is a list of ClusterOperatorProgressInsightList resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ClusterOperatorProgressInsightList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a list of ClusterOperatorProgressInsight resources
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Items []ClusterOperatorProgressInsight `json:"items"`
}
