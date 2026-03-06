/*
Copyright 2024.

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

import "reflect"

// CRDEnabledOption indicates whether the ClusterMonitoring CRD is available.
// +kubebuilder:validation:Enum=True;False
type CRDEnabledOption string

const (
	// CRDEnabledTrue indicates the ClusterMonitoring CRD is installed (e.g. feature gate ClusterMonitoringConfig). CRSpec may be used.
	CRDEnabledTrue CRDEnabledOption = "True"
	// CRDEnabledFalse indicates the ClusterMonitoring CRD is not installed (e.g. non-TechPreview). CRSpec is ignored.
	CRDEnabledFalse CRDEnabledOption = "False"
)

// EffectiveConfigInput holds the configuration sources for computing the effective
// ClusterMonitoring spec. The operator should obtain ConfigMap and CR (when the CRD
// is enabled) and pass them here so that merge is centralized in one place.
//
// When crDEnabled is False (e.g. non-TechPreview), the CR is not available; the
// effective config is computed from Defaults + ConfigMap only, avoiding any
// dependency on the ClusterMonitoring CR.
type EffectiveConfigInput struct {
	// defaults is the platform default spec.
	// +required
	Defaults ClusterMonitoringSpec `json:"defaults,omitzero"`
	// configMapSpec is the spec parsed from the cluster-monitoring ConfigMap, if present. Zero value when not used.
	// +optional
	ConfigMapSpec ClusterMonitoringSpec `json:"configMapSpec,omitempty,omitzero"`
	// crSpec is the spec from the ClusterMonitoring CR, when the CRD is installed. Zero value when CRD is not enabled.
	// +optional
	CRSpec ClusterMonitoringSpec `json:"crSpec,omitempty,omitzero"`
	// crDEnabled indicates whether the ClusterMonitoring CRD is installed. When False, crSpec is ignored.
	// +required
	CRDEnabled CRDEnabledOption `json:"crDEnabled,omitempty"`
}

// ComputeEffectiveConfig returns the effective ClusterMonitoringSpec by merging
// Defaults + ConfigMapSpec + CRSpec (only when CRDEnabled and CRSpec is non-nil).
// This is the single place where config from ConfigMap and CR are merged; the result
// should be the only source of truth for "safe" config the operator applies.
//
// When CRDEnabled is false, the CR is not used, so metrics server (and other) config
// is taken only from Defaults and ConfigMap, and no CR dependency is required.
func ComputeEffectiveConfig(in EffectiveConfigInput) ClusterMonitoringSpec {
	out := (&in.Defaults).DeepCopy()
	if out == nil {
		return in.Defaults
	}

	if !reflect.ValueOf(in.ConfigMapSpec).IsZero() {
		mergeSpecInto(out, &in.ConfigMapSpec)
	}
	if in.CRDEnabled == CRDEnabledTrue && !reflect.ValueOf(in.CRSpec).IsZero() {
		mergeSpecInto(out, &in.CRSpec)
	}

	return *out
}

// mergeSpecInto overlays non-zero fields of src onto dst at the top level.
// Only exported fields of ClusterMonitoringSpec are considered; for each such field,
// if the value in src is not the zero value, it is copied to dst.
func mergeSpecInto(dst, src *ClusterMonitoringSpec) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()
	typ := dstVal.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		srcField := srcVal.Field(i)
		if srcField.IsZero() {
			continue
		}
		dstVal.Field(i).Set(srcField)
	}
}
