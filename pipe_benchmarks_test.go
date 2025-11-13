package uexl_go

import (
	"testing"

	"github.com/maniartech/uexl/vm"
)

// Benchmark data generators for pipe operations
func createPipeTestArray(size int) []any {
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = float64(i + 1)
	}
	return arr
}

func createPipeTestObjects(size int) []any {
	arr := make([]any, size)
	for i := 0; i < size; i++ {
		arr[i] = map[string]any{
			"id":    float64(i + 1),
			"value": float64((i % 10) + 1),
			"name":  "item" + string(rune('A'+i%26)),
		}
	}
	return arr
}

// Helper to compile and run pipe expressions
func runPipeBenchmark(b *testing.B, expr string, params map[string]any) {
	bytecode, err := compileExpression(expr)
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

	// Basic validation
	_ = out
}

// ============================================================================
// MAP PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Map_Arithmetic_Small(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(10)}
	runPipeBenchmark(b, `arr |map: $item * 2.0`, params)
}

func BenchmarkPipe_Map_Arithmetic_Medium(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |map: $item * 2.0`, params)
}

func BenchmarkPipe_Map_Arithmetic_Large(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(1000)}
	runPipeBenchmark(b, `arr |map: $item * 2.0`, params)
}

func BenchmarkPipe_Map_ObjectAccess(b *testing.B) {
	params := map[string]any{"arr": createPipeTestObjects(100)}
	runPipeBenchmark(b, `arr |map: $item.value`, params)
}

// ============================================================================
// FILTER PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Filter_Simple(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |filter: $item > 50.0`, params)
}

func BenchmarkPipe_Filter_Complex(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |filter: $item > 25.0 && $item < 75.0`, params)
}

func BenchmarkPipe_Filter_ObjectProperty(b *testing.B) {
	params := map[string]any{"arr": createPipeTestObjects(100)}
	runPipeBenchmark(b, `arr |filter: $item.value > 5.0`, params)
}

// ============================================================================
// REDUCE PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Reduce_Sum(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |reduce: ($acc || 0) + $item`, params)
}

func BenchmarkPipe_Reduce_WithInitial(b *testing.B) {
	params := map[string]any{
		"arr":     createPipeTestArray(100),
		"initial": 0.0,
	}
	runPipeBenchmark(b, `arr |reduce: ($acc || initial) + $item`, params)
}

// ============================================================================
// FIND PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Find_First(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |find: $item > 5.0`, params)
}

func BenchmarkPipe_Find_Last(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |find: $item > 95.0`, params)
}

func BenchmarkPipe_Find_NotFound(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |find: $item > 1000.0`, params)
}

// ============================================================================
// SOME PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Some_True(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |some: $item > 50.0`, params)
}

func BenchmarkPipe_Some_False(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |some: $item > 1000.0`, params)
}

// ============================================================================
// EVERY PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Every_True(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |every: $item > 0.0`, params)
}

func BenchmarkPipe_Every_False(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |every: $item > 50.0`, params)
}

// ============================================================================
// UNIQUE PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Unique_NoDuplicates(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |unique: $item`, params)
}

func BenchmarkPipe_Unique_ManyDuplicates(b *testing.B) {
	// Array with values repeating 10 times each
	arr := make([]any, 100)
	for i := 0; i < 100; i++ {
		arr[i] = float64((i % 10) + 1)
	}
	params := map[string]any{"arr": arr}
	runPipeBenchmark(b, `arr |unique: $item`, params)
}

// ============================================================================
// SORT PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Sort_Ascending(b *testing.B) {
	// Create reversed array for worst-case scenario
	arr := make([]any, 100)
	for i := 0; i < 100; i++ {
		arr[i] = float64(100 - i)
	}
	params := map[string]any{"arr": arr}
	runPipeBenchmark(b, `arr |sort: $item`, params)
}

func BenchmarkPipe_Sort_WithPredicate(b *testing.B) {
	arr := make([]any, 100)
	for i := 0; i < 100; i++ {
		arr[i] = float64(100 - i)
	}
	params := map[string]any{"arr": arr}
	runPipeBenchmark(b, `arr |sort: $item`, params)
}

// ============================================================================
// GROUPBY PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_GroupBy_Simple(b *testing.B) {
	params := map[string]any{"arr": createPipeTestObjects(100)}
	runPipeBenchmark(b, `arr |groupBy: $item.value`, params)
}

func BenchmarkPipe_GroupBy_Expression(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |groupBy: $item % 10.0`, params)
}

// ============================================================================
// WINDOW PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Window_Size3(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |window: 3.0`, params)
}

func BenchmarkPipe_Window_Size10(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |window: 10.0`, params)
}

// ============================================================================
// CHUNK PIPE BENCHMARKS
// ============================================================================

func BenchmarkPipe_Chunk_Size10(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |chunk: 10.0`, params)
}

func BenchmarkPipe_Chunk_Size5(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |chunk: 5.0`, params)
}

// ============================================================================
// CHAINED PIPE BENCHMARKS (Real-world scenarios)
// ============================================================================

func BenchmarkPipe_Chain_FilterMap(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |filter: $item > 50.0 |map: $item * 2.0`, params)
}

func BenchmarkPipe_Chain_MapFilterReduce(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |map: $item * 2.0 |filter: $item > 100.0 |reduce: ($acc || 0) + $item`, params)
}

func BenchmarkPipe_Chain_Complex(b *testing.B) {
	params := map[string]any{"arr": createPipeTestArray(100)}
	runPipeBenchmark(b, `arr |filter: $item > 25.0 |map: $item * 2.0 |filter: $item < 150.0 |reduce: ($acc || 0) + $item`, params)
}
