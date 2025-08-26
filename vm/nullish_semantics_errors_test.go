package vm_test

import (
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

// runVmExpectError compiles and runs the expression and expects a non-nil error.
func runVmExpectError(t *testing.T, input string) {
	t.Helper()
	// Any error at any stage (parse/compile/runtime) is acceptable for these negative tests.
	p := parser.NewParser(input)
	ast, err := p.Parse()
	if err != nil {
		return // parse error counts as expected error
	}
	comp := compiler.New()
	if err := comp.Compile(ast); err != nil {
		return // compile error counts as expected error
	}
	machine := vm.New(vm.LibContext{Functions: vm.Builtins, PipeHandlers: vm.DefaultPipeHandlers})
	if _, err := machine.Run(comp.ByteCode()); err != nil {
		return // runtime error as expected
	}
	t.Fatalf("expected error for input %q, got nil", input)
}

func TestNullishCoalescing_StrictErrorCases(t *testing.T) {
	// Earlier links must remain strict; only the immediate final access on the left of ?? is softened.
	cases := []string{
		// Missing intermediate link
		`{}` + `?.u?.name ?? "anon"`,
		// Intermediate exists but is nullish -> still error (earlier link strict)
		`{"u": null}` + `.u.name ?? "anon"`,
		// Index on nullish base before final access
		`{"arr": null}` + `.arr[0] ?? 1`,
		// Optional chaining guards only base nullish; missing key still errors
		`{}` + `?.u.name ?? "anon"`,
	}
	for _, in := range cases {
		runVmExpectError(t, in)
	}
}
