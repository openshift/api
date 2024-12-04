package config

// LintersConfig contains configuration for individual linters.
type LintersConfig struct {
	// jsonTags contains configuration for the jsontags linter.
	JSONTags JSONTagsConfig `json:"jsonTags"`

	// optionalOrRequired contains configuration for the optionalorrequired linter.
	OptionalOrRequired OptionalOrRequiredConfig `json:"optionalOrRequired"`
}

// JSONTagsConfig contains configuration for the jsontags linter.
type JSONTagsConfig struct {
	// jsonTagRegex is the regular expression used to validate that json tags are in a particular format.
	// By default, the regex used is "^[a-z][a-z0-9]*(?:[A-Z][a-z0-9]*)*$" and is used to check for
	// camel case like string.
	JSONTagRegex string `json:"jsonTagRegex"`
}

// OptionalOrRequiredConfig contains configuration for the optionalorrequired linter.
type OptionalOrRequiredConfig struct {
	// preferredOptionalMarker is the preferred marker to use for optional fields.
	// If this field is not set, the default value is "optional".
	// Valid values are "optional" and "kubebuilder:validation:Optional".
	PreferredOptionalMarker string `json:"preferredOptionalMarker"`

	// preferredRequiredMarker is the preferred marker to use for required fields.
	// If this field is not set, the default value is "required".
	// Valid values are "required" and "kubebuilder:validation:Required".
	PreferredRequiredMarker string `json:"preferredRequiredMarker"`
}
