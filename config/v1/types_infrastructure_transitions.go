package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TopologyState describes the control-plane and infrastructure topology at one
// end of a topology transition.
type TopologyState struct {
	// controlPlaneTopology is the topology of the control-plane nodes. Valid values
	// are SingleReplica and HighlyAvailable. When set to SingleReplica, operators
	// avoid spending resources for high availability. When set to HighlyAvailable,
	// operators configure high availability as much as possible.
	// controlPlaneTopology is required.
	// +kubebuilder:validation:Enum=SingleReplica;HighlyAvailable
	// +required
	ControlPlaneTopology TopologyMode `json:"controlPlaneTopology,omitempty"`

	// infrastructureTopology is the topology of infrastructure services. Valid
	// values are SingleReplica and HighlyAvailable. When set to SingleReplica,
	// operators avoid spending resources for high availability. When set to
	// HighlyAvailable, operators configure high availability as much as possible.
	// infrastructureTopology is required.
	// +kubebuilder:validation:Enum=SingleReplica;HighlyAvailable
	// +required
	InfrastructureTopology TopologyMode `json:"infrastructureTopology,omitempty"`
}

type TopologyTransitionStatus struct {
	// conditions reports whether supported transitions have been evaluated.
	// TopologyTransitionsEvaluated is Unknown before evaluation, True when evaluation
	// succeeds (even if no transitions are supported), and False when evaluation fails.
	// An absent condition means evaluation has not completed.
	// At most one condition is present.
	// +kubebuilder:validation:MaxItems=1
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// supportedTransitions represents the transitions that are valid for this cluster.
	// An empty list means that no transitions are currently supported from the current topology.
	// At most one transition is supported currently (SNO to HA Compact)
	// +kubebuilder:validation:MaxItems=1
	// +required
	// +listType=atomic
	SupportedTransitions []TopologyTransition `json:"supportedTransitions"`

	// currentTransition is omitted until a topology transition starts.
	// +optional
	CurrentTransition *TopologyTransitionProgress `json:"currentTransition,omitempty"`
}

const (
	// TopologyTransitionsEvaluatedConditionType indicates whether available transitions have been evaluated.
	TopologyTransitionsEvaluatedConditionType = "TopologyTransitionsEvaluated"
)

// TopologyTransitionProgress describes a topology transition that has started.
type TopologyTransitionProgress struct {
	// state indicates the current state of a triggered transition.
	// Valid values are "Completed" when the transition was successfully applied,
	// "Partial" when it was not completely applied or is still in progress, and
	// "Failed" when it failed to apply.
	// +kubebuilder:validation:Enum=Completed;Partial;Failed
	// +required
	State TransitionState `json:"state,omitempty"`

	// reason indicates why current state is as reported.
	// It must be between 1 and 128 characters long.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +required
	Reason string `json:"reason,omitempty"`

	// message is human-readable information about the reason for the current state.
	// It must be between 1 and 2048 characters long.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +required
	Message string `json:"message,omitempty"`

	// startedTime is the time at which the transition was started. When omitted, the start time is not available.
	// +optional
	StartedTime *metav1.Time `json:"startedTime,omitempty"`

	// completionTime is when the transition was fully applied. It is omitted while a transition is being applied.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

// TransitionState tracks the last observed state of a requested transition.
type TransitionState string

const (
	// CompletedTransition indicates an transition was successfully
	// completed on the cluster.
	CompletedTransition TransitionState = "Completed"
	// PartialTransition indicates a transition was never completely applied
	// or is currently being applied.
	PartialTransition TransitionState = "Partial"
	// FailedTransition indicates a transition failed to be applied.
	FailedTransition TransitionState = "Failed"
)

type TopologyTransition struct {
	// source is the control-plane and infrastructure topology this transition starts
	// from. It must equal the current topology in the corresponding status fields.
	// source is required.
	// +required
	Source TopologyState `json:"source,omitempty,omitzero"`

	// target is the control-plane and infrastructure topology this transition would
	// move to. target is required.
	// +required
	Target TopologyState `json:"target,omitempty,omitzero"`

	// reason is a CamelCase machine-readable explanation of the availability, e.g.
	// PreflightCheckFailed. The set of reasons is diagnostic and not exhaustive.
	// When omitted, no machine-readable explanation is available.
	// Must start with an uppercase letter and contain only alphanumeric characters,
	// and must be between 1 and 128 characters long.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:XValidation:rule=`self.matches('^[A-Z][A-Za-z0-9]*$')`,message="reason must be CamelCase, matching ^[A-Z][A-Za-z0-9]*$"
	// +optional
	Reason string `json:"reason,omitempty"`

	// message is a human-readable explanation, primarily for Unavailable
	// transitions (e.g. a concise summary of the failing preconditions). It is for
	// humans only and must not be parsed. It may be truncated by the controller.
	// When omitted, no human-readable explanation is available for the transition.
	// When set, it must be between 1 and 2048 characters long.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +optional
	Message string `json:"message,omitempty"`
}
