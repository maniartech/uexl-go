package vm_test

import (
	"math"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

// TestIEEE754Arithmetic tests arithmetic operations with NaN and Inf values
// following the IEEE-754 semantics defined in ieee754-semantics.md
func TestIEEE754Arithmetic(t *testing.T) {
	tests := []vmTestCase{
		// NaN propagation in arithmetic
		{"NaN + 1", math.NaN()},
		{"1 + NaN", math.NaN()},
		{"NaN + NaN", math.NaN()},
		{"NaN - 1", math.NaN()},
		{"1 - NaN", math.NaN()},
		{"NaN * 1", math.NaN()},
		{"1 * NaN", math.NaN()},
		{"NaN / 1", math.NaN()},
		{"1 / NaN", math.NaN()},
		{"NaN % 1", math.NaN()},
		{"1 % NaN", math.NaN()},
		{"NaN ** 1", math.NaN()},
		{"1 ** NaN", math.NaN()},

		// Infinity arithmetic - addition/subtraction
		{"Inf + 1", math.Inf(1)},
		{"1 + Inf", math.Inf(1)},
		{"Inf + Inf", math.Inf(1)},
		{"-Inf + (-Inf)", math.Inf(-1)},
		{"Inf + (-Inf)", math.NaN()}, // Indeterminate form
		{"Inf - Inf", math.NaN()},    // Indeterminate form
		{"Inf - (-Inf)", math.Inf(1)},
		{"-Inf - Inf", math.Inf(-1)},

		// Infinity arithmetic - multiplication
		{"Inf * 2", math.Inf(1)},
		{"2 * Inf", math.Inf(1)},
		{"Inf * (-2)", math.Inf(-1)},
		{"-Inf * 2", math.Inf(-1)},
		{"-Inf * (-2)", math.Inf(1)},
		{"Inf * Inf", math.Inf(1)},
		{"Inf * (-Inf)", math.Inf(-1)},
		{"0 * Inf", math.NaN()}, // 0 * infinity = NaN
		{"Inf * 0", math.NaN()},

		// Infinity arithmetic - division
		{"Inf / 2", math.Inf(1)},
		{"-Inf / 2", math.Inf(-1)},
		{"Inf / (-2)", math.Inf(-1)},
		{"2 / Inf", 0.0},
		{"2 / (-Inf)", 0.0}, // Note: -0.0 in IEEE-754, but Go treats as 0.0
		{"Inf / Inf", math.NaN()},
		{"Inf / (-Inf)", math.NaN()},
		{"-Inf / Inf", math.NaN()},

		// Power operations with infinity
		{"Inf ** 2", math.Inf(1)},
		{"Inf ** (-1)", 0.0},
		{"-Inf ** 2", math.Inf(1)},  // Even power of -Inf is +Inf
		{"-Inf ** 3", math.Inf(-1)}, // Odd power of -Inf is -Inf
		{"2 ** Inf", math.Inf(1)},
		{"0.5 ** Inf", 0.0},
		{"1 ** Inf", math.NaN()}, // Special case in math.Pow

		// Modulo operations with special values
		{"5 % Inf", 5.0},
		{"Inf % 5", math.NaN()},
		{"5 % 0", math.NaN()}, // Note: Go math.Mod returns NaN, no error in VM

		// Unary minus with special values
		{"-NaN", math.NaN()},
		{"-Inf", math.Inf(-1)},
		{"-(-Inf)", math.Inf(1)},
	}

	// Run tests with custom validation for NaN values
	runVmTestsWithIEEE754(t, tests)
}

// TestIEEE754DivisionByZero tests that division by zero remains an error
// even though IEEE-754 would normally produce infinity
func TestIEEE754DivisionByZero(t *testing.T) {
	tests := []vmTestCase{
		{"1 / 0", "division by zero"},
		{"5.5 / 0", "division by zero"},
		{"Inf / 0", "division by zero"},
		{"-Inf / 0", "division by zero"},
		{"NaN / 0", "division by zero"},
	}

	runVmErrorTestsWithIEEE754(t, tests)
}

// TestIEEE754Comparisons tests comparison operations with NaN and Inf values
func TestIEEE754Comparisons(t *testing.T) {
	tests := []vmTestCase{
		// NaN comparisons - all comparisons with NaN are false except !=
		{"NaN == NaN", false},
		{"NaN == 1", false},
		{"1 == NaN", false},
		{"NaN != NaN", true},
		{"NaN != 1", true},
		{"1 != NaN", true},
		{"NaN > 1", false},
		{"1 > NaN", false},
		{"NaN >= 1", false},
		{"1 >= NaN", false},

		// Infinity comparisons
		{"Inf == Inf", true},
		{"-Inf == -Inf", true},
		{"Inf == -Inf", false},
		{"Inf != -Inf", true},
		{"Inf > 1000", true},
		{"1000 > Inf", false},
		{"-Inf > -1000", false},
		{"-1000 > -Inf", true},
		{"Inf >= Inf", true},
		{"-Inf >= -Inf", true},
		{"Inf >= 1000", true},
		{"-Inf >= -1000", false},

		// Mixed infinity and finite comparisons
		{"Inf > -Inf", true},
		{"-Inf > Inf", false},
		{"1000000 > Inf", false},
		{"-1000000 > -Inf", true},
	}

	runVmTestsWithIEEE754(t, tests)
}

// TestIEEE754Truthiness tests truthiness of special values
func TestIEEE754Truthiness(t *testing.T) {
	tests := []vmTestCase{
		// NaN and Inf are truthy (non-zero)
		{"!NaN", false},    // NaN is truthy, so !NaN is false
		{"!Inf", false},    // Inf is truthy, so !Inf is false
		{"!(-Inf)", false}, // -Inf is truthy, so !(-Inf) is false
		{"!!NaN", true},    // Double negation converts to boolean true
		{"!!Inf", true},
		{"!!(-Inf)", true},
		{"!0", true}, // Only 0 is falsy among numbers
	}

	runVmTestsWithIEEE754(t, tests)
}

// TestIEEE754BitwiseErrors tests that bitwise operations with NaN/Inf produce errors
func TestIEEE754BitwiseErrors(t *testing.T) {
	tests := []vmTestCase{
		{"NaN & 1", "bitwise requires finite integers"},
		{"1 & NaN", "bitwise requires finite integers"},
		{"Inf | 1", "bitwise requires finite integers"},
		{"1 | Inf", "bitwise requires finite integers"},
		{"NaN ~ 1", "bitwise requires finite integers"}, // changed from ^
		{"1 ~ NaN", "bitwise requires finite integers"}, // changed from ^
		{"Inf << 1", "bitwise requires finite integers"},
		{"1 << Inf", "bitwise requires finite integers"},
		{"NaN >> 1", "bitwise requires finite integers"},
		{"1 >> NaN", "bitwise requires finite integers"},
	}

	runVmErrorTestsWithIEEE754(t, tests)
}

// Helper functions

// parseWithIEEE754 creates a parser with IEEE-754 specials enabled
func parseWithIEEE754(input string) parser.Node {
	opts := parser.DefaultOptions()
	opts.EnableIeeeSpecials = true
	p := parser.NewParserWithOptions(input, opts)
	node, err := p.Parse()
	if err != nil {
		return nil
	}
	return node
}

// runVmTestsWithIEEE754 runs VM tests with IEEE-754 support enabled
func runVmTestsWithIEEE754(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for i, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseWithIEEE754(tt.input)
			if program == nil {
				t.Fatalf("Failed to parse: %s", tt.input)
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("[case %d] compiler error for %s: %s", i+1, tt.input, err)
			}

			vm := vm.New(vm.LibContext{
				Functions:    vm.Builtins,
				PipeHandlers: vm.DefaultPipeHandlers,
			})
			bytecode := comp.ByteCode()
			result, err := vm.Run(bytecode, nil)
			if err != nil {
				t.Fatalf("[case %d] vm error for %s: %s", i+1, tt.input, err)
			}

			// Handle NaN comparison specially
			if expectedFloat, ok := tt.expected.(float64); ok && math.IsNaN(expectedFloat) {
				if resultFloat, ok := result.(float64); ok && math.IsNaN(resultFloat) {
					// Both are NaN, test passes
					return
				} else {
					t.Errorf("[case %d] expected NaN, got %v for input: %s", i+1, result, tt.input)
				}
			} else if expectedFloat, ok := tt.expected.(float64); ok && math.IsInf(expectedFloat, 0) {
				expectedSign := 1
				if math.IsInf(expectedFloat, -1) {
					expectedSign = -1
				}
				if resultFloat, ok := result.(float64); ok && math.IsInf(resultFloat, expectedSign) {
					// Both are the same infinity, test passes
					return
				} else {
					t.Errorf("[case %d] expected Inf(%d), got %v for input: %s", i+1, expectedSign, result, tt.input)
				}
			} else {
				err := testExpectedObject(t, tt.expected, result)
				if err != nil {
					t.Errorf("[case %d] input: %s, error: %s", i+1, tt.input, err)
				}
			}
		})
	}
}

// runVmErrorTestsWithIEEE754 runs VM error tests with IEEE-754 support enabled
func runVmErrorTestsWithIEEE754(t *testing.T, tests []vmTestCase) {
	t.Helper()
	for i, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseWithIEEE754(tt.input)
			if program == nil {
				t.Fatalf("Failed to parse: %s", tt.input)
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("[case %d] compiler error for %s: %s", i+1, tt.input, err)
			}

			vm := vm.New(vm.LibContext{
				Functions:    vm.Builtins,
				PipeHandlers: vm.DefaultPipeHandlers,
			})
			bytecode := comp.ByteCode()
			_, err = vm.Run(bytecode, nil)

			if err == nil {
				t.Errorf("[case %d] expected error for: %s", i+1, tt.input)
				return
			}

			expectedErr, ok := tt.expected.(string)
			if !ok {
				t.Fatalf("[case %d] expected error message to be a string, got %T", i+1, tt.expected)
			}

			if !strings.Contains(err.Error(), expectedErr) {
				t.Errorf("[case %d] expected error containing '%s', got '%s' for input: %s",
					i+1, expectedErr, err.Error(), tt.input)
			}
		})
	}
}
