# Chapter 10: Understanding Pipes

> "Pipes are the heart of UExL's data transformation model. Instead of nesting function calls, you chain stages — and each stage gets the result of the previous one as structured input."

---

## 10.1 The Core Idea

A pipe expression chains transformations. Each stage receives the output of the previous stage as its *input*, and applies a *predicate* — a sub-expression — to produce the output.

```uexl
orders
  |filter: $item.total > 100
  |map:    $item.total * 0.9
  |reduce: ($acc ?? 0) + $item
```

This reads naturally: take `orders`, filter to those with total above 100, transform each total to 90% (10% discount), then sum them all. The result is a single number — the discounted sum of qualifying orders.

Each stage is independent. Each predicate sees the element(s) it receives via scope variables like `$item` and `$index`. No state leaks between stages.

---

## 10.2 Pipe Syntax

```
input |pipetype: predicate
```

Where:
- **input** — any expression producing the value that flows into the pipe
- **`|`** — literal pipe character (not bitwise OR, which needs a space `a | b`)
- **pipetype** — the name of the pipe handler (e.g., `map`, `filter`, `reduce`)
- **`:`** — required separator; the predicate follows immediately after

```uexl
[1, 2, 3, 4, 5] |map: $item * 2
//                ^^^^ ^^^^^^^^
//                type  predicate
```

### The colon is required

A bare `|` alone is always bitwise OR. The pipe operator is always `|word:` — the keyword and colon together form a single token.

```uexl
a | b           // bitwise OR
a |map: expr    // pipe operator (map type)
a |: expr       // passthrough pipe (keyword is empty)
```

### Pipe chaining

Multiple pipes chain left-to-right. The output of each stage is the input of the next:

```uexl
arr |filter: $item > 2 |map: $item * 10
```

1. `arr` → `|filter:` → filtered array
2. filtered array → `|map:` → mapped array

### The `:` passthrough pipe

The passthrough pipe `|:` (sometimes written `|pipe:`) has no keyword between `|` and `:`. Its predicate receives the entire input as `$last` and returns the expression result:

```uexl
[1, 2, 3] |map: $item * 2 |: join($last, ", ")
// maps to [2, 4, 6], then joins → "2, 4, 6"
```

`$last` holds the value flowing into the `|:` stage. Inside `|:`, you can inspect and transform that value with any expression.

---

## 10.3 Scope Variables in Pipes

Pipe predicates execute in an isolated scope. Scope variables are set by the pipe handler before each predicate evaluation.

### `$item` and `$index`

Set by most array-iterating pipes (`map`, `filter`, `find`, `some`, `every`, `sort`, `groupBy`, `flatMap`):

| Variable | Type | Meaning |
|----------|------|---------|
| `$item` | any | Current element |
| `$index` | number | Zero-based position of the element |

```uexl
["a", "b", "c"] |map: $item + str($index)
// => ["a0", "b1", "c2"]
```

### `$acc` (reduce only)

Set by `|reduce:`. Holds the accumulated result from previous iterations:

| Variable | First iteration | Subsequent |
|----------|----------------|------------|
| `$acc` | `null` | Result of previous predicate |
| `$item` | First element | Each element in order |
| `$index` | `0` | Current index |

```uexl
[1, 2, 3, 4, 5] |reduce: ($acc ?? 0) + $item
// Step 0: $acc=null, $item=1, $index=0 → ($acc ?? 0) + 1 = 1
// Step 1: $acc=1,    $item=2, $index=1 → 1 + 2 = 3
// Step 2: $acc=3,    $item=3, $index=2 → 3 + 3 = 6
// Step 3: $acc=6,    $item=4, $index=3 → 6 + 4 = 10
// Step 4: $acc=10,   $item=5, $index=4 → 10 + 5 = 15
// => 15
```

> **KEY POINT:** `$acc` is `null` on the *first* iteration. Always guard with `$acc ?? defaultValue` when the accumulator starts from nothing.

### `$last` (passthrough `|:` only)

Set by the `|:` handler. Holds the entire input value:

```uexl
[1, 2, 3] |filter: $item % 2 == 0 |: $last[0] ?? null
// filter → [2], then $last=[2], result=2
```

### `$window` (window pipe)

Set by `|window:`. Holds the current sliding window as an array of consecutive elements. Default window size is 2.

```uexl
[1, 2, 3, 4] |window: $window[0] + $window[1]
// windows: [1,2],[2,3],[3,4] → results: [3, 5, 7]
```

### `$chunk` (chunk pipe)

Set by `|chunk:`. Holds the current fixed-size chunk as an array. Default chunk size is 2.

```uexl
[1, 2, 3, 4, 5, 6] |chunk: $chunk[0] + $chunk[1]
// chunks: [1,2],[3,4],[5,6] → results: [3, 7, 11]
```

---

## 10.4 Pipe Aliases

Any pipe that uses `$item` can declare an alias — a custom name for the current element:

```uexl
orders |map as $order: $order.total * 1.1
```

The alias is set alongside `$item` — both refer to the same value. Aliases are useful when:
- The predicate is long and `$item` is ambiguous in context
- You want to be explicit for reviewers of the expression

```uexl
users
  |filter as $user: $user.active && $user.tier == 'premium'
  |map as $user: $user.name + " (" + $user.tier + ")"
```

---

## 10.5 Access to Context Variables in Pipe Predicates

Pipe predicates are not isolated from the outer context — they can freely read context variables:

```uexl
// `threshold` is a context variable
orders |filter: $item.total > threshold
```

This is intentional. The predicate executes in the context of the same expression, with `$item` / `$index` layered on top.

---

## 10.6 How the Pipe Mechanism Works (Conceptual)

When the VM encounters `|map: predicate`:

1. The **preceding expression** is evaluated — its result becomes the pipe input.
2. The **predicate** is compiled to a bytecode block (an `InstructionBlock`).
3. The **pipe handler** is looked up by name in `PipeHandlers` (e.g., `MapPipeHandler`).
4. The handler iterates the input (for array-based pipes) and calls `EvalItem(element, index)` for each element.
5. `EvalItem` sets `$item` and `$index` in the pipe scope, then runs the bytecode block in an isolated VM frame.
6. The handler collects the results and returns the output.

The predicate is compiled once and re-executed per element. This is more efficient than re-parsing the predicate expression for each element.

---

## 10.7 What Pipes Accept as Input

| Pipe | Required input type |
|------|-------------------|
| `map`, `filter`, `find`, `some`, `every`, `unique`, `sort`, `groupBy`, `flatMap` | `[]any` (array) |
| `reduce` | `[]any` (non-empty array) |
| `window` | `[]any` |
| `chunk` | `[]any` |
| `\|:` (passthrough) | any value |

Passing the wrong type (e.g., a string to `|map:`) is a runtime error. The passthrough `|:` accepts anything.

---

## 10.8 Nest Limitation: One Expression at a Time

UExL evaluates **one expression at a time, producing one result**. A pipe chain is still a single expression. You cannot define intermediate named values:

```uexl
// WRONG — UExL has no assignment / let / where clauses:
discounted = orders |filter: $item.total > 100
sum = discounted |reduce: ($acc ?? 0) + $item.total
sum * taxRate

// CORRECT — chain everything in one expression:
orders |filter: $item.total > 100 |map: $item.total |reduce: ($acc ?? 0) + $item
```

When a result genuinely needs to be named, pass it as an additional context variable from Go — compute the intermediate on the Go side and inject it.

---

## 10.9 ShopLogic: First Pipe Examples

**Total revenue from completed orders:**

```uexl
orders |filter: $item.status == 'completed' |map: $item.total |reduce: ($acc ?? 0) + $item
```

**Names of active premium customers, sorted alphabetically:**

```uexl
customers
  |filter: $item.active && $item.tier == 'premium'
  |map:    $item.name
  |sort:   $item
```

**Is any product in the cart tagged "sale"?**

```uexl
cart.items |some: $item.tags |some: $item == 'sale'
```

Note the nested pipe: the outer pipe sets `$item` to each cart item, then the inner pipe on `$item.tags` scopes a new `$item` to each tag in that item's tags array.

---

## 10.10 Summary

- The pipe operator is `|word:` — a keyword followed by a colon. The bare `|` is always bitwise OR.
- Pipe predicates are compiled once and executed per element in an isolated VM frame.
- Scope variables: `$item`/`$index` for most pipes, `$acc` for reduce, `$last` for passthrough, `$window`/`$chunk` for windowing pipes.
- `$acc` starts as `null` on the first reduce iteration — always use `$acc ?? defaultValue`.
- Pipe predicates can read outer context variables freely.
- All inputs must be arrays except the passthrough `|:`, which accepts anything.
- UExL evaluates one expression at a time — there are no intermediate named values; chain pipes or use Go to inject computed context.

---

## Exercises

**10.1 — Recall.** What is the difference between `a | b` and `a |filter: b`? What is the role of the colon in the pipe syntax?

**10.2 — Apply.** Given `orders` as an array of objects with fields `.total` (number) and `.status` (string), write a single UExL expression that:
1. Keeps only completed orders
2. Extracts each order's total
3. Returns the count of those totals

**10.3 — Extend.** Write a UExL expression that finds the order with the highest total using only pipe operators and no host functions. (Hint: use `|reduce:` with a conditional predicate. What does the result look like on the first iteration when `$acc` is `null`?)
