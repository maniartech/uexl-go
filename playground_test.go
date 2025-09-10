package uexl_go

import (
	"fmt"
	"testing"
)

func TestPlayground(t *testing.T) {
	result, err := EvalExpr(`{
		"a": {
			"b": 123
		}
	}?["a"].b["c"] ?? 'abc'`)
	if err != nil {
		t.Fatalf("failed to evaluate expression: %s", err)
	}
	fmt.Printf("result: %v\n", result)
}
