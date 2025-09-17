# Review: Edge Cases in Bitwise Operations for PR #25

## Summary

PR #25 successfully implements NaN/Infinity handling for bitwise operations, correctly returning "bitwise requires finite integers" errors. However, through comprehensive testing, several **critical edge cases** have been identified that should be addressed to improve robustness and prevent runtime crashes.

## ✅ What Works Well

The current implementation correctly:

1. **NaN/Inf Detection**: Properly detects and rejects NaN and Infinity values before bitwise operations
2. **Error Message**: Returns clear, consistent error message: "bitwise requires finite integers"
3. **IEEE-754 Compliance**: Maintains proper IEEE-754 semantics for arithmetic operations
4. **Basic Bitwise Operations**: Handles normal bitwise operations with finite numbers correctly

## ⚠️ Critical Edge Cases Identified

### 1. **Negative Shift Amount Panics** (🚨 HIGH PRIORITY)

**Issue**: Negative shift amounts cause **runtime panics** that crash the VM:

```go
// These expressions cause RUNTIME PANICS, not graceful errors:
"8 << -1"    // panic: runtime error: negative shift amount
"16 >> -2"   // panic: runtime error: negative shift amount
"1 << -1e20" // panic: runtime error: negative shift amount
```

**Root Cause**: The current implementation doesn't validate shift amounts before casting to `int`:

```go
case code.OpShiftLeft:
    vm.Push(float64(int(leftValue) << int(rightValue))) // Panics if rightValue < 0
```

**Recommendation**: Add validation for shift operations:

```go
if isBitwiseOp && isNanOrInf {
    return fmt.Errorf("bitwise requires finite integers")
}

// Add this check for shift operations specifically
if (operator == code.OpShiftLeft || operator == code.OpShiftRight) {
    if rightValue < 0 {
        return fmt.Errorf("shift amount must be non-negative")
    }
}
```

### 2. **Silent Integer Overflow** (⚠️ MEDIUM PRIORITY)

**Issue**: Large finite numbers silently overflow when cast to `int`, producing unexpected results:

```go
"1e20 & 1"    // Result: 0 (1e20 overflows to MinInt64, which is even)
"1e19 & 1"    // Result: 0 (also overflows)
"1e20 | 0"    // Result: -9223372036854775808 (MinInt64)
```

**Analysis**: 
- Numbers like `1e20` are finite (pass the NaN/Inf check)
- But when cast to `int`, they overflow to `math.MinInt64`
- This produces mathematically incorrect results without any warning

**Recommendation**: Consider adding range validation:

```go
if isBitwiseOp && isNanOrInf {
    return fmt.Errorf("bitwise requires finite integers")
}

// Add range validation for bitwise operations
if isBitwiseOp {
    if leftValue > float64(math.MaxInt64) || leftValue < float64(math.MinInt64) ||
       rightValue > float64(math.MaxInt64) || rightValue < float64(math.MinInt64) {
        return fmt.Errorf("bitwise operands exceed integer range")
    }
}
```

### 3. **Precision Loss Beyond 2^53** (ℹ️ LOW PRIORITY)

**Issue**: Float64 cannot represent integers exactly beyond 2^53, leading to precision loss:

```go
"9007199254740993 & 1"  // Expected: 1, Actual: 0
// 9007199254740993 rounds to 9007199254740992 in float64
```

**Analysis**: This is a fundamental limitation of IEEE-754 double precision, but users might not expect it.

**Recommendation**: Document this behavior or consider adding a warning for numbers beyond safe integer precision.

### 4. **Large Shift Amount Edge Cases** (ℹ️ LOW PRIORITY)

**Documented Behavior**: The current implementation exhibits Go's native shift behavior:

```go
"1 << 63"  // Result: -9223372036854775808 (signed overflow)
"1 << 64"  // Result: 0 (wraps around)
```

This behavior is consistent with Go but might be unexpected for users. Consider documenting this in the IEEE-754 semantics.

## 🔧 Recommended Implementation

Here's the suggested enhanced validation:

```go
func (vm *VM) executeBinaryArithmeticOperation(operator code.Opcode, left, right any) error {
    leftValue := left.(float64)
    rightValue := right.(float64)
    
    // Existing NaN/Inf check
    isNanOrInf := math.IsNaN(leftValue) || math.IsInf(leftValue, 0) || 
                  math.IsNaN(rightValue) || math.IsInf(rightValue, 0)
    
    isBitwiseOp := operator == code.OpBitwiseAnd || operator == code.OpBitwiseOr || 
                   operator == code.OpBitwiseXor || operator == code.OpShiftLeft || 
                   operator == code.OpShiftRight
    
    if isBitwiseOp && isNanOrInf {
        return fmt.Errorf("bitwise requires finite integers")
    }
    
    // NEW: Additional validations for bitwise operations
    if isBitwiseOp {
        // Check for negative shift amounts (prevents runtime panics)
        if (operator == code.OpShiftLeft || operator == code.OpShiftRight) && rightValue < 0 {
            return fmt.Errorf("shift amount must be non-negative")
        }
        
        // Optional: Check for integer overflow (prevents silent overflow)
        if leftValue > float64(math.MaxInt64) || leftValue < float64(math.MinInt64) ||
           rightValue > float64(math.MaxInt64) || rightValue < float64(math.MinInt64) {
            return fmt.Errorf("bitwise operands exceed integer range")
        }
    }
    
    // ... rest of switch statement
}
```

## 🧪 Test Coverage

Comprehensive edge case tests have been added in `vm/bitwise_edge_cases_test.go` covering:

- Integer overflow scenarios
- Precision loss cases
- Shift operation edge cases
- Negative shift amount documentation
- Fractional number truncation behavior

## 📋 Summary of Recommendations

1. **🚨 CRITICAL**: Add negative shift amount validation to prevent runtime panics
2. **⚠️ CONSIDER**: Add integer range validation to prevent silent overflow
3. **ℹ️ DOCUMENT**: Document precision loss limitations and large shift behaviors
4. **✅ MAINTAIN**: Keep existing NaN/Inf handling (works correctly)

The current implementation is solid for its intended use case, but these edge cases could cause unexpected behavior or crashes in production environments.