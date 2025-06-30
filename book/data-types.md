# Data Types

UExL supports several fundamental data types, each with its own characteristics and usage patterns. Understanding these types is essential for writing robust expressions.

## Numbers
Numbers can be integers or floating-point values. UExL supports decimal, scientific notation, and negative numbers.
- Examples: `1`, `-42`, `3.14`, `1e3`
- Edge cases: Leading zeros are not allowed (e.g., `01` is invalid).

## Strings
Strings are sequences of characters enclosed in single or double quotes. Both styles are supported, but quotes must match.
- Examples: `"hello"`, `'world'`, `"He said, 'hi'"`
- Edge cases: Escape sequences (e.g., `\"`, `\\`) are not currently supported; use matching quotes to include quotes inside strings.

## Booleans
Boolean values represent logical truth.
- Values: `true`, `false`
- Usage: Used in logical expressions, conditions, and comparisons.

## Null
`null` represents the absence of a value.
- Usage: Used to indicate missing or undefined data.

## Arrays
Arrays are ordered collections of values, enclosed in square brackets. Elements can be of any type, including nested arrays and objects.
- Examples: `[1, 2, 3]`, `["a", true, null, [1, 2]]`
- Access: Zero-based indexing, e.g., `arr[0]`
- Edge cases: Arrays can be empty (`[]`).

## Objects
Objects are collections of key-value pairs, enclosed in curly braces. Keys are strings (quoted or unquoted if valid identifiers), and values can be any type.
- Examples: `{"name": "UExL", "version": 1.0}`, `{id: 123, values: [1,2,3]}`
- Access: Dot or bracket notation, e.g., `obj.key`, `obj["key"]`
- Edge cases: Keys must be unique within an object.

## Examples
```
42
3.14
"hello"
true
null
[1, 2, 3]
{"name": "UExL", "features": ["pipes", "functions"]}
```

Arrays and objects can be nested and can contain any supported data type. Understanding these types is the foundation for mastering UExL.