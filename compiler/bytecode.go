package compiler

import (
	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/types"
)

type ByteCode struct {
	Instructions code.Instructions
	Constants    []types.Value
	ContextVars  []string
	SystemVars   []any
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
		ContextVars:  c.contextVars,
		SystemVars:   c.SystemVars,
	}
}
