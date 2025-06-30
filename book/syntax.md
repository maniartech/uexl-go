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

### Examples
```
42
"hello"
x + 10
min(1, 2, 3)
[1, 2, 3] |map: $1 * 2
user.name
```

Whitespace is generally ignored except where necessary to separate tokens. Comments are not currently supported in UExL.