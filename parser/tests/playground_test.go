package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	p, err := parser.ParseString("arr[0]['address']['street']")
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	fmt.Println("Parsed AST for p:")
	utils.PrintJSON(p)

	p2, err := parser.ParseString("arr.0.address.street.name")
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	fmt.Println("Parsed AST for p2:")
	utils.PrintJSON(p2)
}
