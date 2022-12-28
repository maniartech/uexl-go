package ast

type Node interface {
	Eval(Map) (any, error)

	GetBaseNode() *BaseNode

	GetType() NodeType

	SetPipeType(string)
}
