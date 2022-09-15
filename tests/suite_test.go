package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testScheme *runtime.Scheme
var ctx = context.Background()

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "API Integration Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	testScheme = scheme.Scheme

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// CEL requires Kube 1.25 and above, so check for the minimum server version.
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())

	serverVersion, err := discoveryClient.ServerVersion()
	Expect(err).ToNot(HaveOccurred())

	Expect(serverVersion.Major).To(Equal("1"))

	minorInt, err := strconv.Atoi(serverVersion.Minor)
	Expect(err).ToNot(HaveOccurred())
	Expect(minorInt).To(BeNumerically(">=", 25), fmt.Sprintf("This test suite requires a Kube API server of at least version 1.25, current version is 1.%s", serverVersion.Minor))
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
