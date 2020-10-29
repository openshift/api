package main

import (
	"fmt"
	"os"

	"github.com/openshift/api"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

func main() {
	scheme := runtime.NewScheme()
	if err := api.Install(scheme); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	if err := api.InstallKube(scheme); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	report := NewOpenShiftCompatibilityReport(scheme)
	fmt.Println(report)
}
