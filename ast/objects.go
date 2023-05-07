package ast

import "github.com/maniartech/uexl_go/types"

type ObjectItem struct {
	Key   types.String
	Value Node
}

type ObjectNode struct {
	*BaseNode

	Items []ObjectItem
}

func NewObjectNode(token string, items []ObjectItem, offset, line, col int) *ObjectNode {
	node := &ObjectNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeObject,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Items: items,
	}

	return node
}

// Eval parses the array node and returns the array node.
func (n *ObjectNode) Eval(env types.Context) (types.Value, error) {
	panic("implement me")
}
