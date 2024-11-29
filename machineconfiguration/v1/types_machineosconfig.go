package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=machineosconfigs,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/1773
// +openshift:enable:FeatureGate=OnClusterBuild
// +openshift:file-pattern=cvoRunLevel=0000_80,operatorName=machine-config,operatorOrdering=01
// +kubebuilder:metadata:labels=openshift.io/operator-managed=

// MachineOSConfig describes the configuration for a build process managed by the MCO
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineOSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the configuration of the machineosconfig
	// +kubebuilder:validation:Required
	Spec MachineOSConfigSpec `json:"spec"`

	// status describes the status of the machineosconfig
	// +optional
	Status *MachineOSConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineOSConfigList describes all configurations for image builds on the system
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineOSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// items contains a collection of MachineOSConfig resources.
	Items []MachineOSConfig `json:"items"`
}

// MachineOSConfigSpec describes user-configurable options as well as information about a build process.
type MachineOSConfigSpec struct {
	// machineConfigPool is the pool which the build is for
	// +kubebuilder:validation:Required
	MachineConfigPool MachineConfigPoolReference `json:"machineConfigPool"`
	// buildInputs is where user input options for the build live
	// +kubebuilder:validation:Required
	BuildInputs BuildInputs `json:"buildInputs"`
	// buildOutputs holds all information needed to handle booting the image after a build
	// This currently contains a currentImagePullSecret field, which should be provided if the final pull secret used to pull the image to nodes from the registry
	// is different than the one used for pushing the image to the registry during the build.
	// +optional
	BuildOutputs *BuildOutputs `json:"buildOutputs,omitempty"`
}

// MachineOSConfigStatus describes the status this config object and relates it to the builds associated with this MachineOSConfig
type MachineOSConfigStatus struct {
	// observedGeneration represents the generation of the MachineOSConfig object observed by the Machine Config Operator's build controller.
	// +kubebuilder:validation:XValidation:rule="self >= oldSelf || (self == 0 && oldSelf > 0)", message="observedGeneration must not move backwards except to zero"
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// currentImagePullspec is the fully qualified image pull spec used by the MCO to pull down the new OSImage. This must include sha256.
	// The format of the image pullspec is:
	// host[:port][/namespace]/name@sha256:<digest>
	// The digest must be 64 characters long, and consist only of lowercase hexadecimal characters, a-f and 0-9.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=447
	// +kubebuilder:validation:XValidation:rule=`(self.split('@').size() == 2 && self.split('@')[1].matches('^sha256:[a-f0-9]{64}$'))`,message="the OCI Image reference must end with a valid '@sha256:<digest>' suffix, where '<digest>' is 64 characters long"
	// +kubebuilder:validation:XValidation:rule=`(self.split('@')[0].matches('^([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]+(:[0-9]{2,5})?/([a-zA-Z0-9-_]{0,61}/)?[a-zA-Z0-9-_.]*?$'))`,message="the OCI Image name should follow the host[:port][/namespace]/name format, resembling a valid URL without the scheme"
	// +optional
	CurrentImagePullspec string `json:"currentImagePullspec,omitempty"`
	// machineOSBuild is a reference to the MachineOSBuild object for this MachineOSConfig, which contains the status for the image build
	// +optional
	MachineOSBuild *ObjectReference `json:"machineOSBuild,omitempty"`
}

// BuildInputs holds all of the information needed to trigger a build
type BuildInputs struct {
	// baseOSExtensionsImagePullspec is the base Extensions image used in the build process
	// The MachineOSConfig object will use the in cluster image registry configuration.
	// If you wish to use a mirror or any other settings specific to registries.conf, please specify those in the cluster wide registries.conf.
	// The format of the image pullspec is:
	// host[:port][/namespace]/name@sha256:<digest>
	// The digest must be 64 characters long, and consist only of lowercase hexadecimal characters, a-f and 0-9.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=447
	// +kubebuilder:validation:XValidation:rule=`(self.split('@').size() == 2 && self.split('@')[1].matches('^sha256:[a-f0-9]{64}$'))`,message="the OCI Image reference must end with a valid '@sha256:<digest>' suffix, where '<digest>' is 64 characters long"
	// +kubebuilder:validation:XValidation:rule=`(self.split('@')[0].matches('^([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]+(:[0-9]{2,5})?/([a-zA-Z0-9-_]{0,61}/)?[a-zA-Z0-9-_.]*?$'))`,message="the OCI Image name should follow the host[:port][/namespace]/name format, resembling a valid URL without the scheme"
	// +optional
	BaseOSExtensionsImagePullspec string `json:"baseOSExtensionsImagePullspec,omitempty"`
	// baseOSImagePullspec is the base OSImage we use to build our custom image.
	// The MachineOSConfig object will use the in cluster image registry configuration.
	// If you wish to use a mirror or any other settings specific to registries.conf, please specify those in the cluster wide registries.conf.
	// The format of the image pullspec is:
	// host[:port][/namespace]/name@sha256:<digest>
	// The digest must be 64 characters long, and consist only of lowercase hexadecimal characters, a-f and 0-9.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=447
	// +kubebuilder:validation:XValidation:rule=`(self.split('@').size() == 2 && self.split('@')[1].matches('^sha256:[a-f0-9]{64}$'))`,message="the OCI Image reference must end with a valid '@sha256:<digest>' suffix, where '<digest>' is 64 characters long"
	// +kubebuilder:validation:XValidation:rule=`(self.split('@')[0].matches('^([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]+(:[0-9]{2,5})?/([a-zA-Z0-9-_]{0,61}/)?[a-zA-Z0-9-_.]*?$'))`,message="the OCI Image name should follow the host[:port][/namespace]/name format, resembling a valid URL without the scheme"
	// +optional
	BaseOSImagePullspec string `json:"baseOSImagePullspec,omitempty"`
	// baseImagePullSecret is the secret used to pull the base image.
	// Must live in the openshift-machine-config-operator namespace if provided.
	// Defaults to using the cluster-wide pull secret if not specified. This is provided during install time of the cluster, and lives in the openshift-config namespace as a secret.
	// +optional
	BaseImagePullSecret *ImageSecretObjectReference `json:"baseImagePullSecret,omitempty"`
	// machineOSImageBuilder describes which image builder will be used in each build triggered by this MachineOSConfig.
	// Currently supported type(s): JobImageBuilder
	// +kubebuilder:validation:Required
	ImageBuilder MachineOSImageBuilder `json:"imageBuilder"`
	// renderedImagePushSecret is the secret used to connect to a user registry.
	// The final image push and pull secrets should be separate and assume the principal of least privilege.
	// The push secret with write privilege is only required to be present on the node hosting the MachineConfigController pod.
	// The pull secret with read only privileges is required on all nodes.
	// By separating the two secrets, the risk of write credentials becoming compromised is reduced.
	// +kubebuilder:validation:Required
	RenderedImagePushSecret ImageSecretObjectReference `json:"renderedImagePushSecret"`
	// renderedImagePushSpec describes the location of the final image.
	// The MachineOSConfig object will use the in cluster image registry configuration.
	// If you wish to use a mirror or any other settings specific to registries.conf, please specify those in the cluster wide registries.conf via the cluster image.config, ImageContentSourcePolicies, ImageDigestMirrorSet, or ImageTagMirrorSet objects.
	// The format of the image pushspec is:
	// host[:port][/namespace]/name:<tag> or svc_name.namespace.svc[:port]/repository/name:<tag>
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=447
	// +kubebuilder:validation:XValidation:rule=`self.matches('^([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]+(:[0-9]{2,5})?(/[a-zA-Z0-9-_]{1,61})*/[a-zA-Z0-9-_.]+:[a-zA-Z0-9._-]+$') || self.matches('^[^.]+\\.[^.]+\\.svc:\\d+\\/[^\\/]+\\/[^\\/]+:[^\\/]+$')`,message="the OCI Image name should follow the host[:port][/namespace]/name format, resembling a valid URL without the scheme. Or it must be a valid .svc followed by a port, repository, image name, and tag."
	// +kubebuilder:validation:Required
	RenderedImagePushSpec string `json:"renderedImagePushSpec"`
	// releaseVersion is an Openshift release version which the base OS image is associated with.
	// This field is populated from the machine-config-osimageurl configmap in the openshift-machine-config-operator namespace.
	// It will come in the format: 4.16.0-0.nightly-2024-04-03-065948 or any valid release. The MachineOSBuilder populates this field and validates that this is a valid stream.
	// This is used as a label in the Containerfile that builds the OS image.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +optional
	ReleaseVersion string `json:"releaseVersion,omitempty"`
	// containerFile describes the custom data the user has specified to build into the image.
	// This is also commonly called a Dockerfile and you can treat it as such. The content is the content of your Dockerfile.
	// See https://github.com/containers/common/blob/main/docs/Containerfile.5.md for the spec reference.
	// you can specify up to 7 containerFiles
	// +patchMergeKey=containerfileArch
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=containerfileArch
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=7
	// +optional
	Containerfile []MachineOSContainerfile `json:"containerFile" patchStrategy:"merge" patchMergeKey:"containerfileArch"`
}

// BuildOutputs holds all information needed to handle booting the image after a build
type BuildOutputs struct {
	// currentImagePullSecret is the secret used to pull the final produced image.
	// Must live in the openshift-machine-config-operator namespace,
	// the final image push and pull secrets should be separate for security concerns. If the final image push secret is somehow exfiltrated,
	// that gives someone the power to push images to the image repository. By comparison, if the final image pull secret gets exfiltrated,
	// that only gives someone to pull images from the image repository. It's basically the principle of least permissions.
	// This pull secret will be used on all nodes in the pool. These nodes will need to pull the final OS image and boot into it using rpm-ostree or bootc.
	// +optional
	CurrentImagePullSecret *ImageSecretObjectReference `json:"currentImagePullSecret,omitempty"`
}

type MachineOSImageBuilder struct {
	// imageBuilderType specifies the backend to be used to build the image.
	// +kubebuilder:validation:Enum:=JobImageBuilder
	// Valid options are: JobImageBuilder
	// +required
	ImageBuilderType MachineOSImageBuilderType `json:"imageBuilderType"`
}

// MachineOSContainerfile contains all custom content the user wants built into the image
type MachineOSContainerfile struct {
	// containerfileArch describes the architecture this containerfile is to be built for.
	// This arch is optional. If the user does not specify an architecture, it is assumed
	// that the content can be applied to all architectures, or in a single arch cluster: the only architecture.
	// +kubebuilder:validation:Enum:=ARM64;AMD64;PPC64LE;S390X;AArch64;x86_64;NoArch
	// +kubebuilder:default:=NoArch
	// +optional
	ContainerfileArch ContainerfileArch `json:"containerfileArch"`
	// content is an embedded Containerfile/Dockerfile that defines the contents to be built into your image.
	// See https://github.com/containers/common/blob/main/docs/Containerfile.5.md for the spec reference.
	// for example, this would add the tree package to your hosts:
	//   FROM configs AS final
	//   RUN rpm-ostree install tree && \
	//     ostree container commit
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=4096
	Content string `json:"content"`
}

// +enum
type ContainerfileArch string

const (
	// describes the arm64 architecture
	Arm64 ContainerfileArch = "ARM64"
	// describes the amd64 architecture
	Amd64 ContainerfileArch = "AMD64"
	// describes the ppc64le architecture
	Ppc ContainerfileArch = "PPC64LE"
	// describes the s390x architecture
	S390 ContainerfileArch = "S390X"
	// describes the aarch64 architecture
	Aarch64 ContainerfileArch = "AArch64"
	// describes the fx86_64 architecture
	X86_64 ContainerfileArch = "x86_64"
	// describes a containerfile that can be applied to any arch
	NoArch ContainerfileArch = "NoArch"
)

// Refers to the name of a MachineConfigPool (e.g., "worker", "infra", etc.):
// the MachineOSBuilder pod validates that the user has provided a valid pool
type MachineConfigPoolReference struct {
	// name of the MachineConfigPool object.
	// Must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// Refers to the name of an image registry push/pull secret needed in the build process.
type ImageSecretObjectReference struct {
	// name is the name of the secret used to push or pull this MachineOSConfig object.
	// Must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// This secret must be in the openshift-machine-config-operator namespace.
	// +kubebuilder:validation:MaxLength:=253
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// +enum
type MachineOSImageBuilderType string

const (
	// describes that the machine-os-builder will use a Job to spin up a custom pod builder that uses buildah
	JobBuilder MachineOSImageBuilderType = "JobImageBuilder"
)
