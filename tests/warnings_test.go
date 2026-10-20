package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestWarningRecorder(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		expected []string
		messages []string
		wantErr  bool
	}{
		{name: "unspecified", messages: []string{"legacy warning"}},
		{name: "explicit empty", expected: []string{}},
		{name: "unexpected warning", expected: []string{}, messages: []string{"unexpected"}, wantErr: true},
		{name: "missing warning", expected: []string{"unknown field"}, wantErr: true},
		{name: "stable substring", expected: []string{`unknown field "spec.canary"`}, messages: []string{`server prefix: unknown field "spec.canary"`}},
		{name: "retry duplicates", expected: []string{"unknown field"}, messages: []string{"unknown field", "unknown field"}},
		{name: "all expectations required", expected: []string{"unknown", "other"}, messages: []string{"unknown field"}, wantErr: true},
		{name: "all messages matched", expected: []string{"unknown"}, messages: []string{"unknown field", "other warning"}, wantErr: true},
		{name: "unordered", expected: []string{"second", "first"}, messages: []string{"first warning", "second warning"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			recorder := &warningRecorder{}
			var workers sync.WaitGroup
			for _, message := range test.messages {
				workers.Go(func() { recorder.HandleWarningHeader(299, "test", message) })
			}
			workers.Wait()
			if err := recorder.check(test.expected); (err != nil) != test.wantErr {
				t.Fatalf("check() = %v, want error %t", err, test.wantErr)
			}
		})
	}
}

func TestExpectedWarningsYAML(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		yaml    string
		wantNil bool
	}{
		{name: "omitted", yaml: "name: legacy", wantNil: true},
		{name: "null", yaml: "expectedWarnings: null", wantNil: true},
		{name: "empty", yaml: "expectedWarnings: []"},
		{name: "substring", yaml: "expectedWarnings: [unknown field]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var create OnCreateTestSpec
			var update OnUpdateTestSpec
			if err := yaml.Unmarshal([]byte(test.yaml), &create); err != nil {
				t.Fatal(err)
			}
			if err := yaml.Unmarshal([]byte(test.yaml), &update); err != nil {
				t.Fatal(err)
			}
			if (create.ExpectedWarnings == nil) != test.wantNil || (update.ExpectedWarnings == nil) != test.wantNil {
				t.Fatalf("nil expectations: create=%v update=%v, want %t", create.ExpectedWarnings == nil, update.ExpectedWarnings == nil, test.wantNil)
			}
		})
	}
}

func TestRequestWarnings(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		method string
		status int
		retry  bool
	}{
		{name: "create", method: http.MethodPost, status: http.StatusCreated},
		{name: "update", method: http.MethodPut, status: http.StatusOK},
		{name: "create error", method: http.MethodPost, status: http.StatusBadRequest},
		{name: "update error", method: http.MethodPut, status: http.StatusBadRequest},
		{name: "create retry", method: http.MethodPost, status: http.StatusCreated, retry: true},
		{name: "update retry", method: http.MethodPut, status: http.StatusOK, retry: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				response.Header().Set("Content-Type", "application/json")
				if request.Method == http.MethodGet {
					response.Header().Add("Warning", `299 test "unrelated read warning"`)
				} else {
					if request.Method != test.method {
						t.Errorf("method = %s, want %s", request.Method, test.method)
					}
					attempts++
					response.Header().Add("Warning", `299 first "unknown field \"spec.canary\"", 299 second "another, quoted warning"`)
					response.Header().Add("Warning", `299 retry "unknown field \"spec.canary\""`)
					if test.retry && attempts == 1 {
						response.Header().Set("Warning", `299 retry "transient retry warning"`)
						response.Header().Set("Retry-After", "0")
						response.WriteHeader(http.StatusTooManyRequests)
						fmt.Fprint(response, `{"kind":"Status","apiVersion":"v1","status":"Failure","reason":"TooManyRequests","code":429}`)
						return
					}
					response.WriteHeader(test.status)
				}
				if test.status == http.StatusBadRequest && request.Method != http.MethodGet {
					fmt.Fprint(response, `{"kind":"Status","apiVersion":"v1","status":"Failure","reason":"BadRequest","code":400}`)
					return
				}
				fmt.Fprint(response, `{"apiVersion":"test.example/v1","kind":"Example","metadata":{"name":"sample","resourceVersion":"1"}}`)
			}))
			defer server.Close()
			groupVersion := schema.GroupVersion{Group: "test.example", Version: "v1"}
			mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{groupVersion})
			mapper.Add(groupVersion.WithKind("Example"), meta.RESTScopeRoot)
			configuration := &rest.Config{Host: server.URL, WarningHandler: rest.NoWarnings{}, WarningHandlerWithContext: rest.NoWarnings{}}
			baseClient, err := client.New(configuration, client.Options{Scheme: runtime.NewScheme(), Mapper: mapper})
			if err != nil {
				t.Fatal(err)
			}
			legacyClient, legacyRecorder, err := clientForWarnings(configuration, baseClient, nil)
			if err != nil || legacyClient != baseClient || legacyRecorder != nil || legacyRecorder.check(nil) != nil {
				t.Fatalf("legacy client behavior changed: %v", err)
			}
			expected := []string{`unknown field "spec.canary"`, "another, quoted warning"}
			requestClient, recorder, err := clientForWarnings(configuration, baseClient, expected)
			if err != nil {
				t.Fatal(err)
			}
			if _, unchanged := configuration.WarningHandler.(rest.NoWarnings); !unchanged {
				t.Fatal("base warning handler changed")
			}
			if _, unchanged := configuration.WarningHandlerWithContext.(rest.NoWarnings); !unchanged {
				t.Fatal("base contextual warning handler changed")
			}
			object := &unstructured.Unstructured{Object: map[string]interface{}{
				"apiVersion": groupVersion.String(), "kind": "Example", "metadata": map[string]interface{}{"name": "sample", "resourceVersion": "1"},
			}}
			if test.method == http.MethodPost {
				err = requestClient.Create(context.Background(), object)
			} else {
				err = requestClient.Update(context.Background(), object)
			}
			if (err != nil) != (test.status == http.StatusBadRequest) {
				t.Fatalf("request error = %v, status %d", err, test.status)
			}
			if err := recorder.check(expected); err != nil {
				t.Fatal(err)
			}
			if test.retry && attempts != 2 {
				t.Fatalf("attempts = %d, want 2", attempts)
			}
			if err := baseClient.Get(context.Background(), client.ObjectKey{Name: "sample"}, object); err != nil {
				t.Fatal(err)
			}
			if err := recorder.check(expected); err != nil {
				t.Fatalf("read leaked warnings: %v", err)
			}
			_, nextRecorder, err := clientForWarnings(configuration, baseClient, []string{})
			if err != nil {
				t.Fatal(err)
			}
			if err := nextRecorder.check([]string{}); err != nil {
				t.Fatalf("warnings leaked between requests: %v", err)
			}
		})
	}
}
