package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

// TestSimpleBinaryExpression tests a simple binary expression
func TestSimpleBinaryExpression(t *testing.T) {
	input := "1 + 2"
	expected := &parser.BinaryExpression{
		Left:     &parser.NumberLiteral{Value: "1", Line: 1, Column: 1},
		Operator: "+",
		Right:    &parser.NumberLiteral{Value: "2", Line: 1, Column: 5},
		Line:     1,
		Column:   3,
	}
	testParser(t, input, expected)
}

// TestNestedBinaryExpression tests a nested binary expression
func TestNestedBinaryExpression(t *testing.T) {
	input := "1 + 2 * 3"
	expected := &parser.BinaryExpression{
		Left:     &parser.NumberLiteral{Value: "1", Line: 1, Column: 1},
		Operator: "+",
		Right: &parser.BinaryExpression{
			Left:     &parser.NumberLiteral{Value: "2", Line: 1, Column: 5},
			Operator: "*",
			Right:    &parser.NumberLiteral{Value: "3", Line: 1, Column: 9},
			Line:     1,
			Column:   7,
		},
		Line:   1,
		Column: 3,
	}
	testParser(t, input, expected)
}

// TestPipeExpression tests a pipe expression

func TestPipeExpressionMapOperation(t *testing.T) {
	input := "[1, 2, 3] |map: $1 * 2 |: $1 + 1"
	p := parser.NewParser(input)
	ast, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error for input: %s", input)

	fmt.Printf("%+v", ast)
}

func TestPipeExpressionMultipleFunctions(t *testing.T) {
	input := "x + y |: func($1) |: otherFunc($1, 2)"
	p := parser.NewParser(input)
	_, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error for input: %s", input)

}

func TestPipeExpressionMethodChaining(t *testing.T) {
	input := "method(1, 2) |: $1.property |: upperCase"
	p := parser.NewParser(input)
	_, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error for input: %s", input)
}

func TestPipeExpressionMultipleOperations(t *testing.T) {
	input := "[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2"
	p := parser.NewParser(input)
	_, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error for input: %s", input)
}

func TestPipeExpressionStringOperations(t *testing.T) {
	input := "'hello' |: splitChars |filter: $1 != 'l' |join: ''"
	p := parser.NewParser(input)
	_, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error for input: %s", input)
}

// TestObjectLiteral tests an object literal
func TestObjectLiteral(t *testing.T) {
	input := `{"a": 1, "b": true}`
	expected := &parser.ObjectLiteral{
		Properties: map[string]parser.Expression{
			"a": &parser.NumberLiteral{Value: "1", Line: 1, Column: 7},
			"b": &parser.BooleanLiteral{Value: true, Line: 1, Column: 15},
		},
		Line:   1,
		Column: 1,
	}
	testParser(t, input, expected)
}

// TestArrayLiteralWithPipe tests an array literal with a pipe expression
func TestArrayLiteralWithPipe(t *testing.T) {
	input := "[1, 2, 3] |map: (1 + 2)"
	expected := &parser.PipeExpression{
		Left: &parser.ArrayLiteral{
			Elements: []parser.Expression{
				&parser.NumberLiteral{Value: "1", Line: 1, Column: 2},
				&parser.NumberLiteral{Value: "2", Line: 1, Column: 5},
				&parser.NumberLiteral{Value: "3", Line: 1, Column: 8},
			},
			Line:   1,
			Column: 1,
		},
		PipeType: "map",
		Right: &parser.BinaryExpression{
			Left:     &parser.NumberLiteral{Value: "1", Line: 1, Column: 17},
			Operator: "+",
			Right:    &parser.NumberLiteral{Value: "2", Line: 1, Column: 21},
			Line:     1,
			Column:   19,
		},
		Line:   1,
		Column: 11,
	}
	testParser(t, input, expected)
}

// testParser is a helper function to run parser tests
func testParser(t *testing.T, input string, expected parser.Expression) {
	t.Helper()
	p := parser.NewParser(input)
	actual, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error")
	assert.Equal(t, expected, actual, "For input %q, AST should match", input)
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{
			input:       "3 + * 4",
			expectedErr: "unexpected token",
		},
		{
			input:       "a.b |: (1, 2)",
			expectedErr: "unexpected token",
		},
		{
			input:       "(1 + 2",
			expectedErr: "expected ')'",
		},
		{
			input:       "[1, 2,]",
			expectedErr: "expected ']'",
		},
		{
			input:       `{"a": 1,,}`,
			expectedErr: "expected '}'",
		},
		{
			input:       "1 + + 2",
			expectedErr: "unexpected token",
		},
		{
			input:       "a.",
			expectedErr: "expected identifier after '.'",
		},
		{
			input:       `{"a": }`,
			expectedErr: "unexpected token",
		},
		{
			input:       "func(1,)",
			expectedErr: "unexpected token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("Expected error, but got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing %q, but got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"(a + b) * (c - d) / e"},
		{"x(1, 2, 3) |: func(a, b)"},
		{"[1, 2, 3] * 2"},
		{`{"a": 1, "b": [true, false], "c": {"d": null}}`},
		{"a && b || c && d || e && f"},
		{"x << 2 + y >> 3 - z & a | b"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			// If we reach here, the parser successfully parsed the complex expression
		})
	}
}

func TestPipeExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"[1, 2, 3] |map: $1 * 2"},
		{"x + y |: func($1) |: otherFunc($1, 2)"},
		{"methodA(1, 2) |: $1.property |: upperCase"},
		{"[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2"},
		{"'hello' |: splitChars |filter: $1 != 'l' |join: ''"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			// If we reach here, the parser successfully parsed the pipe expression
		})
	}
}
