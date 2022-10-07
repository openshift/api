package schemapatch

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-tools/pkg/genall"
)

// Options contains the configuration required for the schemapatch generator.
type Options struct {
	// ControllerGen is the path to a controller-gen binary to use for the generation.
	// When omitted, we will use the generator directly from the code.
	ControllerGen string

	// RequiredFeatureSets is used to filter the feature set manifests that
	// should be generated.
	// When omitted, any manifest with a feature set annotation will be ignored.
	RequiredFeatureSets []string
}

// generator implements the generation.Generator interface.
// It is designed to generate schemapatch updates for a particular API group.
type generator struct {
	controllerGen       string
	requiredFeatureSets sets.String
}

// NewGenerator builds a new schemapatch generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		controllerGen:       opts.ControllerGen,
		requiredFeatureSets: sets.NewString(opts.RequiredFeatureSets...),
	}
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "schemapatch"
}

// GenGroup runs the schemapatch generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	var groupRuntime *genall.Runtime
	versionPaths := allVersionPaths(groupCtx.Versions)

	for _, version := range groupCtx.Versions {
		versionRequired, err := shouldProcessGroupVersion(version, g.requiredFeatureSets)
		if err != nil {
			return fmt.Errorf("could not determine if version %s is required: %w", version.Name, err)
		}

		if !versionRequired {
			continue
		}

		// Lazy load the roots for the group once we know at least one
		// version needs to be generated.
		// Only load the roots if we are using internal generation.
		if groupRuntime == nil && g.controllerGen == "" {
			rt, err := loadGroupRuntime(versionPaths)
			if err != nil {
				return fmt.Errorf("error loading group runtime: %w", err)
			}

			groupRuntime = rt
		}

		if err := g.genGroupVersion(groupCtx.Name, version, groupRuntime, versionPaths); err != nil {
			return fmt.Errorf("could not run schemapatch generator for group/version %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil
}

// genGroupVersion runs the schemapatch generator against a particular version of the API group.
func (g *generator) genGroupVersion(group string, version generation.APIVersionContext, rt *genall.Runtime, versionPaths []string) error {
	if g.controllerGen != "" {
		if err := executeSchemaPatchForGroupVersionWithBinary(g.controllerGen, group, version, versionPaths, g.requiredFeatureSets); err != nil {
			return fmt.Errorf("error executing controller-gen binary: %w", err)
		}
	} else {
		if err := executeSchemaPatchForGroupVersion(rt, group, version, g.requiredFeatureSets); err != nil {
			return fmt.Errorf("error executing schemapatch: %w", err)
		}
	}

	if err := executeYAMLPatchForGroupVersion(version, g.requiredFeatureSets); err != nil {
		return fmt.Errorf("error executing yaml patches: %w", err)
	}

	if err := formatManifestsForGroupVersion(version, g.requiredFeatureSets); err != nil {
		return fmt.Errorf("error formatting manifests: %w", err)
	}

	return nil
}

// loadGroupRuntime builds a genall.Runtime based on the package paths for all version that are passed.
// This allows the runtime to be shared between each version when it's generated.
func loadGroupRuntime(paths []string) (*genall.Runtime, error) {
	generators := &genall.Generators{}
	return generators.ForRoots(paths...)
}

// allVersionPaths creates a list of all version paths for the group.
func allVersionPaths(versions []generation.APIVersionContext) []string {
	out := []string{}

	for _, version := range versions {
		out = append(out, version.Path)
	}

	return out
}
