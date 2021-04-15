// Package v1 is an api version in the apiserver.openshift.io group
package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster"
// +kubebuilder:subresource:status
// +genclient:nonNamespaced

// APIRequestCount tracts requests made to a deprecated API. The instance name should
// be of the form `resource.version.group`, matching the deprecated resource.
type APIRequestCount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec defines the characteristics of the resource.
	// +kubebuilder:validation:Required
	// +required
	Spec APIRequestCountSpec `json:"spec"`

	// status contains the observed state of the resource.
	Status APIRequestCountStatus `json:"status,omitempty"`
}

type APIRequestCountSpec struct {
	// removedRelease is when the API will be removed.
	// +kubebuilder:validation:Pattern=^[0-9][0-9]*\.[0-9][0-9]*$
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=64
	// +required
	RemovedRelease string `json:"removedRelease"`
}

// +k8s:deepcopy-gen=true
type APIRequestCountStatus struct {

	// conditions contains details of the current status of this API Resource.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions"`

	// requestsLastHour contains request history for the current hour. This is porcelain to make the API
	// easier to read by humans seeing if they addressed a problem. This field is reset on the hour.
	RequestsLastHour PerResourceAPIRequestLog `json:"requestsLastHour"`

	// requestsLast24h contains request history for the last 24 hours, indexed by the hour, so
	// 12:00AM-12:59 is in index 0, 6am-6:59am is index 6, etc. The index of the current hour
	// is updated live and then duplicated into the requestsLastHour field.
	RequestsLast24h []PerResourceAPIRequestLog `json:"requestsLast24h"`
}

// PerResourceAPIRequestLog logs request for various nodes.
type PerResourceAPIRequestLog struct {

	// nodes contains logs of requests per node.
	Nodes []PerNodeAPIRequestLog `json:"nodes"`
}

// PerNodeAPIRequestLog contains logs of requests to a certain node.
type PerNodeAPIRequestLog struct {

	// nodeName where the request are being handled.
	NodeName string `json:"nodeName"`

	// lastUpdate should *always* being within the hour this is for.  This is a time indicating
	// the last moment the server is recording for, not the actual update time.
	LastUpdate metav1.Time `json:"lastUpdate"`

	// users contains request details by top 10 users. Note that because in the case of an apiserver
	// restart the list of top 10 users is determined on a best-effort basis, the list might be imprecise.
	Users []PerUserAPIRequestCount `json:"users"`
}

// PerUserAPIRequestCount contains logs of a user's requests.
type PerUserAPIRequestCount struct {

	// userName that made the request.
	UserName string `json:"username"`

	// count of requests.
	Count int `json:"count"`

	// requests details by verb.
	Requests []PerVerbAPIRequestCount `json:"requests"`
}

// PerVerbAPIRequestCount counts requests by API request verb.
type PerVerbAPIRequestCount struct {

	// verb of API request (get, list, create, etc...)
	Verb string `json:"verb"`

	// count of requests for verb.
	Count int `json:"count"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIRequestCountList is a list of APIRequestCount resources.
type APIRequestCountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []APIRequestCount `json:"items"`
}
