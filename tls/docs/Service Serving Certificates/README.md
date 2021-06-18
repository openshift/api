# Service Serving Certificates

Used to secure inter-service communication on the local cluster.

![PKI Graph](cert-flow.png)

- [Signing Certificate/Key Pairs](#signing-certificatekey-pairs)
    - [service-serving-signer](#service-serving-signer)
- [Serving Certificate/Key Pairs](#serving-certificatekey-pairs)
    - [alertmanager-main.openshift-monitoring.svc](#alertmanager-main.openshift-monitoring.svc)
    - [api.openshift-apiserver.svc](#api.openshift-apiserver.svc)
    - [api.openshift-oauth-apiserver.svc](#api.openshift-oauth-apiserver.svc)
    - [catalog-operator-metrics.openshift-operator-lifecycle-manager.svc](#catalog-operator-metrics.openshift-operator-lifecycle-manager.svc)
    - [cco-metrics.openshift-cloud-credential-operator.svc](#cco-metrics.openshift-cloud-credential-operator.svc)
    - [cluster-autoscaler-operator.openshift-machine-api.svc](#cluster-autoscaler-operator.openshift-machine-api.svc)
    - [cluster-monitoring-operator.openshift-monitoring.svc](#cluster-monitoring-operator.openshift-monitoring.svc)
    - [cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc](#cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc)
    - [console.openshift-console.svc](#console.openshift-console.svc)
    - [controller-manager.openshift-controller-manager.svc](#controller-manager.openshift-controller-manager.svc)
    - [csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc](#csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc)
    - [csi-snapshot-webhook.openshift-cluster-storage-operator.svc](#csi-snapshot-webhook.openshift-cluster-storage-operator.svc)
    - [dns-default.openshift-dns.svc](#dns-default.openshift-dns.svc)
    - [etcd.openshift-etcd.svc](#etcd.openshift-etcd.svc)
    - [grafana.openshift-monitoring.svc](#grafana.openshift-monitoring.svc)
    - [image-registry-operator.openshift-image-registry.svc](#image-registry-operator.openshift-image-registry.svc)
    - [image-registry.openshift-image-registry.svc](#image-registry.openshift-image-registry.svc)
    - [kube-controller-manager.openshift-kube-controller-manager.svc](#kube-controller-manager.openshift-kube-controller-manager.svc)
    - [kube-state-metrics.openshift-monitoring.svc](#kube-state-metrics.openshift-monitoring.svc)
    - [machine-api-controllers.openshift-machine-api.svc](#machine-api-controllers.openshift-machine-api.svc)
    - [machine-api-operator-webhook.openshift-machine-api.svc](#machine-api-operator-webhook.openshift-machine-api.svc)
    - [machine-api-operator.openshift-machine-api.svc](#machine-api-operator.openshift-machine-api.svc)
    - [machine-approver.openshift-cluster-machine-approver.svc](#machine-approver.openshift-cluster-machine-approver.svc)
    - [machine-config-daemon.openshift-machine-config-operator.svc](#machine-config-daemon.openshift-machine-config-operator.svc)
    - [marketplace-operator-metrics.openshift-marketplace.svc](#marketplace-operator-metrics.openshift-marketplace.svc)
    - [metrics.openshift-apiserver-operator.svc](#metrics.openshift-apiserver-operator.svc)
    - [metrics.openshift-authentication-operator.svc](#metrics.openshift-authentication-operator.svc)
    - [metrics.openshift-cluster-samples-operator.svc](#metrics.openshift-cluster-samples-operator.svc)
    - [metrics.openshift-config-operator.svc](#metrics.openshift-config-operator.svc)
    - [metrics.openshift-console-operator.svc](#metrics.openshift-console-operator.svc)
    - [metrics.openshift-controller-manager-operator.svc](#metrics.openshift-controller-manager-operator.svc)
    - [metrics.openshift-dns-operator.svc](#metrics.openshift-dns-operator.svc)
    - [metrics.openshift-etcd-operator.svc](#metrics.openshift-etcd-operator.svc)
    - [metrics.openshift-ingress-operator.svc](#metrics.openshift-ingress-operator.svc)
    - [metrics.openshift-insights.svc](#metrics.openshift-insights.svc)
    - [metrics.openshift-kube-apiserver-operator.svc](#metrics.openshift-kube-apiserver-operator.svc)
    - [metrics.openshift-kube-controller-manager-operator.svc](#metrics.openshift-kube-controller-manager-operator.svc)
    - [metrics.openshift-kube-scheduler-operator.svc](#metrics.openshift-kube-scheduler-operator.svc)
    - [metrics.openshift-kube-storage-version-migrator-operator.svc](#metrics.openshift-kube-storage-version-migrator-operator.svc)
    - [metrics.openshift-service-ca-operator.svc](#metrics.openshift-service-ca-operator.svc)
    - [multus-admission-controller.openshift-multus.svc](#multus-admission-controller.openshift-multus.svc)
    - [network-metrics-service.openshift-multus.svc](#network-metrics-service.openshift-multus.svc)
    - [node-exporter.openshift-monitoring.svc](#node-exporter.openshift-monitoring.svc)
    - [node-tuning-operator.openshift-cluster-node-tuning-operator.svc](#node-tuning-operator.openshift-cluster-node-tuning-operator.svc)
    - [oauth-openshift.openshift-authentication.svc](#oauth-openshift.openshift-authentication.svc)
    - [olm-operator-metrics.openshift-operator-lifecycle-manager.svc](#olm-operator-metrics.openshift-operator-lifecycle-manager.svc)
    - [openshift-state-metrics.openshift-monitoring.svc](#openshift-state-metrics.openshift-monitoring.svc)
    - [prometheus-adapter.openshift-monitoring.svc](#prometheus-adapter.openshift-monitoring.svc)
    - [prometheus-k8s-thanos-sidecar.openshift-monitoring.svc](#prometheus-k8s-thanos-sidecar.openshift-monitoring.svc)
    - [prometheus-k8s.openshift-monitoring.svc](#prometheus-k8s.openshift-monitoring.svc)
    - [prometheus-operator.openshift-monitoring.svc](#prometheus-operator.openshift-monitoring.svc)
    - [router-internal-default.openshift-ingress.svc](#router-internal-default.openshift-ingress.svc)
    - [scheduler.openshift-kube-scheduler.svc](#scheduler.openshift-kube-scheduler.svc)
    - [sdn.openshift-sdn.svc](#sdn.openshift-sdn.svc)
    - [telemeter-client.openshift-monitoring.svc](#telemeter-client.openshift-monitoring.svc)
    - [thanos-querier.openshift-monitoring.svc](#thanos-querier.openshift-monitoring.svc)
- [Client Certificate/Key Pairs](#client-certificatekey-pairs)
- [Certificates Without Keys](#certificates-without-keys)
- [Certificate Authority Bundles](#certificate-authority-bundles)
    - [service-ca](#service-ca)

## Signing Certificate/Key Pairs


### service-serving-signer
![PKI Graph](subcert-openshift-service-serving-signer16221335705446463666206287945.png)

Signer used by service-ca to sign serving certificates for internal service DNS names.

| Property | Value |
| ----------- | ----------- |
| Type | Signer |
| CommonName | openshift-service-serving-signer@1622133570 |
| SerialNumber | 5446463666206287945 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y60d |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment<br/>- KeyUsageCertSign |
| ExtendedUsages |  |


#### service-serving-signer Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-service-ca | signing-key |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Serving Certificate/Key Pairs


### alertmanager-main.openshift-monitoring.svc
![PKI Graph](subcert-alertmanager-main.openshift-monitoring.svc7889519590428116393.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | alertmanager-main.openshift-monitoring.svc |
| SerialNumber | 7889519590428116393 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - alertmanager-main.openshift-monitoring.svc<br/>- alertmanager-main.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### alertmanager-main.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | alertmanager-main-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### api.openshift-apiserver.svc
![PKI Graph](subcert-api.openshift-apiserver.svc2115297822024425807.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | api.openshift-apiserver.svc |
| SerialNumber | 2115297822024425807 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - api.openshift-apiserver.svc<br/>- api.openshift-apiserver.svc.cluster.local |
| IP Addresses |  |


#### api.openshift-apiserver.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-apiserver | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### api.openshift-oauth-apiserver.svc
![PKI Graph](subcert-api.openshift-oauth-apiserver.svc485864516996010702.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | api.openshift-oauth-apiserver.svc |
| SerialNumber | 485864516996010702 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - api.openshift-oauth-apiserver.svc<br/>- api.openshift-oauth-apiserver.svc.cluster.local |
| IP Addresses |  |


#### api.openshift-oauth-apiserver.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-oauth-apiserver | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### catalog-operator-metrics.openshift-operator-lifecycle-manager.svc
![PKI Graph](subcert-catalog-operator-metrics.openshift-operator-lifecycle-manager.svc5069841447494153423.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | catalog-operator-metrics.openshift-operator-lifecycle-manager.svc |
| SerialNumber | 5069841447494153423 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - catalog-operator-metrics.openshift-operator-lifecycle-manager.svc<br/>- catalog-operator-metrics.openshift-operator-lifecycle-manager.svc.cluster.local |
| IP Addresses |  |


#### catalog-operator-metrics.openshift-operator-lifecycle-manager.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-operator-lifecycle-manager | catalog-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### cco-metrics.openshift-cloud-credential-operator.svc
![PKI Graph](subcert-cco-metrics.openshift-cloud-credential-operator.svc6893942648061066111.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | cco-metrics.openshift-cloud-credential-operator.svc |
| SerialNumber | 6893942648061066111 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - cco-metrics.openshift-cloud-credential-operator.svc<br/>- cco-metrics.openshift-cloud-credential-operator.svc.cluster.local |
| IP Addresses |  |


#### cco-metrics.openshift-cloud-credential-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cloud-credential-operator | cloud-credential-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### cluster-autoscaler-operator.openshift-machine-api.svc
![PKI Graph](subcert-cluster-autoscaler-operator.openshift-machine-api.svc8305498258745803921.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | cluster-autoscaler-operator.openshift-machine-api.svc |
| SerialNumber | 8305498258745803921 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - cluster-autoscaler-operator.openshift-machine-api.svc<br/>- cluster-autoscaler-operator.openshift-machine-api.svc.cluster.local |
| IP Addresses |  |


#### cluster-autoscaler-operator.openshift-machine-api.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-machine-api | cluster-autoscaler-operator-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### cluster-monitoring-operator.openshift-monitoring.svc
![PKI Graph](subcert-cluster-monitoring-operator.openshift-monitoring.svc8944425219050445342.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | cluster-monitoring-operator.openshift-monitoring.svc |
| SerialNumber | 8944425219050445342 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - cluster-monitoring-operator.openshift-monitoring.svc<br/>- cluster-monitoring-operator.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### cluster-monitoring-operator.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | cluster-monitoring-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc
![PKI Graph](subcert-cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc5294915743012167366.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc |
| SerialNumber | 5294915743012167366 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc<br/>- cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc.cluster.local |
| IP Addresses |  |


#### cluster-storage-operator-metrics.openshift-cluster-storage-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-storage-operator | cluster-storage-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### console.openshift-console.svc
![PKI Graph](subcert-console.openshift-console.svc2317112508926355245.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | console.openshift-console.svc |
| SerialNumber | 2317112508926355245 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - console.openshift-console.svc<br/>- console.openshift-console.svc.cluster.local |
| IP Addresses |  |


#### console.openshift-console.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-console | console-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### controller-manager.openshift-controller-manager.svc
![PKI Graph](subcert-controller-manager.openshift-controller-manager.svc8478506552664981432.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | controller-manager.openshift-controller-manager.svc |
| SerialNumber | 8478506552664981432 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - controller-manager.openshift-controller-manager.svc<br/>- controller-manager.openshift-controller-manager.svc.cluster.local |
| IP Addresses |  |


#### controller-manager.openshift-controller-manager.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-controller-manager | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc
![PKI Graph](subcert-csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc5025200834724127258.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc |
| SerialNumber | 5025200834724127258 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc<br/>- csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc.cluster.local |
| IP Addresses |  |


#### csi-snapshot-controller-operator-metrics.openshift-cluster-storage-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-storage-operator | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### csi-snapshot-webhook.openshift-cluster-storage-operator.svc
![PKI Graph](subcert-csi-snapshot-webhook.openshift-cluster-storage-operator.svc1282769300244468729.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | csi-snapshot-webhook.openshift-cluster-storage-operator.svc |
| SerialNumber | 1282769300244468729 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - csi-snapshot-webhook.openshift-cluster-storage-operator.svc<br/>- csi-snapshot-webhook.openshift-cluster-storage-operator.svc.cluster.local |
| IP Addresses |  |


#### csi-snapshot-webhook.openshift-cluster-storage-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-storage-operator | csi-snapshot-webhook-secret |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### dns-default.openshift-dns.svc
![PKI Graph](subcert-dns-default.openshift-dns.svc495137039081958925.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | dns-default.openshift-dns.svc |
| SerialNumber | 495137039081958925 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - dns-default.openshift-dns.svc<br/>- dns-default.openshift-dns.svc.cluster.local |
| IP Addresses |  |


#### dns-default.openshift-dns.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-dns | dns-default-metrics-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### etcd.openshift-etcd.svc
![PKI Graph](subcert-etcd.openshift-etcd.svc1695572914480243966.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | etcd.openshift-etcd.svc |
| SerialNumber | 1695572914480243966 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - etcd.openshift-etcd.svc<br/>- etcd.openshift-etcd.svc.cluster.local |
| IP Addresses |  |


#### etcd.openshift-etcd.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-etcd | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### grafana.openshift-monitoring.svc
![PKI Graph](subcert-grafana.openshift-monitoring.svc5127637701693466147.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | grafana.openshift-monitoring.svc |
| SerialNumber | 5127637701693466147 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - grafana.openshift-monitoring.svc<br/>- grafana.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### grafana.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | grafana-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### image-registry-operator.openshift-image-registry.svc
![PKI Graph](subcert-image-registry-operator.openshift-image-registry.svc4967320171357519668.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | image-registry-operator.openshift-image-registry.svc |
| SerialNumber | 4967320171357519668 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - image-registry-operator.openshift-image-registry.svc<br/>- image-registry-operator.openshift-image-registry.svc.cluster.local |
| IP Addresses |  |


#### image-registry-operator.openshift-image-registry.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-image-registry | image-registry-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### image-registry.openshift-image-registry.svc
![PKI Graph](subcert-image-registry.openshift-image-registry.svc7911780555769594156.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | image-registry.openshift-image-registry.svc |
| SerialNumber | 7911780555769594156 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - image-registry.openshift-image-registry.svc<br/>- image-registry.openshift-image-registry.svc.cluster.local |
| IP Addresses |  |


#### image-registry.openshift-image-registry.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-image-registry | image-registry-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### kube-controller-manager.openshift-kube-controller-manager.svc
![PKI Graph](subcert-kube-controller-manager.openshift-kube-controller-manager.svc4706016511474554482.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | kube-controller-manager.openshift-kube-controller-manager.svc |
| SerialNumber | 4706016511474554482 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - kube-controller-manager.openshift-kube-controller-manager.svc<br/>- kube-controller-manager.openshift-kube-controller-manager.svc.cluster.local |
| IP Addresses |  |


#### kube-controller-manager.openshift-kube-controller-manager.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-controller-manager | serving-cert |
| openshift-kube-controller-manager | serving-cert-2 |
| openshift-kube-controller-manager | serving-cert-3 |
| openshift-kube-controller-manager | serving-cert-4 |
| openshift-kube-controller-manager | serving-cert-5 |
| openshift-kube-controller-manager | serving-cert-6 |
| openshift-kube-controller-manager | serving-cert-7 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-3/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-3/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-4/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-4/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-6/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-6/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### kube-state-metrics.openshift-monitoring.svc
![PKI Graph](subcert-kube-state-metrics.openshift-monitoring.svc2719079659670312610.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | kube-state-metrics.openshift-monitoring.svc |
| SerialNumber | 2719079659670312610 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - kube-state-metrics.openshift-monitoring.svc<br/>- kube-state-metrics.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### kube-state-metrics.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | kube-state-metrics-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### machine-api-controllers.openshift-machine-api.svc
![PKI Graph](subcert-machine-api-controllers.openshift-machine-api.svc7828335248087138693.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | machine-api-controllers.openshift-machine-api.svc |
| SerialNumber | 7828335248087138693 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - machine-api-controllers.openshift-machine-api.svc<br/>- machine-api-controllers.openshift-machine-api.svc.cluster.local |
| IP Addresses |  |


#### machine-api-controllers.openshift-machine-api.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-machine-api | machine-api-controllers-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### machine-api-operator-webhook.openshift-machine-api.svc
![PKI Graph](subcert-machine-api-operator-webhook.openshift-machine-api.svc2625486651396182955.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | machine-api-operator-webhook.openshift-machine-api.svc |
| SerialNumber | 2625486651396182955 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - machine-api-operator-webhook.openshift-machine-api.svc<br/>- machine-api-operator-webhook.openshift-machine-api.svc.cluster.local |
| IP Addresses |  |


#### machine-api-operator-webhook.openshift-machine-api.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-machine-api | machine-api-operator-webhook-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### machine-api-operator.openshift-machine-api.svc
![PKI Graph](subcert-machine-api-operator.openshift-machine-api.svc5923902699021639505.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | machine-api-operator.openshift-machine-api.svc |
| SerialNumber | 5923902699021639505 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - machine-api-operator.openshift-machine-api.svc<br/>- machine-api-operator.openshift-machine-api.svc.cluster.local |
| IP Addresses |  |


#### machine-api-operator.openshift-machine-api.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-machine-api | machine-api-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### machine-approver.openshift-cluster-machine-approver.svc
![PKI Graph](subcert-machine-approver.openshift-cluster-machine-approver.svc1831298527724831562.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | machine-approver.openshift-cluster-machine-approver.svc |
| SerialNumber | 1831298527724831562 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - machine-approver.openshift-cluster-machine-approver.svc<br/>- machine-approver.openshift-cluster-machine-approver.svc.cluster.local |
| IP Addresses |  |


#### machine-approver.openshift-cluster-machine-approver.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-machine-approver | machine-approver-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### machine-config-daemon.openshift-machine-config-operator.svc
![PKI Graph](subcert-machine-config-daemon.openshift-machine-config-operator.svc894043062816778974.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | machine-config-daemon.openshift-machine-config-operator.svc |
| SerialNumber | 894043062816778974 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - machine-config-daemon.openshift-machine-config-operator.svc<br/>- machine-config-daemon.openshift-machine-config-operator.svc.cluster.local |
| IP Addresses |  |


#### machine-config-daemon.openshift-machine-config-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-machine-config-operator | proxy-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### marketplace-operator-metrics.openshift-marketplace.svc
![PKI Graph](subcert-marketplace-operator-metrics.openshift-marketplace.svc9089605778427485628.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | marketplace-operator-metrics.openshift-marketplace.svc |
| SerialNumber | 9089605778427485628 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - marketplace-operator-metrics.openshift-marketplace.svc<br/>- marketplace-operator-metrics.openshift-marketplace.svc.cluster.local |
| IP Addresses |  |


#### marketplace-operator-metrics.openshift-marketplace.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-marketplace | marketplace-operator-metrics |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-apiserver-operator.svc
![PKI Graph](subcert-metrics.openshift-apiserver-operator.svc7400757488192498955.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-apiserver-operator.svc |
| SerialNumber | 7400757488192498955 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-apiserver-operator.svc<br/>- metrics.openshift-apiserver-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-apiserver-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-apiserver-operator | openshift-apiserver-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-authentication-operator.svc
![PKI Graph](subcert-metrics.openshift-authentication-operator.svc6301067254016768819.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-authentication-operator.svc |
| SerialNumber | 6301067254016768819 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-authentication-operator.svc<br/>- metrics.openshift-authentication-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-authentication-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-authentication-operator | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-cluster-samples-operator.svc
![PKI Graph](subcert-metrics.openshift-cluster-samples-operator.svc2686164050326297964.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-cluster-samples-operator.svc |
| SerialNumber | 2686164050326297964 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-cluster-samples-operator.svc<br/>- metrics.openshift-cluster-samples-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-cluster-samples-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-samples-operator | samples-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-config-operator.svc
![PKI Graph](subcert-metrics.openshift-config-operator.svc1282522901324235330.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-config-operator.svc |
| SerialNumber | 1282522901324235330 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-config-operator.svc<br/>- metrics.openshift-config-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-config-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-config-operator | config-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-console-operator.svc
![PKI Graph](subcert-metrics.openshift-console-operator.svc4841722672729242428.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-console-operator.svc |
| SerialNumber | 4841722672729242428 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-console-operator.svc<br/>- metrics.openshift-console-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-console-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-console-operator | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-controller-manager-operator.svc
![PKI Graph](subcert-metrics.openshift-controller-manager-operator.svc7895726074443218984.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-controller-manager-operator.svc |
| SerialNumber | 7895726074443218984 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-controller-manager-operator.svc<br/>- metrics.openshift-controller-manager-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-controller-manager-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-controller-manager-operator | openshift-controller-manager-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-dns-operator.svc
![PKI Graph](subcert-metrics.openshift-dns-operator.svc7601847597278589785.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-dns-operator.svc |
| SerialNumber | 7601847597278589785 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-dns-operator.svc<br/>- metrics.openshift-dns-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-dns-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-dns-operator | metrics-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-etcd-operator.svc
![PKI Graph](subcert-metrics.openshift-etcd-operator.svc3978805962409959490.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-etcd-operator.svc |
| SerialNumber | 3978805962409959490 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-etcd-operator.svc<br/>- metrics.openshift-etcd-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-etcd-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-etcd-operator | etcd-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-ingress-operator.svc
![PKI Graph](subcert-metrics.openshift-ingress-operator.svc8279326142192924063.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-ingress-operator.svc |
| SerialNumber | 8279326142192924063 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-ingress-operator.svc<br/>- metrics.openshift-ingress-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-ingress-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-ingress-operator | metrics-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-insights.svc
![PKI Graph](subcert-metrics.openshift-insights.svc6511383752257504790.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-insights.svc |
| SerialNumber | 6511383752257504790 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-insights.svc<br/>- metrics.openshift-insights.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-insights.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-insights | openshift-insights-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-kube-apiserver-operator.svc
![PKI Graph](subcert-metrics.openshift-kube-apiserver-operator.svc29949893468932305.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-kube-apiserver-operator.svc |
| SerialNumber | 29949893468932305 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-kube-apiserver-operator.svc<br/>- metrics.openshift-kube-apiserver-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-kube-apiserver-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-apiserver-operator | kube-apiserver-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-kube-controller-manager-operator.svc
![PKI Graph](subcert-metrics.openshift-kube-controller-manager-operator.svc1704350015219690809.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-kube-controller-manager-operator.svc |
| SerialNumber | 1704350015219690809 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-kube-controller-manager-operator.svc<br/>- metrics.openshift-kube-controller-manager-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-kube-controller-manager-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-controller-manager-operator | kube-controller-manager-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-kube-scheduler-operator.svc
![PKI Graph](subcert-metrics.openshift-kube-scheduler-operator.svc5130875085637374527.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-kube-scheduler-operator.svc |
| SerialNumber | 5130875085637374527 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-kube-scheduler-operator.svc<br/>- metrics.openshift-kube-scheduler-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-kube-scheduler-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-scheduler-operator | kube-scheduler-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-kube-storage-version-migrator-operator.svc
![PKI Graph](subcert-metrics.openshift-kube-storage-version-migrator-operator.svc5764379107559423336.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-kube-storage-version-migrator-operator.svc |
| SerialNumber | 5764379107559423336 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-kube-storage-version-migrator-operator.svc<br/>- metrics.openshift-kube-storage-version-migrator-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-kube-storage-version-migrator-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-storage-version-migrator-operator | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### metrics.openshift-service-ca-operator.svc
![PKI Graph](subcert-metrics.openshift-service-ca-operator.svc8465006101141555491.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | metrics.openshift-service-ca-operator.svc |
| SerialNumber | 8465006101141555491 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - metrics.openshift-service-ca-operator.svc<br/>- metrics.openshift-service-ca-operator.svc.cluster.local |
| IP Addresses |  |


#### metrics.openshift-service-ca-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-service-ca-operator | serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### multus-admission-controller.openshift-multus.svc
![PKI Graph](subcert-multus-admission-controller.openshift-multus.svc4660313081021989101.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | multus-admission-controller.openshift-multus.svc |
| SerialNumber | 4660313081021989101 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - multus-admission-controller.openshift-multus.svc<br/>- multus-admission-controller.openshift-multus.svc.cluster.local |
| IP Addresses |  |


#### multus-admission-controller.openshift-multus.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-multus | multus-admission-controller-secret |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### network-metrics-service.openshift-multus.svc
![PKI Graph](subcert-network-metrics-service.openshift-multus.svc1889672894063829328.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | network-metrics-service.openshift-multus.svc |
| SerialNumber | 1889672894063829328 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - network-metrics-service.openshift-multus.svc<br/>- network-metrics-service.openshift-multus.svc.cluster.local |
| IP Addresses |  |


#### network-metrics-service.openshift-multus.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-multus | metrics-daemon-secret |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### node-exporter.openshift-monitoring.svc
![PKI Graph](subcert-node-exporter.openshift-monitoring.svc6997256025626337803.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | node-exporter.openshift-monitoring.svc |
| SerialNumber | 6997256025626337803 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - node-exporter.openshift-monitoring.svc<br/>- node-exporter.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### node-exporter.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | node-exporter-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### node-tuning-operator.openshift-cluster-node-tuning-operator.svc
![PKI Graph](subcert-node-tuning-operator.openshift-cluster-node-tuning-operator.svc7832843275082956404.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | node-tuning-operator.openshift-cluster-node-tuning-operator.svc |
| SerialNumber | 7832843275082956404 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - node-tuning-operator.openshift-cluster-node-tuning-operator.svc<br/>- node-tuning-operator.openshift-cluster-node-tuning-operator.svc.cluster.local |
| IP Addresses |  |


#### node-tuning-operator.openshift-cluster-node-tuning-operator.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-cluster-node-tuning-operator | node-tuning-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### oauth-openshift.openshift-authentication.svc
![PKI Graph](subcert-oauth-openshift.openshift-authentication.svc455687257200358236.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | oauth-openshift.openshift-authentication.svc |
| SerialNumber | 455687257200358236 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - oauth-openshift.openshift-authentication.svc<br/>- oauth-openshift.openshift-authentication.svc.cluster.local |
| IP Addresses |  |


#### oauth-openshift.openshift-authentication.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-authentication | v4-0-config-system-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### olm-operator-metrics.openshift-operator-lifecycle-manager.svc
![PKI Graph](subcert-olm-operator-metrics.openshift-operator-lifecycle-manager.svc8800549647405875654.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | olm-operator-metrics.openshift-operator-lifecycle-manager.svc |
| SerialNumber | 8800549647405875654 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - olm-operator-metrics.openshift-operator-lifecycle-manager.svc<br/>- olm-operator-metrics.openshift-operator-lifecycle-manager.svc.cluster.local |
| IP Addresses |  |


#### olm-operator-metrics.openshift-operator-lifecycle-manager.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-operator-lifecycle-manager | olm-operator-serving-cert |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### openshift-state-metrics.openshift-monitoring.svc
![PKI Graph](subcert-openshift-state-metrics.openshift-monitoring.svc7882046295044958152.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | openshift-state-metrics.openshift-monitoring.svc |
| SerialNumber | 7882046295044958152 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - openshift-state-metrics.openshift-monitoring.svc<br/>- openshift-state-metrics.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### openshift-state-metrics.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | openshift-state-metrics-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### prometheus-adapter.openshift-monitoring.svc
![PKI Graph](subcert-prometheus-adapter.openshift-monitoring.svc2945501834265381842.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | prometheus-adapter.openshift-monitoring.svc |
| SerialNumber | 2945501834265381842 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - prometheus-adapter.openshift-monitoring.svc<br/>- prometheus-adapter.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### prometheus-adapter.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | prometheus-adapter-8tkqrsmu9afpe |
| openshift-monitoring | prometheus-adapter-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### prometheus-k8s-thanos-sidecar.openshift-monitoring.svc
![PKI Graph](subcert-prometheus-k8s-thanos-sidecar.openshift-monitoring.svc8595435866243050124.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | prometheus-k8s-thanos-sidecar.openshift-monitoring.svc |
| SerialNumber | 8595435866243050124 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - prometheus-k8s-thanos-sidecar.openshift-monitoring.svc<br/>- prometheus-k8s-thanos-sidecar.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### prometheus-k8s-thanos-sidecar.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | prometheus-k8s-thanos-sidecar-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### prometheus-k8s.openshift-monitoring.svc
![PKI Graph](subcert-prometheus-k8s.openshift-monitoring.svc8353222847585982884.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | prometheus-k8s.openshift-monitoring.svc |
| SerialNumber | 8353222847585982884 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - prometheus-k8s.openshift-monitoring.svc<br/>- prometheus-k8s.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### prometheus-k8s.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | prometheus-k8s-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### prometheus-operator.openshift-monitoring.svc
![PKI Graph](subcert-prometheus-operator.openshift-monitoring.svc4178324176483245369.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | prometheus-operator.openshift-monitoring.svc |
| SerialNumber | 4178324176483245369 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - prometheus-operator.openshift-monitoring.svc<br/>- prometheus-operator.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### prometheus-operator.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | prometheus-operator-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### router-internal-default.openshift-ingress.svc
![PKI Graph](subcert-router-internal-default.openshift-ingress.svc5035056584329763135.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | router-internal-default.openshift-ingress.svc |
| SerialNumber | 5035056584329763135 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - router-internal-default.openshift-ingress.svc<br/>- router-internal-default.openshift-ingress.svc.cluster.local |
| IP Addresses |  |


#### router-internal-default.openshift-ingress.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-ingress | router-metrics-certs-default |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### scheduler.openshift-kube-scheduler.svc
![PKI Graph](subcert-scheduler.openshift-kube-scheduler.svc6657891906279326247.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | scheduler.openshift-kube-scheduler.svc |
| SerialNumber | 6657891906279326247 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - scheduler.openshift-kube-scheduler.svc<br/>- scheduler.openshift-kube-scheduler.svc.cluster.local |
| IP Addresses |  |


#### scheduler.openshift-kube-scheduler.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-kube-scheduler | serving-cert |
| openshift-kube-scheduler | serving-cert-2 |
| openshift-kube-scheduler | serving-cert-3 |
| openshift-kube-scheduler | serving-cert-4 |
| openshift-kube-scheduler | serving-cert-5 |
| openshift-kube-scheduler | serving-cert-6 |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-3/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-3/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-5/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-5/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-6/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-6/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-2/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-2/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-4/secrets/serving-cert/tls.crt/tls.crt | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-scheduler-pod-4/secrets/serving-cert/tls.crt/tls.key | -rw-------. | root | root | system_u:object_r:kubernetes_file_t:s0 |



### sdn.openshift-sdn.svc
![PKI Graph](subcert-sdn.openshift-sdn.svc2625658655664643767.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | sdn.openshift-sdn.svc |
| SerialNumber | 2625658655664643767 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - sdn.openshift-sdn.svc<br/>- sdn.openshift-sdn.svc.cluster.local |
| IP Addresses |  |


#### sdn.openshift-sdn.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-sdn | sdn-metrics-certs |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### telemeter-client.openshift-monitoring.svc
![PKI Graph](subcert-telemeter-client.openshift-monitoring.svc1112524969442289252.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | telemeter-client.openshift-monitoring.svc |
| SerialNumber | 1112524969442289252 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - telemeter-client.openshift-monitoring.svc<br/>- telemeter-client.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### telemeter-client.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | telemeter-client-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |




### thanos-querier.openshift-monitoring.svc
![PKI Graph](subcert-thanos-querier.openshift-monitoring.svc2348697830836353775.png)



| Property | Value |
| ----------- | ----------- |
| Type | Serving |
| CommonName | thanos-querier.openshift-monitoring.svc |
| SerialNumber | 2348697830836353775 |
| Issuer CommonName | [service-serving-signer](#service-serving-signer) |
| Validity | 2y |
| Signature Algorithm | SHA256-RSA |
| PublicKey Algorithm | RSA 2048 bit |
| Usages | - KeyUsageDigitalSignature<br/>- KeyUsageKeyEncipherment |
| ExtendedUsages | - ExtKeyUsageServerAuth |
| DNS Names | - thanos-querier.openshift-monitoring.svc<br/>- thanos-querier.openshift-monitoring.svc.cluster.local |
| IP Addresses |  |


#### thanos-querier.openshift-monitoring.svc Locations
| Namespace | Secret Name |
| ----------- | ----------- |
| openshift-monitoring | thanos-querier-tls |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |



## Client Certificate/Key Pairs

## Certificates Without Keys

These certificates are present in certificate authority bundles, but do not have keys in the cluster.
This happens when the installer bootstrap clusters with a set of certificate/key pairs that are deleted during the
installation process.

## Certificate Authority Bundles


### service-ca
![PKI Graph](subca-3983882995.png)

CA for recognizing serving certificates for services that were signed by our service-ca controller.

**Bundled Certificates**

| CommonName | Issuer CommonName | Validity | PublicKey Algorithm |
| ----------- | ----------- | ----------- | ----------- |
| [service-serving-signer](#service-serving-signer) | [service-serving-signer](#service-serving-signer) | 2y60d | RSA 2048 bit |

#### service-ca Locations
| Namespace | ConfigMap Name |
| ----------- | ----------- |
| openshift-config-managed | service-ca |
| openshift-kube-controller-manager | service-ca |
| openshift-kube-controller-manager | service-ca-2 |
| openshift-kube-controller-manager | service-ca-3 |
| openshift-kube-controller-manager | service-ca-4 |
| openshift-kube-controller-manager | service-ca-5 |
| openshift-kube-controller-manager | service-ca-6 |
| openshift-kube-controller-manager | service-ca-7 |
| openshift-service-ca | signing-cabundle |

| File | Permissions | User | Group | SE Linux |
| ----------- | ----------- | ----------- | ----------- | ----------- |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-3/configmaps/service-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-4/configmaps/service-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-5/configmaps/service-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-6/configmaps/service-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |
| /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-7/configmaps/service-ca/ca-bundle.crt/ca-bundle.crt | -rw-r--r--. | root | root | system_u:object_r:kubernetes_file_t:s0 |


