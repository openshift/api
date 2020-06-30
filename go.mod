module github.com/openshift/api

go 1.13

require (
	github.com/gogo/protobuf v1.3.1
	github.com/openshift/build-machinery-go v0.0.0-20200424080330-082bf86082cc
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200115044656-831fdb1e1868
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/code-generator v0.18.3
	k8s.io/klog v1.0.0
)

replace k8s.io/apimachinery => github.com/damemi/apimachinery v0.19.0-beta.2.0.20200630191726-2cf00ecaa7ef
