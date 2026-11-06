package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// GlobalPullCredential is a cluster-scoped resource declaring a single registry
// credential (one hostname per instance). A controller in machine-config-controller
// merges it into the controller-owned `openshift-config/pull-secret-merged` Secret
// without modifying the original `openshift-config/pull-secret`.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=globalpullcredentials,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/3126
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=config-operator,operatorOrdering=01
// +openshift:enable:FeatureGate=MultiSecretGlobalPullConfig
// +openshift:compatibility-gen:level=4
type GlobalPullCredential struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds the registry credential contributed by this resource.
	// +required
	Spec GlobalPullCredentialSpec `json:"spec,omitzero"`

	// status holds observed values from the cluster. They may not be overridden by user input.
	// +optional
	Status GlobalPullCredentialStatus `json:"status,omitempty,omitzero"`
}

// GlobalPullCredentialSpec holds the specification for a GlobalPullCredential resource.
type GlobalPullCredentialSpec struct {
	// hostname is the registry host that this credential applies to, for example
	// "quay.io" or "registry.example.com:5000". Each hostname may only be claimed
	// by one GlobalPullCredential instance across the cluster.
	//
	// The value must be a DNS hostname of one or more labels separated by dots.
	// Each label must start and end with an alphanumeric character and may contain
	// only alphanumeric characters or '-'. An optional numeric port may be appended
	// as ":{port}" and must be in the range 1 to 65535.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self.split(':')[0].matches('^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\\\\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$')",message="hostname must start and end with an alphanumeric character, and must only contain alphanumeric characters, '-' or '.'"
	// +kubebuilder:validation:XValidation:rule="!self.contains(':') || (self.split(':').size() == 2 && self.split(':')[1].matches('^[0-9]+$') && int(self.split(':')[1]) >= 1 && int(self.split(':')[1]) <= 65535)",message="port must be between 1 and 65535"
	Hostname string `json:"hostname,omitempty"`

	// auth holds the registry credential for this hostname.
	// +required
	Auth RegistryAuth `json:"auth,omitzero"`
}

// RegistryAuth holds the base64-encoded username and password for a single registry.
// The controller decodes these fields and constructs the Docker config.json auth
// value base64("username:password") at reconcile time.
type RegistryAuth struct {
	// username is the base64-encoded registry username. The decoded value must
	// not contain a colon because the controller constructs "username:password".
	// The encoded value is limited to 684 characters (512 decoded bytes).
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=684
	Username []byte `json:"username,omitempty"`

	// password is the base64-encoded registry password or token. The encoded
	// value is limited to 684 characters (512 decoded bytes).
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=684
	Password []byte `json:"password,omitempty"`
}

// GlobalPullCredentialStatus holds observed state for a GlobalPullCredential resource.
//
// +kubebuilder:validation:MinProperties=1
type GlobalPullCredentialStatus struct {
	// conditions represent the latest available observations of this resource's state.
	// The known condition type is "Applied".
	// Applied=True means the controller merged the credential declared by this
	// generation into openshift-config/pull-secret-merged.
	// Applied=False means the credential was not merged. If the hostname overlaps
	// with the baseline pull-secret, the reason is BaselinePrecedence and the
	// baseline credential is used instead.
	// Applied=Unknown or an absent condition means no reconciliation result is available.
	// Applied does not indicate successful registry authentication or propagation to nodes.
	// When omitted, no observations are available. When present, the list must contain one condition.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=1
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// GlobalPullCredentialList is a collection of GlobalPullCredential resources.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +openshift:compatibility-gen:level=4
type GlobalPullCredentialList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is a list of GlobalPullCredential resources.
	Items []GlobalPullCredential `json:"items"`
}
