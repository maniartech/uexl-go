package uexl

import (
	"context"
	"sort"

	"github.com/maniartech/uexl/compiler"
	"github.com/maniartech/uexl/vm"
)

// CompiledExpr is an immutable pre-compiled expression.
// It is goroutine-safe — multiple goroutines may call Eval concurrently without
// any external synchronization.
type CompiledExpr struct {
	bytecode *compiler.ByteCode
	env      *Env
}

// Eval executes the pre-compiled bytecode against vars, honoring ctx for
// cancellation and deadline. A *vm.VM is borrowed from the env's pool for each
// call and returned via defer — zero allocation on the hot path when the pool
// has a spare VM.
//
// vars may be nil — treated as an empty map.
func (c *CompiledExpr) Eval(ctx context.Context, vars map[string]any) (any, error) {
	// Check for cancellation before borrowing from pool.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	machine := c.env.pool.Get().(*vm.VM)
	defer c.env.pool.Put(machine)

	machine.SetContext(ctx)
	return machine.Run(c.bytecode, mergeVars(c.env.globals, vars))
}

// Variables returns the sorted list of variable names (without $ prefix) that
// the expression references. Returns an empty slice (never nil) for expressions
// with no variable references. The returned slice is a copy — safe to mutate.
func (c *CompiledExpr) Variables() []string {
	if len(c.bytecode.ContextVars) == 0 {
		return []string{}
	}
	cp := make([]string, len(c.bytecode.ContextVars))
	copy(cp, c.bytecode.ContextVars)
	sort.Strings(cp)
	return cp
}

// Env returns the *Env the expression was compiled against. Allocation-free.
func (c *CompiledExpr) Env() *Env {
	return c.env
}

// mergeVars produces a merged variable map with eval vars shadowing env globals.
//
// Fast paths:
//   - globals empty → return vars directly (no allocation)
//   - vars empty/nil → return globals directly (VM treats it as read-only)
//   - both non-empty → allocate a new merged map
func mergeVars(globals map[string]any, vars map[string]any) map[string]any {
	if len(globals) == 0 {
		return vars
	}
	if len(vars) == 0 {
		return globals
	}
	merged := make(map[string]any, len(globals)+len(vars))
	for k, v := range globals {
		merged[k] = v
	}
	for k, v := range vars {
		merged[k] = v
	}
	return merged
}
