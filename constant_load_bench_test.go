package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl_go/vm"
)

// Benchmark ONLY constant loading to verify zero allocations
func BenchmarkVM_ConstantLoad(b *testing.B) {
	// This tests if pushValue (loading constants) allocates
	// Expression: just push one constant
	bytecode, err := compileExpression("42.0")
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{})
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		out, err := machine.Run(bytecode, nil)
		if err != nil {
			b.Fatal(err)
		}
		if out.(float64) != 42.0 {
			b.Fatalf("Expected 42.0, got %v", out)
		}
	}
}
