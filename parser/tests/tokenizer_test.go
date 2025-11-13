package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	. "github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/constants"
)

func TestTokenizerBasic(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []Token
	}{
		{
			name:  "String literals",
			input: "\"hello\" + 'world'",
			tokens: []Token{
				{Type: constants.TokenString, Value: parser.TokenValue{Kind: parser.TVKString, Str: "hello"}, Token: "\"hello\"", Line: 1, Column: 1},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "+"}, Token: "+", Line: 1, Column: 9},
				{Type: constants.TokenString, Value: parser.TokenValue{Kind: parser.TVKString, Str: "world"}, Token: "'world'", Line: 1, Column: 11},
			},
		},
		{
			name:  "Nullish and ternary tokens",
			input: "a ?? b ? c : d",
			tokens: []Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "a"}, Token: "a", Line: 1, Column: 1},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "??"}, Token: "??", Line: 1, Column: 3},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "b"}, Token: "b", Line: 1, Column: 6},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "?"}, Token: "?", Line: 1, Column: 8},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "c"}, Token: "c", Line: 1, Column: 10},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":", Line: 1, Column: 12},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "d"}, Token: "d", Line: 1, Column: 14},
			},
		},
		{
			name:  "Numbers and operators",
			input: "3.14 + x * (y - z)",
			tokens: []Token{
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: 3.14}, Token: "3.14", Line: 1, Column: 1},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "+"}, Token: "+", Line: 1, Column: 6},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "x"}, Token: "x", Line: 1, Column: 8},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "*"}, Token: "*", Line: 1, Column: 10},
				{Type: constants.TokenLeftParen, Value: parser.TokenValue{Kind: parser.TVKString, Str: "("}, Token: "(", Line: 1, Column: 12},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "y"}, Token: "y", Line: 1, Column: 13},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "-"}, Token: "-", Line: 1, Column: 15},
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "z"}, Token: "z", Line: 1, Column: 17},
				{Type: constants.TokenRightParen, Value: parser.TokenValue{Kind: parser.TVKString, Str: ")"}, Token: ")", Line: 1, Column: 18},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewTokenizer(test.input)
			tokens := tokenizer.PreloadTokens()

			// Remove EOF token for comparison
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == constants.TokenEOF {
				tokens = tokens[:len(tokens)-1]
			}

			if len(tokens) != len(test.tokens) {
				t.Errorf("Expected %d tokens, got %d", len(test.tokens), len(tokens))
				return
			}

			for i, expected := range test.tokens {
				actual := tokens[i]
				if actual.Type != expected.Type {
					t.Errorf("Token %d: expected type %v, got %v", i, expected.Type, actual.Type)
				}
				if actual.Value != expected.Value {
					t.Errorf("Token %d: expected value %v, got %v", i, expected.Value, actual.Value)
				}
				if actual.Token != expected.Token {
					t.Errorf("Token %d: expected token %v, got %v", i, expected.Token, actual.Token)
				}
			}
		})
	}
}
