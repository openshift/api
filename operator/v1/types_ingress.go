package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ingress contains configuration options specific to the Ingress Operator itself,
// including how it manages Gateway API integration.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +openshift:file-pattern=cvoRunLevel=0000_50,operatorName=ingress,operatorOrdering=02
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ingresses,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2890
// +openshift:capability=Ingress
// +openshift:enable:FeatureGate=GatewayAPIManagementMode
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="Ingress is a singleton; the .metadata.name field must be 'cluster'"
type Ingress struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ObjectMeta `json:"metadata"`

	// spec holds user settable values for configuration.
	// +required
	Spec IngressSpec `json:"spec,omitzero"`

	// status holds observed values from the cluster.
	// +optional
	Status IngressStatus `json:"status,omitzero"`
}

// IngressSpec is the specification of the desired behavior of the Ingress Operator.
// +kubebuilder:validation:MinProperties=1
type IngressSpec struct {
	// gatewayAPI holds configuration for Gateway API integration, including how the
	// ingress operator manages Gateway API CRDs, the OpenShift Gateway API
	// implementation, and its Gateway API controllers.
	//
	// +optional
	GatewayAPI GatewayAPIIngressConfig `json:"gatewayAPI,omitzero"`

	// canary holds configuration for the ingress canary daemonset (the
	// ingress-canary DaemonSet in the openshift-ingress-canary namespace).
	//
	// canary is optional. When omitted, the ingress canary daemonset keeps its
	// default configuration. When specified, canary must set at least one of
	// its fields; an empty object is rejected.
	//
	// +optional
	// +openshift:enable:FeatureGate=IngressCanaryNodePlacement
	Canary IngressCanary `json:"canary,omitzero"`
}

// IngressStatus defines the observed status of the Ingress Operator.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.observedGeneration) || (has(self.observedGeneration) && self.observedGeneration >= oldSelf.observedGeneration)",message="observedGeneration must remain set and only increase once set"
type IngressStatus struct {
	// conditions is a list of conditions and their status.
	//
	// Gateway API CRD management conditions are reported here with the "GatewayAPI" prefix:
	//
	// * "GatewayAPICRDsManaged" indicates whether the ingress operator is actively
	//   managing Gateway API CRDs.
	// * "GatewayAPICRDsPresent" indicates whether Gateway API CRDs exist on the
	//   cluster.
	// * "GatewayAPICRDsCompliant" indicates whether the installed CRDs match the
	//   version expected by this ingress operator release.
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// observedGeneration represents the most recent generation observed by the operator and specifies the version of
	// the spec field currently being synced.
	//
	// When omitted, the operator has not yet observed the resource.
	// +optional
	// +kubebuilder:validation:Minimum=1
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// GatewayAPIIngressConfig holds configuration for Gateway API integration in the
// Cluster Ingress Operator.
// +kubebuilder:validation:MinProperties=1
type GatewayAPIIngressConfig struct {
	// managementMode specifies how the Cluster Ingress Operator manages Gateway API
	// Custom Resource Definitions (CRDs), the OpenShift Gateway API implementation,
	// and its Gateway API controllers.
	//
	// Allowed values are "Managed" and "Unmanaged".
	//
	// When omitted or set to "Managed", the ingress operator installs, owns, and
	// upgrades the Gateway API CRDs, protects them with a Validating Admission
	// Policy, and deploys the OpenShift Gateway API implementation and its Gateway
	// API controllers.
	//
	// When set to "Unmanaged", the ingress operator does not install or manage
	// Gateway API CRDs and does not deploy the OpenShift Gateway API implementation
	// or its Gateway API controllers. The cluster administrator or a third-party
	// product is responsible for providing their own CRDs and Gateway controller.
	// The ingress operator reports observational status only.
	//
	// +optional
	ManagementMode GatewayAPIManagementMode `json:"managementMode,omitempty"`
}

// GatewayAPIManagementMode describes how the Cluster Ingress Operator manages
// Gateway API Custom Resource Definitions.
// +kubebuilder:validation:Enum=Managed;Unmanaged
type GatewayAPIManagementMode string

const (
	// GatewayAPIManagementModeManaged means the ingress operator installs, owns,
	// protects (via a Validating Admission Policy), and upgrades the Gateway API
	// CRDs, deploys the OpenShift Gateway API implementation, and runs its Gateway
	// API controllers. This is the default mode and the only fully supported
	// configuration.
	GatewayAPIManagementModeManaged GatewayAPIManagementMode = "Managed"

	// GatewayAPIManagementModeUnmanaged means the ingress operator does not
	// install or manage Gateway API CRDs, does not deploy the OpenShift Gateway
	// API implementation, and does not run its Gateway API controllers. The
	// cluster administrator or a third-party product is responsible for bringing
	// their own CRDs and Gateway controller. The ingress operator reports
	// observational status only.
	GatewayAPIManagementModeUnmanaged GatewayAPIManagementMode = "Unmanaged"
)

// IngressCanary holds configuration for the ingress canary daemonset.
// +kubebuilder:validation:MinProperties=1
type IngressCanary struct {
	// nodePlacement enables explicit control over the scheduling of the ingress
	// canary daemonset pods.
	//
	// nodePlacement is optional. To leave the canary daemonset with its default
	// scheduling configuration, which tolerates the
	// node-role.kubernetes.io/infra taint, omit the parent canary field
	// entirely: because canary requires at least one field, nodePlacement can
	// only be absent when canary itself is absent. When specified,
	// nodePlacement must set at least one of its fields; an empty object is
	// rejected.
	//
	// +optional
	NodePlacement CanaryNodePlacement `json:"nodePlacement,omitzero"`
}

// CanaryNodePlacement describes node placement configuration for the ingress
// canary daemonset pods.
// +kubebuilder:validation:MinProperties=1
type CanaryNodePlacement struct {
	// tolerations is a list of tolerations that are applied to the ingress
	// canary daemonset pods in addition to the tolerations the canary daemonset
	// has by default. The default toleration for the
	// node-role.kubernetes.io/infra taint is always preserved.
	//
	// If unset, the canary daemonset has only its default tolerations. When
	// specified, the list must contain at least 1 toleration; to remove all
	// additional tolerations, omit the field rather than setting an empty list.
	//
	// Each toleration has the following constraints:
	//
	// * operator, if specified, must be "Exists", "Equal", or the empty string;
	//   an omitted or empty operator means "Equal". The "Lt" and "Gt" operators
	//   are not supported.
	// * effect, if specified, must be "NoSchedule", "PreferNoSchedule",
	//   "NoExecute", or the empty string; an omitted or empty effect means the
	//   toleration matches all taint effects.
	// * value must be empty when operator is "Exists".
	// * key must be non-empty. A toleration with an empty key is rejected; in
	//   particular, a toleration with an empty key and operator "Exists", which
	//   tolerates all taints, is rejected.
	// * tolerationSeconds may be specified only when effect is "NoExecute".
	//
	// At least 1 and no more than 32 tolerations may be specified.
	//
	// Duplicate tolerations are allowed.
	//
	// See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:XValidation:rule=`self.all(t, !has(t.operator) || t.operator in ['', 'Exists', 'Equal'])`,message="tolerations operator must be 'Exists' or 'Equal'"
	// +kubebuilder:validation:XValidation:rule=`self.all(t, !has(t.effect) || t.effect in ['', 'NoSchedule', 'PreferNoSchedule', 'NoExecute'])`,message="tolerations effect must be 'NoSchedule', 'PreferNoSchedule', or 'NoExecute'"
	// +kubebuilder:validation:XValidation:rule=`self.all(t, !has(t.operator) || t.operator != 'Exists' || !has(t.value) || t.value == '')`,message="tolerations value must be empty when operator is 'Exists'"
	// +kubebuilder:validation:XValidation:rule=`self.all(t, has(t.key) && t.key != '')`,message="tolerations key must be non-empty"
	// +kubebuilder:validation:XValidation:rule=`self.all(t, !has(t.tolerationSeconds) || (has(t.effect) && t.effect == 'NoExecute'))`,message="tolerations effect must be 'NoExecute' when tolerationSeconds is set"
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressList is a collection of Ingresses.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type IngressList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	// items is a list of Ingresses.
	// +optional
	Items []Ingress `json:"items,omitempty"`
}
