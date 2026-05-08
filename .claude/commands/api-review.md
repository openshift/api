---
name: api-review
description: Run strict OpenShift API review workflow for PR changes or local changes
parameters:
  - name: pr_url
    description: GitHub PR URL to review (optional - if not provided, reviews local changes against upstream master)
    required: false
---

I'll run a comprehensive API review for OpenShift API changes. This can review either a specific GitHub PR or local changes against upstream master.

## Step 1: Pre-flight checks and determine review mode

First, I'll check the arguments and determine whether to review a PR or local changes:

```bash
# Save current branch
CURRENT_BRANCH=$(git branch --show-current)
echo "📍 Current branch: $CURRENT_BRANCH"

# Check if a PR URL was provided
if [ -n "$ARGUMENTS" ] && [[ "$ARGUMENTS" =~ github\.com.*pull ]]; then
    REVIEW_MODE="pr"
    PR_NUMBER=$(echo "$ARGUMENTS" | grep -oE '[0-9]+$')
    echo "🔍 PR review mode: Reviewing PR #$PR_NUMBER"

    # For PR review, check for uncommitted changes
    if ! git diff --quiet || ! git diff --cached --quiet; then
        echo "❌ ERROR: Uncommitted changes detected. Cannot proceed with PR review."
        echo "Please commit or stash your changes before running the API review."
        git status --porcelain
        exit 1
    fi
    echo "✅ No uncommitted changes detected. Safe to proceed with PR review."
else
    REVIEW_MODE="local"
    echo "🔍 Local review mode: Reviewing local changes against upstream master"

    # Find a remote pointing to openshift/api repository
    OPENSHIFT_REMOTE=""
    for remote in $(git remote); do
        remote_url=$(git remote get-url "$remote" 2>/dev/null || echo "")
        if [[ "$remote_url" =~ github\.com[/:]openshift/api(\.git)?$ ]]; then
            OPENSHIFT_REMOTE="$remote"
            echo "✅ Found OpenShift API remote: '$remote' -> $remote_url"
            break
        fi
    done

    # If no existing remote found, add upstream
    if [ -z "$OPENSHIFT_REMOTE" ]; then
        echo "⚠️  No remote pointing to openshift/api found. Adding upstream remote..."
        git remote add upstream https://github.com/openshift/api.git
        OPENSHIFT_REMOTE="upstream"
    fi

    # Fetch latest changes from the OpenShift API remote
    echo "🔄 Fetching latest changes from $OPENSHIFT_REMOTE..."
    git fetch "$OPENSHIFT_REMOTE" master
fi
```

## Step 2: Get changed files based on review mode

```bash
if [ "$REVIEW_MODE" = "pr" ]; then
    # PR Review: Checkout the PR and get changed files
    echo "🔄 Checking out PR #$PR_NUMBER..."
    gh pr checkout "$PR_NUMBER"

    echo "📁 Analyzing changed files in PR..."
    CHANGED_FILES=$(gh pr view "$PR_NUMBER" --json files --jq '.files[].path' | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/')
else
    # Local Review: Get changed files compared to openshift remote master
    echo "📁 Analyzing locally changed files compared to $OPENSHIFT_REMOTE/master..."
    CHANGED_FILES=$(git diff --name-only "$OPENSHIFT_REMOTE/master...HEAD" | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/')

    # Also include staged changes
    STAGED_FILES=$(git diff --cached --name-only | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/' || true)
    if [ -n "$STAGED_FILES" ]; then
        CHANGED_FILES=$(echo -e "$CHANGED_FILES\n$STAGED_FILES" | sort -u)
    fi
fi

echo "Changed API files:"
echo "$CHANGED_FILES"

if [ -z "$CHANGED_FILES" ]; then
    echo "ℹ️  No API files changed. Nothing to review."
    if [ "$REVIEW_MODE" = "pr" ]; then
        git checkout "$CURRENT_BRANCH"
    fi
    exit 0
fi
```

## Step 3: Run linting checks on changes

```bash
echo "⏳ Running linting checks on changes..."
make lint

if [ $? -ne 0 ]; then
    echo "❌ Linting checks failed. Please fix the issues before proceeding."
    if [ "$REVIEW_MODE" = "pr" ]; then
        echo "🔄 Switching back to original branch: $CURRENT_BRANCH"
        git checkout "$CURRENT_BRANCH"
    fi
    exit 1
fi

echo "✅ Linting checks passed."
```

## Step 4: Documentation validation

**CRITICAL: Only review new or modified lines (the `+` lines in the diff). Do NOT flag pre-existing issues in unchanged context lines. There is significant tech debt in existing APIs and reviewing it is out of scope.**

For each changed API file, I'll validate only the new/modified lines for:

1. **Field Documentation**: New struct fields must have documentation comments
2. **Optional Field Behavior**: New optional fields must explain what happens when the field is omitted (not provided). This is about omission, not about empty values — empty strings or empty lists are a separate concern handled by validation markers.
3. **Validation Documentation**: New validation rules must be documented and match markers. For `+kubebuilder:validation:Enum` fields, each enum value must be listed AND its meaning explained using the "When set to X, ..." pattern (e.g., "When set to Vault, HashiCorp Vault is used as the secret store"). Simply listing values without explaining what they do is insufficient.
4. **Contradicting Validation vs Documentation**: Check if validation markers contradict the comment prose. For example, if `+kubebuilder:validation:MinItems=1` is set but the comment says "an empty list means no items are excluded", the validation prevents the behavior the documentation describes. The validation markers are the source of truth for runtime behavior, so the documentation must accurately describe what the validation permits. Note: "omitted" (field not provided) is different from "empty" (field present with zero items or empty string). `MinItems=1` preventing empty lists does NOT contradict documentation about omitted behavior — only flag a contradiction when documentation describes behavior for a value that the validation markers explicitly prevent.
5. **Missing Cross-field Validation**: When documentation states a relationship between fields (e.g., "mutually exclusive with FieldX", "required when FieldY is set", "cannot be used together with FieldZ"), there MUST be a corresponding `+kubebuilder:validation:XValidation` rule enforcing that relationship. Documentation alone is not sufficient — the cluster will not enforce undocumented-in-code relationships.
6. **Undocumented Constraints**: ALL kubebuilder constraint markers (`MinLength`, `MaxLength`, `MinItems`, `MaxItems`, `Minimum`, `Maximum`, `MaxProperties`, `Pattern`) MUST be documented in the field's comment. If a field has `+kubebuilder:validation:MinLength=5` but the comment does not mention the minimum length requirement, that is an issue. Exception: `MinProperties` does not need to be documented — it is a structural constraint (\"don't send an empty object\") that is self-evident from the type definition.
7. **CEL Expression Review**: For `+kubebuilder:validation:XValidation` rules, check that CEL expressions are logically correct: no unreachable branches, correct enum value references that match the actual `+kubebuilder:validation:Enum` values, proper use of `has()` guards for optional fields, and no tautological or contradictory conditions. Prefer combined ternary expressions for related required/forbidden checks (e.g., `self.type == 'X' ? has(self.x) : !has(self.x)`). Flag verbose or overly complex CEL that could be simplified.

**IMPORTANT: Only report issues that violate one of the 7 numbered rules above. Do not report design suggestions, code cleanup, type choice recommendations (e.g. `*int32` vs `int32`), missing `enhancementPR` calls, or best practices that fall outside these rules. A false positive is worse than a missed issue.**

```thinking
For EACH changed field, check ALL of the following:
1. Field documentation present?
2. Optional fields explain omitted behavior?
3. Validation markers documented and match?
4. Any validation marker contradicting the comment prose?
5. Documented field relationships enforced with XValidation?
6. ALL constraint markers (MinLength, MaxLength, MinItems, MaxItems, Minimum, Maximum, MinProperties, MaxProperties, Pattern) documented in comment?
7. CEL expressions logically correct? Prefer combined ternary for required/forbidden checks?
```

## Step 5: Generate comprehensive review report

I'll provide a comprehensive report showing:
- ✅ Files that pass all checks
- ❌ Files with documentation issues
- 📋 Specific lines that need attention
- 📚 Guidance on fixing any issues

The review will fail if any documentation requirements are not met for the changed files.

## Step 6: Switch back to original branch (PR mode only)

After completing the review, if we were reviewing a PR, I'll switch back to the original branch:

```bash
if [ "$REVIEW_MODE" = "pr" ]; then
    echo "🔄 Switching back to original branch: $CURRENT_BRANCH"
    git checkout "$CURRENT_BRANCH"
    echo "✅ API review complete. Back on branch: $(git branch --show-current)"
else
    echo "✅ Local API review complete."
fi
```

**CRITICAL WORKFLOW REQUIREMENTS:**

**For PR Review Mode:**
1. MUST check for uncommitted changes before starting
2. MUST abort if uncommitted changes are detected
3. MUST save current branch name before switching
4. MUST checkout the PR before running `make lint`
5. MUST switch back to original branch when complete
6. If any step fails, MUST attempt to switch back to original branch before exiting

**For Local Review Mode:**
1. MUST detect existing remotes pointing to openshift/api repository (supports any remote name)
2. MUST add upstream remote only if no existing openshift/api remote is found
3. MUST fetch latest changes from the detected openshift/api remote
4. MUST compare against the detected remote's master branch
5. MUST include both committed and staged changes in analysis
6. No branch switching required since we're reviewing local changes

## FINAL OUTPUT — MANDATORY

**After completing all analysis, you MUST produce a text response listing every issue found.** If you do not output text, the review is lost. Do not end on a tool call.

Use this EXACT format for EACH issue:

+LineNumber: Brief description
**Current (problematic) code:**
```go
[exact code from the PR diff]
```

**Suggested change:**
```diff
- [old code line]
+ [new code line]
```

**Explanation:** [Why this change is needed]

Every issue must be enumerated individually. Do NOT summarize into tables or counts. If no issues are found, say "No issues found."