package v1alpha1

import (
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
//   2) Creates an OpenShift Release: `ocp/release:4.9.0-0.nightly-2021-09-27-105859`
//   3) Launches: 4.9.0-0.nightly-2021-09-27-105859-aggregated-<name>-analysis-<count>
//   4) Launches: 4.9.0-0.nightly-2021-09-27-105859-aggregated-<name>-aggregator
//   5) Launches: 4.9.0-0.nightly-2021-09-27-105859-<name>
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

// ReleasePayloadSpec has the information to represent a PromotionTest
type ReleasePayloadSpec struct {
	PayloadCoordinates PayloadCoordinates `json:"payloadCoordinates,omitempty"`
}

type PayloadCoordinates struct {
	Namespace          string `json:"namespace,omitempty"`
	ImagestreamName    string `json:"imagestreamName,omitempty"`
	ImagestreamTagName string `json:"imagestreamTagName,omitempty"`
}

// ReleasePayloadStatus the status of all the promotion test jobs
type ReleasePayloadStatus struct {
	// Conditions communicates the state of the ReleasePayload.
	//
	// Supported conditions include PayloadCreated, PayloadCreationFailed, Accepted, and Rejected.
	//
	// If PayloadCreated is false the ReleasePayload is waiting for a release image to be created and pushed to the
	// TargetImageStream.  If PayloadCreated is true a release image has been created and pushed to the TargetImageStream.
	// Verification jobs should begin and will update the status as they complete.
	//
	// If PayloadCreationFailed is true a ReleasePayload image cannot be created for the given set of image mirrors
	// This condition is terminal
	//
	// If Accepted is true the ReleasePayload has passed its verification criteria and can safely
	// be promoted to an external location
	// This condition is terminal
	//
	// if Rejected is true the ReleasePayload has failed one or more of its verification criteria
	// The release-controller will take no more action in this phase.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// BlockingJobResults stores the results of all blocking jobs
	BlockingJobResults []JobStatus `json:"blockingJobResults,omitempty"`

	// InformingJobResults stores the results of all informing jobs
	InformingJobResults []JobStatus `json:"informingJobResults,omitempty"`

	// AnalysisJobResults stores the results of all analysis jobs
	AnalysisJobResults []JobStatus `json:"analysisJobResults,omitempty"`
}

// JobState the aggregate state of the job
// Supported values include Pending, Failed, and Success.
type JobState string

const (
	// JobStatePending not all job runs have completed
	// Transitions to Failed or Success
	JobStatePending JobState = "Pending"

	// JobStateFailed failed job aggregation
	JobStateFailed JobState = "Failed"

	// JobStateSuccess successful job aggregation
	JobStateSuccess JobState = "Success"
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
)

// JobRunResult the results of a prowjob run
// The release-controller creates ProwJobs (prowv1.ProwJob) during the sync_ready control loop and relies on an informer
// to process jobs, that it created, as they are completed. The JobRunResults will be created, by the release-controller
// during the sync_ready loop and updated whenever any changes, to the respective job is received by the informer.
type JobRunResult struct {
	// Name unique name for the job run
	Name string `json:"name,omitempty"`

	// Namespace location where the job ran
	Namespace string `json:"namespace,omitempty"`

	// Cluster is which Kubernetes cluster is used to run the job
	Cluster string `json:"cluster,omitempty"`

	// RunID the unique identifier of the job
	RunId int `json:"runId"`

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
