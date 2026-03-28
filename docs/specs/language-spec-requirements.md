# UExL Language Specification Requirements

Status: Draft
Audience: Maintainers, language designers, and implementers of UExL ports
Scope: Requirements for authoring the UExL language specification, not the language specification itself

## 1. Purpose and Vision

UExL stands for Universal Expression Language.

Its long-term goal is to provide a cross-platform expression execution environment in the same spirit that RegEx provides a cross-platform pattern language: write an expression once, then run that same expression consistently across backends, frontends, API endpoints, edge runtimes, and other host environments.

UExL is now mature enough that its behavior should no longer live only in the Go implementation and explanatory books. To enable compatible ports in JavaScript, Rust, Dart, and other languages, the project needs a language specification that is independent, precise, portable, and testable.

This document defines the requirements for producing that specification.

The final specification must allow a competent team to build a conforming UExL implementation without reading the Go source code as a primary reference.

The specification effort must therefore treat source-level portability as a first-class concern. A user should be able to move the same UExL expression between environments with minimal or no rewriting, provided the target implementations conform to the same language version, profile, and mandatory core feature set.

That portability guarantee is version-scoped, not universal across all UExL releases. An expression authored for a newer language version, such as UExL 1.1, may not be valid or runnable on a host that only supports an older version, such as UExL 1.0.

## 2. Scope of This Requirements Document

This document governs the work of creating the UExL specification set.

It does not itself define the final language grammar, semantics, or built-in library in normative detail. Instead, it defines:

1. What specification artifacts must exist.
2. What subject areas those artifacts must cover.
3. What portability and conformance guarantees the finished specification must support.
4. What boundaries must be maintained between the core language, profiles, and host-specific extensions.

The portability guarantee applies to the core language and the mandatory core standard library only. It does not automatically extend to user-defined functions, custom pipes, host objects, or other injected constructs unless those are later standardized by the specification or by a named extension profile.

## 3. Goals

The UExL specification effort must achieve the following goals:

1. Define the language independently of any single implementation.
2. Preserve the current observable behavior of UExL unless an intentional language change is approved.
3. Make parser, evaluator, and standard-library behavior portable across runtimes.
4. Support at least one future non-Go implementation without requiring Go VM knowledge.
5. Separate core language semantics from host integration concerns.
6. Provide a foundation for automated conformance testing.
7. Make versioning, compatibility, and feature lifecycle rules explicit.
8. Ensure that expressions using the core built-in functions and built-in pipes evaluate the same way across all conforming ports.
9. Allow the same expression source to be reused across backend, frontend, endpoint, edge, and service environments without environment-specific rewrites in the core language.
10. Provide a clear extension boundary so host-specific functionality can exist without weakening the portability contract of the core language.

## 4. Non-Goals

The first specification effort should not attempt to do the following:

1. Standardize bytecode, VM frames, compiler internals, AST layouts, or optimization strategies for any implementation language.
2. Require identical internal representations across implementations.
3. Freeze the exact parse, compile, evaluate, or registration APIs of every host language.
4. Standardize performance targets as language semantics.
5. Introduce major new language features while the specification is being written.
6. Guarantee cross-platform portability for user-defined functions, custom pipes, or other host-injected constructs that are not part of the standardized core language profile.
7. Standardize host operational policies such as timeouts, memory limits, sandboxing, or cancellation behavior unless those policies become observable language behavior.

Performance expectations, benchmarks, and implementation architecture may be documented separately, but they are not part of the normative language contract unless explicitly promoted into the specification.

## 5. Authority and Decision Model

### 5.1 Canonical Authority

The specification must become the canonical description of UExL language behavior.

Until version 1.0 of the specification is complete, the current Go implementation plus approved tests remain the temporary behavioral reference. After version 1.0, the specification and conformance suite must be treated as the authoritative standard, and implementations must be evaluated against them.

### 5.2 Authority Hierarchy

When sources disagree, the following precedence must apply:

1. Approved specification text and approved errata
2. Approved conformance suite
3. Approved versioning, profile, and compatibility addenda
4. Current implementation behavior and approved implementation tests, but only where the specification is not yet explicit
5. User-facing books and reference documentation
6. Benchmarks, profiling notes, and implementation commentary

Any disagreement discovered during the specification effort must be resolved explicitly and recorded. Once the specification and conformance suite are updated, implementations and books must be brought into alignment.

## 6. Required Specification Artifacts

The specification effort must produce the following artifacts.

| Artifact | Required | Purpose |
|----------|----------|---------|
| Specification roadmap | Yes | Tracks sections, source materials, known ambiguities, ownership, and review state during authoring |
| Language overview | Yes | Defines scope, positioning, conformance model, and the boundaries of the language |
| Terminology and conformance glossary | Yes | Defines the canonical vocabulary used across the specification set |
| Lexical specification | Yes | Defines tokens, literals, identifiers, whitespace, and invalid forms |
| Grammar specification | Yes | Defines syntactically valid UExL expressions in one formal notation |
| Semantic specification | Yes | Defines evaluation order, observability, determinism, and result semantics |
| Data model and numeric semantics | Yes | Defines values, nullish rules, truthiness, numeric behavior, string behavior, and Unicode behavior |
| Standard library specification | Yes | Defines built-in functions, built-in pipes, profiles, and their semantics |
| Host interoperability and extension model | Yes | Defines host value mapping, capability declarations, extension boundaries, and extension guidance |
| Error model and diagnostics | Yes | Defines parse-time, runtime, capability, and version-related error categories and required conditions |
| Configuration, profiles, and capabilities | Yes | Defines core vs optional features, named profiles, and host capability declarations |
| Versioning and compatibility process | Yes | Defines language versions, compatibility guarantees, deprecation, removal, and host support declarations |
| Conformance suite and harness contract | Yes | Defines the machine-readable test format and how conformance is judged |
| Implementation guidance appendix | Optional | Provides non-normative guidance for implementers and extension authors |

## 7. Writing and Organization Rules

The specification must follow these writing rules:

1. Use normative language consistently: `MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT`, and `MAY`.
2. Clearly separate normative content from informative notes, rationale, examples, and implementation guidance.
3. Mark the normative status of each artifact clearly.
4. Define each semantic rule in one canonical place only.
5. Use one canonical terminology set and maintain it in the glossary artifact.
6. Avoid implementation-language terms unless explicitly marked as a mapping example.
7. Prefer language-neutral terms such as number, string, array, object, null, absent, profile, capability, and extension.
8. Choose one formal grammar notation and use it consistently throughout the grammar artifact.
9. Provide a machine-readable grammar companion if practical.
10. Every example that demonstrates a result must either show the expected value or reference a matching conformance case.
11. Every configuration, profile, capability, or version dependency must be named explicitly.
12. Every undefined, implementation-defined, or host-defined area must be labeled as such.
13. Error categories and feature names must be stable, even if human-readable error messages vary between implementations.
14. Avoid restating the same rule in multiple artifacts; use cross-references instead.

## 8. Required Coverage Areas

The completed specification set must cover all currently shipped language behavior that affects observable results.

At minimum, it must cover the following areas.

### 8.1 Lexical Layer

The specification must define:

- Valid tokens and token boundaries
- Numeric literal forms
- String literal forms and escapes
- Identifier rules
- Reserved keywords
- Scope variable syntax
- Whitespace handling
- Unsupported constructs such as comments and multi-statement programs

### 8.2 Grammar Layer

The specification must define:

- The top-level expression model
- Operator precedence and associativity
- Parenthesized expressions
- Arrays and objects
- Property access and index access
- Optional chaining and nullish operators
- Function calls
- Pipe syntax, including parameterized pipes such as `|window(n):` and `|chunk(n):`
- Any syntactic restrictions on where pipes may appear
- The formal grammar notation used by the specification
- Any ambiguity resolution rules that affect parsing

### 8.3 Core Semantic Model

The specification must define:

- Evaluation order
- Short-circuit behavior
- Truthiness rules
- Nullish behavior and absent-variable behavior
- Strict vs optional access behavior
- Equality and comparison semantics
- Pipe evaluation model, including system variables such as `$item`, `$index`, `$acc`, `$last`, `$window`, and `$chunk`
- Which behaviors are deterministic and must match exactly across implementations
- Which operational behaviors are host-defined rather than language-defined
- The absence of hidden dependence on host locale, host clock, randomness, or other external state unless explicitly standardized

### 8.4 Data Model, Numeric Semantics, and String Semantics

The specification must define:

- The UExL value model
- Null, absent, boolean, number, string, array, and object semantics
- Host-to-UExL and UExL-to-host value mapping expectations at the conceptual level
- The numeric representation model required for conformance
- Integer and floating-point behavior, if both are exposed semantically
- Division, modulo, rounding, overflow, underflow, and exceptional numeric cases
- IEEE-754 `NaN` and `Inf` behavior, when enabled
- String indexing and slicing behavior
- The Unicode level used by default and any named alternate views or profiles
- Locale-independent string behavior unless a locale-aware mode is explicitly standardized

Numeric behavior is a portability-critical area and must not be left implicit.

### 8.5 Standard Library

The specification must define:

- The mandatory core built-in functions
- The mandatory core built-in pipe types
- Any optional built-ins or built-in profiles
- The exact name, arity, accepted input types, result type, nullish behavior, and error behavior for each built-in
- Any deterministic ordering guarantees, short-circuit guarantees, or stability guarantees
- Any profile or capability gating that affects built-in availability
- The rules for deprecating, replacing, or removing built-ins across language versions

For cross-platform portability, the specification must define a mandatory core standard library.

Optional extensions must not change the semantics of the mandatory core built-ins.

### 8.6 Extension Model and Host Interoperability

The specification must define:

- The boundary between the core language, named profiles, and host-specific extensions
- Whether and how hosts may expose user-defined functions, custom pipes, host objects, or other injected constructs
- The conceptual contract for injected constructs: argument passing, result passing, nullish handling, and error propagation
- Whether injected constructs may access system variables or host capabilities
- The naming and collision model for extensions, if standardized
- The minimum host capability declaration model, including supported language version, supported profiles, and supported named extensions
- The host value categories that may be passed into a UExL context and the categories that must be rejected or transformed
- Which extension behaviors may be standardized later through named extension profiles

User-defined functions, custom pipes, and other injected runtime constructs are not covered by the core portability guarantee unless they are separately standardized by the specification or by a named extension profile.

### 8.7 Error Model and Diagnostics

The specification must define:

- Parse-time error categories
- Runtime error categories
- Capability and profile mismatch categories
- Version mismatch categories
- The required conditions under which each error category must be produced
- Which aspects of an error must be stable for conformance, such as category or code
- Which aspects of an error may vary by implementation, such as human-readable message wording
- Whether source locations, spans, or structured error metadata are required or optional

### 8.8 Versioning, Compatibility, Profiles, and Capability Declarations

The specification must define:

- The language versioning scheme
- The compatibility policy between versions
- The difference between language versions, named profiles, and host extensions
- The feature lifecycle model, including core, optional, experimental, deprecated, and removed states
- How a host declares the UExL version or version range that it supports
- How a host declares supported profiles and capabilities
- What happens when an expression targets unsupported features, profiles, or versions
- Whether and how expressions may declare or imply a target version or required capability set

Portability guarantees apply only between implementations that support the same language version and compatible profiles.

### 8.9 Conformance Suite and Harness Contract

The specification effort must produce a conformance suite in a machine-readable format.

Each test case must include, at minimum:

- Stable test ID
- Category
- Target language version
- Required profile or capabilities, if any
- Expression
- Input context
- Expected result or expected error
- Notes or rationale for unusual cases

The conformance framework must also define:

- How exact result matching works for primitive values, arrays, objects, nullish values, and special numeric values
- How floating-point-sensitive cases are judged
- How error matching works, including exact fields vs implementation-specific text
- How optional-profile or capability-gated tests are selected or skipped
- How unsupported-version cases are represented

The suite must cover:

- Successful evaluation cases
- Parse failures
- Runtime failures
- Version mismatch cases
- Profile and capability mismatch cases
- Boundary conditions
- Cross-feature interaction cases

The suite must not depend on Go bytecode, Go AST internals, or Go-specific error messages.

## 9. Portability Requirements

The specification must be written so that it can be implemented consistently in languages with different runtime models.

To achieve that, the specification must:

1. Avoid dependence on Go compiler, bytecode, or VM implementation details.
2. Define observable behavior only, unless a lower-level mechanism is explicitly standardized.
3. Define numeric behavior clearly enough that JavaScript, Rust, Dart, and Go implementations can align.
4. Define string and Unicode behavior clearly enough that differences in host string representation do not leak into the language.
5. Distinguish the core language from optional host extensions.
6. State which behaviors are deterministic and must match exactly across ports.
7. State which areas are implementation-defined or host-defined.
8. Require that all expressions using only the core language plus the mandatory core built-ins produce equivalent observable results across conforming implementations.
9. State explicitly that portability guarantees do not automatically extend to user-defined functions, custom pipes, host objects, or other injected constructs.
10. Permit non-normative guidance for authors who want to design portable extensions, while making clear that such guidance is not itself a conformance guarantee.
11. State explicitly that portability guarantees apply only between implementations that support the same language version and compatible profiles.
12. Require every host to declare the language version, profiles, and capabilities relevant to conformance.
13. Define the required behavior when an expression targets features that are not supported by the host version, profile, or capability set.

## 10. Relationship to Existing Documentation and Implementations

The books in `docs/book` and `docs/mastering-uexl` are valuable inputs to the specification effort, but they are not themselves the normative specification.

The specification effort should use the following sources:

1. Existing user documentation for terminology, examples, and intent
2. Existing tests for current observable behavior
3. Existing implementation code for clarification where behavior is ambiguous
4. Benchmarks only for non-normative performance guidance

Where the books, tests, and implementation disagree, the discrepancy must be resolved explicitly and recorded before the relevant section of the specification is marked complete.

## 11. Acceptance Criteria for Specification v1

Version 1 of the UExL specification is complete only when all of the following are true:

1. All currently shipped language features with observable behavior are covered normatively.
2. All major edge cases are specified explicitly.
3. Numeric semantics, string semantics, and Unicode behavior are defined strongly enough for independent non-Go implementations to align.
4. The boundary between the core language, named profiles, and host extensions is defined clearly.
5. The conformance suite exists, is machine-readable, and can be run independently of the Go implementation internals.
6. The conformance framework defines how results, errors, versions, profiles, and capabilities are compared.
7. At least one non-Go prototype implementation can use the specification and conformance suite as its primary reference.
8. The maintainers can identify which behaviors are guaranteed by the language and which are implementation details or host-defined policies.
9. The specification defines how language versions are declared, how compatibility is evaluated, and what happens when a newer expression is presented to an older host.
10. The books can safely reference the specification for authoritative behavior.

## 12. Recommended Work Plan

The specification effort should be executed in phases.

### Phase 0: Roadmap and Authority Setup

- Create the specification roadmap artifact.
- Record the authority hierarchy and dispute-resolution rule.
- Define section ownership and review responsibilities.

### Phase 1: Inventory

- Catalog all observable language features.
- Inventory every built-in function and built-in pipe.
- Inventory every configuration-sensitive, profile-sensitive, and capability-sensitive behavior.
- Inventory every host extension point.
- Identify mismatches between implementation, tests, and documentation.

### Phase 2: Normalize Current Behavior

- Resolve disagreements in docs and tests.
- Convert ambiguous behaviors into explicit rules.
- Mark implementation details that must not leak into the specification.
- Decide which behaviors are core, optional, experimental, deprecated, or host-defined.

### Phase 3: Author Normative Artifacts

- Write the lexical, grammar, semantic, data model, standard library, host interoperability, error, versioning, and profile artifacts.
- Mark each artifact or section as draft, reviewed, or approved.
- Add cross-references between related rules rather than duplicating them.

### Phase 4: Build the Conformance Suite

- Translate approved rules into machine-readable test cases.
- Add edge-case coverage for every operator, access mode, function, and pipe.
- Add version, profile, and capability mismatch coverage.
- Define the result-matching and error-matching rules used by the harness.

### Phase 5: Port Validation

- Use the specification to guide at least one alternate implementation.
- Record ambiguities discovered during the port.
- Tighten the specification until the port no longer relies on guesswork.

### Phase 6: Stabilization and Publication

- Freeze the version 1.0 scope.
- Publish the approved specification set and conformance suite together.
- Align the Go implementation and books with the approved specification.
- Define the errata and post-1.0 change process.

## 13. Recommended Repository Layout

The `docs/specs` directory should evolve into a structured specification set rather than a single monolithic document.

A recommended layout is:

```text
docs/specs/
  README.md
  roadmap.md
  language-overview.md
  terminology-glossary.md
  lexical-spec.md
  grammar-spec.md
  semantic-spec.md
  data-model-and-numeric-semantics.md
  standard-library.md
  host-interoperability.md
  error-model.md
  profiles-and-capabilities.md
  versioning-and-compatibility.md
  conformance-suite.md
  change-process.md
```

This file is the entry requirements document for that effort.

## 14. Immediate Next Artifact

The next artifact to create after this requirements document should be `docs/specs/roadmap.md`.

That roadmap should list:

- The artifacts to be authored
- The source files that currently describe each area
- Known conflicts or ambiguities
- Ownership and review status
- Open decisions that block precise specification text

That roadmap should be used to track progress until the full UExL specification set is complete.