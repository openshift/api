package protobuf

import (
	"fmt"
	"os"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/klog/v2"
)

const (
	// DefaultOutputFileBaseName is the default output file base name for the generated protobuf functions.
	DefaultOutputFileBaseName = "generated.pb"
)

// Options contains the configuration required for the protobuf generator.
type Options struct {
	// Disabled indicates whether the deepcopy generator is enabled or not.
	// This defaults to true as the protobuf generator is disabled by default.
	Disabled *bool

	// HeaderFilePath is the path to the file containing the boilerplate header text.
	// When omitted, no header is added to the generated files.
	HeaderFilePath string

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate deepcopy function for a particular API group.
type generator struct {
	disabled       bool
	headerFilePath string
	verify         bool
}

// NewGenerator builds a new protobuf generator.
func NewGenerator(opts Options) generation.Generator {
	disabled := true
	if opts.Disabled != nil {
		disabled = *opts.Disabled
	}

	return &generator{
		disabled:       disabled,
		headerFilePath: opts.HeaderFilePath,
		verify:         opts.Verify,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the deepcopy configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.Protobuf == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled:       config.Protobuf.Disabled,
		HeaderFilePath: config.Protobuf.HeaderFilePath,
		Verify:         g.verify,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "go-to-protobuf"
}

// GenGroup runs the go-to-protobuf generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) ([]generation.Result, error) {
	if g.disabled {
		return nil, nil
	}

	// If there is no header file, create an empty file and pass that through.
	headerFilePath := g.headerFilePath
	if headerFilePath == "" {
		tmpFile, err := os.CreateTemp("", "protobuf-header-*.txt")
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

		klog.V(1).Infof("%s protobuf functions for for %s/%s", action, groupCtx.Name, version.Name)

		if err := generateProtobufFunctions(version.Path, version.PackagePath, headerFilePath, g.verify); err != nil {
			return nil, fmt.Errorf("could not generate protobuf functions for %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil, nil
}
