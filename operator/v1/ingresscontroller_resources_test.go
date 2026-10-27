package v1

import (
	"encoding/json"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"
)

func TestIngressControllerResourcesRoundTrip(t *testing.T) {
	t.Parallel()
	for _, target := range []string{"router", "haproxy", "both"} {
		t.Run(target, func(t *testing.T) {
			t.Parallel()
			configuration := &IngressControllerResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("500m")},
				Limits:   corev1.ResourceList{"qat.intel.com/cy": resource.MustParse("1")},
			}
			original := IngressController{Spec: IngressControllerSpec{Resources: &IngressControllerResources{}}}
			if target == "router" || target == "both" {
				original.Spec.Resources.Router = configuration.DeepCopy()
			}
			if target == "haproxy" || target == "both" {
				original.Spec.Resources.HAProxy = configuration.DeepCopy()
			}
			for _, format := range []string{"json", "yaml"} {
				t.Run(format, func(t *testing.T) {
					var data []byte
					var err error
					if format == "json" {
						data, err = json.Marshal(original)
					} else {
						data, err = yaml.Marshal(original)
					}
					if err != nil {
						t.Fatal(err)
					}
					var decoded IngressController
					if format == "json" {
						err = json.Unmarshal(data, &decoded)
					} else {
						err = yaml.Unmarshal(data, &decoded)
					}
					if err != nil {
						t.Fatal(err)
					}
					if !equality.Semantic.DeepEqual(original.Spec.Resources, decoded.Spec.Resources) {
						t.Fatalf("resources changed across %s round trip: %#v", format, decoded.Spec.Resources)
					}
				})
			}
		})
	}
}

func TestIngressControllerResourcesDeepCopy(t *testing.T) {
	t.Parallel()
	configuration := &IngressControllerResourceRequirements{
		Requests: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("1.234567890123")},
		Limits:   corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("1Gi")},
	}
	original := &IngressController{Spec: IngressControllerSpec{Resources: &IngressControllerResources{
		Router: configuration.DeepCopy(), HAProxy: configuration.DeepCopy(),
	}}}
	copied := original.DeepCopy()
	if copied.Spec.Resources == original.Spec.Resources || copied.Spec.Resources.Router == original.Spec.Resources.Router || copied.Spec.Resources.HAProxy == original.Spec.Resources.HAProxy {
		t.Fatal("deepcopy reused resource configuration pointers")
	}
	for _, component := range []*IngressControllerResourceRequirements{copied.Spec.Resources.Router, copied.Spec.Resources.HAProxy} {
		quantity := component.Requests[corev1.ResourceCPU]
		quantity.Add(resource.MustParse("1"))
		component.Requests[corev1.ResourceCPU] = quantity
		component.Limits[corev1.ResourceMemory] = resource.MustParse("2Gi")
	}
	for _, component := range []*IngressControllerResourceRequirements{original.Spec.Resources.Router, original.Spec.Resources.HAProxy} {
		if !equality.Semantic.DeepEqual(component, configuration) {
			t.Fatal("deepcopy mutation changed original quantities or claims")
		}
	}
}

func TestIngressControllerResourcesOmission(t *testing.T) {
	t.Parallel()
	for _, input := range []string{`{}`, `{"spec":{}}`, `{"spec":{"resources":null}}`} {
		var object IngressController
		if err := json.Unmarshal([]byte(input), &object); err != nil {
			t.Fatal(err)
		}
		if object.Spec.Resources != nil || object.DeepCopy().Spec.Resources != nil {
			t.Fatal("omission materialized a resource configuration")
		}
	}
}
