package generation

// Generator is an interface for running a generator against a particular API group.
type Generator interface {
	// Name returns a name identifier for the generator.
	Name() string

	// GenGroup runs the generator against the given APIGroupContext.
	GenGroup(APIGroupContext) error

	// ApplyConfig creates a new generator instance with the given configuration.
	ApplyConfig(*Config) Generator
}
