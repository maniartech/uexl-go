package ast

import (
	"github.com/maniartech/uexl_go/functions"
	"github.com/maniartech/uexl_go/types"
)

type FunctionNode struct {
	*BaseNode

	// Name is the name of the function
	Name string `json:"name"`

	// Params is the list of parameters for the function
	Params []Node `json:"params"`
}

func NewFunctionNode(token string, name string, params []Node, offset, line, col int) *FunctionNode {
	node := &FunctionNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeFunc,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
	}

	node.Name = name
	node.Params = params

	return node
}

// SetPipeType sets the pipe type for the expression node
// This is used to determine if the expression is running
// as a pipe or not. This function is called by the parser
// when it detects a pipe. It performs the recursive call
// to set the pipe type for all child nodes.
func (n *FunctionNode) SetPipeType(pipeType string) {
	n.pipeType = pipeType

	for _, param := range n.Params {
		param.SetPipeType(pipeType)
	}
}

func (n *FunctionNode) Eval(m types.Context) (res types.Value, err error) {
	args := make([]any, len(n.Params))
	for i, param := range n.Params {
		args[i], err = param.Eval(m)
		if err != nil {
			return nil, err
		}
	}

	res, err = functions.InvokeFunction(n.Name, args)
	return
}
