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

// ClusterMonitoring is the Custom Resource object which holds the current status of Cluster Monitoring Operator. CMO is a central component of the monitoring stack.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:internal
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/1929
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clustermonitorings,scope=Cluster
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

// ClusterMonitoringStatus defines the observed state of ClusterMonitoring
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
// +kubebuilder:validation:MinProperties=1
type ClusterMonitoringSpec struct {
	// userDefined set the deployment mode for user-defined monitoring in addition to the default platform monitoring.
	// userDefined is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The current default value is `Disabled`.
	// +optional
	UserDefined UserDefinedMonitoring `json:"userDefined,omitempty,omitzero"`
	// alertmanagerConfig allows users to configure how the default Alertmanager instance
	// should be deployed in the `openshift-monitoring` namespace.
	// alertmanagerConfig is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	// The current default value is `DefaultConfig`.
	// +optional
	AlertmanagerConfig AlertmanagerConfig `json:"alertmanagerConfig,omitempty,omitzero"`
	// metricsServerConfig is an optional field that can be used to configure the Kubernetes Metrics Server that runs in the openshift-monitoring namespace.
	// Specifically, it can configure how the Metrics Server instance is deployed, pod scheduling, its audit policy and log verbosity.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	MetricsServerConfig MetricsServerConfig `json:"metricsServerConfig,omitempty,omitzero"`
	// prometheusOperatorConfig is an optional field that can be used to configure the Prometheus Operator component.
	// Specifically, it can configure how the Prometheus Operator instance is deployed, pod scheduling, and resource allocation.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	PrometheusOperatorConfig PrometheusOperatorConfig `json:"prometheusOperatorConfig,omitempty,omitzero"`
	// prometheusOperatorAdmissionWebhookConfig is an optional field that can be used to configure the
	// admission webhook component of Prometheus Operator that runs in the openshift-monitoring namespace.
	// The admission webhook validates PrometheusRule and AlertmanagerConfig objects to ensure they are
	// semantically valid, mutates PrometheusRule annotations, and converts AlertmanagerConfig objects
	// between API versions.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	PrometheusOperatorAdmissionWebhookConfig PrometheusOperatorAdmissionWebhookConfig `json:"prometheusOperatorAdmissionWebhookConfig,omitempty,omitzero"`
	// nodeExporterConfig is an optional field that can be used to configure the node-exporter agent
	// that runs as a DaemonSet in the openshift-monitoring namespace. The node-exporter agent collects
	// hardware and OS-level metrics from every node in the cluster.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	NodeExporterConfig NodeExporterConfig `json:"nodeExporterConfig,omitempty,omitzero"`
}

// UserDefinedMonitoring config for user-defined projects.
type UserDefinedMonitoring struct {
	// mode defines the different configurations of UserDefinedMonitoring
	// Valid values are Disabled and NamespaceIsolated
	// Disabled disables monitoring for user-defined projects. This restricts the default monitoring stack, installed in the openshift-monitoring project, to monitor only platform namespaces, which prevents any custom monitoring configurations or resources from being applied to user-defined namespaces.
	// NamespaceIsolated enables monitoring for user-defined projects with namespace-scoped tenancy. This ensures that metrics, alerts, and monitoring data are isolated at the namespace level.
	// The current default value is `Disabled`.
	// +required
	// +kubebuilder:validation:Enum=Disabled;NamespaceIsolated
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

// alertmanagerConfig provides configuration options for the default Alertmanager instance
// that runs in the `openshift-monitoring` namespace. Use this configuration to control
// whether the default Alertmanager is deployed, how it logs, and how its pods are scheduled.
// +kubebuilder:validation:XValidation:rule="self.deploymentMode == 'CustomConfig' ? has(self.customConfig) : !has(self.customConfig)",message="customConfig is required when deploymentMode is CustomConfig, and forbidden otherwise"
type AlertmanagerConfig struct {
	// deploymentMode determines whether the default Alertmanager instance should be deployed
	// as part of the monitoring stack.
	// Allowed values are Disabled, DefaultConfig, and CustomConfig.
	// When set to Disabled, the Alertmanager instance will not be deployed.
	// When set to DefaultConfig, the platform will deploy Alertmanager with default settings.
	// When set to CustomConfig, the Alertmanager will be deployed with custom configuration.
	//
	// +unionDiscriminator
	// +required
	DeploymentMode AlertManagerDeployMode `json:"deploymentMode,omitempty"`

	// customConfig must be set when deploymentMode is CustomConfig, and must be unset otherwise.
	// When set to CustomConfig, the Alertmanager will be deployed with custom configuration.
	// +optional
	CustomConfig AlertmanagerCustomConfig `json:"customConfig,omitempty,omitzero"`
}

// AlertmanagerCustomConfig represents the configuration for a custom Alertmanager deployment.
// alertmanagerCustomConfig provides configuration options for the default Alertmanager instance
// that runs in the `openshift-monitoring` namespace. Use this configuration to control
// whether the default Alertmanager is deployed, how it logs, and how its pods are scheduled.
// +kubebuilder:validation:MinProperties=1
type AlertmanagerCustomConfig struct {
	// logLevel defines the verbosity of logs emitted by Alertmanager.
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
	LogLevel LogLevel `json:"logLevel,omitempty"`
	// nodeSelector defines the nodes on which the Pods are scheduled
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// resources defines the compute resource requests and limits for the Alertmanager container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 4m
	//      limit: null
	//    - name: memory
	//      request: 40Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Each resource name must be unique within this list.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// secrets defines a list of secrets that need to be mounted into the Alertmanager.
	// The secrets must reside within the same namespace as the Alertmanager object.
	// They will be added as volumes named secret-<secret-name> and mounted at
	// /etc/alertmanager/secrets/<secret-name> within the 'alertmanager' container of
	// the Alertmanager Pods.
	//
	// These secrets can be used to authenticate Alertmanager with endpoint receivers.
	// For example, you can use secrets to:
	// - Provide certificates for TLS authentication with receivers that require private CA certificates
	// - Store credentials for Basic HTTP authentication with receivers that require password-based auth
	// - Store any other authentication credentials needed by your alert receivers
	//
	// This field is optional.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Entries in this list must be unique.
	// +optional
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	Secrets []SecretName `json:"secrets,omitempty"`
	// tolerations defines tolerations for the pods.
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// Defaults are empty/unset.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// topologySpreadConstraints defines rules for how Alertmanager Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Default is empty list.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// volumeClaimTemplate Defines persistent storage for Alertmanager. Use this setting to
	// configure the persistent volume claim, including storage class, volume
	// size, and name.
	// If omitted, the Pod uses ephemeral storage and alert data will not persist
	// across restarts.
	// This field is optional.
	// +optional
	VolumeClaimTemplate *v1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
}

// AlertManagerDeployMode defines the deployment state of the platform Alertmanager instance.
//
// Possible values:
// - "Disabled": The Alertmanager instance will not be deployed.
// - "DefaultConfig": The Alertmanager instance will be deployed with default settings.
// - "CustomConfig": The Alertmanager instance will be deployed with custom configuration.
// +kubebuilder:validation:Enum=Disabled;DefaultConfig;CustomConfig
type AlertManagerDeployMode string

const (
	// AlertManagerModeDisabled means the Alertmanager instance will not be deployed.
	AlertManagerDeployModeDisabled AlertManagerDeployMode = "Disabled"
	// AlertManagerModeDefaultConfig means the Alertmanager instance will be deployed with default settings.
	AlertManagerDeployModeDefaultConfig AlertManagerDeployMode = "DefaultConfig"
	// AlertManagerModeCustomConfig means the Alertmanager instance will be deployed with custom configuration.
	AlertManagerDeployModeCustomConfig AlertManagerDeployMode = "CustomConfig"
)

// logLevel defines the verbosity of logs emitted by Alertmanager.
// Valid values are Error, Warn, Info and Debug.
// +kubebuilder:validation:Enum=Error;Warn;Info;Debug
type LogLevel string

const (
	// Error only errors will be logged.
	LogLevelError LogLevel = "Error"
	// Warn, both warnings and errors will be logged.
	LogLevelWarn LogLevel = "Warn"
	// Info, general information, warnings, and errors will all be logged.
	LogLevelInfo LogLevel = "Info"
	// Debug, detailed debugging information will be logged.
	LogLevelDebug LogLevel = "Debug"
)

// ContainerResource defines a single resource requirement for a container.
// +kubebuilder:validation:XValidation:rule="has(self.request) || has(self.limit)",message="at least one of request or limit must be set"
// +kubebuilder:validation:XValidation:rule="!(has(self.request) && has(self.limit)) || quantity(self.limit).compareTo(quantity(self.request)) >= 0",message="limit must be greater than or equal to request"
type ContainerResource struct {
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

// SecretName is a type that represents the name of a Secret in the same namespace.
// It must be at most 253 characters in length.
// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
// +kubebuilder:validation:MaxLength=63
type SecretName string

// MetricsServerConfig provides configuration options for the Metrics Server instance
// that runs in the `openshift-monitoring` namespace. Use this configuration to control
// how the Metrics Server instance is deployed, how it logs, and how its pods are scheduled.
// +kubebuilder:validation:MinProperties=1
type MetricsServerConfig struct {
	// audit defines the audit configuration used by the Metrics Server instance.
	// audit is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	//The current default sets audit.profile to Metadata
	// +optional
	Audit Audit `json:"audit,omitempty,omitzero"`
	// nodeSelector defines the nodes on which the Pods are scheduled
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// tolerations defines tolerations for the pods.
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// Defaults are empty/unset.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// verbosity defines the verbosity of log messages for Metrics Server.
	// Valid values are Errors, Info, Trace, TraceAll and omitted.
	// When set to Errors, only critical messages and errors are logged.
	// When set to Info, only basic information messages are logged.
	// When set to Trace, information useful for general debugging is logged.
	// When set to TraceAll, detailed information about metric scraping is logged.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, that is subject to change over time.
	// The current default value is `Errors`
	// +optional
	Verbosity VerbosityLevel `json:"verbosity,omitempty,omitzero"`
	// resources defines the compute resource requests and limits for the Metrics Server container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 4m
	//      limit: null
	//    - name: memory
	//      request: 40Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Each resource name must be unique within this list.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// topologySpreadConstraints defines rules for how Metrics Server Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Default is empty list.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// PrometheusOperatorConfig provides configuration options for the Prometheus Operator instance
// Use this configuration to control how the Prometheus Operator instance is deployed, how it logs, and how its pods are scheduled.
// +kubebuilder:validation:MinProperties=1
type PrometheusOperatorConfig struct {
	// logLevel defines the verbosity of logs emitted by Prometheus Operator.
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
	LogLevel LogLevel `json:"logLevel,omitempty"`
	// nodeSelector defines the nodes on which the Pods are scheduled
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// When specified, nodeSelector must contain at least 1 entry and must not contain more than 10 entries.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// resources defines the compute resource requests and limits for the Prometheus Operator container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 4m
	//      limit: null
	//    - name: memory
	//      request: 40Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Each resource name must be unique within this list.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// tolerations defines tolerations for the pods.
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// Defaults are empty/unset.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// topologySpreadConstraints defines rules for how Prometheus Operator Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Default is empty list.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// PrometheusOperatorAdmissionWebhookConfig provides configuration options for the admission webhook
// component of Prometheus Operator that runs in the `openshift-monitoring` namespace. The admission
// webhook validates PrometheusRule and AlertmanagerConfig objects, mutates PrometheusRule annotations,
// and converts AlertmanagerConfig objects between API versions.
// +kubebuilder:validation:MinProperties=1
type PrometheusOperatorAdmissionWebhookConfig struct {
	// resources defines the compute resource requests and limits for the
	// prometheus-operator-admission-webhook container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 5m
	//      limit: null
	//    - name: memory
	//      request: 30Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Each resource name must be unique within this list.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// topologySpreadConstraints defines rules for how admission webhook Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Default is empty list.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// NodeExporterConfig provides configuration options for the node-exporter agent
// that runs as a DaemonSet in the `openshift-monitoring` namespace. The node-exporter agent collects
// hardware and OS-level metrics from every node in the cluster, including CPU, memory, disk, and
// network statistics.
// +kubebuilder:validation:MinProperties=1
type NodeExporterConfig struct {
	// nodeSelector defines the nodes on which the Pods are scheduled.
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// When specified, nodeSelector must contain at least 1 entry and must not contain more than 10 entries.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// resources defines the compute resource requests and limits for the node-exporter container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// This field is optional.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// This is a simplified API that maps to Kubernetes ResourceRequirements.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 8m
	//      limit: null
	//    - name: memory
	//      request: 32Mi
	//      limit: null
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// Each resource name must be unique within this list.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// tolerations defines tolerations for the pods.
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default is to tolerate all taints (operator: Exists without any key),
	// which is typical for DaemonSets that must run on every node.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// collectors configures which node-exporter metric collectors are enabled.
	// collectors is optional.
	// Each collector can be individually enabled or disabled. Some collectors may have
	// additional configuration options.
	//
	// When omitted, this means no opinion and the platform is left to choose a reasonable
	// default, which is subject to change over time.
	// +optional
	Collectors NodeExporterCollectorConfig `json:"collectors,omitempty,omitzero"`
	// maxProcs sets the target number of CPUs on which the node-exporter process will run.
	// maxProcs is optional.
	// Use this setting to override the default value, which is set either to 4 or to the number
	// of CPUs on the host, whichever is smaller.
	// The default value is computed at runtime and set via the GOMAXPROCS environment variable before
	// node-exporter is launched.
	// If a kernel deadlock occurs or if performance degrades when reading from sysfs concurrently,
	// you can change this value to 1, which limits node-exporter to running on one CPU.
	// For nodes with a high CPU count, setting the limit to a low number saves resources by preventing
	// Go routines from being scheduled to run on all CPUs. However, I/O performance degrades if the
	// maxProcs value is set too low and there are many metrics to collect.
	// The minimum value is 1.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is min(4, number of host CPUs).
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxProcs int32 `json:"maxProcs,omitempty"`
	// ignoredNetworkDevices is a list of regular expression patterns that match network devices
	// to be excluded from the relevant collector configuration such as netdev, netclass, and ethtool.
	// ignoredNetworkDevices is optional.
	//
	// When omitted, the Cluster Monitoring Operator uses a predefined list of devices to be excluded
	// to minimize the impact on memory usage.
	// When set as an empty list, no devices are excluded.
	// If you modify this setting, monitor the prometheus-k8s deployment closely for excessive memory usage.
	// Maximum length for this list is 50.
	// Minimum length for this list is 1.
	// Each entry must be at most 1024 characters long.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	// +optional
	IgnoredNetworkDevices []NodeExporterIgnoredNetworkDevice `json:"ignoredNetworkDevices,omitempty"`
}

// NodeExporterIgnoredNetworkDevice is a regular expression pattern that matches a network device name
// to be excluded from node-exporter metric collection for collectors such as netdev, netclass, and ethtool.
// Must be a valid regular expression and at most 1024 characters.
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=1024
type NodeExporterIgnoredNetworkDevice string

// NodeExporterCollectorState defines whether a node-exporter collector is enabled or disabled.
// Valid values are "Enabled" and "Disabled".
// +kubebuilder:validation:Enum=Enabled;Disabled
type NodeExporterCollectorState string

const (
	// NodeExporterCollectorEnabled means the collector is active and will produce metrics.
	NodeExporterCollectorEnabled NodeExporterCollectorState = "Enabled"
	// NodeExporterCollectorDisabled means the collector is inactive and will not produce metrics.
	NodeExporterCollectorDisabled NodeExporterCollectorState = "Disabled"
)

// NodeExporterNetlinkState defines whether the netlink implementation of the netclass
// collector is used. Valid values are "Enabled" and "Disabled".
// +kubebuilder:validation:Enum=Enabled;Disabled
type NodeExporterNetlinkState string

const (
	// NodeExporterNetlinkEnabled activates the netlink implementation.
	NodeExporterNetlinkEnabled NodeExporterNetlinkState = "Enabled"
	// NodeExporterNetlinkDisabled deactivates the netlink implementation, falling back to the default sysfs implementation.
	NodeExporterNetlinkDisabled NodeExporterNetlinkState = "Disabled"
)

// NodeExporterCollectorConfig defines settings for individual collectors
// of the node-exporter agent. Each collector can be individually enabled or disabled.
// +kubebuilder:validation:MinProperties=1
type NodeExporterCollectorConfig struct {
	// cpuFreq configures the cpufreq collector, which collects CPU frequency statistics.
	// cpuFreq is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// Under certain circumstances, enabling the cpufreq collector increases CPU usage on machines
	// with many cores. If you enable this collector and have machines with many cores, monitor your
	// systems closely for excessive CPU usage.
	// +optional
	CpuFreq NodeExporterCollectorCpufreqConfig `json:"cpuFreq,omitempty,omitzero"`
	// tcpStat configures the tcpstat collector, which collects TCP connection statistics.
	// tcpStat is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// +optional
	TcpStat NodeExporterCollectorTcpStatConfig `json:"tcpStat,omitempty,omitzero"`
	// ethtool configures the ethtool collector, which collects ethernet device statistics.
	// ethtool is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// +optional
	Ethtool NodeExporterCollectorEthtoolConfig `json:"ethtool,omitempty,omitzero"`
	// netDev configures the netdev collector, which collects network device statistics.
	// netDev is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is enabled.
	// +optional
	NetDev NodeExporterCollectorNetDevConfig `json:"netDev,omitempty,omitzero"`
	// netClass configures the netclass collector, which collects information about network devices.
	// netClass is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is enabled with netlink mode active.
	// +optional
	NetClass NodeExporterCollectorNetClassConfig `json:"netClass,omitempty,omitzero"`
	// buddyInfo configures the buddyinfo collector, which collects statistics about memory
	// fragmentation from the node_buddyinfo_blocks metric. This metric collects data from /proc/buddyinfo.
	// buddyInfo is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// +optional
	BuddyInfo NodeExporterCollectorBuddyInfoConfig `json:"buddyInfo,omitempty,omitzero"`
	// mountStats configures the mountstats collector, which collects statistics about NFS volume
	// I/O activities.
	// mountStats is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// Enabling this collector may produce metrics with high cardinality. If you enable this
	// collector, closely monitor the prometheus-k8s deployment for excessive memory usage.
	// +optional
	MountStats NodeExporterCollectorMountStatsConfig `json:"mountStats,omitempty,omitzero"`
	// ksmd configures the ksmd collector, which collects statistics from the kernel same-page
	// merger daemon.
	// ksmd is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// +optional
	Ksmd NodeExporterCollectorKSMDConfig `json:"ksmd,omitempty,omitzero"`
	// processes configures the processes collector, which collects statistics from processes and
	// threads running in the system.
	// processes is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// +optional
	Processes NodeExporterCollectorProcessesConfig `json:"processes,omitempty,omitzero"`
	// systemd configures the systemd collector, which collects statistics on the systemd daemon
	// and its managed services.
	// systemd is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is disabled.
	// Enabling this collector with a long list of selected units may produce metrics with high
	// cardinality. If you enable this collector, closely monitor the prometheus-k8s deployment
	// for excessive memory usage.
	// +optional
	Systemd NodeExporterCollectorSystemdConfig `json:"systemd,omitempty,omitzero"`
}

// NodeExporterCollectorCpufreqConfig provides configuration for the cpufreq collector
// of the node-exporter agent. The cpufreq collector collects CPU frequency statistics.
// It is disabled by default.
type NodeExporterCollectorCpufreqConfig struct {
	// enabled enables or disables the cpufreq collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the cpufreq collector is active and CPU frequency statistics are collected.
	// When set to "Disabled", the cpufreq collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorTcpStatConfig provides configuration for the tcpstat collector
// of the node-exporter agent. The tcpstat collector collects TCP connection statistics.
// It is disabled by default.
type NodeExporterCollectorTcpStatConfig struct {
	// enabled enables or disables the tcpstat collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the tcpstat collector is active and TCP connection statistics are collected.
	// When set to "Disabled", the tcpstat collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorEthtoolConfig provides configuration for the ethtool collector
// of the node-exporter agent. The ethtool collector collects ethernet device statistics.
// It is disabled by default.
type NodeExporterCollectorEthtoolConfig struct {
	// enabled enables or disables the ethtool collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the ethtool collector is active and ethernet device statistics are collected.
	// When set to "Disabled", the ethtool collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorNetDevConfig provides configuration for the netdev collector
// of the node-exporter agent. The netdev collector collects network device statistics
// such as bytes, packets, errors, and drops per device.
// It is enabled by default.
type NodeExporterCollectorNetDevConfig struct {
	// enabled enables or disables the netdev collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the netdev collector is active and network device statistics are collected.
	// When set to "Disabled", the netdev collector is inactive and the corresponding metrics become unavailable.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorNetClassConfig provides configuration for the netclass collector
// of the node-exporter agent. The netclass collector collects information about network devices
// such as network speed, MTU, and carrier status.
// It is enabled by default.
type NodeExporterCollectorNetClassConfig struct {
	// enabled enables or disables the netclass collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the netclass collector is active and network class information is collected.
	// When set to "Disabled", the netclass collector is inactive and the corresponding metrics become unavailable.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
	// useNetlink activates the netlink implementation of the netclass collector.
	// useNetlink is optional.
	// This implementation improves the performance of the netclass collector.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the netlink implementation is used for improved performance.
	// When set to "Disabled", the default sysfs implementation is used.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default,
	// which is subject to change over time. The current default is "Enabled".
	// +optional
	UseNetlink NodeExporterNetlinkState `json:"useNetlink,omitempty"`
}

// NodeExporterCollectorBuddyInfoConfig provides configuration for the buddyinfo collector
// of the node-exporter agent. The buddyinfo collector collects statistics about memory fragmentation
// from the node_buddyinfo_blocks metric using data from /proc/buddyinfo.
// It is disabled by default.
type NodeExporterCollectorBuddyInfoConfig struct {
	// enabled enables or disables the buddyinfo collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the buddyinfo collector is active and memory fragmentation statistics are collected.
	// When set to "Disabled", the buddyinfo collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorMountStatsConfig provides configuration for the mountstats collector
// of the node-exporter agent. The mountstats collector collects statistics about NFS volume I/O activities.
// It is disabled by default.
// Enabling this collector may produce metrics with high cardinality. If you enable this
// collector, closely monitor the prometheus-k8s deployment for excessive memory usage.
type NodeExporterCollectorMountStatsConfig struct {
	// enabled enables or disables the mountstats collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the mountstats collector is active and NFS volume I/O statistics are collected.
	// When set to "Disabled", the mountstats collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorKSMDConfig provides configuration for the ksmd collector
// of the node-exporter agent. The ksmd collector collects statistics from the kernel
// same-page merger daemon.
// It is disabled by default.
type NodeExporterCollectorKSMDConfig struct {
	// enabled enables or disables the ksmd collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the ksmd collector is active and kernel same-page merger statistics are collected.
	// When set to "Disabled", the ksmd collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorProcessesConfig provides configuration for the processes collector
// of the node-exporter agent. The processes collector collects statistics from processes and threads
// running in the system.
// It is disabled by default.
type NodeExporterCollectorProcessesConfig struct {
	// enabled enables or disables the processes collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the processes collector is active and process/thread statistics are collected.
	// When set to "Disabled", the processes collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
}

// NodeExporterCollectorSystemdConfig provides configuration for the systemd collector
// of the node-exporter agent. The systemd collector collects statistics on the systemd daemon
// and its managed services.
// It is disabled by default.
// Enabling this collector with a long list of selected units may produce metrics with high
// cardinality. If you enable this collector, closely monitor the prometheus-k8s deployment
// for excessive memory usage.
type NodeExporterCollectorSystemdConfig struct {
	// enabled enables or disables the systemd collector.
	// This field is required.
	// Valid values are "Enabled" and "Disabled".
	// When set to "Enabled", the systemd collector is active and systemd unit statistics are collected.
	// When set to "Disabled", the systemd collector is inactive.
	// +required
	Enabled NodeExporterCollectorState `json:"enabled,omitempty"`
	// units is a list of regular expression patterns that match systemd units to be included
	// by the systemd collector.
	// units is optional.
	// By default, the list is empty, so the collector exposes no metrics for systemd units.
	// Each entry must be a valid regular expression and at most 1024 characters.
	// Maximum length for this list is 50.
	// Minimum length for this list is 1.
	// Entries in this list must be unique.
	// +kubebuilder:validation:MaxItems=50
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	// +optional
	Units []NodeExporterSystemdUnit `json:"units,omitempty"`
}

// NodeExporterSystemdUnit is a regular expression pattern that matches a systemd unit name.
// Must be at most 1024 characters.
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=1024
type NodeExporterSystemdUnit string

// AuditProfile defines the audit log level for the Metrics Server.
// +kubebuilder:validation:Enum=None;Metadata;Request;RequestResponse
type AuditProfile string

const (
	// AuditProfileNone disables audit logging
	AuditProfileNone AuditProfile = "None"
	// AuditProfileMetadata logs request metadata (requesting user, timestamp, resource, verb, etc.) but not request or response body
	AuditProfileMetadata AuditProfile = "Metadata"
	// AuditProfileRequest logs event metadata and request body but not response body
	AuditProfileRequest AuditProfile = "Request"
	// AuditProfileRequestResponse logs event metadata, request and response bodies
	AuditProfileRequestResponse AuditProfile = "RequestResponse"
)

// VerbosityLevel defines the verbosity of log messages for Metrics Server.
// +kubebuilder:validation:Enum=Errors;Info;Trace;TraceAll
type VerbosityLevel string

const (
	// VerbosityLevelErrors means only critical messages and errors are logged.
	VerbosityLevelErrors VerbosityLevel = "Errors"
	// VerbosityLevelInfo means basic informational messages are logged.
	VerbosityLevelInfo VerbosityLevel = "Info"
	// VerbosityLevelTrace means extended information useful for general debugging is logged.
	VerbosityLevelTrace VerbosityLevel = "Trace"
	// VerbosityLevelTraceAll means detailed information about metric scraping operations is logged.
	VerbosityLevelTraceAll VerbosityLevel = "TraceAll"
)

// Audit profile configurations
type Audit struct {
	// profile is a required field for configuring the audit log level of the Kubernetes Metrics Server.
	// Allowed values are None, Metadata, Request, or RequestResponse.
	// When set to None, audit logging is disabled and no audit events are recorded.
	// When set to Metadata, only request metadata (such as requesting user, timestamp, resource, verb, etc.) is logged, but not the request or response body.
	// When set to Request, event metadata and the request body are logged, but not the response body.
	// When set to RequestResponse, event metadata, request body, and response body are all logged, providing the most detailed audit information.
	//
	// See: https://kubernetes.io/docs/tasks/debug-application-cluster/audit/#audit-policy
	// for more information about auditing and log levels.
	// +required
	Profile AuditProfile `json:"profile,omitempty"`
}
