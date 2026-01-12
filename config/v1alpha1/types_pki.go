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
	// - Unmanaged: Components use their existing hardcoded certificate generation behavior, exactly as if this feature did not exist.
	//   Each component generates certificates using whatever parameters it was using before this feature.
	//   While most components use RSA 2048, some may use different parameters.
	//   Use of this mode might prevent upgrading to the next major OpenShift release.
	// + Default when upgrading from a version without this feature to ensure zero behavior change.
	//
	// - Default: Use OpenShift-recommended best practices for certificate generation.
	//   The specific parameters may evolve across OpenShift releases to adopt improved cryptographic standards.
	//   In the initial release, this matches Unmanaged behavior for each component.
	//   In future releases, this may adopt ECDSA or larger RSA keys based on industry best practices.
	//   Recommended for most customers who want to benefit from security improvements automatically.
	// + Default when installing a fresh cluster.
	//
	// - Custom: Administrator explicitly configures cryptographic parameters.
	//   Use the custom field to specify certificate generation parameters.
	// + Recommended for customers with specific compliance requirements or organizational PKI policies.
	//
	// + When upgrading from a version without this feature:
	// + - The PKI resource is created with mode: Unmanaged to ensure zero behavior change.
	//
	// + When installing a fresh cluster:
	// + - If no PKI configuration is provided in install-config.yaml, mode: Default is used.
	// + - If PKI configuration is provided in install-config.yaml, mode: Custom is used with the specified configuration.
	//
	// +required
	// +kubebuilder:validation:Enum=Unmanaged;Default;Custom
	// +unionDiscriminator
	Mode PKICertificateManagementMode `json:"mode,omitempty"`

	// custom contains administrator-specified cryptographic configuration.
	// Use the defaults, categories, and overrides fields to specify certificate generation parameters.
	// Required when mode is Custom, and forbidden otherwise.
	//
	// +optional
	// +unionMember
	Custom CustomPKIPolicy `json:"custom,omitzero"`
}

// CustomPKIPolicy contains administrator-specified cryptographic configuration.
// Administrators can specify defaults for all certificates, configure specific categories
// (SignerCertificate, ServingCertificate, ClientCertificate), or override specific named certificates.
type CustomPKIPolicy struct {
	PKIProfile `json:",inline"`
}

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
// to create certificates. Configuration can be specified at three hierarchical levels:
// defaults apply to all certificates, categories apply to certificate types (SignerCertificate,
// ServingCertificate, ClientCertificate), and overrides apply to specific named certificates.
// More specific levels take precedence over general ones.
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

	// overrides allows configuration of certificate parameters
	// for specific named certificates.
	// Override configuration takes precedence over both category
	// and default configuration.
	//
	// +optional
	// +listType=map
	// +listMapKey=certificateName
	// +kubebuilder:validation:MaxItems=256
	Overrides []CertificateOverride `json:"overrides,omitempty"`
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
	// +kubebuilder:validation:Enum=RSA;ECDSA
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
	// Valid values are 2048, 3072, and 4096.
	// +required
	// +kubebuilder:validation:Enum=2048;3072;4096
	// +kubebuilder:validation:Minimum=2048
	KeySize int32 `json:"keySize,omitempty"`
}

// ECDSAKeyConfig specifies parameters for ECDSA key generation.
type ECDSAKeyConfig struct {
	// curve specifies the elliptic curve for ECDSA keys.
	// Valid values are "P256", "P384", and "P521".
	// +required
	// +kubebuilder:validation:Enum=P256;P384;P521
	Curve ECDSACurve `json:"curve,omitempty"`
}

type CategoryCertificateConfig struct {
	// category identifies the certificate category.
	// Valid values are "SignerCertificate", "ServingCertificate", and "ClientCertificate".
	// +required
	// +kubebuilder:validation:Enum=SignerCertificate;ServingCertificate;ClientCertificate
	Category CertificateCategory `json:"category,omitempty"`

	// certificate specifies the configuration for this category
	// +required
	Certificate CertificateConfig `json:"certificate,omitzero"`
}

// CertificateOverride allows configuration of certificate parameters for specific named certificates.
//
// +kubebuilder:validation:XValidation:rule=`self.certificateName in ['admin-kubeconfig-signer','kubelet-bootstrap-kubeconfig-signer','kube-apiserver-localhost-signer','kube-apiserver-service-network-signer','kube-apiserver-lb-signer','root-ca','kube-apiserver-to-kubelet-signer','kube-control-plane-signer','aggregator-signer','kubelet-signer','aggregator-ca','etcd-signer','etcd-metrics-signer','service-ca','csr-signer-ca','kube-apiserver-localhost-server','kube-apiserver-service-network-server','kube-apiserver-lb-server','kube-apiserver-internal-lb-server','machine-config-server','ironic-server','etcd-peer-server','etcd-server','etcd-metrics-server','admin-kubeconfig-client','kubelet-client','kube-apiserver-to-kubelet-client','kube-control-plane-kube-controller-manager-client','kube-control-plane-kube-scheduler-client','aggregator-client','apiserver-proxy-client','journal-gatewayd-client']`,message="certificateName must be a well-known certificate name"
type CertificateOverride struct {
	// certificateName identifies a specific certificate to configure.
	// The name must match a well-known certificate name in the cluster.
	// Examples: "kube-apiserver-to-kubelet-signer", "kube-apiserver-localhost-server",
	// "admin-kubeconfig-client", "etcd-signer", "service-ca"
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	CertificateName string `json:"certificateName,omitempty"`

	// certificate specifies the configuration for this certificate
	// +required
	Certificate CertificateConfig `json:"certificate,omitzero"`
}

type KeyAlgorithm string

const (
	KeyAlgorithmRSA   KeyAlgorithm = "RSA"
	KeyAlgorithmECDSA KeyAlgorithm = "ECDSA"
)

type ECDSACurve string

const (
	ECDSACurveP256 ECDSACurve = "P256"
	ECDSACurveP384 ECDSACurve = "P384"
	ECDSACurveP521 ECDSACurve = "P521"
)

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
