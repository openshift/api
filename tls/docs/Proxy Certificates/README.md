# Proxy Certificates

Used by the OpenShift platform to recognize the proxy.  Other usages are side-effects which work by accident and not by principled design.

![PKI Graph](cert-flow.png)

- [Signing Certificate/Key Pairs](#signing-certificatekey-pairs)
    - [](#)
- [Serving Certificate/Key Pairs](#serving-certificatekey-pairs)
    - [](#)
- [Client Certificate/Key Pairs](#client-certificatekey-pairs)
    - [](#)
- [Certificates Without Keys](#certificates-without-keys)
- [Certificate Authority Bundles](#certificate-authority-bundles)
    - [proxy-ca](#proxy-ca)

## Signing Certificate/Key Pairs


### 
![PKI Graph](subcert-2885142665935399532.png)



| Property | Value |
| ----------- | ----------- |
| Type | Signer,Serving,Client |
| CommonName |  |
| SerialNumber | 2885142665935399532 |
| Issuer CommonName | [](#) |
| Validity | 728d |
| Signature Algorithm | ECDSA-SHA256 |
| PublicKey Algorithm | ECDSA 256 bit, P-256 curve |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageCertSign |
| ExtendedUsages | - ExtKeyUsageClientAuth<br/>- ExtKeyUsageServerAuth |
| Organizations (User Groups) | - Red Hat, Inc. |
| DNS Names | - packageserver-service.openshift-operator-lifecycle-manager<br/>- packageserver-service.openshift-operator-lifecycle-manager.svc |
| IP Addresses |  |


####  Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-operator-lifecycle-manager | packageserver-service-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Serving Certificate/Key Pairs


### 
![PKI Graph](subcert-2885142665935399532.png)



| Property | Value |
| ----------- | ----------- |
| Type | Signer,Serving,Client |
| CommonName |  |
| SerialNumber | 2885142665935399532 |
| Issuer CommonName | [](#) |
| Validity | 728d |
| Signature Algorithm | ECDSA-SHA256 |
| PublicKey Algorithm | ECDSA 256 bit, P-256 curve |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageCertSign |
| ExtendedUsages | - ExtKeyUsageClientAuth<br/>- ExtKeyUsageServerAuth |
| Organizations (User Groups) | - Red Hat, Inc. |
| DNS Names | - packageserver-service.openshift-operator-lifecycle-manager<br/>- packageserver-service.openshift-operator-lifecycle-manager.svc |
| IP Addresses |  |


####  Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-operator-lifecycle-manager | packageserver-service-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Client Certificate/Key Pairs


### 
![PKI Graph](subcert-2885142665935399532.png)



| Property | Value |
| ----------- | ----------- |
| Type | Signer,Serving,Client |
| CommonName |  |
| SerialNumber | 2885142665935399532 |
| Issuer CommonName | [](#) |
| Validity | 728d |
| Signature Algorithm | ECDSA-SHA256 |
| PublicKey Algorithm | ECDSA 256 bit, P-256 curve |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageCertSign |
| ExtendedUsages | - ExtKeyUsageClientAuth<br/>- ExtKeyUsageServerAuth |
| Organizations (User Groups) | - Red Hat, Inc. |
| DNS Names | - packageserver-service.openshift-operator-lifecycle-manager<br/>- packageserver-service.openshift-operator-lifecycle-manager.svc |
| IP Addresses |  |


####  Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-operator-lifecycle-manager | packageserver-service-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Certificates Without Keys

These certificates are present in certificate authority bundles, but do not have keys in the cluster.
This happens when the installer bootstrap clusters with a set of certificate/key pairs that are deleted during the
installation process.

## Certificate Authority Bundles


### proxy-ca
![PKI Graph](subca-772742968.png)

CA used to recognize proxy servers.  By default this will contain standard root CAs on the cluster-network-operator pod.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |


#### proxy-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-apiserver | trusted-ca-bundle |
| openshift-apiserver-operator | trusted-ca-bundle |
| openshift-authentication | v4-0-config-system-trusted-ca-bundle |
| openshift-authentication-operator | trusted-ca-bundle |
| openshift-cloud-credential-operator | cco-trusted-ca |
| openshift-cluster-node-tuning-operator | trusted-ca |
| openshift-config-managed | trusted-ca-bundle |
| openshift-console | trusted-ca-bundle |
| openshift-controller-manager | openshift-global-ca |
| openshift-image-registry | trusted-ca |
| openshift-ingress-operator | trusted-ca |
| openshift-insights | trusted-ca-bundle |
| openshift-kube-apiserver | trusted-ca-bundle |
| openshift-kube-controller-manager | trusted-ca-bundle |
| openshift-machine-api | cbo-trusted-ca |
| openshift-machine-api | mao-trusted-ca |
| openshift-marketplace | marketplace-trusted-ca |
| openshift-monitoring | alertmanager-trusted-ca-bundle |
| openshift-monitoring | alertmanager-trusted-ca-bundle-d34s91lhv300e |
| openshift-monitoring | grafana-trusted-ca-bundle |
| openshift-monitoring | grafana-trusted-ca-bundle-d34s91lhv300e |
| openshift-monitoring | prometheus-trusted-ca-bundle |
| openshift-monitoring | prometheus-trusted-ca-bundle-d34s91lhv300e |
| openshift-monitoring | telemeter-trusted-ca-bundle |
| openshift-monitoring | telemeter-trusted-ca-bundle-d34s91lhv300e |
| openshift-monitoring | thanos-querier-trusted-ca-bundle |
| openshift-monitoring | thanos-querier-trusted-ca-bundle-d34s91lhv300e |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-apiserver-certs/configmaps/trusted-ca-bundle/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-certs/configmaps/trusted-ca-bundle/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |


