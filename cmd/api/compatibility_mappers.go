package main

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// OpenShiftCompatibilityMapper returns the OpenShift compatibility level as defined
// by the "OpenShift 4: Application Compatibility GUIDE".
//
// API groups that end with the suffix *.openshift.io are governed by the OpenShift
// deprecation policy and follow a general mapping between API version exposed and
// corresponding compatibility level unless otherwise specified.
//
// 	API Version  Compatibility Level
//	v1           Level 1
//	v1beta1      Level 2
//	v1alpha1     Level 4
func OpenShiftAPICompatibilityMapper(gv schema.GroupVersion) OpenShiftCompatibilityLevel {
	if !strings.HasSuffix(gv.Group, ".openshift.io") {
		return CompatibilityLevelNone
	}
	switch {
	case VersionIsGenerallyAvailable(gv.Version):
		return CompatibilityLevel1
	case VersionIsPrerelease(gv.Version):
		return CompatibilityLevel2
	case VersionIsExperimental(gv.Version):
		return CompatibilityLevel4
	default:
		return CompatibilityLevelNone
	}
}

// KubernetesAPICompatibilityMapper returns the OpenShift compatibility level for API
// groups that end with `.k8s.io`, as defined by the "OpenShift 4: Application
// Compatibility GUIDE".
//
// API groups that end with the suffix *.k8s.io or have the form version.<name> with
// no suffix are governed by the Kubernetes deprecation policy and follow a general
// mapping between API version exposed and corresponding compatibility level unless
// otherwise specified in the "OpenShift 4: Application Compatibility GUIDE".
//
// 	API Version  Compatibility Level
//	v1           Level 1
//	v1beta1      Level 2
func KubernetesAPICompatibilityMapper(gv schema.GroupVersion) OpenShiftCompatibilityLevel {
	if !strings.HasSuffix(gv.Group, ".k8s.io") && strings.Contains(gv.Group, ".") {
		return CompatibilityLevelNone
	}
	switch {
	case VersionIsGenerallyAvailable(gv.Version):
		return CompatibilityLevel1
	case VersionIsPrerelease(gv.Version):
		return CompatibilityLevel2
	default:
		return CompatibilityLevelNone
	}
}

// MonitoringAPICompatibilityMapper returns the OpenShift compatibility level for API
// groups that end with `monitoring.coreos.com`, as defined by the "OpenShift 4:
// Application Compatibility GUIDE".
//
// API groups that end with suffix monitoring.coreos.com have the following mapping:.
//
// 	API Version  Compatibility Level
//	v1           Level 1
func CoreOSMonitoringAPICompatibilityMapper(gv schema.GroupVersion) OpenShiftCompatibilityLevel {
	if gv.Group != "monitoring.coreos.com" {
		return CompatibilityLevelNone
	}
	switch {
	case VersionIsGenerallyAvailable(gv.Version):
		return CompatibilityLevel1
	default:
		return CompatibilityLevelNone
	}
}

// MonitoringAPICompatibilityMapper returns the OpenShift compatibility level for API
// groups that end with `monitoring.coreos.com`, as defined by the "OpenShift 4:
// Application Compatibility GUIDE".
//
// API groups that end with suffix monitoring.coreos.com have the following mapping:.
//
// 	API Version  Compatibility Level
//	v1           Level 1
//	v1beta1      Level 3
//	v1alpha1     Level 3
func CoreOSOperatorsAPICompatibilityMapper(gv schema.GroupVersion) OpenShiftCompatibilityLevel {
	if gv.Group != "operators.coreos.com" {
		return CompatibilityLevelNone
	}
	switch {
	case VersionIsGenerallyAvailable(gv.Version):
		return CompatibilityLevel1
	case VersionIsPrerelease(gv.Version):
		// TODO verify, the compatibility guide might of had a typo
		return CompatibilityLevel3
	case VersionIsExperimental(gv.Version):
		return CompatibilityLevel3
	default:
		return CompatibilityLevelNone
	}
}

// BareMetalAPICompatibilityMapper returns the OpenShift compatibility level for API
// groups that end with `monitoring.coreos.com`, as defined by the "OpenShift 4:
// Application Compatibility GUIDE".
//
// API groups that end with suffix monitoring.coreos.com have the following mapping:.
//
// 	API Version  Compatibility Level
//	v1alpha1     Level 4
func BareMetalAPICompatibilityMapper(gv schema.GroupVersion) OpenShiftCompatibilityLevel {
	if gv.Group != "metal3.io" {
		return CompatibilityLevelNone
	}
	switch {
	case VersionIsExperimental(gv.Version):
		return CompatibilityLevel4
	default:
		return CompatibilityLevelNone
	}
}
