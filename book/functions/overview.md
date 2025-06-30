# Functions Overview

Functions in UExL allow you to perform operations, calculations, and data transformations. Both built-in and user-defined functions can be called in expressions.

## Calling Functions
- Functions are called using the syntax: `functionName(arg1, arg2, ...)`
- Arguments can be literals, variables, or expressions.
- Functions can be nested and used within pipes.

## Examples
```
min(1, 2, 3)           // 1
concat("a", "b")      // "ab"
len([1, 2, 3])         // 3
max(10, 20, 5)         // 20
sum([1, 2, 3, 4])      // 10
```

## Function Signatures
- Each function has a specific signature (number and type of arguments).
- Some functions accept a variable number of arguments (e.g., `min`, `max`).
- Argument types are automatically converted if possible (see type conversion).

## Argument Rules
- Arguments are evaluated before being passed to the function.
- If an argument is missing or invalid, an error is thrown.
- Functions can return any data type (number, string, array, object, etc.).

## Advanced Usage
- Functions can be used inside pipes:
  `[1, 2, 3] |map: double($1)`
- Functions can be composed:
  `sum(map([1, 2, 3], double))`
- Functions can be used as arguments to other functions (if supported).

## Edge Cases
- Calling a function with the wrong number of arguments results in an error.
- Passing `null` as an argument may produce `null` or an error, depending on the function.

Refer to the next chapter for user-defined functions.