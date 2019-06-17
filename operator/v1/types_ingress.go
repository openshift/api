package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.availableReplicas,selectorpath=.status.labelSelector

// IngressController describes a managed ingress controller for the cluster. The
// controller can service OpenShift Route and Kubernetes Ingress resources.
//
// When an IngressController is created, a new ingress controller deployment is
// created to allow external traffic to reach the services that expose Ingress
// or Route resources. Updating this resource may lead to disruption for public
// facing network connections as a new ingress controller revision may be rolled
// out.
//
// https://kubernetes.io/docs/concepts/services-networking/ingress-controllers
//
// Whenever possible, sensible defaults for the platform are used. See each
// field for more details.
type IngressController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the IngressController.
	Spec IngressControllerSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the IngressController.
	Status IngressControllerStatus `json:"status,omitempty"`
}

// IngressControllerSpec is the specification of the desired behavior of the
// IngressController.
type IngressControllerSpec struct {
	// domain is a DNS name serviced by the ingress controller and is used to
	// configure multiple features:
	//
	// * For the LoadBalancerService endpoint publishing strategy, domain is
	//   used to configure DNS records. See endpointPublishingStrategy.
	//
	// * When using a generated default certificate, the certificate will be valid
	//   for domain and its subdomains. See defaultCertificate.
	//
	// * The value is published to individual Route statuses so that end-users
	//   know where to target external DNS records.
	//
	// domain must be unique among all IngressControllers, and cannot be
	// updated.
	//
	// If empty, defaults to ingress.config.openshift.io/cluster .spec.domain.
	//
	// +optional
	Domain string `json:"domain,omitempty"`

	// replicas is the desired number of ingress controller replicas. If unset,
	// defaults to 2.
	//
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// endpointPublishingStrategy is used to publish the ingress controller
	// endpoints to other networks, enable load balancer integrations, etc.
	//
	// If unset, the default is based on
	// infrastructure.config.openshift.io/cluster .status.platform:
	//
	//   AWS:      LoadBalancerService
	//   Azure:    LoadBalancerService
	//   GCP:      LoadBalancerService
	//   Libvirt:  HostNetwork
	//
	// Any other platform types (including None) default to HostNetwork.
	//
	// endpointPublishingStrategy cannot be updated.
	//
	// +optional
	EndpointPublishingStrategy *EndpointPublishingStrategy `json:"endpointPublishingStrategy,omitempty"`

	// defaultCertificate is a reference to a secret containing the default
	// certificate served by the ingress controller. When Routes don't specify
	// their own certificate, defaultCertificate is used.
	//
	// The secret must contain the following keys and data:
	//
	//   tls.crt: certificate file contents
	//   tls.key: key file contents
	//
	// If unset, a wildcard certificate is automatically generated and used. The
	// certificate is valid for the ingress controller domain (and subdomains) and
	// the generated certificate's CA will be automatically integrated with the
	// cluster's trust store.
	//
	// The in-use certificate (whether generated or user-specified) will be
	// automatically integrated with OpenShift's built-in OAuth server.
	//
	// +optional
	DefaultCertificate *corev1.LocalObjectReference `json:"defaultCertificate,omitempty"`

	// namespaceSelector is used to filter the set of namespaces serviced by the
	// ingress controller. This is useful for implementing shards.
	//
	// If unset, the default is no filtering.
	//
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// routeSelector is used to filter the set of Routes serviced by the ingress
	// controller. This is useful for implementing shards.
	//
	// If unset, the default is no filtering.
	//
	// +optional
	RouteSelector *metav1.LabelSelector `json:"routeSelector,omitempty"`

	// nodePlacement enables explicit control over the scheduling of the ingress
	// controller.
	//
	// If unset, defaults are used. See NodePlacement for more details.
	//
	// +optional
	NodePlacement *NodePlacement `json:"nodePlacement,omitempty"`

	// securitySpec specifies settings for securing an IngressController.
	//
	// +optional
	SecuritySpec *SecuritySpec `json:"securitySpec,omitempty"`
}

// NodePlacement describes node scheduling configuration for an ingress
// controller.
type NodePlacement struct {
	// nodeSelector is the node selector applied to ingress controller
	// deployments.
	//
	// If unset, the default is:
	//
	//   beta.kubernetes.io/os: linux
	//   node-role.kubernetes.io/worker: ''
	//
	// If set, the specified selector is used and replaces the default.
	//
	// +optional
	NodeSelector *metav1.LabelSelector `json:"nodeSelector,omitempty"`

	// tolerations is a list of tolerations applied to ingress controller
	// deployments.
	//
	// The default is an empty list.
	//
	// See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// EndpointPublishingStrategyType is a way to publish ingress controller endpoints.
type EndpointPublishingStrategyType string

const (
	// LoadBalancerService publishes the ingress controller using a Kubernetes
	// LoadBalancer Service.
	LoadBalancerServiceStrategyType EndpointPublishingStrategyType = "LoadBalancerService"

	// HostNetwork publishes the ingress controller on node ports where the
	// ingress controller is deployed.
	HostNetworkStrategyType EndpointPublishingStrategyType = "HostNetwork"

	// Private does not publish the ingress controller.
	PrivateStrategyType EndpointPublishingStrategyType = "Private"
)

// LoadBalancerScope is the scope at which a load balancer is exposed.
type LoadBalancerScope string

var (
	// InternalLoadBalancer is a load balancer that is exposed only on the
	// cluster's private network.
	InternalLoadBalancer LoadBalancerScope = "Internal"

	// ExternalLoadBalancer is a load balancer that is exposed on the
	// cluster's public network (which is typically on the Internet).
	ExternalLoadBalancer LoadBalancerScope = "External"
)

// LoadBalancerStrategy holds parameters for a load balancer.
type LoadBalancerStrategy struct {
	// scope indicates the scope at which the load balancer is exposed.
	// Possible values are "External" and "Internal".  The default is
	// "External".
	// +optional
	Scope LoadBalancerScope `json:"scope"`
}

// HostNetworkStrategy holds parameters for the HostNetwork endpoint publishing
// strategy.
type HostNetworkStrategy struct {
}

// PrivateStrategy holds parameters for the Private endpoint publishing
// strategy.
type PrivateStrategy struct {
}

// EndpointPublishingStrategy is a way to publish the endpoints of an
// IngressController, and represents the type and any additional configuration
// for a specific type.
// +union
type EndpointPublishingStrategy struct {
	// type is the publishing strategy to use. Valid values are:
	//
	// * LoadBalancerService
	//
	// Publishes the ingress controller using a Kubernetes LoadBalancer Service.
	//
	// In this configuration, the ingress controller deployment uses container
	// networking. A LoadBalancer Service is created to publish the deployment.
	//
	// See: https://kubernetes.io/docs/concepts/services-networking/#loadbalancer
	//
	// If domain is set, a wildcard DNS record will be managed to point at the
	// LoadBalancer Service's external name. DNS records are managed only in DNS
	// zones defined by dns.config.openshift.io/cluster .spec.publicZone and
	// .spec.privateZone.
	//
	// Wildcard DNS management is currently supported only on the AWS platform.
	//
	// * HostNetwork
	//
	// Publishes the ingress controller on node ports where the ingress controller
	// is deployed.
	//
	// In this configuration, the ingress controller deployment uses host
	// networking, bound to node ports 80 and 443. The user is responsible for
	// configuring an external load balancer to publish the ingress controller via
	// the node ports.
	//
	// * Private
	//
	// Does not publish the ingress controller.
	//
	// In this configuration, the ingress controller deployment uses container
	// networking, and is not explicitly published. The user must manually publish
	// the ingress controller.
	// +unionDiscriminator
	// +optional
	Type EndpointPublishingStrategyType `json:"type"`

	// loadBalancer holds parameters for the load balancer. Present only if
	// type is LoadBalancerService.
	// +optional
	// +nullable
	LoadBalancer *LoadBalancerStrategy `json:"loadBalancer,omitempty"`

	// hostNetwork holds parameters for the HostNetwork endpoint publishing
	// strategy. Present only if type is HostNetwork.
	// +optional
	// +nullable
	HostNetwork *HostNetworkStrategy `json:"hostNetwork,omitempty"`

	// private holds parameters for the Private endpoint publishing
	// strategy. Present only if type is Private.
	// +optional
	// +nullable
	Private *PrivateStrategy `json:"private,omitempty"`
}

var (
	// Available indicates the ingress controller deployment is available.
	IngressControllerAvailableConditionType = "Available"
	// LoadBalancerManaged indicates the management status of any load balancer
	// service associated with an ingress controller.
	LoadBalancerManagedIngressConditionType = "LoadBalancerManaged"
	// LoadBalancerReady indicates the ready state of any load balancer service
	// associated with an ingress controller.
	LoadBalancerReadyIngressConditionType = "LoadBalancerReady"
	// DNSManaged indicates the management status of any DNS records for the
	// ingress controller.
	DNSManagedIngressConditionType = "DNSManaged"
	// DNSReady indicates the ready state of any DNS records for the ingress
	// controller.
	DNSReadyIngressConditionType = "DNSReady"
)

// SecuritySpec defines the schema for securing an IngressController.
type SecuritySpec struct {
	// profile defines the schema for a security profile.
	Profile SecurityProfileSpec `json:"profile"`
}

// SecurityProfileSpec defines the schema for a security profile.
// +union
type SecurityProfileSpec struct {
	// type is one of Old, Intermediate, Modern or Custom. Custom provides
	// the ability to specify individual security profile parameters. Old,
	// Intermediate and Modern are security profiles based on:
	//
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations
	//
	// If unset, Intermediate is used.
	//
	// Updating any fields of a SecurityProfileSpec will trigger a Rolling Update
	// of the IngressController. For more on Rolling Updates see:
	//
	// https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/
	//
	// +unionDiscriminator
	// +optional
	Type SecurityProfileType `json:"type"`
	// old is a security profile based on:
	//
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Old_backward_compatibility
	//
	// and looks like this (yaml):
	//
	//   ciphers:
	//     - ECDHE-ECDSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-AES128-GCM-SHA256
	//     - ECDHE-ECDSA-AES128-GCM-SHA256
	//     - ECDHE-RSA-AES256-GCM-SHA384
	//     - ECDHE-ECDSA-AES256-GCM-SHA384
	//     - DHE-RSA-AES128-GCM-SHA256
	//     - DHE-DSS-AES128-GCM-SHA256
	//     - kEDH+AESGCM
	//     - ECDHE-RSA-AES128-SHA256
	//     - ECDHE-ECDSA-AES128-SHA256
	//     - ECDHE-RSA-AES128-SHA
	//     - ECDHE-ECDSA-AES128-SHA
	//     - ECDHE-RSA-AES256-SHA384
	//     - ECDHE-ECDSA-AES256-SHA384
	//     - ECDHE-RSA-AES256-SHA
	//     - ECDHE-ECDSA-AES256-SHA
	//     - DHE-RSA-AES128-SHA256
	//     - DHE-RSA-AES128-SHA
	//     - DHE-DSS-AES128-SHA256
	//     - DHE-RSA-AES256-SHA256
	//     - DHE-DSS-AES256-SHA
	//     - DHE-RSA-AES256-SHA
	//     - ECDHE-RSA-DES-CBC3-SHA
	//     - ECDHE-ECDSA-DES-CBC3-SHA
	//     - EDH-RSA-DES-CBC3-SHA
	//     - AES128-GCM-SHA256
	//     - AES256-GCM-SHA384
	//     - AES128-SHA256
	//     - AES256-SHA256
	//     - AES128-SHA
	//     - AES256-SHA
	//     - AES
	//     - DES-CBC3-SHA
	//     - HIGH
	//     - SEED
	//     - !aNULL
	//     - !eNULL
	//     - !EXPORT
	//     - !DES
	//     - !RC4
	//     - !MD5
	//     - !PSK
	//     - !RSAPSK
	//     - !aDH
	//     - !aECDH
	//     - !EDH-DSS-DES-CBC3-SHA
	//     - !KRB5-DES-CBC3-SHA
	//     - !SRP
	//   securityProtocol:
	//     minimumVersion: TLSv1.0
	//     maximumVersion: TLSv1.2
	//   dhParamSize: 1024
	//
	// +optional
	// +nullable
	Old *OldSecurityProfile `json:"old,omitempty"`
	// intermediate is a security profile based on:
	//
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	// and looks like this (yaml):
	//
	//   ciphers:
	//     - ECDHE-ECDSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-CHACHA20-POLY1305
	//     - ECDHE-ECDSA-AES128-GCM-SHA256
	//     - ECDHE-RSA-AES128-GCM-SHA256
	//     - ECDHE-ECDSA-AES256-GCM-SHA384
	//     - ECDHE-RSA-AES256-GCM-SHA384
	//     - DHE-RSA-AES128-GCM-SHA256
	//     - DHE-RSA-AES256-GCM-SHA384
	//     - ECDHE-ECDSA-AES128-SHA256
	//     - ECDHE-RSA-AES128-SHA256
	//     - ECDHE-ECDSA-AES128-SHA
	//     - ECDHE-RSA-AES256-SHA384
	//     - ECDHE-RSA-AES128-SHA
	//     - ECDHE-ECDSA-AES256-SHA384
	//     - ECDHE-ECDSA-AES256-SHA
	//     - ECDHE-RSA-AES256-SHA
	//     - DHE-RSA-AES128-SHA256
	//     - DHE-RSA-AES128-SHA
	//     - DHE-RSA-AES256-SHA256
	//     - DHE-RSA-AES256-SHA
	//     - ECDHE-ECDSA-DES-CBC3-SHA
	//     - ECDHE-RSA-DES-CBC3-SHA
	//     - EDH-RSA-DES-CBC3-SHA
	//     - AES128-GCM-SHA256
	//     - AES256-GCM-SHA384
	//     - AES128-SHA256
	//     - AES256-SHA256
	//     - AES128-SHA
	//     - AES256-SHA
	//     - DES-CBC3-SHA
	//     - !DSS
	//   securityProtocol:
	//     minimumVersion: TLSv1.0
	//     maximumVersion: TLSv1.2
	//   dhParamSize: 2048
	//
	// +optional
	// +nullable
	Intermediate *IntermediateSecurityProfile `json:"intermediate,omitempty"`
	// modern is a security profile based on:
	//
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
	//
	// and looks like this (yaml):
	//
	//   ciphers:
	//     - ECDHE-ECDSA-AES256-GCM-SHA384
	//     - ECDHE-RSA-AES256-GCM-SHA384
	//     - ECDHE-ECDSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-CHACHA20-POLY1305
	//     - ECDHE-ECDSA-AES128-GCM-SHA256
	//     - ECDHE-RSA-AES128-GCM-SHA256
	//     - ECDHE-ECDSA-AES256-SHA384
	//     - ECDHE-RSA-AES256-SHA384
	//     - ECDHE-ECDSA-AES128-SHA256
	//     - ECDHE-RSA-AES128-SHA256
	//   securityProtocol:
	//     minimumVersion: TLSv1.2
	//     maximumVersion: TLSv1.2
	//   dhParamSize: 2048
	//
	// +optional
	// +nullable
	Modern *ModernSecurityProfile `json:"modern,omitempty"`
	// custom is a user-defined security profile. An example custom profile
	// looks like this:
	//
	//   ciphers:
	//     - ECDHE-ECDSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-CHACHA20-POLY1305
	//     - ECDHE-RSA-AES128-GCM-SHA256
	//     - ECDHE-ECDSA-AES128-GCM-SHA256
	//   securityProtocol:
	//     minimumVersion: TLSv1.1
	//     maximumVersion: TLSv1.2
	//   dhParamSize: 1024
	//
	// Any unset fields of a custom security profile default to the "Intermediate"
	// security profile.
	//
	// +optional
	// +nullable
	Custom *CustomSecurityProfile `json:"custom,omitempty"`
}

type OldSecurityProfile struct{}
type IntermediateSecurityProfile struct{}
type ModernSecurityProfile struct{}

// CustomSecurityProfile defines the schema for a custom security profile.
type CustomSecurityProfile struct {
	SecurityProfile `json:",inline"`
}

// SecurityProfileType defines a security profile type.
type SecurityProfileType string

const (
	// Old is a security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Old_backward_compatibility
	SecurityProfileOldType SecurityProfileType = "Old"
	// Intermediate is a security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	SecurityProfileIntermediateType SecurityProfileType = "Intermediate"
	// Modern is a security profile based on:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
	SecurityProfileModernType SecurityProfileType = "Modern"
	// Custom is a security profile that allows for user-defined parameters.
	SecurityProfileCustomType SecurityProfileType = "Custom"
)

// SecurityProfile defines the schema for a security profile.
type SecurityProfile struct {
	// ciphers is used to specify the cipher algorithms that are negotiated
	// during the SSL/TLS handshake with an IngressController. Preface a cipher
	// with a "!" to disable a specific cipher from being negotiated. For example,
	// to use DES-CBC3-SHA but not EDH-DSS-DES-CBC3-SHA (yaml):
	//
	//   ciphers:
	//     - DES-CBC3-SHA
	//     - !EDH-DSS-DES-CBC3-SHA
	//
	// If unset, Ciphersuites are determined by the security profile type.
	//
	// +optional
	Ciphers *[]string `json:"ciphers,omitempty"`
	// securityProtocol is used to specify one or more encryption protocols
	// that are negotiated during the SSL/TLS handshake with the IngressController.
	//
	// If unset, the security protocol is determined by the security profile type.
	//
	// +optional
	SecurityProtocol *SecurityProtocol `json:"securityProtocol,omitempty"`
	// dhParamSize sets the maximum size of the Diffie-Hellman parameters used for generating
	// the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange. The final size
	// will try to match the size of the server's RSA (or DSA) key (e.g, a 2048 bits temporary
	// DH key for a 2048 bits RSA key), but will not exceed this maximum value.
	//
	// If unset, the DH parameter size is determined by the security profile type.
	//
	//   SecurityProfileType Modern:       DHParamSize2048
	//   SecurityProfileType Intermediate: DHParamSize2048
	//   SecurityProfileType Old:          DHParamSize1024
	//
	// Supported DH Parameter sizes are:
	//
	//   "2048": A Diffie-Hellman parameter of 2048 bits.
	//   "1024": A Diffie-Hellman parameter of 1024 bits.
	//
	// +optional
	DHParamSize *DHParamSize `json:"dhParamSize,omitempty"`
}

// SecurityProtocol defines one or more security protocols used by
// an IngressController to secure network connections.
type SecurityProtocol struct {
	// minimumVersion enforces use of the specified SecurityProtocolVersion or newer
	// on SSL connections initiated by an IngressController. minimumVersion must be
	// lower than or equal to maximumVersion.
	//
	// If unset and maximumVersion is set, minimumVersion will be set
	// to maximumVersion. If minimumVersion and maximumVersion are unset,
	// the minimum version is determined by the security profile type.
	//
	//   SecurityProfileType Modern:       SecurityProtocolTLS12Version
	//   SecurityProfileType Intermediate: SecurityProtocolTLS10Version
	//   SecurityProfileType Old:          SecurityProtocolTLS10Version
	//
	// Supported minimum versions are:
	//
	//   "TLSv1.3": Version 1.3 of the TLS security protocol used for securing
	//   IngressController network connections.
	//   "TLSv1.2": Version 1.2 of the TLS security protocol used for securing
	//   IngressController network connections.
	//   "TLSv1.1": Version 1.1 of the TLS security protocol used for securing
	//   IngressController network connections.
	//   "TLSv1.0": Version 1.0 of the TLS security protocol used for securing
	//   IngressController network connections.
	//
	// +optional
	MinimumVersion *SecurityProtocolVersion `json:"minimumVersion,omitempty"`
	// maximumVersion enforces use of the specified SecurityProtocolVersion or older
	// on SSL connections initiated by an IngressController. maximumVersion must be
	// higher than or equal to minimumVersion.
	//
	// If unset and minimumVersion is set, maximumVersion will be set
	// to minimumVersion. If minimumVersion and maximumVersion are unset,
	// the maximum version is determined by the security profile type.
	//
	//   SecurityProfileType Modern:       SecurityProtocolTLS12Version
	//   SecurityProfileType Intermediate: SecurityProtocolTLS12Version
	//   SecurityProfileType Old:          SecurityProtocolTLS12Version
	//
	// Supported maximum versions are the same as minimum versions.
	//
	// +optional
	MaximumVersion *SecurityProtocolVersion `json:"maximumVersion,omitempty"`
}

// SecurityProtocolVersion is a way to specify the TLS security protocol version
// used for securing IngressController network connections.
type SecurityProtocolVersion string

const (
	// TLSv1.0 is version 1.0 of the TLS security protocol used for securing
	// IngressController network connections.
	SecurityProtocolTLS10Version SecurityProtocolVersion = "TLSv1.0"
	// TLSv1.1 is version 1.1 of the TLS security protocol used for securing
	// IngressController network connections.
	SecurityProtocolTLS11Version SecurityProtocolVersion = "TLSv1.1"
	// TLSv1.2 is version 1.2 of the TLS security protocol used for securing
	// IngressController network connections.
	SecurityProtocolTLS12Version SecurityProtocolVersion = "TLSv1.2"
	// TLSv1.3 is version 1.3 of the TLS security protocol used for securing
	// IngressController network connections.
	SecurityProtocolTLS13Version SecurityProtocolVersion = "TLSv1.3"
)

// DHParamSize sets the maximum size of the Diffie-Hellman parameters used for
// generating the ephemeral/temporary Diffie-Hellman key.
type DHParamSize string

const (
	// 1024 is a Diffie-Hellman parameter of 1024 bits.
	DHParamSize1024 DHParamSize = "1024"
	// 2048 is a Diffie-Hellman parameter of 2048 bits.
	DHParamSize2048 DHParamSize = "2048"
)

// IngressControllerStatus defines the observed status of the IngressController.
type IngressControllerStatus struct {
	// availableReplicas is number of observed available replicas according to the
	// ingress controller deployment.
	AvailableReplicas int32 `json:"availableReplicas"`

	// selector is a label selector, in string format, for ingress controller pods
	// corresponding to the IngressController. The number of matching pods should
	// equal the value of availableReplicas.
	Selector string `json:"selector"`

	// domain is the actual domain in use.
	Domain string `json:"domain"`

	// endpointPublishingStrategy is the actual strategy in use.
	EndpointPublishingStrategy *EndpointPublishingStrategy `json:"endpointPublishingStrategy,omitempty"`

	// securityProfile is the actual security profile in use.
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`

	// conditions is a list of conditions and their status.
	//
	// Available means the ingress controller deployment is available and
	// servicing route and ingress resources (i.e, .status.availableReplicas
	// equals .spec.replicas)
	//
	// There are additional conditions which indicate the status of other
	// ingress controller features and capabilities.
	//
	//   * LoadBalancerManaged
	//   - True if the following conditions are met:
	//     * The endpoint publishing strategy requires a service load balancer.
	//   - False if any of those conditions are unsatisfied.
	//
	//   * LoadBalancerReady
	//   - True if the following conditions are met:
	//     * A load balancer is managed.
	//     * The load balancer is ready.
	//   - False if any of those conditions are unsatisfied.
	//
	//   * DNSManaged
	//   - True if the following conditions are met:
	//     * The endpoint publishing strategy and platform support DNS.
	//     * The ingress controller domain is set.
	//     * dns.config.openshift.io/cluster configures DNS zones.
	//   - False if any of those conditions are unsatisfied.
	//
	//   * DNSReady
	//   - True if the following conditions are met:
	//     * DNS is managed.
	//     * DNS records have been successfully created.
	//   - False if any of those conditions are unsatisfied.
	Conditions []OperatorCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// IngressControllerList contains a list of IngressControllers.
type IngressControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngressController `json:"items"`
}
