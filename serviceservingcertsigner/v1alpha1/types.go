package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceServingCertSignerConfig provides information to configure a serving serving cert signing controller
type ServiceServingCertSignerConfig struct {
	metav1.TypeMeta `json:",inline"`

	// bootstrapService is the location of a service representing the controller. This controller is special, it
	// can mint it's own, valid serving certificate.  Only used if specified.
	BootstrapService ServiceLocation `json:"bootstrapService"`

	// ServingInfo is the HTTP serving information for the controller's endpoints
	ServingInfo configv1.HTTPServingInfo `json:"servingInfo" protobuf:"bytes,1,opt,name=servingInfo"`

	// authentication allows configuration of authentication for the endpoints
	Authentication DelegatedAuthentication `json:"authentication,omitempty" protobuf:"bytes,2,opt,name=authentication"`
	// authorization allows configuration of authentication for the endpoints
	Authorization DelegatedAuthorization `json:"authorization,omitempty" protobuf:"bytes,3,opt,name=authorization"`

	// Signer holds the signing information used to automatically sign serving certificates.
	Signer configv1.CertInfo `json:"signer" protobuf:"bytes,4,opt,name=signer"`
}

// ServiceLocation is the location of a service
type ServiceLocation struct {
	// namespace of the service to sign the self-serving certificate for
	Namespace string `json:"namespace"`
	// name of the service to sign the self-serving certificate for
	Name string `json:"name"`
}

// DelegatedAuthentication allows authentication to be disabled.
type DelegatedAuthentication struct {
	// disabled indicates that authentication should be disabled.  By default it will use delegated authentication.
	Disabled bool `json:"disabled,omitempty" protobuf:"varint,1,opt,name=disabled"`
}

// DelegatedAuthorization allows authorization to be disabled.
type DelegatedAuthorization struct {
	// disabled indicates that authorization should be disabled.  By default it will use delegated authorization.
	Disabled bool `json:"disabled,omitempty" protobuf:"varint,1,opt,name=disabled"`
}
