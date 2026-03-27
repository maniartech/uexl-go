# Pipe Types

UExL ships with 12 built-in pipe types covering the full spectrum of collection-processing patterns. Every pipe stage emits `$last` (the result of the previous stage) in addition to its own scope variables.

## Quick Reference

| Pipe | Scope variables | Description |
|------|----------------|-------------|
| `\|:` | `$last` | Passthrough — transform a single value |
| `\|map:` | `$item`, `$index` | Transform each element; returns a new array |
| `\|filter:` | `$item`, `$index` | Keep elements where predicate is truthy |
| `\|reduce:` | `$acc`, `$item`, `$index` | Fold array to a single value |
| `\|find:` | `$item`, `$index` | First matching element, or `null` |
| `\|some:` | `$item`, `$index` | `true` if any element matches (short-circuits) |
| `\|every:` | `$item`, `$index` | `true` if all elements match (short-circuits) |
| `\|sort:` | `$item`, `$index` | Sort by predicate key (ascending) |
| `\|unique:` | `$item`, `$index` | Deduplicate by predicate key |
| `\|groupBy:` | `$item`, `$index` | Group into `map[string][]any` by key |
| `\|flatMap:` | `$item`, `$index` | Map then flatten one level |
| `\|chunk:` | `$chunk`, `$index` | Split into fixed-size sub-arrays |
| `\|window:` | `$window`, `$index` | Sliding window sub-arrays |

## Passthrough: `|:`

Passes the value forward as `$last`. Use for single-value transformations or to rename a result.

```uexl
10 |: $last * 2                          // 20
"hello" |: upper($last)                  // "HELLO"
orders |filter: $item.paid |: $last      // alias the filtered result
```

## Transformation: `|map:`

Applies the predicate to each element and returns a new array.

```uexl
[1, 2, 3] |map: $item * 2               // [2, 4, 6]
users |map: $item.name                   // ["Alice", "Bob"]
[1, 2, 3] |map: { val: $item, i: $index } // [{val:1,i:0}, ...]
```

## Selection: `|filter:`

Keeps only elements for which the predicate evaluates to truthy.

```uexl
[1, 2, 3, 4, 5] |filter: $item > 2     // [3, 4, 5]
users |filter: $item.active             // active users only
```

## Accumulation: `|reduce:`

Folds the array into a single value. `$acc` starts as the first element; `$item` starts at index 1. Use `??` to provide a safe default when the accumulator starts as `null`.

```uexl
[1, 2, 3, 4, 5] |reduce: ($acc ?? 0) + $item     // 15
[1, 2, 3, 4, 5] |reduce: ($acc ?? 1) * $item     // 120
["a","b","c"] |reduce: ($acc ?? "") + $item       // "abc"
```

- Reduces over an empty array return `null` (no elements, no accumulator).
- Non-array input is an error.

## Search: `|find:`

Returns the first element for which the predicate is truthy, or `null` if none match.

```uexl
[1, 2, 3, 4] |find: $item > 2          // 3
users |find: $item.id == targetId       // first matching user or null
```

## Boolean checks: `|some:` and `|every:`

Return a boolean. Both short-circuit: `some` stops at the first truthy, `every` stops at the first falsy.

```uexl
[1, 2, 3] |some: $item > 2             // true
[1, 2, 3] |every: $item > 0            // true
[1, 2, 3] |every: $item > 2            // false
```

## Ordering: `|sort:`

Sorts the array by the value the predicate returns for each element (ascending).

```uexl
[3, 1, 2] |sort: $item                 // [1, 2, 3]
users |sort: $item.name                // alphabetical by name
```

## Deduplication: `|unique:`

Returns a new array with duplicates removed. The predicate selects the key used for comparison.

```uexl
[1, 2, 2, 3] |unique: $item            // [1, 2, 3]
users |unique: $item.id                // users with distinct ids
```

## Grouping: `|groupBy:`

Returns a `map[string][]any` grouping elements by the string the predicate returns.

```uexl
products |groupBy: $item.category
// → { "electronics": [...], "clothing": [...] }
```

## Flat mapping: `|flatMap:`

Maps each element and then flattens the result by one level. Useful when the predicate returns an array for each element.

```uexl
[[1, 2], [3, 4]] |flatMap: $item       // [1, 2, 3, 4]
orders |flatMap: $item.lineItems       // flat list of all line items
```

## Chunking: `|chunk:`

Splits the array into sub-arrays of the given fixed size. The predicate must evaluate to a positive integer.

```uexl
[1, 2, 3, 4, 5] |chunk: 2             // [[1,2],[3,4],[5]]
```

Inside a following pipe stage, `$chunk` is the current sub-array and `$index` is the chunk number.

## Sliding window: `|window:`

Produces overlapping sub-arrays of the given size, sliding one element at a time.

```uexl
[1, 2, 3, 4] |window: 3               // [[1,2,3],[2,3,4]]
arr |window: 2 |map: $window[0] + $window[1]  // pairwise sums
```

Inside a following pipe stage, `$window` is the current window array and `$index` is the window start position.