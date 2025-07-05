package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

func TestNumberTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parser.Token
	}{
		{
			name:  "Simple integer",
			input: "42",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  float64(42),
				Token:  "42",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Zero",
			input: "0",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  float64(0),
				Token:  "0",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Floating point",
			input: "3.14",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  3.14,
				Token:  "3.14",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Floating point with leading zero",
			input: "0.5",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  0.5,
				Token:  "0.5",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Scientific notation lowercase e",
			input: "1e3",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  1000.0,
				Token:  "1e3",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Scientific notation uppercase E",
			input: "2E4",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  20000.0,
				Token:  "2E4",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Scientific notation with positive exponent",
			input: "1.5e+2",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  150.0,
				Token:  "1.5e+2",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Scientific notation with negative exponent",
			input: "1.5e-2",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  0.015,
				Token:  "1.5e-2",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Large integer",
			input: "123456789",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  123456789.0,
				Token:  "123456789",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Very small decimal",
			input: "0.000001",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  0.000001,
				Token:  "0.000001",
				Line:   1,
				Column: 1,
			},
		},
		{
			name:  "Number with multiple decimal places",
			input: "123.456789",
			expected: parser.Token{
				Type:   constants.TokenNumber,
				Value:  123.456789,
				Token:  "123.456789",
				Line:   1,
				Column: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token := tokenizer.NextToken()

			assert.Equal(t, tt.expected.Type, token.Type, "Token type should match")
			assert.Equal(t, tt.expected.Value, token.Value, "Token value should match")
			assert.Equal(t, tt.expected.Token, token.Token, "Token string should match")
			assert.Equal(t, tt.expected.Line, token.Line, "Line should match")
			assert.Equal(t, tt.expected.Column, token.Column, "Column should match")
		})
	}
}

func TestNumberParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple integer",
			input:    "42",
			expected: "42",
		},
		{
			name:     "Zero",
			input:    "0",
			expected: "0",
		},
		{
			name:     "Floating point",
			input:    "3.14",
			expected: "3.14",
		},
		{
			name:     "Scientific notation",
			input:    "1e3",
			expected: "1e3",
		},
		{
			name:     "Scientific notation with decimal",
			input:    "1.5e-2",
			expected: "1.5e-2",
		},
		{
			name:     "Large number",
			input:    "123456789.987654321",
			expected: "123456789.987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err, "Parsing should not produce an error")

			numLit, ok := expr.(*parser.NumberLiteral)
			assert.True(t, ok, "Expected NumberLiteral")
			assert.Equal(t, tt.expected, numLit.Value, "Number value should match")
		})
	}
}

func TestNumberInExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Addition with integers",
			input: "1 + 2",
		},
		{
			name:  "Subtraction with floats",
			input: "3.14 - 2.71",
		},
		{
			name:  "Multiplication with scientific notation",
			input: "1e3 * 2e-2",
		},
		{
			name:  "Division with mixed types",
			input: "10 / 3.0",
		},
		{
			name:  "Complex expression",
			input: "(1.5 + 2.5) * 3e2 / 1.2e1",
		},
		{
			name:  "Negative numbers with unary minus",
			input: "-42 + 10",
		},
		{
			name:  "Large numbers",
			input: "123456789 + 987654321.123456",
		},
		{
			name:  "Very small numbers",
			input: "0.000001 + 1e-6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()
			assert.NoError(t, err, "Parsing should not produce an error for: %s", tt.input)
			assert.NotNil(t, expr, "Expression should not be nil")
		})
	}
}

func TestNumberErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Multiple decimal points",
			input:       "3.14.15",
			expectError: true,
		},
		{
			name:        "Number followed by identifier (not error case)",
			input:       "1e",
			expectError: false, // This should parse as "1" followed by identifier "e"
		},
		{
			name:        "Valid edge case - number followed by identifier",
			input:       "123abc",
			expectError: false, // Should parse as number "123" followed by identifier "abc"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token := tokenizer.NextToken()

			if tt.expectError {
				assert.Equal(t, constants.TokenError, token.Type, "Expected error token")
			} else {
				// For valid cases, just check it's not an error token
				assert.NotEqual(t, constants.TokenError, token.Type, "Should not be error token")
			}
		})
	}
}

func TestNumberPrecision(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Integer precision",
			input:    "9007199254740991", // Max safe integer in JavaScript (2^53 - 1)
			expected: 9007199254740991.0,
		},
		{
			name:     "Float precision",
			input:    "0.1",
			expected: 0.1,
		},
		{
			name:     "Small scientific notation",
			input:    "1e-10",
			expected: 1e-10,
		},
		{
			name:     "Large scientific notation",
			input:    "1e10",
			expected: 1e10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token := tokenizer.NextToken()

			assert.Equal(t, constants.TokenNumber, token.Type, "Should be a number token")
			assert.Equal(t, tt.expected, token.Value, "Value should match expected precision")
		})
	}
}

func TestNumberInArrays(t *testing.T) {
	input := "[1, 2.5, 3e2, 0.001, -42]"
	p := parser.NewParser(input)
	expr, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error")

	arrayLit, ok := expr.(*parser.ArrayLiteral)
	assert.True(t, ok, "Expected ArrayLiteral")
	assert.Len(t, arrayLit.Elements, 5, "Array should have 5 elements")

	// Check first element (1)
	numLit1, ok := arrayLit.Elements[0].(*parser.NumberLiteral)
	assert.True(t, ok, "First element should be NumberLiteral")
	assert.Equal(t, "1", numLit1.Value)

	// Check second element (2.5)
	numLit2, ok := arrayLit.Elements[1].(*parser.NumberLiteral)
	assert.True(t, ok, "Second element should be NumberLiteral")
	assert.Equal(t, "2.5", numLit2.Value)

	// Check third element (3e2)
	numLit3, ok := arrayLit.Elements[2].(*parser.NumberLiteral)
	assert.True(t, ok, "Third element should be NumberLiteral")
	assert.Equal(t, "3e2", numLit3.Value)
}

func TestNumberInObjects(t *testing.T) {
	input := `{"int": 42, "float": 3.14, "sci": 1e3}`
	p := parser.NewParser(input)
	expr, err := p.Parse()
	assert.NoError(t, err, "Parsing should not produce an error")

	objLit, ok := expr.(*parser.ObjectLiteral)
	assert.True(t, ok, "Expected ObjectLiteral")
	assert.Len(t, objLit.Properties, 3, "Object should have 3 properties")

	// Check int property
	intVal, exists := objLit.Properties["int"]
	assert.True(t, exists, "int property should exist")
	numLit1, ok := intVal.(*parser.NumberLiteral)
	assert.True(t, ok, "int value should be NumberLiteral")
	assert.Equal(t, "42", numLit1.Value)

	// Check float property
	floatVal, exists := objLit.Properties["float"]
	assert.True(t, exists, "float property should exist")
	numLit2, ok := floatVal.(*parser.NumberLiteral)
	assert.True(t, ok, "float value should be NumberLiteral")
	assert.Equal(t, "3.14", numLit2.Value)

	// Check sci property
	sciVal, exists := objLit.Properties["sci"]
	assert.True(t, exists, "sci property should exist")
	numLit3, ok := sciVal.(*parser.NumberLiteral)
	assert.True(t, ok, "sci value should be NumberLiteral")
	assert.Equal(t, "1e3", numLit3.Value)
}
