package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/constants"
	"github.com/stretchr/testify/assert"
)

// TestTokenizerAdvancedFeatures tests advanced tokenizer functionality
func TestTokenizerAdvancedFeatures(t *testing.T) {
	// Test raw strings
	tokenizer := parser.NewTokenizer(`r"hello\nworld"`)
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenString, token.Type)
	assert.Equal(t, "hello\\nworld", token.Value.Str) // raw string preserves backslashes

	// Test single quoted strings
	tokenizer = parser.NewTokenizer("'hello world'")
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenString, token.Type)
	assert.Equal(t, "hello world", token.Value.Str)
	assert.True(t, token.IsSingleQuoted)

	// Test double quoted strings
	tokenizer = parser.NewTokenizer(`"hello world"`)
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenString, token.Type)
	assert.Equal(t, "hello world", token.Value.Str)
	assert.False(t, token.IsSingleQuoted)
}

// TestTokenizerEscapeSequences tests string escape sequences
func TestTokenizerEscapeSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"newline", `"hello\nworld"`, "hello\nworld"},
		{"tab", `"hello\tworld"`, "hello\tworld"},
		{"carriage return", `"hello\rworld"`, "hello\rworld"},
		{"backslash", `"hello\\world"`, "hello\\world"},
		{"quote", `"hello\"world"`, "hello\"world"},
		{"unicode", `"hello\u0041world"`, "helloAworld"},
		{"unicode long", `"hello\U00000041world"`, "helloAworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, constants.TokenString, token.Type)
			assert.Equal(t, tt.expected, token.Value.Str)
		})
	}
}

// TestTokenizerNumbers tests various number formats
func TestTokenizerNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"integer", "42", 42.0},
		{"float", "3.14", 3.14},
		{"scientific lowercase", "1e5", 100000.0},
		{"scientific uppercase", "1E5", 100000.0},
		{"scientific with plus", "1e+5", 100000.0},
		{"scientific with minus", "1e-2", 0.01},
		{"scientific with decimal", "1.5e2", 150.0},
		{"zero", "0", 0.0},
		{"zero float", "0.0", 0.0},
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

// TestTokenizerOperators tests all operator types
func TestTokenizerOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"plus", "+", "+"},
		{"minus", "-", "-"},
		{"multiply", "*", "*"},
		{"divide", "/", "/"},
		{"modulo", "%", "%"},
		{"power", "**", "**"},
		{"equal", "==", "=="},
		{"not equal", "!=", "!="},
		{"less than", "<", "<"},
		{"greater than", ">", ">"},
		{"less equal", "<=", "<="},
		{"greater equal", ">=", ">="},
		{"logical and", "&&", "&&"},
		{"logical or", "||", "||"},
		{"bitwise and", "&", "&"},
		{"bitwise or", "|", "|"},
		{"bitwise xor", "^", "^"},
		{"left shift", "<<", "<<"},
		{"right shift", ">>", ">>"},
		{"nullish coalescing", "??", "??"},
		{"question", "?", "?"},
		{"not", "!", "!"},
		{"increment", "++", "++"},
		{"decrement", "--", "--"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, constants.TokenOperator, token.Type)
			assert.Equal(t, tt.expected, token.Value.Str)
		})
	}
}

// TestTokenizerOptionalChaining tests optional chaining tokens
func TestTokenizerOptionalChaining(t *testing.T) {
	// Test ?.
	tokenizer := parser.NewTokenizer("?.")
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenQuestionDot, token.Type)
	assert.Equal(t, "?.", token.Value.Str)

	// Test ?[
	tokenizer = parser.NewTokenizer("?[")
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenQuestionLeftBracket, token.Type)
	assert.Equal(t, "?[", token.Value.Str)
}

// TestTokenizerPipes tests pipe token variations
func TestTokenizerPipes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"default pipe", "|:", ":"},
		{"named pipe", "|map:", "map"},
		{"filter pipe", "|filter:", "filter"},
		{"reduce pipe", "|reduce:", "reduce"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, constants.TokenPipe, token.Type)
			assert.Equal(t, tt.expected, token.Value.Str)
		})
	}
}

// TestTokenizerKeywords tests keyword recognition
func TestTokenizerKeywords(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType constants.TokenType
		expectedVal  interface{}
	}{
		{"true", "true", constants.TokenBoolean, true},
		{"false", "false", constants.TokenBoolean, false},
		{"null", "null", constants.TokenNull, nil},
		{"as", "as", constants.TokenAs, "as"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, token.Type)

			switch tt.expectedType {
			case constants.TokenBoolean:
				assert.Equal(t, tt.expectedVal, token.Value.Bool)
			case constants.TokenNull:
				assert.Equal(t, parser.TVKNull, token.Value.Kind)
			case constants.TokenAs:
				assert.Equal(t, tt.expectedVal, token.Value.Str)
			}
		})
	}
}

// TestTokenizerDelimiters tests delimiter tokens
func TestTokenizerDelimiters(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType constants.TokenType
	}{
		{"left paren", "(", constants.TokenLeftParen},
		{"right paren", ")", constants.TokenRightParen},
		{"left bracket", "[", constants.TokenLeftBracket},
		{"right bracket", "]", constants.TokenRightBracket},
		{"left brace", "{", constants.TokenLeftBrace},
		{"right brace", "}", constants.TokenRightBrace},
		{"comma", ",", constants.TokenComma},
		{"dot", ".", constants.TokenDot},
		{"colon", ":", constants.TokenColon},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, token.Type)
			assert.Equal(t, tt.input, token.Value.Str)
		})
	}
}

// TestTokenizerIdentifiers tests identifier tokenization
func TestTokenizerIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "hello"},
		{"with underscore", "hello_world", "hello_world"},
		{"with dollar", "$variable", "$variable"},
		{"with numbers", "var123", "var123"},
		{"mixed", "_$test123", "_$test123"},
		{"single letter", "x", "x"},
		{"single underscore", "_", "_"},
		{"single dollar", "$", "$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, constants.TokenIdentifier, token.Type)
			assert.Equal(t, tt.expected, token.Value.Str)
		})
	}
}

// TestTokenizerErrorCases tests error conditions in tokenizer
func TestTokenizerErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unterminated string double", `"hello`},
		{"unterminated string single", `'hello`},
		{"unterminated raw string", `r"hello`},
		{"invalid character", "@"},
		{"invalid unicode escape", `"hello\u123"`}, // incomplete unicode
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			_, err := tokenizer.NextToken()
			assert.Error(t, err)
		})
	}
}

// TestTokenizerWhitespace tests whitespace handling
func TestTokenizerWhitespace(t *testing.T) {
	// Test that whitespace is properly skipped
	tokenizer := parser.NewTokenizer("  \t\n  42  \r\n  ")
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenNumber, token.Type)
	assert.Equal(t, 42.0, token.Value.Num)

	// Next token should be EOF
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenEOF, token.Type)
}

// TestTokenizerComplexExpressions tests tokenizing complex expressions
func TestTokenizerComplexExpressions(t *testing.T) {
	input := `obj.prop[0]?.method("test", 42) ?? defaultValue`
	tokenizer := parser.NewTokenizer(input)

	expectedTokens := []constants.TokenType{
		constants.TokenIdentifier,   // obj
		constants.TokenDot,          // .
		constants.TokenIdentifier,   // prop
		constants.TokenLeftBracket,  // [
		constants.TokenNumber,       // 0
		constants.TokenRightBracket, // ]
		constants.TokenQuestionDot,  // ?.
		constants.TokenIdentifier,   // method
		constants.TokenLeftParen,    // (
		constants.TokenString,       // "test"
		constants.TokenComma,        // ,
		constants.TokenNumber,       // 42
		constants.TokenRightParen,   // )
		constants.TokenOperator,     // ??
		constants.TokenIdentifier,   // defaultValue
		constants.TokenEOF,          // EOF
	}

	for i, expectedType := range expectedTokens {
		token, err := tokenizer.NextToken()
		assert.NoError(t, err, "Token %d failed", i)
		assert.Equal(t, expectedType, token.Type, "Token %d type mismatch", i)
	}
}

// TestTokenizerLineAndColumn tests line and column tracking
func TestTokenizerLineAndColumn(t *testing.T) {
	input := "hello\nworld\n  test"
	tokenizer := parser.NewTokenizer(input)

	// First token: "hello" at line 1, column 1
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, token.Type)
	assert.Equal(t, "hello", token.Value.Str)
	assert.Equal(t, 1, token.Line)
	assert.Equal(t, 1, token.Column)

	// Second token: "world" at line 2, column 1
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, token.Type)
	assert.Equal(t, "world", token.Value.Str)
	assert.Equal(t, 2, token.Line)
	assert.Equal(t, 1, token.Column)

	// Third token: "test" at line 3, column 3 (after 2 spaces)
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenIdentifier, token.Type)
	assert.Equal(t, "test", token.Value.Str)
	assert.Equal(t, 3, token.Line)
	assert.Equal(t, 3, token.Column)
}
