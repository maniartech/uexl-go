# Pipes

Pipes in UExL allow chaining of operations, passing the result of one expression to the next. They enable powerful data transformation and functional programming patterns, making complex logic more readable and maintainable.

## What Are Pipes?
A pipe takes the output of one expression and passes it as input to the next stage. This allows you to build pipelines for processing data step by step.

- The value from the previous stage is accessible as `$1` in the next stage.
- Pipes can be chained to perform multiple transformations in sequence.
- Pipes are especially useful for working with arrays and collections.

## Syntax
```
expression |: next_expression
expression |map: next_expression
```

- `|:` is the default pipe, passing the value as `$1`.
- `|map:`, `|filter:`, and `|reduce:` are specialized pipes for array processing.

## Practical Scenarios
- Transforming data: `[1, 2, 3] |map: $1 * 2`
- Filtering: `users |filter: $1.active`
- Aggregating: `[1, 2, 3] |reduce: $1 + $2`

The value of the first expression is accessible in the next stage as `$1` (and `$2`, etc. for reduce pipes).

See the following chapters for pipe types and chaining.