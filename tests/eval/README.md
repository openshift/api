# API Review Evals

Eval suite for the `/api-review` Claude command. Golden tests isolate single review concerns; integration tests combine multiple issues in realistic code modeled on real PR review findings.

## Running

```bash
# Golden tests (single-concern, fast)
make eval-golden

# Integration tests (multi-issue, slower)
make eval-integration

# All evals
make eval
```

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `EVAL_GOLDEN_MODEL` | `sonnet` | Model for golden reviews |
| `EVAL_INTEGRATION_MODEL` | `opus` | Model for integration reviews |
| `EVAL_JUDGE_MODEL` | `haiku` | Model for judging review output |
| `EVAL_RUNS` | `1` | Number of runs per test case |
| `EVAL_THRESHOLD` | `0.8` | Minimum pass rate (0.0–1.0) |
| `EVAL_BASE_REF` | pinned SHA in `eval_test.go` | Git ref to create worktrees from |
| `EVAL_VERBOSE` | unset | Show raw reviewer output when set |

#### Thinking/effort tuning

These are inherited by the Claude subprocess:

```bash
CLAUDE_CODE_EFFORT_LEVEL=max \
CLAUDE_CODE_DISABLE_AUTO_MEMORY=1 \
CLAUDE_CODE_SIMPLE_SYSTEM_PROMPT=1 \
CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1 \
make eval-golden
```

## Patch Baseline

All patches target `config/v1/types_console.go` and depend on specific file content for their context lines. Upstream merges that modify this file shift line numbers and change context, breaking patch application. To avoid this, patches are pinned to a specific commit via `defaultEvalBaseRef` in `eval_test.go`.

Patches are applied with `git apply --3way`, which absorbs minor line shifts from upstream changes. Eventually patches need rebasing when upstream diverges too far.

### Updating the baseline

1. Check which patches still apply against the new base:
   ```bash
   for d in tests/eval/testdata/golden/*/; do
     echo "=== $(basename $d) ===";
     git apply --3way --check "$d/patch.diff" 2>&1;
   done
   ```
2. Fix any patches that fail (regenerate context lines from the new base).
3. Run evals to confirm pass rates: `EVAL_BASE_REF=upstream/master make eval-golden`
4. Update `defaultEvalBaseRef` in `eval_test.go` to the new commit SHA.

## Adding a new golden test

1. Create `testdata/golden/<test-name>/patch.diff` targeting `config/v1/types_console.go`.
2. Create `testdata/golden/<test-name>/expected.txt` with one expected issue per line.
3. Verify the patch applies: `git apply --check testdata/golden/<test-name>/patch.diff`
4. Run the single test to confirm it passes at the desired threshold.

## Writing test patches — known footguns

### FeatureGate markers

All patches add fields to `ConsoleSpec`, a stable (level 1) API. CLAUDE.md says new stable fields need a FeatureGate. Opus reads CLAUDE.md and will flag every new field that lacks `+openshift:enable:FeatureGate=...`, even though FeatureGate checking isn't one of the `/api-review` rules.

**Integration patches** must include a FeatureGate marker on the top-level field added to ConsoleSpec. Without it, opus will report a false positive that poisons the whole review.

**Golden patches** do not need FeatureGate markers. Goldens run on sonnet, which stays scoped to the command's rules and doesn't pull in CLAUDE.md guidance.

### Hunk line counts

When editing a patch, the `@@` header line counts must match exactly. If you add a line to a hunk, update the count (e.g., `+33,12` becomes `+33,13`). `git apply --check` catches this — always run it after editing.

### Pre-existing fields in context

The existing `Authentication` field in ConsoleSpec has no doc comment (just `// +optional`). Opus will flag this as a pre-existing issue despite the command saying "only review new or modified lines." Don't add it to expected issues — it's a false positive.

### Judge strictness

The judge fails the test if the reviewer reports ANY issue not in the expected list. This means:
- Every legitimate issue the reviewer could find must be in `expected.txt`
- Issues from removed rules must not exist in the patch code (remove the problematic code, don't just remove the expected issue)
- Opus is more aggressive than sonnet — test with the model the eval actually uses

### Model differences

- **Sonnet**: Stays focused on the command's numbered rules. Reliable for mechanical checks (rules 1-7). Misses judgment calls.
- **Opus**: Pulls in guidance from CLAUDE.md and general conventions. Finds more issues but also generates false positives (FeatureGate, named types, design opinions). Worse for goldens, needed for integrations.
