# UExL Performance Optimization Documentation

## 🎯 Start Here

**New to performance optimization?** → **[0-optimization-guidelines.md](0-optimization-guidelines.md)**

This comprehensive guide contains:
- **System-wide optimization scope** (100+ targets across 10 categories)
- Daily development workflow
- Validation checklists
- Measurement protocols
- Documentation roadmap

---

## Overview

This directory documents the **COMPREHENSIVE SYSTEM-WIDE PERFORMANCE OVERHAUL** of UExL: Parser → Compiler → VM. This is **NOT** just about boolean expressions or strings - this is **military-grade performance optimization** covering:

- **Every operator** (arithmetic, comparison, logical, bitwise, unary, string, array)
- **Every expression type** (binary, unary, index, member access, function calls, pipes)
- **All 50+ built-in functions** in `vm/builtins.go`
- **VM core** (stack ops, frame management, instruction dispatch)
- **Memory management** (allocations, pooling, scope reuse)
- **Type system** (checking, dispatch, conversion, coercion)
- **Pipe operations** (map, filter, reduce, and 8+ more)

**Goal:** Transform UExL from **~10x slower** than industry leaders to **faster across EVERY operation** while maintaining zero allocations and zero panics.

**Current Status:** 62ns/op (boolean), targeting **20-35ns/op across ALL operations**
**Scope:** **100+ optimization targets** across 10 categories

See **[optimization-rollout-plan.md](optimization-rollout-plan.md#-complete-optimization-scope)** for the complete inventory.

---

## Documentation Structure

### **Essential Documents (Use Daily)**

0. **[0-optimization-guidelines.md](0-optimization-guidelines.md)** - **START HERE** - Your daily workflow guide
1. **[optimization-rollout-plan.md](optimization-rollout-plan.md)** - **PRIMARY ROADMAP** - Phase-by-phase implementation plan
2. **[OPTIMIZATION_SCOPE_SUMMARY.md](OPTIMIZATION_SCOPE_SUMMARY.md)** - **COMPLETE INVENTORY** - All 100+ optimization targets tracked
3. **[dos-and-donts.md](dos-and-donts.md)** - **QUICK REFERENCE** - Code patterns & decisions

### **Measurement & Analysis**

3. **[profiling-guide.md](profiling-guide.md)** - CPU profiling walkthrough (use before/after each optimization)
4. **[benchmarking-guide.md](benchmarking-guide.md)** - Benchmark best practices & statistical analysis

### **Reference Material**

5. **[optimization-techniques.md](optimization-techniques.md)** - Pattern library with reusable code examples
6. **[best-practices.md](best-practices.md)** - Philosophy, guidelines, design patterns
7. **[optimization-journey.md](optimization-journey.md)** - Historical record (update after each optimization)
8. **[pending-optimizations.md](pending-optimizations.md)** - Future work & research opportunities

## Quick Stats

### Before Optimizations (Baseline)
```
Boolean Expression: 106 ns/op, 0 allocs
vs expr:    ~20% slower
vs cel-go:  ~20% faster (but cel-go allocates)
```

### After Optimizations (Current)
```
Boolean Expression: 62 ns/op, 0 allocs/op
vs expr:    41% FASTER (expr: 105 ns/op, 1 alloc)
vs cel-go:  51% FASTER (cel-go: 127 ns/op, 1 alloc)

Performance Improvement: 41% faster than baseline (106ns → 62ns)
```

## Key Achievements

✅ **Zero Allocations** - Maintained throughout all optimizations
✅ **Sub-65ns Execution** - Approaching theoretical limits for bytecode VMs
✅ **Industry-Leading Performance** - Faster than expr and cel-go
✅ **Type Safety Preserved** - No compromises on robustness
✅ **Zero Panics** - All optimizations maintain error handling integrity

## Major Optimization Categories

### 1. Type System Optimizations (3% gain)
- Eliminated redundant type assertions
- Direct type dispatch in comparison operations

### 2. Context Variable Caching (10% gain)
- Pre-resolved context vars to slice for O(1) access
- Eliminated expensive map lookups in hot path

### 3. Smart Cache Invalidation (30% gain)
- Pointer-based cache invalidation
- Avoid rebuilding cache when map is reused

### 4. Map Operations (4% gain)
- Use built-in `clear()` instead of iterative deletion
- Conditional clearing to avoid empty map iteration

## How to Use This Documentation

### For Performance Work
1. Read **[optimization-journey.md](optimization-journey.md)** to understand what was done
2. Consult **[profiling-guide.md](profiling-guide.md)** to identify new bottlenecks
3. Apply patterns from **[optimization-techniques.md](optimization-techniques.md)**
4. Follow **[best-practices.md](best-practices.md)** and **[dos-and-donts.md](dos-and-donts.md)**

### For New Features
1. Check **[pending-optimizations.md](pending-optimizations.md)** for related work
2. Follow **[benchmarking-guide.md](benchmarking-guide.md)** to measure impact
3. Adhere to **[best-practices.md](best-practices.md)** guidelines

### For Code Review
1. Reference **[dos-and-donts.md](dos-and-donts.md)** for quick checks
2. Verify benchmarks per **[benchmarking-guide.md](benchmarking-guide.md)**
3. Ensure alignment with **[best-practices.md](best-practices.md)**

## Performance Philosophy

UExL's performance philosophy is built on three pillars:

1. **Zero-Cost Abstractions** - Pay only for what you use
2. **Data-Oriented Design** - Optimize for cache locality and memory access patterns
3. **Profile-Driven Optimization** - Never guess, always measure

See **[best-practices.md](best-practices.md)** for detailed philosophy discussion.

## Contributing

When adding new optimizations:
1. Profile before and after
2. Document the change in **[optimization-journey.md](optimization-journey.md)**
3. Add technique to **[optimization-techniques.md](optimization-techniques.md)** if novel
4. Update **[pending-optimizations.md](pending-optimizations.md)** if completing a planned item
5. Ensure benchmarks show improvement per **[benchmarking-guide.md](benchmarking-guide.md)**

## External Resources

- Main codebase: `../../uexl-go/`
- Comparison benchmarks: `../`
- Architecture docs: `../../uexl-go/wip-notes/`
- Design philosophy: `../../uexl-go/book/design-philosophy.md`

---

**Last Updated:** October 17, 2025
**Performance Status:** ✅ Industry-Leading (62ns/op, 0 allocs)
