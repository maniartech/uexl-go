// Package parser implements tokenization and parsing for the uexl expression language.
//
// Design goals:
//   - Correctness first: a comprehensive test suite covers grammar and edge cases.
//   - Performance: ASCII fast paths, cached rune decoding, and minimal allocations.
//   - Clear errors: structured error types with line/column context.
//
// The tokenizer is optimized for high throughput and low allocations. It maintains
// a cached current rune and size to avoid repeated utf8.DecodeRuneInString calls,
// and provides ASCII fast paths for common code paths. The parser is a hand-written
// recursive-descent parser that follows the language precedence and associativity
// rules, with careful handling of optional chaining, nullish coalescing and pipes.
package parser
