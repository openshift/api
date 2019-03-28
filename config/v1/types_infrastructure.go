package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import corev1 "k8s.io/api/core/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Infrastructure holds cluster-wide information about Infrastructure.  The canonical name is `cluster`
type Infrastructure struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec InfrastructureSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status InfrastructureStatus `json:"status"`
}

// InfrastructureSpec contains settings that apply to the cluster infrastructure.
type InfrastructureSpec struct {
	// VSphereInfrastructureConfig specifies configuration when the cluster is installed on a VSphere platform.
	// This is only settable if the platform type is VSphere.
	// +optional
	VSphereInfrastructureConfig *VSphereInfrastructureConfig `json:"vsphereInfrastructureConfig,omitempty"`
}

// VSphereInfrastructureConfig specifies the configuration of one or more Virtual Centers where
// the control plane and worker nodes are running in the cluster.  This information is used to
// control the configuration of the vsphere cloud provider support in Kubernetes as supported
// by OpenShift.
type VSphereInfrastructureConfig struct {
	// secretRef is name of the authentication secret for access to the VSphere infrastructure.
	SecretRef corev1.SecretReference `json:"secretRef"`

	// port is the vCenter Server Port.
	// defaults to 443 if not specified.
	// +optional
	Port *int32 `json:"port,omitempty"`

	// insecure if true means vCenter is configured with a self-signed cert.
	// defaults to false.
	// +optional
	Insecure bool `json:"insecure,omitempty"`

	// datacenters in which VMs are located.
	// +optional
	Datacenters []string `json:"datacenters,omitempty"`

	// Per virtual center coniguration overrides.
	// +optional
	VirtualCenters []VSphereVirtualCenterConfig `json:"virtualCenters,omitempty"`

	// Workspace describes the endpoint to create volumes.
	// +optional
	Workspace VSphereWorkspaceConfig `json:"workspace,omitempty"`

	// Network is the configuration information for networking.
	// +optional
	Network VSphereNetworkConfig `json:"network,omitempty"`

	// Disk is the configuration for disks
	// +optional
	Disk VSphereDiskConfig `json:"disk,omitempty"`
}

// VSphereVirtualCenterConfig specifies the configuration for a single
// Virtual Center.  It overrides the default configuration options
// specific to the Virtual Center with the specified name.
type VSphereVirtualCenterConfig struct {
	// server is the name of the virtual center instance (e.g. 1.1.1.1)
	Server string `json:"server"`
	// port is the vCenter Server Port for this virtual center instance.
	// defaults to 443 if not specified.
	// +optional
	Port *int32 `json:"port,omitempty"`
	// Datacenters in which VMs are located.
	Datacenters []string `json:"datacenters,omitempty"`
}

// VSphereNetworkConfig specifies configuration for networking.
type VSphereNetworkConfig struct {
	// publicNetwork is the name of the network the VMs are joined to.
	PublicNetwork string `json:"publicNetwork"`
}

// VSphereDiskConfig specifies configuration for
type VSphereDiskConfig struct {
	// scsiControllerType defines SCSI controller to be used.
	SCSIControllerType string `json:"scsiControllerType"`
}

// VSphereWorkspaceConfig describes the endpoint used to create volumes.
type VSphereWorkspaceConfig struct {
	// server is the virtual center server, for example 1.1.1.1
	Server string `json:"server"`
	// datacenter is the name of datacenter in virtual center server.
	Datacenter string `json:"datacenter"`
	// folder is the virtual center VM folder path under the datacenter.
	Folder string `json:"folder"`
	// resourcePoolPath is the path to the resource pool under the datacenter.
	ResourcePoolPath string `json:"resourcePoolPath"`
}

// InfrastructureStatus describes the infrastructure the cluster is leveraging.
type InfrastructureStatus struct {
	// platform is the underlying infrastructure provider for the cluster. This
	// value controls whether infrastructure automation such as service load
	// balancers, dynamic volume provisioning, machine creation and deletion, and
	// other integrations are enabled. If None, no infrastructure automation is
	// enabled. Allowed values are "AWS", "Azure", "GCP", "Libvirt",
	// "OpenStack", "VSphere", and "None". Individual components may not support
	// all platforms, and must handle unrecognized platforms as None if they do
	// not support that platform.
	Platform PlatformType `json:"platform,omitempty"`

	// etcdDiscoveryDomain is the domain used to fetch the SRV records for discovering
	// etcd servers and clients.
	// For more info: https://github.com/etcd-io/etcd/blob/329be66e8b3f9e2e6af83c123ff89297e49ebd15/Documentation/op-guide/clustering.md#dns-discovery
	EtcdDiscoveryDomain string `json:"etcdDiscoveryDomain"`

	// apiServerURL is a valid URL with scheme(http/https), address and port.
	// apiServerURL can be used by components like kubelet on machines, to contact the `apisever`
	// using the infrastructure provider rather than the kubernetes networking.
	APIServerURL string `json:"apiServerURL"`
}

// PlatformType is a specific supported infrastructure provider.
type PlatformType string

const (
	// AWSPlatform represents Amazon Web Services infrastructure.
	AWSPlatform PlatformType = "AWS"

	// AzurePlatform represents Microsoft Azure infrastructure.
	AzurePlatform PlatformType = "Azure"

	// GCPPlatform represents Google Cloud Platform infrastructure.
	GCPPlatform PlatformType = "GCP"

	// LibvirtPlatform represents libvirt infrastructure.
	LibvirtPlatform PlatformType = "Libvirt"

	// OpenStackPlatform represents OpenStack infrastructure.
	OpenStackPlatform PlatformType = "OpenStack"

	// NonePlatform means there is no infrastructure provider.
	NonePlatform PlatformType = "None"

	// VSpherePlatform represents VMWare vSphere infrastructure.
	VSpherePlatform PlatformType = "VSphere"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureList is
type InfrastructureList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata"`
	Items           []Infrastructure `json:"items"`
}
