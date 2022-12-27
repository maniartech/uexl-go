package ast

import (
	"fmt"
)

type PipeNode struct {
	*BaseNode
	Expressions []Node `json:"expressions"`
}

func NewPipeNode(token string, pType []string, left Node, buffer []Node, offset, line, col int) (Node, error) {
	node := PipeNode{
		BaseNode: &BaseNode{
			Type:   NodeTypePipe,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
	}

	node.Expressions = append(node.Expressions, left)
	node.Expressions = append(node.Expressions, buffer...)

	totalExprs := len(node.Expressions)
	for i := 1; i < totalExprs; i++ {
		node.Expressions[i].GetBaseNode().PipeType = pType[i-1]
	}

	return node, nil
}

func (n PipeNode) String() string {
	return fmt.Sprintf("PipeNode %v", n.Expressions)
}
func (n PipeNode) Eval(Map) (any, error) {
	return nil, nil
}
