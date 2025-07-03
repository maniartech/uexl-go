package parser

import "github.com/maniartech/uexl_go/ast"

func parseArray(vals interface{}, text []byte, offset, line, col int) (ast.Node, error) {
	// Handle empty array case - vals is nil when array is empty
	if vals == nil {
		return ast.NewArrayNode(string(text), []ast.Node{}, offset, line, col), nil
	}

	valsSl := toIfaceSlice(vals)

	if len(valsSl) == 0 {
		return ast.NewArrayNode(string(text), []ast.Node{}, offset, line, col), nil
	}

	// Check if the first element is an empty identifier (indicates empty array)
	if len(valsSl) >= 1 {
		if node, ok := valsSl[0].(ast.Node); ok {
			if id, isIdentifier := node.(*ast.IdentifierNode); isIdentifier {
				if id.Name == "" { // Empty identifier means empty array
					return ast.NewArrayNode(string(text), []ast.Node{}, offset, line, col), nil
				}
			}
		}
	}

	res := []interface{}{valsSl[0]}
	restSl := toIfaceSlice(valsSl[1])

	for _, v := range restSl {
		vSl := toIfaceSlice(v)
		res = append(res, vSl[2])
	}

	return ast.NewArrayNode(string(text), toNodesSlice(res), offset, line, col), nil
}
