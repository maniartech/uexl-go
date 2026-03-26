package uexl

import (
	"context"
	"fmt"
	"sync"

	"github.com/maniartech/uexl/code"
	"github.com/maniartech/uexl/compiler"
	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/vm"
)

// Env is an immutable, goroutine-safe evaluation environment.
// All fields are frozen after construction; Extend always produces a new Env with its own pool.
type Env struct {
	functions    vm.VMFunctions
	pipeHandlers vm.PipeHandlers
	globals      map[string]any
	pool         sync.Pool // per-Env — never copied by Extend
}

// newEnvFromConfig creates an Env from a finalized envConfig.
func newEnvFromConfig(cfg *envConfig) *Env {
	e := &Env{
		functions:    cfg.functions,
		pipeHandlers: cfg.pipeHandlers,
		globals:      cfg.globals,
	}
	// Capture e in the closure; safe because Env is heap-allocated and never moved.
	e.pool.New = func() any {
		return vm.New(vm.LibContext{
			Functions:    e.functions,
			PipeHandlers: e.pipeHandlers,
		})
	}
	return e
}

// NewEnv creates an Env from a blank slate — no built-ins, no pipes, no globals —
// then applies all provided options left-to-right (later call wins on conflict).
func NewEnv(opts ...Option) *Env {
	cfg := &envConfig{
		functions:    make(vm.VMFunctions),
		pipeHandlers: make(vm.PipeHandlers),
		globals:      make(map[string]any),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return newEnvFromConfig(cfg)
}

var (
	defaultEnvOnce sync.Once
	defaultEnv     *Env
)

// Default returns the singleton *Env pre-loaded with vm.Builtins as the function set
// and vm.DefaultPipeHandlers as the pipe handler set. No globals.
// The same pointer is returned on every call (initialized via sync.Once).
func Default() *Env {
	defaultEnvOnce.Do(func() {
		defaultEnv = NewEnv(
			WithFunctions(vm.Builtins),
			WithPipeHandlers(vm.DefaultPipeHandlers),
		)
	})
	return defaultEnv
}

// Extend creates a new *Env inheriting all functions, pipe handlers, and globals
// from the receiver, then applies opts on top. The receiver is never mutated.
func (e *Env) Extend(opts ...Option) *Env {
	cfg := &envConfig{
		functions:    copyMap(e.functions),
		pipeHandlers: copyMap(e.pipeHandlers),
		globals:      copyMap(e.globals),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return newEnvFromConfig(cfg)
}

// Compile parses and compiles expr into a *CompiledExpr bounded to this Env.
// All function call sites are validated against the env's registered functions
// at compile time — unknown functions are caught here, not at eval time.
// No VM is allocated during Compile.
func (e *Env) Compile(expr string) (*CompiledExpr, error) {
	node, err := parser.ParseString(expr)
	if err != nil {
		return nil, err
	}
	comp := compiler.New()
	if err := comp.Compile(node); err != nil {
		return nil, err
	}
	bc := comp.ByteCode()
	if err := e.validateFunctionNames(bc); err != nil {
		return nil, err
	}
	return &CompiledExpr{bytecode: bc, env: e}, nil
}

// validateFunctionNames walks the bytecode (main stream + InstructionBlock pipe predicates)
// and ensures every OpCallFunction references a function registered in e.functions.
func (e *Env) validateFunctionNames(bc *compiler.ByteCode) error {
	if err := e.walkInstructions(bc.Instructions, bc); err != nil {
		return err
	}
	// Also validate function calls inside InstructionBlock constants (pipe predicates).
	for _, cv := range bc.Constants {
		blk, ok := cv.ToAny().(*compiler.InstructionBlock)
		if !ok || blk == nil || blk.Instructions == nil {
			continue
		}
		if err := e.walkInstructions(blk.Instructions, bc); err != nil {
			return err
		}
	}
	return nil
}

// walkInstructions iterates over an instruction stream and validates any OpCallFunction sites.
func (e *Env) walkInstructions(ins code.Instructions, bc *compiler.ByteCode) error {
	for i := 0; i < len(ins); {
		def, err := code.Lookup(ins[i])
		if err != nil {
			// Unknown opcode — skip one byte defensively (should never happen in valid bytecode).
			i++
			continue
		}
		if code.Opcode(ins[i]) == code.OpCallFunction && i+2 < len(ins) {
			funcIdx := int(code.ReadUint16(ins[i+1 : i+3]))
			if funcIdx < len(bc.Constants) {
				if name, ok := bc.Constants[funcIdx].AsString(); ok {
					if _, exists := e.functions[name]; !exists {
						return fmt.Errorf("compile error: unknown function %q — not registered in this environment", name)
					}
				}
			}
		}
		// Advance past this opcode and its operands.
		offset := 1
		for _, w := range def.OperandWidths {
			offset += w
		}
		i += offset
	}
	return nil
}

// MustCompile compiles expr within this env and panics on failure.
// Intended exclusively for package-level var declarations with known-valid expressions.
func (e *Env) MustCompile(expr string) *CompiledExpr {
	c, err := e.Compile(expr)
	if err != nil {
		panic(fmt.Sprintf("uexl: Env.MustCompile: %v", err))
	}
	return c
}

// Validate parses, compiles, and validates expr within this environment.
// Returns nil if valid. No *CompiledExpr artifact is allocated on success.
func (e *Env) Validate(expr string) error {
	_, err := e.Compile(expr)
	return err
}

// Eval is a one-shot parse + compile + run within the environment.
// Context is forwarded to the VM for cancellation and deadline enforcement.
// For expressions evaluated repeatedly, prefer Compile + CompiledExpr.Eval.
func (e *Env) Eval(ctx context.Context, expr string, vars map[string]any) (any, error) {
	ce, err := e.Compile(expr)
	if err != nil {
		return nil, err
	}
	return ce.Eval(ctx, vars)
}

// Info returns a sorted, read-only snapshot of everything registered in this environment.
// The returned EnvInfo is independent — mutating its slices has no effect on the Env.
func (e *Env) Info() EnvInfo {
	return EnvInfo{
		Functions:    sortedKeys(e.functions),
		PipeHandlers: sortedKeys(e.pipeHandlers),
		Globals:      sortedKeys(e.globals),
	}
}

// HasFunction reports whether a function with the given name is registered.
// Returns false for empty string.
func (e *Env) HasFunction(name string) bool {
	if name == "" {
		return false
	}
	_, ok := e.functions[name]
	return ok
}

// HasPipe reports whether a pipe handler with the given name is registered.
// Returns false for empty string.
func (e *Env) HasPipe(name string) bool {
	if name == "" {
		return false
	}
	_, ok := e.pipeHandlers[name]
	return ok
}

// HasGlobal reports whether a global variable with the given name is registered.
// Returns false for empty string.
func (e *Env) HasGlobal(name string) bool {
	if name == "" {
		return false
	}
	_, ok := e.globals[name]
	return ok
}

// copyMap creates a shallow copy of a map.
func copyMap[V any](m map[string]V) map[string]V {
	cp := make(map[string]V, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
