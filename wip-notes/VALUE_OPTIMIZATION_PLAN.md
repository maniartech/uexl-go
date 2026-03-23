# Value Struct Optimization Analysis

## Current Layout (56 bytes - WASTEFUL!)

```
Offset  Field      Size  Padding
------  ---------  ----  -------
0       Typ        1     7 bytes padding ❌
8       FloatVal   8
16      StrVal     16
32      BoolVal    1     7 bytes padding ❌
40      AnyVal     16
-------------------------------
Total: 56 bytes (14 bytes wasted!)
```

## Optimized Layout (40 bytes - BETTER!)

Reorder fields to minimize padding:
```
type Value struct {
	// 16-byte aligned fields first
	AnyVal   any      // offset 0, size 16
	StrVal   string   // offset 16, size 16
	FloatVal float64  // offset 32, size 8

	// Small fields together at end
	Typ      valueType // offset 40, size 1
	BoolVal  bool      // offset 41, size 1
	// 6 bytes padding to align to 8-byte boundary
}
-------------------------------
Total: 48 bytes (vs 56 current = 14% smaller!)
```

## Even Better: Union Approach (24 bytes - OPTIMAL!)

Use unsafe union pattern:
```go
type Value struct {
	typ valueType  // 1 byte
	_   [7]byte    // padding

	data [16]byte  // Union of:
	               // - float64 (8 bytes)
	               // - string (16 bytes)
	               // - bool (1 byte)
	               // - any (16 bytes)
}
```

Access via unsafe pointers:
```go
func (v *Value) FloatVal() float64 {
	return *(*float64)(unsafe.Pointer(&v.data[0]))
}
```

**Total: 24 bytes (57% smaller than current!)**

## Performance Impact Analysis

### Current: 56-byte Value
- Stack copy cost: 56 bytes per pop/push
- Cache line usage: fills 1 cache line (64 bytes) almost completely
- popValue taking 16% of CPU time

### Optimized: 40-byte Value
- Stack copy cost: 40 bytes per pop/push (29% faster)
- Expected speedup: 15-20 ns/op reduction
- **Low risk, high reward** ✅

### Union: 24-byte Value
- Stack copy cost: 24 bytes per pop/push (57% faster)
- Expected speedup: 40-50 ns/op reduction
- **High risk, complex code** ⚠️

## Recommendation: Optimized Layout (40 bytes)

### Implementation Steps:
1. Reorder fields in `types/value.go`
2. Run all tests to ensure no breakage
3. Benchmark to measure improvement
4. Expected result: 210-215 ns/op (from 228 ns)

### Code Change:
```go
// types/value.go
type Value struct {
	AnyVal   any       // 16 bytes, offset 0
	StrVal   string    // 16 bytes, offset 16
	FloatVal float64   // 8 bytes, offset 32
	Typ      valueType // 1 byte, offset 40
	BoolVal  bool      // 1 byte, offset 41
	// implicit 6 bytes padding to 48
}
```

## Why Not Union Approach?

**Pros:**
- 57% smaller (24 vs 56 bytes)
- Potentially 50ns faster

**Cons:**
- Unsafe pointer manipulation
- Platform-specific behavior
- Complex maintenance
- Harder to debug
- May break reflection/JSON marshaling

**Verdict:** Not worth the complexity for 30-40ns gain when we already beat competitors on allocations.

## Final Strategy

1. ✅ **Phase 1**: Reorder fields → 40 bytes (safe, 15-20ns gain)
2. 🤔 **Phase 2**: Add inline directives (safe, 10-15ns gain)
3. ❌ **Phase 3**: Union approach (risky, 30-40ns gain) - **SKIP**

**Target after Phase 1+2: ~200 ns/op, 0 allocs/op**

Still slower than expr (113ns) but:
- expr has 1 alloc, we have 0 ✅
- We're faster on pipes (3x) ✅
- Better GC behavior ✅
- More maintainable code ✅
