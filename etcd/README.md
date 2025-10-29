# etcd.openshift.io API Group

This API group contains CRDs related to etcd cluster management. Specifically, this is only used for TNF (Two Node Fencing)
for gathering status updates from the node to ensure the cluster-admin is warned about unhealthy setups.

## API Versions

### v1alpha1

Contains the `PacemakerCluster` custom resource for monitoring Pacemaker cluster health in TNF (Two Node Fencing) deployments.

#### PacemakerCluster

- **Feature Gate**: None - this CRD is gated by cluster-etcd-operator start-up. It will only be created once a TNF cluster has transitioned to external etcd.
- **Component**: `two-node-fencing`
- **Scope**: Cluster-scoped singleton resource named "cluster"
- **Resource Path**: `pacemakerclusters.etcd.openshift.io`

The `PacemakerCluster` resource provides visibility into the health and status of a Pacemaker-managed cluster. It is periodically updated by the cluster-etcd-operator's status collector running as a privileged CronJob.

**Status Fields:**
- `lastUpdated` (required): Timestamp when status was last collected - used to detect stale data
- `summary`: High-level cluster health metrics
  - `pacemakerDaemonState`: Running state (enum: `Running`, `KnownNotRunning`)
  - `quorumStatus`: Quorum state (enum: `Quorate`, `NoQuorum`)
  - `nodesOnline`, `nodesTotal`: Node counts (0-2)
  - `resourcesStarted`, `resourcesTotal`: Resource counts (0-16)
- `nodes`: Detailed status of each node (1-2 nodes)
  - Name, IPv4/IPv6 addresses, online status (enum), mode (enum: `Active`, `Standby`)
- `resources`: Detailed status of each resource (1-16 resources)
  - Name, resource agent, role (enum: `Started`, `Stopped`), active status (enum), node assignment
- `nodeHistory`: Recent operation failures for troubleshooting (up to 16 entries, last 5 minutes)
- `fencingHistory`: Recent fencing events (up to 16 events, last 24 hours)
  - Target node, action (enum: `reboot`, `off`, `on`), status (enum: `success`, `failed`, `pending`), completion timestamp
- `collectionError`: Any errors encountered during status collection (max 2KB)
- `rawXML`: Full XML output from `pcs status xml` for debugging (max 256KB)

**Design Principles:**
The API follows "Act on Deterministic Information":
- All fields except `lastUpdated` are optional
- Missing data indicates unknown state, not error
- Operator only acts on definitive information
- Unknown state preserves the last known health condition

**Usage:**
The cluster-etcd-operator healthcheck controller watches this resource and updates operator conditions based on the cluster state.
