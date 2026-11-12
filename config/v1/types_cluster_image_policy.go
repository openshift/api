package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImagePolicy holds cluster-wide configuration for image signature verification
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterimagepolicies,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2310
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:compatibility-gen:level=1
type ClusterImagePolicy struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// spec contains the configuration for the cluster image policy.
	// +required
	Spec ClusterImagePolicySpec `json:"spec"`
	// status contains the observed state of the resource.
	// +optional
	Status ClusterImagePolicyStatus `json:"status"`
}

// CLusterImagePolicySpec is the specification of the ClusterImagePolicy custom resource.
// +openshift:validation:FeatureGateAwareXValidation:featureGate="",rule="self.policy.rootOfTrust.policyType == 'Reject' || (has(self.scopes) && size(self.scopes) > 0)",message="scopes must be specified unless policyType is Reject"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=ImagePolicyDenyByDefault,rule="self.policy.rootOfTrust.policyType != 'Reject' || has(self.appliesTo)",message="policyType Reject requires appliesTo Default"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=ImagePolicyDenyByDefault,rule="!has(self.appliesTo) || self.appliesTo != 'Default' || self.policy.rootOfTrust.policyType == 'Reject'",message="appliesTo Default requires policyType Reject"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=ImagePolicyDenyByDefault,rule="!has(self.appliesTo) || self.appliesTo != 'Default' || !has(self.scopes)",message="scopes must be empty when appliesTo is Default"
type ClusterImagePolicySpec struct {
	// appliesTo specifies whether this policy applies to images that do not match any registry-specific policy.
	// The only supported value is "Default".
	// The "Default" value is available only when the ImagePolicyDenyByDefault feature gate is enabled.
	// When omitted, the policy applies to the registry scopes listed in scopes. Omission is the v1 representation of registry-scoped policies.
	// When set to "Default", scopes must be omitted and policyType must be "Reject".
	// +optional
	// +openshift:enable:FeatureGate=ImagePolicyDenyByDefault
	AppliesTo PolicyAppliesTo `json:"appliesTo,omitempty"`
	// scopes is an optional field that defines the list of image identities assigned to a policy. Each item refers to a scope in a registry implementing the "Docker Registry HTTP API V2".
	// Scopes matching individual images are named Docker references in the fully expanded form, either using a tag or digest. For example, docker.io/library/busybox:latest (not busybox:latest).
	// More general scopes are prefixes of individual-image scopes, and specify a repository (by omitting the tag or digest), a repository
	// namespace, or a registry host (by only specifying the host name and possibly a port number) or a wildcard expression starting with `*.`, for matching all subdomains (not including a port number).
	// Wildcards are only supported for subdomain matching, and may not be used in the middle of the host, i.e.  *.example.com is a valid case, but example*.*.com is not.
	// This supports no more than 256 scopes in one object. If multiple scopes match a given image, only the policy requirements for the most specific scope apply. The policy requirements for more general scopes are ignored unless intersection merging is enabled.
	// In addition to setting a policy appropriate for your own deployed applications, make sure that a policy on the OpenShift image repositories
	// quay.io/openshift-release-dev/ocp-release, quay.io/openshift-release-dev/ocp-v4.0-art-dev (or on a more general scope) allows deployment of the OpenShift images required for cluster operation.
	// If a scope is configured in both the ClusterImagePolicy and the ImagePolicy, or if the scope in ImagePolicy is nested under one of the scopes from the ClusterImagePolicy, only the policy from the ClusterImagePolicy will be applied unless intersection merging is enabled in the namespace.
	// When specified, scopes must contain between 1 and 256 items. At least one scope is required when appliesTo is omitted; scopes must be omitted when appliesTo is "Default".
	// For additional details about the format, please refer to the document explaining the docker transport field,
	// which can be found at: https://github.com/containers/image/blob/main/docs/containers-policy.json.5.md#docker
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=256
	// +listType=set
	Scopes []ImageScope `json:"scopes,omitempty"`
	// policy is a required field that contains configuration to allow scopes to be verified, and defines how
	// images not matching the verification policy will be treated.
	// +required
	Policy ImageSigstoreVerificationPolicy `json:"policy"`
}

// +k8s:deepcopy-gen=true
type ClusterImagePolicyStatus struct {
	// conditions provide details on the status of this API Resource.
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImagePolicyList is a list of ClusterImagePolicy resources
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ClusterImagePolicyList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ListMeta `json:"metadata"`

	// items is a list of ClusterImagePolices
	// +kubebuilder:validation:MaxItems=1000
	// +required
	Items []ClusterImagePolicy `json:"items"`
}
