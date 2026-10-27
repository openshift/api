package v1alpha1

import (
	"encoding/json"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerResourceQuantityPresence(t *testing.T) {
	zero := resource.MustParse("0")
	for _, tc := range []struct {
		name    string
		request *resource.Quantity
		limit   *resource.Quantity
	}{
		{name: "omitted"},
		{name: "zero request", request: &zero},
		{name: "zero limit", limit: &zero},
		{name: "both zero", request: &zero, limit: &zero},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := ContainerResource{Name: "cpu", Request: tc.request, Limit: tc.limit}
			data, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			for name, quantity := range map[string]*resource.Quantity{"request": tc.request, "limit": tc.limit} {
				value, present := fields[name]
				if present != (quantity != nil) {
					t.Fatalf("%s: present=%t, want %t", name, present, quantity != nil)
				}
				if quantity != nil && string(value) != `"0"` {
					t.Fatalf("%s: got %s, want explicit zero", name, value)
				}
			}
			var roundTrip ContainerResource
			if err := json.Unmarshal(data, &roundTrip); err != nil {
				t.Fatal(err)
			}
			copied := roundTrip.DeepCopy()
			if (copied.Request != nil) != (tc.request != nil) || (copied.Limit != nil) != (tc.limit != nil) {
				t.Fatal("round trip or DeepCopy changed field presence")
			}
			if copied.Request != nil {
				copied.Request.Set(1)
				if !roundTrip.Request.IsZero() {
					t.Fatal("DeepCopy aliases request")
				}
			}
			if copied.Limit != nil {
				copied.Limit.Set(1)
				if !roundTrip.Limit.IsZero() {
					t.Fatal("DeepCopy aliases limit")
				}
			}
		})
	}
}
