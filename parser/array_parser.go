package parser

import "github.com/maniartech/uexl_go/ast"

func parseArray(vals interface{}, text []byte, offset, line, col int) (ast.Node, error) {
	valsSl := ast.ToIfaceSlice(vals)

	if len(valsSl) == 0 {
		return ast.NewArrayNode(string(text), []ast.Node{}, offset, line, col)
	}
	res := []interface{}{valsSl[0]}
	restSl := ast.ToIfaceSlice(valsSl[1])
	for _, v := range restSl {
		vSl := ast.ToIfaceSlice(v)
		res = append(res, vSl[2])
	}
	return ast.NewArrayNode(string(text), ast.ToNodesSlice(res), offset, line, col)
}
