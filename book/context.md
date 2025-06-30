# Context

In UExL, context refers to the set of variables and values available during expression evaluation. Context determines how identifiers are resolved and how data flows through pipes and functions.

## Variable Scope
- Variables are resolved from the current context.
- In pipe expressions, `$1`, `$2`, etc. refer to values passed between stages.
- Functions and pipes may introduce their own local context for arguments.

## Example
```
user = {"name": "Alice", "age": 30}
user.name // Resolves to "Alice"

[1, 2, 3] |map: $1 * 2 // $1 is each element in the array
```

Understanding context is key to writing correct and powerful UExL expressions.