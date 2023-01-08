package uexl_go

import (
	"github.com/maniartech/uexl_go/parser"
)

func EvalExpr(expr string) (any, error) {
	node, err := parser.ParseString(expr)
	if err != nil {
		return nil, err
	}

	return node.Eval(nil)
}
