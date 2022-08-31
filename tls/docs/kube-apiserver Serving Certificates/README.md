# kube-apiserver Serving Certificates

Used by kube-apiserver clients to recognize the kube-apiserver.

![PKI Graph](cert-flow.png)

- [Signing Certificate/Key Pairs](#signing-certificatekey-pairs)
    - [ingress-operator@1661780067](#ingress-operator1661780067)
    - [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer)
    - [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer)
    - [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer)
    - [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer)
- [Serving Certificate/Key Pairs](#serving-certificatekey-pairs)
    - [*.apps.kostrows-202208291521.group-b.devcluster.openshift.com](#*.apps.kostrows-202208291521.group-b.devcluster.openshift.com)
    - [kube-apiserver-internal-load-balancer-serving](#kube-apiserver-internal-load-balancer-serving)
    - [kube-apiserver-external-load-balancer-serving](#kube-apiserver-external-load-balancer-serving)
    - [kube-apiserver-localhost-serving](#kube-apiserver-localhost-serving)
    - [kube-apiserver-service-network-serving](#kube-apiserver-service-network-serving)
    - [kube-apiserver-recovery-serving](#kube-apiserver-recovery-serving)
- [Client Certificate/Key Pairs](#client-certificatekey-pairs)
- [Certificates Without Keys](#certificates-without-keys)
- [Certificate Authority Bundles](#certificate-authority-bundles)
    - [router-wildcard-serving-ca](#router-wildcard-serving-ca)
    - [kube-apiserver-load-balancer-serving-ca](#kube-apiserver-load-balancer-serving-ca)
    - [kube-apiserver-total-serving-ca](#kube-apiserver-total-serving-ca)
    - [service-account-token-ca.crt](#service-account-token-ca.crt)
    - [kube-apiserver-localhost-serving-ca](#kube-apiserver-localhost-serving-ca)
    - [kube-apiserver-service-network-serving-ca](#kube-apiserver-service-network-serving-ca)
    - [kube-apiserver-recovery-serving-ca](#kube-apiserver-recovery-serving-ca)

## Signing Certificate/Key Pairs


### ingress-operator@1661780067
![PKI Graph](subcert-ingress-operator16617800671.png)



| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | ingress-operator@1661780067 |
| SerialNumber | 1 |
| Issuer CommonName | [ingress-operator@1661780067](#ingress-operator1661780067) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### ingress-operator@1661780067 Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-ingress-operator | router-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-load-balancer-signer
![PKI Graph](subcert-kube-apiserver-lb-signer7595672453481561385.png)

Signer used by the kube-apiserver operator to create serving certificates for the kube-apiserver via internal and external load balancers.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-apiserver-lb-signer |
| SerialNumber | 7595672453481561385 |
| Issuer CommonName | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) |
| Validity | 10y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-apiserver-load-balancer-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | loadbalancer-serving-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-8/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-8/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-5/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-5/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-7/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-7/secrets/localhost-recovery-client-token/service-ca.crt/service-ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |



### kube-apiserver-localhost-signer
![PKI Graph](subcert-kube-apiserver-localhost-signer8281335405268975279.png)

Signer used by the kube-apiserver to create serving certificates for the kube-apiserver via localhost.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-apiserver-localhost-signer |
| SerialNumber | 8281335405268975279 |
| Issuer CommonName | [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) |
| Validity | 10y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-apiserver-localhost-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | localhost-serving-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-4/secrets/localhost-recovery-client-token/ca.crt/ca.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
|  |  |  |  |  |



### kube-apiserver-service-network-signer
![PKI Graph](subcert-kube-apiserver-service-network-signer8689246760707875037.png)

Signer used by the kube-apiserver to create serving certificates for the kube-apiserver via the service network.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | kube-apiserver-service-network-signer |
| SerialNumber | 8689246760707875037 |
| Issuer CommonName | [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) |
| Validity | 10y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-apiserver-service-network-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | service-network-serving-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-recovery-signer
![PKI Graph](subcert-openshift-kube-apiserver-operator_localhost-recovery-serving-signer16617799877484385607974147373.png)

Signer used by the kube-apiserver to create serving certificates for the kube-apiserver via the localhost recovery SNI ServerName

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | openshift-kube-apiserver-operator_localhost-recovery-serving-signer@1661779987 |
| SerialNumber | 7484385607974147373 |
| Issuer CommonName | [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) |
| Validity | 10y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### kube-apiserver-recovery-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | localhost-recovery-serving-signer |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Serving Certificate/Key Pairs


### *.apps.kostrows-202208291521.group-b.devcluster.openshift.com
![PKI Graph](subcert-*.apps.kostrows-202208291521.group-b.devcluster.openshift.com5596896335820314549.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | *.apps.kostrows-202208291521.group-b.devcluster.openshift.com |
| SerialNumber | 5596896335820314549 |
| Issuer CommonName | [ingress-operator@1661780067](#ingress-operator1661780067) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - *.apps.kostrows-202208291521.group-b.devcluster.openshift.com |
| IP Addresses |  |


#### *.apps.kostrows-202208291521.group-b.devcluster.openshift.com Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-ingress | router-certs-default |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-internal-load-balancer-serving
![PKI Graph](subcert-api-int.kostrows-202208291521.group-b.devcluster.openshift.com2451083771403934178.png)

Serving certificate used by the kube-apiserver to terminate requests via the internal load balancer.

| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | api-int.kostrows-202208291521.group-b.devcluster.openshift.com |
| SerialNumber | 2451083771403934178 |
| Issuer CommonName | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - api-int.kostrows-202208291521.group-b.devcluster.openshift.com |
| IP Addresses |  |


#### kube-apiserver-internal-load-balancer-serving Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | internal-loadbalancer-serving-certkey |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/internal-loadbalancer-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/internal-loadbalancer-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-external-load-balancer-serving
![PKI Graph](subcert-api.kostrows-202208291521.group-b.devcluster.openshift.com3223170914555462336.png)

Serving certificate used by the kube-apiserver to terminate requests via the external load balancer.

| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | api.kostrows-202208291521.group-b.devcluster.openshift.com |
| SerialNumber | 3223170914555462336 |
| Issuer CommonName | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - api.kostrows-202208291521.group-b.devcluster.openshift.com |
| IP Addresses |  |


#### kube-apiserver-external-load-balancer-serving Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | external-loadbalancer-serving-certkey |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/external-loadbalancer-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/external-loadbalancer-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-localhost-serving
![PKI Graph](subcert-127.0.0.13009370111951769082.png)

Serving certificate used by the kube-apiserver to terminate requests via localhost.

| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | 127.0.0.1 |
| SerialNumber | 3009370111951769082 |
| Issuer CommonName | [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - localhost<br/>- 127.0.0.1 |
| IP Addresses | - 127.0.0.1 |


#### kube-apiserver-localhost-serving Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | localhost-serving-cert-certkey |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/localhost-serving-cert-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/localhost-serving-cert-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-service-network-serving
![PKI Graph](subcert-172.30.0.15653002920292026104.png)

Serving certificate used by the kube-apiserver to terminate requests via the service network.

| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | 172.30.0.1 |
| SerialNumber | 5653002920292026104 |
| Issuer CommonName | [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) |
| Validity | 30d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - kubernetes<br/>- kubernetes.default<br/>- kubernetes.default.svc<br/>- kubernetes.default.svc.cluster.local<br/>- openshift<br/>- openshift.default<br/>- openshift.default.svc<br/>- openshift.default.svc.cluster.local<br/>- 172.30.0.1 |
| IP Addresses | - 172.30.0.1 |


#### kube-apiserver-service-network-serving Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | service-network-serving-certkey |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/service-network-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/service-network-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-recovery-serving
![PKI Graph](subcert-localhost-recovery3845049330654169764.png)

Serving certificate used by the kube-apiserver to terminate requests via the localhost recovery SNI ServerName.

| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | localhost-recovery |
| SerialNumber | 3845049330654169764 |
| Issuer CommonName | [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) |
| Validity | 9y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - localhost-recovery |
| IP Addresses |  |


#### kube-apiserver-recovery-serving Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver | localhost-recovery-serving-certkey |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-10 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-11 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-12 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-13 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-14 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-5 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-6 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-7 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-8 |
| openshift-kube-apiserver | localhost-recovery-serving-certkey-9 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/secrets/localhost-recovery-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/secrets/localhost-recovery-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/secrets/localhost-recovery-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/secrets/localhost-recovery-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/secrets/localhost-recovery-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/secrets/localhost-recovery-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/secrets/localhost-recovery-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/secrets/localhost-recovery-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/secrets/localhost-recovery-serving-certkey/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/secrets/localhost-recovery-serving-certkey/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |


## Client Certificate/Key Pairs

## Certificates Without Keys

These certificates are present in certificate authority bundles, but do not have keys in the cluster.
This happens when the installer bootstrap clusters with a set of certificate/key pairs that are deleted during the
installation process.

## Certificate Authority Bundles


### router-wildcard-serving-ca
![PKI Graph](subca-3485746206.png)

REVIEW: CA for recognizing the default router wildcard serving certificate.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [*.apps.kostrows-202208291521.group-b.devcluster.openshift.com](#*.apps.kostrows-202208291521.group-b.devcluster.openshift.com) | [ingress-operator@1661780067](#ingress-operator1661780067) | 2y | RSA 2048 bit |
| [ingress-operator@1661780067](#ingress-operator1661780067) | [ingress-operator@1661780067](#ingress-operator1661780067) | 2y | RSA 2048 bit |

#### router-wildcard-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | default-ingress-cert |
| openshift-config-managed | oauth-serving-cert |
| openshift-console | default-ingress-cert |
| openshift-console | oauth-serving-cert |
| openshift-kube-scheduler | serviceaccount-ca-3 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-load-balancer-serving-ca
![PKI Graph](subca-3666101852.png)

CA for recognizing the kube-apiserver when connecting via the internal or external load balancers.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-load-balancer-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | loadbalancer-serving-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-total-serving-ca
![PKI Graph](subca-703958823.png)

CA for recognizing the kube-apiserver when connecting via any means.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-total-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | kube-apiserver-server-ca |
| openshift-kube-apiserver | kube-apiserver-server-ca |
| openshift-kube-apiserver | kube-apiserver-server-ca-10 |
| openshift-kube-apiserver | kube-apiserver-server-ca-11 |
| openshift-kube-apiserver | kube-apiserver-server-ca-12 |
| openshift-kube-apiserver | kube-apiserver-server-ca-13 |
| openshift-kube-apiserver | kube-apiserver-server-ca-14 |
| openshift-kube-apiserver | kube-apiserver-server-ca-5 |
| openshift-kube-apiserver | kube-apiserver-server-ca-6 |
| openshift-kube-apiserver | kube-apiserver-server-ca-7 |
| openshift-kube-apiserver | kube-apiserver-server-ca-8 |
| openshift-kube-apiserver | kube-apiserver-server-ca-9 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-10/configmaps/kube-apiserver-server-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-12/configmaps/kube-apiserver-server-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-13/configmaps/kube-apiserver-server-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-14/configmaps/kube-apiserver-server-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-apiserver-pod-9/configmaps/kube-apiserver-server-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### service-account-token-ca.crt
![PKI Graph](subca-2123940722.png)

CA for recognizing kube-apiserver.  This is injected into each service account token secret at ca.crt.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | [kube-apiserver-load-balancer-signer](#kube-apiserver-load-balancer-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | 10y | RSA 2048 bit |
| [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | 10y | RSA 2048 bit |
| [*.apps.kostrows-202208291521.group-b.devcluster.openshift.com](#*.apps.kostrows-202208291521.group-b.devcluster.openshift.com) | [ingress-operator@1661780067](#ingress-operator1661780067) | 2y | RSA 2048 bit |
| [ingress-operator@1661780067](#ingress-operator1661780067) | [ingress-operator@1661780067](#ingress-operator1661780067) | 2y | RSA 2048 bit |

#### service-account-token-ca.crt Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-controller-manager | serviceaccount-ca |
| openshift-kube-controller-manager | serviceaccount-ca-4 |
| openshift-kube-controller-manager | serviceaccount-ca-5 |
| openshift-kube-controller-manager | serviceaccount-ca-6 |
| openshift-kube-controller-manager | serviceaccount-ca-7 |
| openshift-kube-controller-manager | serviceaccount-ca-8 |
| openshift-kube-scheduler | serviceaccount-ca |
| openshift-kube-scheduler | serviceaccount-ca-4 |
| openshift-kube-scheduler | serviceaccount-ca-5 |
| openshift-kube-scheduler | serviceaccount-ca-6 |
| openshift-kube-scheduler | serviceaccount-ca-7 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-4/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-8/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-5/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-7/configmaps/serviceaccount-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-apiserver-localhost-serving-ca
![PKI Graph](subca-537014043.png)

CA for recognizing the kube-apiserver when connecting via localhost.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | [kube-apiserver-localhost-signer](#kube-apiserver-localhost-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-localhost-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | localhost-serving-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-service-network-serving-ca
![PKI Graph](subca-1272615924.png)

CA for recognizing the kube-apiserver when connecting via the service network (kuberentes.default.svc).

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | [kube-apiserver-service-network-signer](#kube-apiserver-service-network-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-service-network-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | service-network-serving-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-apiserver-recovery-serving-ca
![PKI Graph](subca-586916855.png)

CA for recognizing the kube-apiserver when connecting via the localhost recovery SNI ServerName.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | [kube-apiserver-recovery-signer](#kube-apiserver-recovery-signer) | 10y | RSA 2048 bit |

#### kube-apiserver-recovery-serving-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | localhost-recovery-serving-ca |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



