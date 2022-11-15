package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/deepcopy"
	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/spf13/cobra"
)

// deepcopyCmd represents the deepcopy command
var deepcopyCmd = &cobra.Command{
	Use:   "deepcopy",
	Short: "deepcopy generates deepcopy functions for API types",
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newDeepcopyGenerator()

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(deepcopyCmd)
}

// newDeepcopyhGenerator builds a new deepcopy generator.
func newDeepcopyGenerator() generation.Generator {
	return deepcopy.NewGenerator(deepcopy.Options{
		Verify: verify,
	})
}
