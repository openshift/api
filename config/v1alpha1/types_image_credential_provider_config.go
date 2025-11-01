package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRIOCredentialProviderConfig holds cluster-wide configurations for CRI-O credential provider. CRI-O credential provider is a binary shipped with CRI-O that provides a way to obtain container image pull credentials from external sources.
// For example, it can be used to fetch mirror registry credentials from secrets resources in the cluster within the same namespace the pod will be running in.
// CRIOCredentialProviderConfig configuration specifies the pod image sources registries that should trigger the CRI-O credential provider execution, which will resolve the CRI-O mirror configurations and obtain the necessary credentials for pod creation.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=criocredentialproviderconfigs,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/1929
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=CRIOCredentialProviderConfig
// +openshift:compatibility-gen:level=4
type CRIOCredentialProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// spec defines the desired configuration of the CRIO Credential Provider.
	// +required
	Spec CRIOCredentialProviderConfigSpec `json:"spec,omitzero"`

	// status represents the current state of the CRIOCredentialProviderConfig.
	// +optional
	Status *CRIOCredentialProviderConfigStatus `json:"status,omitempty"`
}

// CRIOCredentialProviderConfigSpec defines the desired configuration of the CRI-O Credential Provider.
type CRIOCredentialProviderConfigSpec struct {
	// matchImages is a required list of string patterns used to determine whether
	// the CRI-O credential provider should be invoked for a given image. This list is
	// passed to the kubelet CredentialProviderConfig, and if any pattern matches
	// the requested image, CRI-O credential provider will be invoked to obtain credentials for pulling
	// that image or its mirrors.
	//
	// For more details, see:
	// - https://kubernetes.io/docs/tasks/administer-cluster/kubelet-credential-provider/
	// - https://github.com/cri-o/crio-credential-provider#architecture
	//
	// Each entry in matchImages is a pattern which can optionally contain a port and a path.
	// Wildcards ('*') are supported for full subdomain labels, such as '*.k8s.io' or 'k8s.*.io',
	// and for top-level domains, such as 'k8s.*' (which matches 'k8s.io' or 'k8s.net').
	// Wildcards are not allowed in the port or path, nor may they appear in the middle of a hostname label.
	// For example, '*.example.com' is valid, but 'example*.*.com' is not.
	// Each wildcard matches only a single domain label,
	// so '*.io' does **not** match '*.k8s.io'.
	//
	// A match exists between an image and a matchImage when all of the below are true:
	// - Both contain the same number of domain parts and each part matches.
	// - The URL path of an matchImages must be a prefix of the target image URL path.
	// - If the matchImages contains a port, then the port must match in the image as well.
	//
	// Example values of matchImages:
	// - 123456789.dkr.ecr.us-east-1.amazonaws.com
	// - *.azurecr.io
	// - gcr.io
	// - *.*.registry.io
	// - registry.io:8080/path
	//
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	// +required
	MatchImages []MatchImage `json:"matchImages,omitempty"`
}

// +kubebuilder:validation:MaxLength=512
// +kubebuilder:validation:XValidation:rule=`self.matches('^((\\*|[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)(\\.(\\*|[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?))*)(:[0-9]+)?(/[-a-zA-Z0-9_/]*)?$')`,message="invalid matchImages value, must be a valid fully qualified domain name with optional wildcard, port, and path"
type MatchImage string

// +k8s:deepcopy-gen=true
// CRIOCredentialProviderConfigStatus defines the observed state of CRIOCredentialProviderConfig
type CRIOCredentialProviderConfigStatus struct {
	// conditions represent the latest available observations of the configuration state
	// +optional
	// +kubebuilder:validation:MaxItems=4
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRIOCredentialProviderConfigList contains a list of CRIOCredentialProviderConfig resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type CRIOCredentialProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	Items []CRIOCredentialProviderConfig `json:"items"`
}

const (
	// ConditionTypeValidated indicates whether the configuration is failed, or partially valid
	ConditionTypeValidated = "Validated"

	// ReasonValidationFailed indicates the MatchImages configuration contains invalid patterns
	ReasonValidationFailed = "ValidationFailed"

	// ReasonConfigurationPartiallyApplied indicates some matchImage entries were ignored due to conflicts
	ReasonConfigurationPartiallyApplied = "ConfigurationPartiallyApplied"
)
