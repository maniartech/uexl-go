package uexl

import (
	"context"
	"fmt"

	parsererrors "github.com/maniartech/uexl/parser/errors"
	"github.com/maniartech/uexl/vm"
)

// Type aliases — callers only need to import "github.com/maniartech/uexl";
// no sub-package imports are required.

// Function is a UExL callable: a built-in or user-registered function.
type Function = vm.VMFunction

// Functions is a registry mapping function names to their implementations.
type Functions = vm.VMFunctions

// PipeHandler is a function that implements a custom pipe operator.
type PipeHandler = vm.PipeHandler

// PipeHandlers is a registry mapping pipe names to their handler functions.
type PipeHandlers = vm.PipeHandlers

// PipeContext provides pipe handlers with access to predicate evaluation and the
// evaluation context. See §3.25 of the design spec for full interface semantics.
type PipeContext = vm.PipeContext

// ParserError is a single structured parse error (Line, Column, Code, Message).
// Re-exported so callers never need to import github.com/maniartech/uexl/parser/errors.
type ParserError = parsererrors.ParserError

// ParseErrors is a collection of parse errors produced by the parser.
// Re-exported so callers never need to import github.com/maniartech/uexl/parser/errors.
type ParseErrors = parsererrors.ParseErrors

// Option is an opaque functional option applied to an Env during construction.
// Create options via WithFunctions, WithPipeHandlers, WithGlobals, or WithLib.
type Option func(*envConfig)

// WithFunctions returns an Option that merges fns into the env's function registry.
// Later calls for the same key win. Panics if fns is nil.
func WithFunctions(fns Functions) Option {
	if fns == nil {
		panic("uexl: WithFunctions: fns must not be nil")
	}
	return func(cfg *envConfig) {
		for k, v := range fns {
			cfg.functions[k] = v
		}
	}
}

// WithPipeHandlers returns an Option that merges pipes into the env's pipe handler registry.
// Later calls for the same key win. Panics if pipes is nil.
func WithPipeHandlers(pipes PipeHandlers) Option {
	if pipes == nil {
		panic("uexl: WithPipeHandlers: pipes must not be nil")
	}
	return func(cfg *envConfig) {
		for k, v := range pipes {
			cfg.pipeHandlers[k] = v
		}
	}
}

// WithGlobals returns an Option that merges vars into the env's global variables.
// Global vars are shadowed by per-call vars of the same name. Panics if vars is nil.
func WithGlobals(vars map[string]any) Option {
	if vars == nil {
		panic("uexl: WithGlobals: vars must not be nil")
	}
	return func(cfg *envConfig) {
		for k, v := range vars {
			cfg.globals[k] = v
		}
	}
}

// WithLib returns an Option that calls lib.Apply during env construction, allowing
// the lib to register functions, pipe handlers, and globals in a single step.
// Panics if lib is nil.
func WithLib(lib Lib) Option {
	if lib == nil {
		panic("uexl: WithLib: lib must not be nil")
	}
	return func(cfg *envConfig) {
		lib.Apply(&EnvConfig{cfg: cfg})
	}
}

// DefaultWith returns a new *Env with all stdlib built-ins and default pipe handlers,
// plus any additional options applied on top. Does not modify the Default() singleton.
func DefaultWith(opts ...Option) *Env {
	return Default().Extend(opts...)
}

// Eval evaluates expr using the Default environment and context.Background().
// vars may be nil. For custom environments use NewEnv, DefaultWith, or Env.Extend.
func Eval(expr string, vars map[string]any) (any, error) {
	return Default().Eval(context.Background(), expr, vars)
}

// Validate parses, compiles, and validates expr against the default stdlib environment.
// Returns nil if valid; otherwise returns the first error encountered.
func Validate(expr string) error {
	return Default().Validate(expr)
}

// MustCompile compiles expr using Default() and panics on failure.
// Intended exclusively for package-level var declarations with known-valid expressions.
func MustCompile(expr string) *CompiledExpr {
	c, err := Default().Compile(expr)
	if err != nil {
		panic(fmt.Sprintf("uexl: MustCompile: %v", err))
	}
	return c
}
