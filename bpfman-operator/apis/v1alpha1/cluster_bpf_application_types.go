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
)

// EBPFProgType defines the supported eBPF program types
type EBPFProgType string

const (
	// ProgTypeXDP refers to the XDP program type.
	ProgTypeXDP EBPFProgType = "xdp"

	// ProgTypeTC refers to the TC program type.
	ProgTypeTC EBPFProgType = "tc"

	// ProgTypeTCX refers to the TCX program type.
	ProgTypeTCX EBPFProgType = "tcx"

	// ProgTypeFentry refers to the Fentry program type.
	ProgTypeFentry EBPFProgType = "fentry"

	// ProgTypeFexit refers to the Fexit program type.
	ProgTypeFexit EBPFProgType = "fexit"

	// ProgTypeKprobe refers to the Kprobe program type.
	ProgTypeKprobe EBPFProgType = "kprobe"

	// ProgTypeKretprobe refers to the Kretprobe program type.
	ProgTypeKretprobe EBPFProgType = "kretprobe"

	// ProgTypeUprobe refers to the Uprobe program type.
	ProgTypeUprobe EBPFProgType = "uprobe"

	// ProgTypeUretprobe refers to the Uretprobe program type.
	ProgTypeUretprobe EBPFProgType = "uretprobe"

	// ProgTypeTracepoint refers to the Tracepoint program type.
	ProgTypeTracepoint EBPFProgType = "tracepoint"
)

// ClBpfApplicationProgram defines the desired state of BpfApplication
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
type ClBpfApplicationProgram struct {
	// Name is the name of the function that is the entry point for the BPF
	// program
	Name string `json:"name"`

	// Type specifies the bpf program type
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:="xdp";"tc";"tcx";"fentry";"fexit";"kprobe";"kretprobe";"uprobe";"uretprobe";"tracepoint"
	Type EBPFProgType `json:"type,omitempty"`

	// xdp defines the desired state of the application's XdpPrograms.
	// +unionMember
	// +optional
	XDPInfo *ClXdpProgramInfo `json:"xdpInfo,omitempty"`

	// tc defines the desired state of the application's TcPrograms.
	// +unionMember
	// +optional
	TCInfo *ClTcProgramInfo `json:"tcInfo,omitempty"`

	// tcx defines the desired state of the application's TcxPrograms.
	// +unionMember
	// +optional
	TCXInfo *ClTcxProgramInfo `json:"tcxInfo,omitempty"`

	// fentry defines the desired state of the application's FentryPrograms.
	// +unionMember
	// +optional
	FentryInfo *ClFentryProgramInfo `json:"fentryInfo,omitempty"`

	// fexit defines the desired state of the application's FexitPrograms.
	// +unionMember
	// +optional
	FexitInfo *ClFexitProgramInfo `json:"fexitInfo,omitempty"`

	// kprobe defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KprobeInfo *ClKprobeProgramInfo `json:"kprobeInfo,omitempty"`

	// kretprobe defines the desired state of the application's KretprobePrograms.
	// +unionMember
	// +optional
	KretprobeInfo *ClKprobeProgramInfo `json:"kretprobeInfo,omitempty"`

	// uprobe defines the desired state of the application's UprobePrograms.
	// +unionMember
	// +optional
	UprobeInfo *ClUprobeProgramInfo `json:"uprobeInfo,omitempty"`

	// uretprobe defines the desired state of the application's UretprobePrograms.
	// +unionMember
	// +optional
	UretprobeInfo *ClUprobeProgramInfo `json:"uretprobeInfo,omitempty"`

	// tracepoint defines the desired state of the application's TracepointPrograms.
	// +unionMember
	// +optional
	TracepointInfo *ClTracepointProgramInfo `json:"tracepointInfo,omitempty"`
}

// ClBpfApplicationSpec defines the desired state of BpfApplication
type ClBpfApplicationSpec struct {
	BpfAppCommon `json:",inline"`
	// Programs is the list of bpf programs in the BpfApplication that should be
	// loaded. The application can selectively choose which program(s) to run
	// from this list based on the optional attach points provided.
	// +kubebuilder:validation:MinItems:=1
	Programs []ClBpfApplicationProgram `json:"programs,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// ClusterBpfApplication is the Schema for the bpfapplications API
// +kubebuilder:printcolumn:name="NodeSelector",type=string,JSONPath=`.spec.nodeselector`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[0].reason`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClusterBpfApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClBpfApplicationSpec `json:"spec,omitempty"`
	Status BpfAppStatus         `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ClusterBpfApplicationList contains a list of BpfApplications
type ClusterBpfApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterBpfApplication `json:"items"`
}
