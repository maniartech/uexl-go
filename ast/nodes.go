package ast

import "github.com/maniartech/uexl_go/core"

type Node interface {
	core.Evaluator

	GetBaseNode() *BaseNode

	GetType() NodeType

	PipeType() string

	SetPipeType(string)
}
