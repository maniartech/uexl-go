# Syntax

UExL syntax is concise, readable, and expressive. It supports a variety of constructs for defining values, operations, and transformations.

## Expressions
An expression in UExL can be:
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
```

### Whitespace and Formatting
Whitespace (spaces, tabs, newlines) is generally ignored except where necessary to separate tokens. You can format expressions across multiple lines for readability:
```
result = (
    a + b
    + c
)
```

### Comments
Comments are not currently supported in UExL. All content is treated as part of the expression.

### Edge Cases
- String literals can use single or double quotes, but must match.
- Identifiers are case-sensitive: `Value` and `value` are different.
- Parentheses can be used to control evaluation order.

Understanding these basics will help you write clear and correct UExL expressions.