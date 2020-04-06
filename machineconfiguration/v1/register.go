package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupName is the group name of this api
	GroupName = "machineconfiguration.openshift.io"
	// GroupVersion is the version of this api group
	GroupVersion  = schema.GroupVersion{Group: GroupName, Version: "v1"}
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// Install is a function which adds this version to a scheme
	Install = schemeBuilder.AddToScheme

	// SchemeGroupVersion generated code relies on this name
	// Deprecated
	SchemeGroupVersion = GroupVersion
	// AddToScheme exists solely to keep the old generators creating valid code
	// DEPRECATED
	AddToScheme = Install
)

// Resource generated code relies on this being here, but it logically belongs to the group
// DEPRECATED
func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

// addKnownTypes adds types to API group
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&ContainerRuntimeConfig{},
		&ContainerRuntimeConfigList{},
		&KubeletConfig{},
		&KubeletConfigList{},
		&MachineConfig{},
		&MachineConfigList{},
		&MachineConfigPool{},
		&MachineConfigPoolList{},
	)

	metav1.AddToGroupVersion(scheme, GroupVersion)

	return nil
}
