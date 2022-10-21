package main

import (
	"flag"
	"fmt"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/spf13/cobra"

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
)

var (
	apiGroupVersions []string
	baseDir          string
)

// rootCmd represents the base command when called without any subcommands.
// This will run all generators in the preferred order for OpenShift APIs.
var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Codegen runs code generators for the OpenShift API definitions",
	RunE: func(cmd *cobra.Command, args []string) error {
		genCtx, err := generation.NewContext(generation.Options{
			BaseDir:          baseDir,
			APIGroupVersions: apiGroupVersions,
		})
		if err != nil {
			return fmt.Errorf("could not build generation context: %w", err)
		}

		return executeGenerators(genCtx, allGenerators()...)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		klog.Fatalf("Error running codegen: %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&apiGroupVersions, "api-group-versions", []string{}, "A list of API group versions in the form <group>/<version>. The group should be fully qualified, e.g. machine.openshift.io/v1. The generator will generate against all group versions found within the base directory when no specific group versions are provided.")
	rootCmd.PersistentFlags().StringVar(&baseDir, "base-dir", ".", "Base directory to search for API group versions")

	klog.InitFlags(flag.CommandLine)
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

// executeGenerators runs each generator for each group in the generation context.
// If an error occurs for a generator within a group, the rest of the generators are ignored for that group.
// Subsequent groups will continue to generate.
func executeGenerators(genCtx generation.Context, generators ...generation.Generator) error {
	errs := []error{}

	for _, group := range genCtx.APIGroups {
		klog.Infof("Running generators for %s", group.Name)

		for _, gen := range generators {
			if err := gen.GenGroup(group); err != nil {
				errs = append(errs, fmt.Errorf("error running generator %s on group %s: %w", gen.Name(), group.Name, err))

				// Don't run any later generators for this group.
				break
			}
		}
	}

	if len(errs) > 0 {
		return kerrors.NewAggregate(errs)
	}

	return nil
}

// allGenerators returns an ordered list of generators to run when
// the root command is executed.
func allGenerators() []generation.Generator {
	return []generation.Generator{
		newSchemaPatchGenerator(),
	}
}
