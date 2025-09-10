# Compiler and VM Upgrade

This document lays out a concrete, low‑risk plan to bring the compiler and VM to industry‑standard quality with strong robustness guarantees and high performance. It is structured to land small, verifiable improvements first, then optional deeper refactors only if benchmarks indicate need. No code changes are implied by this document alone.

## goals and success criteria

- Robustness
	- No panics from untrusted inputs; all failures return structured errors from compiler/VM.
	- Clear invariants documented where correctness depends on instruction layout or stack contracts.
- Performance
	- Maintain or improve current throughput and allocation profile in end‑to‑end parse → compile → execute benchmarks.
	- Keep hot paths branch‑predictable and allocation‑light.
- Maintainability
	- Eliminate typo’d/public internals; tighten contracts between `compiler` and `vm`.
	- Add comments and small tests for tricky invariants (backpatching, jump semantics, safe mode).

## scope (what stays the same)

- No language semantics changes.
- No opcode layout changes (unless a follow‑up RFC explicitly proposes it).
- All current tests must remain green.

## phased plan

### Phase 0 — Baseline and guardrails

- Capture baseline with the current tree:
	- Tests: `go test ./... -race` and `go vet ./...`.
	- Add a tiny end‑to‑end benchmark (representative expressions: arithmetic, object access, optional/nullish chains, function calls) to track ops/sec and allocations.
- Document current invariants that tests depend on:
	- `OpJumpIfFalsy`/`OpJumpIfTruthy`/`OpJumpIf(Not)Nullish` stack behavior.
	- `replaceOperand` addressing convention (operand starts at `pos+1`).

### Phase 1 — Safety hardening (no behavior change)

- Compiler (`compiler/compiler.go`)
	- Function call callee check: replace direct `node.Function.(*parser.Identifier)` with a checked assertion; return a compile error if callee isn’t an identifier. Prevents panics on edge inputs.
- VM (`vm/vm.go`)
	- Replace direct type assertions on constants/system vars (e.g., `.(string)`) with checked assertions and descriptive errors. Defensive against malformed bytecode or future regressions.
	- Ensure all stack pops/reads guard against underflow in one place (helpers) to avoid scattered checks and accidental panics.
- Optional: Add unit tests covering these error paths (non‑identifier call, bad constant types, stack underflow).

### Phase 2 — Clarity and API hygiene (no behavior change)

- Naming and exports
	- Rename `EmmittedInstruction` → `EmittedInstruction` and unexport if not needed outside the package.
	- Consider unexporting `SystemVars` or providing read‑only accessors if external consumers don’t require mutation.
- Invariant comments
	- Annotate `replaceOperand` usage at each backpatch site with a brief note: “operand index begins at `pos+1` (big‑endian uint16).”
	- In ternary compilation, document why `OpPop` is emitted after patching else: reference VM’s `OpJumpIfFalsy` stack contract.
- Byte slice aliasing
	- When creating instruction blocks (e.g., predicate blocks), copy the instruction slice before storing in constants to avoid accidental aliasing if scopes mutate.

### Phase 3 — Micro‑performance (guarded by benchmarks)

- Constant interning
	- Intern repeated strings (identifiers, property names) in the compiler constant pool to reduce duplicate allocations and cache pressure.
- Preallocation and reuse
	- Continue using `make(..., cap)` for hot slices (short‑circuit chains, function args, keys). Pre‑size when upper bounds are known.
	- VM stack sizing: ensure pre‑sized stack for typical programs; reuse between runs where safe.
- Keep interface churn low on hot paths by preferring concrete types where feasible (without semantic changes).

### Phase 4 — Optional deeper improvement (only if perf warrants)

- Typed constant pool
	- Replace `[]any` with a small tagged union for constants (number/string/bool/null/block). Benefits:
		- Removes many runtime type assertions in the VM.
		- Allows simpler, faster VM dispatch.
	- This is an internal contract; no public API change. Defer unless e2e benchmarks show VM type assertions on hot paths are a material cost.

## acceptance tests and quality gates

- Build: `go vet ./...` must pass with no new issues.
- Tests: `go test ./... -race` must pass; add tests for new error paths.
- Benchmarks: end‑to‑end benchmark must not regress (>2% regression triggers investigation). Track allocations/op.
- Lint (optional): run `staticcheck` on `compiler/` and `vm/` for additional signals.

## risks and mitigations

- Risk: Tightening type assertions might surface previously silent issues.
	- Mitigation: Gate changes behind tests; improve messages, not behavior.
- Risk: Instruction slice copying in predicate blocks could add allocations.
	- Mitigation: Measure; copy only when a block is emitted; amortized cost should be minimal.
- Risk: Interning could increase complexity.
	- Mitigation: Start with a simple string‑intern map; add tests; profile before/after.

## rollout plan

1) Land Phase 1 and 2 together (tiny, safe, mostly mechanical). Re‑run quality gates and e2e benchmarks.
2) If benchmarks show room, land Phase 3 tweaks incrementally (interning first).
3) Consider Phase 4 only if profiles indicate VM type assertions dominate.

## tracking checklist

- [ ] Callee type check in compiler function calls
- [ ] VM checked assertions for constants/system vars
- [ ] Centralized stack bounds checks
- [ ] `EmittedInstruction` rename and unexport (as applicable)
- [ ] Invariant comments for backpatch/pop
- [ ] Instruction block slice copy (no aliasing)
- [ ] String interning for constants
- [ ] End‑to‑end benchmark added and tracked

---

Notes
- This plan does not alter language semantics or opcode contracts.
- Each step must keep all tests green; any regression requires immediate rollback or fix.
