# Performance and Build Configuration

UExL is designed to be fast for the common case: pure-ASCII strings, simple field access, and short expressions. This page covers the runtime optimizations built into the Unicode string functions and how to enable optional SIMD acceleration.

## ASCII Fast-path

Every Unicode function in UExL (`runeLen`, `runeSubstr`, `graphemeLen`, `graphemeSubstr`, `runes`, `graphemes`) includes an ASCII fast-path. Before doing any Unicode segmentation, the function scans the string byte-by-byte for values ≥ 0x80. If none are found, the string is pure ASCII and all three Unicode levels (bytes, runes, graphemes) are identical — so the function returns immediately without invoking the UTF-8 decoder or the UAX #29 segmenter.

This matters in practice because most UExL strings are ASCII:

- Identifier names and field accessors: `user.firstName`, `order.total`
- Pipe expressions: `items |map: $item.price * 1.1`
- Numeric/boolean literals passed as strings
- Template keys and configuration values

For these, the fast-path brings Unicode-aware functions down to roughly the same cost as their byte-level counterparts.

### Complexity summary

| Input type | `len` / `substr` / `s[i]` / `s[i:j]` | `runeLen` / `runeSubstr` | `graphemeLen` / `graphemeSubstr` |
|---|---|---|---|
| ASCII string | O(1) or O(n) byte scan | O(n) ASCII scan → early exit | O(n) ASCII scan → early exit |
| Non-ASCII string | O(n) bytes | O(n) UTF-8 decode | O(n) UAX #29 segmentation |

## SIMD-Accelerated ASCII Detection (Go 1.26+)

When compiled with `GOEXPERIMENT=simd`, UExL uses the experimental `simd/archsimd` package (amd64 only, Go 1.26+) to replace the scalar byte-by-byte ASCII scan with a 16-byte-at-a-time SSE2 scan. This affects the hot inner loop of every Unicode function.

### How it works

The SIMD implementation loads 16 bytes at a time into an `Int8x16` register. Since ASCII bytes are 0x00–0x7F, they are non-negative as signed `int8`. Non-ASCII bytes (≥ 0x80) appear negative. A single `Less(zero).ToBits()` comparison checks all 16 bytes at once; any non-zero result means at least one non-ASCII byte is present and the scalar path takes over.

```
isASCII("hello, world!!!")  →  load 16 bytes → compare → ToBits() == 0 → ASCII ✓
```

Strings shorter than 16 bytes use `LoadInt8x16SlicePart`, which zero-pads the remaining lanes (zero is non-negative, so padding never produces false positives).

### Build configuration

| Build | Command | `isASCII` implementation |
|---|---|---|
| Standard (any Go ≥ 1.22) | `go build ./...` | Scalar byte scan |
| SIMD (Go 1.26+, amd64) | `GOEXPERIMENT=simd go build ./...` | 16-byte SSE2 scan |

The correct implementation is selected automatically by the Go build tag `goexperiment.simd` — **library consumers do not need to do anything**. If you build without the experiment, you get the scalar path; if you build with it, you get SIMD. Both produce identical results.

The minimum required Go version remains **1.22**. The `GOEXPERIMENT=simd` build requires Go 1.26 because that is the earliest version with the `simd/archsimd` package.

> **Stability note:** `simd/archsimd` is experimental and not covered by the Go 1 compatibility promise. Its API may change in future Go versions. UExL isolates it behind a build tag in `internal/utils/ascii_simd.go`; if the API breaks in a future Go release, only that file needs updating.

### Benchmark results (AMD Ryzen 7 5700G, amd64)

These numbers measure `isASCII` in isolation — a microbenchmark of the hot path called by every Unicode function.

| Input | Scalar | SIMD | Speedup |
|---|---|---|---|
| 5-byte ASCII (`hello`) | 4.0 ns | 8.8 ns | −2.2× (SIMD overhead dominates) |
| 11-byte ASCII (`hello world`) | 7.1 ns | 10.0 ns | −1.4× |
| 15-byte ASCII | 13.3 ns | 9.8 ns | **+26%** ← crossover |
| 16-byte ASCII | 12.9 ns | **6.5 ns** | **2.0×** |
| 29-byte ASCII (`items \|map: $item.price * 1.1`) | 19 ns | **10 ns** | **1.9×** |
| 45-byte ASCII (full sentence) | 32.8 ns | **12.0 ns** | **2.7×** |
| 88-byte ASCII (long expression) | 57.6 ns | **15.4 ns** | **3.7×** |
| Non-ASCII byte at position 2 | 2.9 ns | 4.5 ns | ≈ same (both exit early) |
| All non-ASCII, short | 2.0 ns | 7.8 ns | −3.9× (SIMD setup cost) |

**Crossover point:** ~14–15 bytes. SIMD wins for all strings ≥ 15 bytes of all-ASCII content, which is the common case for UExL field-access expressions and pipe predicates.

**Early-exit case:** When the first byte is non-ASCII (e.g. `données.prénom`), scalar and SIMD both exit after the first byte/chunk and are indistinguishable.

### End-to-end impact

The ASCII scan is one step inside `runeLen`, `graphemeLen`, etc. — it saves the UTF-8 decode and UAX #29 segmenter call, not the overall VM overhead. For a complete `graphemeLen("hello world")` expression including VM dispatch, the SIMD saving is a small fraction of the total. The benefit is most visible in tight loops over large arrays of ASCII strings:

```javascript
// Each element goes through the isASCII check; SIMD ~2× faster per element:
names |map: graphemeLen($item)
```

## Dependency

Grapheme functions (`graphemeLen`, `graphemeSubstr`, `graphemes`) depend on [`github.com/rivo/uniseg`](https://github.com/rivo/uniseg) (v0.4.7+) for UAX #29 segmentation. This is the only non-test external dependency added in v2. All other functions — including `runeLen`, `runeSubstr`, `runes`, `bytes`, `join` — use only the Go standard library.

The dependency is vendored (`vendor/github.com/rivo/uniseg`) and does not affect library consumers who vendor their dependencies.
