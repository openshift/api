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
	// prometheusK8sConfig provides configuration options for the Prometheus instance
	// Prometheus deployment, pod scheduling, resource allocation, retention policies, and external integrations.
	// prometheusK8sConfig is optional.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	PrometheusK8sConfig PrometheusK8sConfig `json:"prometheusK8sConfig,omitempty,omitzero"`
	// metricsServerConfig is an optional field that can be used to configure the Kubernetes Metrics Server that runs in the openshift-monitoring namespace.
	// Specifically, it can configure how the Metrics Server instance is deployed, pod scheduling, its audit policy and log verbosity.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	MetricsServerConfig MetricsServerConfig `json:"metricsServerConfig,omitempty,omitzero"`
	// kubeStateMetricsConfig is an optional field that can be used to configure the kube-state-metrics agent that runs in the openshift-monitoring namespace.
	// KubeStateMetrics (KSM) is a service that listens to the Kubernetes API server
	// and generates metrics about the state of Kubernetes objects. Metrics reflect
	// the raw, unmodified data of objects such as Deployments, Nodes, and Pods.
	// Specifically, it can configure how the kube-state-metrics instance is deployed, pod scheduling, and resource allocation.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// +optional
	KubeStateMetricsConfig KubeStateMetricsConfig `json:"kubeStateMetricsConfig,omitempty,omitzero"`
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

// PrometheusK8sConfig provides configuration options for the Prometheus instance
// Use this configuration to control
// Prometheus deployment, pod scheduling, resource allocation, retention policies, and external integrations.
// +kubebuilder:validation:MinProperties=1
type PrometheusK8sConfig struct {
	// additionalAlertmanagerConfigs configures additional Alertmanager instances that receive alerts from
	// the Prometheus component. By default, no additional Alertmanager instances are configured.
	// +optional
	AdditionalAlertmanagerConfigs []AdditionalAlertmanagerConfig `json:"additionalAlertmanagerConfigs,omitempty"`
	// enforcedBodySizeLimit enforces a body size limit for Prometheus scraped metrics. If a scraped
	// target's body response is larger than the limit, the scrape will fail.
	// The following values are valid:
	// an empty value to specify no limit,
	// a numeric value in Prometheus size format (such as `64MB`), or
	// the string `automatic`, which indicates that the limit will be
	// automatically calculated based on cluster capacity.
	// The default value is empty, which indicates no limit.
	// +optional
	EnforcedBodySizeLimit string `json:"enforcedBodySizeLimit,omitempty"`
	// externalLabels defines labels to be added to any time series or alerts when
	// communicating with external systems such as federation, remote storage,
	// and Alertmanager. By default, no labels are added.
	// +optional
	ExternalLabels ExternalLabels `json:"externalLabels,omitempty"`
	// logLevel defines the log level setting for Prometheus.
	// The possible values are: `error`, `warn`, `info`, and `debug`.
	// The default value is `info`.
	// +optional
	LogLevel string `json:"logLevel,omitempty"`
	// nodeSelector defines the nodes on which the Pods are scheduled.
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
	// +optional
	QueryLogFile string `json:"queryLogFile,omitempty"`
	// remoteWrite defines the remote write configuration, including URL, authentication,
	// and relabeling settings.
	// +optional
	RemoteWrite []RemoteWriteSpec `json:"remoteWrite,omitempty"`
	// resources defines resource requests and limits for the Prometheus container.
	// +optional
	Resources []ContainerResource `json:"resources,omitempty"`
	// retention defines the duration for which Prometheus retains data.
	// This definition must be specified using the following regular
	// expression pattern: `[0-9]+(ms|s|m|h|d|w|y)` (ms = milliseconds,
	// s= seconds,m = minutes, h = hours, d = days, w = weeks, y = years).
	// The default value is `15d`.
	// +optional
	Retention string `json:"retention,omitempty"`
	// retentionSize defines the maximum amount of disk space used by data blocks plus the
	// write-ahead log (WAL).
	// Supported values are `B`, `KB`, `KiB`, `MB`, `MiB`, `GB`, `GiB`, `TB`,
	// `TiB`, `PB`, `PiB`, `EB`, and `EiB`.
	// By default, no limit is defined.
	// +optional
	RetentionSize string `json:"retentionSize,omitempty"`
	// tolerations defines tolerations for the pods.
	// +optional
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// topologySpreadConstraints defines the pod's topology spread constraints.
	// +optional
	// +kubebuilder:validation:MaxItems=10
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// collectionProfile defines the metrics collection profile that Prometheus uses to collect
	// metrics from the platform components. Supported values are `full` or
	// `minimal`. In the `full` profile (default), Prometheus collects all
	// metrics that are exposed by the platform components. In the `minimal`
	// profile, Prometheus only collects metrics necessary for the default
	// platform alerts, recording rules, telemetry and console dashboards.
	// +optional
	CollectionProfile CollectionProfile `json:"collectionProfile,omitempty"`
	// volumeClaimTemplate defines persistent storage for Prometheus. Use this setting to
	// configure the persistent volume claim, including storage class,
	// volume size and name.
	// +optional
	VolumeClaimTemplate *v1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
}

// AdditionalAlertmanagerConfig represents configuration for additional Alertmanager instances.
// The `AdditionalAlertmanagerConfig` resource defines settings for how a
// component communicates with additional Alertmanager instances.
type AdditionalAlertmanagerConfig struct {
	// apiVersion defines the API version of Alertmanager.
	// `v1` is no longer supported, `v2` is set as the default value.
	// +required
	APIVersion string `json:"apiVersion"`
	// bearerToken defines the secret key reference containing the bearer token
	// to use when authenticating to Alertmanager.
	// +optional
	BearerToken *v1.SecretKeySelector `json:"bearerToken,omitempty"`
	// pathPrefix defines the path prefix to add in front of the push endpoint path.
	// +optional
	PathPrefix string `json:"pathPrefix,omitempty"`
	// scheme defines the URL scheme to use when communicating with Alertmanager
	// instances.
	// Possible values are `http` or `https`. The default value is `http`.
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// staticConfigs is a list of statically configured Alertmanager endpoints in the form
	// of `<hosts>:<port>`.
	// +optional
	StaticConfigs []string `json:"staticConfigs,omitempty"`
	// timeout defines the timeout value used when sending alerts.
	// +optional
	Timeout *string `json:"timeout,omitempty"`
	// tlsConfig defines the TLS settings to use for Alertmanager connections.
	// +optional
	TLSConfig TLSConfig `json:"tlsConfig,omitempty"`
}

// ExternalLabels represents labels to be added to time series and alerts.
type ExternalLabels struct {
	// labels is a map of label names to label values.
	// +required
	Labels map[string]string `json:"labels"`
}

// RemoteWriteSpec represents configuration for remote write endpoints.
type RemoteWriteSpec struct {
	// url is the URL of the remote write endpoint.
	// +required
	URL string `json:"url"`
	// name is the name of the remote write configuration.
	// +optional
	Name string `json:"name,omitempty"`
	// remoteTimeout is the timeout for requests to the remote write endpoint.
	// +optional
	RemoteTimeout string `json:"remoteTimeout,omitempty"`
	// writeRelabelConfigs is a list of relabeling rules to apply before sending data to the remote endpoint.
	// +optional
	WriteRelabelConfigs []RelabelConfig `json:"writeRelabelConfigs,omitempty"`
}

// RelabelConfig represents a relabeling rule.
type RelabelConfig struct {
	// sourceLabels is a list of source label names.
	// +optional
	SourceLabels []string `json:"sourceLabels,omitempty"`
	// separator is the separator used to join source label values.
	// +optional
	Separator string `json:"separator,omitempty"`
	// regex is the regular expression to match against the concatenated source label values.
	// +optional
	Regex string `json:"regex,omitempty"`
	// targetLabel is the target label name.
	// +optional
	TargetLabel string `json:"targetLabel,omitempty"`
	// replacement is the replacement value for the target label.
	// +optional
	Replacement string `json:"replacement,omitempty"`
	// action is the action to perform.
	// +optional
	Action string `json:"action,omitempty"`
}

// TLSConfig represents TLS configuration for Alertmanager connections.
type TLSConfig struct {
	// ca is the CA certificate to use for TLS connections.
	// +optional
	CA *v1.SecretKeySelector `json:"ca,omitempty"`
	// cert is the client certificate to use for TLS connections.
	// +optional
	Cert *v1.SecretKeySelector `json:"cert,omitempty"`
	// key is the client key to use for TLS connections.
	// +optional
	Key *v1.SecretKeySelector `json:"key,omitempty"`
	// serverName is the server name to use for TLS connections.
	// +optional
	ServerName string `json:"serverName,omitempty"`
	// insecureSkipVerify determines whether to skip TLS certificate verification.
	// +optional
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

// CollectionProfile defines the metrics collection profile for Prometheus.
// +kubebuilder:validation:Enum=full;minimal
type CollectionProfile string

const (
	// CollectionProfileFull means Prometheus collects all metrics that are exposed by the platform components.
	CollectionProfileFull CollectionProfile = "full"
	// CollectionProfileMinimal means Prometheus only collects metrics necessary for the default
	// platform alerts, recording rules, telemetry and console dashboards.
	CollectionProfileMinimal CollectionProfile = "minimal"
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

// kubeStateMetricsConfig (KSM) is a service that listens to the Kubernetes API server
// and generates metrics about the state of Kubernetes objects. Metrics reflect
// the raw, unmodified data of objects such as Deployments, Nodes, and Pods.
// Specifically, it can configure how the kube-state-metrics instance is deployed, pod scheduling, and resource allocation.
// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
// +optional
// +kubebuilder:validation:MinProperties=1
type KubeStateMetricsConfig struct {
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
	// resources defines the compute resource requests and limits for the KubeStateMetrics container.
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
	// topologySpreadConstraints defines rules for how KubeStateMetrics Pods should be distributed
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
