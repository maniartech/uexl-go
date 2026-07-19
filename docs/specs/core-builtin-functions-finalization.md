# UExL Core Built-in Functions Finalization

Status: Draft for finalization
Audience: Maintainers, language designers, implementers of UExL ports
Scope: Canonical built-in function names and semantics for the mandatory core standard library. Built-in pipes are out of scope here.

## 1. Purpose

This document freezes the core built-in function surface that a conforming UExL host MUST provide out of the box.

The goal is to keep the standard library small, portable, and easy to implement while still being self-sufficient for common arithmetic, string processing, Unicode-safe text handling, type inspection, and non-mutating object updates.

This document exists because the current runtime, README, and books do not yet agree on the core function set. After approval, this document should become the canonical source of truth for built-in function names, arity, accepted argument types, result types, and error behavior.

## 2. Final decisions

1. `Default()` MUST expose the mandatory core functions listed in Section 4.
2. Books, README examples, tests, and future ports MUST use the canonical names defined here.
3. The mandatory core standard library MUST NOT rely on host-provided helper functions for basic expression authoring.
4. The core library MUST stay pure by default. Built-ins MUST NOT mutate caller-owned arrays, objects, globals, or context values.
5. `set(obj, key, value)` is part of the mandatory core, but it MUST return a new object and MUST NOT mutate `obj`.
6. The mandatory core conversion family is `str(value)`, `num(value)`, and `bool(value)`.
7. `str(value)` is the canonical string conversion function name. `toString(value)` is not part of the mandatory core.
8. Dedicated typecheck helpers such as `isString(value)` and `isNumber(value)` are not part of the mandatory core because `typeof(value)` already provides the canonical, portable inspection path.
9. `concat(...)` is not part of the mandatory core because string concatenation already exists via `+` and `join(...)`.
10. `sum`, `average`, `count`, and `count_if` are not mandatory core built-ins. They MAY exist as host extensions or a future optional profile because equivalent behavior already exists via pipes plus `len(...)`.
11. `isTruthy` and `isFalsy` are not mandatory core built-ins. Their behavior is already expressible with `!!x` and `!x`.
12. No alias names are part of the mandatory core. Hosts MAY provide aliases, but aliases are extensions and are not portability-safe.

## 3. Design constraints

The mandatory core function set MUST satisfy all of the following:

- Small enough to port without a large runtime dependency surface.
- Complete enough that a host can evaluate ordinary business expressions without injecting its own helpers first.
- Explicit about string level, nullish handling, and mutation behavior.
- Consistent in naming and function pairing.
- Portable across implementations without relying on Go-specific formatting or data structures.

## 3.1 Current Go runtime snapshot vs final target

The current Go runtime already ships the following subset in `vm.Builtins`:

- `len`, `substr`, `contains`, `set`, `str`
- `runeLen`, `runeSubstr`
- `graphemeLen`, `graphemeSubstr`
- `runes`, `graphemes`, `bytes`
- `join`

This document intentionally defines a larger final mandatory core. Anything listed in Section 4 but missing from the runtime should be treated as required completion work, not as an optional enhancement.

## 4. Mandatory core function set

### 4.1 Overview table

| Category | Function | Signature | Return | Required behavior summary |
|----------|----------|-----------|--------|---------------------------|
| General | `len` | `len(value)` | `number` | Byte length for strings, element count for arrays |
| General | `str` | `str(value)` | `string` | Explicit scalar-to-string conversion |
| General | `num` | `num(value)` | `number` | Explicit conversion to number |
| General | `bool` | `bool(value)` | `boolean` | Explicit conversion using truthiness rules |
| General | `typeof` | `typeof(value)` | `string` | Returns canonical UExL type name |
| Number | `abs` | `abs(x)` | `number` | Absolute value |
| Number | `min` | `min(x1, x2, ...)` or `min(array)` | `number` | Minimum of non-empty numeric inputs |
| Number | `max` | `max(x1, x2, ...)` or `max(array)` | `number` | Maximum of non-empty numeric inputs |
| Number | `floor` | `floor(x)` | `number` | Round toward negative infinity |
| Number | `ceil` | `ceil(x)` | `number` | Round toward positive infinity |
| Number | `round` | `round(x)` | `number` | Round to nearest integer, ties away from zero |
| Number | `isNaN` | `isNaN(x)` | `boolean` | True only for numeric `NaN` |
| Number | `isFinite` | `isFinite(x)` | `boolean` | True only for finite numbers |
| String | `substr` | `substr(s, start, length)` | `string` | Byte-level substring |
| String | `contains` | `contains(s, needle)` | `boolean` | Substring containment |
| String | `startsWith` | `startsWith(s, prefix)` | `boolean` | Prefix test |
| String | `endsWith` | `endsWith(s, suffix)` | `boolean` | Suffix test |
| String | `split` | `split(s, sep)` | `array<string>` | Split into string elements |
| String | `join` | `join(arr)` or `join(arr, sep)` | `string` | Join string elements |
| String | `upper` | `upper(s)` | `string` | Locale-independent uppercasing |
| String | `lower` | `lower(s)` | `string` | Locale-independent lowercasing |
| String | `trim` | `trim(s)` | `string` | Trim leading and trailing Unicode whitespace |
| String | `isAlpha` | `isAlpha(s)` | `boolean` | True when `s` is a non-empty alphabetic string or grapheme |
| Unicode | `runeLen` | `runeLen(s)` | `number` | Code point count |
| Unicode | `runeSubstr` | `runeSubstr(s, start, length)` | `string` | Rune-level substring |
| Unicode | `graphemeLen` | `graphemeLen(s)` | `number` | Grapheme cluster count |
| Unicode | `graphemeSubstr` | `graphemeSubstr(s, start, length)` | `string` | Grapheme-level substring |
| Unicode | `runes` | `runes(s)` | `array<string>` | Explode into single-code-point strings |
| Unicode | `graphemes` | `graphemes(s)` | `array<string>` | Explode into grapheme-cluster strings |
| Unicode | `bytes` | `bytes(s)` | `array<number>` | Explode into UTF-8 byte values |
| Object update | `set` | `set(obj, key, value)` | `object` | Return a copy of `obj` with `key` assigned |

### 4.2 Canonical type names for `typeof(...)`

`typeof(value)` MUST return exactly one of the following strings:

- `"null"`
- `"boolean"`
- `"number"`
- `"string"`
- `"array"`
- `"object"`

Absent context variables are already treated as `null` by the evaluator before function dispatch, so `typeof(missingVar)` MUST return `"null"` if the expression is otherwise valid.

## 5. Detailed function rules

### 5.1 General functions

#### `len(value)`

- Arity: exactly 1
- Accepted input types:
  - `string` -> UTF-8 byte length
  - `array` -> number of top-level elements
- Returns: `number`
- Errors:
  - wrong arity
  - unsupported type
- Notes:
  - `len` is byte-based for strings by design.
  - `len` does not count object keys.

#### `str(value)`

- Arity: exactly 1
- Accepted input types:
  - `null`
  - `boolean`
  - `number`
  - `string`
- Returns: `string`
- Errors:
  - wrong arity
  - array or object input
- Required formatting:
  - `null` -> `"null"`
  - `true` -> `"true"`
  - `false` -> `"false"`
  - string -> identity
  - number -> shortest round-trip decimal form
  - integer-valued numbers MUST omit a trailing `.0`
  - `NaN` -> `"NaN"`
  - positive infinity -> `"+Inf"`
  - negative infinity -> `"-Inf"`

#### `num(value)`

- Arity: exactly 1
- Accepted input types:
  - `number`
  - `boolean`
  - `string`
- Returns: `number`
- Errors:
  - wrong arity
  - `null`, array, or object input
  - invalid numeric string
- Conversion rules:
  - number -> identity
  - `true` -> `1`
  - `false` -> `0`
  - string -> parse after trimming leading and trailing Unicode whitespace
- Accepted string forms:
  - integer decimal forms such as `"42"`
  - decimal floating forms such as `"3.14"`
  - scientific notation such as `"1e6"`
  - `"NaN"`, `"Inf"`, and `"-Inf"`
- Rejected string forms:
  - empty strings after trimming
  - non-numeric content such as `"abc"`
  - implementation-defined formats not accepted by UExL numeric literals

#### `bool(value)`

- Arity: exactly 1
- Accepted input types: any UExL value
- Returns: `boolean`
- Errors: wrong arity only
- Semantics:
  - explicit truthiness conversion
  - MUST be exactly equivalent to `!!value`
- Truthiness conversion rules:
  - `null` -> `false`
  - `false` -> `false`
  - `true` -> `true`
  - numeric zero -> `false`
  - numeric non-zero, `NaN`, and infinities -> `true`
  - empty string -> `false`
  - non-empty string -> `true`
  - arrays and objects -> `true`

#### `typeof(value)`

- Arity: exactly 1
- Accepted input types: any UExL value
- Returns: one of the canonical type strings listed in Section 4.2
- Errors: wrong arity only

### 5.2 Numeric functions

#### Common numeric argument rules

Unless a function says otherwise:

- numeric arguments MUST be numbers
- wrong arity is an error
- non-number input is a type error
- `NaN` and infinities follow the underlying numeric semantics of the operation unless explicitly rejected

#### `abs(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `number`

#### `min(...)` and `max(...)`

- Supported forms:
  - `min(x1, x2, ...)`
  - `min(array)`
  - `max(x1, x2, ...)`
  - `max(array)`
- Accepted input:
  - variadic numeric arguments, at least 1
  - or one non-empty array containing only numbers
- Returns: `number`
- Errors:
  - zero arguments
  - empty array
  - mixed or non-numeric values

#### `floor(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `number`
- Semantics: round toward negative infinity

#### `ceil(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `number`
- Semantics: round toward positive infinity

#### `round(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `number`
- Semantics: round to the nearest integer, ties away from zero

#### `isNaN(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `boolean`
- Errors:
  - wrong arity
  - non-number input

#### `isFinite(x)`

- Arity: exactly 1
- Accepted input: `number`
- Returns: `boolean`
- Errors:
  - wrong arity
  - non-number input

### 5.3 String functions

#### Common string rules

- All string functions are locale-independent unless explicitly stated otherwise.
- String comparison and search behavior MUST be deterministic and MUST NOT depend on host locale.

#### `substr(s, start, length)`

- Arity: exactly 3
- Accepted input:
  - `s`: `string`
  - `start`: integer-valued `number`
  - `length`: integer-valued `number`
- Returns: `string`
- Semantics:
  - byte-level substring
  - `start < 0` -> error
  - `length < 0` -> error
  - `start > len(s)` -> error
  - end is clamped to the end of the string

#### `contains(s, needle)`

- Arity: exactly 2
- Accepted input: two strings
- Returns: `boolean`

#### `startsWith(s, prefix)`

- Arity: exactly 2
- Accepted input: two strings
- Returns: `boolean`

#### `endsWith(s, suffix)`

- Arity: exactly 2
- Accepted input: two strings
- Returns: `boolean`

#### `split(s, sep)`

- Arity: exactly 2
- Accepted input:
  - `s`: `string`
  - `sep`: non-empty `string`
- Returns: `array<string>`
- Errors:
  - wrong arity
  - non-string input
  - empty separator
- Notes:
  - empty-separator splitting is intentionally excluded from the core to avoid implicit byte/rune/grapheme ambiguity
  - use `runes(s)` or `graphemes(s)` when decomposition by Unicode level is required

#### `join(arr)` and `join(arr, sep)`

- Arity: 1 or 2
- Accepted input:
  - `arr`: array whose elements are all strings
  - `sep`: optional string, default `""`
- Returns: `string`
- Errors:
  - non-array input
  - any non-string element
  - non-string separator

#### `upper(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `string`
- Semantics: locale-independent Unicode uppercasing

#### `lower(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `string`
- Semantics: locale-independent Unicode lowercasing

#### `trim(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `string`
- Semantics: trim leading and trailing Unicode whitespace

#### `isAlpha(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `boolean`
- Errors:
  - wrong arity
  - non-string input
- Semantics:
  - locale-independent Unicode alphabetic classification
  - returns `true` only when `s` is non-empty and all code points in `s` are alphabetic content
  - to preserve grapheme-safe workflows, Unicode combining marks are permitted only when attached to an alphabetic base character
  - equivalently: all code points in `s` MUST be in Unicode Letter or Mark categories, and at least one code point MUST be in a Unicode Letter category
- Examples:
  - `isAlpha("A")` -> `true`
  - `isAlpha("é")` -> `true`
  - `isAlpha("e\u0301")` -> `true`
  - `isAlpha("7")` -> `false`
  - `isAlpha("A7")` -> `false`
  - `isAlpha("")` -> `false`

This function exists because expressions such as `graphemes(lower(userName)) |filter: isAlpha($item)` need a portable, Unicode-aware way to keep letters without forcing every host to inject a custom helper.

### 5.4 Explicit Unicode-level functions

These functions exist specifically so the core language can remain explicit about byte, rune, and grapheme behavior.

#### `runeLen(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `number`
- Semantics: number of Unicode code points in `s`

#### `runeSubstr(s, start, length)`

- Arity: exactly 3
- Accepted input:
  - `s`: `string`
  - `start`: integer-valued `number`
  - `length`: integer-valued `number`
- Returns: `string`
- Semantics:
  - rune-level substring
  - `length < 0` -> error
  - `start < 0` is clamped to `0`
  - `start >= runeLen(s)` returns `""`
  - end is clamped to available runes

#### `graphemeLen(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `number`
- Semantics: number of grapheme clusters in `s`

#### `graphemeSubstr(s, start, length)`

- Arity: exactly 3
- Accepted input:
  - `s`: `string`
  - `start`: integer-valued `number`
  - `length`: integer-valued `number`
- Returns: `string`
- Semantics:
  - grapheme-level substring
  - `length < 0` -> error
  - `start < 0` is clamped to `0`
  - `start >= graphemeLen(s)` returns `""`
  - end is clamped to available graphemes

#### `runes(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `array<string>`
- Semantics: one string element per Unicode code point

#### `graphemes(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `array<string>`
- Semantics: one string element per grapheme cluster

#### `bytes(s)`

- Arity: exactly 1
- Accepted input: `string`
- Returns: `array<number>`
- Semantics: one numeric element per UTF-8 byte, each element in the inclusive range `0..255`

### 5.5 Non-mutating object update

#### `set(obj, key, value)`

- Arity: exactly 3
- Accepted input:
  - `obj`: object
  - `key`: string or integer-valued number
  - `value`: any UExL value
- Returns: `object`
- Semantics:
  - returns a new object containing all properties of `obj`
  - the returned object MUST contain `key` assigned to `value`
  - when `key` is numeric, it is stringified without a trailing `.0`
  - the input `obj` MUST NOT be mutated
- Errors:
  - wrong arity
  - non-object first argument
  - non-string and non-integer key

## 6. Pairing and naming rules

To keep the core library coherent, the following pairings are intentional and MUST remain aligned:

- `contains` pairs with `startsWith` and `endsWith`
- `split` pairs with `join`
- `str`, `num`, and `bool` are the explicit scalar conversion family
- `upper` pairs with `lower`
- `isAlpha` is the minimal character-classification helper required for portable text filtering; additional helpers can be added later as an optional text-classification profile
- `runeLen` pairs with `runeSubstr`
- `graphemeLen` pairs with `graphemeSubstr`
- `runes`, `graphemes`, and `bytes` are the three explicit explode functions for the three string levels
- `min` pairs with `max`
- `floor` pairs with `ceil`, with `round` completing the rounding set

Canonical naming rules:

- Short established names are preferred when already idiomatic and unambiguous: `len`, `str`, `min`, `max`
- Unicode-level names use lowerCamel with explicit units: `runeLen`, `graphemeSubstr`
- Core MUST avoid duplicate aliases such as `str` plus `toString`, or `substr` plus `substring`

## 7. Explicitly not in the mandatory core

The following are intentionally excluded from the mandatory core either because they are redundant, host-specific, or better handled by pipes and future profiles:

- `toString(...)` alias of `str(...)`
- `toNumber(...)` alias of `num(...)`
- `toBoolean(...)` alias of `bool(...)`
- `concat(...)`
- `sum(...)`
- `average(...)`
- `count(...)`
- `count_if(...)`
- dedicated typecheck helpers such as `isString(...)`, `isNumber(...)`, `isBoolean(...)`, `isArray(...)`, `isObject(...)`, and `isNull(...)`
- `isNullish(...)`
- `isTruthy(...)`
- `isFalsy(...)`
- additional character-classification helpers such as `isDigit(...)`, `isAlnum(...)`, and `isSpace(...)`
- date/time helpers; see `datetime-handling.md`
- regex helpers; see `regex-handling.md`

Notes:

- `sum`, `average`, `count`, and `count_if` remain good candidates for an optional convenience profile, especially for Excel-style hosts.
- dedicated typecheck helpers are intentionally excluded because `typeof(value)` already provides a portable alternative without exploding the core surface area.
- `isNullish(x)` is not essential in the mandatory core because `x == null` already expresses the same check under current absent-variable semantics.
- `isTruthy(x)` and `isFalsy(x)` are redundant with `!!x` and `!x`.
- date/time and regex are intentionally excluded from the current core not because they are unimportant, but because they require separate future profile/capability documents to preserve portability.

## 8. Required alignment work after approval

After this document is approved, the following changes should be made to the implementation and docs:

1. Add the missing mandatory core implementations to the Go runtime.
2. Change `set(...)` to return a copy instead of mutating the input object.
3. Update README and books to use canonical core names only.
4. Remove or relabel examples that currently assume non-core helpers such as `concat`, `toString`, `toNumber`, `toBoolean`, `sum`, `average`, and `count_if` unless those helpers are promoted into a named optional profile.
5. Add conformance tests for every function covering:
   - correct result
   - wrong arity
   - wrong type
   - null input where relevant
   - integer-index validation where relevant
   - purity and non-mutation where relevant

## 9. Open implementation note

The current Go runtime already ships a partial subset of this document. That runtime subset should be treated as incomplete, not authoritative, until it matches this finalization document.