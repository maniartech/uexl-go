package parser_test

import (
	"regexp"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

func BenchmarkParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseString("AVERAGE(10, 20)")
	}
}

func BenchmarkRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		regexp.Compile(`^([A-Z]+)\(([0-9]+), ([0-9]+)\)$`)
	}
}

func BenchmarkEvaluation(b *testing.B) {
	node, err := parser.ParseString("AVERAGE(10, 20) + 30 - 40")
	if err != nil {
		b.Errorf("ParseString() error = %v", err)
	}

	for i := 0; i < b.N; i++ {
		node.Eval(nil)
	}
}

func TestParseString(t *testing.T) {
	node, err := parser.ParseString(`r'tes\nting⭐' + 123`)
	if err != nil {
		t.Errorf("ParseString() error = %v", err)
	}

	if node == nil {
		t.Errorf("ParseString() node = %v", node)
	}

	res, err := node.Eval(nil)
	if err != nil {
		t.Errorf("ParseString() error = %v", err)
	}

	if res != types.String("tes\\nting⭐123") {
		t.Errorf("ParseString() res = %v", res)
	}

}
