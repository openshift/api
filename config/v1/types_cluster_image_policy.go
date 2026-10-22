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
// +kubebuilder:validation:XValidation:rule="(!has(self.appliesTo) || self.appliesTo == 'Registries') ? ((has(self.registries) && size(self.registries) > 0) || (has(self.scopes) && size(self.scopes) > 0)) : (!has(self.registries) && !has(self.scopes))",message="registries or scopes must be specified when appliesTo is Registries or omitted, and both must be empty when appliesTo is Default"
// +kubebuilder:validation:XValidation:rule="!(has(self.registries) && size(self.registries) > 0 && has(self.scopes) && size(self.scopes) > 0)",message="registries and scopes are mutually exclusive, specify only one"
type ClusterImagePolicySpec struct {
	// +openshift:enable:FeatureGate=ImagePolicyDenyByDefault
	// appliesTo specifies whether this policy applies to all images as a default fallback or to specific registry scopes.
	// Valid values are "Default" and "Registries".
	// When omitted, defaults to "Registries" for backward compatibility with existing resources.
	// When set to "Default", this policy acts as the fallback for all images that don't match any other registry-specific policy.
	// This maps to the "default" entry in policy.json and the registries field (or deprecated scopes field) must be empty.
	// When set to "Registries", this policy applies only to images matching the registries specified in the registries field
	// (or scopes field for backward compatibility). This maps to entries under transports.docker in policy.json.
	//
	// Conflict resolution when multiple ClusterImagePolicy resources target the same scope:
	//
	// Policy Type Compatibility:
	// - Signature verification types (PublicKey, FulcioCAWithRekor, PKI) are compatible and can be merged.
	//   Multiple policies with these types targeting the same scope will all be applied, allowing images to be
	//   verified against any of the configured signing authorities.
	// - Reject and InsecureAcceptAnything types are only compatible with themselves.
	// - Mixing incompatible policy types (e.g., Reject + PublicKey, InsecureAcceptAnything + FulcioCAWithRekor)
	//   for the same scope will cause the machine config rollout to fail with a validation error.
	//
	// Precedence for appliesTo: Default:
	// - Multiple ClusterImagePolicy resources can have appliesTo: Default if they have compatible policy types.
	//   Compatible policies are merged, allowing images to satisfy any of the configured policies.
	//   For example, two PublicKey policies will allow images signed by either key.
	// - If multiple ClusterImagePolicy resources have appliesTo: Default with incompatible policy types
	//   (e.g., Reject + PublicKey), the machine config rollout fails with a validation error listing the
	//   conflicting policies.
	// - A ClusterImagePolicy with appliesTo: Default takes precedence over any namespace-scoped ImagePolicy with
	//   appliesTo: Default. The namespace ImagePolicy default is ignored when a cluster default exists.
	// +optional
	// +kubebuilder:default=Registries
	AppliesTo *PolicyAppliesTo `json:"appliesTo,omitempty"`
	// +openshift:enable:FeatureGate=ImagePolicyDenyByDefault
	// registries defines the list of registry scopes this policy applies to.
	// This is the preferred field for specifying registry scopes.
	// Each item refers to a scope in a registry implementing the "Docker Registry HTTP API V2".
	// This field is required when appliesTo is "Registries" or omitted, unless the deprecated scopes field is used.
	// This field must be empty when appliesTo is "Default".
	// This field is mutually exclusive with the deprecated scopes field.
	// Scopes matching individual images are named Docker references in the fully expanded form, either using a tag or digest.
	// For example, docker.io/library/busybox:latest (not busybox:latest).
	// More general scopes are prefixes of individual-image scopes, and specify a repository (by omitting the tag or digest),
	// a repository namespace, or a registry host (by only specifying the host name and possibly a port number)
	// or a wildcard expression starting with `*.`, for matching all subdomains (not including a port number).
	// Wildcards are only supported for subdomain matching, and may not be used in the middle of the host,
	// i.e. *.example.com is a valid case, but example*.*.com is not.
	// This supports no more than 256 registry scopes. If multiple scopes match a given image, only the policy
	// requirements for the most specific scope apply. The policy requirements for more general scopes are ignored.
	// In addition to setting a policy appropriate for your own deployed applications, make sure that a policy
	// on the OpenShift image repositories quay.io/openshift-release-dev/ocp-release,
	// quay.io/openshift-release-dev/ocp-v4.0-art-dev (or on a more general scope) allows deployment
	// of the OpenShift images required for cluster operation.
	// If a scope is configured in both the ClusterImagePolicy and the ImagePolicy, or if the scope in ImagePolicy
	// is nested under one of the scopes from the ClusterImagePolicy, only the policy from the ClusterImagePolicy
	// will be applied.
	// For additional details about the format, please refer to the document explaining the docker transport field,
	// which can be found at: https://github.com/containers/image/blob/main/docs/containers-policy.json.5.md#docker
	// +optional
	// +kubebuilder:validation:MaxItems=256
	// +listType=set
	Registries []ImageScope `json:"registries,omitempty"`
	// scopes defines the list of registry scopes this policy applies to.
	// Deprecated: Use the registries field instead. This field is maintained for backward compatibility only
	// and will be removed in a future API version. When both registries and scopes are specified, validation will fail
	// as they are mutually exclusive. When appliesTo is "Registries" or omitted and registries is empty, this field
	// will be used for backward compatibility with existing resources created before the registries field was introduced.
	// This field is mutually exclusive with the registries field.
	// For the same validation rules and format details, see the registries field documentation.
	// +optional
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
