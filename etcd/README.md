# etcd.openshift.io API Group

This API group contains CRDs related to etcd cluster management in Two Node OpenShift with Fencing deployments.

## API Versions

### v1alpha1

Contains the `PacemakerCluster` custom resource for monitoring Pacemaker cluster health in Two Node OpenShift with Fencing deployments.

#### PacemakerCluster

- **Feature Gate**: `DualReplica`
- **Component**: `two-node-fencing`
- **Scope**: Cluster-scoped singleton resource (must be named "cluster")
- **Resource Path**: `pacemakerclusters.etcd.openshift.io`

The `PacemakerCluster` resource provides visibility into the health and status of a Pacemaker-managed cluster.
It is periodically updated by the cluster-etcd-operator's status collector.

### Pacemaker Resources

A **pacemaker resource** is a unit of work managed by pacemaker. In pacemaker terminology, resources are services
or applications that pacemaker monitors, starts, stops, and moves between nodes to maintain high availability.

For Two Node OpenShift with Fencing, we manage three resources:
- **Kubelet**: The Kubernetes node agent and a prerequisite for etcd
- **Etcd**: The distributed key-value store
- **FencingAgent**: Used to isolate failed nodes during a quorum loss event

### Status Structure

> **Note on optional fields:** Per Kubernetes API conventions, all status fields are marked `+optional` in the Go
> type definitions to allow for partial updates and backward compatibility. However, CEL validation rules enforce
> that certain fields (like conditions and lastUpdated) cannot be removed once set. This means these fields are
> "optional on initial creation, but required once present."
> See: [kube-api-linter statusoptional](https://github.com/kubernetes-sigs/kube-api-linter/blob/main/pkg/analysis/statusoptional/doc.go)

```yaml
status:
  conditions:              # Cluster-level conditions (required once set, min 3 items)
    - type: Healthy
    - type: InService
    - type: NodeCountAsExpected
  lastUpdated: <timestamp> # When status was last updated (required once set, cannot decrease)
  nodes:                   # Per-node status (0-32 nodes, expects 2)
    - name: <hostname>     # RFC 1123 subdomain name
      ipAddresses:         # List of global unicast IPv4 or IPv6 addresses in canonical form
        - <ip>             # First address used for etcd peer URLs
      conditions:          # Node-level conditions (required once set, min 7 items)
        - type: Healthy
        - type: Online
        - type: InService
        - type: Active
        - type: Ready
        - type: Clean
        - type: Member
      resources:           # Array of pacemaker resources on this node (required once set, min 3 items)
        - name: Kubelet    # All three resources (Kubelet, Etcd, FencingAgent) must be present
          conditions:      # Resource-level conditions (required once set, min 8 items)
            - type: Healthy
            - type: InService
            - type: Managed
            - type: Enabled
            - type: Operational
            - type: Active
            - type: Started
            - type: Schedulable
        - name: Etcd
          conditions: []
        - name: FencingAgent
          conditions: []
```

### Cluster-Level Conditions

**All three conditions are required once the status is populated.**

| Condition | True | False |
|-----------|------|-------|
| `Healthy` | Cluster is healthy (`ClusterHealthy`) | Cluster has issues (`ClusterUnhealthy`) |
| `InService` | In service (`InService`) | In maintenance (`InMaintenance`) |
| `NodeCountAsExpected` | Node count is as expected (`AsExpected`) | Wrong count (`InsufficientNodes`, `ExcessiveNodes`) |

### Node-Level Conditions

**All seven conditions are required for each node once the status is populated.**

| Condition | True | False |
|-----------|------|-------|
| `Healthy` | Node is healthy (`NodeHealthy`) | Node has issues (`NodeUnhealthy`) |
| `Online` | Node is online (`Online`) | Node is offline (`Offline`) |
| `InService` | In service (`InService`) | In maintenance (`InMaintenance`) |
| `Active` | Node is active (`Active`) | Node is in standby (`Standby`) |
| `Ready` | Node is ready (`Ready`) | Node is pending (`Pending`) |
| `Clean` | Node is clean (`Clean`) | Node is unclean (`Unclean`) |
| `Member` | Node is a member (`Member`) | Not a member (`NotMember`) |

### Resource-Level Conditions

Each resource in the `resources` array has its own conditions. **All eight conditions are required for each resource once the status is populated.**

| Condition | True | False |
|-----------|------|-------|
| `Healthy` | Resource is healthy (`ResourceHealthy`) | Resource has issues (`ResourceUnhealthy`) |
| `InService` | In service (`InService`) | In maintenance (`InMaintenance`) |
| `Managed` | Managed by pacemaker (`Managed`) | Not managed (`Unmanaged`) |
| `Enabled` | Resource is enabled (`Enabled`) | Resource is disabled (`Disabled`) |
| `Operational` | Resource is operational (`Operational`) | Resource has failed (`Failed`) |
| `Active` | Resource is active (`Active`) | Resource is not active (`Inactive`) |
| `Started` | Resource is started (`Started`) | Resource is stopped (`Stopped`) |
| `Schedulable` | Resource is schedulable (`Schedulable`) | Resource is not schedulable (`Unschedulable`) |

### Validation Rules

**Resource naming:**
- Resource name must be "cluster" (singleton)

**Node name validation:**
- Must be a lowercase RFC 1123 subdomain name
- Consists of lowercase alphanumeric characters, '-' or '.'
- Must start and end with an alphanumeric character
- Maximum 253 characters

**IP address validation:**
- Pacemaker allows multiple IP addresses for Corosync communication between nodes (1-8 addresses)
- The first address in the list is used for IP-based peer URLs for etcd membership
- Each address must be in canonical form (e.g., `192.168.1.1` not `192.168.001.001`, or `2001:db8::1` not `2001:0db8::1`)
- Each address must be a valid global unicast IPv4 or IPv6 address (including private/RFC1918 addresses)
- Excludes loopback, link-local, and multicast addresses

**Timestamp validation:**
- `lastUpdated` is optional on initial creation, but once set it cannot be removed and must always increase (prevents stale updates)

**Required conditions (once status is populated):**
- All cluster-level conditions must be present (MinItems=3)
- All node-level conditions must be present for each node (MinItems=7)
- All resource-level conditions must be present for each resource in the `resources` array (MinItems=8)

**Resource names:**
- Valid values are: `Kubelet`, `Etcd`, `FencingAgent`
- All three resources must be present in each node's `resources` array

### Usage

The cluster-etcd-operator healthcheck controller watches this resource and updates operator conditions based on
the cluster state. The aggregate `Healthy` conditions at each level (cluster, node, resource) provide a quick
way to determine overall health.
