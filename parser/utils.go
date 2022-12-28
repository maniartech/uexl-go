package parser

import (
	"github.com/maniartech/uexl_go/ast"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	return v.([]interface{})
}

func toNodesSlice(v interface{}) []ast.Node {
	if v == nil {
		return nil
	}

	islice := v.([]interface{})

	// convert iSlice to []Node
	nodes := make([]ast.Node, len(islice))
	for i, node := range islice {
		nodes[i] = node.(ast.Node)
	}

	return nodes
}
