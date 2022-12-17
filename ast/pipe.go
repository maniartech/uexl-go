package ast

import (
	"fmt"
)

type PipeNode struct {
	BaseNode
	PipeType    string `json:"pipeType"`
	Expressions []Node `json:"expressions"`
}

func NewPipeNode(token, pType string, left, right Node, offset, line, col int) (Node, error) {
	node := PipeNode{
		BaseNode: BaseNode{
			Type:   NodeTypePipe,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		PipeType: pType,
	}

	node.Expressions = append(node.Expressions, left)
	node.Expressions = append(node.Expressions, right)

	return node, nil
}

func (n PipeNode) String() string {
	return fmt.Sprintf("PipeNode %v", n.PipeType)
}
func (n PipeNode) Eval(Map) (any, error) {
	return nil, nil
}
