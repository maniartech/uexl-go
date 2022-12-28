package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/maniartech/uexl_go/ast"
	"github.com/maniartech/uexl_go/parser"
)

// for testing purpose

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: calculator 'EXPR'")
	}

	got, err := parser.ParseReader("", strings.NewReader(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Evaluation pending, printing the AST")
	fmt.Println("====================================")
	ast.PrintNode(got.(ast.Node))
}
