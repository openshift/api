package tests

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/ghodss/yaml"
)

func TestLoadSuiteAndDefaulting(t *testing.T) {
    // Load the test suite files from testdata
    // Read the test suite YAML directly and unmarshal into SuiteSpec. We avoid
    // calling LoadTestSuiteSpecs because it requires CRD files to be present.
    rawPath := filepath.Join("testdata", "minimal.testsuite.yaml")
    raw, err := os.ReadFile(rawPath)
    if err != nil {
        t.Fatalf("could not read testdata file: %v", err)
    }

    var s SuiteSpec
    if err := yaml.Unmarshal(raw, &s); err != nil {
        t.Fatalf("could not unmarshal testsuite YAML: %v", err)
    }

    if len(s.Tests.OnCreate) == 0 {
        t.Fatalf("expected onCreate tests in suite")
    }

    // Verify defaulting for onCreate: expected should default to initial when empty
    onCreate := s.Tests.OnCreate[0]
    expectedBytes := defaultOnCreateExpected(onCreate.Initial, onCreate.Expected)
    if string(expectedBytes) == "" {
        t.Fatalf("expected defaulted expected bytes to be non-empty")
    }

    // Verify defaulting for onUpdate
    if len(s.Tests.OnUpdate) == 0 {
        t.Fatalf("expected onUpdate tests in suite")
    }
    onUpdate := s.Tests.OnUpdate[0]
    expectedBytes2 := defaultOnUpdateExpected(onUpdate.Updated, onUpdate.Expected)
    if string(expectedBytes2) == "" {
        t.Fatalf("expected defaulted expected bytes for onUpdate to be non-empty")
    }
}
