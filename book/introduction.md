# Introduction

UExL (Universal Expression Language) is a modern, cross-platform, embeddable expression language designed for flexibility, clarity, and power. It enables users to define, evaluate, and chain expressions in a concise and readable syntax, making it ideal for configuration, scripting, and dynamic logic in applications.

## Why UExL?

UExL was created to address the need for a lightweight, expressive, and extensible language that can be embedded in any application. Whether you are building configuration systems, data transformation pipelines, or dynamic business logic, UExL provides a robust foundation with a gentle learning curve.

The following examples demonstrate how UExL expressions can be embedded in different languages to compute dynamic values based on runtime data.

**YAML:**
```yaml
threshold: uexl! "100"                # number
isActive: uexl! "user.score > 80 && user.isVerified"   # boolean
welcomeMessage: uexl! "concat('Hello, ', user.name)"   # string
items: uexl! "filter(products, $1.price < 50)"         # array
```

**JSON:**
```json
{
  "threshold": "uexl! 100",
  "isActive": "uexl! user.score > 80 && user.isVerified",
  "welcomeMessage": "uexl! concat('Hello, ', user.name)",
  "items": "uexl! filter(products, $1.price < 50)"
}
```

**Python:**

```python
from uexl import evaluate

products = [
    {"name": "Book 1", "price": 25, "category": "Books"},
    {"name": "Gadget", "price": 75, "category": "Electronics"}
]

# Evaluate UExL expressions directly
cheap_products = evaluate("filter(products, $1.price < 50)", locals())
book_count = evaluate("count(filter(products, $1.category == 'Books'))", locals())
```

**JavaScript:**

```javascript
const config = {
  threshold: uexl("100"),
  isActive: uexl("user.score > 80 && user.isVerified"),
  welcomeMessage: uexl("concat('Hello, ', user.name)"),
  items: uexl("filter(products, $1.price < 50)")
};

const context = {
  user: { score: 95, isVerified: true },
  products: [
    { name: "Product A", price: 45 },
    { name: "Product B", price: 75 },
    { name: "Product C", price: 25 }
  ]
}

const threshold = config.threshold.eval(context);
const isActive = config.isActive.eval(context);
const welcomeMessage = config.welcomeMessage.eval(context);
const items = config.items.eval(context);

console.log(threshold, isActive, welcomeMessage, items);
```
Currently, we are working on the Golang library with Golang, YAML, and JSON processing support. Support for other languages will follow soon!

## Applications of UExL

UExL is designed to empower a wide range of use cases, making your applications more dynamic, flexible, and maintainable. Here are some of the most impactful ways UExL can be applied:

- **Dynamic Configuration**: Define configuration settings that adapt at runtime, allowing changes without redeploying code. For example, feature flags or environment-specific settings can be controlled by expressions.
- **Data Transformation and Pipelines**: Transform, filter, and aggregate data on the fly in ETL pipelines, analytics dashboards, or reporting tools, all using concise UExL expressions.
- **Dynamic Logic and Business Rules**: Implement business logic or conditional workflows that can be updated by non-developers, such as pricing rules, eligibility checks, or approval flows.
- **Validation**: Validate user input, API payloads, or configuration files with expressive rules that are easy to update as requirements evolve.
- **Visualization and Analytics**: Drive dashboards and visualizations with expressions that compute metrics, filter datasets, or trigger alerts based on live data.
- **No Code / Low Code Platforms**: Enable end-users or administrators to define custom logic and automation without writing traditional code, accelerating development and reducing errors.
- **Notifications and Alerts**: Trigger notifications or send alerts based on complex conditions, such as system health, security breaches, or usage patterns.

UExL's versatility means it can be embedded wherever dynamic evaluation is needed, from configuration files to user interfaces, automation scripts, and beyond.

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
