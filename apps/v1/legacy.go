package v1

import (
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	LegacyGroupName               = ""
	LegacyGroupVersion            = schema.GroupVersion{Group: LegacyGroupName, Version: "v1"}
	LegacySchemeBuilder           = runtime.NewSchemeBuilder(addLegacyKnownTypes, corev1.AddToScheme, extensionsv1beta1.AddToScheme)
	DeprecatedInstallWithoutGroup = LegacySchemeBuilder.AddToScheme

	LegacySchemeGroupVersion = LegacyGroupVersion
)

func addLegacyKnownTypes(scheme *runtime.Scheme) error {
	types := []runtime.Object{
		&DeploymentConfig{},
		&DeploymentConfigList{},
		&DeploymentConfigRollback{},
		&DeploymentRequest{},
		&DeploymentLog{},
		&DeploymentLogOptions{},
		&extensionsv1beta1.Scale{},
	}
	scheme.AddKnownTypes(LegacyGroupVersion, types...)
	return nil
}
