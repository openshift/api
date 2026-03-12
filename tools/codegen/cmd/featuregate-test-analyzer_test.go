package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
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
			//t.Skip("this is for ease of manual testing")

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
