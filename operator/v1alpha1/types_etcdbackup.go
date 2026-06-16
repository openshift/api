package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// # EtcdBackup provides configuration options and status for a one-time backup attempt of the etcd cluster
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=etcdbackups,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name=Storage Type,JSONPath=.spec.storage.type,type=string,description="Type of storage used for the backup"
// +kubebuilder:printcolumn:name=Status,JSONPath=.status.conditions[?(@.status=="True")].type,type=string,description="Status of the EtcdBackup"
// +kubebuilder:printcolumn:name=Age,JSONPath=.metadata.creationTimestamp,type=date,description="Age of the EtcdBackup"
// +kubebuilder:printcolumn:name=Node,JSONPath=.status.nodeName,type=string,description="Name of the node where the backup was executed",priority=1
// +kubebuilder:printcolumn:name=PVC,JSONPath=.spec.storage.pvc.name,type=string,description="Name of the PVC where the backup is stored",priority=1
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2952
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=etcd,operatorOrdering=01
// +openshift:enable:FeatureGate=AutomatedEtcdBackup
type EtcdBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec EtcdBackupSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status EtcdBackupStatus `json:"status"`
}

type EtcdBackupSpec struct {
	// nodeSelector specifies which master node(s) to run the backup job on.
	// If no selector is specified, the default node-role.kubernetes.io/control-plane label will be used.
	// If no nodes are matched, then no backup will run.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="nodeSelector is immutable once set"
	// +kubebuilder:validation:Optional
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// storage specifies the location where etcd backup files will be saved.
	// +kubebuilder:validation:Required
	// +required
	Storage EtcdBackupStorage `json:"storage"`
}

// +kubebuilder:validation:XValidation:rule="self.type == 'PVC' ? has(self.pvc) : !has(self.pvc)",message="pvc is required when type is PVC, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Local' ? has(self.local) : !has(self.local)",message="local is required when type is Local, and forbidden otherwise"
// +union
type EtcdBackupStorage struct {
	// +kubebuilder:validation:Enum=PVC;Local
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="type is immutable once set"
	// +required
	// +unionDiscriminator
	Type EtcdBackupStorageType `json:"type"`

	// pvc specifies the PersistentVolumeClaim (PVC) which binds a PersistentVolume where the etcd backup file will be saved.
	// The PVC must always be created in the "openshift-etcd" namespace.
	// This field is required when the storage type is "PVC"
	// +kubebuilder:validation:Optional
	// +optional
	// +unionMember
	PVC *EtcdBackupStoragePvc `json:"pvc,omitempty"`

	// local specifies a host path directory on the master node where the etcd backup file will be saved.
	// This field is required when storage type is "Local"
	// +kubebuilder:validation:Optional
	// +optional
	// +unionMember
	Local *EtcdBackupStorageLocal `json:"local,omitempty"`
}

// EtcdBackupStorageType is an enum of the supported storage backends for backup files
type EtcdBackupStorageType string

const (
	EtcdBackupStorageTypePVC   EtcdBackupStorageType = "PVC"
	EtcdBackupStorageTypeLocal EtcdBackupStorageType = "Local"
)

type EtcdBackupStoragePvc struct {
	// name is a reference to a PVC in the "openshift-etcd" namespace where the etcd backup file will be saved.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="name is immutable once set"
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// path is a directory on the volume where the etcd backup file will be saved.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="path is immutable once set"
	// +kubebuilder:validation:Optional
	// +optional
	Path string `json:"path,omitempty"`
}

type EtcdBackupStorageLocal struct {
	// hostPath is a local directory on the master node where the etcd backup file will be saved.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="hostPath is immutable once set"
	// +kubebuilder:validation:Required
	// +required
	HostPath string `json:"hostPath"`
}

type EtcdBackupStatus struct {
	// conditions provide details on the status of the etcd backup job.
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:Optional
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// job is a reference to the Job created for the backup.
	// +kubebuilder:validation:Optional
	// +optional
	Job *EtcdBackupJobReference `json:"job,omitempty"`

	// nodeName is used for Local backups to bind the control plane node where the snapshot will be taken.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Optional
	// +optional
	NodeName string `json:"nodeName,omitempty"`

	// files tracks the path and size of files generated by the etcd backup.
	// Includes both etcd snapshots and static manifests.
	// +listType=map
	// +listMapKey=path
	// +kubebuilder:validation:Optional
	// +optional
	Files []EtcdBackupFile `json:"files,omitempty"`
}

type EtcdBackupJobReference struct {
	// name of the backup job
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
	// namespace of the backup job
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Required
	// +required
	Namespace string `json:"namespace"`
	// uid of the backup job
	// +kubebuilder:validation:MinLength=36
	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:Format=uuid
	// +kubebuilder:validation:Required
	// +required
	UID string `json:"uid"`
}

type EtcdBackupFile struct {
	// path to the backup file on the storage backend.
	// +kubebuilder:validation:Required
	// +required
	Path string `json:"path"`

	// sizeBytes is the size of the backup file on the storage backend in bytes.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	// +required
	SizeBytes int64 `json:"sizeBytes"`
}

// BackupConditionType enumerates the Condition types added to EtcdBackupStatus at different points in its lifecycle
type BackupConditionType string

var (
	// BackupPending means the backup is ready to start.
	BackupPending BackupConditionType = "Pending"
	// BackupPending means the backup job has started.
	BackupRunning BackupConditionType = "Running"
	// BackupCompleted means the backup completed successfully.
	BackupCompleted BackupConditionType = "Completed"
	// BackupFailed means the backup failed.
	BackupFailed BackupConditionType = "Failed"
	// BackupGarbageCollectionRequired indicates whether or not garbage collection is required
	// on a failed backup to cleanup partially created files.
	BackupGarbageCollectionRequired BackupConditionType = "GarbageCollectionRequired"
)

// BackupConditionReason enumerates the Condition reasons associated with BackupConditionTypes
type BackupConditionReason string

var (
	// BackupReasonReadyToStart means the backup has been queued to start.
	BackupReasonReadyToStart BackupConditionReason = "ReadyToStart"

	// BackupReasonJobStarted means the backup job is currently running.
	BackupReasonJobStarted BackupConditionReason = "JobStarted"

	// BackupReasonJobCompleted means the backup job completed successfully.
	BackupReasonJobCompleted BackupConditionReason = "JobCompleted"

	// BackupReasonPVCNotFound means the backup failed due to a missing PVC.
	BackupReasonPVCNotFound BackupConditionReason = "PVCNotFound"
	// BackupReasonNotNotFound means the backup failed due to a missing node.
	BackupReasonNodeNotFound BackupConditionReason = "NodeNotFound"
	// BackupReasonJobFailed means the backup job failed.
	BackupReasonJobFailed BackupConditionReason = "JobFailed"
	// BackupReasonDeleted means the backup was deleted while in progress.
	BackupReasonDeleted BackupConditionReason = "Deleted"

	// BackupReasonFilesPartiallyCreated means the backup job failed after partially creating some files.
	BackupReasonFilesPartiallyCreated BackupConditionReason = "FilesPartiallyCreated"
	// BackupReasonFilesNotCreated means the backup job failed without creating any files.
	BackupReasonFilesNotCreated BackupConditionReason = "FilesNotCreated"
	// BackupReasonFileStateUnknown means the backup job failed without indicating if it had created any files.
	BackupReasonFileStateUnknown BackupConditionReason = "FileStateUnknown"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdBackupList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type EtcdBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EtcdBackup `json:"items"`
}
