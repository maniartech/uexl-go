# Type Conversion

UExL supports both implicit and explicit type conversion to make expressions flexible and robust.

## Implicit Type Conversion
- UExL automatically converts types in certain contexts:
  - Arithmetic operations: strings and booleans are converted to numbers if possible.
  - Logical expressions: non-boolean values are converted to booleans.
  - Equality checks: types are converted as needed for comparison.
- If conversion fails, the result is `null` or an error.

### Examples
```
"10" + 5         // 15 (string converted to number)
true + 1          // 2 (true is 1)
"abc" * 2        // null ("abc" cannot be converted to number)
1 == "1"          // true (string converted to number)
0 && "hello"      // false (0 is false)
```

## Explicit Type Conversion
- Use built-in functions to convert values:
  - `toNumber(value)`
  - `toString(value)`
  - `toBoolean(value)`
- If conversion is not possible, the result is `null`.

### Examples
```
toNumber("42")      // 42
toNumber("abc")     // null
toString(123)       // "123"
toBoolean(0)        // false
toBoolean("hello")  // true
```

## Edge Cases
- Converting `null` to any type returns `null`.
- Converting arrays or objects to numbers returns `null`.
- Converting arrays or objects to strings returns a string representation (implementation-defined).
- Converting non-empty arrays/objects to boolean returns `true`.

Understanding type conversion is key to writing correct and predictable UExL expressions.