# Chapter 2: Setting Up and Your First Expression

> "The best way to understand an expression language is to evaluate something. Let's start in thirty seconds."

---

## 2.1 The UExL Playground

Before writing a single line of Go, open the UExL playground. It runs the same evaluation engine as the Go library via WebAssembly — no setup required.

**Playground URL:** `https://playground.uexl.dev` *(or check the project README for the current URL)*

Type this into the expression box:

```uexl
10 + 20 * 3
```

You'll see the result: `70`. Notice UExL follows standard operator precedence — multiplication before addition, just like mathematics. Now try:

```uexl
(10 + 20) * 3
```

Result: `90`. Parentheses override precedence, exactly as you'd expect.

The playground is useful throughout this book for quickly testing expressions without switching to Go. We'll refer back to it in Chapter 16 when we discuss debugging strategies.

---

## 2.2 Installing the Go Library

Add UExL to your Go module:

```bash
go get github.com/maniartech/uexl-go
```

Import in your Go file:

```go
import "github.com/maniartech/uexl"
```

> **NOTE:** Make sure you're using Go 1.21 or later. UExL uses generics-adjacent patterns in its internal type handling that require a recent toolchain.

---

## 2.3 The Simplest Path: `uexl.Eval()`

For quick, one-off evaluation, use the top-level `Eval` function:

```go
package main

import (
    "fmt"
    "github.com/maniartech/uexl"
)

func main() {
    result, err := uexl.Eval("10 + 20 * 3", nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result) // 70
}
```

`Eval` takes an expression string and a context map (`nil` means empty context), and returns `(any, error)`. The return type `any` holds the computed value — a `float64` for numbers, `string` for strings, `bool` for booleans, `[]any` for arrays, `map[string]any` for objects, or `nil` for null.

### Providing context variables

Pass data into expressions through the context map:

```go
ctx := map[string]any{
    "price":    99.99,
    "quantity": 3,
    "discount": 0.10,
}

result, err := uexl.Eval("price * quantity * (1 - discount)", ctx)
// result: 269.97300000000004
```

The expression accesses `price`, `quantity`, and `discount` directly by name.

---

## 2.4 The Three-Stage Pipeline

`Eval` is convenient, but production systems almost never use it directly. The reason: `Eval` parses, compiles, **and** executes on every call. If you're evaluating the same expression against thousands of different contexts — say, applying a pricing rule to every product in a catalog — you're doing redundant work.

UExL's architecture separates these three stages explicitly:

```
Expression String  ──►  Parser  ──►  AST  ──►  Compiler  ──►  ByteCode  ──►  VM  ──►  Result
```

**Stage 1 — Parser:** Tokenizes and parses the expression string into an Abstract Syntax Tree (AST). This is the most expensive step linguistically — it checks syntax and builds the tree.

**Stage 2 — Compiler:** Transforms the AST into bytecode — a compact sequence of instructions optimized for fast execution. Constants are extracted and stored separately.

**Stage 3 — VM (Virtual Machine):** Executes the bytecode against a given context. This is a tight loop over a small instruction set — very fast.

The key insight: **Stages 1 and 2 only depend on the expression, not the data.** They can run once at startup. Stage 3 is the only part that touches runtime data.

### Using the three-stage API

```go
package main

import (
    "fmt"
    "log"

    "github.com/maniartech/uexl/compiler"
    "github.com/maniartech/uexl/parser"
    "github.com/maniartech/uexl/vm"
)

func main() {
    expr := "price * quantity * (1 - discount)"

    // Stage 1: Parse
    ast, err := parser.ParseString(expr)
    if err != nil {
        log.Fatal("Parse error:", err)
    }

    // Stage 2: Compile
    comp := compiler.New()
    if err := comp.Compile(ast); err != nil {
        log.Fatal("Compile error:", err)
    }
    bytecode := comp.ByteCode()

    // Stage 3: Execute (can be called many times with different contexts)
    machine := vm.New(vm.LibContext{
        Functions:    vm.Builtins,
        PipeHandlers: vm.DefaultPipeHandlers,
    })

    orders := []map[string]any{
        {"price": 99.99, "quantity": 2, "discount": 0.10},
        {"price": 49.99, "quantity": 5, "discount": 0.05},
        {"price": 199.99, "quantity": 1, "discount": 0.20},
    }

    for _, order := range orders {
        result, err := machine.Run(bytecode, order)
        if err != nil {
            log.Println("Runtime error:", err)
            continue
        }
        fmt.Printf("Order total: %.2f\n", result)
    }
}
```

Output:
```
Order total: 179.98
Order total: 237.45
Order total: 159.99
```

The expression is parsed and compiled **once**. The VM executes it **three times** — once per order — each with a different context but the same bytecode.

> **TIP:** In benchmarks, the compile-once/run-many pattern is typically 5–20× faster than calling `Eval` in a loop for large batches, because parsing and compilation have non-trivial fixed cost per expression.

---

## 2.5 Anatomy of a UExL Expression

Let's look at a slightly richer expression and identify its parts:

```uexl
customer.tier == 'gold' ? price * 0.85 : price
```

| Part | Type | Role |
|------|------|------|
| `customer.tier` | Property access | Reads `tier` from the `customer` context object |
| `==` | Equality operator | Compares two values |
| `'gold'` | String literal | A literal string value |
| `? ... : ...` | Ternary operator | Conditional branch |
| `price * 0.85` | Arithmetic | 15% discounted price |
| `price` | Identifier | The base price from context |

Expressions are composed of: **literals** (raw values), **identifiers** (context lookups), **operators** (compute relationships), **function calls**, and **pipe stages**. We'll cover each in the coming chapters.

---

## 2.6 Your First ShopLogic Expression

Let's write the first expression for the ShopLogic project. The scenario: calculate the effective price for a product, applying a tier-based discount.

Business rule: *"Gold customers get 15% off, Silver customers get 8% off, everyone else pays full price."*

```uexl
product.basePrice * (1 - (
    customer.tier == 'gold'   ? 0.15 :
    customer.tier == 'silver' ? 0.08 :
    0
))
```

Test it in Go:

```go
ctx := map[string]any{
    "product":  map[string]any{"basePrice": 120.00},
    "customer": map[string]any{"tier": "gold"},
}

result, err := uexl.Eval(`
    product.basePrice * (1 - (
        customer.tier == 'gold'   ? 0.15 :
        customer.tier == 'silver' ? 0.08 :
        0
    ))
`, ctx)
// result: 102.0 (15% off 120)
```

This expression reads naturally: multiply the base price by one minus the applicable discount fraction. The nested ternary handles the tier hierarchy cleanly.

We'll return to this expression in Chapter 5 when we explore operators, and in Chapter 15 when we discuss how to design the `product` and `customer` context objects properly.

---

## 2.7 Error Handling from the Start

Every call in the three-stage API returns an error. Never discard them. Here's what each stage's errors tell you:

```go
// Stage 1 parse error — syntax problem in the expression string
ast, err := parser.ParseString("product.price ?? ")
// err: "unexpected end of input at position 18"
// Type: *errors.ParserError — contains position information

// Stage 2 compile error — structural problem that passes the parser
// (rare; usually indicates a parser bug or unsupported AST node)
if err := comp.Compile(ast); err != nil {
    // err: plain error string
}

// Stage 3 runtime error — type problem or missing data
result, err := machine.Run(bytecode, ctx)
// err: "TypeError: cannot multiply string and number"
// err: "ReferenceError: 'discount' is not defined"
```

A common pattern for user-facing APIs is to validate expressions at authoring time (parse + compile only, no execution) and return structured errors to the UI:

```go
func ValidateExpression(expr string) error {
    ast, err := parser.ParseString(expr)
    if err != nil {
        return fmt.Errorf("syntax error: %w", err)
    }
    comp := compiler.New()
    if err := comp.Compile(ast); err != nil {
        return fmt.Errorf("compilation error: %w", err)
    }
    return nil // expression is structurally valid
}
```

We'll build out a full validation API in Chapter 16.

---

## 2.8 Summary

- The UExL playground lets you evaluate expressions instantly in the browser — use it for learning and debugging.
- `uexl.Eval()` is the simplest integration path: parse, compile, and execute in one call.
- The three-stage API (`parser.ParseString` → `compiler.New().Compile()` → `vm.New().Run()`) is the production pattern: compile once, execute with different contexts as many times as needed.
- Expressions are composed of literals, identifiers, operators, function calls, and pipe stages — we cover each in upcoming chapters.
- Always handle errors from every stage; UExL never panics on malformed expressions.

---

## Exercises

**2.1 — Recall.** What are the three stages of the UExL pipeline? Which stage is the most expensive and why? Which stage is called multiple times per request in the compile-once/run-many pattern?

**2.2 — Apply.** Using `uexl.Eval()`, write a Go program that evaluates the expression `len(name) > 0 ? name : 'Anonymous'` with the context `{"name": "alice"}`. What is the expected result? What happens when `name` is `""` (empty string)?

**2.3 — Extend.** For the ShopLogic project, you need to add a new tier: "Platinum" customers get 20% off. Update the discount expression from Section 2.6 to include this tier (Platinum should be checked first). Then convert the expression to use the three-stage API so it can be compiled once and run for a list of five test customers.
