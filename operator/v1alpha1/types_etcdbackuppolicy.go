package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// # EtcdBackupPolicy sets an automated schedule for taking backups of the etcd cluster
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=etcdbackuppolicies,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name=Storage Type,JSONPath=.spec.storage.type,type=string,description="Type of storage used for the backup"
// +kubebuilder:printcolumn:name=Schedule,JSONPath=.spec.schedule,type=string,description="Cron schedule for executing backups"
// +kubebuilder:printcolumn:name=Time Zone,JSONPath=.spec.timeZone,type=string,description="Time zone in which the schedule is evaluated"
// +kubebuilder:printcolumn:name=Last Schedule,JSONPath=.status.lastScheduleTime,type=date,description="Last time the schedule was executed"
// +kubebuilder:printcolumn:name=Age,JSONPath=.metadata.creationTimestamp,type=date,description="Age of the EtcdBackupPolicy"
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2952
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=etcd,operatorOrdering=01
// +openshift:enable:FeatureGate=AutomatedEtcdBackup
type EtcdBackupPolicy struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec EtcdBackupPolicySpec `json:"spec,omitzero"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status EtcdBackupPolicyStatus `json:"status,omitzero"`
}

type EtcdBackupPolicySpec struct {
	// schedule sets the backup schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +required
	Schedule string `json:"schedule,omitempty"`

	// timeZone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// If not specified, this will default to the time zone of the cluster-etcd-operator process.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// nodeSelector specifies which control plane nodes to select from for running backup jobs.
	// The default node-role.kubernetes.io/control-plane label will always be required in addition to any labels set here.
	// If no nodes are matched, then no EtcdBackups will be created.
	// For Local storage type, an EtcdBackup will be created for each selected control plane node every time the schedule is triggered. This is a special case to provide some resiliancy in the event of control plane node loss.
	// For PVC storage type, a single EtcdBackup will be created with the given nodeSelector every time the schedule is triggered.
	// When specified, nodeSelector must contain at least 1 entry and must not contain more than 10 entries.
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// storage specifies the location where etcd backup files will be saved.
	// +required
	Storage EtcdBackupStorage `json:"storage,omitzero"`

	// retentionRules defines the policy for retaining and deleting existing backups.
	// Backups are deleted from the oldest first until all rules are satisfied.
	// If no rules are specified then backups created by this policy will not be automatically deleted.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +listType=map
	// +listMapKey=type
	// +optional
	RetentionRules []EtcdBackupPolicyRetentionRule `json:"retentionRules,omitempty"`

	// failedBackupsHistoryLimit defined the number of failed etcdbackups to retain. Value must be non-negative integer. Defaults to 1.
	// +kubebuilder:validation:Minimum=0
	// +optional
	FailedBackupsHistoryLimit *int32 `json:"failedBackupsHistoryLimit,omitempty"`
}

// +union
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxQuantity') ? has(self.maxQuantity) : !has(self.maxQuantity)",message="maxQuantity is required when type is MaxQuantity, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxSize') ? has(self.maxSize) : !has(self.maxSize)",message="maxSize is required when type is MaxSize, and forbidden otherwise"
type EtcdBackupPolicyRetentionRule struct {
	// type defined which rule field is set
	// +unionDiscriminator
	// +kubebuilder:validation:Enum:=MaxQuantity;MaxSize
	// +required
	Type EtcdBackupPolicyRetentionRuleType `json:"type,omitempty"`

	// maxQuantity enforces the deletion of backups that exceed the given count.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MaxQuantity int32 `json:"maxQuantity,omitempty"`

	// maxSize enforces the deletion of backups by the total size of backups on the storage backend.
	// This is a soft threshold. The total size of backups may temporarily exceed the limit when new backups are created.
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:XValidation:rule="quantity(self).isGreaterThan(quantity('0'))",message="request must be a positive, non-zero quantity"
	// +optional
	MaxSize *resource.Quantity `json:"maxSize,omitempty"`
}

type EtcdBackupPolicyRetentionRuleType string

const (
	EtcdBackupPolicyRetentionRuleMaxQuantity EtcdBackupPolicyRetentionRuleType = "MaxQuantity"
	EtcdBackupPolicyRetentionRuleMaxSize     EtcdBackupPolicyRetentionRuleType = "MaxSize"
)

// +kubebuilder:validation:MinProperties=1
type EtcdBackupPolicyStatus struct {
	// conditions provide details on the status of the etcd backup policy.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// active is a list of references to in progress backups controlled by this policy
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=20
	// +listType=map
	// +listMapKey=name
	// +optional
	Active []EtcdBackupReference `json:"active,omitempty"`

	// lastScheduleTime is the time when the last scheduled backup was triggered.
	// This is used by the controller to track when backups have been executed
	// and to prevent duplicate executions on controller restart.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

type EtcdBackupReference struct {
	// name of the backup
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	Name string `json:"name,omitempty"`
	// uid of the backup
	// +kubebuilder:validation:MinLength=36
	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:Format=uuid
	// +required
	UID string `json:"uid,omitempty"`
}

type EtcdBackupPolicyConditionType string

const (
	// BackupPolicyError indicates that something is invalid on the EtcdBackupPolicy.
	BackupPolicyError EtcdBackupPolicyConditionType = "Error"
)

type EtcdBackupPolicyConditionReason string

const (
	// BackupPolicyInvalidSchedule is set when parsing the schedule or timezone fails.
	BackupPolicyInvalidSchedule EtcdBackupPolicyConditionReason = "InvalidSchedule"
	// BackupPolicyInvalidSelector is set when no control plane nodes are selected by the nodeSelector.
	BackupPolicyInvalidSelector EtcdBackupPolicyConditionReason = "InvalidNodeSelector"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdBackupPolicyList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type EtcdBackupPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EtcdBackupPolicy `json:"items"`
}
