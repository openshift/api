# kube-apiserver Client Certificates

Used by the kube-apiserver to recognize clients using mTLS.

![PKI Graph](cert-flow.png)

- [Signing Certificate/Key Pairs](#signing-certificatekey-pairs)
    - [admin-kubeconfig-signer](#admin-kubeconfig-signer)
    - [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer)
    - [kube-control-plane-signer](#kube-control-plane-signer)
    - [per-master-debugging-signer](#per-master-debugging-signer)
    - [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer)
    - [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer)
- [Serving Certificate/Key Pairs](#serving-certificatekey-pairs)
- [Client Certificate/Key Pairs](#client-certificatekey-pairs)
    - [per-master-debugging-client](#per-master-debugging-client)
    - [system:control-plane-node-admin](#systemcontrol-plane-node-admin)
    - [kube-apiserver-to-kubelet-client](#kube-apiserver-to-kubelet-client)
    - [kube-controller-manager-client](#kube-controller-manager-client)
    - [kube-scheduler-client](#kube-scheduler-client)
    - [kube-apiserver-check-endpoints](#kube-apiserver-check-endpoints)
    - [system:serviceaccount:openshift-monitoring:prometheus-k8s](#systemserviceaccountopenshift-monitoringprometheus-k8s)
    - [system:serviceaccount:openshift-oauth-apiserver:openshift-authenticator](#systemserviceaccountopenshift-oauth-apiserveropenshift-authenticator)
- [Certificates Without Keys](#certificates-without-keys)
    - [kube-csr-signer_@1661779985](#kube-csr-signer_1661779985)
    - [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer)
    - [kubelet-signer](#kubelet-signer)
- [Certificate Authority Bundles](#certificate-authority-bundles)
    - [kube-apiserver-admin-kubeconfig-client-ca](#kube-apiserver-admin-kubeconfig-client-ca)
    - [kube-apiserver-total-client-ca](#kube-apiserver-total-client-ca)
    - [kube-apiserver-from-installer-client-ca](#kube-apiserver-from-installer-client-ca)
    - [kubelet-to-recognize-kube-apiserver-client-ca](#kubelet-to-recognize-kube-apiserver-client-ca)
    - [kube-apiserver-kcm-and-ks-client-ca](#kube-apiserver-kcm-and-ks-client-ca)
    - [kube-csr-signer_@1661779985|kubelet-signer](#kube-csr-signer_1661779985|kubelet-signer)
    - [kube-csr-signer_@1661779985|kubelet-signer|kube-csr-signer_@1661848562|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428](#kube-csr-signer_1661779985|kubelet-signer|kube-csr-signer_1661848562|openshift-kube-controller-manager-operator_csr-signer-signer1661848428)
    - [kube-csr-signer_@1661779985|kubelet-signer|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428](#kube-csr-signer_1661779985|kubelet-signer|openshift-kube-controller-manager-operator_csr-signer-signer1661848428)
    - [kube-controller-manager-csr-ca](#kube-controller-manager-csr-ca)
    - [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer)
    - [kube-apiserver-per-master-debugging-client-ca](#kube-apiserver-per-master-debugging-client-ca)
    - [kube-controller-manager-csr-signer-signer-ca](#kube-controller-manager-csr-signer-signer-ca)

## Signing Certificate/Key Pairs


### admin-kubeconfig-signer
![PKI Graph](subcert-admin-kubeconfig-signer275315413158891686.png)



| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | admin-kubeconfig-signer |
| SerialNumber | 275315413158891686 |
| Issuer CommonName | [admin-kubeconfig-signer](#admin-kubeconfig-signer) |
| Validity | 10y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### admin-kubeconfig-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |


| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/kubelet-ca.crt/kubelet-ca.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |



### kube-apiserver-to-kubelet-signer
![PKI Graph](subcert-kube-apiserver-to-kubelet-signer7567803998700114589.png)

Signer for the kube-apiserver-to-kubelet-client so kubelets can recognize the kube-apiserver.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-apiserver-to-kubelet-signer |
| SerialNumber | 7567803998700114589 |
| Issuer CommonName | [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) |
| Validity | 365d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-apiserver-to-kubelet-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | kube-apiserver-to-kubelet-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-control-plane-signer
![PKI Graph](subcert-kube-control-plane-signer4770144894411613670.png)

Signer for kube-controller-manager and kube-scheduler client certificates.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-control-plane-signer |
| SerialNumber | 4770144894411613670 |
| Issuer CommonName | [kube-control-plane-signer](#kube-control-plane-signer) |
| Validity | 365d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-control-plane-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | kube-control-plane-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### per-master-debugging-signer
![PKI Graph](subcert-openshift-kube-apiserver-operator_node-system-admin-signer16617799857870832362846967023.png)

Signer for the per-master-debugging-client.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | openshift-kube-apiserver-operator_node-system-admin-signer@1661779985 |
| SerialNumber | 7870832362846967023 |
| Issuer CommonName | [per-master-debugging-signer](#per-master-debugging-signer) |
| Validity | 365d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### per-master-debugging-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | node-system-admin-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-controller-manager-csr-signer-signer
![PKI Graph](subcert-openshift-kube-controller-manager-operator_csr-signer-signer16618484284302238975370079477.png)

Signer used by the kube-controller-manager-operator to sign signing certificates for the CSR API.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | openshift-kube-controller-manager-operator_csr-signer-signer@1661848428 |
| SerialNumber | 4302238975370079477 |
| Issuer CommonName | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) |
| Validity | 60d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-controller-manager-csr-signer-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-controller-manager-operator | csr-signer-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-controller-manager-csr-signer
![PKI Graph](subcert-kube-csr-signer_16618485622140840849824701156.png)

Signer used by the kube-controller-manager to sign CSR API requests.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-csr-signer_@1661848562 |
| SerialNumber | 2140840849824701156 |
| Issuer CommonName | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-controller-manager-csr-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-controller-manager | csr-signer |
| openshift-kube-controller-manager-operator | csr-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/secrets/csr-signer/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/secrets/csr-signer/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |


## Serving Certificate/Key Pairs

## Client Certificate/Key Pairs


### per-master-debugging-client
![PKI Graph](subcert-systemadmin9135701793828137188.png)

Client certificate (system:masters) placed on each master to allow communication to kube-apiserver for debugging.

| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:admin |
| SerialNumber | 9135701793828137188 |
| Issuer CommonName | [per-master-debugging-signer](#per-master-debugging-signer) |
| Validity | 120d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) | - system:masters |


#### per-master-debugging-client Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | node-system-admin-client |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### system:control-plane-node-admin
![PKI Graph](subcert-systemcontrol-plane-node-admin6517587000910001957.png)



| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:control-plane-node-admin |
| SerialNumber | 6517587000910001957 |
| Issuer CommonName | [kube-control-plane-signer](#kube-control-plane-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) | - system:masters |


#### system:control-plane-node-admin Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | control-plane-node-admin-client-cert-key |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/control-plane-node-admin-client-cert-key/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/control-plane-node-admin-client-cert-key/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-to-kubelet-client
![PKI Graph](subcert-systemkube-apiserver4982362661315851091.png)

Client certificate used by the kube-apiserver to authenticate to the kubelet for requests like exec and logs.

| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:kube-apiserver |
| SerialNumber | 4982362661315851091 |
| Issuer CommonName | [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) | - kube-master |


#### kube-apiserver-to-kubelet-client Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | kubelet-client |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/kubelet-client/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/kubelet-client/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-controller-manager-client
![PKI Graph](subcert-systemkube-controller-manager2778999888853223777.png)

Client certificate used by the kube-controller-manager to authenticate to the kube-apiserver.

| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:kube-controller-manager |
| SerialNumber | 2778999888853223777 |
| Issuer CommonName | [kube-control-plane-signer](#kube-control-plane-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) |  |


#### kube-controller-manager-client Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-config-managed | kube-controller-manager-client-cert-key |
| openshift-kube-controller-manager | kube-controller-manager-client-cert-key |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/secrets/kube-controller-manager-client-cert-key/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/secrets/kube-controller-manager-client-cert-key/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-scheduler-client
![PKI Graph](subcert-systemkube-scheduler4371142502064353895.png)

Client certificate used by the kube-scheduler to authenticate to the kube-apiserver.

| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:kube-scheduler |
| SerialNumber | 4371142502064353895 |
| Issuer CommonName | [kube-control-plane-signer](#kube-control-plane-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) |  |


#### kube-scheduler-client Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-config-managed | kube-scheduler-client-cert-key |
| openshift-kube-scheduler | kube-scheduler-client-cert-key |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-scheduler-certs/secrets/kube-scheduler-client-cert-key/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-certs/secrets/kube-scheduler-client-cert-key/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-check-endpoints
![PKI Graph](subcert-systemserviceaccountopenshift-kube-apiservercheck-endpoints8598399163418852136.png)

Client certificate used by the network connectivity checker of the kube-apiserver.

| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:serviceaccount:openshift-kube-apiserver:check-endpoints |
| SerialNumber | 8598399163418852136 |
| Issuer CommonName | [kube-control-plane-signer](#kube-control-plane-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) |  |


#### kube-apiserver-check-endpoints Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | check-endpoints-client-cert-key |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/check-endpoints-client-cert-key/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/check-endpoints-client-cert-key/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### system:serviceaccount:openshift-monitoring:prometheus-k8s
![PKI Graph](subcert-systemserviceaccountopenshift-monitoringprometheus-k8s218688643794807735746504063600150770281.png)



| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:serviceaccount:openshift-monitoring:prometheus-k8s |
| SerialNumber | 218688643794807735746504063600150770281 |
| Issuer CommonName | [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer) |
| Validity | 29d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | ECDSA 256 bit, P-256 curve |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) |  |


#### system:serviceaccount:openshift-monitoring:prometheus-k8s Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | metrics-client-certs |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### system:serviceaccount:openshift-oauth-apiserver:openshift-authenticator
![PKI Graph](subcert-systemserviceaccountopenshift-oauth-apiserveropenshift-authenticator322659822661055883015959461509013472419.png)



| Property | Value |
| ----------- | ----------- |
| Type | Client |
| CommonName | system:serviceaccount:openshift-oauth-apiserver:openshift-authenticator |
| SerialNumber | 322659822661055883015959461509013472419 |
| Issuer CommonName | [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer) |
| Validity | 29d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | ECDSA 256 bit, P-256 curve |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageClientAuth |
| Organizations (User Groups) |  |


#### system:serviceaccount:openshift-oauth-apiserver:openshift-authenticator Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-oauth-apiserver | openshift-authenticator-certs |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Certificates Without Keys

These certificates are present in certificate authority bundles, but do not have keys in the cluster.
This happens when the installer bootstrap clusters with a set of certificate/key pairs that are deleted during the
installation process.


### kube-csr-signer_@1661779985
![PKI Graph](subcert-kube-csr-signer_1661779985.png)



| Property | Value |
| ----------- | ----------- |
| Type |  |
| CommonName | kube-csr-signer_@1661779985 |
| SerialNumber |  |
| Issuer CommonName | None |
| Validity |  |
| Signature Algorithm |  |
| PublicKey Algorithm |   |
| Usages |  |
| ExtendedUsages |  |


#### kube-csr-signer_@1661779985 Locations
| Namespace | Secret Name |
| ----------- | ----------- |


| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kubelet-bootstrap-kubeconfig-signer
![PKI Graph](subcert-kubelet-bootstrap-kubeconfig-signer.png)



| Property | Value |
| ----------- | ----------- |
| Type |  |
| CommonName | kubelet-bootstrap-kubeconfig-signer |
| SerialNumber |  |
| Issuer CommonName | None |
| Validity |  |
| Signature Algorithm |  |
| PublicKey Algorithm |   |
| Usages |  |
| ExtendedUsages |  |


#### kubelet-bootstrap-kubeconfig-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |


| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kubelet-signer
![PKI Graph](subcert-kubelet-signer.png)



| Property | Value |
| ----------- | ----------- |
| Type |  |
| CommonName | kubelet-signer |
| SerialNumber |  |
| Issuer CommonName | None |
| Validity |  |
| Signature Algorithm |  |
| PublicKey Algorithm |   |
| Usages |  |
| ExtendedUsages |  |


#### kubelet-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |


| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Certificate Authority Bundles


### kube-apiserver-admin-kubeconfig-client-ca
![PKI Graph](subca-1227797221.png)

CA for kube-apiserver to recognize the system:master created by the installer.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [admin-kubeconfig-signer](#admin-kubeconfig-signer) | [admin-kubeconfig-signer](#admin-kubeconfig-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-admin-kubeconfig-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config | admin-kubeconfig-client-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-total-client-ca
![PKI Graph](subca-1736285941.png)

CA for kube-apiserver to recognize all known certificate based clients.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [admin-kubeconfig-signer](#admin-kubeconfig-signer) | [admin-kubeconfig-signer](#admin-kubeconfig-signer) | 10y | RSA 2048 bit |
| [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 30d | RSA 2048 bit |
| [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 60d | RSA 2048 bit |
| [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | 365d | RSA 2048 bit |
| [kube-control-plane-signer](#kube-control-plane-signer) | [kube-control-plane-signer](#kube-control-plane-signer) | 365d | RSA 2048 bit |
| [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | 10y | RSA 2048 bit |
| [per-master-debugging-signer](#per-master-debugging-signer) | [per-master-debugging-signer](#per-master-debugging-signer) | 365d | RSA 2048 bit |

#### kube-apiserver-total-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | kube-apiserver-client-ca |
| openshift-controller-manager | client-ca |
| openshift-kube-apiserver | client-ca |
| openshift-kube-controller-manager | client-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/configmaps/client-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/configmaps/client-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-from-installer-client-ca
![PKI Graph](subca-3228260836.png)

CA for the kube-apiserver to recognize clients created by the installer.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [admin-kubeconfig-signer](#admin-kubeconfig-signer) | [admin-kubeconfig-signer](#admin-kubeconfig-signer) | 10y | RSA 2048 bit |
| [kubelet-signer](#kubelet-signer) | [kubelet-signer](#kubelet-signer) | 24h | RSA 2048 bit |
| [kube-control-plane-signer](#kube-control-plane-signer) | [kube-control-plane-signer](#kube-control-plane-signer) | 365d | RSA 2048 bit |
| [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | 365d | RSA 2048 bit |
| [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-from-installer-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config | initial-kube-apiserver-server-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kubelet-to-recognize-kube-apiserver-client-ca
![PKI Graph](subca-1996114968.png)

CA for the kubelet to recognize the kube-apiserver client certificate.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | [kube-apiserver-to-kubelet-signer](#kube-apiserver-to-kubelet-signer) | 365d | RSA 2048 bit |

#### kubelet-to-recognize-kube-apiserver-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | kube-apiserver-to-kubelet-client-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-kcm-and-ks-client-ca
![PKI Graph](subca-3938274656.png)

CA for kube-apiserver to recognize the kube-controller-manager and kube-scheduler client certificates.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-control-plane-signer](#kube-control-plane-signer) | [kube-control-plane-signer](#kube-control-plane-signer) | 365d | RSA 2048 bit |

#### kube-apiserver-kcm-and-ks-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | kube-control-plane-signer-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-csr-signer_@1661779985|kubelet-signer
![PKI Graph](subca-769125443.png)



**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-csr-signer_@1661779985](#kube-csr-signer_1661779985) | [kubelet-signer](#kubelet-signer) | 23h | RSA 2048 bit |
| [kubelet-signer](#kubelet-signer) | [kubelet-signer](#kubelet-signer) | 24h | RSA 2048 bit |

#### kube-csr-signer_@1661779985|kubelet-signer Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver | kubelet-serving-ca-10 |
| openshift-kube-apiserver | kubelet-serving-ca-5 |
| openshift-kube-apiserver | kubelet-serving-ca-6 |
| openshift-kube-apiserver | kubelet-serving-ca-7 |
| openshift-kube-apiserver | kubelet-serving-ca-8 |
| openshift-kube-apiserver | kubelet-serving-ca-9 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/configmaps/kubelet-serving-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/configmaps/kubelet-serving-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-csr-signer_@1661779985|kubelet-signer|kube-csr-signer_@1661848562|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428
![PKI Graph](subca-3092952037.png)



**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-csr-signer_@1661779985](#kube-csr-signer_1661779985) | [kubelet-signer](#kubelet-signer) | 23h | RSA 2048 bit |
| [kubelet-signer](#kubelet-signer) | [kubelet-signer](#kubelet-signer) | 24h | RSA 2048 bit |
| [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 30d | RSA 2048 bit |
| [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 60d | RSA 2048 bit |

#### kube-csr-signer_@1661779985|kubelet-signer|kube-csr-signer_@1661848562|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428 Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver | kubelet-serving-ca-12 |
| openshift-kube-apiserver | kubelet-serving-ca-13 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/configmaps/kubelet-serving-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/configmaps/kubelet-serving-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-csr-signer_@1661779985|kubelet-signer|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428
![PKI Graph](subca-2348335898.png)



**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-csr-signer_@1661779985](#kube-csr-signer_1661779985) | [kubelet-signer](#kubelet-signer) | 23h | RSA 2048 bit |
| [kubelet-signer](#kubelet-signer) | [kubelet-signer](#kubelet-signer) | 24h | RSA 2048 bit |
| [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 60d | RSA 2048 bit |

#### kube-csr-signer_@1661779985|kubelet-signer|openshift-kube-controller-manager-operator_csr-signer-signer@1661848428 Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver | kubelet-serving-ca-11 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-controller-manager-csr-ca
![PKI Graph](subca-3252709373.png)

CA to recognize the CSRs (both serving and client) signed by the kube-controller-manager.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-controller-manager-csr-signer](#kube-controller-manager-csr-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 30d | RSA 2048 bit |
| [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 60d | RSA 2048 bit |

#### kube-controller-manager-csr-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | csr-controller-ca |
| openshift-config-managed | kubelet-serving-ca |
| openshift-kube-apiserver | kubelet-serving-ca |
| openshift-kube-apiserver | kubelet-serving-ca-14 |
| openshift-kube-controller-manager-operator | csr-controller-ca |
| openshift-kube-controller-manager-operator | csr-signer-ca |
| openshift-monitoring | kubelet-serving-ca-bundle |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/configmaps/kubelet-serving-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kubelet-bootstrap-kubeconfig-signer
![PKI Graph](subca-2747639387.png)



**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | [kubelet-bootstrap-kubeconfig-signer](#kubelet-bootstrap-kubeconfig-signer) | 10y | RSA 2048 bit |

#### kubelet-bootstrap-kubeconfig-signer Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | kubelet-bootstrap-kubeconfig |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-per-master-debugging-client-ca
![PKI Graph](subca-2978137741.png)

CA for kube-apiserver to recognize local system:masters rendered to each master.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [per-master-debugging-signer](#per-master-debugging-signer) | [per-master-debugging-signer](#per-master-debugging-signer) | 365d | RSA 2048 bit |

#### kube-apiserver-per-master-debugging-client-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | node-system-admin-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-controller-manager-csr-signer-signer-ca
![PKI Graph](subca-17075970.png)

CA to recognize the kube-controller-manager's signer for signing new CSR signing certificates.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | [kube-controller-manager-csr-signer-signer](#kube-controller-manager-csr-signer-signer) | 60d | RSA 2048 bit |

#### kube-controller-manager-csr-signer-signer-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-controller-manager-operator | csr-controller-signer-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



