# Mastering UExL — Book Development Strategy

**Working Title:** Mastering UExL: The Universal Expression Language for Modern Applications
**Series fit:** O'Reilly / Packt "Mastering" series
**Estimated length:** ~380–420 pages (print), ~65,000–75,000 words
**Audience level:** Intermediate developers who embed or build expression-powered systems

---

## 1. Positioning Statement

> "UExL is to dynamic expression evaluation what SQL is to data querying — a portable, embeddable language with clear semantics that any application layer can understand. Mastering UExL teaches you to write it fluently, embed it safely, and extend it powerfully."

The book is not a reference manual (the docs already cover that). It is a guided journey from *why expressions matter* through *production-grade UExL architecture*, built around real-world cases and transferable mental models.

---

## 2. Target Readers

| Reader | Background | Goal |
|--------|-----------|------|
| **Primary** | Go developer building rule engines, config systems, ETL pipelines | Embed and extend UExL in production code |
| **Secondary** | Non-Go developer (Python/JS) using UExL via WASM or a host API | Write sophisticated expressions confidently |
| **Tertiary** | Tech lead / architect | Evaluate UExL fit, understand trade-offs |

---

## 3. Guiding Principles for Writing Style

### 3.1 Voice and Tone
- **Authoritative but welcoming.** We speak as an experienced guide, not a spec document.
- **Show before tell.** Every concept opens with a concrete code example, *then* explains why it works.
- **Honest about trade-offs.** Explicitly name what UExL is *not* designed for (Turing-completeness, I/O, mutation). Readers trust authors who admit limits.
- **Progressive disclosure.** Simple rules first, exceptions and edge cases later — never in the same paragraph.

### 3.2 Formatting Conventions
- **Code blocks labelled by language** (`uexl`, `go`, `json`, `yaml`).
- **Callout boxes** for three purposes only: `NOTE` (context), `WARNING` (pitfall), `TIP` (shortcut). Avoid overuse.
- **Table-first for reference.** When comparing multiple options, use a table rather than bullet lists.
- **One concept per section.** Sections should be skimmable; a reader who already knows the concept must be able to skip safely.
- **Real names for examples.** Use `orders`, `users`, `products`, `config` — not `foo`, `bar`, `x`.

### 3.3 Code Listing Standards
- Expressions ≤ 3 lines: inline in paragraph with backticks.
- Expressions > 3 lines or with explanatory annotation: numbered code listing with `// comment` annotations.
- Every listing that produces a value should show the result as a `// => value` comment.
- Go integration listings always include error handling.

### 3.4 Exercise Design
- Each chapter ends with **3–5 progressive exercises**:
  1. *Recall* — verify understanding (fill in the blank / predict output)
  2. *Apply* — write an expression for a given spec
  3. *Extend* — open-ended real-world scenario
- An online solutions appendix (GitHub repo) accompanies the book.

---

## 4. Chapter Architecture

The book is organized into **five parts**, moving from foundations to mastery.

```
PART I   — Foundations           (Chapters 1–4)    ~90 pages
PART II  — Core Language         (Chapters 5–9)    ~100 pages
PART III — Pipes & Transformation (Chapters 10–12)  ~75 pages
PART IV  — Embedding in Go       (Chapters 13–15)  ~70 pages
PART V   — Production & Mastery  (Chapters 16–18)  ~60 pages

Appendices                                          ~25 pages
```

---

## 5. Detailed Chapter Plan

---

### PART I — FOUNDATIONS

---

#### Chapter 1 — The Expression Evaluation Problem (Why UExL Exists)

**Purpose:** Create the "pain before the cure." Readers must feel the problem before embracing the solution.

**Key topics:**
- The recurring need for runtime-configurable logic (business rules, config, ETL)
- Why scripting languages (Lua, JavaScript) are overkill for embedded evaluation
- The partial solutions: jsep, expr, cel-go — their trade-offs and gaps
- UExL's position: portable, embeddable, explicit-semantics expression engine

**Table: Embedded Expression Options Compared**

| Solution | Language | Sandboxed | Pipe-native | Nullish-explicit | WASM-portable |
|----------|----------|-----------|-------------|------------------|---------------|
| cel-go | Go | Yes | No | No | No |
| expr | Go | Partial | No | No | No |
| gval | Go | Yes | No | No | No |
| **UExL** | Go | **Yes** | **Yes** | **Yes** | **Yes** |

**Closing section:** "What you'll build in this book" — a concrete multi-part project introduced in chapter 1 and extended through the book (e.g., a product filtering/pricing engine used in e-commerce backend).

---

#### Chapter 2 — Setting Up and Your First Expression

**Purpose:** Zero-to-evaluated as fast as possible.

**Key topics:**
- Installation: `go get github.com/maniartech/uexl`
- The three-stage pipeline (parse → compile → run) as a mental model
- First expressions in the playground (UExL WASM playground link)
- Using `uexl.Eval()` for quick one-shot evaluation
- Using the three-stage API (`parser.ParseString` → `compiler.New().Compile()` → `vm.New().Run()`)

**Code walkthrough:** Annotated Go listing showing parse/compile/run with error handling at each stage.

**Sidebar:** "The Three-Stage Pipeline: Why Not Interpret Directly?" — explains the bytecode approach and its performance + safety benefits.

---

#### Chapter 3 — Data Types and Literals

**Purpose:** Build the reader's mental model of UExL's type system before any operators are introduced.

**Key topics:**
- Numbers: integers, floats, scientific notation; `NaN` and `Inf` (IEEE-754 opt-in)
- Strings: single vs. double quotes; immutability; byte-level default
- Booleans: `true` / `false`; truthiness rules
- Null: what it means to be nullish (null vs. absent)
- Arrays: ordered, zero-indexed, heterogeneous
- Objects: key-value pairs; identifier vs. quoted keys
- Type coercion policy: UExL is NOT implicitly coercive — explicit conversion required

**Table: Truthiness Reference**

| Value | Truthy? |
|-------|---------|
| `1`, `3.14` | Yes |
| `0` | No |
| `"hello"` | Yes |
| `""` | No |
| `true` | Yes |
| `false` | No |
| `null` | No |
| `[]` | No |
| `{}` | No |
| `[0]` | Yes |

**Warning box:** NaN/Inf behavior — `NaN != NaN` is true; arithmetic comparisons with NaN always false.

---

#### Chapter 4 — Identifiers, Variables, and Context

**Purpose:** Explain how data enters UExL expressions, and how the reader controls it.

**Key topics:**
- What "context" means: the `map[string]any` passed at runtime
- Identifier resolution order
- Context variables vs. system variables (`$item`, `$index`, `$acc`, `$last`)
- Absent variables: treated as null (not an error)
- Strict variables: how to make referencing an undefined var an error (host config)
- Naming rules for identifiers: case sensitivity, Unicode characters allowed

**Practical example:** Building a pricing expression context from a Go struct:
```go
ctx := map[string]any{
    "product": product,
    "customer": customer,
    "today":    time.Now().Format("2006-01-02"),
}
result, err := machine.Run(bytecode, ctx)
```

---

### PART II — CORE LANGUAGE

---

#### Chapter 5 — Operators In Depth

**Purpose:** Comprehensive, structured coverage of all operators with emphasis on the surprising/unique ones.

**Key topics:**
- Arithmetic: `+`, `-`, `*`, `/`, `%`; integer vs. float results
- Power: `**` and `^`; right-associativity explained (`2**3**2 == 512`)
- Comparison: `<`, `<=`, `>`, `>=`, `==`, `!=`, `<>` (Excel alias)
- Logical: `&&`, `||`, `!`; short-circuit evaluation
- Bitwise: `&`, `|`, `~` (XOR), `~x` (NOT) — Lua-style
- String concatenation: `+` with strings
- Ternary: `condition ? then : else`; nesting patterns
- Unary patterns: `!!value` for bool conversion, `--x` for double negation

**Operator Precedence Table** (full, authoritative — extracted from `advanced-concepts.md` and expanded):

| Level | Operators | Associativity |
|-------|-----------|---------------|
| 1 (highest) | `!`, unary `-`, `~` | Right |
| 2 | `**`, `^` | Right |
| 3 | `*`, `/`, `%` | Left |
| 4 | `+`, `-` | Left |
| 5 | `<`, `<=`, `>`, `>=` | Left |
| 6 | `==`, `!=`, `<>` | Left |
| 7 | `&` (bitwise AND) | Left |
| 8 | `~` (bitwise XOR) | Left |
| 9 | `\|` (bitwise OR) | Left |
| 10 | `&&` | Left |
| 11 | `\|\|` | Left |
| 12 | `??` | Left |
| 13 (lowest) | `? :` | Right |

**Warning box:** `|` is bitwise OR, NOT the pipe operator. The pipe operator is `|type:` (with a colon).

---

#### Chapter 6 — Property and Index Access

**Purpose:** The access operators are where most runtime errors occur; master them here.

**Key topics:**
- Dot access: `obj.key` — strict by default
- Bracket access: `obj["key"]`, `arr[i]` — strict by default
- Optional chaining: `obj?.key`, `arr?.[i]` — guards null base AND missing member
- Chained access: `a.b.c?.d?.[e]`
- The exact semantics of `?.` — it short-circuits the *rest* of the chain, not just one step
- Practical patterns: partial object data, API responses with optional fields

**Decision tree diagram:** "Which access operator should I use?"
```
Is the base guaranteed non-null?
  YES → use .key or [i]
  NO  → Does a null base mean skip (not error)?
          YES → use ?.key or ?.[i]
          NO  → guard explicitly: (base ?? throw("...")).key
```

---

#### Chapter 7 — Nullish and Optional Semantics (The Heart of UExL)

**Purpose:** This is UExL's most distinctive feature area. Readers who master this chapter write safer, more intentional expressions.

**Key topics:**
- What "nullish" means: `null` and absent are equivalent
- `??` — nullish coalescing: falls back on null/absent, preserves falsy values (`0`, `""`, `false`, `[]`, `{}`)
- Why `||` is wrong for defaults when falsy values are meaningful
- Safe mode: `x.a.b ?? c` — only `b` is softened, not `a`
- Interaction of `?.` and `??`: `user?.address?.city ?? "unknown"`
- Patterns: layered defaults, safe navigation, conditional property building

**Comparison table: `??` vs `||`**

| Expression | Value of `x` | `x ?? "default"` | `x \|\| "default"` |
|-----------|--------------|-------------------|---------------------|
| | `null` | `"default"` | `"default"` |
| | absent | `"default"` | `"default"` |
| | `0` | `0` | `"default"` ⚠️ |
| | `""` | `""` | `"default"` ⚠️ |
| | `false` | `false` | `"default"` ⚠️ |
| | `"hello"` | `"hello"` | `"hello"` |

**Warning:** Use `||` for control flow logic, not data defaults. Treating falsy as "missing" is the most common source of subtle bugs.

---

#### Chapter 8 — Strings and Unicode

**Purpose:** Teach the three Unicode levels explicitly; build habits that prevent silent data corruption.

**Key topics:**
- The three levels: bytes (default), runes (code points), grapheme clusters
- Why the default is byte-level (Go compatibility, UDF alignment, O(1) len)
- When each level is appropriate:
  - Bytes → protocols, storage, regexp match offsets
  - Runes → ASCII processing, identifier-safe operations
  - Graphemes → UI display, safe truncation, user-visible length
- Function reference: `len`, `substr`, `runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`
- Explode-to-array: `runes(s)`, `graphemes(s)`, `bytes(s)`
- String concatenation with `+`; `join()` for arrays
- Practical patterns: safe truncation for UI, byte counting for protocol headers

**Visual: The three-layer model for `"café\u0301"`**
```
String: café + combining-acute
Bytes:  [63 61 66 C3 A9 CC 81] = 7
Runes:  [c  a  f  é  ́      ] = 5
Graphs: [c  a  f  é         ] = 4
```

---

#### Chapter 9 — Functions

**Purpose:** Comprehensive built-in function reference *and* the patterns for calling, composing, and replacing them with host extensions.

**Subsections:**

**9.1 Calling functions**
- Syntax: `fn(arg1, arg2, …)`
- Arguments are fully evaluated before the call
- Multiple return values: UExL functions return one value (no destructuring)

**9.2 Built-in function categories**

| Category | Functions |
|----------|-----------|
| Math | `abs`, `ceil`, `floor`, `round`, `min`, `max`, `sum` |
| String | `concat`, `upper`, `lower`, `trim`, `contains`, `startsWith`, `endsWith`, `replace`, `split`, `join` |
| Array | `len`, `flat`, `keys`, `values` |
| Type | `string`, `number`, `bool`, `typeof` |
| Object | `set`, `keys`, `values` |
| Unicode | `runes`, `graphemes`, `bytes`, `runeLen`, `graphemeLen`, `runeSubstr`, `graphemeSubstr` |

**9.3 Composing functions**
```uexl
join(runes("naïve") |filter: $item != "ï", "")  // "nave"
upper(trim("  hello  "))                          // "HELLO"
```

**9.4 Purity policy** — no mutation, no side effects, explicit returns

**9.5 Registering host functions in Go** — `LibContext.Functions` map, type contracts, error propagation

---

### PART III — PIPES AND DATA TRANSFORMATION

---

#### Chapter 10 — Understanding Pipes

**Purpose:** Build the mental model for pipes from first principles before introducing all pipe types.

**Key topics:**
- The pipe as a Unix-pipeline analogy: `cmd1 | cmd2 | cmd3`
- Difference from function composition: pipes are *sequential and scopeful*, not nested
- The `|:` passthrough pipe — the simplest case
- `$last`: always available, always the previous stage's output
- Why pipe predicates are compiled to bytecode (not interpreted): performance + safety
- Reading multi-stage pipelines top-to-bottom

**Side-by-side:** Nested function call vs. pipe chain — same logic, different readability:
```uexl
// Nested (inside-out reading)
join(map(filter(users, $item.active), $item.name), ", ")

// Piped (left-to-right reading)
users |filter: $item.active |map: $item.name |: join($last, ", ")
```

---

#### Chapter 11 — All Pipe Types In Depth

**Purpose:** Complete, example-rich reference for all 13 pipe types, organized by use case rather than alphabetically.

**Sections:**

**11.1 Pass-through: `|:`**
**11.2 Transformation: `|map:` and `|flatMap:`**
**11.3 Selection: `|filter:` and `|find:`**
**11.4 Aggregation: `|reduce:`** — accumulator semantics, starting value with `??`, common patterns (sum, concat, group-by-manual)
**11.5 Boolean checks: `|some:` and `|every:`** — short-circuit behavior
**11.6 Ordering and uniqueness: `|sort:` and `|unique:`**
**11.7 Grouping: `|groupBy:`** — result structure (`map[string][]any`), practical reporting patterns
**11.8 Windowing and chunking: `|window:` and `|chunk:`** — moving averages, pagination, batch processing

**Per-type layout:**
1. One-line description
2. Scope variables emitted
3. Minimal example → expected result
4. Real-world example (from the book's running e-commerce project)
5. Edge cases (empty array, single element, type mismatch)

---

#### Chapter 12 — Advanced Pipe Patterns

**Purpose:** Combine pipe primitives into sophisticated pipelines; teach debugging and readability techniques.

**Key topics:**
- Multi-stage annotation pattern (using `// comments` for each stage)
- Pipe aliasing: `|map as $processed:`
- Nested pipes: `outer |map: (inner |filter: ...)`
- `$last` threading through nested scopes
- Referencing outer context inside a pipe predicate
- Performance considerations: short-circuit with `|some:`/`|every:` vs. `|filter:|: len()`
- Building a complete data pipeline: ingest → validate → transform → aggregate → format

**Capstone exercise:** E-commerce dashboard — given raw order data, produce a summary object with totals, top products, and customer tiers in a single chained expression.

---

### PART IV — EMBEDDING UEXL IN GO

---

#### Chapter 13 — The Three-Stage API

**Purpose:** Go developers get the complete integration picture: how to use each stage independently and when to separate them.

**Key topics:**
- `parser.ParseString()` — when to use standalone (expression validation in APIs)
- `compiler.New().Compile()` — when to cache bytecode (hot-path expressions compiled once)
- `vm.New(LibContext{...}).Run(bytecode, ctx)` — per-request execution with fresh context
- Error types at each stage: `errors.ParserError`, compiler errors, VM runtime errors
- Structuring error messages for end users (not exposing internals)

**Architecture pattern: Compile-once, run-many**
```go
// At startup
ast, _ := parser.ParseString(expr)
comp := compiler.New()
comp.Compile(ast)
bc := comp.ByteCode()

// Per request (concurrency-safe)
machine := vm.New(vm.LibContext{Functions: vm.Builtins, PipeHandlers: vm.DefaultPipeHandlers})
result, err := machine.Run(bc, ctx)
```

**Note:** Thread safety considerations — VM instances are not goroutine-safe; use sync.Pool or create per-goroutine instances.

---

#### Chapter 14 — Registering Custom Functions and Pipes

**Purpose:** Teach the extension model fully, with patterns for production-grade custom functions.

**Key topics:**

**14.1 Custom functions**
- Signature: `func(args []any) (any, error)`
- Type assertion patterns (always check `ok`)
- Returning errors vs. returning null
- Documenting arg count and types (since UExL has no static typing)
- Registering in `LibContext.Functions`

**14.2 Custom pipes**
- Pipe handler signature
- When to build a custom pipe vs. composing built-in pipes
- Example: `|paginate:`, `|validate:`, `|audit:`

**14.3 Security boundary**
- Never expose filesystem, network, or shell functions as UExL builtins
- Sandboxing principles: restrict what the embedding functions can do
- Rate-limiting expression evaluation (inflight counter, timeout context)

---

#### Chapter 15 — Context Design and the UExL Contract

**Purpose:** How to design the context API between your Go application and UExL expressions — the most critical design decision for maintainable embeddings.

**Key topics:**
- Context as a public API: treat it like a REST API (versioning, backward compat)
- Naming conventions for context keys (`camelCase`, no `$` prefix)
- Nested vs. flat context structures: trade-offs
- Providing computed values vs. raw structs
- Go struct → `map[string]any` conversion patterns (reflection vs. manual vs. `encoding/json` marshal-then-unmarshal)
- Environment configuration: `EnvConfig` options (NaN/Inf, strict mode)

**Anti-pattern:** Exposing raw DB models as context — leaks schema, breaks expressions on schema changes.
**Pattern:** Expose a DTO/view model purpose-built for expression consumption.

---

### PART V — PRODUCTION AND MASTERY

---

#### Chapter 16 — Error Handling, Validation, and Debugging

**Purpose:** Production systems need graceful error handling; this chapter covers the full lifecycle.

**Key topics:**
- Parser errors: syntax, unexpected token, unterminated string
- Compiler errors: unsupported node type, malformed AST
- VM runtime errors: type error, reference error, argument error, division by zero
- Wrapping errors for end users: hide internals, show the expression snippet + position
- Expression validation endpoint pattern (parse only, return structured error)
- Debugging tips: using the WASM playground; logging compiled bytecode; isolation testing

**Error taxonomy table:**

| Error Type | When | Example Expression | Go Error Type |
|-----------|------|-------------------|---------------|
| SyntaxError | Parse | `{a: 1,, b}` | `errors.ParserError` |
| TypeError | Runtime | `"abc" * 2` | VM error string |
| ReferenceError | Runtime | `undeclaredVar + 1` | VM error string |
| ArgumentError | Runtime | `min()` | VM error string |
| DivisionByZero | Runtime | `10 / 0` | VM error string |

---

#### Chapter 17 — Performance Tuning and Benchmarking

**Purpose:** Production UExL at scale — measure, profile, and optimize.

**Key topics:**
- The compile-once/run-many pattern's performance impact
- VM benchmark baseline: simple boolean expressions ~100 ns/op target
- Benchmarking methodology: `go test -bench=. -benchtime=20s`
- CPU profiling with `go tool pprof`
- Allocation analysis: `go test -memprofile=mem.prof`
- Expression complexity and cost: pipe chains vs. nested functions
- When to pre-aggregate in Go vs. letting UExL do it
- Build configuration: `env_config.go` options

**Table: Performance characteristics by expression type**

| Expression Pattern | Relative Cost | Notes |
|-------------------|---------------|-------|
| Simple arithmetic/comparison | Very low (1×) | Baseline |
| Short-circuit `&&` / `\|\|` | Very low (1×) | Stops early |
| Single-stage `|map:` over 1k items | Low (~10×) | Linear in input |
| Multi-stage pipe chain | Medium (2–5× per stage) | Each stage is a frame |
| `|groupBy:` | Medium–High | Map allocation per group |
| `|window:` / `|chunk:` | Medium | Slice allocation per window |

---

#### Chapter 18 — Real-World Architectures and Patterns

**Purpose:** Close the book with architecture blueprints, reinforcing that UExL is a serious production tool.

**Case Studies:**

**18.1 Rule Engine for a Lending Platform**
- Expressions stored in database, compiled at load time, executed per loan application
- Context: applicant data, credit score, product parameters
- Pattern: safe compilation sandbox, per-request VM pool

**18.2 Dynamic Pricing Engine**
- Price rules authored by business analysts (non-developers)
- Expression validation UI using WASM playground
- Go backend: `sync.Pool` for VMs, struct-to-map context helper

**18.3 ETL Data Transformation Pipeline**
- Pipe-heavy expressions for record transformation
- Batch execution: one compiled expression, thousands of context maps
- Error accumulation: process all records, collect errors, report summary

**18.4 Configuration-Driven Feature Flags**
- YAML config files with embedded UExL expressions
- Schema: `enabled: "uexl! user.tier == 'premium' && featureDate <= today"`
- Safe parsing layer that validates all expressions at config load time

---

## 6. Appendices

### Appendix A — Complete Operator Quick Reference
Full table: operator, syntax, precedence, associativity, example.

### Appendix B — Complete Built-in Function Reference
Alphabetical list: signature, description, example, edge cases, error conditions.

### Appendix C — All Pipe Types Quick Reference
One-page spread: pipe name, scope vars, 1-line description, 2-line example.

### Appendix D — UExL for Excel Users
Condensed migration guide (expanded from existing doc).

### Appendix E — UExL for Python/JavaScript Users
Condensed migration guide (expanded from existing doc).

### Appendix F — Grammar Reference (EBNF)
Formal grammar for language tooling authors.

### Appendix G — Glossary
~40 terms: nullish, boolish, safe mode, pipe, predicate, context, bytecode, etc.

---

## 7. Running Project Thread

The book should introduce a single realistic project in Chapter 1 and evolve it throughout:

**Project: "ShopLogic" — A configurable e-commerce pricing and filtering engine**

| Part | Project milestone |
|------|------------------|
| Part I | Context setup; first price calculation expression |
| Part II | Operator-based discount logic; nullish-safe customer tier lookup |
| Part III | Full order pipeline: filter → transform → aggregate → format |
| Part IV | Go embedding; cached compilation; custom `currency()` function |
| Part V | Production hardening; benchmarks; rule versioning |

This gives every code example a realistic grounding and lets readers build something tangible.

---

## 8. Front Matter

- **Foreword:** Written by a domain expert (rule engine architect or language designers from Go ecosystem)
- **Preface:** Author's story — why UExL was built, what problems it solves
- **How to read this book:** Paths for different reader types
  - *Just expressions (no Go):* Chapters 1, 3–12
  - *Go integration focus:* Chapters 1–2, 5–7, 13–16
  - *Reference reading:* Appendices A–C + relevant chapter sections

---

## 9. Back Matter

- **Index:** Comprehensive; includes operator symbols, function names, error types
- **Bibliography:** Papers on expression evaluation, PEG parsing, bytecode VMs if relevant
- **About the Author**

---

## 10. Content Gaps to Fill (From Source → Book)

The existing docs are a solid reference base. The book needs to *add* these:

| Gap | Where |
|-----|-------|
| Deep bytecode architecture explanation (parser.Node → InstructionBlock → OpPipe) | Chapter 13 sidebar |
| Thread safety of VM instances | Chapter 13 |
| `sync.Pool` VM pooling pattern | Chapter 13, 17 |
| Context design best practices (DTO vs. raw struct) | Chapter 15 |
| Real VM error taxonomy with Go error types | Chapter 16 |
| `|groupBy:` result structure and iteration patterns | Chapter 11 |
| `|window:` moving-average use case | Chapter 11 |
| Custom pipe registration API | Chapter 14 |
| EnvConfig / build-tag options explained | Chapter 15, 17 |
| `set()` function patterns for immutable updates | Chapter 9 |
| WASM playground usage for rapid dev/debug | Chapter 16 |

---

## 11. Production Schedule (Suggested Milestones)

| Milestone | Deliverable |
|-----------|-------------|
| M1 | Part I (Chapters 1–4) — first draft |
| M2 | Part II (Chapters 5–9) — first draft |
| M3 | Part III (Chapters 10–12) — first draft |
| M4 | Part IV (Chapters 13–15) — first draft |
| M5 | Part V + Appendices — first draft |
| M6 | Full manuscript review + technical accuracy pass |
| M7 | Code listings verified against latest UExL source |
| M8 | Exercise solutions and companion repo published |
| M9 | Final copyedit and index |

---

## 12. Companion Materials

- **GitHub repo:** `mastering-uexl-book` — all code listings, exercise solutions, ShopLogic project
- **WASM Playground link:** Featured in Chapter 2 and 16 as a debugging/learning tool
- **Chapter-by-chapter README:** Each chapter folder contains a `README.md` with setup instructions for runnable examples

---

## 13. Key Differentiators vs. Existing Documentation

| Existing docs | Book adds |
|--------------|-----------|
| Feature reference | *Why* each feature was designed this way |
| Parallel API listing | Progressive skill building |
| Isolated examples | Unified project thread (ShopLogic) |
| Internal architecture notes | Accessible architecture narrative for embedders |
| Developer-targeted | Accessible also to non-Go users (WASM/API consumers) |
| Markdown files | Narrative prose with exercises, diagrams, callout boxes |
