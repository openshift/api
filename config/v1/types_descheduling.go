package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Descheduler holds cluster-wide config information to run the Kubernetes Descheduler
// and influence its placement decisions. The canonical name for this config is `cluster`.
type Descheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec DeschedulerSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status DeschedulerStatus `json:"status"`
}

type DeschedulerSpec struct {
	// Policy contains the name of an upstream Descheduler policy configmap in the openshift-kube-descheduler-operator project
	Policy ConfigMapNameReference `json:"policy"`
	// DeschedulingIntervalSeconds is the number of seconds between descheduler runs
	// +optional
	DeschedulingIntervalSeconds *int32 `json:"deschedulingIntervalSeconds,omitempty"`
	// DryRun sets whether to run the Descheduler in dry-run mode (with no actual evictions)
	// +optional
	DryRun bool `json:"dryRun,omitempty"`
}

type DeschedulerStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeschedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Descheduler `json:"items"`
}
