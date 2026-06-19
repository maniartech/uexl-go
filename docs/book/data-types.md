# Data Types

Every expression language needs a solid foundation of data types, and UExL is no exception. In this chapter, you'll discover the core types that power UExL, how to use them, and practical examples for each.

## Why Data Types Matter
Data types define what kind of information you can work with—numbers, text, logical values, collections, and more. Understanding these types helps you write correct, expressive, and efficient expressions.

## Numbers
Numbers in UExL can be integers or floating-point values. Use them for calculations, comparisons, and more.
- **Examples:**
  - `1` (integer)
  - `-42` (negative integer)
  - `3.14` (floating-point)
  - `1e3` (scientific notation, equals 1000)
- **Edge Case:** Leading zeros are not allowed: `01` is invalid.

### Special numeric values (enabled by default)
UExL supports IEEE-754 special numeric values as literals (enabled by default, configurable):
- `NaN` (Not-a-Number)
- `Inf` (positive infinity)
- `-Inf` (negative infinity, via unary minus)

Notes:
- Enabled by default; can be disabled via parser options (EnableIeeeSpecials). When disabled, `NaN` and `Inf` are treated as identifiers.
- Only `NaN` and `Inf` are recognized; there is no `Infinity` literal and no unary plus form.
- Equality and comparisons follow IEEE-754 rules: for example, `NaN != NaN` is true, and any comparison with `NaN` is false.
- Runtime operator behavior (arithmetic, comparisons, bitwise, etc.) with `NaN`/`Inf` is specified in `vm/ieee754-semantics.md`. Division by zero remains an error by design.

## Strings
Strings are sequences of characters, enclosed in single or double quotes. Use them for text, messages, and keys.
- **Examples:**
  - `"hello"`
  - `'world'`
  - `"He said, 'hi'"`
- **Edge Case:** Escape sequences (like `\"`, `\\`) are not supported; use matching quotes to include quotes inside strings.
- **Indexing:** Use square brackets to access a character by zero-based **byte** index: `"abc"[1] // "b"`. Out-of-bounds returns `null`. Strings are immutable.
- **Slicing:** `s[i:j]` returns the substring from byte `i` (inclusive) to byte `j` (exclusive).
- **`len(s)`** returns the number of **UTF-8 bytes**, not runes or visible characters. For a pure-ASCII string the three counts are identical.
- **Unicode:** For multi-byte code points or composed characters (e.g. emoji, accented letters), byte-level operations may split a codepoint. Use the explicit Unicode functions `runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`, `runes()`, `graphemes()` when character-level or display-level semantics are required. See [Strings and Unicode](strings-unicode.md).

## Booleans
Booleans represent logical truth—`true` or `false`. Use them in conditions, filters, and logical operations.
- **Examples:**
  - `true`
  - `false`
- **Usage:**
  - `x > 10 && y < 20` (returns a boolean)

## Null
`null` means "no value" or "missing data." Use it to indicate absence or undefined values.
- **Example:**
  - `null`
- **Usage:**
  - `user.middleName` might be `null` if not set.

## Arrays
Arrays are ordered lists of values, enclosed in square brackets. They can hold any type, including other arrays and objects.
- **Examples:**
  - `[1, 2, 3]`
  - `["a", true, null, [1, 2]]`
- **Access:**
  - `arr[0]` (first element)
- **Edge Case:** Arrays can be empty: `[]`.

## Objects
Objects are collections of key-value pairs, enclosed in curly braces. Keys are strings (quoted or unquoted if valid identifiers), and values can be any type.
- **Examples:**
  - `{"name": "UExL", "version": 1.0}`
  - `{id: 123, values: [1,2,3]}`
- **Access:**
  - `obj.key` or `obj["key"]`
- **Edge Case:** Keys must be unique within an object.

## Datetime and Duration

UExL's type system includes two temporal value types:

- **`datetime`** — an instant in time, represented as milliseconds since the Unix epoch (`1970-01-01T00:00:00.000Z`, UTC). It is zoneless: it identifies an absolute instant, not a wall-clock reading in a particular zone.
- **`duration`** — an exact span of time in milliseconds, which may be negative (for example, the result of subtracting a later instant from an earlier one).

**Datetime literals** are written with a `d` prefix followed by an ISO 8601 / RFC 3339 string. A missing time-of-day defaults to midnight UTC; an offset is used to compute the UTC instant and then discarded:
- `d"2024-12-01"` — date only → `2024-12-01T00:00:00.000Z`
- `d"2024-12-01T10:30:00Z"` — explicit UTC instant
- `d"2024-12-01T10:30:00+05:30"` — the **same** instant as `d"2024-12-01T05:00:00Z"`
- An invalid literal such as `d"2024-13-01"` or `d"2024-02-30"` is a **parse-time error**.

**Duration literals** are a number with a fixed-unit suffix and no whitespace, where `ms`=millisecond, `s`=second, `m`=minute, `h`=hour, `d`=day, `w`=week. There is deliberately no month or year suffix (those are not fixed-length):
- `7d` (7 days), `1.5h` (90 minutes), `30ms`, `45s`, `2w`
- Fractional magnitudes are allowed and truncate to whole milliseconds (`1.5h` → `5400000` ms).

**Truthiness:** like every other type, a temporal value is falsy when it is the *zero* of its type — a `datetime` at the epoch (`0` ms) and a `duration` of `0` ms are falsy; every other instant or span is truthy. Temporal values are never nullish (so `??` does not treat them as missing), exactly as with `0`, `""`, and `false`.

> **Status:** datetime/duration are fully implemented and conformance-verified — the types, literals, truthiness, the temporal **operators** (`date − date` → duration, `date ± duration`, comparisons), the datetime **function library** (`parseDate`/`formatDate`/`addMonths`/`datePart`/… attached via `uexl.WithDatetime()`), and clock-injected `now()`/`today()`. See the [DateTime Specification](../specs/datetime-spec.md).

## Putting It All Together: Examples
```
42
3.14
"hello"
true
null
[1, 2, 3]
{"name": "UExL", "features": ["pipes", "functions"]}
```

Arrays and objects can be nested and can contain any supported data type. Mastering data types is essential for writing expressive UExL code. In the next chapter, we'll explore how to use variables and identifiers to work with your data.