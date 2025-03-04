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
			input: "\"hello\" + 'world'",
			expected: []parser.Token{
				{Type: parser.TokenString, Value: "\"hello\"", Line: 1, Column: 1},
				{Type: parser.TokenOperator, Value: "+", Line: 1, Column: 9},
				{Type: parser.TokenString, Value: "'world'", Line: 1, Column: 11},
				{Type: parser.TokenEOF, Line: 1, Column: 18},
			},
		},
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
				{Type: parser.TokenIdentifier, Value: "a.b", Line: 1, Column: 1},
				{Type: parser.TokenPipe, Value: ":", Line: 1, Column: 5},
				{Type: parser.TokenIdentifier, Value: "func", Line: 1, Column: 8},
				{Type: parser.TokenLeftParen, Value: "(", Line: 1, Column: 12},
				{Type: parser.TokenNumber, Value: "1", Line: 1, Column: 13},
				{Type: parser.TokenComma, Value: ",", Line: 1, Column: 14},
				{Type: parser.TokenString, Value: "'test'", Line: 1, Column: 16},
				{Type: parser.TokenRightParen, Value: ")", Line: 1, Column: 22},
				{Type: parser.TokenEOF, Line: 1, Column: 23},
			},
		},
		{
			// All cases of | operators, pipe, logical or and bitwise or
			input: "a | b |map: c || d |filter: x | y || z",
			expected: []parser.Token{
				{Type: parser.TokenIdentifier, Value: "a", Line: 1, Column: 1},
				{Type: parser.TokenOperator, Value: "|", Line: 1, Column: 3},
				{Type: parser.TokenIdentifier, Value: "b", Line: 1, Column: 5},
				{Type: parser.TokenPipe, Value: "map", Line: 1, Column: 7},
				{Type: parser.TokenIdentifier, Value: "c", Line: 1, Column: 13},
				{Type: parser.TokenOperator, Value: "||", Line: 1, Column: 15},
				{Type: parser.TokenIdentifier, Value: "d", Line: 1, Column: 18},
				{Type: parser.TokenPipe, Value: "filter", Line: 1, Column: 20},
				{Type: parser.TokenIdentifier, Value: "x", Line: 1, Column: 29},
				{Type: parser.TokenOperator, Value: "|", Line: 1, Column: 31},
				{Type: parser.TokenIdentifier, Value: "y", Line: 1, Column: 33},
				{Type: parser.TokenOperator, Value: "||", Line: 1, Column: 35},
				{Type: parser.TokenIdentifier, Value: "z", Line: 1, Column: 38},
				{Type: parser.TokenEOF, Line: 1, Column: 39},
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

func TestPipe(t *testing.T) {
	input := "a.b * 2 |map: $1.x.y * 2"
	expected := []parser.Token{
		{Type: parser.TokenIdentifier, Value: "a.b", Line: 1, Column: 1},
		{Type: parser.TokenOperator, Value: "*", Line: 1, Column: 5},
		{Type: parser.TokenNumber, Value: "2", Line: 1, Column: 7},
		{Type: parser.TokenPipe, Value: "map", Line: 1, Column: 9},
		{Type: parser.TokenIdentifier, Value: "$1.x.y", Line: 1, Column: 15},
		{Type: parser.TokenOperator, Value: "*", Line: 1, Column: 22},
		{Type: parser.TokenNumber, Value: "2", Line: 1, Column: 24},
		{Type: parser.TokenEOF, Line: 1, Column: 25},
	}
	tokenizer := parser.NewTokenizer(input)

	for _, expected := range expected {
		actual := tokenizer.NextToken()
		if actual.Type != expected.Type || actual.Value != expected.Value ||
			actual.Line != expected.Line || actual.Column != expected.Column {
			t.Errorf("For input %q, expected token %+v, but got %+v", input, expected, actual)
		}
	}
}

func TestTrial(t *testing.T) {
	input := "\"hello\" + 'world'"
	tokenizer := parser.NewTokenizer(input)

	tokenizer.PrintTokens()
}
