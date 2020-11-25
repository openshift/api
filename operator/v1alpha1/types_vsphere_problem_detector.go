package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereProblemDetector is used to manage and report health of Openshift
// cluster running in vSphere environment.

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VSphereProblemDetector objects allows configuration and health reporting of
// Vsphere problem detector operator.
type VSphereProblemDetector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec VSphereProblemDetectorSpec `json:"spec"`

	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status VSphereProblemDetectorStatus `json:"status"`
}

// VSphereProblemDetectorSpec is the desired configuration of vSphere problem detector operator.
type VSphereProblemDetectorSpec struct {
	operatorv1.OperatorSpec `json:",inline"`
}

// VSphereProblemDetectorStatus is the status of vSphere problem detector operator and
// health of Openshift and vSphere integration.
type VSphereProblemDetectorStatus struct {
	operatorv1.OperatorStatus `json:",inline"`
	// LastCheckTime reports timestamp of last check
	LastCheckTime metav1.Time `json:"lastCheckTime,omitempty"`
	// List of checks and their status
	Checks []VSphereHealthCheck `json:"checks,omitempty"`
}

// VSphereHealthCheck is used for reporting various health checks the operator
// performs on VSphere and OpenShift integration.
type VSphereHealthCheck struct {
	Name    string             `json:"name"`
	Result  VSphereCheckResult `json:"result"`
	Message string             `json:"message"`
}

// VSphereCheckResult represents status of a health check.
// Possible values are - Passed, Failed and Warning.
type VSphereCheckResult string

var (
	// VSphereCheckPassed indicates health check for vsphere and openshift has passed.
	VSphereCheckPassed VSphereCheckResult = "Passed"
	// VSphereCheckFailed indicates health check for vsphere and openshift has failed.
	VSphereCheckFailed VSphereCheckResult = "Failed"
	// VSphereCheckWarning indicates there were warnings while verifying vsphere and openshift
	// integration.
	VSphereCheckWarning VSphereCheckResult = "Warning"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// VSphereProblemDetectorList contains a list of VSphereProblemDetector
type VSphereProblemDetectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereProblemDetector `json:"items"`
}
