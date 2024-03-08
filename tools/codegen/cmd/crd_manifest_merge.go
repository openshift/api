package main

import (
	"fmt"
	"github.com/openshift/api/tools/codegen/pkg/manifestmerge"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/spf13/cobra"
)

// schemapatchCmd represents the schemapatch command
var crdManifestMerge = &cobra.Command{
	Use:   "crd-manifest-merge",
	Short: "crd-manifest-merge takes all CRD manifests with the same name and merges them together",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		gen := newCRDManifestMerger()

		return executeGenerators(genCtx, gen)
	},
}

func init() {
	rootCmd.AddCommand(crdManifestMerge)
}

// newSchemaPatchGenerator builds a new schemapatch generator.
func newCRDManifestMerger() generation.Generator {
	return manifestmerge.NewGenerator(manifestmerge.Options{
		Verify: verify,
	})
}
