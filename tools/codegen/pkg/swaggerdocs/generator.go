package swaggerdocs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

// Options contains the configuration required for the swaggerdocs generator.
type Options struct {
	// OutputFileName is the file name to use for writing the generated swagger
	// docs to. This file will be created for each group version.
	OutputFileName string

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate swaggerdocs documentation for a particular API group.
type generator struct {
	outputFileName string
	verify         bool
}

// NewGenerator builds a new schemapatch generator.
func NewGenerator(opts Options) generation.Generator {
	return &generator{
		outputFileName: opts.OutputFileName,
		verify:         opts.Verify,
	}
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "swaggerdocs"
}

// GenGroup runs the schemapatch generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) error {
	for _, version := range groupCtx.Versions {
		if err := g.generateGroupVersion(groupCtx.Name, version); err != nil {
			return fmt.Errorf("error generating swagger docs for %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil
}

// generateGroupVersion generates swagger docs for the group version.
func (g *generator) generateGroupVersion(groupName string, version generation.APIVersionContext) error {
	outFilePath := filepath.Join(version.Path, g.outputFileName)

	versionGlob := filepath.Join(version.Path, typesGlob)
	files, err := filepath.Glob(versionGlob)
	if err != nil {
		return fmt.Errorf("could not read types*.go files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no match for types*.go glob in path %s", version.Path)
	}

	docsForTypes := []kruntime.KubeTypes{}
	for _, file := range files {
		docsForTypes = append(docsForTypes, kruntime.ParseDocumentationFrom(file)...)
	}

	if g.verify {
		klog.V(2).Infof("Verifiying swagger docs for %s/%s", groupName, version.Name)

		return verifySwaggerDocs(version.Name, outFilePath, docsForTypes)
	}

	klog.V(2).Infof("Generating swagger docs for %s/%s", groupName, version.Name)

	generatedDocs, err := generateSwaggerDocs(version.Name, docsForTypes)
	if err != nil {
		return fmt.Errorf("error generating swagger docs: %w", err)
	}

	if err := ioutil.WriteFile(outFilePath, generatedDocs, 0644); err != nil {
		return fmt.Errorf("error writing swagger docs output: %w", err)
	}

	return nil
}
