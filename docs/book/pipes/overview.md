# Pipes

Pipes are one of the most powerful features of UExL, enabling you to chain operations and build expressive data transformation pipelines. It allows you to process the expression in a sequential manner, through series of pipes. Each pipe takes the output of the previous stage as input and performs a specific operation, such as mapping, filtering, or reducing the data. This allows your expressions to remain clean and readable, without the need for complex nested function calls or complex control structures.

In this chapter, you'll learn what pipes are, how they work, and how to use them effectively with practical examples.

## What Are Pipes?

A pipe takes the output of one expression and passes it as input to the next stage. This lets you build pipelines that process data step by step, making your logic more readable and maintainable.

- The value from the previous stage is accessible as `$last` in the next stage for simple pipes (`|:`).
- Pipes can be chained to perform multiple transformations in sequence.
- Pipes are especially useful for working with arrays and collections.

## Pipe Syntax

```uexl
expression |: next_expression
expression |map: next_expression
```

- `|:` is the default pipe, passing the value as `$last`.
- `|map:` and `|filter:` are specialized pipes for array processing, exposing `$item` (current element) and `$index` (current index).
- `|reduce:` exposes `$acc` (accumulator), `$item` (current element), and `$index` (current index).

### Emitted Context Variables

Every pipe stage emits a set of context variables that the predicate (the expression to the right of the colon) can consume.

All pipe types ALWAYS emit:

- `$last` – the result value produced by the immediately previous stage (or the initial left-hand expression for the first stage).

Specialized pipes additionally emit:

- `map`, `filter`, `unique`, `find`, `some`, `every`, `flatMap`, `sort`, `groupBy`:
  - `$item` – current element
  - `$index` – current zero-based index
- `reduce`:
  - `$acc` – accumulator value so far
  - `$item` – current element
  - `$index` – current index
- `window`:
  - `$window` – the current sliding window (array)
  - `$index` – window start index
- `chunk`:
  - `$chunk` – the current chunk (array)
  - `$index` – chunk index

Notes:

- `$last` is always available, even inside specialized pipes (e.g., inside a `map` stage `$last` is the entire array input to that `map`).
- Predicate expressions must only reference variables emitted by the current or previous stages.

#### Context Variable Reference (Quick Lookup)

| Pipe Type / Category | Always Emitted | Additional Variables | Notes |
| -------------------- | -------------- | -------------------- | ----- |
| All pipes            | `$last`        | —                    | Previous stage result (initially the left expression) |
| `|:` default         | `$last`        | —                    | Pass-through; write any expression using `$last` |
| map, flatMap         | `$last`        | `$item`, `$index`    | Transform each element; flatMap flattens one level |
| filter               | `$last`        | `$item`, `$index`    | Keep elements where predicate is truthy |
| unique               | `$last`        | `$item`, `$index`    | Optional key/path logic (by predicate expression) |
| find                 | `$last`        | `$item`, `$index`    | Returns first matching element (or null) |
| some / every         | `$last`        | `$item`, `$index`    | Boolean result (short-circuit semantics) |
| reduce               | `$last`        | `$acc`, `$item`, `$index` | Accumulates; predicate returns new accumulator |
| sort                 | `$last`        | `$item`, `$index`    | Comparator key expression per element |
| groupBy              | `$last`        | `$item`, `$index`    | Key expression groups elements |
| window               | `$last`        | `$window`, `$index`  | Sliding window; index = window start |
| chunk                | `$last`        | `$chunk`, `$index`   | Fixed-size subsets |

> `$last` inside a specialized stage references the whole incoming collection (or value) given to that stage.

#### Multi-Stage Pipeline Example (Annotated)

```uexl
orders
  |filter: $item.status == 'paid'          # Keep only paid orders
  |map: { id: $item.id, total: $item.total }
  |sort: $item.total                       # Sort by total ascending
  |reduce: $acc + $item.total              # Sum totals
```

Explanation:

1. `filter` exposes `$item`, `$index`; `$last` is original `orders`.
2. `map` transforms each element into a simplified object.
3. `sort` evaluates `$item.total` as key for ordering.
4. `reduce` receives `$acc` and `$item`.
5. Final result is the numeric grand total.

## Practical Examples

- **Transforming data:**
  - `[1, 2, 3] |map: $item * 2` // Returns `[2, 4, 6]`
- **Filtering:**
  - `users |filter: $item.active` // Returns only active users
- **Aggregating:**
  - `[1, 2, 3] |reduce: $acc + $item` // Sums the array
- **Chaining:**
  - `products |filter: $item.price < 50 |map: $item.name` // Gets names of affordable products

## Tips for Using Pipes

- Use pipes to break complex logic into clear, sequential steps.
- Remember that `$last` refers to the value from the previous stage for all pipes; specialized pipes add more variables.
- Combine pipes with functions and operators for powerful transformations.
- Use the nullish coalescing operator `??` to provide safe defaults for missing values without overriding valid falsy ones:
  - `data |: ($last.items ?? []) |map: $item.name`
  - `values |reduce: ($acc ?? 0) + $item`

## Practice: Try It Yourself

```uexl
[10, 20, 30] |map: $item / 10
users |filter: $item.isAdmin |map: $item.email
[1, 2, 3, 4] |filter: $item % 2 == 0 |reduce: $acc + $item
```

Mastering pipes will help you write concise, readable, and powerful UExL code. In the next chapter, we'll explore the different types of pipes and how to chain them for advanced data processing.

The value of the first expression is accessible in the next stage as `$last` for simple pipes, and as `$item`, `$index`, `$acc` for specialized pipes.

See the following chapters for pipe types and chaining.

## Developer policy: custom pipe stages

When implementing custom pipe handlers:

- Pure by default: do not mutate the incoming collection or external state. Produce new values.
- Copy-returning updates: if a stage conceptually updates an accumulator (e.g., reduce), ensure the helper used (like `set`) returns a copy rather than mutating the original.
- Deterministic evaluation: process items left-to-right; avoid hidden re-evaluation or dependence on prior side effects.
- Short-circuit semantics: honor standard UExL short-circuiting in your predicate expressions; optional chaining must prevent evaluating index/key expressions when base is nullish.
- Clear naming: if a stage can cause side effects or is stateful, make it explicit in the name and documentation.

See also: Mutability and Purity.