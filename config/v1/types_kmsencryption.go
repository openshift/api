package v1

// KMSConfig defines the configuration for the KMS instance
// that will be used with KMSEncryptionProvider encryption
type KMSConfig struct {
	// managementModel defines how KMS plugins are managed.
	// Valid values are "External".
	// When set to External, encryption keys are managed by a user-deployed
	// KMS plugin that communicates via unix domain socket using KMS V2 API.
	//
	// +kubebuilder:validation:Enum=External
	// +kubebuilder:default=External
	// +optional
	ManagementModel ManagementModel `json:"managementModel,omitempty"`

	// endpoint specifies the unix domain socket endpoint for communicating with the external KMS plugin.
	// The endpoint must follow the format "unix:///path/to/socket" for filesystem-based sockets
	// or "unix:///@abstractname" for abstract sockets.
	//
	// +kubebuilder:validation:MaxLength=120
	// +kubebuilder:validation:MinLength=9
	// +kubebuilder:validation:XValidation:rule="self.matches('^unix:///(@[^/ ]+|[^@ ][^ ]*)$')",message="endpoint must follow the format 'unix:///path' or 'unix:///@abstractname'"
	// +required
	Endpoint string `json:"endpoint,omitempty"`
}

// ManagementModel describes how the KMS plugin is managed.
// Valid values are "External".
type ManagementModel string

const (
	// External represents a KMS plugin that is managed externally and accessed via unix domain socket
	External ManagementModel = "External"
)
