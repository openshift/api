module github.com/openshift/api/tools

go 1.20

require (
	github.com/dave/dst v0.27.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-git/v5 v5.8.1
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.9
	github.com/mikefarah/yq/v4 v4.35.1
	// Tagged versions produce an incorrect diff, see https://github.com/sergi/go-diff/issues/123, this is the merge for the fix.
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/vmware-archive/yaml-patch v0.0.11
	golang.org/x/tools v0.12.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apiextensions-apiserver v0.27.4
	k8s.io/apimachinery v0.27.4
	k8s.io/code-generator v0.27.4
	k8s.io/gengo v0.0.0-20230306165830-ab3349d207d4
	k8s.io/klog/v2 v2.100.1
	k8s.io/kube-openapi v0.0.0-20230718181711-3c0fae5ee9fd
	sigs.k8s.io/controller-tools v0.12.1
	sigs.k8s.io/yaml v1.3.0
)

replace sigs.k8s.io/controller-tools => github.com/openshift/controller-tools v0.11.2-0.20230707132950-4180f051f655

require (
	github.com/a8m/envsubst v1.4.2 // indirect
	github.com/alecthomas/participle/v2 v2.0.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/elliotchance/orderedmap v1.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/gobuffalo/flect v1.0.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.11.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/krishicks/yaml-patch v0.0.10 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.27.4 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.3.0 // indirect
)
