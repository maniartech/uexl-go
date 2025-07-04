package main

import (
	"fmt"

	"github.com/maniartech/uexl_go/parser"
)

func main() {
	testCases := []string{
		"obj.method().chain()",
		"user.getName().length",
		"obj.func().prop",
	}

	for _, testCase := range testCases {
		fmt.Printf("Testing: %s\n", testCase)
		p := parser.NewParser(testCase)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Parsed successfully as: %T\n", expr)
		}
		fmt.Println()
	}
}
