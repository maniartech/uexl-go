package parser_test

import (
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
		{
			name:  "dot expression with function call",
			input: "a.b.method(c.d, e.f)",
		},
		{
			name:  "function return with dot access",
			input: "func(a, b).property",
		},
		{
			name:  "array element with dot access",
			input: "[1, 2, 3][0].property",
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
