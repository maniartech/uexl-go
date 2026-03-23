// Package uexl provides a bytecode-compiled expression evaluator for Go.
//
// UExL parses expressions into an AST, compiles them into bytecode, and executes
// them on a small VM with explicit nullish semantics and pipe-based transforms.
//
// Quick start:
//
//	result, err := uexl.EvalExpr("10 + 20 |: $1 * 2")
//	if err != nil {
//		// handle error
//	}
//	_ = result
//
// For language and syntax docs, see the book directory in this repository.
package uexl
