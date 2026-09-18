package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerManager holds cluster-wide configuration shared by the controller managers
// in the system, among them especially kube-controller-manager.
// The resource is a singleton named "cluster".
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2668
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=controllermanagers,scope=Cluster
// +openshift:enable:FeatureGate=ControllerManagerConfig
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="controllermanager is a singleton, .metadata.name must be 'cluster'"
type ControllerManager struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`
	// spec holds user settable values for configuration.
	// The only way to express no opinion in the spec is to not create the ControllerManager object at all.
	// +required
	Spec ControllerManagerSpec `json:"spec,omitzero"`
}

// ControllerManagerSpec defines the desired state of the controller managers
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
	// While this is the only field in spec, it must be set, because spec cannot be empty.
	// +optional
	// +openshift:enable:FeatureGate=DisableForceDetachOnTimeout
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

	// items is a list of ControllerManager resources
	Items []ControllerManager `json:"items"`
}
