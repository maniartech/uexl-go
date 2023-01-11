package ast

import (
	"strconv"

	"github.com/maniartech/uexl_go/types"
)

// NumberNode represent a number literal.
type NumberNode struct {
	*BaseNode

	// Value is the value set to the NumberNode.
	Value types.Number `json:"value"`
}

// NewNumberNode creates a new NumberNode.
func NewNumberNode(token string, offset, line, col int) (*NumberNode, error) {
	// Convert the token to a float64 using fasted method.
	//
	// This is the fastest method to convert a string to a float64.

	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, err
	}

	node := &NumberNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeNumber,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: types.Number(value),
	}
	return node, nil
}

// Eval evalues the NumberNode.
func (n NumberNode) Eval(types.Map) (any, error) {
	return n.Value, nil
}
