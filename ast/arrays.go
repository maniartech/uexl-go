package ast

import "fmt"

type Array []Node

type ArrayNode struct {
	*BaseNode

	Value Array `json:"value"`
}

func NewArrayNode(token string, items []Node, offset, line, col int) (Node, error) {
	node := ArrayNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeArray,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Value: items,
	}

	return node, nil
}

func (n ArrayNode) String() string {
	return fmt.Sprintf("ArrayNode %v", n.Value)
}

func (n ArrayNode) Eval(Map) (any, error) {
	return n.Value, nil
}
