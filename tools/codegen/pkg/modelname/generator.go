package modelname

import (
	"fmt"
	"os"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
)

const (
	// DefaultOutputFileName is the default output file name for the generated model name functions.
	DefaultOutputFileName = "zz_generated.model_name.go"
)

// Options contains the configuration required for the model name generator.
type Options struct {
	// Disabled indicates whether the model name generator is enabled or not.
	// This defaults to false as the model name generator is enabled by default.
	Disabled bool

	// HeaderFilePath is the path to the file containing the boilerplate header text.
	// When omitted, no header is added to the generated files.
	HeaderFilePath string

	// OutputFileName is the file name to use for writing the generated model
	// names to. This file will be created for each group version.
	// When omitted, this will default to `zz_generated.model_name.go`.
	OutputFileName string

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool

	// GlobalParser is the parser for the global package.
	// This loads all packages found in the base directory.
	GlobalParser *parser.Parser

	// Universe is the universe for the global package.
	Universe types.Universe
}

// generator implements the generation.Generator interface.
// It is designed to generate model name functions for a particular API group.
type generator struct {
	disabled       bool
	headerFilePath string
	outputFileName string
	verify         bool
	globalParser   *parser.Parser
	universe       types.Universe
}

// NewGenerator builds a new model name generator.
func NewGenerator(opts Options) generation.Generator {
	outputFileName := DefaultOutputFileName
	if opts.OutputFileName != "" {
		outputFileName = opts.OutputFileName
	}

	return &generator{
		disabled:       opts.Disabled,
		headerFilePath: opts.HeaderFilePath,
		outputFileName: outputFileName,
		verify:         opts.Verify,
		globalParser:   opts.GlobalParser,
		universe:       opts.Universe,
	}
}

// ApplyConfig returns a new generator based on the configuration passed.
// If the model name configuration is empty, the existing generator is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.ModelName == nil {
		return g
	}

	outputFileName := DefaultOutputFileName
	if config.ModelName.OutputFileName != "" {
		outputFileName = config.ModelName.OutputFileName
	}

	return NewGenerator(Options{
		Disabled:       config.ModelName.Disabled,
		HeaderFilePath: g.headerFilePath,
		OutputFileName: outputFileName,
		Verify:         g.verify,
		GlobalParser:   g.globalParser,
		Universe:       g.universe,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "modelname"
}

// GenGroup runs the model name generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) ([]generation.Result, error) {
	if g.disabled {
		klog.V(2).Infof("Skipping model name generation for %s", groupCtx.Name)
		return nil, nil
	}

	// If there is no header file, create an empty file and pass that through.
	headerFilePath := g.headerFilePath
	if headerFilePath == "" {
		tmpFile, err := os.CreateTemp("", "modelname-header-*.txt")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary file: %w", err)
		}
		tmpFile.Close()

		defer os.Remove(tmpFile.Name())

		headerFilePath = tmpFile.Name()
	}

	for _, version := range groupCtx.Versions {
		action := "Generating"
		if g.verify {
			action = "Verifying"
		}

		klog.V(1).Infof("%s model name functions for %s/%s", action, groupCtx.Name, version.Name)

		if err := generateModelNames(g.globalParser, g.universe, version.PackagePath, g.outputFileName, headerFilePath, g.verify); err != nil {
			return nil, fmt.Errorf("could not generate model name functions for %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil, nil
}
