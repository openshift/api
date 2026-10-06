package tests

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/komega"
)

const ingressCRDName = "ingresses.operator.openshift.io"

var _ = Describe("[operator.openshift.io/v1] Ingress two-version storage compatibility", func() {
	It("round-trips spec and status through v1alpha1 and v1 with v1 storage", func() {
		crd, err := loadCRDFromFile(filepath.Join(
			"..",
			"operator",
			"v1",
			"zz_generated.crd-manifests",
			"0000_50_ingress_02_ingresses.crd.yaml",
		))
		Expect(err).NotTo(HaveOccurred())
		expectIngressVersionConfiguration(crd)

		crdOptions := envtest.CRDInstallOptions{
			CRDs: []*apiextensionsv1.CustomResourceDefinition{crd.DeepCopy()},
		}
		crds, err := envtest.InstallCRDs(cfg, crdOptions)
		Expect(err).NotTo(HaveOccurred())
		Expect(crds).To(HaveLen(1))
		Expect(envtest.WaitForCRDs(cfg, crds, crdOptions)).To(Succeed())

		DeferCleanup(func() {
			for _, version := range []string{"v1", "v1alpha1"} {
				Expect(k8sClient.DeleteAllOf(ctx, ingressForVersion(version))).To(Succeed())
			}
			Expect(envtest.UninstallCRDs(cfg, crdOptions)).To(Succeed())
			Eventually(komega.Get(crd)).Should(Not(Succeed()))
		})

		installedCRD := &apiextensionsv1.CustomResourceDefinition{}
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: ingressCRDName}, installedCRD)).To(Succeed())
		expectIngressVersionConfiguration(installedCRD)

		alphaIngress := ingressForVersion("v1alpha1")
		alphaIngress.SetName("cluster")
		Expect(unstructured.SetNestedField(alphaIngress.Object, "Managed", "spec", "gatewayAPI", "managementMode")).To(Succeed())
		Expect(k8sClient.Create(ctx, alphaIngress)).To(Succeed())

		v1Ingress := ingressForVersion("v1")
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "cluster"}, v1Ingress)).To(Succeed())
		expectIngressManagementMode(v1Ingress, "Managed")

		Expect(unstructured.SetNestedField(v1Ingress.Object, "Unmanaged", "spec", "gatewayAPI", "managementMode")).To(Succeed())
		Expect(k8sClient.Update(ctx, v1Ingress)).To(Succeed())

		alphaIngress = ingressForVersion("v1alpha1")
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "cluster"}, alphaIngress)).To(Succeed())
		expectIngressManagementMode(alphaIngress, "Unmanaged")

		Expect(unstructured.SetNestedField(alphaIngress.Object, int64(1), "status", "observedGeneration")).To(Succeed())
		Expect(k8sClient.Status().Update(ctx, alphaIngress)).To(Succeed())

		v1Ingress = ingressForVersion("v1")
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "cluster"}, v1Ingress)).To(Succeed())
		expectIngressObservedGeneration(v1Ingress, 1)

		Expect(unstructured.SetNestedField(v1Ingress.Object, int64(2), "status", "observedGeneration")).To(Succeed())
		Expect(k8sClient.Status().Update(ctx, v1Ingress)).To(Succeed())

		alphaIngress = ingressForVersion("v1alpha1")
		Expect(k8sClient.Get(ctx, client.ObjectKey{Name: "cluster"}, alphaIngress)).To(Succeed())
		expectIngressObservedGeneration(alphaIngress, 2)

		Eventually(func(g Gomega) {
			storedCRD := &apiextensionsv1.CustomResourceDefinition{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: ingressCRDName}, storedCRD)).To(Succeed())
			g.Expect(storedCRD.Status.StoredVersions).To(ConsistOf("v1"))
		}).Should(Succeed())
	})
})

func ingressForVersion(version string) *unstructured.Unstructured {
	ingress := &unstructured.Unstructured{}
	ingress.SetAPIVersion("operator.openshift.io/" + version)
	ingress.SetKind("Ingress")
	return ingress
}

func expectIngressVersionConfiguration(crd *apiextensionsv1.CustomResourceDefinition) {
	GinkgoHelper()

	Expect(crd.Spec.Versions).To(HaveLen(2))
	versions := map[string]apiextensionsv1.CustomResourceDefinitionVersion{}
	for _, version := range crd.Spec.Versions {
		versions[version.Name] = version
	}

	Expect(versions).To(HaveKey("v1"))
	Expect(versions["v1"].Served).To(BeTrue())
	Expect(versions["v1"].Storage).To(BeTrue())
	Expect(versions).To(HaveKey("v1alpha1"))
	Expect(versions["v1alpha1"].Served).To(BeTrue())
	Expect(versions["v1alpha1"].Storage).To(BeFalse())
}

func expectIngressManagementMode(ingress *unstructured.Unstructured, expected string) {
	GinkgoHelper()

	actual, found, err := unstructured.NestedString(ingress.Object, "spec", "gatewayAPI", "managementMode")
	Expect(err).NotTo(HaveOccurred())
	Expect(found).To(BeTrue())
	Expect(actual).To(Equal(expected))
}

func expectIngressObservedGeneration(ingress *unstructured.Unstructured, expected int64) {
	GinkgoHelper()

	actual, found, err := unstructured.NestedInt64(ingress.Object, "status", "observedGeneration")
	Expect(err).NotTo(HaveOccurred())
	Expect(found).To(BeTrue())
	Expect(actual).To(Equal(expected))
}
