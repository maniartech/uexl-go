package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
)

func TestTokenizationDebugNew(t *testing.T) {
	testCases := []string{
		"obj.method().chain()",
		"func(a, b).property",
		"user.profile.name",
		"getValue().length",
	}

	for _, testCase := range testCases {
		t.Logf("=== Tokenizing: %s ===", testCase)
		tokenizer := parser.NewTokenizer(testCase)

		for {
			token := tokenizer.NextToken()
			if token.Type == constants.TokenEOF {
				break
			}
			t.Logf("Token: '%s', Type: %v", token.Token, token.Type)
		}
		t.Log("")
	}
}
