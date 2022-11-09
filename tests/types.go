package tests

// SuiteSpec defines a test suite specification.
type SuiteSpec struct {
	// Name is the name of the test suite.
	Name string `json:"name"`

	// CRD is a CRD file path that should be installed as a part of this test.
	CRD string `json:"crd"`

	// Version is the version of the CRD under test in this file.
	// When omitted, if there is a single version in the CRD, this is assumed to be the correct version.
	// If there are multiple versions within the CRD, an educated guess is made based on the directory structure.
	Version string `json:"version,omitempty"`

	// Tests defines the test cases to run for this test suite.
	Tests TestSpec `json:"tests"`
}

// TestSpec defines the test specs for individual tests in this suite.
type TestSpec struct {
	// OnCreate defines a list of on create style tests.
	OnCreate []OnCreateTestSpec `json:"onCreate"`

	// OnUpdate defines a list of on create style tests.
	OnUpdate []OnUpdateTestSpec `json:"onUpdate"`
}

// OnCreateTestSpec defines an individual test case for the on create style tests.
type OnCreateTestSpec struct {
	// Name is the name of this test case.
	Name string `json:"name"`

	// Initial is a literal string containing the initial YAML content from which to
	// create the resource.
	// Note `apiVersion` and `kind` fields are required though `metadata` can be omitted.
	// Typically this will vary in `spec` only test to test.
	Initial string `json:"initial"`

	// ExpectedError defines the error string that should be returned when the initial resourec is invalid.
	// This will be matched as a substring of the actual error when non-empty.
	ExpectedError string `json:"expectedError"`

	// Expected is a literal string containing the expected YAML content that should be
	// persisted when the resource is created.
	// Note `apiVersion` and `kind` fields are required though `metadata` can be omitted.
	// Typically this will vary in `spec` only test to test.
	Expected string `json:"expected"`
}

// OnUpdateTestSpec defines an individual test case for the on update style tests.
type OnUpdateTestSpec struct {
	// Name is the name of this test case.
	Name string `json:"name"`

	// Initial is a literal string containing the initial YAML content from which to
	// create the resource.
	// Note `apiVersion` and `kind` fields are required though `metadata` can be omitted.
	// Typically this will vary in `spec` only test to test.
	Initial string `json:"initial"`

	// Updated is a literal string containing the updated YAML content from which to
	// update the resource.
	// Note `apiVersion` and `kind` fields are required though `metadata` can be omitted.
	// Typically this will vary in `spec` only test to test.
	Updated string `json:"updated"`

	// ExpectedError defines the error string that should be returned when the initial resourec is invalid.
	// This will be matched as a substring of the actual error when non-empty.
	ExpectedError string `json:"expectedError"`

	// Expected is a literal string containing the expected YAML content that should be
	// persisted when the resource is updated.
	// Note `apiVersion` and `kind` fields are required though `metadata` can be omitted.
	// Typically this will vary in `spec` only test to test.
	Expected string `json:"expected"`
}
