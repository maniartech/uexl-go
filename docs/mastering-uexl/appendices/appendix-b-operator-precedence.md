# Appendix B: Operator Precedence Table

Operators are listed from **highest** (16) to **lowest** (1) precedence. Operators at the same level have equal precedence and associate left-to-right unless noted.

---

| Level | Operator(s) | Name | Associativity | Example |
|-------|------------|------|---------------|---------|
| 16 | `.` `?.` `[]` `?.[]` | Member access / index / optional chaining | Left | `a.b.c`, `a?.b`, `arr[0]`, `arr?.[0]` |
| 15 | `-x` `!x` `~x` | Unary minus, logical NOT, bitwise NOT | Right (prefix) | `-a`, `!flag`, `~bits` |
| 14 | `**` and `^` | Exponentiation (both operators equivalent) | **Right** | `2**10`, `2^10` |
| 13 | `*` `/` `%` | Multiplication, division, modulo | Left | `a * b / c` |
| 12 | `+` `-` | Addition, subtraction | Left | `a + b - c` |
| 11 | `<<` `>>` | Bitwise left shift, right shift | Left | `a << 2`, `b >> 1` |
| 10 | `??` | Nullish coalescing | Left | `a ?? b ?? c` |
| 9 | `<` `>` `<=` `>=` | Relational comparison | Left | `a < b`, `score >= 90` |
| 8 | `==` `!=` `<>` | Equality (`<>` is alias for `!=`) | Left | `a == b`, `a != b`, `a <> b` |
| 7 | `&` | Bitwise AND | Left | `flags & 0xFF` |
| 6 | `~` | Bitwise XOR | Left | `a ~ b` |
| 5 | `\|` | Bitwise OR | Left | `a \| b` |
| 4 | `&&` | Logical AND (short-circuit) | Left | `a && b` |
| 3 | `\|\|` | Logical OR (short-circuit) | Left | `a \|\| b` |
| 2 | `? :` | Ternary conditional | **Right** | `a ? b : c ? d : e` |
| 1 | `\|map:` `\|filter:` etc. | Pipe operators | Left (chain) | `arr \|map: $item * 2 \|filter: $item > 3` |

---

## Notes

### Exponentiation (`**` / `^`)
Both `**` and `^` produce exponentiation. `^` is provided for Excel compatibility. Note that `~` is the XOR operator in UExL (not `^`).

```uexl
2 ** 8    # 256.0
2 ^ 8     # also 256.0
```

### Bitwise XOR (`~`)
In UExL, `~` serves dual purpose:
- **Unary `~x`**: bitwise NOT (flip all bits)
- **Binary `a ~ b`**: bitwise XOR

```uexl
~0        # -1 (bitwise NOT)
5 ~ 3     # 6  (bitwise XOR: 101 XOR 011 = 110)
```

### Nullish coalescing (`??`) tighter than comparison
`??` binds more tightly than comparison operators:

```uexl
a ?? b < c      # parsed as: (a ?? b) < c
a < b ?? c      # parsed as: a < (b ?? c)
```

Parenthesize explicitly when mixing `??` with comparisons to avoid confusion.

### Short-circuit evaluation
`&&` and `||` do not evaluate the right operand if the left operand determines the result:

```uexl
false && expensive()   # expensive() is never called
true  || expensive()   # expensive() is never called
```

### Pipe operator precedence
The pipe operators (`|map:`, `|filter:`, etc.) have the lowest precedence of all. The entire left-hand expression becomes the pipe input:

```uexl
a + b |map: $item * 2   # (a + b) is the input to |map:
```

### Ternary is right-associative
Nested ternaries associate from right to left:

```uexl
a ? b : c ? d : e
# is parsed as:
a ? b : (c ? d : e)
```

---

## Precedence Examples

```uexl
# Arithmetic binds tighter than comparison
1 + 2 > 2 + 0         # (1 + 2) > (2 + 0) → 3 > 2 → true

# Comparison binds tighter than logical AND
a > 0 && b > 0        # (a > 0) && (b > 0)

# Nullish coalescing tighter than comparison
score ?? 0 >= 60      # (score ?? 0) >= 60

# Unary binds very tight
-2 ** 2               # (-2) ** 2 = 4  (UExL applies unary before power)
                      # (this differs from Python where it would be -(2**2) = -4)

# Optional chaining with nullish fallback
profile?.rating ?? 0  # (profile?.rating) ?? 0
```
