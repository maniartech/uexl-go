package compiler_test

import (
	"math"
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

// TestIEEE754CompilerConstantFolding tests that the compiler correctly handles
// constant folding with NaN and Inf values
func TestIEEE754CompilerConstantFolding(t *testing.T) {
	tests := []struct {
		input    string
		expected string // Expected bytecode pattern or constant
	}{
		// Simple constants should be folded
		{"NaN", "LOAD_CONST NaN"},
		{"Inf", "LOAD_CONST +Inf"},
		{"-Inf", "LOAD_CONST -Inf"},

		// Arithmetic with constants should be folded
		{"NaN + 1", "LOAD_CONST NaN"}, // NaN + anything = NaN
		{"1 + NaN", "LOAD_CONST NaN"},
		{"Inf + 1", "LOAD_CONST +Inf"}, // Inf + finite = Inf
		{"Inf * 2", "LOAD_CONST +Inf"},
		{"0 * Inf", "LOAD_CONST NaN"}, // 0 * Inf = NaN

		// Operations that should NOT be folded (require runtime)
		{"x + NaN", "LOAD_VAR x; LOAD_CONST NaN; ADD"}, // Variable involved
		{"NaN + y", "LOAD_CONST NaN; LOAD_VAR y; ADD"},
	}

	for i, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseWithIEEE754(tt.input)
			if program == nil {
				t.Fatalf("[case %d] Failed to parse: %s", i+1, tt.input)
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("[case %d] Compiler error for %s: %s", i+1, tt.input, err)
			}

			// For now, just ensure compilation succeeds
			// When constant folding is implemented, we can check bytecode patterns
			bytecode := comp.ByteCode()
			if len(bytecode.Constants) == 0 && containsIEEE754Literal(tt.input) {
				t.Errorf("[case %d] Expected constants to contain IEEE-754 values for: %s", i+1, tt.input)
			}
		})
	}
}

// TestIEEE754CompilerConstants tests that IEEE-754 constants are properly
// stored in the bytecode constant pool
func TestIEEE754CompilerConstants(t *testing.T) {
	tests := []struct {
		input     string
		hasNaN    bool
		hasInf    bool
		hasNegInf bool
	}{
		{"NaN", true, false, false},
		{"Inf", false, true, false},
		{"-Inf", false, false, true},
		{"NaN + Inf", true, true, false},
		{"Inf - (-Inf)", false, true, true},
		{"[NaN, Inf, -Inf]", true, true, true},
		{`{"nan": NaN, "inf": Inf, "neginf": -Inf}`, true, true, true},
	}

	for i, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseWithIEEE754(tt.input)
			if program == nil {
				t.Fatalf("[case %d] Failed to parse: %s", i+1, tt.input)
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("[case %d] Compiler error for %s: %s", i+1, tt.input, err)
			}

			bytecode := comp.ByteCode()
			constants := bytecode.Constants

			foundNaN := false
			foundInf := false
			foundNegInf := false

			for _, constant := range constants {
				if f, ok := constant.(float64); ok {
					if math.IsNaN(f) {
						foundNaN = true
					} else if math.IsInf(f, 1) {
						foundInf = true
					} else if math.IsInf(f, -1) {
						foundNegInf = true
					}
				}
			}

			if tt.hasNaN && !foundNaN {
				t.Errorf("[case %d] Expected NaN in constants for: %s", i+1, tt.input)
			}
			if tt.hasInf && !foundInf {
				t.Errorf("[case %d] Expected +Inf in constants for: %s", i+1, tt.input)
			}
			if tt.hasNegInf && !foundNegInf {
				t.Errorf("[case %d] Expected -Inf in constants for: %s", i+1, tt.input)
			}
		})
	}
}

// TestIEEE754CompilerComplexExpressions tests compilation of complex expressions
// involving IEEE-754 special values
func TestIEEE754CompilerComplexExpressions(t *testing.T) {
	tests := []string{
		// Conditional expressions
		"true ? NaN : Inf",
		"false ? Inf : -Inf",
		"x > 0 ? Inf : NaN",

		// Function calls with special values
		"abs(NaN)",
		"min(Inf, -Inf)",
		"max(NaN, 5)",

		// Array/object access with special values
		"arr[NaN]", // Should compile but likely error at runtime
		"obj[Inf]", // Should compile but likely error at runtime

		// Pipe operations with special values
		"NaN | round",
		"Inf | abs",
		"[NaN, Inf, -Inf] | map(x => x + 1)",

		// Complex arithmetic chains
		"(NaN + 1) * (Inf - 2) / (-Inf + 3)",
		"Inf ** 2 + NaN * 0 - (-Inf)",

		// Nested conditionals with special values
		"NaN == NaN ? 'equal' : (Inf > -Inf ? 'greater' : 'other')",
	}

	for i, input := range tests {
		t.Run(input, func(t *testing.T) {
			program := parseWithIEEE754(input)
			if program == nil {
				t.Fatalf("[case %d] Failed to parse: %s", i+1, input)
			}

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("[case %d] Compiler error for %s: %s", i+1, input, err)
			}

			// Just verify compilation succeeds
			bytecode := comp.ByteCode()
			if len(bytecode.Instructions) == 0 {
				t.Errorf("[case %d] Expected non-empty instructions for: %s", i+1, input)
			}
		})
	}
}

// TestIEEE754CompilerErrorHandling tests that the compiler properly handles
// expressions that should produce compile-time errors
func TestIEEE754CompilerErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
		errorMsg    string
	}{
		// These should compile fine (runtime errors)
		{"NaN & 1", false, ""},  // Bitwise with NaN - runtime error
		{"Inf << 2", false, ""}, // Shift with Inf - runtime error

		// These might be compile-time errors depending on implementation
		// For now, assume they compile and error at runtime
		{"NaN && true", false, ""},  // Type error - could be compile-time
		{"Inf || false", false, ""}, // Type error - could be compile-time
	}

	for i, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseWithIEEE754(tt.input)
			if program == nil {
				if tt.shouldError {
					return // Parse error expected
				}
				t.Fatalf("[case %d] Failed to parse: %s", i+1, tt.input)
			}

			comp := compiler.New()
			err := comp.Compile(program)

			if tt.shouldError {
				if err == nil {
					t.Errorf("[case %d] Expected compiler error for: %s", i+1, tt.input)
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("[case %d] Expected error '%s', got '%s' for: %s",
						i+1, tt.errorMsg, err.Error(), tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("[case %d] Unexpected compiler error for %s: %s", i+1, tt.input, err)
				}
			}
		})
	}
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

// containsIEEE754Literal checks if input contains NaN, Inf, or -Inf literals
func containsIEEE754Literal(input string) bool {
	return containsSubstring(input, "NaN") ||
		containsSubstring(input, "Inf") ||
		containsSubstring(input, "-Inf")
}

// containsSubstring is a simple substring check
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
