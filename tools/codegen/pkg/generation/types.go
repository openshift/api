package generation

// Config represents the configuration of a API group version
// and the configuration for each generator within it.
type Config struct {
	// SchemaPatch represents the configuration for the schemapatch generator.
	// When omitted, the default configuration will be used.
	// When provided, any equivalent flag provided values are ignored.
	SchemaPatch *SchemaPatchConfig `json:"schemapatch,omitempty"`
}

// SchemaPatchConfig is the configuration for the schemapatch generator.
type SchemaPatchConfig struct {
	// Disabled determines whether the schemapatch generator should be run or not.
	// This generator is enabled by default so this field defaults to false.
	Disabled bool `json:"disabled,omitempty"`

	// RequiredFeatureSets is a list of feature sets combinations that should be
	// generated for this API group.
	// Each entry in this list is a comma separated list of feature set names
	// which will be matched with the `release.openshift.io/feature-set` annotation
	// on the CRD definition.
	// When omitted, any manifest with a feature set annotation will be ignored.
	// Example entries are `""` (empty string), `"TechPreviewNoUpgrade"` or `"TechPreviewNoUpgrade,CustomNoUpgrade"`.
	RequiredFeatureSets []string `json:"requiredFeatureSets,omitempty"`
}
