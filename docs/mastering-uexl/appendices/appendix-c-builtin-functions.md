# Appendix C: Built-in Function Reference

UExL has **13 core built-in functions**, registered in `vm.Builtins` and always available in every embedding (documented in full below). Beyond those, UExL ships an **attachable standard library** — families of pure helpers (`upper`, `round`, `split`, …) that the host opts into with `uexl.WithStdlib()` (or per-family options like `WithMath()`); these are summarized in the [Standard Library](#the-standard-library) section. Anything else is a host-registered function via `WithFunctions`.

---

## `len(v)`

Returns the length of a string (in bytes) or an array (in elements).

**Signature**: `len(value) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `v` | string or array | The value to measure |

**Returns**: Float64 count

**Errors**: Throws if `v` is neither a string nor an array.

```uexl
len('hello')          # 5
len('héllo')          # 6  (é is 2 bytes in UTF-8)
len([1, 2, 3])        # 3
```

> Use `runeLen` or `graphemeLen` for human-visible character counts.

---

## `substr(s, start, length)`

Returns a byte-level substring.

**Signature**: `substr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based start byte index |
| `length` | number (integer) | Number of bytes to extract |

**Errors**: Throws if `start` or `length` are not non-negative integers, or if the range is out-of-bounds.

```uexl
substr('hello', 0, 3)    # 'hel'
substr('hello', 2, 3)    # 'llo'
```

> For code-point-safe slicing use `runeSubstr`. For grapheme-safe slicing use `graphemeSubstr`.

---

## `contains(s, sub)`

Returns `true` if `s` contains `sub` as a byte-level substring.

**Signature**: `contains(s, sub) → bool`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to search in |
| `sub` | string | The substring to search for |

```uexl
contains('hello world', 'world')    # true
contains('hello world', 'xyz')      # false
contains('hello', '')               # true  (empty string is contained in any string)
```

---

## `set(obj, key, value)`

Sets a key on an object (mutates in-place) and returns the object.

**Signature**: `set(obj, key, value) → obj`

| Argument | Type | Description |
|----------|------|-------------|
| `obj` | object (`map[string]any`) | The object to modify |
| `key` | string or number | The key to set |
| `value` | any | The value to assign |

> **WARNING**: `set` mutates the original Go map. If you need immutability, clone the map in Go before passing it to the expression.

```uexl
set({}, 'name', 'Alice')                    # {name: 'Alice'}
set(product, 'discountedPrice', price * 0.9)  # modifies product AND returns it
```

---

## `str(v)`

Converts any value to its string representation using Go's `fmt.Sprintf("%v", v)` format.

**Signature**: `str(v) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `v` | any | The value to convert |

```uexl
str(42)        # '42'
str(3.14)      # '3.14'
str(true)      # 'true'
str(false)     # 'false'
str(null)      # '<nil>'
str([1, 2])    # '[1 2]'
```

> `str(null)` returns `'<nil>'` (Go's default nil format), not `'null'`. If you need JSON-style null representation, handle nulls before conversion: `v == null ? 'null' : str(v)`.

---

## `runeLen(s)`

Returns the number of Unicode code points (runes) in a string.

**Signature**: `runeLen(s) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to measure |

```uexl
runeLen('hello')    # 5
runeLen('héllo')    # 5  (é is 1 code point, even though it's 2 bytes)
runeLen('日本語')   # 3
```

---

## `runeSubstr(s, start, length)`

Returns a substring measured in Unicode code points.

**Signature**: `runeSubstr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based code point index |
| `length` | number (integer) | Number of code points to extract |

```uexl
runeSubstr('héllo', 0, 2)    # 'hé'
runeSubstr('日本語', 1, 2)   # '本語'
```

---

## `graphemeLen(s)`

Returns the number of user-perceived characters (extended grapheme clusters) in a string.

**Signature**: `graphemeLen(s) → number`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to measure |

```uexl
graphemeLen('café')    # 4  (é is one grapheme cluster)
graphemeLen('👨‍👩‍👧‍👦')  # 1  (family emoji is one grapheme cluster)
```

---

## `graphemeSubstr(s, start, length)`

Returns a substring measured in grapheme clusters (user-perceived characters).

**Signature**: `graphemeSubstr(s, start, length) → string`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The source string |
| `start` | number (integer) | Zero-based grapheme index |
| `length` | number (integer) | Number of graphemes to extract |

---

## `runes(s)`

Explodes a string into an array of single-rune strings.

**Signature**: `runes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
runes('hi')      # ['h', 'i']
runes('日本')   # ['日', '本']
```

---

## `graphemes(s)`

Explodes a string into an array of grapheme cluster strings.

**Signature**: `graphemes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
graphemes('café')   # ['c', 'a', 'f', 'é']
```

---

## `bytes(s)`

Explodes a string into an array of byte values (as float64).

**Signature**: `bytes(s) → array`

| Argument | Type | Description |
|----------|------|-------------|
| `s` | string | The string to explode |

```uexl
bytes('hi')    # [104, 105]  (ASCII codes for 'h' and 'i')
```

---

## `join(arr)` / `join(arr, sep)`

Joins an array of strings into a single string, with an optional separator.

**Signature**:
- `join(arr) → string` — joins with empty string separator
- `join(arr, sep) → string` — joins with `sep` between elements

| Argument | Type | Description |
|----------|------|-------------|
| `arr` | array of strings | The strings to join |
| `sep` | string (optional) | Separator string; defaults to `""` |

**Errors**: Throws if any element of `arr` is not a string.

```uexl
join(['a', 'b', 'c'])          # 'abc'
join(['a', 'b', 'c'], ', ')    # 'a, b, c'
join(['hello', 'world'], ' ')  # 'hello world'
```

---

## The Standard Library

The functions above are the **core** built-ins. UExL also ships an **attachable standard library** of pure helper families. Attach all of them with `uexl.WithStdlib()`, or pick families individually (`WithMath()`, `WithStrings()`, …). Every function is pure (no mutation), and string parsers follow the strict/safe convention: `parseX` errors on bad input, `tryParseX` returns `null`.

```go
env := uexl.DefaultWith(uexl.WithStdlib())            // everything
env := uexl.DefaultWith(uexl.WithMath(), uexl.WithStrings())  // à la carte
```

### Math — `WithMath()`

`min`/`max`/`sum`/`avg` accept **either a single array or several numeric arguments**.

| Function | Returns | Description |
|----------|---------|-------------|
| `abs(x)` | number | Absolute value. |
| `sign(x)` | number | `1`, `-1`, or `0`. |
| `round(x)` | number | Nearest integer, halves away from zero. **One argument** — no decimals param (use `formatNum`). |
| `floor(x)` / `ceil(x)` / `trunc(x)` | number | Round down / up / toward zero. |
| `sqrt(x)` | number | Square root; negative → `NaN`. |
| `min` / `max` / `sum` / `avg` | number | Aggregations; `sum()` is `0`, the others need ≥1 number. |
| `mod(a, b)` | number | Remainder with the sign of `a`. |
| `pow(base, exp)` | number | Power; fractional exponents allowed. |
| `clamp(x, lo, hi)` | number | Constrain to `[lo, hi]`; errors if `lo > hi`. |

```uexl
round(2.5)        # 3     (away from zero)
min([3, 1, 2])    # 1
mod(-7, 3)        # -1
clamp(15, 0, 10)  # 10
```

### Conversion — `WithConversion()`

| Function | Returns | Description |
|----------|---------|-------------|
| `parseNum(s)` / `tryParseNum(s)` | number \| null | Parse a number (strict / safe). |
| `parseBool(s)` / `tryParseBool(s)` | boolean \| null | Parse `"true"`/`"false"` (strict / safe). |
| `formatNum(x[, decimals])` | string | Format a number; fixed `decimals` round half-to-even. |

```uexl
parseNum("  -7.5  ")   # -7.5
tryParseNum("abc")     # null
formatNum(3.14159, 2)  # "3.14"
```

### Introspection — `WithIntrospection()`

| Function | Returns | Description |
|----------|---------|-------------|
| `typeOf(v)` | string | `"null"`/`"number"`/`"string"`/`"boolean"`/`"array"`/`"object"`/`"datetime"`/`"duration"`. |
| `isNull` `isNumber` `isString` `isBool` `isArray` `isObject` `isDate` `isDuration` | boolean | Single-type tests. |
| `isEmpty(v)` | boolean | `true` for `null`, `""`, `[]`, `{}`. |

```uexl
typeOf(d"2024-01-15")  # "datetime"
isNumber("42")         # false
isEmpty([])            # true
```

### Strings — `WithStrings()`

| Function | Returns | Description |
|----------|---------|-------------|
| `upper(s)` / `lower(s)` | string | Case folding. |
| `trim` / `trimStart` / `trimEnd` `(s[, cutset])` | string | Strip whitespace or `cutset` chars. |
| `replace(s, old, new)` | string | Replace **all** occurrences. |
| `split(s, sep)` | array | Split (`""` → characters). |
| `startsWith` / `endsWith` `(s, x)` | bool | Prefix / suffix test. |
| `indexOf(s, sub)` | number | **Byte** index, or `-1`. |
| `repeat(s, count)` | string | `count` copies. |
| `padStart` / `padEnd` `(s, length, pad)` | string | Pad to `length` **runes**. |

```uexl
trim("xxhixx", "x")         # "hi"
split("a,b,c", ",")         # ["a", "b", "c"]
padStart("5", 3, "0")       # "005"
```

### Collections — `WithCollections()`

All pure; `remove`/`merge` return new objects.

| Function | Returns | Description |
|----------|---------|-------------|
| `get(c, key[, default])` | any | Object/array read with optional default. |
| `has(c, key)` | bool | Key/index present. |
| `keys(obj)` / `values(obj)` | array | Sorted keys / values by sorted key. |
| `remove(obj, key)` | object | Copy without `key`. |
| `merge(a, b)` | object | Shallow merge; `b` wins on collisions. |

```uexl
get({a: 1}, "z", 99)               # 99
keys({b: 2, a: 1})                 # ["a", "b"]
merge({a: 1}, {a: 9, b: 2})        # {a: 9, b: 2}
```

### JSON — `WithJSON()`

| Function | Returns | Description |
|----------|---------|-------------|
| `parseJson(text)` | any | JSON string → value. |
| `toJson(value[, pretty])` | string | Value → JSON; `datetime`/`duration` serialize to ISO 8601. |

```uexl
parseJson("[1,2,3]")     # [1, 2, 3]
toJson({a: 1, b: 2})     # "{\"a\":1,\"b\":2}"
```

### Date/Time — `WithDatetime()`

Construction (`date`, `datetime`, `time`), parsing/formatting (`parseDate`/`tryParseDate`, `formatDate`, `parseDur`/`tryParseDur`, `formatDur`), calendar arithmetic (`addMonths`/`addYears`, `diffMonths`/`diffYears`, `datePart`, `duration`/`durationIn`), and epoch conversion (`to`/`fromEpochMillis`, `to`/`fromEpochSeconds`). `now()`/`today()` additionally need an injected clock (`WithClock(ms)`). Formatting/parsing use the gotime **NITES** dialect (24-hour `hhhh`, minute `ii`; named layouts `iso`/`sql`/`rfc`/…). See the [DateTime Specification](../../specs/datetime-spec.md).

```uexl
formatDate(parseDate("19/06/2026", "dd/mm/yyyy"))  # "2026-06-19T00:00:00"
formatDate(addMonths(date(2026, 1, 31), 1))        # "2026-02-28T00:00:00"  (clamped)
durationIn(parseDur("PT12H"), "day")               # 0.5
```

---

## Function Availability Summary

| Group | Functions | Availability |
|-------|-----------|--------------|
| **Core** | `len`, `substr`, `contains`, `set`, `str`, `runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`, `runes`, `graphemes`, `bytes`, `join` | Always (13, `vm.Builtins`) |
| Math | `abs`, `sign`, `round`, `floor`, `ceil`, `trunc`, `sqrt`, `min`, `max`, `sum`, `avg`, `mod`, `pow`, `clamp` | `WithMath()` |
| Conversion | `parseNum`, `tryParseNum`, `parseBool`, `tryParseBool`, `formatNum` | `WithConversion()` |
| Introspection | `typeOf`, `isNull`, `isNumber`, `isString`, `isBool`, `isArray`, `isObject`, `isDate`, `isDuration`, `isEmpty` | `WithIntrospection()` |
| Strings | `upper`, `lower`, `trim`, `trimStart`, `trimEnd`, `replace`, `split`, `startsWith`, `endsWith`, `indexOf`, `repeat`, `padStart`, `padEnd` | `WithStrings()` |
| Collections | `get`, `has`, `keys`, `values`, `remove`, `merge` | `WithCollections()` |
| JSON | `parseJson`, `toJson` | `WithJSON()` |
| Date/Time | `now`, `today`, `date`, `datetime`, `time`, `parseDate`, `tryParseDate`, `formatDate`, `parseDur`, `tryParseDur`, `formatDur`, `duration`, `durationIn`, `addMonths`, `addYears`, `diffMonths`, `diffYears`, `datePart`, `toEpochMillis`, `fromEpochMillis`, `toEpochSeconds`, `fromEpochSeconds` | `WithDatetime()` (+ `WithClock` for `now`/`today`) |
| All standard-library families at once | — | `WithStdlib()` |
