# Expressions

Expressions are the core building blocks in UExL. They represent values, computations, function calls, or data transformations, and can be combined to form complex logic.

## Types of Expressions
- **Literals**: Direct values such as numbers, strings, booleans, or null.
  - Examples: `42`, `"hello"`, `true`, `null`
- **Variables**: Named references to values or data in the current context.
  - Examples: `x`, `user.name`, `arr[0]`
- **Binary expressions**: Operations involving two operands and an operator.
  - Examples: `a + b`, `x && y`, `price > 100`
- **Unary expressions**: Operations involving a single operand and an operator.
  - Examples: `-x`, `!flag`
- **Function calls**: Invoking built-in or user-defined functions with arguments.
  - Examples: `sum(1, 2, 3)`, `max([1, 2, 3])`, `toString(value)`
- **Pipe expressions**: Chaining operations where the output of one stage is passed as input to the next.
  - Examples: `[1, 2, 3] |map: $1 * 2`, `data |filter: $1.active`
- **Object/array member access**: Accessing properties or elements using dot or bracket notation.
  - Examples: `users[0].name`, `data["key"]`
- **Grouped expressions**: Using parentheses to control evaluation order.
  - Example: `(a + b) * c`

## Practical Usage
Expressions can be nested and combined to perform complex computations:
```
users |filter: $1.age >= 18 |map: $1.name |: join(", ")
// Filters users by age, extracts names, and joins them with commas

result = (a + b) * max(1, c)
```

## Edge Cases
- Division by zero results in an error.
- Accessing a property or index that does not exist returns `null`.
- Function calls with the wrong number of arguments may result in an error.

Mastering expressions is key to writing powerful UExL logic.