package ast

import "fmt"

type ExpressionNode struct {
	BaseNode

	OperatorType OperatorType `json:"operatorType"`

	Operator string `json:"operator"`

	Left Node `json:"left"`

	Right Node `json:"right"`
}

func NewExpressionNode(operator string, operatorType OperatorType, left, right Node, offset, line, col int) ExpressionNode {
	node := ExpressionNode{
		BaseNode: BaseNode{
			Type:   NodeTypeExpression,
			Line:   line,
			Column: col,
			Offset: offset,
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

func parseExpression(first, rest interface{}, offset, line, col int) (ExpressionNode, error) {

	// fmt.Printf("EVAL => First:%+v Rest:%+v\n", first, rest)
	// return nil

	l, ok := first.(Node)

	// TODO: handle error
	if l != nil && !ok { // when l is nil, ignore it!
		panic("invalid-expression - assertion-failed")
	}

	// restSl := ToNodesSlice(rest)

	// for _, v := range restSl {
	// restExpr := ToNodesSlice(v)
	// r, ok := restExpr[3].(Node)
	// // TODO: handle error
	// if !ok {
	// 	panic("invalid-expression - assertion-failed!")
	// }
	// op := ""

	// // if o, ok := restExpr[1].([]byte); ok {
	// // 	op = string(o)
	// // } else if o, ok := restExpr[1].(string); ok {
	// // 	op = o
	// // }

	// if op != "" {
	// 	opType := getOperatorType(op)
	// 	l = NewExpressionNode(op, opType, l, r, offset, line, col)
	// }
	// }

	return l.(ExpressionNode), nil
}
