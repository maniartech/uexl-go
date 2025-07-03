package parser

import (
	"github.com/maniartech/uexl_go/ast"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return []interface{}{} // Return empty slice instead of nil
	}

	return v.([]interface{})
}

// toNodesSlice converts an interface{} to a slice of ast.Node.
// If the input is nil, it returns nil.
// Otherwise, it assumes that the input is a slice of interfaces and iterates over it,
// converting each element to an ast.Node and storing it in a new slice.
// The resulting slice of ast.Node is then returned.
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
