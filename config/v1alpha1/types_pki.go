package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PKI configures cryptographic parameters for certificates generated
// internally by OpenShift components.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pkis,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2645
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=ConfigurablePKI
// +openshift:compatibility-gen:level=4
type PKI struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec PKISpec `json:"spec,omitzero"`
}

// PKISpec holds the specification for PKI configuration.
type PKISpec struct {
	// certificateManagement specifies how PKI configuration is managed for internally-generated certificates.
	// This controls the certificate generation approach for all OpenShift components that create
	// certificates internally, including certificate authorities, serving certificates, and client certificates.
	//
	// +required
	CertificateManagement PKICertificateManagement `json:"certificateManagement,omitzero"`
}

// PKICertificateManagement determines whether components use hardcoded defaults (Unmanaged), follow
// OpenShift best practices (Default), or use administrator-specified cryptographic parameters (Custom).
// This provides flexibility for organizations with specific compliance requirements or security policies
// while maintaining backwards compatibility for existing clusters.
//
// +kubebuilder:validation:XValidation:rule="self.mode == 'Custom' ? has(self.custom) : !has(self.custom)",message="custom is required when mode is Custom, and forbidden otherwise"
// +union
type PKICertificateManagement struct {
	// mode determines how PKI configuration is managed.
	// Valid values are "Unmanaged", "Default", and "Custom".
	//
	// When set to Unmanaged, components use their existing hardcoded certificate
	// generation behavior, exactly as if this feature did not exist. Each component
	// generates certificates using whatever parameters it was using before this
	// feature. While most components use RSA 2048, some may use different
	// parameters. Use of this mode might prevent upgrading to the next major
	// OpenShift release.
	//
	// When set to Default, OpenShift-recommended best practices for certificate
	// generation are applied. The specific parameters may evolve across OpenShift
	// releases to adopt improved cryptographic standards. In the initial release,
	// this matches Unmanaged behavior for each component. In future releases, this
	// may adopt ECDSA or larger RSA keys based on industry best practices.
	// Recommended for most customers who want to benefit from security improvements
	// automatically.
	//
	// When set to Custom, the certificate management parameters can be set
	// explicitly. Use the custom field to specify certificate generation parameters.
	//
	// +required
	// +unionDiscriminator
	Mode PKICertificateManagementMode `json:"mode,omitempty"`

	// custom contains administrator-specified cryptographic configuration.
	// Use the defaults and categories fields to specify certificate generation parameters.
	// Required when mode is Custom, and forbidden otherwise.
	//
	// +optional
	// +unionMember
	Custom CustomPKIPolicy `json:"custom,omitzero"`
}

// CustomPKIPolicy contains administrator-specified cryptographic configuration.
// Administrators can specify defaults for all certificates or configure specific categories
// (SignerCertificate, ServingCertificate, ClientCertificate).
type CustomPKIPolicy struct {
	PKIProfile `json:",inline"`
}

// +kubebuilder:validation:Enum=Unmanaged;Default;Custom
type PKICertificateManagementMode string

const (
	// PKICertificateManagementModeUnmanaged uses hardcoded defaults (RSA 2048) for all certificates.
	// Behavior is frozen and will never change across OpenShift releases.
	PKICertificateManagementModeUnmanaged PKICertificateManagementMode = "Unmanaged"

	// PKICertificateManagementModeDefault uses OpenShift-recommended best practices.
	// Specific parameters may evolve across OpenShift releases.
	PKICertificateManagementModeDefault PKICertificateManagementMode = "Default"

	// PKICertificateManagementModeCustom uses administrator-specified configuration.
	PKICertificateManagementModeCustom PKICertificateManagementMode = "Custom"
)

// PKIProfile defines the certificate generation parameters that OpenShift components use
// to create certificates. Configuration can be specified at two hierarchical levels:
// defaults apply to all certificates and categories apply to certificate types (SignerCertificate,
// ServingCertificate, ClientCertificate).
// Category configuration takes precedence over defaults.
// +kubebuilder:validation:MinProperties=1
type PKIProfile struct {
	// defaults specifies the default certificate configuration
	// for all certificates unless overridden by category or specific
	// certificate configuration.
	// If not specified, uses platform defaults (typically RSA 2048).
	//
	// +optional
	Defaults CertificateConfig `json:"defaults,omitzero"`

	// categories allows configuration of certificate parameters
	// for categories of certificates (SignerCertificate, ServingCertificate, ClientCertificate).
	// Category configuration takes precedence over defaults.
	//
	// +optional
	// +listType=map
	// +listMapKey=category
	// +kubebuilder:validation:MaxItems=3
	Categories []CategoryCertificateConfig `json:"categories,omitempty"`
}

// CertificateConfig specifies configuration parameters for certificates.
// +kubebuilder:validation:MinProperties=1
type CertificateConfig struct {
	// key specifies the cryptographic parameters for the certificate's key pair.
	// +optional
	Key KeyConfig `json:"key,omitempty,omitzero"`

	// Future extensibility: fields like Lifetime, Rotation, Extensions
	// can be added here without restructuring the API.
}

// KeyConfig specifies cryptographic parameters for key generation.
//
// +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'RSA' ?  has(self.rsa) : !has(self.rsa)",message="rsa is required when algorithm is RSA, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'ECDSA' ?  has(self.ecdsa) : !has(self.ecdsa)",message="ecdsa is required when algorithm is ECDSA, and forbidden otherwise"
// +union
type KeyConfig struct {
	// algorithm specifies the key generation algorithm.
	// Valid values are "RSA" and "ECDSA".
	// +required
	// +unionDiscriminator
	Algorithm KeyAlgorithm `json:"algorithm,omitempty"`

	// rsa specifies RSA key parameters.
	// Required when algorithm is RSA, and forbidden otherwise.
	// +optional
	// +unionMember
	RSA RSAKeyConfig `json:"rsa,omitzero"`

	// ecdsa specifies ECDSA key parameters.
	// Required when algorithm is ECDSA, and forbidden otherwise.
	// +optional
	// +unionMember
	ECDSA ECDSAKeyConfig `json:"ecdsa,omitzero"`
}

// RSAKeyConfig specifies parameters for RSA key generation.
type RSAKeyConfig struct {
	// keySize specifies the size of RSA keys in bits.
	// Valid values are multiples of 1024 from 2048 to 8192.
	// +required
	// +kubebuilder:validation:Minimum=2048
	// +kubebuilder:validation:Maximum=8192
	// +kubebuilder:validation:MultipleOf=1024
	// +kubebuilder:default=2048
	KeySize int32 `json:"keySize,omitempty"`
}

// ECDSAKeyConfig specifies parameters for ECDSA key generation.
type ECDSAKeyConfig struct {
	// curve specifies the elliptic curve for ECDSA keys.
	// Valid values are "P256", "P384", and "P521".
	// +required
	Curve ECDSACurve `json:"curve,omitempty"`
}

type CategoryCertificateConfig struct {
	// category identifies the certificate category.
	// Valid values are "SignerCertificate", "ServingCertificate", and "ClientCertificate".
	//
	// When set to SignerCertificate, the configuration applies to certificate authority (CA) certificates
	// that sign other certificates.
	//
	// When set to ServingCertificate, the configuration applies to TLS server certificates
	// used to serve HTTPS endpoints.
	//
	// When set to ClientCertificate, the configuration applies to client authentication certificates
	// used to authenticate to servers.
	//
	// +required
	Category CertificateCategory `json:"category,omitempty"`

	// certificate specifies the configuration for this category
	// +required
	Certificate CertificateConfig `json:"certificate,omitzero"`
}

// +kubebuilder:validation:Enum=RSA;ECDSA
type KeyAlgorithm string

const (
	KeyAlgorithmRSA   KeyAlgorithm = "RSA"
	KeyAlgorithmECDSA KeyAlgorithm = "ECDSA"
)

// +kubebuilder:validation:Enum=P256;P384;P521
type ECDSACurve string

const (
	ECDSACurveP256 ECDSACurve = "P256"
	ECDSACurveP384 ECDSACurve = "P384"
	ECDSACurveP521 ECDSACurve = "P521"
)

// +kubebuilder:validation:Enum=SignerCertificate;ServingCertificate;ClientCertificate
type CertificateCategory string

const (
	CertificateCategorySignerCertificate  CertificateCategory = "SignerCertificate"
	CertificateCategoryServingCertificate CertificateCategory = "ServingCertificate"
	CertificateCategoryClientCertificate  CertificateCategory = "ClientCertificate"
)

// PKIList is a collection of PKI resources.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4
type PKIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of PKI resources
	Items []PKI `json:"items"`
}
