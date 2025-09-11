package parser_test

import (
	"strings"
	"testing"

	p "github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/parser/constants"
)

// runBench runs the tokenizer over the entire input once per iteration and reports allocs and bytes.
func runBench(b *testing.B, name string, input string) {
	b.Helper()
	if len(input) == 0 {
		b.Fatalf("empty input for benchmark %s", name)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tz := p.NewTokenizer(input)
			for {
				tok, err := tz.NextToken()
				if err != nil {
					b.Fatalf("tokenize error: %v", err)
				}
				if tok.Type == constants.TokenEOF {
					break
				}
			}
		}
	})
}

// runBenchParallel runs the tokenizer in parallel, one tokenizer per goroutine.
func runBenchParallel(b *testing.B, name string, input string) {
	b.Helper()
	if len(input) == 0 {
		b.Fatalf("empty input for benchmark %s", name)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.Run(name+"/parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				tz := p.NewTokenizer(input)
				for {
					tok, err := tz.NextToken()
					if err != nil {
						b.Fatalf("tokenize error: %v", err)
					}
					if tok.Type == constants.TokenEOF {
						break
					}
				}
			}
		})
	})
}

func repeatN(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s) * n)
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b.String()
}

func BenchmarkTokenizer_ScalarAndOps(b *testing.B) {
	base := `1 + 2*3 - 4/5 % 6 == 7 && 8 != 9 || 10 < 11 && 12 >= 13 << 2 >> 1`
	runBench(b, "scalar-ops/64B", base)
	runBench(b, "scalar-ops/4KB", repeatN(base+"\n", 64))
	runBenchParallel(b, "scalar-ops/4KB", repeatN(base+"\n", 64))
}

func BenchmarkTokenizer_Identifiers(b *testing.B) {
	base := `$alpha + $beta_123 - $γδθ + $with_unicode_变量`
	runBench(b, "idents/64B", base)
	runBench(b, "idents/8KB", repeatN(base+" ", 128))
	runBenchParallel(b, "idents/8KB", repeatN(base+" ", 128))
}

func BenchmarkTokenizer_Strings(b *testing.B) {
	// Mix of double, single and raw strings including escapes and doubled quotes
	base := `"hello\\nworld" 'single\\tquote' r"He said ""hello""" r'It''s raw'` +
		" " + `"unicode 😀 \"quote\"" 'emoji 😺 and \\ backslash'`
	runBench(b, "strings/128B", base)
	runBench(b, "strings/16KB", repeatN(base+" ", 128))
	runBenchParallel(b, "strings/16KB", repeatN(base+" ", 128))
}

func BenchmarkTokenizer_Pipes(b *testing.B) {
	// Pipes: simple, named, and mixed with operators
	base := `$x|:|map:foo|filter:bar|join:"," || $y|reduce:sum ?? 'n/a'`
	runBench(b, "pipes/128B", base)
	runBench(b, "pipes/16KB", repeatN(base+"\n", 128))
	runBenchParallel(b, "pipes/16KB", repeatN(base+"\n", 128))
}

func BenchmarkTokenizer_NullishOptional(b *testing.B) {
	base := `$obj?.prop?[i] ?? ($alt?.call() ?? null)`
	runBench(b, "nullish/64B", base)
	runBench(b, "nullish/8KB", repeatN(base+" ", 128))
	runBenchParallel(b, "nullish/8KB", repeatN(base+" ", 128))
}

func BenchmarkTokenizer_UnicodeHeavy(b *testing.B) {
	// Quote non-ASCII digits to avoid invalid-number errors; we want tokenizer cost, not errors.
	base := `"你好, мир, hello 😀" $π_变量 + '١٢٣٤٥' `
	runBench(b, "unicode/64B", base)
	runBench(b, "unicode/8KB", repeatN(base+"\n", 128))
	runBenchParallel(b, "unicode/8KB", repeatN(base+"\n", 128))
}

func BenchmarkTokenizer_LongNumbersAndSci(b *testing.B) {
	base := `1234567890.123456e-10 + 9876543210.987654E+20 - .5 + 42`
	runBench(b, "numbers/64B", base)
	runBench(b, "numbers/8KB", repeatN(base+" ", 128))
	runBenchParallel(b, "numbers/8KB", repeatN(base+" ", 128))
}

func BenchmarkTokenizer_WhitespaceHeavy(b *testing.B) {
	base := "\t  \n  $a   +   $b\n\n\r\n   ? .  ? [  | :  "
	runBench(b, "whitespace/64B", base)
	runBench(b, "whitespace/8KB", repeatN(base, 128))
	runBenchParallel(b, "whitespace/8KB", repeatN(base, 128))
}
