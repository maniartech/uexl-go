package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/parser"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected parser.Expression
	}{
		// ... (keep existing test cases)
		{
			input: "a.b |: func(1, 'test')",
			expected: &parser.PipeExpression{
				Left: &parser.MemberAccess{
					Object:   &parser.Identifier{Name: "a", Line: 1, Column: 1},
					Property: "b",
					Line:     1,
					Column:   2,
				},
				Right: &parser.FunctionCall{
					Function: &parser.Identifier{Name: "func", Line: 1, Column: 8},
					Arguments: []parser.Expression{
						&parser.NumberLiteral{Value: "1", Line: 1, Column: 13},
						&parser.StringLiteral{Value: "test", Line: 1, Column: 16},
					},
					Line:   1,
					Column: 12,
				},
				Line:   1,
				Column: 5,
			},
		},
		{
			input: `{"a": 1, "b": true}`,
			expected: &parser.ObjectLiteral{
				Properties: map[string]parser.Expression{
					"a": &parser.NumberLiteral{Value: "1", Line: 1, Column: 7},
					"b": &parser.BooleanLiteral{Value: true, Line: 1, Column: 15},
				},
				Line:   1,
				Column: 1,
			},
		},
		{
			input: "[1, 2, 3] |map: (1 + 2)",
			expected: &parser.PipeExpression{
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
			},
		},
		{
			input: "x + 2 |: $1 + 3",
			expected: &parser.PipeExpression{
				Left: &parser.BinaryExpression{
					Left:     &parser.Identifier{Name: "x", Line: 1, Column: 1},
					Operator: "+",
					Right:    &parser.NumberLiteral{Value: "2", Line: 1, Column: 5},
					Line:     1,
					Column:   3,
				},
				Right: &parser.BinaryExpression{
					Left:     &parser.Identifier{Name: "$1", Line: 1, Column: 10},
					Operator: "+",
					Right:    &parser.NumberLiteral{Value: "3", Line: 1, Column: 15},
					Line:     1,
					Column:   13,
				},
				Line:   1,
				Column: 7,
			},
		},
		{
			input: "[1, 2, 3, 4, 5, 6, 7, 8, 9] |filter: $1 % 2 == 0",
			expected: &parser.PipeExpression{
				Left: &parser.ArrayLiteral{
					Elements: []parser.Expression{
						&parser.NumberLiteral{Value: "1", Line: 1, Column: 2},
						&parser.NumberLiteral{Value: "2", Line: 1, Column: 5},
						&parser.NumberLiteral{Value: "3", Line: 1, Column: 8},
						&parser.NumberLiteral{Value: "4", Line: 1, Column: 11},
						&parser.NumberLiteral{Value: "5", Line: 1, Column: 14},
						&parser.NumberLiteral{Value: "6", Line: 1, Column: 17},
						&parser.NumberLiteral{Value: "7", Line: 1, Column: 20},
						&parser.NumberLiteral{Value: "8", Line: 1, Column: 23},
						&parser.NumberLiteral{Value: "9", Line: 1, Column: 26},
					},
					Line:   1,
					Column: 1,
				},
				PipeType: "filter",
				Right: &parser.BinaryExpression{
					Left: &parser.BinaryExpression{
						Left:     &parser.Identifier{Name: "$1", Line: 1, Column: 37},
						Operator: "%",
						Right:    &parser.NumberLiteral{Value: "2", Line: 1, Column: 42},
						Line:     1,
						Column:   40,
					},
					Operator: "==",
					Right:    &parser.NumberLiteral{Value: "0", Line: 1, Column: 47},
					Line:     1,
					Column:   44,
				},
				Line:   1,
				Column: 29,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			actual, err := p.Parse()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("For input %q, expected AST %+v, but got %+v", tt.input, tt.expected, actual)
			}
		})
	}
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
			input:       `{"a": 1,}`,
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
		{"x.y.z(1, 2, 3) |: func(a, b)"},
		{"[1, 2, 3].map(x => x * 2).filter(x => x > 3)"},
		{`{"a": 1, "b": [true, false], "c": {"d": null}}`},
		{"a && b || c && d || e && f"},
		{"x << 2 + y >> 3 - z & 0xFF | 0x0F"},
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
		{"obj.method(1, 2) |: $1.property |: upperCase"},
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
