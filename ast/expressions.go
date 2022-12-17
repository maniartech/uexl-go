package ast

import "fmt"

type ExpressionNode struct {
	BaseNode

	OperatorType OperatorType `json:"operatorType"`

	Operator string `json:"operator"`

	Left Node `json:"left"`

	Right Node `json:"right"`

	PipeType string `json:"pipeType"`
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

func ParseExpression(token string, first, rest interface{}, offset, line, col int) (Node, error) {

	l, ok := first.(Node)

	// TODO: handle error
	if l != nil && !ok { // when l is nil, ignore it!
		panic("invalid-expression - assertion-failed")
	}

	restSl := ToIfaceSlice(rest)

	for _, v := range restSl {
		restExpr := ToIfaceSlice(v)
		r, ok := restExpr[3].(Node)
		// TODO: handle error
		if !ok {
			panic("invalid-expression - assertion-failed!")
		}
		op := ""

		if o, ok := restExpr[1].([]byte); ok {
			op = string(o)
		} else if o, ok := restExpr[1].(string); ok {
			op = o
		}

		if op != "" {
			opType := getOperatorType(op)
			l = NewExpressionNode(token, op, opType, l, r, offset, line, col)
		}
	}

	// PrintNode(l)

	return l, nil
}
