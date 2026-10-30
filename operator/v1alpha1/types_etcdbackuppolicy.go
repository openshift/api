package v1alpha1

import (
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
// +kubebuilder:printcolumn:name=TimeZone,JSONPath=.spec.timeZone,type=string,description="Cron time zone for executing backups"
// +kubebuilder:printcolumn:name=Status,JSONPath=.status.conditions[?(@.type=="Scheduled")].reason,type=string,description="Status of the EtcdBackupPolicy"
// +kubebuilder:printcolumn:name=Age,JSONPath=.metadata.creationTimestamp,type=date,description="Age of the EtcdBackupPolicy"
// +kubebuilder:validation:XValidation:rule="self.metadata.name.size() <= 52",message="name cannot exceed 52 characters"
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
	// The value must be one of the following forms:
	// a predefined macro "@yearly", "@annually", "@monthly", "@weekly", "@daily" or "@hourly";
	// an "@every <duration>" expression where duration is a sequence of decimal numbers each with
	// a unit of "ns", "us", "ms", "s", "m" or "h" (for example "@every 2h" or "@every 1h30m");
	// or a standard cron expression of five space separated fields (minute, hour, day-of-month,
	// month and day-of-week), where each field may only contain digits, the names used for months
	// and days of the week, and the cron special characters "*", "/", "-" and ",".
	// This validation only checks the overall shape of the value; the individual field values are
	// validated when the schedule is parsed by the controller, which reports parsing failures via
	// the schedule's status condition.
	// Must be between 1 and 1024 characters long.
	// This field is required.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:XValidation:rule="self.matches('^(@(annually|yearly|monthly|weekly|daily|hourly)|@every (\d+(\.\d+)?(s|m|h))+|(([-0-9*/,A-Za-z]+ ){4}[-0-9*/,A-Za-z]+))$')",message="must be a valid cron expression, or macro"
	// +required
	Schedule string `json:"schedule,omitempty"`

	// timeZone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// If not specified, this will default to the time zone of the cluster-etcd-operator process.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// assets specifies the type of files to be backed up.
	// If Snapshot is present, the etcd database will be snapshotted and saved to the storage backend.
	// If StaticKubeResources is present, static pod manifests for the control plane will be bundled and saved to the storage backend. If etcd is encrypted, this will include the encryption keys.
	// When specified, must contain between 1 and 2 entries, with at most 1 of each type.
	// When omitted, all assets types will be backed up.
	// This field is optional.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +listType=set
	// +optional
	Assets []EtcdBackupAssetType `json:"assets,omitempty"`

	// nodeSelector specifies which control plane nodes to select from for running backup jobs.
	// The default node-role.kubernetes.io/control-plane label will always be required in addition to any labels set here.
	// If no nodes are matched, then no EtcdBackups will be created.
	// When specified, must contain between 1 and 10 entries.
	// This field is optional.
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// storage specifies the location where etcd backup files will be saved.
	// This field is required.
	// +required
	Storage EtcdBackupPolicyStorage `json:"storage,omitzero"`

	// retentionRules defines the policy for retaining and deleting existing backups.
	// Backups are deleted from the oldest first until all rules are satisfied.
	// If no rules are specified then backups created by this policy will not be automatically deleted.
	// When an EtcdBackup is deleted the files created by it will be deleted as well, as long as the storage backend is still accessible.
	// When specified, must contain between 1 and 2 entries.
	// Must not have more than 1 entry of each rule type.
	// This field is optional.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=2
	// +listType=map
	// +listMapKey=type
	// +optional
	RetentionRules []EtcdBackupPolicyRetentionRule `json:"retentionRules,omitempty"`

	// failedBackupsHistoryLimit defined the number of failed etcdbackups to retain. Failed backups are deleted from the oldest first.
	// Must be non-negative integer.
	// If set to to 0, then failed backups will be deleted immediately.
	// If unset, defaults to 1.
	// Must be a positive integer no greater than 100.
	// This field is optional.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	FailedBackupsHistoryLimit *int32 `json:"failedBackupsHistoryLimit,omitempty"`
}

// +kubebuilder:validation:XValidation:rule="self.type == 'PVC' ? has(self.pvc) : !has(self.pvc)",message="pvc is required when type is PVC, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Local' ? has(self.local) : !has(self.local)",message="local is required when type is Local, and forbidden otherwise"
// +union
type EtcdBackupPolicyStorage struct {
	// type of storage backend to use for storing the etcd backup.
	// Allowed values are Local and PVC.
	// When set to Local, the backup will be saved to a host path directory on the control plane node it was taken on.
	// When set to PVC, the backup will be saved to a volume bound by a PersistentVolumeClaim.
	// This field is required.
	// +kubebuilder:validation:Enum=PVC;Local
	// +required
	// +unionDiscriminator
	Type EtcdBackupStorageType `json:"type,omitempty"`

	// pvc specifies the PersistentVolumeClaim (PVC) which binds a PersistentVolume where the etcd backup file will be saved.
	// The PVC must always be created in the "openshift-etcd" namespace.
	// This field is required when the storage type is "PVC", and forbidden otherwise.
	// +optional
	// +unionMember
	PVC EtcdBackupPolicyStoragePvc `json:"pvc,omitzero"`

	// local specifies a host path directory on the control plane node where the etcd backup file will be saved.
	// This field is required when storage type is "Local", and forbidden otherwise.
	// +optional
	// +unionMember
	Local EtcdBackupPolicyStorageLocal `json:"local,omitzero"`
}

type EtcdBackupPolicyStoragePvc struct {
	// +required
	EtcdBackupStoragePvc `json:",inline"`
}

type EtcdBackupPolicyStorageLocal struct {
	// +required
	EtcdBackupStorageLocal `json:",inline"`

	// nodeCount enables creating multiple local backups per schedule distributed across control plane nodes.
	// Control plane nodes will be iterated in lexicographical order when backups are created, with preference for nodes that already have backups for this policy.
	// The number of created EtcdBackups may be less than nodeCount if there are fewer control plane nodes selected.
	// When specified, must be an integer between 1 and 10.
	// If omitted, defaults to 1.
	// This field is optional.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +optional
	NodeCount int32 `json:"nodeCount,omitempty"`
}

// +union
// +kubebuilder:validation:XValidation:rule="(self.type == 'Quantity') ? has(self.quantity) : !has(self.quantity)",message="quantity is required when type is Quantity, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="(self.type == 'Size') ? has(self.size) : !has(self.size)",message="size is required when type is Size, and forbidden otherwise"
type EtcdBackupPolicyRetentionRule struct {
	// type defined which rule field is set
	// Allowed values are Quantity and Size.
	// When set to Quantity, this rule deletes old EtcdBackups based on the total count.
	// When set to Size, this rule deletes old EtcdBackups based on the total aggregate size of backup files.
	// This field is required.
	// +unionDiscriminator
	// +kubebuilder:validation:Enum:=Quantity;Size
	// +required
	Type EtcdBackupPolicyRetentionRuleType `json:"type,omitempty"`

	// quantity enforces the deletion of EtcdBackups belonging to the EtcdBackupPolicy based on total count.
	// This field is required when the rule type is "Quantity", and forbidden otherwise.
	// +optional
	Quantity EtcdBackupPolicyRetentionQuantity `json:"quantity,omitzero"`

	// size enforces the deletion of EtcdBackups belonging to the EtcdBackupPolicy based on total backup file size.
	// This field is required when the rule type is "Size", and forbidden otherwise.
	// +optional
	Size EtcdBackupPolicyRetentionSize `json:"size,omitzero"`
}

type EtcdBackupPolicyRetentionRuleType string

const (
	EtcdBackupPolicyRetentionRuleTypeQuantity EtcdBackupPolicyRetentionRuleType = "Quantity"
	EtcdBackupPolicyRetentionRuleTypeSize     EtcdBackupPolicyRetentionRuleType = "Size"
)

type EtcdBackupPolicyRetentionQuantity struct {
	// maxQuantity enforces the deletion of backups that exceed the total count.
	// Must be between 1 and 1000.
	// Thius field is required.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +required
	MaxQuantity int32 `json:"maxQuantity,omitempty"`
}

type EtcdBackupPolicyRetentionSize struct {
	// maxTotalSizeGiB enforces the deletion of backups by the total size of backups on the storage backend (in GiB).
	// This is a soft threshold. The total size of backups may temporarily exceed the limit when new backups are created.
	// Must be between 1 and 10240 (10 TiB).
	// This field is required.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10240
	// +required
	MaxTotalSizeGiB int32 `json:"maxTotalSizeGiB,omitempty"`
}

// +kubebuilder:validation:MinProperties=1
type EtcdBackupPolicyStatus struct {
	// conditions provide details on the status of the etcd backup policy.
	// Supported conditions include Scheduled.
	// If Scheduled is true, scheduling the policy succeeded.
	// If Scheduled is false, an error occurred while attempting to schedule the policy. The condition reason and message will provide addition details.
	//
	// When specified, must contain between 1 and 10 entries, with no more than 1 of each type.
	// This field is optional.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type EtcdBackupPolicyConditionType string

const (
	// BackupPolicyScheduled is True if the EtcdBackupPolicy was successfully scheduled, or False if there is an error.
	BackupPolicyScheduled EtcdBackupPolicyConditionType = "Scheduled"
)

type EtcdBackupPolicyConditionReason string

const (
	// BackupPolicyInvalidSchedule is set when parsing the schedule or timeZone fails.
	BackupPolicyInvalidSchedule EtcdBackupPolicyConditionReason = "InvalidSchedule"
	// BackupPolicySchedulingFailed is set when there was an error applying a valid schedule.
	BackupPolicySchedulingFailed EtcdBackupPolicyConditionReason = "SchedulingFailed"
	// BackupPolicySchedulingSucceeded is set when the policy has been successfully scheduled.
	BackupPolicySchedulingSucceeded EtcdBackupPolicyConditionReason = "SchedulingSucceeded"
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
