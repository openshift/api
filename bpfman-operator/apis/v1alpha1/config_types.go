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

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Config holds the configuration for bpfman-operator.
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Progressing",type="string",JSONPath=".status.conditions[?(@.type=='Progressing')].status"
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.conditions[?(@.type=='Available')].status"
type Config struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the desired state of the bpfman-operator.
	// +required
	Spec ConfigSpec `json:"spec,omitzero"`
	// status reflects the observed state of the bpfman-operator.
	// +optional
	Status ConfigStatus `json:"status,omitempty"`
}

// spec defines the desired state of the bpfman-operator.
type ConfigSpec struct {
	// agent specifies the configuration for the bpfman agent DaemonSet.
	// +required
	Agent AgentSpec `json:"agent,omitzero"`
	// configuration specifies the content of bpfman.toml configuration file used by the bpfman DaemonSet.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=65536
	Configuration string `json:"configuration,omitempty"`
	// image specifies the container image for the bpfman DaemonSet.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1023
	Image string `json:"image,omitempty"`
	// logLevel specifies the log level for the bpfman DaemonSet via the RUST_LOG environment variable.
	// The RUST_LOG environment variable controls logging with the syntax: RUST_LOG=[target][=][level][,...].
	// For further information, see https://docs.rs/env_logger/latest/env_logger/.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	LogLevel string `json:"logLevel,omitempty"`
	// namespace specifies the namespace where bpfman-operator resources will be deployed.
	// If not specified, resources will be deployed in the default bpfman namespace.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Namespace string `json:"namespace,omitempty"`
}

// AgentSpec defines the desired state of the bpfman agent.
type AgentSpec struct {
	// healthProbePort specifies the port on which the bpfman agent's health probe endpoint will listen.
	// If unspecified, the default port will be used.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	HealthProbePort int32 `json:"healthProbePort,omitempty"`
	// image specifies the container image for the bpfman agent DaemonSet.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1023
	Image string `json:"image,omitempty"`
	// logLevel specifies the verbosity of logs produced by the bpfman agent.
	// Valid values are: "info", "debug", "trace".
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Enum=info;debug;trace
	LogLevel string `json:"logLevel,omitempty"`
}

// status reflects the status of the bpfman-operator configuration.
type ConfigStatus struct {
	// conditions represents the current state conditions of the bpfman-operator and its components.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=1023
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`

	// components represents the operational status of each individual bpfman-operator component such as the deployed
	// DaemonSets.
	// +optional
	Components map[string]ConfigComponentStatus `json:"components,omitempty"`
}

// +kubebuilder:object:root=true
// ConfigList contains a list of Configs.
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

// ConfigComponentStatus holds the status of a single Config component.
type ConfigComponentStatus string

const (
	// ConfigStatusUnknown indicates the component state cannot be determined.
	ConfigStatusUnknown ConfigComponentStatus = "Unknown"
	// ConfigStatusProgressing indicates the component is being updated or reconciled.
	ConfigStatusProgressing ConfigComponentStatus = "Progressing"
	// ConfigStatusReady indicates the component is fully operational and ready.
	ConfigStatusReady ConfigComponentStatus = "Ready"
)
