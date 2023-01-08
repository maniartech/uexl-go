package ast

import "github.com/maniartech/uexl_go/types"

type StringNode struct {
	*BaseNode

	Value types.String `json:"value"`
}

func NewStringNode(token string, offset, line, col int) *StringNode {
	finalToken := token
	if finalToken[0] == '\'' && finalToken[len(finalToken)-1] == '\'' {
		finalToken = finalToken[1 : len(token)-1]
	}
	node := &StringNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: types.String(finalToken),
	}

	return node
}

func (n StringNode) Eval(Map) (any, error) {
	return n.Value, nil
}
