# Tokenizer Upgrade Plan

## Review findings (2025-08-28)

### Strengths

- Allocation-aware hot paths: ASCII fast paths, cached rune/size, reusable `strBuf`, no regex.
- Manual string unescape for double/single quotes with `\u`/`\U` support; fallback only on invalid sequences.
- Fast number parsing without exponent; precomputed `pow10` table; correct line/column tracking.
- Clean operator and pipe handling; per-instance state only (race-safe).

### Recommended improvements (priority order)

1. Unify `peek` and `peekNext` (keep one helper) to reduce duplication and drift.
1. Make pipe name scan truly byte-only in `readPipeOrBitwiseOr` (ASCII `[A-Za-z]+:`) to avoid rune decodes.
1. Provide a way to avoid retaining large input via `Token.Token`:
	- Store start/end offsets and derive lexeme on demand, or
	- Offer a configurable “no-retain” mode for long-lived token streams.
1. Identifier scanning: extend the ASCII byte loop, fallback to `advance()` only for non-ASCII to shave more cost.
1. Add `Reset(input string)` to reuse a `Tokenizer` instance and keep buffers hot.
1. Tests: add coverage for mixed `\u`/`\U`, invalid surrogate cases, `\U0010FFFF` boundary, and raw-strings with many doubled quotes.
1. Minor nits: remove redundant inner bound check in `skipWhitespace`; add brief doc comments to exported types/helpers; document policy for `.5` (dot-leading) numeric literals in `LANGUAGE.md`.

### Notes and expectations

- Escaped string literals inherently allocate once to materialize the final string; recent changes minimize extra temporaries. Further reductions require lazy-string representations or different semantics.
- `isDigit` is ASCII-only by design for speed; non-ASCII digits are treated as non-digits.
- Unary sign separation and exponent-in-number behavior are enforced by tests.

### Actionable checklist

- [x] Remove `peekNext` and update call sites to `peek`.
- [ ] Switch pipe-name scan to a pure byte loop.
- [ ] Add a token lexeme-retention strategy (offsets or configurable no-retain mode).
- [ ] Expand ASCII-first identifier loop; fallback only on non-ASCII.
- [ ] Add `Reset(input string)` method on `Tokenizer`.
- [ ] Add tests for `\u`/`\U` mix, invalid surrogates, `\U0010FFFF`, and dense doubled quotes in raw strings.
- [ ] Trim redundant check in `skipWhitespace` and add doc comments.
- [ ] Clarify the `.5` numeric literal policy in documentation.
