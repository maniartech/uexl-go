# Operators

UExL provides a variety of operators for arithmetic, logical, bitwise, comparison, and pipe operations. Operators are organized by type and precedence. Understanding how each operator works is essential for writing correct expressions.

## Categories and Examples
- **Primary**: Grouping with parentheses, property/index access.
  - Examples: `(a + b)`, `user.name`, `arr[0]`
- **Unary**: Negation and logical NOT.
  - Examples: `-x`, `!flag`
- **Arithmetic**: Addition, subtraction, multiplication, division, modulo.
  - Examples: `a + b`, `x - 1`, `y * 2`, `z / 3`, `n % 2`
- **Bitwise**: AND, OR, XOR, shift left/right.
  - Examples: `a & b`, `x | y`, `n ^ m`, `a << 2`, `b >> 1`
- **Comparison**: Less than, greater than, less than or equal, greater than or equal.
  - Examples: `x < 10`, `score >= 50`
- **Equality**: Equal, not equal.
  - Examples: `a == b`, `x != y`
- **Logical**: Logical AND, logical OR.
  - Examples: `x && y`, `a || b`
- **Pipe**: Pipe and data transformation operators.
  - Examples: `|:`, `|map:`, `|filter:`, `|reduce:`

## Practical Usage
Operators can be combined to form complex expressions. Parentheses can be used to control evaluation order:
```
(a + b) * c > 100 && isActive
[1, 2, 3] |map: $1 * 2 |filter: $1 > 2
```

Refer to the next section for operator precedence and associativity.

See the next section for operator precedence and details.