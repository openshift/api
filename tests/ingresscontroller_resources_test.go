package tests

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"
)

func ingressResourceDefinition(featureSet string) *apiextensionsv1.CustomResourceDefinition {
	data, err := os.ReadFile(fmt.Sprintf("../operator/v1/zz_generated.crd-manifests/0000_50_ingress_00_ingresscontrollers-%s.crd.yaml", featureSet))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	definition := &apiextensionsv1.CustomResourceDefinition{}
	ExpectWithOffset(1, yaml.Unmarshal(data, definition)).To(Succeed())
	return definition
}

func ingressResourceObject(resources interface{}) *unstructured.Unstructured {
	return (&unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "operator.openshift.io/v1", "kind": "IngressController",
		"metadata": map[string]interface{}{"generateName": "resources-", "namespace": "ingress-resources"},
		"spec":     map[string]interface{}{"resources": resources},
	}}).DeepCopy()
}

func ingressResourceRoundTrip(object *unstructured.Unstructured) interface{} {
	stored := &unstructured.Unstructured{}
	stored.SetGroupVersionKind(object.GroupVersionKind())
	ExpectWithOffset(1, k8sClient.Get(ctx, client.ObjectKeyFromObject(object), stored)).To(Succeed())
	resources, _, err := unstructured.NestedFieldNoCopy(stored.Object, "spec", "resources")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return resources
}

var _ = Describe("NE2926 IngressController resource configuration", Ordered, func() {
	BeforeAll(func() {
		discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		serverVersion, err := discoveryClient.ServerVersion()
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "NE2926 actual API server: %#v\n", *serverVersion)
		definition := ingressResourceDefinition("TechPreviewNoUpgrade")
		resources := definition.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["resources"]
		for _, target := range []string{"router", "haproxy"} {
			component, found := resources.Properties[target]
			Expect(found).To(BeTrue())
			Expect(component.Properties).To(HaveKey("requests"))
			Expect(component.Properties).To(HaveKey("limits"))
			Expect(component.Properties).NotTo(HaveKey("claims"))
			Expect(component.XValidations).To(HaveLen(3))
			Expect(component.XValidations[0].Rule).To(ContainSubstring("self.requests.size() <= 8"))
			Expect(component.XValidations[1].Rule).To(ContainSubstring("[has(self.requests) ? self.requests : {}].all(m, m.all(k,"))
			Expect(component.XValidations[2].Rule).To(ContainSubstring("[has(self.limits) ? self.limits : {}].all(m, m.all(k,"))
		}
		_, err = envtest.InstallCRDs(cfg, envtest.CRDInstallOptions{CRDs: []*apiextensionsv1.CustomResourceDefinition{definition}})
		Expect(err).NotTo(HaveOccurred(), "complete generated schema must register within static CEL cost limits")
		fmt.Fprintln(GinkgoWriter, "NE2926 complete generated preview schema registered within static cost limits")
		Expect(k8sClient.Create(ctx, &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]interface{}{"name": "ingress-resources"},
		}})).To(Succeed())
	})

	AfterAll(func() {
		definition := ingressResourceDefinition("TechPreviewNoUpgrade")
		Expect(envtest.UninstallCRDs(cfg, envtest.CRDInstallOptions{CRDs: []*apiextensionsv1.CustomResourceDefinition{definition}})).To(Succeed())
		Eventually(func() bool {
			return apierrors.IsNotFound(k8sClient.Get(ctx, client.ObjectKeyFromObject(definition), &apiextensionsv1.CustomResourceDefinition{}))
		}, timeout).Should(BeTrue())
	})

	for _, target := range []string{"router", "haproxy"} {
		DescribeTable("stores explicit resources for "+target, func(component map[string]interface{}) {
			resources := map[string]interface{}{target: component}
			object := ingressResourceObject(resources)
			Expect(k8sClient.Create(ctx, object, client.FieldValidation(metav1.FieldValidationStrict))).To(Succeed())
			Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
		},
			Entry("native requests", map[string]interface{}{"requests": map[string]interface{}{"cpu": "500m", "memory": "1Gi", "ephemeral-storage": "1e9"}}),
			Entry("limit only without request defaulting", map[string]interface{}{"limits": map[string]interface{}{"qat.intel.com/cy": "1"}}),
			Entry("integer JSON quantity", map[string]interface{}{"requests": map[string]interface{}{"cpu": int64(1)}, "limits": map[string]interface{}{"cpu": "1000m"}}),
			Entry("zero and tiny quantities", map[string]interface{}{"limits": map[string]interface{}{"cpu": "0", "memory": "0.0000000001"}}),
			Entry("64 character quantity", map[string]interface{}{"limits": map[string]interface{}{"cpu": strings.Repeat("0", 63) + "1"}}),
			Entry("empty maps", map[string]interface{}{"requests": map[string]interface{}{}, "limits": map[string]interface{}{}}),
			Entry("unsupported native name is not validated here", map[string]interface{}{"limits": map[string]interface{}{"storage": "1Gi"}}),
			Entry("reserved namespace is not validated here", map[string]interface{}{"limits": map[string]interface{}{"kubernetes.io/device": "1"}}),
			Entry("quota prefix is not validated here", map[string]interface{}{"limits": map[string]interface{}{"requests.example.com/device": "1"}}),
			Entry("qualified name syntax is not validated here", map[string]interface{}{"limits": map[string]interface{}{"example.com/bad+name": "1"}}),
			Entry("key length is not validated here", map[string]interface{}{"limits": map[string]interface{}{strings.Repeat("a", 318): "1"}}),
			Entry("hugepage size is not validated here", map[string]interface{}{"limits": map[string]interface{}{"hugepages-0": "1", "hugepages-invalid": "2"}}),
			Entry("request-only hugepages", map[string]interface{}{"requests": map[string]interface{}{"hugepages-2Mi": "2Mi"}}),
			Entry("misaligned hugepages", map[string]interface{}{"requests": map[string]interface{}{"hugepages-2Mi": "3Mi"}, "limits": map[string]interface{}{"hugepages-2Mi": "3145728"}}),
			Entry("decimal binary exponent spellings", map[string]interface{}{"requests": map[string]interface{}{"memory": "1Gi", "example.com/device": "1e0"}, "limits": map[string]interface{}{"memory": "1073741824", "example.com/device": "1"}}),
		)

		DescribeTable("rejects invalid resource relationships for "+target, func(component map[string]interface{}) {
			resources := map[string]interface{}{target: component}
			object := ingressResourceObject(resources)
			Expect(apierrors.IsInvalid(k8sClient.Create(ctx, object))).To(BeTrue())
			valid := map[string]interface{}{target: map[string]interface{}{"requests": map[string]interface{}{"cpu": "500m"}, "limits": map[string]interface{}{"cpu": "1"}}}
			stored := ingressResourceObject(valid)
			Expect(k8sClient.Create(ctx, stored)).To(Succeed())
			Expect(unstructured.SetNestedMap(stored.Object, resources, "spec", "resources")).To(Succeed())
			Expect(apierrors.IsInvalid(k8sClient.Update(ctx, stored))).To(BeTrue())
			Expect(ingressResourceRoundTrip(stored)).To(Equal(valid))
		},
			Entry("request exceeds limit", map[string]interface{}{"requests": map[string]interface{}{"cpu": "2"}, "limits": map[string]interface{}{"cpu": "1"}}),
			Entry("request exceeds limit above double precision", map[string]interface{}{"requests": map[string]interface{}{"memory": "9007199254740993"}, "limits": map[string]interface{}{"memory": "9007199254740992"}}),
			Entry("request-only extended resource", map[string]interface{}{"requests": map[string]interface{}{"example.com/device": "1"}}),
			Entry("unequal extended pair", map[string]interface{}{"requests": map[string]interface{}{"example.com/device": "1"}, "limits": map[string]interface{}{"example.com/device": "2"}}),
			Entry("fractional device request", map[string]interface{}{"requests": map[string]interface{}{"example.com/device": "0.5"}, "limits": map[string]interface{}{"example.com/device": "1"}}),
			Entry("fractional limit-only device", map[string]interface{}{"limits": map[string]interface{}{"example.com/device": "0.5"}}),
		)

		for _, resourceMap := range []string{"requests", "limits"} {
			DescribeTable("rejects invalid "+target+" "+resourceMap, func(value interface{}) {
				object := ingressResourceObject(map[string]interface{}{target: map[string]interface{}{resourceMap: map[string]interface{}{"cpu": value}}})
				Expect(apierrors.IsInvalid(k8sClient.Create(ctx, object))).To(BeTrue())
				validResources := map[string]interface{}{target: map[string]interface{}{resourceMap: map[string]interface{}{"cpu": "1"}}}
				stored := ingressResourceObject(validResources)
				Expect(k8sClient.Create(ctx, stored)).To(Succeed())
				Expect(unstructured.SetNestedMap(stored.Object, map[string]interface{}{target: map[string]interface{}{resourceMap: map[string]interface{}{"cpu": value}}}, "spec", "resources")).To(Succeed())
				Expect(apierrors.IsInvalid(k8sClient.Update(ctx, stored))).To(BeTrue())
				Expect(ingressResourceRoundTrip(stored)).To(Equal(validResources))
			},
				Entry("negative string", "-1m"),
				Entry("negative integer", int64(-1)),
				Entry("malformed string", "garbage"),
				Entry("65 character string", strings.Repeat("0", 64)+"1"),
				Entry("boolean is not a quantity", true),
			)
			It("bounds "+target+" "+resourceMap+" to eight entries", func() {
				entries := map[string]interface{}{}
				for index := 0; index < 8; index++ {
					entries[fmt.Sprintf("resource-%d", index)] = "1"
				}
				object := ingressResourceObject(map[string]interface{}{target: map[string]interface{}{resourceMap: entries}})
				Expect(k8sClient.Create(ctx, object)).To(Succeed())
				validResources := ingressResourceRoundTrip(object)
				entries["ninth-resource"] = "1"
				Expect(apierrors.IsInvalid(k8sClient.Create(ctx, ingressResourceObject(map[string]interface{}{target: map[string]interface{}{resourceMap: entries}})))).To(BeTrue())
				Expect(unstructured.SetNestedMap(object.Object, map[string]interface{}{target: map[string]interface{}{resourceMap: entries}}, "spec", "resources")).To(Succeed())
				Expect(apierrors.IsInvalid(k8sClient.Update(ctx, object))).To(BeTrue())
				Expect(ingressResourceRoundTrip(object)).To(Equal(validResources))
			})
		}
	}

	It("admits all four maximum-sized maps within runtime cost budgets", func() {
		resources := map[string]interface{}{}
		for _, target := range []string{"router", "haproxy"} {
			component := map[string]interface{}{}
			for _, resourceMap := range []string{"requests", "limits"} {
				entries := map[string]interface{}{}
				for index := 0; index < 8; index++ {
					entries[strings.Repeat("x", 318)+fmt.Sprint(index)] = strings.Repeat("0", 63) + "1"
				}
				component[resourceMap] = entries
			}
			resources[target] = component
		}
		object := ingressResourceObject(resources)
		Expect(k8sClient.Create(ctx, object)).To(Succeed())
		Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
		fmt.Fprintln(GinkgoWriter, "NE2926 admitted all four 8-entry maps, 64-character quantities and long keys without runtime cost failure")
	})

	It("prunes optional null values and null quantities instead of usable overrides", func() {
		for _, resources := range []interface{}{
			nil,
			map[string]interface{}{"router": nil, "haproxy": nil},
			map[string]interface{}{"router": map[string]interface{}{"requests": nil, "limits": nil, "claims": nil}},
			map[string]interface{}{"haproxy": map[string]interface{}{"requests": map[string]interface{}{"cpu": nil}}},
		} {
			object := ingressResourceObject(resources)
			Expect(k8sClient.Create(ctx, object)).To(Succeed())
			stored := ingressResourceRoundTrip(object)
			if stored != nil {
				for _, component := range stored.(map[string]interface{}) {
					Expect(component).NotTo(BeNil())
					for field, value := range component.(map[string]interface{}) {
						Expect(value).NotTo(BeNil())
						if field == "requests" || field == "limits" {
							Expect(value.(map[string]interface{})).NotTo(HaveKey("cpu"))
						}
					}
				}
			}
		}
	})

	It("updates independent partial maps and resets overrides without API defaults", func() {
		object := ingressResourceObject(map[string]interface{}{})
		Expect(k8sClient.Create(ctx, object)).To(Succeed())
		resources := map[string]interface{}{
			"router":  map[string]interface{}{"requests": map[string]interface{}{"cpu": "500m"}},
			"haproxy": map[string]interface{}{"limits": map[string]interface{}{"qat.intel.com/cy": "1"}},
		}
		Expect(unstructured.SetNestedMap(object.Object, resources, "spec", "resources")).To(Succeed())
		Expect(k8sClient.Update(ctx, object)).To(Succeed())
		Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
		resources["router"] = map[string]interface{}{}
		Expect(unstructured.SetNestedMap(object.Object, resources, "spec", "resources")).To(Succeed())
		Expect(k8sClient.Update(ctx, object)).To(Succeed())
		Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
		unstructured.RemoveNestedField(object.Object, "spec", "resources")
		Expect(k8sClient.Update(ctx, object)).To(Succeed())
		Expect(ingressResourceRoundTrip(object)).To(BeNil())
	})

	It("prunes unknown fields while Strict rejects them", func() {
		resources := map[string]interface{}{"router": map[string]interface{}{"unknown": true, "requests": map[string]interface{}{"cpu": "1"}}}
		Expect(k8sClient.Create(ctx, ingressResourceObject(resources), client.FieldValidation(metav1.FieldValidationStrict))).NotTo(Succeed())
		object := ingressResourceObject(resources)
		Expect(k8sClient.Create(ctx, object, client.FieldValidation(metav1.FieldValidationIgnore))).To(Succeed())
		delete(resources["router"].(map[string]interface{}), "unknown")
		Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
	})

	It("supports old objects and exposes complete resources only in preview schemas", func() {
		for _, featureSet := range []string{"DevPreviewNoUpgrade", "Default", "OKD", "TechPreviewNoUpgrade"} {
			definition := ingressResourceDefinition(featureSet)
			_, enabled := definition.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["resources"]
			Expect(enabled).To(Equal(featureSet == "DevPreviewNoUpgrade" || featureSet == "TechPreviewNoUpgrade"))
			_, err := envtest.InstallCRDs(cfg, envtest.CRDInstallOptions{CRDs: []*apiextensionsv1.CustomResourceDefinition{definition}})
			Expect(err).NotTo(HaveOccurred())
			resources := map[string]interface{}{"router": map[string]interface{}{"requests": map[string]interface{}{"cpu": "1"}}}
			Eventually(func() bool {
				probe := ingressResourceObject(resources)
				if err := k8sClient.Create(ctx, probe, client.DryRunAll, client.FieldValidation(metav1.FieldValidationIgnore)); err != nil {
					return !enabled
				}
				_, present, err := unstructured.NestedFieldNoCopy(probe.Object, "spec", "resources")
				return err == nil && present
			}, timeout).Should(Equal(enabled))
			object := ingressResourceObject(resources)
			delete(object.Object["spec"].(map[string]interface{}), "resources")
			Expect(k8sClient.Create(ctx, object)).To(Succeed())
			Expect(ingressResourceRoundTrip(object)).To(BeNil())
			Expect(unstructured.SetNestedMap(object.Object, resources, "spec", "resources")).To(Succeed())
			Expect(k8sClient.Update(ctx, object, client.FieldValidation(metav1.FieldValidationIgnore))).To(Succeed())
			if enabled {
				Expect(ingressResourceRoundTrip(object)).To(Equal(resources))
				Expect(k8sClient.Create(ctx, ingressResourceObject(resources), client.FieldValidation(metav1.FieldValidationStrict))).To(Succeed())
			} else {
				Expect(ingressResourceRoundTrip(object)).To(BeNil())
				Expect(k8sClient.Create(ctx, ingressResourceObject(resources), client.FieldValidation(metav1.FieldValidationStrict))).NotTo(Succeed())
			}
		}
	})
})
