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

// TcProgramInfo defines the tc program details
type TcProgramInfo struct {
	// links is the list of points to which the program should be attached.  The list items
	// are optional and may be updated after the bpf program has been loaded
	// +optional
	// +kubebuilder:default:={}
	Links []TcAttachInfo `json:"links"`
}

type TcAttachInfo struct {
	// interfaceSelector to determine the network interface (or interfaces)
	InterfaceSelector InterfaceSelector `json:"interfaceSelector"`

	// containers identify the set of containers in which to attach the eBPF
	// program.
	Containers ContainerSelector `json:"containers"`

	// direction specifies the direction of traffic the tc program should
	// attach to for a given network device.
	// +kubebuilder:validation:Enum=Ingress;Egress
	Direction TCDirectionType `json:"direction"`

	// priority specifies the priority of the tc program in relation to
	// other programs of the same type with the same attach point. It is a value
	// from 0 to 1000 where lower values have higher precedence.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	Priority int32 `json:"priority"`

	// proceedOn allows the user to call other tc programs in chain on this exit code.
	// Multiple values are supported by repeating the parameter.
	// +optional
	// +kubebuilder:default:={Pipe,DispatcherReturn}
	ProceedOn []TcProceedOnValue `json:"proceedOn"`
}

type TcProgramInfoState struct {
	// links is the List of attach points for the BPF program on the given node. Each entry
	// in *AttachInfoState represents a specific, unique attach point that is
	// derived from *AttachInfo by fully expanding any selectors.  Each entry
	// also contains information about the attach point required by the
	// reconciler
	// +optional
	// +kubebuilder:default:={}
	Links []TcAttachInfoState `json:"links"`
}

type TcAttachInfoState struct {
	AttachInfoStateCommon `json:",inline"`

	// interfaceName is the Interface name to attach the tc program to.
	InterfaceName string `json:"interfaceName"`

	// containerPid is the Container pid to attach the tc program in.
	ContainerPid int32 `json:"containerPid"`

	// direction specifies the direction of traffic the tc program should
	// attach to for a given network device.
	// +kubebuilder:validation:Enum=Ingress;Egress
	Direction TCDirectionType `json:"direction"`

	// priority specifies the priority of the tc program in relation to
	// other programs of the same type with the same attach point. It is a value
	// from 0 to 1000 where lower values have higher precedence.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	Priority int32 `json:"priority"`

	// proceedOn allows the user to call other tc programs in chain on this exit code.
	// Multiple values are supported by repeating the parameter.
	ProceedOn []TcProceedOnValue `json:"proceedOn"`
}
