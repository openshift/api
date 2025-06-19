package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRDCompatibilityRequirement expresses a set of requirements on a target CRD.
// It is used to ensure compatibility between different actors using the same
// CRD.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:file-pattern=cvoRunLevel=0000_20,operatorName=crd-compatibility-checker,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=crdcompatibilityrequirements,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/XXXX
type CRDCompatibilityRequirement struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// spec is the specification of the desired behavior of the CRD Compatibility Requirement.
	// +required
	Spec CRDCompatibilityRequirementSpec `json:"spec,omitzero"`

	// status is the most recently observed status of the CRD Compatibility Requirement.
	// +optional
	Status CRDCompatibilityRequirementStatus `json:"status,omitzero"`
}

// CRDCompatibilityRequirementSpec is the specification of the desired behavior of the CRD Compatibility Requirement.
type CRDCompatibilityRequirementSpec struct {
	// crdRef is the name of the target CRD. The target CRD is not required to
	// exist, as we may legitimately place requirements on it before it is
	// created.  The observed CRD is given in status.observedCRD, which will be
	// empty if no CRD is observed.
	// This field is required.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	CRDRef string `json:"crdRef,omitempty"`

	// creatorDescription is a string describing the owner of this CRDCompatibilityRequirement. It will be printed in any error or
	// warning emitted by any of the CRDCompatibilityRequirement's webhooks. It should indicate to the recipient who they need to coordinate
	// with in order to safely update the target CRD. The message emitted will be: "This requirement was added by <creatorDescription>".
	// This field is required.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	CreatorDescription string `json:"creatorDescription,omitempty"`

	// compatibilityCRD contains the CRD which is required by the creator of this CRDCompatibilityRequirement. CRD Compatibility Checker will
	// ensure that only a target CRD compatible with compatibilityCRD may be admitted.
	// This field is required.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000000
	// +required
	CompatibilityCRD string `json:"compatibilityCRD,omitempty"`

	// crdAdmitAction determines whether the CRD admission controller will Enforce or Warn if the CRD presented is not compatible.
	// This field is required.
	// +required
	CRDAdmitAction CRDAdmitAction `json:"crdAdmitAction,omitempty"`
}

// CRDAdmitAction determines the action taken when a CRD is not compatible.
// +kubebuilder:validation:Enum=Enforce;Warn
type CRDAdmitAction string

const (
	// CRDAdmitActionEnforce means that incompatible CRDs will be rejected.
	CRDAdmitActionEnforce CRDAdmitAction = "Enforce"

	// CRDAdmitActionWarn means that incompatible CRDs will be allowed but a warning will be generated.
	CRDAdmitActionWarn CRDAdmitAction = "Warn"
)

// CRDCompatibilityRequirementStatus defines the observed status of the CRD Compatibility Requirement.
// +kubebuilder:validation:MinProperties=1
type CRDCompatibilityRequirementStatus struct {
	// conditions is a list of conditions and their status.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=16
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// observedCRD documents the uid and generation of the CRD object when the current status was written.
	// This field will not be emitted if the target CRD does not exist or could not be retrieved.
	// +optional
	ObservedCRD ObservedCRD `json:"observedCRD,omitzero"`
}

// ObservedCRD contains information about the observed target CRD.
// +kubebuilder:validation:MinProperties=1
type ObservedCRD struct {
	// uid is the uid of the observed CRD.
	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Format=uuid
	// +required
	UID string `json:"uid,omitempty"`

	// generation is the observed generation of the CRD.
	// +kubebuilder:validation:Minimum=1
	// +required
	Generation int64 `json:"generation,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CRDCompatibilityRequirementList is a collection of CRDCompatibilityRequirements.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type CRDCompatibilityRequirementList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a list of CRDCompatibilityRequirements.
	// +kubebuilder:validation:MaxItems=1000
	// +optional
	Items []CRDCompatibilityRequirement `json:"items,omitempty"`
}
