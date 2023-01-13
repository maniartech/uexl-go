package ast

import "github.com/maniartech/uexl_go/types"

// ArrayNode represents a node that contains an array of nodes.
// It implements the Node interface and can be evaluated.
type ArrayNode struct {
	*BaseNode

	// Value is the array of nodes.
	Items []Node `json:"value"`
}

func NewArrayNode(token string, items []Node, offset, line, col int) *ArrayNode {
	node := &ArrayNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeArray,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Items: items,
	}

	return node
}

// SetPipeType sets the pipe type of the node.
func (n *ArrayNode) SetPipeType(pipeType string) {
	n.pipeType = pipeType

	for _, item := range n.Items {
		item.SetPipeType(pipeType)
	}
}

// Eval evaluates the node.
func (n *ArrayNode) Eval(m types.Map) (any, error) {
	var array []any = make([]any, 0, len(n.Items))
	for _, item := range n.Items {
		val, err := item.Eval(m)
		if err != nil {
			return nil, err
		}

		array = append(array, val)
	}

	return array, nil
}
