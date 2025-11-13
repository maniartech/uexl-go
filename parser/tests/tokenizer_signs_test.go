package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/constants"
)

func collectTypesTokens(tz *parser.Tokenizer) ([]constants.TokenType, []string, error) {
	types := []constants.TokenType{}
	toks := []string{}
	for {
		tok, err := tz.NextToken()
		if err != nil {
			return nil, nil, err
		}
		types = append(types, tok.Type)
		toks = append(toks, tok.Token)
		if tok.Type == constants.TokenEOF {
			break
		}
	}
	return types, toks, nil
}

func TestTokenizer_LeadingSignSeparated(t *testing.T) {
	cases := []struct {
		input      string
		wantTokens []string
		wantTypes  []constants.TokenType
	}{
		{
			input:      "-123",
			wantTokens: []string{"-", "123", ""},
			wantTypes:  []constants.TokenType{constants.TokenOperator, constants.TokenNumber, constants.TokenEOF},
		},
		{
			input:      "+456",
			wantTokens: []string{"+", "456", ""},
			wantTypes:  []constants.TokenType{constants.TokenOperator, constants.TokenNumber, constants.TokenEOF},
		},
		{
			input:      "(-789)",
			wantTokens: []string{"(", "-", "789", ")", ""},
			wantTypes:  []constants.TokenType{constants.TokenLeftParen, constants.TokenOperator, constants.TokenNumber, constants.TokenRightParen, constants.TokenEOF},
		},
		{
			input: "--10",
			// Intentionally two '-' tokens then number
			wantTokens: []string{"-", "-", "10", ""},
			wantTypes:  []constants.TokenType{constants.TokenOperator, constants.TokenOperator, constants.TokenNumber, constants.TokenEOF},
		},
	}

	for _, tc := range cases {
		tz := parser.NewTokenizer(tc.input)
		types, toks, err := collectTypesTokens(tz)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if len(types) != len(tc.wantTypes) {
			t.Fatalf("input %q: got %d tokens, want %d; got %v", tc.input, len(types), len(tc.wantTypes), toks)
		}
		for i := range types {
			if types[i] != tc.wantTypes[i] || toks[i] != tc.wantTokens[i] {
				t.Fatalf("input %q: token %d mismatch: got (%v,%q) want (%v,%q); all: %v", tc.input, i, types[i], toks[i], tc.wantTypes[i], tc.wantTokens[i], toks)
			}
		}
	}
}

func TestTokenizer_ExponentSignWithinNumber(t *testing.T) {
	cases := []struct {
		input     string
		wantToken string
	}{
		{input: "1e-5", wantToken: "1e-5"},
		{input: "1E+6", wantToken: "1E+6"},
		{input: "123.45e-2", wantToken: "123.45e-2"},
	}

	for _, tc := range cases {
		tz := parser.NewTokenizer(tc.input)
		tok, err := tz.NextToken()
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", tc.input, err)
		}
		if tok.Type != constants.TokenNumber || tok.Token != tc.wantToken {
			t.Fatalf("%q: got (%v,%q), want (TokenNumber,%q)", tc.input, tok.Type, tok.Token, tc.wantToken)
		}
		tok2, err := tz.NextToken()
		if err != nil {
			t.Fatalf("%q: unexpected error 2: %v", tc.input, err)
		}
		if tok2.Type != constants.TokenEOF {
			t.Fatalf("%q: expected EOF, got %v", tc.input, tok2.Type)
		}
	}
}
