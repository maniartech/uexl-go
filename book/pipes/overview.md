# Pipes

Pipes in UExL allow chaining of operations, passing the result of one expression to the next. They enable powerful data transformation and functional programming patterns.

## Syntax
```
expression |: next_expression
expression |map: next_expression
```

The value of the first expression is accessible in the next stage as `$1` (and `$2`, etc. for reduce pipes).

See the following chapters for pipe types and chaining.