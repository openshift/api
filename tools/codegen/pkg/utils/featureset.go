package utils

import (
	"strings"
)

var (
	clusterProfileToShortName = map[string]string{
		"include.release.openshift.io/ibm-cloud-managed":              "Hypershift",
		"include.release.openshift.io/self-managed-high-availability": "SelfManagedHA",
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
