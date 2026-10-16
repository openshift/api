/*
Copyright Coraza Kubernetes Operator contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// -----------------------------------------------------------------------------
// RuleData - Schema Registration
// -----------------------------------------------------------------------------

func init() {
	SchemeBuilder.Register(&RuleData{}, &RuleDataList{})
}

// -----------------------------------------------------------------------------
// RuleData
// -----------------------------------------------------------------------------

// RuleData holds data file content (e.g. for @pmFromFile) for consumption by
// RuleSet resources. Each entry in spec.files maps a filename to its content.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:validation:XValidation:rule="has(self.spec.files) && size(self.spec.files) > 0",message="files must be non-empty"
// +kubebuilder:validation:XValidation:rule="has(self.spec.files) ? self.spec.files.all(k, k.matches('^[-._a-zA-Z0-9]+$') && size(k) <= 253) : true",message="files keys must be valid data file names (alphanumeric, '-', '_', '.'; max 253 chars)"
type RuleData struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata.
	//
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the data file content.
	//
	// +required
	Spec RuleDataSpec `json:"spec,omitzero"`

	// status defines the observed state of RuleData.
	//
	// +optional
	Status RuleDataStatus `json:"status,omitempty,omitzero"`
}

// RuleDataList contains a list of RuleData resources.
//
// +kubebuilder:object:root=true
type RuleDataList struct {
	metav1.TypeMeta `json:",inline"`

	// ListMeta is standard list metadata.
	//
	// +optional
	metav1.ListMeta `json:"metadata,omitzero"`

	// Items is the list of RuleData.
	//
	// +required
	Items []RuleData `json:"items"`
}

// -----------------------------------------------------------------------------
// RuleData - Spec
// -----------------------------------------------------------------------------

// RuleDataSpec defines the content of a RuleData resource.
type RuleDataSpec struct {
	// files maps filenames to file content, used for @pmFromFile data.
	//
	// +required
	// +kubebuilder:validation:MinProperties=1
	Files map[string]string `json:"files,omitempty"`
}

// -----------------------------------------------------------------------------
// RuleData - Status
// -----------------------------------------------------------------------------

// RuleDataStatus defines the observed state of RuleData.
// +kubebuilder:validation:MinProperties=1
type RuleDataStatus struct {
	// conditions represent the current state of the RuleData resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Ready": the RuleData has been loaded by at least one RuleSet
	//
	// The status of each condition is one of True, False, or Unknown.
	//
	// +listType=map
	// +listMapKey=type
	// +patchStrategy=merge
	// +patchMergeKey=type
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}
