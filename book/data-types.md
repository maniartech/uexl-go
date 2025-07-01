# Data Types

Every expression language needs a solid foundation of data types, and UExL is no exception. In this chapter, you'll discover the core types that power UExL, how to use them, and practical examples for each.

## Why Data Types Matter
Data types define what kind of information you can work with—numbers, text, logical values, collections, and more. Understanding these types helps you write correct, expressive, and efficient expressions.

## Numbers
Numbers in UExL can be integers or floating-point values. Use them for calculations, comparisons, and more.
- **Examples:**
  - `1` (integer)
  - `-42` (negative integer)
  - `3.14` (floating-point)
  - `1e3` (scientific notation, equals 1000)
- **Edge Case:** Leading zeros are not allowed: `01` is invalid.

## Strings
Strings are sequences of characters, enclosed in single or double quotes. Use them for text, messages, and keys.
- **Examples:**
  - `"hello"`
  - `'world'`
  - `"He said, 'hi'"`
- **Edge Case:** Escape sequences (like `\"`, `\\`) are not supported; use matching quotes to include quotes inside strings.

## Booleans
Booleans represent logical truth—`true` or `false`. Use them in conditions, filters, and logical operations.
- **Examples:**
  - `true`
  - `false`
- **Usage:**
  - `x > 10 && y < 20` (returns a boolean)

## Null
`null` means "no value" or "missing data." Use it to indicate absence or undefined values.
- **Example:**
  - `null`
- **Usage:**
  - `user.middleName` might be `null` if not set.

## Arrays
Arrays are ordered lists of values, enclosed in square brackets. They can hold any type, including other arrays and objects.
- **Examples:**
  - `[1, 2, 3]`
  - `["a", true, null, [1, 2]]`
- **Access:**
  - `arr[0]` (first element)
- **Edge Case:** Arrays can be empty: `[]`.

## Objects
Objects are collections of key-value pairs, enclosed in curly braces. Keys are strings (quoted or unquoted if valid identifiers), and values can be any type.
- **Examples:**
  - `{"name": "UExL", "version": 1.0}`
  - `{id: 123, values: [1,2,3]}`
- **Access:**
  - `obj.key` or `obj["key"]`
- **Edge Case:** Keys must be unique within an object.

## Putting It All Together: Examples
```
42
3.14
"hello"
true
null
[1, 2, 3]
{"name": "UExL", "features": ["pipes", "functions"]}
```

Arrays and objects can be nested and can contain any supported data type. Mastering data types is essential for writing expressive UExL code. In the next chapter, we'll explore how to use variables and identifiers to work with your data.