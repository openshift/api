package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +openshift:compatibility-gen:level=1

// StableConfigType is a stable config type that may include TechPreviewNoUpgrade fields.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type StableConfigType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the StableConfigType.
	Spec StableConfigTypeSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the StableConfigType.
	Status StableConfigTypeStatus `json:"status,omitempty"`
}

// StableConfigTypeSpec is the desired state
type StableConfigTypeSpec struct {
	// coolNewField is a field that is for tech preview only.  On normal clusters this shouldn't be present
	//
	// +kubebuilder:validation:Optional
	// +openshift:enable:FeatureSets=TechPreviewNoUpgrade
	// +optional
	CoolNewField string `json:"coolNewField"`

	// stableField is a field that is present on default clusters and on tech preview clusters
	//
	// If empty, the platform will choose a good default, which may change over time without notice.
	//
	// +optional
	StableField string `json:"stableField"`

	// evolvingUnion demonstrates how to phase in new values into discriminated union
	// +optional
	EvolvingUnion EvolvingUnion `json:"evolvingUnion"`

	// celUnion demonstrates how to validate a discrminated union using CEL
	// +optional
	CELUnion CELUnion `json:"celUnion,omitempty"`
}

type EvolvingUnion struct {
	// type is the discriminator. It has different values for Default and for TechPreviewNoUpgrade
	// +kubebuilder:validation:Required
	Type EvolvingDiscriminator `json:"type,omitempty"`
}

// EvolvingDiscriminator defines the audit policy profile type.
// +openshift:validation:FeatureSetAwareEnum:featureSet=Default,enum="";StableValue
// +openshift:validation:FeatureSetAwareEnum:featureSet=TechPreviewNoUpgrade,enum="";StableValue;TechPreviewOnlyValue
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
	// +kubebuilder:validation:Required
	// +unionDiscriminator
	Type CELUnionDiscriminator `json:"type,omitempty"`

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
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +openshift:compatibility-gen:level=1

// StableConfigTypeList contains a list of StableConfigTypes.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type StableConfigTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StableConfigType `json:"items"`
}
