package v1

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

const (
	// NutanixCredentialsSecretName is the name of the secret holding the credentials for PC client
	NutanixCredentialsSecretName = "nutanix-credentials"
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

	// clusterReferenceUuid is the UUID of the PE/cluster the Machine's VM will be created
	ClusterReferenceUUID string `json:"clusterReferenceUuid,omitempty"`

	// imageUuid is the UUID of the rhcos image uploaded to the PC.
	// If the imageUUID is configured in the Machine CR, it will be used to create the VM.
	// Otherwise, the imageName will be used to obtain the imageUUID, before creating the VM.
	ImageUUID string `json:"imageUuid,omitempty"`

	// imageName is the name of the rhcos image uploaded to the PC
	ImageName string `json:"imageName,omitempty"`

	// subnetUuid is the UUID of the network subnet to use for the Machine's VM
	SubnetUUID string `json:"subnetUuid,omitempty"`

	// numVcpusPerSocket is the number of vCPUs per socket of the VM to create
	NumVcpusPerSocket int64 `json:"numVcpusPerSocket,omitempty"`

	// numSockets is the number of sockets of the VM to create
	NumSockets int64 `json:"numSockets,omitempty"`

	// memorySizeMib is the memory size in megabytes of the VM to create
	MemorySizeMib int64 `json:"memorySizeMib,omitempty"`

	// diskSizeMib is the disk size in megabytes of the VM to create
	DiskSizeMib int64 `json:"diskSizeMib,omitempty"`

	// powerState is the expected power state of the VM to create
	PowerState string `json:"powerState,omitempty"`

	// userDataSecret is a local reference to a secret that contains the
	// UserData to apply to the VM
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// credentialsSecret is a local reference to a secret that contains the
	// credentials data to access Nutanix PC client
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NutanixMachineProviderConfigList contains a list of NutanixMachineProviderConfig
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type NutanixMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlibabaCloudMachineProviderConfig `json:"items"`
}

// NutanixMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains nutanix-specific status information.
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NutanixMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// vmUUID is the VM's UUID
	// +optional
	VmUUID *string `json:"vmUUID,omitempty"`

	// conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	// +optional
	Conditions []NutanixMachineProviderCondition `json:"conditions,omitempty"`
}

// NutanixMachineProviderConditionType is a valid value for NutanixMachineProviderCondition.Type
type NutanixMachineProviderConditionType string

// Valid conditions for an Nutanix VM instance.
const (
	// NutanixMachineCreation indicates whether the machine's VM has been created or not. If not,
	// it should include a reason and message for the failure.
	NutanixMachineCreation NutanixMachineProviderConditionType = "MachineCreation"

	// NutanixMachineUpdate indicates whether the machine's VM has been updated or not. If not,
	// it should include a reason and message for the failure.
	NutanixMachineUpdate NutanixMachineProviderConditionType = "MachineUpdate"

	// NutanixMachineDeletion indicates whether the machine's VM has been deleted or not. If not,
	// it should include a reason and message for the failure.
	NutanixMachineDeletion NutanixMachineProviderConditionType = "MachineDeletion"
)

// NutanixMachineProviderConditionReason is reason for the condition's last transition.
type NutanixMachineProviderConditionReason string

const (
	// NutanixMachineCreationSucceeded indicates machine creation success.
	NutanixMachineCreationSucceeded NutanixMachineProviderConditionReason = "MachineCreationSucceeded"
	// NutanixMachineCreationFailed indicates machine creation failure.
	NutanixMachineCreationFailed NutanixMachineProviderConditionReason = "MachineCreationFailed"
	// NutanixMachineCreationSucceeded indicates machine update success.
	NutanixMachineUpdateSucceeded NutanixMachineProviderConditionReason = "MachineUpdateSucceeded"
	// NutanixMachineCreationFailed indicates machine update failure.
	NutanixMachineUpdateFailed NutanixMachineProviderConditionReason = "MachineUpdateFailed"
	// NutanixMachineCreationSucceeded indicates machine deletion success.
	NutanixMachineDeletionSucceeded NutanixMachineProviderConditionReason = "MachineDeletionSucceeded"
	// NutanixMachineCreationFailed indicates machine deletion failure.
	NutanixMachineDeletionFailed NutanixMachineProviderConditionReason = "MachineDeletionFailed"
)

// NutanixMachineProviderCondition is a condition in a NutanixMachineProviderStatus.
type NutanixMachineProviderCondition struct {
	// Type is the type of the condition.
	Type NutanixMachineProviderConditionType `json:"type"`

	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`

	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason NutanixMachineProviderConditionReason `json:"reason,omitempty"`

	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// RawExtensionFromNutanixMachineProviderSpec marshals the machine provider spec.
func RawExtensionFromNutanixMachineProviderSpec(spec *NutanixMachineProviderConfig) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, fmt.Errorf("error marshalling providerSpec: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// RawExtensionFromNutanixMachineProviderStatus marshals the machine provider status
func RawExtensionFromNutanixMachineProviderStatus(status *NutanixMachineProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, fmt.Errorf("error marshalling providerStatus: %v", err)
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// NutanixMachineProviderSpecFromRawExtension unmarshals a raw extension into an NutanixMachineProviderConfig type
func NutanixMachineProviderSpecFromRawExtension(rawExtension *runtime.RawExtension) (*NutanixMachineProviderConfig, error) {
	if rawExtension == nil {
		return &NutanixMachineProviderConfig{}, nil
	}

	spec := new(NutanixMachineProviderConfig)
	if err := json.Unmarshal(rawExtension.Raw, &spec); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerSpec: %v", err)
	}

	klog.V(5).Infof("Got provider Spec from raw extension: %+v", spec)
	return spec, nil
}

// NutanixMachineProviderStatusFromRawExtension unmarshals a raw extension into an NutanixMachineProviderStatus type
func NutanixMachineProviderStatusFromRawExtension(rawExtension *runtime.RawExtension) (*NutanixMachineProviderStatus, error) {
	if rawExtension == nil {
		return &NutanixMachineProviderStatus{}, nil
	}

	providerStatus := new(NutanixMachineProviderStatus)
	if err := json.Unmarshal(rawExtension.Raw, providerStatus); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerStatus: %v", err)
	}

	klog.V(5).Infof("Got provider Status from raw extension: %+v", providerStatus)
	return providerStatus, nil
}

// The expected resource description and category (key/value) for the Nutanix resources (ex. vms, images)
// created for the cluster
/*const (
	NutanixExpectedResourceDescription = "Created By OpenShift Installer"
	NutanixExpectedCategoryKeyPrefix   = "openshift-" // format: "openshift-<cluster-id>"
	NutanixExpectedCategoryValue       = "openshift-ipi-installations"
)

// NutanixExpectedCategory holds a category key/value for the Nutanix resources (ex. vms, images)
type NutanixExpectedCategory struct {
	Key   string
	Value string
}

func CreateNutanixExpectedCategory(infraID string) *NutanixExpectedCategory {
	return &NutanixExpectedCategory{
		Key:   fmt.Sprintf("%s%s", NutanixExpectedCategoryKeyPrefix, infraID),
		Value: NutanixExpectedCategoryValue,
	}
}*/
