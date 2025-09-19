/*
Copyright 2024.

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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressControllerConfig is the Custom Resource object which holds the current configuration of Ingress Controllers.
// This provides a cluster-level configuration API for managing ingress controller operational settings.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/XXXX
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=IngressControllerConfig
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ingresscontrollerconfigs,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:metadata:annotations="description=Ingress Controller configuration API"
type IngressControllerConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user configuration for the Ingress Controllers
	// +required
	Spec IngressControllerConfigSpec `json:"spec,omitzero"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status *IngressControllerConfigStatus `json:"status,omitempty"`
}

// IngressControllerConfigStatus defines the observed state of IngressControllerConfig
type IngressControllerConfigStatus struct {
	// conditions represent the latest available observations of the IngressControllerConfig's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type IngressControllerConfigList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// items is a list of IngressControllerConfig
	// +optional
	Items []IngressControllerConfig `json:"items"`
}

// IngressControllerConfigSpec defines the desired state of Ingress Controller operational configuration
// +kubebuilder:validation:MinProperties=1
type IngressControllerConfigSpec struct {
	// defaultControllerConfig allows users to configure how the default ingress controller instance
	// should be deployed and managed.
	// defaultControllerConfig is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	// +optional
	DefaultControllerConfig DefaultIngressControllerConfig `json:"defaultControllerConfig,omitempty,omitzero"`

	// performanceTuning provides configuration options for performance optimization of ingress controllers.
	// performanceTuning is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	// +optional
	PerformanceTuning IngressControllerPerformanceTuning `json:"performanceTuning,omitempty,omitzero"`
}

// DefaultIngressControllerConfig represents the configuration for the default ingress controller deployment.
// defaultIngressControllerConfig provides configuration options for the default ingress controller instance
// that runs in the `openshift-ingress` namespace. Use this configuration to control
// how the default ingress controller is deployed, how it logs, and how its pods are scheduled.
// +kubebuilder:validation:MinProperties=1
type DefaultIngressControllerConfig struct {
	// logLevel defines the verbosity of logs emitted by the ingress controller.
	// This field allows users to control the amount and severity of logs generated, which can be useful
	// for debugging issues or reducing noise in production environments.
	// Allowed values are Error, Warn, Info, and Debug.
	// When set to Error, only errors will be logged.
	// When set to Warn, both warnings and errors will be logged.
	// When set to Info, general information, warnings, and errors will all be logged.
	// When set to Debug, detailed debugging information will be logged.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	// The current default value is `Info`.
	// +optional
	LogLevel IngressControllerLogLevel `json:"logLevel,omitempty"`

	// nodeSelector defines the nodes on which the ingress controller Pods are scheduled
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// resources defines the compute resource requests and limits for the ingress controller container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 100m
	//      limit: null
	//    - name: memory
	//      request: 256Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []IngressControllerContainerResource `json:"resources,omitempty"`

	// replicas defines the desired number of ingress controller replicas.
	// This field allows users to control the availability and load distribution of the ingress controller.
	// When not specified, defaults are used by the platform based on the cluster topology.
	// The current default behavior is:
	// - SingleReplica topology: 1 replica
	// - HighlyAvailable topology: 2 replicas
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=20
	Replicas int32 `json:"replicas,omitempty"`

	// tolerations defines the tolerations for ingress controller pods.
	// This allows the ingress controller to be scheduled on nodes with matching taints.
	// When not specified, no tolerations are applied.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=50
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// affinity defines the affinity rules for ingress controller pods.
	// This allows users to control pod placement for high availability or performance optimization.
	// When not specified, no affinity rules are applied.
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
}

// IngressControllerPerformanceTuning provides configuration options for performance optimization
// of ingress controllers. Use this configuration to control connection limits, timeouts,
// and other performance-related settings.
// +kubebuilder:validation:MinProperties=1
type IngressControllerPerformanceTuning struct {
	// connectionLimits defines limits on connections handled by the ingress controller.
	// connectionLimits is optional.
	// When omitted, this means no opinion and the platform is left to choose reasonable defaults.
	// +optional
	ConnectionLimits *IngressControllerConnectionLimits `json:"connectionLimits,omitempty"`

	// timeouts defines timeout settings for the ingress controller.
	// timeouts is optional.
	// When omitted, this means no opinion and the platform is left to choose reasonable defaults.
	// +optional
	Timeouts *IngressControllerTimeouts `json:"timeouts,omitempty"`

	// bufferSizes defines buffer size settings for the ingress controller.
	// bufferSizes is optional.
	// When omitted, this means no opinion and the platform is left to choose reasonable defaults.
	// +optional
	BufferSizes *IngressControllerBufferSizes `json:"bufferSizes,omitempty"`
}

// IngressControllerConnectionLimits defines connection-related limits for ingress controllers.
type IngressControllerConnectionLimits struct {
	// maxConnections defines the maximum number of concurrent connections.
	// This helps prevent resource exhaustion under high load.
	// When not specified, the platform default is used.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000000
	MaxConnections int32 `json:"maxConnections,omitempty"`

	// maxConnectionsPerBackend defines the maximum number of connections per backend server.
	// This helps distribute load evenly across backend servers.
	// When not specified, the platform default is used.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10000
	MaxConnectionsPerBackend int32 `json:"maxConnectionsPerBackend,omitempty"`

	// maxRequestsPerConnection defines the maximum number of requests per connection.
	// This controls connection reuse behavior.
	// When not specified, the platform default is used.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	MaxRequestsPerConnection int32 `json:"maxRequestsPerConnection,omitempty"`
}

// IngressControllerTimeouts defines timeout settings for ingress controllers.
type IngressControllerTimeouts struct {
	// clientTimeout defines the timeout for client connections.
	// This is the maximum time to wait for a client to send a request.
	// When not specified, the platform default is used.
	// +optional
	ClientTimeout *metav1.Duration `json:"clientTimeout,omitempty"`

	// serverTimeout defines the timeout for backend server connections.
	// This is the maximum time to wait for a response from a backend server.
	// When not specified, the platform default is used.
	// +optional
	ServerTimeout *metav1.Duration `json:"serverTimeout,omitempty"`

	// connectTimeout defines the timeout for establishing connections to backend servers.
	// This is the maximum time to wait when establishing a connection to a backend.
	// When not specified, the platform default is used.
	// +optional
	ConnectTimeout *metav1.Duration `json:"connectTimeout,omitempty"`
}

// IngressControllerBufferSizes defines buffer size settings for ingress controllers.
type IngressControllerBufferSizes struct {
	// requestHeaderBufferSize defines the size of the buffer for request headers.
	// This affects the maximum size of request headers that can be processed.
	// When not specified, the platform default is used.
	// +optional
	RequestHeaderBufferSize *resource.Quantity `json:"requestHeaderBufferSize,omitempty"`

	// responseBufferSize defines the size of the buffer for responses.
	// This affects buffering behavior for responses from backend servers.
	// When not specified, the platform default is used.
	// +optional
	ResponseBufferSize *resource.Quantity `json:"responseBufferSize,omitempty"`
}

// IngressControllerContainerResource defines a single resource requirement for an ingress controller container.
// +kubebuilder:validation:XValidation:rule="has(self.request) || has(self.limit)",message="at least one of request or limit must be set"
// +kubebuilder:validation:XValidation:rule="!(has(self.request) && has(self.limit)) || quantity(self.limit).compareTo(quantity(self.request)) >= 0",message="limit must be greater than or equal to request"
type IngressControllerContainerResource struct {
	// name of the resource (e.g. "cpu", "memory", "hugepages-2Mi").
	// This field is required.
	// name must consist only of alphanumeric characters, `-`, `_` and `.` and must start and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.qualifiedName().validate(self).hasValue()",message="name must consist only of alphanumeric characters, `-`, `_` and `.` and must start and end with an alphanumeric character"
	Name string `json:"name,omitempty"`

	// request is the minimum amount of the resource required (e.g. "2Mi", "1Gi").
	// This field is optional.
	// When limit is specified, request cannot be greater than limit.
	// +optional
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="isQuantity(self) && quantity(self).isGreaterThan(quantity('0'))",message="request must be a positive, non-zero quantity"
	Request resource.Quantity `json:"request,omitempty"`

	// limit is the maximum amount of the resource allowed (e.g. "2Mi", "1Gi").
	// This field is optional.
	// When request is specified, limit cannot be less than request.
	// The value must be greater than 0 when specified.
	// +optional
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="isQuantity(self) && quantity(self).isGreaterThan(quantity('0'))",message="limit must be a positive, non-zero quantity"
	Limit resource.Quantity `json:"limit,omitempty"`
}

// IngressControllerLogLevel defines the log level for ingress controllers
// +kubebuilder:validation:Enum="Error";"Warn";"Info";"Debug"
type IngressControllerLogLevel string

const (
	// IngressControllerLogLevelError only errors will be logged.
	IngressControllerLogLevelError IngressControllerLogLevel = "Error"
	// IngressControllerLogLevelWarn, both warnings and errors will be logged.
	IngressControllerLogLevelWarn IngressControllerLogLevel = "Warn"
	// IngressControllerLogLevelInfo, general information, warnings, and errors will all be logged.
	IngressControllerLogLevelInfo IngressControllerLogLevel = "Info"
	// IngressControllerLogLevelDebug, detailed debugging information will be logged.
	IngressControllerLogLevelDebug IngressControllerLogLevel = "Debug"
)

