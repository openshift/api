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
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec EtcdBackupPolicySpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status EtcdBackupPolicyStatus `json:"status"`
}

type EtcdBackupPolicySpec struct {
	// schedule sets the backup schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:Required
	// +required
	Schedule string `json:"schedule"`

	// timeZone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// If not specified, this will default to the time zone of the cluster-etcd-operator process.
	// +kubebuilder:validation:Optional
	// +optional
	TimeZone string `json:"timeZone,omitempty"`

	// nodeSelector specifies which master node(s) to run backup jobs on.
	// If no selector is specified, the default node-role.kubernetes.io/master label will be used.
	// If no nodes are matched, then no backups will run.
	// +kubebuilder:validation:Optional
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// storage specifies the location where etcd backup files will be saved.
	// +kubebuilder:validation:Required
	// +required
	Storage EtcdBackupStorage `json:"storage"`

	// retentionRules defines the policy for retaining and deleting existing backups.
	// Backups are deleted from the oldest first until all rules are satisfied.
	// If no rules are specified then backups created by this policy will not be automatically deleted.
	// +kubebuilder:validation:Optional
	// +optional
	RetentionRules []EtcdBackupPolicyRetentionRule `json:"retentionRules,omitzero"`

	// failedBackupsHistoryLimit defined the number of failed etcdbackups to retain. Value must be non-negative integer. Defaults to 1.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Default=1
	// +kubebuilder:validation:Optional
	// +optional
	FailedBackupsHistoryLimit int `json:"failedBackupsHistoryLimit"`
}

// +union
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxQuantity') ? has(self.maxQuantity) : !has(self.maxQuantity)",message="maxQuantity is required when type is MaxQuantity, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="(self.type == 'MaxSize') ? has(self.maxSize) : !has(self.maxSize)",message="maxSize is required when type is MaxSize, and forbidden otherwise"
type EtcdBackupPolicyRetentionRule struct {
	// type defined which rule field is set
	// +unionDiscriminator
	// +kubebuilder:validation:Enum:=MaxQuantity;MaxSize
	// +kubebuilder:validation:Required
	// +required
	Type EtcdBackupPolicyRetentionRuleType `json:"type"`

	// maxQuantity enforces the deletion of backups that exceed the given count.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Optional
	// +optional
	MaxQuantity int `json:"maxQuantity,omitzero"`

	// maxSize enforces the deletion of backups by the total size of backups on the storage backend.
	// This is a soft threshold. The total size of backups may temporarily exceed the limit when new backups are created.
	// +kubebuilder:validation:Optional
	// +optional
	MaxSize resource.Quantity `json:"maxSize,omitzero"`
}

type EtcdBackupPolicyRetentionRuleType string

const (
	EtcdBackupPolicyRetentionRuleMaxQuantity EtcdBackupPolicyRetentionRuleType = "MaxQuantity"
	EtcdBackupPolicyRetentionRuleMaxSize     EtcdBackupPolicyRetentionRuleType = "MaxSize"
)

type EtcdBackupPolicyStatus struct {
	// active is a list of references to in progress backups controlled by this policy
	// +kubebuilder:validation:Optional
	// +optional
	Active []EtcdBackupReference `json:"active,omitempty"`

	// lastScheduleTime is the time when the last scheduled backup was triggered.
	// This is used by the controller to track when backups have been executed
	// and to prevent duplicate executions on controller restart.
	// +kubebuilder:validation:Optional
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

type EtcdBackupReference struct {
	// name of the backup
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
	// uid of the backup
	// +kubebuilder:validation:Required
	// +required
	UID string `json:"uid"`
}

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
