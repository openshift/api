package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WebConsoleConfiguration holds the necessary configuration options for serving the web console
type WebConsoleConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// ServingInfo is the HTTP serving information for these assets
	ServingInfo HTTPServingInfo `json:"servingInfo" protobuf:"bytes,1,opt,name=servingInfo"`

	// PublicURL is where you can find the asset server (TODO do we really need this?)
	PublicURL string `json:"publicURL" protobuf:"bytes,2,opt,name=publicURL"`

	// LogoutURL is an optional, absolute URL to redirect web browsers to after logging out of the web console.
	// If not specified, the built-in logout page is shown.
	LogoutURL string `json:"logoutURL" protobuf:"bytes,3,opt,name=logoutURL"`

	// MasterPublicURL is how the web console can access the OpenShift v1 server
	MasterPublicURL string `json:"masterPublicURL" protobuf:"bytes,4,opt,name=masterPublicURL"`

	// LoggingPublicURL is the public endpoint for logging (optional)
	LoggingPublicURL string `json:"loggingPublicURL" protobuf:"bytes,5,opt,name=loggingPublicURL"`

	// MetricsPublicURL is the public endpoint for metrics (optional)
	MetricsPublicURL string `json:"metricsPublicURL" protobuf:"bytes,6,opt,name=metricsPublicURL"`

	// ExtensionScripts are file paths on the asset server files to load as scripts when the Web
	// Console loads
	ExtensionScripts []string `json:"extensionScripts" protobuf:"bytes,7,rep,name=extensionScripts"`

	// ExtensionProperties are key(string) and value(string) pairs that will be injected into the console under
	// the global variable OPENSHIFT_EXTENSION_PROPERTIES
	ExtensionProperties map[string]string `json:"extensionProperties" protobuf:"bytes,8,rep,name=extensionProperties"`

	// ExtensionStylesheets are file paths on the asset server files to load as stylesheets when
	// the Web Console loads
	ExtensionStylesheets []string `json:"extensionStylesheets" protobuf:"bytes,9,rep,name=extensionStylesheets"`

	// Extensions are files to serve from the asset server filesystem under a subcontext
	Extensions []AssetExtensionsConfig `json:"extensions" protobuf:"bytes,10,rep,name=extensions"`

	// ExtensionDevelopment when true tells the asset server to reload extension scripts and
	// stylesheets for every request rather than only at startup. It lets you develop extensions
	// without having to restart the server for every change.
	ExtensionDevelopment bool `json:"extensionDevelopment" protobuf:"varint,11,opt,name=extensionDevelopment"`
}

// AssetExtensionsConfig holds the necessary configuration options for asset extensions
type AssetExtensionsConfig struct {
	// SubContext is the path under /<context>/extensions/ to serve files from SourceDirectory
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// SourceDirectory is a directory on the asset server to serve files under Name in the Web
	// Console. It may have nested folders.
	SourceDirectory string `json:"sourceDirectory" protobuf:"bytes,2,opt,name=sourceDirectory"`
	// HTML5Mode determines whether to redirect to the root index.html when a file is not found.
	// This is needed for apps that use the HTML5 history API like AngularJS apps with HTML5
	// mode enabled. If HTML5Mode is true, also rewrite the base element in index.html with the
	// Web Console's context root. Defaults to false.
	HTML5Mode bool `json:"html5Mode" protobuf:"varint,3,opt,name=html5Mode"`
}

// HTTPServingInfo holds configuration for serving HTTP
type HTTPServingInfo struct {
	// ServingInfo is the HTTP serving information
	ServingInfo `json:",inline" protobuf:"bytes,1,opt,name=servingInfo"`
	// MaxRequestsInFlight is the number of concurrent requests allowed to the server. If zero, no limit.
	MaxRequestsInFlight int64 `json:"maxRequestsInFlight" protobuf:"varint,2,opt,name=maxRequestsInFlight"`
	// RequestTimeoutSeconds is the number of seconds before requests are timed out. The default is 60 minutes, if
	// -1 there is no limit on requests.
	RequestTimeoutSeconds int64 `json:"requestTimeoutSeconds" protobuf:"varint,3,opt,name=requestTimeoutSeconds"`
}

// ServingInfo holds information about serving web pages
type ServingInfo struct {
	// BindAddress is the ip:port to serve on
	BindAddress string `json:"bindAddress" protobuf:"bytes,1,opt,name=bindAddress"`
	// BindNetwork is the type of network to bind to - defaults to "tcp4", accepts "tcp",
	// "tcp4", and "tcp6"
	BindNetwork string `json:"bindNetwork" protobuf:"bytes,2,opt,name=bindNetwork"`
	// CertInfo is the TLS cert info for serving secure traffic.
	// this is anonymous so that we can inline it for serialization
	CertInfo `json:",inline" protobuf:"bytes,3,opt,name=certInfo"`
	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	ClientCA string `json:"clientCA" protobuf:"bytes,4,opt,name=clientCA"`
	// NamedCertificates is a list of certificates to use to secure requests to specific hostnames
	NamedCertificates []NamedCertificate `json:"namedCertificates" protobuf:"bytes,5,rep,name=namedCertificates"`
	// MinTLSVersion is the minimum TLS version supported.
	// Values must match version names from https://golang.org/pkg/crypto/tls/#pkg-constants
	MinTLSVersion string `json:"minTLSVersion,omitempty" protobuf:"bytes,6,opt,name=minTLSVersion"`
	// CipherSuites contains an overridden list of ciphers for the server to support.
	// Values must match cipher suite IDs from https://golang.org/pkg/crypto/tls/#pkg-constants
	CipherSuites []string `json:"cipherSuites,omitempty" protobuf:"bytes,7,rep,name=cipherSuites"`
}

// CertInfo relates a certificate with a private key
type CertInfo struct {
	// CertFile is a file containing a PEM-encoded certificate
	CertFile string `json:"certFile" protobuf:"bytes,1,opt,name=certFile"`
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string `json:"keyFile" protobuf:"bytes,2,opt,name=keyFile"`
}

// NamedCertificate specifies a certificate/key, and the names it should be served for
type NamedCertificate struct {
	// Names is a list of DNS names this certificate should be used to secure
	// A name can be a normal DNS name, or can contain leading wildcard segments.
	Names []string `json:"names" protobuf:"bytes,1,rep,name=names"`
	// CertInfo is the TLS cert info for serving secure traffic
	CertInfo `json:",inline" protobuf:"bytes,2,opt,name=certInfo"`
}
