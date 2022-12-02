package openapi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/klog/v2"
)

const (
	// DefaultOutputFileBaseName is the default output file base name for the generated openapi functions.
	DefaultOutputFileBaseName = "zz_generated.openapi"
)

var (
	// DefaultOutputPackagePath is the default output package path for the generated openapi functions.
	DefaultOutputPackagePath = filepath.Join("openapi", "generated_openapi")
)

// Options contains the configuration required for the compatibility generator.
type Options struct {
	// HeaderFilePath is the path to the file containing the boilerplate header text.
	// When omitted, no header is added to the generated files.
	HeaderFilePath string

	// OutputFileBaseName is the base name of the output file.
	// When omitted, DefaultOutputFileBaseName is used.
	// The current value of DefaultOutputFileBaseName is "zz_generated.openapi".
	OutputFileBaseName string

	// OutputPackagePath is the package path where the generated golang files will be written.
	OutputPackagePath string

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate openapi function for a particular API group.
type generator struct {
	headerFilePath     string
	outputBaseFileName string
	outputPackagePath  string
	verify             bool
}

// NewGenerator builds a new openapi generator.
func NewGenerator(opts Options) generation.MultiGroupGenerator {
	outputFileBaseName := DefaultOutputFileBaseName
	if opts.OutputFileBaseName != "" {
		outputFileBaseName = opts.OutputFileBaseName
	}

	outputPackagePath := DefaultOutputPackagePath
	if opts.OutputPackagePath != "" {
		outputPackagePath = opts.OutputPackagePath
	}

	return &generator{
		headerFilePath:     opts.HeaderFilePath,
		outputBaseFileName: outputFileBaseName,
		outputPackagePath:  outputPackagePath,
		verify:             opts.Verify,
	}
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "openapi"
}

// GenGroup runs the openapi generator against the given group context.
func (g *generator) GenGroups(groupCtxs []generation.APIGroupContext) error {
	return nil
}
