# Advanced Concepts

This chapter covers advanced features and patterns in UExL.

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