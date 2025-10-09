package v1alpha1

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
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
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec is the specification of the desired behavior of the CRD Compatibility Requirement.
	// +required
	Spec CRDCompatibilityRequirementSpec `json:"spec,omitzero"`

	// status is the most recently observed status of the CRD Compatibility Requirement.
	// +optional
	Status CRDCompatibilityRequirementStatus `json:"status,omitzero"`
}

// CRDCompatibilityRequirementSpec is the specification of the desired behavior of the CRD Compatibility Requirement.
type CRDCompatibilityRequirementSpec struct {
	// compatibilitySchema defines the schema used by crdSchemaValidation and objectSchemaValidation.
	// This field is required.
	// +required
	CompatibilitySchema CompatibilitySchema `json:"compatibilitySchema,omitempty,omitzero"`

	// crdSchemaValidation ensures that updates to the installed CRD are compatible with this compatibility requirement.
	// This field is optional.
	// +optional
	CRDSchemaValidation CRDSchemaValidation `json:"crdSchemaValidation,omitempty,omitzero"`

	// objectSchemaValidation ensures that matching objects conform to compatibilitySchema.
	// This field is optional.
	// +optional
	ObjectSchemaValidation ObjectSchemaValidation `json:"objectSchemaValidation,omitempty,omitzero"`
}

// CompatibilitySchema defines the schema used by crdSchemaValidation and objectSchemaValidation.
type CompatibilitySchema struct {
	// crdYAML contains the complete YAML document of the CRD for schema and object validation purposes.
	// This field is required.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1572864
	// +required
	CRDYAML string `json:"crdYAML,omitempty"`

	// requireVersions specifies which versions we will automatically extract from the yaml and require.
	// Valid options are:
	//   StorageOnly - only storage version(s) required for compatibility. Users can create/update
	//     objects using any served version. additionalVersions are applied on top of this.
	//   All - all versions defined in the CRD are required for compatibility.
	// This field is required.
	// +required
	RequireVersions RequireVersions `json:"requireVersions,omitempty"`

	// additionalVersions is a set of versions to require in addition to those discovered by requireVersions.
	// Overlap with requireVersions is explicitly permitted.
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=255
	// +kubebuilder:validation:MaxItems=255
	// +listType=set
	// +optional
	AdditionalVersions []string `json:"additionalVersions,omitempty"`

	// excludeFields is a set of fields in the yaml which will not be validated by either
	// crdSchemaValidation or objectSchemaValidation.
	// FIXME(chrischdi): explain the format which is
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=8192
	// +kubebuilder:validation:MaxItems=1024
	// +listType=set
	// +optional
	ExcludeFields []string `json:"excludeFields,omitempty"`
}

// CRDSchemaValidation ensures that updates to the installed CRD are compatible with this compatibility requirement.
type CRDSchemaValidation struct {
	// action determines whether violations are not admitted (Enforce) or admitted with an API warning (Warn).
	// Valid options are:
	//   Enforce - incompatible CRDs will be rejected and not admitted to the cluster.
	//   Warn - incompatible CRDs will be allowed but a warning will be generated in the API response.
	// This field is required.
	// +required
	Action CRDAdmitAction `json:"action,omitempty"`
}

// ObjectSchemaValidation ensures that matching objects conform to compatibilitySchema.
type ObjectSchemaValidation struct {
	// action determines whether violations are not admitted (Enforce) or admitted with an API warning (Warn).
	// Valid options are:
	//   Enforce - incompatible Objects will be rejected and not admitted to the cluster.
	//   Warn - incompatible Objects will be allowed but a warning will be generated in the API response.
	// This field is required.
	// +required
	Action CRDAdmitAction `json:"action,omitempty"`

	// namespaceSelector defines the namespaceSelector field of the resulting ValidatingWebhookConfiguration.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// objectSelector defines the objectSelector field of the resulting ValidatingWebhookConfiguration.
	// +optional
	ObjectSelector *metav1.LabelSelector `json:"objectSelector,omitempty"`

	// matchConditions defines the matchConditions field of the resulting ValidatingWebhookConfiguration.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	// +optional
	MatchConditions []admissionregistrationv1.MatchCondition `json:"matchConditions,omitempty"`
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

// RequireVersions specifies which versions we will automatically extract from the yaml and require.
// +kubebuilder:validation:Enum=StorageOnly;All
type RequireVersions string

const (
	// RequireVersionsStorageOnly means only storage versions will be required.
	RequireVersionsStorageOnly RequireVersions = "StorageOnly"

	// RequireVersionsAll means all versions will be required.
	RequireVersionsAll RequireVersions = "All"
)

// CRDCompatibilityRequirement's Progressing condition and corresponding reasons.
const (
	// CRDCompatibilityRequirementProgressing is false if the spec has been completely reconciled.
	// True indicates that reconciliation is still in progress and the current status does not represent
	// a stable state. Progressing false with an error reason indicates that the object cannot be reconciled.
	CRDCompatibilityRequirementProgressing string = "Progressing"

	// CRDCompatibilityRequirementConfigurationErrorReason surfaces when reconciliation cannot progress due to an invalid spec.
	CRDCompatibilityRequirementConfigurationErrorReason string = "ConfigurationError"

	// CRDCompatibilityRequirementTransientErrorReason surfaces when reconciliation failed due to an error that can be retried.
	CRDCompatibilityRequirementTransientErrorReason string = "TransientError"

	// CRDCompatibilityRequirementUpToDateReason surfaces when reconciliation completed successfully.
	CRDCompatibilityRequirementUpToDateReason string = "UpToDate"
)

// CRDCompatibilityRequirement's Admitted condition and corresponding reasons.
const (
	// CRDCompatibilityRequirementAdmitted is true if the requirement has been configured in the validating webhook,
	// otherwise false.
	CRDCompatibilityRequirementAdmitted string = "Admitted"

	// CRDCompatibilityRequirementAdmittedReason surfaces when the requirement has been configured in the validating webhook.
	CRDCompatibilityRequirementAdmittedReason string = "Admitted"

	// CRDCompatibilityRequirementNotAdmittedReason surfaces when the requirement has not been configured in the validating webhook.
	CRDCompatibilityRequirementNotAdmittedReason string = "NotAdmitted"
)

// CRDCompatibilityRequirement's Compatible condition and corresponding reasons.
const (
	// CRDCompatibilityRequirementCompatible is true if the observed CRD is compatible with the requirement,
	// otherwise false. Note that Compatible may be false when adding a new requirement which the existing
	// CRD does not meet.
	CRDCompatibilityRequirementCompatible string = "Compatible"

	// CRDCompatibilityRequirementRequirementsNotMetReason surfaces when a CRD exists, and it is not compatible with this requirement.
	CRDCompatibilityRequirementRequirementsNotMetReason string = "RequirementsNotMet"

	// CRDCompatibilityRequirementCRDDoesNotExistReason surfaces when the referenced CRD does not exist.
	CRDCompatibilityRequirementCRDDoesNotExistReason string = "CRDDoesNotExist"

	// CRDCompatibilityRequirementCompatibleWithWarningsReason surfaces when the CRD exists and is compatible with this requirement, but Message contains one or more warning messages.
	CRDCompatibilityRequirementCompatibleWithWarningsReason string = "CompatibleWithWarnings"

	// CRDCompatibilityRequirementCompatibleReason surfaces when the CRD exists and is compatible with this requirement.
	CRDCompatibilityRequirementCompatibleReason string = "Compatible"
)

// CRDCompatibilityRequirementStatus defines the observed status of the CRD Compatibility Requirement.
// +kubebuilder:validation:MinProperties=1
type CRDCompatibilityRequirementStatus struct {
	// conditions is a list of conditions and their status.
	// Known condition types are Progressing, Admitted, Compatible.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:MinItems=1
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// observedCRD documents the uid and generation of the CRD object when the current status was written.
	// This field will not be emitted if the target CRD does not exist or could not be retrieved.
	// +optional
	ObservedCRD ObservedCRD `json:"observedCRD,omitzero"`

	// crdName is the name of the target CRD. The target CRD is not required to
	// exist, as we may legitimately place requirements on it before it is
	// created.  The observed CRD is given in status.observedCRD, which will be
	// empty if no CRD is observed.
	// This field is optional.
	// crdRef must be at most 253 characters in length and must consist only of lower-case alphanumeric characters, periods (.) and hyphens (-). Each period separated label must start and end with an alphanumeric character and be at most 63 characters in length.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +optional
	CRDName string `json:"crdName,omitempty"`
}

// ObservedCRD contains information about the observed target CRD.
// +kubebuilder:validation:MinProperties=1
type ObservedCRD struct {
	// uid is the uid of the observed CRD.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:Format=uuid
	// +kubebuilder:validation:XValidation:rule="!format.uuid().validate(self).hasValue()",message="uid must be a valid UUID. It must consist only of lower-case hexadecimal digits, in 5 hyphenated blocks, where the blocks are of length 8-4-4-4-12 respectively."
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
	metav1.ListMeta `json:"metadata,omitzero"`

	// items is a list of CRDCompatibilityRequirements.
	// +kubebuilder:validation:MaxItems=1000
	// +optional
	Items []CRDCompatibilityRequirement `json:"items,omitempty"`
}
