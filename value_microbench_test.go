package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl/vm"
)

// Micro-benchmarks to verify zero-allocation primitive operations
// These test ONLY the VM execution without context variable overhead

func BenchmarkVM_PureArithmetic(b *testing.B) {
	// Expression with only literals - no context variables
	bytecode, err := compileExpression("(10.0 + 20.0) * 5.0")
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
		if out.(float64) != 150.0 {
			b.Fatalf("Expected 150.0, got %v", out)
		}
	}
}

func BenchmarkVM_PureBoolean(b *testing.B) {
	// Boolean expression with only literals
	bytecode, err := compileExpression("true && false || true")
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
		if !out.(bool) {
			b.Fail()
		}
	}
}

func BenchmarkVM_PureComparison(b *testing.B) {
	// Comparison with literals
	bytecode, err := compileExpression("10.0 > 5.0 && 20.0 == 20.0")
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
		if !out.(bool) {
			b.Fail()
		}
	}
}

func BenchmarkVM_PureString(b *testing.B) {
	// String concatenation with literals
	bytecode, err := compileExpression(`"hello" + ", " + "world"`)
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
		if out.(string) != "hello, world" {
			b.Fatalf("Expected 'hello, world', got %v", out)
		}
	}
}
