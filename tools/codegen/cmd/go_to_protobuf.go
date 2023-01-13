package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/protobuf"
	"github.com/spf13/cobra"
)

// goToProtobufCmd represents the go-to-protobuf command
var goToProtobufCmd = &cobra.Command{
	Use:   "go-to-protobuf",
	Short: "go-to-protobuf generates protofbuf definitions for API types",
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newGoToProtobufGenerator()

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(goToProtobufCmd)

}

// newGoToProtobufGenerator builds a new go-to-protobuf generator.
func newGoToProtobufGenerator() generation.Generator {
	return protobuf.NewGenerator(protobuf.Options{
		Verify:             verify,
	})
}
