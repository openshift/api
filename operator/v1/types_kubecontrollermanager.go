package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeControllerManager provides information to configure an operator to manage kube-controller-manager.
type KubeControllerManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   KubeControllerManagerSpec   `json:"spec"`
	Status KubeControllerManagerStatus `json:"status"`
}

type KubeControllerManagerSpec struct {
	StaticPodOperatorSpec `json:",inline"`

	// forceRedeploymentReason can be used to force the redeployment of the kube-controller-manager by providing a unique string.
	// This provides a mechanism to kick a previously failed deployment and provide a reason why you think it will work
	// this time instead of failing again on the same config.
	ForceRedeploymentReason string `json:"forceRedeploymentReason"`

	// unsupportedCertificateRotation contains information for overriding the default rotation times on particular certificates.
	// You should only change settings here if you are very sure you know what you're doing. Changing defaults can
	// cause cluster instability
	CertificateRotation *KubeControllerManagerCertificateRotationOverrides `json:"unsupportedCertificateRotation,omitempty"`
}

type KubeControllerManagerCertificateRotationOverrides struct {
	// csrSignerCert is how the CSR signing controller signs certificates
	CSRSignerCert *CertificateRotation `json:"csrSignerCert,omitempty"`
}

type KubeControllerManagerStatus struct {
	StaticPodOperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeControllerManagerList is a collection of items
type KubeControllerManagerList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items contains the items
	Items []KubeControllerManager `json:"items"`
}
