package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterCSIDriver is used to manage and configure CSI driver installed by default
// in OpenShift. An example configuration may look like:
//   apiVersion: operator.openshift.io/v1
//   kind: "ClusterCSIDriver"
//   metadata:
//     name: "ebs.csi.aws.com"
//   spec:
//     logLevel: Debug

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clustercsidrivers,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:validation:FeatureGateAwareXValidation:featureGate=CinderCSIDirectAttach,rule="self.spec.?driverConfig.?driverType.orValue('') == 'Cinder' ? self.metadata.name == 'cinder.csi.openstack.org' : true",message="driverType Cinder requires metadata.name cinder.csi.openstack.org"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=CinderCSIDirectAttach,rule="self.metadata.name == 'cinder.csi.openstack.org' ? self.spec.?driverConfig.?driverType.orValue('') in ['', 'Cinder'] : true",message="the Cinder resource permits only omitted or Cinder driver configuration"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=CinderCSIDirectAttach,rule="self.?status.?cinder.hasValue() ? self.metadata.name == 'cinder.csi.openstack.org' : true",message="Cinder status is valid only on the Cinder resource"
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/701
// +openshift:file-pattern=cvoRunLevel=0000_50,operatorName=csi-driver,operatorOrdering=01
// +kubebuilder:validation:XValidation:rule="self.spec.?driverConfig.driverType.orValue('') == 'SecretsStore' ? self.metadata.name == 'secrets-store.csi.k8s.io' : true",message="driverType 'SecretsStore' requires metadata.name 'secrets-store.csi.k8s.io'"
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'secrets-store.csi.k8s.io' ? (!has(self.spec.driverConfig) || self.spec.driverConfig.driverType == 'SecretsStore') : true",message="metadata.name 'secrets-store.csi.k8s.io' requires driverType 'SecretsStore'"
// +kubebuilder:validation:XValidation:rule="oldSelf.spec.?driverConfig.?secretsStore.?tokenRequests.?type.orValue('') != 'Managed' || self.spec.?driverConfig.?secretsStore.?tokenRequests.?type.orValue('') == 'Managed'",message="tokenRequests type cannot be changed from Managed"

// ClusterCSIDriver object allows management and configuration of a CSI driver operator
// installed by default in OpenShift. Name of the object must be name of the CSI driver
// it operates. See CSIDriverName type for list of allowed values.
// With CinderCSIDirectAttach enabled, Cinder configuration and status are valid
// only on cinder.csi.openstack.org, which permits only omitted or Cinder configuration.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ClusterCSIDriver struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec ClusterCSIDriverSpec `json:"spec"`

	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ClusterCSIDriverStatus `json:"status"`
}

// CSIDriverName is the name of the CSI driver
type CSIDriverName string

// +kubebuilder:validation:Enum="";Managed;Unmanaged;Removed
// StorageClassStateName defines various configuration states for storageclass management
// and reconciliation by CSI operator.
type StorageClassStateName string

const (
	// ManagedStorageClass means that the operator is actively managing its storage classes.
	// Most manual changes made by cluster admin to storageclass will be wiped away by CSI
	// operator if StorageClassState is set to Managed.
	ManagedStorageClass StorageClassStateName = "Managed"
	// UnmanagedStorageClass means that the operator is not actively managing storage classes.
	// If StorageClassState is Unmanaged then CSI operator will not be actively reconciling storage class
	// it previously created. This can be useful if cluster admin wants to modify storage class installed
	// by CSI operator.
	UnmanagedStorageClass StorageClassStateName = "Unmanaged"
	// RemovedStorageClass instructs the operator to remove the storage class.
	// If StorageClassState is Removed - CSI operator will delete storage classes it created
	// previously. This can be useful in clusters where cluster admins want to prevent
	// creation of dynamically provisioned volumes but still need rest of the features
	// provided by CSI operator and driver.
	RemovedStorageClass StorageClassStateName = "Removed"
)

// If you are adding a new driver name here, ensure that 0000_50_cluster_csi_driver_01_config.crd.yaml-merge-patch file is also updated with new driver name.
const (
	AWSEBSCSIDriver          CSIDriverName = "ebs.csi.aws.com"
	AWSEFSCSIDriver          CSIDriverName = "efs.csi.aws.com"
	AzureDiskCSIDriver       CSIDriverName = "disk.csi.azure.com"
	AzureFileCSIDriver       CSIDriverName = "file.csi.azure.com"
	GCPFilestoreCSIDriver    CSIDriverName = "filestore.csi.storage.gke.io"
	GCPPDCSIDriver           CSIDriverName = "pd.csi.storage.gke.io"
	CinderCSIDriver          CSIDriverName = "cinder.csi.openstack.org"
	VSphereCSIDriver         CSIDriverName = "csi.vsphere.vmware.com"
	ManilaCSIDriver          CSIDriverName = "manila.csi.openstack.org"
	KubevirtCSIDriver        CSIDriverName = "csi.kubevirt.io"
	SharedResourcesCSIDriver CSIDriverName = "csi.sharedresource.openshift.io"
	AlibabaDiskCSIDriver     CSIDriverName = "diskplugin.csi.alibabacloud.com"
	IBMVPCBlockCSIDriver     CSIDriverName = "vpc.block.csi.ibm.io"
	IBMPowerVSBlockCSIDriver CSIDriverName = "powervs.csi.ibm.com"
	SecretsStoreCSIDriver    CSIDriverName = "secrets-store.csi.k8s.io"
	SambaCSIDriver           CSIDriverName = "smb.csi.k8s.io"
)

// ClusterCSIDriverSpec is the desired behavior of CSI driver operator
type ClusterCSIDriverSpec struct {
	OperatorSpec `json:",inline"`
	// storageClassState determines if CSI operator should create and manage storage classes.
	// If this field value is empty or Managed - CSI operator will continuously reconcile
	// storage class and create if necessary.
	// If this field value is Unmanaged - CSI operator will not reconcile any previously created
	// storage class.
	// If this field value is Removed - CSI operator will delete the storage class it created previously.
	// When omitted, this means the user has no opinion and the platform chooses a reasonable default,
	// which is subject to change over time.
	// The current default behaviour is Managed.
	// +optional
	StorageClassState StorageClassStateName `json:"storageClassState,omitempty"`

	// driverConfig can be used to specify platform specific driver configuration.
	// When omitted, this means no opinion and the platform is left to choose reasonable
	// defaults. These defaults are subject to change over time.
	// +optional
	DriverConfig CSIDriverConfigSpec `json:"driverConfig"`
}

// CSIDriverType indicates type of CSI driver being configured.
// +openshift:validation:FeatureGateAwareEnum:featureGate="",enum="";AWS;Azure;GCP;IBMCloud;vSphere;SecretsStore
// +openshift:validation:FeatureGateAwareEnum:featureGate=CinderCSIDirectAttach,enum="";AWS;Azure;GCP;IBMCloud;vSphere;SecretsStore;Cinder
type CSIDriverType string

const (
	AWSDriverType          CSIDriverType = "AWS"
	AzureDriverType        CSIDriverType = "Azure"
	GCPDriverType          CSIDriverType = "GCP"
	IBMCloudDriverType     CSIDriverType = "IBMCloud"
	VSphereDriverType      CSIDriverType = "vSphere"
	SecretsStoreDriverType CSIDriverType = "SecretsStore"
	CinderDriverType       CSIDriverType = "Cinder"
)

// CSIDriverConfigSpec defines configuration spec that can be
// used to optionally configure a specific CSI Driver.
// Selecting Cinder requires cinder and prohibits all other driver configuration
// members. The cinder member must be omitted for all other driver types.
// +openshift:validation:FeatureGateAwareXValidation:featureGate=CinderCSIDirectAttach,rule="has(self.driverType) && self.driverType == 'Cinder' ? has(self.cinder) : !has(self.cinder)",message="cinder must be set if driverType is Cinder, but remain unset otherwise"
// +openshift:validation:FeatureGateAwareXValidation:featureGate=CinderCSIDirectAttach,rule="has(self.driverType) && self.driverType == 'Cinder' ? !has(self.aws) && !has(self.azure) && !has(self.gcp) && !has(self.ibmcloud) && !has(self.vSphere) && !has(self.secretsStore) : true",message="other driver configuration must be unset when driverType is Cinder"
// +kubebuilder:validation:XValidation:rule="has(self.driverType) && self.driverType == 'IBMCloud' ? has(self.ibmcloud) : !has(self.ibmcloud)",message="ibmcloud must be set if driverType is 'IBMCloud', but remain unset otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.driverType) && self.driverType == 'SecretsStore' ? has(self.secretsStore) : !has(self.secretsStore)",message="secretsStore must be set if driverType is 'SecretsStore', but remain unset otherwise"
// +union
type CSIDriverConfigSpec struct {
	// driverType indicates type of CSI driver for which the
	// driverConfig is being applied to.
	// Valid values are: AWS, Azure, GCP, IBMCloud, vSphere, SecretsStore and empty.
	// Cinder is also valid when CinderCSIDirectAttach is enabled.
	// Consumers should treat unknown values as a NO-OP.
	// +required
	// +unionDiscriminator
	DriverType CSIDriverType `json:"driverType"`

	// aws is used to configure the AWS CSI driver.
	// +optional
	AWS *AWSCSIDriverConfigSpec `json:"aws,omitempty"`

	// azure is used to configure the Azure CSI driver.
	// +optional
	Azure *AzureCSIDriverConfigSpec `json:"azure,omitempty"`

	// gcp is used to configure the GCP CSI driver.
	// +optional
	GCP *GCPCSIDriverConfigSpec `json:"gcp,omitempty"`

	// ibmcloud is used to configure the IBM Cloud CSI driver.
	// +optional
	IBMCloud *IBMCloudCSIDriverConfigSpec `json:"ibmcloud,omitempty"`

	// vSphere is used to configure the vsphere CSI driver.
	// +optional
	VSphere *VSphereCSIDriverConfigSpec `json:"vSphere,omitempty"`

	// secretsStore is used to configure the Secrets Store CSI driver.
	// +optional
	SecretsStore SecretsStoreCSIDriverConfigSpec `json:"secretsStore,omitzero"`

	// cinder is required when driverType is Cinder and forbidden otherwise.
	// Omission preserves Nova attachment unless the operator has locked Direct,
	// in which case the operator rejects removal without changing the mode.
	// +openshift:enable:FeatureGate=CinderCSIDirectAttach
	// +optional
	Cinder CinderCSIDriverConfigSpec `json:"cinder,omitzero"`
}

// CinderCSIDriverConfigSpec configures Cinder volume attachment and backend files.
// +kubebuilder:validation:XValidation:rule="self.attachMode == 'Direct' || !has(self.backendConfig) || size(self.backendConfig) == 0",message="backendConfig is valid only in Direct mode"
type CinderCSIDriverConfigSpec struct {
	// attachMode is required. Valid values are Nova (Nova-mediated attachment)
	// and Direct (worker-side os-brick attachment without Nova on all workers).
	// The operator permits Direct only on self-managed OpenStack clusters before
	// any Cinder use and when no Cinder PVCs or PVs exist.
	// The operator locks Nova on first Cinder use and Direct on acceptance,
	// before rollout. The lock survives rollout failures and resource recreation.
	// The operator rejects changes or removal that would change the locked mode.
	// HyperShift rejects Direct without changing workloads. These restrictions
	// are enforced by the operator, not API admission.
	// +required
	AttachMode CinderAttachMode `json:"attachMode,omitempty"`

	// backendConfig optionally supplies up to 16 read-only Secret file sources
	// to each os-brick sidecar. Source names must be unique; mount paths may be
	// shared if filenames differ. Non-empty configuration requires Direct mode.
	// Omission or an empty list mounts no additional files. Before rollout, the
	// operator rejects missing or empty Secrets, unsafe filenames, and duplicate
	// destination filenames. Sources may be updated without changing the mode.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=16
	// +optional
	BackendConfig []CinderBackendConfig `json:"backendConfig,omitempty"`
}

// CinderAttachMode selects the Cinder volume attachment mechanism.
// +kubebuilder:validation:Enum=Nova;Direct
type CinderAttachMode string

const (
	// CinderAttachModeNova uses Nova-mediated volume attachment.
	CinderAttachModeNova CinderAttachMode = "Nova"
	// CinderAttachModeDirect uses worker-side os-brick attachment without Nova.
	CinderAttachModeDirect CinderAttachMode = "Direct"
)

// CinderBackendConfig projects a Secret's files at a supported sidecar directory.
type CinderBackendConfig struct {
	// name is required and uniquely identifies this source within backendConfig.
	// It must contain 1 to 63 lowercase alphanumeric characters or hyphens,
	// starting and ending with an alphanumeric character.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?$')",message="name must be a DNS label"
	// +required
	Name string `json:"name,omitempty"`

	// secret is required and names an administrator-managed Secret in
	// openshift-cluster-csi-drivers. All data keys become read-only sidecar files.
	// The operator validates the Secret's existence and contents.
	// +required
	Secret CinderBackendSecret `json:"secret,omitzero"`

	// mountPath is required. The only supported value is /etc/ceph, the read-only
	// sidecar directory for Ceph configuration and keyrings.
	// +required
	MountPath CinderBackendMountPath `json:"mountPath,omitempty"`
}

// CinderBackendSecret names a backend configuration Secret in
// openshift-cluster-csi-drivers.
type CinderBackendSecret struct {
	// name is required and identifies the Secret. It must contain 1 to 253
	// lowercase alphanumeric characters, hyphens, or dots. Each dot-separated
	// component must start and end with an alphanumeric character.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?([.][a-z0-9]([-a-z0-9]*[a-z0-9])?)*$')",message="secret name must be a DNS subdomain of 1 to 253 characters"
	// +required
	Name string `json:"name,omitempty"`
}

// CinderBackendMountPath is an allowed backend file destination in the sidecar.
// +kubebuilder:validation:Enum="/etc/ceph"
type CinderBackendMountPath string

const (
	// CinderBackendMountPathCeph is the directory for Ceph configuration and keyrings.
	CinderBackendMountPathCeph CinderBackendMountPath = "/etc/ceph"
)

// AWSCSIDriverConfigSpec defines properties that can be configured for the AWS CSI driver.
type AWSCSIDriverConfigSpec struct {
	// kmsKeyARN sets the cluster default storage class to encrypt volumes with a user-defined KMS key,
	// rather than the default KMS key used by AWS.
	// The value may be either the ARN or Alias ARN of a KMS key.
	//
	// The ARN must follow the format: arn:<partition>:kms:<region>:<account-id>:(key|alias)/<key-id-or-alias>, where:
	// <partition> is the AWS partition (aws, aws-cn, aws-us-gov, aws-iso, aws-iso-b, aws-iso-e, aws-iso-f, or aws-eusc),
	// <region> is the AWS region,
	// <account-id> is a 12-digit numeric identifier for the AWS account,
	// <key-id-or-alias> is the KMS key ID or alias name.
	//
	// +openshift:validation:FeatureGateAwareXValidation:featureGate="",rule=`matches(self, '^arn:(aws|aws-cn|aws-us-gov|aws-iso|aws-iso-b|aws-iso-e|aws-iso-f):kms:[a-z0-9-]+:[0-9]{12}:(key|alias)/.*$')`,message=`kmsKeyARN must be a valid AWS KMS key ARN in the format: arn:<partition>:kms:<region>:<account-id>:(key|alias)/<key-id-or-alias>`
	// +openshift:validation:FeatureGateAwareXValidation:featureGate=AWSEuropeanSovereignCloudInstall,rule=`matches(self, '^arn:(aws|aws-cn|aws-us-gov|aws-iso|aws-iso-b|aws-iso-e|aws-iso-f|aws-eusc):kms:[a-z0-9-]+:[0-9]{12}:(key|alias)/.*$')`,message=`kmsKeyARN must be a valid AWS KMS key ARN in the format: arn:<partition>:kms:<region>:<account-id>:(key|alias)/<key-id-or-alias>`
	// +optional
	KMSKeyARN string `json:"kmsKeyARN,omitempty"`

	// efsVolumeMetrics sets the configuration for collecting metrics from EFS volumes used by the EFS CSI Driver.
	// +optional
	EFSVolumeMetrics *AWSEFSVolumeMetrics `json:"efsVolumeMetrics,omitempty"`
}

// AWSEFSVolumeMetricsState defines the modes for collecting volume metrics in the AWS EFS CSI Driver.
// This can either enable recursive collection of volume metrics or disable metric collection entirely.
// +kubebuilder:validation:Enum:="RecursiveWalk";"Disabled"
type AWSEFSVolumeMetricsState string

const (
	// AWSEFSVolumeMetricsRecursiveWalk indicates that volume metrics collection in the AWS EFS CSI Driver
	// is performed by recursively walking through the files in the volume.
	AWSEFSVolumeMetricsRecursiveWalk AWSEFSVolumeMetricsState = "RecursiveWalk"

	// AWSEFSVolumeMetricsDisabled indicates that volume metrics collection in the AWS EFS CSI Driver is disabled.
	AWSEFSVolumeMetricsDisabled AWSEFSVolumeMetricsState = "Disabled"
)

// AWSEFSVolumeMetrics defines the configuration for volume metrics in the EFS CSI Driver.
// +union
type AWSEFSVolumeMetrics struct {
	// state defines the state of metric collection in the AWS EFS CSI Driver.
	// This field is required and must be set to one of the following values: Disabled or RecursiveWalk.
	// Disabled means no metrics collection will be performed. This is the default value.
	// RecursiveWalk means the AWS EFS CSI Driver will recursively scan volumes to collect metrics.
	// This process may result in high CPU and memory usage, depending on the volume size.
	// +unionDiscriminator
	// +required
	State AWSEFSVolumeMetricsState `json:"state"`

	// recursiveWalk provides additional configuration for collecting volume metrics in the AWS EFS CSI Driver
	// when the state is set to RecursiveWalk.
	// +unionMember
	// +optional
	RecursiveWalk *AWSEFSVolumeMetricsRecursiveWalkConfig `json:"recursiveWalk,omitempty"`
}

// AWSEFSVolumeMetricsRecursiveWalkConfig defines options for volume metrics in the EFS CSI Driver.
type AWSEFSVolumeMetricsRecursiveWalkConfig struct {
	// refreshPeriodMinutes specifies the frequency, in minutes, at which volume metrics are refreshed.
	// When omitted, this means no opinion and the platform is left to choose a reasonable
	// default, which is subject to change over time. The current default is 240.
	// The valid range is from 1 to 43200 minutes (30 days).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=43200
	// +optional
	RefreshPeriodMinutes int32 `json:"refreshPeriodMinutes,omitempty"`

	// fsRateLimit defines the rate limit, in goroutines per file system, for processing volume metrics.
	// When omitted, this means no opinion and the platform is left to choose a reasonable
	// default, which is subject to change over time. The current default is 5.
	// The valid range is from 1 to 100 goroutines.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +optional
	FSRateLimit int32 `json:"fsRateLimit,omitempty"`
}

// AzureDiskEncryptionSet defines the configuration for a disk encryption set.
type AzureDiskEncryptionSet struct {
	// subscriptionID defines the Azure subscription that contains the disk encryption set.
	// The value should meet the following conditions:
	// 1. It should be a 128-bit number.
	// 2. It should be 36 characters (32 hexadecimal characters and 4 hyphens) long.
	// 3. It should be displayed in five groups separated by hyphens (-).
	// 4. The first group should be 8 characters long.
	// 5. The second, third, and fourth groups should be 4 characters long.
	// 6. The fifth group should be 12 characters long.
	// An Example SubscrionID: f2007bbf-f802-4a47-9336-cf7c6b89b378
	// +required
	// +kubebuilder:validation:MaxLength:=36
	// +kubebuilder:validation:Pattern:=`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`
	SubscriptionID string `json:"subscriptionID"`

	// resourceGroup defines the Azure resource group that contains the disk encryption set.
	// The value should consist of only alphanumberic characters,
	// underscores (_), parentheses, hyphens and periods.
	// The value should not end in a period and be at most 90 characters in
	// length.
	// +required
	// +kubebuilder:validation:MaxLength:=90
	// +kubebuilder:validation:Pattern:=`^[\w\.\-\(\)]*[\w\-\(\)]$`
	ResourceGroup string `json:"resourceGroup"`

	// name is the name of the disk encryption set that will be set on the default storage class.
	// The value should consist of only alphanumberic characters,
	// underscores (_), hyphens, and be at most 80 characters in length.
	// +required
	// +kubebuilder:validation:MaxLength:=80
	// +kubebuilder:validation:Pattern:=`^[a-zA-Z0-9\_-]+$`
	Name string `json:"name"`
}

// AzureCSIDriverConfigSpec defines properties that can be configured for the Azure CSI driver.
type AzureCSIDriverConfigSpec struct {
	// diskEncryptionSet sets the cluster default storage class to encrypt volumes with a
	// customer-managed encryption set, rather than the default platform-managed keys.
	// +optional
	DiskEncryptionSet *AzureDiskEncryptionSet `json:"diskEncryptionSet,omitempty"`
}

// GCPKMSKeyReference gathers required fields for looking up a GCP KMS Key
type GCPKMSKeyReference struct {
	// name is the name of the customer-managed encryption key to be used for disk encryption.
	// The value should correspond to an existing KMS key and should
	// consist of only alphanumeric characters, hyphens (-) and underscores (_),
	// and be at most 63 characters in length.
	// +kubebuilder:validation:Pattern:=`^[a-zA-Z0-9\_-]+$`
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +required
	Name string `json:"name"`

	// keyRing is the name of the KMS Key Ring which the KMS Key belongs to.
	// The value should correspond to an existing KMS key ring and should
	// consist of only alphanumeric characters, hyphens (-) and underscores (_),
	// and be at most 63 characters in length.
	// +kubebuilder:validation:Pattern:=`^[a-zA-Z0-9\_-]+$`
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=63
	// +required
	KeyRing string `json:"keyRing"`

	// projectID is the ID of the Project in which the KMS Key Ring exists.
	// It must be 6 to 30 lowercase letters, digits, or hyphens.
	// It must start with a letter. Trailing hyphens are prohibited.
	// +kubebuilder:validation:Pattern:=`^[a-z][a-z0-9-]+[a-z0-9]$`
	// +kubebuilder:validation:MinLength:=6
	// +kubebuilder:validation:MaxLength:=30
	// +required
	ProjectID string `json:"projectID"`

	// location is the GCP location in which the Key Ring exists.
	// The value must match an existing GCP location, or "global".
	// Defaults to global, if not set.
	// +kubebuilder:validation:Pattern:=`^[a-zA-Z0-9\_-]+$`
	// +optional
	Location string `json:"location,omitempty"`
}

// GCPCSIDriverConfigSpec defines properties that can be configured for the GCP CSI driver.
type GCPCSIDriverConfigSpec struct {
	// kmsKey sets the cluster default storage class to encrypt volumes with customer-supplied
	// encryption keys, rather than the default keys managed by GCP.
	// +optional
	KMSKey *GCPKMSKeyReference `json:"kmsKey,omitempty"`
}

// IBMCloudCSIDriverConfigSpec defines the properties that can be configured for the IBM Cloud CSI driver.
type IBMCloudCSIDriverConfigSpec struct {
	// encryptionKeyCRN is the IBM Cloud CRN of the customer-managed root key to use
	// for disk encryption of volumes for the default storage classes.
	// +required
	// +kubebuilder:validation:MaxLength:=154
	// +kubebuilder:validation:MinLength:=144
	// +kubebuilder:validation:Pattern:=`^crn:v[0-9]+:bluemix:(public|private):(kms|hs-crypto):[a-z-]+:a/[0-9a-f]+:[0-9a-f-]{36}:key:[0-9a-f-]{36}$`
	EncryptionKeyCRN string `json:"encryptionKeyCRN"`
}

// VSphereCSIDriverConfigSpec defines properties that
// can be configured for vsphere CSI driver.
type VSphereCSIDriverConfigSpec struct {
	// topologyCategories indicates tag categories with which
	// vcenter resources such as hostcluster or datacenter were tagged with.
	// If cluster Infrastructure object has a topology, values specified in
	// Infrastructure object will be used and modifications to topologyCategories
	// will be rejected.
	// +listType=atomic
	// +optional
	TopologyCategories []string `json:"topologyCategories,omitempty"`

	// globalMaxSnapshotsPerBlockVolume is a global configuration parameter that applies to volumes on all kinds of
	// datastores. If omitted, the platform chooses a default, which is subject to change over time, currently that default is 3.
	// Snapshots can not be disabled using this parameter.
	// Increasing number of snapshots above 3 can have negative impact on performance, for more details see: https://kb.vmware.com/s/article/1025279
	// Volume snapshot documentation: https://docs.vmware.com/en/VMware-vSphere-Container-Storage-Plug-in/3.0/vmware-vsphere-csp-getting-started/GUID-E0B41C69-7EEB-450F-A73D-5FD2FF39E891.html
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=32
	// +optional
	GlobalMaxSnapshotsPerBlockVolume *uint32 `json:"globalMaxSnapshotsPerBlockVolume,omitempty"`

	// granularMaxSnapshotsPerBlockVolumeInVSAN is a granular configuration parameter on vSAN datastore only. It
	// overrides GlobalMaxSnapshotsPerBlockVolume if set, while it falls back to the global constraint if unset.
	// Snapshots for VSAN can not be disabled using this parameter.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=32
	// +optional
	GranularMaxSnapshotsPerBlockVolumeInVSAN *uint32 `json:"granularMaxSnapshotsPerBlockVolumeInVSAN,omitempty"`

	// granularMaxSnapshotsPerBlockVolumeInVVOL is a granular configuration parameter on Virtual Volumes datastore only.
	// It overrides GlobalMaxSnapshotsPerBlockVolume if set, while it falls back to the global constraint if unset.
	// Snapshots for VVOL can not be disabled using this parameter.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=32
	// +optional
	GranularMaxSnapshotsPerBlockVolumeInVVOL *uint32 `json:"granularMaxSnapshotsPerBlockVolumeInVVOL,omitempty"`

	// maxAllowedBlockVolumesPerNode is an optional configuration parameter that allows setting a custom value for the
	// limit of the number of PersistentVolumes attached to a node. In vSphere version 7 this limit was set to 59 by
	// default, however in vSphere version 8 this limit was increased to 255.
	// Before increasing this value above 59 the cluster administrator needs to ensure that every node forming the
	// cluster is updated to ESXi version 8 or higher and that all nodes are running the same version.
	// The limit must be between 1 and 255, which matches the vSphere version 8 maximum.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to
	// change over time.
	// The current default is 59, which matches the limit for vSphere version 7.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=255
	// +openshift:enable:FeatureGate=VSphereConfigurableMaxAllowedBlockVolumesPerNode
	// +optional
	MaxAllowedBlockVolumesPerNode int32 `json:"maxAllowedBlockVolumesPerNode,omitempty"`
}

// SecretsStoreCSIDriverConfigSpec defines properties that can be configured for the Secrets Store CSI driver.
// +kubebuilder:validation:MinProperties=1
type SecretsStoreCSIDriverConfigSpec struct {
	// secretRotation controls automatic secret rotation behavior.
	// When omitted, secret rotation is enabled with a default poll interval of 2 minutes.
	// +optional
	SecretRotation SecretsStoreSecretRotation `json:"secretRotation,omitzero"`

	// tokenRequests controls service account token configuration for
	// workload identity federation (WIF) with cloud providers.
	// When omitted, the operator preserves any existing tokenRequests
	// already configured on the CSIDriver object without modification.
	// +optional
	TokenRequests SecretsStoreTokenRequests `json:"tokenRequests,omitzero"`
}

// TokenRequestsType determines how the operator manages the tokenRequests
// field on the storage.k8s.io CSIDriver object.
// +kubebuilder:validation:Enum=Managed;Unmanaged
type TokenRequestsType string

const (
	// TokenRequestsManaged means the operator uses the audiences list
	// as the sole source of truth for the CSIDriver.spec.tokenRequests field.
	TokenRequestsManaged TokenRequestsType = "Managed"

	// TokenRequestsUnmanaged means the operator preserves any existing
	// tokenRequests already configured on the CSIDriver object and does not
	// overwrite them.
	TokenRequestsUnmanaged TokenRequestsType = "Unmanaged"
)

// SecretsStoreTokenRequests configures how service account tokens are
// provided to the Secrets Store CSI driver for workload identity federation.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Managed' ? has(self.managed) : !has(self.managed)",message="managed must be set when type is 'Managed', and must not be set otherwise"
// +union
type SecretsStoreTokenRequests struct {
	// type determines how the operator manages tokenRequests on the CSIDriver object.
	// When "Unmanaged", existing tokenRequests on the CSIDriver are preserved
	// and the managed field is not used.
	// When "Managed", the operator sets tokenRequests from the audiences
	// specified in the managed field, replacing any previously configured values.
	// Once set to "Managed", type cannot be reverted back to "Unmanaged".
	// +unionDiscriminator
	// +required
	Type TokenRequestsType `json:"type,omitempty"`

	// managed holds configuration for operator-managed tokenRequests.
	// Only valid when type is "Managed".
	// +optional
	Managed ManagedTokenRequests `json:"managed,omitzero"`
}

// ManagedTokenRequests holds the configuration for operator-managed
// service account token requests.
// +kubebuilder:validation:MinProperties=1
type ManagedTokenRequests struct {
	// audiences specifies service account token audiences that kubelet will
	// provide to the CSI driver during NodePublishVolume calls. These tokens
	// enable workload identity federation (WIF) with cloud providers such as
	// AWS, Azure, and GCP.
	// When empty, the operator clears all tokenRequests from the CSIDriver object.
	// +optional
	// +listType=map
	// +listMapKey=audience
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=10
	Audiences *[]SecretsStoreTokenRequest `json:"audiences,omitempty"`
}

// SecretRotationType determines the secret rotation behavior for the
// Secrets Store CSI driver.
// +kubebuilder:validation:Enum=None;Custom
type SecretRotationType string

const (
	// SecretRotationNone disables automatic secret rotation. Secrets are only
	// fetched at initial pod mount time.
	SecretRotationNone SecretRotationType = "None"

	// SecretRotationCustom enables automatic secret rotation with the
	// configuration specified in the custom field.
	SecretRotationCustom SecretRotationType = "Custom"
)

// SecretsStoreSecretRotation configures the automatic secret rotation behavior
// for the Secrets Store CSI driver.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Custom' ? has(self.custom) : !has(self.custom)",message="custom must be set when type is 'Custom', and must not be set otherwise"
// +union
type SecretsStoreSecretRotation struct {
	// type determines the secret rotation behavior.
	// When "None", secret rotation is disabled and secrets are only fetched at
	// initial pod mount time.
	// When "Custom", secret rotation is enabled with the configuration specified
	// in the custom field.
	// +unionDiscriminator
	// +required
	Type SecretRotationType `json:"type,omitempty"`

	// custom holds the custom rotation configuration.
	// Only valid when type is "Custom".
	// +optional
	Custom CustomSecretRotation `json:"custom,omitzero"`
}

// CustomSecretRotation holds configuration for custom secret rotation behavior.
// +kubebuilder:validation:MinProperties=1
type CustomSecretRotation struct {
	// minimumRefreshAge is the minimum time in seconds between secret
	// rotation attempts. Each time kubelet calls NodePublishVolume, the driver
	// checks whether this interval has elapsed since the last successful provider
	// call. If it has, the driver contacts the secret provider to fetch the latest
	// secret values and updates the mounted volume.
	// Setting this value below the kubelet syncFrequency (default: 1 minute)
	// has no additional effect on the actual rotation cadence.
	// Must be at least 1 second and no more than 31560000 seconds (~1 year).
	// When omitted, this means no opinion and the platform is left to choose a
	// reasonable default, which is subject to change over time.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=31560000
	// +optional
	MinimumRefreshAge int32 `json:"minimumRefreshAge,omitempty"`

	// --- TOMBSTONE ---
	// rotationPollIntervalSeconds was the previous name for minimumRefreshAge.
	// The field has been renamed to better reflect its semantics.
	// The JSON key is reserved to prevent reuse.
	//
	// +optional
	// RotationPollIntervalSeconds int32 `json:"rotationPollIntervalSeconds,omitempty"`
}

// SecretsStoreTokenRequest specifies a service account token audience configuration
// for workload identity federation (WIF) with the Secrets Store CSI driver.
type SecretsStoreTokenRequest struct {
	// audience is the intended audience of the service account token.
	// An empty string means the issued token will use the kube-apiserver's default APIAudiences.
	// +kubebuilder:validation:MinLength=0
	// +kubebuilder:validation:MaxLength=253
	// +required
	Audience *string `json:"audience,omitempty"`

	// expirationSeconds is the requested duration of validity of the service account token.
	// The token issuer may return a token with a different validity duration.
	// When omitted, the token expiration is determined by the kube-apiserver.
	// Must be at least 600 seconds (10 minutes) and no more than 315360000 seconds (~10 years).
	// +kubebuilder:validation:Minimum=600
	// +kubebuilder:validation:Maximum=315360000
	// +optional
	ExpirationSeconds int32 `json:"expirationSeconds,omitempty"`
}

// ClusterCSIDriverStatus is the observed status of CSI driver operator
type ClusterCSIDriverStatus struct {
	OperatorStatus `json:",inline"`

	// cinder optionally reports Cinder configuration on cinder.csi.openstack.org.
	// When omitted, no completed Cinder rollout has been reported.
	// +openshift:enable:FeatureGate=CinderCSIDirectAttach
	// +optional
	Cinder *CinderCSIDriverStatus `json:"cinder,omitempty"`
}

// CinderCSIDriverStatus reports the last successfully applied attachment mode.
// An empty status is valid before a rollout has completed.
// +kubebuilder:validation:MinProperties=0
type CinderCSIDriverStatus struct {
	// observedAttachMode optionally reports the last successfully completed
	// controller and node rollout. Valid values are Nova (Nova-mediated attachment)
	// and Direct (worker-side os-brick attachment without Nova).
	// Omission means no completed rollout has been observed. Pending, failed, or
	// rejected changes retain the previous value. This is neither the mode lock
	// nor proof that the current generation has rolled out; check conditions and
	// observedGeneration.
	// +optional
	ObservedAttachMode CinderAttachMode `json:"observedAttachMode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCSIDriverList contains a list of ClusterCSIDriver
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ClusterCSIDriverList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterCSIDriver `json:"items"`
}
