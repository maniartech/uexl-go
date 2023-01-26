package ast

import (
	"strconv"

	"github.com/maniartech/uexl_go/types"
)

type StringNode struct {
	*BaseNode

	Value types.String `json:"value"`
}

func NewStringNode(token string, offset, line, col int) (*StringNode, error) {
	// finalToken := token

	// if finalToken[0] == '\'' && finalToken[len(finalToken)-1] == '\'' {
	// 	finalToken = finalToken[1 : len(token)-1]
	// }

	value, err := strconv.Unquote(token)
	if err != nil {
		return nil, err
	}

	node := &StringNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: types.String(value),
	}

	return node, nil
}

func (n StringNode) Eval(types.Map) (any, error) {
	return n.Value, nil
}
