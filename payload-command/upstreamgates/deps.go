// This file is used to register all kubernetes feature gates via all kubernetes packages that specify
// a "features" package with an init function that will register the gates with the default mutable
//
package upstreamgates

import (
	_ "k8s.io/kubernetes/pkg/features"
	_ "github.com/openshift/api/payload-command/upstreamgates/overrides"
)
