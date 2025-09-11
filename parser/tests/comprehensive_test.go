package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

// TestTokenizerInternalMethods tests tokenizer functionality through public interface
func TestTokenizerInternalMethods(t *testing.T) {
	tokenizer := parser.NewTokenizer("hello world 123")

	// Test through NextToken which uses internal methods
	token1, err1 := tokenizer.NextToken()
	assert.NoError(t, err1)
	assert.Equal(t, constants.TokenIdentifier, token1.Type)
	assert.Equal(t, "hello", token1.Value.Str)

	token2, err2 := tokenizer.NextToken()
	assert.NoError(t, err2)
	assert.Equal(t, constants.TokenIdentifier, token2.Type)
	assert.Equal(t, "world", token2.Value.Str)

	token3, err3 := tokenizer.NextToken()
	assert.NoError(t, err3)
	assert.Equal(t, constants.TokenNumber, token3.Type)
	assert.Equal(t, 123.0, token3.Value.Num)
}

// TestTokenizerStringProcessing tests string processing through public interface
func TestTokenizerStringProcessing(t *testing.T) {
	// Test escape sequences through actual string parsing
	tokenizer := parser.NewTokenizer(`"hello\nworld\t!"`)
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenString, token.Type)
	assert.Equal(t, "hello\nworld\t!", token.Value.Str)

	// Test raw string with doubled quotes
	tokenizer2 := parser.NewTokenizer(`r"hello""world"`)
	token2, err2 := tokenizer2.NextToken()
	assert.NoError(t, err2)
	assert.Equal(t, constants.TokenString, token2.Type)
	assert.Equal(t, "hello\"world", token2.Value.Str)

	// Test invalid escape sequence
	tokenizer3 := parser.NewTokenizer(`"hello\zworld"`)
	_, err3 := tokenizer3.NextToken()
	assert.Error(t, err3) // Should error on invalid escape
}

// TestParserInternalMethods tests parser error handling through public interface
func TestParserInternalMethods(t *testing.T) {
	// Test error accumulation through invalid syntax
	p := parser.NewParser("1 + + 2")
	expr, err := p.Parse()
	assert.Error(t, err)
	assert.Nil(t, expr)

	// Test that errors are properly formatted
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unexpected")
}

// TestParserPipeHandling tests pipe-related parsing
func TestParserPipeHandling(t *testing.T) {
	// Test simple expression that should parse successfully
	p := parser.NewParser("$x + 1")
	expr, err := p.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, expr)

	// Test binary expression
	binary, ok := expr.(*parser.BinaryExpression)
	assert.True(t, ok)
	assert.Equal(t, "+", binary.Operator)
}

// TestParserAdvancedErrorCases tests advanced error scenarios
func TestParserAdvancedErrorCases(t *testing.T) {
	// Test parser with tokenizer error
	p := parser.NewParser("@invalid")
	expr, err := p.Parse()
	assert.Error(t, err)
	assert.Nil(t, expr)

	// Test parser with multiple errors
	p2 := parser.NewParser("(((")
	expr2, err2 := p2.Parse()
	assert.Error(t, err2)
	assert.Nil(t, expr2)
}

// TestTokenizerComplexNumbers tests complex number parsing scenarios
func TestTokenizerComplexNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"very large number", "123456789012345", 123456789012345.0},
		{"very small decimal", "0.000000001", 0.000000001},
		{"scientific with decimal", "1.23e-4", 0.000123},
		{"scientific uppercase", "1.23E+4", 12300.0},
		{"zero with decimal", "0.0", 0.0},
		{"integer with trailing zeros", "1000", 1000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, constants.TokenNumber, token.Type)
			assert.Equal(t, tt.expected, token.Value.Num)
		})
	}
}

// TestTokenizerAdvancedOperators tests advanced operator parsing
func TestTokenizerAdvancedOperators(t *testing.T) {
	// Test consecutive minus signs
	tokenizer := parser.NewTokenizer("---5")

	// Should get three separate minus tokens
	for i := 0; i < 3; i++ {
		token, err := tokenizer.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, constants.TokenOperator, token.Type)
		assert.Equal(t, "-", token.Value.Str)
	}

	// Then the number
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, token.Type)
	assert.Equal(t, 5.0, token.Value.Num)
}

// TestTokenizerPrintTokens tests the PrintTokens method
func TestTokenizerPrintTokens(t *testing.T) {
	tokenizer := parser.NewTokenizer("1 + 2")

	// This should not panic
	assert.NotPanics(t, func() {
		tokenizer.PrintTokens()
	})
}
