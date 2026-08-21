package upstreamgates

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/featuregate"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/payload-command/upstreamgates/overrides"
	"github.com/openshift/api/payload-command/upstreamgates/overrides/registry"
)

type UpstreamGatesOptions struct {
	OutputFilename  string
	OutputDirectory string
}

func (o *UpstreamGatesOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.OutputFilename, "output-filename", "zz_generated.upstreamgates.go", "sets the filename that should be used when generating the upstream gate definitions.")
	fs.StringVar(&o.OutputDirectory, "output-directory", "./features/", "sets the directory where the generated upstream gate definitions file should be created.")
}

func (o *UpstreamGatesOptions) Validate() error {
	if o.OutputFilename == "" {
		return errors.New("output-filename is required and must not be empty")
	}

	return nil
}

func (o *UpstreamGatesOptions) Run() error {
	if err := o.Validate(); err != nil {
		return err
	}

	upstreamGates := getMostRecentUpstreamGateStates()

	builderDefinitions := getOpenShiftGateBuilderDefinitionsForUpstreamGates(upstreamGates)

	return writeUpstreamGatesFile(o.OutputDirectory, o.OutputFilename, builderDefinitions...)
}

func getMostRecentUpstreamGateStates() map[featuregate.Feature]featuregate.FeatureSpec {
	upstreamVersionedGates := feature.DefaultMutableFeatureGate.GetAllVersioned()

	gateStateMap := map[featuregate.Feature]featuregate.FeatureSpec{}

	for gate, versionedSpecs := range upstreamVersionedGates {
		slices.SortFunc(versionedSpecs, func(a, b featuregate.FeatureSpec) int {
			if b.Version.LessThan(a.Version) {
				return -1
			}

			if b.Version.GreaterThan(a.Version) {
				return 1
			}

			return 0
		})

		mostRecentVersionedSpec := versionedSpecs[0]

		gateStateMap[gate] = mostRecentVersionedSpec
	}

	return gateStateMap
}

func getOpenShiftGateBuilderDefinitionsForUpstreamGates(upstreamGates map[featuregate.Feature]featuregate.FeatureSpec) []string {
	builderDefinitions := []string{}
	for gate, spec := range upstreamGates {
		featureSets := getOpenShiftFeatureSetsForUpstreamGateSpec(spec)
		gateFields := registry.FieldsForGate(gate)

		if gateFields.EnabledFeatureSets != nil {
			featureSets = sets.New(gateFields.EnabledFeatureSets...)
		}

		builderDefinition := createOpenShiftFeatureGateBuilderDefinition(gate, featureSets, gateFields)

		if len(builderDefinition) == 0 {
			continue
		}

		builderDefinitions = append(builderDefinitions, builderDefinition)
	}

	// Determinism for the gate definition ordering so the generated file doesn't
	// change every single time we perform generation.
	slices.Sort(builderDefinitions)

	return builderDefinitions
}

func getOpenShiftFeatureSetsForUpstreamGateSpec(spec featuregate.FeatureSpec) sets.Set[configv1.FeatureSet] {
	// If the feature gate is pre-alpha or locked to its default setting, we cannot modify the gate state meaning
	// we cannot have an opinion as to the enablement of the gate.
	if spec.PreRelease == featuregate.PreAlpha || spec.LockToDefault {
		return sets.New(overrides.NoOpinion)
	}

	// By defaulting to setting NoOpinion on gates that have been marked as deprecated, we will always
	// follow upstream opinions on these gates. If we identify that we need to have an opinion on a particular gate,
	// we can utilize the overrides mechanism to specify an explicit opinion that we migrate as needed.
	if spec.PreRelease == featuregate.Deprecated {
		return sets.New(overrides.NoOpinion)
	}

	if spec.Default {
		if spec.PreRelease == featuregate.GA {
			return sets.New(configv1.Default, configv1.DevPreviewNoUpgrade, configv1.TechPreviewNoUpgrade)
		}

		return sets.New(configv1.DevPreviewNoUpgrade, configv1.TechPreviewNoUpgrade)
	}

	if spec.PreRelease == featuregate.Beta {
		return sets.New(configv1.DevPreviewNoUpgrade, configv1.TechPreviewNoUpgrade)
	}

	return sets.New(configv1.DevPreviewNoUpgrade)
}

func createOpenShiftFeatureGateBuilderDefinition(name featuregate.Feature, featureSets sets.Set[configv1.FeatureSet], gateFields registry.OpenShiftGateFields) string {
	enableStrings := []string{}

	if featureSets.Has(overrides.NoOpinion) {
		fmt.Println(fmt.Sprintf("feature gate %q is deemed as needing to have no opinion, skipping creating an OpenShift equivalent feature gate definition", name))
		return ""
	}

	if featureSets.Has(configv1.Default) {
		enableStrings = append(enableStrings, "inDefault()", "inOKD()")
	}

	if featureSets.Has(configv1.TechPreviewNoUpgrade) {
		enableStrings = append(enableStrings, "inTechPreviewNoUpgrade()")
	}

	if featureSets.Has(configv1.DevPreviewNoUpgrade) {
		enableStrings = append(enableStrings, "inDevPreviewNoUpgrade()")
	}

	if len(gateFields.GroupKindResources) > 0 {
		template := "withGroupKindResources(%s)"
		gkrStrings := []string{}
		for _, gkr := range gateFields.GroupKindResources {
			gkrStrings = append(gkrStrings, fmt.Sprintf("groupKindResource{Group: %q, Kind: %q, Resource: %q}", gkr.Group, gkr.Kind, gkr.Resource))
		}
		enableStrings = append(enableStrings, fmt.Sprintf(template, strings.Join(gkrStrings, ",")))
	}

	enhancement := gateFields.EnhancementPullRequest

	if gateFields.EnhancementPullRequest != "legacyFeatureGateWithoutEnhancement" {
		enhancement = fmt.Sprintf("%q", gateFields.EnhancementPullRequest)
	}

	return fmt.Sprintf(FeatureGateBuilderTemplate, name, name, gateFields.JiraComponent, gateFields.ContactPerson, enhancement, strings.Join(enableStrings, ","))
}

func writeUpstreamGatesFile(directory string, filename string, gateDefinitions ...string) error {
	fileContents := fmt.Sprintf(FeatureGateFileTemplate, strings.Join(gateDefinitions, "\n"))

	filePath := filepath.Join(directory, filename)

	return os.WriteFile(filePath, []byte(fileContents), os.ModePerm)
}
