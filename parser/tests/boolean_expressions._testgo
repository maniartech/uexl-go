package parser_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
)

func TestBooleanExpressions(t *testing.T) {
	expressions := []string{
		"true",
		"false",
		"true && false",
		"true || false",
		"(true || false) && true",
		"true && (false || true)",
		"(true || true) && (false || true)",
		"!true",
		"!false",
		"!(true && false)",
		"!(true || false)",
	}

	for _, expr := range expressions {
		fmt.Printf("Expression: %s\n", expr)

		// Parse the expression
		node, err := parser.ParseString(expr)
		if err != nil {
			log.Printf("Error converting '%s': %v", expr, err)
			continue
		}

		// Evaluate the expression
		utils.PrintJSON(node)
	}
}
