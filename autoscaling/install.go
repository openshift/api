package autoscaling

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	autoscalingv1 "github.com/openshift/api/autoscaling/v1"
	autoscalingv1beta1 "github.com/openshift/api/autoscaling/v1beta1"
)

const (
	GroupName = "autoscaling.openshift.io"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(
		autoscalingv1.Install,
		autoscalingv1beta1.Install,
	)
	// Install is a function which adds every version of this group to a scheme
	Install = schemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

func Kind(kind string) schema.GroupKind {
	return schema.GroupKind{Group: GroupName, Kind: kind}
}
