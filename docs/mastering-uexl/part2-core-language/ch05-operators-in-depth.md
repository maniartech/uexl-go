# Chapter 5: Operators In Depth

> "Operators are the verbs of an expression language. Learn them well and expressions write themselves."

---

## 5.1 Overview

UExL provides a rich set of operators organized in a clear precedence hierarchy. Unlike some languages where operator behavior surprises you, UExL follows mathematical conventions and makes explicit choices where ambiguity historically causes bugs.

This chapter covers every operator — its syntax, semantics, precedence, and the sharp corners you need to know.

---

## 5.2 Operator Precedence Table

Precedence determines which operators bind more tightly when multiple operators appear in an expression without parentheses. Higher precedence = binds first.

| Level | Operators | Associativity | Category |
|-------|-----------|---------------|----------|
| 1 (highest) | `!`, unary `-`, `~` (unary NOT) | Right | Prefix unary |
| 2 | `**`, `^` (power) | Right | Power |
| 3 | `*`, `/`, `%` | Left | Multiplicative |
| 4 | `+`, `-` | Left | Additive |
| 5 | `<`, `<=`, `>`, `>=` | Left | Relational |
| 6 | `==`, `!=`, `<>` | Left | Equality |
| 7 | `&` (bitwise AND) | Left | Bitwise AND |
| 8 | `~` (binary XOR) | Left | Bitwise XOR |
| 9 | `\|` (bitwise OR) | Left | Bitwise OR |
| 10 | `&&` | Left | Logical AND |
| 11 | `\|\|` | Left | Logical OR |
| 12 | `??` | Left | Nullish coalescing |
| 13 (lowest) | `? :` | Right | Ternary |

Access operators (`.`, `[]`, `?.`, `?[]`) bind above all of these — they are resolved first in any expression.

When in doubt, use parentheses. Parentheses always override precedence.

---

## 5.3 Arithmetic Operators

```uexl
10 + 3      // => 13
10 - 3      // => 7
10 * 3      // => 30
10 / 3      // => 3.3333333333333335
10 % 3      // => 1  (modulo/remainder)
```

### The `+` operator and type rules

`+` is addition for numbers and concatenation for strings. Using it with mixed types is a TypeError:

```uexl
5 + 3           // => 8     (number + number)
"hello" + " world"  // => "hello world"  (string + string)
5 + "3"         // TypeError — no implicit coercion
```

### Division always returns a float

```uexl
10 / 2      // => 5     (not 5.0, but stored as float64)
10 / 3      // => 3.3333333333333335
7 / 2       // => 3.5
```

For integer (truncating) division, there is no built-in `floor()` in UExL. Register a host function (Chapter 14) or use the exact `7 / 2 == 3.5` result directly:

```uexl
// With a host-provided floor() registered in LibContext:
floor(7 / 2)    // => 3
```

### Modulo behaves as remainder (not mathematical modulo for negatives)

```uexl
10 % 3      // => 1
-10 % 3     // => -1   (sign follows dividend, like Go)
10 % -3     // => 1
```

### Division by zero is an error

```uexl
10 / 0      // RuntimeError: division by zero
10 % 0      // RuntimeError: division by zero
```

---

## 5.4 Power Operators

UExL supports two power syntaxes for familiarity across backgrounds:

```uexl
2 ** 8    // => 256  (Python/JavaScript style)
2 ^ 8     // => 256  (Excel style)
```

Both operators are **right-associative**, matching mathematical convention:

```uexl
2 ** 3 ** 2    // => 512  (parsed as 2 ** (3 ** 2) = 2 ** 9 = 512)
(2 ** 3) ** 2  // => 64   (explicit left grouping)
```

> **WARNING:** `^` in most C-family languages is bitwise XOR. In UExL it is power (like Excel). The bitwise XOR role is taken by `~` (the tilde). This design choice was intentional to match Excel users' expectations.

```uexl
2 ^ 3     // => 8    (power — NOT XOR!)
2 ~ 3     // => 1    (XOR: 010 XOR 011 = 001)
```

---

## 5.5 Comparison Operators

```uexl
5 < 10      // => true
5 <= 5      // => true
5 > 10      // => false
5 >= 6      // => false
5 == 5      // => true
5 != 5      // => false
5 <> 5      // => false  (<> is an alias for != — Excel style)
```

### Deep equality for arrays and objects

`==` and `!=` perform **deep structural equality** for arrays and objects:

```uexl
[1, 2, 3] == [1, 2, 3]         // => true  (same structure and values)
[1, 2, 3] == [1, 2, 4]         // => false
{a: 1, b: 2} == {b: 2, a: 1}   // => true  (order doesn't matter for objects)
{a: 1} == {a: 1, b: 2}         // => false (different number of keys)
```

### No implicit coercion in equality

```uexl
5 == "5"        // => false  (number vs string — no coercion)
0 == false      // => false  (number vs boolean)
null == false   // => false
null == 0       // => false
```

This is an intentional departure from JavaScript's `==`. UExL has no loose equality operator — `==` is always strict for primitives. If you need to compare a value that might be a string from a JSON source, ensure it is already the correct type on the Go side before passing to the expression, or register a host conversion function (Chapter 14).

### Comparison with null

Comparing any value to `null` with `<`, `<=`, `>`, `>=` is a TypeError. Use `??` or `?.` for null-guarded access, then compare the non-null value.

---

## 5.6 Logical Operators

```uexl
true && false    // => false   (AND)
true || false    // => true    (OR)
!true            // => false   (NOT)
```

### Short-circuit evaluation

`&&` and `||` use short-circuit evaluation — the right side is only evaluated if needed:

```uexl
false && expensiveCall()    // expensiveCall() is never called
true  || expensiveCall()    // expensiveCall() is never called
```

This is not just a performance optimization — it's also a safety mechanism. You can safely use a `&&` guard to prevent evaluation of subsequent expressions that would fail:

```uexl
// Safe: arr is only indexed if arr has length > 0
len(arr) > 0 && arr[0] == targetValue
```

### Logical operators return the controlling value, not a boolean

Like JavaScript, UExL's `||` and `&&` return the *value* that determined the result, not necessarily `true` or `false`:

```uexl
0 || 42         // => 42    (0 is falsy, so right side is returned)
1 || 42         // => 1     (1 is truthy, left side returned)
"a" && "b"      // => "b"   (both truthy, right side returned)
false && "b"    // => false  (false is falsy, left side returned)
```

> **TIP:** Use `!!expr` to force any expression to a boolean `true`/`false`:
> ```uexl
> !!(0 || 42)     // => true  (!! converts 42 to true)
> !!false         // => false
> !!""            // => false
> !![]            // => false
> !!{}            // => false
> ```

---

## 5.7 The Ternary Operator

The ternary operator is UExL's conditional expression. It selects one of two values based on a condition:

```uexl
condition ? valueIfTrue : valueIfFalse
```

Examples:

```uexl
score >= 60 ? "pass" : "fail"                    // => "pass" if score >= 60
stock > 0 ? "available" : "out_of_stock"
customer.active ? "welcome back" : "account inactive"
```

### Nested ternary for multi-way conditions

```uexl
score >= 90 ? "A" :
score >= 80 ? "B" :
score >= 70 ? "C" :
score >= 60 ? "D" : "F"
```

Format multi-level ternaries across lines for readability. The ternary is right-associative, so `A ? B : C ? D : E` is `A ? B : (C ? D : E)`.

### Ternary is lazy (short-circuit)

Only the selected branch is evaluated:

```uexl
isAdmin ? loadAdminPanel() : loadUserPanel()
// Only one of the two function calls executes
```

---

## 5.8 Bitwise Operators

UExL supports bitwise operations on integer-valued numbers:

```uexl
5 & 3       // => 1   (AND:  101 & 011 = 001)
5 | 3       // => 7   (OR:   101 | 011 = 111)
5 ~ 3       // => 6   (XOR:  101 ~ 011 = 110)
~5          // => -6  (NOT: bitwise complement)
5 << 2      // => 20  (left shift: 101 << 2 = 10100)
5 >> 1      // => 2   (right shift: 101 >> 1 = 10)
```

> **WARNING:** Remember the UExL operator mapping:
> - `~` (tilde, binary) = XOR (not `^`)
> - `~` (tilde, unary) = bitwise NOT (same symbol as `~` in C/Go)
> - `^` = Power (not XOR!)
> - `|` = bitwise OR (not the pipe operator — pipes use `|type:` with a colon)

The `|` vs. pipe distinction is important enough to state clearly:

```uexl
5 | 3               // bitwise OR — produces 7
[1,2,3] |map: $item  // pipe operator — transforms the array
```

The pipe operator **always** has the form `|word:` — a pipe type keyword followed by a colon. A bare `|` is always bitwise OR.

---

## 5.9 Unary Operators

### Unary minus

```uexl
-5          // => -5
-x          // negates the value of x
-(-5)       // => 5   (equivalent to --5)
```

### Double negation

```uexl
--5     // => 5     (-(−5))
--x     // same as x for numeric values
```

Double negation is not commonly useful but is valid syntax. More useful is the boolean double-NOT pattern.

### Logical NOT and double NOT

```uexl
!true       // => false
!false      // => true
!0          // => true  (0 is falsy)
!""         // => true  (empty string is falsy)
!null       // => true  (null is falsy)
!"hello"    // => false (non-empty string is truthy)

!!true      // => true   (boolean conversion pattern)
!!0         // => false
!!""        // => false
!!"hello"   // => true
```

The `!!` double-NOT pattern converts any value to a boolean. It's the standard UExL idiom for this conversion.

### Bitwise NOT (tilde unary)

```uexl
~5      // => -6  (bitwise complement: ~n = -(n+1))
~0      // => -1
~(-1)   // => 0
```

---

## 5.10 Common Expression Patterns

### Range check

```uexl
score >= 0 && score <= 100
age >= 18 && age < 65
```

### Conditional with fallback

```uexl
price > 0 ? price : listPrice    // prefer price, fall back to listPrice
```

### Multi-condition with short-circuit guard

```uexl
user.active && user.tier == 'premium' && user.trialEnd > today
```

### Discount calculation (ShopLogic)

> **NOTE:** `min(a, b)` in the following example is a **host-provided function** registered in `LibContext.Functions` — it is not a UExL built-in. Chapter 14 shows how to register it. Without it, replace `min(a, b)` with the ternary `a <= b ? a : b`.

```uexl
product.basePrice * (1 - min(
    customer.tier == 'platinum' ? 0.20 :
    customer.tier == 'gold'     ? 0.15 :
    customer.tier == 'silver'   ? 0.08 : 0.00,
    maxDiscount
))
```

This expression:
1. Determines the raw discount rate using a nested ternary
2. Caps it with `min()` against the `maxDiscount` ceiling
3. Applies it as a multiplier against the base price

### Power in price calculations

```uexl
// Compound interest: P(1 + r)^n
principal * (1 + rate) ^ years

// Area of circle: πr²
3.14159265358979 * radius ** 2
```

---

## 5.11 Summary

- UExL operators follow a clear 13-level precedence hierarchy; use parentheses when in doubt.
- Power is `**` or `^` (not XOR); XOR is `~` (binary tilde); bitwise NOT is `~` (unary tilde).
- `+` adds numbers and concatenates strings; mixing types is a TypeError.
- `&&` and `||` short-circuit; the ternary `? :` is lazy and returns the selected branch only.
- `==` performs deep equality for arrays/objects and exact equality for primitives — no implicit coercion.
- The pipe operator is `|word:` (with a keyword and colon); bare `|` is always bitwise OR.
- `!!value` converts any value to boolean; `--x` is double negation; these are valid and occasionally useful idioms.

---

## Exercises

**5.1 — Recall.** What is the result of `2 ** 3 ** 2` in UExL? Why? What is the difference between `^` in UExL and `^` in Go? What does `~5` evaluate to?

**5.2 — Apply.** Without running the code, determine the result of each expression:
1. `10 + 3 * 2`
2. `(10 + 3) * 2`
3. `true || false && false`
4. `!!(null || 0 || "" || "hello")`
5. `5 ~ 3` (hint: XOR)
6. `3 ** 2 ** 2` (hint: right-associative)

**5.3 — Extend.** For ShopLogic, write a complete pricing expression that:
- Applies tier-based discount: Platinum 20%, Gold 15%, Silver 8%, standard 0%
- Adds a `loyaltyBonus` discount of 2% if `customer.totalSpent > 1000`
- Caps the total discount at `maxDiscount` (hint: `discount <= maxDiscount ? discount : maxDiscount`, or register a host `min()` function)
- Returns the final price as `product.basePrice * (1 - totalDiscount)`

Hint: Calculate intermediate discounts carefully using `&&` and `+`. Consider how you'd structure the ternary to avoid double-counting.
