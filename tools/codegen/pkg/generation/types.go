package generation

// Config represents the configuration of a API group version
// and the configuration for each generator within it.
type Config struct {
	// Compatibility represents the configuration of the compatiblity generator.
	// When omitted, the default configuration will be used.
	Compatibility *CompatibilityConfig `json:"compatibility,omitempty"`

	// SchemaPatch represents the configuration for the schemapatch generator.
	// When omitted, the default configuration will be used.
	// When provided, any equivalent flag provided values are ignored.
	SchemaPatch *SchemaPatchConfig `json:"schemapatch,omitempty"`

	// SwaggerDocs represents the configuration for the swaggerdocs generator.
	// When omitted, the default configuration will be used.
	// When provided, any equivalent flag provided values are ignored.
	SwaggerDocs *SwaggerDocsConfig `json:"swaggerdocs,omitempty"`
}

// CompatibilityConfig is the configuration for the compatibility generator.
type CompatibilityConfig struct {
	// Disabled determines whether the compatibility generator should be run or not.
	// This generator is enabled by default so this field defaults to false.
	Disabled bool `json:"disabled,omitempty"`
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

// SwaggerDocsConfig is the configuration for the swaggerdocs generator.
type SwaggerDocsConfig struct {
	// Disabled determines whether the swaggerdocs generator should be run or not.
	// This generator is enabled by default so this field defaults to false.
	Disabled bool `json:"disabled,omitempty"`

	// OutputFileName is the file name to use for writing the generated swagger
	// docs to. This file will be created for each group version.
	// Whem omitted, this will default to `zz_generated.swagger_doc_generated.go`.
	OutputFileName string `json:"outputFileName,omitempty"`

	// EnforceComments determines whether the generator should enforce that all
	// fields have a comment or not.
	// When set to true, verification will fail if any field in the API group is
	// missing a comment.
	EnforceComments bool `json:"enforceComments,omitempty"`
}
