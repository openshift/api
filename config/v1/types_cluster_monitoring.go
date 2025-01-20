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

package v1

import (
	v1 "k8s.io/api/core/v1"
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
	// alertmanagerMainConfig defines settings for the Alertmanager component in the `openshift-monitoring` namespace.
	// +required
	AlertmanagerMainConfig AlertmanagerMainConfig `json:"alertmanagerMainConfig"`
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

// The `AlertmanagerMainConfig` resource defines settings for the
// Alertmanager component in the `openshift-monitoring` namespace.
type AlertmanagerMainConfig struct {
	// mode enables or disables the main Alertmanager instance. in the `openshift-monitoring` namespace
	// Allowed values are "Enabled", "Disabled".
	// +kubebuilder:validation:Enum:=Enabled;Disabled;""
	// +required
	Mode AlertManagerMode `json:"mode"`
	// userMode enables or disables user-defined namespaces
	// to be selected for `AlertmanagerConfig` lookups. This setting only
	// applies if the user workload monitoring instance of Alertmanager
	// is not enabled.
	// +required
	UserMode UserAlertManagerMode `json:"userMode"`
	// logLevel Defines the log level setting for Alertmanager.
	// The possible values are: `Error`, `Warn`, `Info`, `Debug`.
	// The default value is `Info`.
	// +optional
	// +kubebuilder:default=Info
	LogLevel string `json:"logLevel,omitempty"`
	// nodeSelector Defines the nodes on which the Pods are scheduled.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// resources Defines resource requests and limits for the Alertmanager container.
	// +optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`
	// secrets Defines a list of secrets that need to be mounted into the Alertmanager.
	// The secrets must reside within the same namespace as the Alertmanager object.
	// They will be added as volumes named secret-<secret-name> and mounted at
	// /etc/alertmanager/secrets/<secret-name> within the 'alertmanager' container of
	// the Alertmanager Pods.
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// tolerations Defines tolerations for the pods.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// topologySpreadConstraints Defines a pod's topology spread constraints.
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// volumeClaimTemplate Defines persistent storage for Alertmanager. Use this setting to
	// configure the persistent volume claim, including storage class, volume
	// size, and name.
	// +optional
	VolumeClaimTemplate v1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
}

// AlertmanagerMode defines mode for AlertManager instance
// +kubebuilder:validation:Enum="";Enabled;Disabled
type AlertManagerMode string

const (
	// AlertManagerEnable enables the main Alertmanager instance. in the `openshift-monitoring` namespace
	AlertManagerEnabled AlertManagerMode = "Enabled"
	// AlertManagerDisabled enables the main Alertmanager instance. in the `openshift-monitoring` namespace
	AlertManagerDisabled AlertManagerMode = "Disabled"
)

// UserAlertManagerMode defines mode for user-defines namespaced
// +kubebuilder:validation:Enum="";Enabled;Disabled
type UserAlertManagerMode string

const (
	// AlertManagerEnabled enables user-defined namespaces to be selected for `AlertmanagerConfig` lookups. This setting only
	// applies if the user workload monitoring instance of Alertmanager is not enabled.
	UserAlertManagerEnabled UserAlertManagerMode = "Enabled"
	// AlertManagerDisabled disables user-defined namespaces to be selected for `AlertmanagerConfig` lookups. This setting only
	// applies if the user workload monitoring instance of Alertmanager is not enabled.
	UserAlertManagerDisabled UserAlertManagerMode = "Disabled"
)

// +kubebuilder:validation:Enum="";Error;Warn;Info;Debug
type LogLevel string

var (
	Error LogLevel = "error"

	Warn LogLevel = "warn"

	Info LogLevel = "info"

	Debug LogLevel = "debug"
)
