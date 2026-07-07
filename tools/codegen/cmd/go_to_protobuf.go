package main

import (
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/protobuf"
	"github.com/spf13/cobra"
)

var (
	goToProtobufHeaderFilePath string
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

	rootCmd.PersistentFlags().StringVar(&goToProtobufHeaderFilePath, "go-to-protobuf:header-file-path", "", "Path to file containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year. When omitted, no header is added to the generated files.")
}

// newGoToProtobufGenerator builds a new go-to-protobuf generator.
func newGoToProtobufGenerator() generation.Generator {
	return protobuf.NewGenerator(protobuf.Options{
		HeaderFilePath: goToProtobufHeaderFilePath,
		Verify:         verify,
	})
}
