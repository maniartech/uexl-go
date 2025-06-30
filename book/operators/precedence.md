# Operator Precedence

Understanding operator precedence is crucial for writing correct and predictable expressions in UExL. The following table summarizes the precedence and associativity of operators supported by UExL.

| Precedence | Operator(s) | Description | Associativity |
|------------|-------------|-------------|---------------|
| 1 | `()` | Grouping | Left-to-right |
| 2 | `.` `[]` | Property/Index access | Left-to-right |
| 3 | `-` (unary), `!` | Negation, Logical NOT | Right-to-left |
| 4 | `%` | Modulo | Left-to-right |
| 5 | `*` `/` | Multiplication, Division | Left-to-right |
| 6 | `+` `-` | Addition, Subtraction | Left-to-right |
| 7 | `<<` `>>` | Bitwise Shift | Left-to-right |
| 8 | `<` `>` `<=` `>=` | Comparison | Left-to-right |
| 9 | `==` `!=` | Equality | Left-to-right |
| 10 | `&` | Bitwise AND | Left-to-right |
| 11 | `^` | Bitwise XOR | Left-to-right |
| 12 | `|` | Bitwise OR | Left-to-right |
| 13 | `&&` | Logical AND | Left-to-right |
| 14 | `||` | Logical OR | Left-to-right |
| 15 | `|:` `|map:` etc. | Pipe | Left-to-right |

Refer to this table when constructing complex expressions to ensure correct evaluation order.