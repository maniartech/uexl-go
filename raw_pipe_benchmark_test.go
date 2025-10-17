package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl_go/vm"
)

// Raw pipe benchmarks - testing pure pipe infrastructure overhead

func createRawTestArray(size int) []any {
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = float64(i + 1)
	}
	return arr
}

// Benchmark Filter with identity operation (just returns $item)
func BenchmarkPipe_Filter_Identity(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(100)}
	bytecode, err := compileExpression(`arr |filter: $item`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}

// Benchmark Map with identity operation (just returns $item)
func BenchmarkPipe_Map_Identity(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(100)}
	bytecode, err := compileExpression(`arr |map: $item`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}

// Benchmark Filter with simple arithmetic (minimal predicate work)
func BenchmarkPipe_Filter_TrueLiteral(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(100)}
	bytecode, err := compileExpression(`arr |filter: true`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}

// Benchmark Map with simple arithmetic
func BenchmarkPipe_Map_Add1(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(100)}
	bytecode, err := compileExpression(`arr |map: $item + 1.0`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}

// Benchmark different array sizes to understand scaling
func BenchmarkPipe_Filter_Identity_10(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(10)}
	bytecode, err := compileExpression(`arr |filter: $item`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}

func BenchmarkPipe_Filter_Identity_1000(b *testing.B) {
	params := map[string]any{"arr": createRawTestArray(1000)}
	bytecode, err := compileExpression(`arr |filter: $item`)
	if err != nil {
		b.Fatal(err)
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	var out any
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	_ = out
}
