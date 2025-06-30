# Context

In UExL, context refers to the set of variables and values available during expression evaluation. Context determines how identifiers are resolved and how data flows through pipes and functions.

## Variable Scope
- Variables are resolved from the current context, which may include global, local, or function-specific variables.
- In pipe expressions, `$1`, `$2`, etc. refer to values passed between stages. `$1` is the output of the previous stage; for reduce pipes, `$2` is the accumulator or second argument.
- Functions and pipes may introduce their own local context for arguments, shadowing variables from outer scopes.

## Context Propagation
- When evaluating nested expressions, the current context is passed down, allowing inner expressions to access variables defined in outer scopes.
- Assignments (if supported) update the current context.
- Context can be extended or overridden in function calls and pipe stages.

## Practical Examples
```
user = {"name": "Alice", "age": 30}
user.name // Resolves to "Alice"

[1, 2, 3] |map: $1 * 2 // $1 is each element in the array

users = [{"name": "Bob", "active": true}, {"name": "Eve", "active": false}]
users |filter: $1.active |map: $1.name // Filters active users and extracts names
```

## Edge Cases
- Referencing an undefined variable returns `null`.
- Shadowing: Inner scopes can define variables with the same name as outer scopes, hiding the outer value.

Understanding context is key to writing correct and powerful UExL expressions.