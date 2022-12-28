package parser

import (
	"github.com/maniartech/uexl_go/ast"
	"github.com/maniartech/uexl_go/ast/constants"
)

func parseInput(token string, firstNode, restNodes interface{}, offset, line, col int) (interface{}, error) {
	pipedNodes := toIfaceSlice(restNodes)
	if len(pipedNodes) == 0 {
		return firstNode, nil
	}

	var nodes []ast.Node = make([]ast.Node, 0, len(pipedNodes)+1)

	firstNode.(ast.Node).SetPipeType(constants.PipeTypeFirst)
	nodes = append(nodes, firstNode.(ast.Node))

	for _, v := range pipedNodes {
		nodes = append(nodes, v.(ast.Node))
	}

	return ast.NewPipeNode(token, nodes, offset, line, col), nil
}
