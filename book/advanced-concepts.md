# Advanced Concepts

This chapter covers advanced features and patterns in UExL.

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