package v1

// KMSConfig defines the configuration for the KMS instance
// that will be used with KMSEncryptionProvider encryption
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'AWS' ?  has(self.aws) : !has(self.aws)",message="aws config is required when kms provider type is AWS, and forbidden otherwise"
// +union
type KMSConfig struct {
	// type defines the kind of platform for the KMS provider
	//
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	Type KMSProviderType `json:"type"`

	// aws defines the key config for using an AWS KMS instance
	// for the encryption. The AWS KMS instance is managed
	// by the user outside the purview of the control plane.
	//
	// +unionMember
	// +optional
	AWS *AWSKMSConfig `json:"aws,omitempty"`
}

// AWSKMSConfig defines the KMS config specific to AWS KMS provider
type AWSKMSConfig struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.matches('^arn:aws:kms:[a-z0-9-]+:[0-9]{12}:key/[a-f0-9-]+$')",message="keyARN must start with `arn:aws:kms:` followed by region, account and key id"
	KeyARN string `json:"keyARN"`
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z]{2}-[a-z]+-[0-9]+$')",message="region must be a valid AWS region"
	Region string `json:"region"`
}

// KMSProviderType is a specific supported KMS provider
// +kubebuilder:validation:Enum="";AWS
type KMSProviderType string

const (
	// AWSKMSProvider represents a supported KMS provider for use with AWS KMS
	AWSKMSProvider KMSProviderType = "AWS"
)
