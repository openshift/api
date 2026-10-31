package tests

import (
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// No release feature set enables ClusterAPIOperator without the checker gate.
// Exercise that combination against its generated schema using the existing
// fixture runner, rather than selecting a merged release schema.
var _ = Describe("ClusterVersion operator-only capability schema", func() {
	raw, err := os.ReadFile(filepath.Join("testdata", "clusterversion-capabilities", "ClusterAPIOperator.yaml"))
	Expect(err).NotTo(HaveOccurred())

	var suite SuiteSpec
	Expect(yaml.Unmarshal(raw, &suite)).To(Succeed())
	suite.Version = "v1"
	suite.PerTestRuntimeInfo = &PerTestRuntimeInfo{
		CRDFilenames: []string{filepath.Join("..", "config", "v1", "zz_generated.featuregated-crd-manifests", "clusterversions.config.openshift.io", "ClusterAPIOperator.yaml")},
	}
	GenerateTestSuite(suite)
})
