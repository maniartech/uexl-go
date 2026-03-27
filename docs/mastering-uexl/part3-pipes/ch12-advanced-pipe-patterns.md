# Chapter 12: Advanced Pipe Patterns

> "Pipes are linear by design, but the patterns they unlock are surprisingly rich. This chapter explores the space beyond the basics: nested scopes, chaining strategies, performance, and where to stop."

---

## 12.1 Scope Stacking in Nested Pipes

When a pipe predicate contains another pipe expression, UExL stacks the inner scope on top of the outer. The inner `$item` shadows the outer `$item` for the duration of the inner predicate:

```uexl
// Outer: $item is each product
// Inner: $item is each tag of product
products |filter: $item.tags |some: $item == 'sale'
```

Here:
- Outer `|filter:` sets `$item = each product`
- Predicate: `$item.tags |some: $item == 'sale'`
  - `$item.tags` uses the outer `$item` (the product)
  - `|some:` creates a new scope where `$item = each tag string`
- The outer `$item` is no longer accessible inside `|some:` unless an alias was declared

Use aliases to preserve outer scope access:

```uexl
products |filter as $product: $product.tags |some: $item == 'sale' && $product.rating > 3
```

Now `$product` persists through the inner `|some:` scope, while `$item` refers to each tag.

---

## 12.2 Chaining Strategy: Keep Chains Readable

Long pipe chains can become hard to follow. Structure them for clarity:

### Rule 1: Filter early

Move `|filter:` before expensive `|map:` or `|reduce:` to minimize the number of elements processed:

```uexl
// Less efficient: map then filter
products |map: expensiveComputation($item) |filter: $item.score > threshold

// Better: filter first
products |filter: $item.eligible |map: expensiveComputation($item)
```

### Rule 2: Extract early with `|map:`

When multiple downstream pipes need the same sub-field, extract once:

```uexl
// Repeated field access
orders |filter: $item.total > 100 |sort: $item.total |reduce: ($acc ?? 0) + $item.total

// Better: extract total first
orders |filter: $item.total > 100 |map: $item.total |sort: $item |reduce: ($acc ?? 0) + $item
```

### Rule 3: Use `|:` as a type adapter

When you want to apply a non-iterating operation to the result of a pipeline:

```uexl
// Get the length of a filtered set
products |filter: $item.inStock |: len($last)

// Join a mapped list of strings
names |map: $item.first + " " + $item.last |: join($last, "; ")
```

---

## 12.3 The Accumulate-and-Test Pattern

A common reduce pattern tests an accumulated condition against each element. Use `|reduce:` with a conditional accumulator that carries forward a complex value:

```uexl
// Find the maximum (no host max() needed)
[5, 2, 9, 1, 7] |reduce: $acc == null || $item > $acc ? $item : $acc
// => 9

// Merge objects into one (reduce as merge)
[{a: 1}, {b: 2}, {c: 3}] |reduce:
    set(set($acc ?? {}, "merged", true), $item.key, $item.value)
```

---

## 12.4 The Extract-Check-Report Pattern

A common business pattern: filter a collection, check a condition, and produce a summary:

```uexl
// How many orders is the customer late on?
orders |filter: $item.dueDate < today && $item.status != 'completed' |: len($last)
```

```uexl
// Are there any high-value at-risk orders?
orders |some: $item.total > 1000 && $item.status == 'at_risk'
```

```uexl
// Find the most recent order (assumes ISO date strings → lexicographic sort works)
orders |sort: $item.createdAt |: $last[len($last) - 1]
```

---

## 12.5 Building Objects in Pipes

`|map:` can produce objects — the predicate can be an object literal:

```uexl
products |map: {
    id:           $item.id,
    name:         $item.name,
    displayPrice: str($item.basePrice) + " USD",
    inStock:      $item.stock > 0
}
// => array of enriched product objects
```

This is UExL's projection pattern — selecting and renaming fields from a collection.

---

## 12.6 Two-Pass vs. Single-Pass

Sometimes what looks like two separate expressions can be collapsed into one. Sometimes it can't. Know the difference:

**Can collapse:**
```uexl
// Two logical passes: filter, then reduce
// One expression:
orders |filter: $item.total > 100 |map: $item.total |reduce: ($acc ?? 0) + $item
```

**Cannot collapse (genuinely separate concerns):**
```uexl
// Expression 1: How many qualifying orders?
orders |filter: $item.total > 100 |: len($last)

// Expression 2: What is their total?
orders |filter: $item.total > 100 |map: $item.total |reduce: ($acc ?? 0) + $item
```

Two questions = two UExL expressions, called from Go with the same compiled bytecode.

---

## 12.7 Pipe Performance Characteristics

Every pipe predicate call requires:
1. Setting scope variables (fast — a field write in the fast-path scope)
2. Resetting the VM instruction pointer (fast)
3. Running the bytecode block (linear in instruction count)

For large arrays, avoid per-element string allocations in predicates:

```uexl
// Slower: builds a string per element just to compare
products |filter: str($item.id) == targetId

// Better: ensure context is already the right type (pass targetId as string from Go)
products |filter: $item.id == targetId
```

For extremely large arrays (thousands of elements), consider doing a first pass in Go before the expression, reducing the size of the input array.

---

## 12.8 Custom Pipe Handlers

Just as you can register custom functions, you can register custom pipe handlers. A pipe handler is any function with this signature:

```go
type PipeHandler func(ctx vm.PipeContext, input any) (any, error)
```

**Example: a `take` pipe that returns the first N elements:**

```go
takePipeHandler := func(ctx vm.PipeContext, input any) (any, error) {
    arr, ok := input.([]any)
    if !ok {
        return nil, fmt.Errorf("take pipe expects an array")
    }
    // Evaluate predicate once to get N
    n, err := ctx.EvalWith(map[string]any{})
    if err != nil {
        return nil, err
    }
    count, ok := n.(float64)
    if !ok || count < 0 {
        return nil, fmt.Errorf("take expects a non-negative number")
    }
    limit := int(count)
    if limit > len(arr) {
        limit = len(arr)
    }
    return arr[:limit], nil
}

handlers := vm.DefaultPipeHandlers
handlers["take"] = takePipeHandler

machine := vm.New(vm.LibContext{
    Functions:    vm.Builtins,
    PipeHandlers: handlers,
})
```

Usage in expressions:
```uexl
products |sort: $item.rating |take: 5     // top 5 by rating
```

Chapter 14 covers custom pipe patterns and the `PipeContext` API in detail.

---

## 12.9 Pipe Cancellation

In long-running pipes (e.g., large arrays), you can pass a context with a deadline or cancellation signal:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
defer cancel()

result, err := machine.RunWithContext(ctx, bytecode, contextVars)
```

The VM checks the context at each opcode boundary. If the context is cancelled or times out, the VM returns the context error immediately. This is the primary mechanism for bounding pipe execution time.

---

## 12.10 When Not to Use Pipes

Pipes are the right tool for **collection processing**. For scalar or single-value logic, plain operators and function calls are cleaner and faster:

```uexl
// BAD: unnecessary pipe for a single value
[product.basePrice] |map: $item * 1.1 |: $last[0]

// GOOD: just multiply directly
product.basePrice * 1.1
```

```uexl
// BAD: pipe to check a scalar
[product.rating] |some: $item >= 4

// GOOD: direct comparison
product.rating >= 4
```

If you find yourself wrapping single values in arrays to apply pipe operators, reconsider — the operator approach is simpler and much faster.

---

## 12.11 ShopLogic: Advanced Pipe Scenarios

**Generate a discount summary for each product:**

```uexl
products
  |filter: $item.stock > 0
  |map: {
      id:       $item.id,
      name:     $item.name,
      price:    $item.basePrice,
      onSale:   $item.tags |some: $item == 'sale',
      lowStock: $item.stock < 10
  }
```

**Compute a running balance from transactions:**

```uexl
transactions |reduce:
    ($acc ?? 0) + ($item.type == 'credit' ? $item.amount : -$item.amount)
```

**Find all products that need restocking:**

```uexl
products |filter: $item.stock < $item.minStock
```

**Tag frequency across all products:**

```uexl
products
  |flatMap: $item.tags
  |reduce: set($acc ?? {}, $item, ($acc[$item] ?? 0) + 1)
```

This builds a frequency object: `{"sale": 5, "new": 3, "featured": 2, ...}`.

---

## 12.12 Summary

- When pipes nest, scopes stack — inner `$item` shadows outer `$item`. Use aliases to preserve outer references.
- Filter early to minimize downstream work; extract fields with `|map:` before sorting or reducing.
- Use `|:` as an adapter between pipe results and non-iterating expressions.
- `|reduce:` can carry complex state — objects, arrays, or single values.
- Custom pipe handlers let you extend the pipe vocabulary; register them in `PipeHandlers`.
- Provide a `context.Context` with a deadline if pipe execution should be time-bounded.
- Don't use pipes for scalar operations — plain operators are cleaner and faster.

---

## Exercises

**12.1 — Recall.** What happens to outer `$item` when a pipe predicate contains an inner pipe? How do you preserve access to the outer element?

**12.2 — Apply.** For the ShopLogic order dataset (orders with `total`, `status`, `customer.tier`), write a single expression that computes the average order total for premium-tier customers only. The expression should use no host functions.

**12.3 — Extend.** Write a custom `skip` pipe handler in Go that skips the first N elements of an array, where N is provided as the pipe predicate. Register it alongside `vm.DefaultPipeHandlers` and write an expression that uses it: `products |sort: $item.rating |skip: 5 |: $last[0]` to get the 6th-highest-rated product.
