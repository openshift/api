package utils

import (
	"bytes"
	"strings"

	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	colourReset = "\x1b[0m"
	colourRed   = "\x1b[31m"
	colourGreen = "\x1b[32m"
)

// Diff returns a string containing the diff between the two strings.
// This is equivalent to running `git diff` on the two strings.
func Diff(a, b []byte) string {
	diffs := diff.Do(string(a), string(b))

	if len(diffs) > 1 {
		return prettyPrintDiff(diffs)
	}

	return ""
}

// prettyPrintDiff prints the diff for the file out as if it were a git diff.
func prettyPrintDiff(diffs []diffmatchpatch.Diff) string {
	diffsByLine := splitDiffsByLine(diffs)
	buff := bytes.NewBuffer(nil)

	for i, diff := range diffsByLine {
		switch diff.diffType {
		case diffmatchpatch.DiffInsert, diffmatchpatch.DiffDelete:
			// Print any Insert or Delete diffs with the previous 3 and next 3 lines
			// for context around the diff.
			printDiffWithSurroundingLines(diffsByLine, i, buff, diff)
		case diffmatchpatch.DiffEqual:
			if i == 0 || i == len(diffsByLine)-1 {
				// If this is the first or last diff, ignore it.
				continue
			}

			// Write a break to indicate there are more lines.
			// that have been omitted between the diffs.
			buff.WriteString("...\n")
		}
	}

	return buff.String()
}

// printDiffWithSurroundingLines prints the diff with the previous and next 3 lines.
// This provides some context around the diff.
func printDiffWithSurroundingLines(diffs []diffLines, diffIndex int, buff *bytes.Buffer, diff diffLines) {
	// If this isn't the first diff, print the previous diff.
	if diffIndex > 0 {
		prevDiff := diffs[diffIndex-1]

		// Only print the previous diff if it an equal diff.
		// Other diff types will be printed by themselves.
		if prevDiff.diffType == diffmatchpatch.DiffEqual {
			printPreviousDiff(buff, prevDiff)
		}
	}

	switch diff.diffType {
	case diffmatchpatch.DiffInsert:
		_, _ = buff.WriteString(colourGreen)
		printDiffLines(buff, "+ ", diff.lines)
		_, _ = buff.WriteString(colourReset)
	case diffmatchpatch.DiffDelete:
		_, _ = buff.WriteString(colourRed)
		printDiffLines(buff, "- ", diff.lines)
		_, _ = buff.WriteString(colourReset)
	}

	// If this isn't the last diff, print the next diff.
	if diffIndex < len(diffs)-1 {
		nextDiff := diffs[diffIndex+1]

		// Only print the next diff if it an equal diff.
		// Other diff types will be printed by themselves.
		if nextDiff.diffType == diffmatchpatch.DiffEqual {
			printNextDiff(buff, nextDiff)
		}
	}
}

// printPreviousDiff prints the previous diff with the last 3 lines.
func printPreviousDiff(buff *bytes.Buffer, diff diffLines) {
	var lines []string

	// Only print up to the last 3 lines of the previous diff.
	for j := 0; j < len(diff.lines) && j < 3; j++ {
		ln := len(diff.lines) - j - 1
		lines = append(lines, diff.lines[ln])
	}

	printDiffLines(buff, "  ", lines)
}

// printNextDiff prints the next diff with the first 3 lines.
func printNextDiff(buff *bytes.Buffer, diff diffLines) {
	var lines []string

	// Only print up to the first 3 lines of the next diff.
	for j := 0; j < len(diff.lines) && j < 3; j++ {
		lines = append(lines, diff.lines[j])
	}

	printDiffLines(buff, " ", lines)
}

// printDiffLines prints each line in the diff as a separate line with the given prefix.
func printDiffLines(buff *bytes.Buffer, prefix string, lines []string) {
	for _, line := range lines {
		_, _ = buff.WriteString(prefix + line + "\n")
	}
}

// diffLines is a diff with the lines split by newline.
type diffLines struct {
	diffType diffmatchpatch.Operation
	lines    []string
}

// splitDiffsByLine splits the diffs by line.
func splitDiffsByLine(in []diffmatchpatch.Diff) []diffLines {
	var out []diffLines

	for _, d := range in {
		data := strings.TrimSuffix(d.Text, "\n")
		lines := strings.Split(data, "\n")

		out = append(out, diffLines{
			diffType: d.Type,
			lines:    lines,
		})
	}

	return out
}
