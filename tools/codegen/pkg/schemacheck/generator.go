package schemacheck

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
)

// Options contains the configuration required for the schemacheck generator.
type Options struct {
	// Disabled indicates whether the schemacheck generator is disabled or not.
	// This defaults to false as the schemacheck generator is enabled by default.
	Disabled bool
}

// generator implements the generation.Generator interface.
// It is designed to verify the CRD schema updates for a particular API group.
type generator struct {
	disabled bool
}

// NewGenerator builds a new schemacheck generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		disabled: opts.Disabled,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the schemacheck configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.SchemaCheck == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled: config.SchemaCheck.Disabled,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "schemacheck"
}

// GenGroup runs the schemacheck generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	if g.disabled {
		klog.V(2).Infof("Skipping API schema check for %s", groupCtx.Name)
		return nil
	}

	errs := []error{}

	for _, version := range groupCtx.Versions {
		klog.V(1).Infof("Verifying API schema for for %s/%s", groupCtx.Name, version.Name)

		if err := g.genGroupVersion(groupCtx.Name, version); err != nil {
			errs = append(errs, fmt.Errorf("could not run schemacheck generator for group/version %s/%s: %w", groupCtx.Name, version.Name, err))
		}
	}

	if len(errs) > 0 {
		return kerrors.NewAggregate(errs)
	}

	return nil
}

// genGroupVersion runs the schemacheck generator against a particular version of the API group.
func (g *generator) genGroupVersion(group string, version generation.APIVersionContext) error {
	return nil
}
