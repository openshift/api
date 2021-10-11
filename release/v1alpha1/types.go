package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleasePayload encapsulates the information for the creation of a ReleasePayload
// and aggregates the results of its respective verification tests.
//
// The release-controller is configured to monitor imagestreams, in a specific namespace, that are annotated with a
// ReleaseConfig.  The ReleaseConfig is a definition of how releases are calculated.  When a ReleasePayload is
// generated, it will be generated in the same namespace as the imagstream that produced it. If/when an update
// occurs, to one of these imagestreams, the release-controller will:
//   1) Create a point-in-time mirror of the updated imagestream
//   2) Create a new Release from the mirror
//        - Any errors before this point will cause the release to marked `Failed`
//   3) Launches a set of release analysis jobs
//   4) Launches an aggregation job
//   5) Launches a set of release verification jobs
//        - These can either be `Blocking Jobs` which will prevent release acceptance or `Informing Jobs` which will
//          not prevent release acceptance.
//   6) Monitors for job completions
//        - If all `Blocking Jobs` complete successfully, then the release is `Accepted`.  If any `Blocking Jobs` fail,
//          the release will be marked `Rejected`
//   7) Publishes all results to the respective webpage
//
// Example:
// ART:
//   1) Publishes an update to the `ocp/4.9-art-latest` imagestream
//
// Release-controller:
//   1) Creates a mirror named: `ocp/4.9-art-latest-2021-09-27-105859`
//   2) Creates a ReleasePayload: `ocp/4.9.0-0.nightly-2021-09-27-105859-<random-string>`
//       -Labels:
//         release.openshift.io/imagestream=release
//         release.openshift.io/imagestreamtag-name=4.9.0-0.nightly-2021-09-27-105859
//   3) Creates an OpenShift Release: `ocp/release:4.9.0-0.nightly-2021-09-27-105859`
//   4) Update ReleasePayload conditions with results of release creation job
//   If the release was created successfully, the release-controller:
//   5) Launches: 4.9.0-0.nightly-2021-09-27-105859-aggregated-<name>-analysis-<count>
//   6) Launches: 4.9.0-0.nightly-2021-09-27-105859-aggregated-<name>-aggregator
//   7) Launches: 4.9.0-0.nightly-2021-09-27-105859-<name>
//
// Mapping from a Release to ReleasePayload:
// A ReleasePayload will always be named after the Release that it corresponds to, with the addition of a
// random string suffix.  Both objects will reside in the same namespace.
//   For a release: `ocp/release:4.9.0-0.nightly-2021-09-27-105859`
//   A corresponding ReleasePayload will exist: `ocp/4.9.0-0.nightly-2021-09-27-105859-<random-string>`
//
// Mapping from ReleasePayload to Release:
// A ReleasePayload is decorated with a couple labels that will point back to the Release that it corresponds to:
//   - release.openshift.io/imagestream=release
//   - release.openshift.io/imagestreamtag-name=4.9.0-0.nightly-2021-09-27-105859
// Because the ReleasePayload and the Release will both reside in the same namespace, the release that created the
// ReleasePayload will be located here:
//   <namespace>/<release.openshift.io/imagestream>:<release.openshift.io/imagestreamtag-name>
// Similarly, the ReleasePayload object itself also has the PayloadCoordinates (.spec.payloadCoordinates) that point
// back to the Release as well:
//   spec:
//     payloadCoordinates:
//       imagestreamName: release
//       imagestreamTagName: 4.9.0-0.nightly-2021-09-27-105859
//       namespace: ocp
// The release that created the ReleasePayload will be located here:
//   <namespace>/<imagestreamName>:<imagestreamTagName>
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ReleasePayload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec the inputs used to create the ReleasePayload
	Spec ReleasePayloadSpec `json:"spec,omitempty"`

	// Status is the current status of the ReleasePayload
	Status ReleasePayloadStatus `json:"status,omitempty"`
}

// ReleasePayloadSpec has the information to represent a ReleasePayload
type ReleasePayloadSpec struct {
	PayloadCoordinates PayloadCoordinates     `json:"payloadCoordinates,omitempty"`
	PayloadOverride    ReleasePayloadOverride `json:"payloadOverride,omitempty"`
}

// PayloadCoordinates houses the information pointing to the location of the imagesteamtag that this ReleasePayload
// is verifying.
//
// Example:
// For a ReleasePayload named: "4.9.0-0.nightly-2021-09-27-105859-<random-string>" in the "ocp" namespace, and configured
// to be written into the "release" imagestream, we expect:
//   1) Namespace to equal "ocp
//   2) ImagestreamName to equal "release"
//   3) ImagestreamTagName to equal "4.9.0-0.nightly-2021-09-27-105859", which will also serves as the prefix of the ReleasePayload
//
// These coordinates can then be used to get the release imagestreamtag itself:
//    # oc -n ocp get imagestreamtag release:4.9.0-0.nightly-2021-09-27-105859
type PayloadCoordinates struct {
	// Namespace must match that of the ReleasePayload
	Namespace string `json:"namespace,omitempty"`

	// ImagestreamName is the location of the configured "release" imagestream
	//   - This is a configurable parameter ("to") passed into the release-controller via the ReleaseConfig's defined here:
	//     https://github.com/openshift/release/blob/master/core-services/release-controller/_releases
	ImagestreamName string `json:"imagestreamName,omitempty"`

	// ImagestreamTagName is the name of the actual release
	ImagestreamTagName string `json:"imagestreamTagName,omitempty"`
}

type ReleasePayloadOverrideType string

// These are the supported ReleasePayloadOverride values.
const (
	// ReleasePayloadOverrideAccepted enables the manual Acceptance of a ReleasePayload.
	ReleasePayloadOverrideAccepted ReleasePayloadOverrideType = "Accepted"

	// ReleasePayloadOverrideRejected enables the manual Rejection of a ReleasePayload.
	ReleasePayloadOverrideRejected ReleasePayloadOverrideType = "Rejected"
)

// ReleasePayloadOverride provides the ability to manually Accept/Reject a ReleasePayload
// ART, occasionally, needs the ability to manually accept/reject a Release that, for some reason or another:
//   - won't pass one or more of it's blocking jobs.
//   - shouldn't proceed with the normal release verification processing
// This would be the one scenario where another party, besides the release-controller, would update a
// ReleasePayload instance.  Upon doing so, the release-controller should see that an update occurred and make all
// the necessary changes to formally accept/reject the respective release.
type ReleasePayloadOverride struct {
	// Override specifies the ReleasePayloadOverride to apply to the ReleasePayload
	Override ReleasePayloadOverrideType `json:"override"`

	// Reason is a human-readable string that specifies the reason for manually overriding the
	// Acceptance/Rejections of a ReleasePayload
	Reason string `json:"reason"`
}

// ReleasePayloadStatus the status of all the promotion test jobs
type ReleasePayloadStatus struct {
	// Conditions communicates the state of the ReleasePayload.
	// Supported conditions include PayloadCreated, PayloadFailed, PayloadAccepted, and PayloadRejected.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// BlockingJobResults stores the results of all blocking jobs
	BlockingJobResults []JobStatus `json:"blockingJobResults,omitempty"`

	// InformingJobResults stores the results of all informing jobs
	InformingJobResults []JobStatus `json:"informingJobResults,omitempty"`

	// AnalysisJobResults stores the results of all analysis jobs
	AnalysisJobResults []JobStatus `json:"analysisJobResults,omitempty"`
}

type ReleasePayloadStatusConditionType string

// These are valid conditions of ReleasePayloadStatus.
const (
	// PayloadCreated if false, the ReleasePayload is waiting for a release image to be created and pushed to the
	// TargetImageStream.  If PayloadCreated is true, a release image has been created and pushed to the TargetImageStream.
	// Verification jobs should begin and will update the status as they complete.
	PayloadCreated ReleasePayloadStatusConditionType = "PayloadCreated"

	// PayloadFailed is true if a ReleasePayload image cannot be created for the given set of image mirrors
	// This condition is terminal
	PayloadFailed ReleasePayloadStatusConditionType = "PayloadFailed"

	// PayloadAccepted is true if the ReleasePayload has passed its verification criteria and can safely
	// be promoted to an external location
	// This condition is terminal
	PayloadAccepted ReleasePayloadStatusConditionType = "PayloadAccepted"

	// PayloadRejected is true if the ReleasePayload has failed one or more of its verification criteria
	// The release-controller will take no more action in this phase.
	PayloadRejected ReleasePayloadStatusConditionType = "PayloadRejected"
)

// ReleasePayloadStatusCondition contains condition information for a tag event.
type ReleasePayloadStatusCondition struct {
	// Type of release payload status condition
	Type ReleasePayloadStatusConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// LastTransitionTIme is the time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a brief machine readable explanation for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Message is a human readable description of the details about last transition, complementing reason.
	Message string `json:"message,omitempty"`
	// Generation is the spec tag generation that this status corresponds to
	Generation int64 `json:"generation"`
}

// JobState the aggregate state of the job
// Supported values include Pending, Failed, Success, and Ignored.
type JobState string

const (
	// JobStatePending not all job runs have completed
	// Transitions to Failed or Success
	JobStatePending JobState = "Pending"

	// JobStateFailed failed job aggregation
	JobStateFailed JobState = "Failed"

	// JobStateSuccess successful job aggregation
	JobStateSuccess JobState = "Success"

	// JobStateIgnored ignored job aggregation
	// This is specifically for Analysis jobs that are not being monitored directly.
	JobStateIgnored JobState = "Ignored"
)

// JobStatus encapsulates the name of the job, all the results of the jobs, and an aggregated
// result of all the jobs
type JobStatus struct {
	// jobName is the name of the job
	JobName string `json:"name,omitempty"`

	// MaxRetries maximum times to retry a job
	MaxRetries int `json:"maxRetries,omitempty"`

	// AggregateState is the overall success/failure of all the executed jobs
	AggregateState JobState `json:"state,omitempty"`

	// JobRunResults contains the links for individual jobs
	JobRunResults []JobRunResult `json:"results,omitempty"`
}

// JobRunState the status of a job
type JobRunState string

const (
	// JobRunStateTriggered job has been created but not scheduled
	JobRunStateTriggered JobRunState = "Triggered"

	// JobRunStatePending job is running and awaiting completion
	JobRunStatePending JobRunState = "Pending"

	// JobRunStateFailure job completed with errors
	JobRunStateFailure JobRunState = "Failure"

	// JobRunStateSuccess job completed without errors
	JobRunStateSuccess JobRunState = "Success"

	// JobRunStateAborted job was terminated early
	JobRunStateAborted JobRunState = "Aborted"

	// JobRunStateError job could not be scheduled
	JobRunStateError JobRunState = "Error"

	// JobRunStateIgnored job that is submitted to the system but is not monitored for success/failure
	// This is specifically for Analysis jobs that are monitored by the JobRunAggregator.
	JobRunStateIgnored JobRunState = "Ignored"
)

// JobRunCoordinates houses the information necessary to locate individual job executions
type JobRunCoordinates struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
}

// JobRunResult the results of a prowjob run
// The release-controller creates ProwJobs (prowv1.ProwJob) during the sync_ready control loop and relies on an informer
// to process jobs, that it created, as they are completed. The JobRunResults will be created, by the release-controller
// during the sync_ready loop and updated whenever any changes, to the respective job is received by the informer.
type JobRunResult struct {
	// Coordinates the location of the job
	Coordinates JobRunCoordinates `json:"coordinates,omitempty"`

	// StartTime timestamp for when the prowjob was created
	StartTime metav1.Time `json:"startTime,omitempty"`

	// CompletionTime timestamp for when the prow pipeline controller observes the final state of the ProwJob
	// For instance, if a client Aborts a ProwJob, the Pipeline controller will receive notification of the change
	// and update the PtowJob's Status accordingly.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// State the current state of the job run
	State JobRunState `json:"state"`

	// HumanProwResultsURL the html link to the prow results
	HumanProwResultsURL string `json:"humanProwResultsURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleasePayloadList is a list of ReleasePayloads
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ReleasePayloadList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of ReleasePayloads
	Items []ReleasePayload `json:"items"`
}
