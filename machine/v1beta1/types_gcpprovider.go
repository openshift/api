package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an GCP virtual machine. It is used by the GCP machine actuator to create a single Machine.
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type GCPMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty" protobuf:"bytes,2,opt,name=userDataSecret"`

	// CredentialsSecret is a reference to the secret with GCP credentials.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty" protobuf:"bytes,3,opt,name=credentialsSecret"`

	CanIPForward       bool                   `json:"canIPForward" protobuf:"varint,4,opt,name=canIPForward"`
	DeletionProtection bool                   `json:"deletionProtection" protobuf:"varint,5,opt,name=deletionProtection"`
	Disks              []*GCPDisk             `json:"disks,omitempty" protobuf:"bytes,6,rep,name=disks"`
	Labels             map[string]string      `json:"labels,omitempty" protobuf:"bytes,7,rep,name=labels"`
	Metadata           []*GCPMetadata         `json:"gcpMetadata,omitempty" protobuf:"bytes,8,rep,name=gcpMetadata"`
	NetworkInterfaces  []*GCPNetworkInterface `json:"networkInterfaces,omitempty" protobuf:"bytes,9,rep,name=networkInterfaces"`
	ServiceAccounts    []GCPServiceAccount    `json:"serviceAccounts" protobuf:"bytes,10,rep,name=serviceAccounts"`
	Tags               []string               `json:"tags,omitempty" protobuf:"bytes,11,rep,name=tags"`
	TargetPools        []string               `json:"targetPools,omitempty" protobuf:"bytes,12,rep,name=targetPools"`
	MachineType        string                 `json:"machineType" protobuf:"bytes,13,opt,name=machineType"`
	Region             string                 `json:"region" protobuf:"bytes,14,opt,name=region"`
	Zone               string                 `json:"zone" protobuf:"bytes,15,opt,name=zone"`
	ProjectID          string                 `json:"projectID,omitempty" protobuf:"bytes,16,opt,name=projectID"`

	// Preemptible indicates if created instance is preemptible
	Preemptible bool `json:"preemptible,omitempty" protobuf:"varint,17,opt,name=preemptible"`
}

// GCPDisk describes disks for GCP.
type GCPDisk struct {
	AutoDelete    bool                       `json:"autoDelete" protobuf:"varint,1,opt,name=autoDelete"`
	Boot          bool                       `json:"boot" protobuf:"varint,2,opt,name=boot"`
	SizeGb        int64                      `json:"sizeGb" protobuf:"varint,3,opt,name=sizeGb"`
	Type          string                     `json:"type" protobuf:"bytes,4,opt,name=type"`
	Image         string                     `json:"image" protobuf:"bytes,5,opt,name=image"`
	Labels        map[string]string          `json:"labels" protobuf:"bytes,6,rep,name=labels"`
	EncryptionKey *GCPEncryptionKeyReference `json:"encryptionKey,omitempty" protobuf:"bytes,7,opt,name=encryptionKey"`
}

// GCPMetadata describes metadata for GCP.
type GCPMetadata struct {
	Key   string  `json:"key" protobuf:"bytes,1,opt,name=key"`
	Value *string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

// GCPNetworkInterface describes network interfaces for GCP
type GCPNetworkInterface struct {
	PublicIP   bool   `json:"publicIP,omitempty" protobuf:"varint,1,opt,name=publicIP"`
	Network    string `json:"network,omitempty" protobuf:"bytes,2,opt,name=network"`
	ProjectID  string `json:"projectID,omitempty" protobuf:"bytes,3,opt,name=projectID"`
	Subnetwork string `json:"subnetwork,omitempty" protobuf:"bytes,4,opt,name=subnetwork"`
}

// GCPServiceAccount describes service accounts for GCP.
type GCPServiceAccount struct {
	Email  string   `json:"email" protobuf:"bytes,1,opt,name=email"`
	Scopes []string `json:"scopes" protobuf:"bytes,2,rep,name=scopes"`
}

// GCPEncryptionKeyReference describes the encryptionKey to use for a disk's encryption.
type GCPEncryptionKeyReference struct {
	KMSKey *GCPKMSKeyReference `json:"kmsKey,omitempty" protobuf:"bytes,1,opt,name=kmsKey"`

	// KMSKeyServiceAccount is the service account being used for the
	// encryption request for the given KMS key. If absent, the Compute
	// Engine default service account is used.
	// See https://cloud.google.com/compute/docs/access/service-accounts#compute_engine_service_account
	// for details on the default service account.
	KMSKeyServiceAccount string `json:"kmsKeyServiceAccount,omitempty" protobuf:"bytes,2,opt,name=kmsKeyServiceAccount"`
}

// GCPKMSKeyReference gathers required fields for looking up a GCP KMS Key
type GCPKMSKeyReference struct {
	// Name is the name of the customer managed encryption key to be used for the disk encryption.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// KeyRing is the name of the KMS Key Ring which the KMS Key belongs to.
	KeyRing string `json:"keyRing" protobuf:"bytes,2,opt,name=keyRing"`

	// ProjectID is the ID of the Project in which the KMS Key Ring exists.
	// Defaults to the VM ProjectID if not set.
	ProjectID string `json:"projectID,omitempty" protobuf:"bytes,3,opt,name=projectID"`

	// Location is the GCP location in which the Key Ring exists.
	Location string `json:"location" protobuf:"bytes,4,opt,name=location"`
}

// GCPMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains GCP-specific status information.
// Compatibility level 2: Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=2
type GCPMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// InstanceID is the ID of the instance in GCP
	// +optional
	InstanceID *string `json:"instanceId,omitempty" protobuf:"bytes,2,opt,name=instanceId"`

	// InstanceState is the provisioning state of the GCP Instance.
	// +optional
	InstanceState *string `json:"instanceState,omitempty" protobuf:"bytes,3,opt,name=instanceState"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []GCPMachineProviderCondition `json:"conditions,omitempty" protobuf:"bytes,4,rep,name=conditions"`
}

// GCPMachineProviderCondition is a condition in a GCPMachineProviderStatus
type GCPMachineProviderCondition struct {
	// Type is the type of the condition.
	Type ConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=GCPMachineProviderConditionType"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}
