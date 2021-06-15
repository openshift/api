package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
//
// ExternalDNS describes a managed ExternalDNS controller instance for a cluster.
// The controller is responsible for creating external DNS records in supported
// DNS providers based off of instances of select Kubernetes resources.
type ExternalDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the ExternalDNS.
	Spec ExternalDNSSpec `json:"spec"`
	// status is the most recently observed status of the ExternalDNS.
	Status ExternalDNSStatus `json:"status"`
}

// ExternalDNSSpec defines the desired state of the ExternalDNS.
type ExternalDNSSpec struct {
	// Domains specifies which domains that ExternalDNS should
	// create DNS records for. Multiple domain values
	// can be specified such that subdomains of an included domain
	// can effectively be ignored using the "Include" and "Exclude"
	// domain filter options.
	//
	// An entry that excludes a domain that isn't a subdomain of some
	// included domain will be ignored.
	//
	// An empty list of domains means ExternalDNS will create
	// DNS records for any included source resource regardless
	// of the resource's desired hostname.
	//
	// +optional
	Domains []ExternalDNSDomain `json:"domains,omitempty"`

	// Provider refers to the DNS provider that ExternalDNS
	// should publish records to. Note that each ExternalDNS
	// is tied to a single provider.
	//
	// +kubebuilder:validation:Required
	// +required
	Provider ExternalDNSProvider `json:"provider"`

	// Source describes which source resource
	// ExternalDNS will be configured to create
	// DNS records for.
	//
	// Multiple ExternalDNS CRs must be
	// created if multiple ExternalDNS source resources
	// are desired.
	//
	// +kubebuilder:validation:Required
	// +required
	Source ExternalDNSSource `json:"source"`

	// Zones describes which DNS Zone IDs
	// ExternalDNS should publish records to.
	//
	// +optional
	Zones []string `json:"zones,omitempty"`
}

// ExternalDNSDomain describes how sets of included
// or excluded domains are to be constructed.
type ExternalDNSDomain struct {
	ExternalDNSDomainUnion `json:",inline"`

	// FilterType marks the Name or Pattern field
	// as an included or excluded set of domains.
	//
	// Note that excluded domains that
	// have not been included will be
	// ignored.
	//
	// This field accepts the following values:
	//
	//  "Include": Include the domain set specified
	//  by name or pattern.
	//
	//  "Exclude": Exclude the domain set specified
	//  by name or pattern.
	//
	// +kubebuilder:validation:Required
	// +required
	FilterType ExternalDNSFilterType `json:"filterType"`
}

// ExternalDNSDomainUnion describes optional fields of an External domain
// that should be captured.
// +union
type ExternalDNSDomainUnion struct {
	// MatchType specifies the type of match to be performed
	// by ExternalDNS when determining whether or not to publish DNS
	// records for a given source resource based on the resource's
	// requested hostname.
	//
	// This field accepts the following values:
	//
	//  "Exact": Explicitly match the full domain string
	//   specified via the name field, including any subdomains
	//   of the given name.
	//
	//  "Pattern": Match potential domains against
	//  the provided regular expression pattern string.
	//
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +required
	MatchType DomainMatchType `json:"matchType"`

	// Name is a string representing a single domain
	// value.
	//
	// e.g. my-app.my-cluster-domain.com
	//
	// +optional
	Name *string `json:"names,omitempty"`

	// Pattern is a regular expression used to
	// match a set of domains.
	//
	// +optional
	Pattern *string `json:"pattern,omitempty"`
}

// +kubebuilder:validation:Enum=Exact;Pattern
type DomainMatchType string

const (
	DomainMatchTypeExact DomainMatchType = "Exact"
	DomainMatchTypeRegex DomainMatchType = "Pattern"
)

// +kubebuilder:validation:Enum=Include;Exclude
type ExternalDNSFilterType string

const (
	FilterTypeInclude ExternalDNSFilterType = "Include"
	FilterTypeExclude ExternalDNSFilterType = "Exclude"
)

// ExternalDNSProvider specifies configuration
// options for the desired ExternalDNS DNS provider.
type ExternalDNSProvider struct {
	ExternalDNSProviderUnion `json:",inline"`

	// Credentials is a reference to a secret containing
	// the relevant credentials for the given provider, to be
	// used by an ExternalDNS controller for creating, updating,
	// and viewing DNS records.
	//
	// +kubebuilder:validation:Required
	// +required
	Credentials *corev1.LocalObjectReference `json:"credentials"`
}

// ExternalDNSProviderUnion describes optional fields for an ExternalDNS
// provider that should be captured.
// +union
type ExternalDNSProviderUnion struct {
	// Type describes which DNS provider
	// ExternalDNS should publish records to.
	// The following DNS providers are supported:
	//
	//  * AWS
	//  * GCP
	//  * Azure
	//  * BlueCat
	//  * Infoblox
	//
	// +kubebuilder:validation:Required
	// +unionDiscriminator
	// +required
	Type ExternalDNSProviderType `json:"type"`
}

// +kubebuilder:validation:Enum=AWS;GCP;Azure;BlueCat;Infoblox
type ExternalDNSProviderType string

const (
	ProviderTypeAWS      ExternalDNSProviderType = "AWS"
	ProviderTypeGCP      ExternalDNSProviderType = "GCP"
	ProviderTypeAzure    ExternalDNSProviderType = "Azure"
	ProviderTypeBlueCat  ExternalDNSProviderType = "BlueCat"
	ProviderTypeInfoblox ExternalDNSProviderType = "Infoblox"
	// More providers will ultimately be added in the future.
)

// ExternalDNSSource describes which Source resource
// the ExternalDNS should create DNS records for.
type ExternalDNSSource struct {
	ExternalDNSSourceUnion `json:",inline"`

	// HostnameAnnotationPolicy specifies whether or not ExternalDNS
	// should ignore the "external-dns.alpha.kubernetes.io/hostname"
	// annotation, which overrides DNS hostnames on a given source resource.
	//
	// The following values are accepted:
	//
	//  "Ignore": Ignore any hostname annotation overrides.
	//  "Allow": Allow all hostname annotation overrides.
	//
	// The default behavior of the ExternalDNS is "Ignore".
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=Ignore
	// +required
	HostnameAnnotationPolicy *HostnameAnnotationPolicy `json:"hostnameAnnotation"`
}

// ExternalDNSSourceUnion describes optional fields for an ExternalDNS source that should
// be captured.
// +union
type ExternalDNSSourceUnion struct {
	// Type specifies an ExternalDNS source resource
	// to create DNS records for.
	//
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +required
	Type ExternalDNSSourceType `json:"type"`

	// AnnotationFilter describes an annotation filter
	// used to filter which source instance resources
	// ExternalDNS publishes records for.
	// The annotation filter uses label selector semantics
	// against source resource annotations.
	//
	// +optional
	AnnotationFilter map[string]string `json:"annotationFilter,omitempty"`

	// Namespace instructs ExternalDNS to only acknowledge
	// source resource instances in a specific namespace.
	//
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Service describes source configuration options specific
	// to the service source resource.
	//
	// +optional
	Service *ExternalDNSServiceSourceOptions `json:"service,omitempty"`

	// CRD describes source configuration options specific
	// to the CRD source resource.
	//
	// +optional
	CRD *ExternalDNSCRDSourceOptions `json:"crd,omitempty"`

	// OpenShiftRoute source currently has no unique configuration options.
}

// +kubebuilder:validation:Enum=OpenShiftRoute;Service;CRD
type ExternalDNSSourceType string

const (
	SourceTypeRoute   ExternalDNSSourceType = "OpenShiftRoute"
	SourceTypeService ExternalDNSSourceType = "Service"
	SourceTypeCRD     ExternalDNSSourceType = "CRD"
)

// +kubebuilder:validation:Enum=Ignore;Allow
type HostnameAnnotationPolicy string

const (
	HostnameAnnotationPolicyIgnore HostnameAnnotationPolicy = "Ignore"
	HostnameAnnotationPolicyAllow  HostnameAnnotationPolicy = "Allow"
)

// ExternalDNSServiceSourceOptions describes options
// specific to the ExternalDNS service source.
type ExternalDNSServiceSourceOptions struct {
	// PublishPolicy determines what kinds of Service
	// resources are published to the given DNS provider.
	//
	// The following publishing policies are available:
	//
	//  "PublishExternal": Only publish DNS records for
	//  service that are externally reachable,
	//  such as NodePort, ExternalName, and LoadBalancer
	//  services. This is the default behavior of ExternalDNS.
	//
	//  "PublishInternalExternal": Publish DNS records
	//   for externally reachable services, in addition to
	//   internally reachable services, such as ClusterIP
	//   services.
	//
	// +optional
	PublishPolicy *ServicePublishPolicy `json:"publishPolicy,omitempty"`

	// ServiceType determines what types of Service resources
	// are watched by ExternalDNS. The following types are
	// available options:
	//
	//  "NodePort"
	//  "ExternalName"
	//  "LoadBalancer"
	//  "ClusterIP" (requires a publishPolicy of "PublishInternalExternal")
	//
	// One or more Service types can be specified, if desired.
	//
	// If no service types are provided, ExternalDNS will be
	// configured to create DNS records for LoadBalancer services
	// only by default.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default:={"LoadBalancer"}
	// +kubebuilder:validation:MinItems=1
	// +required
	ServiceType []corev1.ServiceType `json:"serviceType,omitempty"`
}

type ServicePublishPolicy string

// +kubebuilder:validation:Enum=PublishExternal;PublishInternalExternal
const (
	ServicePublishPolicyExternalOnly        ServicePublishPolicy = "PublishExternal"
	ServicePublishPolicyExternalAndInternal ServicePublishPolicy = "PublishInternalExternal"
)

type ExternalDNSCRDSourceOptions struct {
	// Kind is the kind of the CRD
	// source resource type to be
	// consumed by ExternalDNS.
	//
	// e.g. "DNSEndpoint"
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +required
	Kind string `json:"kind"`

	// Version is the API version
	// of the given resource kind for
	// ExternalDNS to use.
	//
	// e.g. "externaldns.k8s.io/v1alpha1"
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +required
	Version string `json:"version"`

	// LabelFilter specifies a label filter
	// to be used to filter CRD resource instances.
	// Only one label filter can be specified on
	// an ExternalDNS instance.
	//
	// +optional
	LabelFilter map[string]string `json:"labelFilter,omitempty"`
}

// ExternalDNSStatus defines the observed state of ExternalDNS
type ExternalDNSStatus struct {
	// Conditions is a list of operator-specific conditions
	// and their status.
	Conditions []ExternalDNSOperatorCondition `json:"conditions,omitempty"`

	// ObservedGeneration is the most recent generation observed.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Provider is the configured provider in use by ExternalDNS.
	Provider ExternalDNSProvider `json:"provider,omitempty"`

	// Zones is the configured zones in use by ExternalDNS.
	Zones []string `json:"zones,omitempty"`
}

type ExternalDNSOperatorCondition struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +required
	Type string `json:"type"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +required
	Status string `json:"status"`

	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// +optional
	Reason *string `json:"reason,omitempty"`

	// +optional
	Message *string `json:"message,omitempty"`
}

var (
	// Available indicates that the ExternalDNS is available.
	ExternalDNSAvailableConditionType = "Available"

	// AuthenticationFailed indicates that there were issues starting
	// ExternalDNS pods related to the given provider credentials.
	ExternalDNSProviderAuthFailedReasonType = "AuthenticationFailed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//
// ExternalDNSList contains a list of ExternalDNS
type ExternalDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ExternalDNS `json:"items"`
}
