package compiler

import (
	"encoding/binary"
	"fmt"

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
		contextVars: []string{},
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

func (c *Compiler) addContextVar(node string) int {
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

// compileAccessNode compiles a member/index access chain.
// If softenLast is true, only the final access operation (member/index) is executed in safe mode.
func (c *Compiler) compileAccessNode(n parser.Node, softenLast bool) error {
	base, steps := flattenAccessChain(n)
	// Compile base once
	if err := c.Compile(base); err != nil {
		return err
	}

	// Optional chaining: collect jump placeholders; each jump skips to end of chain
	var jumpPositions []int
	for i, step := range steps {
		if step.isOptional {
			pos := len(c.currentInstructions())
			c.emit(code.OpJumpIfNullish, 0) // placeholder to END
			jumpPositions = append(jumpPositions, pos)
		}

		isLast := i == len(steps)-1
		if step.isMember {
			// Push property name constant
			propIdx := c.addConstant(step.property)
			if softenLast && isLast {
				// If base is nullish, error (earlier link strict). Otherwise, soften only the final access.
				// Layout:
				// [base]
				// JumpIfNullish -> ERR
				//   push prop; SafeOn; MemberAccess; SafeOff; Jump -> END
				// ERR:
				//   push prop; MemberAccess (no safe)  // will error on nil base
				// END:
				// Check for nullish without consuming stack
				jErrPos := len(c.currentInstructions())
				c.emit(code.OpJumpIfNullish, 0) // -> ERR

				// Normal softened path
				c.emit(code.OpConstant, propIdx)
				c.emit(code.OpSafeModeOn)
				c.emit(code.OpMemberAccess)
				c.emit(code.OpSafeModeOff)
				jEndPos := len(c.currentInstructions())
				c.emit(code.OpJump, 0) // -> END

				// ERR label
				errPos := len(c.currentInstructions())
				c.replaceOperand(jErrPos+1, errPos)
				c.emit(code.OpConstant, propIdx)
				c.emit(code.OpMemberAccess) // no safe: raises if base is nil

				// END label
				endPos := len(c.currentInstructions())
				c.replaceOperand(jEndPos+1, endPos)
			} else {
				c.emit(code.OpConstant, propIdx)
				c.emit(code.OpMemberAccess)
			}
		} else {
			// Compile index expression
			idxExpr, ok := step.property.(parser.Node)
			if !ok {
				return fmt.Errorf("invalid index expression")
			}
			if softenLast && isLast {
				// Similar to member: guard base nullish to trigger hard error; else soften index access
				jErrPos := len(c.currentInstructions())
				c.emit(code.OpJumpIfNullish, 0) // -> ERR

				// Normal softened path
				if err := c.Compile(idxExpr); err != nil {
					return err
				}
				c.emit(code.OpSafeModeOn)
				c.emit(code.OpIndex, 1) // optional = true
				c.emit(code.OpSafeModeOff)
				jEndPos := len(c.currentInstructions())
				c.emit(code.OpJump, 0) // -> END

				// ERR label: compile index again and execute without safe to raise meaningful error
				errPos := len(c.currentInstructions())
				c.replaceOperand(jErrPos+1, errPos)
				if err := c.Compile(idxExpr); err != nil {
					return err
				}
				c.emit(code.OpIndex, 0) // optional = false

				// END label
				endPos := len(c.currentInstructions())
				c.replaceOperand(jEndPos+1, endPos)
			} else {
				if err := c.Compile(idxExpr); err != nil {
					return err
				}
				c.emit(code.OpIndex, 0) // optional = false
			}
		}
	}

	end := len(c.currentInstructions())
	for _, jp := range jumpPositions {
		c.replaceOperand(jp+1, end)
	}
	return nil
}

// compileNullishChain compiles a flattened sequence of terms for the nullish coalescing operator.
// For each term except the last, we compile the term with softening of only its final access (if any),
// then emit a JumpIfNotNullish to the end. The last term is compiled normally.
func (c *Compiler) compileNullishChain(terms []parser.Node) error {
	jumpPositions := make([]int, 0, len(terms))

	for i, term := range terms {
		isLast := i == len(terms)-1

		// For non-last terms, if term is an access chain, soften only the final access
		if !isLast {
			switch term.(type) {
			case *parser.MemberAccess, *parser.IndexAccess:
				if err := c.compileAccessNode(term, true); err != nil {
					return err
				}
			default:
				if err := c.Compile(term); err != nil {
					return err
				}
			}

			// Jump to END if result is not nullish
			pos := len(c.currentInstructions())
			c.emit(code.OpJumpIfNotNullish, 0) // placeholder to END
			jumpPositions = append(jumpPositions, pos)
			continue
		}

		// Last term: compile normally (no softening)
		if err := c.Compile(term); err != nil {
			return err
		}
	}

	end := len(c.currentInstructions())
	for _, p := range jumpPositions {
		c.replaceOperand(p+1, end)
	}
	return nil
}
