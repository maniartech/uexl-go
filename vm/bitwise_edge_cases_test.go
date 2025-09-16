package vm_test

import (
	"math"
	"testing"
)

// TestBitwiseEdgeCasesFloatToIntConversion tests edge cases in float64 to int conversion
// These tests reveal potential issues with the current implementation  
func TestBitwiseEdgeCasesFloatToIntConversion(t *testing.T) {
	tests := []vmTestCase{
		// Large finite numbers that overflow when cast to int
		{"1e20 & 1", 0.0}, // 1e20 overflows to MinInt64, MinInt64 & 1 = 0
		{"1 & 1e20", 0.0}, // Same but operands swapped
		
		// Numbers near integer limits that also overflow
		{"9223372036854776000 & 1", 0.0}, // Near MaxInt64, overflows to MinInt64, & 1 = 0
		
		// Precision loss in float64: numbers beyond 2^53 lose integer precision
		{"9007199254740992 & 1", 0.0}, // 2^53, even number so & 1 = 0
		{"9007199254740993 & 1", 0.0}, // 2^53 + 1 rounds to 2^53 in float64, so still even
		
		// Fractional numbers (current behavior: truncation toward zero)
		{"42.9 & 7", 2.0},   // 42 & 7 = 2
		{"-42.9 & 7", 6.0},  // -42 & 7 = 6 (truncates toward zero)
		{"5.1 | 2.8", 7.0},  // 5 | 2 = 7
		{"10.7 ^ 3.2", 9.0}, // 10 ^ 3 = 9
		
		// Zero and edge cases
		{"0.0 & 5", 0.0},
		{"5 & 0.0", 0.0},
		{"-0.0 | 5", 5.0}, // -0.0 becomes 0 when cast to int
	}

	runVmTests(t, tests)
}

// TestBitwiseShiftEdgeCases tests edge cases specifically for shift operations
func TestBitwiseShiftEdgeCases(t *testing.T) {
	tests := []vmTestCase{
		// Large shift amounts that could cause issues
		{"1 << 63", -9223372036854775808.0}, // Shift by 63 causes signed overflow
		{"1 << 64", 0.0}, // Shift by 64 wraps around in Go to 0
		{"8 >> 3", 1.0},  // Normal right shift
		
		// Fractional shift amounts (truncated)
		{"8 << 2.9", 32.0}, // 8 << 2 = 32 (2.9 truncated to 2)
		{"16 >> 1.7", 8.0}, // 16 >> 1 = 8 (1.7 truncated to 1)
		
		// Zero shift amounts
		{"42 << 0", 42.0},
		{"42 >> 0", 42.0},
		
		// Negative numbers in shifts
		{"-8 << 2", -32.0}, // -8 << 2 = -32
		{"-8 >> 2", -2.0},  // -8 >> 2 = -2 (arithmetic right shift)
	}

	runVmTests(t, tests)
}

// TestBitwiseOverflowBehavior documents the current overflow behavior
func TestBitwiseOverflowBehavior(t *testing.T) {
	tests := []vmTestCase{
		// These demonstrate the current overflow behavior which may be unexpected
		{"1e20 | 0", float64(math.MinInt64)}, // Large number overflows to MinInt64
		{"1e19 & 1", 0.0},                    // This number also overflows, & 1 = 0
	}

	runVmTests(t, tests)
}