package ast

type Null struct{}

type NullNode struct {
	*BaseNode
}

func NewNullNode(token string, offset, line, col int) *NullNode {
	node := &NullNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeNull,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
	}

	return node
}

func (n NullNode) Eval(Map) (any, error) {
	return nil, nil
}
