package ast

import (
	"fmt"
	"strconv"
)

// Number represents a number literal.
type Number float64

// String returns a string representation of the number.
func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// NumberNode represent a number literal.
type NumberNode struct {
	BaseNode

	// Value is the value set to the NumberNode.
	Value Number `json:"value"`
}

// NewNumberNode creates a new NumberNode.
func NewNumberNode(value Number, offset, line, col int) NumberNode {
	node := NumberNode{
		BaseNode: BaseNode{
			Type:   NodeTypeNumber,
			Line:   line,
			Column: col,
			Offset: offset,
		},
		Value: value,
	}
	return node
}

// String returns a string representation of the NumberNode.
func (n NumberNode) String() string {
	return fmt.Sprintf("NumberNode %s", n.Value)
}

// Eval evalues the NumberNode.
func (n NumberNode) Eval(Map) (any, error) {
	return n.Value, nil
}
