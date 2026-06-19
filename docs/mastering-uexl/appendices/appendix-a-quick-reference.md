# Appendix A: Quick Reference Card

---

## Literals

| Type | Examples |
|------|---------|
| Number | `42` `3.14` `-7` `1e6` `0xFF` (hex) `0b1010` (binary) `0o17` (octal) |
| String | `'hello'` `"world"` |
| Boolean | `true` `false` |
| Null | `null` |
| Array | `[1, 2, 3]` `['a', 'b']` `[]` |
| Object | `{name: 'Alice', age: 30}` `{}` |
| Datetime | `d"2024-12-01"` `d"2024-12-01T10:30:00Z"` |
| Duration | `7d` `1.5h` `30ms` `45s` `2w` |

---

## Operators (highest to lowest precedence)

| Precedence | Operator | Description |
|-----------|----------|-------------|
| 16 | `.` `?.` `[]` `?.[]` | Member access, optional chaining |
| 15 | `-x` `!x` `~x` | Unary minus, logical NOT, bitwise NOT |
| 14 | `**` `^` | Power (right-associative) |
| 13 | `*` `/` `%` | Multiply, divide, modulo |
| 12 | `+` `-` | Add, subtract |
| 11 | `<<` `>>` | Bitwise shift |
| 10 | `??` | Nullish coalescing |
| 9 | `<` `>` `<=` `>=` | Comparison |
| 8 | `==` `!=` `<>` | Equality / inequality |
| 7 | `&` | Bitwise AND |
| 6 | `~` | Bitwise XOR |
| 5 | `\|` | Bitwise OR |
| 4 | `&&` | Logical AND (short-circuit) |
| 3 | `\|\|` | Logical OR (short-circuit) |
| 2 | `? :` | Ternary conditional (right-associative) |
| 1 | `\|map:` `\|filter:` etc. | Pipe operators (lowest) |

---

## Built-in Functions (core, 13)

| Function | Description |
|----------|-------------|
| `len(v)` | Byte count of string, or element count of array |
| `substr(s, start, length)` | Byte-level substring |
| `contains(s, sub)` | Byte-level substring check → bool |
| `set(obj, key, value)` | Mutates object in-place AND returns it |
| `str(v)` | Any value to string |
| `runeLen(s)` | Count Unicode code points |
| `runeSubstr(s, start, length)` | Code-point-level substring |
| `graphemeLen(s)` | Count grapheme clusters |
| `graphemeSubstr(s, start, length)` | Grapheme-level substring |
| `runes(s)` | Explode string to `[]any` of single-rune strings |
| `graphemes(s)` | Explode string to `[]any` of grapheme strings |
| `bytes(s)` | Explode string to `[]any` of byte values (float64) |
| `join(arr)` | Join array of strings with empty separator |
| `join(arr, sep)` | Join array of strings with separator |

---

## Standard Library (attach with `WithStdlib()`)

Opt-in families of pure helpers; see [Appendix C](appendix-c-builtin-functions.md#the-standard-library) for full signatures. String parsers come in strict (`parseX`, errors) / safe (`tryParseX`, `null`) pairs.

| Family | Option | Functions |
|--------|--------|-----------|
| Math | `WithMath()` | `abs` `sign` `round` `floor` `ceil` `trunc` `sqrt` `min` `max` `sum` `avg` `mod` `pow` `clamp` |
| Conversion | `WithConversion()` | `parseNum` `tryParseNum` `parseBool` `tryParseBool` `formatNum` |
| Introspection | `WithIntrospection()` | `typeOf` `isNull` `isNumber` `isString` `isBool` `isArray` `isObject` `isDate` `isDuration` `isEmpty` |
| Strings | `WithStrings()` | `upper` `lower` `trim` `trimStart` `trimEnd` `replace` `split` `startsWith` `endsWith` `indexOf` `repeat` `padStart` `padEnd` |
| Collections | `WithCollections()` | `get` `has` `keys` `values` `remove` `merge` |
| JSON | `WithJSON()` | `parseJson` `toJson` |
| Date/Time | `WithDatetime()` | `now` `today` `date` `datetime` `time` `parseDate` `tryParseDate` `formatDate` `parseDur` `tryParseDur` `formatDur` `duration` `durationIn` `addMonths` `addYears` `diffMonths` `diffYears` `datePart` `to/fromEpochMillis` `to/fromEpochSeconds` |

> `min`/`max`/`sum`/`avg` take either an array or several numeric args. `round(x)` has no decimals parameter — use `formatNum(x, n)`. `now()`/`today()` also need `WithClock(ms)`.

---

## Pipe Operators

| Syntax | Scope Variables | Effect |
|--------|----------------|--------|
| `arr \|map: expr` | `$item`, `$index` | Transform each element |
| `arr \|filter: bool` | `$item`, `$index` | Keep truthy elements |
| `arr \|reduce: expr` | `$acc` (null on first!), `$item`, `$index` | Fold to single value |
| `arr \|find: bool` | `$item`, `$index` | First matching element or null |
| `arr \|some: bool` | `$item`, `$index` | True if any element matches |
| `arr \|every: bool` | `$item`, `$index` | True if all elements match |
| `arr \|unique:` | — | Deduplicated array |
| `arr \|sort: bool` | `$item`, `$index` | Sorted array (asc by default) |
| `arr \|groupBy: key` | `$item`, `$index` | Object keyed by computed value |
| `arr \|window: expr` | `$window`, `$index` | Sliding window (default size 2) |
| `arr \|window(n): expr` | `$window`, `$index` | Sliding window of `n` elements |
| `arr \|chunk: expr` | `$chunk`, `$index` | Chunks of 2 elements each |
| `arr \|chunk(n): expr` | `$chunk`, `$index` | Chunks of `n` elements each |
| `arr \|flatMap: expr` | `$item`, `$index` | Map then flatten |
| `value \|: expr` | `$last` | Passthrough / default pipe |

### Pipe Alias Syntax
```uexl
arr |map as $result: $result.price * 0.9
```

---

## Optional Chaining and Nullish Safety

| Syntax | Behavior |
|--------|---------|
| `a?.b` | Returns null if `a` is null; otherwise `a.b` |
| `a?.[i]` | Returns null if `a` is null; otherwise `a[i]` |
| `a ?? b` | Returns `a` if not null; otherwise `b` |
| `a?.b ?? c` | Optional chain with null fallback |

> `??` falls back ONLY on `null`. It preserves `false`, `0`, and `""`.

---

## Slicing

```uexl
arr[start:end]        # elements at index start..end-1
arr[start:]           # from start to end
arr[:end]             # from beginning to end-1
arr[:]                # full copy / all characters
arr[start:end:step]   # with step
```

Slicing works on both arrays (element-level) and strings (byte-level).

---

## Go API at a Glance

```go
// One-shot (for scripts/tests)
result, err := uexl.Eval(expr, vars)

// Production — compile once
env := uexl.DefaultWith(
    uexl.WithFunctions(myFuncs),
    uexl.WithGlobals(map[string]any{"TAX_RATE": 0.08}),
)
compiled, err := env.Compile(expr)  // validates function names
result, err  := compiled.Eval(ctx, vars)  // goroutine-safe

// Extend for tenant isolation
tenantEnv := env.Extend(uexl.WithGlobals(tenantConfig))
```
