module github.com/openshift/api

go 1.13

require (
	github.com/gogo/protobuf v1.3.1
	github.com/openshift/build-machinery-go v0.0.0-20200731024703-cd7e6e844b55
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200602230032-c00d67ef29d0
	k8s.io/api v0.19.0-rc.2
	k8s.io/apimachinery v0.19.0-rc.2
	k8s.io/code-generator v0.19.0-rc.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-tools v0.3.0
)

replace github.com/openshift/build-machinery-go => github.com/tnozicka/build-machinery-go v0.0.0-20200813151022-40b80b29a377
