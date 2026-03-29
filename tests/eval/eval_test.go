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
	goldenDir        = "golden"
	integrationDir   = "integration"
	patchFileName    = "patch.diff"
	expectedFileName = "expected.txt"

	// Setting everything to haiku for development is cheap and quick.
	// Opus is expensive, and seems best for the integration tests
	// but often hallucinates more.
	sonnetModel = "claude-sonnet-4-5@20250929"
	opusModel   = "claude-opus-4-5@20251101"
	haikuModel  = "claude-haiku-4-5@20251001"

	defaultGoldenModel      = sonnetModel
	defaultIntegrationModel = opusModel
	defaultJudgeModel       = haikuModel

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

var (
	tempDir           string
	localRepoRoot     string
	testCases         []string
	goldenModel       string
	integrationModel  string
	judgeModel        string
	totalReviewerCost float64
	totalJudgeCost    float64
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

	goldenPath := filepath.Join(localRepoRoot, "tests", "eval", testdataDir, goldenDir)
	testCases, err = discoverTestCases(goldenPath)
	Expect(err).NotTo(HaveOccurred())
	Expect(testCases).NotTo(BeEmpty(), "no test cases found in testdata/golden directory")
})

var _ = AfterSuite(func() {
	if tempDir != "" {
		By("cleaning up temp directory")
		os.RemoveAll(tempDir)
	}
	fmt.Printf("\nTotal Cost: $%.4f (Reviewer: $%.4f, Judge: $%.4f)\n", totalReviewerCost+totalJudgeCost, totalReviewerCost, totalJudgeCost)
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

func loadGoldenEntries() []TableEntry {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	goldenPath := filepath.Join(cwd, testdataDir, goldenDir)
	cases, err := discoverTestCases(goldenPath)
	Expect(err).NotTo(HaveOccurred())

	var entries []TableEntry
	for _, tc := range cases {
		entries = append(entries, Entry(tc, tc))
	}
	return entries
}

func loadIntegrationEntries() []TableEntry {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	integrationPath := filepath.Join(cwd, testdataDir, integrationDir)
	cases, err := discoverTestCases(integrationPath)
	if err != nil || len(cases) == 0 {
		return nil
	}

	var entries []TableEntry
	for _, tc := range cases {
		entries = append(entries, Entry(tc, tc))
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

func readAndApplyPatch(patchPath string) {
	By("reading and applying patch")
	patchContent, err := os.ReadFile(patchPath)
	Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("git", "apply", "-")
	cmd.Dir = tempDir
	cmd.Stdin = bytes.NewReader(patchContent)
	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git apply failed: %s", string(output))
}

// runAPIReview and runJudge can probably share some common code.
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

func runTestCase(tier, tc, reviewModel, judgeModelName string) {
	resetRepo()

	testCaseDir := filepath.Join(localRepoRoot, "tests", "eval", testdataDir, tier, tc)
	readAndApplyPatch(filepath.Join(testCaseDir, patchFileName))

	expectedContent, err := os.ReadFile(filepath.Join(testCaseDir, expectedFileName))
	Expect(err).NotTo(HaveOccurred())
	expectedIssues := strings.TrimSpace(string(expectedContent))

	reviewOutput, reviewCost := runAPIReview(reviewModel)
	result, judgeCost := runJudge(judgeModelName, reviewOutput, expectedIssues)

	GinkgoWriter.Printf("Cost: Reviewer=$%.4f, Judge=$%.4f, Total=$%.4f\n", reviewCost, judgeCost, reviewCost+judgeCost)
	GinkgoWriter.Printf("Judge result: pass=%v, reason=%s\n", result.Pass, result.Reason)
	Expect(result.Pass).To(BeTrue(), "API review did not match expected issues.\nJudge reason: %s\nReview output:\n%s\nExpected issues:\n%s", result.Reason, reviewOutput, expectedIssues)
}

var _ = Describe("API Review Evaluation", func() {
	Context("Golden Tests", func() {
		goldenEntries := loadGoldenEntries()

		DescribeTable("should correctly identify single issues",
			func(tc string) {
				runTestCase(goldenDir, tc, goldenModel, judgeModel)
			},
			goldenEntries,
		)
	})

	Context("Integration Tests", func() {
		integrationEntries := loadIntegrationEntries()
		if len(integrationEntries) == 0 {
			return
		}

		DescribeTable("should correctly identify multiple issues",
			func(tc string) {
				runTestCase(integrationDir, tc, integrationModel, judgeModel)
			},
			integrationEntries,
		)
	})
})
