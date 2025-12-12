package v1

// KMSConfig defines the configuration for the KMS instance
// that will be used with KMSEncryptionProvider encryption
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Manual' ? (has(self.manual) && has(self.manual.name) && self.manual.name != '') : !has(self.manual)",message="manual config with non-empty name is required when kms provider type is Manual, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'AWS' ?  has(self.aws) : !has(self.aws)",message="aws config is required when kms provider type is AWS, and forbidden otherwise"
// +union
type KMSConfig struct {
	// type defines the kind of platform for the KMS provider.
	// Available provider types are AWS, Manual.
	//
	// +unionDiscriminator
	// +required
	Type KMSProviderType `json:"type"`

	// manual defines the configuration for manually managed KMS plugins.
	// The KMS plugin must be deployed as a static pod by the cluster admin.
	//
	// +unionMember
	// +optional
	Manual *ManualKMSConfig `json:"manual,omitempty"`

	// aws defines the key config for using an AWS KMS instance
	// for the encryption. The AWS KMS instance is managed
	// by the user outside the purview of the control plane.
	// Deprecated: There is no logic listening to this resource type, we plan to remove it in next release.
	//
	// +unionMember
	// +optional
	AWS *AWSKMSConfig `json:"aws,omitempty"`
}

// ManualKMSConfig defines the configuration for manually managed KMS plugins
type ManualKMSConfig struct {
	// name specifies the KMS plugin name.
	// This name is templated into the UNIX domain socket path: unix:///var/run/kmsplugin/<name>.sock
	// and is between 1 and 80 characters in length.
	// The KMS plugin must listen at this socket path.
	// The name must be a safe socket filename and must not contain '/' or '..'.
	//
	// +kubebuilder:validation:MaxLength=80
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="!self.contains('/') && !self.contains('..')",message="name must be a safe socket filename (must not contain '/' or '..')"
	// +optional
	Name string `json:"name,omitempty"`
}

// AWSKMSConfig defines the KMS config specific to AWS KMS provider
// Deprecated: There is no logic listening to this resource type, we plan to remove it in next release.
type AWSKMSConfig struct {
	// keyARN specifies the Amazon Resource Name (ARN) of the AWS KMS key used for encryption.
	// The value must adhere to the format `arn:aws:kms:<region>:<account_id>:key/<key_id>`, where:
	// - `<region>` is the AWS region consisting of lowercase letters and hyphens followed by a number.
	// - `<account_id>` is a 12-digit numeric identifier for the AWS account.
	// - `<key_id>` is a unique identifier for the KMS key, consisting of lowercase hexadecimal characters and hyphens.
	//
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.matches('^arn:aws:kms:[a-z0-9-]+:[0-9]{12}:key/[a-f0-9-]+$')",message="keyARN must follow the format `arn:aws:kms:<region>:<account_id>:key/<key_id>`. The account ID must be a 12 digit number and the region and key ID should consist only of lowercase hexadecimal characters and hyphens (-)."
	// +required
	KeyARN string `json:"keyARN"`
	// region specifies the AWS region where the KMS instance exists, and follows the format
	// `<region-prefix>-<region-name>-<number>`, e.g.: `us-east-1`.
	// Only lowercase letters and hyphens followed by numbers are allowed.
	//
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]+(-[a-z0-9]+)*$')",message="region must be a valid AWS region, consisting of lowercase characters, digits and hyphens (-) only."
	// +required
	Region string `json:"region"`
}

// KMSProviderType is a specific supported KMS provider
// +kubebuilder:validation:Enum=AWS;Manual
type KMSProviderType string

const (
	// AWSKMSProvider represents a supported KMS provider for use with AWS KMS
	// Deprecated: There is no logic listening to this resource type, we plan to remove it in next release.
	AWSKMSProvider KMSProviderType = "AWS"

	// ManualKMSProvider represents a supported KMS provider is managed by user manually not by OpenShift.
	// KMS plugin is supposed to be run as static pods on each control plane
	ManualKMSProvider KMSProviderType = "Manual"
)
