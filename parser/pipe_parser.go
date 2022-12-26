package parser

import "github.com/maniartech/uexl_go/ast"

func parsePipe(expr, pipeType interface{}) {
	nodeBuffer = append(nodeBuffer, expr.(ast.Node))
	pType = "pipe"
	if len(pipeType.([]interface{})) != 0 {
		pType = resolveAscii(pipeType)
	}
	pTypeBuffer = append(pTypeBuffer, pType)
}
