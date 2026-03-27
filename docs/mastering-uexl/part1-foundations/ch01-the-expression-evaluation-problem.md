# Chapter 1: The Expression Evaluation Problem

> "Every application eventually needs to let the outside world make decisions inside it. The question isn't *whether* you need dynamic expression evaluation — it's *how well* you'll do it."

---

## 1.1 The Moment Every Developer Faces

Picture a product manager walking up to your desk with a reasonable request: "Instead of hard-coding the discount logic, can we let our business team configure it? They want to set rules like *'apply 15% if order total exceeds 200 and customer tier is Gold.'"*

You nod. It sounds simple. Three weeks later, you've built a fragile custom mini-language, added a YAML-based rule format that nobody on the business team understands, or worse — you deployed a JavaScript sandbox that runs untrusted code in production with all the security implications that brings.

This scenario plays out in every organization that builds systems with any configurability. It shows up as:

- **Rule engines** — discount rules, fraud detection, approval workflows
- **Configuration systems** — feature flags that go beyond true/false, environment-specific thresholds
- **ETL pipelines** — field transformation and validation rules authored by analysts, not engineers
- **Dynamic forms** — visibility and validation rules for form fields driven by data
- **Policy engines** — access control, rate limiting, data masking rules

The common thread: *someone needs to embed a decision into a system without rebuilding and redeploying it*.

---

## 1.2 Why Not Just Use a Scripting Language?

The instinct for many developers is to reach for an embedded scripting language. V8 for JavaScript, Lua, Python via an FFI — these are general-purpose, well-documented, and expressive.

The problem is that *general-purpose* is exactly what you don't want.

### The full-scripting trap

When you embed a scripting language, you inherit its full power — including all the things you didn't want to expose:

- **File system access.** A Lua expression can read your `/etc/passwd`.
- **Network calls.** A JavaScript expression can exfiltrate your database credentials.
- **Infinite loops.** A user can write `while(true) {}` and halt your server.
- **Memory bombs.** `Array(1e9).fill(0)` will trigger your OOM killer.

Sandboxing scripting runtimes is a solved problem in theory and a painful one in practice. You need to carefully disable APIs, set timeouts, limit memory, prevent prototype pollution (for JavaScript), intercept syscalls, or run in a separate process. Each of those choices adds complexity, latency, or operational burden.

### The expression evaluator sweet spot

What business logic actually needs is far narrower than what scripting languages provide:

- Read values from provided data
- Perform arithmetic, comparison, and string operations
- Filter, transform, and aggregate collections
- Return a result

No I/O. No loops of unbounded length. No mutable state. No side effects. Just: *take this data, apply this logic, give me a result*.

This is the sweet spot of an **expression evaluator** — a tool deliberately designed for this use case, with a hard boundary between what is allowed and what is not.

---

## 1.3 The Landscape Before UExL

Several Go libraries exist for embedded expression evaluation. Each solves the problem partially. Understanding their trade-offs explains why UExL was built.

| Library | Strengths | Limitations |
|---------|-----------|-------------|
| **expr** (Antonmedv) | Fast, strongly typed, CEL-inspired | No pipe operators, no nullish semantics, type binding required |
| **cel-go** (Google) | CEL standard, protocol-buffer friendly | Complex setup, verbose, no pipe transforms, no WASM |
| **gval** | Extensible operators, simple API | No collection pipelines, limited Unicode, no nullish safety |
| **tengo** | Full scripting language | Too powerful, not sandboxed by default |
| **goja** (V8 JS) | Full JavaScript | Full exposure — full risk |

The gaps are consistent across most solutions:

1. **No pipe-based data transformation.** Collection operations (map, filter, reduce) require nested function calls that are hard to read and write.
2. **No explicit nullish semantics.** The JavaScript confusion between "missing", "null", and "falsy" infects every expression evaluator influenced by JS.
3. **No cross-platform portability.** Expressions written for a Go backend can't be validated or previewed in a browser UI.
4. **Implicit type coercion.** Silently converting `"5"` to `5` in arithmetic hides bugs.

---

## 1.4 Enter UExL

UExL (Universal Expression Language) was built to fill these gaps with deliberate, opinionated design choices.

**The Universal in Universal Expression Language** has a specific meaning: expressions written in UExL should be portable across languages and platforms — usable in a Go backend, a JavaScript frontend, a Python data pipeline, or any WebAssembly-capable environment. Like regular expressions became a universal standard for text pattern matching, UExL aims to be the universal standard for expression evaluation.

### What makes UExL different

**Pipe-native data transformation.** Instead of `sum(map(filter(orders, predicate), transform))`, UExL gives you:

```uexl
orders
  |filter: $item.status == 'paid'
  |map:    $item.total
  |reduce: ($acc ?? 0) + $item
```

Left-to-right, top-to-bottom — matching the mental model of how data actually flows.

**Explicit nullish semantics.** UExL cleanly separates three categories that other systems conflate:

- *Nullish* (`null` or absent) — handled by `??` and `?.`
- *Falsy* (`0`, `""`, `false`, `[]`, `{}`) — handled by `||` and `&&`
- *Missing* (strict access throws) — intentional errors when data structure is wrong

**Three-stage compiled pipeline.** Parse → Compile → Execute. Expressions are compiled to bytecode once and executed many times, making the "compile at startup, run at request time" pattern both ergonomic and fast.

**Zero-panic robustness.** Every error path in UExL returns a structured error. No expression — however malformed or adversarial — can panic the host application.

**WASM-portable.** The same core can run in a browser via WebAssembly, making it possible to build expression editors and validators that use the same evaluation engine as production.

### What UExL deliberately is NOT

UExL is an expression language, not a programming language. It intentionally excludes:

- Assignment operators (`=`, `+=`, `++`)
- Loops (`for`, `while`)
- I/O operations (file, network, console)
- Exception handling (`try`/`catch`)
- Type declarations or class definitions

These omissions are features, not limitations. They define the safe perimeter inside which expressions run.

---

## 1.5 The ShopLogic Project

Throughout this book, we'll build **ShopLogic** — a configurable pricing and filtering engine for an e-commerce platform. It's a realistic system where business stakeholders need to author rules without engineering involvement.

By the end of Part V, ShopLogic will:

- Evaluate pricing rules stored in a database (compile once, run per order)
- Filter and rank products based on configurable criteria
- Transform order data into reporting summaries using pipe chains
- Expose an expression validation API for the rule authoring UI
- Run expressions safely in a concurrent Go service

Here's a preview of what ShopLogic expressions look like in their final form. Each expression below is a single, self-contained UExL expression — UExL evaluates one expression at a time, producing one result.

**Pricing rule** — dynamic discount for gold customers:

```uexl
product.basePrice * (1 - (customer.tier == 'gold' ? 0.15 : customer.tier == 'silver' ? 0.08 : 0))
```

**Product ranking pipeline** — filter, sort, reshape, and take the top 10:

```uexl
products
  |filter: $item.stock > 0 && $item.category == targetCategory
  |sort:   $item.rating
  |map:    { id: $item.id, name: $item.name, price: $item.price * (1 - discount) }
  |: $last[0:10]
```

**Order summary aggregation** — fold completed orders into a totals object:

```uexl
orders
  |filter: $item.status == 'completed'
  |reduce: {
      total: ($acc ?? {total: 0, count: 0}).total + $item.amount,
      count: ($acc ?? {total: 0, count: 0}).count + 1
  }
```

We'll build up to these expressions step by step, starting with the simplest arithmetic in the next chapter.

---

## 1.6 Summary

- Dynamic expression evaluation is a recurring need in production systems: rule engines, configuration, ETL, forms, and policy engines all require it.
- General-purpose scripting languages are the wrong tool — they expose surface area (I/O, memory, loops) that expression evaluation doesn't need and that sandboxing alone can't fully address.
- Existing Go expression libraries address parts of the problem but lack pipe-native collection transforms, explicit nullish semantics, and cross-platform portability.
- UExL fills these gaps with deliberate design: pipe operators, explicit nullish/boolish separation, a compiled bytecode pipeline, zero-panic robustness, and WASM portability.
- This book builds everything around ShopLogic, a realistic e-commerce pricing and filtering engine. You'll build it progressively and have a production-ready embedding by Part V.

---

## Exercises

**1.1 — Recall.** What are three application domains where embedded expression evaluation is commonly needed? Name one security risk introduced by embedding a full scripting language (like JavaScript) instead of a restricted expression evaluator.

**1.2 — Compare.** Pick one of the existing Go expression libraries mentioned in this chapter (expr, cel-go, or gval). Look up its README and identify one feature it provides that UExL also provides, and one feature UExL provides that it does not.

**1.3 — Design.** You're building an e-commerce platform that lets marketing managers configure personalized banner messages using user data. Write a description (in plain English, no code yet) of what data you'd expose in the expression context and what a sample expression rule might look like. What would you *not* expose, and why?
