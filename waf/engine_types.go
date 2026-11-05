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
// Engine - Schema Registration
// -----------------------------------------------------------------------------

func init() {
	SchemeBuilder.Register(&Engine{}, &EngineList{})
}

// -----------------------------------------------------------------------------
// Engine
// -----------------------------------------------------------------------------

// Engine represents an instance of a Web Application Firewall (WAF) engine.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="RuleSet",type=string,JSONPath=`.spec.ruleSet.name`
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.target.provider`
// +kubebuilder:printcolumn:name="Target Type",type=string,JSONPath=`.spec.target.type`
// +kubebuilder:printcolumn:name="Target Name",type=string,JSONPath=`.spec.target.name`
// +kubebuilder:printcolumn:name="Failure Policy",type=string,JSONPath=`.spec.failurePolicy`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type Engine struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata.
	//
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Engine.
	//
	// +required
	Spec EngineSpec `json:"spec,omitzero"`

	// status defines the observed state of Engine.
	//
	// +optional
	Status *EngineStatus `json:"status,omitempty"`
}

// EngineList contains a list of Engine resources.
//
// +kubebuilder:object:root=true
type EngineList struct {
	metav1.TypeMeta `json:",inline"`

	// ListMeta is standard list metadata.
	//
	// +optional
	metav1.ListMeta `json:"metadata,omitzero"`

	// Items is the list of Engines.
	//
	// +required
	Items []Engine `json:"items"`
}

// -----------------------------------------------------------------------------
// Engine - Spec
// -----------------------------------------------------------------------------

// EngineSpec defines the desired state of an Engine.
//
// +kubebuilder:validation:XValidation:rule="!has(self.driver) || !has(self.driver.type) || (self.target.provider == 'Istio' && self.driver.type == 'wasm')",message="driver type must be compatible with the target provider (Istio supports wasm)"
type EngineSpec struct {
	// ruleSet specifies the RuleSet resource that will be used to load rules
	// into the Engine. The referenced RuleSet must be in the same namespace
	// as the Engine.
	//
	// +required
	RuleSet RuleSetReference `json:"ruleSet,omitzero"`

	// target identifies the workload that the Engine protects. The operator
	// derives the workload selector from this reference (e.g., for Gateway
	// targets, the GEP-1762 gateway-name label is used).
	//
	// +required
	Target EngineTarget `json:"target,omitzero"`

	// failurePolicy determines the behavior when the WAF is not ready or
	// encounters errors. Valid values are:
	//
	// - "Fail": Block traffic when the WAF is not ready or encounters errors
	// - "Allow": Allow traffic through when the WAF is not ready or encounters errors
	//
	// When omitted, this means the user has no opinion and the platform
	// will choose a reasonable default, which is subject to change over time.
	//
	// The current default is fail.
	//
	// +optional
	// +default="fail"
	FailurePolicy FailurePolicy `json:"failurePolicy,omitempty"`

	// ruleSetCacheServer contains configuration for the ruleset cache server.
	//
	// When omitted, no cache server will be used and no rulesets will be
	// dynamically loaded. This implies that your Engine will be deployed with
	// all rules statically embedded.
	//
	// +optional
	RuleSetCacheServer *RuleSetCacheServerConfig `json:"ruleSetCacheServer,omitempty"`

	// observability configures WAF dataplane observability independently of the
	// driver used to deploy the Engine. The operator translates this intent into
	// driver-specific settings.
	//
	// +optional
	Observability ObservabilityConfig `json:"observability,omitempty,omitzero"`

	// driver configures the mechanism used to deploy the WAF filter into the
	// target workload. When omitted, the operator uses a default driver for the
	// underlying Engine (eg.: WASM for Istio)
	//
	// +optional
	Driver DriverConfig `json:"driver,omitempty,omitzero"`
}

// ObservabilityConfig configures WAF dataplane observability for an Engine.
type ObservabilityConfig struct {
	// mode controls whether WAF dataplane telemetry is available to the target
	// platform's observability pipeline.
	//
	// +required
	Mode ObservabilityMode `json:"mode,omitempty"`
}

// ObservabilityMode controls whether WAF dataplane observability is enabled.
//
// +kubebuilder:validation:Enum=Enabled;Disabled
type ObservabilityMode string

const (
	// ObservabilityModeEnabled enables WAF dataplane observability.
	ObservabilityModeEnabled ObservabilityMode = "Enabled"

	// ObservabilityModeDisabled disables WAF dataplane observability.
	ObservabilityModeDisabled ObservabilityMode = "Disabled"
)

// -----------------------------------------------------------------------------
// Engine - Status
// -----------------------------------------------------------------------------

// EngineStatus defines the observed state of Engine.
// +kubebuilder:validation:MinProperties=0
type EngineStatus struct {
	// conditions represent the current state of the Engine resource.
	// Each condition has a unique type and reflects the status of a specific
	// aspect of the resource.
	//
	// Standard condition types include:
	// - "Accepted": the target is valid and not contested by another Engine.
	//    Reasons: "Accepted", "TargetNotFound", "TargetConflict"
	// - "Ready": the engine has been successfully deployed and is operational
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	// - "ObservabilityReady": optional Telemetry was created and is available
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

// -----------------------------------------------------------------------------
// Engine - Failure Policy
// -----------------------------------------------------------------------------

// FailurePolicy describes the failure policy for the Engine.
//
// +kubebuilder:validation:Enum=fail;allow
type FailurePolicy string

const (
	// FailurePolicyFail blocks traffic when the Engine is not ready or encounters
	// errors.
	FailurePolicyFail FailurePolicy = "fail"

	// FailurePolicyAllow allows traffic through when the Engine is not ready or
	// encounters errors.
	FailurePolicyAllow FailurePolicy = "allow"
)

// -----------------------------------------------------------------------------
// Engine - Reference Types
// -----------------------------------------------------------------------------

// RuleSetReference is a reference to a RuleSet resource.
type RuleSetReference struct {
	// name is the name of the RuleSet in the same namespace as the Engine.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`
}
