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

	// defaultInputPaths contains the default list of input paths for the OpenAPI generator.
	defaultInputPaths = []string{
		"k8s.io/apimachinery/pkg/apis/meta/v1",
		"k8s.io/apimachinery/pkg/runtime",
		"k8s.io/apimachinery/pkg/util/intstr",
		"k8s.io/apimachinery/pkg/api/resource",
		"k8s.io/apimachinery/pkg/version",
		"k8s.io/api/core/v1",
		"k8s.io/api/rbac/v1",
		"k8s.io/api/authorization/v1",
	}
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
	action := "Generating"
	if g.verify {
		action = "Verifying"
	}

	klog.V(1).Infof("%s openapi definitions", action)

	inputPaths := getInputPaths(groupCtxs)

	// If there is no header file, create an empty file and pass that through.
	headerFilePath := g.headerFilePath
	if headerFilePath == "" {
		tmpFile, err := os.CreateTemp("", "openapi-header-*.txt")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %w", err)
		}
		tmpFile.Close()

		defer os.Remove(tmpFile.Name())

		headerFilePath = tmpFile.Name()
	}

	if err := generateOpenAPIDefinitions(inputPaths, g.outputPackagePath, g.outputBaseFileName, headerFilePath, g.verify); err != nil {
		return fmt.Errorf("could not generate openapi definitions: %w", err)
	}

	return nil
}

// getInputPaths collates the input paths from all of the API groups and versions
// within the given group contexts.
// It also includes a standard list of additional packages.
func getInputPaths(groupCtxs []generation.APIGroupContext) []string {
	inputPaths := append([]string{}, defaultInputPaths...)

	for _, groupCtx := range groupCtxs {
		if groupCtx.Config != nil && groupCtx.Config.OpenAPI != nil && groupCtx.Config.OpenAPI.Disabled {
			klog.V(2).Info("Excluding API group %q from openapi generation", groupCtx.Name)
			continue
		}

		for _, version := range groupCtx.Versions {
			inputPaths = append(inputPaths, version.PackagePath)
		}
	}

	return inputPaths
}
