package ast

import (
	"encoding/json"

	"github.com/maniartech/uexl_go/operators"
	"github.com/maniartech/uexl_go/types"
)

type DotExpressionNode struct {
	*BaseNode

	Expr Node `json:"expr"`

	Key string `json:"key"`
}

func NewDotExpressionNode(token string, expr Node, key string, offset, line, col int) *DotExpressionNode {
	node := &DotExpressionNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeDotExpression,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Expr: expr,
		Key:  key,
	}
	return node
}

// SetPipeType sets the pipe type for the expression node
// This is used to determine if the expression is running
// as a pipe or not. This function is called by the parser
// when it detects a pipe. It performs the recursive call
// to set the pipe type for all child nodes.
func (n *DotExpressionNode) SetPipeType(pipeType string) {
	n.pipeType = pipeType

	if n.Expr != nil {
		n.Expr.SetPipeType(pipeType)
	}

}

func (n *DotExpressionNode) Eval(ctx types.Context) (types.Value, error) {
	return operators.DotEval(n.Expr, n.Key, ctx)
}

func (n *DotExpressionNode) String() string {
	b, _ := json.MarshalIndent(n, "", "  ")
	return string(b)
}
