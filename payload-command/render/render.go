package render

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	configv1 "github.com/openshift/api/config/v1"
	assets "github.com/openshift/api/payload-command/render/renderassets"
	"k8s.io/apimachinery/pkg/util/sets"
)

// RenderOpts holds values to drive the render command.
type RenderOpts struct {
	ImageProvidedManifestDir      string
	RenderedManifestInputFilename string
	PayloadVersion                string
	AssetOutputDir                string
	ClusterProfile                string
}

func (o *RenderOpts) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.RenderedManifestInputFilename, "rendered-manifest-dir", o.RenderedManifestInputFilename,
		"files or directories containing yaml or json manifests that will be created via cluster-bootstrapping.")
	fs.StringVar(&o.ImageProvidedManifestDir, "image-manifests", o.ImageProvidedManifestDir, "Directory containing the manifest templates provided by the image.")
	fs.StringVar(&o.PayloadVersion, "payload-version", o.PayloadVersion, "Version that will eventually be placed into ClusterOperator.status.  This normally comes from the CVO set via env var: OPERATOR_IMAGE_VERSION.")
	fs.StringVar(&o.AssetOutputDir, "asset-output-dir", o.AssetOutputDir, "Output path for rendered manifests.")
	fs.StringVar(&o.ClusterProfile, "cluster-profile", o.ClusterProfile, "self-managed-high-availability, single-node-developer, ibm-cloud-managed")
}

// Validate verifies the inputs.
func (o *RenderOpts) Validate() error {
	switch o.ClusterProfile {
	case "":
		// to be disallowed soonish
	case "self-managed-high-availability", "single-node-developer", "ibm-cloud-managed":
		// ok
	default:
		return fmt.Errorf("--cluster-profile must be one of self-managed-high-availability, single-node-developer, ibm-cloud-managed")
	}

	return nil
}

// Complete fills in missing values before command execution.
func (o *RenderOpts) Complete() error {
	return nil
}

// Run contains the logic of the render command.
func (o *RenderOpts) Run() error {
	featureSet := ""
	featureGateFiles, err := featureGateManifests([]string{o.RenderedManifestInputFilename})
	if err != nil {
		return fmt.Errorf("problem with featuregate manifests: %w", err)
	}
	for _, featureGateFile := range featureGateFiles {
		uncastObj, err := featureGateFile.GetDecodedObj()
		if err != nil {
			return fmt.Errorf("error decoding FeatureGate: %w", err)
		}
		featureGates := &configv1.FeatureGate{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uncastObj.(*unstructured.Unstructured).Object, featureGates)
		if err != nil {
			return fmt.Errorf("error converting FeatureGate: %w", err)
		}
		currentDetails, err := FeaturesGateDetailsFromFeatureSets(configv1.FeatureSets, featureGates, o.PayloadVersion)
		if err != nil {
			return fmt.Errorf("error determining FeatureGates: %w", err)
		}
		featureGates.Status.FeatureGates = []configv1.FeatureGateDetails{*currentDetails}

		featureGateOutBytes := writeFeatureGateV1OrDie(featureGates)
		if err := os.WriteFile(featureGateFile.OriginalFilename, []byte(featureGateOutBytes), 0644); err != nil {
			return fmt.Errorf("error writing FeatureGate manifest: %w", err)
		}
		featureSet = string(featureGates.Spec.FeatureSet)
	}

	err = assets.SubstituteAndCopyFiles(
		o.ImageProvidedManifestDir,
		filepath.Join(o.AssetOutputDir, "manifests"),
		featureSet,
		o.ClusterProfile,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to substitute and copy files: %w", err)
	}

	return nil
}

func featureGateManifests(renderedManifestInputFilenames []string) (assets.RenderedManifests, error) {
	if len(renderedManifestInputFilenames) == 0 {
		return nil, fmt.Errorf("cannot return FeatureGate without rendered manifests")
	}

	inputManifests := assets.RenderedManifests{}
	for _, filename := range renderedManifestInputFilenames {
		manifestContent, err := assets.LoadFilesRecursively(filename)
		if err != nil {
			return nil, fmt.Errorf("failed loading rendered manifest inputs from %q: %w", filename, err)
		}
		for manifestFile, content := range manifestContent {
			inputManifests = append(inputManifests, assets.RenderedManifest{
				OriginalFilename: filepath.Join(filename, manifestFile),
				Content:          content,
			})
		}
	}
	featureGates := inputManifests.ListManifestOfType(configv1.GroupVersion.WithKind("FeatureGate"))
	if len(featureGates) == 0 {
		return nil, fmt.Errorf("no FeatureGates found in manfest dir: %v", renderedManifestInputFilenames)
	}

	return featureGates, nil
}

func FeaturesGateDetailsFromFeatureSets(featureSetMap map[configv1.FeatureSet]*configv1.FeatureGateEnabledDisabled, featureGates *configv1.FeatureGate, currentVersion string) (*configv1.FeatureGateDetails, error) {
	enabled, disabled, err := featuresGatesFromFeatureSets(featureSetMap, featureGates)
	if err != nil {
		return nil, err
	}
	currentDetails := configv1.FeatureGateDetails{
		Version: currentVersion,
	}
	for _, gateName := range enabled {
		currentDetails.Enabled = append(currentDetails.Enabled, configv1.FeatureGateAttributes{
			Name: gateName,
		})
	}
	for _, gateName := range disabled {
		currentDetails.Disabled = append(currentDetails.Disabled, configv1.FeatureGateAttributes{
			Name: gateName,
		})
	}

	// sort for stability
	sort.Sort(byName(currentDetails.Enabled))
	sort.Sort(byName(currentDetails.Disabled))

	return &currentDetails, nil
}

type byName []configv1.FeatureGateAttributes

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool {
	if strings.Compare(string(a[i].Name), string(a[j].Name)) < 0 {
		return true
	}
	return false
}

func featuresGatesFromFeatureSets(knownFeatureSets map[configv1.FeatureSet]*configv1.FeatureGateEnabledDisabled, featureGates *configv1.FeatureGate) ([]configv1.FeatureGateName, []configv1.FeatureGateName, error) {
	if featureGates.Spec.FeatureSet == configv1.CustomNoUpgrade {
		if featureGates.Spec.FeatureGateSelection.CustomNoUpgrade != nil {
			completeEnabled, completeDisabled := completeFeatureGates(knownFeatureSets, featureGates.Spec.FeatureGateSelection.CustomNoUpgrade.Enabled, featureGates.Spec.FeatureGateSelection.CustomNoUpgrade.Disabled)
			return completeEnabled, completeDisabled, nil
		}
		return []configv1.FeatureGateName{}, []configv1.FeatureGateName{}, nil
	}

	featureSet, ok := knownFeatureSets[featureGates.Spec.FeatureSet]
	if !ok {
		return []configv1.FeatureGateName{}, []configv1.FeatureGateName{}, fmt.Errorf(".spec.featureSet %q not found", featureSet)
	}

	completeEnabled, completeDisabled := completeFeatureGates(knownFeatureSets, toFeatureGateNames(featureSet.Enabled), toFeatureGateNames(featureSet.Disabled))
	return completeEnabled, completeDisabled, nil
}

func toFeatureGateNames(in []configv1.FeatureGateDescription) []configv1.FeatureGateName {
	out := []configv1.FeatureGateName{}
	for _, curr := range in {
		out = append(out, curr.FeatureGateAttributes.Name)
	}

	return out
}

// completeFeatureGates identifies every known feature and ensures that is explicitly on or explicitly off
func completeFeatureGates(knownFeatureSets map[configv1.FeatureSet]*configv1.FeatureGateEnabledDisabled, enabled, disabled []configv1.FeatureGateName) ([]configv1.FeatureGateName, []configv1.FeatureGateName) {
	specificallyEnabledFeatureGates := sets.New[configv1.FeatureGateName]()
	specificallyEnabledFeatureGates.Insert(enabled...)

	knownFeatureGates := sets.New[configv1.FeatureGateName]()
	knownFeatureGates.Insert(enabled...)
	knownFeatureGates.Insert(disabled...)
	for _, known := range knownFeatureSets {
		for _, curr := range known.Disabled {
			knownFeatureGates.Insert(curr.FeatureGateAttributes.Name)
		}
		for _, curr := range known.Enabled {
			knownFeatureGates.Insert(curr.FeatureGateAttributes.Name)
		}
	}

	return enabled, knownFeatureGates.Difference(specificallyEnabledFeatureGates).UnsortedList()
}
