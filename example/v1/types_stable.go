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
	EvolvingUnion EvolvingUnion `json:"evolvingUnion,omitempty"`
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
