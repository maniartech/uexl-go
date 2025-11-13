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
		{"42.9 & 7", "error"},   // Should error: non-integer operand
		{"-42.9 & 7", "error"},  // Should error: non-integer operand
		{"5.1 | 2.8", "error"},  // Should error: non-integer operand
		{"10.7 ~ 3.2", "error"}, // Should error: non-integer operand (changed from ^)

		// Zero and edge cases
		{"0.0 & 5", 0.0},
		{"5 & 0.0", 0.0},
		{"-0.0 | 5", 5.0}, // -0.0 becomes 0 when cast to int
	}

	// Split tests into those expecting errors and those expecting results
	var errorTests, okTests []vmTestCase
	for _, tc := range tests {
		if tc.expected == "error" {
			tc.expected = "bitwise operations require integerish operands (no decimals), got 42.9 and 7"
			if tc.input == "-42.9 & 7" {
				tc.expected = "bitwise operations require integerish operands (no decimals), got -42.9 and 7"
			}
			if tc.input == "5.1 | 2.8" {
				tc.expected = "bitwise operations require integerish operands (no decimals), got 5.1 and 2.8"
			}
			if tc.input == "10.7 ~ 3.2" { // changed from ^
				tc.expected = "bitwise operations require integerish operands (no decimals), got 10.7 and 3.2"
			}
			errorTests = append(errorTests, tc)
		} else {
			okTests = append(okTests, tc)
		}
	}
	if len(okTests) > 0 {
		runVmTests(t, okTests)
	}
	if len(errorTests) > 0 {
		runVmErrorTests(t, errorTests)
	}
}

// TestBitwiseShiftEdgeCases tests edge cases specifically for shift operations
func TestBitwiseShiftEdgeCases(t *testing.T) {
	tests := []vmTestCase{
		// Large shift amounts that could cause issues
		{"1 << 63", -9223372036854775808.0}, // Shift by 63 causes signed overflow
		{"1 << 64", "error"},                // Shift by 64 is out of range - correctly rejected
		{"8 >> 3", 1.0},                     // Normal right shift

		// Fractional shift amounts (truncated)
		{"8 << 2.9", "error"},  // Should error: non-integer shift amount
		{"16 >> 1.7", "error"}, // Should error: non-integer shift amount

		// Zero shift amounts
		{"42 << 0", 42.0},
		{"42 >> 0", 42.0},

		// Negative numbers in shifts
		{"-8 << 2", -32.0}, // -8 << 2 = -32
		{"-8 >> 2", -2.0},  // -8 >> 2 = -2 (arithmetic right shift)
	}

	// Split tests into those expecting errors and those expecting results
	var errorTests, okTests []vmTestCase
	for _, tc := range tests {
		if tc.expected == "error" {
			if tc.input == "8 << 2.9" {
				tc.expected = "bitwise operations require integerish operands (no decimals), got 8 and 2.9"
			}
			if tc.input == "16 >> 1.7" {
				tc.expected = "bitwise operations require integerish operands (no decimals), got 16 and 1.7"
			}
			if tc.input == "1 << 64" {
				tc.expected = "shift count 64 out of range [0, 63]"
			}
			errorTests = append(errorTests, tc)
		} else {
			okTests = append(okTests, tc)
		}
	}
	if len(okTests) > 0 {
		runVmTests(t, okTests)
	}
	if len(errorTests) > 0 {
		runVmErrorTests(t, errorTests)
	}
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
