# Advanced Concepts

This chapter covers advanced features and patterns in UExL.

## Operator Precedence and Associativity

UExL follows standard mathematical precedence rules:

### Precedence Levels (highest to lowest)

1. **Prefix operators**: `!`, `-`, `--`, `~` (right-associative)
2. **Power operator**: `**` (right-associative)
3. **Multiplicative**: `*`, `/`, `%` (left-associative)
4. **Additive**: `+`, `-` (left-associative)
5. **Comparison**: `<`, `<=`, `>`, `>=` (left-associative)
6. **Equality**: `==`, `!=` (left-associative)
7. **Bitwise AND**: `&` (left-associative)
8. **Bitwise XOR**: `^` (left-associative)
9. **Bitwise OR**: `|` (left-associative)
10. **Logical AND**: `&&` (left-associative)
11. **Logical OR**: `||` (left-associative)

### Right-Associative Operators

The power operator `**` is right-associative, meaning:
- `2**3**2` evaluates as `2**(3**2)` = `2**9` = `512`
- Not as `(2**3)**2` = `8**2` = `64`

This follows mathematical convention for exponentiation.

## Consecutive Unary Operators
UExL supports chaining multiple unary operators for advanced logic patterns:

### Double Negation Patterns
```
--value           // Mathematical double negation: -(-(value))
---x              // Triple negation: -(-(-(x)))
```

### Boolean Conversion Idioms
```
!!value           // Convert any value to boolean: !(!(value))
!!!flag           // Triple NOT for complex boolean logic
```

### Mixed Unary Operators
```
!-x               // NOT of negative: !(-(x))
-!condition       // Negative of NOT: -(!(condition))
```

These patterns are particularly useful for:
- Type conversions (`!!value` for boolean conversion)
- Mathematical transformations (`--x` for ensuring positive values)
- Complex logical operations (`!!!flag` for triple negation)

## Nested Pipes
Pipes can be nested for complex data flows:
```
[1, 2, 3] |map: ($1 * 2 |: $1 + 1)
```

## Aliasing in Pipes
You can alias pipe values for clarity:
```
[1, 2, 3] |map: $1 as $item |: $item * 2
```

## Custom Functions
Extend UExL by registering custom functions in your host environment.

## Extensibility
UExL is designed to be extensible with new operators, functions, and data types as needed.