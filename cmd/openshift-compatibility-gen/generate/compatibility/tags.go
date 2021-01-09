package compatibility

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
)

const (
	baseTagName     = "openshift:compatibility-gen"
	levelTagName    = baseTagName + ":level"
	internalTagName = baseTagName + ":internal"
)

// enabledTagValue holds parameters from a tagName tag.
type tagValue struct {
	value string
}

func extractTag(tagName string, comments []string) (*tagValue, error) {
	tagVals := types.ExtractCommentTags("+", comments)[tagName]
	if tagVals == nil {
		// No match for the tag.
		return nil, nil
	}
	// If there are multiple values, abort.
	if len(tagVals) > 1 {
		return nil, fmt.Errorf("found %d %s tags: %q", len(tagVals), tagName, tagVals)
	}

	// If we got here we are returning something.
	tag := &tagValue{}

	// Get the primary value.
	parts := strings.Split(tagVals[0], ",")
	if len(parts) >= 1 {
		tag.value = parts[0]
	}

	// Parse extra arguments.
	parts = parts[1:]
	for i := range parts {
		kv := strings.SplitN(parts[i], "=", 2)
		k := kv[0]
		switch k {
		default:
			return nil, fmt.Errorf("unsupported %s param: %q", tagName, parts[i])
		}
	}
	return tag, nil
}

func containsCompatibilityGenTag(t *types.Type) bool {
	comments := append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
	tags := types.ExtractCommentTags("+", comments)
	for tag := range tags {
		if strings.HasPrefix(tag, baseTagName) {
			return true
		}
	}
	return false
}

func isInternal(t *types.Type) bool {
	comments := append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
	tag, err := extractTag(internalTagName, comments)
	if err != nil {
		klog.Fatalf("%s: error extracting %s tag: %v", t.Name, internalTagName, err)
	}

	// no tag: default to false
	if tag == nil {
		return false
	}

	// no value: default value to true
	if len(tag.value) == 0 {
		return true
	}

	// parse value
	internal, err := strconv.ParseBool(tag.value)
	if err != nil {
		klog.Fatalf("%s: error extracting %s tag value: %v", t.Name, internalTagName, err)
	}
	return internal
}

func extractOpenShiftCompatibilityLevelTag(t *types.Type) (int, bool) {
	comments := append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
	rawTag, err := extractTag(levelTagName, comments)
	if err != nil {
		klog.Fatalf("%s: unable to parse value of %s tag: %v", t.Name, levelTagName, err)
	}
	if rawTag == nil || len(rawTag.value) == 0 {
		return 0, false
	}
	level, err := strconv.Atoi(rawTag.value)
	if err != nil {
		klog.Fatalf("%s: unable to parse value of %s tag: %v", t.Name, levelTagName, err)
	}
	switch level {
	case 1, 2, 3, 4:
	default:
		klog.Fatalf("%s: invalid value of %s tag: %v", t.Name, levelTagName, level)
	}
	return level, true
}
