# Bitwise Operator Research Across Languages

## Problem Statement
We need to find bitwise operators that:
1. Don't conflict with Excel operators (`^` → power, `&` → concat)
2. Are intuitive and have precedent in other languages
3. Are easy to type and read
4. Don't conflict with existing UExL syntax (`|:` pipes, etc.)

## Language Survey

### C-Family (C, C++, Java, C#, Go, Rust)
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Conflicts with our Excel compatibility needs

### Python
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Same as C-family

### Ruby
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Same as C-family

### JavaScript / TypeScript
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
USHR: >>> (unsigned right shift)
```
**Problem:** Same as C-family

### PHP
```
AND:  &  (or 'and' for logical)
OR:   |  (or 'or' for logical)
XOR:  ^  (or 'xor' for logical)
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Same as C-family

### Perl
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Same as C-family

### VBA (Visual Basic for Applications - Excel Macros!)
```
AND:  And (keyword)
OR:   Or (keyword)
XOR:  Xor (keyword)
NOT:  Not (keyword)
SHL:  (no operator, use functions)
SHR:  (no operator, use functions)
```
**Interesting:** Keywords instead of symbols!

### Visual Basic .NET
```
AND:  And (keyword)
OR:   Or (keyword)
XOR:  Xor (keyword)
NOT:  Not (keyword)
SHL:  <<
SHR:  >>
```
**Interesting:** Keywords for bitwise, symbols for shifts

### Pascal / Delphi
```
AND:  and (keyword)
OR:   or (keyword)
XOR:  xor (keyword)
NOT:  not (keyword)
SHL:  shl (keyword)
SHR:  shr (keyword)
```
**Very interesting:** All keywords!

### Lua
```
AND:  &  (Lua 5.3+)
OR:   |
XOR:  ~  (NOT is 'not' keyword)
NOT:  ~ (bitwise NOT)
SHL:  <<
SHR:  >>
```
**Interesting:** `~` is XOR in Lua 5.3+, NOT is a keyword

### Kotlin
```
AND:  and (infix function)
OR:   or (infix function)
XOR:  xor (infix function)
NOT:  inv() (function)
SHL:  shl (infix function)
SHR:  shr (infix function)
USHR: ushr (infix function)
```
**Very interesting:** All keywords/infix functions!

### Swift
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
```
**Problem:** Same as C-family

### R (Statistical Language)
```
AND:  &  (vectorized bitAnd for bitwise)
OR:   |  (vectorized bitOr for bitwise)
XOR:  xor() (function, vectorized bitXor for bitwise)
NOT:  !  (bitwNot() function for bitwise)
SHL:  bitwShiftL() (function)
SHR:  bitwShiftR() (function)
```
**Interesting:** Separates logical (&, |) from bitwise (functions)

### MATLAB
```
AND:  &  (logical), bitand() (bitwise)
OR:   |  (logical), bitor() (bitwise)
XOR:  xor() (logical), bitxor() (bitwise)
NOT:  ~  (logical), bitcmp() (bitwise)
SHL:  bitshift(x, n) (function)
SHR:  bitshift(x, -n) (function)
```
**Interesting:** Explicit separation between logical and bitwise

### SQL (Various Dialects)
```
AND:  &  (some dialects)
OR:   |  (some dialects)
XOR:  ^  (some dialects)
NOT:  ~  (some dialects)
SHL:  <<  (some dialects)
SHR:  >>  (some dialects)
```
**Note:** SQL varies widely, some use keywords

### Haskell
```
AND:  .&.  (Data.Bits)
OR:   .|.  (Data.Bits)
XOR:  xor  (Data.Bits)
NOT:  complement (Data.Bits)
SHL:  shiftL (Data.Bits)
SHR:  shiftR (Data.Bits)
```
**VERY INTERESTING:** Uses `.&.` and `.|.` notation!

### OCaml
```
AND:  land (keyword)
OR:   lor (keyword)
XOR:  lxor (keyword)
NOT:  lnot (keyword)
SHL:  lsl (keyword)
SHR:  lsr (keyword, logical shift)
ASR:  asr (keyword, arithmetic shift)
```
**Interesting:** Prefix "l" for "logical" (bitwise) operations

### Erlang
```
AND:  band (keyword)
OR:   bor (keyword)
XOR:  bxor (keyword)
NOT:  bnot (keyword)
SHL:  bsl (keyword)
SHR:  bsr (keyword)
```
**Interesting:** Prefix "b" for "bitwise" operations

### F#
```
AND:  &&&  (triple ampersand!)
OR:   |||  (triple pipe!)
XOR:  ^^^  (triple caret!)
NOT:  ~~~  (triple tilde!)
SHL:  <<<  (triple left angle!)
SHR:  >>>  (triple right angle!)
```
**VERY INTERESTING:** Triple character operators!

### Scala
```
AND:  &
OR:   |
XOR:  ^
NOT:  ~
SHL:  <<
SHR:  >>
USHR: >>>
```
**Problem:** Same as C-family

## Summary of Unique Approaches

### 1. **Keywords (VBA, Pascal, Kotlin, OCaml, Erlang)**
```
5 and 3    // Bitwise AND
5 or 7     // Bitwise OR
5 xor 3    // Bitwise XOR
not 5      // Bitwise NOT
```
**Pros:**
- Clear and explicit
- Used by VBA (Excel's macro language!)
- No symbol conflicts
- Easy to search/grep

**Cons:**
- More verbose
- May conflict with logical operators in some contexts
- Requires keyword reservation

### 2. **Dot Notation (Haskell)**
```
5 .&. 3    // Bitwise AND
5 .|. 3    // Bitwise OR
5 .^. 3    // Bitwise XOR
.~. 5      // Bitwise NOT (or complement function)
```
**Pros:**
- Visual similarity to logical operators
- Clear "this is special" signal
- Mathematically inspired

**Cons:**
- Longer to type
- Unusual syntax for most programmers

### 3. **Triple Characters (F#)**
```
5 &&& 3    // Bitwise AND
5 ||| 3    // Bitwise OR
5 ^^^ 3    // Bitwise XOR
~~~ 5      // Bitwise NOT
5 <<< 2    // Shift left
5 >>> 1    // Shift right
```
**Pros:**
- Clear relationship to logical operators
- No conflicts
- Easy to type (just hold the key)

**Cons:**
- Unusual appearance
- May look like typos

### 4. **Prefix Keywords (OCaml: land/lor, Erlang: band/bor)**
```
5 band 3   // Bitwise AND (Erlang style: "bitwise and")
5 bor 7    // Bitwise OR
5 bxor 3   // Bitwise XOR
bnot 5     // Bitwise NOT
```
**Pros:**
- Extremely clear intent
- "b" prefix = "bitwise"
- No conflicts

**Cons:**
- More verbose
- Less common pattern

## Recommendations for UExL

### ✅ FINAL DECISION: Lua-Style Context-Dependent `~` ⭐

**Chosen approach:**
```uexl
// Bitwise operators (minimal changes from current UExL)
5 & 3      // Bitwise AND = 1 (NO CHANGE)
5 | 3      // Bitwise OR = 7 (NO CHANGE)
5 ~ 3      // Bitwise XOR = 6 (CHANGED from ^)
~5         // Bitwise NOT = -6 (IMPLEMENT - was broken)
5 << 2     // Shift left = 20 (NO CHANGE)
5 >> 1     // Shift right = 2 (NO CHANGE)

// Power operator (BREAKING CHANGE)
2 ^ 3      // Power = 8 (CHANGED from XOR to match Excel)
2 ** 3     // Power = 8 (NO CHANGE - still works)

// String concatenation
"A" + "B"  // Concat = "AB" (NO CHANGE - use + not &)
```

**Why this is best:**
1. ✅ **Minimal breaking changes:** Only `^` changes meaning
2. ✅ **Lua 5.3+ precedent:** Proven design for context-dependent `~`
3. ✅ **Familiar to C-family developers:** Keep `&`, `|`, `~` symbols
4. ✅ **Excel compatibility:** `^` becomes power (critical requirement)
5. ✅ **Concise:** Single-character operators
6. ✅ **No keyword explosion:** Avoid `band`, `bor`, `bxor` verbosity
7. ✅ **Parser precedent:** Like `-` (unary negate vs binary minus)

### Option 1: Keywords (VBA-style) - NOT CHOSEN
```uexl
5 band 3   // "bitwise and" - explicit and clear
5 bor 7    // "bitwise or"
5 bxor 3   // "bitwise xor"
bnot 5     // "bitwise not"
5 << 2     // Shift left (keep symbol)
5 >> 1     // Shift right (keep symbol)
```

**Why this is best:**
1. ✅ **Excel connection:** VBA (Excel's language) uses keywords too
2. ✅ **Explicit:** Aligns with UExL philosophy (explicit > implicit)
3. ✅ **Clear intent:** No confusion about "is this bitwise or logical?"
4. ✅ **Precedent:** Erlang, OCaml, VBA all use this pattern
5. ✅ **Easy to teach:** "Use `band` for bitwise AND, `&&` for logical AND"
6. ✅ **No conflicts:** Keywords don't clash with Excel operators

### Option 2: F#-style Triple Characters
```uexl
5 &&& 3    // Bitwise AND
5 ||| 3    // Bitwise OR
5 ^^^ 3    // Bitwise XOR
~~~5       // Bitwise NOT
5 <<< 2    // Shift left
5 >>> 1    // Shift right
```

**Why this works:**
- Clear relationship to logical operators
- No conflicts with Excel
- Shorter than keywords

**Why it's second choice:**
- Less intuitive for beginners
- Looks unusual

### Option 3: Haskell-style Dot Notation
```uexl
5 .&. 3    // Bitwise AND
5 .|. 3    // Bitwise OR
5 .^. 3    // Bitwise XOR
.~. 5      // Bitwise NOT
```

**Why it's third choice:**
- Longer to type
- Less familiar pattern
- Could conflict with future method call syntax

## Final Implementation: Lua-Style Context-Dependent `~`

```uexl
// LOGICAL (Boolean operations)
true && false    // Logical AND
true || false    // Logical OR
!true            // Logical NOT

// BITWISE (Integer bit operations) - Minimal changes!
5 & 3            // Bitwise AND = 1 (NO CHANGE)
5 | 3            // Bitwise OR = 7 (NO CHANGE)
5 ~ 3            // Bitwise XOR = 6 (CHANGED from ^)
~5               // Bitwise NOT = -6 (IMPLEMENT - was broken)
5 << 2           // Bitwise shift left = 20 (NO CHANGE)
5 >> 1           // Bitwise shift right = 2 (NO CHANGE)

// EXCEL-COMPATIBLE
2 ^ 3            // Power = 8 (CHANGED from XOR)
2 ** 3           // Power = 8 (NO CHANGE - still works)
"A" + "B"        // String concat = "AB" (NO CHANGE)
x <> y           // Not equals (NEW Excel alias)
x != y           // Not equals (NO CHANGE)
```

**This achieves:**
- ✅ Excel `^` for power (critical requirement)
- ✅ Minimal breaking changes (only `^` operator)
- ✅ Familiar C-family symbols (`&`, `|`, `~`)
- ✅ Lua 5.3+ precedent for context-dependent `~`
- ✅ Concise syntax (no keyword explosion)
- ✅ Parser precedent (like `-` for unary/binary)
- ✅ No conflicts with pipe syntax (`|:` has `:` suffix)

## Alternative: Shorter Keywords

If `band`/`bor`/`bxor` feel too long, could use:

```uexl
5 bit_and 3      // More readable?
5 bitand 3       // One word (like Kotlin)
5 bit& 3         // Hybrid approach?
```

But I believe `band`/`bor`/`bxor` strikes the best balance of:
- Short enough (4-5 chars)
- Clear meaning
- Established precedent

## Implementation Notes

### Operator Precedence
Keep bitwise at current levels:
```
Precedence 12: & (bitwise AND) - NO CHANGE
Precedence 13: ~ (bitwise XOR) - CHANGED from ^
Precedence 14: | (bitwise OR) - NO CHANGE
Precedence 8:  <<, >> (bitwise shifts) - NO CHANGE
Precedence 3:  ^ (power) - CHANGED from precedence 13
```

### Parsing Context-Dependent `~`
```go
// In parser, detect unary vs binary context
func (p *Parser) parseUnary() Node {
    if p.match(TokenTilde) {
        // Unary NOT
        return &UnaryExpression{
            Operator: "~",
            Right: p.parseUnary(),
        }
    }
    // ...
}

func (p *Parser) parseBitwiseXor() Node {
    left := p.parseBitwiseAnd()
    for p.match(TokenTilde) {
        // Binary XOR
        right := p.parseBitwiseAnd()
        left = &BinaryExpression{
            Operator: "~",
            Left: left,
            Right: right,
        }
    }
    return left
}
```

### Compiler Mapping
```go
case "~":
    if expr.IsUnary {
        emit(OpBitwiseNot)  // Unary: bitwise NOT
    } else {
        emit(OpBitwiseXor)  // Binary: bitwise XOR
    }
case "^":
    emit(OpPow)  // Power (was OpBitwiseXor)
```

### VM Handler (New)
```go
case OpBitwiseNot:
    // Implement bitwise NOT (currently missing)
    val := vm.pop()
    intVal := int64(val.(float64))
    vm.push(float64(^intVal))
```

### Migration
```bash
# Replace ^ with ~ where it's XOR (requires manual review)
# Context-dependent: distinguish XOR from power usage
sed -i 's/\([0-9a-zA-Z_()]\) \^ \([0-9a-zA-Z_(]\)/\1 ~ \2/g' *.uexl
```
