package v1alpha1

// ScopeType is one of ControlPlane or WorkerPool
// +kubebuilder:validation:Enum=ControlPlane;WorkerPool
type ScopeType string

const (
	// ControlPlane is used for insights that are related to the control plane (including control plane pool or nodes)
	ControlPlaneScope ScopeType = "ControlPlane"
	// WorkerPool is used for insights that are related to a worker pools and nodes (excluding control plane)
	WorkerPoolScope ScopeType = "WorkerPool"
)
