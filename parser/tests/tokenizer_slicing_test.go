package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/stretchr/testify/assert"
)

func TestTokenizer_Slicing(t *testing.T) {
	tests := []struct {
		input    string
		expected []parser.Token
	}{
		{
			input: "arr[1:4]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(1)}, Token: "1"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(4)}, Token: "4"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[:3]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(3)}, Token: "3"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[1:]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(1)}, Token: "1"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[:]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[0:5:2]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(0)}, Token: "0"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(5)}, Token: "5"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(2)}, Token: "2"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[::-1]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "-"}, Token: "-"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(1)}, Token: "1"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
		{
			input: "arr[-3:-1]",
			expected: []parser.Token{
				{Type: constants.TokenIdentifier, Value: parser.TokenValue{Kind: parser.TVKIdentifier, Str: "arr"}, Token: "arr"},
				{Type: constants.TokenLeftBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "["}, Token: "["},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "-"}, Token: "-"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(3)}, Token: "3"},
				{Type: constants.TokenColon, Value: parser.TokenValue{Kind: parser.TVKString, Str: ":"}, Token: ":"},
				{Type: constants.TokenOperator, Value: parser.TokenValue{Kind: parser.TVKOperator, Str: "-"}, Token: "-"},
				{Type: constants.TokenNumber, Value: parser.TokenValue{Kind: parser.TVKNumber, Num: float64(1)}, Token: "1"},
				{Type: constants.TokenRightBracket, Value: parser.TokenValue{Kind: parser.TVKString, Str: "]"}, Token: "]"},
				{Type: constants.TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			tokens := tokenizer.PreloadTokens()

			assert.Equal(t, len(tt.expected), len(tokens), "Number of tokens does not match")

			for i, expectedToken := range tt.expected {
				if i >= len(tokens) {
					t.Errorf("Missing token at index %d. Expected %v", i, expectedToken)
					continue
				}
				actualToken := tokens[i]
				assert.Equal(t, expectedToken.Type, actualToken.Type, "Token type mismatch at index %d", i)
				assert.Equal(t, expectedToken.Value, actualToken.Value, "Token value mismatch at index %d", i)
				assert.Equal(t, expectedToken.Token, actualToken.Token, "Token literal mismatch at index %d", i)
			}
		})
	}
}
