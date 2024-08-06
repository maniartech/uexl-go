package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input    string
		expected []parser.Token
	}{
		{
			input: "3.14 + x * (y - z)",
			expected: []parser.Token{
				{Type: parser.TokenNumber, Value: "3.14", Line: 1, Column: 1},
				{Type: parser.TokenOperator, Value: "+", Line: 1, Column: 6},
				{Type: parser.TokenIdentifier, Value: "x", Line: 1, Column: 8},
				{Type: parser.TokenOperator, Value: "*", Line: 1, Column: 10},
				{Type: parser.TokenLeftParen, Value: "(", Line: 1, Column: 12},
				{Type: parser.TokenIdentifier, Value: "y", Line: 1, Column: 13},
				{Type: parser.TokenOperator, Value: "-", Line: 1, Column: 15},
				{Type: parser.TokenIdentifier, Value: "z", Line: 1, Column: 17},
				{Type: parser.TokenRightParen, Value: ")", Line: 1, Column: 18},
				{Type: parser.TokenEOF, Line: 1, Column: 19},
			},
		},
		{
			input: "true && false || null",
			expected: []parser.Token{
				{Type: parser.TokenBoolean, Value: "true", Line: 1, Column: 1},
				{Type: parser.TokenOperator, Value: "&&", Line: 1, Column: 6},
				{Type: parser.TokenBoolean, Value: "false", Line: 1, Column: 9},
				{Type: parser.TokenOperator, Value: "||", Line: 1, Column: 15},
				{Type: parser.TokenNull, Value: "null", Line: 1, Column: 18},
				{Type: parser.TokenEOF, Line: 1, Column: 22},
			},
		},
		{
			input: "a.b |: func(1, 'test')",
			expected: []parser.Token{
				{Type: parser.TokenIdentifier, Value: "a", Line: 1, Column: 1},
				{Type: parser.TokenDot, Value: ".", Line: 1, Column: 2},
				{Type: parser.TokenIdentifier, Value: "b", Line: 1, Column: 3},
				{Type: parser.TokenPipe, Value: "|:", Line: 1, Column: 5},
				{Type: parser.TokenIdentifier, Value: "func", Line: 1, Column: 8},
				{Type: parser.TokenLeftParen, Value: "(", Line: 1, Column: 12},
				{Type: parser.TokenNumber, Value: "1", Line: 1, Column: 13},
				{Type: parser.TokenComma, Value: ",", Line: 1, Column: 14},
				{Type: parser.TokenString, Value: "'test'", Line: 1, Column: 16},
				{Type: parser.TokenRightParen, Value: ")", Line: 1, Column: 22},
				{Type: parser.TokenEOF, Line: 1, Column: 23},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokenizer := parser.NewTokenizer(tt.input)
			for _, expected := range tt.expected {
				actual := tokenizer.NextToken()
				if actual.Type != expected.Type || actual.Value != expected.Value ||
					actual.Line != expected.Line || actual.Column != expected.Column {
					t.Errorf("For input %q, expected token %+v, but got %+v", tt.input, expected, actual)
				}
			}
		})
	}
}
