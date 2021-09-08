package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AzureMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Azure virtual machine. It is used by the Azure machine actuator to create a single Machine.
// Required parameters such as location that are not specified by this configuration, will be defaulted
// by the actuator.
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type AzureMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.SecretReference `json:"userDataSecret,omitempty" protobuf:"bytes,2,opt,name=userDataSecret"`

	// CredentialsSecret is a reference to the secret with Azure credentials.
	CredentialsSecret *corev1.SecretReference `json:"credentialsSecret,omitempty" protobuf:"bytes,3,opt,name=credentialsSecret"`

	Location     string            `json:"location,omitempty" protobuf:"bytes,4,opt,name=location"`
	VMSize       string            `json:"vmSize,omitempty" protobuf:"bytes,5,opt,name=vmSize"`
	Image        Image             `json:"image" protobuf:"bytes,6,opt,name=image"`
	OSDisk       OSDisk            `json:"osDisk" protobuf:"bytes,7,opt,name=osDisk"`
	SSHPublicKey string            `json:"sshPublicKey,omitempty" protobuf:"bytes,8,opt,name=sshPublicKey"`
	PublicIP     bool              `json:"publicIP" protobuf:"varint,9,opt,name=publicIP"`
	Tags         map[string]string `json:"tags,omitempty" protobuf:"bytes,10,rep,name=tags"`

	// Network Security Group that needs to be attached to the machine's interface.
	// No security group will be attached if empty.
	SecurityGroup string `json:"securityGroup,omitempty" protobuf:"bytes,11,opt,name=securityGroup"`

	// Application Security Groups that need to be attached to the machine's interface.
	// No application security groups will be attached if zero-length.
	ApplicationSecurityGroups []string `json:"applicationSecurityGroups,omitempty" protobuf:"bytes,12,rep,name=applicationSecurityGroups"`

	// Subnet to use for this instance
	Subnet string `json:"subnet" protobuf:"bytes,13,opt,name=subnet"`

	// PublicLoadBalancer to use for this instance
	PublicLoadBalancer string `json:"publicLoadBalancer,omitempty" protobuf:"bytes,14,opt,name=publicLoadBalancer"`

	// InternalLoadBalancerName to use for this instance
	InternalLoadBalancer string `json:"internalLoadBalancer,omitempty" protobuf:"bytes,15,opt,name=internalLoadBalancer"`

	// NatRule to set inbound NAT rule of the load balancer
	NatRule *int64 `json:"natRule,omitempty" protobuf:"varint,16,opt,name=natRule"`

	// ManagedIdentity to set managed identity name
	ManagedIdentity string `json:"managedIdentity,omitempty" protobuf:"bytes,17,opt,name=managedIdentity"`

	// Vnet to set virtual network name
	Vnet string `json:"vnet,omitempty" protobuf:"bytes,18,opt,name=vnet"`

	// Availability Zone for the virtual machine.
	// If nil, the virtual machine should be deployed to no zone
	Zone *string `json:"zone,omitempty" protobuf:"bytes,19,opt,name=zone"`

	NetworkResourceGroup string `json:"networkResourceGroup,omitempty" protobuf:"bytes,20,opt,name=networkResourceGroup"`
	ResourceGroup        string `json:"resourceGroup,omitempty" protobuf:"bytes,21,opt,name=resourceGroup"`

	// SpotVMOptions allows the ability to specify the Machine should use a Spot VM
	SpotVMOptions *SpotVMOptions `json:"spotVMOptions,omitempty" protobuf:"bytes,22,opt,name=spotVMOptions"`

	// SecurityProfile specifies the Security profile settings for a virtual machine.
	// +optional
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty" protobuf:"bytes,23,opt,name=securityProfile"`
}

// SpotVMOptions defines the options relevant to running the Machine on Spot VMs
type SpotVMOptions struct {
	// MaxPrice defines the maximum price the user is willing to pay for Spot VM instances
	MaxPrice *resource.Quantity `json:"maxPrice,omitempty" protobuf:"bytes,1,opt,name=maxPrice"`
}

// AzureMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains Azure-specific status information.
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type AzureMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// VMID is the ID of the virtual machine created in Azure.
	// +optional
	VMID *string `json:"vmId,omitempty" protobuf:"bytes,2,opt,name=vmId"`

	// VMState is the provisioning state of the Azure virtual machine.
	// +optional
	VMState *VMState `json:"vmState,omitempty" protobuf:"bytes,3,opt,name=vmState,casttype=VMState"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status.
	// +optional
	Conditions []AzureMachineProviderCondition `json:"conditions,omitempty" protobuf:"bytes,4,rep,name=conditions"`
}

// VMState describes the state of an Azure virtual machine.
type VMState string

var (
	// ProvisioningState related values
	// VMStateCreating ...
	VMStateCreating = VMState("Creating")
	// VMStateDeleting ...
	VMStateDeleting = VMState("Deleting")
	// VMStateFailed ...
	VMStateFailed = VMState("Failed")
	// VMStateMigrating ...
	VMStateMigrating = VMState("Migrating")
	// VMStateSucceeded ...
	VMStateSucceeded = VMState("Succeeded")
	// VMStateUpdating ...
	VMStateUpdating = VMState("Updating")

	// PowerState related values
	// VMStateStarting ...
	VMStateStarting = VMState("Starting")
	// VMStateRunning ...
	VMStateRunning = VMState("Running")
	// VMStateStopping ...
	VMStateStopping = VMState("Stopping")
	// VMStateStopped ...
	VMStateStopped = VMState("Stopped")
	// VMStateDeallocating ...
	VMStateDeallocating = VMState("Deallocating")
	// VMStateDeallocated ...
	VMStateDeallocated = VMState("Deallocated")

	// VMStateUnknown ...
	VMStateUnknown = VMState("Unknown")
)

// Image is a mirror of azure sdk compute.ImageReference
type Image struct {
	// Fields below refer to os images in marketplace
	Publisher string `json:"publisher" protobuf:"bytes,1,opt,name=publisher"`
	Offer     string `json:"offer" protobuf:"bytes,2,opt,name=offer"`
	SKU       string `json:"sku" protobuf:"bytes,3,opt,name=sku"`
	Version   string `json:"version" protobuf:"bytes,4,opt,name=version"`
	// ResourceID represents the location of OS Image in azure subscription
	ResourceID string `json:"resourceID" protobuf:"bytes,5,opt,name=resourceID"`
}

type OSDisk struct {
	OSType      string                `json:"osType" protobuf:"bytes,1,opt,name=osType"`
	ManagedDisk ManagedDiskParameters `json:"managedDisk" protobuf:"bytes,2,opt,name=managedDisk"`
	DiskSizeGB  int32                 `json:"diskSizeGB" protobuf:"varint,3,opt,name=diskSizeGB"`
}

type ManagedDiskParameters struct {
	StorageAccountType string                       `json:"storageAccountType" protobuf:"bytes,1,opt,name=storageAccountType"`
	DiskEncryptionSet  *DiskEncryptionSetParameters `json:"diskEncryptionSet,omitempty" protobuf:"bytes,2,opt,name=diskEncryptionSet"`
}

type DiskEncryptionSetParameters struct {
	ID string `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`
}

// SecurityProfile specifies the Security profile settings for a
// virtual machine or virtual machine scale set.
type SecurityProfile struct {
	// This field indicates whether Host Encryption should be enabled
	// or disabled for a virtual machine or virtual machine scale
	// set. Default is disabled.
	EncryptionAtHost *bool `json:"encryptionAtHost,omitempty" protobuf:"varint,1,opt,name=encryptionAtHost"`
}

// AzureMachineProviderCondition is a condition in a AzureMachineProviderStatus
type AzureMachineProviderCondition struct {
	// Type is the type of the condition.
	Type ConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ConditionType"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}
