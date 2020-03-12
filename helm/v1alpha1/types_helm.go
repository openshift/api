package v1alpha1

import (
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:plural=helmchartrepositories

// HelmChartRepository holds cluster-wide configuration for proxied Helm chart repository
type HelmChartRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec HelmChartRepositorySpec `json:"spec"`

	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status HelmChartRepositoryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HelmChartRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HelmChartRepository `json:"items"`
}

// Helm chart repository exposed within the cluster
type HelmChartRepositorySpec struct {

	// Chart repository URL
	// +kubebuilder:validation:Pattern=`^https?:\/\/`
	URL string `json:"url"`

	// Optional associated human readable repository name, it can be used by UI for displaying purposes
	// +kubebuilder:validation:MinLength=1
	// +optional
	DisplayName string `json:"name,omitempty"`

	// Optional human readable repository description, it can be used by UI for displaying purposes
	// +kubebuilder:validation:MinLength=1
	// +optional
	Description string `json:"description,omitempty"`

	// ca is an optional reference to a config map by name containing the PEM-encoded CA bundle.
	// It is used as a trust anchor to validate the TLS certificate presented by the remote server.
	// The key "ca.crt" is used to locate the data.
	// If empty, the default system roots are used.
	// The namespace for this config map is openshift-config.
	// +optional
	CA *configv1.ConfigMapNameReference `json:"ca,omitempty"`

	// tlsClientCert is an optional reference to a secret by name that contains the
	// PEM-encoded TLS client certificate to present when connecting to the server.
	// The key "client.crt" is used to locate the data.
	// The namespace for this secret is openshift-config.
	// +optional
	TLSClientCert *configv1.SecretNameReference `json:"tlsClientCert,omitempty"`

	// tlsClientKey is an optional reference to a secret by name that contains the
	// PEM-encoded TLS private key for the client certificate referenced in tlsClientCert.
	// The key "client.key" is used to locate the data.
	// The namespace for this secret is openshift-config.
	// +optional
	TLSClientKey *configv1.SecretNameReference `json:"tlsClientKey,omitempty"`

	// Skip verification of the chart repo certificate
	// +optional
	InsecureSkipTLSVerify bool `json:"insecure_skip_tls_verify,omitempty"`

	// Optional Username used for authenticating access to the chart repository
	// +optional
	Username *string `json:"username,omitempty"`

	// Password is an optional reference to a secret by name that contains
	// the password used for authenticating access to the chart repository
	// The key "password" is used to locate the data.
	// The namespace for this secret is openshift-config.
	// +optional
	Password *configv1.SecretNameReference `json:"password,omitempty"`
}

type HelmChartRepositoryStatus struct {

	// conditions is a list of conditions and their status
	// +optional
	Conditions []HelmChartRepositoryCondition `json:"conditions,omitempty"`
}

// HelmChartRepositoryCondition is just the standard condition fields.
type HelmChartRepositoryCondition struct {
	Type               string          `json:"type"`
	Status             ConditionStatus `json:"status"`
	LastTransitionTime metav1.Time     `json:"lastTransitionTime,omitempty"`
	Reason             string          `json:"reason,omitempty"`
	Message            string          `json:"message,omitempty"`
}

type ConditionStatus string
