---
document_type: technical_reference
primary_audience: ai_agents
purpose: Guide for OpenShift API field validation and test creation
last_updated: 2025-01-06
authoritative_source: /tests/README.md
key_sections:
  - quick_decision_trees
  - validation_hierarchy
  - common_errors
  - test_requirements
  - validator_reference
---

# OpenShift API Validation Testing Guide

## Executive Summary

This document guides AI agents in determining which validation mechanisms to use for OpenShift API fields and which tests to write for those validations.

**Key Decision Hierarchy**: CEL Format Validators > Format Markers > CEL Custom Logic > Patterns

**Core Principle**: All validation constraints must be documented in field godoc comments. Documentation without enforcement fails review; enforcement without documentation also fails review.

**Test Requirements**: Tests required for Pattern, XValidation, cross-field validation, immutability, and format markers. No tests needed for MinLength, MaxLength, MinItems, MaxItems, Enum, Required, Optional.

**Authoritative Reference**: `/tests/README.md` in the OpenShift API repository

---

## How to Update This Document

**Canonical Source for Section 7.2 (CEL Validators)**: `k8s.io/apiserver/pkg/cel/library` in the Kubernetes source tree

The CEL validator documentation in Section 7.2 is derived from the following library files:
- `format.go` - Format validation functions (DNS, UUID, date, byte, URI)
- `urls.go` - URL validation and parsing
- `ip.go` - IP address validation and properties
- `cidr.go` - CIDR validation and containment
- `quantity.go` - Kubernetes resource.Quantity operations
- `semverlib.go` - Semantic version validation and comparison
- `regex.go` - Pattern matching (find, findAll)

When updating validator documentation, refer to these source files for authoritative function signatures, behavior, and examples.

**Canonical Source for Section 7.1 (Format Markers)**: `k8s.io/kube-openapi/pkg/validation` in the [kube-openapi repository](https://github.com/kubernetes/kube-openapi)

The Format Marker documentation in Section 7.1 is derived from:
- `strfmt/default.go` - Canonical list of supported format markers and their validation implementations
  - Defines the default format registry with all standard validators (UUID, date-time, URI, email, etc.)
  - Each format's validation function and behavior

When updating format marker documentation, refer to `strfmt/default.go` for the authoritative list of supported formats and their validation behavior.

---

## Table of Contents

### Part 1: GUIDANCE
1. [Quick Decision Trees](#quick-decision-trees)
2. [Best Practices](#best-practices)
3. [Common Errors & Solutions](#common-errors--solutions)
4. [Acceptable Patterns (Do Not Flag)](#acceptable-patterns-do-not-flag)

### Part 2: REFERENCE
5. [Test Requirements](#5-test-requirements)
6. [Validator Quick Lookup](#6-validator-quick-lookup)
7. [Complete Validator Catalog](#7-complete-validator-catalog)
8. [Test Framework Reference](#8-test-framework-reference)
9. [Complete Examples](#9-complete-examples)
10. [Reference Examples](#10-reference-examples)
11. [CEL Format Validators vs Kubebuilder Format Markers](#11-cel-format-validators-vs-kubebuilder-format-markers)

---

# Part 1: GUIDANCE

## Quick Decision Trees

### 1.1 How to Choose Validation Method

**Security Note**: NEVER use Format markers `ipv4`, `ipv6`, or `cidr` (CVE-2024-24790, CVE-2021-29923). Always use CEL alternatives.

**Decision Hierarchy** (choose first applicable option):

1. **Use Format marker** if exact match exists and not prohibited
   - Example: UUID field → `Format=uuid`
   - Check: [Section 7.1.0](#710-security-prohibited-format-markers) for prohibited markers

2. **Use CEL format validator** if no Format marker exists OR Format is prohibited
   - Example: IPv4 field → `isIP(self) && ip(self).family() == 4` (Format=ipv4 is prohibited)
   - See: [Section 7.2](#72-cel-format-validators)

3. **Use CEL format + custom logic** for partial format matching or format with additional constraints
   - Example: Comma-separated IPs → CEL to split and validate each with `isIP()`
   - **Prefer multiple rules**: Format validation + constraint validation as separate rules
   - Example: DNS label starting with 'app-' → TWO XValidation rules (dns1123Label + startsWith)

4. **Use Pattern** for domain/protocol-specific formats with no standard validator
   - Example: AWS Subnet ID → `Pattern=^subnet-[0-9A-Za-z]{17}$`

**Key Principle**: Multiple validation rules > complex combined logic. See [Section 2.5](#25-multiple-validation-rules-for-format--constraint)

### 1.2 What Tests to Write

```
Does the field have any of these?
├─ Pattern validation → YES: See [5.1.4 Pattern Test Coverage](#514-pattern-test-coverage)
├─ XValidation (CEL) → YES: See [5.1.5 CEL Test Coverage](#515-cel-test-coverage)
├─ Cross-field dependencies → YES: Comprehensive tests
├─ Immutability constraints → YES: onUpdate tests
├─ Format markers → YES: See [5.1.3 Minimal Coverage](#513-minimal-coverage-acceptable-standard-libraries)
│
└─ Only these markers?
    ├─ MinLength, MaxLength, MinItems, MaxItems
    ├─ Enum, Required, Optional
    └─ NO tests needed (framework-enforced)
```

**Note**: "Comprehensive" means testing all logical paths, edge cases, and boundary conditions. See [Section 5.1 Test Coverage Requirements](#51-test-coverage-requirements) for complete coverage rules.

### 1.3 Documentation Requirements Checklist

For each field, verify the godoc comment includes:

```
☐ Human-readable description of the field's purpose
☐ For optional fields: Behavior when omitted ("When not specified, ...")
☐ For required fields: Statement "This field is required"
☐ For enum fields: All valid values and their meanings ("When set to X, ...")
  Exception: Single-value enums - see [4.3 Single-Value Enum Phrasings](#43--single-value-enum-phrasings)
☐ For validated fields: Complete constraint description matching all markers
☐ For cross-field relationships: Both description AND XValidation enforcement (see [3.8 XValidation Placement](#38-error-cross-field-validation-placement) for placement)
☐ For XValidation rules: Message accurately describes failure condition and uses correct enum values/field names (see [3.9 Inaccurate Validation Messages](#39-error-inaccurate-or-inconsistent-validation-messages))
```

**CRITICAL**: Documentation stating "cannot be used with field X" or "mutually exclusive with Y" MUST have corresponding `+kubebuilder:validation:XValidation` rules. Documentation without enforcement fails review.

### 1.4 Common Tasks Index

| Task | Go to Section |
|------|---------------|
| Choose between Format, CEL, or Pattern | [2.1 Validation Decision Hierarchy](#21-validation-decision-hierarchy) |
| Document validation constraints | [2.4 Documentation Requirements](#24-documentation-requirements) |
| Understand test coverage requirements | [5.1 Test Coverage Requirements](#51-test-coverage-requirements) |
| Fix array field validation | [3.1 Array Field Validation](#31-array-field-validation) |
| Understand format marker limitations | [3.2 Silent Format Stripping](#32-silent-format-stripping) |
| Fix cross-field validation placement | [3.8 XValidation Placement](#38-error-cross-field-validation-placement) |
| Fix validation error messages | [3.9 Inaccurate Validation Messages](#39-error-inaccurate-or-inconsistent-validation-messages) |
| Write onCreate tests | [8.3 onCreate Tests](#83-oncreate-tests) |
| Write onUpdate tests | [8.4 onUpdate Tests](#84-onupdate-tests) |
| Test validation ratcheting | [8.5 Validation Ratcheting](#85-validation-ratcheting-with-initialcrdpatches) |
| Look up a validator | [6. Validator Quick Lookup](#6-validator-quick-lookup) |

---

## 2. Best Practices

### 2.1 [CRITICAL] Validation Decision Hierarchy

When implementing field validation, follow this hierarchy to select the most appropriate validation mechanism:

**SECURITY WARNING**: Format markers `ipv4`, `ipv6`, and `cidr` MUST NOT be used (CVE-2024-24790, CVE-2021-29923). See [Section 7.1.0](#710-security-prohibited-format-markers) and use CEL alternatives instead.

#### 2.1.1 First Choice: Kubebuilder Format Markers (Security-Aware)

**When to use**: Field validation matches a supported Format marker exactly AND the marker is not prohibited

**Benefits**:
- Simpler syntax than CEL for basic cases
- Standard, well-tested validation logic
- Consistent behavior across Kubernetes ecosystem
- Clear intent through standardized format names

**Tradeoff**: Invalid format names are silently ignored by the API server (poor developer experience compared to CEL)

**Prohibited Format Markers** (use CEL instead):
- `ipv4`, `ipv6`, `cidr` - Security vulnerabilities CVE-2024-24790, CVE-2021-29923
- See: [Section 7.1.0](#710-security-prohibited-format-markers) for CEL alternatives

**Example**: UUID field
```go
// Good: Uses Format marker (allowed)
// uid is the unique identifier.
// Must be a valid UUID in 8-4-4-4-12 format (e.g., 550e8400-e29b-41d4-a716-446655440000).
// +kubebuilder:validation:Format=uuid
UID string `json:"uid"`

// Bad: IPv4 address with prohibited Format marker
// +kubebuilder:validation:Format=ipv4  // ❌ PROHIBITED - security vulnerability
IPAddress string `json:"ipAddress"`

// Good: IPv4 address with CEL
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).family() == 4"
IPAddress string `json:"ipAddress"`
```

See: [7.1 Kubernetes Format Markers](#71-kubernetes-format-markers) for complete list

#### 2.1.2 Second Choice: CEL Format Validators

**When to use**:
- No Format marker exists for your validation need, OR
- Format marker exists but is prohibited for security reasons

**Benefits**:
- Standard, well-tested validation logic from Kubernetes
- Invalid CEL fails during CRD admission (immediate feedback)
- Access to CEL's expressive capabilities for combining validators
- Clear intent through standardized function names
- Required when Format marker is prohibited (ipv4, ipv6, cidr)

**Example**: IPv4 field (Format marker prohibited)
```go
// Good: Uses CEL format validator (required for IPv4)
// ipAddress is the IPv4 address.
// Must be a valid IPv4 address in dotted-quad notation (e.g., 192.168.1.1).
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).family() == 4",message="must be a valid IPv4 address"
IPAddress string `json:"ipAddress"`

// Bad: Prohibited Format marker
// +kubebuilder:validation:Format=ipv4  // ❌ SECURITY VULNERABILITY
IPAddress string `json:"ipAddress"`

// Also Good: UUID with CEL (alternative to Format marker)
// uid is the unique identifier.
// Must be a valid UUID in 8-4-4-4-12 format (e.g., 550e8400-e29b-41d4-a716-446655440000).
// +kubebuilder:validation:XValidation:rule="!format.uuid().validate(self).hasValue()",message="must be a valid UUID"
UID string `json:"uid"`
```

**Important**: CEL format validators return an optional error. Use the pattern `!format.validator().validate(self).hasValue()` to check validity (the `!` negates the error presence).

See: [7.2 CEL Format Validators](#72-cel-format-validators) for complete list

#### 2.1.3 Third Choice: CEL with Custom Logic

**When to use**: Field requires validation that includes but extends beyond a standard format

**Benefits**:
- Combines standard validation with custom logic
- Invalid CEL fails during CRD admission (immediate feedback)
- Access to CEL's full expressive capabilities

**IMPORTANT**: When combining format validation with additional constraints, prefer **multiple separate validation rules** over a single combined rule. See [Section 2.5](#25-multiple-validation-rules-for-format--constraint) for details.

**Example**: DNS label with custom constraint (PREFERRED - separate rules)
```go
// PREFERRED: Separate validation rules for better clarity and error messages
// name is the application name.
// Must be a valid DNS-1123 label starting with 'app-' (e.g., app-myservice).
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue()",message="must be a valid DNS-1123 label"
// +kubebuilder:validation:XValidation:rule="self.startsWith('app-')",message="must start with 'app-'"
Name string `json:"name"`

// ACCEPTABLE: Combined rule (less preferred)
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue() && self.startsWith('app-')",message="must be a valid DNS-1123 label starting with 'app-'"
Name string `json:"name"`

// Bad: Custom pattern duplicates DNS label validation logic
// +kubebuilder:validation:Pattern=`^app-[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
Name string `json:"name"`
```

**Use Case**: Partial format matching (no Format marker or CEL validator covers complete requirement)
```go
// Example: Comma-separated list of IP addresses
// ips is a comma-separated list of IPv4 addresses.
// +kubebuilder:validation:XValidation:rule="self.split(',').all(ip, isIP(ip.trim()) && ip(ip.trim()).family() == 4)",message="must be comma-separated IPv4 addresses"
IPs string `json:"ips"`
```

#### 2.1.4 Fourth Choice: Custom Patterns

**When to use**: No standard validator covers your use case

**Common scenarios**:
- Domain-specific identifiers (e.g., AWS resource IDs)
- Protocol-specific formats (e.g., HTTP header names per RFC 2616)
- Application-specific syntax requirements

**Example**: AWS Subnet ID
```go
// Appropriate: No standard validator for AWS-specific format
// subnetID is the AWS subnet identifier.
// Must follow AWS subnet ID format: 'subnet-' followed by exactly 17 alphanumeric characters.
// Total length must be exactly 24 characters (e.g., subnet-0123456789abcdef0).
// +kubebuilder:validation:Pattern=`^subnet-[0-9A-Za-z]+$`
// +kubebuilder:validation:MinLength=24
// +kubebuilder:validation:MaxLength=24
SubnetID string `json:"subnetID"`
```

### 2.2 Combining Validators

Multiple validation mechanisms can and should be combined when appropriate:
- Format markers with length constraints
- CEL format functions with cross-field validation
- Pattern validators with numeric range constraints

**Key principle**: Each validator should enforce a distinct aspect of the validation requirement.

### 2.3 Type-Level vs Field-Level Validation Markers

**Preferred Pattern**: Define validation markers on types, not on individual fields

```go
// Type definition with validation
// +kubebuilder:validation:Enum=Deny;Warn
type AdmitAction string

const (
    AdmitActionDeny AdmitAction = "Deny"
    AdmitActionWarn AdmitAction = "Warn"
)

// CORRECT: Field documentation describes constraint, marker is on type
// action determines the admission behavior.
// Valid options are Deny and Warn.
// When set to Deny, requests will be rejected.
// When set to Warn, requests will be admitted with a warning.
// +required
Action AdmitAction `json:"action"`

// INCORRECT: Do not duplicate the Enum marker from the type
// +kubebuilder:validation:Enum=Deny;Warn  // ❌ REDUNDANT
Action AdmitAction `json:"action"`
```

**Benefits**:
- Kubebuilder automatically applies type-level markers to all fields of that type
- Single source of truth (update once, applies everywhere)
- Ensures consistency across all uses of the type
- Reduces maintenance burden

**Key Principle**: Document constraints in field comments, enforce them once on the type definition.

### 2.4 [CRITICAL] Documentation Requirements

Regardless of which validation mechanism you choose, the field's godoc comment must completely describe all validation constraints in human-readable language.

#### Required Elements:

1. **Field purpose and behavior**
2. **Valid format with examples** where helpful
3. **Reference to standards** (RFCs, etc.) when applicable
4. **For optional fields**: Behavior when omitted
5. **For required fields**: Explicit statement "This field is required"
6. **All validation constraints** in natural language
7. **For XValidation rules**: Error messages that accurately describe when validation fails, using exact enum values and JSON field names

#### Documentation Style Preferences:

| Scenario | Preferred Phrasing | Avoid |
|----------|-------------------|-------|
| Required fields | "This field is required" | Omitting this statement |
| Optional fields | "When/If omitted/not specified/present, &lt;behavior&gt;" | "This field is optional" (when behavior is specified) |
| Single-value enums | "The only permitted value is X" / "Must be X" | "Valid values are: X" |
| Multi-value enums | "When set to X, &lt;behavior&gt;" format | Bulleted lists (poor kubectl describe rendering) |
| Standard K8s fields | Minimal documentation | Detailed validation docs |

**Note on Phrasing Equivalence**: Semantically equivalent phrasings are acceptable (e.g., "When not specified" vs "If not specified", "When omitted" vs "If omitted"). Do not flag minor stylistic variations.

**Special Case**: `status.conditions` (standard Kubernetes field) does NOT require detailed validation documentation.

#### Complete Example:

```go
// duration specifies the maximum idle time for connections.
// The value must be a valid duration string parseable by Go's time.ParseDuration
// (e.g., "30s", "5m", "1h30m"). Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// When omitted, connections will not timeout due to idle time.
// +optional
// +kubebuilder:validation:Format=duration
Duration *metav1.Duration `json:"duration,omitempty"`
```

See: [3.5 Missing Documentation](#35-missing-documentation-for-validation-constraints)

### 2.5 Multiple Validation Rules for Format + Constraint

When a field requires both format validation AND an additional constraint, prefer **multiple separate validation rules** over combining them into complex CEL logic.

**Benefits**:
- **Clearer intent**: Each rule validates one specific aspect
- **Better error messages**: Users see which specific constraint failed
- **Easier maintenance**: Modify individual constraints independently
- **Simpler testing**: Test each validation rule separately

**Pattern**: Use separate `+kubebuilder:validation:XValidation` rules for:
1. Format validation (e.g., `format.dns1123Label()`)
2. Additional constraints (e.g., `startsWith()`, length checks, value restrictions)

**Example 1**: DNS label with prefix requirement
```go
// PREFERRED: Two separate validation rules
// name is the application name.
// Must be a valid DNS-1123 label starting with 'app-' (e.g., app-myservice).
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue()",message="must be a valid DNS-1123 label"
// +kubebuilder:validation:XValidation:rule="self.startsWith('app-')",message="must start with 'app-'"
// +kubebuilder:validation:MaxLength=63
Name string `json:"name"`

// AVOID: Single combined rule (less clear, less helpful error messages)
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue() && self.startsWith('app-')",message="must be a valid DNS-1123 label starting with 'app-'"
Name string `json:"name"`
```

**Example 2**: URI with HTTPS scheme requirement
```go
// PREFERRED: Separate rules for URI format and scheme constraint
// endpoint is the HTTPS endpoint URL.
// Must be a valid URI with HTTPS scheme (e.g., https://api.example.com/v1).
// +kubebuilder:validation:Format=uri
// +kubebuilder:validation:XValidation:rule="url(self).getScheme() == 'https'",message="scheme must be https"
Endpoint string `json:"endpoint"`
```

**Example 3**: Email address with domain restriction
```go
// PREFERRED: Separate rules for email format and domain constraint
// email is the contact email address.
// Must be a valid email address ending with @redhat.com (e.g., user@redhat.com).
// +kubebuilder:validation:Format=email
// +kubebuilder:validation:XValidation:rule="self.endsWith('@redhat.com')",message="must end with @redhat.com"
Email string `json:"email"`
```

**When to combine**: Only combine rules when a single expensive operation like `url(self)` provides multiple data points that must be checked.

```go
// ACCEPTABLE: Truly atomic validation
// The url(self) operation is performed once, and its results (port, scheme) are checked.
// It would be inefficient to have two separate rules that both call url(self).
// +kubebuilder:validation:XValidation:rule="url(self).getScheme() == 'https' && int(url(self).getPort()) >= 1024",message="URL must be HTTPS and use port >= 1024"
Endpoint string `json:"endpoint"`
```

**Key Principle**: If the validations are separate operations (e.g., format.dns1123Label() and self.startsWith()), they MUST be in separate rules. If they are attributes of one operation (e.g., url().getScheme() and url().getPort()), they MAY be in one rule.

---

## 3. Common Errors & Solutions

**Note**: This section documents actual problems that require fixes. See [Section 4: Acceptable Patterns](#4-acceptable-patterns-do-not-flag) for patterns that should NOT be flagged as issues.

### 3.1 [ERROR] Array Field Validation

**Problem**: Applying field-level validation markers directly to array types has no effect.

```go
// WRONG: MinLength applies to the array itself, not the items
// names is a list of service names.
// Each name must be between 5 and 63 characters.
// +kubebuilder:validation:MinLength=5  // ❌ IGNORED
// +kubebuilder:validation:MaxLength=63 // ❌ IGNORED
Names []string `json:"names"`
```

**Solution**: Use `+kubebuilder:validation:items:` prefix to apply validation to array elements.

```go
// CORRECT: Validation applies to each string in the array
// names is a list of service names.
// Each name must be between 5 and 63 characters.
// +kubebuilder:validation:items:MinLength=5
// +kubebuilder:validation:items:MaxLength=63
Names []string `json:"names"`
```

**Note**: Array-level validations like `MinItems`, `MaxItems`, and `UniqueItems` apply directly to the array field without the `items:` prefix.

### 3.2 [ERROR] Silent Format Stripping

**Problem**: Unsupported or misspelled format markers are silently removed from the generated CRD at runtime without any build-time warning.

```go
// WRONG: Typo in format name - will be silently ignored
// +kubebuilder:validation:Format=ipv-4  // ❌ Should be "ipv4"
IPAddress string `json:"ipAddress"`

// WRONG: Unsupported format - will be silently stripped
// +kubebuilder:validation:Format=custom-format  // ❌ Not a valid format
CustomField string `json:"customField"`
```

**Solution**:
- Only use documented format values (see [7.1 Kubernetes Format Markers](#71-kubernetes-format-markers))
- Test generated CRDs to verify format appears in the schema
- Consider CEL format functions instead - invalid CEL fails during CRD admission (unlike Format markers)

**Note**: Format markers only apply to string-typed fields. They are silently ignored when applied to other types (integers, booleans, arrays, objects).

### 3.3 [ERROR] CEL Validation Pattern Confusion

**Problem**: Forgetting the double-negative pattern when using CEL format validators.

```go
// WRONG: Returns true when validation FAILS
// +kubebuilder:validation:XValidation:rule="format.dns1123Label().validate(self).hasValue()"
// ❌ hasValue() returns true when there IS an error

// CORRECT: Returns true when validation SUCCEEDS
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue()"
// ✓ The '!' negates, so true when there is NO error
```

**Explanation**: CEL format validators return an optional error. The `.hasValue()` method returns true when there IS an error (validation failed). The leading `!` negates this to make the rule pass when validation succeeds.

### 3.4 [ERROR] Pattern vs CEL .matches() Escaping Requirements

**Problem**: Regex patterns require different escaping in `Pattern` markers versus CEL `.matches()` functions.

**RECOMMENDED APPROACH**: Use the simplest escaping pattern:

1. **For Pattern markers**: Use backticks with single backslash
   ```go
   // CORRECT
   // +kubebuilder:validation:Pattern=`^[a-z]\d+$`
   ```

2. **For CEL matches**: Use backticks for rule + raw string (`r''`) for pattern
   ```go
   // CORRECT
   // +kubebuilder:validation:XValidation:rule=`self.matches(r'^[a-z]\d+$')`,message="must match pattern"
   ```

**Key Rule**: Both Pattern markers and CEL patterns use **single backslash** (`\d`, `\s`, `\w`, etc.) when following this approach.

**Common Mistakes**:
```go
// NOT RECOMMENDED: Using quoted strings instead of backticks + raw strings
// +kubebuilder:validation:XValidation:rule="self.matches('^[a-z]\\\\d+$')"
// ❌ Requires 4 backslashes, error-prone

// WRONG: Forgetting r prefix with backticks
// +kubebuilder:validation:XValidation:rule=`self.matches('^[a-z]\d+$')`
// ❌ CEL compilation fails: "Syntax error: token recognition error at: '\d'"
```

**Note**: While other quoting combinations are technically valid, they require complex multi-level escaping (2 or 4 backslashes) and are error-prone. Prefer the backtick + raw string approach above.

### 3.5 [ERROR] Missing Documentation for Validation Constraints

**Problem**: Validation markers without corresponding documentation in the godoc comment.

```go
// WRONG: Validation not documented
// uid is the unique identifier for this resource.
// +kubebuilder:validation:Format=uuid       // ❌ Format not described
// +kubebuilder:validation:MinLength=36      // ❌ Length not described
// +kubebuilder:validation:MaxLength=36
UID string `json:"uid"`
```

**Solution**: Always document all validation constraints in the field comment.

```go
// CORRECT: Complete documentation
// uid is the unique identifier for this resource.
// Must be a valid UUID in 8-4-4-4-12 format (e.g., 550e8400-e29b-41d4-a716-446655440000).
// Length must be exactly 36 characters.
// +kubebuilder:validation:Format=uuid
// +kubebuilder:validation:MinLength=36
// +kubebuilder:validation:MaxLength=36
UID string `json:"uid"`
```

### 3.6 [ERROR] Cross-Field Constraints Without Enforcement

**Problem**: Documentation describes field relationships but lacks corresponding validation rules.

```go
// WRONG: Documented constraint without enforcement
// mode determines the operating mode.
// Cannot be used together with legacyMode field.  // ❌ Not enforced
// +optional
Mode *string `json:"mode,omitempty"`

// legacyMode determines legacy operating mode.
// Cannot be used together with mode field.  // ❌ Not enforced
// +optional
LegacyMode *string `json:"legacyMode,omitempty"`
```

**Solution**: Add XValidation rule to enforce the documented constraint. The validation must be placed on the parent struct that contains both fields.

```go
// CORRECT: Constraint both documented and enforced
// +kubebuilder:validation:XValidation:rule="!(has(self.mode) && has(self.legacyMode))",message="mode and legacyMode are mutually exclusive"
type MySpec struct {
    // mode determines the operating mode.
    // Cannot be used together with legacyMode field.
    // +optional
    Mode *string `json:"mode,omitempty"`

    // legacyMode determines legacy operating mode.
    // Cannot be used together with mode field.
    // +optional
    LegacyMode *string `json:"legacyMode,omitempty"`
}
```

**CRITICAL**: Any documented field relationship (mutually exclusive, dependent, conditional) MUST have a corresponding validation rule.

**Note**: Cross-field validation must be placed on the parent struct, not individual fields. See [3.8 XValidation Placement](#38-error-cross-field-validation-placement) for details on correct placement.

**Key Principle**: XValidation rules are placed on the parent struct, but documentation of the constraint goes on the affected field(s). Each field that participates in the constraint should document its relationship to other fields.

**Documentation Placement Guidelines**:
- **Mutual exclusion**: Document on both fields ("Cannot be used with X")
- **Conditional requirement**: Document on the dependent field ("Required when X is set to Y")
- **Conditional prohibition**: Document on the prohibited field ("Forbidden when X is set to Y")

### 3.7 [ERROR] Case Sensitivity in Format Names

**Problem**: Format marker names are case-sensitive and must match exactly (lowercase).

```go
// WRONG: Case mismatch
// +kubebuilder:validation:Format=UUID  // ❌ Uppercase
UID string `json:"uid"`

// CORRECT: Exact lowercase name as documented
// +kubebuilder:validation:Format=uuid  // ✓
UID string `json:"uid"`

// CORRECT: Use documented format name
// +kubebuilder:validation:Format=date-time  // ✓
Timestamp string `json:"timestamp"`
```

**Note**: Format marker names are normalized internally by removing dashes (via `DefaultNameNormalizer` in `k8s.io/kube-openapi/pkg/validation/strfmt/format.go`), so `date-time` and `datetime` are equivalent. However, always use the primary documented name with dashes (e.g., `date-time`) from [Section 7.1](#71-kubernetes-format-markers) for consistency and readability.

### 3.8 [ERROR] Cross-Field Validation Placement

**Problem**: Placing cross-field XValidation rules on individual fields instead of the parent struct.

```go
// WRONG: XValidation on individual field
// mode determines the operating mode.
// Cannot be used together with legacyMode field.
// +kubebuilder:validation:XValidation:rule="!(has(self.mode) && has(self.legacyMode))",message="mode and legacyMode are mutually exclusive"
// ❌ The 'mode' field doesn't have 'mode' or 'legacyMode' fields
// +optional
Mode *string `json:"mode,omitempty"`
```

**Solution**: Place XValidation on the parent struct that contains both fields.

```go
// CORRECT: XValidation on parent struct
// +kubebuilder:validation:XValidation:rule="!(has(self.mode) && has(self.legacyMode))",message="mode and legacyMode are mutually exclusive"
type MySpec struct {
    // mode determines the operating mode.
    // Cannot be used together with legacyMode field.
    // +optional
    Mode *string `json:"mode,omitempty"`

    // legacyMode determines legacy operating mode.
    // Cannot be used together with mode field.
    // +optional
    LegacyMode *string `json:"legacyMode,omitempty"`
}
```

**Explanation**: The `has()` function checks whether a field exists on the current object (`self`). When validating cross-field constraints, `self` must refer to the struct containing all relevant fields, not an individual field.

**Documentation vs Validation Placement**:
- **Validation (XValidation)**: MUST be on the parent struct that contains all fields involved in the constraint
- **Documentation**: MUST be on the affected field(s) themselves, describing the constraint in human-readable terms

This separation allows field documentation to be seen when viewing individual fields (e.g., in kubectl explain) while keeping the technical validation logic at the struct level where it can access all necessary fields.

### 3.9 [ERROR] Inaccurate or Inconsistent Validation Messages

**Problem**: XValidation message attributes that don't accurately describe the validation rule or use inconsistent terminology (e.g., wrong enum values, ambiguous field references).

**Common Issues**:
1. Messages referencing wrong enum values
2. Messages describing inverse of actual logic
3. Messages using CEL paths instead of JSON field names
4. Messages inconsistent with field/type names

#### Issue 1: Wrong Enum Value in Message

```go
// WRONG: Message references 'All' but enum value is 'AllServed'
// +kubebuilder:validation:Enum=StorageOnly;AllServed
type SelectionType string

// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additional)",message="additional may not be defined when defaultSelection is 'All'"
// ❌ Message says 'All' but enum value is 'AllServed'
type APIVersions struct { ... }
```

**Solution**: Use exact enum values in messages.

```go
// CORRECT: Message uses exact enum value
// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additional)",message="additionalVersions may not be specified when defaultSelection is 'AllServed'"
type APIVersions struct { ... }
```

#### Issue 2: Message Describes Inverse Logic

```go
// WRONG: Rule checks for presence, message describes absence
// +kubebuilder:validation:XValidation:rule="has(self.mode)",message="mode must not be specified"
// ❌ Rule requires field, message says it must not be specified
```

**Solution**: Ensure message accurately reflects when the rule fails.

```go
// CORRECT: Message matches rule logic
// +kubebuilder:validation:XValidation:rule="has(self.mode)",message="mode is required"
```

#### Issue 3: Incorrect Field References

```go
// WRONG: Message uses Go field instead of JSON field name
// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additionalVersions)",message="AdditionalVersions may not be defined when defaultSelection is 'AllServed'"
// ❌ Uses 'additional' (CEL path) instead of 'additionalVersions' (field name users see)
```

**Solution**: Use the actual JSON field names that users see.

```go
// CORRECT: Message uses JSON field name
// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additionalVersions)",message="additionalVersions may not be defined when defaultSelection is 'AllServed'"
AdditionalVersions []string `json:"additionalVersions,omitempty"`
```

**Validation Checklist for Messages**:
- [ ] Uses exact enum values (not abbreviations or variations)
- [ ] Accurately describes when the rule fails (not when it succeeds)
- [ ] References JSON field names (what users see)
- [ ] Uses consistent terminology with field documentation
- [ ] Clear and actionable (user knows what to fix)

**Review Guidance**: When reviewing XValidation rules, verify that the message attribute accurately reflects the validation logic and uses terminology consistent with the API definition.

---

## 4. Acceptable Patterns (Do Not Flag)

**Purpose**: This section documents patterns that are **acceptable** and should **NOT** be flagged as problems during API reviews.

### 4.1 ✅ Explicit Length Constraints with Format Markers

**Pattern**: Adding MinLength/MaxLength markers alongside Format markers that already enforce length.

```go
// ACCEPTABLE: MinLength/MaxLength with Format=uuid
// uid is the unique identifier.
// Must be a valid UUID in 8-4-4-4-12 format.
// +kubebuilder:validation:Format=uuid
// +kubebuilder:validation:MinLength=36
// +kubebuilder:validation:MaxLength=36
UID string `json:"uid"`
```

**Rationale**: While the Format marker already enforces the length constraint, the explicit MinLength/MaxLength markers serve important purposes:
- Contribute to API cost complexity analysis for resource limits
- Make requirements explicit in documentation
- Provide clear bounds for API server resource planning

### 4.2 ✅ "This field is required" Documentation

**Pattern**: Explicitly documenting that a field is required.

```go
// ACCEPTABLE: Documenting required status aids readability
// name is the resource name.
// This field is required.
// +required
Name string `json:"name"`
```

**Rationale**: The `+required` marker enforces the constraint at the API level but is not included in generated documentation. The explicit "This field is required" statement in the godoc comment informs users of the requirement and is recommended practice.

### 4.3 ✅ Single-Value Enum Phrasings

**Pattern**: Using natural language for single-value enums.

```go
// ACCEPTABLE: Natural phrasing for single-value enum
// type indicates the data format.
// The only supported type is "yaml".
// +required
Type DataType `json:"type"`

// ALSO ACCEPTABLE: Alternative phrasings
// The value must be "yaml".
// Must be "yaml".
```

**Rationale**: For single-value enums, phrases like "the only permitted value is", "the value must be", or "must be" are more natural than "Valid values are: yaml".

### 4.4 ✅ "This field is optional" (Context-Dependent)

**When ACCEPTABLE**:
1. When NOT accompanied by "When/If omitted/present/not specified" clauses
2. When present alongside omission behavior (style preference only - not worth flagging in isolation)

```go
// ACCEPTABLE: Optional status useful when omission behavior isn't specified
// timeout specifies the request timeout.
// This field is optional.
// +optional
Timeout *int32 `json:"timeout,omitempty"`
```

**When to recommend removal**: Only when you are already recommending other substantive changes to the same field's documentation.

```go
// LOW PRIORITY: "This field is optional" is redundant here, but not critical
// timeout specifies the request timeout.
// When not specified, requests never timeout.
// This field is optional.  // <-- REDUNDANT but low priority
// +optional
Timeout *int32 `json:"timeout,omitempty"`

// PREFERRED: More concise without redundant statement
// timeout specifies the request timeout.
// When not specified, requests never timeout.
// +optional
Timeout *int32 `json:"timeout,omitempty"`
```

**Review Priority**: This is a **low-priority style issue**. Only flag "This field is optional" for removal if you are already recommending other substantive changes to the same field's documentation. Do not create review comments solely to remove "This field is optional" statements.

**Rationale**: "When/If omitted/not specified" clauses already convey the optional nature. Adding "This field is optional" provides no additional information, but removing it is not important enough to warrant a standalone review comment.

### 4.5 ✅ Fields Already Following Best Practices

**Pattern**: Well-documented fields with complete validation documentation and proper omission behavior.

```go
// ACCEPTABLE: Exemplary field documentation - DO NOT FLAG
// matchConditions defines filtering rules for the webhook.
// When present, must contain between 1 and 64 match conditions.
// When not specified, the webhook will match all requests.
// +kubebuilder:validation:MinItems=1
// +kubebuilder:validation:MaxItems=64
// +optional
MatchConditions []MatchCondition `json:"matchConditions,omitempty"`
```

**Rationale**: This field is already exemplary. Flagging it with "No change needed" or "Already compliant" wastes review time and creates noise.

### 4.6 ✅ Standard Kubernetes Fields

**Pattern**: Common Kubernetes fields like `status.conditions`, `metadata`, `spec`, etc.

```go
// ACCEPTABLE: Standard field - DO NOT FLAG
// conditions is a list of status conditions.
// +optional
Conditions []metav1.Condition `json:"conditions,omitempty"`
```

**Rationale**: Standard Kubernetes fields like `status.conditions` do not require detailed validation documentation.

### 4.7 ✅ Validation Markers on Field Types

**Pattern**: Validation markers defined on a type are automatically applied to fields of that type.

**IMPORTANT**: This is the correct and preferred pattern. Duplicating type-level validation markers on fields is redundant and should be avoided.

```go
// Type definition with validation
// +kubebuilder:validation:Enum=Deny;Warn
type AdmitAction string

const (
    AdmitActionDeny AdmitAction = "Deny"
    AdmitActionWarn AdmitAction = "Warn"
)

// CORRECT: Field documentation describes constraint, marker is on type
// action determines the admission behavior.
// Valid options are Deny and Warn.
// When set to Deny, requests will be rejected.
// When set to Warn, requests will be admitted with a warning.
// +required
Action AdmitAction `json:"action"`
```

**Array Fields with Typed Elements**:

```go
// Type definition with validation
// +kubebuilder:validation:Format=uuid
type ResourceID string

type MyResource struct {
    // ids is a list of resource identifiers.
    // Each ID must be a valid UUID in 8-4-4-4-12 format.
    // The list may contain between 1 and 100 IDs.
    // +kubebuilder:validation:MinItems=1
    // +kubebuilder:validation:MaxItems=100
    // +listType=set
    IDs []ResourceID `json:"ids"`
    // Note: Format=uuid marker is on ResourceID type, not duplicated here
}
```

**Rationale**:
- Kubebuilder and controller-gen automatically apply type-level validation markers to all fields of that type
- Duplicating markers creates maintenance burden (must update in multiple places)
- Type-level markers ensure consistency across all uses of the type
- The field's documentation should still describe the constraint in human-readable form

### 4.8 ✅ Equivalent Phrasing Variations

**Pattern**: Using semantically equivalent phrases for the same concept.

**ACCEPTABLE - All of these are equivalent**:
```go
// When not specified, defaults to X
// When omitted, defaults to X
// If not specified, defaults to X
// If omitted, defaults to X
```

**Rationale**: Natural language has many acceptable ways to express the same concept. AI agents should not flag minor stylistic variations that don't affect clarity or correctness.

**Examples of acceptable equivalences**:
- "When/If not specified/omitted/present"
- "Must be/Should be/Must consist of"
- "Contains/Includes/Comprises"
- "Cannot/May not/Must not"

**Review Guidance**: Do not create review comments for equivalent phrasings. Focus on substantive documentation issues like missing validation constraints, incorrect behavior descriptions, or unclear field purposes.

### 4.9 ✅ Cross-Field Validation Documentation on Fields

**Pattern**: When a struct has XValidation enforcing cross-field constraints, the documentation explaining those constraints appears on the affected field(s), not on the struct itself.

**CORRECT - Documentation on affected field**:
```go
// APIVersions specifies a set of API versions of a CRD.
// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additional)",message="additional may not be defined when defaultSelection is 'All'"
type APIVersions struct {
    // defaultSelection specifies selection method.
    // +required
    DefaultSelection APIVersionSelectionType `json:"defaultSelection"`

    // additionalVersions specifies additional versions to require.
    // Cannot be specified when defaultSelection is set to 'AllServed'.
    // +optional
    AdditionalVersions []APIVersionString `json:"additionalVersions,omitempty"`
}
```

**WRONG - Documentation on struct instead of field**:
```go
// APIVersions specifies a set of API versions of a CRD.
// The additionalVersions field may not be defined when defaultSelection is set to 'AllServed'.  // ❌ WRONG LOCATION
// +kubebuilder:validation:XValidation:rule="self.defaultSelection != 'AllServed' || !has(self.additional)",message="additional may not be defined when defaultSelection is 'All'"
type APIVersions struct {
    // defaultSelection specifies selection method.
    // +required
    DefaultSelection APIVersionSelectionType `json:"defaultSelection"`

    // additionalVersions specifies additional versions to require.
    // ❌ MISSING: Constraint documentation should be here
    // +optional
    AdditionalVersions []APIVersionString `json:"additionalVersions,omitempty"`
}
```

**Rationale**:
- Field-level documentation appears in `kubectl explain` output for that specific field
- Struct-level comments are for describing the overall purpose of the struct
- Users reading field documentation need to see constraints on that field
- XValidation is a technical implementation detail that belongs at the struct level where `self` has the right scope

**Review Guidance**: When reviewing cross-field validation, check that constraint documentation appears on the affected field(s), not the parent struct. The struct comment should only have XValidation markers and describe the struct's overall purpose. Do not recommend adding constraint documentation to struct comments.

---

# Part 2: REFERENCE

## 5. Test Requirements

### 5.0 Validation-to-Test Quick Lookup

**Purpose**: Quickly determine whether tests are required for a validation marker.

| If you see this validation... | Tests Required? | Coverage Level | Test Section | Details |
|-------------------------------|-----------------|----------------|--------------|---------|
| `MinLength`, `MaxLength` | ❌ No | None | N/A | Framework-enforced ([5.1.1](#511-no-tests-required-framework-enforced)) |
| `MinItems`, `MaxItems` | ❌ No | None | N/A | Framework-enforced ([5.1.1](#511-no-tests-required-framework-enforced)) |
| `Enum`, `Required`, `Optional` | ❌ No | None | N/A | Framework-enforced ([5.1.1](#511-no-tests-required-framework-enforced)) |
| `Format=<marker>` (non-prohibited) | ✅ Yes | **Minimal** (1 pass, 1 fail) | onCreate | Examples: `uuid`, `uri`, `email` ([5.1.3](#513-minimal-coverage-acceptable-standard-libraries)) |
| `Pattern=...` | ✅ Yes | **Comprehensive** | onCreate | Test all pattern requirements ([5.1.4](#514-pattern-test-coverage)) |
| `XValidation` (standard lib only) | ✅ Yes | **Minimal** (1 pass, 1 fail) | onCreate | Examples: `isIP()`, `isCIDR()`, `format.dns1123Label()` ([5.1.3](#513-minimal-coverage-acceptable-standard-libraries)) |
| `XValidation` (custom logic) | ✅ Yes | **Comprehensive** | onCreate | Examples: `self.startsWith()`, conditionals ([5.1.5](#515-cel-test-coverage)) |
| `XValidation` (cross-field) | ✅ Yes | **Comprehensive** | onCreate | Examples: `has(self.a) && has(self.b)` ([5.1.2](#512-comprehensive-coverage-required)) |
| `XValidation` (immutability) | ✅ Yes | **Comprehensive** | onUpdate | Example: `self == oldSelf` ([5.1.2](#512-comprehensive-coverage-required)) |

**Quick Answers**:
- **No tests needed**: Length/count constraints, enums, required/optional markers
- **Minimal tests** (1 valid + 1 invalid): Format markers and CEL standard library validators
- **Comprehensive tests** (all paths/combinations): Patterns, custom CEL logic, cross-field validation, immutability

### 5.0.1 Test Coverage Analysis Methodology

**Purpose**: Reusable process for analyzing test coverage completeness.

**Step-by-step process:**

1. **Extract validation markers** from API type files (`types_*.go`):
   - Search for all `+kubebuilder:validation:*` markers
   - Group by field/struct location

2. **Categorize using Section 5.0** (Validation-to-Test Quick Lookup):
   - For each validation, look up the table to determine: Tests Required? Coverage Level?
   - Separate into three buckets:
     - No tests required (5.1.1)
     - Minimal coverage required (5.1.3)
     - Comprehensive coverage required (5.1.2)

3. **Locate test files**:
   - Path: `<group>/<version>/tests/<crd-name>/*.yaml`
   - Read all test suite YAML files
   - Parse `onCreate` and `onUpdate` sections

4. **Map validations to tests**:
   - For each validation in "Minimal" or "Comprehensive" buckets:
     - Search test cases for matching field name and validation scenario
     - Match `expectedError` messages to validation error messages
     - Count valid (expected) and invalid (expectedError) test cases

5. **Verify coverage completeness**:
   - **Minimal (5.1.3)**: Requires exactly 1 valid + 1 invalid test
   - **Comprehensive (5.1.2, 5.1.4, 5.1.5)**: Requires tests for all logical paths/branches/combinations
     - Pattern: test each component of the regex (see 5.1.4)
     - Cross-field: test all field combinations (see 5.1.2)
     - Custom CEL: test all branches/conditionals (see 5.1.5)

6. **Identify gaps**:
   - **Missing**: Validation requires tests but has zero coverage
   - **Insufficient**: Validation has some tests but doesn't meet requirements
     - Minimal: missing valid OR invalid case
     - Comprehensive: missing logical paths, edge cases, or combinations
   - **Unnecessary**: Tests exist for 5.1.1 framework-enforced validations

7. **Report using standard format** (see below)

### 5.0.2 Standard Test Coverage Report Format

**Missing/Insufficient Test Coverage:**
```
* [Test description] (onCreate) - [accepted/rejected]
  - Example: [value] ([reason per validation rule])

* [Test description] (onUpdate) - [accepted/rejected]
  - Example: [value] → [new value] ([reason per validation rule])
```

**Example (onCreate tests):**
```
* Valid UUID format (onCreate) - accepted
  - Example: 550e8400-e29b-41d4-a716-446655440000 (Format=uuid accepts standard UUID)
* Invalid UUID format (onCreate) - rejected
  - Example: not-a-uuid (Format=uuid requires 8-4-4-4-12 format)
```

**Example (onUpdate tests):**
```
* Changing immutable field (onUpdate) - rejected
  - Example: clusterID: "cluster-123" → clusterID: "cluster-456" (clusterID is immutable)
* Updating other fields without changing immutable field (onUpdate) - accepted
  - Example: clusterID: "cluster-123", otherField: "foo" → clusterID: "cluster-123", otherField: "bar" (immutable field unchanged)
```

**Unnecessary Test Report Format:**
```
* Test: "[test name]" (line X in test file)
  - Reason: Tests [validation marker] which is framework-enforced (5.1.1)
  - Action: Can be removed to reduce test maintenance burden
```

### 5.1 Test Coverage Requirements

This section is the **authoritative source** for all test coverage rules. An API requires tests for all non-trivial validations.

#### 5.1.1 No Tests Required (Framework-Enforced)

The following validations are enforced by the framework and do not require explicit tests:
- `MaxLength`, `MinLength`
- `MaxItems`, `MinItems`
- `MinProperties`, `MaxProperties`
- `Required`, `Optional`
- `Enum`

#### 5.1.2 Comprehensive Coverage Required

The following validations require tests covering **all logical paths, edge cases, and boundary conditions**:

**Pattern Validation:**
- Test all requirements expressed in the pattern (not just 1 valid + 1 invalid)
- Cover boundary conditions, format variations, special characters, optional vs required parts
- Example: URL pattern `^(http|https)://[a-z0-9-]+\.[a-z]{2,}(/.*)?$` needs tests for:
  - Both schemes (http, https)
  - Invalid scheme
  - Hostname format variations
  - Missing/invalid TLD
  - Path present/absent (optional part)
  - Invalid path characters

**CEL Custom Logic:**
- Test all branches, conditions, and field combinations in YOUR code
- Cover edge cases in conditional logic (`if/else`, `&&`, `||`)
- Example: Struct-level rule `!(has(self.mode) && has(self.legacyMode))` needs tests for:
  - Only mode present
  - Only legacyMode present
  - Both present (invalid)
  - Neither present

**Cross-Field Validation:**
- Test all valid and invalid field combinations
- Cover discriminated unions, dependent fields, mutually exclusive fields
- Example: Field A valid when field B="x", invalid when B="y"
- **Important**: Validation must be on parent struct, not individual fields - see [3.8 XValidation Placement](#38-error-cross-field-validation-placement)

**Immutability:**
- Test change attempts at different lifecycle stages
- Test optional immutability (can't change once set, but can be initially empty)
- Test that other fields can still be updated

**State Transitions:**
- Test all allowed state changes
- Test all forbidden state changes
- Test edge cases in transition logic

#### 5.1.3 Minimal Coverage Acceptable (Standard Libraries)

When validation **delegates to standard library functions**, limit testing to verification of correct invocation:

**Applies to:**
- CEL format functions: `format.dns1123Label()`, `format.uuid()`, etc.
- CEL IP/CIDR validators: `isIP()`, `isCIDR()`, `ip().family()`, `cidr().containsIP()`, etc.
- Kubebuilder Format markers (non-prohibited): `Format=uuid`, `Format=uri`, `Format=date`, etc. (NOT ipv4, ipv6, cidr)
- CEL URL functions: `url().getScheme()`, `url().getHostname()`, `url().getPort()`, etc.
- CEL string functions when used purely: `.matches()`, `.contains()`, `.split()`

**Coverage needed:**
- **One positive test**: Verify a valid value is accepted
- **One negative test**: Verify an invalid value is rejected

**Rationale**: Upstream libraries (Kubernetes, Go standard library) handle comprehensive edge case testing. Our tests verify correct integration, not library correctness.

**Example - Minimal Coverage:**
```go
// Just using standard library function
// +kubebuilder:validation:Format=uuid
UID string `json:"uid"`

// Tests needed:
✓ Valid UUID accepted (e.g., 550e8400-e29b-41d4-a716-446655440000)
✗ Invalid value rejected (e.g., "not-a-uuid")
```

#### 5.1.4 Pattern Test Coverage

**Rule of thumb**: If the pattern has N distinct requirements or alternatives, you need ≥N tests to cover them.

**Simple Pattern** (single requirement):
```go
// Only one requirement: lowercase letters
// +kubebuilder:validation:Pattern=`^[a-z]+$`
Field string `json:"field"`

// Tests may suffice:
✓ Valid: "abc"
✗ Invalid: "ABC123"
```

**Complex Pattern** (multiple requirements):
```go
// Multiple requirements: scheme, hostname format, TLD, optional path
// +kubebuilder:validation:Pattern=`^(http|https)://[a-z0-9-]+\.[a-z]{2,}(/.*)?$`
URL string `json:"url"`

// Comprehensive tests needed:
✓ Valid with http scheme
✓ Valid with https scheme
✓ Valid without path (optional part)
✓ Valid with path
✗ Missing scheme
✗ Invalid scheme (ftp)
✗ Invalid hostname format
✗ Missing TLD
✗ Too-short TLD
✗ Invalid path characters
```

**AWS Subnet ID Example:**
```go
// Pattern: prefix + exactly 17 alphanumeric chars
// +kubebuilder:validation:Pattern=`^subnet-[0-9A-Za-z]{17}$`
SubnetID string `json:"subnetID"`

// Comprehensive tests needed:
✓ Valid subnet ID
✗ Missing "subnet-" prefix
✗ Wrong prefix
✗ Too short (< 17 chars after prefix)
✗ Too long (> 17 chars after prefix)
✗ Invalid characters (symbols, spaces)
✗ Just "subnet-" with no ID
```

#### 5.1.5 CEL Test Coverage

**Standard Library Function Only** (minimal coverage):
```go
// Uses format marker
// +kubebuilder:validation:Format=uuid
UID string `json:"uid"`

// Tests needed:
✓ Valid UUID accepted
✗ Invalid value rejected
```

**Standard Library + Custom Logic** (combined coverage):
```go
// Combines format.dns1123Label() with custom startsWith check
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue() && self.startsWith('app-')",message="must be a valid DNS-1123 label starting with 'app-'"
Name string `json:"name"`

// Tests needed:
// Minimal for library function:
✓ Valid DNS label with prefix accepted (e.g., "app-myservice")
✗ Invalid DNS label rejected (e.g., "app-My_Service") - verifies format.dns1123Label() called

// Comprehensive for custom logic:
✗ Valid DNS label WITHOUT prefix rejected (e.g., "myservice") - tests startsWith
✗ Empty after prefix (e.g., "app-") - boundary condition
✗ Just "app" without dash - edge case
```

**Pure Custom Logic** (comprehensive coverage):
```go
// No standard library, pure custom logic on parent struct
// +kubebuilder:validation:XValidation:rule="!(has(self.mode) && has(self.legacyMode))",message="mode and legacyMode are mutually exclusive"
type MySpec struct {
    // mode determines the operating mode.
    // +optional
    Mode *string `json:"mode,omitempty"`

    // legacyMode determines legacy operating mode.
    // +optional
    LegacyMode *string `json:"legacyMode,omitempty"`
}

// Comprehensive tests needed for all combinations:
✓ Only mode present
✓ Only legacyMode present
✓ Neither present (both omitted)
✗ Both present (the invalid case)
```

**Complex CEL with Multiple Conditions:**
```go
// +kubebuilder:validation:XValidation:rule="self.type == 'external' ? has(self.externalURL) : has(self.internalService)"

// Comprehensive tests needed:
✓ type="external" with externalURL present
✓ type="internal" with internalService present
✗ type="external" without externalURL (violates conditional)
✗ type="internal" without internalService (violates conditional)
✗ type="external" with internalService instead
✗ type="other" value (if not enum-constrained)
```

#### 5.1.6 Quick Coverage Decision Tree

```
What kind of validation is this?

Standard library function ONLY?
├─ Format=uuid, isIP(), isCIDR(), format.dns1123Label(), url().getScheme(), etc.
└─ Coverage: Minimal (1 valid + 1 invalid)

Standard library + custom logic?
├─ format.dns1123Label() && self.startsWith('app-')
├─ Format=uri && url(self).getScheme() in ['http', 'https']
└─ Coverage: Minimal for library part + Comprehensive for custom part

Pure custom logic?
├─ Pattern with multiple requirements
├─ CEL with conditions/branches
├─ Cross-field validation
└─ Coverage: Comprehensive (all paths, branches, combinations)

Framework-enforced only?
├─ MinLength, MaxLength, Enum, Required, Optional
└─ Coverage: None needed
```

### 5.2 How to Establish the Intent of a Validation

Validations should be described fully with words in the godoc. The message of a specific validation rule describes the aspect of the complete validation covered by that rule. The human readable text should be considered the intent of the validation when writing tests.

**Action for AI agents**: Highlight to the user any validation which is insufficiently described to write a test, or whose implementation is very clearly different from the description. Do not attempt to validate the implementation in detail: the tests will do this.

### 5.3 Scope of a Test Suite

There should be a test suite for each 'root' object marked with `+kubebuilder:object:root=true`.

---

## 6. Validator Quick Lookup

### 6.1 Quick Reference Table

| Use Case | Validator | Type | Max Length |
|----------|-----------|------|------------|
| **Network & Communication** |
| IPv4 address | `isIP() + ip().family() == 4` | CEL | - |
| IPv6 address | `isIP() + ip().family() == 6` | CEL | - |
| CIDR notation | `isCIDR()` | CEL | - |
| URI | `Format=uri` | Kubebuilder | - |
| Email address | `Format=email` | Kubebuilder | - |
| Hostname | `Format=hostname` | Kubebuilder | 253 |
| MAC address | `Format=mac` | Kubebuilder | 17 |
| **Identifiers** |
| UUID (any version) | `Format=uuid` | Kubebuilder | 36 |
| UUID version 3 | `Format=uuid3` | Kubebuilder | 36 |
| UUID version 4 | `Format=uuid4` | Kubebuilder | 36 |
| UUID version 5 | `Format=uuid5` | Kubebuilder | 36 |
| **Date & Time** |
| Date (YYYY-MM-DD) | `Format=date` | Kubebuilder | 10 |
| Date-time (RFC 3339) | `Format=date-time` | Kubebuilder | 35 |
| Duration | `Format=duration` | Kubebuilder | - |
| **Data Encoding** |
| Base64 data | `Format=byte` | Kubebuilder | - |
| **Kubernetes (1.34+)** |
| K8s short name | `Format=k8s-short-name` | Kubebuilder | 63 |
| K8s long name | `Format=k8s-long-name` | Kubebuilder | 253 |
| **DNS (CEL)** |
| DNS-1123 label | `format.dns1123Label()` | CEL | 63 |
| DNS-1123 subdomain | `format.dns1123Subdomain()` | CEL | 253 |
| DNS-1035 label | `format.dns1035Label()` | CEL | 63 |
| **Kubernetes Names (CEL)** |
| Qualified name | `format.qualifiedName()` | CEL | 316 |
| Label value | `format.labelValue()` | CEL | 63 |

---

## 7. Complete Validator Catalog

### 7.1 Kubernetes Format Markers

Use these markers with `+kubebuilder:validation:Format=<name>`. **All format markers only apply to string-typed fields.**

**Cross-reference**: Many format markers have equivalent CEL validators in [Section 7.2](#72-cel-format-validators). See individual entries for CEL alternatives.

#### 7.1.0 Security: Prohibited Format Markers

The following format markers **MUST NOT be used** due to security vulnerabilities that allow attackers to bypass validation through malformed input:

| Prohibited Marker | CVE References | Required CEL Alternative | Section |
|-------------------|----------------|--------------------------|---------|
| `ipv4` | CVE-2024-24790<br>CVE-2021-29923 | `isIP(self) && ip(self).family() == 4` | [7.2.4](#724-ip-address-validation) |
| `ipv6` | CVE-2024-24790<br>CVE-2021-29923 | `isIP(self) && ip(self).family() == 6` | [7.2.4](#724-ip-address-validation) |
| `cidr` | CVE-2024-24790<br>CVE-2021-29923 | `isCIDR(self)` | [7.2.5](#725-cidr-validation) |

**Security Impact**: These vulnerabilities allow attackers to craft malformed IP addresses or CIDR ranges that pass Format marker validation but are rejected or misinterpreted by downstream systems, potentially leading to security bypasses, denial of service, or other exploitation.

**Enforcement**: API reviews will reject any use of these markers. Always use the CEL alternatives listed above.

**Example Migration**:
```go
// PROHIBITED - Security Vulnerability
// ipAddress is the node IP address.
// +kubebuilder:validation:Format=ipv4  // ❌ CVE-2024-24790, CVE-2021-29923
IPAddress string `json:"ipAddress"`

// REQUIRED - Secure Alternative
// ipAddress is the node IP address.
// Must be a valid IPv4 address in canonical form (e.g., 192.168.1.1).
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).family() == 4",message="must be a valid IPv4 address"
IPAddress string `json:"ipAddress"`
```

#### 7.1.1 Network & Communication Formats

**`uri`** - URI following RFC 3986 syntax
- *Description*: "Must be a valid URI including scheme (e.g., https://example.com/path?query=value)"
- *Constraints*: None (validated by Go's net/url.ParseRequestURI)
- *CEL Alternative*: [format.uri()](#7232-crd-format-validators) or [isURL()](#7233-url-validation-and-parsing)

**`email`** - Email address
- *Description*: "Must be a valid email address (e.g., user@example.com)"
- *Constraints*: None (validated by Go's net/mail.ParseAddress)

**`hostname`** - Internet hostname per RFC 1034
- *Description*: "Must be a valid hostname consisting of labels separated by dots (e.g., api.example.com)"
- *Constraints*: MaxLength: 253

**`ipv4`** - IPv4 address **[PROHIBITED - SECURITY VULNERABILITY]**
- **DO NOT USE**: This format marker is vulnerable to CVE-2024-24790 and CVE-2021-29923
- **Use Instead**: CEL validator `isIP(self) && ip(self).family() == 4` - see [7.2.4 IP Address Validation](#724-ip-address-validation)
- **Security Risk**: Allows malformed IPv4 addresses that bypass validation
- *Historical Description*: Previously validated IPv4 addresses in dotted-quad notation (e.g., 192.168.1.1)
- *Historical Constraints*: MinLength: 7, MaxLength: 15

**`ipv6`** - IPv6 address **[PROHIBITED - SECURITY VULNERABILITY]**
- **DO NOT USE**: This format marker is vulnerable to CVE-2024-24790 and CVE-2021-29923
- **Use Instead**: CEL validator `isIP(self) && ip(self).family() == 6` - see [7.2.4 IP Address Validation](#724-ip-address-validation)
- **Security Risk**: Allows malformed IPv6 addresses that bypass validation
- *Historical Description*: Previously validated IPv6 addresses (e.g., 2001:db8::1 or ::1)
- *Historical Constraints*: MinLength: 2, MaxLength: 45

**`cidr`** - CIDR notation **[PROHIBITED - SECURITY VULNERABILITY]**
- **DO NOT USE**: This format marker is vulnerable to CVE-2024-24790 and CVE-2021-29923
- **Use Instead**: CEL validator `isCIDR(self)` - see [7.2.5 CIDR Validation](#725-cidr-validation)
- **Security Risk**: Allows malformed CIDR ranges that bypass validation
- *Historical Description*: Previously validated CIDR notation IP address ranges (e.g., 10.0.0.0/8 or fd00::/8)
- *Historical Constraints*: MaxLength: 49 (IPv6 address + /128 suffix: 45 + 4 = 49)

**`mac`** - MAC address
- *Description*: "Must be a valid MAC address with consistent separators (e.g., 3D:F2:C9:A6:B3:4F, 3D-F2-C9-A6-B3-4F, or 3DF2.C9A6.B34F). Case-insensitive."
- *Constraints*: MaxLength: 17

#### 7.1.2 Identifier Formats

**`uuid`** - UUID (any version)
- *Description*: "Must be a valid UUID in 8-4-4-4-12 format (e.g., 550e8400-e29b-41d4-a716-446655440000)"
- *Constraints*: Length: exactly 36 characters (MinLength: 36, MaxLength: 36)
- *CEL Alternative*: [format.uuid()](#7232-crd-format-validators)

**`uuid3`** - UUID version 3 specifically
- *Description*: "Must be a valid UUID version 3 (MD5 hash-based)"
- *Constraints*: Length: exactly 36 characters (MinLength: 36, MaxLength: 36)

**`uuid4`** - UUID version 4 specifically
- *Description*: "Must be a valid UUID version 4 (random)"
- *Constraints*: Length: exactly 36 characters (MinLength: 36, MaxLength: 36)

**`uuid5`** - UUID version 5 specifically
- *Description*: "Must be a valid UUID version 5 (SHA-1 hash-based)"
- *Constraints*: Length: exactly 36 characters (MinLength: 36, MaxLength: 36)

#### 7.1.3 Data Encoding Formats

**`byte`** - Base64-encoded data
- *Description*: "Must be a valid base64-encoded string"
- *Constraints*: None
- *CEL Alternative*: [format.byte()](#7232-crd-format-validators)

**`password`** - Password (no actual validation)
- *Description*: "Any string value (marker for sensitive data, not validated)"
- *Constraints*: None

#### 7.1.4 Date & Time Formats

**`date`** - RFC 3339 full-date
- *Description*: "Must be a valid date in RFC 3339 format (YYYY-MM-DD, e.g., 2024-01-15)"
- *Constraints*: MinLength: 10, MaxLength: 10
- *CEL Alternative*: [format.date()](#7232-crd-format-validators)

**`duration`** - Duration string
- *Description*: "Must be a valid duration string parseable by Go's time.ParseDuration (e.g., 30s, 5m, 1h30m). Valid units: ns, us (µs), ms, s, m, h"
- *Constraints*: None

**`date-time`** (or **`datetime`**) - RFC 3339 date-time
- *Description*: "Must be a valid RFC 3339 date-time (e.g., 2024-01-15T14:30:00Z or 2024-01-15T14:30:00+00:00)"
- *Constraints*: MinLength: 20, MaxLength: 35
- *CEL Alternative*: [format.datetime()](#7232-crd-format-validators)

#### 7.1.5 Kubernetes-Specific Formats (1.34+)

**`k8s-short-name`** - Kubernetes short name (DNS label)
- *Description*: "Must be a valid Kubernetes resource name: lowercase alphanumeric characters or '-', starting and ending with alphanumeric"
- *Constraints*: MinLength: 1, MaxLength: 63

**`k8s-long-name`** - Kubernetes long name (DNS subdomain)
- *Description*: "Must be a valid Kubernetes qualified name: lowercase alphanumeric characters, '-', or '.', starting and ending with alphanumeric"
- *Constraints*: MinLength: 1, MaxLength: 253

#### 7.1.6 Deprecated/Discouraged Formats

The following format markers are supported but should **not** be used in new APIs as they are not relevant to Kubernetes use cases: `bsonobjectid`, `isbn`, `isbn10`, `isbn13`, `creditcard`, `ssn`, `hexcolor`, `rgbcolor`.

### 7.2 CEL Format Validators

Use these functions in `+kubebuilder:validation:XValidation` rules.

**Important**: Format validators return an optional error. Use the pattern `!format.validator().validate(self).hasValue()` to check validity (the `!` negates the error presence, so true when validation succeeds).

**Cross-reference**: Many CEL validators have equivalent kubebuilder Format markers in [Section 7.1](#71-kubernetes-format-markers).

#### 7.2.1 DNS and Naming Formats

**`format.dns1123Label()`**
- *MaxLength*: 63
- *Usage*: `!format.dns1123Label().validate(self).hasValue()`
- *Example*: `format.dns1123Label().validate("my-label-name")` → no error
- *Format Marker*: None (use CEL or Pattern)

**`format.dns1123Subdomain()`**
- *MaxLength*: 253
- *Usage*: `!format.dns1123Subdomain().validate(self).hasValue()`
- *Example*: `format.dns1123Subdomain().validate("apiextensions.k8s.io")` → no error
- *Format Marker*: None (use CEL or Pattern)

**`format.dns1035Label()`**
- *MaxLength*: 63
- *Usage*: `!format.dns1035Label().validate(self).hasValue()`
- *Example*: `format.dns1035Label().validate("my-label-name")` → no error
- *Format Marker*: None (use CEL or Pattern)

**`format.qualifiedName()`**
- *MaxLength*: 316
- *Usage*: `!format.qualifiedName().validate(self).hasValue()`
- *Example*: `format.qualifiedName().validate("apiextensions.k8s.io/v1beta1")` → no error
- *Format Marker*: None (use CEL or Pattern)

**`format.labelValue()`**
- *MaxLength*: 63
- *Usage*: `!format.labelValue().validate(self).hasValue()`
- *Example*: `format.labelValue().validate("my-cool-label-Value")` → no error
- *Format Marker*: None (use CEL or Pattern)

**`format.dns1123LabelPrefix()`**
- *MaxLength*: 63
- *Usage*: `!format.dns1123LabelPrefix().validate(self).hasValue()`
- *Example*: `format.dns1123LabelPrefix().validate("my-label-prefix-")` → no error
- *Note*: Allows trailing dash for concatenation
- *Format Marker*: None (use CEL or Pattern)

**`format.dns1123SubdomainPrefix()`**
- *MaxLength*: 253
- *Usage*: `!format.dns1123SubdomainPrefix().validate(self).hasValue()`
- *Example*: `format.dns1123SubdomainPrefix().validate("mysubdomain.prefix.-")` → no error
- *Note*: Allows trailing dash
- *Format Marker*: None (use CEL or Pattern)

**`format.dns1035LabelPrefix()`**
- *MaxLength*: 63
- *Usage*: `!format.dns1035LabelPrefix().validate(self).hasValue()`
- *Example*: `format.dns1035LabelPrefix().validate("my-label-prefix-")` → no error
- *Note*: Allows trailing dash
- *Format Marker*: None (use CEL or Pattern)

#### 7.2.2 CRD Format Validators

**`format.uri()`**
- *MaxLength*: None (typically constrained by application)
- *Usage*: `!format.uri().validate(self).hasValue()`
- *Example*: `format.uri().validate("http://example.com")` → no error
- *Format Marker*: [uri](#7111-network--communication-formats)

**`format.uuid()`**
- *MaxLength*: 36
- *Usage*: `!format.uuid().validate(self).hasValue()`
- *Example*: `format.uuid().validate("123e4567-e89b-12d3-a456-426614174000")` → no error
- *Format Marker*: [uuid](#7112-identifier-formats)

**`format.byte()`**
- *MaxLength*: None (depends on encoded data size)
- *Usage*: `!format.byte().validate(self).hasValue()`
- *Example*: `format.byte().validate("aGVsbG8=")` → no error
- *Format Marker*: [byte](#7113-data-encoding-formats)

**`format.date()`**
- *MaxLength*: 10
- *Usage*: `!format.date().validate(self).hasValue()`
- *Example*: `format.date().validate("2021-01-01")` → no error
- *Format Marker*: [date](#7114-date--time-formats)

**`format.datetime()`**
- *MaxLength*: 35
- *Usage*: `!format.datetime().validate(self).hasValue()`
- *Example*: `format.datetime().validate("2021-01-01T00:00:00Z")` → no error
- *Format Marker*: [date-time](#7114-date--time-formats)

**`format.named(name: string)`**
- *Returns*: The Format validator with the given name, or `optional.none` if not found
- *Allowed Names*: `dns1123Label`, `dns1123Subdomain`, `dns1035Label`, `qualifiedName`, `dns1123LabelPrefix`, `dns1123SubdomainPrefix`, `dns1035LabelPrefix`, `labelValue`, `uri`, `uuid`, `byte`, `date`, `datetime`
- *Usage*: `!format.named("dns1123Label").validate(self).hasValue()`
- *Example*: `format.named("dns1123Label").hasValue()` → check if format exists

#### 7.2.3 URL Validation and Parsing

**`isURL(string)`**
- *Returns*: Boolean (true if valid URL)
- *Usage*: `isURL(self)`
- *Examples*:
  - `isURL('https://user:pass@example.com:80/path?query=val#fragment')` → true
  - `isURL('/absolute-path')` → true
  - `isURL('https://a:b:c/')` → false (invalid port)
  - `isURL('../relative-path')` → false

**`url(string)`**
- *Returns*: URL type or error if invalid
- *Usage*: `url(self).getScheme()`
- *Format Marker*: [uri](#7111-network--communication-formats)

**URL Object Methods**:
- `url().getScheme()` - Example: `url('https://example.com').getScheme()` → "https"
- `url().getHost()` - Example: `url('https://example.com:443').getHost()` → "example.com:443"
- `url().getHostname()` - Example: `url('https://example.com:443').getHostname()` → "example.com"
- `url().getPort()` - Example: `url('https://example.com:443').getPort()` → "443"
- `url().getEscapedPath()` - Example: `url('https://example.com/api/v1').getEscapedPath()` → "/api/v1" (URL-encoded path)
- `url().getQuery()` - Example: `url('https://example.com?foo=bar&foo=baz').getQuery()` → {"foo": ["bar", "baz"]} (returns map<string, list<string>> without '?')

#### 7.2.4 IP Address Validation

**`isIP(string)`**
- *Returns*: Boolean (true if valid IP)
- *Rejects*: IPv4-mapped IPv6 addresses (e.g., `::ffff:1.2.3.4`), IP addresses with zones (e.g., `fe80::1%eth0`), leading zeros in IPv4 octets
- *Usage*: `isIP(self)`
- *Examples*:
  - `isIP('127.0.0.1')` → true
  - `isIP('::1')` → true
  - `isIP('127.0.0.256')` → false
- *Format Markers*: [ipv4](#7111-network--communication-formats), [ipv6](#7111-network--communication-formats)

**`ip(string)`**
- *Returns*: IP type or error if invalid
- *Usage*: `ip(self).family()`

**`ip.isCanonical(string)`**
- *Returns*: Boolean
- *Note*: All valid IPv4 addresses are canonical; IPv6 must use lowercase and minimal form
- *Usage*: `ip.isCanonical(self)`
- *Examples*:
  - `ip.isCanonical('127.0.0.1')` → true
  - `ip.isCanonical('2001:db8::abcd')` → true
  - `ip.isCanonical('2001:DB8::ABCD')` → false (uppercase not canonical)

**IP Object Methods**:
- `ip().family()` - Returns 4 or 6 (int). Example: `ip('127.0.0.1').family()` → 4
- `ip().isUnspecified()` - Returns true for 0.0.0.0 or ::
- `ip().isLoopback()` - Returns true for 127.x.x.x or ::1
- `ip().isLinkLocalMulticast()` - Returns true for 224.0.0.x or ff02::/16
- `ip().isLinkLocalUnicast()` - Returns true for 169.254.x.x or fe80::/10
- `ip().isGlobalUnicast()` - Returns true for global unicast addresses

#### 7.2.5 CIDR Validation

**`isCIDR(string)`**
- *Returns*: Boolean (true if valid CIDR)
- *Rejects*: IPv4-mapped IPv6 addresses, leading zeros in IPv4 octets
- *Usage*: `isCIDR(self)`
- *Examples*:
  - `isCIDR('192.168.0.0/16')` → true
  - `isCIDR('::1/128')` → true
  - `isCIDR('192.168.0.0/33')` → false (invalid prefix length)
- *Format Marker*: [cidr](#7111-network--communication-formats)

**`cidr(string)`**
- *Returns*: CIDR type or error if invalid
- *Usage*: `cidr(self).prefixLength()`

**CIDR Object Methods**:
- `cidr().containsIP(ip)` - Returns true if CIDR contains the IP (accepts IP object or string)
- `cidr().containsCIDR(cidr)` - Returns true if CIDR contains another CIDR (accepts CIDR object or string)
- `cidr().ip()` - Returns the IP address representation of the CIDR
- `cidr().masked()` - Returns canonical form with host bits zeroed. Example: `cidr('192.168.1.5/24').masked()` → "192.168.1.0/24"
- `cidr().prefixLength()` - Returns prefix length in bits (int)

#### 7.2.6 Kubernetes Quantity

**`isQuantity(string)`**
- *Returns*: Boolean (true if valid Quantity)
- *Accepts*: Patterns like "1.5G", "200k", "1.3Gi" (binary and decimal suffixes)
- *Usage*: `isQuantity(self)`
- *Examples*:
  - `isQuantity('1.3G')` → true
  - `isQuantity('1.3Gi')` → true
  - `isQuantity('200K')` → false (uppercase K not allowed)

**`quantity(string)`**
- *Returns*: Quantity type or error if invalid
- *Usage*: `quantity(self).asInteger()`

**Quantity Object Methods**:
- `quantity().isInteger()` - Checks if asInteger() is safe to call (won't overflow or lose precision)
- `quantity().asInteger()` - Returns int64 value (errors on overflow/precision loss)
- `quantity().asApproximateFloat()` - Returns float64 (may lose precision)
- `quantity().sign()` - Returns 1, -1, or 0
- `quantity().add(quantity|int)` - Returns sum
- `quantity().sub(quantity|int)` - Returns difference
- `quantity().isGreaterThan(quantity)` - Returns boolean
- `quantity().isLessThan(quantity)` - Returns boolean
- `quantity().compareTo(quantity)` - Returns 0, 1, or -1

#### 7.2.7 Semantic Versions

**`isSemver(string)` / `isSemver(string, bool)`**
- *Returns*: Boolean (true if valid semver)
- *Parameters*: Optional second bool parameter enables normalization (removes "v" prefix, adds missing minor/patch zeros, removes leading zeros)
- *Usage*: `isSemver(self)` or `isSemver(self, true)`
- *Examples*:
  - `isSemver('1.0.0')` → true
  - `isSemver('v1.0', true)` → true (with normalization)
  - `isSemver('v1.0', false)` → false (without normalization)

**`semver(string)` / `semver(string, bool)`**
- *Returns*: Semver type or error if invalid
- *Parameters*: Optional second bool parameter enables normalization
- *Usage*: `semver(self).major()`

**Semver Object Methods**:
- `semver().major()` - Returns major version (int)
- `semver().minor()` - Returns minor version (int)
- `semver().patch()` - Returns patch version (int)
- `semver().isGreaterThan(semver)` - Returns boolean
- `semver().isLessThan(semver)` - Returns boolean
- `semver().compareTo(semver)` - Returns 0, 1, or -1

#### 7.2.8 Pattern Matching

**`string.find(pattern)`**
- *Returns*: First substring matching the regex pattern (string)
- *Usage*: `self.find('[0-9]+')`
- *Example*: `"abc 123".find('[0-9]+')` → "123"

**`string.findAll(pattern)` / `string.findAll(pattern, limit)`**
- *Returns*: List of all substrings matching the regex pattern
- *Parameters*: Optional limit parameter restricts result count
- *Usage*: `self.findAll('[0-9]+')`
- *Examples*:
  - `"123 abc 456".findAll('[0-9]+')` → ["123", "456"]
  - `"123 abc 456".findAll('[0-9]+', 1)` → ["123"]

---

## 8. Test Framework Reference

### 8.1 Test Suite File Structure

Test suites are YAML files with the following structure:

```yaml
apiVersion: apiextensions.k8s.io/v1  # Required (controller-gen requirement)
name: <Kind>                         # The Kind name of the API under test
crdName: <plural>.<group>            # Fully qualified CRD name (e.g., routes.route.openshift.io)
featureGates:                        # Optional: Required feature gates
- FeatureGateName                    # Feature gate that must be enabled
- -DisabledFeatureGate               # Prefix with '-' for gates that must be disabled
version: v1                          # Optional: API version (auto-detected if only one version exists)
tests:
  onCreate: [...]                    # Tests for resource creation
  onUpdate: [...]                    # Tests for resource updates
```

### 8.2 Test File Location and Naming

**Location**: `<group>/<version>/tests/<crd-name>/`

**Naming Conventions**:
- **`AAA_ungated.yaml`**: Base test suite that runs without feature gate requirements (the "AAA" prefix ensures it runs first)
- **`<FeatureGateName>.yaml`**: Tests that require the named feature gate to be enabled
- **`<FeatureName>+<FeatureName2>.yaml`**: Tests requiring multiple feature gates
- **`-<FeatureGateName>.yaml`**: Tests that run only when the feature gate is disabled
- **Descriptive suffixes**: Optional context (e.g., `stable.controlplanemachineset.aws.testsuite.yaml`)

**Critical Rule**: There must be one test suite file for every file in `<group>/<version>/zz_generated.featuregated-crd-manifests/<crd-name>/`. Test suite files must have matching names.

### 8.3 onCreate Tests

**Purpose**: Test resource creation, field validation, default value application, and validation rules that apply during creation.

**Limitations**: onCreate tests cannot test status field creation. Use onUpdate tests for status validation.

#### Basic Structure

```yaml
tests:
  onCreate:
    - name: Descriptive test name
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: value
      expected: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: value
          defaultedField: defaultValue  # Shows fields that were defaulted
```

#### Validation Failure Structure

When a test should generate an error, omit `expected` and specify `expectedError`:

```yaml
tests:
  onCreate:
    - name: Should reject invalid pattern
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          emailField: not-an-email
      expectedError: "failed to match pattern"
```

**expectedError Best Practices**:
- Use a distinctive substring of the error message
- Include the validation rule portion (e.g., the part after `[` up to the first `,` for CEL errors)
- Avoid dynamic portions like timestamps or resource versions
- Be specific enough to catch only the right error but flexible enough to handle minor message changes

#### Metadata Handling

- Do **not** include `metadata.name` or `metadata.namespace` in test cases for namespace-scoped resources - these are auto-generated
- For cluster-scoped resources (like ClusterOperator), include `metadata.name` if required by the API

#### Required Tests for onCreate

1. **Minimal creation test**: Every API should have a test showing minimal valid object creation with expected defaults
2. **Pattern validation**: Test invalid patterns
3. **CEL validation**: Test both valid and invalid cases for each CEL rule
4. **Cross-field validation**: Test invalid field combinations
5. **Format validation**: Test invalid formats (IP addresses, CIDRs, etc.)

### 8.4 onUpdate Tests

**Purpose**: Test update operations, immutability constraints, status subresource validation, state transitions, and validation ratcheting.

**Key Constraint**: The `initial` object must always be valid. To test validation on updates, create a valid initial object then attempt an invalid update.

#### Basic Update Structure

```yaml
tests:
  onUpdate:
    - name: Descriptive update test name
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: initialValue
      updated: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: updatedValue
      expected: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: updatedValue
```

#### Status Field Testing

Status fields are applied **separately** from spec fields:

```yaml
tests:
  onUpdate:
    - name: Should allow updating status field
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          specField: value
        status:               # Applied after spec is created
          statusField: initialStatus
      updated: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          specField: value
        status:
          statusField: updatedStatus
      expected: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          specField: value
        status:
          statusField: updatedStatus
```

#### Error Handling: expectedError vs expectedStatusError

- **`expectedError`**: Use for spec field validation failures
- **`expectedStatusError`**: Use for status subresource validation failures

```yaml
tests:
  onUpdate:
    - name: Should reject invalid status value
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: value
        status:
          statusField: validValue
      updated: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: value
        status:
          statusField: invalidValue
      expectedStatusError: "substring of status validation error"
```

#### Required Tests for onUpdate

1. **Immutability**: Test that immutable fields cannot be changed
2. **Optional immutability**: Test fields that cannot be changed once set but can initially be empty
3. **Status immutability**: Test that status fields follow their immutability rules
4. **State transitions**: Test allowed and forbidden state transitions
5. **Validation ratcheting**: For fields with validation rules that may evolve (see below)

### 8.5 Validation Ratcheting with initialCRDPatches

**Purpose**: Test that existing resources with values that were valid under old validation rules can continue to be updated (for fields other than the now-invalid value) even after validation rules are tightened.

**Kubernetes Validation Ratcheting**: Kubernetes allows updates to resources that have fields violating current validation rules, as long as those specific invalid fields are not changed. This allows operators to update other fields while preserving existing invalid values that were valid when created.

#### initialCRDPatches Structure

Use JSON Patch (RFC 6902) operations to modify the CRD before creating the initial object:

```yaml
tests:
  onUpdate:
    - name: Should allow changing other fields when persisted value is no longer valid
      initialCRDPatches:
        - op: remove
          path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/field/minimum
        - op: replace
          path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/field/maximum
          value: 1000
      initial: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: 5              # Valid with patched CRD, invalid with current CRD
          otherField: foo
      updated: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: 5              # Keep the now-invalid value unchanged
          otherField: bar       # Change other fields
      expected: |
        apiVersion: <group>/<version>
        kind: <Kind>
        spec:
          field: 5
          otherField: bar
```

**Supported JSON Patch Operations**:
- `add`: Add a new value at path
- `remove`: Remove value at path
- `replace`: Replace value at path
- `move`: Move value from one path to another
- `copy`: Copy value from one path to another
- `test`: Test that a value exists at path

#### The Three Required Ratcheting Tests

For each field with validation rules that may evolve, create three tests:

1. **Allow updates to other fields while invalid value persists**
   ```yaml
   - name: Should allow changing other fields when persisted value is no longer valid
     initialCRDPatches:
       - op: remove
         path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/field/minimum
     initial: |
       spec:
         field: 1              # Now violates minimum of 8
         otherField: foo
     updated: |
       spec:
         field: 1              # Keep invalid value
         otherField: bar       # But can change other fields
     expected: |
       spec:
         field: 1
         otherField: bar
   ```

2. **Reject updates from invalid to another invalid value**
   ```yaml
   - name: Should reject changing a persisted invalid value to another invalid value
     initialCRDPatches:
       - op: remove
         path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/field/minimum
     initial: |
       spec:
         field: 1              # Violates minimum of 8
     updated: |
       spec:
         field: 5              # Still violates minimum of 8
     expectedError: "Invalid value"
   ```

3. **Allow updates from invalid to valid value**
   ```yaml
   - name: Should allow changing a persisted invalid value to a valid value
     initialCRDPatches:
       - op: remove
         path: /spec/versions/0/schema/openAPIV3Schema/properties/spec/properties/field/minimum
     initial: |
       spec:
         field: 1              # Violates minimum of 8
     updated: |
       spec:
         field: 10             # Now valid
     expected: |
       spec:
         field: 10
   ```

#### When to Use Validation Ratcheting Tests

Write ratcheting tests for:
- Validation rules added to existing fields in stable APIs
- Tightened validation constraints (stricter min/max, more restrictive patterns)
- New enum restrictions
- New CEL validation rules on existing fields

Do **not** write ratcheting tests for:
- New fields (no existing values to preserve)
- Tech preview APIs where breaking changes are acceptable
- Required fields (cannot be omitted)

### 8.6 Test Execution

#### Running Tests

```bash
# Run tests for a specific API group/version (recommended for focused work)
make -C <group>/<version> test

# Examples:
make -C config/v1 test
make -C operator/v1 test
make -C route/v1 test

# Run all integration tests
make integration
```

#### Test Output Handling

Always capture test output to a temporary file for analysis:

```bash
make -C apiextensions/v1alpha1 test > /tmp/test-output.txt 2>&1
```

This allows extraction of both summary and detailed output without re-running tests.

#### Code Regeneration Before Testing

**Critical**: When API types are modified, CRDs must be regenerated before running tests:

```bash
# Regenerate CRDs for specific API group (faster)
make update-codegen API_GROUP_VERSIONS=operator.openshift.io/v1alpha1

# Or regenerate all CRDs
make update-codegen

# Then run tests
make -C operator/v1alpha1 test
```

### 8.7 Test Template Generation

Use the test generation script to create initial test structure:

```bash
./tests/hack/gen-minimal-test.sh <folder> <version>

# Examples:
./tests/hack/gen-minimal-test.sh operator/v1 v1
./tests/hack/gen-minimal-test.sh config/v1alpha1 v1alpha1
```

The script:
1. Scans `zz_generated.featuregated-crd-manifests/` for CRD files
2. Creates corresponding test files in `tests/<crd-name>/`
3. Generates minimal test structure with basic onCreate test
4. Creates Makefile if needed

**Note**: Generated tests are minimal starting points. Add comprehensive validation tests based on the API's validation rules.

---

## 9. Complete Examples

### 9.1 Example 1: UUID Field with Format Marker

**Use Case**: Field requiring UUID validation

**Code**:
```go
// uid is the unique identifier for this resource.
// Must be a valid UUID in 8-4-4-4-12 format (e.g., 550e8400-e29b-41d4-a716-446655440000).
// Length must be exactly 36 characters.
// +kubebuilder:validation:Format=uuid
// +kubebuilder:validation:MinLength=36
// +kubebuilder:validation:MaxLength=36
UID string `json:"uid"`
```

**Tests Required**: Yes (format validation)

**Coverage Strategy**: Minimal (standard library only) - see [5.1.3 Minimal Coverage](#513-minimal-coverage-acceptable-standard-libraries)
- `Format=uuid` is a standard library format marker
- Need only 1 valid + 1 invalid test
- `MinLength`/`MaxLength` are framework-enforced (no tests needed)

**Test Structure**:
```yaml
onCreate:
  - name: Should accept valid UUID
    initial: |
      spec:
        uid: 550e8400-e29b-41d4-a716-446655440000
    expected: |
      spec:
        uid: 550e8400-e29b-41d4-a716-446655440000

  - name: Should reject invalid UUID format
    initial: |
      spec:
        uid: not-a-uuid
    expectedError: "must be a valid UUID"
```

**Common Mistakes**: Using custom pattern instead of Format=uuid

### 9.2 Example 2: DNS Label with Custom Constraints

**Use Case**: Field requiring DNS-1123 label validation plus application-specific prefix

**Code**:
```go
// name is the application name.
// Must be a valid DNS-1123 label starting with 'app-' (e.g., app-myservice).
// Maximum length is 63 characters.
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue() && self.startsWith('app-')",message="must be a valid DNS-1123 label starting with 'app-'"
// +kubebuilder:validation:MaxLength=63
Name string `json:"name"`
```

**Tests Required**: Yes (CEL validation)

**Coverage Strategy**: Combined (standard library + custom logic) - see [5.1.5 CEL Test Coverage](#515-cel-test-coverage)
- Minimal coverage for `format.dns1123Label()` (standard library)
- Comprehensive coverage for `startsWith('app-')` (custom logic)
- No tests for `MaxLength` (framework-enforced)

**Test Structure**:
```yaml
onCreate:
  - name: Should accept valid DNS label with app- prefix
    initial: |
      spec:
        name: app-myservice
    expected: |
      spec:
        name: app-myservice

  - name: Should reject DNS label without app- prefix
    initial: |
      spec:
        name: myservice
    expectedError: "must be a valid DNS-1123 label starting with 'app-'"

  - name: Should reject invalid DNS label format
    initial: |
      spec:
        name: app-My_Service
    expectedError: "must be a valid DNS-1123 label"

  # Additional tests for comprehensive custom logic coverage:
  - name: Should reject empty after prefix
    initial: |
      spec:
        name: app-
    expectedError: "must be a valid DNS-1123 label"
```

**Common Mistakes**: Duplicating DNS validation logic in a custom pattern

### 9.3 Example 3: Mutually Exclusive Fields

**Use Case**: Two fields that cannot be set together

**Code**:
```go
// +kubebuilder:validation:XValidation:rule="!(has(self.mode) && has(self.legacyMode))",message="mode and legacyMode are mutually exclusive"
type MySpec struct {
    // mode determines the operating mode.
    // Cannot be used together with legacyMode field.
    // +optional
    Mode *string `json:"mode,omitempty"`

    // legacyMode determines legacy operating mode.
    // Cannot be used together with mode field.
    // +optional
    LegacyMode *string `json:"legacyMode,omitempty"`
}
```

**Tests Required**: Yes (cross-field validation)

**Coverage Strategy**: Comprehensive (pure custom logic) - see [5.1.2 Comprehensive Coverage](#512-comprehensive-coverage-required)
- Test all field combinations: mode only, legacyMode only, both, neither
- Cross-field validation requires comprehensive coverage

**Test Structure**:
```yaml
onCreate:
  - name: Should accept mode without legacyMode
    initial: |
      spec:
        mode: standard
    expected: |
      spec:
        mode: standard

  - name: Should accept legacyMode without mode
    initial: |
      spec:
        legacyMode: legacy
    expected: |
      spec:
        legacyMode: legacy

  - name: Should reject both mode and legacyMode
    initial: |
      spec:
        mode: standard
        legacyMode: legacy
    expectedError: "mode and legacyMode are mutually exclusive"

  - name: Should accept neither mode nor legacyMode
    initial: |
      spec:
        otherField: value
    expected: |
      spec:
        otherField: value
```

**Common Mistakes**: Documenting mutual exclusivity without XValidation enforcement

### 9.4 Example 4: Immutable Field

**Use Case**: Field that cannot be changed after creation

**Code**:
```go
// clusterID is the unique identifier of the cluster.
// This field is immutable after creation.
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="clusterID is immutable"
// +required
ClusterID string `json:"clusterID"`
```

**Tests Required**: Yes (immutability constraint)

**Coverage Strategy**: Comprehensive (immutability) - see [5.1.2 Comprehensive Coverage](#512-comprehensive-coverage-required)
- Test that field cannot be changed
- Test that other fields can still be updated (doesn't block all updates)

**Test Structure**:
```yaml
onUpdate:
  - name: Should reject changes to immutable clusterID
    initial: |
      spec:
        clusterID: cluster-123
    updated: |
      spec:
        clusterID: cluster-456
    expectedError: "clusterID is immutable"

  - name: Should allow update without changing clusterID
    initial: |
      spec:
        clusterID: cluster-123
        otherField: foo
    updated: |
      spec:
        clusterID: cluster-123
        otherField: bar
    expected: |
      spec:
        clusterID: cluster-123
        otherField: bar
```

### 9.5 Example 5: Type-Level Validation

**Use Case**: Enum validation applied to a custom type

**Code**:
```go
// +kubebuilder:validation:Enum=Deny;Warn
type AdmitAction string

const (
    AdmitActionDeny AdmitAction = "Deny"
    AdmitActionWarn AdmitAction = "Warn"
)

type MySpec struct {
    // action determines the admission behavior.
    // Valid options are Deny and Warn.
    // When set to Deny, requests will be rejected.
    // When set to Warn, requests will be admitted with a warning.
    // +required
    Action AdmitAction `json:"action"`
}
```

**Tests Required**: No (enum is trivial validation)

**Coverage Strategy**: None - see [5.1.1 No Tests Required](#511-no-tests-required-framework-enforced)
- `Enum` is framework-enforced
- `Required` is framework-enforced

**Common Mistakes**: Duplicating the Enum marker on the field when it's already on the type

### 9.6 Example 6: Array with Item Validation

**Use Case**: Array of strings where each string must meet length constraints

**Code**:
```go
// names is a list of service names.
// Each name must be between 5 and 63 characters.
// The list must contain between 1 and 10 names.
// +kubebuilder:validation:MinItems=1
// +kubebuilder:validation:MaxItems=10
// +kubebuilder:validation:items:MinLength=5
// +kubebuilder:validation:items:MaxLength=63
Names []string `json:"names"`
```

**Tests Required**: No (length constraints are trivial)

**Coverage Strategy**: None - see [5.1.1 No Tests Required](#511-no-tests-required-framework-enforced)
- `MinItems`, `MaxItems`, `MinLength`, `MaxLength` are all framework-enforced
- No tests needed for framework-enforced validations

**Common Mistakes**: Using `MinLength` instead of `items:MinLength` for array elements

---

## 10. Reference Examples

The following test files are well-established and serve as excellent references:

### 10.1 Comprehensive Examples

- **`example/v1/tests/stableconfigtypes.example.openshift.io/AAA_ungated.yaml`**: Explicitly designed as reference implementation showing onCreate, onUpdate, immutability, and validation ratcheting patterns
- **`config/v1/tests/infrastructures.config.openshift.io/AAA_ungated.yaml`**: Extensive validation ratcheting examples, complex CEL validation, IP/CIDR validation
- **`operator/v1/tests/ingresscontrollers.operator.openshift.io/AAA_ungated.yaml`**: Discriminated unions, duration validation, domain validation

### 10.2 Pattern-Specific Examples

- **State transition validation**: `config/v1/tests/featuregates.config.openshift.io/AAA_ungated.yaml`
- **CEL cross-field validation**: `route/v1/tests/routes.route.openshift.io/AAA_ungated.yaml`
- **Feature-gated fields**: `example/v1/tests/stableconfigtypes.example.openshift.io/Example.yaml`
- **Format validation**: `config/v1/tests/apiservers.config.openshift.io/KMSEncryptionProvider.yaml`

### 10.3 Authoritative Documentation

- **`tests/README.md`**: Complete test framework documentation with detailed examples

---

## 11. CEL Format Validators vs Kubebuilder Format Markers

This section helps AI agents choose between Kubebuilder Format markers and CEL validators for common validation use cases.

**Quick Decision Rule**:
1. **Security-prohibited formats** (ipv4, ipv6, cidr) → MUST use CEL
2. **Format marker exists and not prohibited** → Prefer Format marker (simpler)
3. **No Format marker exists** → Use CEL validator
4. **Need to inspect format components** (e.g., URL parts, IP properties) → Use CEL

### Comparison Table

| Use Case | Format Marker | CEL Validator | Recommended | Rationale |
|----------|---------------|---------------|-------------|-----------|
| **Network** |
| IPv4 address | `Format=ipv4` ❌ | `isIP(self) && ip(self).family() == 4` | **CEL (Required)** | Format=ipv4 prohibited (CVE-2024-24790, CVE-2021-29923) |
| IPv6 address | `Format=ipv6` ❌ | `isIP(self) && ip(self).family() == 6` | **CEL (Required)** | Format=ipv6 prohibited (CVE-2024-24790, CVE-2021-29923) |
| CIDR notation | `Format=cidr` ❌ | `isCIDR(self)` | **CEL (Required)** | Format=cidr prohibited (CVE-2024-24790, CVE-2021-29923) |
| URI | `Format=uri` | `format.uri()` or `isURL()` | **Format Marker** | Simpler. Use CEL only if inspecting components (e.g., `url(self).getScheme()`) |
| Email address | `Format=email` | *(none)* | **Format Marker** | No CEL equivalent available |
| Hostname | `Format=hostname` | *(none)* | **Format Marker** | No CEL equivalent available |
| MAC address | `Format=mac` | *(none)* | **Format Marker** | No CEL equivalent available |
| **Identifiers** |
| UUID (any) | `Format=uuid` | `format.uuid()` | **Format Marker** | Simpler syntax, clear intent |
| UUID v3 | `Format=uuid3` | *(none)* | **Format Marker** | No CEL equivalent available |
| UUID v4 | `Format=uuid4` | *(none)* | **Format Marker** | No CEL equivalent available |
| UUID v5 | `Format=uuid5` | *(none)* | **Format Marker** | No CEL equivalent available |
| **Date & Time** |
| Date (RFC 3339) | `Format=date` | `format.date()` | **Format Marker** | Simpler syntax |
| Date-time (RFC 3339) | `Format=date-time` | `format.datetime()` | **Format Marker** | Simpler syntax |
| Duration | `Format=duration` | *(none)* | **Format Marker** | No CEL equivalent available |
| **Data Encoding** |
| Base64 | `Format=byte` | `format.byte()` | **Format Marker** | Simpler syntax |
| **DNS & K8s Names** |
| DNS-1123 label | *(none)* | `format.dns1123Label()` | **CEL** | No Format marker exists |
| DNS-1123 subdomain | *(none)* | `format.dns1123Subdomain()` | **CEL** | No Format marker exists |
| DNS-1035 label | *(none)* | `format.dns1035Label()` | **CEL** | No Format marker exists |
| Qualified name | *(none)* | `format.qualifiedName()` | **CEL** | No Format marker exists |
| Label value | *(none)* | `format.labelValue()` | **CEL** | No Format marker exists |
| K8s short name (1.34+) | `Format=k8s-short-name` | `format.dns1123Label()` | **Either** | Equivalent validations |
| K8s long name (1.34+) | `Format=k8s-long-name` | `format.dns1123Subdomain()` | **Either** | Equivalent validations |

### Usage Notes

**Format Markers**:
```go
// Simple and clear
// ipAddress must be a valid URI (e.g., https://example.com).
// +kubebuilder:validation:Format=uri
IPAddress string `json:"ipAddress"`
```

**CEL Validators**:
```go
// Required for prohibited formats
// ipAddress must be a valid IPv4 address (e.g., 192.168.1.1).
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).family() == 4",message="must be a valid IPv4 address"
IPAddress string `json:"ipAddress"`

// Use pattern for format validation
// +kubebuilder:validation:XValidation:rule="!format.dns1123Label().validate(self).hasValue()",message="must be a valid DNS-1123 label"
Name string `json:"name"`
```

**When to use CEL for component inspection**:
```go
// Need to check URL scheme - use CEL
// endpoint must be an HTTPS URL.
// +kubebuilder:validation:Format=uri
// +kubebuilder:validation:XValidation:rule="url(self).getScheme() == 'https'",message="scheme must be https"
Endpoint string `json:"endpoint"`

// Need to check IP properties - use CEL
// ip must be a loopback address.
// +kubebuilder:validation:XValidation:rule="isIP(self) && ip(self).isLoopback()",message="must be a loopback address"
IP string `json:"ip"`
```