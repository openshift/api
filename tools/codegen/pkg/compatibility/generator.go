package compatibility

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/klog/v2"
)

// Options contains the configuration required for the compatibility generator.
type Options struct {
	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate compatibility level comments for a particular API group.
type generator struct {
	verify bool
}

// NewGenerator builds a new compatibility generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		verify: opts.Verify,
	}
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "compatibility"
}

// GenGroup runs the compatibility generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	for _, version := range groupCtx.Versions {
		action := "Generating"
		if g.verify {
			action = "Verifying"
		}

		klog.V(2).Infof("%s compatibility level comments for %s/%s", action, groupCtx.Name, version.Name)

		if err := insertCompatibilityLevelComments(version.Path, g.verify); err != nil {
			return fmt.Errorf("could not insert compatibility level comments for %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil
}
