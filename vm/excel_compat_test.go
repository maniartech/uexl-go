package vm_test

import (
	"testing"

	"github.com/maniartech/uexl/compiler"
	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/vm"
)

// TestExcelCompatibility verifies all Excel-compatible operator changes
// This suite tests the breaking changes and new aliases introduced for Excel compatibility

func TestExcelCompat_PowerOperator_Caret(t *testing.T) {
	// ^ is Excel-compatible power operator (alternative to **)
	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"Basic power", "2 ^ 3", 8.0},
		{"Power with negatives", "-2 ^ 2", 4.0},
		{"Power of zero", "5 ^ 0", 1.0},
		{"Power with decimals", "2 ^ 2.5", 5.65685424949238},
		{"Chained power (right associative)", "2 ^ 3 ^ 2", 512.0}, // 2^(3^2) = 2^9
		{"Power in expression", "10 + 2 ^ 3", 18.0},
		{"Power with parentheses", "(2 + 3) ^ 2", 25.0},
		{"Fractional exponent", "4 ^ 0.5", 2.0},
		{"Large exponent", "2 ^ 10", 1024.0},
		{"Negative exponent", "2 ^ -3", 0.125},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v ^ to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_PowerOperator_DoubleStar(t *testing.T) {
	// ** is an active alternative to ^ (Python/JavaScript style)
	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"Basic power", "2 ** 3", 8.0},
		{"Power with negatives", "-2 ** 2", 4.0},
		{"Power of zero", "5 ** 0", 1.0},
		{"Chained power", "2 ** 3 ** 2", 512.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v ** to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_BitwiseXOR_Tilde(t *testing.T) {
	// ~ for XOR (Lua-style, was previously ^ which now means power)
	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"XOR same values", "5 ~ 5", 0.0},
		{"XOR different values", "5 ~ 3", 6.0},
		{"XOR with zero", "7 ~ 0", 7.0},
		{"XOR large values", "255 ~ 128", 127.0},
		{"XOR negative", "-1 ~ 1", -2.0},
		{"Chained XOR", "5 ~ 3 ~ 1", 7.0},
		{"XOR in expression", "10 + (5 ~ 3)", 16.0},
		{"XOR with parentheses", "(8 ~ 4) ~ 2", 14.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_BitwiseNOT_Tilde(t *testing.T) {
	// Unary ~ for bitwise NOT (Lua-style)
	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"NOT zero", "~0", -1.0},
		{"NOT one", "~1", -2.0},
		{"NOT five", "~5", -6.0},
		{"NOT negative", "~(-1)", 0.0},
		{"Double NOT", "~~5", 5.0},
		{"NOT in expression", "~5 + 10", 4.0},
		{"NOT with parentheses", "~(5 + 3)", -9.0},
		{"NOT of XOR", "~(5 ~ 3)", -7.0},
		{"NOT chain", "~~~5", -6.0},
		{"NOT large value", "~255", -256.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_NotEquals_Diamond(t *testing.T) {
	// <> is Excel-compatible not-equals operator (alternative to !=)
	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"Numbers not equal", "5 <> 3", true},
		{"Numbers equal", "5 <> 5", false},
		{"Strings not equal", `"hello" <> "world"`, true},
		{"Strings equal", `"hello" <> "hello"`, false},
		{"Bools not equal", "true <> false", true},
		{"Bools equal", "true <> true", false},
		{"Zero comparison", "0 <> 0", false},
		{"Negative numbers", "-5 <> -3", true},
		{"Mixed types (coerced)", "5 <> 5.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExprBool(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_NotEquals_BangEquals(t *testing.T) {
	// != is an active alternative to <> (C/Python/JavaScript style)
	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"Numbers not equal", "5 != 3", true},
		{"Numbers equal", "5 != 5", false},
		{"Strings not equal", `"hello" != "world"`, true},
		{"Strings equal", `"hello" != "hello"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExprBool(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

func TestExcelCompat_MixedOperators(t *testing.T) {
	// Complex expressions using multiple new operators
	tests := []struct {
		name     string
		expr     string
		expected interface{}
	}{
		{"Power and XOR", "2 ^ 3 ~ 4", 12.0},            // 8 ~ 4 = 12
		{"NOT and power", "~2 ^ 2", 9.0},                // (-3) ^ 2 = 9
		{"Power and not-equals", "(2 ^ 3) <> 8", false}, // 8 <> 8 = false
		{"XOR and not-equals", "(5 ~ 3) <> 6", false},   // 6 <> 6 = false
		{"All operators", "(~(2 ^ 3) ~ 4) <> -5", true}, // ~8 ~ 4 = -9 ~ 4 = -13, -13 <> -5 = true
		{"Nested power", "(2 ^ 2) ^ 2", 16.0},           // 4 ^ 2 = 16
		{"NOT of NOT", "~~(5 ~ 3)", 6.0},                // ~~6 = 6
		{"Complex arithmetic", "10 + 2 ^ 3 - ~5", 24.0}, // 10 + 8 - (-6) = 24
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch expected := tt.expected.(type) {
			case float64:
				result := evalExpr(t, tt.expr)
				if result != expected {
					t.Errorf("Expected %v to equal %v, got %v", tt.expr, expected, result)
				}
			case bool:
				result := evalExprBool(t, tt.expr)
				if result != expected {
					t.Errorf("Expected %v to equal %v, got %v", tt.expr, expected, result)
				}
			}
		})
	}
}

func TestExcelCompat_ErrorCases(t *testing.T) {
	// Verify proper error handling with new operators
	tests := []struct {
		name        string
		expr        string
		shouldError bool
	}{
		{"Bitwise NOT with decimal", "~5.5", true},
		{"Bitwise XOR with decimal left", "5.5 ~ 3", true},
		{"Bitwise XOR with decimal right", "5 ~ 3.5", true},
		{"Power is OK with decimals", "2.5 ^ 2", false},
		{"NOT is OK with whole decimals", "~5.0", false},
		{"XOR is OK with whole decimals", "5.0 ~ 3.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parser.ParseString(tt.expr)
			if err != nil {
				if !tt.shouldError {
					t.Errorf("Parser error for %v: %v", tt.expr, err)
				}
				return
			}

			comp := compiler.New()
			comp.Compile(node)
			bytecode := comp.ByteCode()

			machine := vm.New(vm.LibContext{
				Functions:    vm.Builtins,
				PipeHandlers: vm.DefaultPipeHandlers,
			})

			_, err = machine.Run(bytecode, nil)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for %v, got nil", tt.expr)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for %v: %v", tt.expr, err)
			}
		})
	}
}

func TestExcelCompat_PrecedenceCorrectness(t *testing.T) {
	// Verify operator precedence is correct after changes

	// Arithmetic precedence tests
	arithmeticTests := []struct {
		name     string
		expr     string
		expected float64
	}{
		// Power has higher precedence than multiplication
		{"Power before multiply", "2 * 3 ^ 2", 18.0}, // 2 * 9 = 18

		// Bitwise XOR has correct precedence (between AND and OR)
		{"XOR with AND", "5 & 3 ~ 1", 0.0}, // (5 & 3) ~ 1 = 1 ~ 1 = 0

		// Unary NOT before binary operations
		{"NOT before XOR", "~5 ~ 3", -7.0}, // (-6) ~ 3 = -7

		// Power is right-associative
		{"Power right associative", "2 ^ 3 ^ 2", 512.0}, // 2 ^ (3 ^ 2) = 2 ^ 9 = 512
	}

	for _, tt := range arithmeticTests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExpr(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}

	// Comparison precedence tests
	comparisonTests := []struct {
		name     string
		expr     string
		expected bool
	}{
		// Comparison has lower precedence than arithmetic
		{"Comparison after arithmetic", "2 + 3 <> 5", false}, // (2 + 3) <> 5 = 5 <> 5 = false
		{"Comparison after power", "2 ^ 3 <> 8", false},      // (2 ^ 3) <> 8 = 8 <> 8 = false
	}

	for _, tt := range comparisonTests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalExprBool(t, tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %v to equal %v, got %v", tt.expr, tt.expected, result)
			}
		})
	}
}

// Helper functions

func evalExpr(t *testing.T, expr string) float64 {
	t.Helper()

	node, err := parser.ParseString(expr)
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	result, err := machine.Run(bytecode, nil)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	floatResult, ok := result.(float64)
	if !ok {
		t.Fatalf("Expected float64 result, got %T: %v", result, result)
	}

	return floatResult
}

func evalExprBool(t *testing.T, expr string) bool {
	t.Helper()

	node, err := parser.ParseString(expr)
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	result, err := machine.Run(bytecode, nil)
	if err != nil {
		t.Fatalf("Runtime error: %v", err)
	}

	boolResult, ok := result.(bool)
	if !ok {
		t.Fatalf("Expected bool result, got %T: %v", result, result)
	}

	return boolResult
}
