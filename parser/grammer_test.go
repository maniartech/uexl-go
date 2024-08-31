package parser_test

import (
	"log"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/parser"
)

func TestParser_old(t *testing.T) {
	// Test for the Modulus (//) Operator.
	_, err := parser.ParseReader("", strings.NewReader("10 + 4"))
	if err != nil {
		log.Fatal(err)
	}
	// expected := 2
	// if got != expected {
	// 	t.Errorf("Expected %v, got, %v\n", expected, got)
	// }
}
