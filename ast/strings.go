package ast

import (
	"fmt"
)

type String string

type StringNode struct {
	BaseNode

	Value String `json:"value"`
}

func NewStringNode(token string, offset, line, col int) (Node, error) {
	node := StringNode{
		BaseNode: BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: String(token[1 : len(token)-1]),
	}

	return node, nil
}

func (n StringNode) String() string {
	return fmt.Sprintf("StringNode %s", n.Value)
}

func (n StringNode) Eval(Map) (any, error) {
	return n.Value, nil
}
