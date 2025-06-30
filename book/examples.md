# Examples

## Basic Arithmetic
```
10 + 20           // 30
5 * (10 - 3)      // 35
```

## Conditional Logic
```
x > 10 && y < 20  // true if x > 10 and y < 20
a == 1 || b == 2  // true if a equals 1 or b equals 2
```

## Working with Arrays
```
[1, 2, 3][1]      // 2 (second element)
len([1, 2, 3])    // 3
```

## Using Pipes
```
10 + 20 |: $1 * 2           // 60
[1, 2, 3] |map: $1 * $1     // [1, 4, 9]
```

## Complex Expressions
```
users |filter: $1.age >= 18 |map: $1.name |: join(", ")
// Filters users by age, extracts names, and joins them with commas
```