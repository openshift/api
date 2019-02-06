package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	configv1 "github.com/openshift/api/config/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceServingCertSignerConfig provides information to configure a serving serving cert signing controller
type ServiceServingCertSignerConfig struct {
	metav1.TypeMeta `json:",inline"`
	configv1.GenericControllerConfig `json:",inline"`

	// signer holds the signing information used to automatically sign serving certificates.
	Signer configv1.CertInfo `json:"signer"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIServiceCABundleInjectorConfig provides information to configure an APIService CA Bundle Injector controller
type APIServiceCABundleInjectorConfig struct {
	metav1.TypeMeta `json:",inline"`
	configv1.GenericControllerConfig `json:",inline"`

	// caBundleFile holds the ca bundle containing the service signer CA to apply to APIServices
	CABundleFile string `json:"caBundleFile"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapCABundleInjectorConfig provides information to configure a ConfigMap CA Bundle Injector controller
type ConfigMapCABundleInjectorConfig struct {
	metav1.TypeMeta `json:",inline"`
	configv1.GenericControllerConfig `json:",inline"`

	// caBundleFile holds the ca bundle to apply to ConfigMaps
	CABundleFile string `json:"caBundleFile"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceCA provides information to configure an operator to manage the service cert controllers
type ServiceCA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ServiceCASpec   `json:"spec"`
	Status ServiceCAStatus `json:"status"`
}

type ServiceCASpec struct {
	OperatorSpec `json:",inline"`

	// serviceServingCertSignerConfig holds a sparse config that the user wants for this component.  It only needs to be the overrides from the defaults
	// it will end up overlaying in the following order:
	// 1. hardcoded default
	// 2. this config
	ServiceServingCertSignerConfig runtime.RawExtension `json:"serviceServingCertSignerConfig"`

	// apiServiceCABundleInjectorConfig holds a sparse config that the user wants for this component.  It only needs to be the overrides from the defaults
	// it will end up overlaying in the following order:
	// 1. hardcoded default
	// 2. this config
	APIServiceCABundleInjectorConfig runtime.RawExtension `json:"apiServiceCABundleInjectorConfig"`

	// configMapCABundleInjectorConfig holds a sparse config that the user wants for this component.  It only needs to be the overrides from the defaults
	// it will end up overlaying in the following order:
	// 1. hardcoded default
	// 2. this config
	ConfigMapCABundleInjectorConfig runtime.RawExtension `json:"configMapCABundleInjectorConfig"`
}

type ServiceCAStatus struct {
	OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceCAList is a collection of items
type ServiceCAList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items contains the items
	Items []ServiceCA `json:"items"`
}
