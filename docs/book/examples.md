# Examples

## Basic Arithmetic
```
10 + 20           // 30
5 * (10 - 3)      // 35
--10              // 10 (double negation)
2**8              // 256 (power)
```

## Conditional Logic
```
x > 10 && y < 20  // true if x > 10 and y < 20
a == 1 || b == 2  // true if a equals 1 or b equals 2
!!value           // Boolean conversion (true if value is truthy)
```

## Working with Arrays
```
[1, 2, 3][1]      // 2 (second element)
len([1, 2, 3])    // 3
```

## String Indexing
```
"world"[0]        // "w"
"world"[10]       // null (out of bounds)
words |map: $1[0]  // take first character of each word
```

## Bitwise Operations
```
5 ^ 3             // 6 (XOR: 101 ^ 011 = 110)
7 & 3             // 3 (AND: 111 & 011 = 011)
~5                // -6 (NOT: bitwise complement)
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

!!user.email && user.isActive
// Check if user has email (truthy) and is active

--score + bonus
// Double negation of score plus bonus

2**3**2
// Right-associative power: 2**(3**2) = 512

area = 3.14 * radius**2
// Calculate circle area using power operator
```