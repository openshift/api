package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/api/tools/codegen/pkg/sippy"
	"github.com/openshift/api/tools/codegen/pkg/utils"
)

const (
	// all features should have at least this many tests
	requiredNumberOfTests = 5

	// all variant should run at least this many times
	requiredNumberOfTestRunsPerVariant = 14

	// required pass rate.
	// nearly all current tests pass 99% of the time, but in a two week window we lack enough data to say.
	requiredPassRateOfTestsPerVariant = 0.95

	// required pass rate for "install should succeed" test
	requiredPassRateForInstallTest = 1.0
)

type FeatureGateTestAnalyzerOptions struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	CurrentFeatureSetDir  string
	PreviousFeatureSetDir string
	OutputDir             string
}

func NewFeatureGateTestAnalyzerOptions(in io.Reader, out, errOut io.Writer) *FeatureGateTestAnalyzerOptions {
	return &FeatureGateTestAnalyzerOptions{
		In:                    in,
		Out:                   out,
		ErrOut:                errOut,
		CurrentFeatureSetDir:  filepath.Join("payload-manifests", "featuregates"),
		PreviousFeatureSetDir: filepath.Join("_tmp", "previous-openshift-api", "payload-manifests", "featuregates"),
	}
}

func NewFeatureGateTestAnalyzerFlagsCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	o := NewFeatureGateTestAnalyzerOptions(in, out, errOut)

	cmd := &cobra.Command{
		Use:   "featuregate-test-analyzer",
		Short: "featuregate-test-analyzer looks to see how well tested a particular FeatureGate is.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()

			if err := o.Complete(); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			return o.Run(ctx)
		},
	}
	o.AddFlags(cmd.Flags())

	return cmd
}

func (o *FeatureGateTestAnalyzerOptions) Validate() error {
	if len(o.CurrentFeatureSetDir) == 0 {
		return fmt.Errorf("--featureset-manifest-path is required")
	}
	if len(o.PreviousFeatureSetDir) == 0 {
		return fmt.Errorf("--previous-featureset-manifest-path is required")
	}
	if _, err := os.ReadDir(o.CurrentFeatureSetDir); err != nil {
		return fmt.Errorf("--featureset-manifest-path cannot be read: %w", err)
	}
	if _, err := os.ReadDir(o.PreviousFeatureSetDir); err != nil {
		return fmt.Errorf("--previous-featureset-manifest-path cannot be read: %w", err)
	}

	return nil
}

func (o *FeatureGateTestAnalyzerOptions) Complete() error {
	artifactDir := os.Getenv("ARTIFACT_DIR")
	if len(artifactDir) > 0 {
		o.OutputDir = artifactDir
	}

	return nil
}

func (o *FeatureGateTestAnalyzerOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.CurrentFeatureSetDir, "featureset-manifest-path", o.CurrentFeatureSetDir, "path to directory containing the FeatureGate YAMLs for each FeatureGateTestAnalyzer,ClusterProfile tuple.")
	flags.StringVar(&o.PreviousFeatureSetDir, "previous-featureset-manifest-path", o.PreviousFeatureSetDir, "path to directory containing the OLD FeatureGate YAMLs for each FeatureGateTestAnalyzer,ClusterProfile tuple.")
}

func init() {
	rootCmd.AddCommand(NewFeatureGateTestAnalyzerFlagsCommand(os.Stdin, os.Stdout, os.Stderr))
}

func (o *FeatureGateTestAnalyzerOptions) Run(ctx context.Context) error {
	allCurrentClusterProfiles, _, _, currentByClusterProfileByFeatureSetTestAnalyzer, err := readFeatureGate(ctx, o.CurrentFeatureSetDir)
	if err != nil {
		return err
	}
	_, _, _, previousByClusterProfileByFeatureSetTestAnalyzer, err := readFeatureGate(ctx, o.PreviousFeatureSetDir)
	if err != nil {
		return err
	}

	md := utils.NewMarkdown("FeatureGate Promotion Summary")

	recentlyEnabledFeatureGatesToClusterProfiles := map[string]sets.Set[string]{}
	errs := []error{}
	for _, clusterProfile := range allCurrentClusterProfiles.List() {
		// we only need to check test coverage for current cluster profiles
		currentByFeatureSet := currentByClusterProfileByFeatureSetTestAnalyzer[clusterProfile]
		currentDefaultFeatureGateInfo := currentByFeatureSet["Default"]

		var previousDefaultFeatureGateInfo *featureGateInfo
		if previousByFeatureSet, ok := previousByClusterProfileByFeatureSetTestAnalyzer[clusterProfile]; ok {
			previousDefaultFeatureGateInfo = previousByFeatureSet["Default"]
		}

		currentFeatureGateNames := sets.StringKeySet(currentDefaultFeatureGateInfo.allFeatureGates)
		for _, featureGateName := range currentFeatureGateNames.List() {
			currentFeatureGateEnabled := currentDefaultFeatureGateInfo.allFeatureGates[featureGateName]
			if !currentFeatureGateEnabled {
				continue
			}

			previousFeatureGateEnabled := false
			if previousDefaultFeatureGateInfo != nil {
				previousFeatureGateEnabled = previousDefaultFeatureGateInfo.allFeatureGates[featureGateName]
			}
			if currentFeatureGateEnabled == previousFeatureGateEnabled {
				continue
			}

			// we've gone from false to true.
			if _, ok := recentlyEnabledFeatureGatesToClusterProfiles[featureGateName]; !ok {
				recentlyEnabledFeatureGatesToClusterProfiles[featureGateName] = sets.Set[string]{}
			}
			recentlyEnabledFeatureGatesToClusterProfiles[featureGateName].Insert(clusterProfile)
		}
	}

	if len(recentlyEnabledFeatureGatesToClusterProfiles) == 0 {
		md.Textf("No new Default FeatureGates found.\n")
		fmt.Fprintf(o.Out, "No new Default FeatureGates found.\n")
	}

	release, err := getRelease()
	if err != nil {
		return fmt.Errorf("couldn't determine release version: %w", err)
	}

	featureGateHTMLData := []utils.HTMLFeatureGate{}
	recentlyEnabledFeatureGates := sets.KeySet(recentlyEnabledFeatureGatesToClusterProfiles)
	for _, enabledFeatureGate := range sets.List(recentlyEnabledFeatureGates) {
		clusterProfiles := recentlyEnabledFeatureGatesToClusterProfiles[enabledFeatureGate]
		md.Title(1, enabledFeatureGate)

		testingResults, installTestLevelData, err := listTestResultFor(enabledFeatureGate, clusterProfiles)
		if err != nil {
			return err
		}

		writeTestingMarkDown(testingResults, md)

		validationResults := checkIfTestingIsSufficient(enabledFeatureGate, testingResults, installTestLevelData)

		// Separate warnings and blocking errors
		blockingErrors := []error{}
		warnings := []error{}
		var blockingResults, warningResults []ValidationResult
		for _, vr := range validationResults {
			if vr.IsWarning {
				warnings = append(warnings, vr.Error)
				warningResults = append(warningResults, vr)
			} else {
				blockingErrors = append(blockingErrors, vr.Error)
				blockingResults = append(blockingResults, vr)
			}
		}

		// For Install feature gates, report "install should succeed: overall" test statistics first
		if strings.Contains(enabledFeatureGate, "Install") {
			md.Text("")
			fmt.Fprintf(o.Out, "\n")
			md.Textf("**Install test statistics for \"install should succeed: overall\":**\n")
			fmt.Fprintf(o.Out, "Install test statistics for \"install should succeed: overall\":\n")
			jobVariants := make([]JobVariant, 0, len(testingResults))
			for jobVariant := range testingResults {
				jobVariants = append(jobVariants, jobVariant)
			}
			sort.Slice(jobVariants, func(i, j int) bool {
				return jobVariants[i].String() < jobVariants[j].String()
			})
			prevCloud := ""
			for _, jobVariant := range jobVariants {
				if prevCloud != "" && jobVariant.Cloud != prevCloud {
					md.Text("")
					fmt.Fprintf(o.Out, "\n")
				}
				prevCloud = jobVariant.Cloud
				installTest := installTestLevelData[jobVariant]
				if installTest == nil {
					md.Textf("  - %v: test not found\n", jobVariant)
					fmt.Fprintf(o.Out, "  %v: test not found\n", jobVariant)
				} else if installTest.TotalRuns > 0 {
					passPercent := float32(installTest.SuccessfulRuns) / float32(installTest.TotalRuns)
					displayActual := int(passPercent * 100)
					md.Textf("  - %v: passed %d%% (%d/%d runs)\n", jobVariant, displayActual, installTest.SuccessfulRuns, installTest.TotalRuns)
					fmt.Fprintf(o.Out, "  %v: passed %d%% (%d/%d runs)\n", jobVariant, displayActual, installTest.SuccessfulRuns, installTest.TotalRuns)
				} else {
					md.Textf("  - %v: 0 runs\n", jobVariant)
					fmt.Fprintf(o.Out, "  %v: 0 runs\n", jobVariant)
				}
			}
			md.Text("")
			fmt.Fprintf(o.Out, "\n")
		}

		if len(validationResults) == 0 {
			md.Textf("Sufficient CI testing for %q.\n", enabledFeatureGate)
			fmt.Fprintf(o.Out, "Sufficient CI testing for %q.\n", enabledFeatureGate)
		} else {
			if len(blockingErrors) > 0 {
				md.Textf("INSUFFICIENT CI testing for %q.\n", enabledFeatureGate)
				fmt.Fprintf(o.Out, "INSUFFICIENT CI testing for %q.\n", enabledFeatureGate)
			} else if len(warnings) > 0 {
				md.Textf("CI testing issues found for %q (non-blocking warnings).\n", enabledFeatureGate)
				fmt.Fprintf(o.Out, "CI testing issues found for %q (non-blocking warnings).\n", enabledFeatureGate)
			} else {
				md.Textf("Sufficient CI testing for %q.\n", enabledFeatureGate)
				fmt.Fprintf(o.Out, "Sufficient CI testing for %q.\n", enabledFeatureGate)
			}

			if len(blockingErrors) > 0 || len(warnings) > 0 {
				md.Textf("* At least five tests are expected for a feature\n")
				md.Textf("* Tests must be be run on every TechPreview platform (ask for an exception if your feature doesn't support a variant)")
				md.Textf("* All tests must run at least 14 times on every platform")
				md.Textf("* All tests must pass at least 95%% of the time")
				md.Textf("* For Install feature gates, the \"install should succeed: overall\" test must pass at least 100%% of the time")
				md.Textf("* JobTier must be one of: standard, informing, blocking, candidate (candidate is allowed but produces a warning as it is not covered by Component Readiness)\n")
				md.Text("")
			}

			if len(warnings) > 0 {
				md.Textf("**Non-blocking warnings (optional variants):**\n")
				writeGroupedValidationResults(warningResults, md)
				md.Text("")
			}

			if len(blockingErrors) > 0 {
				md.Textf("**Blocking errors:**\n")
				writeGroupedValidationResults(blockingResults, md)
				md.Text("")
			}
			md.Text("")
		}

		// Only add blocking errors to the error list (warnings don't fail the job)
		errs = append(errs, groupErrorsByCategory(blockingResults)...)
		featureGateHTMLData = append(featureGateHTMLData, buildHTMLFeatureGateData(enabledFeatureGate, testingResults, blockingErrors, release))

	}

	summaryMarkdown := md.ExactBytes()
	if len(o.OutputDir) > 0 {
		filename := filepath.Join(o.OutputDir, "feature-promotion-summary.md")
		if err := os.WriteFile(filename, summaryMarkdown, 0o644); err != nil {
			errs = append(errs, err)
		}

		htmlFilename := filepath.Join(o.OutputDir, "feature-promotion-summary.html")
		if err := writeHTMLFromTemplate(htmlFilename, featureGateHTMLData); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func topologyDisplayName(topology string) string {
	if topology == "external" {
		return "hypershift"
	}
	return topology
}

func buildHTMLFeatureGateData(name string, testingResults map[JobVariant]*TestingResults, blockingErrors []error, release string) utils.HTMLFeatureGate {
	jobVariantsSet := sets.KeySet(testingResults)
	jobVariants := OrderedJobVariants(jobVariantsSet.UnsortedList())

	sort.Sort(jobVariants)

	variants := make([]utils.HTMLVariantColumn, 0, len(jobVariants))
	for i, jv := range jobVariants {
		variants = append(variants, utils.HTMLVariantColumn{
			Topology:     topologyDisplayName(jv.Topology),
			Cloud:        jv.Cloud,
			Architecture: jv.Architecture,
			NetworkStack: jv.NetworkStack,
			OS:           jv.OS,
			JobTiers:     jv.JobTiers,
			Optional:     jv.Optional,
			ColIndex:     i + 1,
		})
	}

	allTests := sets.Set[string]{}
	for _, variantTestingResults := range testingResults {
		for _, currTestingResult := range variantTestingResults.TestResults {
			allTests.Insert(currTestingResult.TestName)
		}
	}

	tests := make([]utils.HTMLTestRow, 0, len(allTests))
	for _, testName := range sets.List(allTests) {
		row := utils.HTMLTestRow{
			TestName: testName,
			Cells:    make([]utils.HTMLTestCell, len(jobVariants)),
		}
		for i, jobVariant := range jobVariants {
			allTesting := testingResults[jobVariant]
			testResults := testResultByName(allTesting.TestResults, testName)
			cell := utils.HTMLTestCell{
				SippyURL: sippy.BuildSippyTestAnalysisURL(
					release,
					testName,
					jobVariant.Topology,
					jobVariant.Cloud,
					jobVariant.Architecture,
					jobVariant.NetworkStack,
					jobVariant.OS,
				),
			}
			if testResults == nil {
				cell.Failed = true
			} else {
				var passPercent float32
				if testResults.TotalRuns > 0 {
					passPercent = float32(testResults.SuccessfulRuns) / float32(testResults.TotalRuns)
				}
				cell.PassPercent = int(passPercent * 100)
				cell.SuccessfulRuns = testResults.SuccessfulRuns
				cell.TotalRuns = testResults.TotalRuns
				cell.FailedRuns = testResults.FailedRuns
				if testResults.TotalRuns < requiredNumberOfTestRunsPerVariant || passPercent < requiredPassRateOfTestsPerVariant {
					cell.Failed = true
				}
			}
			row.Cells[i] = cell
		}
		tests = append(tests, row)
	}

	return utils.HTMLFeatureGate{
		Name:       name,
		Sufficient: len(blockingErrors) == 0,
		Variants:   variants,
		Tests:      tests,
	}
}

func writeHTMLFromTemplate(filename string, featureGateHTMLData []utils.HTMLFeatureGate) error {
	data := utils.HTMLTemplateData{
		FeatureGates: featureGateHTMLData,
	}

	tmpl, err := template.New("report").Parse(utils.HTMLTemplateSrc)
	if err != nil {
		return fmt.Errorf("error parsing HTML template: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating HTML file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("error executing HTML template: %w", err)
	}

	return nil
}

func checkIfTestingIsSufficient(featureGate string, testingResults map[JobVariant]*TestingResults, installTestLevelData map[JobVariant]*TestResults) []ValidationResult {
	results := []ValidationResult{}

	for jobVariant, testedVariant := range testingResults {
		// Use the Optional field to determine if validation failures are warnings or errors
		// Optional variants (like RHEL 10 in 4.22) have non-blocking warnings
		isOptional := jobVariant.Optional

		// If candidate-tier queries returned results for this variant, emit a warning.
		// Candidate tier jobs are not covered by the Component Readiness main view and
		// do not have our standard regression protection mechanisms. The results are still
		// included in the pass/fail calculation alongside other tiers.
		if testedVariant.HasCandidateTierResults {
			results = append(results, ValidationResult{
				Error: fmt.Errorf("warning: variant %v includes test data from candidate-tier jobs which are not covered by Component Readiness and lack standard regression protection",
					jobVariant),
				IsWarning: true,
				Category:  CategoryCandidateTier,
			})
		}

		if len(testedVariant.TestResults) == 0 {
			continue
		}

		if len(testedVariant.TestResults) < requiredNumberOfTests {
			results = append(results, ValidationResult{
				Error: fmt.Errorf("error: only %d tests found, need at least %d for %q on %v",
					len(testedVariant.TestResults), requiredNumberOfTests, featureGate, jobVariant),
				IsWarning: isOptional,
				Category:  CategoryInsufficientTests,
			})
		}

		for _, testResults := range testedVariant.TestResults {
			// Skip "install should succeed: overall" for Install feature gates - it has special validation below
			if strings.Contains(featureGate, "Install") && testResults.TestName == "install should succeed: overall" {
				continue
			}

			if testResults.TotalRuns < requiredNumberOfTestRunsPerVariant {
				results = append(results, ValidationResult{
					Error: fmt.Errorf("error: %q only has %d runs, need at least %d runs for %q on %v",
						testResults.TestName, testResults.TotalRuns, requiredNumberOfTestRunsPerVariant, featureGate, jobVariant),
					IsWarning:  isOptional,
					Category:  CategoryInsufficientRuns,
				})
			}
			if testResults.TotalRuns == 0 {
				continue
			}
			passPercent := float32(testResults.SuccessfulRuns) / float32(testResults.TotalRuns)
			if passPercent < requiredPassRateOfTestsPerVariant {
				displayExpected := int(requiredPassRateOfTestsPerVariant * 100)
				displayActual := int(passPercent * 100)
				results = append(results, ValidationResult{
					Error: fmt.Errorf("error: %q only passed %d%%, need at least %d%% for %q on %v",
						testResults.TestName, displayActual, displayExpected, featureGate, jobVariant),
					IsWarning: isOptional,
					Category:  CategoryPassRate,
				})
			}
		}

		// For Install feature gates, validate "install should succeed: overall" using test-level data from Sippy
		if strings.Contains(featureGate, "Install") {
			installTest := installTestLevelData[jobVariant]
			if installTest == nil {
				results = append(results, ValidationResult{
					Error: fmt.Errorf("error: \"install should succeed: overall\" test data not found for Install feature gate %q on %v",
						featureGate, jobVariant),
					IsWarning: false,
					Category:  CategoryInstallTest,
				})
			} else {
				if installTest.TotalRuns < requiredNumberOfTestRunsPerVariant {
					results = append(results, ValidationResult{
						Error: fmt.Errorf("error: \"install should succeed: overall\" only has %d runs, need at least %d runs for %q on %v",
							installTest.TotalRuns, requiredNumberOfTestRunsPerVariant, featureGate, jobVariant),
						IsWarning: isOptional,
						Category:  CategoryInstallTest,
					})
				}
				if installTest.TotalRuns > 0 {
					passPercent := float32(installTest.SuccessfulRuns) / float32(installTest.TotalRuns)
					if passPercent < requiredPassRateForInstallTest {
						displayExpected := int(requiredPassRateForInstallTest * 100)
						displayActual := int(passPercent * 100)
						results = append(results, ValidationResult{
							Error: fmt.Errorf("error: \"install should succeed: overall\" only passed %d%%, need at least %d%% for %q on %v",
								displayActual, displayExpected, featureGate, jobVariant),
							IsWarning: isOptional,
							Category:  CategoryInstallTest,
						})
					}
				}
			}
		}
	}

	return results
}

func groupErrorsByCategory(results []ValidationResult) []error {
	categoryOrder := []ValidationCategory{}
	grouped := map[ValidationCategory][]ValidationResult{}
	for _, vr := range results {
		if _, seen := grouped[vr.Category]; !seen {
			categoryOrder = append(categoryOrder, vr.Category)
		}
		grouped[vr.Category] = append(grouped[vr.Category], vr)
	}
	var errs []error
	for _, cat := range categoryOrder {
		for _, vr := range grouped[cat] {
			errs = append(errs, vr.Error)
		}
	}
	return errs
}

func writeGroupedValidationResults(results []ValidationResult, md *utils.Markdown) {
	categoryOrder := []ValidationCategory{}
	grouped := map[ValidationCategory][]ValidationResult{}
	for _, vr := range results {
		if _, seen := grouped[vr.Category]; !seen {
			categoryOrder = append(categoryOrder, vr.Category)
		}
		grouped[vr.Category] = append(grouped[vr.Category], vr)
	}
	for i, cat := range categoryOrder {
		if i > 0 {
			md.Text("")
		}
		for _, vr := range grouped[cat] {
			md.Textf("  - %s\n", vr.Error.Error())
		}
	}
}

func writeTestingMarkDown(testingResults map[JobVariant]*TestingResults, md *utils.Markdown) {
	jobVariantsSet := sets.KeySet(testingResults)
	jobVariants := jobVariantsSet.UnsortedList()
	sort.Sort(OrderedJobVariants(jobVariants))

	md.NextTableColumn()
	md.Exact("Test ")
	for _, jobVariant := range jobVariants {
		md.NextTableColumn()
		columnHeader := fmt.Sprintf("%v <br/> %v <br/> %v ", topologyDisplayName(jobVariant.Topology), jobVariant.Cloud, jobVariant.Architecture)
		if jobVariant.NetworkStack != "" {
			columnHeader = columnHeader + fmt.Sprintf("<br/> %v ", jobVariant.NetworkStack)
		}
		if jobVariant.OS != "" {
			columnHeader = columnHeader + fmt.Sprintf("<br/> OS:%v ", jobVariant.OS)
		}
		if jobVariant.JobTiers != "" {
			columnHeader = columnHeader + fmt.Sprintf("<br/> Tiers:%v ", jobVariant.JobTiers)
		}
		md.Exact(columnHeader)
	}
	md.EndTableRow()
	md.NextTableColumn()
	md.Exact(":------ ")
	for i := 0; i < len(jobVariants); i++ {
		md.NextTableColumn()
		md.Exact(":---: ")
	}
	md.EndTableRow()

	allTests := sets.Set[string]{}
	for _, variantTestingResults := range testingResults {
		for _, currTestingResult := range variantTestingResults.TestResults {
			allTests.Insert(currTestingResult.TestName)
		}
	}

	for _, testName := range sets.List(allTests) {
		md.NextTableColumn()
		md.Exact(fmt.Sprintf("%s ", testName))

		for _, jobVariant := range jobVariants {
			md.NextTableColumn()
			allTesting := testingResults[jobVariant]
			testResults := testResultByName(allTesting.TestResults, testName)
			if testResults == nil {
				md.Exact(fmt.Sprintf("FAIL <br/> %d%% ( %d / %d ) ", 0, 0, 0))
				continue
			}
			failString := ""
			passPercent := float32(testResults.SuccessfulRuns) / float32(testResults.TotalRuns)
			switch {
			case testResults.TotalRuns < requiredNumberOfTestRunsPerVariant:
				failString = "FAIL <br/> "
			case passPercent < requiredPassRateOfTestsPerVariant:
				failString = "FAIL <br/> "
			}
			cellString := fmt.Sprintf("%s%d%% ( %d / %d ) ", failString, int(passPercent*100), testResults.SuccessfulRuns, testResults.TotalRuns)
			md.Exact(cellString)
		}

		md.EndTableRow()
	}
	md.Text("")
	md.Text("")
}

var (
	requiredSelfManagedJobVariants = []JobVariant{
		{
			Cloud:        "aws",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "azure",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "gcp",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "vsphere",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "ha",
			NetworkStack: "ipv4",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "ha",
			NetworkStack: "ipv6",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "ha",
			NetworkStack: "dual",
		},
		{
			Cloud:        "aws",
			Architecture: "amd64",
			Topology:     "single",
		},
		{
			Cloud:        "aws",
			Architecture: "amd64",
			Topology:     "ha",
			OS:           "rhel10",
			Optional:     true, // RHEL 10 is optional in 4.22, will be required in OCP 5
		},

		// TODO restore these once we run TechPreview jobs that contain them
		//{
		//	Cloud:        "metal-ipi",
		//	Architecture: "amd64",
		//	Topology:     "single",
		//},
	}

	// These are only checked if the feature gate is platform specific
	optionalSelfManagedPlatformVariants = []JobVariant{
		{
			Cloud:        "nutanix",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "openstack",
			Architecture: "amd64",
			Topology:     "ha",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-arbiter",
			NetworkStack: "ipv4",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-arbiter",
			NetworkStack: "ipv6",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-arbiter",
			NetworkStack: "dual",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-fencing",
			NetworkStack: "ipv4",
			JobTiers:     "candidate,standard,informing,blocking",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-fencing",
			NetworkStack: "ipv6",
			JobTiers:     "candidate,standard,informing,blocking",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "two-node-fencing",
			NetworkStack: "dual",
			JobTiers:     "candidate,standard,informing,blocking",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "single",
		},
		{
			Cloud:        "metal",
			Architecture: "amd64",
			Topology:     "compact",
		},
	}

	nonHypershiftPlatforms        = regexp.MustCompile("(?i)nutanix|metal|vsphere|openstack|azure|gcp")
	requiredHypershiftJobVariants = []JobVariant{
		{
			Cloud:        "aws",
			Architecture: "amd64",
			Topology:     "external",
		},
		// ibm and powervs?
	}
)

type JobVariant struct {
	Cloud        string
	Architecture string
	Topology     string
	NetworkStack string
	OS           string
	JobTiers     string // Comma-separated tiers (e.g., "standard,informing,blocking"). If empty, defaults to "standard,informing,blocking,candidate"
	Optional     bool   // If true, validation failures for this variant are non-blocking warnings
}

func (jv JobVariant) String() string {
	result := fmt.Sprintf("cloud=%s arch=%s topology=%s", jv.Cloud, jv.Architecture, jv.Topology)
	if jv.NetworkStack != "" {
		result += fmt.Sprintf(" network=%s", jv.NetworkStack)
	}
	if jv.OS != "" {
		result += fmt.Sprintf(" os=%s", jv.OS)
	}
	if jv.Optional {
		result += " optional=true"
	}
	return result
}

type OrderedJobVariants []JobVariant

func (a OrderedJobVariants) Len() int      { return len(a) }
func (a OrderedJobVariants) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a OrderedJobVariants) Less(i, j int) bool {
	if strings.Compare(a[i].Topology, a[j].Topology) < 0 {
		return true
	} else if strings.Compare(a[i].Topology, a[j].Topology) > 0 {
		return false
	}

	if strings.Compare(a[i].Cloud, a[j].Cloud) < 0 {
		return true
	} else if strings.Compare(a[i].Cloud, a[j].Cloud) > 0 {
		return false
	}

	if strings.Compare(a[i].Architecture, a[j].Architecture) < 0 {
		return true
	} else if strings.Compare(a[i].Architecture, a[j].Architecture) > 0 {
		return false
	}

	// Map these to an ordered list of strings so that we can define the order
	// rather than them being alphabetical.
	networkStackOrder := map[string]string{
		"":     "0",
		"ipv4": "1",
		"ipv6": "2",
		"dual": "3",
	}

	if strings.Compare(networkStackOrder[a[i].NetworkStack], networkStackOrder[a[j].NetworkStack]) < 0 {
		return true
	} else if strings.Compare(networkStackOrder[a[i].NetworkStack], networkStackOrder[a[j].NetworkStack]) > 0 {
		return false
	}

	if strings.Compare(a[i].OS, a[j].OS) < 0 {
		return true
	} else if strings.Compare(a[i].OS, a[j].OS) > 0 {
		return false
	}

	if strings.Compare(a[i].JobTiers, a[j].JobTiers) < 0 {
		return true
	} else if strings.Compare(a[i].JobTiers, a[j].JobTiers) > 0 {
		return false
	}

	return false
}

type TestingResults struct {
	JobVariant JobVariant

	TestResults             []TestResults
	HasCandidateTierResults bool // true if candidate-tier queries returned any test data
}

type TestResults struct {
	TestName       string
	TotalRuns      int
	SuccessfulRuns int
	FailedRuns     int
	FlakedRuns     int
}

type ValidationCategory string

const (
	CategoryCandidateTier    ValidationCategory = "candidate-tier"
	CategoryInsufficientTests ValidationCategory = "insufficient-tests"
	CategoryInsufficientRuns  ValidationCategory = "insufficient-runs"
	CategoryPassRate          ValidationCategory = "pass-rate"
	CategoryInstallTest       ValidationCategory = "install-test"
)

// ValidationResult represents a validation error or warning
type ValidationResult struct {
	Error     error
	IsWarning bool               // if true, this is a non-blocking warning (for optional variants)
	IsInfo    bool               // if true, this is informational telemetry (e.g., install test pass percentage)
	Category  ValidationCategory // groups related results together in output
}

func testResultByName(results []TestResults, testName string) *TestResults {
	for _, curr := range results {
		if curr.TestName == testName {
			return &curr
		}
	}
	return nil
}

func validateJobTiers(jobVariant JobVariant) error {
	if jobVariant.JobTiers == "" {
		return nil // Empty is valid - will default to standard,informing,blocking,candidate
	}

	validTiers := map[string]bool{
		"standard":  true,
		"informing": true,
		"blocking":  true,
		"candidate": true,
	}

	hasValidTier := false
	for _, tier := range strings.Split(jobVariant.JobTiers, ",") {
		tier = strings.TrimSpace(tier)
		if tier != "" {
			hasValidTier = true
			if !validTiers[tier] {
				return fmt.Errorf("invalid JobTier %q in variant %+v - must be one of: standard, informing, blocking, candidate", tier, jobVariant)
			}
		}
	}

	// Reject malformed strings like "," or " , " that contain no valid tiers
	if !hasValidTier {
		return fmt.Errorf("JobTiers string %q contains no valid tier names in variant %+v", jobVariant.JobTiers, jobVariant)
	}

	return nil
}

func listTestResultFor(featureGate string, clusterProfiles sets.Set[string]) (map[JobVariant]*TestingResults, map[JobVariant]*TestResults, error) {
	fmt.Printf("Query sippy for all test run results for feature gate %q on clusterProfile %q\n", featureGate, sets.List(clusterProfiles))

	results := map[JobVariant]*TestingResults{}
	installTestLevelData := map[JobVariant]*TestResults{}

	var jobVariantsToCheck []JobVariant
	if clusterProfiles.Has("Hypershift") && !nonHypershiftPlatforms.MatchString(featureGate) {
		jobVariantsToCheck = append(jobVariantsToCheck, filterVariants(featureGate, requiredHypershiftJobVariants)...)
	}
	if clusterProfiles.Has("SelfManagedHA") {
		// See if the feature gate is specific to any platform
		selfManagedPlatformVariants := filterVariants(featureGate, optionalSelfManagedPlatformVariants, requiredSelfManagedJobVariants)

		// If this isn't a platform specific variant, then check all required ones
		if len(selfManagedPlatformVariants) == 0 {
			selfManagedPlatformVariants = requiredSelfManagedJobVariants
		}

		jobVariantsToCheck = append(jobVariantsToCheck, selfManagedPlatformVariants...)

		// Always include metal single and compact variants as optional checks
		for _, v := range optionalSelfManagedPlatformVariants {
			if strings.ToLower(v.Cloud) == "metal" && (strings.ToLower(v.Topology) == "single" || strings.ToLower(v.Topology) == "compact") {
				jobVariantsToCheck = append(jobVariantsToCheck, v)
			}
		}
	}

	// Validate all variants before making expensive API calls
	for _, jobVariant := range jobVariantsToCheck {
		if err := validateJobTiers(jobVariant); err != nil {
			return nil, nil, err
		}
	}

	for _, jobVariant := range jobVariantsToCheck {
		jobVariantResults, err := listTestResultForVariant(featureGate, jobVariant)
		if err != nil {
			return nil, nil, err
		}
		results[jobVariant] = jobVariantResults

		// For Install feature gates, count "install should succeed: overall" results
		// directly from individual job runs, excluding infrastructure failures entirely.
		if strings.Contains(featureGate, "Install") {
			installTestData, err := getInstallTestResultsFromJobRuns(featureGate, jobVariant)
			if err != nil {
				return nil, nil, err
			}
			installTestLevelData[jobVariant] = installTestData
		}
	}

	return results, installTestLevelData, nil
}

func filterVariants(featureGate string, variantsList ...[]JobVariant) []JobVariant {
	var filteredVariants []JobVariant
	normalizedFeatureGate := strings.ToLower(featureGate)

	for _, variants := range variantsList {
		for _, variant := range variants {
			normalizedCloud := strings.ReplaceAll(strings.ToLower(variant.Cloud), "-ipi", "") // The feature gate probably won't include the install type, but some cloud variants do
			normalizedArchitecture := strings.ToLower(variant.Architecture)
			normalizedTopology := strings.ToLower(variant.Topology)

			if strings.Contains(normalizedFeatureGate, normalizedCloud) || strings.Contains(normalizedFeatureGate, normalizedArchitecture) || matchTwoNodeFeatureGates(normalizedFeatureGate, normalizedTopology) {
				filteredVariants = append(filteredVariants, variant)
			}
		}
	}

	return filteredVariants
}

// getLatestRelease returns the latest release from Sippy.
func getLatestRelease() (string, error) {
	releaseAPI := "https://sippy.dptools.openshift.org/api/releases"
	resp, err := http.Get(releaseAPI)
	if err != nil {
		return "", fmt.Errorf("error fetching data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var result struct {
		ReleaseAttrs map[string]struct {
			DevelopmentStart *time.Time `json:"development_start,omitempty"`
			Product          string     `json:"product,omitempty"`
		} `json:"release_attrs,omitempty"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	latestRelease := ""
	latestReleaseStart := time.Time{}

	for release, releaseAttrs := range result.ReleaseAttrs {
		if releaseAttrs.Product != "OCP" {
			// We only want to consider OCP releases.
			continue
		}

		if releaseAttrs.DevelopmentStart != nil && !releaseAttrs.DevelopmentStart.IsZero() && time.Now().Before(*releaseAttrs.DevelopmentStart) {
			// We only want to consider releases that have started development.
			continue
		}

		if releaseAttrs.DevelopmentStart != nil && !releaseAttrs.DevelopmentStart.IsZero() && releaseAttrs.DevelopmentStart.After(latestReleaseStart) {
			latestRelease = release
			latestReleaseStart = *releaseAttrs.DevelopmentStart
		}
	}

	if latestRelease == "" {
		return "", fmt.Errorf("no valid development releases found")
	}

	return latestRelease, nil
}

func getRelease() (string, error) {
	// if its not main branch, then use the ENV var to determine the release version
	currentRelease := os.Getenv("PULL_BASE_REF")
	if strings.Contains(currentRelease, "release-") {
		// example: release-4.18, release-4.17
		return strings.TrimPrefix(currentRelease, "release-"), nil
	}
	// means its main branch
	return getLatestRelease()
}

func getInstallTestResultsFromJobRuns(featureGate string, jobVariant JobVariant) (*TestResults, error) {
	release, err := getRelease()
	if err != nil {
		return nil, fmt.Errorf("couldn't determine release: %w", err)
	}

	defaultTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	sippyClient := &http.Client{
		Timeout:   2 * time.Minute,
		Transport: defaultTransport,
	}

	jobs, err := getJobsForFeatureGateFromSippy(sippyClient, release, featureGate, jobVariant)
	if err != nil {
		return nil, fmt.Errorf("getting jobs for install test results: %w", err)
	}

	testResults := &TestResults{
		TestName: "install should succeed: overall",
	}

	for _, job := range jobs {
		jobRuns, err := getJobRunsFromSippy(sippyClient, release, job.Name)
		if err != nil {
			return nil, fmt.Errorf("getting job runs for %q: %w", job.Name, err)
		}
		for _, jobRun := range jobRuns {
			if jobRun.OverallResult == "N" || jobRun.OverallResult == "n" {
				continue
			}
			testResults.TotalRuns++
			installFailed := false
			for _, failure := range jobRun.FailedTestNames {
				if failure == "install should succeed: overall" {
					installFailed = true
					break
				}
			}
			if installFailed {
				testResults.FailedRuns++
			} else {
				testResults.SuccessfulRuns++
			}
		}
	}

	if testResults.TotalRuns == 0 {
		return nil, nil
	}

	return testResults, nil
}

func getInstallTestLevelData(featureGate string, jobVariant JobVariant) (*TestResults, error) {
	testPattern := "install should succeed: overall"
	queries := sippy.QueriesForWithCapability(jobVariant.Cloud, jobVariant.Architecture, jobVariant.Topology,
		jobVariant.NetworkStack, jobVariant.OS, jobVariant.JobTiers, testPattern, featureGate)

	defaultTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	sippyClient := &http.Client{
		Timeout:   2 * time.Minute,
		Transport: defaultTransport,
	}

	release, err := getRelease()
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch latest release version: %w", err)
	}

	var installTestResult *TestResults
	for _, currQuery := range queries {
		currURL := &url.URL{
			Scheme: "https",
			Host:   "sippy.dptools.openshift.org",
			Path:   "api/tests",
		}
		queryParams := currURL.Query()
		queryParams.Add("release", release)
		queryParams.Add("period", "default")
		filterJSON, err := json.Marshal(currQuery)
		if err != nil {
			return nil, err
		}
		queryParams.Add("filter", string(filterJSON))
		currURL.RawQuery = queryParams.Encode()

		req, err := http.NewRequest(http.MethodGet, currURL.String(), nil)
		if err != nil {
			return nil, err
		}

		response, err := sippyClient.Do(req)
		if err != nil {
			return nil, err
		}
		if response.StatusCode < 200 || response.StatusCode > 299 {
			return nil, fmt.Errorf("error getting sippy results (status=%d) for: %v", response.StatusCode, currURL.String())
		}
		queryResultBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		response.Body.Close()

		testInfos := []sippy.SippyTestInfo{}
		if err := json.Unmarshal(queryResultBytes, &testInfos); err != nil {
			return nil, err
		}

		for _, currTest := range testInfos {
			if installTestResult == nil {
				installTestResult = &TestResults{
					TestName: currTest.Name,
				}
			}

			// Accumulate results across multiple JobTier queries
			if currTest.CurrentRuns >= requiredNumberOfTestRunsPerVariant {
				installTestResult.TotalRuns += currTest.CurrentRuns
				installTestResult.SuccessfulRuns += currTest.CurrentSuccesses
				installTestResult.FailedRuns += currTest.CurrentFailures
				installTestResult.FlakedRuns += currTest.CurrentFlakes
			} else {
				installTestResult.TotalRuns += currTest.CurrentRuns + currTest.PreviousRuns
				installTestResult.SuccessfulRuns += currTest.CurrentSuccesses + currTest.PreviousSuccesses
				installTestResult.FailedRuns += currTest.CurrentFailures + currTest.PreviousFailures
				installTestResult.FlakedRuns += currTest.CurrentFlakes + currTest.PreviousFlakes
			}
		}
	}

	return installTestResult, nil
}

func listTestResultForVariant(featureGate string, jobVariant JobVariant) (*TestingResults, error) {
	// Feature gates used by the installer don't need separate tests, use the overall install tests
	if strings.Contains(featureGate, "Install") {
		return verifyJobBasedFeatureGatePromotion(featureGate, jobVariant)
	}

	var testPattern string
	var queries []*sippy.SippyQueryStruct

	// Substring here matches for both [OCPFeatureGate:...] and [FeatureGate:...]
	testPattern = fmt.Sprintf("FeatureGate:%s]", featureGate)
	queries = sippy.QueriesFor(jobVariant.Cloud, jobVariant.Architecture, jobVariant.Topology,
		jobVariant.NetworkStack, jobVariant.OS, jobVariant.JobTiers, testPattern)
	fmt.Printf("Query sippy for all test run results for pattern %q on variant %#v\n", testPattern, jobVariant)

	defaultTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	sippyClient := &http.Client{
		Timeout:   2 * time.Minute,
		Transport: defaultTransport,
	}

	testNameToResults := map[string]*TestResults{}
	hasCandidateTierResults := false
	release, err := getRelease()
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch latest release version: %w", err)
	}

	for _, currQuery := range queries {
		currURL := &url.URL{
			Scheme: "https",
			Host:   "sippy.dptools.openshift.org",
			Path:   "api/tests",
		}
		queryParams := currURL.Query()
		queryParams.Add("release", release)
		queryParams.Add("period", "default")
		filterJSON, err := json.Marshal(currQuery)
		if err != nil {
			return nil, err
		}
		queryParams.Add("filter", string(filterJSON))
		currURL.RawQuery = queryParams.Encode()

		req, err := http.NewRequest(http.MethodGet, currURL.String(), nil)
		if err != nil {
			return nil, err
		}

		response, err := sippyClient.Do(req)
		if err != nil {
			return nil, err
		}
		if response.StatusCode < 200 || response.StatusCode > 299 {
			return nil, fmt.Errorf("error getting sippy results (status=%d) for: %v", response.StatusCode, currURL.String())
		}
		queryResultBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		response.Body.Close()

		testInfos := []sippy.SippyTestInfo{}
		if err := json.Unmarshal(queryResultBytes, &testInfos); err != nil {
			return nil, err
		}

		if currQuery.TierName == "candidate" && len(testInfos) > 0 {
			hasCandidateTierResults = true
		}

		for _, currTest := range testInfos {
			testResults, ok := testNameToResults[currTest.Name]
			if !ok {
				testResults = &TestResults{
					TestName: currTest.Name,
				}
			}

			// Try to find enough test results in the last week, but if we have to we can extend
			// the window to two weeks.
			// NOTE: Use += to accumulate results across multiple JobTier queries
			if currTest.CurrentRuns >= requiredNumberOfTestRunsPerVariant {
				testResults.TotalRuns += currTest.CurrentRuns
				testResults.SuccessfulRuns += currTest.CurrentSuccesses
				testResults.FailedRuns += currTest.CurrentFailures
				testResults.FlakedRuns += currTest.CurrentFlakes
			} else {
				fmt.Printf("Insufficient results in last 7 days, increasing lookback to 2 weeks...")
				testResults.TotalRuns += currTest.CurrentRuns + currTest.PreviousRuns
				testResults.SuccessfulRuns += currTest.CurrentSuccesses + currTest.PreviousSuccesses
				testResults.FailedRuns += currTest.CurrentFailures + currTest.PreviousFailures
				testResults.FlakedRuns += currTest.CurrentFlakes + currTest.PreviousFlakes
			}
			testNameToResults[currTest.Name] = testResults
		}
	}

	jobVariantResults := &TestingResults{
		JobVariant:              jobVariant,
		TestResults:             nil,
		HasCandidateTierResults: hasCandidateTierResults,
	}
	testNames := sets.StringKeySet(testNameToResults)
	for _, testName := range testNames.List() {
		jobVariantResults.TestResults = append(jobVariantResults.TestResults, *testNameToResults[testName])
	}

	return jobVariantResults, nil
}

// Check for Arbiter and DualReplica or Fencing featureGates as these have special topologies
func matchTwoNodeFeatureGates(featureGate string, topology string) bool {
	if (strings.Contains(featureGate, "dualreplica") || strings.Contains(featureGate, "fencing")) && strings.Contains(topology, "fencing") {
		return true
	}
	return false
}

func verifyJobBasedFeatureGatePromotion(featureGate string, jobVariant JobVariant) (*TestingResults, error) {
	ocpRelease, err := getRelease()
	if err != nil {
		return nil, fmt.Errorf("getting release version: %w", err)
	}

	defaultTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	sippyClient := &http.Client{
		Timeout:   2 * time.Minute,
		Transport: defaultTransport,
	}

	jobs, err := getJobsForFeatureGateFromSippy(sippyClient, ocpRelease, featureGate, jobVariant)
	if err != nil {
		return nil, fmt.Errorf("getting jobs for feature-gate %q for variant %v : %w", featureGate, jobVariant, err)
	}

	testResults := []TestResults{}

	for _, job := range jobs {
		results, err := verifyJobPassRate(sippyClient, ocpRelease, job, jobVariant)
		if err != nil {
			return nil, fmt.Errorf("verifying job pass rate for job %q: %w", job.Name, err)
		}

		testResults = append(testResults, *results)
	}

	return &TestingResults{
		JobVariant:  jobVariant,
		TestResults: testResults,
	}, nil
}

func verifyJobPassRate(client *http.Client, release string, job sippy.SippyJob, variant JobVariant) (*TestResults, error) {
	// Do an early check for 95% pass rate with at least 14 runs
	runs := job.CurrentRuns
	passes := job.CurrentPasses

	if runs < requiredNumberOfTestRunsPerVariant {
		fmt.Printf("Insufficient results in last 7 days, increasing lookback to 2 weeks...")
		runs += job.PreviousRuns
		passes += job.PreviousPasses
	}

	// If we have less than 14 runs, return the current set of results as-is
	// because it doesn't meet promotion criteria.
	//
	// This saves us from unnecessarily making calls out to Sippy to perform a more nuanced
	// failures analysis of the job runs to see if failed runs are true failures or known regressions.
	if runs < requiredNumberOfTestRunsPerVariant {
		return &TestResults{
			TestName: job.Name,
			TotalRuns: runs,
			SuccessfulRuns: passes,
			FailedRuns: runs - passes,
		}, nil
	}

	// If we have greater than or equal to 14 runs AND they are passing at a rate of at least 95%,
	// we can return early because this job has passed the promotion requirements.
	//
	// This saves us from unnecessarily making calls out to Sippy to perform a more nuanced
	// failures analysis of the job runs to see if failed runs are true failures or known regressions.
	if float32(passes) / float32(runs) >= requiredPassRateOfTestsPerVariant {
		return &TestResults{
			TestName: job.Name,
			TotalRuns: runs,
			SuccessfulRuns: passes,
			FailedRuns: runs - passes,
		}, nil
	}
	
	// We haven't passed promotion requirements with this job, but jobs might be impacted
	// by known regressed tests. While important to get fixed, many regressions are either
	// release blockers or require an exception to not be a release blocker.
	//
	// We can be reasonably confident in promoting a feature if the tests that are failing
	// on failed runs are only ones with known regressions for the platform being tested.
	//
	// From here on, we fetch up to the 100 most recent job runs for the job in question from Sippy,
	// fetch the known regressions for the release + platform variant, and compare failing
	// job runs failed tests with the known regressions - only counting failures that have
	// unknown test failures as a true failure.

	jobRuns, err := getJobRunsFromSippy(client, release, job.Name)
	if err != nil {
		return nil, fmt.Errorf("getting job %q results from sippy: %w", job.Name, err)
	}

	triagedTestFailures, err := getTriagedTestFailuresFromSippy(client, release, variant)
	if err != nil {
		return nil, fmt.Errorf("getting triaged test failures from sippy: %w", err)
	}

	fmt.Printf("\nIgnoring job runs that have internal or external infrastructure failures from our analysis.\n\n")

	infraFailures := 0
	testResults := &TestResults{
		TestName:  job.Name,
		TotalRuns: len(jobRuns),
	}

	for _, jobRun := range jobRuns {
		if jobRun.OverallResult == "N" || jobRun.OverallResult == "n" {
			infraFailures++
			continue
		}

		if jobRun.OverallResult == "F" && !jobRun.KnownFailure {

			untriagedTestFailures := []string{}
			for _, failure := range jobRun.FailedTestNames {
				if !triagedTestFailures.Has(failure) {
					untriagedTestFailures = append(untriagedTestFailures, failure)
				}
			}

			if len(untriagedTestFailures) > 0 {
				var writer strings.Builder
				writer.WriteString(fmt.Sprintf("job run %s has untriaged test failures:\n", jobRun.TestGridURL))
				for _, testFailure := range untriagedTestFailures {
					writer.WriteString(fmt.Sprintf("\t- %s\n", testFailure))
				}

				fmt.Println(writer.String())
				testResults.FailedRuns++

				continue
			}
		}

		testResults.SuccessfulRuns++
	}

	testResults.TotalRuns -= infraFailures

	return testResults, nil
}

func getJobsForFeatureGateFromSippy(client *http.Client, release, featureGate string, variant JobVariant) ([]sippy.SippyJob, error) {
	resp, err := client.Get(sippy.BuildSippyJobsForFeatureGateURL(featureGate, release, variant.Topology, variant.Cloud, variant.Architecture, variant.NetworkStack, variant.OS))
	if err != nil {
		return nil, fmt.Errorf("getting job info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected a 200 OK status code but got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}


	jobs := []sippy.SippyJob{}
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return jobs, nil
}

func getJobRunsFromSippy(client *http.Client, release, jobName string) ([]sippy.SippyJobRun, error) {
	resp, err := client.Get(sippy.BuildSippyJobRunsForJobURL(release, jobName, time.Now().Add(-1 * 14 * 24 * time.Hour)))
	if err != nil {
		return nil, fmt.Errorf("getting job info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected a 200 OK status code but got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}


	runResults := &sippy.SippyJobRunsResult{}
	err = json.Unmarshal(body, runResults)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return runResults.Rows, nil
}

func getTriagedTestFailuresFromSippy(client *http.Client, release string, variant JobVariant) (sets.Set[string], error) {
	reqURL, err := url.Parse("https://sippy.dptools.openshift.org/api/component_readiness/triages")
	if err != nil {
		panic(fmt.Sprintf("couldn't parse sippy triages url: %v", err))
	}

	resp, err := client.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("getting sippy triages: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected a 200 OK status code but got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	defer resp.Body.Close()

	triageItems := []sippy.SippyTriageItem{}
	err = json.Unmarshal(body, &triageItems)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %w", err)
	}

	regressedTests := sets.New[string]()

	for _, triageItem := range triageItems {
		for _, regression := range triageItem.Regressions {
			if regression.Release != release {
				continue
			}

			regressionVariants := sets.New(regression.Variants...)

			if !regressionVariants.Has(fmt.Sprintf("Platform:%s", variant.Cloud)) {
				continue
			}

			if !regressionVariants.Has(fmt.Sprintf("Topology:%s", variant.Topology)) {
				continue
			}

			if !regressionVariants.Has(fmt.Sprintf("Architecture:%s", variant.Architecture)) {
				continue
			}

			if variant.NetworkStack != "" && !regressionVariants.Has(fmt.Sprintf("NetworkStack:%s", variant.NetworkStack)) {
				continue
			}

			if variant.OS != "" && !regressionVariants.Has(fmt.Sprintf("OS:%s", variant.OS)) {
				continue
			}

			regressedTests.Insert(regression.TestName)
		}
	}

	return regressedTests, nil
}
