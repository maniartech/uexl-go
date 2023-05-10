package ast

import (
	"fmt"

	"github.com/maniartech/uexl_go/types"
)

type IdentifierNode struct {
	BaseNode

	Name string `json:"name"`
}

func NewIdentifierNode(token string, offset, line, col int) (*IdentifierNode, error) {
	node := &IdentifierNode{
		BaseNode: BaseNode{
			Type:   NodeTypeIdentifier,
			Line:   line,
			Offset: offset,
			Column: col,
			Token:  token,
		},
		Name: token,
	}

	return node, nil
}

func (n IdentifierNode) Eval(m types.Context) (types.Value, error) {
	// If m is a nil, return an error
	if m == nil {
		return nil, fmt.Errorf("cannot access identifier '%s' from nil context", n.Name)
	}

	if val, ok := m[n.Name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key '%s' not found", n.Name)
}
