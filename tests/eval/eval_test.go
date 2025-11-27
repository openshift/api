package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	claudeTimeout    = 5 * time.Minute
	cloneURL         = "https://github.com/openshift/api.git"
	testdataDir      = "testdata"
	patchFileName    = "patch.diff"
	expectedFileName = "expected.txt"
)

var (
	tempDir       string
	localRepoRoot string
	testCases     []string
)

func TestEval(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Review Eval Suite")
}

var _ = BeforeSuite(func() {
	var err error
	localRepoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	Expect(err).NotTo(HaveOccurred())

	By("verifying local AGENTS.md exists")
	_, err = os.Stat(filepath.Join(localRepoRoot, "AGENTS.md"))
	Expect(err).NotTo(HaveOccurred(), "AGENTS.md must exist in repository root")

	By("verifying local .claude/commands/api-review.md exists")
	_, err = os.Stat(filepath.Join(localRepoRoot, ".claude", "commands", "api-review.md"))
	Expect(err).NotTo(HaveOccurred(), ".claude/commands/api-review.md must exist")

	By("creating temp directory for clone")
	tempDir, err = os.MkdirTemp("", "api-review-eval-*")
	Expect(err).NotTo(HaveOccurred())

	By("shallow cloning openshift/api")
	cmd := exec.Command("git", "clone", "--depth", "1", cloneURL, tempDir)
	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git clone failed: %s", string(output))

	copyLocalFiles()

	testdataPath := filepath.Join(localRepoRoot, "tests", "eval", testdataDir)
	testCases, err = discoverTestCases(testdataPath)
	Expect(err).NotTo(HaveOccurred())
	Expect(testCases).NotTo(BeEmpty(), "no test cases found in testdata directory")
})

var _ = AfterSuite(func() {
	if tempDir != "" {
		By("cleaning up temp directory")
		os.RemoveAll(tempDir)
	}
})

func copyLocalFiles() {
	By("copying local AGENTS.md to temp clone")
	src := filepath.Join(localRepoRoot, "AGENTS.md")
	dst := filepath.Join(tempDir, "AGENTS.md")
	data, err := os.ReadFile(src)
	Expect(err).NotTo(HaveOccurred())
	err = os.WriteFile(dst, data, 0644)
	Expect(err).NotTo(HaveOccurred())

	By("copying local .claude directory to temp clone")
	srcClaudeDir := filepath.Join(localRepoRoot, ".claude")
	dstClaudeDir := filepath.Join(tempDir, ".claude")
	os.RemoveAll(dstClaudeDir)
	claudeFS := os.DirFS(srcClaudeDir)
	err = os.CopyFS(dstClaudeDir, claudeFS)
	Expect(err).NotTo(HaveOccurred())
}

func resetRepo() {
	By("verifying origin remote exists")
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		addCmd := exec.Command("git", "remote", "add", "origin", cloneURL)
		addCmd.Dir = tempDir
		output, err := addCmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), "failed to add origin remote: %s", string(output))
	}

	By("resetting repo to clean state")
	cmd = exec.Command("git", "reset", "--hard", "origin/master")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git reset failed: %s", string(output))

	cmd = exec.Command("git", "clean", "-fd")
	cmd.Dir = tempDir
	output, err = cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git clean failed: %s", string(output))

	copyLocalFiles()
}

func discoverTestCases(testdataPath string) ([]string, error) {
	entries, err := os.ReadDir(testdataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read testdata directory: %w", err)
	}

	var cases []string
	for _, entry := range entries {
		if entry.IsDir() {
			patchPath := filepath.Join(testdataPath, entry.Name(), patchFileName)
			expectedPath := filepath.Join(testdataPath, entry.Name(), expectedFileName)

			if _, err := os.Stat(patchPath); err != nil {
				return nil, fmt.Errorf("patch.diff missing in %s: %w", entry.Name(), err)
			}
			if _, err := os.Stat(expectedPath); err != nil {
				return nil, fmt.Errorf("expected.txt missing in %s: %w", entry.Name(), err)
			}

			cases = append(cases, entry.Name())
		}
	}
	return cases, nil
}

func loadTableEntries() []TableEntry {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testdataPath := filepath.Join(cwd, testdataDir)
	cases, err := discoverTestCases(testdataPath)
	Expect(err).NotTo(HaveOccurred())

	var entries []TableEntry
	for _, tc := range cases {
		entries = append(entries, Entry(tc, tc))
	}
	return entries
}

var _ = Describe("API Review Evaluation", func() {
	tableEntries := loadTableEntries()

	DescribeTable("should correctly review",
		func(tc string) {
			resetRepo()

			testCaseDir := filepath.Join(localRepoRoot, "tests", "eval", testdataDir, tc)
			patchPath := filepath.Join(testCaseDir, patchFileName)
			expectedPath := filepath.Join(testCaseDir, expectedFileName)

			By("reading patch file")
			patchContent, err := os.ReadFile(patchPath)
			Expect(err).NotTo(HaveOccurred())

			By("applying patch")
			cmd := exec.Command("git", "apply", "-")
			cmd.Dir = tempDir
			cmd.Stdin = bytes.NewReader(patchContent)
			output, err := cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "git apply failed: %s", string(output))

			By("reading expected issues")
			expectedContent, err := os.ReadFile(expectedPath)
			Expect(err).NotTo(HaveOccurred())

			By("running API review via Claude")
			ctx, cancel := context.WithTimeout(context.Background(), claudeTimeout)
			defer cancel()

			reviewCmd := exec.CommandContext(ctx, "claude",
				"--print",
				"--dangerously-skip-permissions",
				"-p", "/api-review",
				"--allowedTools", "Bash,Read,Grep,Glob,Task",
			)
			reviewCmd.Dir = tempDir

			reviewOutput, err := reviewCmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "claude command failed: %s", string(reviewOutput))

			reviewResult := string(reviewOutput)

			By("comparing results with Claude")
			expectedStr := strings.TrimSpace(string(expectedContent))

			comparePrompt := `You are a judge evaluating an API review output against expected issues.

API review output:
` + reviewResult + `

Expected issues (one per line):
` + expectedStr + `

Compare using SEMANTIC matching - focus on whether the same fundamental problems were identified, not exact wording or action item counts.

You should return pass=true ONLY if BOTH conditions are met:
1. ALL expected issues are semantically covered in the output (the same core problem is identified, even if described differently or split into sub-items)
2. NO unrelated issues are reported - if the review identifies a problem that is NOT semantically related to any expected issue, you should return pass=false

Expanding on an expected issue is OK (e.g., "missing FeatureGate" expanding to include "register in features.go").
Reporting an entirely different issue is NOT OK (e.g., if "missing length validation" is not in expected list, you should return pass=false).

Examples of semantic matches:
- "missing FeatureGate" matches "needs FeatureGate and must register it in features.go"
- "optional field missing omitted behavior" matches "field does not document what happens when not specified"

You should respond with ONLY a JSON object in this exact format (no markdown, no other text):
{"pass": true, "reason": "Brief summary of matched issues"}
or
{"pass": false, "reason": "Explanation of what was missing or what unexpected issue was found"}`

			ctx2, cancel2 := context.WithTimeout(context.Background(), claudeTimeout)
			defer cancel2()

			compareCmd := exec.CommandContext(ctx2, "claude", "--print", "--dangerously-skip-permissions", "-p", comparePrompt)
			compareCmd.Dir = tempDir

			compareOutput, err := compareCmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "claude compare command failed: %s", string(compareOutput))

			var judgeResult struct {
				Pass   bool   `json:"pass"`
				Reason string `json:"reason"`
			}
			compareStr := strings.TrimSpace(string(compareOutput))
			err = json.Unmarshal([]byte(compareStr), &judgeResult)
			Expect(err).NotTo(HaveOccurred(), "failed to parse judge response as JSON: %s", compareStr)

			GinkgoWriter.Printf("Judge result: pass=%v, reason=%s\n", judgeResult.Pass, judgeResult.Reason)
			Expect(judgeResult.Pass).To(BeTrue(), "API review did not match expected issues.\nJudge reason: %s\nReview output:\n%s\nExpected issues:\n%s", judgeResult.Reason, reviewResult, expectedStr)
		},
		tableEntries,
	)
})
