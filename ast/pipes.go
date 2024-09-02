package ast

import (
	"fmt"

	"github.com/maniartech/uexl_go/ast/constants"
	"github.com/maniartech/uexl_go/pipes"
	"github.com/maniartech/uexl_go/types"
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
			Token:  token,
		},
	}

	node.pipeType = constants.PipeTypeRoot
	node.Nodes = nodes

	return node
}

// Eval evaluates the node and returns the result.
func (n *PipeNode) Eval(ctx types.Context) (result types.Value, err error) {

	// copy the context into ctx
	ctx = ctx.ShallowCopy()

	for _, node := range n.Nodes {

		// fmt.Println(
		// 	"PipeNode.Eval() node.PipeType():", node.PipeType(),
		// 	"node.GetType():", node.GetType(),
		// )

		handler, ok := pipes.Get(node.PipeType())
		if !ok {
			return nil, fmt.Errorf("pipe %s not found", node.PipeType())
		}

		result, err = handler(node, ctx, result)
		if err != nil {
			return nil, err
		}
	}

	// Return the result
	return result, nil
}
