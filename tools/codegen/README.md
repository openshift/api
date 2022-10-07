# Codegen

This tool simplifies the process of running API generation across a large number of API groups and versions.
It's aim is to encode specific knowledge about the API generators we use and make generating OpenShift APIs
easier.

## Usage

The tool can be compiled in the normal way using either `go build` or `make codegen` from the root of the tools module.
When using make the tool will output to `$(TOOLS_MODULE)/_output/bin/$(GOOS)/$(GOARCH)/bin/codegen`.

The function has the following arguments:
- `--api-group-versions` - A comma separated list of group versions to generate, all groups must be fully qualified.
  e.g. apps.openshift.io/v1,machine.openshift.io/v1,machine.openshift.io/v1beta1.
- `--base-dir` - The path to the root of the API folders, this directory will be recursively searched to find the group
  versions specified in `--api-group-versions`. When no group versions are specified, all discovered group versions
  will be generated.
- `--controller-gen` - optionally use a particular controller-gen binary. When not specified, the tool will use the
  built in generator.
  Note, you must use a `controller-gen` built from the [OpenShift fork](https://github.com/openshift/kubernetes-sigs-controller-tools) of controller-tools.
- `--required-feature-sets` - optionally generate based on the OpenShift feature sets annotations.
  This will update only CRDs with a matching value for the `release.openshift.io/feature-set` annotation.

As a full example, from the root of the OpenShift API repository, you may run:
```
codegen --base-dir /go/src/github.com/openshift/api --api-group-versions apps.openshift.io/v1,config.openshift.io/v1,operator.openshift.io/v1
```

And to generate only TechPreviewNoUpgrade versions of CRDs:
```
codegen-crds --base-dir /go/src/github.com/openshift/api --api-group-versions apps.openshift.io/v1,config.openshift.io/v1,operator.openshift.io/v1 --require-feature-sets TechPreviewNoUpgrade
```

## Inclusion in other repositories

To use this tool in another repository, you should make sure to add the tools Go submodule to your dependency magnet
or `tools.go` file.

```go
//go:build tools
// +build tools

package tools

import (
  _ "github.com/openshift/api/tools"
  _ "github.com/openshift/api/tools/codegen/cmd"
)
```

You will also need to replace controller-tools with the OpenShift version within your `go.mod`:
```go
replace sigs.k8s.io/controller-tools => github.com/openshift/controller-tools v0.9.3-0.20220912174723-cf3ef054f3dd // v0.9.2+openshift-0.2
```

Then ensure your vendored dependencies are up to date:

```bash
go mod tidy
go mod vendor
```

In your top level Makefile, you can add a target like below, substituting the path to your own APIs base folder and
group versions:
```Make
.PHONY: update-codegen
update-codegen:
  make -C vendor/github.com/openshift/api/tools run-codegen  BASE_DIR="${PWD}/pkg/apis" API_GROUP_VERSIONS="autoscaling.openshift.io/v1,autoscaling.openshift.io/v1beta1"
```

To generate the same group versions but with the TechPreviewNoUpgrade FeatureSet, you would add the FeatureSet to the
end:
```Make
.PHONY: update-codegen-crds
update-codegen-crds:
  make -C vendor/github.com/openshift/api/tools run-codegen-crds  BASE_DIR="${PWD}/pkg/apis" API_GROUP_VERSIONS="autoscaling.openshift.io/v1,autoscaling.openshift.io/v1beta1" OPENSHIFT_REQUIRED_FEATURESETS="TechPreviewNoUpgrade"
```

You may also want to add the `_output` directory to your `.gitignore` to avoid checking in compiled binaries created
by this make target.
