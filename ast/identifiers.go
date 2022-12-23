package ast

import "fmt"

type IdentifierNode struct {
	BaseNode
	Name string `json:"name"`
}

func NewIdentifierNode(token string, offset, line, col int) (Node, error) {
	node := IdentifierNode{
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

func (n IdentifierNode) String() string {
	return fmt.Sprintf("IdentifierNode %v", n.Name)
}

func (n IdentifierNode) Eval(Map) (any, error) {
	return n.Name, nil
}
