package v1

// SecurityProfile defines the schema for a security profile. This object is used
// by operators to apply network security settings to operands.
// +union
type SecurityProfile struct {
	// type is one of Old, Intermediate, Modern or Custom. Custom provides
	// the ability to specify individual security profile parameters. Old,
	// Intermediate and Modern are security profiles based on:
	//
	// https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations
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
	//     - "!aNULL"
	//     - "!eNULL"
	//     - "!EXPORT"
	//     - "!RC4"
	//     - "!MD5"
	//     - "!PSK"
	//     - "!RSAPSK"
	//     - "!aDH"
	//     - "!aECDH"
	//     - "!EDH-DSS-DES-CBC3-SHA"
	//     - "!KRB5-DES-CBC3-SHA"
	//     - "!SRP"
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
	//     - "!DSS"
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
	// custom is a user-defined security profile. Be extremely careful using a custom
	// profile as invalid configurations can be catastrophic. An example custom profile
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
	// +optional
	// +nullable
	Custom *CustomSecurityProfile `json:"custom,omitempty"`
}

// OldSecurityProfile is a security profile based on:
// https://wiki.mozilla.org/Security/Server_Side_TLS#Old_backward_compatibility
type OldSecurityProfile struct{}

// IntermediateSecurityProfile is a security profile based on:
// https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
type IntermediateSecurityProfile struct{}

// ModernSecurityProfile is a security profile based on:
// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
type ModernSecurityProfile struct{}

// CustomSecurityProfile is a user-defined security profile. Be extremely careful using
// a custom profile as invalid configurations can be catastrophic.
type CustomSecurityProfile struct {
	SecurityProfileSpec `json:",inline"`
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

// SecurityProfileSpec is the desired behavior of a SecurityProfile.
type SecurityProfileSpec struct {
	// ciphers is used to specify the cipher algorithms that are negotiated
	// during the SSL/TLS handshake. Preface a cipher with a "!" to disable
	// a specific cipher from being negotiated. Note that disabled ciphers must
	// be quoted due to the leading "!". For example, to use 3DES but not
	// EDH-DSS-DES-CBC3-SHA (yaml):
	//
	//   ciphers:
	//     - 3DES
	//     - "!EDH-DSS-DES-CBC3-SHA"
	//
	Ciphers []string `json:"ciphers"`
	// securityProtocol is used to specify one or more encryption protocols
	// that are negotiated during the SSL/TLS handshake. For example, to use
	// TLS versions 1.1, 1.2 and 1.3 (yaml):
	//
	//   securityProtocol:
	//     minimumVersion: TLSv1.1
	//     maximumVersion: TLSv1.3
	//
	SecurityProtocol SecurityProtocol `json:"securityProtocol"`
	// dhParamSize sets the maximum size of the Diffie-Hellman parameters used for generating
	// the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange. The final size
	// will try to match the size of the server's RSA (or DSA) key (e.g, a 2048 bits temporary
	// DH key for a 2048 bits RSA key), but will not exceed this maximum value.
	//
	// Available DH Parameter sizes are:
	//
	//   "2048": A Diffie-Hellman parameter of 2048 bits.
	//   "1024": A Diffie-Hellman parameter of 1024 bits.
	//
	// For example, to use a Diffie-Hellman parameter of 2048 bits (yaml):
	//
	//   dhParamSize: 2048
	//
	DHParamSize DHParamSize `json:"dhParamSize"`
}

// SecurityProtocol defines one or more security protocols used to secure network connections.
type SecurityProtocol struct {
	// minimumVersion enforces use of the specified SecurityProtocolVersion or newer
	// on SSL connections. minimumVersion must be lower than or equal to maximumVersion.
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
	//   "TLSv1.3": Version 1.3 of the TLS security protocol used for securing network connections.
	//   "TLSv1.2": Version 1.2 of the TLS security protocol used for securing network connections.
	//   "TLSv1.1": Version 1.1 of the TLS security protocol used for securing network connections.
	//   "TLSv1.0": Version 1.0 of the TLS security protocol used for securing network connections.
	//
	MinimumVersion SecurityProtocolVersion `json:"minimumVersion"`
	// maximumVersion enforces use of the specified SecurityProtocolVersion or older
	// on SSL connections. maximumVersion must be higher than or equal to minimumVersion.
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
	MaximumVersion SecurityProtocolVersion `json:"maximumVersion"`
}

// SecurityProtocolVersion is a way to specify the TLS security protocol version used to
// secure network connections.
type SecurityProtocolVersion string

const (
	// TLSv1.0 is version 1.0 of the TLS security protocol used for securing network connections.
	SecurityProtocolTLS10Version SecurityProtocolVersion = "TLSv1.0"
	// TLSv1.1 is version 1.1 of the TLS security protocol used for securing network connections.
	SecurityProtocolTLS11Version SecurityProtocolVersion = "TLSv1.1"
	// TLSv1.2 is version 1.2 of the TLS security protocol used for securing network connections.
	SecurityProtocolTLS12Version SecurityProtocolVersion = "TLSv1.2"
	// TLSv1.3 is version 1.3 of the TLS security protocol used for securing network connections.
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
