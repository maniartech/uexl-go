package parser

import "github.com/maniartech/uexl_go/ast"

func parseArray(vals interface{}, text []byte, offset, line, col int) (ast.Node, error) {
	valsSl := toIfaceSlice(vals)

	if len(valsSl) == 0 {
		return ast.NewArrayNode(string(text), []ast.Node{}, offset, line, col), nil
	}

	res := []interface{}{valsSl[0]}
	restSl := toIfaceSlice(valsSl[1])

	for _, v := range restSl {
		vSl := toIfaceSlice(v)
		res = append(res, vSl[2])
	}

	return ast.NewArrayNode(string(text), toNodesSlice(res), offset, line, col), nil
}
