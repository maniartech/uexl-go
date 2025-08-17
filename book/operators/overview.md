# Operators

Operators are the tools that let you combine, compare, and transform values in UExL. In this chapter, you'll learn about the different types of operators, how to use them, and see practical examples for each.

## Why Operators Matter
Operators allow you to perform calculations, make decisions, and manipulate data. Mastering operators is essential for writing expressive and efficient expressions.

## Categories of Operators
- **Primary:** Used for grouping, property access, and indexing.
  - Example: `(a + b)`, `user.name`, `arr[0]`
- **Unary:** Operate on a single value, such as negation or logical NOT.
  - Example: `-x`, `!flag`, `--x` (double negation), `!!value` (boolean conversion)
- **Arithmetic:** Perform mathematical calculations.
  - Example: `a + b`, `x - 1`, `y * 2`, `z / 3`, `n % 2`, `base**exponent`
- **Bitwise:** Manipulate individual bits in numbers.
  - Example: `a & b`, `x | y`, `n ^ m`, `~value`, `a << 2`, `b >> 1`
- **Comparison:** Compare values for ordering.
  - Example: `x < 10`, `score >= 50`
- **Equality:** Test if values are equal or not.
  - Example: `a == b`, `x != y`
- **Logical:** Combine boolean values.
  - Example: `x && y`, `a || b`
  - Nullish coalescing: `a ?? b` keeps valid falsy values and only falls back for nullish
- **Pipe:** Chain and transform data.
  - Example: `|:`, `|map:`, `|filter:`, `|reduce:`

## Practical Examples
Here are some ways you can use operators in UExL:
```
(a + b) * c > 100 && isActive
[1, 2, 3] |map: $1 * 2 |filter: $1 > 2
user.age >= 18 ? "adult" : "minor"
flag = !isDisabled
--x      // Double negation: -(-(x))
!!value  // Boolean conversion: !(!(value))
2**3     // Power operation: 8
5 ^ 3    // XOR operation: 6
0 || 10  // Logical OR: 10 (replaces falsy 0)
0 ?? 10  // Nullish: 0 (keeps valid falsy)
```

## Tips for Using Operators
- Use parentheses to control the order of operations.
- Combine different operators to create complex logic.
- Remember that operator precedence affects how expressions are evaluated (see the next chapter).

## Practice: Try It Yourself
Experiment with these expressions:
```
5 * (2 + 3)
x > 10 && y < 20
[10, 20, 30] |filter: $1 > 15
user.name == "Alice"
--10         // Double negation
!!false      // Boolean conversion
!!!value     // Triple NOT
2**8         // Power: 256
7 ^ 3        // XOR: 4
```

Understanding operators unlocks the full power of UExL expressions. In the next chapter, we'll dive into operator precedence and associativity.

See the next section for operator precedence and details.