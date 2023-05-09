package parser

import (
	"github.com/maniartech/uexl_go/ast"
)

func parseExpression(token string, first, rest interface{}, offset, line, col int) (ast.Node, error) {

	l, ok := first.(ast.Node)

	// TODO: handle error
	if l != nil && !ok { // when l is nil, ignore it!
		panic("invalid-expression - assertion-failed")
	}

	restSl := toIfaceSlice(rest)

	for _, v := range restSl {
		restExpr := toIfaceSlice(v)
		r, ok := restExpr[3].(ast.Node)
		// TODO: handle error
		if !ok {
			panic("invalid-expression - assertion-failed!")
		}
		op := ""

		if o, ok := restExpr[1].([]byte); ok {
			op = string(o)
		} else if o, ok := restExpr[1].(string); ok {
			op = o
		}

		if op != "" {
			opType := ast.GetOperatorType(op)
			l = ast.NewExpressionNode(token, op, opType, l, r, offset, line, col)
		}
	}

	return l, nil
}

// func parseDotExpression(token string, first interface{}, rest interface{}, offset, line, col int) (ast.Node, error) {
// 	fmt.Println("parseDotExpression", token, first, rest)
// 	return nil, nil
// }
