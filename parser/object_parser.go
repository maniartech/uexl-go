package parser

import "github.com/maniartech/uexl_go/ast"

func parseObject(vals interface{}, text []byte, offset, line, col int) (ast.Node, error) {
	valsSl := toIfaceSlice(vals)
	res := make([]ast.ObjectItem, 0, len(valsSl))
	if len(valsSl) == 0 {
		return ast.NewObjectNode(string(text), res, offset, line, col), nil
	}

	res = append(res, ast.ObjectItem{Key: valsSl[0].(*ast.StringNode).Value, Value: valsSl[4].(ast.Node)})
	restSl := toIfaceSlice(valsSl[5])

	for _, v := range restSl {
		vSl := toIfaceSlice(v)
		res = append(res, ast.ObjectItem{Key: vSl[2].(ast.StringNode).Value, Value: vSl[6].(ast.Node)})
	}

	return ast.NewObjectNode(string(text), res, offset, line, col), nil
}
