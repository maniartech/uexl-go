package ast

import "github.com/maniartech/uexl_go/evaluators"

type Node interface {
	evaluators.Evaluator

	GetBaseNode() *BaseNode

	GetType() NodeType

	PipeType() string

	SetPipeType(string)
}
