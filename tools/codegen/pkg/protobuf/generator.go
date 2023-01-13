package protobuf

import (
	"github.com/openshift/api/tools/codegen/pkg/generation"
)


// Options contains the configuration required for the protobuf generator.
type Options struct {
	// Disabled indicates whether the deepcopy generator is enabled or not.
	// This defaults to true as the protobuf generator is disabled by default.
	Disabled *bool

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate deepcopy function for a particular API group.
type generator struct {
	disabled           bool
	verify             bool
}

// NewGenerator builds a new protobuf generator.
func NewGenerator(opts Options) generation.Generator {
	disabled := true
	if opts.Disabled != nil {
		disabled = *opts.Disabled
	}

	return &generator{
		disabled:           disabled,
		verify:             opts.Verify,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the deepcopy configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.Protobuf == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled:           config.Protobuf.Disabled,
		Verify:             g.verify,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "go-to-protobuf"
}

// GenGroup runs the go-to-protobuf generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	return nil
}
