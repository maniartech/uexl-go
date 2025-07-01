# Introduction

UExL (Universal Expression Language) is a modern, cross-platform, embeddable expression language designed for flexibility, clarity, and power. It enables users to define, evaluate, and chain expressions in a concise and readable syntax, making it ideal for configuration, scripting, and dynamic logic in applications.

## Why UExL?

UExL was created to address the need for a lightweight, expressive, and extensible language that can be embedded in any application. Whether you are building configuration systems, data transformation pipelines, or dynamic business logic, UExL provides a robust foundation with a gentle learning curve.

The following example demonstrates how UExL expressions can be embedded in YAML to compute dynamic values based on runtime data.

```yaml
# Example: Using UExL in a YAML configuration file
threshold: uexl! "100"                # number
isActive: uexl! "user.score > 80 && user.isVerified"   # boolean
welcomeMessage: uexl! "concat('Hello, ', user.name)"   # string
items: uexl! "filter(products, $1.price < 50)"         # array

This demonstrates how UExL expressions can be embedded in YAML to compute dynamic values of different types based on runtime data.
```

## Key Features
- Simple, expressive syntax that is easy to read and write
- Support for numbers, strings, booleans, null, arrays, and objects
- Rich set of operators and built-in functions for common tasks
- Powerful pipe and chaining support for data transformation and functional programming
- Flexible type conversion and coercion
- Clear error handling and debugging support
- Extensible with custom functions and operators

## Who Should Read This Book?
This book is for developers, architects, and technical users who want to:
- Embed a scripting or expression language in their applications
- Write concise and maintainable logic for configuration, validation, or transformation
- Understand the design and implementation of UExL

## What You'll Learn
- The syntax and semantics of UExL
- How to use data types, variables, operators, and expressions
- Advanced features like pipes, custom functions, and extensibility
- Practical examples and best practices
- How to integrate UExL with Go applications

## Book Structure
This book is organized into chapters that progressively introduce UExL concepts, from basic syntax to advanced topics. Each chapter includes detailed explanations, practical examples, and tips for effective usage.

Let's begin your journey into UExL!
