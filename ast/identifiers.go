package ast

import "fmt"

type IdentifierNode struct {
	BaseNode
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewIdentifierNode(token, name, value string, offset, line, col int) (Node, error) {
	node := IdentifierNode{
		BaseNode: BaseNode{
			Type:   NodeTypeIdentifier,
			Line:   line,
			Offset: offset,
			Column: col,
			Token:  token,
		},
		Name:  name,
		Value: value,
	}

	return node, nil
}

func (n IdentifierNode) String() string {
	return fmt.Sprintf("IdentifierNode %v", n.Value)
}

func (n IdentifierNode) Eval(Map) (any, error) {
	return n.Value, nil
}
