package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Project holds cluster-wide information about Project.  The canonical name is `cluster`
type Project struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	Spec ProjectSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	Status ProjectStatus `json:"status"`
}

// TemplateReference references a project request template in a 'openshift-config' namespace.
type TemplateReference struct {
	// name is the metadata.name of the referenced project request template
	Name string `json:"name"`
}

type ProjectSpec struct {
	// projectRequestMessage is the string presented to a user if they are unable to request a project via the projectrequest api endpoint
	ProjectRequestMessage string `json:"projectRequestMessage"`

	// projectRequestTemplate is the template to use for creating projects in response to projectrequest.
	// This must point to a template in 'openshift-config' namespace. It is optional.
	// If it is not specified, a default template is used.
	ProjectRequestTemplate *TemplateReference `json:"projectRequestTemplate"`
}

type ProjectStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}
