package compatibility

import (
	"path"
	"reflect"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// isAPIType indicates whether or not a type could be used to serve an API.  That means, "does it have TypeMeta".
// This doesn't mean the type is served, but we will handle all TypeMeta types.
func isAPIType(t *types.Type) bool {
	// Filter out private types.
	if namer.IsPrivateGoName(t.Name.Name) {
		return false
	}

	if t.Kind != types.Struct {
		return false
	}

	for _, currMember := range t.Members {
		if currMember.Embedded && currMember.Name == "TypeMeta" {
			return true
		}
	}

	if t.Kind == types.Alias {
		return isAPIType(t.Underlying)
	}

	return false
}

// BoilerplatePath uses the boilerplate in code-generator by calculating the relative path to it.
func BoilerplatePath() string {
	return path.Join(reflect.TypeOf(empty{}).PkgPath(), "/../../hack/boilerplate.go.txt")
}

type empty struct{}
