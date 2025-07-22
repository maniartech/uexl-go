package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

const StackSize = 1024

var True = parser.BooleanLiteral{Value: true}
var False = parser.BooleanLiteral{Value: false}
var Null = parser.NullLiteral{}

type VM struct {
	constants []parser.Node
	stack     []parser.Node
	sp        int
}

func New() *VM {
	return &VM{
		constants: []parser.Node{},
		stack:     make([]parser.Node, StackSize),
		sp:        0,
	}
}
func (vm *VM) Push(node parser.Node) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = node
	vm.sp++
	return nil
}

func (vm *VM) Pop() parser.Node {
	if vm.sp == 0 {
		return nil
	}
	vm.sp--
	node := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil // Clear the reference
	return node
}

func (vm *VM) LastPoppedStackElem() parser.Node {
	return vm.stack[vm.sp]
}

func (vm *VM) Top() parser.Node {
	// Check if the stack is empty
	if vm.sp == 0 {
		return nil
	}
	// Return the top element without removing it
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var opcode code.Opcode
	for ip < len(ins) {
		ip++
		ins = ins[ip:]
		opcode = code.Opcode(ins[0])
		switch opcode {
		case code.OpConstant:
			// Read the constant index
			constIndex := code.ReadUint16(ins[1:])
			ins = ins[3:]
			// Push the constant onto the stack
			vm.constants = append(vm.constants, vm.constants[constIndex])
		}
	}
	return nil
}
