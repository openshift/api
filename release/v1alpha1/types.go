package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PromotionTestResults encapsulates the inputs to produce a new set of promotion tests and
// the status of the promotion tests
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type PromotionTestResults struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec the inputs used to create the promotion test
	Spec PromotionTestSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the current status of the promotion test
	Status PromotionTestResultsStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PromotionTestSpec has the information to represent a PromotionTest
type PromotionTestSpec struct {
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
}

// PromotionTestResultsStatus the status of all the promotion test jobs
type PromotionTestResultsStatus struct {
	JobResults []JobResults `json:"results,omitempty" protobuf:"bytes,1,rep,name=results"`
}

// JobAggregateResultState the final state of the job
type JobAggregateResultState string

const (
	// FailedJobAggregateResultState failed job aggregation
	FailedJobAggregateResultState JobAggregateResultState = "Failed"
	// SuccessJobAggregateResultState successful job aggregation
	SuccessJobAggregateResultState JobAggregateResultState = "Success"
)

// JobResults encapsulates the name of the job, all the results of the jobs, and an aggregated
// result of all the jobs
type JobResults struct {
	// jobName is the name of the job
	JobName string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// AggregateState is the overall success/failure of all the executed jobs
	AggregateState JobAggregateResultState `json:"state,omitempty" protobuf:"bytes,2,opt,name=state,casttype=AggregateState"`
	// JobRunResults contains the links for individual jobs
	JobRunResults []JobRunResult `json:"results,omitempty" protobuf:"bytes,3,rep,name=results"`
}

// JobExecutionState the status of a job
type JobExecutionState string

const (
	// PendingJobExecutionState job currently running
	PendingJobExecutionState JobExecutionState = "Pending"
	// FailedJobExecutionState job failed
	FailedJobExecutionState JobExecutionState = "Failed"
	// SuccessJobExecutionState job successful
	SuccessJobExecutionState JobExecutionState = "Success"
)

// JobRunResult the results of a job run
type JobRunResult struct {
	// State the current state of the job run
	State JobExecutionState `json:"state" protobuf:"bytes,1,opt,name=state"`
	// URL the html link to the prow results
	URL string `json:"url" protobuf:"bytes,2,opt,name=url"`
	// Retries the number of times the job has been retried
	Retries int `json:"retries,omitempty" protobuf:"varint,3,opt,name=retries"`
	// TransitionTime the timestamp of the last result change
	TransitionTime *metav1.Time `json:"transitionTime,omitempty" protobuf:"bytes,4,opt,name=transitionTime"`
}
