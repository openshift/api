package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/api/tools/codegen/pkg/sippy"
)

func Test_listTestResultFor(t *testing.T) {
	type args struct {
		clusterProfile string
		featureGate    string
	}
	tests := []struct {
		name    string
		args    args
		want    map[JobVariant]*TestingResults
		wantErr bool
	}{
		{
			name: "test example",
			args: args{
				clusterProfile: "SelfManagedHA",
				featureGate:    "Example",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "platform example",
			args: args{
				clusterProfile: "VSphereGate",
				featureGate:    "Example",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "optional platform example",
			args: args{
				featureGate:    "NutanixGate",
				clusterProfile: "SelfManagedHA",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "install example",
			args: args{
				featureGate:    "FooBarInstall",
				clusterProfile: "Example",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("this is for ease of manual testing")

			got, _, err := listTestResultFor(context.Background(), tt.args.featureGate, sets.New[string](tt.args.clusterProfile))
			if (err != nil) != tt.wantErr {
				t.Errorf("listTestResultFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				for _, jobVariantResult := range got {
					result, serializationErr := json.MarshalIndent(jobVariantResult, "", "  ")
					if serializationErr != nil {
						t.Log(serializationErr.Error())
					}
					t.Log(string(result))
				}
				t.Errorf("listTestResultFor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateJobTiers_candidate(t *testing.T) {
	tests := []struct {
		name    string
		variant JobVariant
		wantErr bool
	}{
		{
			name:    "candidate is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "candidate"},
			wantErr: false,
		},
		{
			name:    "candidate with standard is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "standard,candidate"},
			wantErr: false,
		},
		{
			name:    "invalid tier still rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "bogus"},
			wantErr: true,
		},
		{
			name:    "candidate with invalid tier rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "candidate,bogus"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJobTiers(tt.variant)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJobTiers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkIfTestingIsSufficient_CandidateVariants(t *testing.T) {
	sufficientTests := []TestResults{
		{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
		{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
		{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
		{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
		{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
	}
	insufficientTests := []TestResults{
		{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
		{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
		// Only 2 tests, need 5
	}

	tests := []struct {
		name               string
		featureGate        string
		testingResults     map[JobVariant]*TestingResults
		wantBlockingErrors int
		wantWarnings       int
	}{
		{
			name:        "candidate tier returned results with sufficient tests - warning about component readiness",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults:             sufficientTests,
					HasCandidateTierResults: true,
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       1, // component readiness warning
		},
		{
			name:        "candidate tier returned results with insufficient tests - blocking error plus warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults:             insufficientTests,
					HasCandidateTierResults: true,
				},
			},
			wantBlockingErrors: 1, // insufficient tests is still blocking
			wantWarnings:       1, // component readiness warning
		},
		{
			name:        "no candidate tier results - no warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults:             sufficientTests,
					HasCandidateTierResults: false,
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       0,
		},
		{
			name:        "candidate tier returned results with low pass rate - blocking error plus warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 20, SuccessfulRuns: 18}, // 90%
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
					HasCandidateTierResults: true,
				},
			},
			wantBlockingErrors: 1, // low pass rate is still blocking
			wantWarnings:       1, // component readiness warning
		},
		{
			name:        "mix of variants - one with candidate results one without",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults:             sufficientTests,
					HasCandidateTierResults: false,
				},
				{
					Cloud:        "gcp",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults:             sufficientTests,
					HasCandidateTierResults: true,
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       1, // only gcp variant has candidate warning
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := checkIfTestingIsSufficient(tt.featureGate, tt.testingResults, nil)

			blockingErrors := 0
			warnings := 0
			for _, result := range results {
				if result.IsWarning {
					warnings++
				} else {
					blockingErrors++
				}
			}

			if blockingErrors != tt.wantBlockingErrors {
				t.Errorf("got %d blocking errors, want %d", blockingErrors, tt.wantBlockingErrors)
				for _, result := range results {
					if !result.IsWarning {
						t.Logf("  Blocking error: %v", result.Error)
					}
				}
			}
			if warnings != tt.wantWarnings {
				t.Errorf("got %d warnings, want %d", warnings, tt.wantWarnings)
				for _, result := range results {
					if result.IsWarning {
						t.Logf("  Warning: %v", result.Error)
					}
				}
			}
		})
	}
}

func Test_checkIfTestingIsSufficient_OptionalVariants(t *testing.T) {
	tests := []struct {
		name               string
		featureGate        string
		testingResults     map[JobVariant]*TestingResults
		wantBlockingErrors int
		wantWarnings       int
	}{
		{
			name:        "required variant with insufficient tests - should be blocking error",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					Optional:     false,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						// Only 2 tests, need 5
					},
				},
			},
			wantBlockingErrors: 1,
			wantWarnings:       0,
		},
		{
			name:        "optional variant with insufficient tests - should be warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					OS:           "rhel10",
					Optional:     true,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						// Only 2 tests, need 5
					},
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       1,
		},
		{
			name:        "required variant with insufficient runs - should be blocking error",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					Optional:     false,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 10, SuccessfulRuns: 10}, // Only 10 runs, need 14
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			wantBlockingErrors: 1,
			wantWarnings:       0,
		},
		{
			name:        "optional variant with insufficient runs - should be warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					OS:           "rhel10",
					Optional:     true,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 10, SuccessfulRuns: 10}, // Only 10 runs, need 14
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       1,
		},
		{
			name:        "required variant with low pass rate - should be blocking error",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					Optional:     false,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 20, SuccessfulRuns: 18}, // 90% pass rate, need 95%
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			wantBlockingErrors: 1,
			wantWarnings:       0,
		},
		{
			name:        "optional variant with low pass rate - should be warning",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					OS:           "rhel10",
					Optional:     true,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 20, SuccessfulRuns: 18}, // 90% pass rate, need 95%
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			wantBlockingErrors: 0,
			wantWarnings:       1,
		},
		{
			name:        "mix of required and optional variants - both have issues",
			featureGate: "TestFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					Optional:     false,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						// Only 2 tests, need 5
					},
				},
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
					OS:           "rhel10",
					Optional:     true,
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						// Only 2 tests, need 5
					},
				},
			},
			wantBlockingErrors: 1,
			wantWarnings:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := checkIfTestingIsSufficient(tt.featureGate, tt.testingResults, nil)

			blockingErrors := 0
			warnings := 0
			for _, result := range results {
				if result.IsWarning {
					warnings++
				} else {
					blockingErrors++
				}
			}

			if blockingErrors != tt.wantBlockingErrors {
				t.Errorf("checkIfTestingIsSufficient() got %d blocking errors, want %d", blockingErrors, tt.wantBlockingErrors)
				for _, result := range results {
					if !result.IsWarning {
						t.Logf("  Blocking error: %v", result.Error)
					}
				}
			}
			if warnings != tt.wantWarnings {
				t.Errorf("checkIfTestingIsSufficient() got %d warnings, want %d", warnings, tt.wantWarnings)
				for _, result := range results {
					if result.IsWarning {
						t.Logf("  Warning: %v", result.Error)
					}
				}
			}
		})
	}
}

func Test_checkIfTestingIsSufficient_InstallFeatureGates(t *testing.T) {
	tests := []struct {
		name               string
		featureGate        string
		testingResults     map[JobVariant]*TestingResults
		installTestData    map[JobVariant]*TestResults
		wantBlockingErrors int
		wantWarnings       int
	}{
		{
			name:        "Install feature gate with install should succeed: overall test failing 100% requirement",
			featureGate: "FakeInstallFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "metal",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 19},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData: map[JobVariant]*TestResults{
				{Cloud: "metal", Architecture: "amd64", Topology: "ha"}: {TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 19},
			},
			wantBlockingErrors: 1, // Blocking error: install should succeed: overall must pass at 100%
			wantWarnings:       0,
		},
		{
			name:        "Install feature gate without install should succeed: overall test - warning reported",
			featureGate: "MockInstallGate",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "metal",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test5", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData:    map[JobVariant]*TestResults{},
			wantBlockingErrors: 1, // Blocking error: install test data not found on required variant
			wantWarnings:       0,
		},
		{
			name:        "Non-Install feature gate with install should succeed: overall test - no special validation",
			featureGate: "SomeOtherFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 19},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData:    nil,
			wantBlockingErrors: 0,
			wantWarnings:       0,
		},
		{
			name:        "Install feature gate with 100% pass rate - no errors",
			featureGate: "FakeInstallFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 20},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData: map[JobVariant]*TestResults{
				{Cloud: "aws", Architecture: "amd64", Topology: "ha"}: {TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 20},
			},
			wantBlockingErrors: 0,
			wantWarnings:       0,
		},
		{
			name:        "Install feature gate with multiple variants - one fails 100% requirement",
			featureGate: "FakeInstallFeature",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "metal",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 19},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
				{
					Cloud:        "metal",
					Architecture: "amd64",
					Topology:     "single",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 18, SuccessfulRuns: 18},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData: map[JobVariant]*TestResults{
				{Cloud: "metal", Architecture: "amd64", Topology: "ha"}:     {TestName: "install should succeed: overall", TotalRuns: 20, SuccessfulRuns: 19},
				{Cloud: "metal", Architecture: "amd64", Topology: "single"}: {TestName: "install should succeed: overall", TotalRuns: 18, SuccessfulRuns: 18},
			},
			wantBlockingErrors: 1, // One variant (ha) fails 100% requirement
			wantWarnings:       0,
		},
		{
			name:        "Install feature gate with insufficient runs for install test - blocking error",
			featureGate: "MockInstallGate",
			testingResults: map[JobVariant]*TestingResults{
				{
					Cloud:        "aws",
					Architecture: "amd64",
					Topology:     "ha",
				}: {
					TestResults: []TestResults{
						{TestName: "install should succeed: overall", TotalRuns: 10, SuccessfulRuns: 10},
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test3", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test4", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			installTestData: map[JobVariant]*TestResults{
				{Cloud: "aws", Architecture: "amd64", Topology: "ha"}: {TestName: "install should succeed: overall", TotalRuns: 10, SuccessfulRuns: 10},
			},
			wantBlockingErrors: 1, // Blocking error for insufficient runs (< 14)
			wantWarnings:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := checkIfTestingIsSufficient(tt.featureGate, tt.testingResults, tt.installTestData)

			blockingErrors := 0
			warnings := 0
			for _, result := range results {
				if result.IsWarning {
					warnings++
				} else {
					blockingErrors++
				}
			}

			if blockingErrors != tt.wantBlockingErrors {
				t.Errorf("got %d blocking errors, want %d", blockingErrors, tt.wantBlockingErrors)
				for _, result := range results {
					if !result.IsWarning {
						t.Logf("  Blocking error: %v", result.Error)
					}
				}
			}
			if warnings != tt.wantWarnings {
				t.Errorf("got %d warnings, want %d", warnings, tt.wantWarnings)
				for _, result := range results {
					if result.IsWarning {
						t.Logf("  Warning: %v", result.Error)
					}
				}
			}
		})
	}
}

func Test_defaultQueriesIncludeCandidateTier(t *testing.T) {
	// When JobTiers is empty, QueriesFor should generate queries for all tiers
	// including candidate. This test is added to prevent regressions for candidate-tier
	// jobs being excluded, which is used by some TP jobs.
	queries := sippy.QueriesFor("vsphere", "amd64", "ha", "", "", "", "FeatureGate:TestGate]")

	tierNames := sets.New[string]()
	for _, q := range queries {
		tierNames.Insert(q.TierName)
	}

	expectedTiers := []string{"standard", "informing", "blocking", "candidate"}
	for _, tier := range expectedTiers {
		if !tierNames.Has(tier) {
			t.Errorf("default queries missing tier %q - got tiers: %v", tier, sets.List(tierNames))
		}
	}
}

func Test_allRequiredVariantsQueryCandidateTier(t *testing.T) {
	// Verify that all required variant definitions will query for the candidate
	// tier, either explicitly via JobTiers or via the default.
	allVariants := append(append([]JobVariant{}, requiredSelfManagedJobVariants...), requiredHypershiftJobVariants...)

	for _, variant := range allVariants {
		queries := sippy.QueriesFor(variant.Cloud, variant.Architecture, variant.Topology, variant.NetworkStack, variant.OS, variant.JobTiers, "FeatureGate:Test]")
		hasCandidateQuery := false
		for _, q := range queries {
			if q.TierName == "candidate" {
				hasCandidateQuery = true
				break
			}
		}
		if !hasCandidateQuery {
			t.Errorf("variant %+v does not query candidate tier - some platforms only run TechPreview tests in candidate-tier jobs", variant)
		}
	}
}

func Test_requiredHypershiftJobVariants_TopologyIsExternal(t *testing.T) {
	// Sippy classifies hypershift jobs with Topology:"external", not "hypershift".
	// Using the wrong value causes the verify-feature-promotion tool to find 0
	// tests when querying Sippy for hypershift variants.
	for i, variant := range requiredHypershiftJobVariants {
		if variant.Topology != "external" {
			t.Errorf("requiredHypershiftJobVariants[%d] has Topology=%q, want %q - Sippy classifies hypershift jobs as Topology:external",
				i, variant.Topology, "external")
		}
	}
}

func Test_requiredHypershiftJobVariants_SippyQueriesUseExternalTopology(t *testing.T) {
	for _, variant := range requiredHypershiftJobVariants {
		queries := sippy.QueriesFor(variant.Cloud, variant.Architecture, variant.Topology, variant.NetworkStack, variant.OS, variant.JobTiers, "FeatureGate:Test]")
		for _, query := range queries {
			hasExternalTopology := false
			for _, item := range query.Items {
				if item.ColumnField == "variants" && item.Value == "Topology:external" {
					hasExternalTopology = true
					break
				}
			}
			if !hasExternalTopology {
				t.Errorf("Sippy query for hypershift variant %+v (tier=%s) does not contain Topology:external filter",
					variant, query.TierName)
			}
		}
	}
}

func Test_topologyDisplayName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"external", "hypershift"},
		{"ha", "ha"},
		{"single", "single"},
		{"two-node-fencing", "two-node-fencing"},
		{"two-node-arbiter", "two-node-arbiter"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := topologyDisplayName(tt.input)
			if got != tt.want {
				t.Errorf("topologyDisplayName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func Test_nonHypershiftPlatforms(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"NutanixFeature", true},
		{"nutanixfeature", true},
		{"MetalFeature", true},
		{"VSphereFeature", true},
		{"OpenStackFeature", true},
		{"AzureFeature", true},
		{"GCPFeature", true},
		{"AWSFeature", false},
		{"GenericFeature", false},
		{"HypershiftFeature", false},
		{"ExternalFeature", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := nonHypershiftPlatforms.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("nonHypershiftPlatforms.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func Test_filterVariants(t *testing.T) {
	tests := []struct {
		name        string
		featureGate string
		variants    [][]JobVariant
		want        []JobVariant
	}{
		{
			name:        "AWS feature gate matches aws hypershift variant with external topology",
			featureGate: "AWSServiceLBNetworkSecurityGroup",
			variants: [][]JobVariant{
				requiredHypershiftJobVariants,
			},
			want: []JobVariant{
				{Cloud: "aws", Architecture: "amd64", Topology: "external"},
			},
		},
		{
			name:        "generic feature gate matches no hypershift variants",
			featureGate: "GenericFeature",
			variants: [][]JobVariant{
				requiredHypershiftJobVariants,
			},
			want: nil,
		},
		{
			name:        "VSphere feature gate matches vsphere self-managed variant",
			featureGate: "VSphereControlPlaneMachineSet",
			variants: [][]JobVariant{
				requiredSelfManagedJobVariants,
			},
			want: []JobVariant{
				{Cloud: "vsphere", Architecture: "amd64", Topology: "ha"},
			},
		},
		{
			name:        "Nutanix feature gate matches optional nutanix variant",
			featureGate: "NutanixGate",
			variants: [][]JobVariant{
				optionalSelfManagedPlatformVariants,
			},
			want: []JobVariant{
				{Cloud: "nutanix", Architecture: "amd64", Topology: "ha"},
			},
		},
		{
			name:        "Metal feature gate matches metal variants with network stacks",
			featureGate: "MetalFeature",
			variants: [][]JobVariant{
				requiredSelfManagedJobVariants,
			},
			want: []JobVariant{
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv4"},
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv6"},
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "dual"},
			},
		},
		{
			name:        "feature gate with multiple variant lists",
			featureGate: "MetalFeature",
			variants: [][]JobVariant{
				optionalSelfManagedPlatformVariants,
				requiredSelfManagedJobVariants,
			},
			want: []JobVariant{
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-arbiter", NetworkStack: "ipv4"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-arbiter", NetworkStack: "ipv6"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-arbiter", NetworkStack: "dual"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv4", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv6", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "dual", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv4"},
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv6"},
				{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "dual"},
			},
		},
		{
			name:        "DualReplica feature gate matches fencing topology variants",
			featureGate: "DualReplicaFeature",
			variants: [][]JobVariant{
				optionalSelfManagedPlatformVariants,
			},
			want: []JobVariant{
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv4", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv6", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "dual", JobTiers: "candidate,standard,informing,blocking"},
			},
		},
		{
			name:        "Fencing feature gate matches fencing topology variants",
			featureGate: "FencingFeature",
			variants: [][]JobVariant{
				optionalSelfManagedPlatformVariants,
			},
			want: []JobVariant{
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv4", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "ipv6", JobTiers: "candidate,standard,informing,blocking"},
				{Cloud: "metal", Architecture: "amd64", Topology: "two-node-fencing", NetworkStack: "dual", JobTiers: "candidate,standard,informing,blocking"},
			},
		},
		{
			name:        "amd64 in feature gate name matches amd64 variants only",
			featureGate: "Amd64SpecificFeature",
			variants: [][]JobVariant{
				{
					{Cloud: "aws", Architecture: "amd64", Topology: "ha"},
					{Cloud: "aws", Architecture: "arm64", Topology: "ha"},
				},
			},
			want: []JobVariant{
				{Cloud: "aws", Architecture: "amd64", Topology: "ha"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterVariants(tt.featureGate, tt.variants...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterVariants(%q) =\n  %+v\nwant:\n  %+v", tt.featureGate, got, tt.want)
			}
		})
	}
}

func Test_matchTwoNodeFeatureGates(t *testing.T) {
	tests := []struct {
		featureGate string
		topology    string
		want        bool
	}{
		{"dualreplicafeature", "two-node-fencing", true},
		{"fencingfeature", "two-node-fencing", true},
		{"dualreplicafeature", "ha", false},
		{"fencingfeature", "ha", false},
		{"genericfeature", "two-node-fencing", false},
		{"genericfeature", "ha", false},
		{"dualreplicafeature", "two-node-arbiter", false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s/%s", tt.featureGate, tt.topology)
		t.Run(name, func(t *testing.T) {
			got := matchTwoNodeFeatureGates(tt.featureGate, tt.topology)
			if got != tt.want {
				t.Errorf("matchTwoNodeFeatureGates(%q, %q) = %v, want %v",
					tt.featureGate, tt.topology, got, tt.want)
			}
		})
	}
}

func Test_testResultByName(t *testing.T) {
	results := []TestResults{
		{TestName: "test-alpha", TotalRuns: 10, SuccessfulRuns: 9},
		{TestName: "test-beta", TotalRuns: 20, SuccessfulRuns: 20},
		{TestName: "test-gamma", TotalRuns: 5, SuccessfulRuns: 3},
	}

	tests := []struct {
		name     string
		testName string
		wantNil  bool
		wantName string
	}{
		{"found first", "test-alpha", false, "test-alpha"},
		{"found middle", "test-beta", false, "test-beta"},
		{"found last", "test-gamma", false, "test-gamma"},
		{"not found", "test-delta", true, ""},
		{"empty name", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testResultByName(results, tt.testName)
			if tt.wantNil {
				if got != nil {
					t.Errorf("testResultByName(%q) = %+v, want nil", tt.testName, got)
				}
			} else {
				if got == nil {
					t.Fatalf("testResultByName(%q) = nil, want non-nil", tt.testName)
				}
				if got.TestName != tt.wantName {
					t.Errorf("testResultByName(%q).TestName = %q, want %q",
						tt.testName, got.TestName, tt.wantName)
				}
			}
		})
	}
}

func Test_OrderedJobVariants(t *testing.T) {
	input := []JobVariant{
		{Cloud: "gcp", Architecture: "amd64", Topology: "ha"},
		{Cloud: "aws", Architecture: "amd64", Topology: "ha"},
		{Cloud: "aws", Architecture: "amd64", Topology: "external"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "dual"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv4"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv6"},
		{Cloud: "aws", Architecture: "amd64", Topology: "single"},
	}

	want := []JobVariant{
		{Cloud: "aws", Architecture: "amd64", Topology: "external"},
		{Cloud: "aws", Architecture: "amd64", Topology: "ha"},
		{Cloud: "gcp", Architecture: "amd64", Topology: "ha"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv4"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "ipv6"},
		{Cloud: "metal", Architecture: "amd64", Topology: "ha", NetworkStack: "dual"},
		{Cloud: "aws", Architecture: "amd64", Topology: "single"},
	}

	sorted := make([]JobVariant, len(input))
	copy(sorted, input)
	sort.Sort(OrderedJobVariants(sorted))

	if !reflect.DeepEqual(sorted, want) {
		t.Errorf("OrderedJobVariants sort:\ngot:  %+v\nwant: %+v", sorted, want)
	}
}

func Test_validateJobTiers_comprehensive(t *testing.T) {
	tests := []struct {
		name    string
		variant JobVariant
		wantErr bool
	}{
		{
			name:    "empty job tiers is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha"},
			wantErr: false,
		},
		{
			name:    "standard is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "standard"},
			wantErr: false,
		},
		{
			name:    "informing is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "informing"},
			wantErr: false,
		},
		{
			name:    "blocking is valid",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "blocking"},
			wantErr: false,
		},
		{
			name:    "all valid tiers combined",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "standard,informing,blocking,candidate"},
			wantErr: false,
		},
		{
			name:    "invalid tier rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "invalid"},
			wantErr: true,
		},
		{
			name:    "valid tier with invalid tier rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: "standard,invalid"},
			wantErr: true,
		},
		{
			name:    "only commas rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: ",,,"},
			wantErr: true,
		},
		{
			name:    "whitespace only tiers rejected",
			variant: JobVariant{Cloud: "aws", Architecture: "amd64", Topology: "ha", JobTiers: " , , "},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJobTiers(tt.variant)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJobTiers(%+v) error = %v, wantErr %v", tt.variant, err, tt.wantErr)
			}
		})
	}
}

func Test_allDefinedVariantsHaveValidJobTiers(t *testing.T) {
	allVariants := []struct {
		name     string
		variants []JobVariant
	}{
		{"requiredSelfManagedJobVariants", requiredSelfManagedJobVariants},
		{"optionalSelfManagedPlatformVariants", optionalSelfManagedPlatformVariants},
		{"requiredHypershiftJobVariants", requiredHypershiftJobVariants},
	}

	for _, group := range allVariants {
		for i, variant := range group.variants {
			name := fmt.Sprintf("%s[%d]-%s-%s-%s", group.name, i, variant.Cloud, variant.Architecture, variant.Topology)
			t.Run(name, func(t *testing.T) {
				if err := validateJobTiers(variant); err != nil {
					t.Errorf("variant %+v has invalid JobTiers: %v", variant, err)
				}
			})
		}
	}
}

func Test_buildHTMLFeatureGateData(t *testing.T) {
	tests := []struct {
		name           string
		featureGate    string
		testingResults map[JobVariant]*TestingResults
		blockingErrors []error
		release        string
		wantSufficient bool
		wantVariants   int
		wantTests      int
	}{
		{
			name:        "sufficient testing with external topology hypershift variant",
			featureGate: "AWSFeature",
			testingResults: map[JobVariant]*TestingResults{
				{Cloud: "aws", Architecture: "amd64", Topology: "external"}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
						{TestName: "test2", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			blockingErrors: nil,
			release:        "4.18",
			wantSufficient: true,
			wantVariants:   1,
			wantTests:      2,
		},
		{
			name:        "multiple variants including external topology",
			featureGate: "AWSFeature",
			testingResults: map[JobVariant]*TestingResults{
				{Cloud: "aws", Architecture: "amd64", Topology: "ha"}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
				{Cloud: "aws", Architecture: "amd64", Topology: "external"}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 15, SuccessfulRuns: 15},
					},
				},
			},
			blockingErrors: nil,
			release:        "4.18",
			wantSufficient: true,
			wantVariants:   2,
			wantTests:      1,
		},
		{
			name:        "with blocking errors",
			featureGate: "AWSFeature",
			testingResults: map[JobVariant]*TestingResults{
				{Cloud: "aws", Architecture: "amd64", Topology: "external"}: {
					TestResults: []TestResults{
						{TestName: "test1", TotalRuns: 5, SuccessfulRuns: 5},
					},
				},
			},
			blockingErrors: []error{fmt.Errorf("insufficient tests")},
			release:        "4.18",
			wantSufficient: false,
			wantVariants:   1,
			wantTests:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildHTMLFeatureGateData(tt.featureGate, tt.testingResults, tt.blockingErrors, tt.release)

			if result.Name != tt.featureGate {
				t.Errorf("Name = %q, want %q", result.Name, tt.featureGate)
			}
			if result.Sufficient != tt.wantSufficient {
				t.Errorf("Sufficient = %v, want %v", result.Sufficient, tt.wantSufficient)
			}
			if len(result.Variants) != tt.wantVariants {
				t.Errorf("got %d variants, want %d", len(result.Variants), tt.wantVariants)
			}
			if len(result.Tests) != tt.wantTests {
				t.Errorf("got %d tests, want %d", len(result.Tests), tt.wantTests)
			}

			for _, v := range result.Variants {
				if v.Topology == "hypershift" {
					for _, test := range result.Tests {
						cell := test.Cells[v.ColIndex-1]
						if !strings.Contains(cell.SippyURL, "Topology%3Aexternal") && !strings.Contains(cell.SippyURL, "Topology:external") {
							t.Errorf("Sippy URL for hypershift variant should query Topology:external but got: %s", cell.SippyURL)
						}
					}
				}
				if v.Topology == "external" {
					t.Errorf("HTML variant column should display %q not %q for hypershift variants", "hypershift", "external")
				}
			}
		})
	}
}
