package openapi

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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

const generatedOpenAPI = "generated_openapi"

// generateOpenAPIDefinitions generates the OpenAPI functions for the given API package paths.
// It generates both the Go file and the JSON schema output.
func generateOpenAPIDefinitions(globalParser *parser.Parser, universe types.Universe, inputPaths []string, outputPackagePath, outputFileName, headerFilePath string, verify bool) error {
	// This is the expected path to the output file.
	// This is what we will compare against if verify is true.
	originalOutputPackagePath := outputPackagePath
	goOutputPackagePath := filepath.Join(outputPackagePath, generatedOpenAPI)
	outputFile := filepath.Join(goOutputPackagePath, outputFileName)

	if verify {
		outputPackageBase := filepath.Base(outputPackagePath)

		tmpDir, err := os.MkdirTemp("", "codegen-openapi-verify-*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		outputPackagePath = filepath.Join(tmpDir, outputPackageBase)
	}
	arguments := args.New()
	arguments.OutputDir = goOutputPackagePath
	arguments.OutputPkg = goOutputPackagePath
	arguments.OutputFile = outputFileName
	arguments.GoHeaderFile = headerFilePath

	if err := arguments.Validate(); err != nil {
		return err
	}

	klog.V(2).Infof("Generating openapi into %s", outputPackagePath)

	myTargets := func(context *gengenerator.Context) []gengenerator.Target {
		boilerplate, err := gengo.GoBoilerplate(arguments.GoHeaderFile, gengo.StdBuildTag, gengo.StdGeneratedBy)
		if err != nil {
			log.Fatalf("Failed loading boilerplate: %v", err)
		}

		return generators.GetOpenAPITargets(context, arguments, boilerplate)
	}

	if err := generation.Execute(
		globalParser,
		universe,
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		myTargets,
		inputPaths,
	); err != nil {
		return fmt.Errorf("error executing openapi generator: %w", err)
	}

	if !verify {
		// For normal generation, write JSON schema to disk
		if err := generateJSONSchema(goOutputPackagePath, outputPackagePath); err != nil {
			return fmt.Errorf("error generating JSON schema: %w", err)
		}
	} else {
		// For verification, first verify Go files, then verify JSON schema separately
		if err := verifyGoFile(outputFile, goOutputPackagePath, outputFileName); err != nil {
			return fmt.Errorf("error verifying generated openapi Go file: %w", err)
		}

		// After Go verification passes, verify JSON schema, this is all handled in memory
		return verifyJSONSchema(goOutputPackagePath, originalOutputPackagePath)
	}

	return nil
}

// verifyGoFile compares the generated Go file in the temporary directory
// with the current Go file in the expected location.
// It returns a diff in the error if the files are different.
func verifyGoFile(currentFile, tempOutputPackagePath, outputFileName string) error {
	verifyGoFile := filepath.Join(tempOutputPackagePath, outputFileName)
	verifyGoData, err := os.ReadFile(verifyGoFile)
	if err != nil {
		return fmt.Errorf("failed to read generated Go file: %w", err)
	}

	currentGoData, err := os.ReadFile(currentFile)
	if err != nil {
		return fmt.Errorf("failed to read current Go file: %w", err)
	}

	if !bytes.Equal(currentGoData, verifyGoData) {
		diff := utils.Diff(currentGoData, verifyGoData, currentFile)
		return fmt.Errorf("OpenAPI Go schema for %s is out of date, please regenerate the OpenAPI schema:\n%s", currentFile, diff)
	}

	return nil
}

// verifyJSONSchema generates JSON schema in memory and compares it with the current JSON file.
// This is run separately after Go file verification passes.
func verifyJSONSchema(schemaSourcePackage, outputPackagePath string) error {
	outputFileName := "openapi.json"
	klog.V(2).Infof("Verifying JSON schema for %s", outputFileName)

	// Generate JSON schema in memory for verification
	jsonSchemaData, err := generateJSONSchemaInMemory(schemaSourcePackage)
	if err != nil {
		return fmt.Errorf("failed to generate JSON schema for verification: %w", err)
	}

	// Verify the JSON schema file
	currentJsonFile := filepath.Join(outputPackagePath, outputFileName)
	currentJsonData, err := os.ReadFile(currentJsonFile)
	if err != nil {
		return fmt.Errorf("failed to read current JSON schema file: %w", err)
	}

	if !bytes.Equal(currentJsonData, jsonSchemaData) {
		diff := utils.Diff(currentJsonData, jsonSchemaData, currentJsonFile)
		return fmt.Errorf("OpenAPI JSON schema for %s is out of date, please regenerate the OpenAPI schema:\n%s", currentJsonFile, diff)
	}

	return nil
}

// generateJSONSchemaInMemory creates JSON schema data from the generated Go OpenAPI definitions
// and returns it as a byte slice without writing to disk
func generateJSONSchemaInMemory(schemaSourcePackage string) ([]byte, error) {
	outputFileName := "openapi.json"
	klog.V(2).Infof("Generating JSON schema in memory from %s", outputFileName)

	// Create a temporary Go program that imports the generated package and outputs JSON
	tempDir, err := os.MkdirTemp("", "openapi-json-gen-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Determine the package name from the output package path
	packageName := filepath.Base(schemaSourcePackage)

	// Create temporary main.go that calls the generated function
	tempMainContent := fmt.Sprintf(`package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openshift/api/%s"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func main() {
	refFunc := func(name string) spec.Ref {
		return spec.MustCreateRef(fmt.Sprintf("#/definitions/%%s", name))
	}

	defs := %s.GetOpenAPIDefinitions(refFunc)
	schemaDefs := make(map[string]spec.Schema, len(defs))

	for k, v := range defs {
		// Replace top-level schema with v2 if a v2 schema is embedded
		if schema, ok := v.Schema.Extensions[common.ExtensionV2Schema]; ok {
			if v2Schema, isOpenAPISchema := schema.(spec.Schema); isOpenAPISchema {
				schemaDefs[k] = v2Schema
				continue
			}
		}
		schemaDefs[k] = v.Schema
	}

	// Use a buffer and encoder to control JSON formatting
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // Don't escape HTML characters like & < >
	encoder.SetIndent("", "  ")  // Pretty print with 2-space indent

	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: schemaDefs,
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:   "OpenShift API",
					Version: "unversioned",
				},
			},
			Swagger: "2.0",
		},
	}

	if err := encoder.Encode(swagger); err != nil {
		fmt.Fprintf(os.Stderr, "error serializing api definitions: %%v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(buf.Bytes())
}`, schemaSourcePackage, packageName)

	tempMainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(tempMainFile, []byte(tempMainContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temporary main.go: %w", err)
	}

	// Get the absolute path to the repository root
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create go.mod for the temporary module
	tempGoModContent := fmt.Sprintf(`module temp-json-gen

go 1.21

replace github.com/openshift/api => %s

require github.com/openshift/api v0.0.0
`, repoRoot)
	tempGoModFile := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(tempGoModFile, []byte(tempGoModContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write temporary go.mod: %w", err)
	}

	// Run go mod tidy to fix dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	if err := tidyCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	vendorCmd := exec.Command("go", "mod", "vendor")
	vendorCmd.Dir = tempDir
	if err := vendorCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run go mod vendor: %w", err)
	}

	// Run the temporary program and capture its output
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = tempDir

	// Capture both stdout and stderr separately for better error handling
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run temporary JSON generation program: %w\nstderr: %s", err, stderr.String())
	}

	output := stdout.Bytes()
	if len(output) == 0 {
		return nil, fmt.Errorf("temporary JSON generation program produced no output\nstderr: %s", stderr.String())
	}

	return output, nil
}

// generateJSONSchema creates a JSON schema file from the generated Go OpenAPI definitions
func generateJSONSchema(schemaSourcePackage, outputPackagePath string) error {
	jsonSchemaData, err := generateJSONSchemaInMemory(schemaSourcePackage)
	if err != nil {
		return err
	}

	// Write the JSON schema to the output directory
	jsonSchemaFile := filepath.Join(outputPackagePath, "openapi.json")
	if err := os.WriteFile(jsonSchemaFile, jsonSchemaData, 0644); err != nil {
		return fmt.Errorf("failed to write OpenAPI JSON schema file: %w", err)
	}

	klog.V(1).Infof("Generated OpenAPI JSON schema: %s", jsonSchemaFile)
	return nil
}
