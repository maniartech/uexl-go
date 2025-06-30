# Expressions

Expressions are the core building blocks in UExL. They can represent values, computations, function calls, or data transformations.

## Types of Expressions
- **Literals**: `42`, `"hello"`, `true`, `null`
- **Variables**: `x`, `user.name`
- **Binary expressions**: `a + b`, `x && y`
- **Unary expressions**: `-x`, `!flag`
- **Function calls**: `sum(1, 2, 3)`, `max([1, 2, 3])`
- **Pipe expressions**: `[1, 2, 3] |map: $1 * 2`
- **Object/array member access**: `users[0].name`, `data["key"]`

## Example
```
users |filter: $1.age >= 18 |map: $1.name |: join(", ")
```
This filters users by age, extracts names, and joins them with commas.