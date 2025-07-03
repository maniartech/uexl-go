package parser_test

import (
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
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
func TestPipeExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"[1, 2, 3] |map: $1 * 2"},
		{"x + y |: func($1) |: otherFunc($1, 2)"},
		{"methodA(1, 2) |: $1.property |: upperCase"},
		{"[1, 2, 3, 4, 5] |filter: $1 > 2 |map: $1 * 2 |reduce: $1 + $2"},
		{"'hello' |: splitChars |filter: $1 != 'l' |join: ''"},
		{"[$1.x.y, 2] |map: $1 * 2"}, // From previous TestGroupWithPipe
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			ast, err := p.Parse()
			assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)

			// Verify it's a pipe expression
			_, ok := ast.(*parser.PipeExpression)
			assert.True(t, ok, "Expected a PipeExpression for input: %s", tt.input)
		})
	}
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
			expectedErr: "expected ')'",
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
			input:       "1 + + 2",
			expectedErr: "unexpected token",
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

func TestMemberAccess(t *testing.T) {
	input := "$1.x.y * 2"
	p := parser.NewParser(input)
	ast, err := p.Parse()
	assert.NoError(t, err, "Unexpected error: %v", err)

	_, ok := ast.(*parser.BinaryExpression)
	assert.True(t, ok, "Expected a BinaryExpression")
}

func TestPipeExpressionWithAlias(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"x + 10 as $a |: y + 20 as $b |: $a + $b"},
		{"[1, 2, 3] |map: $1 * 2 as $doubled |filter: $doubled > 4"},
		{"obj.value as $v |: transform($v) as $t |: format($t)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			ast, err := p.Parse()
			assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)

			// Verify it's a pipe expression if it contains a pipe or alias
			if strings.Contains(tt.input, "|") || strings.Contains(tt.input, " as ") {
				pipeExpr, ok := ast.(*parser.PipeExpression)
				assert.True(t, ok, "Expected a PipeExpression for input: %s", tt.input)
				assert.NotNil(t, pipeExpr.Aliases, "Aliases slice should not be nil")
			}
		})
	}
}

// Add negative test cases for 'as' keyword misuse
func TestPipeExpressionWithInvalidAlias(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"x + 10 as y", "expected identifier starting with $"},
		{"x as $a as $b", "unexpected token"},
		{"[1, x as $a, 2]", "pipe expressions cannot be sub-expressions"},
		{"x as $a y as $b", "unexpected token"},
		{"(x + 1 as $a)", "pipe expressions cannot be sub-expressions"},
		{"func(1 as $a)", "pipe expressions cannot be sub-expressions"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()

			assert.Error(t, err, "Expected error for input: %s", tt.input)
			assert.Contains(t, err.Error(), tt.expectedErr,
				"Expected error containing %q, but got %q", tt.expectedErr, err.Error())
		})
	}
}

func TestPipeAliasInFunctionArgs(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"x as $a |: func($a) |: transform($a)"},
		{"obj.value as $v |: transform($v, 42) |: format($v)"},
		{"[1, 2, 3] |map: $1 * 2 as $doubled |filter: check($doubled)"},
		{"x as $a |: func($a + y as $b)"}, // This should fail since aliases can't contain other aliases
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			ast, err := p.Parse()

			if strings.Contains(tt.input, "func($a + y as $b)") {
				assert.Error(t, err, "Should error on nested alias in %s", tt.input)
				assert.Contains(t, err.Error(), "pipe expressions cannot be sub-expressions",
					"Should give correct error message for nested alias")
				return
			}

			assert.NoError(t, err, "Parsing should not produce an error for input: %s", tt.input)

			// Verify it's a pipe expression
			pipeExpr, ok := ast.(*parser.PipeExpression)
			assert.True(t, ok, "Expected a PipeExpression for input: %s", tt.input)
			assert.NotEmpty(t, pipeExpr.Aliases, "Should have at least one alias")
		})
	}
}

func TestParserTrial(t *testing.T) {
	input := "[1, x as $a, 2]" // This input is should fail becuase alias cannot be a sub-expression
	p := parser.NewParser(input)
	ast, err := p.Parse()
	assert.Error(t, err, "Parsing should not produce an error")

	utils.PrintJSON(ast)
}

// TestEmptyPipeExpressions ensures empty pipe expressions are rejected
func TestEmptyPipeExpressions(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"|: x + 1", "empty pipe expression is not allowed"},
		{"x + 1 |: |map: y + 2", "empty pipe expression is not allowed"},
		{"x + 1 |map:", "empty pipe expression is not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			assert.Error(t, err, "Expected error for input: %s", tt.input)
			assert.Contains(t, err.Error(), tt.expectedErr,
				"Expected error containing %q, but got %q", tt.expectedErr, err.Error())
		})
	}
}

// TestEmptyPipeWithAlias ensures empty pipe expressions cannot have aliases
func TestEmptyPipeWithAlias(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"x + 2 |: as $a", "empty pipe expression cannot have an alias"},
		{"|: as $a", "empty pipe expression cannot have an alias"},
		{"x + 1 |map: as $b", "empty pipe expression cannot have an alias"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			_, err := p.Parse()
			assert.Error(t, err, "Expected error for input: %s", tt.input)
			assert.Contains(t, err.Error(), tt.expectedErr,
				"Expected error containing %q, but got %q", tt.expectedErr, err.Error())
		})
	}
}
