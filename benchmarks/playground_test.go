package benchmarks_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl/parser"
)

func TestPlayground(t *testing.T) {
	tokens := parser.NewTokenizer(`-Inf`).PreloadTokens()
	fmt.Printf("tokens: %+v\n", tokens)
}
