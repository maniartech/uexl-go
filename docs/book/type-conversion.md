# Type Conversion

UExL supports both implicit and explicit type conversion to make expressions flexible and robust.

## Implicit Type Conversion
- UExL automatically converts types in certain contexts:
  - Arithmetic operations: strings and booleans are converted to numbers if possible.
  - Logical expressions: non-boolean values are converted to booleans.
  - Equality checks: types are converted as needed for comparison.
- If conversion fails, the result is `null` or an error.

### Examples
```
"10" + 5         // 15 (string converted to number)
true + 1          // 2 (true is 1)
"abc" * 2        // null ("abc" cannot be converted to number)
1 == "1"          // true (string converted to number)
0 && "hello"      // false (0 is false)
```

## Explicit Type Conversion

Value-to-string is the core built-in `str` (always available). String parsing and number formatting are part of the [standard library](functions/standard-library.md#conversion) **conversion** family — attach it with `uexl.WithStdlib()` (or `WithConversion()`). Parsing follows the strict/safe convention: **`parseX`** raises an error on bad input, while **`tryParseX`** returns `null`.

| Conversion | Function | On bad input |
|------------|----------|--------------|
| any → string | `str(v)` | — (core built-in) |
| string → number | `parseNum(s)` / `tryParseNum(s)` | error / `null` |
| string → boolean | `parseBool(s)` / `tryParseBool(s)` | error / `null` |
| number → string | `formatNum(x[, decimals])` | error |

Use the double NOT operator (`!!`) to convert any value to a boolean by truthiness.

### Examples
```
str(123)             // "123"
str(true)            // "true"

parseNum("42")       // 42
parseNum("abc")      // error
tryParseNum("abc")   // null    (safe variant)

parseBool("true")    // true
tryParseBool("yes")  // null

formatNum(3.14159, 2)  // "3.14"

// Boolean conversion with double NOT
!!1                 // true
!!0                 // false
!!"text"            // true
!!""                // false
!!null              // false
!![]                // true (non-empty array)
!!{}                // true (object)
```

## Edge Cases
- `str(null)` returns `"<nil>"` (Go's nil format), not `"null"`. Handle nulls first if you need JSON-style output: `v == null ? "null" : str(v)`.
- `parseBool` accepts only the exact words `true`/`false` (case-insensitive); `"1"`/`"yes"` are rejected.
- `tryParseNum`/`tryParseBool` return `null` for a non-string argument as well as for unparseable text.
- `formatNum` with a fixed `decimals` rounds half-to-even; with no `decimals` it uses the shortest form, so large magnitudes render in scientific notation (`formatNum(1000000)` → `"1e+06"`).

Understanding type conversion is key to writing correct and predictable UExL expressions.