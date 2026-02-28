# UExL Go — Agent Quality Principles

> These principles govern ALL code changes in this project. Every optimization, feature, and refactor MUST satisfy these criteria.

## Core Quality Criteria

### 1. Military-Grade Robust
- **Zero-panic policy** — NEVER use `panic()` in production paths
- **All error paths handled** — no silent failures, no swallowed errors
- **Race-free** — must pass `go test ./... -race`
- **100% test pass rate** maintained at all times
- **Comprehensive test coverage** — critical paths must be covered by tests (coverage > 90% for core packages)
- **Edge cases included** — tests must cover nulls, empty collections, out-of-bounds, type mismatches, etc.

### 2. Super High Performance
- **Beat ALL competitor benchmarks** (expr, cel-go) across every operation
- **Zero allocations** for primitive expression evaluation (0 B/op, 0 allocs/op)
- **Profile-driven** — every optimization backed by CPU profiles and benchstat
- **Statistical rigor** — p-value < 0.05, minimum 10 runs, ≥5% improvement required

### 3. A+ Quality Code
- **Readable** — code is self-documenting; comments only where logic isn't obvious
- **Consistent** — follow existing patterns in the codebase
- **Type-safe** — always check `ok` on type assertions, use typed APIs where possible
- **Well-tested** — new code requires tests; optimizations require benchmarks

### 4. Future-Proof
- **No hardcoded paths** for specific benchmarks or test cases
- **General solutions** — optimizations must work for ALL inputs, not just known cases
- **Extensible** — new pipes, builtins, operators plug in without modifying core

### 5. DRY (Don't Repeat Yourself)
- **One canonical pattern** per optimization type (e.g., the `pop2Values()` → typed dispatch pattern)
- **Reuse proven patterns** — apply what works for comparisons to arithmetic, logical, string ops
- **Shared infrastructure** — `pipeFastScope`, typed push/pop, `Value` type system

### 6. SRP (Single Responsibility Principle)
- **Each function does one thing** — dispatch logic separate from operation logic
- **Typed inner handlers** receive typed arguments (no `any` in hot paths)
- **Outer dispatcher** resolves types once, delegates to typed handlers

### 7. KISS (Keep It Simple, Stupid)
- **Simplest correct solution** — don't add abstraction layers for theoretical future needs
- **Minimum viable optimization** — if `pop2Values()` in the run loop fixes boxing, do that; don't redesign the VM
- **Proven patterns over novel ideas** — extend what already works (comparison dispatch pattern)

### 8. Not Overengineered
- **No premature abstractions** — no generic optimization frameworks, no plugin architectures for optimizations
- **No speculative features** — only optimize what exists in the codebase today (5 builtins, not 50+)
- **No unnecessary indirection** — direct function calls over interface dispatch in hot paths

## Optimization-Specific Rules

- **Always profile BEFORE and AFTER** — no blind optimization
- **Benchmark with `benchstat`** — statistical validation is mandatory
- **Zero allocation target** — for all non-allocating expression types
- **No test-specific code paths** — if it only speeds up one benchmark, it's cheating
- **Clean up after yourself** — remove hardcoded fast paths (e.g., `tryFastMapArithmetic`)
