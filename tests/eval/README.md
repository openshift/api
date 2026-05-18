# API Review Evals

Eval suite for the `/api-review` skill. Tests that the reviewer finds the right issues and nothing else.

## Structure

```
tests/eval/
├── eval_test.go                       # Test harness (ginkgo, build tag: eval)
├── testdata/
│   ├── golden/<name>/                 # Single-concern tests (run on sonnet)
│   │   ├── patch.diff                 # Diff applied to config/v1/types_console.go
│   │   └── expected.txt               # One expected issue per line
│   └── integration/<name>/            # Multi-issue tests (run on opus)
│       ├── patch.diff
│       └── expected.txt
```

## Running

```bash
make eval                # All evals (golden then integration)
make eval-golden         # Golden only
make eval-integration    # Integration only
```

### Environment variables

| Variable | Default | Purpose |
|---|---|---|
| `EVAL_GOLDEN_MODEL` | `sonnet` | Model for golden reviews |
| `EVAL_INTEGRATION_MODEL` | `opus` | Model for integration reviews |
| `EVAL_JUDGE_MODEL` | `haiku` | Model for judging review output |
| `EVAL_RUNS` | `1` | Runs per test case |
| `EVAL_THRESHOLD` | `0.8` | Minimum pass rate (0.0–1.0) |
| `EVAL_BASE_REF` | pinned SHA in `eval_test.go` | Git ref for worktrees |
| `EVAL_GOLDEN_PROCS` | `4` | Max parallel golden tests |
| `EVAL_INTEGRATION_PROCS` | `2` | Max parallel integration tests |
| `EVAL_GINKGO_ARGS` | (empty) | Extra ginkgo flags |
| `EVAL_VERBOSE` | unset | Show raw reviewer output |

## Judging

Both golden and integration tests use a judge LLM to compare reviewer output against `expected.txt`. Judging is semantic — wording doesn't need to match, just the underlying issue.

**Golden judge**: ALL expected issues must be found. ANY issue not in `expected.txt` fails the test.

**Integration judge**: At least 70% of expected issues must be found. ANY issue not in `expected.txt` still fails the test (precision over recall).

Both judges fail on false positives. If the reviewer can legitimately find an issue, it must be in `expected.txt`.

## Adding a test case

1. Create `testdata/{golden,integration}/<name>/patch.diff` targeting `config/v1/types_console.go`.
2. Create `testdata/{golden,integration}/<name>/expected.txt` with one issue per line.
3. Verify the patch applies: `git apply --check testdata/golden/<name>/patch.diff`
4. Run: `make eval-golden` or `make eval-integration`

## Patch baseline

All patches target `config/v1/types_console.go` at a pinned commit (`defaultEvalBaseRef` in `eval_test.go`). Tests create worktrees from this ref and apply patches there.

`git apply --3way` absorbs minor line shifts. Patches need rebasing when upstream diverges too far.

### Updating the baseline

1. Check which patches still apply:
   ```bash
   for d in tests/eval/testdata/golden/*/; do
     echo "=== $(basename $d) ===";
     git apply --3way --check "$d/patch.diff" 2>&1;
   done
   ```
2. Fix any that fail (regenerate context lines from the new base).
3. Run evals: `EVAL_BASE_REF=<new-sha> make eval`
4. Update `defaultEvalBaseRef` in `eval_test.go`.

## Footguns

### FeatureGate markers

Patches add fields to `ConsoleSpec`, a stable API. Opus reads CLAUDE.md and flags fields missing `+openshift:enable:FeatureGate=...` even though that's outside `/api-review` scope.

- **Integration patches**: must include a FeatureGate marker on top-level fields added to ConsoleSpec.
- **Golden patches**: don't need one. Sonnet stays scoped to the command's rules.

### Hunk line counts

The `@@` header line counts must match exactly. If you add a line to a hunk, update the count (e.g., `+33,12` → `+33,13`). Always run `git apply --check` after editing a patch.

### Pre-existing fields in context

The existing `Authentication` field in ConsoleSpec has no doc comment (just `// +optional`). Opus flags this as a pre-existing issue. Don't add it to `expected.txt` — it's a false positive from context, not from the patch.

### Model differences

- **Sonnet**: stays focused on the command's numbered rules. Reliable for mechanical checks. Misses judgment calls.
- **Opus**: pulls in CLAUDE.md and general conventions. Finds more but also generates false positives (FeatureGate, named types, design opinions). Worse for goldens, needed for integrations.
