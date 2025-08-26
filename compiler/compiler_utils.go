package compiler

import (
	"encoding/binary"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)
	c.scopes[c.scopeIndex].instructions = updatedInstructions
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(opcode code.Opcode, position int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmmittedInstruction{Opcode: opcode, Position: position}
	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	pos := c.addInstruction(instruction)
	c.setLastInstruction(op, pos)
	return pos
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmmittedInstruction{},
		previousInstruction: EmmittedInstruction{},
	}
	return &Compiler{
		constants:   []any{},
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		contextVars: []any{},
	}
}

func NewWithState(constants []any) *Compiler {
	compiler := New()
	compiler.constants = constants
	return compiler
}

func (c *Compiler) addConstant(node any) int {
	c.constants = append(c.constants, node)
	return len(c.constants) - 1
}

func (c *Compiler) addContextVar(node any) int {
	for i, existing := range c.contextVars {
		if existing == node {
			return i // Return the index of the existing variable by name
		}
	}
	c.contextVars = append(c.contextVars, node)
	return len(c.contextVars) - 1
}

func (c *Compiler) enterScope() {
	c.scopes = append(c.scopes, CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmmittedInstruction{},
		previousInstruction: EmmittedInstruction{},
	})
	c.scopeIndex++
}

func (c *Compiler) exitScope() {
	if c.scopeIndex == 0 {
		panic("exitScope: already at the global scope")
	}
	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--
}

func (c *Compiler) compilePredicateBlock(expr parser.Node) (int, error) {
	if expr == nil {
		return c.addConstant(&InstructionBlock{Instructions: nil}), nil
	}
	c.enterScope()
	err := c.Compile(expr)
	if err != nil {
		return 0, err
	}
	blockIns := c.ByteCode().Instructions

	c.exitScope()
	return c.addConstant(&InstructionBlock{Instructions: blockIns}), nil
}

func (c *Compiler) addPipeLocalVar(name string) int {
	c.SystemVars = append(c.SystemVars, name)
	return len(c.SystemVars) - 1
}

func isPipeLocalVar(name string) bool {
	return len(name) > 0 && name[0] == '$'
}

func (c *Compiler) replaceOperand(pos int, value int) {
	binary.BigEndian.PutUint16(c.currentInstructions()[pos:], uint16(value))
}

// flattenLogicalChain collects a left-associative chain of the same logical operator into terms.
// Example: a || b || c  => [a, b, c]
func flattenLogicalChain(n parser.Node, op string, out *[]parser.Node) {
	be, ok := n.(*parser.BinaryExpression)
	if !ok || be.Operator != op {
		*out = append(*out, n)
		return
	}
	flattenLogicalChain(be.Left, op, out)
	flattenLogicalChain(be.Right, op, out)
}

func flattenAccessChain(n parser.Node) (base parser.Node, steps []accessStep) {
	var rev []accessStep
	cur := n
	for {
		switch v := cur.(type) {
		case *parser.MemberAccess:
			rev = append(rev, accessStep{
				isMember:   true,
				property:   v.Property, // assume already a string
				isOptional: v.Optional,
			})
			cur = v.Target
		case *parser.IndexAccess:
			rev = append(rev, accessStep{
				isMember:   false,
				property:   v.Index, // parser.Node to be compiled later
				isOptional: v.Optional,
			})
			cur = v.Target
		default:
			base = cur
			// reverse order
			for i := len(rev) - 1; i >= 0; i-- {
				steps = append(steps, rev[i])
			}
			return
		}
	}
}
