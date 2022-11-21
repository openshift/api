package schemapatch

import (
	"fmt"
	"strings"

	"github.com/openshift/api/tools/codegen/pkg/generation"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

// Options contains the configuration required for the schemapatch generator.
type Options struct {
	// ControllerGen is the path to a controller-gen binary to use for the generation.
	// When omitted, we will use the generator directly from the code.
	ControllerGen string

	// Disabled indicates whether the schemapatch generator is disabled or not.
	// This default to false as the schemapatch generator is enabled by default.
	Disabled bool

	// RequiredFeatureSets is used to filter the feature set manifests that
	// should be generated.
	// When omitted, any manifest with a feature set annotation will be ignored.
	RequiredFeatureSets []sets.String
}

// generator implements the generation.Generator interface.
// It is designed to generate schemapatch updates for a particular API group.
type generator struct {
	controllerGen       string
	disabled            bool
	requiredFeatureSets []sets.String
}

// NewGenerator builds a new schemapatch generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		controllerGen:       opts.ControllerGen,
		disabled:            opts.Disabled,
		requiredFeatureSets: opts.RequiredFeatureSets,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the schemapatch configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.SchemaPatch == nil {
		return g
	}

	featureSets := []sets.String{}
	for _, featureSet := range config.SchemaPatch.RequiredFeatureSets {
		featureSets = append(featureSets, sets.NewString(strings.Split(featureSet, ",")...))
	}

	return NewGenerator(Options{
		ControllerGen:       g.controllerGen,
		Disabled:            config.SchemaPatch.Disabled,
		RequiredFeatureSets: featureSets,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "schemapatch"
}

// GenGroup runs the schemapatch generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	if g.disabled {
		klog.V(3).Infof("Skipping API schema generation for %s", groupCtx.Name)
		return nil
	}

	versionPaths := allVersionPaths(groupCtx.Versions)

	for _, version := range groupCtx.Versions {
		versionRequired, err := shouldProcessGroupVersion(version, g.requiredFeatureSets)
		if err != nil {
			return fmt.Errorf("could not determine if version %s is required: %w", version.Name, err)
		}

		if !versionRequired {
			continue
		}

		if err := g.genGroupVersion(groupCtx.Name, version, versionPaths); err != nil {
			return fmt.Errorf("could not run schemapatch generator for group/version %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil
}

// genGroupVersion runs the schemapatch generator against a particular version of the API group.
func (g *generator) genGroupVersion(group string, version generation.APIVersionContext, versionPaths []string) error {
	if len(g.requiredFeatureSets) == 0 {
		klog.V(2).Infof("Generating API schema for %s/%s", group, version.Name)
		if err := g.executeSchemaPatch(group, version, versionPaths, sets.NewString()); err != nil {
			return fmt.Errorf("could not generate schema patch for %s/%s: %w", group, version.Name, err)
		}
	} else {
		for _, requiredFeatureSet := range g.requiredFeatureSets {
			klog.V(2).Infof("Generating API schema for %s/%s with FeatureSets %v", group, version.Name, requiredFeatureSet.List())
			if err := g.executeSchemaPatch(group, version, versionPaths, requiredFeatureSet); err != nil {
				return fmt.Errorf("could not generate schema patch for %s/%s with feature set %v: %w", group, version.Name, requiredFeatureSet, err)
			}
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

func (g *generator) executeSchemaPatch(group string, version generation.APIVersionContext, versionPaths []string, requiredFeatureSet sets.String) error {
	if g.controllerGen != "" {
		if err := executeSchemaPatchForGroupVersionWithBinary(g.controllerGen, group, version, versionPaths, requiredFeatureSet); err != nil {
			return fmt.Errorf("error executing controller-gen binary: %w", err)
		}
	} else {
		if err := executeSchemaPatchForGroupVersion(group, version, requiredFeatureSet, versionPaths); err != nil {
			return fmt.Errorf("error executing schemapatch: %w", err)
		}
	}

	return nil
}

// allVersionPaths creates a list of all version paths for the group.
func allVersionPaths(versions []generation.APIVersionContext) []string {
	out := []string{}

	for _, version := range versions {
		out = append(out, version.Path)
	}

	return out
}
