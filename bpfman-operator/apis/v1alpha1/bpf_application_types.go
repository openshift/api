/*
Copyright 2024 The bpfman Authors.

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

// BpfApplicationProgram defines the desired state of BpfApplication
// +union
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'XDP' ?  has(self.xdpInfo) : !has(self.xdpInfo)",message="xdpInfo configuration is required when type is xdp, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'TC' ?  has(self.tcInfo) : !has(self.tcInfo)",message="tcInfo configuration is required when type is tc, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'TCX' ?  has(self.tcxInfo) : !has(self.tcxInfo)",message="tcxInfo configuration is required when type is tcx, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Uprobe' ?  has(self.uprobeInfo) : !has(self.uprobeInfo)",message="uprobeInfo configuration is required when type is uprobe, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'UretProbe' ?  has(self.uretprobeInfo) : !has(self.upretrobeInfo)",message="uretprobe configuration is required when type is uretprobe, and forbidden otherwise"
type BpfApplicationProgram struct {
	// name is the name of the function that is the entry point for the BPF
	// program
	Name string `json:"name"`

	// type specifies the bpf program type
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum:="XDP";"TC";"TCX";"Uprobe";"UretProbe"
	Type EBPFProgType `json:"type"`

	// xdpInfo defines the desired state of the application's XdpPrograms.
	// +unionMember
	// +optional
	XDPInfo *XdpProgramInfo `json:"xdpInfo,omitempty"`

	// tcInfo defines the desired state of the application's TcPrograms.
	// +unionMember
	// +optional
	TCInfo *TcProgramInfo `json:"tcInfo,omitempty"`

	// tcxInfo defines the desired state of the application's TcxPrograms.
	// +unionMember
	// +optional
	TCXInfo *TcxProgramInfo `json:"tcxInfo,omitempty"`

	// uprobeInfo defines the desired state of the application's UprobePrograms.
	// +unionMember
	// +optional
	UprobeInfo *UprobeProgramInfo `json:"uprobeInfo,omitempty"`

	// uretprobeInfo defines the desired state of the application's UretprobePrograms.
	// +unionMember
	// +optional
	UretprobeInfo *UprobeProgramInfo `json:"uretprobeInfo,omitempty"`
}

// BpfApplicationSpec defines the desired state of BpfApplication
type BpfApplicationSpec struct {
	BpfAppCommon `json:",inline"`

	// programs is the list of bpf programs in the BpfApplication that should be
	// loaded. The application can selectively choose which program(s) to run
	// from this list based on the optional attach points provided.
	// +kubebuilder:validation:MinItems:=1
	Programs []BpfApplicationProgram `json:"programs,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// BpfApplication is the Schema for the bpfapplications API
// +kubebuilder:printcolumn:name="NodeSelector",type=string,JSONPath=`.spec.nodeselector`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[0].reason`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type BpfApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BpfApplicationSpec `json:"spec,omitempty"`
	Status BpfAppStatus       `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// BpfApplicationList contains a list of BpfApplications
type BpfApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BpfApplication `json:"items"`
}
