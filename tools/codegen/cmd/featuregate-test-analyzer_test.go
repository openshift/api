package main

import (
	"encoding/json"
	"reflect"
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

			got, err := listTestResultFor(tt.args.featureGate, sets.New[string](tt.args.clusterProfile))
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
			results := checkIfTestingIsSufficient(tt.featureGate, tt.testingResults)

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
			results := checkIfTestingIsSufficient(tt.featureGate, tt.testingResults)

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
