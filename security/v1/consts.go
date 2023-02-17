package v1

const (
	UIDRangeAnnotation = "openshift.io/sa.scc.uid-range"
	// SupplementalGroupsAnnotation contains a comma delimited list of allocated supplemental groups
	// for the namespace.  Groups are in the form of a Block which supports {start}/{length} or {start}-{end}
	SupplementalGroupsAnnotation = "openshift.io/sa.scc.supplemental-groups"
	MCSAnnotation                = "openshift.io/sa.scc.mcs"
	ValidatedSCCAnnotation       = "openshift.io/scc"
	// SCCHistory holds all prior SCCs that the pod was assigned to, if any, from most recent to least recent.
	SCCHistory = "openshift.io/scc-history"
)
