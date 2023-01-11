package ast

import "github.com/maniartech/uexl_go/types"

type Node interface {
	Eval(types.Map) (any, error)

	GetBaseNode() *BaseNode

	GetType() NodeType

	SetPipeType(string)
}
