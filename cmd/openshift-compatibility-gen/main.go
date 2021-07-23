package main

import (
	"flag"
	"os/exec"
	"path/filepath"

	"github.com/openshift/api/cmd/openshift-compatibility-gen/generate/comment"
	"github.com/openshift/api/cmd/openshift-compatibility-gen/generate/compatibility"
	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := compatibility.NewDefaults()

	// Override defaults.
	genericArgs.GoHeaderFilePath = filepath.Join(defaultSourceTree(), "./hack/boilerplate.go.txt")

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := compatibility.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	// Run it.
	/*	if err := genericArgs.Execute(
			compatibility.NameSystems(),
			compatibility.DefaultNameSystem(),
			compatibility.Packages,
		); err != nil {
			klog.Fatalf("Error: %v", err)
		}
	*/
	if err := comment.GenerateCompatibilityComments(genericArgs.InputDirs); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}

func defaultSourceTree() string {
	gomod, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		klog.Errorln(err)
	}
	if len(gomod) > 0 {
		return filepath.Dir(string(gomod))
	}
	return args.DefaultSourceTree()
}
