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

### Status Subresource Design

This resource uses the standard Kubernetes status subresource pattern (`+kubebuilder:subresource:status`).
The status collector creates the resource without status, then immediately populates it via the `/status` endpoint.

**Why not atomic create-with-status?**

We initially explored removing the status subresource to allow creating the resource with status in a single
atomic operation. This would ensure the resource is never observed in an incomplete state. However:

1. The Kubernetes API server strips the `status` field from create requests when a status subresource is enabled
2. Without the subresource, we cannot use separate RBAC for spec vs status updates
3. The OpenShift API test framework assumes status subresource exists for status update tests

The status collector performs a two-step operation: create resource, then immediately update status.
The brief window where status is empty is acceptable since the healthcheck controller handles missing status gracefully.

### Pacemaker Resources

A **pacemaker resource** is a unit of work managed by pacemaker. In pacemaker terminology, resources are services
or applications that pacemaker monitors, starts, stops, and moves between nodes to maintain high availability.

For Two Node OpenShift with Fencing, we manage three resource types:
- **Kubelet**: The Kubernetes node agent and a prerequisite for etcd
- **Etcd**: The distributed key-value store
- **FencingAgent**: Used to isolate failed nodes during a quorum loss event (tracked separately)

### Status Structure

```yaml
status:                    # Optional on creation, populated via status subresource
  conditions:              # Required when status present (min 3 items)
    - type: Healthy
    - type: InService
    - type: NodeCountAsExpected
  lastUpdated: <timestamp> # Required when status present, cannot decrease
  alertAgents:             # Optional: cluster-wide alert agent registration (0-8 items)
    - name: <agent_name>   # "Taint Alert Agent" or "Untaint Alert Agent"
      conditions:          # Required: Alert agent-level conditions (min 2 items)
        - type: Healthy
        - type: Configured
  nodes:                   # Control-plane nodes (0-5, expects 2 for TNF)
    - nodeName: <hostname> # RFC 1123 subdomain name
      addresses:           # Required: List of node addresses (1-8 items)
        - type: InternalIP # Currently only InternalIP is supported
          address: <ip>    # First address used for etcd peer URLs
      conditions:          # Required: Node-level conditions (min 9 items)
        - type: Healthy
        - type: Online
        - type: InService
        - type: Active
        - type: Ready
        - type: Clean
        - type: Member
        - type: FencingAvailable
        - type: FencingHealthy
      resources:           # Required: Pacemaker resources on this node (min 2)
        - name: Kubelet    # Both Kubelet and Etcd must be present
          conditions:      # Required: Resource-level conditions (min 8 items)
            - type: Healthy
            - type: InService
            - type: Managed
            - type: Enabled
            - type: Operational
            - type: Active
            - type: Started
            - type: Schedulable
        - name: Etcd
          conditions: [...]  # Same 8 conditions as Kubelet (abbreviated)
      fencingAgents:       # Required: Fencing agents for THIS node (1-8)
        - name: <unique_id>  # e.g., "master-0_redfish" (unique, max 300 chars)
          method: <method>   # Fencing method: "Redfish" or "IPMI"
          conditions: [...]  # Same 8 conditions as resources (abbreviated)
      alertAgentScripts:   # Optional: per-node alert agent script presence (0-8 items)
        - name: <agent_name>  # "Taint Alert Agent" or "Untaint Alert Agent"
          conditions:      # Required: Alert agent script-level conditions (min 2 items)
            - type: Healthy
            - type: ScriptPresent
```

### Fencing Agents

Fencing agents are STONITH (Shoot The Other Node In The Head) devices used to isolate failed nodes.
Unlike regular pacemaker resources (Kubelet, Etcd), fencing agents are tracked separately because:

1. **Mapping by target, not schedule**: Resources are mapped to the node where they are scheduled to run.
   Fencing agents are mapped to the node they can *fence* (their target), regardless of which node
   their monitoring operations are scheduled on.

2. **Multiple agents per node**: A node can have multiple fencing agents for redundancy
   (e.g., both Redfish and IPMI). Expected: 1 per node, supported: up to 8.

3. **Health tracking via two node-level conditions**:
   - **FencingAvailable**: True if at least one agent is healthy (fencing works), False if all agents unhealthy (degrades operator)
   - **FencingHealthy**: True if all agents are healthy (ideal state), False if any agent is unhealthy (emits warning events)

### Alert Agents

Alert agents are pacemaker alert handlers used to automatically taint a node after it is fenced, and remove
that taint once the node rejoins the cluster. There are two known alert agents: "Taint Alert Agent" and
"Untaint Alert Agent". Unlike resources and fencing agents, alert agents are registered cluster-wide (a single
CIB object shared by all nodes via Pacemaker's CIB replication), so their registration status is tracked once
in `status.alertAgents` rather than per node.

The script each alert agent invokes, however, must exist locally on whichever node the triggering event
occurs on, since Pacemaker executes it there. Script delivery is handled independently per node by MCO, so
presence is tracked per node in `status.nodes[].alertAgentScripts` to catch delivery gaps between nodes.

Both `alertAgents` and `alertAgentScripts` are optional fields. They are omitted when this status has not yet
been collected by the status collector (including by a collector version that predates these fields), and
`alertAgents` is additionally omitted when no alert agents are configured.

### Cluster-Level Conditions

| Condition | True | False |
|-----------|------|-------|
| `Healthy` | Cluster is healthy (`ClusterHealthy`) | Cluster has issues (`ClusterUnhealthy`) |
| `InService` | In service (`InService`) | In maintenance (`InMaintenance`) |
| `NodeCountAsExpected` | Node count is as expected (`AsExpected`) | Wrong count (`InsufficientNodes`, `ExcessiveNodes`) |

### Node-Level Conditions

| Condition | True | False |
|-----------|------|-------|
| `Healthy` | Node is healthy (`NodeHealthy`) | Node has issues (`NodeUnhealthy`) |
| `Online` | Node is online (`Online`) | Node is offline (`Offline`) |
| `InService` | In service (`InService`) | In maintenance (`InMaintenance`) |
| `Active` | Node is active (`Active`) | Node is in standby (`Standby`) |
| `Ready` | Node is ready (`Ready`) | Node is pending (`Pending`) |
| `Clean` | Node is clean (`Clean`) | Node is unclean (`Unclean`) |
| `Member` | Node is a member (`Member`) | Not a member (`NotMember`) |
| `FencingAvailable` | At least one agent healthy (`FencingAvailable`) | All agents unhealthy (`FencingUnavailable`) - degrades operator |
| `FencingHealthy` | All agents healthy (`FencingHealthy`) | Some agents unhealthy (`FencingUnhealthy`) - emits warnings |

### Resource-Level Conditions

Each resource in the `resources` array and each fencing agent in the `fencingAgents` array has its own conditions.

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

### Alert Agent Conditions

Each entry in the `alertAgents` array has its own conditions.

| Condition | True | False | Unknown |
|-----------|------|-------|---------|
| `Healthy` | Alert agent is healthy (`AlertAgentHealthy`) | Alert agent has issues (`AlertAgentUnhealthy`) | Not yet observed (`Pending`) |
| `Configured` | Registered in the CIB as expected (`Configured`) | Not registered (`Missing`) or registered with an unexpected path/filter (`Misconfigured`) | Not yet observed (`Pending`) |

### Alert Agent Script Conditions

Each entry in a node's `alertAgentScripts` array has its own conditions.

| Condition | True | False | Unknown |
|-----------|------|-------|---------|
| `Healthy` | Alert agent script is healthy (`AlertAgentScriptHealthy`) | Alert agent script has issues (`AlertAgentScriptUnhealthy`) | Not yet observed (`Pending`) |
| `ScriptPresent` | Script is present and executable on this node (`Present`) | Script is missing from this node (`Missing`) | Not yet observed (`Pending`) |

### Validation Rules

**Resource naming:**
- Resource name must be "cluster" (singleton)

**Node name validation:**
- Must be a lowercase RFC 1123 subdomain name
- Consists of lowercase alphanumeric characters, '-' or '.'
- Must start and end with an alphanumeric character
- Maximum 253 characters

**Node addresses:**
- Uses `PacemakerNodeAddress` type (similar to `corev1.NodeAddress` but with IP validation)
- Currently only `InternalIP` type is supported
- Pacemaker allows multiple addresses for Corosync communication between nodes (1-8 addresses)
- The first address in the list is used for IP-based peer URLs for etcd membership
- IP validation:
  - Must be a valid global unicast IPv4 or IPv6 address
  - Must be in canonical form (e.g., `192.168.1.1` not `192.168.001.001`, or `2001:db8::1` not `2001:0db8::1`)
  - Excludes loopback, link-local, and multicast addresses
  - Maximum length is 39 characters (full IPv6 address)

**Timestamp validation:**
- `lastUpdated` is required when status is present
- Once set, cannot be set to an earlier timestamp (validation uses `!has(oldSelf.lastUpdated)` to handle initial creation)
- Timestamps must always increase (prevents stale updates from overwriting newer data)

**Status fields:**
- `status` - Optional on creation (pointer type), populated via status subresource
- When status is present, `conditions`, `lastUpdated`, and `nodes` are required; `alertAgents` is optional:
  - `conditions` - Required array of cluster conditions (min 3 items)
  - `lastUpdated` - Required timestamp for staleness detection
  - `nodes` - Required array of control-plane node statuses (min 0, max 5; empty allowed for catastrophic failures)
  - `alertAgents` - Optional array of cluster-wide alert agent registration status (min 0, max 8 items); omitted when not yet collected or when no alert agents are configured

**Node fields (when node present):**
- `nodeName` - Required, RFC 1123 subdomain
- `addresses` - Required (min 1, max 8 items)
- `conditions` - Required (min 9 items with specific types enforced via XValidation)
- `resources` - Required (min 2 items: Kubelet and Etcd)
- `fencingAgents` - Required (min 1, max 8 items)
- `alertAgentScripts` - Optional (min 0, max 8 items); omitted when not yet collected by the status collector

**Conditions validation:**
- Cluster-level: MinItems=3 (Healthy, InService, NodeCountAsExpected)
- Node-level: MinItems=9 (Healthy, Online, InService, Active, Ready, Clean, Member, FencingAvailable, FencingHealthy)
- Resource-level: MinItems=8 (Healthy, InService, Managed, Enabled, Operational, Active, Started, Schedulable)
- Fencing agent-level: MinItems=8 (same conditions as resources)
- Alert agent-level: MinItems=2, MaxItems=8 (Healthy, Configured)
- Alert agent script-level: MinItems=2, MaxItems=8 (Healthy, ScriptPresent)

All condition arrays have XValidation rules to ensure specific condition types are present.

**Alert agent names:**
- Valid values are: `Taint Alert Agent`, `Untaint Alert Agent`
- Names must be unique within the `alertAgents` and `alertAgentScripts` arrays (enforced via XValidation), but neither array requires both names to be present

**Resource names:**
- Valid values are: `Kubelet`, `Etcd`
- Both resources must be present in each node's `resources` array

**Fencing agent fields:**
- `name`: Unique identifier for the fencing agent (e.g., "master-0_redfish")
  - Must be unique within the `fencingAgents` array
  - May contain alphanumeric characters, dots, hyphens, and underscores (`^[a-zA-Z0-9._-]+$`)
  - Maximum 300 characters (provides headroom beyond 253 node name + underscore + method)
- `method`: Fencing method enum - valid values are `Redfish` or `IPMI`
- `conditions`: Required, same 8 conditions as resources

Note: The target node is implied by the parent `PacemakerClusterNodeStatus` - fencing agents are nested under the node they can fence.

### Usage

The cluster-etcd-operator healthcheck controller watches this resource and updates operator conditions based on
the cluster state. The aggregate `Healthy` conditions at each level (cluster, node, resource) provide a quick
way to determine overall health.
