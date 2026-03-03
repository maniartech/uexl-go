# UExL — Open Source Launch Checklist

Track progress toward the first public release of `github.com/maniartech/uexl`.

Legend: ✅ Done · 🔧 In progress · ⬜ Not started

---

## 1. Legal & Identity

- [ ] **Add LICENSE file** — No license = "all rights reserved" by default; nobody can legally use or contribute. Choose one (MIT recommended for a library) and add a `LICENSE` file to the repo root.
- [x] **Fix module path** — `go.mod`, README, and book now all consistently use `github.com/maniartech/uexl`.
- [ ] **Add `CHANGELOG.md`** — Document what's in the first release so users know what they're getting.
- [ ] **Decide version tag** — Will the first release be `v0.1.0` (unstable API signal) or `v1.0.0`? Given the planned features not yet implemented, `v0.1.0` is more honest.

---

## 2. Public API

- [ ] **Design and implement a top-level API in `uexl.go`** — Currently only `EvalExpr(expr string)` is exported with no way to pass context variables, register functions, or register pipe handlers. Minimum surface needed:
  ```go
  // Simple one-shot evaluation
  func Eval(expr string, vars map[string]any) (any, error)

  // Reusable evaluator with registered functions and pipe handlers
  type Evaluator struct { ... }
  func NewEvaluator(opts ...Option) *Evaluator
  func (e *Evaluator) Eval(expr string, vars map[string]any) (any, error)

  // Options / registration
  func WithFunctions(fns map[string]any) Option
  func WithPipeHandlers(handlers vm.PipeHandlers) Option
  ```
- [ ] **Ensure the public API is stable-ish** — Every exported symbol in a `v1.x` release is a compatibility commitment. Keep the surface minimal and intentional.
- [ ] **Add `doc.go` at package root** — Package-level godoc comment describing UExL, linking to the book, and showing a minimal usage example.

---

## 3. Correctness & Robustness

- [x] **All tests pass with `-race`** — 284 test functions across parser, compiler, VM, and utils. All green.
- [ ] **Replace `panic()` in production code** — Two remaining panics violate the zero-panic policy:
  - `code/code.go:133` — unknown opcode in `String()` method; return `"<unknown opcode>"` instead.
  - `compiler/compiler_utils.go:83` — `exitScope` underflow; return an error instead of panicking.
- [ ] **Review `TODO` comments in production code**:
  - `code/code.go:123` — `// TODO: handle error properly`
  - `parser/constants/operators.go:139` — unused function, remove or use it
  - `vm/vm_handlers.go:174` — strings.Builder optimization (low priority)

---

## 4. Repository Hygiene

- [ ] **Remove `playground_test.go`** — Dev scratch file in the root package; should not ship.
- [ ] **Clean up or remove `cmd/`** — Currently a dev tool that imports internal packages to demo panic fixes. Either make it a proper example (`examples/basic/`) or delete it.
- [ ] **Move WIP performance artifacts** — Root-level `*.md` performance files (`OPTIMIZATION_LEARNINGS.md`, `PHASE2C_RESULTS.md`, `FINAL_PERFORMANCE_RESULTS.md`, etc.) are internal dev notes. Move to `wip-notes/` or exclude via `.gitignore`.
- [ ] **Fix `.gitignore`** — Currently ignores `vendor/` (vendor should be committed for a library) and `*.txt`/`*.prof` (reasonable). Remove the `vendor/` line.
- [ ] **Add `.gitignore` entries** — `*.prof`, `cpu_*.prof`, `mem_*.prof` build artifacts.

---

## 5. Documentation

- [x] **Book structure** — Main book covers all implemented features; "Planned Features" section covers unimplemented ones. All cross-references updated.
- [x] **Import paths correct** — All docs use `github.com/maniartech/uexl`.
- [ ] **Update `book/golang/overview.md`** — Rewrite with the actual top-level API once it's designed (currently shows stubs with `uexl.Eval` and `uexl.RegisterFunction` which don't exist yet).
- [ ] **Update README** — Currently claims `uexl.Eval()` (doesn't exist; actual is `EvalExpr`) and "zero allocations, 227 ns/op" (current benchmark: ~145 ns/op, which is better but not zero-alloc). Update after the public API is settled.
- [ ] **Add contributing guide** — `CONTRIBUTING.md` with: how to run tests, how to run benchmarks, code style notes, PR expectations.

---

## 6. CI / Release Infrastructure

- [ ] **Add GitHub Actions workflow** — `.github/workflows/ci.yml` running `go test ./... -race` on every push and PR across Go 1.22, 1.23, 1.24, 1.25, 1.26.
- [ ] **Add `go vet` to CI** — Catches common correctness issues automatically.
- [ ] **Tag first release** — `git tag v0.1.0 && git push origin v0.1.0` after all blockers are resolved. This is what makes it `go get`-able by version.

---

## 7. Nice-to-Have (Post-Launch)

- [ ] **Add `CodeQL` / `govulncheck` to CI** — Security scanning, important for a library.
- [ ] **Add `golangci-lint`** — Catches unused exports, shadowing, etc.
- [ ] **Pkg.go.dev badge in README** — Auto-generated after first tagged release.
- [ ] **Benchmark CI** — Track performance regressions across commits (e.g., `benchstat` with stored baselines).

---

## Summary

| Category | Status |
|---|---|
| Legal (license) | ⬜ Blocked |
| Public API | ⬜ Needs design |
| Tests passing | ✅ |
| Panics in prod code | 🔧 2 remaining |
| Repo hygiene | 🔧 Several items |
| Documentation | 🔧 Mostly done, needs API update |
| CI | ⬜ Not started |

**Minimum to publish:** Items in sections 1 (license), 2 (API), and 3 (panic removal). Everything else can follow in a patch release.
