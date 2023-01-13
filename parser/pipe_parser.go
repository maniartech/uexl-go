package parser

import "github.com/maniartech/uexl_go/ast"

func parsePipe(pipeType, expr interface{}) ast.Node {
	node := expr.(ast.Node)
	pType := "pipe"
	if pipeType != nil {
		pType = pipeType.(string)
	}

	node.SetPipeType(pType)
	return node
}
