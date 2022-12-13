package ast

import "fmt"

type ExpressionNode struct {
	BaseNode

	OperatorType OperatorType `json:"operatorType"`

	Operator string `json:"operator"`

	Left Node `json:"left"`

	Right Node `json:"right"`
}

func NewExpressionNode(token string, operator string, operatorType OperatorType, left, right Node, offset, line, col int) ExpressionNode {
	node := ExpressionNode{
		BaseNode: BaseNode{
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

func (n ExpressionNode) String() string {
	return fmt.Sprintf("ExpressionNode %s %s %s", n.Left, n.Operator, n.Right)
}

func (n ExpressionNode) Eval(m Map) (any, error) {
	l, err := n.Left.Eval(m)
	if err != nil {
		return nil, err
	}

	r, err := n.Right.Eval(m)
	if err != nil {
		return nil, err
	}

	fmt.Println("EVAL =>", l, n.Operator, r)
	return 0, nil
}
