package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceServingCertSignerConfig provides information to configure a serving serving cert signing controller
type ServiceServingCertSignerConfig struct {
	metav1.TypeMeta `json:",inline"`

	// ServingInfo is the HTTP serving information for the controller's endpoints
	ServingInfo configv1.HTTPServingInfo `json:"servingInfo"`

	// authentication allows configuration of authentication for the endpoints
	Authentication DelegatedAuthentication `json:"authentication,omitempty"`
	// authorization allows configuration of authentication for the endpoints
	Authorization DelegatedAuthorization `json:"authorization,omitempty"`

	// Signer holds the signing information used to automatically sign serving certificates.
	Signer configv1.CertInfo `json:"signer"`
}

// DelegatedAuthentication allows authentication to be disabled.
type DelegatedAuthentication struct {
	// disabled indicates that authentication should be disabled.  By default it will use delegated authentication.
	Disabled bool `json:"disabled,omitempty"`
}

// DelegatedAuthorization allows authorization to be disabled.
type DelegatedAuthorization struct {
	// disabled indicates that authorization should be disabled.  By default it will use delegated authorization.
	Disabled bool `json:"disabled,omitempty"`
}
