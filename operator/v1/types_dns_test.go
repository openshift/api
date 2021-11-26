package v1

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/grpclog"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientschema "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apiextensions-apiserver/test/integration/fixtures"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
)

var etcdURL = ""

const installEtcd = `
Cannot find etcd, cannot run integration tests
Please see https://git.k8s.io/community/contributors/devel/sig-testing/integration-tests.md#install-etcd-dependency for instructions.
You can use 'hack/install-etcd.sh' to install a copy in third_party/.
`

func getAvailablePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("could not bind to a port: %v", err)
	}
	// It is possible but unlikely that someone else will bind this port before we
	// get a chance to use it.
	defer func() { _ = l.Close() }()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// getEtcdPath returns a path to an etcd executable.
func getEtcdPath() (string, error) {
	//return exec.LookPath("etcd")
	return "/tmp/etcd-download-test/etcd", nil
}

// RunCustomEtcd starts a custom etcd instance for test purposes.
func RunCustomEtcd(dataDir string, customFlags []string) (url string, stopFn func(), err error) {
	// TODO: Check for valid etcd version.
	etcdPath, err := getEtcdPath()
	if err != nil {
		fmt.Fprint(os.Stderr, installEtcd)
		return "", nil, fmt.Errorf("could not find etcd in PATH: %v", err)
	}
	etcdPort, err := getAvailablePort()
	if err != nil {
		return "", nil, fmt.Errorf("could not get a port: %v", err)
	}
	customURL := fmt.Sprintf("http://127.0.0.1:%d", etcdPort)

	klog.Infof("starting etcd on %s", customURL)

	etcdDataDir, err := ioutil.TempDir(os.TempDir(), dataDir)
	if err != nil {
		return "", nil, fmt.Errorf("unable to make temp etcd data dir %s: %v", dataDir, err)
	}
	klog.Infof("storing etcd data in: %v", etcdDataDir)

	ctx, cancel := context.WithCancel(context.Background())
	args := []string{
		"--data-dir",
		etcdDataDir,
		"--listen-client-urls",
		customURL,
		"--advertise-client-urls",
		customURL,
		"--listen-peer-urls",
		"http://127.0.0.1:0",
		"-log-level",
		"warn", // set to info or debug for more logs
	}
	args = append(args, customFlags...)
	cmd := exec.CommandContext(ctx, etcdPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stop := func() {
		cancel()
		err := cmd.Wait()
		klog.Infof("etcd exit status: %v", err)
		err = os.RemoveAll(etcdDataDir)
		if err != nil {
			klog.Warningf("error during etcd cleanup: %v", err)
		}
	}

	// Quiet etcd logs for integration tests
	// Comment out to get verbose logs if desired
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stderr))

	if err := cmd.Start(); err != nil {
		return "", nil, fmt.Errorf("failed to run etcd: %v", err)
	}

	var i int32 = 1
	const pollCount = int32(300)

	for i <= pollCount {
		conn, err := net.DialTimeout("tcp", strings.TrimPrefix(customURL, "http://"), 1*time.Second)
		if err == nil {
			conn.Close()
			break
		}

		if i == pollCount {
			stop()
			return "", nil, fmt.Errorf("could not start etcd")
		}

		time.Sleep(100 * time.Millisecond)
		i = i + 1
	}

	return customURL, stop, nil
}

func GetEnvAsStringOrFallback(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func startEtcd() (func(), error) {
	//if runtime.GOARCH == "arm64" {
	//	os.Setenv("ETCD_UNSUPPORTED_ARCH", "arm64")
	//}

	etcdURL = GetEnvAsStringOrFallback("KUBE_INTEGRATION_ETCD_URL", "http://127.0.0.1:2379")
	conn, err := net.Dial("tcp", strings.TrimPrefix(etcdURL, "http://"))
	if err == nil {
		klog.Infof("etcd already running at %s", etcdURL)
		_ = conn.Close()
		return func() {}, nil
	}
	klog.V(1).Infof("could not connect to etcd: %v", err)

	currentURL, stop, err := RunCustomEtcd("integration_test_etcd_data", nil)
	if err != nil {
		return nil, err
	}

	etcdURL = currentURL
	os.Setenv("KUBE_INTEGRATION_ETCD_URL", etcdURL)

	return stop, nil
}

func TestCustomResourceItemsValidation(t *testing.T) {
	// start etcd
	stop, err := startEtcd()
	defer stop()
	if err != nil {
		t.Fatal(err)
	}
	// Read CRD from file
	crdInBytes, err := os.ReadFile("0000_70_dns-operator_00.crd.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Start a client
	tearDown, apiExtensionClient, client, err := fixtures.StartDefaultServerWithClients(t)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	// decode CRD manifest
	obj, _, err := clientschema.Codecs.UniversalDeserializer().Decode([]byte(crdInBytes), nil, nil)
	if err != nil {
		t.Fatalf("failed decoding of: %v\n\n%s", err, crdInBytes)
	}
	crd := obj.(*apiextensionsv1.CustomResourceDefinition)

	// create CRDs
	t.Logf("Creating CRD %s", crd.Name)
	if _, err = fixtures.CreateNewV1CustomResourceDefinition(crd, apiExtensionClient, client); err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	// create CR
	gvr, testCases := schema.GroupVersionResource{
		Group:    crd.Spec.Group,
		Version:  crd.Spec.Versions[0].Name,
		Resource: crd.Spec.Names.Plural,
	}, []struct {
		name                 string
		dns                  *DNS
		expectedErrorMessage string
		expectedType         UpstreamType
		expectedAddress      string
		expectedPort         uint32
	}{
		{
			name: "Dns spec without upstreamResolvers passes",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{},
			},
			expectedErrorMessage: "none",
			expectedType:         SystemResolveConfType,
		},
		{
			name: "Dns spec with upstream typed System passes",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type: SystemResolveConfType,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "none",
			expectedType:         SystemResolveConfType,
		},
		{
			name: "Dns spec with upstream typed System with Address fails",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type:    SystemResolveConfType,
								Address: "1.2.3.6",
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "\"spec.upstreamresolvers.upstreams\" must validate at least one schema (anyOf)",
			expectedType:         SystemResolveConfType,
		},
		{
			name: "Dns spec with upstream typed System with Port fails",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type: SystemResolveConfType,
								Port: 53,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "\"spec.upstreamresolvers.upstreams\" must validate at least one schema (anyOf)",
			expectedType:         SystemResolveConfType,
		},
		{
			name: "Dns spec with type upstream Network without Address fails",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type: NetworkResolverType,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "Unsupported value: \"Network\": supported values: \"\", \"SystemResolvConf\"", //Is it possible to modify crd to have a clearer error message?
			expectedType:         NetworkResolverType,
		},
		{
			name: "Dns spec with network upstream passes",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type:    NetworkResolverType,
								Address: "1.2.3.4",
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "none",
			expectedPort:         53,
			expectedAddress:      "1.2.3.4",
			expectedType:         NetworkResolverType,
		},
		{
			name: "Dns spec with network upstream passes",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type:    NetworkResolverType,
								Address: "1.2.3.4",
								Port:    5354,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "none",
			expectedPort:         5354,
			expectedAddress:      "1.2.3.4",
			expectedType:         NetworkResolverType,
		},
		{
			name: "Dns spec with network upstream with wrong Address fails",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type:    NetworkResolverType,
								Address: "this is no address",
								Port:    5354,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "address in body must be of type ipv4",
		},
		{
			name: "Dns spec with network upstream with wrong Address fails",
			dns: &DNS{
				TypeMeta: metav1.TypeMeta{
					Kind:       "DNS",
					APIVersion: "operator.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: DNSSpec{
					UpstreamResolvers: UpstreamResolvers{
						Upstreams: []Upstream{
							{
								Type:    NetworkResolverType,
								Address: "1.23.4.44",
								Port:    99999,
							},
						},
						Policy: RoundRobinForwardingPolicy,
					},
				},
			},
			expectedErrorMessage: "port in body should be less than or equal to 65535",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			z, err := runtime.DefaultUnstructuredConverter.ToUnstructured(tc.dns)
			if err != nil {
				t.Errorf("toUnstructured unexpected error: %v", err)
			}
			u := unstructured.Unstructured{Object: z}
			_, err = client.Resource(gvr).Create(context.TODO(), &u, metav1.CreateOptions{})

			if err != nil {
				assert.Containsf(t, err.Error(), tc.expectedErrorMessage, "expected error containing %q, got %s", tc.expectedErrorMessage, err)
			} else {
				v, err := client.Resource(gvr).Get(context.TODO(), "default", metav1.GetOptions{})
				if err != nil {
					t.Error(err.Error())
				}
				var savedDNS DNS
				err = runtime.DefaultUnstructuredConverter.FromUnstructured(v.Object, &savedDNS)
				if err != nil {
					t.Error(err.Error())
				}
				if tc.expectedPort > 0 {
					assert.Equal(t, savedDNS.Spec.UpstreamResolvers.Upstreams[0].Port, tc.expectedPort)
				}
				if tc.expectedAddress != "" {
					assert.Equal(t, savedDNS.Spec.UpstreamResolvers.Upstreams[0].Address, tc.expectedAddress)
				}
				if tc.expectedType != "" {
					assert.Equal(t, savedDNS.Spec.UpstreamResolvers.Upstreams[0].Type, tc.expectedType)
					if tc.expectedType == SystemResolveConfType {
						assert.Emptyf(t, savedDNS.Spec.UpstreamResolvers.Upstreams[0].Address, "Address should be empty")
					}
				}
				err = client.Resource(gvr).Delete(context.TODO(), "default", metav1.DeleteOptions{})
				if err != nil {
					t.Error(err.Error())
				}
			}

		})
	}

}
