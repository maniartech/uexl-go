package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	// Test cases for multiple consecutive operators
	testCases := []string{
		"--10",     // double negation
		"!!true",   // double NOT
		"!!!false", // triple NOT
		"!true",    // single NOT
		"!=",       // should fail (incomplete)
		"!!x",      // double NOT with variable
	}

	for i, input := range testCases {
		t.Logf("Testing case %d: %s", i+1, input)
		expr, err := parser.ParseString(input)
		if err != nil {
			t.Logf("Case %d failed: %v", i+1, err)
		} else {
			t.Logf("Case %d succeeded", i+1)
			utils.PrintJSON(expr)
		}
		t.Log("---")
	}
}
