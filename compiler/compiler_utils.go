package compiler

import (
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

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	old := c.currentInstructions()
	new := old[:last.Position]
	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()
	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
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
		constants:   []parser.Node{},
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		contextVars: []parser.Node{},
	}
}

func NewWithState(constants []parser.Node) *Compiler {
	compiler := New()
	compiler.constants = constants
	return compiler
}

func (c *Compiler) addConstant(node parser.Node) int {
	c.constants = append(c.constants, node)
	return len(c.constants) - 1
}

func (c *Compiler) addContextVar(node parser.Node) int {
	ident, ok := node.(*parser.Identifier)
	if !ok {
		panic("addContextVar: node is not *parser.Identifier")
	}
	for i, existing := range c.contextVars {
		if exIdent, ok := existing.(*parser.Identifier); ok && exIdent.Name == ident.Name {
			return i // Return the index of the existing variable by name
		}
	}
	c.contextVars = append(c.contextVars, node)
	return len(c.contextVars) - 1
}

func (c *Compiler) addArray(node *parser.ArrayLiteral) int {
	c.constants = append(c.constants, node)
	return len(c.constants) - 1
}
