module github.com/openshift/api

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.0
	github.com/google/gofuzz v1.0.0
	github.com/json-iterator/go v1.1.8
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/openshift/library-go v0.0.0-20191205152556-73e1fb871a9b
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	golang.org/x/text v0.3.2
	gonum.org/v1/gonum v0.0.0-20190331200053-3d26580ed485
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/code-generator v0.0.0
	k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/klog v1.0.0
	sigs.k8s.io/yaml v1.1.0
)

replace (
	k8s.io/api => k8s.io/api v0.17.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.0
	k8s.io/code-generator => k8s.io/code-generator v0.17.0
)
