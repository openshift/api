package utils

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

var (
	clusterProfileToShortName = map[string]string{
		"include.release.openshift.io/ibm-cloud-managed":              "Hypershift",
		"include.release.openshift.io/self-managed-high-availability": "SelfManagedHA",
		"include.release.openshift.io/single-node-developer":          "SingleNode",
	}
)

func ClusterProfileToShortName(annotation string) string {
	return clusterProfileToShortName[annotation]
}

func HasClusterProfilePreference(annotations map[string]string) bool {
	for k := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") {
			return true
		}
	}

	return false
}

func ClusterProfilesFrom(annotations map[string]string) sets.String {
	ret := sets.NewString()
	for k, v := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") && v == "true" {
			ret.Insert(k)
		}
	}
	return ret
}
