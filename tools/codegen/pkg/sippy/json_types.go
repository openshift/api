package sippy

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type SippyQueryStruct struct {
	Items        []SippyQueryItem `json:"items"`
	LinkOperator string           `json:"linkOperator"`
}

type SippyQueryItem struct {
	ColumnField   string `json:"columnField"`
	Not           bool   `json:"not"`
	OperatorValue string `json:"operatorValue"`
	Value         string `json:"value"`
}

type SippyTestInfo struct {
	Id                        int         `json:"id"`
	Name                      string      `json:"name"`
	SuiteName                 string      `json:"suite_name"`
	Variants                  interface{} `json:"variants"`
	JiraComponent             string      `json:"jira_component"`
	JiraComponentId           int         `json:"jira_component_id"`
	CurrentSuccesses          int         `json:"current_successes"`
	CurrentFailures           int         `json:"current_failures"`
	CurrentFlakes             int         `json:"current_flakes"`
	CurrentPassPercentage     float64     `json:"current_pass_percentage"`
	CurrentFailurePercentage  float64     `json:"current_failure_percentage"`
	CurrentFlakePercentage    float64     `json:"current_flake_percentage"`
	CurrentWorkingPercentage  float64     `json:"current_working_percentage"`
	CurrentRuns               int         `json:"current_runs"`
	PreviousSuccesses         int         `json:"previous_successes"`
	PreviousFailures          int         `json:"previous_failures"`
	PreviousFlakes            int         `json:"previous_flakes"`
	PreviousPassPercentage    float64     `json:"previous_pass_percentage"`
	PreviousFailurePercentage float64     `json:"previous_failure_percentage"`
	PreviousFlakePercentage   float64     `json:"previous_flake_percentage"`
	PreviousWorkingPercentage float64     `json:"previous_working_percentage"`
	PreviousRuns              int         `json:"previous_runs"`
	NetFailureImprovement     float64     `json:"net_failure_improvement"`
	NetFlakeImprovement       float64     `json:"net_flake_improvement"`
	NetWorkingImprovement     float64     `json:"net_working_improvement"`
	NetImprovement            float64     `json:"net_improvement"`
	Watchlist                 bool        `json:"watchlist"`
	Tags                      interface{} `json:"tags"`
	OpenBugs                  int         `json:"open_bugs"`
}

func QueriesFor(cloud, architecture, topology, networkStack, testPattern string) []*SippyQueryStruct {
	queries := []*SippyQueryStruct{
		{
			Items: []SippyQueryItem{
				{
					ColumnField:   "variants",
					Not:           false,
					OperatorValue: "contains",
					Value:         fmt.Sprintf("Platform:%s", cloud),
				},
				{
					ColumnField:   "variants",
					Not:           false,
					OperatorValue: "contains",
					Value:         fmt.Sprintf("Architecture:%s", architecture),
				},
				{
					ColumnField:   "variants",
					Not:           false,
					OperatorValue: "contains",
					Value:         fmt.Sprintf("Topology:%s", topology),
				},
				{
					ColumnField:   "name",
					Not:           false,
					OperatorValue: "contains",
					Value:         testPattern,
				},
			},
			LinkOperator: "and",
		},
	}

	if networkStack != "" {
		queries[0].Items = append(queries[0].Items,
			SippyQueryItem{
				ColumnField:   "variants",
				Not:           false,
				OperatorValue: "contains",
				Value:         fmt.Sprintf("NetworkStack:%s", networkStack),
			},
		)

	}

	return queries
}

func BuildSippyTestAnalysisURL(release, testName, topology, cloud, architecture, networkStack string) string {
	filterItems := []SippyQueryItem{
		{
			ColumnField:   "name",
			OperatorValue: "equals",
			Value:         testName,
		},
		{
			ColumnField:   "variants",
			Not:           true,
			OperatorValue: "has entry",
			Value:         "never-stable",
		},
		{
			ColumnField:   "variants",
			Not:           true,
			OperatorValue: "has entry",
			Value:         "aggregated",
		},
		{
			ColumnField:   "variants",
			OperatorValue: "has entry",
			Value:         fmt.Sprintf("Topology:%s", topology),
		},
		{
			ColumnField:   "variants",
			OperatorValue: "has entry",
			Value:         fmt.Sprintf("Platform:%s", cloud),
		},
		{
			ColumnField:   "variants",
			OperatorValue: "has entry",
			Value:         fmt.Sprintf("Architecture:%s", architecture),
		},
	}
	if networkStack != "" {
		filterItems = append(filterItems, SippyQueryItem{
			ColumnField:   "variants",
			OperatorValue: "has entry",
			Value:         fmt.Sprintf("NetworkStack:%s", networkStack),
		})
	}

	filterObj := SippyQueryStruct{
		Items:        filterItems,
		LinkOperator: "and",
	}
	filterJSON, err := json.Marshal(filterObj)
	if err != nil {
		return ""
	}

	u := &url.URL{
		Scheme: "https",
		Host:   "sippy.dptools.openshift.org",
		Path:   fmt.Sprintf("sippy-ng/tests/%s/analysis", release),
	}
	q := u.Query()
	q.Set("filters", string(filterJSON))
	q.Set("test", testName)
	u.RawQuery = q.Encode()

	return u.String()
}
