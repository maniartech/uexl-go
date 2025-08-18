# Syntax

Understanding the syntax of UExL is the first step to mastering the language. In this chapter, we'll explore the building blocks of UExL syntax, how to write clear expressions, and common patterns you'll use in real-world scenarios.

## What Is Syntax?
Syntax defines the rules for how you write expressions, values, and operations in UExL. A good grasp of syntax helps you avoid errors and write readable, maintainable code.

## Expressions: The Heart of UExL
An expression is any valid combination of values, variables, operators, and functions that produces a result. Expressions can be simple or complex, and they are the foundation of all UExL logic.

- A literal (number, string, boolean, null)
- A variable or identifier
- A binary or unary operation
- A function call
- A pipe operation
- An object or array access
- A grouped expression (using parentheses)

### Examples
```
42                  // Number literal
"hello"             // String literal
x + 10              // Binary operation
min(1, 2, 3)        // Function call
[1, 2, 3] |map: $1 * 2 // Pipe operation
user.name           // Object property access
(a + b) * c         // Grouped expression
"hello"[0]          // String index access
```

## Whitespace and Formatting
Whitespace (spaces, tabs, newlines) is generally ignored except where needed to separate tokens. You can format expressions across multiple lines for clarity:



## Writing Clear Expressions
- Use parentheses to control evaluation order.
- Break long expressions into multiple lines for readability.
- Choose descriptive variable names.

## Edge Cases and Tips
- String literals can use single or double quotes, but quotes must match.
- Identifiers are case-sensitive: `Value` and `value` are different.
- Parentheses help clarify complex logic.

## Practice: Try It Yourself
Here are some practice expressions to try in UExL:
```
(5 + 3) * 2
"UExL" + " rocks!"
[10, 20, 30] |filter: $1 > 15
user.isActive && user.score > 80
```

Mastering syntax is the first step toward writing powerful UExL expressions. In the next chapter, we'll dive into data types and how to use them effectively.