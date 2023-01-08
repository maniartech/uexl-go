package ast

import (
	"github.com/maniartech/uexl_go/operators"
)

type ExpressionNode struct {
	*BaseNode

	OperatorType OperatorType `json:"operatorType"`

	Operator string `json:"operator"`

	Left Node `json:"left"`

	Right Node `json:"right"`
}

func NewExpressionNode(token string, operator string, operatorType OperatorType, left, right Node, offset, line, col int) *ExpressionNode {
	node := &ExpressionNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeExpression,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Operator:     operator,
		OperatorType: operatorType,
		Left:         left,
		Right:        right,
	}
	return node
}

// SetPipeType sets the pipe type for the expression node
// This is used to determine if the expression is running
// as a pipe or not. This function is called by the parser
// when it detects a pipe. It performs the recursive call
// to set the pipe type for all child nodes.
func (n *ExpressionNode) SetPipeType(pipeType string) {
	n.PipeType = pipeType

	if n.Left != nil {
		n.Left.SetPipeType(pipeType)
	}

	if n.Right != nil {
		n.Right.SetPipeType(pipeType)
	}
}

func (n *ExpressionNode) Eval(m Map) (any, error) {
	l, err := n.Left.Eval(m)
	if err != nil {
		return nil, err
	}

	r, err := n.Right.Eval(m)
	if err != nil {
		return nil, err
	}

	return operators.Eval(n.Operator, l, r)
}
