---
name: api-review
description: Run strict OpenShift API review workflow for PR changes
parameters:
  - name: pr_url
    description: GitHub PR URL to review
    required: true
---

# Output Format Requirements
You MUST use this EXACT format for ALL review feedback:


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


I'll run a comprehensive API review for the OpenShift API changes in the specified GitHub PR.

## Step 1: Pre-flight checks and branch management

First, I'll check for uncommitted changes and save the current branch:

```bash
# Save current branch
CURRENT_BRANCH=$(git branch --show-current)
echo "📍 Current branch: $CURRENT_BRANCH"

# Check for uncommitted changes
if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "❌ ERROR: Uncommitted changes detected. Cannot proceed with API review."
    echo "Please commit or stash your changes before running the API review."
    git status --porcelain
    exit 1
fi

echo "✅ No uncommitted changes detected. Safe to proceed."
```

## Step 2: Extract PR number and checkout PR

```bash
# Extract PR number from URL
PR_NUMBER=$(echo "$ARGUMENTS" | grep -oE '[0-9]+$')
echo "🔍 Reviewing PR #$PR_NUMBER"

# Checkout the PR
echo "🔄 Checking out PR #$PR_NUMBER..."
gh pr checkout "$PR_NUMBER"

# Get changed Go files in API directories
echo "📁 Analyzing changed files..."
gh pr view "$PR_NUMBER" --json files --jq '.files[].path' | grep '\.go$' | grep -E '/(v1|v1alpha1|v1beta1)/'
```

## Step 3: Run linting checks on PR changes

```bash
echo "⏳ Running linting checks on PR changes..."
make lint
```

## Step 4: Documentation validation

For each changed API file, I'll validate:

1. **Field Documentation**: All struct fields must have documentation comments
2. **Optional Field Behavior**: Optional fields must explain what happens when they are omitted
3. **Validation Documentation**: Validation rules must be documented and match markers

Let me check each changed file for these requirements:

```thinking
I need to analyze the changed files to:
1. Find struct fields without documentation
2. Find optional fields without behavior documentation
3. Find validation annotations without corresponding documentation

For each Go file, I'll:
- Look for struct field definitions
- Check if they have preceding comment documentation
- For optional fields (those with `+kubebuilder:validation:Optional` or `+optional`), verify behavior is explained
- For fields with validation annotations, ensure the validation is documented
```

## Step 5: Generate comprehensive review report

I'll provide a comprehensive report showing:
- ✅ Files that pass all checks
- ❌ Files with documentation issues
- 📋 Specific lines that need attention
- 📚 Guidance on fixing any issues

The review will fail if any documentation requirements are not met for the changed files.

## Step 6: Switch back to original branch

After completing the review, I'll switch back to the original branch:

```bash
echo "🔄 Switching back to original branch: $CURRENT_BRANCH"
git checkout "$CURRENT_BRANCH"
echo "✅ API review complete. Back on branch: $(git branch --show-current)"
```

**CRITICAL WORKFLOW REQUIREMENTS:**
1. MUST check for uncommitted changes before starting
2. MUST abort if uncommitted changes are detected
3. MUST save current branch name before switching
4. MUST checkout the PR before running `make lint`
5. MUST switch back to original branch when complete
6. If any step fails, MUST attempt to switch back to original branch before exiting
