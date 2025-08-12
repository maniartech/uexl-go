package compiler

import (
	"sort"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

type Compiler struct {
	constants   []parser.Node
	contextVars []parser.Node
	SystemVars  []parser.Node
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

type InstructionBlock struct {
	parser.Node
	Instructions code.Instructions
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
	case *parser.FunctionCall:
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		fnIdx := c.addConstant(&parser.StringLiteral{Value: node.Function.(*parser.Identifier).Name})
		c.emit(code.OpCallFunction, fnIdx, len(node.Arguments))
	case *parser.ProgramNode:
		// First expression is the entry point will just be normal expression from which we get the result
		// We will just compile it like a normal expression
		if len(node.PipeExpressions) > 0 {
			if err := c.Compile(node.PipeExpressions[0].Expression); err != nil {
				return err
			}
			if node.PipeExpressions[0].Alias != "" {
				// If the first pipe has an alias, store it in the context
				aliasVarIdx := c.addPipeLocalVar(node.PipeExpressions[0].Alias)
				c.emit(code.OpStore, aliasVarIdx)
			}
		}
		// Compile each pipe expression
		for _, pipeExpr := range node.PipeExpressions[1:] {
			// Compile the pipe's predicate expression block
			pipeTypeIdx := c.addConstant(&parser.StringLiteral{Value: pipeExpr.PipeType})
			aliasIdx := c.addPipeLocalVar(pipeExpr.Alias)
			blockIdx, err := c.compilePredicateBlock(pipeExpr.Expression)
			if err != nil {
				return err
			}
			c.emit(code.OpPipe, pipeTypeIdx, aliasIdx, blockIdx)

		}

	case *parser.IndexAccess:
		if err := c.Compile(node.Array); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(code.OpIndex)

	case *parser.ObjectLiteral:
		// Ensure deterministic order by sorting keys
		keys := make([]string, 0, len(node.Properties))
		for key := range node.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			stringLiteral := &parser.StringLiteral{Value: key}
			keyIdx := c.addConstant(stringLiteral)
			c.emit(code.OpConstant, keyIdx) // Push key onto stack

			if err := c.Compile(node.Properties[key]); err != nil {
				return err
			}
		}
		c.emit(code.OpObject, len(node.Properties)*2) // Each key-value pair is two stack elements
	case *parser.NumberLiteral:
		// Add the number literal to constants
		c.emit(code.OpConstant, c.addConstant(node))
	case *parser.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *parser.StringLiteral:
		// Add the string literal to constants
		c.emit(code.OpConstant, c.addConstant(node))
	case *parser.Identifier:
		// Identifiers are variables passed via go's environment context. They are "Constant" in a sense that they are not computed at runtime.
		// If identifer begins with a dollar sign, it is a local variable in the pipe context.
		if isPipeLocalVar(node.Name) {
			c.emit(code.OpIdentifier, c.addPipeLocalVar(node.Name))
		} else {
			c.emit(code.OpContextVar, c.addContextVar(node))
		}
	case *parser.ArrayLiteral:
		// Compile each element in the array
		for _, element := range node.Elements {
			if err := c.Compile(element); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	}
	return nil
}
