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
// RuleSource - Schema Registration
// -----------------------------------------------------------------------------

func init() {
	SchemeBuilder.Register(&RuleSource{}, &RuleSourceList{})
}

// -----------------------------------------------------------------------------
// RuleSource - Constants
// -----------------------------------------------------------------------------

const (
	// AnnotationSkipValidation controls per-fragment Coraza rule validation on
	// a RuleSource. When set to "false", per-source validation is skipped
	// (the aggregated RuleSet validation still runs).
	AnnotationSkipValidation = Group + "/rule-validation"
)

// -----------------------------------------------------------------------------
// RuleSource
// -----------------------------------------------------------------------------

// RuleSource holds SecLang WAF rule text for consumption by RuleSet resources.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:validation:XValidation:rule="has(self.spec.rules) && self.spec.rules != \"\"",message="rules must be set and non-empty"
type RuleSource struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata.
	//
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the rule content.
	//
	// +required
	Spec RuleSourceSpec `json:"spec,omitzero"`

	// status defines the observed state of RuleSource.
	//
	// +optional
	Status RuleSourceStatus `json:"status,omitempty,omitzero"`
}

// RuleSourceList contains a list of RuleSource resources.
//
// +kubebuilder:object:root=true
type RuleSourceList struct {
	metav1.TypeMeta `json:",inline"`

	// ListMeta is standard list metadata.
	//
	// +optional
	metav1.ListMeta `json:"metadata,omitzero"`

	// Items is the list of RuleSources.
	//
	// +required
	Items []RuleSource `json:"items"`
}

// -----------------------------------------------------------------------------
// RuleSource - Spec
// -----------------------------------------------------------------------------

// RuleSourceSpec defines the content of a RuleSource.
type RuleSourceSpec struct {
	// rules contains SecLang rule text.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1572864
	Rules string `json:"rules,omitempty"`
}

// -----------------------------------------------------------------------------
// RuleSource - Status
// -----------------------------------------------------------------------------

// RuleSourceStatus defines the observed state of RuleSource.
// +kubebuilder:validation:MinProperties=1
type RuleSourceStatus struct {
	// conditions represent the current state of the RuleSource resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Ready": the operator has validated the RuleSource spec.rules (Coraza), or validation was skipped by annotation
	// - "Degraded": per-source rule validation failed
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
