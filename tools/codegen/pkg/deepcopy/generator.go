package deepcopy

import (
	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/klog/v2"
)

// Options contains the configuration required for the compatibility generator.
type Options struct {
	// Disabled indicates whether the deepcopy generator is enabled or not.
	// This default to false as the deepcopy generator is enabled by default.
	Disabled bool

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate deepcopy function for a particular API group.
type generator struct {
	disabled bool
	verify   bool
}

// NewGenerator builds a new deepcopy generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		disabled: opts.Disabled,
		verify:   opts.Verify,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the deepcopy configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.Deepcopy == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled: config.Deepcopy.Disabled,
		Verify:   g.verify,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "deepcopy"
}

// GenGroup runs the deepcopy generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	if g.disabled {
		klog.V(2).Infof("Skipping deepcopy generation for %s", groupCtx.Name)
		return nil
	}

	return nil
}
