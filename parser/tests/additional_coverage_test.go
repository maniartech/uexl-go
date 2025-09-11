package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

// TestTokenizerEdgeCases tests edge cases in tokenizer
func TestTokenizerEdgeCases(t *testing.T) {
	// Test empty input
	tokenizer := parser.NewTokenizer("")
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, constants.TokenEOF, token.Type)

	// Test whitespace only
	tokenizer2 := parser.NewTokenizer("   \n\t  ")
	token2, err2 := tokenizer2.NextToken()
	assert.NoError(t, err2)
	assert.Equal(t, constants.TokenEOF, token2.Type)

	// Test invalid character
	tokenizer3 := parser.NewTokenizer("@")
	_, err3 := tokenizer3.NextToken()
	assert.Error(t, err3)
}

// TestTokenizerNumberEdgeCases tests number parsing edge cases
func TestTokenizerNumberEdgeCases(t *testing.T) {
	// Test number followed by dot (member access)
	tokenizer := parser.NewTokenizer("123.toString")

	// Should get number token
	token1, err1 := tokenizer.NextToken()
	assert.NoError(t, err1)
	assert.Equal(t, constants.TokenNumber, token1.Type)
	assert.Equal(t, 123.0, token1.Value.Num)

	// Should get dot token
	token2, err2 := tokenizer.NextToken()
	assert.NoError(t, err2)
	assert.Equal(t, constants.TokenDot, token2.Type)

	// Should get identifier token
	token3, err3 := tokenizer.NextToken()
	assert.NoError(t, err3)
	assert.Equal(t, constants.TokenIdentifier, token3.Type)
	assert.Equal(t, "toString", token3.Value.Str)
}

// TestTokenizerStringEdgeCases tests string parsing edge cases
func TestTokenizerStringEdgeCases(t *testing.T) {
	// Test unterminated string
	tokenizer := parser.NewTokenizer(`"hello`)
	_, err := tokenizer.NextToken()
	assert.Error(t, err)

	// Test raw string with doubled quotes
	tokenizer2 := parser.NewTokenizer(`r"hello""world"`)
	token, err2 := tokenizer2.NextToken()
	assert.NoError(t, err2)
	assert.Equal(t, "hello\"world", token.Value.Str)
}

// TestParserErrorPaths tests various error paths in parser
func TestParserErrorPaths(t *testing.T) {
	// Test leading pipe error
	p := parser.NewParser("|: x + 1")
	expr, err := p.Parse()
	assert.Error(t, err)
	assert.Nil(t, expr)

	// Test empty pipe with alias error
	p2 := parser.NewParser("|: as $a")
	expr2, err2 := p2.Parse()
	assert.Error(t, err2)
	assert.Nil(t, expr2)
}
