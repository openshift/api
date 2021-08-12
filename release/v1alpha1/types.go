package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleasePayload encapsulates the information for the creation of the ReleasePayload
// and aggregates the outputs of the verification tests.
// The release-controller monitors a specific set of imagestreams for updates.  If/when
// an update occurs to these imagestreams, the release-controller will:
//   1) Create a new ReleasePayload
//      a: in the namespace configured by the --release-namespace parameter of the release-controller
//      b. named after the name or mirrorPrefix in the aforementioned monitored imagestreams
//   2) Start any number of tests specified as an annotation on the aforementioned monitored imagestreams
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
	// Source the name of the imagestream where this imagestream (aka mirror) was created from
	Source string `json:"source,omitempty"`

	// ReleaseTag is the name of the imagestreamtag in the "Target" imagestream
	ReleaseTag string `json:"releaseTag,omitempty"`

	// Target the name of the imagestream where the "ReleaseTag" to this imagestream will be located
	Target string `json:"target,omitempty"`

	// Hash is a sha256 encoded string of the image or dockerImageReference of the first element
	// in Status.Tags of the Source imagestream.  This is used, by the release-controller, to
	// determine if the Source imagestream has any new images
	Hash string `json:"hash,omitempty"`
}

// ReleasePayloadStatus the status of all the promotion test jobs
type ReleasePayloadStatus struct {
	// Conditions communicates the state of the ReleasePayload.
	//
	// Supported conditions include Pending, Ready, Failed, Accepted, and Rejected.
	//
	// If Pending is true the release tag is waiting for an updated payload to be created and pushed
	// Transitions to Failed or Ready
	//
	// If Failed is true a ReleasePayload image cannot be created for the given set of image mirrors
	// This condition is terminal
	//
	// if Ready is true a ReleasePayload has a valid update payload image created and pushed to the
	// release image stream.  Verification jobs should begin and will update the status as they
	// complete
	// Transitions to Accepted or Rejected
	//
	// if Accepted is true the ReleasePayload has passed its verification criteria and can safely
	// be promoted to an external location
	// This condition is terminal
	//
	// if Rejected is true the ReleasePayload has failed one or more of its verification criteria
	// The release-controller will take no more action in this phase, but a human may set the
	// phase back to Ready to retry and the controller will attempt verification again.
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

	// AggregateState is the overall success/failure of all the executed jobs
	AggregateState JobState `json:"state,omitempty"`

	// JobRunResults contains the links for individual jobs
	JobRunResults []JobRunResult `json:"results,omitempty"`
}

// JobRunState the status of a job
type JobRunState string

const (
	// JobRunStatePending job is not running, failed, or succeeded
	JobRunStatePending JobRunState = "Pending"

	// JobRunStateRunning job running
	JobRunStateRunning JobRunState = "Running"

	// JobRunStateFailed job failed
	JobRunStateFailed JobRunState = "Failed"

	// JobRunStateSuccess job successful
	JobRunStateSuccess JobRunState = "Success"
)

// JobRunResult the results of a job run
// The release-controller monitors for prowjobs, that it creates, and creates/updates these results accordingly
type JobRunResult struct {
	// RunID the id of the job
	RunId int `json:"runId"`

	// State the current state of the job run
	State JobRunState `json:"state"`

	// HumanProwResultsURL the html link to the prow results
	HumanProwResultsURL string `json:"humanProwResultsURL"`

	//TODO: Add field for GCS bucket
}
