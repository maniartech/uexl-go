# Operator Precedence

Understanding operator precedence is crucial for writing correct and predictable expressions in UExL. Operators with higher precedence are evaluated before those with lower precedence. Associativity determines the order in which operators of the same precedence are evaluated.

## Precedence Table
| Precedence | Operator(s) | Description | Associativity |
|------------|-------------|-------------|---------------|
| 1 | `()` | Grouping | Left-to-right |
| 2 | `.` `[]` | Property/Index access | Left-to-right |
| 3 | `-` (unary), `!` | Negation, Logical NOT | Right-to-left |
| 4 | `**` | Exponentiation/Power | Right-to-left |
| 5 | `%` | Modulo | Left-to-right |
| 6 | `*` `/` | Multiplication, Division | Left-to-right |
| 7 | `+` `-` | Addition, Subtraction | Left-to-right |
| 8 | `<<` `>>` | Bitwise Shift | Left-to-right |
| 9 | `<` `>` `<=` `>=` | Comparison | Left-to-right |
| 10 | `==` `!=` | Equality | Left-to-right |
| 11 | `&` | Bitwise AND | Left-to-right |
| 12 | `^` | Bitwise XOR | Left-to-right |
| 13 | `|` | Bitwise OR | Left-to-right |
| 14 | `&&` | Logical AND | Left-to-right |
| 15 | `||` | Logical OR | Left-to-right |
| 16 | `|:` `|map:` etc. | Pipe | Left-to-right |

## Associativity
- **Left-to-right**: Operators are evaluated from left to right (e.g., `a - b - c` is `(a - b) - c`).
- **Right-to-left**: Operators are evaluated from right to left (e.g., `-x` is evaluated before applying to the next operator).

## Consecutive Unary Operators
UExL supports consecutive unary operators, which are particularly useful for:

### Double Negation (`--`)
Double negation can be used to convert a number to its positive form or for mathematical operations:
```
--10      // Results in 10: -(-(10))
--x       // Results in x: -(-(x))
---5      // Results in -5: -(-(-(5)))
```

### Boolean Conversion (`!!`)
Double NOT is a common pattern for converting any value to a boolean:
```
!!1       // Results in true: !(!(1))
!!0       // Results in false: !(!(0))
!!"text"  // Results in true: !(!(true))
!!null    // Results in false: !(!(null))
```

### Multiple Consecutive Operators
You can chain any number of unary operators:
```
!!!value  // Triple NOT: !(!(!(value)))
!-x       // NOT of negation: !(-(x))
-!flag    // Negation of NOT: -(!(flag))
```

All consecutive unary operators are evaluated right-to-left, meaning the rightmost operator is applied first.

## Power Operator (`**`)
The power operator has its own precedence level and is right-associative, which is important for mathematical correctness:

### Right-Associativity of Power
```
2**3**2      // Evaluates as 2**(3**2) = 2**9 = 512
             // NOT as (2**3)**2 = 8**2 = 64
```

### Power vs Other Operators
The power operator has higher precedence than multiplication but lower than unary operators:
```
2*3**2       // Evaluates as 2*(3**2) = 2*9 = 18
-2**3        // Evaluates as -(2**3) = -8
(-2)**3      // Evaluates as (-2)**3 = -8
```

## Bitwise vs Logical Operations
Note the distinction between bitwise and logical operators:
```
5 ^ 3        // Bitwise XOR: 6 (binary: 101 ^ 011 = 110)
5 ** 3       // Power: 125 (5 to the power of 3)
true && false // Logical AND: false
5 & 3        // Bitwise AND: 1 (binary: 101 & 011 = 001)
```

## Practical Tips
- Use parentheses to make complex expressions clear and to override default precedence.
- When in doubt, add parentheses for readability and correctness.

## Examples
```
1 + 2 * 3        // 7 (multiplication before addition)
(a + b) * c      // Parentheses override precedence
x > 10 && y < 20 // Comparison before logical AND
[1, 2, 3] |map: $1 * 2 |filter: $1 > 2 // Pipes are evaluated last
--10             // 10: double negation -(-(10))
!!value          // Boolean conversion !(!(value))
!-x              // NOT of negation: !(-(x))
2**3**2          // 512: right-associative power 2**(3**2)
2*3**2           // 18: power before multiplication 2*(3**2)
5 ^ 3            // 6: bitwise XOR
```

Refer to this table when constructing complex expressions to ensure correct evaluation order.