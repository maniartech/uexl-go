package compiler

import (
	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

type Compiler struct {
	constants   []parser.Node
	contextVars []parser.Node
	scopes      []CompilationScope
	scopeIndex  int
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

		if node.Operator == "<" || node.Operator == "<=" {
			// swap operand sides for less than comparisons because they are right associative
			left, right = right, left
		}

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
		case "%":
			c.emit(code.OpMod)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">", "<":
			c.emit(code.OpGreaterThan)
		case ">=", "<=":
			c.emit(code.OpGreaterThanOrEqual)
		case "&&":
			c.emit(code.OpLogicalAnd)
		case "||":
			c.emit(code.OpLogicalOr)
		case "&":
			c.emit(code.OpBitwiseAnd)
		case "|":
			c.emit(code.OpBitwiseOr)
		case "^":
			c.emit(code.OpBitwiseXor)
		case "<<":
			c.emit(code.OpShiftLeft)
		case ">>":
			c.emit(code.OpShiftRight)

		}
	case *parser.UnaryExpression:
		err := c.Compile(node.Operand)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		case "~":
			c.emit(code.OpBitwiseNot)
		}
	case *parser.NumberLiteral:
		// Add the number literal to constants
		number := &parser.NumberLiteral{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(number))
	case *parser.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *parser.StringLiteral:
		// Add the string literal to constants
		stringLiteral := &parser.StringLiteral{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(stringLiteral))
	case *parser.Identifier:
		// Identifiers are variables passed via go's environment context. They are "Constant" in a sense that they are not computed at runtime.
		c.emit(code.OpContextVar, c.addContextVar(node))
	}
	return nil
}
