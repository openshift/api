package main

import (
	"regexp"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type OpenShiftCompatibilityLevel string

const (
	CompatibilityLevelNone OpenShiftCompatibilityLevel = "None"
	CompatibilityLevel1    OpenShiftCompatibilityLevel = "Level1"
	CompatibilityLevel2    OpenShiftCompatibilityLevel = "Level2"
	CompatibilityLevel3    OpenShiftCompatibilityLevel = "Level3"
	CompatibilityLevel4    OpenShiftCompatibilityLevel = "Level4"
)

type OpenShiftCompatibilityStatement struct {
	GVK           schema.GroupVersionKind
	compatibility OpenShiftCompatibilityLevel
}

type CompatibilityMapper func(gv schema.GroupVersion) OpenShiftCompatibilityLevel

var CompatibilityMappers = []CompatibilityMapper{
	OpenShiftAPICompatibilityMapper,
	KubernetesAPICompatibilityMapper,
	BareMetalAPICompatibilityMapper,
	CoreOSMonitoringAPICompatibilityMapper,
	CoreOSOperatorsAPICompatibilityMapper,
}

func NewOpenShiftCompatibilityStatement(gvk schema.GroupVersionKind) OpenShiftCompatibilityStatement {
	for _, f := range CompatibilityMappers {
		compatibility := f(gvk.GroupVersion())
		if compatibility != CompatibilityLevelNone {
			return OpenShiftCompatibilityStatement{GVK: gvk, compatibility: compatibility}
		}
	}
	return OpenShiftCompatibilityStatement{GVK: gvk, compatibility: CompatibilityLevelNone}
}

func VersionIsGenerallyAvailable(version string) bool {
	return regexp.MustCompile(`^v\d*$`).MatchString(version)
}

func VersionIsPrerelease(version string) bool {
	return regexp.MustCompile(`^v\d*beta\d*$`).MatchString(version)
}

func VersionIsExperimental(version string) bool {
	return regexp.MustCompile(`^v\d*alpha\d*$`).MatchString(version)
}
