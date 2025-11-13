# Quick Start for Excel Users

Already familiar with Excel formulas? This guide maps the most common Excel patterns to UExL so you can get productive quickly.

New to UExL? Start with the core docs first: see [Introduction](introduction.md) and [Syntax](syntax.md). For complete details, refer to the [Operators](operators/overview.md), [Pipes](pipes/overview.md), and [Functions](functions/overview.md) sections.

## Quick Mapping (Excel → UExL)

| Excel | UExL | Notes |
|-------|------|-------|
| `A1^2` | `value ^ 2` | Power (same as Excel) |
| `A1<>B1` | `a <> b` | Not equals |
| `IF(c,a,b)` | `c ? a : b` | Ternary (like nested IFs) |
| `AND(a,b)` | `a && b` | Logical AND |
| `OR(a,b)` | `a \|\| b` | Logical OR |
| `A1 & B1` | `a + b` | String concat |
| `SUM(arr)` | `sum(arr)` or `arr \|reduce: $acc + $item, 0` | Prefer host `sum()` when available |
| `AVERAGE(arr)` | `average(arr)` or `(arr \|reduce: $acc + $item, 0) / len(arr)` | Prefer host `average()` |
| `COUNTA(arr)` | `count(arr)` or `len(arr)` | `count()` if provided by host |
| `COUNTIF(arr, cond)` | `count_if(arr, cond)` or `arr \|filter: cond \|: len($item)` | Prefer host `count_if()` |
| `VLOOKUP` concept | `obj.property` | Use objects; add `?.` for safe access |

## Common Patterns

```javascript
// IF / nested IF
score > 90 ? "A" : score > 75 ? "B" : "C"

// Ranges → arrays + pipes
values |filter: $item > 10 |map: $item ^ 2 |reduce: $acc + $item, 0

// Text
first + " " + last            // concat
status <> ""                  // not empty

// Lookup with fallback
customer?.email ?? "No email"
```

### Host-provided functions (if available)

If your embedding environment exposes Excel-like helpers, prefer them for clarity and performance:

- `sum(arr)` instead of manual reduce
- `average(arr)` instead of `sum(arr) / len(arr)`
- `count(arr)` for counting items
- `count_if(arr, cond)` for conditional counts

Fallbacks using pipes (shown in the table) remain valid and portable when host helpers are not present.

## Key Differences from Excel

- 0-based indexing: `arr[0]` is the first element
- Case-sensitive string comparison: `"A" <> "a"` is true
- No cell references or `$` anchors; use variable names instead
- No implicit type coercion; be explicit in conversions
- Property/index access is strict; use `?.` to guard nullish bases

## Mini Examples

```javascript
// Commission
sales > 10000 ? sales * 0.15 : sales > 5000 ? sales * 0.10 : sales * 0.05

// Discount tier
quantity >= 100 ? price * 0.9 : quantity >= 50 ? price * 0.95 : price

// Validation (all filled)
name <> "" && email <> "" && phone <> ""
```

## See Also

- [Operators](operators/overview.md)
- [Pipes](pipes/overview.md)
- [Functions](functions/overview.md)
- [Examples](examples.md)
