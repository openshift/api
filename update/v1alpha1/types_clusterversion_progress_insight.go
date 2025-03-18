package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterVersionProgressInsight provides summary information about an ongoing cluster control plane update in Standalone clusters
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=clusterversionprogressinsights,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2012
// +openshift:file-pattern=cvoRunLevel=0000_00,operatorName=cluster-version-operator,operatorOrdering=02
// +openshift:enable:FeatureGate=UpgradeStatus
// +kubebuilder:metadata:annotations="description=Provides summary information about an ongoing cluster control plane update in Standalone clusters."
// +kubebuilder:metadata:annotations="displayName=ClusterVersionProgressInsights"
// ClusterVersionProgressInsight reports the state of a ClusterVersion resource (which represents a control plane
// update in standalone clusters), during a cluster update.
type ClusterVersionProgressInsight struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is empty for now, ClusterVersionProgressInsight is purely status-reporting API. In the future spec may be used to hold
	// configuration to drive what information is surfaced and how
	// +required
	Spec ClusterVersionProgressInsightSpec `json:"spec"`
	// status exposes the health and status of the ongoing cluster update
	// +optional
	Status ClusterVersionProgressInsightStatus `json:"status"`
}

// ClusterVersionProgressInsightSpec is empty for now, ClusterVersionProgressInsightSpec is purely status-reporting API. In the future spec may be used
// to hold configuration to drive what information is surfaced and how
type ClusterVersionProgressInsightSpec struct {
}

// ClusterVersionProgressInsightStatus reports the state of a ClusterVersion resource (which represents a control plane
// update in standalone clusters), during the update.
type ClusterVersionProgressInsightStatus struct {
	// conditions provides detailed observed conditions about ClusterVersion. It contains at most 10 items.
	// Known conditions are:
	// - Updating: whether the control plane (represented by this ClusterVersion) is updating
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +optional
	// +kubebuilder:validation:MaxItems=10
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// name is equal to the name of the corresponding clusterversions.config.openshift.io resource, typically 'version'
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Name string `json:"name"`

	// assessment is the assessment of the control plane update process. Valid values are: Unknown, Progressing, Completed, Degraded
	// +required
	Assessment ControlPlaneAssessment `json:"assessment"`

	// versions contains the original and target versions involved in the update
	// +required
	Versions ControlPlaneUpdateVersions `json:"versions"`

	// completion is a percentage of the update completion (0-100)
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Completion int32 `json:"completion"`

	// startedAt is the time when the update started
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	StartedAt metav1.Time `json:"startedAt"`

	// completedAt is the time when the update completed
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`

	// estimatedCompletedAt is the estimated time when the update will complete
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	EstimatedCompletedAt *metav1.Time `json:"estimatedCompletedAt,omitempty"`
}

// ControlPlaneAssessment is the assessment of the control plane update process
// +kubebuilder:validation:Enum=Unknown;Progressing;Completed;Degraded
type ControlPlaneAssessment string

const (
	// Unknown means the update status and health cannot be determined
	ControlPlaneAssessmentUnknown ControlPlaneAssessment = "Unknown"
	// Progressing means the control plane is updating and no problems or slowness are detected
	ControlPlaneAssessmentProgressing ControlPlaneAssessment = "Progressing"
	// Completed means the control plane successfully completed updating and no problems are detected
	ControlPlaneAssessmentCompleted ControlPlaneAssessment = "Completed"
	// Degraded means the process of updating the control plane suffers from an observed problem
	ControlPlaneAssessmentDegraded ControlPlaneAssessment = "Degraded"
)

// ControlPlaneUpdateVersions contains the original and target versions of the upgrade
type ControlPlaneUpdateVersions struct {
	// previous is the version of the control plane before the update. When the cluster is being installed
	// for the first time, the version will have a placeholder value '<none>' and carry 'Installation' metadata
	// +required
	// +kubebuilder:validation:XValidation:rule="self.version == '<none>' ? (has(self.metadata) && self.metadata.exists(m, m.key == 'Installation')) : !(has(self.metadata) && self.metadata.exists(m, m.key == 'Installation'))",message="previous version must be '<none>' iff marked with Installation metadata"
	Previous Version `json:"previous"`

	// target is the version of the control plane after the update. It may never be '<none>' or have `Installation` metadata
	// +required
	// +kubebuilder:validation:XValidation:rule="self.version != '<none>' && !(has(self.metadata) && self.metadata.exists(m, m.key == 'Installation'))",message="target version must not be '<none>' or have Installation metadata"
	Target Version `json:"target"`
}

// Version describes a version involved in an update, typically on one side of an update edge
type Version struct {
	// version is a semantic version string, or a placeholder '<none>' for the special case where this
	// is a "previous" version in a new installation, in which case the metadata must contain an item
	// with key 'Installation'
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^(<none>)|((?:0|[1-9]\d*)[.](?:0|[1-9]\d*)[.](?:0|[1-9]\d*)(?:-(?:(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:[.](?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?)$`
	Version string `json:"version"`

	// metadata is a list of metadata associated with the version. It is a list of key-value pairs. The value is optional
	// and when not provided, the metadata item has boolean semantics (presence indicates true). For example, 'Partial'
	// metadata on a previous version indicates that the previous update was never fully completed. Can contain at most 5 items.
	// +listType=map
	// +listMapKey=key
	// +patchStrategy=merge
	// +patchMergeKey=key
	// +optional
	// +kubebuilder:validation:MaxItems=5
	Metadata []VersionMetadata `json:"metadata,omitempty" patchStrategy:"merge" patchMergeKey:"key"`
}

// VersionMetadata is a key:value item assigned to version involved in the update. Value can be empty, then the metadata
// have boolean semantics (true when present, false when absent)
type VersionMetadata struct {
	// key is the name of this metadata value. Valid values are: Installation, Partial, Architecture
	// +required
	Key VersionMetadataKey `json:"key"`

	// value is the value for the metadata, at most 32 characters long. It is not expected to be provided for Installation
	// and Partial metadata. For Architecture metadata, it is expected to be a string that indicates the architecture of the
	// payload image of the version involved in the upgrade, when relevant.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=32
	Value string `json:"value,omitempty"`
}

// VersionMetadataKey is a key for a metadata value associated with a version
// +kubebuilder:validation:Enum=Installation;Partial;Architecture
type VersionMetadataKey string

const (
	// Installation denotes a boolean that indicates the update was initiated as an installation
	InstallationMetadata VersionMetadataKey = "Installation"
	// Partial denotes a boolean that indicates the update was initiated in a state where the previous upgrade
	// (to the original version) was not fully completed
	PartialMetadata VersionMetadataKey = "Partial"
	// Architecture denotes a string that indicates the architecture of the payload image of the version,
	// when relevant
	ArchitectureMetadata VersionMetadataKey = "Architecture"
)

// ClusterVersionProgressInsightConditionType are types of conditions that can be reported on ClusterVersion status insight
type ClusterVersionProgressInsightConditionType string

const (
	// Updating condition communicates whether the ClusterVersion is updating
	ClusterVersionProgressInsightUpdating ClusterVersionProgressInsightConditionType = "Updating"
)

// ClusterVersionProgressInsightUpdatingReason are well-known reasons for the Updating condition on ClusterVersion status insights
type ClusterVersionProgressInsightUpdatingReason string

const (
	// CannotDetermineUpdating is used with Updating=Unknown
	ClusterVersionCannotDetermineUpdating ClusterVersionProgressInsightUpdatingReason = "CannotDetermineUpdating"
	// ClusterVersionProgressing means that ClusterVersion is considered to be Updating=True because it has a Progressing=True condition
	ClusterVersionProgressing ClusterVersionProgressInsightUpdatingReason = "ClusterVersionProgressing"
	// ClusterVersionNotProgressing means that ClusterVersion is considered to be Updating=False because it has a Progressing=False condition
	ClusterVersionNotProgressing ClusterVersionProgressInsightUpdatingReason = "ClusterVersionNotProgressing"
)

// ResourceRef is a reference to a kubernetes resource, typically involved in an insight
type ResourceRef struct {
	// group of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Group string `json:"group,omitempty"`

	// resource of object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Resource string `json:"resource"`

	// name of the object being referenced
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	Name string `json:"name"`

	// namespace of the object being referenced, if any
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Namespace string `json:"namespace,omitempty"`
}

// ScopeType is one of ControlPlane or WorkerPool
// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
type ScopeType string

const (
	// ControlPlane is used for insights that are related to the control plane (including control plane pool or nodes)
	ControlPlaneScope ScopeType = "ControlPlane"
	// WorkerPool is used for insights that are related to a worker pools and nodes (excluding control plane)
	WorkerPoolScope ScopeType = "WorkerPool"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterVersionProgressInsightList is a list of ClusterVersionProgressInsightList resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ClusterVersionProgressInsightList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is standard Kubernetes object metadata
	// +optional
	metav1.ListMeta `json:"metadata"`

	// items is a list of ClusterVersionProgressInsight resources
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	Items []ClusterVersionProgressInsight `json:"items"`
}
