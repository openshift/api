package extractjsontags

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// StructFieldTags is used to find information about
// json tags on fields within struct.
type StructFieldTags interface {
	FieldTags(*ast.StructType, string) FieldTagInfo
}

type structFieldTags struct {
	structToFieldTags map[*ast.StructType]map[string]FieldTagInfo
}

func newStructFieldTags() StructFieldTags {
	return &structFieldTags{
		structToFieldTags: make(map[*ast.StructType]map[string]FieldTagInfo),
	}
}

func (s *structFieldTags) insertFieldTagInfo(styp *ast.StructType, field string, tagInfo FieldTagInfo) {
	if s.structToFieldTags[styp] == nil {
		s.structToFieldTags[styp] = make(map[string]FieldTagInfo)
	}

	s.structToFieldTags[styp][field] = tagInfo
}

// FieldTags find the tag information for the named field within the given struct.
func (s *structFieldTags) FieldTags(styp *ast.StructType, field string) FieldTagInfo {
	structFields := s.structToFieldTags[styp]

	if structFields != nil {
		return structFields[field]
	}

	return FieldTagInfo{}
}

// Analyzer is the analyzer for the jsontags package.
// It checks that all struct fields in an API are tagged with json tags.
var Analyzer = &analysis.Analyzer{
	Name:       "extractjsontags",
	Doc:        "Iterates over all fields in structs and extracts their json tags.",
	Run:        run,
	Requires:   []*analysis.Analyzer{inspect.Analyzer},
	ResultType: reflect.TypeOf(newStructFieldTags()),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	results := newStructFieldTags().(*structFieldTags)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		s := n.(*ast.StructType)
		styp, ok := pass.TypesInfo.Types[s].Type.(*types.Struct)
		// Type information may be incomplete.
		if !ok {
			return
		}

		for i := 0; i < styp.NumFields(); i++ {
			field := styp.Field(i)
			tag := styp.Tag(i)

			results.insertFieldTagInfo(s, field.Name(), extractTagInfo(tag))
		}
	})

	return results, nil
}

func extractTagInfo(tag string) FieldTagInfo {
	tagValue, ok := reflect.StructTag(tag).Lookup("json")
	if !ok {
		return FieldTagInfo{Missing: true}
	}

	if tagValue == "" {
		return FieldTagInfo{}
	}

	tagValues := strings.Split(tagValue, ",")

	if len(tagValues) == 2 && tagValues[0] == "" && tagValues[1] == "inline" {
		return FieldTagInfo{Inline: true}
	}

	tagName := tagValues[0]
	return FieldTagInfo{Name: tagName, OmitEmpty: len(tagValues) == 2 && tagValues[1] == "omitempty"}
}

// FieldTagInfo contains information about a field's json tag.
// This is used to pass information about a field's json tag between analyzers.
type FieldTagInfo struct {
	// Name is the name of the field extracted from the json tag.
	Name string

	// OmitEmpty is true if the field has the omitempty option in the json tag.
	OmitEmpty bool

	// Inline is true if the field has the inline option in the json tag.
	Inline bool

	// Missing is true when the field had no json tag.
	Missing bool
}
