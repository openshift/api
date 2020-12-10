package comments

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

func GenerateCompatibilityComments(inputPkgs []string) error {
	for _, inputPkg := range inputPkgs {
		output, err := exec.Command("go", "list", "-f", "{{ .Dir }}", inputPkg).Output()
		if err != nil {
			klog.Errorf(string(output))
			return err
		}
		path := string(output)
		path = strings.TrimSpace(path)
		insertCompatibilityLevelComments(path)
	}
	return nil
}

func insertCompatibilityLevelComments(path string) {
	fset := token.NewFileSet()
	var filter func(os.FileInfo) bool
	filter = func(info os.FileInfo) bool {
		if strings.HasPrefix("zz_generated", info.Name()) {
			return false
		}
		switch info.Name() {
		case "doc.go", "register.go", "generated.pb.go":
			return false
		}
		return true
	}
	mode := parser.ParseComments
	pkgs, err := parser.ParseDir(fset, path, filter, mode)
	if err != nil {
		klog.Exit(err)
	}
	commentPattern := regexp.MustCompile(`//\s*Compatibility level [0-9]:`)
	tagPattern := regexp.MustCompile(`//\s*\+`)
	emptyPattern := regexp.MustCompile(`//\s*$`)
	for _, pkg := range pkgs {
		for fname, f := range pkg.Files {
			cmap := ast.NewCommentMap(fset, f, f.Comments)
			fchanged := false
			for _, cgrp := range cmap.Comments() {
				firstTagIndex := -1
				existingCommentIndex := -1
				for i, c := range cgrp.List {
					if commentPattern.MatchString(c.Text) {
						// found existing compatibility comment in godoc
						existingCommentIndex = i
					}
					if firstTagIndex < 0 {
						if tagPattern.MatchString(c.Text) {
							firstTagIndex = i
						}
					}
					if strings.Contains(c.Text, "+openshift:compatibility-gen:level=") {
						level := extractLevelFromComment(c.Text)
						insertIndex := firstTagIndex
						if existingCommentIndex > 0 {
							cgrp.List = append(cgrp.List[:existingCommentIndex], cgrp.List[existingCommentIndex+1:]...)
							insertIndex = existingCommentIndex
						}
						cnew := []*ast.Comment{{Text: fmt.Sprintf("// Compatibility level %d: %s", level, commentForLevel(level))}}
						if insertIndex-1 >= 0 && !emptyPattern.MatchString(cgrp.List[insertIndex-1].Text) {
							cnew = append([]*ast.Comment{{Text: "//"}}, cnew...)
						}
						if !emptyPattern.MatchString(cgrp.List[insertIndex].Text) {
							cnew = append(cnew, &ast.Comment{Text: "//"})
						}
						cgrp.List = append(cgrp.List[:insertIndex], append(cnew, cgrp.List[insertIndex:]...)...)
						fchanged = true
					}
				}
			}
			if !fchanged {
				continue
			}
			var buf bytes.Buffer
			if err := format.Node(&buf, fset, f); err != nil {
				panic(err)
			}
			ioutil.WriteFile(fname, buf.Bytes(), 0777)
		}
	}
}

func commentForLevel(level int) string {
	switch level {
	case 1:
		return "Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer)."
	case 2:
		return "Stable within a major release for a minimum of 9 months or 3 minor releases (whichever is longer)."
	case 3:
		return "Will attempt to be as compatible from version to version as possible, but version to version compatibility is not guaranteed"
	case 4:
		return "No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support."
	default:
		panic(level)
	}
}

func extractLevelFromComment(c string) int {
	p := strings.Split(c, "=")
	if len(p) != 2 {
		klog.Exitf("Could not successfully parse comment tag %s: expected exactly one '='", c)
	}
	level, err := strconv.Atoi(p[1])
	if err != nil {
		klog.Exitf("Could not successfully parse comment tag %s: %v", c, err)
	}
	if level < 1 || level > 4 {
		klog.Exitf("Could not successfully parse comment tag %s: level must be one of 1, 2, 3, or 4", c)
	}
	return level
}
