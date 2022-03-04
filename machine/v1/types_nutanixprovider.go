package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NutanixMachineProviderConfig is the Schema for the nutanixmachineproviderconfigs API
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +k8s:openapi-gen=true
type NutanixMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// clusterReference is to identify the PE/cluster in which the Machine's VM will be created
	// +kubebuilder:validation:Required
	ClusterReference *NutanixReference `json:"clusterReference"`

	// imageReference is to identify the rhcos image uploaded to the Prism Central (PC)
	// +kubebuilder:validation:Required
	ImageReference *NutanixReference `json:"imageReference"`

	// subnetReference is to identify the cluster's network subnet to use for the Machine's VM
	// +kubebuilder:validation:Required
	SubnetReference *NutanixReference `json:"subnetReference"`

	// numVcpusPerSocket is the number of vCPUs per socket of the VM to create
	// +kubebuilder:validation:Required
	NumVcpusPerSocket int64 `json:"numVcpusPerSocket"`

	// numSockets is the number of sockets of the VM to create
	// +kubebuilder:validation:Required
	NumSockets int64 `json:"numSockets"`

	// memorySize is the memory size (in Quantity format) of the VM to create
	// +kubebuilder:validation:Required
	MemorySize resource.Quantity `json:"memorySize"`

	// diskSize is the system disk size (in Quantity format) of the VM to create
	// +kubebuilder:validation:Required
	DiskSize resource.Quantity `json:"diskSize"`

	// userDataSecret is a local reference to a secret that contains the
	// UserData to apply to the VM
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// credentialsSecret is a local reference to a secret that contains the
	// credentials data to access Nutanix PC client
	// +kubebuilder:validation:Required
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret"`
}

// NutanixReference holds the identity of a Nutanix PC resource (cluster, image, subnet, etc.)
type NutanixReference struct {
	// uuid is the UUID of the resource in the PC.
	// If this is configured, it will be used to create the VM.
	// Otherwise, the resource name will be used to obtain the UUID, before creating the VM.
	UUID string `json:"uuid,omitempty"`

	// name is the resource name in the PC
	Name string `json:"name,omitempty"`
}

// NutanixMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains nutanix-specific status information.
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NutanixMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// vmState is the Machine associated VM's current state
	// +optional
	VmState *string `json:"vmState,omitempty"`

	// vmUUID is the Machine associated VM's UUID
	// +optional
	VmUUID *string `json:"vmUUID,omitempty"`
}
