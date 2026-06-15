package v1alpha1

import (
	v1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ingress contains configuration options specific to the Ingress Operator itself.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:api-approved.openshift.io=<TBD>
// +openshift:file-pattern=cvoRunLevel=0000_50,operatorName=ingress,operatorOrdering=02
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ingresses,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:capability=Ingress
// +openshift:enable:FeatureGate=GatewayAPIManagementMode
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="ingress is a singleton; the .metadata.name field must be 'cluster'"
type Ingress struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration.
	// +required
	Spec IngressSpec `json:"spec,omitempty,omitzero"`

	// status holds observed values from the cluster.
	// +optional
	Status IngressStatus `json:"status,omitempty,omitzero"`
}

// IngressSpec defines the desired configuration of the cluster-ingress-operator
// +kubebuilder:validation:MinProperties=1
type IngressSpec struct {
	// Inline OperatorSpec for standard operator fields
	// (managementState, logLevel, etc.)
	v1.OperatorSpec `json:",inline"`

	// gatewayAPI holds configuration for Gateway API
	// integration, including how the Cluster Ingress Operator
	// manages Gateway API CRDs and the Gateway controller stack.
	//
	// When omitted, the operator uses the default Gateway API
	// configuration, which sets managementMode to "Managed".
	//
	// +optional
	// +openshift:enable:FeatureGate=GatewayAPIManagementMode
	GatewayAPI GatewayAPIIngressConfig `json:"gatewayAPI,omitempty,omitzero"`
}

// IngressStatus describe the current state of cluster-ingress-operator
// +kubebuilder:validation:MinProperties=0
type IngressStatus struct {
	// Inline OperatorStatus for standard operator status fields
	// (conditions, version, observedGeneration, etc.).
	// conditions holds a list of conditions representing the
	// operator's current state. Gateway API CRD management
	// conditions are reported here with the "GatewayAPI" prefix:
	//
	// "GatewayAPICRDsManaged" indicates whether CIO is actively
	// managing Gateway API CRDs:
	//   - status: True, reason: "ManagedByCIO" — CIO is
	//     installing, protecting (via VAP), and upgrading CRDs.
	//   - status: False, reason: "Unmanaged" — the administrator
	//     chose Unmanaged mode; CIO does not manage CRDs or the
	//     Gateway controller stack.
	//
	// "GatewayAPICRDsPresent" indicates whether Gateway API CRDs
	// exist on the cluster:
	//   - status: True, reason: "CRDsFound" — Gateway API CRDs
	//     are present on the cluster.
	//   - status: False, reason: "CRDsNotFound" — Gateway API
	//     CRDs are not present on the cluster.
	//
	// "GatewayAPICRDsCompliant" indicates whether the installed
	// CRDs match the version expected by this CIO release:
	//   - status: True, reason: "VersionMatch" — installed CRDs
	//     match the expected version.
	//   - status: False, reason: "VersionMismatch" — installed
	//     CRDs do not match the expected version. The message
	//     includes expected and actual versions and a pointer
	//     to where valid manifests can be obtained.
	//   - status: Unknown, reason: "NotApplicable" — compliance
	//     check is not applicable (e.g., Unmanaged mode with no
	//     CRDs present).
	v1.OperatorStatus `json:",inline"`
}

// GatewayAPIManagementMode describes how the Cluster Ingress
// Operator manages Gateway API Custom Resource Definitions.
//
// +kubebuilder:validation:Enum=Managed;Unmanaged
type GatewayAPIManagementMode string

const (
	// GatewayAPIManagementModeManaged means CIO installs, owns,
	// protects (via VAP), and upgrades the Gateway API CRDs.
	// CIO also deploys the full Gateway controller stack (the
	// Istio instance deployed by CIO, GatewayClass, Gateway).
	// This is the default mode and the only fully supported
	// configuration.
	GatewayAPIManagementModeManaged GatewayAPIManagementMode = "Managed"

	// GatewayAPIManagementModeUnmanaged means CIO does NOT
	// install or manage Gateway API CRDs and does NOT deploy
	// the Gateway controller stack. The customer or a
	// third-party product is responsible for bringing their own
	// CRDs and Gateway controller. CIO reports observational
	// status only. This mode signals to layered products that
	// the installed CRDs may not be the ones supported by the
	// OpenShift Gateway API implementation.
	GatewayAPIManagementModeUnmanaged GatewayAPIManagementMode = "Unmanaged"
)

// GatewayAPIIngressConfig holds configuration for Gateway API
// integration in the Cluster Ingress Operator.
// +kubebuilder:validation:MinProperties=1
type GatewayAPIIngressConfig struct {
	// managementMode specifies how the Cluster Ingress
	// Operator manages Gateway API Custom Resource Definitions
	// (CRDs) and the associated Gateway controller stack.
	//
	// When set to "Managed" (the default), CIO installs, owns,
	// and upgrades the Gateway API CRDs, protects them with a
	// Validating Admission Policy, and deploys the full Gateway
	// controller stack (the Istio instance deployed by CIO,
	// GatewayClass, Gateway resources). This is the only fully
	// supported configuration.
	//
	// When set to "Unmanaged", CIO does not install or manage
	// Gateway API CRDs and does not deploy the Gateway controller
	// stack. The cluster administrator or a third-party product
	// is responsible for providing their own CRDs and Gateway
	// controller. CIO reports observational status only. This
	// mode also serves as a signal to layered products that the
	// installed CRDs may not be the ones supported by the
	// OpenShift Gateway API implementation.
	//
	// When omitted, the field defaults to "Managed".
	//
	// +default="Managed"
	// +optional
	ManagementMode GatewayAPIManagementMode `json:"managementMode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressList is a collection of Ingresses.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type IngressList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	metav1.ListMeta `json:"metadata"`

	// items is a list of Ingresses.
	// +optional
	Items []Ingress `json:"items,omitempty"`
}
