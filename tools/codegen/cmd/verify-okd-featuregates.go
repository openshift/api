package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "sigs.k8s.io/yaml"
)

type verifyOKDFeatureGatesOptions struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	FeatureSetManifestDir string
}

func newVerifyOKDFeatureGatesOptions(in io.Reader, out, errOut io.Writer) *verifyOKDFeatureGatesOptions {
	return &verifyOKDFeatureGatesOptions{
		In:                    in,
		Out:                   out,
		ErrOut:                errOut,
		FeatureSetManifestDir: filepath.Join("payload-manifests", "featuregates"),
	}
}

func NewVerifyOKDFeatureGatesCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	o := newVerifyOKDFeatureGatesOptions(in, out, errOut)

	cmd := &cobra.Command{
		Use:   "verify-okd-featuregates",
		Short: "verify-okd-featuregates verifies that all featuregates enabled in Default are also enabled in OKD",
		Long: `This verifier ensures that the OKD featureset includes all featuregates that are enabled
in the Default featureset. OKD may have additional featuregates beyond Default (e.g., from
TechPreviewNoUpgrade or DevPreviewNoUpgrade), but it must not be missing any Default featuregates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run(ctx)
		},
	}
	o.AddFlags(cmd.Flags())

	return cmd
}

func init() {
	rootCmd.AddCommand(NewVerifyOKDFeatureGatesCommand(os.Stdin, os.Stdout, os.Stderr))
}

func (o *verifyOKDFeatureGatesOptions) Validate() error {
	if len(o.FeatureSetManifestDir) == 0 {
		return fmt.Errorf("--featureset-manifest-path is required")
	}
	if _, err := os.ReadDir(o.FeatureSetManifestDir); err != nil {
		return fmt.Errorf("--featureset-manifest-path cannot be read: %w", err)
	}
	return nil
}

func (o *verifyOKDFeatureGatesOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.FeatureSetManifestDir, "featureset-manifest-path", o.FeatureSetManifestDir, "path to directory containing the FeatureGate YAMLs for each FeatureSet,ClusterProfile tuple.")
}

func (o *verifyOKDFeatureGatesOptions) Run(ctx context.Context) error {
	// Read all featuregate manifests
	featureSetsByProfile, err := readFeatureGateManifests(o.FeatureSetManifestDir)
	if err != nil {
		return err
	}

	allErrors := []string{}

	// Check each cluster profile
	for profile, featureSets := range featureSetsByProfile {
		defaultGates, hasDefault := featureSets["Default"]
		okdGates, hasOKD := featureSets["OKD"]

		// If OKD doesn't exist for this profile, skip
		if !hasOKD {
			continue
		}

		// If Default doesn't exist for this profile, skip
		if !hasDefault {
			continue
		}

		// Check that all Default featuregates are in OKD
		missingInOKD := defaultGates.Difference(okdGates)

		if missingInOKD.Len() > 0 {
			missingList := missingInOKD.List()
			sort.Strings(missingList)

			errorMsg := fmt.Sprintf(
				"ERROR: ClusterProfile %q: OKD featureset is missing %d featuregate(s) that are enabled in Default:\n  - %s\n\nAll featuregates enabled in Default must also be enabled in OKD.",
				profile,
				missingInOKD.Len(),
				strings.Join(missingList, "\n  - "),
			)
			allErrors = append(allErrors, errorMsg)
		}
	}

	if len(allErrors) > 0 {
		fmt.Fprintln(o.ErrOut, strings.Join(allErrors, "\n\n"))
		return fmt.Errorf("OKD featuregate verification failed")
	}

	return nil
}

// readFeatureGateManifests reads the featuregate manifests and returns a map of
// cluster profile -> feature set -> enabled featuregates
func readFeatureGateManifests(manifestDir string) (map[string]map[string]sets.String, error) {
	result := map[string]map[string]sets.String{}

	files, err := os.ReadDir(manifestDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read manifest dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(manifestDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("unable to read %q: %w", filePath, err)
		}

		// Parse the YAML
		obj := map[string]interface{}{}
		if err := kyaml.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("unable to parse %q: %w", filePath, err)
		}
		uncastObj := unstructured.Unstructured{Object: obj}

		// Get cluster profile from annotations
		profile := getClusterProfile(uncastObj.GetAnnotations())
		if profile == "" {
			continue // Skip if no profile found
		}

		// Get feature set name
		featureSet, _, _ := unstructured.NestedString(obj, "spec", "featureSet")
		if featureSet == "" {
			featureSet = "Default"
		}

		// Get enabled featuregates
		enabledGates := sets.NewString()
		featureGateSlice, _, err := unstructured.NestedSlice(obj, "status", "featureGates")
		if err == nil && len(featureGateSlice) > 0 {
			enabledList, _, err := unstructured.NestedSlice(featureGateSlice[0].(map[string]interface{}), "enabled")
			if err == nil {
				for _, gate := range enabledList {
					name, _, _ := unstructured.NestedString(gate.(map[string]interface{}), "name")
					if name != "" {
						enabledGates.Insert(name)
					}
				}
			}
		}

		// Store in result
		if _, ok := result[profile]; !ok {
			result[profile] = map[string]sets.String{}
		}
		result[profile][featureSet] = enabledGates
	}

	return result, nil
}

// getClusterProfile extracts a simplified cluster profile name from annotations
func getClusterProfile(annotations map[string]string) string {
	for k, v := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") && v == "false-except-for-the-config-operator" {
			// Extract short name from annotation
			if strings.Contains(k, "self-managed-high-availability") {
				return "SelfManagedHA"
			}
			if strings.Contains(k, "ibm-cloud-managed") {
				return "Hypershift"
			}
		}
	}
	return ""
}
