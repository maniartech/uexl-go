package ast

import (
	"fmt"
)

type Null struct{}

type NullNode struct {
	BaseNode

	Value Null `json:"value"`
}

func NewNullNode(token string, offset, line, col int) (Node, error) {
	node := NullNode{
		BaseNode: BaseNode{
			Type:   NodeTypeNull,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		// Value: Null,
	}

	return node, nil
}

func (n NullNode) String() string {
	return fmt.Sprintf("NullNode null")
}

func (n NullNode) Eval(Map) (any, error) {
	return n.Value, nil
}
