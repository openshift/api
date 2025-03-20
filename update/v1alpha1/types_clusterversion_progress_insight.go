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
// +kubebuilder:validation:XValidation:rule="!has(self.status) || self.status.name == self.metadata.name",message="When status is present, .status must match .metadata.name"
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
	// +kubebuilder:validation:MaxItems=5
	// +TODO: Add validations to enforce all known conditions are present (CEL+MinItems), once conditions stabilize
	// +TODO: Add validations to enforce that only known Reasons are used in conditions, once conditions stabilize
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// name is equal to the name of the corresponding clusterversions.config.openshift.io resource, typically 'version'
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[a-z0-9-]+$`
	Name string `json:"name"`

	// assessment is a brief summary assessment of the control plane update process. This value is human-oriented, and
	// while it looks like a state/phase enum, it is not meant to be used as such. Assessment is meant as human-oriented
	// brief summary matching the state expressed in conditions (taking into account various relations between them, like
	// ordering or precedence), intended to be directly used in UIs and reports. For machine-oriented conditional behavior
	// depending on the state, the conditions should be used instead.
	//
	// The known values are: Unknown, Progressing, Completed, Degraded. The API is not restricted to these values, and
	// valid values can be even brief phrases, up to 64 characters long.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +required
	Assessment ClusterVersionAssessment `json:"assessment"`

	// versions contains the original and target versions involved in the update, either the ongoing one or the last update
	// completed.
	// +required
	Versions ControlPlaneUpdateVersions `json:"versions"`

	// completionPercent conveys the update completion (0-100). When there is no update in progress, the ClusterVersion
	// Progress Insight represents the last update (or installation, which is considered to be an update to the initial
	// version) that is by definition completed, and therefore the completionPercent is 100.
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Completion int32 `json:"completionPercent"`

	// startedAt is the time when the update started. When there is no update in progress, the Cluster Version
	// Progress Insight represents the last update (or installation, which is considered to be an update to the initial version)
	// that is by definition completed, and this field represents the time when that update was initiated.
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	StartedAt metav1.Time `json:"startedAt"`

	// completedAt is the time when the update completed. This field is only set when the update is completed.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`

	// estimatedCompletedAt is the estimated time when the update will complete, if such estimate is available. When there
	// is no update in progress, this field is either not set at all, or its value is the time when the last update was
	// expected to complete.
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	EstimatedCompletedAt *metav1.Time `json:"estimatedCompletedAt,omitempty"`
}

// ControlPlaneAssessment convey assessments of the control plane update process. These values are human-oriented, and
// while they look like state enums, they are not meant to be used as such. They are meant as human-oriented brief
// summaries that can be directly used in UIs and reports. For machine-oriented conditional behavior depending on the
// state, the conditions should be used instead. The known values are: Unknown, Progressing, Completed, Degraded but
// the API is not restricted to these values, and valid values could be even more than one word.
type ClusterVersionAssessment string

const (
	// Unknown means the update status and health cannot be determined
	ClusterVersionAssessmentUnknown ClusterVersionAssessment = "Unknown"
	// Progressing means the control plane is updating and no problems or slowness are detected
	ClusterVersionAssessmentProgressing ClusterVersionAssessment = "Progressing"
	// Completed means the control plane successfully completed updating and no problems are detected
	ClusterVersionAssessmentCompleted ClusterVersionAssessment = "Completed"
	// Degraded means the process of updating the control plane suffers from an observed problem
	ClusterVersionAssessmentDegraded ClusterVersionAssessment = "Degraded"
)

// ControlPlaneUpdateVersions contains the original and target versions of the upgrade
// +kubebuilder:validation:XValidation:rule="has(self.previous) || (has(self.target.metadata) && self.target.metadata.exists(m, m.key == 'Installation'))",message="target version must have 'Installation' metadata when previous version is empty"
// +kubebuilder:validation:XValidation:rule="!(has(self.previous) && has(self.target.metadata) && self.target.metadata.exists(m, m.key == 'Installation'))",message="target version can only have 'Installation' metadata when previous version is empty"
type ControlPlaneUpdateVersions struct {
	// previous is the desired version of the control plane the before the update, regardless of completion. When empty,
	// it means the cluster was never updated yet, and the target version is the initial version of the cluster. When the
	// current update was triggered in the state where the previous update was not fully completed, the version will carry
	// 'Partial' metadata.
	// +optional
	Previous *Version `json:"previous,omitempty"`

	// target is the version of the control plane after the update. If the cluster was never updated, the version will carry
	// 'Installation' metadata.
	// +required
	// +kubebuilder:validation:XValidation:rule="has(self.version) && size(self.version) > 0",message="Target version must be set and not empty"
	Target Version `json:"target"`
}

// Version describes a version involved in an update, typically on one side of an update edge
type Version struct {
	// version is a semantic version string, or a placeholder '<none>' for the special case where this
	// is a "previous" version in a new installation, in which case the metadata must contain an item
	// with key 'Installation'
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^(?:0|[1-9]\d*)[.](?:0|[1-9]\d*)[.](?:0|[1-9]\d*)(?:-(?:(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:[.](?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?$`
	Version string `json:"version,omitempty"`

	// metadata is a list of metadata associated with the version. It is a list of key-value pairs. The value is optional
	// and when not provided, the metadata item has boolean semantics (presence indicates true). For example, 'Partial'
	// metadata on a previous version indicates that the previous update was never fully completed. Can contain at most 5 items.
	// +listType=map
	// +listMapKey=key
	// +patchStrategy=merge
	// +patchMergeKey=key
	// +optional
	// +kubebuilder:validation:MaxItems=2
	Metadata []VersionMetadata `json:"metadata,omitempty" patchStrategy:"merge" patchMergeKey:"key"`
}

// VersionMetadata is a key:value item assigned to version involved in the update. Value can be empty, then the metadata
// have boolean semantics (true when present, false when absent)
type VersionMetadata struct {
	// key is the name of this metadata value. Valid values are:
	//   Installation (boolean): indicates the target version is also initial version of the cluster
	//   Partial (boolean): indicates the previous update was not fully completed
	//   Architecture: a string that indicates the architecture of the payload image of the version involved in the upgrade, present only when relevant
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

// ClusterVersionProgressInsightConditionType are types of conditions that can be reported on ClusterVersion progress insight
type ClusterVersionProgressInsightConditionType string

const (
	// Updating condition communicates whether the control plane, represented by this ClusterVersion, is updating
	ClusterVersionProgressInsightUpdating ClusterVersionProgressInsightConditionType = "Updating"
	// Healthy condition communicates whether the control plane update process is observed to be healthy
	ClusterVersionProgressInsightHealthy ClusterVersionProgressInsightConditionType = "Healthy"
)

// ClusterVersionProgressInsightUpdatingReason are well-known reasons for the Updating condition on ClusterVersion progress insights
type ClusterVersionProgressInsightUpdatingReason string

const (
	// CannotDetermineUpdating is used with Updating=Unknown
	ClusterVersionCannotDetermineUpdating ClusterVersionProgressInsightUpdatingReason = "CannotDetermineUpdating"
	// ClusterVersionProgressing means that ClusterVersion is considered to be Updating=True because it has a Progressing=True condition
	ClusterVersionProgressing ClusterVersionProgressInsightUpdatingReason = "Progressing"
	// ClusterVersionNotProgressing means that ClusterVersion is considered to be Updating=False because it has a Progressing=False condition
	ClusterVersionNotProgressing ClusterVersionProgressInsightUpdatingReason = "NotProgressing"
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
