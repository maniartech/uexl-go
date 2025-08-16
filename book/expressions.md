# Expressions

Expressions are the core building blocks in UExL. In this chapter, you'll learn what expressions are, how to combine them, and how to use them to create powerful logic for your applications.

## What Is an Expression?
An expression is any combination of values, variables, operators, and functions that produces a result. Expressions can be as simple as a number or as complex as a chained data transformation.

## Types of Expressions
- **Literals:** Direct values such as numbers, strings, booleans, or null.
  - Example: `42`, `"hello"`, `true`, `null`
- **Variables:** Named references to values or data in the current context.
  - Example: `x`, `user.name`, `arr[0]`
- **Binary expressions:** Operations involving two operands and an operator.
  - Example: `a + b`, `x && y`, `price > 100`
- **Unary expressions:** Operations involving a single operand and an operator.
  - Example: `-x`, `!flag`, `--number` (double negation), `!!value` (boolean conversion)
- **Function calls:** Invoking built-in or user-defined functions with arguments.
  - Example: `sum(1, 2, 3)`, `max([1, 2, 3])`, `toString(value)`
- **Pipe expressions:** Chaining operations where the output of one stage is passed as input to the next.
  - Example: `[1, 2, 3] |map: $1 * 2`, `data |filter: $1.active`
- **Object/array member access:** Accessing properties or elements using dot or bracket notation.
  - Example: `users[0].name`, `data["key"]`
- **Grouped expressions:** Using parentheses to control evaluation order.
  - Example: `(a + b) * c`

## Combining Expressions: Practical Usage
Expressions can be nested and combined to perform complex computations. Here are some practical examples:
```
users |filter: $1.age >= 18 |map: $1.name |: join(", ")
// Filters users by age, extracts names, and joins them with commas

result = (a + b) * max(1, c)

score = user.points > 100 ? "high" : "low"
```

## Edge Cases and Tips
- Division by zero results in an error.
- Accessing a property or index that does not exist returns `null`.
- Function calls with the wrong number of arguments may result in an error.
- Use parentheses to clarify complex logic and control evaluation order.

## Practice: Try It Yourself
Try writing your own expressions using different types and combinations:
```
(10 + 5) * 2
user.isActive && user.score > 80
[1, 2, 3] |map: $1 * $1
min(3, 7, 2)
--42         // Double negation
!!user.name  // Boolean conversion
!!!flag      // Triple NOT
```

Mastering expressions is key to writing powerful UExL logic. In the next chapter, we'll explore operators and how they shape your expressions.