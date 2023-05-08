package main

import (
	"fmt"
	"log"
	"os"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

// for testing purpose

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: calculator 'EXPR'")
	}

	expr := os.Args[1]

	node, err := parser.ParseString(expr)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Evaluation pending, printing the AST")
	// fmt.Println("====================================")
	// ast.PrintNode(node)

	result, err := node.Eval(types.Context{
		"x": types.Number(10),
		"y": types.Object{
			"z": types.Number(5),
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(expr, "=", result)

}
