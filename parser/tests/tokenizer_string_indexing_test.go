package parser_test

import (
	"testing"

	. "github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
)

func TestTokenizer_ForStringIndexingTokens(t *testing.T) {
	tokenizer := NewTokenizer(`"hello"[2]`)
	tokens := tokenizer.PreloadTokens()

	// Expect: String, LeftBracket, Number, RightBracket, EOF
	if len(tokens) < 5 {
		t.Fatalf("expected at least 5 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != constants.TokenString || tokens[0].Value != "hello" {
		t.Fatalf("token 0 expected string 'hello', got %v %v", tokens[0].Type, tokens[0].Value)
	}
	if tokens[1].Type != constants.TokenLeftBracket {
		t.Fatalf("token 1 expected LeftBracket, got %v", tokens[1].Type)
	}
	if tokens[2].Type != constants.TokenNumber || tokens[2].Value != 2.0 {
		t.Fatalf("token 2 expected number 2, got %v %v", tokens[2].Type, tokens[2].Value)
	}
	if tokens[3].Type != constants.TokenRightBracket {
		t.Fatalf("token 3 expected RightBracket, got %v", tokens[3].Type)
	}
}
