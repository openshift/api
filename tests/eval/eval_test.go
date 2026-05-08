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
	"syscall"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	claudeTimeout    = 10 * time.Minute
	testdataDir      = "testdata"
	goldenDir        = "golden"
	integrationDir   = "integration"
	patchFileName    = "patch.diff"
	expectedFileName = "expected.txt"

	// Pinned commit that all eval patches are written against.
	// Override with EVAL_BASE_REF env var during development.
	// Update this when rebasing patches — see EVAL_EXPANSION.md § "Updating the Eval Baseline".
	defaultEvalBaseRef = "FILL_BEFORE_MERGE"

	defaultGoldenModel      = "sonnet"
	defaultIntegrationModel = "opus"
	defaultJudgeModel       = "haiku"

	// Adjust once we have baseline data from CI runs.
	defaultThreshold = 0.8

	goldenJudgePrompt = `You are a judge evaluating an API review output against expected issues.

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

	integrationJudgePrompt = `You are a judge evaluating an API review output against expected issues.

API review output:
%s

Expected issues (one per line):
%s

Compare using SEMANTIC matching - focus on whether the same fundamental problems were identified, not exact wording or action item counts.

You should return pass=true ONLY if BOTH conditions are met:
1. At least 70%% of the expected issues are semantically covered in the output (the same core problem is identified, even if described differently or split into sub-items). Missing some expected issues is acceptable.
2. NO false positives - if the review identifies a problem that is NOT semantically related to any expected issue, you MUST return pass=false. Precision matters more than recall.

Expanding on an expected issue is OK (e.g., "missing FeatureGate" expanding to include "register in features.go").
Reporting an entirely different issue is NOT OK (e.g., if "missing enhancementPR" is not in expected list, you should return pass=false).

Examples of semantic matches:
- "missing FeatureGate" matches "needs FeatureGate and must register it in features.go"
- "optional field missing omitted behavior" matches "field does not document what happens when not specified"

You MUST respond with ONLY a raw JSON object. Do NOT wrap in markdown code blocks. Do NOT include any other text.
{"pass": true, "reason": "Brief summary: N of M expected issues found, no false positives"}
or
{"pass": false, "reason": "Explanation of what false positive was found OR that too few expected issues were covered"}`
)

type testCase struct {
	Name           string
	Patch          []byte
	ExpectedIssues string
}

var (
	repoRoot         string
	goldenModel      string
	integrationModel string
	judgeModel       string
	evalRuns         int
	evalThreshold    float64
	evalVerbose      bool
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

var _ = SynchronizedBeforeSuite(func() []byte {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	Expect(err).NotTo(HaveOccurred(), "failed to resolve repository root")

	_, err = os.Stat(filepath.Join(root, "AGENTS.md"))
	Expect(err).NotTo(HaveOccurred(), "AGENTS.md must exist in repository root")

	_, err = os.Stat(filepath.Join(root, ".claude", "commands", "api-review.md"))
	Expect(err).NotTo(HaveOccurred(), ".claude/commands/api-review.md must exist")

	return []byte(root)
}, func(data []byte) {
	repoRoot = string(data)

	goldenModel = envOrDefault("EVAL_GOLDEN_MODEL", defaultGoldenModel)
	integrationModel = envOrDefault("EVAL_INTEGRATION_MODEL", defaultIntegrationModel)
	judgeModel = envOrDefault("EVAL_JUDGE_MODEL", defaultJudgeModel)
	evalVerbose = os.Getenv("EVAL_VERBOSE") != ""

	var err error
	evalRuns, err = strconv.Atoi(envOrDefault("EVAL_RUNS", "1"))
	Expect(err).NotTo(HaveOccurred(), "EVAL_RUNS must be an integer")
	Expect(evalRuns).To(BeNumerically(">", 0), "EVAL_RUNS must be positive")

	evalThreshold, err = strconv.ParseFloat(envOrDefault("EVAL_THRESHOLD", fmt.Sprintf("%g", defaultThreshold)), 64)
	Expect(err).NotTo(HaveOccurred(), "EVAL_THRESHOLD must be a float")
	Expect(evalThreshold).To(BeNumerically(">=", 0.0), "EVAL_THRESHOLD=%g is out of range [0.0, 1.0]", evalThreshold)
	Expect(evalThreshold).To(BeNumerically("<=", 1.0), "EVAL_THRESHOLD=%g is out of range [0.0, 1.0]", evalThreshold)
})

var _ = ReportAfterSuite("Eval Summary", func(report Report) {
	threshold, _ := strconv.ParseFloat(envOrDefault("EVAL_THRESHOLD", fmt.Sprintf("%g", defaultThreshold)), 64)

	type resultRow struct {
		Name string
		Rate float64
	}
	var rows []resultRow
	var totalReviewerCost, totalJudgeCost float64

	for _, spec := range report.SpecReports {
		var name, passRate, reviewerCost, judgeCost string
		for _, entry := range spec.ReportEntries {
			switch entry.Name {
			case "test_name":
				name = entry.StringRepresentation()
			case "pass_rate":
				passRate = entry.StringRepresentation()
			case "reviewer_cost_usd":
				reviewerCost = entry.StringRepresentation()
			case "judge_cost_usd":
				judgeCost = entry.StringRepresentation()
			}
		}
		if name == "" {
			continue
		}
		rc, _ := strconv.ParseFloat(reviewerCost, 64)
		jc, _ := strconv.ParseFloat(judgeCost, 64)
		totalReviewerCost += rc
		totalJudgeCost += jc
		rate, _ := strconv.ParseFloat(passRate, 64)
		rows = append(rows, resultRow{Name: name, Rate: rate})
	}

	if len(rows) > 0 {
		GinkgoWriter.Printf("\n%-35s | %s\n", "Test Case", "Rate")
		GinkgoWriter.Printf("%s\n", strings.Repeat("-", 50))
		for _, r := range rows {
			line := fmt.Sprintf("%-35s | %3.0f%%", r.Name, r.Rate*100)
			if r.Rate < threshold {
				line += " <- FAIL"
			}
			GinkgoWriter.Println(line)
		}
		GinkgoWriter.Printf("\nThreshold: %.0f%%\n", threshold*100)
	}

	GinkgoWriter.Printf("Total Cost: $%.4f (Reviewer: $%.4f, Judge: $%.4f)\n",
		totalReviewerCost+totalJudgeCost, totalReviewerCost, totalJudgeCost)
})

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

func runClaude(ctx context.Context, cmd *exec.Cmd) (stdout, stderr []byte, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
}

// parseClaudeJSON tries a clean unmarshal first, then falls back to
// scanning for the first '{' in case the CLI leaked warnings to stdout.
func parseClaudeJSON(stdout []byte, dest interface{}) error {
	if err := json.Unmarshal(stdout, dest); err == nil {
		return nil
	}
	if idx := bytes.IndexByte(stdout, '{'); idx > 0 {
		return json.Unmarshal(stdout[idx:], dest)
	}
	return json.Unmarshal(stdout, dest)
}

func stripMarkdownCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func createWorktree(patch []byte) string {
	dir, err := os.MkdirTemp("", "eval-worktree-*")
	Expect(err).NotTo(HaveOccurred())

	ref := defaultEvalBaseRef
	if override := os.Getenv("EVAL_BASE_REF"); override != "" {
		ref = override
	}
	By("creating git worktree from " + ref)
	cmd := exec.Command("git", "worktree", "add", "--detach", dir, ref)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git worktree add failed: %s", string(output))

	By("copying .claude/commands from working tree")
	srcCommands := filepath.Join(repoRoot, ".claude", "commands")
	dstCommands := filepath.Join(dir, ".claude", "commands")
	err = os.MkdirAll(dstCommands, 0o755)
	Expect(err).NotTo(HaveOccurred(), "failed to create .claude/commands in worktree")
	cmdFiles, err := os.ReadDir(srcCommands)
	Expect(err).NotTo(HaveOccurred(), "failed to read .claude/commands")
	for _, f := range cmdFiles {
		if f.IsDir() {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(srcCommands, f.Name()))
		Expect(readErr).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(dstCommands, f.Name()), data, 0o644)).To(Succeed())
	}

	By("applying patch in worktree")
	cmd = exec.Command("git", "apply", "--3way", "-")
	cmd.Dir = dir
	cmd.Stdin = bytes.NewReader(patch)
	output, err = cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "git apply failed: %s", string(output))

	return dir
}

func removeWorktree(dir string) {
	cmd := exec.Command("git", "worktree", "remove", "--force", dir)
	cmd.Dir = repoRoot
	cmd.CombinedOutput()
	os.RemoveAll(dir)
}

func runAPIReview(model, workDir string) (string, float64) {
	By(fmt.Sprintf("running API review via Claude (%s)", model))
	ctx, cancel := context.WithTimeout(context.Background(), claudeTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "claude",
		"--print",
		"--dangerously-skip-permissions",
		"--model", model,
		"--max-turns", "30",
		"-p", "/api-review",
		"--allowedTools", "Bash,Read,Grep,Glob,Task",
		"--output-format", "json",
	)
	cmd.Dir = workDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	stdout, stderr, err := runClaude(ctx, cmd)
	if len(stderr) > 0 {
		GinkgoWriter.Printf("[reviewer] stderr: %s\n", string(stderr))
	}
	Expect(err).NotTo(HaveOccurred(), "claude command failed: stdout=%s stderr=%s", string(stdout), string(stderr))

	if evalVerbose {
		GinkgoWriter.Printf("\n--- Raw claude stdout (%d bytes) ---\n%s\n--- End raw ---\n", len(stdout), string(stdout))
	}

	var parsed claudeOutput
	err = parseClaudeJSON(stdout, &parsed)
	Expect(err).NotTo(HaveOccurred(), "failed to parse claude output: %s", string(stdout))

	if parsed.Result == "" {
		GinkgoWriter.Printf("WARNING: empty result from claude. Raw stdout (%d bytes): %s\n", len(stdout), string(stdout))
	}

	return parsed.Result, parsed.TotalCostUSD
}

func runJudge(model, reviewOutput, expectedIssues, judgePromptTemplate string) (evalResult, float64) {
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
	cmd.Dir = repoRoot
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	stdout, stderr, err := runClaude(ctx, cmd)
	if len(stderr) > 0 {
		GinkgoWriter.Printf("[judge] stderr: %s\n", string(stderr))
	}
	Expect(err).NotTo(HaveOccurred(), "claude judge command failed: stdout=%s stderr=%s", string(stdout), string(stderr))

	var parsed claudeOutput
	err = parseClaudeJSON(stdout, &parsed)
	Expect(err).NotTo(HaveOccurred(), "failed to parse judge output: %s", string(stdout))

	var result evalResult
	jsonStr := stripMarkdownCodeBlock(parsed.Result)
	err = json.Unmarshal([]byte(jsonStr), &result)
	Expect(err).NotTo(HaveOccurred(), "failed to parse judge response as JSON: %s", parsed.Result)
	return result, parsed.TotalCostUSD
}

func runTestCase(tc testCase, reviewModel, judgeModelName, judgePromptTemplate string) {
	var passed int
	var failures []string
	var caseReviewerCost, caseJudgeCost float64

	for i := range evalRuns {
		func() {
			By(fmt.Sprintf("[%s] run %d/%d", tc.Name, i+1, evalRuns))

			workDir := createWorktree(tc.Patch)
			defer removeWorktree(workDir)

			reviewOutput, reviewCost := runAPIReview(reviewModel, workDir)
			judgeResult, judgeCost := runJudge(judgeModelName, reviewOutput, tc.ExpectedIssues, judgePromptTemplate)
			caseReviewerCost += reviewCost
			caseJudgeCost += judgeCost

			if evalVerbose || !judgeResult.Pass {
				GinkgoWriter.Printf("\n--- [%s] Reviewer Output (run %d/%d) ---\n%s\n--- End Reviewer Output ---\n\n",
					tc.Name, i+1, evalRuns, reviewOutput)
			}

			GinkgoWriter.Printf("[%s] Run %d/%d: pass=%v, Reviewer=$%.4f, Judge=$%.4f\n",
				tc.Name, i+1, evalRuns, judgeResult.Pass, reviewCost, judgeCost)
			GinkgoWriter.Printf("[%s] Judge reason: %s\n", tc.Name, judgeResult.Reason)

			if judgeResult.Pass {
				passed++
			} else {
				failures = append(failures, fmt.Sprintf("run %d: %s", i+1, judgeResult.Reason))
			}
		}()
	}

	rate := float64(passed) / float64(evalRuns)

	AddReportEntry("test_name", tc.Name)
	AddReportEntry("pass_rate", fmt.Sprintf("%.2f", rate))
	AddReportEntry("reviewer_cost_usd", fmt.Sprintf("%.4f", caseReviewerCost))
	AddReportEntry("judge_cost_usd", fmt.Sprintf("%.4f", caseJudgeCost))
	AddReportEntry("total_cost_usd", fmt.Sprintf("%.4f", caseReviewerCost+caseJudgeCost))

	GinkgoWriter.Printf("[%s] Result: %d/%d passed (%.0f%%), threshold: %.0f%%\n",
		tc.Name, passed, evalRuns, rate*100, evalThreshold*100)

	if rate < evalThreshold {
		msg := fmt.Sprintf("%s: %d/%d passed (%.0f%%), threshold %.0f%%",
			tc.Name, passed, evalRuns, rate*100, evalThreshold*100)
		for _, f := range failures {
			msg += fmt.Sprintf("\n  - %s", f)
		}
		Fail(msg)
	}
}

var _ = Describe("API Review Evaluation", func() {
	Context("Golden Tests", func() {
		DescribeTable("should correctly identify single issues",
			func(tc testCase) {
				runTestCase(tc, goldenModel, judgeModel, goldenJudgePrompt)
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
				runTestCase(tc, integrationModel, judgeModel, integrationJudgePrompt)
			},
			entries,
		)
	})
})
