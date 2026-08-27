package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kmsplugins,scope=Namespaced
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'cluster'",message="kmsplugin is a singleton per namespace, .metadata.name must be 'cluster'"
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/TBD
// +openshift:file-pattern=cvoRunLevel=0000_20,operatorName=kube-apiserver,operatorOrdering=02
// +openshift:compatibility-gen:level=4
// +openshift:enable:FeatureGate=KMSEncryption

// KMSPlugin defines how the platform runs a KMS encryption provider plugin sidecar.
// A KMS provider operator installed via OLM reconciles provider-specific configuration
// and publishes the container runtime configuration in status.runtime.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type KMSPlugin struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +required
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is reserved for future use.
	// Operators publish runtime configuration in status.runtime.
	// +optional
	Spec *KMSPluginSpec `json:"spec,omitempty"`

	// status is the most recently observed status of the KMSPlugin.
	// +optional
	Status *KMSPluginStatus `json:"status,omitempty"`
}

// KMSPluginSpec is reserved for future use.
// +kubebuilder:validation:MinProperties=0
type KMSPluginSpec struct{}

// KMSPluginStatus defines the observed status of KMSPlugin.
// +kubebuilder:validation:MinProperties=0
type KMSPluginStatus struct {
	// runtime describes how the platform should run the KMS plugin sidecar.
	// The KMS provider operator must populate this before the platform can deploy
	// the plugin. When omitted, the platform cannot proceed with KMS encryption.
	//
	// +optional
	Runtime KMSPluginRuntime `json:"runtime,omitempty,omitzero"`
}

// KMSPluginRuntime describes the container configuration for a KMS plugin sidecar.
// The platform injects -listen-address and manages lifecycle, resources, security context,
// and mounting Secrets and ConfigMaps from the operator namespace at well-known injection points.
// Operators must not set -listen-address in args, either as -listen-address=<path>
// or as a separate -listen-address flag.
//
// +kubebuilder:validation:XValidation:rule="self.args.all(a, a != '-listen-address' && !a.startsWith('-listen-address='))",message="args must not include -listen-address; the platform injects this argument"
type KMSPluginRuntime struct {
	// image is the digest-pinned OCI image for the KMS plugin.
	//
	// The image must be a fully qualified OCI image pull spec with a SHA256 digest.
	// The format is: host[:port][/namespace]/name@sha256:<digest>
	// where the digest must be 64 characters long and consist only of lowercase hexadecimal characters, a-f and 0-9.
	// The total length must be between 75 and 447 characters.
	//
	// Short names (e.g., "vault-plugin" or "hashicorp/vault-plugin") are not allowed.
	// The registry hostname must be included and must contain at least one dot.
	// Image tags (e.g., ":latest", ":v1.0.0") are not allowed.
	//
	// +kubebuilder:validation:MinLength=75
	// +kubebuilder:validation:MaxLength=447
	// +kubebuilder:validation:XValidation:rule=`(self.split('@').size() == 2 && self.split('@')[1].matches('^sha256:[a-f0-9]{64}$'))`,message="the OCI Image reference must end with a valid '@sha256:<digest>' suffix, where '<digest>' is 64 characters long"
	// +kubebuilder:validation:XValidation:rule=`(self.split('@')[0].matches('^([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]+(:[0-9]{2,5})?(/[a-zA-Z0-9-_.]+)+$'))`,message="the OCI Image name should follow the host[:port][/namespace]/name format, resembling a valid URL without the scheme. Short names are not allowed, the registry hostname must be included."
	// +required
	Image string `json:"image,omitempty"`

	// args are the command-line arguments passed to the KMS plugin container.
	// The platform prepends -listen-address=<uds-path> before these arguments.
	// Arguments may reference credential files mounted by the platform at well-known
	// injection points under /var/run/kms/ from Secrets and ConfigMaps in the
	// operator namespace.
	//
	// +required
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=4096
	Args []string `json:"args,omitempty"`

	// secrets lists Secrets in the same namespace as the KMSPlugin that the platform
	// mounts at well-known injection points under /var/run/kms/secrets/.
	// Plugin arguments may reference files from these mounted Secrets.
	// When omitted, no Secrets are mounted.
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	Secrets []KMSPluginSecretReference `json:"secrets,omitempty"`

	// configMaps lists ConfigMaps in the same namespace as the KMSPlugin that the platform
	// mounts at well-known injection points under /var/run/kms/config/.
	// Plugin arguments may reference files from these mounted ConfigMaps.
	// When omitted, no ConfigMaps are mounted.
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	ConfigMaps []KMSPluginConfigMapReference `json:"configMaps,omitempty"`
}

// KMSPluginSecretReference references a Secret in the same namespace as the KMSPlugin.
type KMSPluginSecretReference struct {
	// name is the metadata.name of the referenced Secret.
	// The name must be a valid DNS subdomain name: it must contain no more than 253 characters,
	// contain only lowercase alphanumeric characters, '-' or '.', and start and end with an alphanumeric character.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]([a-z0-9\\\\-]*[a-z0-9])?(\\\\.[a-z0-9]([a-z0-9\\\\-]*[a-z0-9])?)*$')",message="name must be a valid DNS subdomain name: contain no more than 253 characters, contain only lowercase alphanumeric characters, '-' or '.', and start and end with an alphanumeric character"
	// +required
	Name string `json:"name,omitempty"`
}

// KMSPluginConfigMapReference references a ConfigMap in the same namespace as the KMSPlugin.
type KMSPluginConfigMapReference struct {
	// name is the metadata.name of the referenced ConfigMap.
	// The name must be a valid DNS subdomain name: it must contain no more than 253 characters,
	// contain only lowercase alphanumeric characters, '-' or '.', and start and end with an alphanumeric character.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z0-9]([a-z0-9\\\\-]*[a-z0-9])?(\\\\.[a-z0-9]([a-z0-9\\\\-]*[a-z0-9])?)*$')",message="name must be a valid DNS subdomain name: contain no more than 253 characters, contain only lowercase alphanumeric characters, '-' or '.', and start and end with an alphanumeric character"
	// +required
	Name string `json:"name,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4

// KMSPluginList contains a list of KMSPlugins.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
type KMSPluginList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of KMSPlugins.
	Items []KMSPlugin `json:"items"`
}
