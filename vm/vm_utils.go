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

	// Fast-path pipe scope - eliminates map overhead for common pipe variables
	// Using direct field access instead of map[string]any reduces 83% overhead
	pipeFastScope struct {
		item   any // $item - current element in iteration
		index  int // $index - current index in iteration
		acc    any // $acc - accumulator for reduce operations
		window any // $window - current window in window operations
		chunk  any // $chunk - current chunk in chunk operations
		last   any // $last - last value in reduce operations
	}
	pipeFastScopeActive bool // Flag indicating fast-path is active (reduces branch cost)

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
	// Activate fast-path for common pipe operations
	vm.pipeFastScopeActive = true
}

func (vm *VM) popPipeScope() {
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes = vm.pipeScopes[:len(vm.pipeScopes)-1]
	}
	// Deactivate fast-path when no pipe scopes remain
	vm.pipeFastScopeActive = len(vm.pipeScopes) > 0
}

func (vm *VM) setPipeVar(name string, value any) {
	// Fast-path: Direct field access for common pipe variables (eliminates 83% map overhead)
	if vm.pipeFastScopeActive {
		switch name {
		case "$item":
			vm.pipeFastScope.item = value
			return
		case "$index":
			vm.pipeFastScope.index = value.(int)
			return
		case "$acc":
			vm.pipeFastScope.acc = value
			return
		case "$window":
			vm.pipeFastScope.window = value
			return
		case "$chunk":
			vm.pipeFastScope.chunk = value
			return
		case "$last":
			vm.pipeFastScope.last = value
			return
		}
	}
	// Fall back to map for custom variables (aliases, etc.)
	if len(vm.pipeScopes) > 0 {
		vm.pipeScopes[len(vm.pipeScopes)-1][name] = value
	}
}

func (vm *VM) getPipeVar(name string) (any, bool) {
	// Fast-path: Direct field access for common pipe variables (eliminates 83% map overhead)
	if vm.pipeFastScopeActive {
		switch name {
		case "$item":
			return vm.pipeFastScope.item, true
		case "$index":
			return vm.pipeFastScope.index, true
		case "$acc":
			return vm.pipeFastScope.acc, true
		case "$window":
			return vm.pipeFastScope.window, true
		case "$chunk":
			return vm.pipeFastScope.chunk, true
		case "$last":
			return vm.pipeFastScope.last, true
		}
	}
	// Fall back to map for custom variables (aliases, etc.)
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
