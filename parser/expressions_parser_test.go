package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

func TestExpressionParsing(t *testing.T) {

	node, err := parser.ParseString("SUM({ 'x': ['a', 'b', 'c', {'d': x}]}.x.3.d.y1.2, 2)  == 3 + 2")

	// node, err := parser.ParseString("'hello world'.1 + 'l'")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	// fmt.Println(node)

	result, err := node.Eval(types.Context{
		"x": types.Object{
			"y": types.Object{
				"z": types.Number(5),
			},
			"y1": types.Array{
				types.Number(1),
				types.Number(2),
				types.Number(3),
			},
		},
	})

	fmt.Println("result", result, err)

	// parsed, err := parser.ParseReader("", strings.NewReader("test.0.abc"))
	// if err != nil {
	// 	t.Errorf("Error: %v", err)
	// 	return
	// }

	// fmt.Println(parsed)

}
