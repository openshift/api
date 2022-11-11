package schemapatch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openshift/api/tools/codegen/pkg/generation"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/schemapatcher"
)

const openshiftFeatureSetEnv = "OPENSHIFT_REQUIRED_FEATURESET"

// executeSchemaPatchForGroupVersionWithBinary runs a schemapatch on the controller-gen binary against the group version
// provided. If any requiredFeatureSets are present it will set the appropriate environment variable to ensure the
// generator only executes the generator on the correct features sets.
func executeSchemaPatchForGroupVersionWithBinary(controllerGen string, group string, version generation.APIVersionContext, versionPaths []string, requiredFeatureSets sets.String) error {
	if requiredFeatureSets.Len() > 0 {
		// The controller generator picks up feature sets from an env var.
		if err := os.Setenv(openshiftFeatureSetEnv, strings.Join(requiredFeatureSets.List(), ",")); err != nil {
			return fmt.Errorf("could not set env %s: %w", openshiftFeatureSetEnv, err)
		}

		defer os.Unsetenv(openshiftFeatureSetEnv)
	}

	args := []string{}

	args = append(args, manifestsArg(version.Path))
	args = append(args, pathsArgs(versionPaths)...)
	args = append(args, outputArg(version.Path))

	cmd := exec.Command(controllerGen, args...)

	// Ensure we get the output from the command.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running controller-gen: %w", err)
	}

	return nil
}

// manifestArg generates the schemapatch:manifests arg for the controller-gen binary.
func manifestsArg(versionPath string) string {
	return fmt.Sprintf("schemapatch:manifests=\"%s\"", versionPath)
}

// pathsArg generates the paths arg for the controller-gen binary.
func pathsArgs(versionPaths []string) []string {
	paths := []string{}

	for _, path := range versionPaths {
		paths = append(paths, fmt.Sprintf("paths=%s", path))
	}

	return paths
}

// outputArg generates the output:dir arg for the controller-gen binary.
func outputArg(versionPath string) string {
	return fmt.Sprintf("output:dir=\"%s\"", versionPath)
}

// executeSchemaPatchForGroupVersion runs the schemapatch code directly for the given group and version.
func executeSchemaPatchForGroupVersion(group string, version generation.APIVersionContext, requiredFeatureSets sets.String, versionPaths []string) error {
	rt, err := loadGroupRuntime(versionPaths)
	if err != nil {
		return fmt.Errorf("error loading group runtime: %w", err)
	}

	markers.RequiredFeatureSets.Insert(requiredFeatureSets.List()...)
	defer func() {
		markers.RequiredFeatureSets = sets.NewString()
	}()

	gen := schemapatcher.Generator{
		ManifestsPath: version.Path,
	}

	ctx := rt.GenerationContext
	ctx.OutputRule = genall.OutputToDirectory(version.Path)

	if err := gen.RegisterMarkers(ctx.Collector.Registry); err != nil {
		return fmt.Errorf("could not register markers: %w", err)
	}

	if err := gen.Generate(&ctx); err != nil {
		return fmt.Errorf("could not run schemapatch generator: %w", err)
	}

	return nil
}

// loadGroupRuntime builds a genall.Runtime based on the package paths for all version that are passed.
// This allows the runtime to be shared between each version when it's generated.
func loadGroupRuntime(paths []string) (*genall.Runtime, error) {
	generators := &genall.Generators{}
	return generators.ForRoots(paths...)
}
