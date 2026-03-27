# Quick Start for Python/JavaScript Users

Already comfortable with Python or JavaScript? This guide shows how familiar language concepts translate to UExL, so you can jump in fast.

If you’re brand new to UExL, read the [Introduction](introduction.md) and [Syntax](syntax.md) first. For complete coverage, see [Operators](operators/overview.md), [Pipes](pipes/overview.md), and [Functions](functions/overview.md).

## Quick Mapping

| Concept | Python | JavaScript | UExL |
|--------|--------|------------|------|
| Power | `**` | `**` | `**` (or `^`) |
| Not equals | `!=` | `!==` or `!=` | `!=` (or `<>`) |
| XOR | `^` | `^` | `~` |
| Bitwise NOT | `~` | `~` | `~` |
| AND / OR / NOT | `and` / `or` / `not` | `&&` / `\|\|` / `!` | `&&` / `\|\|` / `!` |
| Nullish coalescing | N/A | `??` | `??` |
| Optional chaining | N/A | `?.` | `?.` |
| Ternary | `a if c else b` | `c ? a : b` | `c ? a : b` |

## Idioms in UExL

```javascript
// map / filter / reduce with pipes
data |filter: $item > 10 |map: $item ** 2 |reduce: $acc + $item, 0

// ternary chains
score > 90 ? "A" : score > 80 ? "B" : "C"

// nullish coalescing and optional chaining
user?.profile?.email ?? "no-email@example.com"

// string operations
first + " " + last
contains("hello", "ell")
```

### Host-provided helpers (when available)

If your runtime exposes aggregation helpers, prefer them over manual pipes:

```javascript
sum(arr)                    // preferred over arr |reduce: $acc + $item, 0
average(arr)                // preferred over (sum(arr) / len(arr))
count(arr)                  // number of items
count_if(arr, $item > 10)   // conditional count
```

When these helpers aren't available, use the pipe-based fallbacks shown above.

## Differences and Gotchas

- No `===` / `!==`: UExL has `==` and `!=` only
  - `==` is exact for primitives (no coercion) and deep for arrays/objects
- XOR uses `~` (not `^`); `^` means power
- Expression-only: no assignments or declarations inside expressions
- No implicit type coercion: be explicit (e.g., convert strings before math)
- Slicing syntax is not built-in; use helpers like `substr()` for strings

## Mini Examples

```javascript
// Pipeline example
values |filter: $item % 2 == 0 |map: $item ** 2

// Conditional with default
(user?.age ?? 0) >= 18 ? "adult" : "minor"

// Deep equality
[1, 2] == [1, 2]   // true
```

## See Also

- [Operators](operators/overview.md)
- [Pipes](pipes/overview.md)
- [Functions](functions/overview.md)
- [Examples](examples.md)
