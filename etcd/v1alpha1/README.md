# etcd.openshift.io/v1alpha1

This API group contains types related to two-node fencing for etcd cluster management.

## PacemakerCluster

The `PacemakerCluster` CRD provides visibility into the health and status of Pacemaker-managed clusters in dual-replica (two-node) OpenShift deployments.

### Feature Gate

- **Feature Gate**: None - this CRD is gated by cluster-etcd-operator start-up. It will only be created once a TNF cluster has transitioned to external etcd.
- **Component**: `two-node-fencing`

### Usage

The PacemakerCluster resource is a cluster-scoped, status-only singleton named "cluster". It is periodically updated by a privileged controller that runs `pcs status xml` and parses the output into structured fields for health checking.

### Status Fields

- **LastUpdated** (required): Timestamp when status was last collected
- **Summary**: High-level cluster state including:
  - `pacemakerDaemonState`: Running state of the pacemaker daemon (enum: `Running`, `KnownNotRunning`)
  - `quorumStatus`: Whether cluster has quorum (enum: `Quorate`, `NoQuorum`)
  - `nodesOnline`, `nodesTotal`: Node counts
  - `resourcesStarted`, `resourcesTotal`: Resource counts
- **Nodes**: Detailed per-node status (name, IPv4/IPv6 addresses, online status, mode)
- **Resources**: Detailed per-resource status (name, resource agent type, role enum, active status, node assignment)
- **NodeHistory**: Recent operation history for troubleshooting (operation failures within last 5 minutes)
- **FencingHistory**: Recent fencing events (events within last 24 hours)
- **RawXML**: Complete XML output from `pcs status xml` (for debugging only, max 256KB)
- **CollectionError**: Any errors encountered during status collection

### Design Principles

The API follows a "Design Principle: Act on Deterministic Information" approach:
- Almost all fields are optional except `lastUpdated`
- Missing data means "unknown" not "error"
- The operator only transitions between PacemakerHealthy and PacemakerDegraded states based on deterministic information
- When information is unavailable, the last known state is preserved

### Notes

The spec field is reserved but unused - all meaningful data is in the status subresource.
