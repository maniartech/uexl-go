# Chapter 11: All Pipe Types

> "Thirteen pipes ship with UExL. Each does one thing well. Together they cover the full taxonomy of collection transformation."

---

## 11.1 The Thirteen Default Pipes

`vm.DefaultPipeHandlers` registers these pipe names:

| Name | Input | Output | Primary use |
|------|-------|--------|-------------|
| `map` | `[]any` | `[]any` | Transform each element |
| `filter` | `[]any` | `[]any` | Keep elements matching a condition |
| `reduce` | `[]any` | `any` | Fold array to a single value |
| `find` | `[]any` | `any` or `null` | First element matching a condition |
| `some` | `[]any` | `bool` | True if any element matches |
| `every` | `[]any` | `bool` | True if all elements match |
| `unique` | `[]any` | `[]any` | Deduplicate elements |
| `sort` | `[]any` | `[]any` | Order elements by a key |
| `groupBy` | `[]any` | `object` | Partition elements into groups |
| `window` | `[]any` | `[]any` | Sliding window over elements (default size 2; `\|window(n):` for custom size) |
| `chunk` | `[]any` | `[]any` | Fixed-size consecutive slices (default size 2; `\|chunk(n):` for custom size) |
| `flatMap` | `[]any` | `[]any` | Map then flatten one level |
| `pipe` (alias `\|:`) | `any` | `any` | Passthrough — arbitrary transform |

All pipe predicates compile to bytecode at compile time and execute in an isolated VM frame at runtime.

---

## 11.2 `|map:` — Transform Every Element

The most common pipe. Applies the predicate to each element and collects results.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
// Signature
array |map: predicate-expr
```

```uexl
[1, 2, 3, 4, 5] |map: $item * $item           // => [1, 4, 9, 16, 25]
["a", "b", "c"] |map: $item + str($index)      // => ["a0", "b1", "c2"]

// Object transformation
products |map: {
    id:    $item.id,
    price: $item.basePrice * 1.1
}
```

**Result:** A new array of the same length as the input; `null` elements are preserved.

---

## 11.3 `|filter:` — Keep Matching Elements

Keeps only elements for which the predicate evaluates to a truthy value.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
[1, 2, 3, 4, 5] |filter: $item > 2                         // => [3, 4, 5]
orders           |filter: $item.status == 'completed'       // completed orders
products         |filter: $item.rating >= 4 && $item.stock > 0   // 4+ stars, in stock
```

**Result:** A new array that is a subset of the input (possibly empty). Order is preserved.

---

## 11.4 `|reduce:` — Fold to a Single Value

Applies the predicate repeatedly, accumulating a result.

**Scope variables:** `$acc`, `$item`, `$index`
**Initial `$acc`:** `null` (first iteration always gets `null`)

```uexl
// Sum
[1, 2, 3, 4, 5] |reduce: ($acc ?? 0) + $item          // => 15

// Product
[1, 2, 3, 4] |reduce: ($acc ?? 1) * $item              // => 24

// Maximum (no host function needed)
[3, 1, 4, 1, 5, 9, 2, 6] |reduce:
    $acc == null || $item > $acc ? $item : $acc         // => 9

// Concatenate strings
["foo", "bar", "baz"] |reduce: ($acc ?? "") + $item    // => "foobarbaz"
```

> **CRITICAL:** `$acc` is `null` on the first iteration, not the first element. Always use `$acc ?? initialValue`.

**Empty array:** `|reduce:` on an empty array is a runtime error — it requires at least one element. Guard with a check if arrays can be empty:

```uexl
len(orders) > 0 ? orders |reduce: ($acc ?? 0) + $item.total : 0
```

---

## 11.5 `|find:` — First Matching Element

Returns the first element for which the predicate is truthy, or `null` if none match.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
[1, 2, 3, 4, 5] |find: $item > 3                   // => 4  (not [4, 5])
orders           |find: $item.id == targetOrderId   // => the matching order or null
products         |find: $item.tags |some: $item == 'featured'  // first featured product
```

**Result:** Single element (any type) or `null`. Does NOT return an array.

---

## 11.6 `|some:` — Existential Check

Returns `true` if at least one element satisfies the predicate, `false` otherwise.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
[1, 2, 3, 4, 5] |some: $item > 3       // => true
[1, 2, 3]       |some: $item > 10      // => false
products         |some: $item.inStock   // are any products in stock?
```

Short-circuits — stops at the first matching element.

---

## 11.7 `|every:` — Universal Check

Returns `true` if all elements satisfy the predicate, `false` if any fail.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
[2, 4, 6, 8]   |every: $item % 2 == 0   // => true  (all even)
[2, 4, 5, 8]   |every: $item % 2 == 0   // => false (5 is odd)
cart.items      |every: $item.inStock    // can we fulfill the whole cart?
```

Short-circuits — stops at the first failing element.

---

## 11.8 `|unique:` — Deduplicate

Returns a new array containing only the first occurrence of each element (based on string representation).

**Scope variables:** None — `|unique:` takes no predicate.

Wait — actually `unique` does NOT take a predicate. Its behavior is fixed:

```uexl
[1, 2, 2, 3, 1, 4] |unique: null         // Syntax: predicate still required but result ignores it
```

> **NOTE:** The current implementation of `|unique:` deduplicates by converting each element to its `fmt.Sprintf("%v", ...)` string key. This works correctly for numbers, strings, and booleans. For objects (maps), the key is the Go default print format, which may not be stable. For object deduplication by a specific field, use `|groupBy:` and then take the first element of each group.

```uexl
// Unique string tags
tags |unique: $item     // keeps the predicate for type consistency

// Deduplicate by field (workaround):
products |groupBy: $item.category |: ...
```

---

## 11.9 `|sort:` — Order Elements

Returns a new array sorted by the value the predicate returns for each element.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
[3, 1, 4, 1, 5, 9] |sort: $item               // => [1, 1, 3, 4, 5, 9]  ascending numbers
["banana", "apple", "cherry"] |sort: $item    // => ["apple", "banana", "cherry"]  strings

// Sort objects by field
products |sort: $item.basePrice                // ascending by price
orders   |sort: $item.createdAt               // ascending by date (ISO string)
```

**Order:** Ascending by default. Numbers compare numerically; strings compare lexicographically. Mixed-type arrays (number vs string) may not sort predictably — ensure elements have consistent types.

**Descending sort:** Negate numbers or use `|: ...` to reverse afterward:

```uexl
// Descending by price (negate the sort key)
products |sort: -$item.basePrice
```

---

## 11.10 `|groupBy:` — Partition into Groups

Returns an **object** where keys are the predicate results (as strings) and values are arrays of matching elements.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
products |groupBy: $item.category
// => {
//     "electronics": [{...}, {...}],
//     "clothing": [{...}],
//     ...
// }
```

```uexl
orders |groupBy: $item.status
// => { "completed": [...], "pending": [...], "cancelled": [...] }
```

**Result:** An object (map), not an array. Access groups by key:

```uexl
// Count completed orders after grouping
(orders |groupBy: $item.status).completed |: len($last) ?? 0
```

---

## 11.11 `|window:` / `|window(n):` — Sliding Window

Applies the predicate to overlapping windows of consecutive elements. The **window size defaults to 2**. Pass a compile-time literal integer argument to use a different size.

**Scope variables:** `$window` (array of the current window), `$index` (window start index)

```uexl
// Default size 2
[1, 2, 3, 4, 5] |window: $window[0] + $window[1]
// windows: [1,2],[2,3],[3,4],[4,5]
// results: [3, 5, 7, 9]

// Explicit size via args
[1, 2, 3, 4, 5] |window(3): $window
// windows: [1,2,3],[2,3,4],[3,4,5]

[1, 2, 3, 4, 5] |window(3): $window[0] + $window[1] + $window[2]
// results: [6, 9, 12]
```

```uexl
// Detect upward trends (default size 2 is natural for pairs)
prices |window: $window[1] > $window[0]
// => [true, false, true, ...]  (true where next > current)

// 4-period moving average
prices |window(4): ($window[0] + $window[1] + $window[2] + $window[3]) / 4
```

**Arg rules:**
- Must be a literal integer `≥ 2`; fractions, zero, and negatives fall back to the default of 2.
- Whitespace around the argument is allowed: `|window( 3 ):` is valid.
- When no argument is given (`|window:`), size is 2.
- If the array is shorter than the window size, the result is an empty array.

**Boundary behavior:** Every window is always exactly `n` elements — there are no partial windows. The result contains `len(arr) - n + 1` windows. Unlike chunking, every element appears in at least one window when `len(arr) >= n`.

---

## 11.12 `|chunk:` / `|chunk(n):` — Fixed-Size Batches

Divides the array into consecutive non-overlapping chunks. The **chunk size defaults to 2**. Pass a compile-time literal integer argument to use a different size. The last chunk may be smaller than the requested size.

**Scope variables:** `$chunk` (current batch as array), `$index` (batch number, 0-based)

```uexl
// Default size 2
[1, 2, 3, 4, 5] |chunk: $chunk
// chunks: [1,2],[3,4],[5]
// result: [[1,2],[3,4],[5]]

// Explicit size via args
[1, 2, 3, 4, 5, 6] |chunk(3): $chunk
// chunks: [1,2,3],[4,5,6]

[1, 2, 3, 4, 5] |chunk(3): $chunk[0] + $chunk[1] + ($chunk[2] ?? 0)
// chunks: [1,2,3],[4,5]  — last chunk has 2 elements
// results: [6, 9]
```

Use `?? 0` when accessing beyond the end of the last (potentially short) chunk to avoid a runtime error.

```uexl
// Process in batches of 100, keep only full batches
records |chunk(100): $chunk |filter: len($chunk) == 100
```

**Arg rules:**
- Must be a literal integer `≥ 2`; fractions, zero, and negatives fall back to the default of 2.
- Whitespace around the argument is allowed: `|chunk( 4 ):` is valid.
- When no argument is given (`|chunk:`), size is 2.

**Boundary behavior:** The result always contains `⌈len(arr) / n⌉` chunks. The last chunk holds the remaining elements and may be shorter than `n`. If the array is shorter than `n`, exactly one chunk is produced containing all elements. When the length is an exact multiple of `n`, all chunks are the same size.

---

## 11.13 `|flatMap:` — Map and Flatten

Each element maps to a value; if that value is an array, it is flattened into the result. Non-array values are kept as-is.

**Scope variables:** `$item`, `$index`, alias (optional)

```uexl
// Expand tags: each product has a `tags` array; collect all tags
products |flatMap: $item.tags
// => ["sale", "new", "featured", "new", ...]  (all tags from all products)

// Generate pairs
[1, 2, 3] |flatMap: [$item, $item * 10]
// => [1, 10, 2, 20, 3, 30]
```

---

## 11.14 `|:` — Passthrough Transform

The passthrough pipe applies a single expression to the entire input value, accessible as `$last`.

```uexl
|: predicate
```

**Scope variables:** `$last` (the entire pipe input)

```uexl
[1, 2, 3] |: len($last)               // => 3
[1, 2, 3] |: $last[0]                 // => 1

runes("hello") |: join($last, "-")    // => "h-e-l-l-o"

// Chain: filter then join
runes("hello world") |filter: $item != " " |: join($last, "")
// => "helloworld"
```

The passthrough is the adapter between pipes and non-array values. Any time you need to finish a pipe chain with a non-iterating operation (join, index, count), use `|:`.

---

## 11.15 Pipe Combination Patterns

### Map then reduce (common aggregation)

```uexl
orders |map: $item.total |reduce: ($acc ?? 0) + $item  // sum of totals
```

### Filter then count

```uexl
orders |filter: $item.status == 'completed' |: len($last)
```

### Map to objects (projection)

```uexl
products |map: {id: $item.id, name: $item.name, price: $item.basePrice}
```

### Nested pipes (pipe inside pipe predicate)

```uexl
// Products where any tag is "sale"
products |filter: $item.tags |some: $item == 'sale'
```

The inner `$item` refers to each tag inside the predicate of `|some:`, not the outer product. Pipe scopes are stacked.

### Sort then take first (top/bottom element)

```uexl
products |sort: $item.basePrice |: $last[0]    // cheapest product
```

---

## 11.16 ShopLogic: Complete Pipe Showcase

**Revenue breakdown by customer tier:**

```uexl
orders |groupBy: $item.customer.tier
```

**Top 3 best-selling products by orders received:**

```uexl
orders
  |groupBy: $item.productId
  |: {id: $last.id, count: len($last.items ?? [])}
```

Wait — `groupBy` returns an object, not directly iterable. For complex post-groupBy work, the result is typically consumed on the Go side. The expression returns the grouped object; Go iterates it.

**All unique tags across the product catalogue:**

```uexl
products |flatMap: $item.tags |unique: $item
```

**Average order total:**

```uexl
orders |map: $item.total |: (($last |reduce: ($acc ?? 0) + $item) / len($last))
```

Wait — nested pipe in passthrough: the inner `($last |reduce: ...)` evaluates `$last` (the array from the outer `|map:`) in a new reduce. This is valid.

**Check if every item in a cart is in stock:**

```uexl
cart.items |every: $item.product.stock > 0
```

---

## 11.17 Summary

- UExL ships 13 default pipe types: `map`, `filter`, `reduce`, `find`, `some`, `every`, `unique`, `sort`, `groupBy`, `window`, `chunk`, `flatMap`, and the passthrough `|:`.
- Array pipes require `[]any` input; the passthrough accepts anything.
- `$acc` starts as `null` in `|reduce:` — always guard with `$acc ?? initial`.
- `|find:` returns `null`, not an empty array, when nothing matches.
- `|sort:` sorts ascending; negate numeric keys for descending.
- `|window(n):` and `|chunk(n):` accept a compile-time literal integer argument for the window/chunk size; both default to 2 when no argument is provided.
- `|groupBy:` returns an object, not an array.
- `|:` gives you `$last`, the full input — use it to apply a single expression to a pipe result.
- Pipe scopes stack — nested pipes each get their own `$item`/`$index`.

---

## Exercises

**11.1 — Recall.** Which pipes modify the length of the result? Which always preserve length? Which return a non-array result?

**11.2 — Apply.** Using only built-in pipes and functions, write expressions for:
1. The sum of all `price` fields from `items` where `price > 50`.
2. All unique `category` strings from `products`.
3. The first product alphabetically by `name`.
4. Whether all orders in `orders` have a defined (non-null) `customerId`.

**11.3 — Extend.** Write a UExL expression that:
- Takes `transactions` (array of `{amount, type}` objects where type is `"credit"` or `"debit"`)
- Returns the net balance: sum of credits minus sum of debits
- Each stage should be a single pipe (no host functions needed)

(Hint: separate credits and debits with two passes, or use `|reduce:` with a conditional accumulator.)
