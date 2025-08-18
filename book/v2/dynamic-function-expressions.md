# Dynamic Function Expressions (v2)

Dynamic Function Expressions let you define new functions using plain UExL, avoiding host-language (Go) code. They are analogous to pipe macros but for function calls.

> Experimental: API and behavior may evolve.

## Why Dynamic Functions?
- Share reusable logic as first-class functions
- Keep configuration and business logic close to UExL
- Reduce boilerplate for simple utility functions

## API (Proposed)

```go
package uexl

// Registers a function implemented by a UExL expression template.
// Template is a normal UExL expression that can reference positional arguments $1, $2, ...
func RegisterFunctionExpression(name string, template string)

func UnregisterFunctionExpression(name string)

func HasFunctionExpression(name string) bool
```

### Argument Passing
- Use `$1`, `$2`, `$3`, ... to reference call-site arguments.
- Arguments are not pre-evaluated specially; they are normal expression arguments and can be expressions themselves.
- Arity is not enforced by the template system; you may add validation in your implementation.

## Authoring Templates
Templates are ordinary UExL expressions. They can call functions, access operators, and even start a pipeline from an argument via `|:` if desired.

Examples:

### 1) concatStr (function form)
Equivalent to a function that concatenates array items into a comma-separated string via a pipeline.

Registration:

```go
uexl.RegisterFunctionExpression(
  "concatStr",
  "($1 |reduce: ($acc || '') + str($item) + ',' |: substr($last, 0, len($last)-1))",
)
```

Usage:

```uexl
concatStr([1,2,3])           # "1,2,3"
concatStr(map(users, $1.name))
```

### 2) between (predicate-style numeric check)

```go
uexl.RegisterFunctionExpression("between", "$1 >= $2 && $1 <= $3")

42 |: between($last, 10, 100)    # true
```

### 3) coalesceWith

```go
uexl.RegisterFunctionExpression("coalesceWith", "$1 != null ? $1 : $2")

coalesceWith(null, "fallback")   # "fallback"
```

### 4) project (map over array)

```go
uexl.RegisterFunctionExpression(
  "project",
  "$1 |map: $2",
)

project(users, $item.name)
```

## Design Notes (How it Works)

- When compiler encounters a function call `name(args...)` and `name` is registered:
  - The template is parsed once and a substitution pass replaces `$1`, `$2`, ... with the actual argument ASTs.
  - The result is compiled as if written inline.
- No special VM changes required; bytecode is emitted for the expanded expression.

## Best Practices
- Keep templates concise and pure.
- Use clear function names and document expected argument shapes.
- Prefer pipelines inside the template when transforming collections.
