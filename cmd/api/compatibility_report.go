package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
)

type OpenshiftCompatibilityReport struct {
	Statements []OpenShiftCompatibilityStatement
}

func NewOpenShiftCompatibilityReport(scheme *runtime.Scheme) *OpenshiftCompatibilityReport {
	report := &OpenshiftCompatibilityReport{}
	for gvk := range scheme.AllKnownTypes() {
		report.Statements = append(report.Statements, NewOpenShiftCompatibilityStatement(gvk))
	}
	sort.Slice(report.Statements, func(i, j int) bool {
		c := strings.Compare(report.Statements[i].GVK.Group, report.Statements[j].GVK.Group)
		if c != 0 {
			return c < 0
		}
		c = strings.Compare(report.Statements[i].GVK.Version, report.Statements[j].GVK.Version)
		if c != 0 {
			return c < 0
		}
		return strings.Compare(report.Statements[i].GVK.Kind, report.Statements[j].GVK.Kind) < 0
	})
	return report
}

func (r *OpenshiftCompatibilityReport) String() string {
	var buf bytes.Buffer
	for _, s := range r.Statements {
		buf.WriteString(fmt.Sprintf("%-7s %-40s %s\n", s.compatibility, s.GVK.Kind, s.GVK.GroupVersion()))
	}
	return buf.String()

}
