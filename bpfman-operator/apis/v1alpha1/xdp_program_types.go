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

// All fields are required unless explicitly marked optional
package v1alpha1

// XdpProgramInfo contains the xdp program details
type XdpProgramInfo struct {
	// links is the list of points to which the program should be attached.  The list items
	// are optional and may be udated after the bpf program has been loaded
	// +optional
	// +kubebuilder:default:={}
	Links []XdpAttachInfo `json:"links"`
}

type XdpAttachInfo struct {
	// interfaceSelector to determine the network interface (or interfaces)
	InterfaceSelector InterfaceSelector `json:"interfaceSelector"`

	// containers identify the set of containers in which to attach the eBPF
	// program.
	Containers ContainerSelector `json:"containers"`

	// priority specifies the priority of the bpf program in relation to
	// other programs of the same type with the same attach point. It is a value
	// from 0 to 1000 where lower values have higher precedence.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	Priority int32 `json:"priority"`

	// proceedOn allows the user to call other xdp programs in chain on this exit code.
	// Multiple values are supported by repeating the parameter.
	// +optional
	// +kubebuilder:default:={Pass,DispatcherReturn}
	ProceedOn []XdpProceedOnValue `json:"proceedOn"`
}

type XdpProgramInfoState struct {
	// links is the list of attach points for the BPF program on the given node. Each entry
	// in *AttachInfoState represents a specific, unique attach point that is
	// derived from *AttachInfo by fully expanding any selectors.  Each entry
	// also contains information about the attach point required by the
	// reconciler
	// +optional
	// +kubebuilder:default:={}
	Links []XdpAttachInfoState `json:"links"`
}

type XdpAttachInfoState struct {
	AttachInfoStateCommon `json:",inline"`

	// interfaceName is the interface name to attach the xdp program to.
	InterfaceName string `json:"interfaceName"`

	// containerPid Container pid to attach the xdp program in.
	ContainerPid int32 `json:"containerPid"`

	// priority specifies the priority of the xdp program in relation to
	// other programs of the same type with the same attach point. It is a value
	// from 0 to 1000 where lower values have higher precedence.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	Priority int32 `json:"priority"`

	// proceedOn allows the user to call other xdp programs in chain on this exit code.
	// Multiple values are supported by repeating the parameter.
	ProceedOn []XdpProceedOnValue `json:"proceedOn"`
}
