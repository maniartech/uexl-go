package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

const StackSize = 1024
const MaxFrames = 1024

var True = parser.BooleanLiteral{Value: true}
var False = parser.BooleanLiteral{Value: false}
var Null = parser.NullLiteral{}

type VMFunctions map[string]func(args ...parser.Node) (parser.Node, error)
type PipeHandler func(
	input parser.Node,
	block any,
	alias string,
	vm *VM) (parser.Node, error)

type PipeHandlers map[string]PipeHandler

type Frame struct {
	instructions code.Instructions
	ip           int
	basePointer  int
}

type VM struct {
	constants       []parser.Node
	contextVars     []parser.Node
	systemVars      []parser.Node
	aliasVars       map[string]parser.Node
	functionContext VMFunctions
	pipeHandlers    PipeHandlers             // Add pipe handlers registry
	pipeScopes      []map[string]parser.Node // Add scope stack for pipe variables

	instructions code.Instructions
	stack        []parser.Node
	sp           int
	frames       []*Frame
	framesIdx    int
}

func New(bytecode *compiler.ByteCode, functionContext VMFunctions, pipeHandlers PipeHandlers) *VM {
	mainFrame := NewFrame(bytecode.Instructions, 0)
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	if pipeHandlers == nil {
		pipeHandlers = make(PipeHandlers)
	}

	return &VM{
		constants:       bytecode.Constants,
		contextVars:     bytecode.ContextVars,
		systemVars:      bytecode.SystemVars,
		instructions:    bytecode.Instructions,
		aliasVars:       make(map[string]parser.Node),
		functionContext: functionContext,
		pipeHandlers:    pipeHandlers,
		pipeScopes:      make([]map[string]parser.Node, 0),
		stack:           make([]parser.Node, StackSize),
		sp:              0,
		frames:          frames,
		framesIdx:       1, // Start with the main frame at index 0
	}
}

func NewFrame(instructions code.Instructions, basePointer int) *Frame {
	return &Frame{
		instructions: instructions,
		ip:           0,
		basePointer:  basePointer,
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
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Top() parser.Node {
	// Check if the stack is empty
	if vm.sp == 0 {
		return nil
	}
	// Return the top element without removing it
	return vm.stack[vm.sp-1]
}

func (vm *VM) pushPipeScope() {
	vm.pipeScopes = append(vm.pipeScopes, make(map[string]parser.Node))
}

func (vm *VM) popPipeScope() {
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes = vm.pipeScopes[:len(vm.pipeScopes)-1]
	}
}

func (vm *VM) setPipeVar(name string, value parser.Node) {
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes[len(vm.pipeScopes)-1][name] = value
	}
}

func (vm *VM) getPipeVar(name string) (parser.Node, bool) {
	for i := len(vm.pipeScopes) - 1; i >= 0; i-- {
		if val, ok := vm.pipeScopes[i][name]; ok {
			return val, true
		}
	}
	return nil, false
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIdx] = f
	vm.framesIdx++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIdx--
	return vm.frames[vm.framesIdx]
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIdx-1]
}
