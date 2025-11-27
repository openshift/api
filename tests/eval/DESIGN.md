# API Review Command Eval Test Suite

Design document for a Go/Ginkgo-based evaluation framework to test the `/api-review` Claude command against known API review scenarios.

## Overview

This test suite validates that the `/api-review` Claude command correctly identifies API documentation issues. Each test case consists of a patch file and expected issues. The suite applies patches to a clean clone of the repository, runs the API review command, and verifies the output matches expectations exactly (no missing issues, no hallucinated issues).

## Directory Structure

```
tests/eval/
├── eval_test.go              # Main Ginkgo test suite
├── DESIGN.md                 # This file
├── testdata/
│   ├── golden/               # Single-issue tests
│   │   ├── missing-optional-doc/
│   │   │   ├── patch.diff
│   │   │   └── expected.txt
│   │   ├── undocumented-enum/
│   │   │   ├── patch.diff
│   │   │   └── expected.txt
│   │   └── valid-api-change/
│   │       ├── patch.diff
│   │       └── expected.txt  # Empty file = no issues expected
│   └── integration/          # Multi-issue tests
│       └── new-field-multiple-issues/
│           ├── patch.diff
│           └── expected.txt
```

## Test Case Format

### patch.diff

Standard git diff format:

```diff
diff --git a/config/v1/types.go b/config/v1/types.go
--- a/config/v1/types.go
+++ b/config/v1/types.go
@@ -10,0 +11,5 @@
+// MyField does something
+// +optional
+// +kubebuilder:validation:Enum=Foo;Bar
+MyField string `json:"myField"`
```

### expected.txt

One expected issue per line:

```
enum values Foo and Bar not documented in comment
optional field does not explain behavior when omitted
```

Empty file means the API change should pass review with no issues.

**Note**: Order of issues in `expected.txt` does not matter. Comparison uses semantic matching, not exact string matching.

## Test Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Pre-flight:                                               │
│    a. Verify local AGENTS.md and                             │
│       .claude/commands/api-review.md exist                   │
│    b. These will be copied to temp dir after clone           │
│       (ensures local changes are tested)                     │
├─────────────────────────────────────────────────────────────┤
│ 2. Setup (once):                                             │
│    a. Shallow clone openshift/api to temp dir                │
│    b. Copy local AGENTS.md and .claude/ to temp dir          │
├─────────────────────────────────────────────────────────────┤
│ 3. For each test case (sequential):                          │
│    a. Reset repo to clean state                              │
│    b. Apply patch.diff                                       │
│    c. Run claude api-review on changed files                 │
│    d. Run claude to compare output vs expected.txt           │
│    e. Parse true/false response, assert                      │
├─────────────────────────────────────────────────────────────┤
│ 4. Teardown: Remove temp dir                                 │
└─────────────────────────────────────────────────────────────┘
```

## Reset Between Tests

```bash
git reset --hard origin/master && git clean -fd
```

- `git reset --hard origin/master`: Resets all tracked files to match the remote master branch, discarding any local commits and staged/unstaged changes
- `git clean -f`: Force remove untracked files (files not in git)
- `git clean -d`: Also remove untracked directories

After reset, re-copy `AGENTS.md` and `.claude/` from local source to preserve local modifications being tested.

**Remote origin handling**: The shallow clone creates `origin` automatically. Before reset, verify remote exists:

```bash
git remote get-url origin || git remote add origin https://github.com/openshift/api.git
```

## Claude Invocations

### Step 1 - Run API Review

```bash
claude --print -p "/api-review" --allowedTools "Bash,Read,Grep,Glob,Task" <files>
```

### Step 2 - Compare Results

```bash
claude --print -p "Given this API review output:
<output>

Expected issues (one per line):
<contents of expected.txt>

Compare the review output against the expected issues list.
The result is 'true' ONLY if:
1. ALL expected issues are identified in the output
2. NO additional issues are reported beyond what is expected

If any expected issue is missing, reply 'false'.
If any issue is reported that is NOT in the expected list, reply 'false'.

Reply with exactly 'true' or 'false' (no other text)."
```

Parse response, trim whitespace, check for `true` or `false`.

## Pre-flight Check

Before cloning, verify these local files exist in the source repo:

- `AGENTS.md`
- `.claude/commands/api-review.md`

These are copied into the temp clone so that any local modifications to the review command are tested, not the remote versions.

## Configuration

| Setting | Value |
|---------|-------|
| Timeout per Claude call | 5 minutes |
| Execution mode | Sequential |
| Clone depth | 1 (shallow) |
| Clone source | `https://github.com/openshift/api.git` |
| Reset between tests | Verify origin remote exists, `git reset --hard origin/master && git clean -fd`, re-copy local AGENTS.md and .claude/ |

---

## Phase 2

### Cost Tracking

Use `--output-format json` to capture `total_cost_usd` from each Claude invocation. Accumulate across all calls (review + judge) and print the total in `AfterSuite`.

### Test Structure Reorganization ✅ IMPLEMENTED

Reorganize `testdata/` into two categories:

```
tests/eval/testdata/
├── golden/                     # Base truth tests - single isolated issues
│   ├── missing-optional-doc/
│   │   ├── patch.diff          # Triggers ONLY missing-optional-doc
│   │   └── expected.txt
│   ├── undocumented-enum/
│   │   ├── patch.diff          # Triggers ONLY undocumented-enum
│   │   └── expected.txt
│   ├── missing-featuregate/
│   │   ├── patch.diff          # Triggers ONLY missing-featuregate
│   │   └── expected.txt
│   └── valid-api-change/
│       ├── patch.diff          # Triggers NO issues
│       └── expected.txt
└── integration/                # Complex scenarios - multiple issues
    ├── new-field-all-issues/
    │   ├── patch.diff          # Triggers multiple issues together
    │   └── expected.txt
    └── partial-documentation/
        ├── patch.diff
        └── expected.txt
```

**Golden tests**: Each patch is carefully crafted to trigger exactly one issue type. These validate that the review command correctly identifies individual issue categories in isolation.

**Integration tests**: Patches that trigger multiple issues, testing the review command's ability to identify combinations of problems in realistic scenarios.

### Model Selection ✅ IMPLEMENTED

Each test tier has a default model, overridable via environment variable:

| Test Type | Default Model | Override Env Var |
|-----------|---------------|------------------|
| Golden tests | Sonnet | `EVAL_GOLDEN_MODEL` |
| Integration tests | Opus | `EVAL_INTEGRATION_MODEL` |
| Judge LLM | Haiku | `EVAL_JUDGE_MODEL` |

The test suite reads these at startup and applies per-tier:

```go
goldenModel := getEnvOrDefault("EVAL_GOLDEN_MODEL", "claude-sonnet-4-5@20250929")
integrationModel := getEnvOrDefault("EVAL_INTEGRATION_MODEL", "claude-opus-4-5@20251101")
judgeModel := getEnvOrDefault("EVAL_JUDGE_MODEL", "claude-haiku-4-5-20251001")
```

Usage:
```bash
# Use defaults
go test ./tests/eval/...

# Override golden tests to use Haiku
EVAL_GOLDEN_MODEL=claude-3-haiku-20240307 go test ./tests/eval/...

# Override all models
EVAL_GOLDEN_MODEL=claude-3-haiku-20240307 \
EVAL_INTEGRATION_MODEL=claude-sonnet-4-20250514 \
go test ./tests/eval/...
```

### Patch Stability

Patches may fail to apply as `origin/master` evolves over time. Strategies:

- Pin to a specific commit SHA in the clone step
- Use `git apply --3way` for better conflict handling
- Periodic patch refresh CI job

### Error Handling

Current design does not address failure scenarios:

- Patch application failures
- Resource cleanup on test failures

Using `--output-format json` also enables better error handling in future phases:

- Claude CLI timeouts or crashes (detect via JSON parse failure or missing fields)
- Empty or malformed output (validate JSON structure)
- Authentication failures (check for error fields in JSON response)

### Performance Optimizations

The API review step is the slowest part of the eval suite. Options to improve:

1. **Skip linting by default** - Update api-review command to skip `make lint` unless explicitly requested. Linting adds significant time.

2. **Cache review outputs** - For development, cache the review output keyed by patch hash. Skip re-running if cached result exists. Clear cache on command changes.

3. **Parallel test execution** - Run golden tests in parallel (requires separate repo clones per test).

4. **Smaller/faster model for development** - Use Haiku for rapid iteration, Sonnet/Opus for CI validation.