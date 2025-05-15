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
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterMonitoring is the Custom Resource object which holds the current status of Cluster Monitoring Operator. CMO is a central component of the monitoring stack.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:internal
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/1929
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clustermonitoring,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:metadata:annotations="description=Cluster Monitoring Operators configuration API"
// +openshift:enable:FeatureGate=ClusterMonitoringConfig
// ClusterMonitoring is the Schema for the Cluster Monitoring Operators API
type ClusterMonitoring struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user configuration for the Cluster Monitoring Operator
	// +required
	Spec ClusterMonitoringSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ClusterMonitoringStatus `json:"status,omitempty"`
}

// MonitoringOperatorStatus defines the observed state of MonitoringOperator
type ClusterMonitoringStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:internal
type ClusterMonitoringList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// items is a list of ClusterMonitoring
	// +optional
	Items []ClusterMonitoring `json:"items"`
}

// ClusterMonitoringSpec defines the desired state of Cluster Monitoring Operator
// +required
type ClusterMonitoringSpec struct {
	// userDefined set the deployment mode for user-defined monitoring in addition to the default platform monitoring.
	// +required
	UserDefined UserDefinedMonitoring `json:"userDefined"`

	// metricsServer defines the configuration for the Metrics Server component.
	// The Metrics Server provides container resource metrics for use in autoscaling pipelines.
	// When omitted, this means no opinion and the platform is left to choose a default,
	// which is subject to change over time.
	// +optional
	MetricsServerConfig MetricsServerConfig `json:"metricsServer,omitempty"`
}

// UserDefinedMonitoring config for user-defined projects.
// +required
type UserDefinedMonitoring struct {
	// mode defines the different configurations of UserDefinedMonitoring
	// Valid values are Disabled and NamespaceIsolated
	// Disabled disables monitoring for user-defined projects. This restricts the default monitoring stack, installed in the openshift-monitoring project, to monitor only platform namespaces, which prevents any custom monitoring configurations or resources from being applied to user-defined namespaces.
	// NamespaceIsolated enables monitoring for user-defined projects with namespace-scoped tenancy. This ensures that metrics, alerts, and monitoring data are isolated at the namespace level.
	// +kubebuilder:validation:Enum:="Disabled";"NamespaceIsolated"
	// +required
	Mode UserDefinedMode `json:"mode"`
}

// UserDefinedMode specifies mode for UserDefine Monitoring
// +enum
type UserDefinedMode string

const (
	// UserDefinedDisabled disables monitoring for user-defined projects. This restricts the default monitoring stack, installed in the openshift-monitoring project, to monitor only platform namespaces, which prevents any custom monitoring configurations or resources from being applied to user-defined namespaces.
	UserDefinedDisabled UserDefinedMode = "Disabled"
	// UserDefinedNamespaceIsolated enables monitoring for user-defined projects with namespace-scoped tenancy. This ensures that metrics, alerts, and monitoring data are isolated at the namespace level.
	UserDefinedNamespaceIsolated UserDefinedMode = "NamespaceIsolated"
)

// The MetricsServerConfig resource defines settings for the Metrics Server component.
// This configuration allows users to customize how the Metrics Server is deployed
// and how it operates within the cluster.
type MetricsServerConfig struct {
	// audit defines the audit configuration used by the Metrics Server instance.
	// When omitted, this means no opinion and the platform is left to choose a default,
	// which is subject to change over time. The current default is "metadata".
	// The audit field is optional.
	// +optional
	Audit *Audit `json:"audit,omitempty"`

	// nodeSelector is the node selector applied to metrics server pods.
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default is `kubernetes.io/os: linux` so that Pods can be scheduled onto any available node.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// tolerations is a list of tolerations applied to metrics server components
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// Maximum length for this list is 10
	// +kubebuilder:validation:MaxItems=10
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// resources defines the compute resource requests and limits for the Metrics Server container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// Resources is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// +optional
	Resources *MetricServerConfigContainerResources `json:"resources,omitempty"`

	// topologySpreadConstraints defines rules for how Metrics Server Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Maximum length for this list is 10
	// +kubebuilder:validation:MaxItems=10
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// Audit defines the configuration for Metrics Server audit logging.
type Audit struct {
	// profile specifies the audit log level to use.
	// Valid values are:
	// - "metadata" - log metadata about requests (default)
	// - "request" - log metadata and request payloads
	// - "requestresponse" - log metadata, requests, and responses
	// - "none" - don't log requests
	//
	// See: https://kubernetes.io/docs/tasks/debug-application-cluster/audit/#audit-policy
	// for more details about audit logging.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=metadata;request;requestresponse;none
	// +required
	Profile auditv1.Level `json:"profile"`
}

// ResourceSpec defines the requested and limited value of a resource.
type ResourceSpec struct {
	// request is the minimum amount of the resource required (e.g. "2Mi", "1Gi").
	// This field is optional.
	// +optional
	Request resource.Quantity `json:"request,omitempty"`

	// limit is the maximum amount of the resource allowed (e.g. "2Mi", "1Gi").
	// This field is optional.
	// +optional
	Limit resource.Quantity `json:"limit,omitempty"`
}

// MetricServerConfigContainerResources defines simplified resource requirements for a container.
type MetricServerConfigContainerResources struct {
	// cpu defines the CPU resource limits and requests.
	// This filed is optional
	// +optional
	CPU *ResourceSpec `json:"cpu,omitempty"`

	// memory defines the memory resource limits and requests.
	// This filed is optional
	// +optional
	Memory *ResourceSpec `json:"memory,omitempty"`

	// hugepages is a list of hugepage resource specifications by page size.
	// defines an optional list of unique configurations identified by their `size` field.
	// A maximum of 10 items is allowed.
	// The list is treated as a map, using `size` as the key
	// +optional
	// +listType=map
	// +listMapKey=size
	// +kubebuilder:validation:MaxItems=10
	HugePages []HugePageResource `json:"hugepages,omitempty"`
}

// HugePageResource describes hugepages resources by page size (e.g. 2Mi, 1Gi).
type HugePageResource struct {
	// size of the hugepage (e.g. "2Mi", "1Gi").
	// This field is required.
	// +required
	Size resource.Quantity `json:"size"`

	// request amount for this hugepage size.
	// This filed is optional
	// +optional
	Request resource.Quantity `json:"request,omitempty"`

	// limit amount for this hugepage size.
	// This filed is optional
	// +optional
	Limit resource.Quantity `json:"limit,omitempty"`
}
