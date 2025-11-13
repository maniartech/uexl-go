# UExL Design Documentation

## Overview

This directory contains comprehensive design documentation for UExL's architecture, performance optimizations, and development guidelines.

## Documentation Structure

### Performance & Optimization

#### [value-migration/](value-migration/) - **Zero-Allocation Architecture** (Completed ✅)
**Status**: Complete - November 2025
**Achievement**: 266ns/4allocs → 227ns/0allocs

Documents the Value type system migration that achieved zero allocations:
- **[README.md](value-migration/README.md)** - Overview and quick start
- **[value-type-system.md](value-migration/value-type-system.md)** - Architecture deep dive
- **[performance-optimization.md](value-migration/performance-optimization.md)** - Complete journey
- **[development-guidelines.md](value-migration/development-guidelines.md)** - Critical rules

**Key Achievement**: Only framework in benchmarks with 0 allocs/op

#### [performance/](performance/) - **Comprehensive Optimization Plan** (Future Work 🔴)
**Status**: Planning complete, implementation pending
**Target**: 227ns → 20-35ns across ALL operations

100+ optimization targets across 10 categories:
- **[README.md](performance/README.md)** - Overview and navigation
- **[0-optimization-guidelines.md](performance/0-optimization-guidelines.md)** - Daily workflow guide
- **[optimization-rollout-plan.md](performance/optimization-rollout-plan.md)** - Phase-by-phase plan
- **[OPTIMIZATION_SCOPE_SUMMARY.md](performance/OPTIMIZATION_SCOPE_SUMMARY.md)** - Complete inventory
- **[dos-and-donts.md](performance/dos-and-donts.md)** - Quick reference patterns
- **[profiling-guide.md](performance/profiling-guide.md)** - CPU profiling walkthrough
- **[benchmarking-guide.md](performance/benchmarking-guide.md)** - Benchmark best practices
- **[optimization-techniques.md](performance/optimization-techniques.md)** - Pattern library
- **[best-practices.md](performance/best-practices.md)** - Philosophy & guidelines
- **[optimization-journey.md](performance/optimization-journey.md)** - Historical record
- **[pending-optimizations.md](performance/pending-optimizations.md)** - Future work

**Scope**: VM core, operators, pipes, built-ins, compiler, type system, memory management

---

## Quick Navigation

### For Current Development
**Maintaining zero allocations**: [value-migration/development-guidelines.md](value-migration/development-guidelines.md)

### For Performance Work
**Planning optimizations**: [performance/0-optimization-guidelines.md](performance/0-optimization-guidelines.md)

### For Architecture Understanding
**Value type system**: [value-migration/value-type-system.md](value-migration/value-type-system.md)

---

## Performance Status

### Current State (November 2025)

```
Expression: (Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)

UExL:   227.1 ns/op    0 B/op    0 allocs/op  ← ZERO ALLOCATIONS ✅
expr:   132.1 ns/op   32 B/op    1 allocs/op
celgo:  174.1 ns/op   16 B/op    1 allocs/op
```

### Completed Work
- ✅ **Zero allocations** achieved (Value migration)
- ✅ **14.7% speed improvement** (266ns → 227ns)
- ✅ **3× faster pipes** than competitors
- ✅ **Type-safe architecture** maintained
- ✅ **Zero panics** preserved

### Planned Work
- 🔴 VM core optimizations (instruction dispatch, stack ops, frames)
- 🔴 Operator optimizations (arithmetic, string, bitwise)
- 🔴 Pipe optimizations (filter, reduce, etc.)
- 🔴 Built-in function optimizations (50+ functions)
- 🔴 Compiler optimizations (constant folding, type hints)

**Target**: 20-35ns/op across ALL operations

---

## Contributing

### Adding New Features
1. Read [value-migration/development-guidelines.md](value-migration/development-guidelines.md)
2. Follow zero-allocation patterns
3. Test allocations (must be 0 for primitives)
4. Profile if performance-sensitive

### Performance Optimization
1. Review [performance/0-optimization-guidelines.md](performance/0-optimization-guidelines.md)
2. Profile to identify bottlenecks
3. Follow systematic rollout plan
4. Document results in optimization-journey.md

### Code Review Checklist
- [ ] No `Pop()` in opcode handlers (use `popValue()`)
- [ ] No boxing in hot paths (use Value directly)
- [ ] Allocations tested (0 for primitives)
- [ ] Benchmarks show no regression
- [ ] Documentation updated

---

## Document Maintenance

### value-migration/
**Frozen**: This documentation is complete and should only be updated for:
- Critical corrections
- New pitfalls discovered
- Breaking changes to Value system (rare)

### performance/
**Active**: Update as optimizations are implemented:
- Mark targets complete in OPTIMIZATION_SCOPE_SUMMARY.md
- Document results in optimization-journey.md
- Update benchmark baselines
- Add new techniques to optimization-techniques.md

---

## External Resources

- **Main codebase**: `../`
- **Comparison benchmarks**: `../../golang-expression-evaluation-comparison/`
- **User documentation**: `../book/`
- **WIP notes**: `../wip-notes/`

---

**Last Updated:** November 13, 2025
**Status**: Value migration complete ✅, comprehensive optimization planned 🔴
