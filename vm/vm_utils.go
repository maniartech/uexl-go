package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

const StackSize = 1024
const MaxFrames = 1024

var True = parser.BooleanLiteral{Value: true}
var False = parser.BooleanLiteral{Value: false}
var Null = parser.NullLiteral{}

// Sentinel value to distinguish "variable not provided" from "variable is nil"
type contextVarMissing struct{}

var contextVarNotProvided = contextVarMissing{}

type VMFunctions map[string]func(args ...any) (any, error)
type PipeHandler func(
	input any,
	block any,
	alias string,
	vm *VM) (any, error)

type PipeHandlers map[string]PipeHandler

type LibContext struct {
	Functions    VMFunctions
	PipeHandlers PipeHandlers
}

type Frame struct {
	instructions code.Instructions
	ip           int
	basePointer  int
}

type VM struct {
	constants         []any
	contextVars       []string
	contextVarsValues map[string]any
	contextVarCache   []any          // Pre-resolved context var values for O(1) access
	lastContextValues map[string]any // Cache the last context values pointer to detect changes
	systemVars        []any
	aliasVars         map[string]any
	functionContext   VMFunctions
	pipeHandlers      PipeHandlers     // Add pipe handlers registry
	pipeScopes        []map[string]any // Add scope stack for pipe variables

	stack     []any
	sp        int
	frames    []*Frame
	framesIdx int
	safeMode  bool
}

func New(libCtx LibContext) *VM {
	if libCtx.PipeHandlers == nil {
		libCtx.PipeHandlers = make(PipeHandlers)
	}
	if libCtx.Functions == nil {
		libCtx.Functions = make(VMFunctions)
	}

	return &VM{
		functionContext: libCtx.Functions,
		pipeHandlers:    libCtx.PipeHandlers,
		frames:          make([]*Frame, MaxFrames),
		pipeScopes:      make([]map[string]any, 0),
		stack:           make([]any, StackSize),
		aliasVars:       make(map[string]any),
	}

}

func NewFrame(instructions code.Instructions, basePointer int) *Frame {
	return &Frame{
		instructions: instructions,
		ip:           0,
		basePointer:  basePointer,
	}
}

func (vm *VM) Push(node any) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = node
	vm.sp++
	return nil
}

// Type-specific push methods to avoid interface boxing overhead
// These eliminate runtime.convT* calls by storing typed values directly

func (vm *VM) pushFloat64(val float64) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = val
	vm.sp++
	return nil
}

func (vm *VM) pushString(val string) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = val
	vm.sp++
	return nil
}

func (vm *VM) pushBool(val bool) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = val
	vm.sp++
	return nil
}

func (vm *VM) Pop() any {
	if vm.sp == 0 {
		return nil
	}
	vm.sp--
	node := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil // Clear the reference
	return node
}

func (vm *VM) LastPoppedStackElem() any {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Top() any {
	// Check if the stack is empty
	if vm.sp == 0 {
		return nil
	}
	// Return the top element without removing it
	return vm.stack[vm.sp-1]
}

func (vm *VM) pushPipeScope() {
	vm.pipeScopes = append(vm.pipeScopes, make(map[string]any))
}

func (vm *VM) popPipeScope() {
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes = vm.pipeScopes[:len(vm.pipeScopes)-1]
	}
}

func (vm *VM) setPipeVar(name string, value any) {
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes[len(vm.pipeScopes)-1][name] = value
	}
}

func (vm *VM) getPipeVar(name string) (any, bool) {
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
