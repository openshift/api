package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/swaggerdocs"
	"github.com/spf13/cobra"
)

// swaggerDocsCmd represents the swaggerdocs command
var swaggerDocsCmd = &cobra.Command{
	Use:   "swaggerdocs",
	Short: "swaggerdocs generates swagger documentation from API definitions",
	Long: `swaggerdocs generates swagger documentation from API definitions.
	
	The generator creates a SwaggerDoc method on each type in the API which
	returns a map of fields to their documentation. The documentation is sourced
	from the godoc on each field.
	
	A warning will be produced whenever a field is missing documentation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newSwaggerDocsGenerator()

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(swaggerDocsCmd)
}

// newSwaggerDocsGenerator builds a new swaggerdocs generator.
func newSwaggerDocsGenerator() generation.Generator {
	return swaggerdocs.NewGenerator(swaggerdocs.Options{})
}
