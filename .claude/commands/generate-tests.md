---
name: generate-tests
description: Generate comprehensive integration test suites for OpenShift API type definitions
parameters:
  - name: path
    description: Path to a types file or API group/version directory (e.g., config/v1/types_infrastructure.go or operator/v1)
    required: true
---

# Generate Integration Tests for OpenShift API Types

Generate `.testsuite.yaml` integration test files for the specified API types.
Tests validate CRD schemas using the envtest-based integration test runner.

## Step 1: Prechecks

```bash
REPO_ROOT=$(git rev-parse --show-toplevel)
TARGET_PATH="$ARGUMENTS"
if [ -z "$TARGET_PATH" ]; then
  echo "PRECHECK FAILED: No target path provided."
  echo "Usage: /generate-tests <path-to-types-file-or-api-directory>"
  echo "Examples:"
  echo "  /generate-tests config/v1/types_infrastructure.go"
  echo "  /generate-tests operator/v1"
  echo "  /generate-tests route/v1/types_route.go"
  exit 1
fi

echo "Repository root: $REPO_ROOT"
echo "Target: $TARGET_PATH"
```

```bash
if [ -d "$TARGET_PATH" ]; then
  API_DIR="$TARGET_PATH"
else
  API_DIR=$(dirname "$TARGET_PATH")
fi
```

## Step 2: Identify Target Types and CRDs

```thinking
I need to determine the target types:
1. If the user provided a specific types file → read it directly
2. If the user provided an API group/version directory → find all types files

For each types file I must extract:
- The API group (from the package or doc.go: e.g., config.openshift.io)
- The version (from the directory: e.g., v1, v1alpha1)
- All CRD types (structs with +kubebuilder:object:root=true or +genclient)
- For each CRD: kind, resource plural, scope (cluster vs namespaced), singleton name constraints
- Every spec and status field with types, validation markers, godoc
- Enum types and their allowed values
- Discriminated unions (discriminator + members)
- Immutable fields (XValidation rules with oldSelf)
- Default values
- Feature gate annotations (+openshift:enable:FeatureGate=...)
- CEL validation rules (+kubebuilder:validation:XValidation)
- All constraints: min/max, minLength/maxLength, minItems/maxItems, pattern, format

I also need to check zz_generated.featuregated-crd-manifests/ to understand:
- Which CRD variants exist (AAA_ungated.yaml, <FeatureGate>.yaml, compound gates like A+B.yaml)
- The full CRD name (<plural>.<group>)
- Which featureGates control which variants
- The OpenAPI v3 schema for each variant
```

Read the types file(s) and extract all information. Also read the corresponding CRD manifests
from `zz_generated.featuregated-crd-manifests/<plural>.<group>/` to get:
- The list of CRD variants (ungated, per-featureGate, compound)
- The full CRD name
- Default values applied by the schema
- The full OpenAPI v3 validation tree

## Step 3: Understand Existing Tests

Check the `tests/` directory for existing test suites:

```bash
# List existing test files for this CRD
TEST_DIR="$API_DIR/tests"
if [ -d "$TEST_DIR" ]; then
  echo "=== Existing test directories ==="
  ls -la "$TEST_DIR/"
fi
```

If test files already exist for this CRD, read them to understand existing coverage.
Do NOT duplicate tests that already exist. Only add tests for NEW or MODIFIED fields/types.

## Step 4: Generate Test Suites

For each CRD variant found in `zz_generated.featuregated-crd-manifests/<plural>.<group>/`,
generate a corresponding test file at:

```text
<group>/<version>/tests/<plural>.<group>/<VariantName>.yaml
```

Where `<VariantName>` matches the CRD manifest filename (e.g., `AAA_ungated.yaml`,
`OnPremDNSRecords.yaml`, `FeatureA+FeatureB.yaml`).

### Test File Format

Every test file MUST follow this exact structure:

```yaml
apiVersion: apiextensions.k8s.io/v1 # Hack because controller-gen complains if we don't have this
name: "<Kind>"
crdName: <plural>.<group>
```

For feature-gated variants, add the gate(s):

```yaml
# Single feature gate:
featureGate: <FeatureGateName>

# Multiple feature gates:
featureGates:
- <FeatureGate1>
- <FeatureGate2>

# Negative gates (enabled in ungated but not when this gate is active):
featureGates:
- -<FeatureGateName>
```

Then add the test sections:

```yaml
tests:
  onCreate:
    - name: <test name>
      initial: |
        <yaml resource>
      expected: |
        <yaml resource with defaults applied>
    - name: <test name for error case>
      initial: |
        <yaml resource with invalid values>
      expectedError: "<error substring>"
  onUpdate:
    - name: <test name>
      initial: |
        <yaml resource>
      updated: |
        <yaml resource with changes>
      expected: |
        <yaml resource with expected result>
    - name: <test name for update error>
      initial: |
        <yaml resource>
      updated: |
        <yaml resource with invalid change>
      expectedError: "<error substring>"
    - name: <test name for status error>
      initial: |
        <yaml resource>
      updated: |
        <yaml resource with invalid status change>
      expectedStatusError: "<error substring>"
```

For validation ratcheting tests (where a previously invalid value is persisted and we test
update behavior), use `initialCRDPatches`:

```yaml
    - name: <ratcheting test name>
      initialCRDPatches:
        - op: remove
          path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/<field>/x-kubernetes-validations
      initial: |
        <yaml with previously invalid value>
      updated: |
        <yaml with change>
      expected: |
        <yaml with expected result>
```

### Test Priority

**Focus on custom validation logic.** The primary value of integration tests is verifying
custom behavior that Kubernetes does not enforce out of the box — CEL rules, cross-field
constraints, immutability, discriminated unions, defaulting, and ratcheting.

**Do NOT generate tests for basic OpenAPI schema enforcement.** Testing that an enum rejects
an unlisted value, that a string respects maxLength, or that a required field is required
adds little value — Kubernetes already enforces these. Only test OpenAPI schema constraints
when they interact with custom logic (e.g., a CEL rule that references an enum field).

### Test Categories

Generate tests covering the following categories that apply to the target types.
Derive specific test values from the actual Go types, markers, and CRD schema.

#### 1. Minimal Valid Create

Every test suite MUST include at least one test that creates a minimal valid instance:

```yaml
- name: Should be able to create a minimal <Kind>
  initial: |
    apiVersion: <group>/<version>
    kind: <Kind>
    spec: {} # No spec is required for a <Kind>
  expected: |
    apiVersion: <group>/<version>
    kind: <Kind>
    spec: {}
```

**IMPORTANT**: The `expected` block must include any default values the CRD schema applies.
Read the CRD manifest's OpenAPI schema to determine what defaults are applied to an empty spec.
For example, if a field has `default: Normal`, the expected output must include that field.

#### 2. CEL Validation Rules

This is the **highest-value** category. For every `+kubebuilder:validation:XValidation` rule:
- Test the valid case → `expected`
- Test each way the rule can be violated → `expectedError`
- For CEL rules with `oldSelf`: test update scenarios
- For multi-field CEL rules: test each field combination that triggers/passes the rule

#### 3. Discriminated Unions

If the type has discriminated unions (one-of pattern):
- Test each valid discriminator + required member → `expected`
- Test missing required member → `expectedError`
- Test forbidden member when discriminator doesn't match → `expectedError`

#### 4. Immutability and Update Constraints

For fields with immutability rules (CEL rules referencing `oldSelf`):
- Test that the value can be set initially → `expected`
- Test that changes are rejected → `expectedError` or `expectedStatusError`
- Test that retaining the same value on update succeeds → `expected`

#### 5. Defaulting Behavior

For fields with `default:` values in the CRD schema:
- Verify correct defaults are applied on creation
- Test interactions between defaults and custom validation

#### 6. Status Subresource

If the type has status fields with custom validation:
- Test valid status updates
- Test invalid status updates → `expectedStatusError`

#### 7. Feature-Gated Fields

For fields gated behind a FeatureGate:
- In the ungated (`AAA_ungated.yaml`) test: the field should NOT be available (if applicable)
- In the feature-gate-specific test file: test that the field works correctly

#### 8. Singleton Name Validation

If the CRD is a cluster-scoped singleton with a name constraint:
- Test creation with the required name → success
- Test creation with a different name → `expectedError`

#### 9. Validation Ratcheting

For fields where validation rules have been tightened:
- Use `initialCRDPatches` to relax validation on the initial object
- Create an initial object with the old (now invalid) value
- Test that updating OTHER fields while retaining the invalid value succeeds (ratcheting)
- Test that changing the invalid value to another invalid value fails

#### 10. Additional Coverage

```thinking
I have generated tests for the categories above. Now I must re-examine every CEL rule,
cross-field constraint, custom annotation etc. from the types and ask: is there any
validation behavior or edge case NOT already covered?

Examples:
- Cross-field dependencies (field B required only when field A is set)
- Mutually exclusive fields outside formal unions
- Nested object custom validation (embedded structs with CEL rules)
- Complex multi-field CEL rules
- Defaulting interactions between fields
- Metadata constraints (if enforced by the CRD)
- Edge cases around zero values vs nil for pointer fields

For each uncovered scenario I find, I will generate appropriate test cases.
Do NOT add tests for basic OpenAPI schema enforcement (enum, maxLength, pattern, min/max)
unless they are very critical to add.
```

## Step 5: Write Test Files

Write the generated `.testsuite.yaml` files to the correct location:

```text
<group>/<version>/tests/<plural>.<group>/<VariantName>.yaml
```

If a test file already exists for a variant, read it and **append** new tests — do NOT
overwrite or duplicate existing tests. Preserve all existing tests exactly as they are.

If a test file does not exist, create it with the full test suite.

Also ensure the group/version has a `Makefile` for running tests:

```bash
MAKEFILE="$API_DIR/Makefile"
if [ ! -f "$MAKEFILE" ]; then
  echo "NOTE: No Makefile found in $API_DIR. You can create one with:"
  echo "  make -C ../../tests test GINKGO_EXTRA_ARGS=--focus=\"<group>/<version>\""
fi
```

## Step 6: Verify

After writing the test files, run verification:

```bash
# Quick syntax check — ensure test files are valid YAML
for f in $(find "$API_DIR/tests" -name '*.yaml' -type f); do
  if ! python3 -c "import yaml; yaml.safe_load(open('$f'))" 2>/dev/null; then
    echo "WARNING: Invalid YAML in $f"
  fi
done

echo ""
echo "=== Next Steps ==="
echo "1. Review the generated test suites for correctness"
echo "2. Run targeted integration tests:"
echo "   make -C $API_DIR test"
echo "3. Run full verification:"
echo "   make verify"
echo "4. If you modified types, regenerate CRDs first:"
echo "   make update-codegen API_GROUP_VERSIONS=<group>/<version>"
```

## Step 7: Summary

Output a summary of all generated tests:

```text
=== Integration Test Generation Summary ===

Target: <group>/<version>
CRD: <plural>.<group>

Generated/Updated Test Files:
  - <path/to/test1.yaml> — <N> onCreate tests, <M> onUpdate tests
  - <path/to/test2.yaml> — <N> onCreate tests, <M> onUpdate tests

Test Coverage:
  onCreate:
    - Minimal valid create ✓
    - <field>: valid values (<count> tests)
    - <field>: invalid values (<count> tests)
    - Discriminated union: <field> (<count> tests)
    - ...
  onUpdate:
    - Immutable field <field>: rejected
    - Ratcheting: <field> (<count> tests)
    - Status validation: <field> (<count> tests)
    - ...

Total: <X> onCreate tests, <Y> onUpdate tests across <Z> files
```

---

## Behavioral Rules

1. **Prioritize custom validations over schema enforcement**: Focus test effort on CEL rules,
   cross-field constraints, immutability, discriminated unions, and defaulting. Do NOT generate
   tests that merely verify basic OpenAPI schema behavior (enum rejection, maxLength, pattern
   matching, min/max).

2. **Derive everything from source**: All test values, error messages, and expected behaviors
   MUST come from the actual Go types, validation markers, CRD OpenAPI schema, and CEL rules.
   Never hardcode assumed behavior.

3. **Match existing style exactly**: Read existing test files in the same `tests/` directory
   and match their formatting, indentation, naming patterns, and level of detail precisely.

4. **Expected blocks must include defaults**: When testing `onCreate` with `expected`, always
   include any default values the CRD schema applies. Read the CRD manifest OpenAPI schema
   to find `default:` declarations. Compare with existing tests to see what defaults are
   expected (e.g., `logLevel: Normal`, `httpEmptyRequestsPolicy: Respond`).

5. **Error messages must be accurate**: For `expectedError` and `expectedStatusError`, use the
   exact error message substrings from CEL `message` fields or Kubernetes validation output.

6. **Minimal YAML**: Only include fields relevant to each specific test case. Don't pad
   tests with unrelated fields.

7. **Preserve existing tests**: When adding to an existing file, keep all existing tests
   unchanged and only append new tests for new or modified fields.

8. **Feature gate consistency**: Each test file must correspond to exactly one CRD variant
   in the generated CRD manifests directory. The `featureGate`/`featureGates` field must
   match the variant. For compound gates (e.g., `FeatureA+FeatureB.yaml`), list all gates.

9. **Ratcheting tests**: Only generate ratcheting tests when there is evidence that validation
   rules have been tightened. Use `initialCRDPatches` to relax the specific validation that
   was tightened.
