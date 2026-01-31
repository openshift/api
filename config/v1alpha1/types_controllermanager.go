package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerManager holds cluster-wide config information to run the Kubernetes controller manager
// and influence its placement decisions. The canonical name for this config is `cluster`.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2668
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=controllermanagers,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:metadata:annotations=release.openshift.io/bootstrap-required=true
// +openshift:enable:FeatureGate=DisableForceDetachOnTimeout
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
	// forceDetachOnTimeout expresses whether to allow kube-controller-manager
	// to force detach volumes when unmount takes longer than the timeout.
	// Valid values are Enabled and Disabled. If omitted, the default is Enabled.
	// +default="Enabled"
	// +optional
	ForceDetachOnTimeout ForceDetachOnTimeoutPolicy `json:"forceDetachOnTimeout,omitempty"`
}

// +kubebuilder:validation:Enum=Enabled;Disabled
type ForceDetachOnTimeoutPolicy string

const (
	// ForceDetachOnTimeoutEnabled will allow kube-controller-manager to
	// force detach volumes based on maximum unmount time and node status.
	ForceDetachOnTimeoutEnabled  ForceDetachOnTimeoutPolicy = "Enabled"
	// ForceDetachOnTimeoutDisabled will prevent kube-controller-manager
	// from force detaching volumes.
	ForceDetachOnTimeoutDisabled ForceDetachOnTimeoutPolicy = "Disabled"
)

// ControllerManagerStatus defines the observed state of the Kubernetes controller manager
// +kubebuilder:validation:MinProperties=1
type ControllerManagerStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ControllerManagerList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	Items []ControllerManager `json:"items"`
}
