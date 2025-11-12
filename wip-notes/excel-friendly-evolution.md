# Excel-Friendly Language Evolution

**Goal:** Make UExL Excel-compatible by default without maintaining dual modes.

**Philosophy:** "Excel-compatible where it matters, modern where it helps"

## Breaking Changes (Required for Excel Compatibility)

### 1. Power Operator Remapping ⚠️ CRITICAL

**Current State:**
- `^` = Bitwise XOR (OpBitwiseXor) ← Working, tested
- `**` = Power (OpPow) ← Working, tested
- `~` = Bitwise NOT (OpBitwiseNot) ← Defined but NOT IMPLEMENTED (no VM handler)

**New State (Lua 5.3+ inspired):**
- `^` = Power (OpPow) ← Excel standard
- `**` = Power (OpPow) ← Keep as alias for backward compat
- `~` = Context-dependent:
  - **Binary**: Bitwise XOR (`5 ~ 3` → 6)
  - **Unary**: Bitwise NOT (`~5` → -6)

**Rationale for `~` approach:**
- ✅ Lua 5.3+ precedent (proven design)
- ✅ Single character (concise)
- ✅ Context-dependent parsing (like `-` for minus/negate)
- ✅ No new keywords needed
- ✅ Bitwise NOT already uses `~` in most languages

**Migration Path for Existing UExL Users:**
```
OLD: 5 ^ 3   → 6 (XOR)
NEW: 5 ~ 3   → 6 (XOR - use tilde)
     5 ^ 3   → 125 (POWER - breaking!)
     5 ** 3  → 125 (still works)

OLD: ~5      → (broken - not implemented)
NEW: ~5      → -6 (unary NOT - now works)
```

**Impact:** HIGH - Any existing UExL code using `^` for XOR will break
**Justification:** Silent wrong mathematical results are unacceptable

### 2. Bitwise Operators (Minimal Changes)

**Current State:**
- `&` = Bitwise AND (OpBitwiseAnd) ← Working
- `|` = Bitwise OR (OpBitwiseOr) ← Working
- `^` = Bitwise XOR (OpBitwiseXor) ← Working (needs to change)
- `~` = Bitwise NOT (OpBitwiseNot) ← Defined but NOT implemented

**New State (Keep & and | as-is):**
- `&` = Bitwise AND ← **NO CHANGE** (keep as bitwise)
- `|` = Bitwise OR ← **NO CHANGE** (keep as bitwise)
- `~` = Context-dependent:
  - **Binary**: Bitwise XOR (`5 ~ 3` → 6)
  - **Unary**: Bitwise NOT (`~5` → -6)
- `^` = Power (moved from XOR)

**String Concatenation:**
- Use `+` operator (already works in UExL)
- Excel users can use `&` via custom function if needed

**Rationale for keeping symbols:**
- ✅ Minimal breaking changes (only `^` changes)
- ✅ Standard C-family syntax for `&` and `|`
- ✅ Pipe operator `|:` already disambiguated by `:` suffix
- ✅ String concat with `+` is common (JavaScript, Python, etc.)
- ❌ Excel's `&` for concat is non-critical (users can adapt)

**Migration Path:**
```
OLD: 5 & 3  → 1 (bitwise AND)
NEW: 5 & 3  → 1 (NO CHANGE)

OLD: 5 ^ 3  → 6 (bitwise XOR)
NEW: 5 ~ 3  → 6 (use tilde for XOR)
     5 ^ 3  → 125 (POWER - breaking!)

OLD: ~5     → (broken - not implemented)
NEW: ~5     → -6 (unary NOT - now works)

String concat:
     "A" + "B"  → "AB" (already works)
```

### 3. Accept `<>` as Not-Equals (Additive)

**Current State:**
- `!=` = Not equals (OpNotEqual)

**New State:**
- `!=` = Not equals (still works)
- `<>` = Not equals (Excel alias)

**Impact:** NONE - purely additive

### 4. Case-Insensitive Language Literals (Additive)

**Current State:**
- `true`, `false`, `null` (lowercase only, case-sensitive)
- `NaN`, `Inf` (exact case, when IEEE specials enabled)

**New State - All case-insensitive:**
- Accept: `TRUE`, `True`, `true`, `tRuE`, etc.
- Accept: `FALSE`, `False`, `false`, `FaLsE`, etc.
- Accept: `NULL`, `Null`, `null`, `nUlL`, etc.
- Accept: `NAN`, `Nan`, `nan`, `NaN`, etc. (when IEEE enabled)
- Accept: `INF`, `Inf`, `inf`, `InF`, etc. (when IEEE enabled)

**User-defined remain case-sensitive:**
- Functions: `sum` ≠ `SUM` ≠ `Sum` (user must match case)
- Variables: `userName` ≠ `UserName`
- Pipe names: `customPipe` ≠ `CustomPipe`

**Impact:** NONE - purely additive, makes language more forgiving

## Implementation Checklist

### Phase 1: Parser Changes
- [ ] Tokenizer: `^` → TokenPower (was TokenBitwiseXor)
- [ ] Tokenizer: `~` binary → TokenBitwiseXor (context-dependent)
- [ ] Tokenizer: `~` unary → TokenBitwiseNot (context-dependent)
- [ ] Tokenizer: `<>` → TokenNotEqual (new alias)
- [ ] Tokenizer: Case-insensitive keyword normalization
- [ ] Update token constants in `parser/constants/tokens.go`

### Phase 2: Compiler Changes
- [ ] Compiler: `^` emits OpPow (was OpBitwiseXor)
- [ ] Compiler: `~` binary emits OpBitwiseXor (was from `^`)
- [ ] Compiler: `~` unary emits OpBitwiseNot (needs implementation)
- [ ] Compiler: `**` still emits OpPow (backward compat)
- [ ] Compiler: Keep `&` as OpBitwiseAnd (no change)
- [ ] Compiler: Keep `|` as OpBitwiseOr (no change)
- [ ] Update operator precedence if needed

### Phase 3: VM Changes
- [ ] VM: Implement OpBitwiseNot handler (currently missing)
- [ ] VM: Keep `&` as bitwise AND (no change)
- [ ] VM: Keep `|` as bitwise OR (no change)

### Phase 4: Testing
- [ ] Update all tests using `^` for XOR
- [ ] Add Excel compatibility test suite
- [ ] Add migration examples
- [ ] Performance benchmarks

### Phase 5: Documentation
- [ ] Update operator documentation
- [ ] Migration guide for existing UExL users
- [ ] Excel user quickstart guide
- [ ] Update examples

## Migration Guide for Existing UExL Users

### Breaking Change: `^` Operator (ONLY Breaking Change)

**Before:**
```uexl
5 ^ 3        // XOR = 6
flags ^ mask // Flip bits
```

**After:**
```uexl
5 ~ 3        // XOR = 6 (use tilde)
flags ~ mask // Flip bits (use tilde)
5 ^ 3        // POWER = 125 ⚠️ CHANGED!
```

**Auto-migration script:**
```bash
# Replace ^ with ~ where it's XOR (context-dependent)
# Requires manual review to distinguish XOR from power usage
sed -i 's/\([0-9a-zA-Z_)]\) \^ \([0-9a-zA-Z_(]\)/\1 ~ \2/g' *.uexl
```

### Non-Breaking: Everything Else Still Works

```uexl
// Bitwise operators (NO CHANGES)
5 & 3            // Bitwise AND = 1 (still works)
5 | 3            // Bitwise OR = 7 (still works)
~5               // Bitwise NOT = -6 (now works!)
5 ~ 3            // Bitwise XOR = 6 (new syntax)

// Power operators
2 ** 3           // Power = 8 (still works)
2 ^ 3            // Power = 8 (NEW - was XOR!)

// Other operators
x != y           // Not equals (still works)
x <> y           // Not equals (new Excel alias)
"Hello" + " "    // String concat (still works)
true && false    // Logical AND (still works)
```

## Risk Assessment

| Change | Risk Level | Impact | Mitigation |
|--------|-----------|--------|------------|
| `^` for power | 🔴 HIGH | Breaks XOR code (use `~` instead) | Migration script, clear docs |
| `~` for XOR/NOT | 🟡 MEDIUM | Implements unary `~` (was broken) | Only affects if users tried `~` before |
| `<>` not-equals | 🟢 LOW | None (additive) | None needed |
| Case-insensitive | 🟢 LOW | None (additive) | None needed |

## Timeline

**Week 1: Parser + Tokenizer**
- Implement operator changes
- Add case-insensitive keywords
- Update token constants

**Week 2: Compiler + VM**
- Remap opcodes
- Implement OpBitwiseNot handler
- Update handlers

**Week 3: Testing + Migration**
- Update test suite
- Write migration tools
- Performance validation
- Documentation

**Week 4: Release**
- Release notes
- Migration guide
- Examples
- Community communication

## Success Criteria

- ✅ Excel formulas paste and work correctly (^, &, <>)
- ✅ Existing UExL code has clear migration path
- ✅ Performance neutral or better
- ✅ All tests pass
- ✅ Documentation complete
- ✅ Zero silent failures

## Open Questions

1. ~~**Keep `|` for bitwise OR or only as pipe prefix?**~~
   - ✅ **RESOLVED:** Keep `|` as bitwise OR, pipe syntax `|:` already disambiguated

2. ~~**Keep `~` for bitwise NOT?**~~
   - ✅ **RESOLVED:** Use `~` for both XOR (binary) and NOT (unary), Lua-style

3. **Deprecation period for old syntax?**
   - Recommendation: Major version bump, no deprecation (clean break)

4. ~~**Should `+` still work for string concat?**~~
   - ✅ **RESOLVED:** Yes, `+` is primary string concat operator

## Future Considerations

- Add more Excel functions (VLOOKUP, SUMIF, etc.)
- Range syntax support (A1:A10)
- Cell reference helpers
- Excel error types (#N/A, #VALUE!, #DIV/0!)
