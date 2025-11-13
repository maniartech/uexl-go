package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/vm"
)

// Benchmarks for Excel-compatible operators to verify performance targets

func BenchmarkExcel_Power_CaretOperator(b *testing.B) {
	// Test ^ power operator (Excel-compatible)
	expr := "2 ^ 10"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_Power_DoubleStar(b *testing.B) {
	// Test ** power operator (legacy)
	expr := "2 ** 10"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_BitwiseXor_Tilde(b *testing.B) {
	// Test ~ XOR operator (Lua-style)
	expr := "5 ~ 3"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_BitwiseNot_Tilde(b *testing.B) {
	// Test ~ NOT operator (unary)
	expr := "~5"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_NotEquals_DiamondOperator(b *testing.B) {
	// Test <> not-equals operator (Excel-compatible)
	expr := "5 <> 3"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}

func BenchmarkExcel_NotEquals_BangEquals(b *testing.B) {
	// Test != not-equals operator (legacy)
	expr := "5 != 3"
	node, _ := parser.ParseString(expr)
	comp := compiler.New()
	comp.Compile(node)
	bytecode := comp.ByteCode()

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		machine.Run(bytecode, nil)
	}
}
