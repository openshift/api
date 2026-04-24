//go:build eval

package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	goldenDir        = "golden"
	integrationDir   = "integration"
	patchFileName    = "patch.diff"
	expectedFileName = "expected.txt"

	sonnetModel = "claude-sonnet-4-6"
	opusModel   = "claude-opus-4-6"
	haikuModel  = "claude-haiku-4-5-20251001"

	defaultGoldenModel      = sonnetModel
	defaultIntegrationModel = opusModel
	defaultJudgeModel       = haikuModel

	// Adjust once we have baseline data from CI runs.
	defaultThreshold = 0.8

	judgePromptTemplate = `You are a judge evaluating an API review output against expected issues.

API review output:
%s

Expected issues (one per line):
%s

Compare using SEMANTIC matching - focus on whether the same fundamental problems were identified, not exact wording or action item counts.

You should return pass=true ONLY if BOTH conditions are met:
1. ALL expected issues are semantically covered in the output (the same core problem is identified, even if described differently or split into sub-items)
2. NO unrelated issues are reported - if the review identifies a problem that is NOT semantically related to any expected issue, you should return pass=false

Expanding on an expected issue is OK (e.g., "missing FeatureGate" expanding to include "register in features.go").
Reporting an entirely different issue is NOT OK (e.g., if "missing length validation" is not in expected list, you should return pass=false).

Examples of semantic matches:
- "missing FeatureGate" matches "needs FeatureGate and must register it in features.go"
- "optional field missing omitted behavior" matches "field does not document what happens when not specified"

You MUST respond with ONLY a raw JSON object. Do NOT wrap in markdown code blocks. Do NOT include any other text.
{"pass": true, "reason": "Brief summary of matched issues"}
or
{"pass": false, "reason": "Explanation of what was missing or what unexpected issue was found"}`
)

type testCase struct {
	Name           string
	Patch          []byte
	ExpectedIssues string
}

type testCaseResult struct {
	Name     string
	Passed   int
	Runs     int
	Rate     float64
	Failures []string
}

var (
	tempDir           string
	localRepoRoot     string
	goldenModel       string
	integrationModel  string
	judgeModel        string
	evalRuns          int
	evalThreshold     float64
	totalReviewerCost float64
	totalJudgeCost    float64
	allResults        []testCaseResult
)

type claudeOutput struct {
	Type         string  `json:"type"`
	Result       string  `json:"result"`
	TotalCostUSD float64 `json:"total_cost_usd"`
}

func TestEval(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Review Eval Suite")
}

func envOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

var _ = BeforeSuite(func() {
	goldenModel = envOrDefault("EVAL_GOLDEN_MODEL", defaultGoldenModel)
	integrationModel = envOrDefault("EVAL_INTEGRATION_MODEL", defaultIntegrationModel)
	judgeModel = envOrDefault("EVAL_JUDGE_MODEL", defaultJudgeModel)

	var err error
	evalRuns, err = strconv.Atoi(envOrDefault("EVAL_RUNS", "1"))
	Expect(err).NotTo(HaveOccurred(), "EVAL_RUNS must be an integer")
	Expect(evalRuns).To(BeNumerically(">", 0), "EVAL_RUNS must be positive")

	evalThreshold, err = strconv.ParseFloat(envOrDefault("EVAL_THRESHOLD", fmt.Sprintf("%g", defaultThreshold)), 64)
	Expect(err).NotTo(HaveOccurred(), "EVAL_THRESHOLD must be a float")

	localRepoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	Expect(err).NotTo(HaveOccurred(), "failed to resolve repository root")

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
})

var _ = AfterSuite(func() {
	if tempDir != "" {
		By("cleaning up temp directory")
		os.RemoveAll(tempDir)
	}

	if len(allResults) > 0 {
		fmt.Printf("\n%-35s | %6s | %4s | %s\n", "Test Case", "Passed", "Runs", "Rate")
		fmt.Printf("%s\n", strings.Repeat("-", 65))
		for _, r := range allResults {
			line := fmt.Sprintf("%-35s | %6d | %4d | %3.0f%%", r.Name, r.Passed, r.Runs, r.Rate*100)
			if r.Rate < evalThreshold {
				line += " <- FAIL"
			}
			fmt.Println(line)
		}
		fmt.Printf("\nThreshold: %.0f%%\n", evalThreshold*100)
	}

	fmt.Printf("Total Cost: $%.4f (Reviewer: $%.4f, Judge: $%.4f)\n", totalReviewerCost+totalJudgeCost, totalReviewerCost, totalJudgeCost)
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

func discoverTestCases(testdataPath string) ([]testCase, error) {
	entries, err := os.ReadDir(testdataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read testdata directory: %w", err)
	}

	var cases []testCase
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		patch, err := os.ReadFile(filepath.Join(testdataPath, entry.Name(), patchFileName))
		if err != nil {
			return nil, fmt.Errorf("patch.diff missing in %s: %w", entry.Name(), err)
		}

		expected, err := os.ReadFile(filepath.Join(testdataPath, entry.Name(), expectedFileName))
		if err != nil {
			return nil, fmt.Errorf("expected.txt missing in %s: %w", entry.Name(), err)
		}

		cases = append(cases, testCase{
			Name:           entry.Name(),
			Patch:          patch,
			ExpectedIssues: strings.TrimSpace(string(expected)),
		})
	}
	return cases, nil
}

func loadEntries(dir string) []TableEntry {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	cases, err := discoverTestCases(filepath.Join(cwd, testdataDir, dir))
	if err != nil || len(cases) == 0 {
		return nil
	}

	var entries []TableEntry
	for _, tc := range cases {
		entries = append(entries, Entry(tc.Name, tc))
	}
	return entries
}

type evalResult struct {
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

func stripMarkdownCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func applyPatch(patch []byte) {
	By("applying patch")
	cmd := exec.Command("git", "apply", "-")
	cmd.Dir = tempDir
	cmd.Stdin = bytes.NewReader(patch)
	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git apply failed: %s", string(output))
}

func runAPIReview(model string) (string, float64) {
	By(fmt.Sprintf("running API review via Claude (%s)", model))
	ctx, cancel := context.WithTimeout(context.Background(), claudeTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude",
		"--print",
		"--dangerously-skip-permissions",
		"--model", model,
		"-p", "/api-review",
		"--allowedTools", "Bash,Read,Grep,Glob,Task",
		"--output-format", "json",
	)
	cmd.Dir = tempDir

	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "claude command failed: %s", string(output))

	var parsed claudeOutput
	err = json.Unmarshal(output, &parsed)
	Expect(err).NotTo(HaveOccurred(), "failed to parse claude output: %s", string(output))

	totalReviewerCost += parsed.TotalCostUSD
	return parsed.Result, parsed.TotalCostUSD
}

func runJudge(model, reviewOutput, expectedIssues string) (evalResult, float64) {
	By(fmt.Sprintf("comparing results with Claude judge (%s)", model))
	ctx, cancel := context.WithTimeout(context.Background(), claudeTimeout)
	defer cancel()

	prompt := fmt.Sprintf(judgePromptTemplate, reviewOutput, expectedIssues)
	cmd := exec.CommandContext(ctx, "claude",
		"--print",
		"--dangerously-skip-permissions",
		"--model", model,
		"-p", prompt,
		"--output-format", "json",
	)
	cmd.Dir = tempDir

	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "claude judge command failed: %s", string(output))

	var parsed claudeOutput
	err = json.Unmarshal(output, &parsed)
	Expect(err).NotTo(HaveOccurred(), "failed to parse judge output: %s", string(output))

	totalJudgeCost += parsed.TotalCostUSD

	var result evalResult
	jsonStr := stripMarkdownCodeBlock(parsed.Result)
	err = json.Unmarshal([]byte(jsonStr), &result)
	Expect(err).NotTo(HaveOccurred(), "failed to parse judge response as JSON: %s", parsed.Result)
	return result, parsed.TotalCostUSD
}

func runTestCase(tc testCase, reviewModel, judgeModelName string) {
	result := testCaseResult{Name: tc.Name, Runs: evalRuns}

	for i := range evalRuns {
		By(fmt.Sprintf("run %d/%d", i+1, evalRuns))
		resetRepo()
		applyPatch(tc.Patch)

		reviewOutput, reviewCost := runAPIReview(reviewModel)
		judgeResult, judgeCost := runJudge(judgeModelName, reviewOutput, tc.ExpectedIssues)

		GinkgoWriter.Printf("Run %d/%d: pass=%v, Reviewer=$%.4f, Judge=$%.4f\n",
			i+1, evalRuns, judgeResult.Pass, reviewCost, judgeCost)
		GinkgoWriter.Printf("Judge reason: %s\n", judgeResult.Reason)

		if judgeResult.Pass {
			result.Passed++
		} else {
			result.Failures = append(result.Failures, fmt.Sprintf("run %d: %s", i+1, judgeResult.Reason))
		}
	}

	result.Rate = float64(result.Passed) / float64(result.Runs)
	allResults = append(allResults, result)

	GinkgoWriter.Printf("Result: %d/%d passed (%.0f%%), threshold: %.0f%%\n",
		result.Passed, result.Runs, result.Rate*100, evalThreshold*100)

	Expect(result.Rate).To(BeNumerically(">=", evalThreshold),
		"pass rate %.0f%% below threshold %.0f%% for %s.\nFailures:\n%s",
		result.Rate*100, evalThreshold*100, tc.Name, strings.Join(result.Failures, "\n"))
}

var _ = Describe("API Review Evaluation", func() {
	Context("Golden Tests", func() {
		DescribeTable("should correctly identify single issues",
			func(tc testCase) {
				runTestCase(tc, goldenModel, judgeModel)
			},
			loadEntries(goldenDir),
		)
	})

	Context("Integration Tests", func() {
		entries := loadEntries(integrationDir)
		if len(entries) == 0 {
			return
		}

		DescribeTable("should correctly identify multiple issues",
			func(tc testCase) {
				runTestCase(tc, integrationModel, judgeModel)
			},
			entries,
		)
	})
})
