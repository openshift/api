---
name: api-review
description: Run strict OpenShift API review workflow for PR changes or local changes
parameters:
  - name: pr_url
    description: GitHub PR URL to review (optional - if not provided, reviews local changes against upstream master)
    required: false
---

## Step 1: Run preflight checks

Run the preflight script to discover changed API files and run linting:

```bash
bash .claude/skills/api-review/scripts/preflight.sh $ARGUMENTS
```

Read the output:
- If it says `NO_API_FILES_CHANGED`, respond "No API files changed. Nothing to review." and stop.
- If lint says `FAILED`, report the lint failures and stop.
- The `Changed API Files` section lists the files to review in Step 2.
- If the mode is `pr`, note the `Original-Branch` for cleanup in Step 3.

## Step 2: Documentation validation

**CRITICAL: Only review new or modified lines (the `+` lines in the diff). Do NOT flag pre-existing issues in unchanged context lines. There is significant tech debt in existing APIs and reviewing it is out of scope.**

For each changed API file listed in the preflight output, validate only the new/modified lines for:

1. **Field Documentation**: New struct fields must have documentation comments
2. **Optional Field Behavior**: New optional fields must explain what happens when the field is omitted (not provided). This is about omission, not about empty values — empty strings or empty lists are a separate concern handled by validation markers.
3. **Validation Documentation**: New validation rules must be documented and match markers. For `+kubebuilder:validation:Enum` fields, each enum value must be listed AND its meaning explained in the **field's own comment** using the "When set to X, ..." pattern (e.g., "When set to Vault, HashiCorp Vault is used as the secret store"). Simply listing values without explaining what they do is insufficient. Note: explanations on the type definition do NOT satisfy this rule — the field comment is what users see in generated docs.
4. **Validation / Documentation Mismatch**: This rule covers two directions:
   - **Markers contradict docs**: A validation marker prevents behavior the comment describes. For example, `+kubebuilder:validation:MinItems=1` is set but the comment says "an empty list means no items are excluded" — the validation prevents the documented behavior. Note: "omitted" (field not provided) is different from "empty" (field present with zero items or empty string). `MinItems=1` preventing empty lists does NOT contradict documentation about omitted behavior — only flag when documentation describes behavior for a value that validation markers explicitly prevent.
   - **Docs claim constraints with no enforcement**: The comment describes a constraint but there is no validation marker (`Pattern`, `XValidation`, `MaxLength`, etc.) that enforces it. The API will silently accept values that violate the documented constraint. Either add the marker or remove the claim. Only flag when enforcement is completely absent — do not flag when a marker enforces the documented constraint but you think the enforcement could be stricter or more thorough.
5. **Missing Cross-field Validation**: When documentation states a relationship between fields (e.g., "mutually exclusive with FieldX", "required when FieldY is set", "cannot be used together with FieldZ"), there MUST be a corresponding `+kubebuilder:validation:XValidation` rule enforcing that relationship. Documentation alone is not sufficient — the cluster will not enforce undocumented-in-code relationships.
6. **Undocumented Constraints**: ALL kubebuilder constraint markers (`MinLength`, `MaxLength`, `MinItems`, `MaxItems`, `Minimum`, `Maximum`, `MaxProperties`, `Pattern`) MUST be documented in the field's comment. If a field has `+kubebuilder:validation:MinLength=5` but the comment does not mention the minimum length requirement, that is an issue. This includes `MinLength=1` on required fields — even though "required + MinLength=1" may seem redundant, the MinLength constraint is separately enforced at the validation layer and must be documented so users know empty strings are rejected. Check every field with constraint markers, including fields on supporting/referenced types added in the same diff (e.g., a `SecretNameReference` type with a `name` field). Exception: `MinProperties` does not need to be documented — it is a structural constraint ("don't send an empty object") that is self-evident from the type definition.
7. **CEL Expression Review**: For `+kubebuilder:validation:XValidation` rules, check that CEL expressions are logically correct: no unreachable branches, correct enum value references that match the actual `+kubebuilder:validation:Enum` values, proper use of `has()` guards for optional fields, and no tautological or contradictory conditions. Prefer combined ternary expressions for related required/forbidden checks (e.g., `self.type == 'X' ? has(self.x) : !has(self.x)`). Flag verbose or overly complex CEL that could be simplified.

**IMPORTANT: Only report issues that violate one of the 7 numbered rules above. Do not report design suggestions, code cleanup, type choice recommendations (e.g. `*int32` vs `int32`), missing `enhancementPR` calls, or best practices that fall outside these rules. A false positive is worse than a missed issue.**

```thinking
For EACH changed field (including fields on supporting types introduced in the diff), check ALL of the following:
1. Field documentation present?
2. Optional fields explain omitted behavior?
3. Enum values explained in the FIELD's own comment (not just the type definition)?
4. Any validation marker contradicting the comment prose? Any comment claiming a constraint that has no marker to enforce it?
5. Documented field relationships enforced with XValidation?
6. ALL constraint markers (MinLength, MaxLength, MinItems, MaxItems, Minimum, Maximum, MaxProperties, Pattern) documented in comment? Including MinLength=1? (Exception: MinProperties does not need documentation)
7. CEL expressions logically correct? Prefer combined ternary for required/forbidden checks?
```

## Step 3: Cleanup (PR mode only)

If the preflight output showed `Mode: pr`, switch back to the original branch:

```bash
git checkout <Original-Branch from preflight output>
```

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