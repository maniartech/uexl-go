package compiler

import (
	"fmt"
	"sort"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

type Compiler struct {
	constants   []any
	contextVars []string
	SystemVars  []any
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
	Instructions code.Instructions
}

type accessStep struct {
	safe         bool
	propertyStr  string      // member name
	propertyExpr parser.Node // index expression
}

func (c *Compiler) compileShortCircuitChain(terms []parser.Node, jumpOpcode code.Opcode) error {
	jumpPositions := make([]int, 0, len(terms))
	for i, term := range terms {
		if err := c.Compile(term); err != nil {
			return err
		}
		if i < len(terms)-1 {
			pos := len(c.currentInstructions())
			c.emit(jumpOpcode, 0) // placeholder to END
			jumpPositions = append(jumpPositions, pos)
		}
	}
	end := len(c.currentInstructions())
	for _, p := range jumpPositions {
		c.replaceOperand(p+1, end)
	}
	return nil
}

func (c *Compiler) Compile(node parser.Node) error {
	switch node := node.(type) {
	case *parser.BinaryExpression:
		left := node.Left
		operator := node.Operator
		right := node.Right

		// Handle right-associative operators and operand swapping
		switch operator {
		case "<", "<=":
			left, right = right, left
		}

		// Short-circuit logical operators (flattened chains + single backpatch)
		if operator == "||" || operator == "&&" || operator == "??" {
			terms := make([]parser.Node, 0, 4)
			flattenLogicalChain(node, operator, &terms)
			switch operator {
			case "||":
				return c.compileShortCircuitChain(terms, code.OpJumpIfTruthy)
			case "&&":
				return c.compileShortCircuitChain(terms, code.OpJumpIfFalsy)
			case "??":
				return c.compileNullishChain(terms)
			}
			return nil
		}

		// Compile operands for other operators
		if err := c.Compile(left); err != nil {
			return err
		}
		if err := c.Compile(right); err != nil {
			return err
		}

		// Emit instruction based on operator
		switch operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "**":
			c.emit(code.OpPow)
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
		default:
			return fmt.Errorf("unsupported binary operator: %s", operator)
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
	case *parser.GroupedExpression:
		// GroupedExpression is just a wrapper - compile the inner expression
		return c.Compile(node.Expression)
	case *parser.FunctionCall:
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		fnIdx := c.addConstant(node.Function.(*parser.Identifier).Name)
		c.emit(code.OpCallFunction, fnIdx, len(node.Arguments))
	case *parser.ConditionalExpression:
		// condition ? consequent : alternate
		// Compile condition
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		// Jump to else (alternate) if condition is falsy
		jumpIfFalsyPos := len(c.currentInstructions())
		c.emit(code.OpJumpIfFalsy, 0) // placeholder
		// Compile consequent
		if err := c.Compile(node.Consequent); err != nil {
			return err
		}
		// After consequent, jump to end
		jumpToEndPos := len(c.currentInstructions())
		c.emit(code.OpJump, 0) // placeholder
		// Patch first jump to here (start of alternate)
		elsePos := len(c.currentInstructions())
		c.replaceOperand(jumpIfFalsyPos+1, elsePos)
		// Remove the re-pushed falsy condition value
		c.emit(code.OpPop)
		// Compile alternate
		if err := c.Compile(node.Alternate); err != nil {
			return err
		}
		// Patch end jump
		endPos := len(c.currentInstructions())
		c.replaceOperand(jumpToEndPos+1, endPos)
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
			pipeTypeIdx := c.addConstant(pipeExpr.PipeType)
			aliasIdx := c.addPipeLocalVar(pipeExpr.Alias)
			blockIdx, err := c.compilePredicateBlock(pipeExpr.Expression)
			if err != nil {
				return err
			}
			c.emit(code.OpPipe, pipeTypeIdx, aliasIdx, blockIdx)

		}
	case *parser.MemberAccess, *parser.IndexAccess:
		return c.compileAccessNode(node, false)
	case *parser.ObjectLiteral:
		// Ensure deterministic order by sorting keys
		keys := make([]string, 0, len(node.Properties))
		for key := range node.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			keyIdx := c.addConstant(key)
			c.emit(code.OpConstant, keyIdx) // Push key onto stack

			if err := c.Compile(node.Properties[key]); err != nil {
				return err
			}
		}
		c.emit(code.OpObject, len(node.Properties)*2) // Each key-value pair is two stack elements
	case *parser.NumberLiteral:
		// Add the number literal to constants
		c.emit(code.OpConstant, c.addConstant(node.Value))
	case *parser.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *parser.StringLiteral:
		// Add the string literal to constants
		c.emit(code.OpConstant, c.addConstant(node.Value))
	case *parser.NullLiteral:
		c.emit(code.OpNull)
	case *parser.Identifier:
		// Identifiers are variables passed via go's environment context. They are "Constant" in a sense that they are not computed at runtime.
		// If identifer begins with a dollar sign, it is a local variable in the pipe context.
		if isPipeLocalVar(node.Name) {
			c.emit(code.OpIdentifier, c.addPipeLocalVar(node.Name))
		} else {
			c.emit(code.OpContextVar, c.addContextVar(node.Name))
		}
	case *parser.ArrayLiteral:
		// Compile each element in the array
		for _, element := range node.Elements {
			if err := c.Compile(element); err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *parser.SliceExpression:
		if err := c.Compile(node.Target); err != nil {
			return err
		}
		if node.Start != nil {
			if err := c.Compile(node.Start); err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}
		if node.End != nil {
			if err := c.Compile(node.End); err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}
		if node.Step != nil {
			if err := c.Compile(node.Step); err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
		}

		optional := 0
		if node.Optional {
			optional = 1
		}
		c.emit(code.OpSlice, optional)
	}
	return nil
}
