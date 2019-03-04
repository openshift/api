package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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
// Additional features are available by default and through explicit
// configuration.
//
// 1. To publish the ingress controller endpoints to the Internet, other
//    networks, or external load balancers, configure an endpoint publishing
//    strategy.
//
// 2. On certain cloud platforms, when publishing an ingress controller, managed
//    wildcard DNS is automatically enabled. OpenShift will manage a DNS record
//    pointing to the ingress controller. DNS records are managed only in DNS
//    zones defined in the DNS cluster configuration resource.
//
// 3. If an ingress controller does not specify a default certificate, a new
//    self-signed certificate valid for the specified domain is generated for
//    the ingress controller.
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
	// domain is a DNS name used to configure various features that help publish
	// the ingress controller and enable external integrations.
	//
	// * The value is published to individual Route statuses so that end-users
	//   know where to target external DNS records.
	//
	// * When wildcard DNS management is enabled, domain is the base domain used
	//   to construct the wildcard host name.
	//
	// * If a generated default certificate is used for the ingress controller,
	//   the certificate will be valid for domain.
	//
	// domain must be unique among all IngressControllers, and cannot be
	// updated.
	//
	// If empty, defaults to the cluster Ingress config domain.
	//
	// +optional
	Domain string `json:"domain,omitempty"`

	// replicas is the desired number of ingress controller replicas. If unset,
	// defaults to 2.
	//
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// endpointPublishingStrategy is used to publish the ingress controller
	// endpoints to other netowrks, enable load balancer integrations, etc.
	//
	// If empty, the default is based on the cluster platform:
	//
	//   AWS: LoadBalancerService
	//   All other platform types: Private
	//
	// endpointPublishingStrategy cannot be updated.
	//
	// +optional
	EndpointPublishingStrategy EndpointPublishingStrategy `json:"endpointPublishingStrategy,omitempty"`

	// defaultCertificate is a reference to a secret containing the default
	// certificate served by the ingress controller. The secret must contain the
	// following data:
	//
	//   tls.crt: the certificate file
	//   tls.key: the certificate secret file
	//
	// If unset, a wildcard certificate is automatically generated and used. The
	// certificate is valid for the domain (and subdomains) and the certificate's
	// CA will be automatically integrated with the cluster's trust store.
	//
	// Whatever certificate is used (whether the generated default or explicitly
	// provided), the certificate will be automatically integrated with the
	// built-in authentication service.
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
}

// EndpointPublishingStrategyType is a way to publish ingress controller endpoints.
type EndpointPublishingStrategyType string

const (
	// LoadBalancerService publishes the ingress controler using a Kubernetes
	// LoadBalancer Service.
	//
	// In this configuration, the ingress controller deployment uses container
	// networking. A LoadBalancer Service is created to publish the deployment. If
	// domain is set, a wildcard DNS record will point to the LoadBalancer
	// Service's external name.
	//
	// See: https://kubernetes.io/docs/concepts/services-networking/#loadbalancer
	LoadBalancerServiceStrategyType EndpointPublishingStrategyType = "LoadBalancerService"

	// HostNetwork publishes the ingress controller on node ports where the
	// ingress controller is deployed.
	//
	// In this configuration, the ingress controller deployment uses host
	// networking, bound to node ports 80 and 443. The user is responsible for
	// configuring an external load balancer to publish the ingress controller via
	// the node ports.
	HostNetworkStrategyType EndpointPublishingStrategyType = "HostNetwork"

	// Private does not publish the ingress controller.
	//
	// In this configuration, the ingress controller deployment uses container
	// networking, and is not explicitly published. The user must manually publish
	// the ingress controller.
	PrivateStrategyType = "Private"
)

// EndpointPublishingStrategy is the a way to publish the endpoints of an
// IngressController, and represents the type and any additional configuration
// for a specific type.
type EndpointPublishingStrategy struct {
	// type is the publishing strategy to use. Valid values are
	// LoadBalancerService, HostNetwork, or Private.
	Type EndpointPublishingStrategyType `json:"type"`
}

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
	EndpointPublishingStrategy EndpointPublishingStrategy `json:"endpointPublishingStrategy"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressControllerList contains a list of IngressControllers.
type IngressControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IngressController `json:"items"`
}
