// Package uexl provides a bytecode-compiled, embedded expression evaluation engine
// with explicit nullish/boolish semantics and pipe-based data transformations.
//
// # Quick Start
//
// Evaluate a simple expression:
//
//	result, err := uexl.Eval("price * qty", map[string]any{"price": 9.99, "qty": 3.0})
//
// # Environments
//
// An Env bundles functions, pipe handlers, and globals into an immutable,
// goroutine-safe evaluation context:
//
//	env := uexl.DefaultWith(uexl.WithFunctions(uexl.Functions{"discount": discountFn}))
//
// # Compile Once, Evaluate Many Times
//
//	rule, err := env.Compile("price * qty * (1 - discount)")
//	for _, order := range orders { total, _ := rule.Eval(ctx, order) }
//
// For full language, syntax, and API docs, see the book/ directory.
package uexl
