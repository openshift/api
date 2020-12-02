package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
)

// Packages makes the package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)), boilerplate...)

	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		for _, t := range pkg.Types {
			apiType := isAPIType(t)
			hasTag := containsCompatibilityLevelTag(t)
			switch {
			case apiType && !hasTag:
				klog.Exitf("Type %v is missing the required %s tag.", t, levelTagName)
			case !apiType && hasTag:
				klog.Exitf("Type %v requests API compatibility generation but is not an API", t)
			}
		}

		fmt.Printf("Generating package %q\n", pkg.Path)
		path := pkg.Path
		// if the source path is within a /vendor/ directory (for example,
		// k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/apis/meta/v1), allow
		// generation to output to the proper relative path (under vendor).
		// Otherwise, the generator will create the file in the wrong location
		// in the output directory.
		// TODO: build a more fundamental concept in gengo for dealing with modifications
		// to vendored packages.
		if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
			expandedPath := strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
			if strings.Contains(expandedPath, "/vendor/") {
				path = expandedPath
			}
		}
		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						NewCompatibilityGen(arguments.OutputFileBaseName, pkg.Path),
					}
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})
	}
	return packages
}

// CustomArgs is used tby the go2idl framework to pass args specific to this generator.
type CustomArgs struct {
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}
