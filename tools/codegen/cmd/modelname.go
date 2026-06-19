package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/modelname"
	"github.com/spf13/cobra"
)

var (
	modelNameHeaderFilePath string
	modelNameOutputFileName string
)

// modelNameCmd represents the modelname command
var modelNameCmd = &cobra.Command{
	Use:   "modelname",
	Short: "modelname generates OpenAPI model name accessor functions from API definitions",
	Long: `modelname generates OpenAPI model name accessor functions from API definitions.

The generator creates an OpenAPIModelName() method on each type in the API which
returns the fully qualified OpenAPI model name for the type. The model name is
determined by the +k8s:openapi-model-package tag in the package doc.go file.

These model names are used in OpenAPI specs as schema references instead of
Go type names, allowing for explicit control over the OpenAPI schema naming.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newModelNameGenerator(genCtx)

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(modelNameCmd)

	rootCmd.PersistentFlags().StringVar(&modelNameHeaderFilePath, "modelname:header-file-path", "", "Path to file containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year. When omitted, no header is added to the generated files.")
	rootCmd.PersistentFlags().StringVar(&modelNameOutputFileName, "modelname:output-file-name", modelname.DefaultOutputFileName, "Defines the file name to use for the model name generated functions for each group version.")
}

// newModelNameGenerator builds a new model name generator.
func newModelNameGenerator(genCtx generation.Context) generation.Generator {
	return modelname.NewGenerator(modelname.Options{
		HeaderFilePath: modelNameHeaderFilePath,
		OutputFileName: modelNameOutputFileName,
		Verify:         verify,
		GlobalParser:   genCtx.GlobalParser,
		Universe:       genCtx.Universe,
	})
}
