package jsontags

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"regexp"

	"github.com/JoelSpeed/kal/pkg/analysis/helpers/extractjsontags"
	"github.com/JoelSpeed/kal/pkg/config"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// camelCaseRegex is a regular expression that matches camel case strings.
	camelCaseRegex = "^[a-z][a-z0-9]*(?:[A-Z][a-z0-9]*)*$"

	name = "jsontags"
)

var (
	errCouldNotGetInspector = errors.New("could not get inspector")
	errCouldNotGetJSONTags  = errors.New("could not get json tags")
)

type analyzer struct {
	jsonTagRegex *regexp.Regexp
}

// newAnalyzer creates a new analyzer with the given json tag regex.
func newAnalyzer(cfg config.JSONTagsConfig) (*analysis.Analyzer, error) {
	defaultConfig(&cfg)

	jsonTagRegex, err := regexp.Compile(cfg.JSONTagRegex)
	if err != nil {
		return nil, fmt.Errorf("could not compile json tag regex: %w", err)
	}

	a := &analyzer{
		jsonTagRegex: jsonTagRegex,
	}

	return &analysis.Analyzer{
		Name:     name,
		Doc:      "Check that all struct fields in an API are tagged with json tags",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, extractjsontags.Analyzer},
	}, nil
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errCouldNotGetInspector
	}

	jsonTags, ok := pass.ResultOf[extractjsontags.Analyzer].(extractjsontags.StructFieldTags)
	if !ok {
		return nil, errCouldNotGetJSONTags
	}

	// Filter to structs so that we can iterate over fields in a struct.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		s, ok := n.(*ast.StructType)
		if !ok {
			return
		}

		styp, ok := pass.TypesInfo.Types[s].Type.(*types.Struct)
		// Type information may be incomplete.
		if !ok {
			return
		}

		for i := 0; i < styp.NumFields(); i++ {
			field := styp.Field(i)

			a.checkField(pass, s, field, jsonTags)
		}
	})

	return nil, nil //nolint:nilnil
}

func (a *analyzer) checkField(pass *analysis.Pass, sTyp *ast.StructType, field *types.Var, jsonTags extractjsontags.StructFieldTags) {
	tagInfo := jsonTags.FieldTags(sTyp, field.Name())

	if tagInfo.Missing {
		pass.Reportf(field.Pos(), "field %s is missing json tag", field.Name())
		return
	}

	if tagInfo.Inline {
		return
	}

	if tagInfo.Name == "" {
		pass.Reportf(field.Pos(), "field %s has empty json tag", field.Name())
		return
	}

	matched := a.jsonTagRegex.Match([]byte(tagInfo.Name))
	if !matched {
		pass.Reportf(field.Pos(), "field %s json tag does not match pattern %q: %s", field.Name(), a.jsonTagRegex.String(), tagInfo.Name)
	}
}

func defaultConfig(cfg *config.JSONTagsConfig) {
	if cfg.JSONTagRegex == "" {
		cfg.JSONTagRegex = camelCaseRegex
	}
}
