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

// -----------------------------------------------------------------------------
// Engine - Driver Config
// -----------------------------------------------------------------------------

// DriverConfig configures how the WAF filter is deployed into the target.
// When omitted from the Engine spec, the operator uses a default driver
// (currently wasm for Istio).
//
// TODO: When using a Gateway resource, the engine reconciler MUST recognize
// what GatewayAPI controller was used and set the better default driver.
//
// Exactly one driver-specific configuration must match the selected type.
//
// +kubebuilder:validation:XValidation:rule="self.type == 'wasm' ? has(self.wasm) : true",message="wasm config is required when type is wasm"
// +kubebuilder:validation:MinProperties=0
type DriverConfig struct {
	// type selects the driver mechanism used to deploy the WAF filter.
	//
	// +required
	Type DriverType `json:"type,omitempty"`

	// wasm contains configuration specific to the WASM driver.
	//
	// +optional
	Wasm *WasmDriverConfig `json:"wasm,omitempty"`
}

// -----------------------------------------------------------------------------
// Engine - Driver Type
// -----------------------------------------------------------------------------

// DriverType specifies the mechanism used to deploy the WAF filter.
//
// +kubebuilder:validation:Enum=wasm
type DriverType string

const (
	// DriverTypeWasm deploys the WAF as a WebAssembly plugin.
	DriverTypeWasm DriverType = "wasm"
)

// -----------------------------------------------------------------------------
// Engine - WASM Driver Config
// -----------------------------------------------------------------------------

// WasmDriverConfig defines configuration for deploying the Engine as a WASM
// plugin.
//
// +kubebuilder:validation:MinProperties=0
// +kubebuilder:validation:XValidation:rule="!has(self.image) || self.image.matches('^oci://')",message="image must start with oci:// when set"
// +kubebuilder:validation:XValidation:rule="!has(self.image) || size(self.image) <= 1024",message="image must be at most 1024 characters when set"
type WasmDriverConfig struct {
	// image is the OCI image reference for the Coraza WASM plugin.
	// If omitted the operator uses its configured default WASM OCI reference
	// (--default-wasm-image / CORAZA_DEFAULT_WASM_IMAGE).
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Image string `json:"image,omitempty"`

	// imagePullSecret is the name of a Kubernetes Secret in the same namespace
	// as the Engine that contains Docker registry credentials for pulling the
	// WASM OCI image.
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
}

// MaxImageLen must match the CEL size constraint on WasmDriverConfig.Image.
const MaxImageLen = 1024
