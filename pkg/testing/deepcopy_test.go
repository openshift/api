package testing

import (
	"math/rand"
	"testing"

	apps "github.com/openshift/api/apps/v1"
	authorization "github.com/openshift/api/authorization/v1"
	build "github.com/openshift/api/build/v1"
	image "github.com/openshift/api/image/v1"
	network "github.com/openshift/api/network/v1"
	oauth "github.com/openshift/api/oauth/v1"
	project "github.com/openshift/api/project/v1"
	quota "github.com/openshift/api/quota/v1"
	route "github.com/openshift/api/route/v1"
	security "github.com/openshift/api/security/v1"
	template "github.com/openshift/api/template/v1"
	user "github.com/openshift/api/user/v1"

	"k8s.io/apimachinery/pkg/api/testing/fuzzer"
	"k8s.io/apimachinery/pkg/api/testing/roundtrip"
	genericfuzzer "k8s.io/apimachinery/pkg/apis/meta/fuzzer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
)

var groups = map[schema.GroupVersion]runtime.SchemeBuilder{
	apps.SchemeGroupVersion:          apps.SchemeBuilder,
	authorization.SchemeGroupVersion: authorization.SchemeBuilder,
	build.SchemeGroupVersion:         build.SchemeBuilder,
	image.SchemeGroupVersion:         image.SchemeBuilder,
	network.SchemeGroupVersion:       network.SchemeBuilder,
	oauth.SchemeGroupVersion:         oauth.SchemeBuilder,
	project.SchemeGroupVersion:       project.SchemeBuilder,
	quota.SchemeGroupVersion:         quota.SchemeBuilder,
	route.SchemeGroupVersion:         route.SchemeBuilder,
	security.SchemeGroupVersion:      security.SchemeBuilder,
	template.SchemeGroupVersion:      template.SchemeBuilder,
	user.SchemeGroupVersion:          user.SchemeBuilder,
}

// TestRoundTripTypesWithoutProtobuf applies the round-trip test to all round-trippable Kinds
// in the scheme.
func TestRoundTripTypesWithoutProtobuf(t *testing.T) {
	// TODO upstream this loop
	for gv, builder := range groups {
		t.Logf("starting group %q", gv)
		scheme := runtime.NewScheme()
		codecs := serializer.NewCodecFactory(scheme)

		builder.AddToScheme(scheme)
		seed := rand.Int63()
		// I'm only using the generic fuzzer funcs, but at some point in time we might need to
		// switch to specialized. For now we're happy with the current serialization test.
		fuzzer := fuzzer.FuzzerFor(genericfuzzer.Funcs, rand.NewSource(seed), codecs)
		kinds := scheme.KnownTypes(gv)

		for kind := range kinds {
			if globalNonRoundTrippableTypes.Has(kind) {
				continue
			}
			gvk := gv.WithKind(kind)
			// FIXME: this is explicitly testing w/o protobuf which was failing
			// The RoundTripSpecificKindWithoutProtobuf performs the following sequence of actions:
			// 1. object := original.DeepCopyObject()
			// 2. obj3 := Decode(Encode(object)
			// 3. Fuzz(object)
			// 4. check semantic.DeepEqual(original, obj3)
			roundtrip.RoundTripSpecificKindWithoutProtobuf(t, gvk, scheme, codecs, fuzzer, nil)
		}
		t.Logf("finished group %q", gv)
	}
}

// TODO: externalize this upstream
var globalNonRoundTrippableTypes = sets.NewString(
	"ExportOptions",
	"GetOptions",
	// WatchEvent does not include kind and version and can only be deserialized
	// implicitly (if the caller expects the specific object). The watch call defines
	// the schema by content type, rather than via kind/version included in each
	// object.
	"WatchEvent",
	// ListOptions is now part of the meta group
	"ListOptions",
	// Delete options is only read in metav1
	"DeleteOptions",
)
