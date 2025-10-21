# etcd.openshift.io API Group

This API group contains CRDs related to etcd cluster management in Two Node OpenShift with Fencing (TNF) deployments.

## API Versions

### v1alpha1

Contains the `PacemakerCluster` custom resource for monitoring Pacemaker cluster health in TNF deployments.

#### PacemakerCluster

- **Feature Gate**: `DualReplica`
- **Component**: `two-node-fencing`
- **Scope**: Cluster-scoped singleton resource (must be named "cluster")
- **Resource Path**: `pacemakerclusters.etcd.openshift.io`

The `PacemakerCluster` resource provides visibility into the health and status of a Pacemaker-managed cluster. It is
periodically updated by the cluster-etcd-operator's status collector.

### Pacemaker Resources

A **pacemaker resource** is a unit of work managed by pacemaker. In pacemaker terminology, resources are services
or applications that pacemaker monitors, starts, stops, and moves between nodes to maintain high availability.

For TNF, we manage three resources:
- **Kubelet**: The Kubernetes node agent and a prerequisite for etcd
- **Etcd**: The distributed key-value store
- **Fencing Agent**: Used to isolate failed nodes during a quorum loss event

### Status Structure

```yaml
status:
  conditions: []           # Cluster-level conditions (required)
  lastUpdated: <timestamp> # When status was last updated
  nodes:                   # Per-node status (1-32 nodes, TNF expects 2)
    - name: <hostname>     # RFC 1123 subdomain name
      ipAddress: <ip>      # Global unicast IPv4 or IPv6
      conditions: []       # Node-level conditions (required)
      kubelet:             # Kubelet resource status
        conditions: []     # Resource-level conditions (required)
      etcd:                # Etcd resource status
        conditions: []
      fencingAgent:        # Fencing agent resource status
        conditions: []
```

### Cluster-Level Conditions

All conditions follow positive polarity (True = healthy). **All three conditions are required.**

| Condition | True | False |
|-----------|------|-------|
| `PacemakerClusterHealthy` | Cluster is healthy (`ClusterHealthy`) | Cluster has issues (`ClusterUnhealthy`) |
| `PacemakerClusterNotInMaintenanceMode` | Not in maintenance (`NotInMaintenance`) | In maintenance (`InMaintenance`) |
| `PacemakerClusterExpectedNodeCount` | Expected nodes present (`ExpectedNodeCount`) | Wrong count (`InsufficientNodes`, `ExcessiveNodes`) |

### Node-Level Conditions

**All seven conditions are required for each node.**

| Condition | True | False |
|-----------|------|-------|
| `PacemakerClusterNodeHealthy` | Node is healthy (`NodeHealthy`) | Node has issues (`NodeUnhealthy`) |
| `PacemakerClusterNodeOnline` | Node is online (`Online`) | Node is offline (`Offline`) |
| `PacemakerClusterNodeNotInMaintenance` | Not in maintenance (`NotInMaintenance`) | In maintenance (`InMaintenance`) |
| `PacemakerClusterNodeActive` | Node is active (`NodeActive`) | Node is in standby (`NodeStandby`) |
| `PacemakerClusterNodeReady` | Node is ready (`NodeReady`) | Node is pending (`NodePending`) |
| `PacemakerClusterNodeClean` | Node is clean (`NodeClean`) | Node is unclean (`NodeUnclean`) |
| `PacemakerClusterNodeMember` | Node is a member (`Member`) | Not a member (`NotMember`) |

### Resource-Level Conditions

Each resource (kubelet, etcd, fencingAgent) has its own conditions. **All eight conditions are required for each resource.**

| Condition | True | False |
|-----------|------|-------|
| `PacemakerClusterResourceHealthy` | Resource is healthy (`ResourceHealthy`) | Resource has issues (`ResourceUnhealthy`) |
| `PacemakerClusterResourceNotInMaintenance` | Not in maintenance (`NotInMaintenance`) | In maintenance (`InMaintenance`) |
| `PacemakerClusterResourceManaged` | Managed by pacemaker (`Managed`) | Not managed (`NotManaged`) |
| `PacemakerClusterResourceEnabled` | Resource is enabled (`Enabled`) | Resource is disabled (`Disabled`) |
| `PacemakerClusterResourceOperational` | Resource is operational (`Operational`) | Resource has failed (`Failed`) |
| `PacemakerClusterResourceActive` | Resource is active (`Active`) | Resource is not active (`NotActive`) |
| `PacemakerClusterResourceStarted` | Resource is started (`Started`) | Resource is stopped (`Stopped`) |
| `PacemakerClusterResourceUnblocked` | Resource is unblocked (`Unblocked`) | Resource is blocked (`Blocked`) |

### Validation Rules

**Resource naming:**
- Resource name must be "cluster" (singleton)

**Node name validation:**
- Must be a lowercase RFC 1123 subdomain name
- Consists of lowercase alphanumeric characters, '-' or '.'
- Must start and end with an alphanumeric character
- Maximum 253 characters

**IP address validation:**
- Must be a valid canonical global unicast IPv4 or IPv6 address
- Excludes loopback, link-local, and multicast addresses

**Timestamp validation:**
- `lastUpdated` timestamp must always increase (prevents stale updates)

**Required conditions:**
- All cluster-level conditions must be present
- All node-level conditions must be present for each node
- All resource-level conditions must be present for each resource (kubelet, etcd, fencingAgent)

### Usage

The cluster-etcd-operator healthcheck controller watches this resource and updates operator conditions based on
the cluster state. The aggregate `*Healthy` conditions at each level (cluster, node, resource) provide a quick
way to determine overall health.
