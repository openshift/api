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
	// prometheusConfig provides configuration options for the default platform Prometheus instance
	// that runs in the `openshift-monitoring` namespace. This configuration applies only to the
	// platform Prometheus instance; user-workload Prometheus instances are configured separately.
	//
	// This field allows you to customize how the platform Prometheus is deployed and operated, including:
	//   - Pod scheduling (node selectors, tolerations, topology spread constraints)
	//   - Resource allocation (CPU, memory requests/limits)
	//   - Retention policies (how long metrics are stored)
	//   - External integrations (remote write, additional alertmanagers)
	//
	// This field is optional. When omitted, the platform chooses reasonable defaults, which may change over time.
	// +optional
	PrometheusConfig PrometheusConfig `json:"prometheusConfig,omitempty,omitzero"`
	// metricsServerConfig is an optional field that can be used to configure the Kubernetes Metrics Server that runs in the openshift-monitoring namespace.
	// Specifically, it can configure how the Metrics Server instance is deployed, pod scheduling, its audit policy and log verbosity.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	MetricsServerConfig MetricsServerConfig `json:"metricsServerConfig,omitempty,omitzero"`
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
	// Maximum length for this list is 10
	// Minimum length for this list is 1
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
	// Minimum length for this list is 1
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// volumeClaimTemplate defines persistent storage for Alertmanager. Use this setting to
	// configure the persistent volume claim, including storage class and volume size.
	// If omitted, the Pod uses ephemeral storage and alert data will not persist
	// across restarts.
	// +optional
	VolumeClaimTemplate VolumeClaimConfig `json:"volumeClaimTemplate,omitempty,omitzero"`
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

// LogLevel defines the verbosity of logs emitted by Alertmanager.
// Valid values are Error, Warn, Info and Debug.
// +kubebuilder:validation:Enum=Error;Warn;Info;Debug
type LogLevel string

const (
	// LogLevelError only errors will be logged.
	LogLevelError LogLevel = "Error"
	// LogLevelWarn, both warnings and errors will be logged.
	LogLevelWarn LogLevel = "Warn"
	// LogLevelInfo, general information, warnings, and errors will all be logged.
	LogLevelInfo LogLevel = "Info"
	// LogLevelDebug, detailed debugging information will be logged.
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
	// Maximum length for this list is 10
	// Minimum length for this list is 1
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
	// Minimum length for this list is 1
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// PrometheusConfig provides configuration options for the Prometheus instance.
// Use this configuration to control
// Prometheus deployment, pod scheduling, resource allocation, retention policies, and external integrations.
// +kubebuilder:validation:MinProperties=1
type PrometheusConfig struct {
	// additionalAlertmanagerConfigs configures additional Alertmanager instances that receive alerts from
	// the Prometheus component. This is useful for organizations that need to:
	//   - Send alerts to external monitoring systems (like PagerDuty, Slack, or custom webhooks)
	//   - Route different types of alerts to different teams or systems
	//   - Integrate with existing enterprise alerting infrastructure
	//   - Maintain separate alert routing for compliance or organizational requirements
	// By default, no additional Alertmanager instances are configured.
	// When omitted, no additional Alertmanager instances are configured (default behavior).
	// When provided, at least one configuration must be specified (minimum 1, maximum 10 items).
	// Each entry must have a unique name field, which serves as the map key for server-side apply.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=name
	AdditionalAlertmanagerConfigs []AdditionalAlertmanagerConfig `json:"additionalAlertmanagerConfigs,omitempty"`
	// enforcedBodySizeLimitBytes enforces a body size limit (in bytes) for Prometheus scraped metrics.
	// If a scraped target's body response is larger than the limit, the scrape will fail.
	// This helps protect Prometheus from targets that return excessively large responses.
	// The value is specified in bytes (e.g., 4194304 for 4MB, 1073741824 for 1GB).
	// When omitted, the Cluster Monitoring Operator automatically calculates an appropriate
	// limit based on cluster capacity. Set an explicit value to override the automatic calculation.
	// Minimum value is 10240 (10kB).
	// Maximum value is 1073741824 (1GB).
	// +kubebuilder:validation:Minimum=10240
	// +kubebuilder:validation:Maximum=1073741824
	// +optional
	EnforcedBodySizeLimitBytes int64 `json:"enforcedBodySizeLimitBytes,omitempty"`
	// externalLabels defines labels to be attached to time series and alerts
	// when communicating with external systems such as federation, remote storage,
	// and Alertmanager. These labels are not stored with metrics on disk; they are
	// only added when data leaves Prometheus (e.g., during federation queries,
	// remote write, or alert notifications).
	// At least 1 label must be specified when set, with a maximum of 50 labels allowed.
	// Each label key must be unique within this list.
	// When omitted, no external labels are applied.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=50
	// +listType=map
	// +listMapKey=key
	ExternalLabels []Label `json:"externalLabels,omitempty"`
	// logLevel defines the verbosity of logs emitted by Prometheus.
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
	// nodeSelector defines the nodes on which the Pods are scheduled.
	// nodeSelector is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default value is `kubernetes.io/os: linux`.
	// Maximum of 10 node selector key-value pairs can be specified.
	// +optional
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=10
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// queryLogFile specifies the file to which PromQL queries are logged.
	// This setting can be either a filename, in which
	// case the queries are saved to an `emptyDir` volume
	// at `/var/log/prometheus`, or a full path to a location where
	// an `emptyDir` volume will be mounted and the queries saved.
	// Writing to `/dev/stderr`, `/dev/stdout` or `/dev/null` is supported, but
	// writing to any other `/dev/` path is not supported. Relative paths are
	// also not supported.
	// By default, PromQL queries are not logged.
	// Must be an absolute path starting with `/` or a simple filename without path separators.
	// Must not contain consecutive slashes, end with a slash, or include '..' path traversal.
	// Must contain only alphanumeric characters, '.', '_', '-', or '/'.
	// Must be between 1 and 255 characters in length.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-zA-Z0-9._/-]+$')",message="must contain only alphanumeric characters, '.', '_', '-', or '/'"
	// +kubebuilder:validation:XValidation:rule="self.startsWith('/') || !self.contains('/')",message="must be an absolute path starting with '/' or a simple filename without '/'"
	// +kubebuilder:validation:XValidation:rule="!self.startsWith('/dev/') || self in ['/dev/stdout', '/dev/stderr', '/dev/null']",message="only /dev/stdout, /dev/stderr, and /dev/null are allowed as /dev/ paths"
	// +kubebuilder:validation:XValidation:rule="!self.contains('//') && !self.endsWith('/') && !self.contains('..')",message="must not contain '//', end with '/', or contain '..'"
	QueryLogFile string `json:"queryLogFile,omitempty"`
	// remoteWrite defines the remote write configuration, including URL, authentication, and relabeling settings.
	// Remote write allows Prometheus to send metrics it collects to external long-term storage systems.
	// When omitted, no remote write endpoints are configured.
	// When provided, at least one configuration must be specified (minimum 1, maximum 10 items).
	// Each entry must have a unique URL.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=url
	// +optional
	RemoteWrite []RemoteWriteSpec `json:"remoteWrite,omitempty"`
	// resources defines the compute resource requests and limits for the Prometheus container.
	// This includes CPU, memory and HugePages constraints to help control scheduling and resource usage.
	// When not specified, defaults are used by the platform. Requests cannot exceed limits.
	// Each entry must have a unique resource name.
	// Minimum of 1 and maximum of 10 resource entries can be specified.
	// The current default values are:
	//   resources:
	//    - name: cpu
	//      request: 4m
	//    - name: memory
	//      request: 40Mi
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	Resources []ContainerResource `json:"resources,omitempty"`
	// retention configures how long Prometheus retains metrics data and how much storage it can use.
	// When omitted, the platform chooses reasonable defaults (currently 15 days retention, no size limit).
	// +optional
	Retention Retention `json:"retention,omitempty,omitzero"`
	// tolerations defines tolerations for the pods.
	// tolerations is optional.
	//
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// Defaults are empty/unset.
	// Maximum length for this list is 10
	// Minimum length for this list is 1
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// topologySpreadConstraints defines rules for how Prometheus Pods should be distributed
	// across topology domains such as zones, nodes, or other user-defined labels.
	// topologySpreadConstraints is optional.
	// This helps improve high availability and resource efficiency by avoiding placing
	// too many replicas in the same failure domain.
	//
	// When omitted, this means no opinion and the platform is left to choose a default, which is subject to change over time.
	// This field maps directly to the `topologySpreadConstraints` field in the Pod spec.
	// Default is empty list.
	// Maximum length for this list is 10.
	// Minimum length for this list is 1
	// Entries must have unique topologyKey and whenUnsatisfiable pairs.
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// collectionProfile defines the metrics collection profile that Prometheus uses to collect
	// metrics from the platform components. Supported values are `Full` or
	// `Minimal`. In the `Full` profile (default), Prometheus collects all
	// metrics that are exposed by the platform components. In the `Minimal`
	// profile, Prometheus only collects metrics necessary for the default
	// platform alerts, recording rules, telemetry and console dashboards.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is `Full`.
	// +optional
	CollectionProfile CollectionProfile `json:"collectionProfile,omitempty"`
	// volumeClaimTemplate defines persistent storage for Prometheus. Use this setting to
	// configure the persistent volume claim, including storage class and volume size.
	// If omitted, the Pod uses ephemeral storage and Prometheus data will not persist
	// across restarts.
	// +optional
	VolumeClaimTemplate *v1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty,omitzero"`
}

// AlertmanagerScheme defines the URL scheme to use when communicating with Alertmanager instances.
// +kubebuilder:validation:Enum=HTTP;HTTPS
type AlertmanagerScheme string

const (
	AlertmanagerSchemeHTTP  AlertmanagerScheme = "HTTP"
	AlertmanagerSchemeHTTPS AlertmanagerScheme = "HTTPS"
)

// AdditionalAlertmanagerConfig represents configuration for additional Alertmanager instances.
// The `AdditionalAlertmanagerConfig` resource defines settings for how a
// component communicates with additional Alertmanager instances.
type AdditionalAlertmanagerConfig struct {
	// name is a unique identifier for this Alertmanager configuration entry.
	// This field is used as the map key for server-side apply, allowing multiple
	// controllers (e.g., ACM, cluster admins) to independently manage their own entries.
	// The name must be a valid DNS subdomain (RFC 1123): lowercase alphanumeric characters,
	// hyphens, or periods, and must start and end with an alphanumeric character.
	// Maximum length is 253 characters.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +required
	Name string `json:"name,omitempty"`
	// authorization configures the authentication method for Alertmanager connections.
	// Supports bearer token authentication. When omitted, no authentication is used.
	// +optional
	Authorization AuthorizationConfig `json:"authorization,omitempty,omitzero"`
	// pathPrefix defines an optional URL path prefix to prepend to the Alertmanager API endpoints.
	// For example, if your Alertmanager is behind a reverse proxy at "/alertmanager/",
	// set this to "/alertmanager" so requests go to "/alertmanager/api/v1/alerts" instead of "/api/v1/alerts".
	// This is commonly needed when Alertmanager is deployed behind ingress controllers or load balancers.
	// Must start with "/" and not end with "/" (unless it is the root path "/").
	// Must not contain query strings ("?") or fragments ("#").
	// +kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self.startsWith('/')",message="pathPrefix must start with '/'"
	// +kubebuilder:validation:XValidation:rule="!self.endsWith('/') || self == '/'",message="pathPrefix must not end with '/' (except for root '/')"
	// +kubebuilder:validation:XValidation:rule="!self.contains('?') && !self.contains('#')",message="pathPrefix must not contain '?' or '#'"
	// +optional
	PathPrefix string `json:"pathPrefix,omitempty"`
	// scheme defines the URL scheme to use when communicating with Alertmanager
	// instances.
	// Possible values are `HTTP` or `HTTPS`.
	// When omitted, defaults to `HTTP`.
	// +kubebuilder:default=HTTP
	// +optional
	Scheme AlertmanagerScheme `json:"scheme,omitempty"`
	// staticConfigs is a list of statically configured Alertmanager endpoints in the form
	// of `<host>:<port>`. Each entry must be a valid hostname, IPv4 address, or IPv6 address
	// (in brackets) followed by a colon and a valid port number (1-65535).
	// Examples: "alertmanager.example.com:9093", "192.168.1.100:9093", "[::1]:9093"
	// At least one endpoint must be specified (minimum 1, maximum 10 endpoints).
	// Each entry must be unique.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:items:MaxLength=255
	// +kubebuilder:validation:items:XValidation:rule="isURL('http://' + self) && size(url('http://' + self).getHostname()) > 0 && size(url('http://' + self).getPort()) > 0 && int(url('http://' + self).getPort()) >= 1 && int(url('http://' + self).getPort()) <= 65535",message="must be a valid 'host:port' where host is a DNS name, IPv4, or IPv6 address (in brackets), and port is 1-65535"
	// +listType=set
	// +required
	StaticConfigs []string `json:"staticConfigs,omitempty"`
	// timeoutSeconds defines the timeout in seconds for requests to Alertmanager.
	// When omitted, the default is 10 seconds.
	// Minimum value is 1 second.
	// Maximum value is 600 seconds (10 minutes).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=600
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// tlsConfig defines the TLS settings to use for Alertmanager connections.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	TLSConfig TLSConfig `json:"tlsConfig,omitempty,omitzero"`
}

// Label represents a key/value pair for external labels.
type Label struct {
	// key is the name of the label.
	// Prometheus supports UTF-8 label names, so any valid UTF-8 string is allowed.
	// Must be between 1 and 128 characters in length.
	// +required
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key,omitempty"`
	// value is the value of the label.
	// Must be between 1 and 128 characters in length.
	// +required
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value,omitempty"`
}

// RemoteWriteSpec represents configuration for remote write endpoints.
type RemoteWriteSpec struct {
	// url is the URL of the remote write endpoint.
	// Must be a valid URL with http or https scheme.
	// Must be between 1 and 2048 characters in length.
	// +required
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="isURL(self) && (url(self).getScheme() == 'http' || url(self).getScheme() == 'https')",message="must be a valid URL with http or https scheme"
	URL string `json:"url,omitempty"`
	// name is an optional identifier for this remote write configuration.
	// When omitted, Prometheus generates a unique name automatically.
	// Must contain only alphanumeric characters, hyphens, and underscores.
	// Must be between 1 and 63 characters in length when specified.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-zA-Z0-9_-]+$')",message="must contain only alphanumeric characters, hyphens, and underscores"
	Name string `json:"name,omitempty"`
	// remoteTimeoutSeconds is the timeout in seconds for requests to the remote write endpoint.
	// When omitted, the default is 30 seconds.
	// Minimum value is 1 second.
	// Maximum value is 600 seconds (10 minutes).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=600
	// +optional
	RemoteTimeoutSeconds int32 `json:"remoteTimeoutSeconds,omitempty"`
	// writeRelabelConfigs is a list of relabeling rules to apply before sending data to the remote endpoint.
	// When omitted, no relabeling is performed and all metrics are sent as-is.
	// Minimum of 1 and maximum of 10 relabeling rules can be specified.
	// Each rule must have a unique name.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=map
	// +listMapKey=name
	WriteRelabelConfigs []RelabelConfig `json:"writeRelabelConfigs,omitempty"`
}

// RelabelConfig represents a relabeling rule.
// +kubebuilder:validation:XValidation:rule="self.action in ['Replace', 'HashMod', 'Lowercase', 'Uppercase', 'KeepEqual', 'DropEqual'] ? (has(self.targetLabel) && size(self.targetLabel) > 0) : !has(self.targetLabel)",message="targetLabel is required when action is Replace, HashMod, Lowercase, Uppercase, KeepEqual or DropEqual, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.action in ['Replace', 'LabelMap'] || !has(self.replacement)",message="replacement is only valid when action is Replace or LabelMap"
type RelabelConfig struct {
	// name is a unique identifier for this relabel configuration.
	// Must contain only alphanumeric characters, hyphens, and underscores.
	// Must be between 1 and 63 characters in length.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-zA-Z0-9_-]+$')",message="must contain only alphanumeric characters, hyphens, and underscores"
	Name string `json:"name,omitempty"`
	// sourceLabels specifies which label names to extract from each series for this relabeling rule.
	// Each entry must be a valid label name (non-empty).
	// The values of these labels are joined together using the configured separator,
	// and the resulting string is then matched against the regular expression for
	// the replace, keep, or drop actions.
	// If a referenced label does not exist on a series, Prometheus substitutes an empty string.
	// When omitted, the rule operates without extracting source labels (useful for actions like labelmap).
	// Minimum of 1 and maximum of 10 source labels can be specified, each between 1 and 128 characters.
	// Each entry must be unique.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=128
	// +listType=set
	SourceLabels []string `json:"sourceLabels,omitempty"`
	// separator is the character sequence used to join source label values.
	// Common examples: ";", ",", "::", "|||".
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is ";".
	// Must be between 1 and 5 characters in length when specified.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=5
	Separator string `json:"separator,omitempty"`
	// regex is the regular expression to match against the concatenated source label values.
	// Must be a valid RE2 regular expression (https://github.com/google/re2/wiki/Syntax).
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is "(.*)" to match everything.
	// Must be between 1 and 1000 characters in length when specified.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000
	Regex string `json:"regex,omitempty"`
	// targetLabel is the target label name where the result is written.
	// Required for `Replace`, `HashMod`, `Lowercase`, `Uppercase`,`KeepEqual` and `DropEqual` actions.
	// Must be between 1 and 128 characters in length when specified.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	TargetLabel string `json:"targetLabel,omitempty"`
	// replacement value against which a Replace action is performed if the
	// regular expression matches. Regex capture groups are available (e.g., $1, $2).
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is "$1" (the first capture group).
	// Setting to an empty string ("") explicitly clears the target label value.
	// Must be at most 255 characters in length.
	// +optional
	// +kubebuilder:validation:MaxLength=255
	Replacement *string `json:"replacement,omitempty"`
	// action is the action to perform on the matched labels.
	// Valid actions are:
	//   - Replace: Replaces the value of targetLabel with replacement, using regex capture groups.
	//   - Keep: Keeps only metrics where regex matches the source labels.
	//   - Drop: Drops metrics where regex matches the source labels.
	//   - HashMod: Sets targetLabel to the hash modulus of the source labels.
	//   - LabelMap: Copies labels matching regex to new label names derived from replacement.
	//   - LabelDrop: Drops labels matching regex.
	//   - LabelKeep: Keeps only labels matching regex.
	// +required
	Action RelabelAction `json:"action,omitempty"`
}

// TLSConfig represents TLS configuration for Alertmanager connections.
// At least one TLS configuration option must be specified.
// For mutual TLS (mTLS), both cert and key must be specified together, or both omitted.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="(has(self.cert) && has(self.key)) || (!has(self.cert) && !has(self.key))",message="cert and key must both be specified together for mutual TLS, or both be omitted"
type TLSConfig struct {
	// ca is an optional CA certificate to use for TLS connections.
	// When omitted, the system's default CA bundle is used.
	// +optional
	CA SecretKeySelector `json:"ca,omitempty,omitzero"`
	// cert is an optional client certificate to use for mutual TLS connections.
	// When omitted, no client certificate is presented.
	// +optional
	Cert SecretKeySelector `json:"cert,omitempty,omitzero"`
	// key is an optional client key to use for mutual TLS connections.
	// When omitted, no client key is used.
	// +optional
	Key SecretKeySelector `json:"key,omitempty,omitzero"`
	// serverName is an optional server name to use for TLS connections.
	// When specified, must be a valid DNS subdomain as per RFC 1123.
	// When omitted, the server name is derived from the URL.
	// Must be between 1 and 253 characters in length.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="must be a valid DNS subdomain (lowercase alphanumeric characters, '-' or '.', start and end with alphanumeric)"
	ServerName string `json:"serverName,omitempty"`
	// certificateVerification determines the policy for TLS certificate verification.
	// Allowed values are "Verify" (performs certificate verification, secure) and "SkipVerify" (skips verification, insecure).
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is "Verify".
	// +optional
	CertificateVerification CertificateVerificationType `json:"certificateVerification,omitempty"`
}

// CertificateVerificationType defines the TLS certificate verification policy.
// +kubebuilder:validation:Enum=Verify;SkipVerify
type CertificateVerificationType string

const (
	// CertificateVerificationVerify performs certificate verification (secure, recommended).
	CertificateVerificationVerify CertificateVerificationType = "Verify"
	// CertificateVerificationSkipVerify skips certificate verification (insecure, use with caution).
	CertificateVerificationSkipVerify CertificateVerificationType = "SkipVerify"
)

// AuthorizationType defines the type of authentication to use.
// +kubebuilder:validation:Enum=None;BearerToken
type AuthorizationType string

const (
	// AuthorizationTypeNone indicates no authentication.
	AuthorizationTypeNone AuthorizationType = "None"
	// AuthorizationTypeBearerToken indicates bearer token authentication.
	AuthorizationTypeBearerToken AuthorizationType = "BearerToken"
)

// AuthorizationConfig defines the authentication method for Alertmanager connections.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'BearerToken' ? has(self.bearerToken) : !has(self.bearerToken)",message="bearerToken is required when type is BearerToken, and forbidden otherwise"
// +union
type AuthorizationConfig struct {
	// type specifies the authentication type to use.
	// Valid values are "None" (no authentication) and "BearerToken" (bearer token authentication).
	// When set to None, no authentication credentials are sent.
	// When set to BearerToken, the bearerToken field must be specified.
	// +unionDiscriminator
	// +required
	Type AuthorizationType `json:"type,omitempty"`
	// bearerToken defines the secret reference containing the bearer token.
	// Required when type is "BearerToken", forbidden otherwise.
	// The secret must exist in the openshift-monitoring namespace.
	// +optional
	BearerToken SecretKeySelector `json:"bearerToken,omitempty,omitzero"`
}

// SecretKeySelector selects a key of a Secret in the `openshift-monitoring` namespace.
// +structType=atomic
type SecretKeySelector struct {
	// name is the name of the secret in the `openshift-monitoring` namespace to select from.
	// Must be a valid Kubernetes secret name (lowercase alphanumeric, '-' or '.', start/end with alphanumeric).
	// Must be between 1 and 253 characters in length.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="must be a valid secret name (lowercase alphanumeric characters, '-' or '.', start and end with alphanumeric)"
	Name string `json:"name,omitempty"`
	// key is the key of the secret to select from.
	// Must consist of alphanumeric characters, '-', '_', or '.'.
	// Must be between 1 and 253 characters in length.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-zA-Z0-9._-]+$')",message="must contain only alphanumeric characters, '-', '_', or '.'"
	Key string `json:"key,omitempty"`
}

// VolumeClaimConfig defines the configuration for a PersistentVolumeClaim used for storage.
type VolumeClaimConfig struct {
	// storageClassName is the name of the StorageClass to use for the PersistentVolumeClaim.
	// When omitted, the default StorageClass is used.
	// Must be a valid DNS subdomain name (RFC 1123) when specified.
	// Maximum length is 253 characters.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="storageClassName must be a valid DNS subdomain name"
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
	// sizeInGiB specifies the storage capacity in gibibytes (GiB) for the volume.
	// This determines how much data can be stored.
	// Minimum value is 1 GiB.
	// Maximum value is 16384 GiB (16 TiB).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=16384
	// +required
	SizeInGiB int32 `json:"sizeInGiB,omitempty"`
}

// Retention configures how long Prometheus retains metrics data and how much storage it can use.
// +kubebuilder:validation:MinProperties=1
type Retention struct {
	// durationInDays specifies how many days Prometheus will retain metrics data.
	// Prometheus automatically deletes data older than this duration.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The default value is 15.
	// Minimum value is 1 day.
	// Maximum value is 365 days (1 year).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=365
	// +optional
	DurationInDays int32 `json:"durationInDays,omitempty"`
	// sizeInGiB specifies the maximum storage size in gibibytes (GiB) that Prometheus
	// can use for data blocks and the write-ahead log (WAL).
	// When the limit is reached, Prometheus will delete oldest data first.
	// When omitted, no size limit is enforced and Prometheus uses available PersistentVolume capacity.
	// Minimum value is 1 GiB.
	// Maximum value is 16384 GiB (16 TiB).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=16384
	// +optional
	SizeInGiB int32 `json:"sizeInGiB,omitempty"`
}

// RelabelAction defines the action to perform in a relabeling rule.
// +kubebuilder:validation:Enum=Replace;Keep;Drop;HashMod;LabelMap;LabelDrop;LabelKeep
type RelabelAction string

const (
	// RelabelActionReplace replaces the target label with the replacement value.
	RelabelActionReplace RelabelAction = "Replace"
	// RelabelActionKeep keeps metrics that match the regex.
	RelabelActionKeep RelabelAction = "Keep"
	// RelabelActionDrop drops metrics that match the regex.
	RelabelActionDrop RelabelAction = "Drop"
	// RelabelActionHashMod sets the target label to the modulus of a hash of the source labels.
	RelabelActionHashMod RelabelAction = "HashMod"
	// RelabelActionLabelMap maps label names based on regex matching.
	RelabelActionLabelMap RelabelAction = "LabelMap"
	// RelabelActionLabelDrop removes labels that match the regex.
	RelabelActionLabelDrop RelabelAction = "LabelDrop"
	// RelabelActionLabelKeep removes labels that do not match the regex.
	RelabelActionLabelKeep RelabelAction = "LabelKeep"
)

// CollectionProfile defines the metrics collection profile for Prometheus.
// +kubebuilder:validation:Enum=Full;Minimal
type CollectionProfile string

const (
	// CollectionProfileFull means Prometheus collects all metrics that are exposed by the platform components.
	CollectionProfileFull CollectionProfile = "Full"
	// CollectionProfileMinimal means Prometheus only collects metrics necessary for the default
	// platform alerts, recording rules, telemetry and console dashboards.
	CollectionProfileMinimal CollectionProfile = "Minimal"
)

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
