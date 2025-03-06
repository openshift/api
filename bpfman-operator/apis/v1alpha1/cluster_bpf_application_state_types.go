/*
Copyright 2023 The bpfman Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClBpfApplicationProgramState defines the desired state of BpfApplication
// +union
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'xdp' ?  has(self.xdpInfo) : !has(self.xdpInfo)",message="xdpInfo configuration is required when type is xdp, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'tc' ?  has(self.tcInfo) : !has(self.tcInfo)",message="tcInfo configuration is required when type is tc, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'tcx' ?  has(self.tcxInfo) : !has(self.tcxInfo)",message="tcxInfo configuration is required when type is tcx, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'fentry' ?  has(self.fentryInfo) : !has(self.fentryInfo)",message="fentryInfo configuration is required when type is fentry, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'fexit' ?  has(self.fexitInfo) : !has(self.fexitInfo)",message="fexitInfo configuration is required when type is fexit, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'kprobe' ?  has(self.kprobeInfo) : !has(self.kprobeInfo)",message="kprobeInfo configuration is required when type is kprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'kretprobe' ?  has(self.kretprobeInfo) : !has(self.kretprobeInfo)",message="kretprobeInfo configuration is required when type is kretprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'uprobe' ?  has(self.uprobeInfo) : !has(self.uprobeInfo)",message="uprobeInfo configuration is required when type is uprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'uretprobe' ?  has(self.uretprobeInfo) : !has(self.uretprobeInfo)",message="uretprobeInfo configuration is required when type is uretprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'tracepoint' ?  has(self.tracepointInfo) : !has(self.tracepointInfo)",message="tracepointInfo configuration is required when type is tracepoint, and forbidden otherwise"
type ClBpfApplicationProgramState struct {
	BpfProgramStateCommon `json:",inline"`
	// Type specifies the bpf program type
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:="xdp";"tc";"tcx";"fentry";"fexit";"kprobe";"kretprobe";"uprobe";"uretprobe";"tracepoint"

	Type EBPFProgType `json:"type,omitempty"`

	// xdp defines the desired state of the application's XdpPrograms.
	// +unionMember
	// +optional
	XDPInfo *ClXdpProgramInfoState `json:"xdpInfo,omitempty"`

	// tc defines the desired state of the application's TcPrograms.
	// +unionMember
	// +optional
	TCInfo *ClTcProgramInfoState `json:"tcInfo,omitempty"`

	// tcx defines the desired state of the application's TcxPrograms.
	// +unionMember
	// +optional
	TCXInfo *ClTcxProgramInfoState `json:"tcxInfo,omitempty"`

	// fentry defines the desired state of the application's FentryPrograms.
	// +unionMember
	// +optional
	FentryInfo *ClFentryProgramInfoState `json:"fentryInfo,omitempty"`

	// fexit defines the desired state of the application's FexitPrograms.
	// +unionMember
	// +optional
	FexitInfo *ClFexitProgramInfoState `json:"fexitInfo,omitempty"`

	// kprobe defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KprobeInfo *ClKprobeProgramInfoState `json:"kprobeInfo,omitempty"`

	// kprobe defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KretprobeInfo *ClKprobeProgramInfoState `json:"kretprobeInfo,omitempty"`

	// uprobe defines the desired state of the application's UprobePrograms.
	// +unionMember
	// +optional
	UprobeInfo *ClUprobeProgramInfoState `json:"uprobeInfo,omitempty"`

	// uretprobe defines the desired state of the application's UretprobePrograms.
	// +unionMember
	// +optional
	UretprobeInfo *ClUprobeProgramInfoState `json:"uretprobeInfo,omitempty"`

	// tracepoint defines the desired state of the application's TracepointPrograms.
	// +unionMember
	// +optional
	TracepointInfo *ClTracepointProgramInfoState `json:"tracepointInfo,omitempty"`
}

// BpfApplicationSpec defines the desired state of BpfApplication
type ClBpfApplicationStateSpec struct {
	// Node is the name of the node for this BpfApplicationStateSpec.
	Node string `json:"node"`
	// The number of times the BpfApplicationState has been updated.  Set to 1
	// when the object is created, then it is incremented prior to each update.
	// This allows us to verify that the API server has the updated object prior
	// to starting a new Reconcile operation.
	UpdateCount int64 `json:"updateCount"`
	// AppLoadStatus reflects the status of loading the bpf application on the
	// given node.
	AppLoadStatus AppLoadStatus `json:"appLoadStatus"`
	// Programs is a list of bpf programs contained in the parent application.
	// It is a map from the bpf program name to BpfApplicationProgramState
	// elements.
	Programs []ClBpfApplicationProgramState `json:"programs,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// ClusterBpfApplicationState contains the per-node state of a BpfApplication.
// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=".spec.node"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[0].reason`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClusterBpfApplicationState struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClBpfApplicationStateSpec `json:"spec,omitempty"`
	Status BpfAppStatus              `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ClusterBpfApplicationStateList contains a list of BpfApplicationState objects
type ClusterBpfApplicationStateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterBpfApplicationState `json:"items"`
}

func (an ClusterBpfApplicationState) GetName() string {
	return an.Name
}

func (an ClusterBpfApplicationState) GetUID() metav1types.UID {
	return an.UID
}

func (an ClusterBpfApplicationState) GetAnnotations() map[string]string {
	return an.Annotations
}

func (an ClusterBpfApplicationState) GetLabels() map[string]string {
	return an.Labels
}

func (an ClusterBpfApplicationState) GetStatus() *BpfAppStatus {
	return &an.Status
}

func (an ClusterBpfApplicationState) GetClientObject() client.Object {
	return &an
}

func (anl ClusterBpfApplicationStateList) GetItems() []ClusterBpfApplicationState {
	return anl.Items
}
