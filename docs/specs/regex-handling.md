# UExL Regular Expression Handling

Status: Draft for future version
Audience: Maintainers, language designers, implementers of UExL ports
Scope: Proposed future handling for regular-expression support, capability declaration, portability boundaries, and possible standardization strategies. This document does not make regex part of the current mandatory core standard library.

## 1. Purpose

Regular expressions are a common requirement for production expression hosts: validation, extraction, routing rules, text normalization, and policy matching all benefit from regex support.

UExL does not currently include regex helpers in the mandatory core built-in set. This document records the future design space so regex support can be added without weakening UExL's portability story.

## 2. Current status

- Regex helpers are not part of the current mandatory core built-in set.
- Hosts MAY inject custom regex helpers today, but those are host extensions and are not portability-safe.
- Regex support should not enter the mandatory core unless the dialect and behavior are clearly standardized.

## 3. Why regex is harder than it looks

Regex support seems smaller than date/time support, but portability is harder.

Different hosts use materially different regex engines and dialects:

- Go commonly uses RE2-style semantics
- JavaScript uses ECMAScript regular expressions
- Python, Java, .NET, and Rust all differ in supported features and edge cases

The differences are not cosmetic. They affect observable behavior:

- lookbehind support
- backreferences
- named capture groups
- Unicode character classes
- greedy vs lazy details in corner cases
- replacement syntax
- catastrophic backtracking risk in some engines but not others

Because UExL is aiming for portability, regex support must be capability-aware from the start.

## 4. Core recommendation

Regex SHOULD NOT be part of the mandatory core built-in function set in its first standardized form.

Instead, future regex support SHOULD use one of these models:

1. a named regex profile with a fixed UExL-defined dialect
2. a capability-gated regex extension where hosts must declare the supported regex engine or regex feature subset

The first model gives stronger portability. The second is easier to ship but weaker as a cross-host contract.

## 5. Recommended standardization strategy

### 5.1 Preferred direction: fixed-profile dialect

The best long-term path is a named profile such as `regex-core` that defines a fixed accepted syntax and required behavior.

That profile SHOULD:

- choose one constrained regex feature set
- define exactly what syntax is valid
- define matching behavior and flags
- define Unicode behavior
- forbid host-specific extensions that change semantics within the profile

Hosts that cannot implement the chosen dialect faithfully MUST NOT claim support for that profile.

### 5.2 Acceptable interim direction: capability-gated extension

If a fixed profile is too much for the first pass, regex MAY be introduced as a named capability-gated extension.

In that model, hosts MUST declare at least:

- whether regex is supported at all
- which dialect or engine family is used
- whether Unicode classes are supported
- whether capture groups are supported
- whether replacement helpers are supported

This is less portable, but still honest.

## 6. Recommended first-pass surface

If regex support is added in a future version, the initial surface SHOULD be deliberately small.

Recommended first-pass helpers:

- `matches(text, pattern)`
  - boolean full or search match behavior must be defined explicitly
- `find(text, pattern)`
  - returns the first match or `null`
- `findIndex(text, pattern)`
  - returns the first match start index or `-1`

These are safer to standardize than full replacement and group APIs.

Helpers that SHOULD be deferred initially:

- `replaceRegex(...)`
- `findAll(...)`
- capture-group extraction APIs
- named-group maps
- regex splitting

Those require much tighter standardization of replacement syntax and group semantics.

## 7. Pattern representation

For the first version, patterns SHOULD be plain strings passed to regex helpers.

Regex literal syntax SHOULD be deferred initially because it raises extra questions:

- escaping rules
- inline flags
- grammar ambiguity
- serialization concerns across hosts and config formats

Using strings first keeps the language grammar simpler.

## 8. Flags and matching mode

If regex support is standardized, the profile MUST define matching mode explicitly.

At minimum it must answer:

- is `matches(...)` full-string match or substring match?
- how are case-insensitive matches requested?
- are multiline and dotall flags supported?
- are flags embedded in the pattern, passed as a separate argument, or both?

The first version SHOULD prefer a small fixed flag model rather than inheriting every host regex option.

## 9. Unicode and text model

Regex support interacts with UExL's string model, so the profile MUST be explicit about what level regex operates on.

Recommended rule:

- regex operates on strings as strings, not on explicit grapheme arrays or byte arrays
- regex indices, if exposed, MUST define whether they are byte offsets, code-point offsets, or something else

To stay consistent with current UExL string semantics, byte offsets are the most natural first choice, but this must be stated explicitly.

## 10. Performance and safety

Regex support can introduce denial-of-service and predictability issues if the chosen engine allows catastrophic backtracking.

For that reason, a future standardized regex profile SHOULD strongly prefer a linear-time feature set or engine model.

If a host supports a potentially backtracking engine, it SHOULD expose regex as a capability-gated extension rather than claiming a portability-safe profile unless it can guarantee equivalent behavior.

## 11. Error model

Regex support SHOULD define stable error categories for at least:

- invalid pattern syntax
- unsupported regex feature for the claimed profile
- unsupported capability on the current host
- wrong arity
- wrong type

If replacement or capture APIs are ever standardized later, additional categories will be needed.

## 12. Capability declaration guidance

Hosts that support future regex features SHOULD declare support explicitly.

Recommended declaration fields:

- language version
- whether regex is supported
- whether support is via a named profile or host extension
- regex dialect identifier, if extension-based
- supported flags
- Unicode-character-class support
- replacement support, if any

## 13. Non-goals for the first regex version

The first standardized version SHOULD NOT attempt to solve all of these at once:

- host-native regex dialect passthrough
- full replacement syntax portability
- named-group extraction objects
- lookbehind portability
- backreference portability
- regex literals in the grammar

## 14. Recommended rollout order

1. define whether regex is profile-based or capability-gated
2. standardize a minimal search/match surface
3. standardize pattern and flag rules
4. only then consider replacement and capture-group APIs

## 15. Relationship to current core builtin finalization

Regex helpers are intentionally excluded from the current mandatory core built-in set. Their future standardization should happen through the profile or capability model described here, not by expanding the current core ad hoc.