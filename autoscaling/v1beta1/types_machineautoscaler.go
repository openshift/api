package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoscaler targets machine-api managed machines for autoscaling.
//
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ref Kind",type="string",JSONPath=".spec.scaleTargetRef.kind",description="Kind of object scaled"
// +kubebuilder:printcolumn:name="Ref Name",type="string",JSONPath=".spec.scaleTargetRef.name",description="Name of object scaled"
// +kubebuilder:printcolumn:name="Min",type="integer",JSONPath=".spec.minReplicas",description="Min number of replicas"
// +kubebuilder:printcolumn:name="Max",type="integer",JSONPath=".spec.maxReplicas",description="Max number of replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="MachineAutoscaler resoruce age"
type MachineAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of constraints of a scalable resource
	Spec MachineAutoscalerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of a scalable resource
	Status MachineAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MachineAutoscalerSpec represents the desired state of a MachineAutoscaler.
type MachineAutoscalerSpec struct {
	// MinReplicas constrains the minimal number of replicas of a scalable resource
	// +kubebuilder:validation:Minimum=0
	MinReplicas int32 `json:"minReplicas" protobuf:"varint,1,opt,name=minReplicas"`

	// MaxReplicas constrains the maximal number of replicas of a scalable resource
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int32 `json:"maxReplicas" protobuf:"varint,2,opt,name=maxReplicas"`

	// ScaleTargetRef holds reference to a scalable resource
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef" protobuf:"bytes,3,opt,name=scaleTargetRef"`
}

// CrossVersionObjectReference identifies another object by name, API version,
// and kind.
type CrossVersionObjectReference struct {
	// APIVersion defines the versioned schema of this representation of an
	// object. Servers should convert recognized schemas to the latest internal
	// value, and may reject unrecognized values. More info:
	// http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#resources
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,1,opt,name=apiVersion"`

	// Kind is a string value representing the REST resource this object
	// represents. Servers may infer this from the endpoint the client submits
	// requests to. Cannot be updated. In CamelCase. More info:
	// http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#types-kinds
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind" protobuf:"bytes,2,opt,name=kind"`

	// Name specifies a name of an object, e.g. worker-us-east-1a.
	// Scalable resources are expected to exist under a single namespace.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
}

// MachineAutoscalerStatus defines the observed state of a MachineAutoscaler.
type MachineAutoscalerStatus struct {
	// LastTargetRef holds reference to the recently observed scalable resource
	LastTargetRef *CrossVersionObjectReference `json:"lastTargetRef,omitempty" protobuf:"bytes,1,opt,name=lastTargetRef"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoscalerList contains a list of MachineAutoscaler resources.
type MachineAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MachineAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
