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
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:XValidation:rule="self.matches('^(@(annually|yearly|monthly|weekly|daily|hourly)|@every (\d+(\.\d+)?(s|m|h))+|(([-0-9*/,A-Za-z]+ ){4}[-0-9*/,A-Za-z]+))$')",message="must be a valid cron expression, or macro"
	// +required
	Schedule string `json:"schedule,omitempty"`

	// nodeSelector specifies which control plane nodes to select from for running backup jobs.
	// The default node-role.kubernetes.io/control-plane label will always be required in addition to any labels set here.
	// If no nodes are matched, then no EtcdBackups will be created.
	// For Local storage type, an EtcdBackup will be created for each selected control plane node every time the schedule is triggered. This is a special case to provide some resiliancy in the event of control plane node loss.
	// For PVC storage type, a single EtcdBackup will be created with the given nodeSelector every time the schedule is triggered.
	// When specified, must contain between 1 and 10 entries.
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
	// When an EtcdBackup is deleted the files created by it will be deleted as well, as long as the storage backend is still accessible.
	// When specified, must contain between 1 and 2 entries.
	// Must not have more than 1 entry of each rule type.
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
	// Must be a positive integer no greater than 1000.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	// +optional
	FailedBackupsHistoryLimit *int32 `json:"failedBackupsHistoryLimit,omitempty"`
}

// +union
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxQuantity') ? has(self.maxQuantity) : !has(self.maxQuantity)",message="maxQuantity is required when type is MaxQuantity, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxTotalSizeGb') ? has(self.maxTotalSizeGb) : !has(self.maxTotalSizeGb)",message="maxTotalSizeGb is required when type is MaxTotalSizeGb, and forbidden otherwise"
type EtcdBackupPolicyRetentionRule struct {
	// type defined which rule field is set
	// Allowed values are MaxQuantity and MaxTotalSizeGb.
	// When set to MaxQuantity, this rule deletes old EtcdBackups when the quantity is exceeded.
	// When set to MaxTotalSizeGb, this rule deletes old EtcdBackups when the total aggregate size of backup files exceeds the value.
	// +unionDiscriminator
	// +kubebuilder:validation:Enum:=MaxQuantity;MaxTotalSizeGb
	// +required
	Type EtcdBackupPolicyRetentionRuleType `json:"type,omitempty"`

	// maxQuantity enforces the deletion of backups that exceed the given count.
	// This field is required when the rule type is "MaxQuantity", and forbidden otherwise.
	// Must be between 1 and 1000.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +optional
	MaxQuantity int32 `json:"maxQuantity,omitempty"`

	// maxTotalSizeGb enforces the deletion of backups by the total size of backups on the storage backend (in GB).
	// This is a soft threshold. The total size of backups may temporarily exceed the limit when new backups are created.
	// This field is required when the rule type is "MaxTotalSizeGb", and forbidden otherwise.
	// Must be between 1 and 10000.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	// +optional
	MaxTotalSizeGb int32 `json:"maxTotalSizeGb,omitempty"`
}

type EtcdBackupPolicyRetentionRuleType string

const (
	EtcdBackupPolicyRetentionRuleMaxQuantity    EtcdBackupPolicyRetentionRuleType = "MaxQuantity"
	EtcdBackupPolicyRetentionRuleMaxTotalSizeGb EtcdBackupPolicyRetentionRuleType = "MaxTotalSizeGb"
)

// +kubebuilder:validation:MinProperties=1
type EtcdBackupPolicyStatus struct {
	// conditions provide details on the status of the etcd backup policy.
	// When specified, must contain between 1 and 3 entries, with no more than 1 of each type.
	//
	// Supported conditions include Error, LastScheduled, and LastSuccessfulTime.
	// If Error is true, the EtcdBackupPolicy is invalid and the condition reason and message will include details on the error.
	// If LastScheduled is true, the most recently scheduled EtcdBackup was created at the condition's lastTransitionTime.
	// If LastSuccessfulTime is true, the most recently completed EtcdBackup finished at the condition's lastTransitionTime.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=3
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// active is a list of references to in progress backups controlled by this policy
	// When specified, must contain between 1 and 20 entries.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=20
	// +listType=map
	// +listMapKey=name
	// +optional
	Active []EtcdBackupReference `json:"active,omitempty"`
}

type EtcdBackupReference struct {
	// name of the backup
	// Must be between 1 and 253 characters and conform to RFC 1123 subdomain format:
	// lowercase alphanumeric characters, '-' or '.', starting and ending with alphanumeric characters.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="name must be a lowercase RFC 1123 subdomain consisting of lowercase alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character"
	// +required
	Name string `json:"name,omitempty"`
	// uid of the backup.
	// Must be a valid UUID formatted string with 36 characters.
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
	// BackupPolicyLastScheduled is set true with the timestamp that the last scheduled backup was triggered.
	BackupPolicyLastScheduled EtcdBackupPolicyConditionType = "LastScheduled"
	// BackupPolicyLastSuccessfulTime is set true with the timestmap that the last scheduled backup completed successfully.
	BackupPolicyLastSuccessfulTime EtcdBackupPolicyConditionType = "LastSuccessfulTime"
)

type EtcdBackupPolicyConditionReason string

const (
	// BackupPolicyInvalidSchedule is set when parsing the schedule fails.
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
