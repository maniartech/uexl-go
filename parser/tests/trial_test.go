package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	. "github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	input := "obj.a.b.c + func('10')"

	parser := NewParser(input)
	expr, err := parser.Parse()

	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	utils.PrintJSON(expr)
}
