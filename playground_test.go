package uexl_go

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/parser"
)

func TestPlayground(t *testing.T) {
	tokens := parser.NewTokenizer(`-Inf`).PreloadTokens()
	fmt.Printf("tokens: %+v\n", tokens)
	// result, err := EvalExpr(tokens)
	// if err != nil {
	// 	t.Fatalf("failed to evaluate expression: %s", err)
	// }
	// fmt.Printf("result: %v\n", result)
}
