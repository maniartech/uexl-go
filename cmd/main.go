package main

import (
	"log"
	"os"

	"github.com/maniartech/uexl_go/internal/utils"
	"github.com/maniartech/uexl_go/parser"
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

	utils.PrintJSON(node)
}
