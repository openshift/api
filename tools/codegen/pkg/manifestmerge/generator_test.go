package manifestmerge

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"reflect"
	kyaml "sigs.k8s.io/yaml"
	"testing"
)

func Test_mergeCRD2(t *testing.T) {
	type args struct {
		obj   *unstructured.Unstructured
		patch []byte
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "version-overlay",
			args: args{
				obj:   readCRDYamlOrDie2(t, "testdata/field-overlay/01.yaml"),
				patch: readFileOrDie(t, "testdata/field-overlay/02.yaml"),
			},
		},
		{
			name: "version-invert",
			args: args{
				obj:   readCRDYamlOrDie2(t, "testdata/field-overlay/02.yaml"),
				patch: readFileOrDie(t, "testdata/field-overlay/01.yaml"),
			},
		},
		{
			name: "cluster-version-signature-store",
			args: args{
				obj:   readCRDYamlOrDie2(t, "/home/deads/workspaces/api/src/github.com/openshift/api/config/v1/crd-manifests-by-featuregate/clusterversions/aaa_no-gates.yaml"),
				patch: readFileOrDie(t, "/home/deads/workspaces/api/src/github.com/openshift/api/config/v1/crd-manifests-by-featuregate/clusterversions/SignatureStores.yaml"),
			},
		},
		{
			name: "cluster-version-signature-store-accidental-removal?",
			args: args{
				obj:   readCRDYamlOrDie2(t, "/home/deads/workspaces/api/src/github.com/openshift/api/config/v1/crd-manifests-by-featuregate/clusterversions/SignatureStores.yaml"),
				patch: readFileOrDie(t, "/home/deads/workspaces/api/src/github.com/openshift/api/config/v1/crd-manifests-by-featuregate/clusterversions/SigstoreImageVerification.yaml"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeCRD(tt.args.obj, tt.args.patch); !reflect.DeepEqual(got, tt.want) {
				outBytes, err := kyaml.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(string(outBytes))
				t.Errorf("mergeCRD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readCRDYamlOrDie2(t *testing.T, path string) *unstructured.Unstructured {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := readCRDYaml(data)
	if err != nil {
		t.Fatal(err)
	}
	return ret
}

func readFileOrDie(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
