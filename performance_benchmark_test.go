package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl/compiler"
	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/vm"
)

// Test data similar to the comparison project
const benchmarkBooleanExpr = `(Origin == "MOW" || Country == "RU") && (Value >= 100.0 || Adults == 1.0)`
const benchmarkArithmeticExpr = `(a + b) * c - d / e`
const benchmarkStringExpr = `"hello" + ", world"`
const benchmarkStringCompareExpr = `name == "/groups/" + group + "/bar"`
const benchmarkMapExpr = `array |map: $item * 2.0`

func createBenchmarkParams() map[string]any {
	return map[string]any{
		"Origin":  "MOW",
		"Country": "RU",
		"Adults":  1.0,
		"Value":   100.0,
	}
}

func createBenchmarkArithmeticParams() map[string]any {
	return map[string]any{
		"a": 10.0,
		"b": 20.0,
		"c": 5.0,
		"d": 100.0,
		"e": 4.0,
	}
}

func createBenchmarkStringParams() map[string]any {
	return map[string]any{
		"name":  "/groups/foo/bar",
		"group": "foo",
	}
}

func createBenchmarkArrayParams() map[string]any {
	// Create array of float64 for UExL compatibility
	array := make([]any, 100)
	for i := 0; i < 100; i++ {
		array[i] = float64(i + 1)
	}
	return map[string]any{
		"array": array,
	}
}

func compileExpression(expr string) (*compiler.ByteCode, error) {
	node, err := parser.ParseString(expr)
	if err != nil {
		return nil, err
	}

	comp := compiler.New()
	err = comp.Compile(node)
	if err != nil {
		return nil, err
	}

	return comp.ByteCode(), nil
}

// Current implementation benchmark (creates VM per iteration)
func BenchmarkVM_Boolean_Current(b *testing.B) {
	params := createBenchmarkParams()
	bytecode, err := compileExpression(benchmarkBooleanExpr)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if !out.(bool) {
		b.Fail()
	}
}

func BenchmarkVM_Arithmetic_Current(b *testing.B) {
	params := createBenchmarkArithmeticParams()
	bytecode, err := compileExpression(benchmarkArithmeticExpr)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	// (10 + 20) * 5 - 100 / 4 = 30 * 5 - 25 = 150 - 25 = 125
	if out.(float64) != 125.0 {
		b.Fatalf("Expected 125.0, got %v", out)
	}
}

func BenchmarkVM_String_Current(b *testing.B) {
	bytecode, err := compileExpression(benchmarkStringExpr)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if out.(string) != "hello, world" {
		b.Fail()
	}
}

func BenchmarkVM_StringCompare_Current(b *testing.B) {
	params := createBenchmarkStringParams()
	bytecode, err := compileExpression(benchmarkStringCompareExpr)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if !out.(bool) {
		b.Fail()
	}
}

func BenchmarkVM_Map_Current(b *testing.B) {
	params := createBenchmarkArrayParams()
	bytecode, err := compileExpression(benchmarkMapExpr)
	if err != nil {
		b.Fatal(err)
	}

	var out any

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = machine.Run(bytecode, params)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if outArray, ok := out.([]any); !ok || len(outArray) == 0 || outArray[0] != 2.0 {
		b.Fail()
	}
}

// Baseline compilation benchmark (to separate compilation from execution cost)
func BenchmarkCompilation_Boolean(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := compileExpression(benchmarkBooleanExpr)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func BenchmarkCompilation_String(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := compileExpression(benchmarkStringExpr)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

// Optimized benchmarks will be added here as we implement improvements
// These will serve as comparison points to track our progress

// TODO: Add these benchmarks as we implement optimizations:
// - BenchmarkVM_Boolean_FastPath (optimized execution path)
// - BenchmarkVM_Boolean_JIT (JIT compilation if implemented)
