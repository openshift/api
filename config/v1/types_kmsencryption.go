package v1

// KMSConfig defines the configuration for the KMS instance
// that will be used with KMS encryption
// +union
type KMSConfig struct {
	// type defines the kind of platform for the KMS provider.
	//
	// +unionDiscriminator
	// +required
	Type KMSProviderType `json:"type"`

	// --- TOMBSTONE ---
	// aws was a field that allowed configuring AWS KMS.
	// It was never implemented and has been removed.
	// The field name is reserved to prevent reuse.
	//
	// +optional
	// AWS *AWSKMSConfig `json:"aws,omitempty"`
}

// --- TOMBSTONE ---
// AWSKMSConfig was a type for AWS KMS configuration that was never implemented.
// The type name is reserved to prevent reuse.
//
// type AWSKMSConfig struct {
// 	KeyARN string `json:"keyARN"`
// 	Region string `json:"region"`
// }

// KMSProviderType is a specific supported KMS provider
type KMSProviderType string

// const (
	// --- TOMBSTONE ---
	// AWSKMSProvider was a constant for AWS KMS support that was never implemented.
	// The constant name is reserved to prevent reuse.
//	AWSKMSProvider KMSProviderType = "AWS"
// )
