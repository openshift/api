package sippy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

type SippyQueryStruct struct {
	Items        []SippyQueryItem `json:"items"`
	LinkOperator string           `json:"linkOperator"`
	TierName     string           `json:"-"` // not serialized, used to track which tier this query is for
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

func QueriesFor(cloud, architecture, topology, networkStack, os, jobTiers, testPattern string) []*SippyQueryStruct {
	// Build base query items that are common to all JobTier queries
	baseItems := []SippyQueryItem{
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
	}

	if networkStack != "" {
		baseItems = append(baseItems, SippyQueryItem{
			ColumnField:   "variants",
			Not:           false,
			OperatorValue: "contains",
			Value:         fmt.Sprintf("NetworkStack:%s", networkStack),
		})
	}

	if os != "" {
		baseItems = append(baseItems, SippyQueryItem{
			ColumnField:   "variants",
			Not:           false,
			OperatorValue: "contains",
			Value:         fmt.Sprintf("OS:%s", os),
		})
	}

	// Parse JobTiers - comma-separated string, default to standard/informing/blocking/candidate if empty
	var jobTiersList []string
	if jobTiers == "" {
		jobTiersList = []string{"standard", "informing", "blocking", "candidate"}
	} else {
		// Split by comma, trim whitespace, and deduplicate using sets
		tierSet := sets.New[string]()
		for _, tier := range strings.Split(jobTiers, ",") {
			if trimmed := strings.TrimSpace(tier); trimmed != "" {
				tierSet.Insert(trimmed)
			}
		}
		// If all tiers were whitespace/empty after trimming, use defaults
		if tierSet.Len() == 0 {
			jobTiersList = []string{"standard", "informing", "blocking", "candidate"}
		} else {
			jobTiersList = sets.List(tierSet)
		}
	}

	// Generate one query per JobTier (to work around API's single LinkOperator limitation)
	var queries []*SippyQueryStruct
	for _, jobTier := range jobTiersList {
		// Copy base items for this query
		items := make([]SippyQueryItem, len(baseItems))
		copy(items, baseItems)

		// Add JobTier filter
		items = append(items, SippyQueryItem{
			ColumnField:   "variants",
			Not:           false,
			OperatorValue: "contains",
			Value:         fmt.Sprintf("JobTier:%s", jobTier),
		})

		queries = append(queries, &SippyQueryStruct{
			Items:        items,
			LinkOperator: "and",
			TierName:     jobTier,
		})
	}

	return queries
}

func BuildSippyTestAnalysisURL(release, testName, topology, cloud, architecture, networkStack, os string) string {
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
	if os != "" {
		filterItems = append(filterItems, SippyQueryItem{
			ColumnField:   "variants",
			OperatorValue: "has entry",
			Value:         fmt.Sprintf("OS:%s", os),
		})
	}
	// Note: We don't filter by JobTier in the URL so the link shows all tiers
	// The actual queries filter by each tier separately and merge results

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
