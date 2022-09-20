package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ghodss/yaml"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/komega"
)

// LoadTestSuiteSpecs recursively walks the given paths looking for any file with the suffix `.testsuite.yaml`.
// It then loads these files in SuiteSpec structs ready for the generator to generate the test cases.
func LoadTestSuiteSpecs(paths ...string) ([]SuiteSpec, error) {
	suiteFiles := make(map[string]struct{})

	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".testsuite.yaml") {
				suiteFiles[path] = struct{}{}
			}

			return nil
		}); err != nil {
			return nil, fmt.Errorf("could not load files from path %q: %w", path, err)
		}
	}

	out := []SuiteSpec{}
	for path := range suiteFiles {
		suite, err := loadSuiteFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not set up test suite: %w", err)
		}

		out = append(out, suite)
	}

	return out, nil
}

// loadSuiteFile loads an individual SuiteSpec from the given file name.
func loadSuiteFile(path string) (SuiteSpec, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return SuiteSpec{}, fmt.Errorf("could not read file %q: %w", path, err)
	}

	s := SuiteSpec{}
	if err := yaml.Unmarshal(raw, &s); err != nil {
		return SuiteSpec{}, fmt.Errorf("could not unmarshal YAML file %q: %w", path, err)
	}

	if s.CRD == "" {
		return SuiteSpec{}, fmt.Errorf("test suite spec %q is invalid: missing required field `crd`", path)
	}

	// If the CRD path isn't absolute, generate the absolute path from the testsuite file location.
	if err := setAbsolutePath(path, &s.CRD); err != nil {
		return SuiteSpec{}, fmt.Errorf("could not set absolute path for CRD: %w", err)
	}

	return s, nil
}

// setAbsolutePath overwrites the given path with the absolute path if it isn't already
// an absolute path.
// The path is expected to be relative to the suite file.
func setAbsolutePath(suitePath string, path *string) error {
	if path == nil {
		return nil
	}

	if filepath.IsAbs(*path) {
		return nil
	}
	dir := filepath.Dir(suitePath)
	relPath := filepath.Join(dir, *path)

	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return fmt.Errorf("could not generate absolute path for %q: %w", relPath, err)
	}

	*path = absPath

	return nil
}

// GenerateTestSuite generates a Ginkgo test suite from the provided SuiteSpec.
func GenerateTestSuite(suiteSpec SuiteSpec) {
	suiteName, err := generateSuiteName(suiteSpec)
	Expect(err).ToNot(HaveOccurred())

	Describe(suiteName, func() {
		var crdOptions envtest.CRDInstallOptions
		var crd *apiextensionsv1.CustomResourceDefinition

		BeforeEach(func() {
			Expect(k8sClient).ToNot(BeNil(), "Kuberentes client is not initialised")

			crdOptions = envtest.CRDInstallOptions{
				Paths: []string{
					suiteSpec.CRD,
				},
				ErrorIfPathMissing: true,
			}

			crds, err := envtest.InstallCRDs(cfg, crdOptions)
			Expect(err).ToNot(HaveOccurred())

			Expect(crds).To(HaveLen(1), "Only one CRD should have been installed")
			crd = crds[0]

			Expect(envtest.WaitForCRDs(cfg, crds, crdOptions)).To(Succeed())
		})

		AfterEach(func() {
			// Remove all of the resources we created during the test.
			for _, u := range newUnstructuredsFor(crd) {
				Expect(k8sClient.DeleteAllOf(ctx, u, client.InNamespace("default")))
			}

			// Remove the CRD and wait for it to be removed from the API.
			// If we don't wait then subsequent tests may fail.
			Expect(envtest.UninstallCRDs(cfg, crdOptions)).ToNot(HaveOccurred())
			Eventually(komega.Get(crd)).Should(Not(Succeed()))
		})

		generateOnCreateTable(suiteSpec.Tests.OnCreate)
	})
}

// generateOnCreateTable generates a table of tests from the defined OnCreate tests
// within the test suite test spec.
func generateOnCreateTable(onCreateTests []OnCreateTestSpec) {
	type onCreateTableInput struct {
		initial       []byte
		expected      []byte
		expectedError string
	}

	// assertOnCreate runs the actual test for each table entry
	var assertOnCreate interface{} = func(in onCreateTableInput) {
		initialObj, err := newUnstructuredFrom(in.initial)
		Expect(err).ToNot(HaveOccurred(), "initial data should be a valid Kubernetes YAML resource")

		err = k8sClient.Create(ctx, initialObj)
		if in.expectedError != "" {
			Expect(err).To(MatchError(ContainSubstring(in.expectedError)))
			return
		}
		Expect(err).ToNot(HaveOccurred())

		// Fetch the object we just created from the API.
		gotObj := newEmptyUnstructuredFrom(initialObj)
		Expect(k8sClient.Get(ctx, objectKey(initialObj), gotObj))

		expectedObj, err := newUnstructuredFrom(in.expected)
		Expect(err).ToNot(HaveOccurred(), "expected data should be a valid Kubernetes YAML resource when no expected error is provided")

		// Ensure the name and namespace match.
		// The IgnoreAutogeneratedMetadata will ignore any additional meta set in the API.
		expectedObj.SetName(gotObj.GetName())
		expectedObj.SetNamespace(gotObj.GetNamespace())

		Expect(gotObj).To(komega.EqualObject(expectedObj, komega.IgnoreAutogeneratedMetadata))
	}

	// First argument to the table is the test function.
	tableEntries := []interface{}{assertOnCreate}

	// Convert the test specs into table entries
	for _, testEntry := range onCreateTests {
		tableEntries = append(tableEntries, Entry(testEntry.Name, onCreateTableInput{
			initial:       []byte(testEntry.Initial),
			expected:      []byte(testEntry.Expected),
			expectedError: testEntry.ExpectedError,
		}))
	}

	DescribeTable("On Create", tableEntries...)
}

// newUnstructuredsFor creates a set of unstructured resources for each version of the CRD.
// This allows us to ensure all CR instances are deleted after each test.
func newUnstructuredsFor(crd *apiextensionsv1.CustomResourceDefinition) []*unstructured.Unstructured {
	out := []*unstructured.Unstructured{}

	for _, version := range crd.Spec.Versions {
		out = append(out, newUnstructuredsForVersion(crd, version.Name))
	}

	return out
}

// newUnstructuredsForVersion creates an unstructured resource for the CRD at a given version.
func newUnstructuredsForVersion(crd *apiextensionsv1.CustomResourceDefinition, version string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}

	u.SetAPIVersion(fmt.Sprintf("%s/%s", crd.Spec.Group, version))
	u.SetKind(crd.Spec.Names.Kind)

	return u
}

// newUnstructuredFrom unmarshals the raw YAML data into an unstructured,
// and then sets the namespace and generateName ahead of the test.
func newUnstructuredFrom(raw []byte) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}

	if err := k8syaml.Unmarshal(raw, &u.Object); err != nil {
		return nil, fmt.Errorf("could not unmarshal raw YAML: %w", err)
	}

	// Names should be unique for each test so ensure we generate a name
	u.SetGenerateName("test-")
	// We need to have a namespace, use the default.
	u.SetNamespace("default")

	return u, nil
}

// newEmptyUnstructuredFrom creates a new unstructured with the same GVK as the input object,
// all other fields are cleared.
func newEmptyUnstructuredFrom(initial *unstructured.Unstructured) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}

	if initial != nil {
		u.GetObjectKind().SetGroupVersionKind(initial.GetObjectKind().GroupVersionKind())
	}

	return u
}

// objectKey extracts a client.ObjectKey from the given object.
func objectKey(obj client.Object) client.ObjectKey {
	return client.ObjectKey{Namespace: obj.GetNamespace(), Name: obj.GetName()}
}

// loadCRD loads the CustomResourceDefinition defined in the suite spec.
func loadCRD(suiteSpec SuiteSpec) (*apiextensionsv1.CustomResourceDefinition, error) {
	raw, err := ioutil.ReadFile(suiteSpec.CRD)
	if err != nil {
		return nil, fmt.Errorf("could not load CRD: %w", err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := yaml.Unmarshal(raw, crd); err != nil {
		return nil, fmt.Errorf("could not unmarshal CRD: %w", err)
	}

	return crd, nil
}

// generateSuiteName prepends the specified suite name with the GVR string
// for the CRD under test.
func generateSuiteName(suiteSpec SuiteSpec) (string, error) {
	crd, err := loadCRD(suiteSpec)
	if err != nil {
		return "", fmt.Errorf("could not load CRD: %w", err)
	}

	gvr := schema.GroupVersionResource{
		Group: crd.Spec.Group,
		Resource: crd.Spec.Names.Plural,
	}

	if len(crd.Spec.Versions) == 1 {
		// When there's only one version it's easy to know which we are testing.
		gvr.Version = crd.Spec.Versions[0].Name
	} else {
		// Otherwise we need to guess the version we are testing, it's probably the package/folder name.
		packageDir := filepath.Dir(suiteSpec.CRD)
		gvr.Version = filepath.Base(packageDir)
	}

	return fmt.Sprintf("[%s] %s", gvr.String(), suiteSpec.Name), nil
}
