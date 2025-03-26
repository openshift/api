package protobuf

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"k8s.io/code-generator/cmd/go-to-protobuf/protobuf"
	"k8s.io/klog/v2"
)

const (
	defaultDropEmbeddedFields   = "k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta"
	defaultAPIMachineryPackages = "-k8s.io/apimachinery/pkg/util/intstr,-k8s.io/apimachinery/pkg/api/resource,-k8s.io/apimachinery/pkg/runtime/schema,-k8s.io/apimachinery/pkg/runtime,-k8s.io/apimachinery/pkg/apis/meta/v1,-k8s.io/apimachinery/pkg/apis/meta/v1beta1,-k8s.io/api/core/v1,-k8s.io/api/rbac/v1"
)

// generateProtobufFunctions generates the Protobuf functions for the given API package paths.
func generateProtobufFunctions(path, packagePath, headerFilePath string, verify bool) error {
	// Include the current executable dir in the PATH so that the generator can find the
	// protoc-gen-gogo binary. It's likely it was built into the same directory.
	currentExDir, err := currentExecutableDir()
	if err != nil {
		return fmt.Errorf("failed to get current executable directory: %w", err)
	}
	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s", currentExDir, originalPath))
	defer os.Setenv("PATH", originalPath)

	// Check the pre-requisite binaries exist.
	if err := checkBinaries(); err != nil {
		if os.Getenv("PROTO_OPTIONAL") != "" {
			klog.Warningf("Skipping protobuf generation: %v", err)
			return nil
		}

		return fmt.Errorf("could not verify required binaries: set PROTO_OPTIONAL to skip protobuf generation: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	pathPrefix, err := getPathPrefix(wd, path, packagePath)
	if err != nil {
		return fmt.Errorf("failed to get path prefix: %w", err)
	}

	vendor, err := findDirectoryInParent(path, "vendor")
	if err != nil {
		return fmt.Errorf("failed to find vendor directory: %w", err)
	}

	generator := protobuf.Generator{
		GoHeaderFile:         headerFilePath,
		APIMachineryPackages: defaultAPIMachineryPackages,
		DropEmbeddedFields:   defaultDropEmbeddedFields,
		Packages:             packagePath,
		OutputDir:            strings.TrimSuffix(wd, pathPrefix),
		ProtoImport:          []string{vendor},
	}

	// The generator doesn't return an error an instead `log.Fatal`s when there's an issue.
	// Short of rewriting it, this is the best we can do.
	protobuf.Run(&generator)

	return nil
}

// getPathPrefix calculates the pathPrefix that needs to be trimmed from the current working directory.
// The generator will generate the output file to the current working directory plus the
// package path name.
// This function calculates what is needed to be trimmed from the working directory path name to
// make sure the output ends up in the correct directory.
// Eg. if the package is github.com/openshift.io/api/machine/v1,
//   - the current working directory is /home/user/go/src, then the path would be
//     github.com/openshift.io/api/machine/v1 and so the output would be the empty string.
//   - the current working directory is /home/user/go/src/github.com/openshift.io/api, then
//     the path would be machine/v1 and so the output would be github.com/openshift.io/api.
func getPathPrefix(wd, path, packagePath string) (string, error) {
	relPath, err := filepath.Rel(wd, path)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	if strings.HasPrefix(relPath, "../") {
		return "", errors.New("cannot generate deepcopy functions for a path outside of the working directory")
	}

	if !strings.HasSuffix(packagePath, relPath) {
		return "", fmt.Errorf("package path %s does not match with input path %s, expected package path to end with input path", packagePath, relPath)
	}

	return filepath.Clean(strings.TrimSuffix(packagePath, relPath)), nil
}

// findDirectoryInParent finds the directory with the given name in the parent directories of the given path.
// It returns the path to the directory if found, otherwise it returns an error.
func findDirectoryInParent(path string, target string) (string, error) {
	for path != "." && path != "/" {
		if _, err := os.Stat(filepath.Join(path, target)); err == nil {
			return filepath.Join(path, target), nil
		}
		path = filepath.Dir(path)
	}

	return "", errors.New("could not find vendor directory")
}

// currentExecutableDir returns the absolute path to the directory containing the current executable.
func currentExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Dir(ex))
}

// checkBinaries checks that the required binaries are available.
// It returns an error if any of the binaries are missing.
// It looks for both the protoc and protoc-gen-gogo binaries.
// It will also check that protoc is version 3.0.0 or higher.
func checkBinaries() error {
	if _, err := exec.LookPath("protoc-gen-gogo"); err != nil {
		return errors.New("protoc-gen-gogo is required to generate protobuf files")
	}

	protoc, err := exec.LookPath("protoc")
	if err != nil {
		return errors.New("protoc is required to generate protobuf files")
	}

	// Check that the protoc version is at least 3.0.0.
	// The generator will fail with a cryptic error if the version is too old.
	out, err := exec.Command(protoc, "--version").Output()
	if err != nil {
		return fmt.Errorf("failed to get protoc version: %w", err)
	}

	// The output is of the form "libprotoc 3.0.0".
	// We only care about the version number.
	version := strings.Split(string(out), " ")[1]
	if version == "" {
		return errors.New("failed to get protoc version")
	}

	// Check that the major version is at least 3.
	// The minor version is not important.
	majorVersion, err := strconv.Atoi(strings.Split(version, ".")[0])
	if err != nil {
		return fmt.Errorf("failed to parse protoc version: %w", err)
	}

	if majorVersion < 3 {
		return fmt.Errorf("protoc version %s is too old, version 3.0.0 or newer is required", version)
	}

	return nil
}
