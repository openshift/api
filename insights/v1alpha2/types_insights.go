package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=datagathers,scope=Cluster
// +kubebuilder:subresource:status
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/2248
// +openshift:file-pattern=cvoRunLevel=0000_10,operatorName=insights,operatorOrdering=01
// +openshift:enable:FeatureGate=InsightsOnDemandDataGather
// +kubebuilder:printcolumn:name=State,type=string,JSONPath=.status.dataGatherState,description=DataGather job state
// +kubebuilder:printcolumn:name=StartTime,type=date,JSONPath=.status.startTime,description=DataGather start time
// +kubebuilder:printcolumn:name=FinishTime,type=date,JSONPath=.status.finishTime,description=DataGather finish time
//
// DataGather provides data gather configuration options and status for the particular Insights data gathering.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type DataGather struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// spec holds user settable values for configuration
	// +required
	Spec DataGatherSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status DataGatherStatus `json:"status"`
}

// DataGatherSpec contains the configuration for the DataGather.
type DataGatherSpec struct {
	// dataPolicy is an optional list of DataPolicyOptions that allows user to enable additional obfuscation of the Insights archive data.
	// It may not exceed 2 items and must not contain duplicates.
	// Valid values are ObfuscateNetworking and WorkloadNames.
	// When set to ObfuscateNetworking the IP addresses and the cluster domain name are obfuscated.
	// When set to WorkloadNames, the gathered data about cluster resources will not contain the workload names for your deployments. Resources UIDs will be used instead.
	// When omitted no obfuscation is applied.
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:XValidation:rule="self.all(x, self.exists_one(y, x == y))",message="dataPolicy items must be unique"
	// +listType=atomic
	// +optional
	DataPolicy []DataPolicyOption `json:"dataPolicy"`
	// gatherers is an optional field that specifies the configuration of the gatherers.
	// +optional
	Gatherers Gatherers `json:"gatherers,omitempty"`
	// storage is an optional field that allows user to define persistent storage for gathering jobs to store the Insights data archive.
	// If omitted, the gathering job will use ephemeral storage.
	// +optional
	Storage *Storage `json:"storage,omitempty"`
}

// Gathereres specifies the configuration of the gatherers
// +kubebuilder:validation:XValidation:rule="has(self.mode) && self.mode == 'Custom' ?  has(self.custom) : !has(self.custom)",message="custom is required when mode is Custom, and forbidden otherwise"
type Gatherers struct {
	// mode is a required field that specifies the mode for gatherers. Allowed values are All, None, and Custom.
	// When set to All, all gatherers wil run and gather data.
	// When set to None, all gatherers will be disabled and no data will be gathered.
	// When set to Custom, the custom configuration from the custom field will be applied.
	// +required
	Mode GatheringMode `json:"mode"`
	// custom provides gathering configuration.
	// It is required when mode is Custom, and forbidden otherwise.
	// Custom configuration allows user to disable only a subset of gatherers.
	// Gatherers that are not explicitly disabled in custom configuration will run.
	// +optional
	Custom *Custom `json:"custom,omitempty"`
}

// custom provides the custom configuration of gatherers
type Custom struct {
	// configs is a required list of gatherers configurations that can be used to enable or disable specific gatherers.
	// It may not exceed 100 items and each gatherer can be present only once.
	// It is possible to disable an entire set of gatherers while allowing a specific function within that set.
	// The particular gatherers IDs can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// Run the following command to get the names of last active gatherers:
	// "oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'"
	// +kubebuilder:validation:MaxItems=100
	// +listType=map
	// +listMapKey=name
	// +required
	Configs []GathererConfig `json:"configs"`
}

// gatheringMode defines the valid gathering modes.
// +kubebuilder:validation:Enum=All;None;Custom
type GatheringMode string

const (
	// Enabled enables all gatherers
	GatheringModeAll GatheringMode = "All"
	// Disabled disables all gatherers
	GatheringModeNone GatheringMode = "None"
	// Custom applies the configuration from GatheringConfig.
	GatheringModeCustom GatheringMode = "Custom"
)

// storage provides persistent storage configuration options for gathering jobs.
// If the type is set to PersistentVolume, then the PersistentVolume must be defined.
// If the type is set to Ephemeral, then the PersistentVolume must not be defined.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'PersistentVolume' ?  has(self.persistentVolume) : !has(self.persistentVolume)",message="persistentVolume is required when type is PersistentVolume, and forbidden otherwise"
type Storage struct {
	// type is a required field that specifies the type of storage that will be used to store the Insights data archive.
	// Valid values are "PersistentVolume" and "Ephemeral".
	// When set to Ephemeral, the Insights data archive is stored in the ephemeral storage of the gathering job.
	// When set to PersistentVolume, the Insights data archive is stored in the PersistentVolume that is
	// defined by the PersistentVolume field.
	// +required
	Type StorageType `json:"type"`
	// persistentVolume is an optional field that specifies the PersistentVolume that will be used to store the Insights data archive.
	// The PersistentVolume must be created in the openshift-insights namespace.
	// +optional
	PersistentVolume *PersistentVolumeConfig `json:"persistentVolume,omitempty"`
}

// storageType declares valid storage types
// +kubebuilder:validation:Enum=PersistentVolume;Ephemeral
type StorageType string

const (
	// StorageTypePersistentVolume storage type
	StorageTypePersistentVolume StorageType = "PersistentVolume"
	// StorageTypeEphemeral storage type
	StorageTypeEphemeral StorageType = "Ephemeral"
)

// persistentVolumeConfig provides configuration options for PersistentVolume storage.
type PersistentVolumeConfig struct {
	// claim is a required field that specifies the configuration of the PersistentVolumeClaim that will be used to store the Insights data archive.
	// The PersistentVolumeClaim must be created in the openshift-insights namespace.
	// +required
	Claim PersistentVolumeClaimReference `json:"claim"`
	// mountPath is an optional field specifying the directory where the PVC will be mounted inside the Insights data gathering Pod.
	// When omitted, this means no opinion and the platform is left to choose a reasonable default, which is subject to change over time.
	// The current default mount path is /var/lib/insights-operator
	// The path may not exceed 1024 characters and must not contain a colon.
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:XValidation:rule="!self.contains(':')",message="mountPath must not contain a colon"
	// +optional
	MountPath string `json:"mountPath,omitempty"`
}

// persistentVolumeClaimReference is a reference to a PersistentVolumeClaim.
type PersistentVolumeClaimReference struct {
	// name is a string that follows the DNS1123 subdomain format.
	// It must be at most 253 characters in length, and must consist only of lower case alphanumeric characters, '-' and '.', and must start and end with an alphanumeric character.
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	// +required
	Name string `json:"name"`
}

// dataPolicyOption declares valid data policy types
// +kubebuilder:validation:Enum=ObfuscateNetworking;WorkloadNames
type DataPolicyOption string

const (
	// IP addresses and cluster domain name are obfuscated
	DataPolicyOptionObfuscateNetworking DataPolicyOption = "ObfuscateNetworking"
	// Data from Deployment Validation Operator are obfuscated
	DataPolicyOptionObfuscateWorkloadNames DataPolicyOption = "WorkloadNames"
)

// gathererConfig allows to configure specific gatherers
type GathererConfig struct {
	// name is the required name of a specific gatherer
	// It may not exceed 256 characters.
	// The format for a gatherer name is: {gatherer}/{function} where the function is optional.
	// Gatherer consists of a lowercase letters only that may include underscores (_).
	// Function consists of a lowercase letters only that may include underscores (_) and is separated from the gatherer by a forward slash (/).
	// The particular gatherers can be found at https://github.com/openshift/insights-operator/blob/master/docs/gathered-data.md.
	// Run the following command to get the names of last active gatherers:
	// "oc get insightsoperators.operator.openshift.io cluster -o json | jq '.status.gatherStatus.gatherers[].name'"
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:XValidation:rule=`self.matches("^[a-z]+[_a-z]*[a-z]([/a-z][_a-z]*)?[a-z]$")`,message=`gatherer name must be in the format of {gatherer}/{function} where the gatherer and function are lowercase letters only that may include underscores (_) and are separated by a forward slash (/) if the function is provided`
	// +required
	Name string `json:"name"`
	// state is a required field that allows you to configure specific gatherer. Valid values are "Enabled" and "Disabled".
	// When set to Enabled the gatherer will run.
	// When set to Disabled the gatherer will not run.
	// +required
	State GathererState `json:"state"`
}

// state declares valid gatherer state types.
// +kubebuilder:validation:Enum=Enabled;Disabled
type GathererState string

const (
	// GathererStateEnabled gatherer state, which means that the gatherer will run.
	GathererStateEnabled GathererState = "Enabled"
	// GathererStateDisabled gatherer state, which means that the gatherer will not run.
	GathererStateDisabled GathererState = "Disabled"
)

// dataGatherState declares valid gathering state types
// +kubebuilder:validation:Optional
// +kubebuilder:validation:Enum=Running;Completed;Failed;Pending
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Running' && self == 'Pending')", message="dataGatherState cannot transition from Running to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Completed' && self == 'Pending')", message="dataGatherState cannot transition from Completed to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Failed' && self == 'Pending')", message="dataGatherState cannot transition from Failed to Pending"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Completed' && self == 'Running')", message="dataGatherState cannot transition from Completed to Running"
// +kubebuilder:validation:XValidation:rule="!(oldSelf == 'Failed' && self == 'Running')", message="dataGatherState cannot transition from Failed to Running"
type DataGatherState string

const (
	// Data gathering is running
	DataGatherStateRunning DataGatherState = "Running"
	// Data gathering is completed
	DataGatherStateCompleted DataGatherState = "Completed"
	// Data gathering failed
	DataGatherStateFailed DataGatherState = "Failed"
	// Data gathering is pending
	DataGatherStatePending DataGatherState = "Pending"
)

// DataGatherStatus contains information relating to the DataGather state.
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.insightsRequestID) || has(self.insightsRequestID))",message="cannot remove insightsRequestID attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.startTime) || has(self.startTime))",message="cannot remove startTime attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.finishTime) || has(self.finishTime))",message="cannot remove finishTime attribute from status"
// +kubebuilder:validation:XValidation:rule="(!has(oldSelf.dataGatherState) || has(self.dataGatherState))",message="cannot remove dataGatherState attribute from status"
// +kubebuilder:validation:Optional
type DataGatherStatus struct {
	// conditions is an optional field that provides details on the status of the gatherer job.
	// It may not exceed 100 items and must not contain duplicates.
	//
	// The current condition types are DataUploaded, DataRecorded, DataProcessed, RemoteConfigurationNotAvailable, RemoteConfigurationInvalid
	//
	// The DataUploaded condition is used to represent whether or not the archive was successfully uploaded for further processing.
	// When it has a status of True and a reason of HttpStatus200, the archive was successfully uploaded.
	// When it has a status of Unknown and a reason of NoUploadYet, the upload has not occurred yet.
	// When it has a status of False and a reason like HttpStatusXXX, the upload failed and the reason reflects the returned HTTP status code. The accompanying message will include the specific error encountered.
	//
	// The DataRecorded condition is used to represent whether or not the archive was successfully recorded.
	// When it has a status of True and a reason of AsExpected, the archive was recorded successfully.
	// When it has a status of Unknown and a reason of NoDataGatheringYet, the data gathering process has not started yet.
	// When it has a status of False and a reason of RecordingFailed, the recording failed and a message will include the specific error encountered.
	//
	// The DataProcessed condition is used to represent whether or not the archive was processed by the processing service.
	// When it has a status of True and a reason of Processed, the data was processed successfully.
	// When it has a status of Unknown and a reason of NothingToProcessYet, there is no data to process at the moment.
	// When it has a status of False and a reason of Failure, processing failed and a message will include the specific error encountered.
	//
	// The RemoteConfigurationNotAvailable condition is used to represent whether the remote configuration is available.
	// When it has a status of Unknown and a reason of Unknown, the state of the remote configuration is unknown—typically at startup.
	// When is has a satatus of True and a reason of AsExpected, the configuration is available.
	//
	// The RemoteConfigurationInvalid condition is used to represent whether the remote configuration is valid.
	// When it has a status of Unknown and a reason of Unknown, the validity of the remote configuration is unknown—typically at startup.
	// When is has a satatus of True and a reason of AsExpected, the configuration is valid.
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=100
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// dataGatherState reflects the current state of the data gathering process.
	// +optional
	State DataGatherState `json:"dataGatherState,omitempty"`
	// gatherers is a list of active gatherers (and their statuses) in the last gathering.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=100
	// +optional
	Gatherers []GathererStatus `json:"gatherers,omitempty"`
	// startTime is the time when Insights data gathering started.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="startTime is immutable once set"
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// finishTime is the time when Insights data gathering finished.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="finishTime is immutable once set"
	// +optional
	FinishTime metav1.Time `json:"finishTime,omitempty"`
	// relatedObjects is an optional list of resources which are useful when debugging or inspecting the data gathering Pod
	// It may not exceed 100 items and must not contain duplicates.
	// +listType=map
	// +listMapKey=name
	// +listMapKey=namespace
	// +kubebuilder:validation:MaxItems=100
	// +optional
	RelatedObjects []ObjectReference `json:"relatedObjects,omitempty"`
	// insightsRequestID is an optional Insights request ID to track the status of the Insights analysis (in console.redhat.com processing pipeline) for the corresponding Insights data archive.
	// It may not exceed 256 characters and is immutable once set.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="insightsRequestID is immutable once set"
	// +kubebuilder:validation:MaxLength=256
	// +optional
	InsightsRequestID string `json:"insightsRequestID,omitempty"`
	// insightsReport provides general Insights analysis results.
	// When omitted, this means no data gathering has taken place yet or the
	// corresponding Insights analysis (identified by "insightsRequestID") is not available.
	// +optional
	InsightsReport InsightsReport `json:"insightsReport,omitempty"`
}

// gathererStatus represents information about a particular
// data gatherer.
type GathererStatus struct {
	// conditions provide details on the status of each gatherer.
	//
	// The current condition type is DataGathered
	//
	// The DataGathered condition is used to represent whether or not the data was gathered by a gatherer specified by name.
	// When it has a status of True and a reason of GatheredOK, the data has been successfully gathered as expected.
	// When it has a status of False and a reason of NoData, no data was gathered—for example, when the resource is not present in the cluster.
	// When it has a status of False and a reason of GatherError, an error occurred and no data was gathered.
	// When it has a status of False and a reason of GatherPanic, a panic occurred during gathering and no data was collected.
	// When it has a status of False and a reason of GatherWithErrorReason, data was partially gathered or gathered with an error message.
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// name is the required name of the gatherer.
	// It must contain at least 5 characters and may not exceed 256 characters.
	// +required
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:MinLength=5
	Name string `json:"name"`
	// lastGatherDuration represents the time spent gathering.
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^(([0-9]+(?:\\.[0-9]+)?(ns|us|µs|μs|ms|s|m|h))+)$"
	LastGatherDuration metav1.Duration `json:"lastGatherDuration"`
}

// insightsReport provides Insights health check report based on the most
// recently sent Insights data.
type InsightsReport struct {
	// downloadedTime is an optional time when the last Insights report was downloaded.
	// An empty value means that there has not been any Insights report downloaded yet and
	// it usually appears in disconnected clusters (or clusters when the Insights data gathering is disabled).
	// +optional
	DownloadedTime metav1.Time `json:"downloadedTime,omitempty"`
	// healthChecks provides basic information about active Insights health checks
	// in a cluster.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=100
	// +optional
	HealthChecks []HealthCheck `json:"healthChecks,omitempty"`
	// uri is optional field that provides the URL link from which the report was downloaded.
	// The link must be a valid HTTPS URL and the maximum length is 2048 characters.
	// +kubebuilder:validation:XValidation:rule=`isURL(self) && url(self).getScheme() == "https"`,message=`URI must be a valid HTTPS URL (e.g., https://example.com)`
	// +kubebuilder:validation:MaxLength=2048
	// +optional
	URI string `json:"uri,omitempty"`
}

// healthCheck represents an Insights health check attributes.
type HealthCheck struct {
	// description is required field that provides basic description of the healtcheck.
	// It must contain at least 10 characters and may not exceed 2048 characters.
	// +required
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:MinLength=10
	Description string `json:"description"`
	// totalRisk is the required field of the healthcheck.
	// It is indicator of the total risk posed by the detected issue; combination of impact and likelihood.
	// The values can be from 1 to 4, and the higher the number, the more important the issue.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=4
	TotalRisk int32 `json:"totalRisk"`
	// advisorURI is required field that provides the URL link to the Insights Advisor.
	// The link must be a valid HTTPS URL and the maximum length is 2048 characters.
	// +kubebuilder:validation:XValidation:rule=`isURL(self) && url(self).getScheme() == "https"`,message=`advisorURI must be a valid HTTPS URL (e.g., https://example.com)`
	// +kubebuilder:validation:MaxLength=2048
	// +required
	AdvisorURI string `json:"advisorURI"`
	// state is required field that determines what the current state of the health check is.
	// Health check is enabled by default and can be disabled by the user in the Insights advisor user interface.
	// +required
	State HealthCheckState `json:"state"`
}

// healthCheckState provides information about the status of the
// health check (for example, the health check may be marked as disabled by the user).
// +kubebuilder:validation:Enum:=Enabled;Disabled
type HealthCheckState string

const (
	// enabled marks the health check as enabled
	HealthCheckEnabled HealthCheckState = "Enabled"
	// disabled marks the health check as disabled
	HealthCheckDisabled HealthCheckState = "Disabled"
)

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// group is required field that specifies the API Group of the Resource.
	// Enter empty string for the core group.
	// This value is empty or it should follow the DNS1123 subdomain format.
	// It must be at most 253 characters in length, and must consist only of lower case alphanumeric characters, '-' and '.', and must start and end with an alphanumeric character.
	// Example: "", "apps", "build.openshift.io", etc.
	// +kubebuilder:validation:XValidation:rule="self.size() == 0 || !format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength:=253
	// +required
	Group string `json:"group"`
	// resource is required field of the type that is being referenced.
	// It is normally the plural form of the resource kind in lowercase.
	// It must be at most 63 characters in length, and must must consist of only lowercase alphanumeric characters and hyphens
	// Example: "deployments", "deploymentconfigs", "pods", etc.
	// +kubebuilder:validation:XValidation:rule=`!format.dns1123Label().validate(self).hasValue()`,message="the value must consist of only lowercase alphanumeric characters and hyphens"
	// +kubebuilder:validation:MaxLength=63
	// +required
	Resource string `json:"resource"`
	// name is required field that specifies the referent that follows the DNS1123 subdomain format.
	// It must be at most 253 characters in length, and must consist only of lower case alphanumeric characters, '-' and '.', and must start and end with an alphanumeric character.
	// +kubebuilder:validation:XValidation:rule="!format.dns1123Subdomain().validate(self).hasValue()",message="a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
	// +kubebuilder:validation:MaxLength=253
	// +required
	Name string `json:"name"`
	// namespace if required field of the referent that follows the DNS1123 labels format.
	// It must be at most 63 characters in length, and must must consist of only lowercase alphanumeric characters and hyphens
	// +kubebuilder:validation:XValidation:rule=`!format.dns1123Label().validate(self).hasValue()`,message="the value must consist of only lowercase alphanumeric characters and hyphens"
	// +kubebuilder:validation:MaxLength=63
	// +required
	Namespace string `json:"namespace"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataGatherList is a collection of items
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type DataGatherList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// items contains a list of DataGather resources.
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=100
	// +optional
	Items []DataGather `json:"items,omitempty"`
}
