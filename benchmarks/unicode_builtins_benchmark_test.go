package benchmarks_test

import (
	"testing"

	"github.com/maniartech/uexl/vm"
)

// ---- Expressions ------------------------------------------------------------

const (
	benchmarkRuneLenASCII       = `runeLen("hello world")`
	benchmarkRuneLenUnicode     = `runeLen("naïve café résumé")`
	benchmarkGraphemeLenASCII   = `graphemeLen("hello world")`
	benchmarkGraphemeLenUnicode = `graphemeLen("naïve café résumé")`
	benchmarkGraphemeLenEmoji   = `graphemeLen("👨‍👩‍👧‍👦 hello 🎉")`

	benchmarkRuneSubstrASCII   = `runeSubstr("hello world", 3, 5)`
	benchmarkRuneSubstrUnicode = `runeSubstr("naïve café résumé", 2, 8)`

	benchmarkGraphemeSubstrASCII   = `graphemeSubstr("hello world", 3, 5)`
	benchmarkGraphemeSubstrUnicode = `graphemeSubstr("naïve café résumé", 2, 8)`
	benchmarkGraphemeSubstrEmoji   = `graphemeSubstr("👨‍👩‍👧‍👦 hello 🎉 world", 0, 5)`

	benchmarkRunesASCII   = `runes("hello")`
	benchmarkRunesUnicode = `runes("naïve")`

	benchmarkGraphemesASCII   = `graphemes("hello")`
	benchmarkGraphemesUnicode = `graphemes("naïve")`
	benchmarkGraphemesEmoji   = `graphemes("👨‍👩‍👧‍👦 hello 🎉")`

	benchmarkBytesASCII   = `bytes("hello")`
	benchmarkBytesUnicode = `bytes("naïve")`

	benchmarkJoinSmall = `join(["a", "b", "c", "d", "e"], "-")`
	benchmarkJoinNoSep = `join(["hello", " ", "world"])`

	// Round-trip: explode + transform + join
	benchmarkExplodeJoinASCII   = `join(runes("hello"), "")`
	benchmarkExplodeJoinUnicode = `join(runes("naïve"), "")`

	// Realistic: filter combining marks then rejoin
	benchmarkFilterRunes = `join(runes("naïve") |map: $item, "")`
)

// ---- Helpers ----------------------------------------------------------------

func newUnicodeMachine() *vm.VM {
	return vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
}

func runUnicodeBenchmark(b *testing.B, expr string) {
	b.Helper()
	bytecode, err := compileExpression(expr)
	if err != nil {
		b.Fatalf("compile: %v", err)
	}
	machine := newUnicodeMachine()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = machine.Run(bytecode, nil)
		if err != nil {
			b.Fatalf("run: %v", err)
		}
	}
}

// ---- runeLen ----------------------------------------------------------------

func BenchmarkBuiltin_RuneLen_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRuneLenASCII)
}

func BenchmarkBuiltin_RuneLen_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRuneLenUnicode)
}

// ---- graphemeLen ------------------------------------------------------------

func BenchmarkBuiltin_GraphemeLen_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeLenASCII)
}

func BenchmarkBuiltin_GraphemeLen_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeLenUnicode)
}

func BenchmarkBuiltin_GraphemeLen_Emoji(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeLenEmoji)
}

// ---- runeSubstr -------------------------------------------------------------

func BenchmarkBuiltin_RuneSubstr_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRuneSubstrASCII)
}

func BenchmarkBuiltin_RuneSubstr_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRuneSubstrUnicode)
}

// ---- graphemeSubstr ---------------------------------------------------------

func BenchmarkBuiltin_GraphemeSubstr_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeSubstrASCII)
}

func BenchmarkBuiltin_GraphemeSubstr_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeSubstrUnicode)
}

func BenchmarkBuiltin_GraphemeSubstr_Emoji(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemeSubstrEmoji)
}

// ---- runes ------------------------------------------------------------------

func BenchmarkBuiltin_Runes_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRunesASCII)
}

func BenchmarkBuiltin_Runes_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkRunesUnicode)
}

// ---- graphemes --------------------------------------------------------------

func BenchmarkBuiltin_Graphemes_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemesASCII)
}

func BenchmarkBuiltin_Graphemes_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemesUnicode)
}

func BenchmarkBuiltin_Graphemes_Emoji(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkGraphemesEmoji)
}

// ---- bytes ------------------------------------------------------------------

func BenchmarkBuiltin_Bytes_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkBytesASCII)
}

func BenchmarkBuiltin_Bytes_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkBytesUnicode)
}

// ---- join -------------------------------------------------------------------

func BenchmarkBuiltin_Join_WithSep(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkJoinSmall)
}

func BenchmarkBuiltin_Join_NoSep(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkJoinNoSep)
}

// ---- round-trips ------------------------------------------------------------

func BenchmarkBuiltin_ExplodeJoin_ASCII(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkExplodeJoinASCII)
}

func BenchmarkBuiltin_ExplodeJoin_Unicode(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkExplodeJoinUnicode)
}

func BenchmarkBuiltin_FilterRunes(b *testing.B) {
	runUnicodeBenchmark(b, benchmarkFilterRunes)
}
