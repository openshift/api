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
type ClBpfApplicationProgram struct {
	// name is the name of the function that is the entry point for the BPF
	// program
	Name string `json:"name"`

	// type specifies the bpf program type
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum:="XDP";"TC";"TCX";"Fentry";"Fexit";"Kprobe";"KretProbe";"Uprobe";"UretProbe";"TracePoint"
	Type EBPFProgType `json:"type"`

	// xdpInfo defines the desired state of the application's XdpPrograms.
	// +unionMember
	// +optional
	XDPInfo *ClXdpProgramInfo `json:"xdpInfo,omitempty"`

	// tcInfo defines the desired state of the application's TcPrograms.
	// +unionMember
	// +optional
	TCInfo *ClTcProgramInfo `json:"tcInfo,omitempty"`

	// tcxInfo defines the desired state of the application's TcxPrograms.
	// +unionMember
	// +optional
	TCXInfo *ClTcxProgramInfo `json:"tcxInfo,omitempty"`

	// fentryInfo defines the desired state of the application's FentryPrograms.
	// +unionMember
	// +optional
	FentryInfo *ClFentryProgramInfo `json:"fentryInfo,omitempty"`

	// fexitInfo defines the desired state of the application's FexitPrograms.
	// +unionMember
	// +optional
	FexitInfo *ClFexitProgramInfo `json:"fexitInfo,omitempty"`

	// kprobeInfo defines the desired state of the application's KprobePrograms.
	// +unionMember
	// +optional
	KprobeInfo *ClKprobeProgramInfo `json:"kprobeInfo,omitempty"`

	// kretprobeInfo defines the desired state of the application's KretprobePrograms.
	// +unionMember
	// +optional
	KretprobeInfo *ClKprobeProgramInfo `json:"kretprobeInfo,omitempty"`

	// uprobeInfo defines the desired state of the application's UprobePrograms.
	// +unionMember
	// +optional
	UprobeInfo *ClUprobeProgramInfo `json:"uprobeInfo,omitempty"`

	// uretprobeInfo defines the desired state of the application's UretprobePrograms.
	// +unionMember
	// +optional
	UretprobeInfo *ClUprobeProgramInfo `json:"uretprobeInfo,omitempty"`

	// tracepointInfo defines the desired state of the application's TracepointPrograms.
	// +unionMember
	// +optional
	TracepointInfo *ClTracepointProgramInfo `json:"tracepointInfo,omitempty"`
}

// ClBpfApplicationSpec defines the desired state of BpfApplication
type ClBpfApplicationSpec struct {
	BpfAppCommon `json:",inline"`
	// programs is the list of bpf programs in the BpfApplication that should be
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
