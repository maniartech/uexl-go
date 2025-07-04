package main

import (
	"fmt"
	"log"

	"github.com/maniartech/uexl_go/parser"
)

func main() {
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
		p := parser.NewParser(expr)
		ast, err := p.Parse()
		if err != nil {
			log.Printf("Error parsing '%s': %v", expr, err)
			continue
		}

		// Convert to evaluable AST
		node, err := parser.ParseString(expr)
		if err != nil {
			log.Printf("Error converting '%s': %v", expr, err)
			continue
		}

		// Evaluate the expression
		result, err := node.Eval(nil)
		if err != nil {
			log.Printf("Error evaluating '%s': %v", expr, err)
			continue
		}

		fmt.Printf("Result: %v\n", result)
		fmt.Printf("AST Type: %T\n", ast)
		fmt.Println("---")
	}
}
