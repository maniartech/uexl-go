package compiler

import (
	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

type Compiler struct {
	constants  []parser.Node
	scopes     []CompilationScope
	scopeIndex int
}

type EmmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmmittedInstruction
	previousInstruction EmmittedInstruction
}

func (c *Compiler) Compile(node parser.Node) error {
	switch node := node.(type) {
	case *parser.BinaryExpression:
		left := node.Left
		operator := node.Operator
		right := node.Right

		// Compile the left operand
		if err := c.Compile(left); err != nil {
			return err
		}
		
		// Compile the right operand
		if err := c.Compile(right); err != nil {
			return err
		}

		// Emit the appropriate instruction for the operator
		switch operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		}


	case *parser.NumberLiteral:
		// Add the number literal to constants
		number := &parser.NumberLiteral{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(number))
	}
	return nil
}
