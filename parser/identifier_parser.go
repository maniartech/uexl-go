package parser

import "github.com/maniartech/uexl_go/ast"

func parseIdentifier(token []byte, name, value interface{}, offset, line, col int) (ast.Node, error) {
	Name := resolveAscii1(name)
	Value := ""
	if len(value.([]interface{})) != 0 {
		Value = resolveAscii1(value.([]interface{})[0].([]interface{})[1:][0])
	}
	return ast.NewIdentifierNode(string(token), Name, Value, offset, line, col)
}
