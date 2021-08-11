package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PromotionTestResults aggregates the outputs of the promotion tests
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type PromotionTestResults struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec the inputs used to create the promotion test
	Spec PromotionTestSpec `json:"spec,omitempty"`

	// Status is the current status of the promotion test
	Status PromotionTestResultsStatus `json:"status,omitempty"`
}

// PromotionTestSpec has the information to represent a PromotionTest
type PromotionTestSpec struct {
	//TODO: What are the required pieces for this...
	Name string `json:"name,omitempty"`
}

// PromotionTestResultsStatus the status of all the promotion test jobs
type PromotionTestResultsStatus struct {
	// BlockingJobResults stores the results of all blocking jobs
	BlockingJobResults []JobResults `json:"blockingJobResults,omitempty"`
	// InformingJobResults stores the results of all informing jobs
	InformingJobResults []JobResults `json:"informingJobResults,omitempty"`
	// AnalysisJobResults stores the results of all analysis jobs
	AnalysisJobResults []JobResults `json:"analysisJobResults,omitempty"`
}

// JobAggregateResultState the final state of the job
type JobAggregateResultState string

const (
	// PendingJobAggregateResultState job currently running
	PendingJobAggregateResultState JobAggregateResultState = "Pending"
	// FailedJobAggregateResultState failed job aggregation
	FailedJobAggregateResultState JobAggregateResultState = "Failed"
	// SuccessJobAggregateResultState successful job aggregation
	SuccessJobAggregateResultState JobAggregateResultState = "Success"
)

// JobResults encapsulates the name of the job, all the results of the jobs, and an aggregated
// result of all the jobs
type JobResults struct {
	// jobName is the name of the job
	JobName string `json:"name,omitempty"`
	// AggregateState is the overall success/failure of all the executed jobs
	AggregateState JobAggregateResultState `json:"state,omitempty"`
	// JobRunResults contains the links for individual jobs
	JobRunResults []JobRunResult `json:"results,omitempty"`
}

// JobExecutionState the status of a job
type JobExecutionState string

const (
	// PendingJobExecutionState job currently running
	PendingJobExecutionState JobExecutionState = "Pending"
	// RunningJobExecutionState job currently running
	RunningJobExecutionState JobExecutionState = "Running"
	// FailedJobExecutionState job failed
	FailedJobExecutionState JobExecutionState = "Failed"
	// SuccessJobExecutionState job successful
	SuccessJobExecutionState JobExecutionState = "Success"
)

// JobRunResult the results of a job run
type JobRunResult struct {
	// RunID the id of the job
	RunId int `json:"runId"`
	// State the current state of the job run
	State JobExecutionState `json:"state"`
	// URL the html link to the prow results
	URL string `json:"url"`
	// TransitionTime the timestamp of the last result change
	TransitionTime *metav1.Time `json:"transitionTime,omitempty"`

	//TODO: Add field for GCS bucket
}
