# UExL (Universal Expression Language) in Golang

## Introduction

UExL, short for Universal Expression Language, is a robust and versatile expression evaluator written in Go. Designed for efficiency and simplicity, UExL offers an intuitive approach to handling and evaluating expressions in various formats.

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

(List the key features of UExL.)

## Operator Precedence

| Operators | Type             | Associativity   |
|-----------|------------------|-----------------|
| `(` `)`   | Parentheses      | Left to Right   |
| `.`       | Dot              | Left to Right   |
| `%`       | Modulus          | Left to Right   |
| `*` `/`   | Multiplicative   | Left to Right   |
| `+` `-`   | Additive         | Left to Right   |
| `<<` `>>` | Bitwise Shift    | Left to Right   |
| `<` `>` `<=` `>=` | Comparison | Left to Right |
| `==` `!=` | Equality         | Left to Right   |
| `&` `\|` `^` | Bitwise       | Left to Right   |
| `&&` `\|\|` | Logical        | Left to Right   |
| `\|:`     | Pipe             | Left to Right   |

## Examples

```go
result, err := uexl.Eval("10 + 20 |: $1 * 2") // Returns 60
```
