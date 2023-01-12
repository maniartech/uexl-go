package parser

import (
	"strings"

	"github.com/maniartech/uexl_go/ast"
)

// ParseString parses the given expression and returns the AST Node.
// It allows you to parse the expression without having to create a file.
// For example:
//
//	node, err := ParseString("1 + 2")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(node.Eval(nil))
func ParseString(expr string) (ast.Node, error) {
	node, err := ParseReader("", strings.NewReader(expr))
	if err != nil {
		return nil, err
	}

	return node.(ast.Node), nil
}
