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
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'XDP' ?  has(self.xdpInfo) : !has(self.xdpInfo)",message="xdpInfo configuration is required when type is xdp, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'TC' ?  has(self.tcInfo) : !has(self.tcInfo)",message="tcInfo configuration is required when type is tc, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'TCX' ?  has(self.tcxInfo) : !has(self.tcxInfo)",message="tcxInfo configuration is required when type is tcx, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Fentry' ?  has(self.fentryInfo) : !has(self.fentryInfo)",message="fentryInfo configuration is required when type is fentry, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Fexit' ?  has(self.fexitInfo) : !has(self.fexitInfo)",message="fexitInfo configuration is required when type is fexit, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Kprobe' ?  has(self.kprobeInfo) : !has(self.kprobeInfo)",message="kprobeInfo configuration is required when type is kprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'KretProbe' ?  has(self.kretprobeInfo) : !has(self.kretprobeInfo)",message="kretprobeInfo configuration is required when type is kretprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Uprobe' ?  has(self.uprobeInfo) : !has(self.uprobeInfo)",message="uprobeInfo configuration is required when type is uprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'UretProbe' ?  has(self.uretprobeInfo) : !has(self.uretprobeInfo)",message="uretprobeInfo configuration is required when type is uretprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'TracePoint' ?  has(self.tracepointInfo) : !has(self.tracepointInfo)",message="tracepointInfo configuration is required when type is tracepoint, and forbidden otherwise"
type ClBpfApplicationProgramState struct {
	BpfProgramStateCommon `json:",inline"`
	// type specifies the bpf program type
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum:="XDP";"TC";"TCX";"Fentry";"Fexit";"Kprobe";"KretProbe";"Uprobe";"UretProbe";"TracePoint"

	Type EBPFProgType `json:"type"`

	// xdpInfo defines the desired state of the application's XdpPrograms.
	// +unionMember
	// +optional
	XDPInfo *ClXdpProgramInfoState `json:"xdpInfo,omitempty"`

	// tcInfo defines the desired state of the application's TcPrograms.
	// +unionMember
	// +optional
	TCInfo *ClTcProgramInfoState `json:"tcInfo,omitempty"`

	// tcxInfo defines the desired state of the application's TcxPrograms.
	// +unionMember
	// +optional
	TCXInfo *ClTcxProgramInfoState `json:"tcxInfo,omitempty"`

	// fentryInfo defines the desired state of the application's FentryPrograms.
	// +unionMember
	// +optional
	FentryInfo *ClFentryProgramInfoState `json:"fentryInfo,omitempty"`

	// fexitInfo defines the desired state of the application's FexitPrograms.
	// +unionMember
	// +optional
	FexitInfo *ClFexitProgramInfoState `json:"fexitInfo,omitempty"`

	// kprobeInfo defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KprobeInfo *ClKprobeProgramInfoState `json:"kprobeInfo,omitempty"`

	// kretprobeInfo defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KretprobeInfo *ClKprobeProgramInfoState `json:"kretprobeInfo,omitempty"`

	// uprobeInfo defines the desired state of the application's UprobePrograms.
	// +unionMember
	// +optional
	UprobeInfo *ClUprobeProgramInfoState `json:"uprobeInfo,omitempty"`

	// uretprobeInfo defines the desired state of the application's UretprobePrograms.
	// +unionMember
	// +optional
	UretprobeInfo *ClUprobeProgramInfoState `json:"uretprobeInfo,omitempty"`

	// tracepointInfo defines the desired state of the application's TracepointPrograms.
	// +unionMember
	// +optional
	TracepointInfo *ClTracepointProgramInfoState `json:"tracepointInfo,omitempty"`
}

// BpfApplicationSpec defines the desired state of BpfApplication
type ClBpfApplicationStateSpec struct {
	// node is the name of the node for this BpfApplicationStateSpec.
	Node string `json:"node"`
	// updateCount is the number of times the BpfApplicationState has been updated. Set to 1
	// when the object is created, then it is incremented prior to each update.
	// This allows us to verify that the API server has the updated object prior
	// to starting a new Reconcile operation.
	UpdateCount int64 `json:"updateCount"`
	// appLoadStatus reflects the status of loading the bpf application on the
	// given node.
	AppLoadStatus AppLoadStatus `json:"appLoadStatus"`
	// programs is a list of bpf programs contained in the parent application.
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
