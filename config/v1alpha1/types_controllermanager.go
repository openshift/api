package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerManager holds cluster-wide configuration for the Kubernetes controller manager.
// The resource is a singleton named "cluster".
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2668
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=controllermanagers,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:enable:FeatureGate=DisableForceDetachOnTimeout
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="controllermanager is a singleton, .metadata.name must be 'cluster'"
type ControllerManager struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`
	// spec holds user settable values for configuration
	// +required
	Spec ControllerManagerSpec `json:"spec,omitzero"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ControllerManagerStatus `json:"status,omitzero"`
}

// ControllerManagerSpec defines the desired state of the Kubernetes controller manager
// +kubebuilder:validation:MinProperties=1
type ControllerManagerSpec struct {
	// forceDetachOnTimeout controls whether kube-controller-manager force detaches
	// volumes from a node that is not healthy once the volumes have not been
	// unmounted within the maximum unmount time (6 minutes).
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", volumes are force detached from unhealthy nodes after
	// the maximum unmount time, so that workloads using them can start on other nodes.
	// Force detaching a volume that is still in use by the node can corrupt its data.
	// When set to "Disabled", volumes are not force detached based on the maximum
	// unmount time. Volumes remain attached to an unhealthy node until it recovers,
	// or until the node is tainted with "node.kubernetes.io/out-of-service"
	// as part of the non-graceful node shutdown procedure.
	// When omitted, this means the user has no opinion and the platform is left
	// to choose a reasonable default, which is subject to change over time.
	// The current default is "Enabled".
	// +optional
	ForceDetachOnTimeout ForceDetachOnTimeoutPolicy `json:"forceDetachOnTimeout,omitempty"`
}

// ForceDetachOnTimeoutPolicy describes the policy for force detaching volumes
// when the maximum unmount time is exceeded.
// Valid values are "Enabled" and "Disabled".
// +enum
// +kubebuilder:validation:Enum=Enabled;Disabled
type ForceDetachOnTimeoutPolicy string

const (
	// ForceDetachOnTimeoutEnabled allows kube-controller-manager to force detach
	// volumes from unhealthy nodes once the maximum unmount time is exceeded.
	ForceDetachOnTimeoutEnabled ForceDetachOnTimeoutPolicy = "Enabled"
	// ForceDetachOnTimeoutDisabled prevents kube-controller-manager from force
	// detaching volumes based on the maximum unmount time.
	ForceDetachOnTimeoutDisabled ForceDetachOnTimeoutPolicy = "Disabled"
)

// ControllerManagerStatus defines the observed state of the Kubernetes controller manager
// +kubebuilder:validation:MinProperties=1
type ControllerManagerStatus struct {
	// conditions represent the latest available observations of the configuration state.
	// When omitted, it indicates that no conditions have been reported yet.
	// The maximum number of conditions is 16.
	// When set, at least one condition must be present.
	// Conditions are stored as a map keyed by condition type, ensuring uniqueness.
	// +optional
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerManagerList is a collection of ControllerManager resources.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ControllerManagerList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	Items []ControllerManager `json:"items"`
}
