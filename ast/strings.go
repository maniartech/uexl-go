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
	finalToken := token
	if finalToken[0] == '\'' && finalToken[len(finalToken)-1] == '\'' {
		finalToken = finalToken[1 : len(token)-1]
	}
	node := StringNode{
		BaseNode: BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: String(finalToken),
	}

	return node, nil
}

func (n StringNode) String() string {
	return fmt.Sprintf("StringNode %s", n.Value)
}

func (n StringNode) Eval(Map) (any, error) {
	return n.Value, nil
}
