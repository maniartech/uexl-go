package ast

import (
	"github.com/maniartech/uexl_go/ast/constants"
)

// PipeNode represents a Node that is a pipe.
type PipeNode struct {
	*BaseNode

	// Nodes is the list of nodes that are piped.
	Nodes []Node `json:"nodes"`
}

// NewPipeNode creates a new PipeNode.
func NewPipeNode(token string, nodes []Node, offset, line, col int) *PipeNode {
	node := &PipeNode{
		BaseNode: &BaseNode{
			Type:   NodeTypePipe,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
	}

	node.PipeType = constants.PipeTypeRoot
	node.Nodes = nodes

	return node
}

// Eval evaluates the node and returns the result.
func (n *PipeNode) Eval(Map) (interface{}, error) {
	panic("implement me")
}
