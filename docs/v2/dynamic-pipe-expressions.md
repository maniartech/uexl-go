# Dynamic Pipe Expressions (v2)

Dynamic Pipe Expressions allow you to define custom pipe stages using only UExL — no Go code required. Think of them as "pipe macros" that expand into one or more standard pipe stages at compile-time.

They mirror the ergonomics of built-in pipes and integrate with existing stage variables like `$last`, `$item`, `$index`, and `$acc`.

> Experimental: API and semantics may evolve.

## Why Dynamic Pipes?
- Reuse common pipeline patterns without copy-pasting
- Share logic across expressions and teams
- Keep pipelines readable and composable

## API (Proposed)

```go
package uexl

// Registers a pipe macro under a name.
// template is a UExL pipeline fragment that must start with a pipe (e.g., "|:", "|map:").
// Use the placeholder `$PRED` if the macro requires a predicate from the call site.
func RegisterPipeExpression(name string, template string)

// Removes a registered macro.
func UnregisterPipeExpression(name string)

// Checks if a pipe macro exists.
func HasPipeExpression(name string) bool
```

Rules:
- If the template contains `$PRED`, the pipe macro requires a predicate at the call site (`|macro: <expr>`). Omitting it is a compile-time error.
- If `$PRED` is absent, call the macro without a predicate (`|macro`).
- Templates can chain multiple stages and may reference standard stage variables.
- `$last` refers to the value from the immediately previous stage.
- In specialized pipes, `$item`, `$index`, `$acc` behave as usual.

## Authoring Templates

Templates are pipeline fragments starting with a pipe stage, for example:

```uexl
|filter: $PRED
|map: { id: $item.id, name: upper($item.name) }
|reduce: ($acc || 0) + $item.value
|: someFunc($last)
```

Use `$PRED` inside the template if your macro needs a predicate expression provided by the caller.

## Usage

### concatStr (predicate required)
Concatenate mapped string values with commas and trim the trailing comma.

Registration:

```go
uexl.RegisterPipeExpression(
  "concatStr",
  "|reduce: ($acc || '') + ($PRED) + ',' |: substr($last, 0, len($last)-1)",
)
```

Examples:

```uexl
[1,2,3] |concatStr: str($item)              # "1,2,3"
[{name:'a'},{name:'b'}] |concatStr: $item.name  # "a,b"
```

Error (predicate required):

```uexl
[1,2,3] |concatStr        # ERROR: pipe 'concatStr' requires a predicate
```

### where (alias for filter)

```go
uexl.RegisterPipeExpression("where", "|filter: $PRED")
```

```uexl
[1,2,3,4,5] |where: $item > 2   # [3,4,5]
```

### select (alias for map)

```go
uexl.RegisterPipeExpression("select", "|map: $PRED")
```

```uexl
[1,2,3] |select: $item * 2   # [2,4,6]
```

### sumBy (reduce over a derived value)

```go
uexl.RegisterPipeExpression("sumBy", "|reduce: ($acc || 0) + ($PRED)")
```

```uexl
[{price:10},{price:25}] |sumBy: $item.price   # 35
```

### Composing Macros
Macros can call other macros by name like any pipe stage.

```go
uexl.RegisterPipeExpression("selectAndConcat", "|map: $PRED |concatStr: str($item)")
```

```uexl
[1,2,3] |selectAndConcat: ($item * 2)   # "2,4,6"
```

## Design Notes (How it Works)

- At compile-time, when the compiler encounters `|name: <predicate?>` and `name` is a registered macro:
  - If the template contains `$PRED` but the `<predicate?>` is missing → error.
  - The compiler substitutes `$PRED` in the template with the predicate AST from the call site.
  - The template is parsed as if it were written inline at that location.
  - The resulting stages are then compiled normally.
- No special runtime machinery is required beyond what pipes already use.

## Best Practices
- Keep templates small and focused; compose for complexity.
- Use clear names (`where`, `select`, `sumBy`, `concatStr`).
- Avoid capturing context that templates shouldn’t rely on (keep them pure).
- Validate presence/absence of `$PRED` early to surface helpful errors.

## FAQ

- Q: Can a macro be used both with and without predicate?
  - A: In v2, we keep strictness simple: presence of `$PRED` in the template means predicate required; otherwise, disallow predicate.
- Q: Can macros access host variables?
  - A: They see the same variables the surrounding pipeline sees, including `$last`, and specialized variables within the stages they expand into.
