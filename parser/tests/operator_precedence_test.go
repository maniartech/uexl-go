package parser_test

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

// TestOperatorPrecedence tests that the parser correctly follows
// the operator precedence rules previously defined in grammar.peg (now replaced by custom tokenizer-based parser)
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "multiplication has higher precedence than addition",
			input: "1 + 2 * 3",
			// Should be equivalent to: 1 + (2 * 3)
		},
		{
			name:  "division has higher precedence than subtraction",
			input: "10 - 6 / 2",
			// Should be equivalent to: 10 - (6 / 2)
		},
		{
			name:  "modulo has higher precedence than multiplication",
			input: "5 * 3 % 2",
			// Should be equivalent to: 5 * (3 % 2)
		},
		{
			name:  "dot operator has higher precedence than modulo",
			input: "a.b % 2",
			// Should be equivalent to: (a.b) % 2
		},
		{
			name:  "comparison has higher precedence than equality",
			input: "a == b > c",
			// Should be equivalent to: a == (b > c)
		},
		{
			name:  "equality has higher precedence than bitwise AND",
			input: "a & b == c",
			// Should be equivalent to: a & (b == c)
		},
		{
			name:  "bitwise AND has higher precedence than bitwise XOR",
			input: "a ^ b & c",
			// Should be equivalent to: a ^ (b & c)
		},
		{
			name:  "bitwise XOR has higher precedence than bitwise OR",
			input: "a | b ^ c",
			// Should be equivalent to: a | (b ^ c)
		},
		{
			name:  "bitwise OR has higher precedence than logical AND",
			input: "a && b | c",
			// Should be equivalent to: a && (b | c)
		},
		{
			name:  "logical AND has higher precedence than logical OR",
			input: "a || b && c",
			// Should be equivalent to: a || (b && c)
		},
		{
			name:  "nullish coalescing same precedence as logical OR",
			input: "a ?? b || c && d",
			// Should be equivalent grouping: (a ?? b) || (c && d)
		},
		{
			name:  "conditional lower than || and ??",
			input: "a ?? b ? c : d",
			// Should parse as (a ?? b) ? c : d
		},
		{
			name:  "nested conditional right associative",
			input: "a ? b : c ? d : e",
			// Should parse as a ? b : (c ? d : e)
		},
		{
			name:  "shift operators have correct precedence",
			input: "a + b << c * d",
			// Should be equivalent to: (a + b) << (c * d)
		},
		{
			name:  "complex expression with multiple precedence levels",
			input: "a || b && c | d ^ e & f == g > h + i * j % k.l",
			// Should parse according to precedence rules
		},
		{
			name:  "parentheses override precedence",
			input: "(a + b) * c",
			// Should be equivalent to: (a + b) * c, not a + (b * c)
		},
		{
			name:  "nested parentheses",
			input: "a * (b + (c - d))",
			// Should respect all parentheses
		},
		{
			name:  "member access chain",
			input: "a.b.c.d",
			// Should be parsed as ((a.b).c).d
		},
		{
			name:  "member access with operations",
			input: "a.b + c.d * e.f",
			// Should be equivalent to: (a.b) + ((c.d) * (e.f))
		},
		{
			name:  "pipe with binary expressions",
			input: "a + b |map: c * d",
			// Should be equivalent to: (a + b) |map: (c * d)
		},
		{
			name:  "multiple pipes",
			input: "a |map: b + c |filter: d > e",
			// Should be equivalent to: (a |map: (b + c)) |filter: (d > e)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			ast, err := p.Parse()

			assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)

			// For visual inspection during test development
			// Uncomment to see the AST structure
			// json.PrintJSON(ast)

			// This test mainly verifies that parsing succeeds
			// More specific tests for the structure would be added based on inspecting the output
			assert.NotNil(t, ast, "AST should not be nil")
		})
	}
}

// TestDotExpressionPrecedence specifically tests dot expressions
func TestDotExpressionPrecedence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple dot expression",
			input: "a.b",
		},
		{
			name:  "chained dot expressions",
			input: "a.b.c.d",
		},
		{
			name:  "dot expressions with operations",
			input: "a.b.c + d.e.f",
		},
		// Removed invalid test: property access after function call is valid in UExL
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()

			if tt.name == "dot expression with function call (should error)" {
				assert.Error(t, err, "Parsing should produce an error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)
			}
		})
	}
}

// TestInvalidDotExpressions tests invalid dot expressions that should be rejected
func TestInvalidDotExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Removed: property access after function call is valid in UExL
		// Note: Array indexing followed by member access is valid
		// [1, 2, 3][0].property is a valid pattern
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()

			if tt.input == "func(a, b).property" {
				if err == nil {
					t.Errorf("Expected error for invalid pattern: %s", tt.input)
				} else {
					errStr := err.Error()
					if !strings.Contains(errStr, "function calls are only allowed after identifiers or function calls") {
						t.Errorf("Expected error about function call chaining, got: %s", errStr)
					}
				}
			} else {
				assert.Error(t, err, "Parsing should produce an error for invalid pattern: %s", tt.input)
			}
		})
	}
}

// TestPipeExpressionPrecedence specifically tests pipe expressions
func TestPipeExpressionPrecedence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple pipe",
			input: "a |: b",
		},
		{
			name:  "pipe with type",
			input: "a |map: b",
		},
		{
			name:  "multiple pipes",
			input: "a |map: b |filter: c |reduce: d",
		},
		{
			name:  "pipes with complex expressions",
			input: "a + b * c |map: d.e + f |filter: g && h",
		},
		{
			name:  "parenthesized pipe source",
			input: "(a + b) |map: c",
		},
		{
			name:  "pipe with function calls",
			input: "list |map: transform(item) |filter: isValid(item)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()

			assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)
		})
	}
}

// TestConsecutiveOperators tests parsing of consecutive unary operators like -- and !!
func TestConsecutiveOperators(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "double_negation_number",
			input:       "--10",
			expectError: false,
			description: "Double negation with number should parse as -(-(10))",
		},
		{
			name:        "double_negation_variable",
			input:       "--x",
			expectError: false,
			description: "Double negation with variable should parse as -(-(x))",
		},
		{
			name:        "triple_negation",
			input:       "---5",
			expectError: false,
			description: "Triple negation should parse as -(-(-(5)))",
		},
		{
			name:        "double_not_true",
			input:       "!!true",
			expectError: false,
			description: "Double NOT with boolean should parse as !(!(true))",
		},
		{
			name:        "double_not_false",
			input:       "!!false",
			expectError: false,
			description: "Double NOT with boolean should parse as !(!(false))",
		},
		{
			name:        "triple_not",
			input:       "!!!true",
			expectError: false,
			description: "Triple NOT should parse as !(!(!(true)))",
		},
		{
			name:        "double_not_variable",
			input:       "!!x",
			expectError: false,
			description: "Double NOT with variable should parse as !(!(x))",
		},
		{
			name:        "mixed_consecutive_operators",
			input:       "!-5",
			expectError: false,
			description: "Mixed unary operators should parse as !(-(5))",
		},
		{
			name:        "mixed_consecutive_operators_reverse",
			input:       "-!true",
			expectError: false,
			description: "Mixed unary operators should parse as -(!(true))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)
				assert.NotNil(t, expr, "Expression should not be nil for input: %s", tt.input)

				// Verify that we get a UnaryExpression for consecutive operators
				unaryExpr, ok := expr.(*parser.UnaryExpression)
				assert.True(t, ok, "Expected UnaryExpression for input: %s", tt.input)
				assert.NotNil(t, unaryExpr.Operand, "Operand should not be nil for input: %s", tt.input)
			}
		})
	}
}

// TestConsecutiveOperatorsParsing tests the AST structure of consecutive operators
func TestConsecutiveOperatorsParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOp     string
		expectedNested int // How many nested unary expressions we expect
	}{
		{
			name:           "double_minus",
			input:          "--10",
			expectedOp:     "-",
			expectedNested: 2,
		},
		{
			name:           "triple_minus",
			input:          "---10",
			expectedOp:     "-",
			expectedNested: 3,
		},
		{
			name:           "double_not",
			input:          "!!true",
			expectedOp:     "!",
			expectedNested: 2,
		},
		{
			name:           "triple_not",
			input:          "!!!true",
			expectedOp:     "!",
			expectedNested: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)

			// Check nesting depth
			current := expr
			depth := 0
			for {
				unaryExpr, ok := current.(*parser.UnaryExpression)
				if !ok {
					break
				}
				depth++
				assert.Equal(t, tt.expectedOp, unaryExpr.Operator, "Expected operator %s at depth %d", tt.expectedOp, depth)
				current = unaryExpr.Operand
			}

			assert.Equal(t, tt.expectedNested, depth, "Expected nesting depth %d for input %s", tt.expectedNested, tt.input)
		})
	}
}

// TestPowerOperator tests the ** power operator
func TestPowerOperator(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "simple_power",
			input:       "2**3",
			expectError: false,
			description: "Simple power operation: 2**3",
		},
		{
			name:        "power_with_spaces",
			input:       "2 ** 3",
			expectError: false,
			description: "Power operation with spaces: 2 ** 3",
		},
		{
			name:        "right_associative_power",
			input:       "2**3**2",
			expectError: false,
			description: "Right-associative power: 2**(3**2) = 2**9",
		},
		{
			name:        "power_with_parentheses",
			input:       "(2**3)**2",
			expectError: false,
			description: "Power with parentheses: (2**3)**2 = 8**2",
		},
		{
			name:        "power_with_negative",
			input:       "(-2)**3",
			expectError: false,
			description: "Power with negative base: (-2)**3",
		},
		{
			name:        "power_with_decimal",
			input:       "2.5**2",
			expectError: false,
			description: "Power with decimal: 2.5**2",
		},
		{
			name:        "zero_power",
			input:       "5**0",
			expectError: false,
			description: "Zero exponent: 5**0",
		},
		{
			name:        "power_vs_multiplication",
			input:       "2*3**2",
			expectError: false,
			description: "Power has higher precedence than multiplication: 2*(3**2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)
				assert.NotNil(t, expr, "Expression should not be nil for input: %s", tt.input)
			}
		})
	}
}

// TestPowerOperatorPrecedence tests that power operator has correct precedence
func TestPowerOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
		description  string
	}{
		{
			name:         "power_vs_multiplication",
			input:        "2*3**2",
			expectedType: "BinaryExpression",
			description:  "Should be parsed as 2*(3**2), not (2*3)**2",
		},
		{
			name:         "power_vs_addition",
			input:        "1+2**3",
			expectedType: "BinaryExpression",
			description:  "Should be parsed as 1+(2**3), not (1+2)**3",
		},
		{
			name:         "right_associative",
			input:        "2**3**4",
			expectedType: "BinaryExpression",
			description:  "Should be parsed as 2**(3**4), not (2**3)**4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err)

			binExpr, ok := expr.(*parser.BinaryExpression)
			assert.True(t, ok, "Expected BinaryExpression for input: %s", tt.input)

			// Verify the structure based on the specific test
			switch tt.name {
			case "power_vs_multiplication":
				// Should be: * operator at root, with left=2 and right=(3**2)
				assert.Equal(t, "*", binExpr.Operator)
				rightExpr, ok := binExpr.Right.(*parser.BinaryExpression)
				assert.True(t, ok, "Right side should be a BinaryExpression")
				assert.Equal(t, "**", rightExpr.Operator)

			case "power_vs_addition":
				// Should be: + operator at root, with left=1 and right=(2**3)
				assert.Equal(t, "+", binExpr.Operator)
				rightExpr, ok := binExpr.Right.(*parser.BinaryExpression)
				assert.True(t, ok, "Right side should be a BinaryExpression")
				assert.Equal(t, "**", rightExpr.Operator)

			case "right_associative":
				// Should be: ** operator at root, with left=2 and right=(3**4)
				assert.Equal(t, "**", binExpr.Operator)
				rightExpr, ok := binExpr.Right.(*parser.BinaryExpression)
				assert.True(t, ok, "Right side should be a BinaryExpression")
				assert.Equal(t, "**", rightExpr.Operator)
			}
		})
	}
}
