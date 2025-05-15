package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=1
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/xxx
// +openshift:file-pattern=cvoRunLevel=0000_50,operatorName=my-operator,operatorOrdering=01

// StableConfigType is a stable config type that may include TechPreviewNoUpgrade fields.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=stableconfigtypes,scope=Cluster
// +kubebuilder:subresource:status
type StableConfigType struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the StableConfigType.
	Spec StableConfigTypeSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the StableConfigType.
	Status StableConfigTypeStatus `json:"status,omitempty"`
}

// StableConfigTypeSpec is the desired state
// +openshift:validation:FeatureGateAwareXValidation:featureGate=Example,rule="has(oldSelf.coolNewField) ? has(self.coolNewField) : true",message="coolNewField may not be removed once set"
// +openshift:validation:FeatureGateAwareXValidation:requiredFeatureGate=Example;Example2,rule="has(oldSelf.stableField) ? has(self.stableField) : true",message="stableField may not be removed once set (this should only show up with both the Example and Example2 feature gates)"
type StableConfigTypeSpec struct {
	// coolNewField is a field that is for tech preview only.  On normal clusters this shouldn't be present
	//
	// +openshift:enable:FeatureGate=Example
	// +optional
	CoolNewField string `json:"coolNewField"`

	// stableField is a field that is present on default clusters and on tech preview clusters
	//
	// If empty, the platform will choose a good default, which may change over time without notice.
	//
	// +optional
	StableField string `json:"stableField"`

	// immutableField is a field that is immutable once the object has been created.
	// It is required at all times.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="immutableField is immutable"
	// +required
	ImmutableField string `json:"immutableField"`

	// optionalImmutableField is a field that is immutable once set.
	// It is optional but may not be changed once set.
	// +kubebuilder:validation:XValidation:rule="oldSelf == '' || self == oldSelf",message="optionalImmutableField is immutable once set"
	// +optional
	OptionalImmutableField string `json:"optionalImmutableField"`

	// evolvingUnion demonstrates how to phase in new values into discriminated union
	// +optional
	EvolvingUnion EvolvingUnion `json:"evolvingUnion"`

	// celUnion demonstrates how to validate a discrminated union using CEL
	// +optional
	CELUnion CELUnion `json:"celUnion,omitempty"`

	// nonZeroDefault is a demonstration of creating an integer field that has a non zero default.
	// It required two default tags (one for CRD generation, one for client generation) and must have `omitempty` and be optional.
	// A minimum value is added to demonstrate that a zero value would not be accepted.
	// +kubebuilder:default:=8
	// +default=8
	// +kubebuilder:validation:Minimum:=8
	// +optional
	NonZeroDefault int32 `json:"nonZeroDefault,omitempty"`

	// evolvingCollection demonstrates how to have a collection where the maximum number of items varies on cluster type.
	// For default clusters, this will be "1" but on TechPreview clusters, this value will be "3".
	// +openshift:validation:FeatureGateAwareMaxItems:featureGate="",maxItems=1
	// +openshift:validation:FeatureGateAwareMaxItems:featureGate=Example,maxItems=3
	// +optional
	// +listType=atomic
	EvolvingCollection []string `json:"evolvingCollection,omitempty"`

	// set demonstrates how to define and validate set of strings
	// +optional
	Set StringSet `json:"set,omitempty"`

	// subdomainNameField represents a kubenetes name field.
	// The intention is that it validates the name in the same way metadata.Name is validated.
	// That is, it is a DNS-1123 subdomain.
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	// +optional
	SubdomainNameField string `json:"subdomainNameField,omitempty"`

	// subnetsWithExclusions demonstrates how to validate a list of subnets with exclusions
	// +optional
	SubnetsWithExclusions SubnetsWithExclusions `json:"subnetsWithExclusions,omitempty"`
}

// SetValue defines the types allowed in string set type
// +kubebuilder:validation:Enum:=Foo;Bar;Baz;Qux;Corge
type SetValue string

// StringSet defines the set of strings
// +listType=set
// +kubebuilder:validation:XValidation:rule="self.all(x,self.exists_one(y,x == y))"
// +kubebuilder:validation:MaxItems=5
type StringSet []SetValue

type EvolvingUnion struct {
	// type is the discriminator. It has different values for Default and for TechPreviewNoUpgrade
	// +required
	Type EvolvingDiscriminator `json:"type"`
}

// EvolvingDiscriminator defines the audit policy profile type.
// +openshift:validation:FeatureGateAwareEnum:featureGate="",enum="";StableValue
// +openshift:validation:FeatureGateAwareEnum:featureGate=Example,enum="";StableValue;TechPreviewOnlyValue
type EvolvingDiscriminator string

const (
	// "StableValue" is always present.
	StableValue EvolvingDiscriminator = "StableValue"

	// "TechPreviewOnlyValue" should only be allowed when TechPreviewNoUpgrade is set in the cluster
	TechPreviewOnlyValue EvolvingDiscriminator = "TechPreviewOnlyValue"
)

// CELUnion demonstrates how to use a discriminated union and how to validate it using CEL.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'RequiredMember' ?  has(self.requiredMember) : !has(self.requiredMember)",message="requiredMember is required when type is RequiredMember, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'OptionalMember' ?  true : !has(self.optionalMember)",message="optionalMember is forbidden when type is not OptionalMember"
// +union
type CELUnion struct {
	// type determines which of the union members should be populated.
	// +required
	// +unionDiscriminator
	Type CELUnionDiscriminator `json:"type"`

	// requiredMember is a union member that is required.
	// +unionMember
	RequiredMember *string `json:"requiredMember,omitempty"`

	// optionalMember is a union member that is optional.
	// +unionMember,optional
	OptionalMember *string `json:"optionalMember,omitempty"`
}

// CELUnionDiscriminator is a union discriminator for the CEL union.
// +kubebuilder:validation:Enum:="RequiredMember";"OptionalMember";"EmptyMember"
type CELUnionDiscriminator string

const (
	// RequiredMember represents a required union member.
	RequiredMember CELUnionDiscriminator = "RequiredMember"

	// OptionalMember represents an optional union member.
	OptionalMember CELUnionDiscriminator = "OptionalMember"

	// EmptyMember represents an empty union member.
	EmptyMember CELUnionDiscriminator = "EmptyMember"
)

// StableConfigTypeStatus defines the observed status of the StableConfigType.
type StableConfigTypeStatus struct {
	// Represents the observations of a foo's current state.
	// Known .status.conditions.type are: "Available", "Progressing", and "Degraded"
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// immutableField is a field that is immutable once the object has been created.
	// It is required at all times.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="immutableField is immutable"
	// +optional
	ImmutableField string `json:"immutableField,omitempty"`
}

// SubnetsWithExclusions is used to validate a list of subnets with exclusions.
// It demonstrates how exclusions should be validated as subnetworks of the networks listed in the subnets field.
// +kubebuilder:validation:XValidation:rule="!has(self.excludeSubnets) || self.excludeSubnets.all(e, self.subnets.exists(s, cidr(s).containsCIDR(cidr(e))))",message="excludeSubnets must be subnetworks of the networks specified in the subnets field",fieldPath=".excludeSubnets"
type SubnetsWithExclusions struct {
	// subnets is a list of subnets.
	// It may contain up to 2 subnets.
	// The list may be either 1 IPv4 subnet, 1 IPv6 subnet, or 1 of each.
	// +kubebuilder:validation:XValidation:rule="size(self) != 2 || !isCIDR(self[0]) || !isCIDR(self[1]) || cidr(self[0]).ip().family() != cidr(self[1]).ip().family()",message="subnets must not contain 2 subnets of the same IP family"
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +listType=atomic
	// +required
	Subnets []CIDR `json:"subnets"`

	// excludeSubnets is a list of CIDR exclusions.
	// The subnets in this list must be subnetworks of the subnets in the subnets list.
	// +kubebuilder:validation:MaxItems=25
	// +optional
	ExcludeSubnets []CIDR `json:"excludeSubnets,omitempty"`
}

// CIDR is used to validate a CIDR notation network.
// The longest CIDR notation is 43 characters.
// +kubebuilder:validation:XValidation:rule="isCIDR(self)",message="value must be a valid CIDR"
// +kubebuilder:validation:MaxLength:=43
type CIDR string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=1

// StableConfigTypeList contains a list of StableConfigTypes.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type StableConfigTypeList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []StableConfigType `json:"items"`
}
