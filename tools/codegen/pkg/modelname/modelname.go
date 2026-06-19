package modelname

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"github.com/openshift/api/tools/codegen/pkg/utils"
	"k8s.io/gengo/v2"
	gengenerator "k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/parser"
	"k8s.io/gengo/v2/types"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/cmd/openapi-gen/args"
	"k8s.io/kube-openapi/pkg/generators"
)

// generateModelNames generates the model name functions for the given API package path.
func generateModelNames(globalParser *parser.Parser, universe types.Universe, packagePath, outputFileName, headerFilePath string, verify bool) error {
	// Get the package information from the universe
	pkg := universe[packagePath]
	if pkg == nil {
		return fmt.Errorf("package %s not found in universe", packagePath)
	}

	outputFile := filepath.Join(pkg.Dir, outputFileName)
	originalDir := pkg.Dir

	// In verify mode, generate to a temporary directory to avoid modifying the working tree
	var tmpDir string
	if verify {
		var err error
		tmpDir, err = os.MkdirTemp("", "codegen-modelname-verify-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		// Modify the package directory to point to temp dir.
		// This is safe because we're modifying a local copy of the universe.
		pkg.Dir = tmpDir
	}

	arguments := args.New()
	arguments.OutputModelNameFile = outputFileName
	arguments.GoHeaderFile = headerFilePath
	arguments.OutputDir = pkg.Dir
	arguments.OutputPkg = packagePath
	// OpenAPI args require OutputFile to be set even though we're only generating model names
	arguments.OutputFile = "zz_generated.openapi.go" // Required by validation but not used for model names

	if err := arguments.Validate(); err != nil {
		return err
	}

	klog.V(2).Infof("Generating model names for package %s", packagePath)

	// Load boilerplate before creating the targets function so we can properly return errors
	boilerplate, err := gengo.GoBoilerplate(arguments.GoHeaderFile, gengo.StdBuildTag, gengo.StdGeneratedBy)
	if err != nil {
		return fmt.Errorf("failed loading boilerplate: %w", err)
	}

	myTargets := func(context *gengenerator.Context) []gengenerator.Target {
		return generators.GetModelNameTargets(context, arguments, boilerplate)
	}

	if err := generation.Execute(
		globalParser,
		universe,
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		myTargets,
		[]string{packagePath},
	); err != nil {
		return fmt.Errorf("error executing model name generator: %w", err)
	}

	if verify {
		// Restore the original directory in the universe
		pkg.Dir = originalDir

		verifyFile := filepath.Join(tmpDir, outputFileName)
		generatedData, err := os.ReadFile(verifyFile)
		if err != nil {
			if os.IsNotExist(err) {
				// File wasn't generated - this package doesn't have types that need model names
				return nil
			}
			return fmt.Errorf("failed to read generated file for verification: %w", err)
		}

		currentData, err := os.ReadFile(outputFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("model names for %s do not exist, please run without --verify to generate them", packagePath)
			}
			return fmt.Errorf("failed to read existing file for verification: %w", err)
		}

		if !bytes.Equal(currentData, generatedData) {
			// Get the relative path for better error messages
			wd, _ := os.Getwd()
			relPath, _ := filepath.Rel(wd, outputFile)
			if relPath == "" {
				relPath = outputFile
			}

			diff := utils.Diff(currentData, generatedData, relPath)
			return fmt.Errorf("model names for %s are out of date, please regenerate the model names:\n%s", packagePath, diff)
		}
	}

	return nil
}
