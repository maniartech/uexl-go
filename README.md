# UExL (Universal Expression Language) in Golang

## Introduction

UExL (Universal Expression Language) is an embeddable, platform-independent expression evaluation engine with a unique zero-allocation architecture. Designed for efficiency and simplicity, UExL offers an intuitive approach to handling and evaluating expressions in various formats.

UExL is perfect for applications where expressions are not known at compile time, or where you need flexibility through configuration files or databases. With industry-leading performance (zero allocations, 227ns/op), comprehensive pipe operations, and Excel-compatible syntax, UExL makes dynamic expression evaluation both powerful and production-ready.

## Table of Contents

- [UExL (Universal Expression Language) in Golang](#uexl-universal-expression-language-in-golang)
  - [Introduction](#introduction)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Installing UExL](#installing-uexl)
  - [Getting Started](#getting-started)
    - [Importing the Library](#importing-the-library)
    - [Basic Usage](#basic-usage)
  - [Features](#features)
  - [Performance](#performance)
    - [Benchmark Comparison](#benchmark-comparison)
    - [Why UExL is Better](#why-uexl-is-better)
    - [Performance Details](#performance-details)
  - [Operator Precedence](#operator-precedence)
  - [Examples](#examples)

## Installation

UExL is *not yet released and ready* for use. However, whenever it is, you can
install it using the following instructions.

### Installing UExL

To install UExL, run the following command in your terminal:

```bash
go get github.com/maniartech/uexl-go
```

This command downloads and installs the UExL package along with its dependencies.

---

## Getting Started

### Importing the Library

First, include UExL in your Go project by importing it:

```go
import "github.com/maniartech/uexl-go"
```

### Basic Usage

Here’s how you can quickly start using UExL to evaluate an expression:

```go
package main

import (
    "fmt"
    "github.com/maniartech/uexl-go"
    "log"
)

func main() {
    // Evaluating a simple expression
    result, err := uexl.Eval("10 + 20")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Result: %v\n", result) // Output: Result: 30

    // Using pipe operator in expression
    result, err = uexl.Eval("10 + 20 |: $1 * 2")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Result: %v\n", result) // Output: Result: 60
}
```

This basic example demonstrates how to use UExL to evaluate simple arithmetic expressions.

## Features

- **Excel-Compatible Operators**: Supports both traditional programming syntax (`**`, `!=`) and Excel-style operators (`^` for power, `<>` for not-equals)
- **Lua-Style Bitwise Operators**: Use `~` for XOR and bitwise NOT operations
- **Flexible Syntax**: Choose operator styles based on your background (Excel, Python, JavaScript, C, Lua)
- **Pipe Operations**: Transform data using intuitive pipe syntax with operators like `|map:`, `|filter:`, `|reduce:`, and 10+ more
- **Type Safety**: Strong type checking with explicit nullish/boolish semantics and robust error handling
- **Zero Allocations**: Industry-leading performance with 0 allocations for primitive operations
- **High Performance**: Optimized bytecode VM with competitive speed (227ns/op for complex expressions)
- **Zero Panics**: Production-ready with comprehensive error handling, never crashes
- **Comprehensive Testing**: 248+ tests with race detection ensuring correctness and reliability

## Performance

UExL achieves **industry-leading performance** through its unique Value type system and zero-allocation architecture:

### Benchmark Comparison

Expression: `(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`

| Framework | Time (ns/op) | Memory (B/op) | Allocs/op | Status |
|-----------|--------------|---------------|-----------|---------|
| **UExL** | **227** | **0** | **0** | ✅ **Zero Allocations** |
| expr | 132 | 32 | 1 | Fastest raw speed |
| cel-go | 174 | 16 | 1 | Good performance |

### Why UExL is Better

1. **Zero Allocations = Better GC Behavior**
   - Only framework with 0 allocations per operation
   - No garbage collection pressure in hot paths
   - Predictable, consistent performance
   - Better for high-throughput applications

2. **Superior Pipe Performance**
   - Map operations: **3× faster** than competitors (3,428ns vs 10,588ns)
   - Optimized scope reuse and frame management
   - Efficient data transformation pipelines

3. **Explicit Semantics**
   - Clear nullish handling (`??` operator)
   - Explicit optional chaining (`?.`, `?.[`)
   - No JavaScript-style type coercion surprises
   - Predictable behavior in production

4. **Production-Ready Architecture**
   - Zero panics guarantee (never crashes)
   - Comprehensive error handling
   - Type-safe operations
   - Well-tested (248+ tests with race detection)

5. **Competitive Speed**
   - Within 72% of fastest competitor on simple expressions
   - 3× faster on complex pipe operations
   - Trade-off: slightly slower on primitives for zero allocations

### Performance Details

For detailed performance documentation, see:
- [Value Migration Architecture](designdocs/value-migration/README.md) - How zero allocations were achieved
- [Performance Optimization](designdocs/value-migration/performance-optimization.md) - Complete optimization journey
- [Future Optimizations](designdocs/performance/README.md) - Planned improvements (20-35ns target)

## Operator Precedence

| Operators | Type             | Associativity   | Notes |
|-----------|------------------|-----------------|-------|
| `(` `)`   | Parentheses      | Left to Right   | |
| `.` `[]` `?.` `?[]` | Access | Left to Right | Property/Index/Optional |
| `**` `^`  | Power            | Right to Left   | Excel: `^`, Python/JS: `**` |
| `-` `!` `~` (unary) | Unary | Right to Left | Negation, NOT, Bitwise NOT |
| `%`       | Modulus          | Left to Right   | |
| `*` `/`   | Multiplicative   | Left to Right   | |
| `+` `-`   | Additive         | Left to Right   | |
| `<<` `>>` | Bitwise Shift    | Left to Right   | |
| `??`      | Nullish Coalescing | Left to Right | |
| `<` `>` `<=` `>=` | Comparison | Left to Right | |
| `==` `!=` `<>` | Equality    | Left to Right   | Excel: `<>`, C/Python/JS: `!=` |
| `&`       | Bitwise AND      | Left to Right   | |
| `~`       | Bitwise XOR      | Left to Right   | Lua-style |
| `\|`      | Bitwise OR       | Left to Right   | |
| `&&`      | Logical AND      | Left to Right   | |
| `\|\|`    | Logical OR       | Left to Right   | |
| `?:`      | Conditional      | Right to Left   | Ternary |
| `\|:`     | Pipe             | Left to Right   | |

## Examples

```go
// Basic arithmetic
result, err := uexl.Eval("10 + 20")  // Returns 30

// Power operators (both styles work)
result, err := uexl.Eval("2 ^ 8")    // Excel style: Returns 256
result, err := uexl.Eval("2 ** 8")   // Python/JS style: Returns 256

// Bitwise operations
result, err := uexl.Eval("5 ~ 3")    // XOR: Returns 6
result, err := uexl.Eval("~5")       // NOT: Returns -6

// Not-equals (both styles work)
result, err := uexl.Eval("5 <> 3")   // Excel style: Returns true
result, err := uexl.Eval("5 != 3")   // C/Python/JS style: Returns true

// Pipe operations
result, err := uexl.Eval("10 + 20 |: $1 * 2") // Returns 60

if err != nil {
    log.Fatal(err)
}

// Compile and run with context
exprc, err := uexl.Compile("base ^ 2 * rate")
if err != nil {
    log.Fatal(err)
}

result, err := uexl.Run(exprc, map[string]any{
  "base": 10,
  "rate": 2,
})

fmt.Println("Result:", result)  // Output: Result: 200
```
