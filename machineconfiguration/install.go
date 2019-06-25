package machineconfiguration

import (
	machineconfigurationv1 "github.com/openshift/api/machineconfiguration/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "machineconfiguration.openshift.io"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(machineconfigurationv1.Install)
	// Install is a function which adds every version of this group to a scheme
	Install = schemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

func Kind(kind string) schema.GroupKind {
	return schema.GroupKind{Group: GroupName, Kind: kind}
}
