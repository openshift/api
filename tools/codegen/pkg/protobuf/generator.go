package protobuf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/openshift/api/tools/codegen/pkg/generation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
)

const (
	// DefaultOutputFileBaseName is the default output file base name for the generated protobuf functions.
	DefaultOutputFileBaseName = "generated.pb"

	// minimumProtocVersion is the minimum version of protoc that is supported.
	minimumProtocVersion = 23
)

// Options contains the configuration required for the protobuf generator.
type Options struct {
	// Disabled indicates whether the protobuf generator is enabled or not.
	// This defaults to true as the protobuf generator is disabled by default.
	Disabled *bool

	// DisabledVersions allows you to explicitly disable the generation of protobuf for
	// specific versions of an API.
	// This is a list of version names.
	// When omitted, no versions are disabled.
	DisabledVersions []string

	// HeaderFilePath is the path to the file containing the boilerplate header text.
	// When omitted, no header is added to the generated files.
	HeaderFilePath string

	// Verify determines whether the generator should verify the content instead
	// of updating the generated file.
	Verify bool
}

// generator implements the generation.Generator interface.
// It is designed to generate protobuf functions for a particular API group.
type generator struct {
	disabled         bool
	disabledVersions sets.String
	headerFilePath   string
	verify           bool
}

// NewGenerator builds a new protobuf generator.
func NewGenerator(opts Options) generation.Generator {
	disabled := true
	if opts.Disabled != nil {
		disabled = *opts.Disabled
	}

	return &generator{
		disabled:         disabled,
		disabledVersions: sets.NewString(opts.DisabledVersions...),
		headerFilePath:   opts.HeaderFilePath,
		verify:           opts.Verify,
	}
}

// ApplyConfig creates returns a new generator based on the configuration passed.
// If the deepcopy configuration is empty, the existing generation is returned.
func (g *generator) ApplyConfig(config *generation.Config) generation.Generator {
	if config == nil || config.Protobuf == nil {
		return g
	}

	return NewGenerator(Options{
		Disabled:         config.Protobuf.Disabled,
		DisabledVersions: config.Protobuf.DisabledVersions,
		HeaderFilePath:   config.Protobuf.HeaderFilePath,
		Verify:           g.verify,
	})
}

// Name returns the name of the generator.
func (g *generator) Name() string {
	return "go-to-protobuf"
}

// GenGroup runs the go-to-protobuf generator against the given group context.
func (g *generator) GenGroup(groupCtx generation.APIGroupContext) ([]generation.Result, error) {
	if g.disabled {
		return nil, nil
	}

	if proto := os.Getenv("PROTO_OPTIONAL"); proto != "" {
		klog.Warningf("Skipping protobuf generation: PROTO_OPTIONAL set to a non-empty value: %s", proto)
		return nil, nil
	}

	// Include the current executable dir in the PATH so that the generator can find the
	// protoc-gen-gogo binary. It's likely it was built into the same directory.
	currentExDir, err := currentExecutableDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get current executable directory: %w", err)
	}

	originalPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", currentExDir, originalPath)); err != nil {
		return nil, fmt.Errorf("failed to set PATH: %w", err)
	}

	defer func() {
		if err := os.Setenv("PATH", originalPath); err != nil {
			klog.Warningf("failed to restore PATH: %v", err)
		}
	}()

	// Check the pre-requisite binaries exist.
	if err := checkBinaries(); err != nil {
		return nil, fmt.Errorf("could not verify required binaries: set PROTO_OPTIONAL to skip protobuf generation: %w", err)
	}

	// If there is no header file, create an empty file and pass that through.
	headerFilePath := g.headerFilePath
	if headerFilePath == "" {
		tmpFile, err := os.CreateTemp("", "protobuf-header-*.txt")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary file: %w", err)
		}
		if err := tmpFile.Close(); err != nil {
			return nil, fmt.Errorf("failed to close temporary file: %w", err)
		}

		defer func() {
			if err := os.Remove(tmpFile.Name()); err != nil {
				klog.Warningf("failed to remove temporary file: %v", err)
			}
		}()

		headerFilePath = tmpFile.Name()
	}

	for _, version := range groupCtx.Versions {
		if g.disabledVersions.Has(version.Name) {
			klog.V(1).Infof("Skipping generation of protobuf functions for %s/%s", groupCtx.Name, version.Name)
			continue
		}

		action := "Generating"
		if g.verify {
			action = "Verifying"
		}

		klog.V(1).Infof("%s protobuf functions for for %s/%s", action, groupCtx.Name, version.Name)

		if err := generateProtobufFunctions(version.Path, version.PackagePath, headerFilePath, g.verify); err != nil {
			return nil, fmt.Errorf("could not generate protobuf functions for %s/%s: %w", groupCtx.Name, version.Name, err)
		}
	}

	return nil, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check that the protoc version is at least 3.0.0.
	// The generator will fail with a cryptic error if the version is too old.
	out, err := exec.CommandContext(ctx, protoc, "--version").Output()
	if err != nil {
		return fmt.Errorf("failed to get protoc version: %w", err)
	}

	// The output is of the form "libprotoc 3.0.0".
	// We only care about the version number.
	versionStrings := strings.Split(string(out), " ")
	if len(versionStrings) < 2 {
		return fmt.Errorf("failed to get protoc version: unexpected output format: %s", string(out))
	}

	version := versionStrings[1]
	if version == "" {
		return errors.New("failed to get protoc version")
	}

	// Check that the major version is at least 3.
	// The minor version is not important.
	majorVersion, err := strconv.Atoi(strings.Split(version, ".")[0])
	if err != nil {
		return fmt.Errorf("failed to parse protoc version: %w", err)
	}

	if majorVersion < minimumProtocVersion {
		return fmt.Errorf("protoc version %s is too old, version %d.0.0 or newer is required", version, minimumProtocVersion)
	}

	return nil
}
