package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Console holds cluster-wide configuration for the web console, including the
// logout URL, and reports the public URL of the console. The canonical name is
// `cluster`.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/470
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=consoles,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:metadata:annotations=release.openshift.io/bootstrap-required=true
type Console struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec ConsoleSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ConsoleStatus `json:"status"`
}

// ConsoleSpec is the specification of the desired behavior of the Console.
type ConsoleSpec struct {
	// authentication configures console authentication behavior.
	// When omitted, default authentication settings are used.
	// +optional
	Authentication ConsoleAuthentication `json:"authentication"`

	// externalSecretStore configures integration with an external secret store
	// for console secret management. When omitted, no external secret store is configured.
	// +optional
	// +openshift:enable:FeatureGate=ExternalSecretStore
	ExternalSecretStore ExternalSecretStoreConfig `json:"externalSecretStore,omitempty,omitzero"`
}

// ConsoleStatus defines the observed status of the Console.
type ConsoleStatus struct {
	// The URL for the console. This will be derived from the host for the route that
	// is created for the console.
	// +optional
	ConsoleURL string `json:"consoleURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ConsoleList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	Items []Console `json:"items"`
}

// ConsoleAuthentication defines a list of optional configuration for console authentication.
type ConsoleAuthentication struct {
	// An optional, absolute URL to redirect web browsers to after logging out of
	// the console. If not specified, it will redirect to the default login page.
	// This is required when using an identity provider that supports single
	// sign-on (SSO) such as:
	// - OpenID (Keycloak, Azure)
	// - RequestHeader (GSSAPI, SSPI, SAML)
	// - OAuth (GitHub, GitLab, Google)
	// Logging out of the console will destroy the user's token. The logoutRedirect
	// provides the user the option to perform single logout (SLO) through the identity
	// provider to destroy their single sign-on session.
	// +optional
	// +kubebuilder:validation:Pattern=`^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$`
	LogoutRedirect string `json:"logoutRedirect,omitempty"`
}

// ExternalSecretStoreType defines the type of external secret store.
// When set to Vault, HashiCorp Vault is used as the external secret store.
// +kubebuilder:validation:Enum=Vault
type ExternalSecretStoreType string

// ExternalSecretStoreConfig defines the configuration for integration with an
// external secret store.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Vault' ? has(self.vault) : !has(self.vault)",message="vault configuration is required when type is Vault, and forbidden otherwise"
type ExternalSecretStoreConfig struct {
	// type specifies the type of external secret store to use.
	// Currently supported values:
	// - Vault
	// +required
	Type ExternalSecretStoreType `json:"type"`

	// vault contains the configuration for a HashiCorp Vault secret store.
	// This field is required when type is Vault, and forbidden otherwise.
	// +optional
	Vault VaultSecretStoreConfig `json:"vault,omitempty,omitzero"`
}

// VaultSecretStoreConfig defines the configuration for HashiCorp Vault integration.
type VaultSecretStoreConfig struct {
	// serverAddress specifies the address of the Vault server.
	// Must be a valid URL starting with https://.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	ServerAddress string `json:"serverAddress"`

	// transitKeyName specifies the name of the transit encryption key in Vault.
	// The name must not exceed 253 characters.
	// +required
	// +kubebuilder:validation:MinLength=1
	TransitKeyName string `json:"transitKeyName"`

	// transitMountPath specifies the mount path for the transit secrets engine.
	// +required
	// +kubebuilder:validation:MinLength=1
	TransitMountPath string `json:"transitMountPath"`

	// namespace specifies the Vault namespace to use.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Namespace string `json:"namespace,omitempty"`

	// caCertificate contains the PEM-encoded CA certificate for TLS verification
	// of the Vault server connection.
	// When omitted, the system's trusted CA certificates are used.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=65536
	CACertificate string `json:"caCertificate,omitempty"`

	// authSecret references a secret containing the authentication credentials
	// for Vault. The secret must exist in the openshift-config namespace.
	// +required
	AuthSecret SecretNameReference `json:"authSecret"`

	// refreshInterval specifies how often secrets are re-fetched from Vault,
	// in seconds. The value must be between 30 and 3600.
	// When omitted, the platform chooses a reasonable default.
	// +optional
	// +kubebuilder:validation:Minimum=30
	// +kubebuilder:validation:Maximum=3600
	RefreshInterval int32 `json:"refreshInterval,omitempty"`
}
