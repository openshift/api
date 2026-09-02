package overrides

import (
	"github.com/openshift/api/payload-command/upstreamgates/overrides/registry"
	apiserverfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/component-base/featuregate"
	genericfeatures "k8s.io/kubernetes/pkg/features"

	configv1 "github.com/openshift/api/config/v1"
)

func init() {
	for gate, override := range overrides {
		registry.RegisterGateOverride(gate, override)
	}
}

// NoOpinion signals that an upstream feature gate should not have an equivalent OpenShift feature gate created for it.
// This is not a real OpenShift FeatureSet and is a utility constant specifically for the tooling used to generate
// OpenShift-specific feature gates that represent upstream feature gates.
const NoOpinion configv1.FeatureSet = "NoOpinion"

var overrides = map[featuregate.Feature]registry.OpenShiftGateFields{
	genericfeatures.ServiceAccountTokenNodeBinding: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/4193",
		JiraComponent:          "apiserver-auth",
		ContactPerson:          "ibihim",
	},

	apiserverfeatures.MutatingAdmissionPolicy: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/3962",
		JiraComponent:          "kube-apiserver",
		ContactPerson:          "benluddy",
	},

	genericfeatures.MaxUnavailableStatefulSet: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/961",
		JiraComponent:          "apps",
		ContactPerson:          "atiratree",
	},

	genericfeatures.EventedPLEG: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/3386",
		JiraComponent:          "node",
		ContactPerson:          "sairameshv",
	},

	apiserverfeatures.KMSv1: {
		EnhancementPullRequest: "legacyFeatureGateWithoutEnhancement",
		JiraComponent:          "kube-apiserver",
		ContactPerson:          "dgrisonnet",
		EnabledFeatureSets:     []configv1.FeatureSet{configv1.Default, configv1.TechPreviewNoUpgrade, configv1.DevPreviewNoUpgrade},
	},

	genericfeatures.SELinuxMount: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/1710",
		JiraComponent:          "Storage / Kubernetes",
		ContactPerson:          "jsafrane",
	},

	genericfeatures.MutableCSINodeAllocatableCount: {
		EnhancementPullRequest: "https://github.com/kubernetes/enhancements/issues/4876",
		JiraComponent:          "Storage / Kubernetes External Components",
		ContactPerson:          "jsafrane",
	},

	// Default hardcoded "special" gates for enabling all alpha and all beta feature gates.
	// We don't want these, so always ignore them.
	featuregate.Feature("AllAlpha"): {
		EnabledFeatureSets: []configv1.FeatureSet{NoOpinion},
	},

	featuregate.Feature("AllBeta"): {
		EnabledFeatureSets: []configv1.FeatureSet{NoOpinion},
	},

	// TODO: these overrides are for example purposes. Remove or populate correctly before
	// merging.
	genericfeatures.DRADeviceTaintRules: {
		EnhancementPullRequest: "https://github.com/openshift/enhancements/pull/2084",
		JiraComponent:          "kube-apiserver",
		ContactPerson:          "bpalmer",
		GroupKindResources: []registry.GroupKindResource{
			{
				Group:    "resource.k8s.io",
				Kind:     "DeviceTaintRule",
				Resource: "devicetaintrules",
			},
		},
	},

	genericfeatures.PodCertificateRequest: {
		EnhancementPullRequest: "https://github.com/openshift/enhancements/pull/2084",
		JiraComponent:          "kube-apiserver",
		ContactPerson:          "bpalmer",
		GroupKindResources: []registry.GroupKindResource{
			{
				Group:    "certificates.k8s.io",
				Kind:     "PodCertificateRequest",
				Resource: "podcertificaterequests",
			},
		},
	},

	genericfeatures.CoordinatedLeaderElection: {
		EnhancementPullRequest: "https://github.com/openshift/enhancements/pull/2084",
		JiraComponent:          "kube-apiserver",
		ContactPerson:          "bpalmer",
		GroupKindResources: []registry.GroupKindResource{
			{
				Group:    "coordination.k8s.io",
				Kind:     "LeaseCandidate",
				Resource: "leasecandidates",
			},
		},
	},

	// Kubelet bootstrapping is failing due to our usage of the serverTLSBootstrap field in the Kubelet config
	// when we disable the RotateKubeletServerCertificate feature-gate in the default feature-set due.
	// It is enabled by default and beta upstream, so since we have built an explicit dependency on it when we
	// had no opinion, let's go back to that.
	genericfeatures.RotateKubeletServerCertificate: {
		EnabledFeatureSets: []configv1.FeatureSet{NoOpinion},
	},
}
