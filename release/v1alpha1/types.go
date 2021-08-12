package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleasePayload aggregates the outputs of the promotion tests
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type ReleasePayload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec the inputs used to create the release payload
	Spec ReleasePayloadSpec `json:"spec,omitempty"`

	// Status is the current status of the release payload
	Status ReleasePayloadStatus `json:"status,omitempty"`
}

// ReleasePayloadSpec has the information to represent a PromotionTest
type ReleasePayloadSpec struct {
	//TODO: What are the required pieces for this...
}

// ReleasePayloadStatus the status of all the promotion test jobs
type ReleasePayloadStatus struct {
	// BlockingJobResults stores the results of all blocking jobs
	BlockingJobResults []JobResults `json:"blockingJobResults,omitempty"`
	// InformingJobResults stores the results of all informing jobs
	InformingJobResults []JobResults `json:"informingJobResults,omitempty"`
	// AnalysisJobResults stores the results of all analysis jobs
	AnalysisJobResults []JobResults `json:"analysisJobResults,omitempty"`

	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ReleasePayloadPromotionConditionType specifies the current condition type of the ReleasePayload
type ReleasePayloadPromotionConditionType string

const (
        Passed ReleasePayloadPromotionConditionType = "Passed"
		Failed ReleasePayloadPromotionConditionType = "Failed"
)

// JobState the final state of the job
type JobState string

const (
	// PendingJobState not all job runs have completed
	PendingJobState JobState = "Pending"
	// FailedJobState failed job aggregation
	FailedJobState JobState = "Failed"
	// SuccessJobState successful job aggregation
	SuccessJobState JobState = "Success"
)

// JobResults encapsulates the name of the job, all the results of the jobs, and an aggregated
// result of all the jobs
type JobResults struct {
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
	// PendingJobRunState job currently running
	PendingJobRunState JobRunState = "Pending"
	// RunningJobRunState job currently running
	RunningJobRunState JobRunState = "Running"
	// FailedJobRunState job failed
	FailedJobRunState JobRunState = "Failed"
	// SuccessJobRunState job successful
	SuccessJobRunState JobRunState = "Success"
)

// JobRunResult the results of a job run
type JobRunResult struct {
	// RunID the id of the job
	RunId int `json:"runId"`
	// State the current state of the job run
	State JobRunState `json:"state"`
	// HumanProwResultsURL the html link to the prow results
	HumanProwResultsURL string `json:"humanProwResultsURL"`

	//TODO: Add field for GCS bucket
}
