package openapi

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/api/tools/codegen/pkg/utils"
	"k8s.io/gengo/v2"
	gengogenerator "k8s.io/gengo/v2/generator"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/cmd/openapi-gen/args"
	"k8s.io/kube-openapi/pkg/generators"
)

// generateDeepcopyFunctions generates the OpenAPI functions for the given API package paths.
func generateOpenAPIDefinitions(inputPaths []string, outputPackagePath, outputBaseFileName, headerFilePath string, verify bool) error {
	// This is the expected path to the output file.
	// This is what we will compare against if verify is true.
	outputFile := filepath.Join(outputPackagePath, outputBaseFileName+".go")

	if verify {
		outputPackageBase := filepath.Base(outputPackagePath)

		tmpDir, err := os.MkdirTemp("", "codegen-openapi-verify-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		outputPackagePath = filepath.Join(tmpDir, outputPackageBase)
	}

	args := &args.Args{
		GoHeaderFile: headerFilePath,
		OutputFile:   outputBaseFileName,
		OutputPkg:    outputPackagePath,
		OutputDir:    outputPackagePath,
	}

	// Temporary to prevent me merging this
	if verify {
		return errors.New("verify is not supported")
	}

	klog.V(2).Infof("Generating openapi into %s", outputPackagePath)

	myTargets := func(context *gengogenerator.Context) []gengogenerator.Target {
		return generators.GetTargets(context, args)
	}

	if err := gengo.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		myTargets,
		gengo.StdBuildTag,
		inputPaths,
	); err != nil {
		return fmt.Errorf("error executing openapi generator: %w", err)
	}

	if verify {
		return verifyDiff(outputFile, outputPackagePath, outputBaseFileName)
	}

	return nil
}

// verifyDiff compares the generated file we put in the temporary directory
// with the current file in the expected location.
// It returns a diff in the error if the files are different.
func verifyDiff(currentFile, outputPackagePath, outputBaseFileName string) error {
	verifyFile := filepath.Join(outputPackagePath, outputBaseFileName+".go")

	verifyData, err := os.ReadFile(verifyFile)
	if err != nil {
		return fmt.Errorf("failed to read generated file: %w", err)
	}

	currentData, err := os.ReadFile(currentFile)
	if err != nil {
		return fmt.Errorf("failed to read current file: %w", err)
	}

	if !bytes.Equal(currentData, verifyData) {
		diff := utils.Diff(currentData, verifyData, currentFile)

		return fmt.Errorf("OpenAPI schema for %s is out of date, please regenerate the OpenAPI schema:\n%s", currentFile, diff)
	}

	return nil
}
