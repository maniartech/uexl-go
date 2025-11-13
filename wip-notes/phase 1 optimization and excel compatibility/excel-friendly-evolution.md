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

**Note on Case-Insensitive Literals:**
Originally planned but **REMOVED** due to performance concerns. All language literals (`true`, `false`, `null`, `NaN`, `Inf`) remain **case-sensitive** to avoid string comparison overhead in hot paths.

## Implementation Checklist

> **⚠️ CRITICAL: Apply Optimization Patterns Throughout**
>
> Based on optimization-progress-tracker.md lessons learned:
> - ✅ **Type-specific function signatures** eliminate 30-44% overhead (proven in arithmetic/string ops)
> - ✅ **Type-specific push methods** (`pushFloat64`, `pushString`, `pushBool`) eliminate interface boxing
> - ✅ **Scope reuse patterns** reduce allocations by 96%+ (proven in pipe operations)
> - ✅ **Fast-path direct field access** eliminates 83% overhead (map string hashing bottleneck)
> - 🎯 **Target:** Maintain 0 allocs on boolean operations, minimize allocations elsewhere

---

### Phase 0: Preparation & Constants (30-45 minutes)

**Goal:** Set up symbols and constants correctly from the start

**Files to modify:**
- `parser/constants/language.go`
- `parser/constants/operators.go`

**Tasks:**
- [ ] **Add new symbol constants:**
  ```go
  // In parser/constants/language.go
  SymbolPower = "^"  // NEW - was XOR
  SymbolBitwiseXor = "~"  // CHANGED - was "^"
  SymbolBitwiseNot = "~"  // SAME symbol, context-dependent
  SymbolNotEqual2 = "<>"  // NEW - Excel alias
  ```

- [ ] **Update operator precedence constants if needed:**
  ```go
  // Verify precedence levels (likely no changes needed)
  PrecedencePower   = 10  // ^ now power (was 13 for XOR)
  PrecedenceBitXor  = 13  // ~ binary XOR (keep same level)
  PrecedencePrefix  = 11  // ~ unary NOT (keep same level)
  ```

- [ ] **Verify `isOperatorChar()` includes `~`:**
  ```go
  // In parser/constants/operators.go line 142
  // Should already include '~' - verify
  ch == '&' || ch == '|' || ch == '^' || ch == '~' || ...
  ```

**Validation:**
```bash
# Verify constants compile
go build ./parser/constants/
```

**Expected outcome:** All symbol constants ready, no breaking changes yet

---

### Phase 1: Parser/Tokenizer Changes (2-3 hours)

**Goal:** Recognize new operators and case-insensitive literals

**Files to modify:**
- `parser/tokenizer.go`
- `parser/parser.go`

#### 1.1: Tokenizer Operator Recognition (45 min)

**Tasks:**

- [ ] **Handle `<>` as not-equals alias:**
  ```go
  // In readOperator() after handling "!=" (around line 605)

  // Handle <> operator (Excel not-equals)
  if t.current() == '<' && t.peek() == '>' {
      t.advance()
      t.advance()
      operator := "!="  // Map to canonical form
      return Token{Type: constants.TokenOperator, Value: TokenValue{Kind: TVKOperator, Str: operator}, Token: "<>", Line: t.line, Column: startColumn}, nil
  }
  ```

- [ ] **Keep `^` as single-char operator** (no changes needed)
  - Already recognized in `isOperatorChar()` line 759
  - Will be reinterpreted by parser as power

- [ ] **Keep `~` as single-char operator** (already recognized)
  - Already in `isOperatorChar()` line 759
  - Parser will handle context-dependent logic

- [ ] **Case-insensitive keyword normalization:**
  ```go
  // In readIdentifierOrKeyword() after reading identifier (around line 385)

  // Normalize language literals to lowercase (Excel compatibility)
  lowerToken := strings.ToLower(tokenStr)
  switch lowerToken {
  case "true":
      return Token{Type: constants.TokenBoolean, Value: TokenValue{Kind: TVKBoolean, Bool: true}, Token: tokenStr, Line: t.line, Column: startColumn}, nil
  case "false":
      return Token{Type: constants.TokenBoolean, Value: TokenValue{Kind: TVKBoolean, Bool: false}, Token: tokenStr, Line: t.line, Column: startColumn}, nil
  case "null":
      return Token{Type: constants.TokenNull, Value: TokenValue{Kind: TVKNull}, Token: tokenStr, Line: t.line, Column: startColumn}, nil
  case "nan":
      if t.options.EnableIEEE754Specials {
          return Token{Type: constants.TokenNumber, Value: TokenValue{Kind: TVKNumber, Num: math.NaN()}, Token: tokenStr, Line: t.line, Column: startColumn}, nil
      }
  case "inf":
      if t.options.EnableIEEE754Specials {
          return Token{Type: constants.TokenNumber, Value: TokenValue{Kind: TVKNumber, Num: math.Inf(1)}, Token: tokenStr, Line: t.line, Column: startColumn}, nil
      }
  }

  // Continue with normal identifier handling
  ```

**Validation:**
```bash
# Test tokenizer changes
go test ./parser/tests -run TestTokenizer -v
```

#### 1.2: Parser Operator Precedence (45 min)

**Tasks:**

- [ ] **Update `parsePower()` to recognize `^` operator:**
  ```go
  // In parsePower() around line 294

  func (p *Parser) parsePower() Expression {
      left := p.parseUnary()

      if p.current.Type == constants.TokenOperator {
          if p.current.Value.Kind == TVKOperator {
              // Accept both ** and ^ for power (Excel compatibility)
              if p.current.Value.Str == "**" || p.current.Value.Str == "^" {
                  op := p.current
                  p.advance()
                  right := p.parsePower()  // Right-associative
                  return &BinaryExpression{Left: left, Operator: op.Value.Str, Right: right, Line: op.Line, Column: op.Column}
              }
          }
      }
      return left
  }
  ```

- [ ] **Update `parseUnary()` to handle `~` operator:**
  ```go
  // In parseUnary() around line 310

  func (p *Parser) parseUnary() Expression {
      if p.current.Type == constants.TokenOperator {
          if p.current.Value.Kind == TVKOperator {
              // Unary operators: -, !, ~ (NOT)
              if p.current.Value.Str == "-" || p.current.Value.Str == "!" || p.current.Value.Str == "~" {
                  op := p.current
                  opStr := op.Value.Str
                  p.advance()
                  expr := p.parseUnary()
                  return &UnaryExpression{Operator: opStr, Operand: expr, Line: op.Line, Column: op.Column}
              }
          }
      }
      return p.parseMemberAccess()
  }
  ```

- [ ] **Update `parseBitwiseXor()` to recognize `~` operator:**
  ```go
  // In parseBitwiseXor() around line 250

  func (p *Parser) parseBitwiseXor() Expression {
      // Now parses ~ for binary XOR (was ^)
      return p.parseBinaryOp(p.parseBitwiseAnd, "~")  // Changed from "^"
  }
  ```

- [ ] **Update `parseEquality()` to accept `<>` operator:**
  ```go
  // In parseEquality() around line 258

  func (p *Parser) parseEquality() Expression {
      // Accept both != and <> for not-equals
      return p.parseBinaryOp(p.parseComparison, constants.SymbolEqual, constants.SymbolNotEqual, "<>")
  }
  ```

**Validation:**
```bash
# Test parser changes
go test ./parser/tests -run TestParser -v

# Test specific cases
go test ./parser/tests -run "TestPower|TestBitwise|TestEquality" -v
```

**Expected outcome:** Parser recognizes all new operators correctly

---

### Phase 2: Compiler Changes (1-2 hours)

**Goal:** Emit correct opcodes for new operator mappings

**Files to modify:**
- `compiler/compiler.go`

#### 2.1: Binary Operator Compilation (45 min)

**Tasks:**

- [ ] **Update binary `^` to emit OpPow:**
  ```go
  // In compileBinaryExpression() around line 118

  case "^":
      c.emit(code.OpPow)  // Changed from OpBitwiseXor - Excel compatibility
  ```

- [ ] **Update binary `~` to emit OpBitwiseXor:**
  ```go
  // In compileBinaryExpression() around line 136

  case "~":
      c.emit(code.OpBitwiseXor)  // Moved from ^ - Lua-style context-dependent
  ```

- [ ] **Keep `**` emitting OpPow (no changes):**
  ```go
  // Already at line 118 - verify it's there
  case "**":
      c.emit(code.OpPow)
  ```

- [ ] **Handle `<>` as OpNotEqual:**
  ```go
  // In compileBinaryExpression() equality section

  case "==":
      c.emit(code.OpEqual)
  case "!=", "<>":  // Add "<>" to existing != case
      c.emit(code.OpNotEqual)
  ```

#### 2.2: Unary Operator Compilation (30 min)

**Tasks:**

- [ ] **Add unary `~` compilation:**
  ```go
  // In compileUnaryExpression() around line 155

  case "~":
      c.emit(code.OpBitwiseNot)  // NEW - context-dependent unary NOT
  ```

**Validation:**
```bash
# Test compiler changes
go test ./compiler/tests -run TestCompiler -v

# Verify bytecode output
go test ./compiler/tests -run "TestPower|TestBitwise" -v
```

**Expected outcome:** Compiler emits correct opcodes for all operators

---

### Phase 3: VM Changes with Optimization Patterns (2-3 hours)

**Goal:** Implement OpBitwiseNot handler using proven optimization patterns

**Files to modify:**
- `vm/vm_handlers.go`
- `vm/vm.go`

#### 3.1: Implement OpBitwiseNot Handler (1 hour)

> **🎯 CRITICAL: Follow type-specific pattern from arithmetic optimization**
>
> Lessons learned from optimization-progress-tracker.md:
> - Type-specific parameters eliminate 30-44% overhead
> - Type-specific push methods eliminate interface boxing
> - Result: 44.48% speed improvement for arithmetic

**Tasks:**

- [ ] **Create type-specific bitwise NOT function:**
  ```go
  // In vm/vm_handlers.go, add near executeNumberArithmetic()

  // executeBitwiseNot handles bitwise NOT with type-specific parameter
  // Follows optimization pattern: typed input → typed computation → typed push
  // Expected: 30-40% faster than any-based implementation
  func (vm *VM) executeBitwiseNot(value float64) error {
      // Validate integerish (no decimals)
      if value != math.Trunc(value) {
          return fmt.Errorf("bitwise operations require integerish operands (no decimals), got %v", value)
      }

      // Bitwise NOT operation
      intVal := int64(value)
      result := ^intVal

      // Use type-specific push (eliminates interface boxing)
      return vm.pushFloat64(float64(result))
  }
  ```

- [ ] **Add OpBitwiseNot case to VM opcode handler:**
  ```go
  // In vm/vm.go, in the main opcode switch (around line 131)

  case code.OpBitwiseNot:
      // Pop operand
      val := vm.pop()

      // Type check once, then use type-specific handler
      numVal, ok := val.(float64)
      if !ok {
          return fmt.Errorf("bitwise NOT requires number, got %T", val)
      }

      // Execute with type-specific function (optimization pattern)
      if err := vm.executeBitwiseNot(numVal); err != nil {
          return err
      }
  ```

- [ ] **Verify existing bitwise operations use type-specific pattern:**
  ```go
  // In vm_handlers.go executeNumberArithmetic() - already optimized!
  // Bitwise AND, OR, XOR, shifts all use pushFloat64()
  // No changes needed - pattern already applied ✅
  ```

#### 3.2: Update Opcode Dispatcher (30 min)

**Tasks:**

- [ ] **Ensure OpBitwiseNot is in unary operators list:**
  ```go
  // In vm/vm.go main switch, verify it's listed with other unary ops
  case code.OpNegate, code.OpNot, code.OpBitwiseNot:
      // Handle as unary (pops 1, pushes 1)
  ```

**Validation:**
```bash
# Test VM changes
go test ./vm -run TestBitwise -v

# Verify optimization (should be 0 allocations for boolean results)
go test -bench BenchmarkVM_Bitwise -benchmem

# Run all tests
go test ./... -v
```

**Expected outcome:**
- OpBitwiseNot works correctly
- Performance comparable to other optimized operations
- Zero allocations for operations returning booleans
- 30-40% faster than naive implementation

---

### Phase 4: Testing & Migration (2-3 hours)

**Goal:** Update all tests and create comprehensive test coverage

**Files to modify:**
- All `*_test.go` files using `^` operator
- New test files for Excel compatibility

#### 4.1: Update Existing Tests (1-2 hours)

**Strategy:** Use grep to find all uses of `^` and update

**Tasks:**

- [ ] **Find all test cases using `^` for XOR:**
  ```bash
  # Search all test files
  grep -r '\^' --include="*_test.go" vm/ parser/ compiler/

  # Expected locations based on earlier grep:
  # - vm/vm_test.go line 262, 267, 270, 278
  # - vm/ieee754_vm_test.go line 163, 164
  # - vm/bitwise_edge_cases_test.go line 27, 46
  # - parser/tests/tokenizer_coverage_test.go line 116
  # - parser/tests/parser_coverage_test.go line 52, 53
  # - compiler/tests/compiler_test.go line 177, 200, 207, 211, 231
  ```

- [ ] **Update VM tests:**
  ```go
  // In vm/vm_test.go, change all XOR tests:

  // OLD: {"5 ^ 3", 6.0}
  // NEW:
  {"5 ~ 3", 6.0},             // XOR: 0101 ~ 0011 = 0110
  {"15 ~ 7", 8.0},            // XOR: 1111 ~ 0111 = 1000
  {"(5 & 3) | (2 ~ 1)", 3.0}, // XOR in expression

  // Add new power tests:
  {"2 ^ 3", 8.0},             // Power (was XOR)
  {"5 ^ 2", 25.0},            // Power
  {"10 ^ 0", 1.0},            // Power edge case

  // Add bitwise NOT tests:
  {"~5", -6.0},               // NOT: ~0101 = ...11111010 (-6 in two's complement)
  {"~0", -1.0},               // NOT: ~0 = -1
  {"~~5", 5.0},               // Double NOT
  ```

- [ ] **Update parser tests:**
  ```go
  // In parser/tests/parser_coverage_test.go:

  // OLD: {"bitwise xor", "a ^ b ^ c"}
  // NEW:
  {"bitwise xor", "a ~ b ~ c"},
  {"power", "2 ^ 3"},
  {"power vs xor", "2 ^ 3 != 2 ~ 3"},  // Verify both work
  ```

- [ ] **Update compiler tests:**
  ```go
  // In compiler/tests/compiler_test.go:

  // Change all OpBitwiseXor expectations to use ~ not ^
  // Add OpPow expectations for ^ operator
  ```

#### 4.2: Create Excel Compatibility Test Suite (1 hour)

**Tasks:**

- [ ] **Create new test file `vm/excel_compat_test.go`:**
  ```go
  package vm

  import "testing"

  func TestExcelCompatibility(t *testing.T) {
      tests := []vmTestCase{
          // Power operator (Excel standard)
          {"2 ^ 3", 8.0},
          {"10 ^ 2", 100.0},
          {"5 ^ 0", 1.0},
          {"2 ^ -1", 0.5},

          // Bitwise XOR (Lua-style ~)
          {"5 ~ 3", 6.0},
          {"15 ~ 7", 8.0},

          // Bitwise NOT (unary ~)
          {"~5", -6.0},
          {"~0", -1.0},

          // Not-equals alias
          {"5 <> 3", true},
          {"5 <> 5", false},
          {"'a' <> 'b'", true},

          // Case-insensitive literals
          {"TRUE", true},
          {"True", true},
          {"FALSE", false},
          {"False", false},
          {"NULL", nil},
          {"Null", nil},

          // Mixed operations
          {"2 ^ 3 + 5", 13.0},           // Power precedence
          {"(2 + 3) ^ 2", 25.0},         // Parentheses
          {"5 ~ 3 & 1", 0.0},            // XOR then AND
          {"~5 | 3", -5.0},              // NOT then OR
      }

      runVmTests(t, tests)
  }

  func TestExcelCompatIEEE754(t *testing.T) {
      // Test case-insensitive NaN/Inf (requires EnableIEEE754Specials)
      // ... tests for NaN, Nan, nan, INF, Inf, inf
  }
  ```

**Validation:**
```bash
# Run new Excel compat tests
go test ./vm -run TestExcelCompat -v

# Run all tests
go test ./... -v

# Verify no regressions
go test ./... -race
```

**Expected outcome:**
- All existing tests updated and passing
- New Excel compatibility tests passing
- No test regressions

---

### Phase 5: Performance Benchmarking (1-2 hours)

**Goal:** Verify performance meets optimization targets

**Files to create:**
- `excel_benchmark_test.go`

#### 5.1: Create Excel Operator Benchmarks (30 min)

**Tasks:**

- [ ] **Create benchmark file:**
  ```go
  // File: excel_benchmark_test.go (in root or vm/)

  package vm

  import "testing"

  func BenchmarkExcel_Power(b *testing.B) {
      // Test ^ power operator performance
      runBenchmark(b, "2 ^ 10", map[string]any{})
      // Target: <100 ns/op, 0 allocs (boolean result)
  }

  func BenchmarkExcel_BitwiseXor(b *testing.B) {
      // Test ~ XOR operator performance
      runBenchmark(b, "flags ~ mask", map[string]any{
          "flags": 15.0,
          "mask": 7.0,
      })
      // Target: <120 ns/op (similar to arithmetic)
  }

  func BenchmarkExcel_BitwiseNot(b *testing.B) {
      // Test ~ NOT operator performance
      runBenchmark(b, "~value", map[string]any{"value": 5.0})
      // Target: <10 ns/op, 0 allocs (inline stack operation, must match other bitwise ops)
  }

  func BenchmarkExcel_NotEquals(b *testing.B) {
      // Test <> not-equals performance
      runBenchmark(b, "a <> b", map[string]any{"a": 5.0, "b": 3.0})
      // Target: <5 ns/op, 0 allocs (boolean result, already optimized)
  }

  func BenchmarkExcel_CaseInsensitive(b *testing.B) {
      // Test case-insensitive literal performance
      runBenchmark(b, "TRUE && FALSE", map[string]any{})
      // Target: Same as lowercase (no performance penalty)
  }

  func BenchmarkExcel_MixedOperations(b *testing.B) {
      // Test realistic Excel-like expression
      runBenchmark(b, "(a ^ 2 + b ^ 2) ~ mask <> 0", map[string]any{
          "a": 3.0,
          "b": 4.0,
          "mask": 7.0,
      })
      // Target: <300 ns/op
  }
  ```

#### 5.2: Run Performance Validation (30 min)

**Tasks:**

- [ ] **Establish baselines:**
  ```bash
  # Before changes (current performance)
  go test -bench BenchmarkExcel -benchmem -count=10 > excel_baseline.txt
  ```

- [ ] **After implementation:**
  ```bash
  # After all changes
  go test -bench BenchmarkExcel -benchmem -count=10 > excel_after.txt

  # Statistical comparison
  benchstat excel_baseline.txt excel_after.txt
  ```

- [ ] **Verify targets met:**
  ```
  Expected results (based on optimization patterns):

  ✅ Power (^):        <10 ns/op, 0 allocs (inline value operation)
  ✅ Bitwise XOR (~):  <10 ns/op, 0 allocs (matches existing bitwise ops)
  ✅ Bitwise NOT (~):  <10 ns/op, 0 allocs (unary bitwise operation)
  ✅ Not-equals (<>):  <5 ns/op, 0 allocs (boolean result, already optimized)

  CRITICAL: All operations MUST be zero allocations
  - Use direct stack manipulation (no interface boxing)
  - Use type-specific push methods (pushFloat64, pushBool)
  - Profile confirms 0 B/op, 0 allocs/op

  If any operation shows allocations or >20 ns/op:
  - Profile with -cpuprofile
  - Apply type-specific optimization pattern
  - Verify using pushFloat64/pushBool (not push(any))
  - Check for interface boxing in hot path
  ```

**Validation:**
```bash
# Full benchmark suite
go test -bench=. -benchmem -benchtime=3s

# Compare against competitors (if available)
cd ../golang-expression-evaluation-comparison/
./run_benchmarks.sh
```

**Expected outcome:**
- All operations meet performance targets
- No regressions vs current performance
- Maintain competitive advantage (2-30x faster than expr/cel-go)

---

### Phase 6: Documentation (1-2 hours)

**Goal:** Update all documentation for Excel compatibility

**Files to modify:**
- `book/operators/*.md`
- `book/syntax.md`
- `README.md`

#### 6.1: Update Operator Documentation

**Tasks:**

- [ ] **Update power operator docs:**
  ```markdown
  ## Power Operator

  **Symbols:** `^` (Excel-compatible), `**` (alternative)

  **Precedence:** 10 (right-associative)

  **Examples:**
  ```
  2 ^ 3     // 8 (Excel standard)
  2 ** 3    // 8 (also works)
  10 ^ -2   // 0.01
  ```

  **Migration:** In UExL <2.0, `^` was bitwise XOR. Use `~` for XOR now.
  ```

- [ ] **Update bitwise operator docs:**
  ```markdown
  ## Bitwise Operators

  **Symbols:**
  - `&` - Bitwise AND
  - `|` - Bitwise OR
  - `~` - Bitwise XOR (binary), Bitwise NOT (unary)
  - `<<` - Left shift
  - `>>` - Right shift

  **Context-Dependent ~ (Lua-style):**
  ```
  5 ~ 3   // 6 (binary XOR)
  ~5      // -6 (unary NOT)
  ```

  **Migration:** In UExL <2.0, `^` was XOR. Now use `~`.
  ```

- [ ] **Update comparison operator docs:**
  ```markdown
  ## Not-Equals Operators

  **Symbols:** `!=` (standard), `<>` (Excel-compatible)

  **Examples:**
  ```
  5 != 3    // true (standard)
  5 <> 3    // true (Excel alias)
  "a" <> "b" // true
  ```
  ```

- [ ] **Add Excel compatibility section to README:**
  ```markdown
  ## Excel Compatibility

  UExL is designed to be Excel-compatible where it matters:

  - ✅ `^` is power operator (like Excel, not XOR)
  - ✅ `<>` is not-equals (Excel alias for `!=`)
  - ✅ Case-insensitive literals: `TRUE`, `False`, `NULL`, `NaN`, `Inf`
  - ✅ Familiar precedence and semantics

  For bitwise operations, UExL uses Lua-style `~` operator:
  - `5 ~ 3` - XOR
  - `~5` - NOT
  ```

#### 6.2: Create Migration Guide

**Tasks:**

- [ ] **Create `book/migration/v1-to-v2.md`:**
  ```markdown
  # Migration Guide: UExL v1.x to v2.0

  ## Breaking Changes

  ### Power Operator (`^` changed from XOR to power)

  **Impact:** HIGH - any code using `^` for XOR will break

  **Before (v1.x):**
  ```
  5 ^ 3        // XOR = 6
  flags ^ mask // Flip bits
  ```

  **After (v2.0):**
  ```
  5 ~ 3        // XOR = 6 (use tilde)
  flags ~ mask // Flip bits (use tilde)
  5 ^ 3        // POWER = 125 (CHANGED!)
  ```

  **Auto-migration:**
  ```bash
  # Find all uses of ^ in expressions
  grep -r '\^' your_project/

  # Manual review required - distinguish XOR from power intent
  # XOR: replace ^ with ~
  # Power: leave as ^ (but verify it's intentional power, not XOR!)
  ```

  ## Non-Breaking Additions

  ### Bitwise NOT (`~` unary operator)
  - NEW in v2.0, was not implemented in v1.x
  - `~5` → -6 (bitwise NOT)

  ### Excel Aliases
  - `<>` now accepted as not-equals (same as `!=`)
  - Case-insensitive: `TRUE`, `FALSE`, `NULL`, `NaN`, `Inf`

  ## Everything Else Works the Same
  - Bitwise AND (`&`), OR (`|`) - no changes
  - Power operator `**` - still works
  - All other operators - no changes
  ```

**Expected outcome:**
- Complete documentation for all changes
- Clear migration path for existing users
- Excel users have quickstart guide

---

## Summary: Implementation Phases with Time Estimates

| Phase | Focus | Time | Validation |
|-------|-------|------|------------|
| **Phase 0** | Preparation & Constants | 30-45 min | `go build ./parser/constants/` |
| **Phase 1** | Parser/Tokenizer | 2-3 hours | `go test ./parser/tests -v` |
| **Phase 2** | Compiler | 1-2 hours | `go test ./compiler/tests -v` |
| **Phase 3** | VM (Optimized) | 2-3 hours | `go test ./vm -v` + benchmarks |
| **Phase 4** | Testing | 2-3 hours | `go test ./... -v -race` |
| **Phase 5** | Performance | 1-2 hours | `benchstat` comparisons |
| **Phase 6** | Documentation | 1-2 hours | Manual review |
| **TOTAL** | **All Phases** | **10-16 hours** | Full test suite + benchmarks |

**Total Estimated Time:** 10-16 hours (1.5 - 2 days)

---

## Success Criteria Checklist

**Functionality:**
- ✅ `^` compiles to OpPow (power operator)
- ✅ `~` binary compiles to OpBitwiseXor
- ✅ `~` unary compiles to OpBitwiseNot (NEW implementation)
- ✅ `<>` works as not-equals
- ✅ All existing tests updated and passing
- ✅ New Excel compatibility tests passing

**Performance (Based on Optimization Patterns):**
- ✅ OpBitwiseNot uses type-specific pattern (30-40% faster)
- ✅ All bitwise ops use `pushFloat64()` (eliminate boxing)
- ✅ Power operator comparable to other arithmetic (100-120 ns/op)
- ✅ Zero allocations on boolean operations
- ✅ No performance regressions

**Quality:**
- ✅ Statistical significance: p < 0.05 on benchmarks
- ✅ Zero panics, robust error handling
- ✅ All race conditions checked: `go test ./... -race`
- ✅ Documentation complete and accurate

**Competitive Position:**
- ✅ Maintain 2-30x performance advantage over expr/cel-go
- ✅ Excel compatibility achieved without performance penalty
- ✅ Best-in-class optimization patterns applied throughout

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

// NOTE: Literals remain case-sensitive for performance
true             // Valid
false            // Valid
null             // Valid
TRUE             // Invalid - not recognized
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
