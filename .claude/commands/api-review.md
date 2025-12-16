---
name: api-review
description: Run strict OpenShift API review workflow for PR changes or local changes
parameters:
  - name: pr_url
    description: GitHub PR URL to review (optional - if not provided, reviews local changes against upstream master)
    required: false
---

I will assume the role of an experienced software engineer. I will perform a comprehensive review of a proposed OpenShift API change, which may be either a specific GitHub PR or local changes against upstream master.

## Step 1: Pre-flight checks and determine review mode

I will directly execute the `pre-flight-checks.sh` script from this command, passing any command arguments as arguments to the script.

The output of this script clearly indicates:
* The changed API files I must review
* Whether the lint check passed

Do not continue if:
* An error occurred
* There are no files to review
* The lint check failed

## Step 2: Documentation and Validation Review

Consult **AGENTS-validations.md** for detailed guidance:
- **Quick Decision Trees** - Fast lookup for validation method selection
- **Best Practices** - Validation patterns, documentation requirements, combining rules
- **Common Errors & Solutions** - Problems to detect and fix during review
- **Acceptable Patterns** - Patterns that should NOT be flagged as issues
- **Test Requirements** - Test coverage rules and requirements
- **Complete Validator Catalog** - Reference for all Format markers and CEL validators

I will now analyze each changed API file according to the requirements in AGENTS-validations.md:

### Documentation Requirements (1.3 Documentation Requirements Checklist)
- All fields have documentation comments (2.4 Documentation Requirements)
- Optional fields explain behavior when omitted (2.4 Documentation Requirements, 4.4 "This field is optional")
- Validation rules are documented in human-readable form (2.4 Documentation Requirements, 3.5 Missing Documentation for Validation Constraints)
- Cross-field relationships are both documented AND enforced with XValidation (3.6 Cross-Field Constraints Without Enforcement, 3.8 Cross-Field Validation Placement)
- XValidation rules have accurate message fields (3.9 Inaccurate or Inconsistent Validation Messages)

### Validation Requirements (Section 2: Best Practices)
- Follow validation decision hierarchy (2.1 Validation Decision Hierarchy)
- Prohibited format markers not used (7.1.0 Security: Prohibited Format Markers)
- Cross-field validation placed on parent structs (3.8 Cross-Field Validation Placement)
- Validation messages use exact enum values and JSON field names (3.9 Inaccurate or Inconsistent Validation Messages)

## Step 3: Comprehensive Test Coverage Analysis

### Process Reference:

**CRITICAL:** Execute this for *all* validation markers, not just those with other documentation issues.

Execute all 7 steps from Section 5.0.1:
1. Extract validation markers from changed API files
2. Categorize using Section 5.0 lookup table
3. Locate test files in `<group>/<version>/tests/<crd-name>/`
4. Map validations to existing tests
5. Verify coverage completeness (minimal vs comprehensive)
6. Identify gaps (missing, insufficient, unnecessary)
7. Report using standard format from 5.0.1

## Step 4: Generate comprehensive review report

I'll provide a comprehensive report showing:
- ✅ Files that pass all checks
- ❌ Files with documentation issues
- 📋 Specific lines that need attention
- 🧪 Any API validations with insufficient test coverage
- 🧪 Any redundant test coverage
- 📚 Guidance on fixing any issues

**CRITICAL:** You MUST use this EXACT format for documentation issues, validation problems, and lint errors (but not test coverage):

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

**Test coverage output must follow the Standard Test Coverage Report Format** defined in Section 5.0.2.

**Do not provide output for anything that does not require corrective action.**

The review will fail if any documentation or test coverage requirements are not met for the changed files.

## Step 5: Switch back to original branch (PR mode only)

After completing the review, if we were reviewing a PR, I'll switch back to the original branch:

```bash
cat > /tmp/api_review_step5.sh << 'STEP5_EOF'
#!/bin/bash
source /tmp/api_review_vars.sh

if [ "$REVIEW_MODE" = "pr" ]; then
    echo "🔄 Switching back to original branch: $CURRENT_BRANCH"
    git checkout "$CURRENT_BRANCH"
    echo "✅ API review complete. Back on branch: $(git branch --show-current)"
else
    echo "✅ Local API review complete."
fi
STEP5_EOF

bash /tmp/api_review_step5.sh
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