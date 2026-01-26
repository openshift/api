package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PKICertificateDefinition registers certificates managed by a component,
// enabling dynamic validation of certificate names in PKI overrides.
// Components create PKICertificateDefinition resources to declare which
// certificates they manage, allowing administrators to configure those
// certificates via the PKI resource.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pkicertificatedefinitions,scope=Namespaced
// +kubebuilder:subresource:status
// +kubebuilder:validation:XValidation:rule="self.metadata.namespace == 'openshift-config'",message="pkicertificatedefinitions must be created in the openshift-config namespace"
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2645
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=ConfigurablePKI
// +openshift:compatibility-gen:level=4
type PKICertificateDefinition struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds the certificate registration specification
	// +required
	Spec PKICertificateDefinitionSpec `json:"spec,omitzero"`

	// status holds observed state
	// +optional
	Status PKICertificateDefinitionStatus `json:"status,omitempty"`
}

// PKICertificateDefinitionSpec defines certificates managed by a component.
type PKICertificateDefinitionSpec struct {
	// component identifies the operator or component managing these certificates.
	// This should typically be the name of the operator (e.g., "etcd-operator", "kube-apiserver-operator").
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Component string `json:"component,omitempty"`

	// certificates is a list of certificate definitions managed by this component.
	// Each certificate must have a unique name within the cluster.
	// +required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=256
	// +listType=map
	// +listMapKey=name
	Certificates []CertificateDefinition `json:"certificates,omitempty"`
}

// CertificateDefinition describes a single certificate managed by a component.
// +kubebuilder:validation:XValidation:rule="self.name.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?$')",message="name must be a valid DNS subdomain (lowercase alphanumeric with hyphens)"
type CertificateDefinition struct {
	// name is the unique identifier for this certificate.
	// This name is used in PKI.spec.overrides[].certificateName to configure this certificate.
	// Must be a valid DNS subdomain (lowercase letters, numbers, and hyphens).
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name,omitempty"`

	// category specifies the certificate category.
	// This helps administrators understand the certificate's role and select appropriate
	// cryptographic parameters.
	// +required
	Category CertificateCategory `json:"category,omitempty"`

	// description provides a human-readable explanation of this certificate's purpose.
	// Examples: "CA for etcd peer and server certificates", "Server certificate for API server localhost endpoint"
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Description string `json:"description,omitempty"`
}

// PKICertificateDefinitionStatus contains observed state of the certificate registration.
type PKICertificateDefinitionStatus struct {
	// conditions represent the latest available observations of the PKICertificateDefinition's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// registeredAt is the timestamp when this definition was first successfully validated.
	// +optional
	RegisteredAt *metav1.Time `json:"registeredAt,omitempty"`
}

// PKICertificateDefinitionList is a collection of PKICertificateDefinition resources.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4
type PKICertificateDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of PKICertificateDefinition resources
	Items []PKICertificateDefinition `json:"items"`
}
