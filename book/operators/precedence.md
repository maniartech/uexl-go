# Operator Precedence

Understanding operator precedence is crucial for writing correct and predictable expressions in UExL. Operators with higher precedence are evaluated before those with lower precedence. Associativity determines the order in which operators of the same precedence are evaluated.

## Precedence Table
| Precedence | Operator(s) | Description | Associativity |
|------------|-------------|-------------|---------------|
| 1 | `()` | Grouping | Left-to-right |
| 2 | `.` `[]` `?.` `?[ ]` | Property/Index/Optional access (objects, arrays, strings) | Left-to-right |
| 3 | `**` `^` | Exponentiation/Power | Right-to-left |
| 4 | `-` (unary), `!`, `~` (unary) | Negation, Logical NOT, Bitwise NOT | Right-to-left |
| 5 | `%` | Modulo | Left-to-right |
| 6 | `*` `/` | Multiplication, Division | Left-to-right |
| 7 | `+` `-` | Addition, Subtraction | Left-to-right |
| 8 | `<<` `>>` | Bitwise Shift | Left-to-right |
| 9 | `??` | Nullish coalescing | Left-to-right |
| 10 | `<` `>` `<=` `>=` | Comparison | Left-to-right |
| 11 | `==` `!=` `<>` | Equality (exact for primitives; deep for arrays/objects) | Left-to-right |
| 12 | `&` | Bitwise AND | Left-to-right |
| 13 | `~` | Bitwise XOR | Left-to-right |
| 14 | `|` | Bitwise OR | Left-to-right |
| 15 | `&&` | Logical AND | Left-to-right |
| 16 | `||` | Logical OR | Left-to-right |
| 17 | `?:` | Conditional (ternary) | Right-to-left |
| 18 | `|:` `|map:` etc. | Pipe | Left-to-right |

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

## Power Operator (`**` and `^`)
The power operator has two forms: `**` (Python/JavaScript style) and `^` (Excel style). Both have the same precedence and are right-associative, which is important for mathematical correctness:

### Right-Associativity of Power
```
2**3**2      // Evaluates as 2**(3**2) = 2**9 = 512
2^3^2        // Evaluates as 2^(3^2) = 2^9 = 512
             // NOT as (2**3)**2 = 8**2 = 64
```

### Power vs Other Operators
The power operator has higher precedence than multiplication and also higher precedence than unary operators:
```
2*3**2       // Evaluates as 2*(3**2) = 2*9 = 18
2*3^2        // Evaluates as 2*(3^2) = 2*9 = 18
-2**3        // Evaluates as -(2**3) = -8
-2^3         // Evaluates as -(2^3) = -8
(-2)**3      // Evaluates as (-2)**3 = -8
(-2)^3       // Evaluates as (-2)^3 = -8
```

## Bitwise vs Logical Operations
Note the distinction between bitwise and logical operators:
```
5 ~ 3        // Bitwise XOR: 6 (binary: 101 ~ 011 = 110)
5 ** 3       // Power: 125 (5 to the power of 3)
5 ^ 3        // Power: 125 (Excel style)
true && false // Logical AND: false
5 & 3        // Bitwise AND: 1 (binary: 101 & 011 = 001)
~5           // Bitwise NOT: -6 (binary: ~0101 = 1010 in two's complement)

## Equality notes
- Only `==` and `!=`/`<>` exist; there is no `===`/`!==`.
- Both `!=` and `<>` are equivalent not-equals operators (C/Python/JS style vs Excel style).
- Equality is exact for primitives (no cross-type coercion) and deep for arrays/objects.
- Use `!!x` to convert any value to a boolean via truthiness.
- Use `??` and `?.` for nullish flows; equality does not treat "missing" specially.

## Nullish Coalescing (??)
`a ?? b` evaluates to `a` if it is not nullish; otherwise it evaluates to `b`. Use it to provide defaults only for missing values, not for all falsy values.

Precedence placement in UExL:
- `??` binds tighter than comparisons, equality, bitwise, and logical operators.
- `??` is looser than arithmetic (`%`, `*`, `/`, `+`, `-`) and bitwise shifts (`<<`, `>>`).
- `??` is left-associative: `a ?? b ?? c` → `(a ?? b) ?? c`.

- `0 ?? 10` → `0` (keeps valid falsy)
- `"" ?? "(empty)"` → `""`
- `null ?? 10` → `10`

When combining `??` with `&&` or `||`, UExL parses without requiring parentheses because `??` binds tighter. These are equivalent:

```
a || b ?? c   // same as a || (b ?? c)
a && b ?? d   // same as a && (b ?? d)
```

Parentheses are still recommended for readability in complex expressions.
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
2^3^2            // 512: right-associative power 2^(3^2)
2*3**2           // 18: power before multiplication 2*(3**2)
2*3^2            // 18: power before multiplication 2*(3^2)
5 ~ 3            // 6: bitwise XOR
~5               // -6: bitwise NOT
5 <> 3           // true: not equals (Excel style)
5 != 3           // true: not equals (C/Python/JS style)
"hello"[1]       // "e": string index access
0 || 10          // 10: logical OR replaces falsy 0
0 ?? 10          // 0: nullish keeps valid falsy 0
a || (b ?? c)    // ?? binds tighter than ||; parentheses for clarity only
```

Refer to this table when constructing complex expressions to ensure correct evaluation order.