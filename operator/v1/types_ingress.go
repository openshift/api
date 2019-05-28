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

	// securitySpec specifies settings for securing IngressController connections.
	//
	// If unset, the "Intermediate" security profile is used.
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

// EndpointPublishingStrategy is a way to publish the endpoints of an
// IngressController, and represents the type and any additional configuration
// for a specific type.
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
	Type EndpointPublishingStrategyType `json:"type"`
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

// SecuritySpec defines the settings for securing IngressController connections.
type SecuritySpec struct {
	// profile is one of "Old", "Intermediate", "Modern" or "Custom". "Old",
	// "Intermediate" and "Modern" profiles map to security configurations from [1]:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations
	//
	// When a profile of type "Old", "Intermediate" or "Modern" is set, the CustomSettings
	// field is forbidden.
	//
	// customSettings must be provided if, and only if, profile is "Custom".
	//
	// If unset, the "Intermediate" profile is used.
	//
	Profile SecurityProfileType `json:"profile"`
	// customSettings defines the schema for settings of a "Custom" profile
	// and is ignored unless a "Custom" profile is specified.
	//
	// +optional
	CustomSettings *CustomProfileSettings `json:"customSettings,omitempty"`

}

// SecurityProfileType defines a security profile.
type SecurityProfileType string

const (
	// SecurityProfileOldType is a security profile that maps to:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Old_backward_compatibility
	SecurityProfileOldType SecurityProfileType = "Old"
	// SecurityProfileIntermediateType is a security profile that maps to:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	SecurityProfileIntermediateType SecurityProfileType = "Intermediate"
	// SecurityProfileModernType is a security profile that maps to:
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
	SecurityProfileModernType SecurityProfileType = "Modern"
	// SecurityProfileCustomType is a security profile that allows for custom settings
	// through type CustomProfileSettings.
	SecurityProfileCustomType SecurityProfileType = "Custom"

)

// CustomProfileSettings defines the schema for a custom security profile.
type CustomProfileSettings struct {
	// ciphers is used to specify the cipher algorithms that are negotiated
	// during the SSL/TLS handshake with an IngressController. Each cipher must
	// be an explicit, colon-delimited list of ciphers.
	//
	// If unset, the "Intermediate" Ciphersuites [1] are used:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	// +optional
	Ciphers *string `json:"ciphers,omitempty"`
	// securityProtocol is used to specify one or more encryption protocols
	// that are negotiated during the SSL/TLS handshake with the IngressController.
	//
	// If unset, the "Intermediate" Versions [1] are used:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	//
	// +optional
	SecurityProtocol *SecurityProtocol `json:"securityProtocol,omitempty"`
	// dhParamSize sets the maximum size of the Diffie-Hellman parameters used for generating
	// the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange. The final size
	// will try to match the size of the server's RSA (or DSA) key (e.g, a 2048 bits temporary
	// DH key for a 2048 bits RSA key), but will not exceed this maximum value. Only 1024 or 2048
	// values are allowed.
	//
	// If unset, the "Intermediate" DH Parameter size [1] is used:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	// +optional
	DHParamSize *DHParamSize `json:"dHParamSize,omitempty"`
}

// SecurityProtocol defines one or more security protocols used by
// an IngressController to secure network connections.
type SecurityProtocol struct {
	// minimumVersion enforces use of SecurityProtocolVersion or newer on
	// SSL connections instantiated from an IngressController. minimumVersion
	// must be lower than maximumVersion.
	//
	// If unset and maximumVersion is set, minimumVersion will be set
	// to maximumVersion. If minimumVersion and maximumVersion are unset,
	// the minimum version in "Intermediate" Versions [1] is used:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	// +optional
	MinimumVersion *SecurityProtocolVersion `json:"minimumVersion,omitempty"`
	// maximumVersion enforces use of SecurityProtocolVersion or older on
	// SSL connections instantiated from an IngressController. maximumVersion
	// must be higher than minimumVersion.
	//
	// If unset and minimumVersion is set, maximumVersion will be set
	// to minimumVersion. If minimumVersion and maximumVersion are unset,
	// the maximum version in "Intermediate" Versions [1] is used:
	//
	// [1] https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
	//
	// +optional
	MaximumVersion *SecurityProtocolVersion `json:"maximumVersion,omitempty"`
}

// SecurityProtocolVersion is a way to specify an IngressController security protocol.
type SecurityProtocolVersion string

const (
	// SecurityProtocolTLS10Version is v1.0 of the TLS security protocol.
	SecurityProtocolTLS10Version SecurityProtocolVersion = "TLSv1.0"
	// SecurityProtocolTLS11Version is v1.1 of the TLS security protocol.
	SecurityProtocolTLS11Version SecurityProtocolVersion = "TLSv1.1"
	// SecurityProtocolTLS12Version is v1.2 of the TLS security protocol.
	SecurityProtocolTLS12Version SecurityProtocolVersion = "TLSv1.2"
	// SecurityProtocolTLS13Version is v1.3 of the TLS security protocol.
	SecurityProtocolTLS13Version SecurityProtocolVersion = "TLSv1.3"
)

// DHParamSize sets the maximum size of the Diffie-Hellman parameters used for generating
// the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange.
type DHParamSize string

const (
	// DHParamSize1024 is a Diffie-Hellman parameter of 1024 bits.
	DHParamSize1024 DHParamSize = "1024"
	// DHParamSize2048 is a Diffie-Hellman parameter of 2048 bits.
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

	// securityProfileType is the actual security profile in use.
	SecurityProfile *SecurityProfileType `json:"securityProfile,omitempty"`

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
