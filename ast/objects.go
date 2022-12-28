package ast

type ObjectItem struct {
	Key   String
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
func (n *ObjectNode) Eval(env Map) (any, error) {
	panic("implement me")
}
