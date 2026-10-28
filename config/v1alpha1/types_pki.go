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
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="pki is a singleton, .metadata.name must be 'cluster'"
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
	// certificates specifies the cluster's certificate policy for internally-generated
	// certificates, resolved most-specific-wins: a named certificate override, then a
	// category override (CA, server, client), then the required defaults.
	// This field is required so that every PKI resource specifies a certificate policy.
	// This field names a configuration policy; it is unrelated to the X.509
	// certificate-policies extension.
	//
	// +required
	Certificates PKICertificatePolicy `json:"certificates,omitzero"`

	// keys optionally configures keys generated outside of any certificate, such as
	// the bound service-account token signing key. Entries are named and opt-in with no
	// inherited default. When omitted, no keys are configured by this resource.
	// When set, at least one named key must be specified; an empty object is rejected.
	// Unlike certificates, keys have no resolution chain: a key is either configured
	// here or it is not, and a component whose key is not configured decides for
	// itself (an internal default, or an error if it has none).
	//
	// +optional
	Keys PKIKeyPolicy `json:"keys,omitempty,omitzero"`

	// --- TOMBSTONE ---
	// certificateManagement was replaced by certificates (a net-new PKICertificatePolicy
	// type; the section declares a policy, and "management" named an activity rather than
	// the artifact being configured). See the PKICertificateManagement tombstone below for
	// the former structure of the fields it configured.
	// The field name "certificateManagement" and its json tag are reserved to prevent
	// reuse with a different backing type.
	//
	// +optional
	// CertificateManagement PKICertificateManagement `json:"certificateManagement,omitzero"`
}

// PKICertificatePolicy is the cluster's policy for internally-generated certificates:
// a required defaults, optional per-category overrides, and optional overrides for
// specific named certificates. The most-specific setting wins.
type PKICertificatePolicy struct {
	// defaults specifies the default certificate configuration applied to any certificate
	// that has no more-specific (category or named) override. This field is required so
	// that every certificate has a well-defined configuration.
	//
	// +required
	Defaults PKIDefaultCertificateConfig `json:"defaults,omitzero"`

	// caCertificates optionally overrides certificate parameters for certificate
	// authority (CA) certificates that sign other certificates. When set, these
	// parameters take precedence over defaults for all CA certificates. When omitted,
	// defaults are used for CA certificates.
	//
	// +optional
	CACertificates PKICertificateConfig `json:"caCertificates,omitempty,omitzero"`

	// serverCertificates optionally overrides certificate parameters for TLS server
	// certificates. When set, these parameters take precedence over defaults for all
	// server certificates. When omitted, defaults are used for server certificates.
	//
	// +optional
	ServerCertificates PKICertificateConfig `json:"serverCertificates,omitempty,omitzero"`

	// clientCertificates optionally overrides certificate parameters for client
	// authentication certificates. When set, these parameters take precedence over
	// defaults for all client certificates. When omitted, defaults are used for client
	// certificates.
	//
	// +optional
	ClientCertificates PKICertificateConfig `json:"clientCertificates,omitempty,omitzero"`

	// namedCertificates optionally overrides certificate parameters for specific
	// certificates identified by their well-known name. These overrides have the highest
	// precedence (named, then category, then defaults). Entries are keyed by name, so a
	// given name may appear at most once. When omitted, no named overrides are applied.
	// When present, the list must contain at least one entry. At most 512 named
	// certificate overrides may be specified.
	//
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=512
	// +optional
	NamedCertificates []PKINamedCertificateConfig `json:"namedCertificates,omitempty"`
}

// PKINamedCertificateConfig overrides certificate parameters for a specific certificate
// identified by its well-known name.
type PKINamedCertificateConfig struct {
	// name is the well-known override name of the targeted certificate. It is a DNS
	// subdomain prefix, a "/", then a name made of alphanumerics, "-", "_" and "."
	// (for example "openshift.io/service-ca.service-serving-signer"). The prefix is at
	// most 253 characters and the name after the "/" at most 63, so the full name
	// must be between 1 and 317 characters. The name is validated
	// for format only; there is no registry of valid names, because certificates can be
	// registered dynamically. Core-platform certificates use the bare "openshift.io/"
	// prefix; names outside the core platform use an owner subdomain prefix (for example
	// "hypershift.openshift.io/").
	// This field is required so that every named certificate override identifies the
	// certificate it applies to.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=317
	// +kubebuilder:validation:XValidation:rule="self.contains('/') && !format.qualifiedName().validate(self).hasValue()",message="name must be a prefixed qualified name, for example openshift.io/etcd.signer"
	Name string `json:"name,omitempty"`

	// PKICertificateConfig holds the overridden parameters for this certificate.
	PKICertificateConfig `json:",inline"`
}

// PKIKeyPolicy configures keys generated outside of any certificate, by well-known name.
// When set, at least one named key must be specified; an empty object is rejected.
type PKIKeyPolicy struct {
	// namedKeys configures keys generated outside of any certificate, such as the bound
	// service-account token signing key, identified by their well-known name. Entries are
	// opt-in with no inherited default and are keyed by name, so a given name may appear at
	// most once. This field is required so that a keys object is never empty.
	// At most 16 named keys may be specified.
	//
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +required
	NamedKeys []PKINamedKeyConfig `json:"namedKeys,omitempty,omitzero"`
}

// PKINamedKeyConfig configures the cryptographic parameters for a single key generated
// outside of any certificate, identified by its well-known name.
type PKINamedKeyConfig struct {
	// name is the well-known name of the key. It is a DNS subdomain prefix, a
	// "/", then a name made of alphanumerics, "-", "_" and "." (for example
	// "openshift.io/kube-apiserver.bound-service-account-token-signer"). The prefix is at
	// most 253 characters and the name after the "/" at most 63, so the full name
	// must be between 1 and 317 characters. The name is validated
	// for format only, using the same rules as a named certificate override.
	// This field is required so that every named key entry identifies the key it configures.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=317
	// +kubebuilder:validation:XValidation:rule="self.contains('/') && !format.qualifiedName().validate(self).hasValue()",message="name must be a prefixed qualified name, for example openshift.io/etcd.signer"
	Name string `json:"name,omitempty"`

	// key specifies the cryptographic parameters for this key. It is required
	// because keys configured here have no inherited default.
	//
	// +required
	Key PKIKeyConfig `json:"key,omitzero"`
}

// PKIDefaultCertificateConfig specifies the default certificate configuration
// parameters. All fields are required to ensure that defaults are fully
// specified for all certificates.
type PKIDefaultCertificateConfig struct {
	// key specifies the cryptographic parameters for the certificate's key pair.
	// This field is required in defaults to ensure all certificates have a
	// well-defined key configuration.
	// +required
	Key PKIKeyConfig `json:"key,omitzero"`
}

// PKICertificateConfig specifies configuration parameters for certificates.
type PKICertificateConfig struct {
	// key specifies the cryptographic parameters for the certificate's key pair.
	// +required
	Key PKIKeyConfig `json:"key,omitzero"`
}

// PKIKeyConfig specifies cryptographic parameters for key generation.
//
// +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'RSA' ?  has(self.rsa) : !has(self.rsa)",message="rsa is required when algorithm is RSA, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'ECDSA' ?  has(self.ecdsa) : !has(self.ecdsa)",message="ecdsa is required when algorithm is ECDSA, and forbidden otherwise"
// +union
type PKIKeyConfig struct {
	// algorithm specifies the key generation algorithm.
	// Valid values are "RSA" and "ECDSA".
	//
	// When set to RSA, the rsa field must be specified and the generated key
	// will be an RSA key with the configured key size.
	//
	// When set to ECDSA, the ecdsa field must be specified and the generated key
	// will be an ECDSA key using the configured elliptic curve.
	//
	// +required
	// +unionDiscriminator
	Algorithm PKIKeyAlgorithm `json:"algorithm,omitempty"`

	// rsa specifies RSA key parameters.
	// Required when algorithm is RSA, and forbidden otherwise.
	// +optional
	// +unionMember
	RSA PKIRSAKeyConfig `json:"rsa,omitzero"`

	// ecdsa specifies ECDSA key parameters.
	// Required when algorithm is ECDSA, and forbidden otherwise.
	// +optional
	// +unionMember
	ECDSA PKIECDSAKeyConfig `json:"ecdsa,omitzero"`
}

// PKIRSAKeyConfig specifies parameters for RSA key generation.
type PKIRSAKeyConfig struct {
	// keySize specifies the size of RSA keys in bits.
	// Valid values are multiples of 1024 from 2048 to 4096.
	// +required
	// +kubebuilder:validation:Minimum=2048
	// +kubebuilder:validation:Maximum=4096
	// +kubebuilder:validation:MultipleOf=1024
	KeySize int32 `json:"keySize,omitempty"`
}

// PKIECDSAKeyConfig specifies parameters for ECDSA key generation.
type PKIECDSAKeyConfig struct {
	// curve specifies the NIST elliptic curve for ECDSA keys.
	// Valid values are "P256", "P384", and "P521".
	//
	// When set to P256, the NIST P-256 curve (also known as secp256r1) is used,
	// providing 128-bit security.
	//
	// When set to P384, the NIST P-384 curve (also known as secp384r1) is used,
	// providing 192-bit security.
	//
	// When set to P521, the NIST P-521 curve (also known as secp521r1) is used,
	// providing 256-bit security.
	//
	// +required
	Curve PKIECDSACurve `json:"curve,omitempty"`
}

// PKIKeyAlgorithm specifies the cryptographic algorithm used for key generation.
//
// +kubebuilder:validation:Enum=RSA;ECDSA
type PKIKeyAlgorithm string

const (
	// PKIKeyAlgorithmRSA specifies the RSA (Rivest-Shamir-Adleman) algorithm for key generation.
	PKIKeyAlgorithmRSA PKIKeyAlgorithm = "RSA"

	// PKIKeyAlgorithmECDSA specifies the ECDSA (Elliptic Curve Digital Signature Algorithm) for key generation.
	PKIKeyAlgorithmECDSA PKIKeyAlgorithm = "ECDSA"
)

// PKIECDSACurve specifies the elliptic curve used for ECDSA key generation.
//
// +kubebuilder:validation:Enum=P256;P384;P521
type PKIECDSACurve string

const (
	// PKIECDSACurveP256 specifies the NIST P-256 curve (also known as secp256r1), providing 128-bit security.
	PKIECDSACurveP256 PKIECDSACurve = "P256"

	// PKIECDSACurveP384 specifies the NIST P-384 curve (also known as secp384r1), providing 192-bit security.
	PKIECDSACurveP384 PKIECDSACurve = "P384"

	// PKIECDSACurveP521 specifies the NIST P-521 curve (also known as secp521r1), providing 256-bit security.
	PKIECDSACurveP521 PKIECDSACurve = "P521"
)

// Start: TOMBSTONE
//
// The following types made up the former PKICertificateManagement chain that the
// spec.certificateManagement tombstone above pointed to. They are preserved here,
// commented out in their entirety, so the exact former serialization and validation
// can be followed and revived if needed. Note the validation differences relative
// to their replacements: RSAKeyConfig allowed keySize up to 8192 (PKIRSAKeyConfig
// caps at 4096), and CertificateConfig enforced its single key property via
// MinProperties=1 (PKICertificateConfig marks key required).
//
// PKICertificateManagement determined whether components use hardcoded defaults
// (Unmanaged), follow OpenShift best practices (Default), or use
// administrator-specified cryptographic parameters (Custom). Removed; replaced by
// the net-new PKICertificatePolicy configured via the certificates field.
//
// // +kubebuilder:validation:XValidation:rule="self.mode == 'Custom' ? has(self.custom) : !has(self.custom)",message="custom is required when mode is Custom, and forbidden otherwise"
// // +union
// type PKICertificateManagement struct {
// 	// mode determined how PKI configuration was managed.
// 	// Valid values were "Unmanaged", "Default", and "Custom".
// 	//
// 	// // +required
// 	// // +unionDiscriminator
// 	Mode PKICertificateManagementMode `json:"mode,omitempty"`
//
// 	// custom contained administrator-specified cryptographic configuration.
// 	// Required when mode was Custom, and forbidden otherwise.
// 	//
// 	// // +optional
// 	// // +unionMember
// 	Custom CustomPKIPolicy `json:"custom,omitzero"`
// }
//
// CustomPKIPolicy contained administrator-specified cryptographic configuration.
// Removed; its wrapped fields were flattened into PKICertificatePolicy.
//
// // +kubebuilder:validation:MinProperties=1
// type CustomPKIPolicy struct {
// 	PKIProfile `json:",inline"`
// }
//
// PKICertificateManagementMode specified the mode for PKI certificate management.
// Removed along with the mode union.
//
// // +kubebuilder:validation:Enum=Unmanaged;Default;Custom
// type PKICertificateManagementMode string
//
// const (
// 	// PKICertificateManagementModeUnmanaged used each component's existing hardcoded defaults.
// 	PKICertificateManagementModeUnmanaged PKICertificateManagementMode = "Unmanaged"
//
// 	// PKICertificateManagementModeDefault used OpenShift-recommended best practices.
// 	PKICertificateManagementModeDefault PKICertificateManagementMode = "Default"
//
// 	// PKICertificateManagementModeCustom used administrator-specified configuration.
// 	PKICertificateManagementModeCustom PKICertificateManagementMode = "Custom"
// )
//
// PKIProfile defined the certificate generation parameters that OpenShift components
// used to create certificates. Removed; its fields were flattened into
// PKICertificatePolicy (signerCertificates -> caCertificates,
// servingCertificates -> serverCertificates).
//
// type PKIProfile struct {
// 	// defaults specified the default certificate configuration that applied
// 	// to all certificates unless overridden by a category override.
// 	//
// 	// // +required
// 	Defaults DefaultCertificateConfig `json:"defaults,omitzero"`
//
// 	// signerCertificates optionally overrode certificate parameters for
// 	// certificate authority (CA) certificates that sign other certificates.
// 	// Renamed to caCertificates on PKICertificatePolicy.
// 	//
// 	// // +optional
// 	SignerCertificates CertificateConfig `json:"signerCertificates,omitempty,omitzero"`
//
// 	// servingCertificates optionally overrode certificate parameters for
// 	// TLS server certificates used to serve HTTPS endpoints.
// 	// Renamed to serverCertificates on PKICertificatePolicy.
// 	//
// 	// // +optional
// 	ServingCertificates CertificateConfig `json:"servingCertificates,omitempty,omitzero"`
//
// 	// clientCertificates optionally overrode certificate parameters for
// 	// client authentication certificates used to authenticate to servers.
// 	// Flattened into clientCertificates on PKICertificatePolicy.
// 	//
// 	// // +optional
// 	ClientCertificates CertificateConfig `json:"clientCertificates,omitempty,omitzero"`
// }
//
// DefaultCertificateConfig specified the default certificate configuration
// parameters. Renamed to PKIDefaultCertificateConfig.
//
// type DefaultCertificateConfig struct {
// 	// key specified the cryptographic parameters for the certificate's key pair.
// 	//
// 	// // +required
// 	Key KeyConfig `json:"key,omitzero"`
// }
//
// CertificateConfig specified configuration parameters for certificates.
// Renamed to PKICertificateConfig.
//
// // +kubebuilder:validation:MinProperties=1
// type CertificateConfig struct {
// 	// key specified the cryptographic parameters for the certificate's key pair.
// 	//
// 	// // +optional
// 	Key KeyConfig `json:"key,omitzero"`
// }
//
// KeyConfig specified cryptographic parameters for key generation.
// Renamed to PKIKeyConfig.
//
// // +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'RSA' ?  has(self.rsa) : !has(self.rsa)",message="rsa is required when algorithm is RSA, and forbidden otherwise"
// // +kubebuilder:validation:XValidation:rule="has(self.algorithm) && self.algorithm == 'ECDSA' ?  has(self.ecdsa) : !has(self.ecdsa)",message="ecdsa is required when algorithm is ECDSA, and forbidden otherwise"
// // +union
// type KeyConfig struct {
// 	// algorithm specified the key generation algorithm.
// 	// Valid values were "RSA" and "ECDSA".
// 	//
// 	// // +required
// 	// // +unionDiscriminator
// 	Algorithm KeyAlgorithm `json:"algorithm,omitempty"`
//
// 	// rsa specified RSA key parameters.
// 	//
// 	// // +optional
// 	// // +unionMember
// 	RSA RSAKeyConfig `json:"rsa,omitzero"`
//
// 	// ecdsa specified ECDSA key parameters.
// 	//
// 	// // +optional
// 	// // +unionMember
// 	ECDSA ECDSAKeyConfig `json:"ecdsa,omitzero"`
// }
//
// RSAKeyConfig specified parameters for RSA key generation.
// Renamed to PKIRSAKeyConfig, which caps keySize at 4096 (was 8192).
//
// type RSAKeyConfig struct {
// 	// keySize specified the size of RSA keys in bits.
// 	// Valid values were multiples of 1024 from 2048 to 8192.
// 	//
// 	// // +required
// 	// // +kubebuilder:validation:Minimum=2048
// 	// // +kubebuilder:validation:Maximum=8192
// 	// // +kubebuilder:validation:MultipleOf=1024
// 	KeySize int32 `json:"keySize,omitempty"`
// }
//
// ECDSAKeyConfig specified parameters for ECDSA key generation.
// Renamed to PKIECDSAKeyConfig.
//
// type ECDSAKeyConfig struct {
// 	// curve specified the NIST elliptic curve for ECDSA keys.
// 	// Valid values were "P256", "P384", and "P521".
// 	//
// 	// // +required
// 	Curve ECDSACurve `json:"curve,omitempty"`
// }
//
// KeyAlgorithm specified the cryptographic algorithm used for key generation.
// Renamed to PKIKeyAlgorithm.
//
// // +kubebuilder:validation:Enum=RSA;ECDSA
// type KeyAlgorithm string
//
// const (
// 	// KeyAlgorithmRSA specified the RSA (Rivest-Shamir-Adleman) algorithm for key generation.
// 	KeyAlgorithmRSA KeyAlgorithm = "RSA"
//
// 	// KeyAlgorithmECDSA specified the ECDSA (Elliptic Curve Digital Signature Algorithm) for key generation.
// 	KeyAlgorithmECDSA KeyAlgorithm = "ECDSA"
// )
//
// ECDSACurve specified the elliptic curve used for ECDSA key generation.
// Renamed to PKIECDSACurve.
//
// // +kubebuilder:validation:Enum=P256;P384;P521
// type ECDSACurve string
//
// const (
// 	// ECDSACurveP256 specified the NIST P-256 curve (also known as secp256r1), providing 128-bit security.
// 	ECDSACurveP256 ECDSACurve = "P256"
//
// 	// ECDSACurveP384 specified the NIST P-384 curve (also known as secp384r1), providing 192-bit security.
// 	ECDSACurveP384 ECDSACurve = "P384"
//
// 	// ECDSACurveP521 specified the NIST P-521 curve (also known as secp521r1), providing 256-bit security.
// 	ECDSACurveP521 ECDSACurve = "P521"
// )
//
// End: TOMBSTONE

// PKIList is a collection of PKI resources.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4
type PKIList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of PKI resources
	Items []PKI `json:"items"`
}
