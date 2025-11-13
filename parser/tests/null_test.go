package parser_test

import (
	"testing"

	"github.com/maniartech/uexl/parser"
)

func TestNull(t *testing.T) {
	parserInstance := parser.NewParser("null")
	got, _ := parserInstance.Parse()

	gotNode, ok := got.(*parser.NullLiteral)
	if !ok {
		t.Fatalf("Expected *parser.NullLiteral, got %T", got)
	}

	if gotNode == nil {
		t.Fatalf("NullLiteral is nil")
	}
}
