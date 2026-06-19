# Standard Library

Beyond the small set of [core built-in functions](overview.md) that every embedding always has, UExL ships an **attachable standard library**: families of pure helper functions for math, string manipulation, type introspection, collections, conversion, JSON, and date/time. They live in the `builtins/` packages and are opt-in, so a host that wants a minimal surface pays nothing for them.

## Attaching the library

In Go, attach everything at once with `WithStdlib()`, or pick individual families:

```go
import "github.com/maniartech/uexl"

// Everything: math, conversion, introspection, strings, collections, json, datetime.
env := uexl.DefaultWith(uexl.WithStdlib())

// Or à la carte:
env := uexl.DefaultWith(uexl.WithMath(), uexl.WithStrings())
```

| Family | Option | What it adds |
|--------|--------|--------------|
| [Math](#math) | `WithMath()` | `abs` `sign` `round` `floor` `ceil` `trunc` `sqrt` `min` `max` `sum` `avg` `mod` `pow` `clamp` |
| [Conversion](#conversion) | `WithConversion()` | `parseNum` `tryParseNum` `parseBool` `tryParseBool` `formatNum` |
| [Introspection](#introspection) | `WithIntrospection()` | `typeOf` `isNull` `isNumber` `isString` `isBool` `isArray` `isObject` `isDate` `isDuration` `isEmpty` |
| [Strings](#strings) | `WithStrings()` | `upper` `lower` `trim` `trimStart` `trimEnd` `replace` `split` `startsWith` `endsWith` `indexOf` `repeat` `padStart` `padEnd` |
| [Collections](#collections) | `WithCollections()` | `get` `has` `keys` `values` `remove` `merge` |
| [JSON](#json) | `WithJSON()` | `parseJson` `toJson` |
| [Date/Time](#datetime) | `WithDatetime()` | construction, parsing, formatting, calendar & epoch arithmetic (see the [DateTime spec](../../specs/datetime-spec.md)) |

Every standard-library function is **pure**: it never mutates its arguments. The strict/safe convention runs throughout — `parseX` raises an error on bad input, while `tryParseX` returns `null` instead.

---

## Math

Numeric helpers operating on numbers. `min`/`max`/`sum`/`avg` accept **either a single array or several numeric arguments**.

| Function | Returns | Description |
|----------|---------|-------------|
| `abs(x)` | number | Absolute value of `x`. |
| `sign(x)` | number | `1`, `-1`, or `0` by the sign of `x`. |
| `round(x)` | number | Nearest integer; halves round away from zero. **One argument — no decimals parameter** (use `formatNum` to fix decimals). |
| `floor(x)` | number | Largest integer ≤ `x`. |
| `ceil(x)` | number | Smallest integer ≥ `x`. |
| `trunc(x)` | number | Drops the fractional part (rounds toward zero). |
| `sqrt(x)` | number | Non-negative square root; a negative argument yields `NaN`. |
| `min(arr \| a, b, …)` | number | Smallest value; needs at least one number. |
| `max(arr \| a, b, …)` | number | Largest value; needs at least one number. |
| `sum(arr \| a, b, …)` | number | Total; an empty call returns `0`. |
| `avg(arr \| a, b, …)` | number | Arithmetic mean; needs at least one number. |
| `mod(a, b)` | number | Floating-point remainder; result takes the sign of `a`. |
| `pow(base, exp)` | number | `base ** exp`; fractional exponents work (`pow(x, 0.5)` is `sqrt`). |
| `clamp(x, lo, hi)` | number | Constrains `x` to `[lo, hi]`; errors if `lo > hi`. |

```uexl
round(2.5)          // 3      (halves round away from zero)
round(-2.5)         // -3
min([3, 1, 2])      // 1      (array form)
max(3, 1, 2)        // 3      (variadic form)
sum([1, 2, 3])      // 6
avg(2, 4, 6)        // 4
mod(-7, 3)          // -1     (sign of the dividend)
pow(2, 10)          // 1024
clamp(15, 0, 10)    // 10
```

---

## Conversion

String-to-value parsing in **strict** (`parseX`, errors) and **safe** (`tryParseX`, `null`) pairs, plus number formatting. Value-to-string is the core `str`.

| Function | Returns | Description |
|----------|---------|-------------|
| `parseNum(s)` | number | Parses a number (int, decimal, scientific); trims whitespace. Errors on bad input. |
| `tryParseNum(s)` | number \| null | Safe `parseNum`: `null` on a non-string or unparseable text. |
| `parseBool(s)` | boolean | Parses `"true"`/`"false"` case-insensitively (trimmed). Errors otherwise. |
| `tryParseBool(s)` | boolean \| null | Safe `parseBool`: `null` on anything else. |
| `formatNum(x[, decimals])` | string | Renders a number; with `decimals` uses fixed places (round-half-to-even), else the shortest form. |

```uexl
parseNum("  -7.5  ")   // -7.5
tryParseNum("abc")     // null
parseBool("True")      // true
formatNum(3.14159, 2)  // "3.14"
formatNum(1000000)     // "1e+06"   (shortest form is scientific for large magnitudes)
```

> Only the exact words `true`/`false` are recognized by `parseBool` — `"1"`/`"yes"` are rejected. For boolean coercion by truthiness, use `!!value` instead.

---

## Introspection

`typeOf` names a value's type; the `is*` predicates test one type (or emptiness). Each takes exactly one argument.

| Function | Returns | Description |
|----------|---------|-------------|
| `typeOf(v)` | string | One of `"null"`, `"number"`, `"string"`, `"boolean"`, `"array"`, `"object"`, `"datetime"`, `"duration"`. |
| `isNull(v)` | boolean | `true` only for `null`. |
| `isNumber(v)` | boolean | `true` for any number (a numeric *string* is `false`). |
| `isString(v)` | boolean | `true` for strings. |
| `isBool(v)` | boolean | `true` for booleans (`1` is not a boolean). |
| `isArray(v)` | boolean | `true` for arrays. |
| `isObject(v)` | boolean | `true` for objects (arrays are not objects). |
| `isDate(v)` | boolean | `true` for `datetime` values (a date *string* is `false`). |
| `isDuration(v)` | boolean | `true` for `duration` values. |
| `isEmpty(v)` | boolean | `true` for `null`, `""`, `[]`, `{}`; `false` for everything else (incl. `0`, `false`). |

```uexl
typeOf(d"2024-01-15")    // "datetime"
isNumber("42")           // false   (a string, not a number)
isObject([1])            // false   (arrays are not objects)
isEmpty([])              // true
isEmpty(0)               // false
```

---

## Strings

String transformations beyond the core `substr`/`contains`/unicode views.

| Function | Returns | Description |
|----------|---------|-------------|
| `upper(s)` / `lower(s)` | string | Case folding (Unicode-aware). |
| `trim(s[, cutset])` | string | Strip whitespace, or characters in `cutset`, from both ends. |
| `trimStart(s[, cutset])` / `trimEnd(s[, cutset])` | string | Trim one side only. |
| `replace(s, old, new)` | string | Replace **all** occurrences of `old`. |
| `split(s, sep)` | array | Split around `sep`; `""` splits into characters. |
| `startsWith(s, prefix)` / `endsWith(s, suffix)` | bool | Prefix/suffix test. |
| `indexOf(s, sub)` | number | **Byte** index of the first `sub`, or `-1`. |
| `repeat(s, count)` | string | `count` copies (`count ≥ 0`). |
| `padStart(s, length, pad)` / `padEnd(s, length, pad)` | string | Pad to `length` **runes** with `pad`. |

```uexl
upper("hello")              // "HELLO"
trim("xxhixx", "x")         // "hi"      (cutset is a set of chars)
replace("a-b-c", "-", "+")  // "a+b+c"
split("a,b,c", ",")         // ["a", "b", "c"]
indexOf("héllo", "llo")     // 3         (byte index: é is 2 bytes)
padStart("5", 3, "0")       // "005"
```

---

## Collections

Lookup and transformation for objects and arrays. All pure — `remove`/`merge` return **new** objects.

| Function | Returns | Description |
|----------|---------|-------------|
| `get(c, key[, default])` | any | Value at object key / array index, else `default` (`null` if omitted). |
| `has(c, key)` | bool | Whether the key exists / the index is in bounds. |
| `keys(obj)` | array | Keys, **sorted** ascending (deterministic). |
| `values(obj)` | array | Values, ordered by sorted key (`values(o)[i]` ↔ `keys(o)[i]`). |
| `remove(obj, key)` | object | Copy without `key` (input untouched). |
| `merge(a, b)` | object | Shallow merge; on a key collision **`b` wins** (inputs untouched). |

```uexl
get({a: 1}, "z", 99)              // 99     (default for a missing key)
get([10, 20, 30], 1)             // 20
has([10, 20], -1)                // false  (out of bounds)
keys({b: 2, a: 1})               // ["a", "b"]   (sorted)
merge({a: 1, b: 2}, {b: 3, c: 4})  // {a: 1, b: 3, c: 4}
```

> An object key must be a string and an array index must be an integer, otherwise `get`/`has` **error** rather than returning the default. `merge` is shallow — nested objects are replaced wholesale, not deep-merged.

---

## JSON

| Function | Returns | Description |
|----------|---------|-------------|
| `parseJson(text)` | any | Decode a JSON string to a value (all JSON numbers become `number`). |
| `toJson(value[, pretty])` | string | Encode a value to JSON; `pretty = true` indents by two spaces. |

`toJson` deep-converts `datetime` values to RFC 3339 strings and `duration` values to ISO 8601, including inside nested arrays/objects, so they round-trip portably.

```uexl
parseJson("[1,2,3]")                  // [1, 2, 3]
parseJson("{\"a\":1}")                // {a: 1}
toJson({a: 1, b: 2})                  // "{\"a\":1,\"b\":2}"
toJson(date(2024, 1, 15))            // "\"2024-01-15T00:00:00Z\""
toJson(90m)                          // "\"PT1H30M\""
```

---

## Date/Time {#datetime}

The `datetime` and `duration` *types*, *literals* (`d"2024-12-01"`, `7d`, `1.5h`), and *operators* are part of the core language ([Data Types](../data-types.md)). The functions below — construction, parsing, formatting, and calendar/epoch arithmetic — are the attachable `datetime` family (`WithDatetime()`); `now()`/`today()` additionally need an injected clock (`WithClock(ms)`).

| Function | Returns | Description |
|----------|---------|-------------|
| `now()` / `today()` | datetime | Injected clock instant / its UTC midnight; stable within one evaluation. |
| `date(y, m, d)` | datetime | Midnight-UTC datetime from components. |
| `datetime(y, m, d, h?, mi?, s?, ms?)` | datetime | Datetime from up to seven components (UTC wall clock). |
| `time(h, mi, s?, ms?)` | datetime | Time-of-day on `1970-01-01`. |
| `parseDate(s[, pattern])` | datetime | ISO 8601, or a NITES `pattern`/named layout. |
| `tryParseDate(s[, pattern])` | datetime \| null | Safe `parseDate`. |
| `formatDate(d[, pattern][, offset])` | string | NITES render; default layout `iso`; optional fixed `offset`. |
| `parseDur(s)` / `tryParseDur(s)` | duration \| null | ISO 8601 duration parse (strict / safe). |
| `formatDur(d)` | string | ISO 8601 duration string. |
| `duration(amount, unit)` | duration | Exact duration (`millisecond`…`week`). |
| `durationIn(d, unit)` | number | Express a duration in a unit (may be fractional). |
| `addMonths(d, n)` / `addYears(d, n)` | datetime | Calendar add, end-of-month clamped. |
| `diffMonths(a, b)` / `diffYears(a, b)` | number | Whole calendar months/years between, signed. |
| `datePart(d, component[, offset])` | number | Extract `year`…`millisecond`, or `weekday` (0=Sun). |
| `toEpochMillis(d)` / `fromEpochMillis(n)` | number / datetime | Unix-millisecond conversion. |
| `toEpochSeconds(d)` / `fromEpochSeconds(n)` | number / datetime | Unix-second conversion. |

```uexl
formatDate(date(2026, 6, 19), "yyyy/mm/dd")          // "2026/06/19"
formatDate(parseDate("19/06/2026", "dd/mm/yyyy"))    // "2026-06-19T00:00:00"
formatDate(date(2026, 6, 19), "iso", "+05:30")       // "2026-06-19T05:30:00"
formatDur(duration(90, "minute"))                    // "PT1H30M"
durationIn(parseDur("PT12H"), "day")                 // 0.5
formatDate(addMonths(date(2026, 1, 31), 1))          // "2026-02-28T00:00:00"  (clamped)
datePart(datetime(2026, 6, 19, 14, 30), "weekday")   // 5
```

> **NITES format dialect:** the 24-hour hour is `hhhh`, the minute is `ii`; named layouts are `iso`/`isodate`/`date`/`time`/`sql`/`rfc`. Formatting and parsing are backed by the [gotime](https://github.com/maniartech/gotime) reference implementation. See the [DateTime Specification](../../specs/datetime-spec.md) for the full grammar.
