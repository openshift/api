package v1

import (
	v1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=dnses,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/475
// +openshift:file-pattern=cvoRunLevel=0000_70,operatorName=dns,operatorOrdering=00

// DNS manages the CoreDNS component to provide a name resolution service
// for pods and services in the cluster.
//
// This supports the DNS-based service discovery specification:
// https://github.com/kubernetes/dns/blob/master/docs/specification.md
//
// More details: https://kubernetes.io/docs/tasks/administer-cluster/coredns
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type DNS struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the DNS.
	Spec DNSSpec `json:"spec,omitempty"`
	// status is the most recently observed status of the DNS.
	Status DNSStatus `json:"status,omitempty"`
}

// DNSSpec is the specification of the desired behavior of the DNS.
type DNSSpec struct {
	// servers is a list of DNS resolvers that provide name query delegation for one or
	// more subdomains outside the scope of the cluster domain. If servers consists of
	// more than one Server, longest suffix match will be used to determine the Server.
	//
	// For example, if there are two Servers, one for "foo.com" and another for "a.foo.com",
	// and the name query is for "www.a.foo.com", it will be routed to the Server with Zone
	// "a.foo.com".
	//
	// If this field is nil, no servers are created.
	//
	// +optional
	Servers []Server `json:"servers,omitempty"`

	// upstreamResolvers defines a schema for configuring CoreDNS
	// to proxy DNS messages to upstream resolvers for the case of the
	// default (".") server
	//
	// If this field is not specified, the upstream used will default to
	// /etc/resolv.conf, with policy "sequential"
	//
	// +optional
	UpstreamResolvers UpstreamResolvers `json:"upstreamResolvers"`

	// nodePlacement provides explicit control over the scheduling of DNS
	// pods.
	//
	// Generally, it is useful to run a DNS pod on every node so that DNS
	// queries are always handled by a local DNS pod instead of going over
	// the network to a DNS pod on another node.  However, security policies
	// may require restricting the placement of DNS pods to specific nodes.
	// For example, if a security policy prohibits pods on arbitrary nodes
	// from communicating with the API, a node selector can be specified to
	// restrict DNS pods to nodes that are permitted to communicate with the
	// API.  Conversely, if running DNS pods on nodes with a particular
	// taint is desired, a toleration can be specified for that taint.
	//
	// If unset, defaults are used. See nodePlacement for more details.
	//
	// +optional
	NodePlacement DNSNodePlacement `json:"nodePlacement,omitempty"`

	// managementState indicates whether the DNS operator should manage cluster
	// DNS
	// +optional
	ManagementState ManagementState `json:"managementState,omitempty"`

	// operatorLogLevel controls the logging level of the DNS Operator.
	// Valid values are: "Normal", "Debug", "Trace".
	// Defaults to "Normal".
	// setting operatorLogLevel: Trace will produce extremely verbose logs.
	// +optional
	// +kubebuilder:default=Normal
	OperatorLogLevel DNSLogLevel `json:"operatorLogLevel,omitempty"`

	// logLevel describes the desired logging verbosity for CoreDNS.
	// Any one of the following values may be specified:
	// * Normal logs errors from upstream resolvers.
	// * Debug logs errors, NXDOMAIN responses, and NODATA responses.
	// * Trace logs errors and all responses.
	//  Setting logLevel: Trace will produce extremely verbose logs.
	// Valid values are: "Normal", "Debug", "Trace".
	// Defaults to "Normal".
	// +optional
	// +kubebuilder:default=Normal
	LogLevel DNSLogLevel `json:"logLevel,omitempty"`

	// cache describes the caching configuration that applies to all server blocks listed in the Corefile.
	// This field allows a cluster admin to optionally configure:
	// * positiveTTL which is a duration for which positive responses should be cached.
	// * negativeTTL which is a duration for which negative responses should be cached.
	// If this is not configured, OpenShift will configure positive and negative caching with a default value that is
	// subject to change. At the time of writing, the default positiveTTL is 900 seconds and the default negativeTTL is
	// 30 seconds or as noted in the respective Corefile for your version of OpenShift.
	// +optional
	Cache DNSCache `json:"cache,omitempty"`

	// template is an optional configuration for custom DNS query handling via the CoreDNS template plugin.
	// The template defines how to handle queries matching specific zones and query types.
	//
	// The template applies to all domains (custom domains from spec.servers and the cluster domain)
	// to ensure consistent DNS resolution across all paths.
	//
	// When this field is not set, no template plugin configuration is added to CoreDNS.
	//
	// +optional
	// +openshift:enable:FeatureGate=DNSTemplatePlugin
	Template DNSTemplate `json:"template,omitzero"`

	// responseRules is an optional list of rules that configure CoreDNS to
	// return custom A, CNAME, or NXDOMAIN responses for specific DNS names.
	//
	// Rules are evaluated independently; an exact-name rule always takes
	// precedence over a wildcard rule for the same query, and among
	// wildcard rules the most specific match wins. Queries that do not
	// match any rule continue through normal CoreDNS resolution.
	//
	// responseRules apply only to the default server block; they do not
	// affect spec.servers or spec.template.
	//
	// When this field is not set, no response rule configuration is added
	// to CoreDNS.
	// At most 100 rules may be configured. Each query name may appear only
	// once. An exact A rule cannot overlap a wildcard CNAME rule because
	// the wildcard would also answer AAAA and CNAME queries for that name.
	//
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=100
	// +kubebuilder:validation:XValidation:rule="self.all(r, self.exists_one(other, r.query.name == other.query.name))",message="responseRules must have distinct query names"
	// +kubebuilder:validation:XValidation:rule="self.all(a, a.action.type != 'A' || a.query.name.startsWith('*.') || self.all(c, c.action.type != 'CNAME' || !c.query.name.startsWith('*.') || !(a.query.name.endsWith(c.query.name.substring(1)) && !a.query.name.substring(0, size(a.query.name) - size(c.query.name) + 1).contains('.'))))",message="an exact A rule must not overlap a wildcard CNAME rule"
	// +openshift:enable:FeatureGate=DNSTemplatePlugin
	ResponseRules []DNSResponseRule `json:"responseRules,omitempty"`
}

// DNSCache defines the fields for configuring DNS caching.
type DNSCache struct {
	// positiveTTL is optional and specifies the amount of time that a positive response should be cached.
	//
	// If configured, it must be a value of 1s (1 second) or greater up to a theoretical maximum of several years. This
	// field expects an unsigned duration string of decimal numbers, each with optional fraction and a unit suffix,
	// e.g. "100s", "1m30s", "12h30m10s". Values that are fractions of a second are rounded down to the nearest second.
	// If the configured value is less than 1s, the default value will be used.
	// If not configured, the value will be 0s and OpenShift will use a default value of 900 seconds unless noted
	// otherwise in the respective Corefile for your version of OpenShift. The default value of 900 seconds is subject
	// to change.
	// +kubebuilder:validation:Pattern=^(0|([0-9]+(\.[0-9]+)?(ns|us|µs|μs|ms|s|m|h))+)$
	// +kubebuilder:validation:Type:=string
	// +optional
	PositiveTTL metav1.Duration `json:"positiveTTL,omitempty"`

	// negativeTTL is optional and specifies the amount of time that a negative response should be cached.
	//
	// If configured, it must be a value of 1s (1 second) or greater up to a theoretical maximum of several years. This
	// field expects an unsigned duration string of decimal numbers, each with optional fraction and a unit suffix,
	// e.g. "100s", "1m30s", "12h30m10s". Values that are fractions of a second are rounded down to the nearest second.
	// If the configured value is less than 1s, the default value will be used.
	// If not configured, the value will be 0s and OpenShift will use a default value of 30 seconds unless noted
	// otherwise in the respective Corefile for your version of OpenShift. The default value of 30 seconds is subject
	// to change.
	// +kubebuilder:validation:Pattern=^(0|([0-9]+(\.[0-9]+)?(ns|us|µs|μs|ms|s|m|h))+)$
	// +kubebuilder:validation:Type:=string
	// +optional
	NegativeTTL metav1.Duration `json:"negativeTTL,omitempty"`
}

// +kubebuilder:validation:Enum:=Normal;Debug;Trace
type DNSLogLevel string

var (
	// Normal is the default.  Normal, working log information, everything is fine, but helpful notices for auditing or common operations.  In kube, this is probably glog=2.
	DNSLogLevelNormal DNSLogLevel = "Normal"

	// Debug is used when something went wrong.  Even common operations may be logged, and less helpful but more quantity of notices.  In kube, this is probably glog=4.
	DNSLogLevelDebug DNSLogLevel = "Debug"

	// Trace is used when something went really badly and even more verbose logs are needed.  Logging every function call as part of a common operation, to tracing execution of a query.  In kube, this is probably glog=6.
	DNSLogLevelTrace DNSLogLevel = "Trace"
)

// Server defines the schema for a server that runs per instance of CoreDNS.
type Server struct {
	// name is required and specifies a unique name for the server. Name must comply
	// with the Service Name Syntax of rfc6335.
	Name string `json:"name"`
	// zones is required and specifies the subdomains that Server is authoritative for.
	// Zones must conform to the rfc1123 definition of a subdomain. Specifying the
	// cluster domain (i.e., "cluster.local") is invalid.
	Zones []string `json:"zones"`
	// forwardPlugin defines a schema for configuring CoreDNS to proxy DNS messages
	// to upstream resolvers.
	ForwardPlugin ForwardPlugin `json:"forwardPlugin"`
}

// DNSTransport indicates what type of connection should be used.
// +kubebuilder:validation:Enum=TLS;Cleartext;""
type DNSTransport string

const (
	// TLSTransport indicates that TLS should be used for the connection.
	TLSTransport DNSTransport = "TLS"

	// CleartextTransport indicates that no encryption should be used for
	// the connection.
	CleartextTransport DNSTransport = "Cleartext"
)

// DNSTransportConfig groups related configuration parameters used for configuring
// forwarding to upstream resolvers that support DNS-over-TLS.
// +union
type DNSTransportConfig struct {
	// transport allows cluster administrators to opt-in to using a DNS-over-TLS
	// connection between cluster DNS and an upstream resolver(s). Configuring
	// TLS as the transport at this level without configuring a CABundle will
	// result in the system certificates being used to verify the serving
	// certificate of the upstream resolver(s).
	//
	// Possible values:
	// "" (empty) - This means no explicit choice has been made and the platform chooses the default which is subject
	// to change over time. The current default is "Cleartext".
	// "Cleartext" - Cluster admin specified cleartext option. This results in the same functionality
	// as an empty value but may be useful when a cluster admin wants to be more explicit about the transport,
	// or wants to switch from "TLS" to "Cleartext" explicitly.
	// "TLS" - This indicates that DNS queries should be sent over a TLS connection. If Transport is set to TLS,
	// you MUST also set ServerName. If a port is not included with the upstream IP, port 853 will be tried by default
	// per RFC 7858 section 3.1; https://datatracker.ietf.org/doc/html/rfc7858#section-3.1.
	//
	// +optional
	// +unionDiscriminator
	Transport DNSTransport `json:"transport,omitempty"`

	// tls contains the additional configuration options to use when Transport is set to "TLS".
	TLS *DNSOverTLSConfig `json:"tls,omitempty"`
}

// DNSOverTLSConfig describes optional DNSTransportConfig fields that should be captured.
type DNSOverTLSConfig struct {
	// serverName is the upstream server to connect to when forwarding DNS queries. This is required when Transport is
	// set to "TLS". ServerName will be validated against the DNS naming conventions in RFC 1123 and should match the
	// TLS certificate installed in the upstream resolver(s).
	//
	// + ---
	// + Inspired by the DNS1123 patterns in Kubernetes: https://github.com/kubernetes/kubernetes/blob/7c46f40bdf89a437ecdbc01df45e235b5f6d9745/staging/src/k8s.io/apimachinery/pkg/util/validation/validation.go#L178-L218
	// +required
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
	ServerName string `json:"serverName"`

	// caBundle references a ConfigMap that must contain either a single
	// CA Certificate or a CA Bundle. This allows cluster administrators to provide their
	// own CA or CA bundle for validating the certificate of upstream resolvers.
	//
	// 1. The configmap must contain a `ca-bundle.crt` key.
	// 2. The value must be a PEM encoded CA certificate or CA bundle.
	// 3. The administrator must create this configmap in the openshift-config namespace.
	// 4. The upstream server certificate must contain a Subject Alternative Name (SAN) that matches ServerName.
	//
	// +optional
	CABundle v1.ConfigMapNameReference `json:"caBundle,omitempty"`
}

// ForwardingPolicy is the policy to use when forwarding DNS requests.
// +kubebuilder:validation:Enum=Random;RoundRobin;Sequential
type ForwardingPolicy string

const (
	// RandomForwardingPolicy picks a random upstream server for each query.
	RandomForwardingPolicy ForwardingPolicy = "Random"

	// RoundRobinForwardingPolicy picks upstream servers in a round-robin order, moving to the next server for each new query.
	RoundRobinForwardingPolicy ForwardingPolicy = "RoundRobin"

	// SequentialForwardingPolicy tries querying upstream servers in a sequential order until one responds, starting with the first server for each new query.
	SequentialForwardingPolicy ForwardingPolicy = "Sequential"
)

// ForwardPlugin defines a schema for configuring the CoreDNS forward plugin.
type ForwardPlugin struct {
	// upstreams is a list of resolvers to forward name queries for subdomains of Zones.
	// Each instance of CoreDNS performs health checking of Upstreams. When a healthy upstream
	// returns an error during the exchange, another resolver is tried from Upstreams. The
	// Upstreams are selected in the order specified in Policy. Each upstream is represented
	// by an IP address or IP:port if the upstream listens on a port other than 53.
	//
	// A maximum of 15 upstreams is allowed per ForwardPlugin.
	//
	// +kubebuilder:validation:MaxItems=15
	Upstreams []string `json:"upstreams"`

	// policy is used to determine the order in which upstream servers are selected for querying.
	// Any one of the following values may be specified:
	//
	// * "Random" picks a random upstream server for each query.
	// * "RoundRobin" picks upstream servers in a round-robin order, moving to the next server for each new query.
	// * "Sequential" tries querying upstream servers in a sequential order until one responds, starting with the first server for each new query.
	//
	// The default value is "Random"
	//
	// +optional
	// +kubebuilder:default:="Random"
	Policy ForwardingPolicy `json:"policy,omitempty"`

	// transportConfig is used to configure the transport type, server name, and optional custom CA or CA bundle to use
	// when forwarding DNS requests to an upstream resolver.
	//
	// The default value is "" (empty) which results in a standard cleartext connection being used when forwarding DNS
	// requests to an upstream resolver.
	//
	// +optional
	TransportConfig DNSTransportConfig `json:"transportConfig,omitempty"`

	// protocolStrategy specifies the protocol to use for upstream DNS
	// requests.
	// Valid values for protocolStrategy are "TCP" and omitted.
	// When omitted, this means no opinion and the platform is left to choose
	// a reasonable default, which is subject to change over time.
	// The current default is to use the protocol of the original client request.
	// "TCP" specifies that the platform should use TCP for all upstream DNS requests,
	// even if the client request uses UDP.
	// "TCP" is useful for UDP-specific issues such as those created by
	// non-compliant upstream resolvers, but may consume more bandwidth or
	// increase DNS response time. Note that protocolStrategy only affects
	// the protocol of DNS requests that CoreDNS makes to upstream resolvers.
	// It does not affect the protocol of DNS requests between clients and
	// CoreDNS.
	//
	// +optional
	ProtocolStrategy ProtocolStrategy `json:"protocolStrategy"`
}

// UpstreamResolvers defines a schema for configuring the CoreDNS forward plugin in the
// specific case of the default (".") server.
// It defers from ForwardPlugin in the default values it accepts:
// * At least one upstream should be specified.
// * the default policy is Sequential
type UpstreamResolvers struct {
	// upstreams is a list of resolvers to forward name queries for the "." domain.
	// Each instance of CoreDNS performs health checking of Upstreams. When a healthy upstream
	// returns an error during the exchange, another resolver is tried from Upstreams. The
	// Upstreams are selected in the order specified in Policy.
	//
	// A maximum of 15 upstreams is allowed per ForwardPlugin.
	// If no Upstreams are specified, /etc/resolv.conf is used by default
	//
	// +optional
	// +kubebuilder:validation:MaxItems=15
	// +kubebuilder:default={{"type":"SystemResolvConf"}}
	Upstreams []Upstream `json:"upstreams"`

	// policy is used to determine the order in which upstream servers are selected for querying.
	// Any one of the following values may be specified:
	//
	// * "Random" picks a random upstream server for each query.
	// * "RoundRobin" picks upstream servers in a round-robin order, moving to the next server for each new query.
	// * "Sequential" tries querying upstream servers in a sequential order until one responds, starting with the first server for each new query.
	//
	// The default value is "Sequential"
	//
	// +optional
	// +kubebuilder:default="Sequential"
	Policy ForwardingPolicy `json:"policy,omitempty"`

	// transportConfig is used to configure the transport type, server name, and optional custom CA or CA bundle to use
	// when forwarding DNS requests to an upstream resolver.
	//
	// The default value is "" (empty) which results in a standard cleartext connection being used when forwarding DNS
	// requests to an upstream resolver.
	//
	// +optional
	TransportConfig DNSTransportConfig `json:"transportConfig,omitempty"`

	// protocolStrategy specifies the protocol to use for upstream DNS
	// requests.
	// Valid values for protocolStrategy are "TCP" and omitted.
	// When omitted, this means no opinion and the platform is left to choose
	// a reasonable default, which is subject to change over time.
	// The current default is to use the protocol of the original client request.
	// "TCP" specifies that the platform should use TCP for all upstream DNS requests,
	// even if the client request uses UDP.
	// "TCP" is useful for UDP-specific issues such as those created by
	// non-compliant upstream resolvers, but may consume more bandwidth or
	// increase DNS response time. Note that protocolStrategy only affects
	// the protocol of DNS requests that CoreDNS makes to upstream resolvers.
	// It does not affect the protocol of DNS requests between clients and
	// CoreDNS.
	//
	// +optional
	ProtocolStrategy ProtocolStrategy `json:"protocolStrategy"`
}

// Upstream can either be of type SystemResolvConf, or of type Network.
//
//   - For an Upstream of type SystemResolvConf, no further fields are necessary:
//     The upstream will be configured to use /etc/resolv.conf.
//   - For an Upstream of type Network, a NetworkResolver field needs to be defined
//     with an IP address or IP:port if the upstream listens on a port other than 53.
type Upstream struct {

	// type defines whether this upstream contains an IP/IP:port resolver or the local /etc/resolv.conf.
	// Type accepts 2 possible values: SystemResolvConf or Network.
	//
	// * When SystemResolvConf is used, the Upstream structure does not require any further fields to be defined:
	//   /etc/resolv.conf will be used
	// * When Network is used, the Upstream structure must contain at least an Address
	//
	// +required
	Type UpstreamType `json:"type"`

	// address must be defined when Type is set to Network. It will be ignored otherwise.
	// It must be a valid ipv4 or ipv6 address.
	//
	// +optional
	Address string `json:"address,omitempty"`

	// port may be defined when Type is set to Network. It will be ignored otherwise.
	// Port must be between 65535
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=53
	Port uint32 `json:"port,omitempty"`
}

// +kubebuilder:validation:Enum=SystemResolvConf;Network;""
type UpstreamType string

const (
	SystemResolveConfType UpstreamType = "SystemResolvConf"
	NetworkResolverType   UpstreamType = "Network"
)

// ProtocolStrategy is a preference for the protocol to use for DNS queries.
// + ---
// + When consumers observe an unknown value, they should use the default strategy.
// +kubebuilder:validation:Enum:=TCP;""
type ProtocolStrategy string

var (
	// ProtocolStrategyDefault specifies no opinion for DNS protocol.
	// If empty, the default behavior of CoreDNS is used. Currently, this means that CoreDNS uses the protocol of the
	// originating client request as the upstream protocol.
	// Note that the default behavior of CoreDNS is subject to change.
	ProtocolStrategyDefault ProtocolStrategy = ""

	// ProtocolStrategyTCP instructs CoreDNS to always use TCP, regardless of the originating client's request protocol.
	ProtocolStrategyTCP ProtocolStrategy = "TCP"
)

// DNSNodePlacement describes the node scheduling configuration for DNS pods.
type DNSNodePlacement struct {
	// nodeSelector is the node selector applied to DNS pods.
	//
	// If empty, the default is used, which is currently the following:
	//
	//   kubernetes.io/os: linux
	//
	// This default is subject to change.
	//
	// If set, the specified selector is used and replaces the default.
	//
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// tolerations is a list of tolerations applied to DNS pods.
	//
	// If empty, the DNS operator sets a toleration for the
	// "node-role.kubernetes.io/master" taint.  This default is subject to
	// change.  Specifying tolerations without including a toleration for
	// the "node-role.kubernetes.io/master" taint may be risky as it could
	// lead to an outage if all worker nodes become unavailable.
	//
	// Note that the daemon controller adds some tolerations as well.  See
	// https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

const (
	// Available indicates the DNS controller daemonset is available.
	DNSAvailable = "Available"
)

// QueryType represents DNS query types supported by templates.
// +kubebuilder:validation:Enum=AAAA
type QueryType string

const (
	// QueryTypeAAAA represents IPv6 address records (AAAA).
	QueryTypeAAAA QueryType = "AAAA"
)

// QueryClass represents DNS query classes supported by templates.
// Valid value is "IN".
// +kubebuilder:validation:Enum=IN
type QueryClass string

const (
	// QueryClassIN represents the Internet class.
	QueryClassIN QueryClass = "IN"
)

// ResponseCode represents DNS response codes.
// +kubebuilder:validation:Enum=NOERROR
type ResponseCode string

const (
	// ResponseCodeNOERROR indicates a successful DNS query with or without answer records.
	ResponseCodeNOERROR ResponseCode = "NOERROR"
)

// Zone is a DNS zone name. It must be either "." (catch-all) or a valid RFC1123 subdomain.
// Valid RFC1123 subdomains consist of lowercase alphanumeric characters, hyphens, and dots.
// Labels cannot start or end with hyphens and must be 1-63 characters each.
// The entire value must be 1-253 characters long.
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=253
// +kubebuilder:validation:XValidation:rule="self == '.' || !format.dns1123Subdomain().validate(self).hasValue()",message="zone must be '.' or a valid lowercase DNS subdomain"
// +kubebuilder:validation:XValidation:rule="self.split('.').all(label, size(label) <= 63)",message="each DNS label must be at most 63 characters"
type Zone string

// DNSTemplate defines a template for custom DNS query handling via the CoreDNS template plugin.
// DNSTemplate enables filtering or custom responses for DNS queries matching specific zones and query types.
type DNSTemplate struct {
	// zones specifies the DNS zones this template applies to.
	// Each zone must be a valid DNS name as defined in RFC 1123.
	// The special zone "." matches all domains (catch-all).
	// Multiple zones can be specified to apply the same template actions to multiple domains.
	// At least 1 and at most 15 zones may be specified.
	//
	// Note: root zone (".") includes cluster domain (cluster.local); use specific zones to avoid impacting IPv6 queries in IPv6 or dual-stack clusters.
	//
	// Examples:
	// - ["."] matches all domains (catch-all for global AAAA filtering)
	// - ["example.com"] matches only example.com and its subdomains
	// - ["example.com", "test.com"] matches both domains and their subdomains
	//
	// +listType=set
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=15
	// +required
	Zones []Zone `json:"zones,omitempty"`

	// queryType specifies the DNS query type to match.
	// Valid values are "AAAA".
	//
	// +required
	QueryType QueryType `json:"queryType,omitempty"`

	// queryClass specifies the DNS query class to match.
	// Valid values are "IN".
	//
	// +required
	QueryClass QueryClass `json:"queryClass,omitempty"`

	// action defines how to handle queries matching this template's zones and query type.
	// The action builds a single DNS response by specifying the response code and may be
	// extended by additional fields in the future.
	//
	// +required
	Action DNSTemplateAction `json:"action,omitzero"`
}

// DNSTemplateAction defines how to construct a DNS response for queries matching the template.
type DNSTemplateAction struct {
	// rcode is the DNS response code to return.
	// Valid values are "NOERROR".
	//
	// The template returns a response with no answer records. For AAAA filtering,
	// this means IPv6 address queries return successfully but with no IPv6 addresses,
	// causing clients to fall back to IPv4 (A record) queries.
	//
	// +required
	Rcode ResponseCode `json:"rcode,omitempty"`
}

// DNSResponseRuleName identifies a DNSResponseRule within spec.responseRules.
// It must be 1-63 lowercase alphanumeric characters or hyphens, beginning
// and ending with an alphanumeric character.
// +kubebuilder:validation:MinLength=1
// +kubebuilder:validation:MaxLength=63
// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]([-a-z0-9]*[a-z0-9])?$')",message="rule name must be a lowercase DNS label"
type DNSResponseRuleName string

// DNSResponseRule configures CoreDNS to return a custom A, CNAME, or
// NXDOMAIN response for queries matching query.name.
// An exact-name CNAME cannot target itself because CoreDNS would repeatedly
// resolve the same name.
//
// +kubebuilder:validation:XValidation:rule="self.action.type == 'NXDOMAIN' ? !self.query.name.startsWith('*.') : true",message="nxdomain is only valid when query.name is an exact name, not a wildcard"
// +kubebuilder:validation:XValidation:rule="!has(self.action.cname) || self.query.name.startsWith('*.') || self.action.cname.target != self.query.name",message="an exact-name CNAME rule must not target its own query name"
type DNSResponseRule struct {
	// name uniquely identifies this rule within spec.responseRules.
	//
	// +required
	Name DNSResponseRuleName `json:"name,omitempty"`

	// query selects the DNS name this rule answers: an exact name, or a
	// name with exactly one leading wildcard label (e.g. "*.example.com").
	//
	// +required
	Query DNSResponseQuery `json:"query,omitzero"`

	// ttlSeconds is the TTL, in seconds, included in synthesized answers.
	//
	// When omitted, this means no opinion and the platform chooses a
	// reasonable default, which is subject to change over time. The
	// current default is 60. Valid values are 30-3600.
	//
	// +optional
	// +kubebuilder:validation:Minimum=30
	// +kubebuilder:validation:Maximum=3600
	// +default=60
	TTLSeconds int32 `json:"ttlSeconds,omitempty"`

	// action specifies how matching queries are answered. Exactly one of
	// a, cname, or nxdomain is set, as indicated by action.type.
	//
	// +required
	Action DNSResponseAction `json:"action,omitzero"`
}

// DNSResponseQuery selects the DNS name a DNSResponseRule answers.
type DNSResponseQuery struct {
	// name is the DNS name this rule matches: either an exact name, or a
	// name with exactly one leading wildcard label replacing the first
	// name segment (e.g. "*.example.com"). It must be a lowercase DNS name
	// with at least two labels, at most 253 characters overall, and at most
	// 63 characters per label. It must not include a trailing dot.
	// Rules cannot cover cluster.local, in-addr.arpa, or ip6.arpa, including
	// a wildcard that matches one of those zones.
	//
	// +required
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.startsWith('*.') ? size(self.split('.')) >= 3 && !format.dns1123Subdomain().validate(self.substring(2)).hasValue() : size(self.split('.')) >= 2 && !format.dns1123Subdomain().validate(self).hasValue()",message="query name must be a lowercase DNS name with an optional single leading wildcard label"
	// +kubebuilder:validation:XValidation:rule="self.split('.').all(label, size(label) <= 63)",message="each DNS label must be at most 63 characters"
	// +kubebuilder:validation:XValidation:rule="!(self == 'cluster.local' || self.endsWith('.cluster.local') || self == '*.local' || self == 'in-addr.arpa' || self.endsWith('.in-addr.arpa') || self == 'ip6.arpa' || self.endsWith('.ip6.arpa') || self == '*.arpa')",message="response rules must not cover cluster.local or reverse DNS zones"
	Name string `json:"name,omitempty"`
}

// DNSResponseActionType specifies which kind of synthesized response a
// DNSResponseRule produces.
// +kubebuilder:validation:Enum=A;CNAME;NXDOMAIN
type DNSResponseActionType string

const (
	// DNSResponseActionTypeA indicates that the rule synthesizes an A response.
	DNSResponseActionTypeA DNSResponseActionType = "A"

	// DNSResponseActionTypeCNAME indicates that the rule synthesizes a CNAME response.
	DNSResponseActionTypeCNAME DNSResponseActionType = "CNAME"

	// DNSResponseActionTypeNXDOMAIN indicates that the rule synthesizes an NXDOMAIN response.
	DNSResponseActionTypeNXDOMAIN DNSResponseActionType = "NXDOMAIN"
)

// DNSResponseAction specifies how CoreDNS answers a query matching a
// DNSResponseRule's query. Exactly one of a, cname, or nxdomain is set,
// as indicated by type.
//
// +union
// +kubebuilder:validation:XValidation:rule="self.type == 'A' ? has(self.a) : !has(self.a)",message="a must be set if type is 'A', but remain unset otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'CNAME' ? has(self.cname) : !has(self.cname)",message="cname must be set if type is 'CNAME', but remain unset otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'NXDOMAIN' ? has(self.nxdomain) : !has(self.nxdomain)",message="nxdomain must be set if type is 'NXDOMAIN', but remain unset otherwise"
type DNSResponseAction struct {
	// type specifies which kind of response this rule synthesizes.
	// Valid values are "A", "CNAME", and "NXDOMAIN".
	//
	// +unionDiscriminator
	// +required
	Type DNSResponseActionType `json:"type,omitempty"`

	// a configures a synthesized A response. Required when type is A,
	// forbidden otherwise.
	//
	// +unionMember
	// +optional
	A DNSResponseActionA `json:"a,omitzero"`

	// cname configures a synthesized CNAME response. Required when type
	// is CNAME, forbidden otherwise.
	//
	// +unionMember
	// +optional
	CNAME DNSResponseActionCNAME `json:"cname,omitzero"`

	// nxdomain configures a synthesized NXDOMAIN response. Required when
	// type is NXDOMAIN, forbidden otherwise. nxdomain is only valid when
	// the rule's query.name is an exact name, not a wildcard.
	//
	// +unionMember
	// +optional
	//nolint:kubeapilinter // NXDOMAIN deliberately has an empty object payload.
	NXDOMAIN *DNSResponseActionNXDOMAIN `json:"nxdomain,omitzero"`
}

// IPv4Address is an IPv4 address in dotted-decimal notation, with four
// decimal octets from 0 to 255 and no leading zeroes in multi-digit octets.
// It must be 7-15 characters long.
// +kubebuilder:validation:MinLength=7
// +kubebuilder:validation:MaxLength=15
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).family() == 4 && ip.isCanonical(self)",message="address must be an IPv4 address in dotted-decimal notation"
type IPv4Address string

// DNSResponseActionA configures a synthesized A response.
type DNSResponseActionA struct {
	// addresses is the list of IPv4 addresses returned in the synthesized
	// answer. At least 1 and at most 10 addresses may be specified.
	//
	// +required
	// +listType=set
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	Addresses []IPv4Address `json:"addresses,omitempty"`
}

// DNSResponseActionCNAME configures a synthesized CNAME response.
type DNSResponseActionCNAME struct {
	// target is the lowercase DNS name the CNAME response points to. It must
	// have at least two labels, at most 253 characters overall, and at most
	// 63 characters per label. It must not include a trailing dot.
	//
	// +required
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="size(self.split('.')) >= 2 && !format.dns1123Subdomain().validate(self).hasValue()",message="CNAME target must be a lowercase DNS name"
	// +kubebuilder:validation:XValidation:rule="self.split('.').all(label, size(label) <= 63)",message="each DNS label must be at most 63 characters"
	Target string `json:"target,omitempty"`
}

// DNSResponseActionNXDOMAIN configures a synthesized NXDOMAIN response.
type DNSResponseActionNXDOMAIN struct {
}

// DNSStatus defines the observed status of the DNS.
type DNSStatus struct {
	// clusterIP is the service IP through which this DNS is made available.
	//
	// In the case of the default DNS, this will be a well known IP that is used
	// as the default nameserver for pods that are using the default ClusterFirst DNS policy.
	//
	// In general, this IP can be specified in a pod's spec.dnsConfig.nameservers list
	// or used explicitly when performing name resolution from within the cluster.
	// Example: dig foo.com @<service IP>
	//
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	//
	// +required
	ClusterIP string `json:"clusterIP"`

	// clusterDomain is the local cluster DNS domain suffix for DNS services.
	// This will be a subdomain as defined in RFC 1034,
	// section 3.5: https://tools.ietf.org/html/rfc1034#section-3.5
	// Example: "cluster.local"
	//
	// More info: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service
	//
	// +required
	ClusterDomain string `json:"clusterDomain"`

	// conditions provide information about the state of the DNS on the cluster.
	//
	// These are the supported DNS conditions:
	//
	//   * Available
	//   - True if the following conditions are met:
	//     * DNS controller daemonset is available.
	//   - False if any of those conditions are unsatisfied.
	//
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []OperatorCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSList contains a list of DNS
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type DNSList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DNS `json:"items"`
}
