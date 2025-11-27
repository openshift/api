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
│   ├── missing-optional-doc/
│   │   ├── patch.diff        # Git patch to apply
│   │   └── expected.txt      # Expected issues (one per line)
│   ├── undocumented-enum/
│   │   ├── patch.diff
│   │   └── expected.txt
│   └── valid-api-change/
│       ├── patch.diff
│       └── expected.txt      # Empty file = no issues expected
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

## Phase 2 (Future Work)

### Patch Stability

Patches may fail to apply as `origin/master` evolves over time. Need a strategy to handle this (e.g., pinning to a specific commit).

### Error Handling

Current design does not address failure scenarios:

- Patch application failures
- Claude CLI timeouts or crashes
- Empty or malformed output from Claude
- Authentication failures
- Resource cleanup on test failures