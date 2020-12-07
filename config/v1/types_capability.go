package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Capability holds the expected high-availability mode of the cluster. The canonical name is `cluster`
type Capability struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec CapabilitySpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status CapabilityStatus `json:"status"`
}

type CapabilitySpec struct {
	// HighAvailabilityMode express the high-availability expectations.
	//
	// The default is 'Full', which represents the behavior operators have in a \"normal\" cluster.
	// The 'None' mode will be used in single-node deployments (developer and production) for example,
	// and the operators should not configure the operand for highly-available operation
	//
	// Once set, this field cannot be changed.
	// +kubebuilder:default=Full
	HighAvailabilityMode HighAvailabilityMode `json:"highAvailabilityMode"`
}

// HighAvailabilityMode defines the high-availability mode of the cluster.
// +kubebuilder:validation:Enum=Full;None
type HighAvailabilityMode string

const (
	// "Full" is for operators to configure high-availability as much as possible.
	FullHighAvailabilityMode HighAvailabilityMode = "Full"

	// "None" is for operators to avoid spending resources for high-availability purpose.
	NoneHighAvailabilityMode HighAvailabilityMode = "None"
)

type CapabilityStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CapabilityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Capability `json:"items"`
}
