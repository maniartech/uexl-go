package ast

import "fmt"

type Identifier string

type IdentifierNode struct {
	BaseNode

	Value Identifier `json:"value"`
}

func NewIdentifierNode(token string, offset, line, col int) (Node, error) {
	node := IdentifierNode{
		BaseNode: BaseNode{
			Type:   NodeTypeIdentifier,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Value: Identifier(token),
	}

	return node, nil

}

func (n IdentifierNode) String() string {
	return fmt.Sprintf("IdentifierNode %s", n.Value)
}

func (n IdentifierNode) Eval(Map) (any, error) {
	return n.Value, nil
}
