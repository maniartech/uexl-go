package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

func TestExpressionParsing(t *testing.T) {

	node, _ := parser.ParseString("x . y1 . 2 ")

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
