package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupName     = "machineconfiguration.openshift.io"
	GroupVersion  = schema.GroupVersion{Group: GroupName, Version: "v1"}
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// Install is a function which adds this version to a scheme
	Install = schemeBuilder.AddToScheme
)

// addKnownTypes adds types to API group
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&ContainerRuntimeConfig{},
		&ContainerRuntimeConfigList{},
		&ControllerConfig{},
		&ControllerConfigList{},
		&KubeletConfig{},
		&KubeletConfigList{},
		&MachineConfigPool{},
		&MachineConfigPoolList{},
	)

	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
