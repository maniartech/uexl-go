package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/vm"
)

// Internal benchmarks to isolate VM handler performance (excludes Run() interface boxing)

func BenchmarkExcel_Internal_BitwiseNot(b *testing.B) {
	// Create a minimal bytecode for ~5
	constants := []any{5.0}
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 0), // Load 5
		code.Make(code.OpBitwiseNot),  // ~5
		code.Make(code.OpReturn),
	}

	var flatInstructions code.Instructions
	for _, inst := range instructions {
		flatInstructions = append(flatInstructions, inst...)
	}

	bytecode := &compiler.ByteCode{
		Instructions: flatInstructions,
		Constants:    constants,
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Run the VM but don't capture result (avoids interface boxing)
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_Internal_Power(b *testing.B) {
	// Create a minimal bytecode for 2 ^ 10
	constants := []any{2.0, 10.0}
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 0), // Load 2
		code.Make(code.OpConstant, 1), // Load 10
		code.Make(code.OpPow),         // 2 ^ 10
		code.Make(code.OpReturn),
	}

	var flatInstructions code.Instructions
	for _, inst := range instructions {
		flatInstructions = append(flatInstructions, inst...)
	}

	bytecode := &compiler.ByteCode{
		Instructions: flatInstructions,
		Constants:    constants,
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_Internal_BitwiseXor(b *testing.B) {
	// Create a minimal bytecode for 5 ~ 3
	constants := []any{5.0, 3.0}
	instructions := []code.Instructions{
		code.Make(code.OpConstant, 0), // Load 5
		code.Make(code.OpConstant, 1), // Load 3
		code.Make(code.OpBitwiseXor),  // 5 ~ 3
		code.Make(code.OpReturn),
	}

	var flatInstructions code.Instructions
	for _, inst := range instructions {
		flatInstructions = append(flatInstructions, inst...)
	}

	bytecode := &compiler.ByteCode{
		Instructions: flatInstructions,
		Constants:    constants,
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}
