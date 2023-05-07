package ast

import (
	"strconv"

	"github.com/maniartech/uexl_go/types"
)

// BooleanNode implements the Node interface and represents a boolean value.
type BooleanNode struct {
	*BaseNode

	// Value is the boolean value.
	Value bool `json:"value"`
}

// NewBooleanNode creates a new BooleanNode.
func NewBooleanNode(token string, offset, line, col int) (*BooleanNode, error) {
	value, err := strconv.ParseBool(token)
	if err != nil {
		return nil, err
	}

	node := &BooleanNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeBoolean,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: value,
	}

	return node, nil
}

// Eval evaluates the BooleanNode and returns the value.
func (n *BooleanNode) Eval(types.Context) (types.Value, error) {
	return types.Boolean(n.Value), nil
}
