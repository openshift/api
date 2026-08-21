package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openshift/api/payload-command/upstreamgates"
)

// TODO: look into adding group-resources for feature gates

func main() {
	o := &upstreamgates.UpstreamGatesOptions{}
	o.AddFlags(flag.CommandLine)
	flag.Parse()

	if err := o.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}
