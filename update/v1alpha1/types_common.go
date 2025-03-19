package v1alpha1

// ResourceRef is a reference to a kubernetes resource, typically involved in an insight
type ResourceRef struct {
	// group of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Group string `json:"group,omitempty"`

	// resource of object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Resource string `json:"resource"`

	// name of the object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Name string `json:"name"`

	// namespace of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Namespace string `json:"namespace,omitempty"`
}

// ScopeType is one of ControlPlane or WorkerPool
// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
type ScopeType string

const (
	// ControlPlane is used for insights that are related to the control plane (including control plane pool or nodes)
	ControlPlaneScope ScopeType = "ControlPlane"
	// WorkerPool is used for insights that are related to a worker pools and nodes (excluding control plane)
	WorkerPoolScope ScopeType = "WorkerPool"
)
