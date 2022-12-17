package ast

type PipeNode struct {
	BaseNode

	Expressions ExpressionNode `json:"expressions"`
}

func NewPipeNode(token, pipeType string, offset, line, col int) (Node, error) {
	node := PipeNode{
		BaseNode: BaseNode{
			Type:   NodeTypePipe,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Expressions: ExpressionNode{
			PipeType: pipeType,
		},
	}

	return node, nil
}

func (n PipeNode) Eval(Map) (any, error) {
	return n.Expressions, nil
}
