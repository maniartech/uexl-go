package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	p, err := parser.ParseString("arr.0.2.3.4")
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	fmt.Println("Parsed AST for p:")
	utils.PrintJSON(p)
}
