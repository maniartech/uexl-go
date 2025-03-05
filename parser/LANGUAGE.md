# UExl Language Guide

UExl (Universal Expression Language) is an embeddable, platform-independent expression evaluation engine. It's designed to evaluate expressions that are not known at compile time, making applications more flexible by allowing users to define expressions through configuration files or databases.

## Table of Contents

- [Data Types](#data-types)
- [Literals](#literals)
- [Variables and Identifiers](#variables-and-identifiers)
- [Operators](#operators)
- [Operator Precedence](#operator-precedence)
- [Expressions](#expressions)
- [Functions](#functions)
- [Pipes](#pipes)
- [Objects and Arrays](#objects-and-arrays)
- [Examples](#examples)

## Data Types

UExl supports the following data types:

- **Numbers**: Integer and floating-point values (`1`, `3.14`)
- **Strings**: Text enclosed in single or double quotes (`"hello"`, `'world'`)
- **Booleans**: Logical values (`true`, `false`)
- **Null**: Represents the absence of a value (`null`)
- **Arrays**: Ordered collections of values (`[1, 2, 3]`)
- **Objects**: Collections of key-value pairs (`{"name": "UExl", "version": 1.0}`)

## Literals

### Number Literals

Numbers can be written as integers or floating-point values:

```
42        // Integer
3.14      // Floating-point
```

Scientific notation is also supported:

```
1e3       // 1000
1.5e-2    // 0.015
```

### String Literals

Strings can be enclosed in either single or double quotes:

```
"hello"   // Double quotes
'world'   // Single quotes
```

Escape sequences are supported within strings:

```
"hello\nworld"   // Contains a newline
"tab\tcharacter" // Contains a tab
```

### Raw String Literals

Raw strings are prefixed with `r` and do not process escape sequences:

```
r"C:\path\to\file"   // Backslashes are treated as literal characters
r'path\to\file'      // Equivalent to the above
```

### Boolean Literals

```
true
false
```

### Null Literal

```
null
```

### Array Literals

Arrays are ordered collections of values enclosed in square brackets:

```
[1, 2, 3]
["apple", "banana", "cherry"]
[true, false, null]
[1, "mixed", true, [1, 2]]
```

### Object Literals

Objects are collections of key-value pairs enclosed in curly braces:

```
{"name": "UExl", "version": 1.0}
{"values": [1, 2, 3], "enabled": true}
{"nested": {"a": 1, "b": 2}}
```

## Variables and Identifiers

Variables and identifiers in UExl can contain letters, numbers, underscores, and the dollar sign:

```
x
count
user_name
$value
```

Dot notation is used to access object properties:

```
user.name
data.values[0]
```

Within pipe operations, `$1` refers to the input value passed to the pipe.

## Operators

UExl provides a variety of operators organized by type and precedence (highest to lowest):

### Operator Table

| Precedence | Operator | Description | Type | Associativity |
|------------|----------|-------------|------|--------------|
| **Primary Operators** |||||
| 1 | `()` | Grouping | Grouping | Left-to-right |
| 2 | `.` | Property access | Member access | Left-to-right |
| 2 | `[]` | Index access | Member access | Left-to-right |
| **Unary Operators** |||||
| 3 | `-` (unary) | Negation | Unary | Right-to-left |
| 3 | `!` | Logical NOT | Unary | Right-to-left |
| **Arithmetic Operators** |||||
| 4 | `%` | Modulo (remainder) | Arithmetic | Left-to-right |
| 5 | `*` | Multiplication | Arithmetic | Left-to-right |
| 5 | `/` | Division | Arithmetic | Left-to-right |
| 6 | `+` | Addition | Arithmetic | Left-to-right |
| 6 | `-` | Subtraction | Arithmetic | Left-to-right |
| **Bitwise Shift Operators** |||||
| 7 | `<<` | Left shift | Bitwise | Left-to-right |
| 7 | `>>` | Right shift | Bitwise | Left-to-right |
| **Comparison Operators** |||||
| 8 | `<` | Less than | Comparison | Left-to-right |
| 8 | `>` | Greater than | Comparison | Left-to-right |
| 8 | `<=` | Less than or equal | Comparison | Left-to-right |
| 8 | `>=` | Greater than or equal | Comparison | Left-to-right |
| **Equality Operators** |||||
| 9 | `==` | Equal | Equality | Left-to-right |
| 9 | `!=` | Not equal | Equality | Left-to-right |
| **Bitwise Operators** |||||
| 10 | `&` | Bitwise AND | Bitwise | Left-to-right |
| 11 | `^` | Bitwise XOR | Bitwise | Left-to-right |
| 12 | `\|` | Bitwise OR | Bitwise | Left-to-right |
| **Logical Operators** |||||
| 13 | `&&` | Logical AND | Logical | Left-to-right |
| 14 | `\|\|` | Logical OR | Logical | Left-to-right |
| **Pipe Operators** |||||
| 15 | `\|:` | Simple pipe | Pipe | Left-to-right |
| 15 | `\|name:` | Named pipe | Pipe | Left-to-right |

## Expressions

Expressions in UExl can be:

1. **Literals**: `42`, `"hello"`, `true`
2. **Variables**: `x`, `user.name`
3. **Binary expressions**: `a + b`, `x && y`
4. **Function calls**: `sum(1, 2, 3)`, `max([1, 2, 3])`
5. **Pipe expressions**: `[1, 2, 3] |map: $1 * 2`
6. **Object/array member access**: `users[0].name`, `data["key"]`

## Functions

Functions are called by name followed by parentheses containing arguments:

```
min(1, 2, 3)
concat("hello", " ", "world")
len([1, 2, 3])
```

## Pipes

Pipes allow for the chaining of operations, passing the result of one expression to the next. The syntax is:

```
expression |[pipe_type]: next_expression
```

Where `pipe_type` is optional and, if provided, specifies the type of pipe operation.

The value of the first expression is accessible in the second expression as `$1`.

### Example

```
[1, 2, 3] |map: $1 * 2
```

In this example, the array `[1, 2, 3]` is passed to the `map` pipe, which multiplies each element by 2, resulting in `[2, 4, 6]`.

### Common Pipe Types

| Pipe Type | Description | Example |
|-----------|-------------|---------|
| `\|:` | Simple pipe (default) | `10 \|: $1 * 2` |
| `\|map:` | Maps each element in an array | `[1, 2, 3] \|map: $1 * 2` |
| `\|filter:` | Filters elements in an array | `[1, 2, 3, 4] \|filter: $1 > 2` |
| `\|reduce:` | Reduces an array to a single value | `[1, 2, 3, 4, 5] \|reduce: $1 + $2` |

### Examples

```
[1, 2, 3] |map: $1 * 2                  // Results in [2, 4, 6]
[1, 2, 3, 4] |filter: $1 > 2            // Results in [3, 4]
[1, 2, 3, 4, 5] |reduce: $1 + $2        // Results in 15
```

Multiple pipes can be chained:

```
[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2
// This filters to [3, 4, 5], maps to [6, 8, 10], then reduces to 24
```

## Objects and Arrays

### Object Access

Properties of objects can be accessed using dot notation or bracket notation:

```
user.name
user["name"]
```

### Array Access

Elements of arrays can be accessed using bracket notation with a zero-based index:

```
numbers[0]    // First element
names[2]      // Third element
```

## Examples

### Basic Arithmetic

```
10 + 20           // 30
5 * (10 - 3)      // 35
```

### Conditional Logic

```
x > 10 && y < 20  // true if x > 10 and y < 20
a == 1 || b == 2  // true if a equals 1 or b equals 2
```

### Working with Arrays

```
[1, 2, 3][1]      // 2 (second element)
len([1, 2, 3])    // 3
```

### Using Pipes

```
10 + 20 |: $1 * 2           // 60
[1, 2, 3] |map: $1 * $1     // [1, 4, 9]
```

### Complex Expressions

```
users |filter: $1.age >= 18 |map: $1.name |: join(", ")
// Filters users by age, extracts names, and joins them with commas
```

This guide covers the core syntax and features of UExl. For more advanced usage, refer to the official documentation and examples.