package compiler

import (
	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

type ByteCode struct {
	Instructions code.Instructions
	Constants    []parser.Node
	ContextVars  []parser.Node
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}
