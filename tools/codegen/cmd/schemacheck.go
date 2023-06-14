package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/schemacheck"
	"github.com/spf13/cobra"
)

// schemacheckCmd represents the schemacheck command
var schemacheckCmd = &cobra.Command{
	Use:   "schemacheck",
	Short: "schemacheck validates CRD API schemas based on the best practices",
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newSchemaCheckGenerator()

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(schemacheckCmd)
}

// newSchemaCheckGenerator builds a new schemacheck generator.
func newSchemaCheckGenerator() generation.Generator {
	return schemacheck.NewGenerator(schemacheck.Options{})
}
