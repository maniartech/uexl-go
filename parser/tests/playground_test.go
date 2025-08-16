package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	input := "!!true"

	expr, err := parser.ParseString(input)
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	utils.PrintJSON(expr)
}
