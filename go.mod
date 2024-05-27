module github.com/openshift/api

go 1.22.0

toolchain go1.22.1

require (
	github.com/gogo/protobuf v1.3.2
	golang.org/x/tools v0.18.0
	k8s.io/api v0.30.1
	k8s.io/apimachinery v0.30.1
	k8s.io/klog/v2 v2.120.1
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)

// v3.9.0 is the only tag in openshift/api and it was created before go.mod was
// introduced. We retract it so that go command don't select it automatically
// when resolving versions like @latest.
retract v3.9.0+incompatible

// To make go aware of the retraction, we need to tag a new version that can be
// retracted by itself.
retract v0.0.1
