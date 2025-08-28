# Tokenizer Upgrade Plan

This document tracks the performance refactor of `parser/tokenizer.go` toward lower latency and allocations without changing public behavior.

## Goals

- Keep semantic value parsing (numbers, strings, booleans, null) intact.
- Minimize allocations in hot paths (numbers, identifiers, strings, pipes).
- Ensure race-safety and deterministic behavior.
- Provide benchmarks to quantify improvements.

## Constraints

- No API breaking changes for `Tokenizer` and `Token`.
- Maintain error messages, positions (Line/Column), and tokenization semantics.

## Phases & TODOs

1. Strings fast path and unescape

- [ ] Add manual double-quoted unescape (ASCII fast path; handle `\\`, `\"`, `\n`, `\t`, `\r`, `\uXXXX`, `\UXXXXXXXX`).
- [ ] Only fall back to `strconv.Unquote` for rare/invalid cases.
- [ ] Keep raw-string doubled-quote collapse using single pass with reusable buffer.

1. Pipe scanning and operators

- [ ] Ensure `readPipeOrBitwiseOr` uses pure byte scan for `[A-Za-z]+:` (no rune decode) and returns correct columns.
- [ ] Validate `?., ?[, ??` and operator pairs remain correct; add testcases if missing.

1. Number parsing

- [ ] Retain manual parse for ints/decimals without exponent; fallback to `ParseFloat` otherwise.
- [ ] Ensure column/rune cache updated correctly after ASCII fast path loops.
- [x] Unary signs are not part of numeric tokens: `-123` → `-` and `123`, `+456` → `+` and `456`.
- [x] Exponent signs are part of numbers: `1e-5` is a single number token.
	- Acceptance tests live in `parser/tokenizer_signs_test.go`.

1. Whitespace and identifier scan

- [ ] Keep ASCII fast path; fallback to unicode for non-ASCII.
- [ ] Confirm `$` identifiers and underscores remain valid.

1. Benchmarks & acceptance criteria

- [ ] Add/keep comprehensive benchmarks in `parser/bench_tokenizer_test.go` covering:
	- Scalars/ops, identifiers, strings, pipes, nullish/optional, unicode heavy, long numbers.
- [ ] Target improvements (Windows dev box as reference):
	- identifiers 64B: ≤ 500ns/op, ≤ 8 allocs/op
	- strings 128B: ≤ 1.2µs/op, ≤ 18 allocs/op
	- pipes 128B: ≤ 0.85µs/op, ≤ 15 allocs/op
	- numbers 64B: ≤ 600ns/op, ≤ 9 allocs/op
- [ ] All parser tests and vm tests pass; parser tests pass with `-race`.

## Notes

- Prefer ASCII byte loops where possible; update cached rune via `setCur()` only when position advanced by rune-size.
- Avoid regex and large temporary strings in hot paths.

## Done (to update as we progress)

- [x] Rune caching and ASCII fast paths.
- [x] Manual pipe parsing (no regex).
- [x] Unary sign separation policy and tests.
