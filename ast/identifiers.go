package ast

type IdentifierNode struct {
	BaseNode

	Name string `json:"name"`
}

func NewIdentifierNode(token string, offset, line, col int) (*IdentifierNode, error) {
	node := &IdentifierNode{
		BaseNode: BaseNode{
			Type:   NodeTypeIdentifier,
			Line:   line,
			Offset: offset,
			Column: col,
			Token:  token,
		},
		Name: token,
	}

	return node, nil
}

func (n IdentifierNode) Eval(Map) (any, error) {
	return n.Name, nil
}
