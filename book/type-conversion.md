# Type Conversion

UExL provides both implicit and explicit type conversion mechanisms.

## Implicit Type Conversion
- Arithmetic operations may convert strings to numbers if possible.
- Logical expressions use truthiness: non-zero numbers and non-empty strings are truthy; `false` and `null` are falsy.
- Equality checks may perform type coercion.

## Explicit Type Conversion
Use built-in functions for explicit conversion:
- `toNumber(value)` — Converts to number, errors if not possible.
- `toString(value)` — Converts to string.
- `toBoolean(value)` — Converts to boolean.

## Example
```
toNumber("42") // 42
toString(123)  // "123"
toBoolean(0)   // false
```