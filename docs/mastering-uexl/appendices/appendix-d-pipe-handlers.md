# Appendix D: Pipe Handler Reference

UExL ships with 13 built-in pipe handlers registered in `vm.DefaultPipeHandlers`. Custom handlers can be added via `WithPipeHandlers`.

---

## Pipe Syntax

```uexl
input |pipeName: predicate
input |pipeName as $alias: predicate
```

- `input` must be an array for all pipes except `|:` (passthrough)
- `predicate` is an expression evaluated once per element (or once for the whole input for `unique`)
- `$alias` gives the scope variable a custom name for chaining

---

## `|map:`

Transforms each element into a new value. Returns an array of the same length.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: array (same length)

```uexl
prices |map: $item * 0.9
products |map: $item.name
[1, 2, 3] |map: $item * $item
```

---

## `|filter:`

Keeps elements for which the predicate is truthy. Returns a subset.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: array (subset)

```uexl
orders |filter: $item.total > 100
products |filter: $item.inStock
users |filter: $item.active && $item.age >= 18
```

---

## `|reduce:`

Folds the array into a single value by applying the predicate to each element, accumulating a result.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$acc` | any | Accumulated value (`null` on first iteration!) |
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: non-empty array (empty array → runtime error)
**Output**: any (the final accumulated value)

> **CRITICAL**: `$acc` is `null` on the very first iteration. Always guard it: `($acc ?? initial) + $item`.

```uexl
numbers |reduce: ($acc ?? 0) + $item
orders  |reduce: ($acc ?? 0) + $item.total
items   |reduce: ($acc ?? '') + $item.name + ' '
```

---

## `|find:`

Returns the first element for which the predicate is truthy, or `null` if none matches.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: element or `null` (never throws on "not found")

```uexl
orders |find: $item.id == targetId
users  |find: $item.email == 'alice@example.com'
```

---

## `|some:`

Returns `true` if at least one element satisfies the predicate, `false` otherwise. Short-circuits on first match.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: bool

```uexl
products |some: $item.inStock
orders   |some: $item.status == 'overdue'
```

---

## `|every:`

Returns `true` if all elements satisfy the predicate, `false` otherwise. Short-circuits on first non-matching element.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: bool

```uexl
orders |every: $item.isPaid
items  |every: $item.quantity > 0
```

---

## `|unique:`

Removes duplicate elements. Uniqueness is determined by the string representation of each element (`fmt.Sprintf("%v", v)`).

| Scope variable | None |
|----------------|------|

**Input**: array
**Output**: deduplicated array (first occurrence preserved)

```uexl
categories |unique:
[1, 2, 1, 3, 2] |unique:    # [1, 2, 3]
```

> Object uniqueness is based on `fmt.Sprintf("%v", obj)`. Two objects with the same keys and values are considered equal. Two objects with the same keys but in different insertion order may compare differently.

---

## `|sort:`

Returns a sorted copy of the array. Default: ascending. Provide a predicate to extract a sort key.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: sorted array

```uexl
prices   |sort: $item           # ascending numeric sort
products |sort: $item.name      # sort by name field
orders   |sort: $item.createdAt # sort by date string (lexicographic for ISO format)
```

---

## `|groupBy:`

Groups elements into an object where keys are the computed predicate values. Each key maps to an array of elements sharing that key.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: object (`map[string][]any`)

```uexl
products |groupBy: $item.category
# Returns: { 'electronics': [...], 'clothing': [...], ... }

orders |groupBy: $item.status
# Returns: { 'pending': [...], 'complete': [...], 'cancelled': [...] }
```

---

## `|window:`

Produces a sliding window of 2 consecutive elements at a time.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$window` | array | Current window (array of 2 elements) |
| `$index` | number | Zero-based window index |

**Input**: array (minimum 2 elements recommended)
**Output**: array of windows, where each window is an array of 2 elements

```uexl
prices |window: {a: $window[0], b: $window[1], delta: $window[1] - $window[0]}
# For [10, 20, 15], produces:
# [{a:10, b:20, delta:10}, {a:20, b:15, delta:-5}]
```

> Window size is fixed at 2. Custom window sizes require a custom pipe handler.

---

## `|chunk:`

Divides the array into fixed-size chunks of 2 elements each.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$chunk` | array | Current chunk (array of up to 2 elements) |
| `$index` | number | Zero-based chunk index |

**Input**: array
**Output**: array of chunks (each chunk is an array of up to 2 elements; last chunk may be shorter)

```uexl
[1, 2, 3, 4, 5] |chunk: $chunk
# Produces: [[1, 2], [3, 4], [5]]
```

> Chunk size is fixed at 2. For other sizes, register a custom pipe handler.

---

## `|flatMap:`

Maps each element to an array, then flattens one level.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$item` | any | Current element |
| `$index` | number | Zero-based element index |

**Input**: array
**Output**: flattened array

```uexl
orders |flatMap: $item.lineItems
# For [{lineItems:[a,b]},{lineItems:[c]}], produces: [a, b, c]
```

---

## `|:` (Passthrough / Default Pipe)

The default pipe passes the input through unchanged, exposing it as `$last`. Most useful for chaining without transformation, or as a named alias point.

| Scope variable | Type | Value |
|----------------|------|-------|
| `$last` | any | The entire input value |

**Input**: any
**Output**: the predicate result

```uexl
total |: $last * TAX_RATE + $last
price |: $last > MAX_DISCOUNT ? MAX_DISCOUNT : $last
```

---

## Custom Pipe Handlers

Register additional pipes with `WithPipeHandlers`:

```go
env := uexl.DefaultWith(
    uexl.WithPipeHandlers(map[string]vm.PipeHandler{
        "take": takePipe,
        "skip": skipPipe,
    }),
)
```

See Chapter 14 for the full `PipeHandler` function signature and implementation pattern.
