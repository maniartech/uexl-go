package parser_test

import (
	"fmt"
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/maniartech/uexl_go/types"
)

func TestExpressionParsing(t *testing.T) {

	node, err := parser.ParseString("SUM({ 'x': ['a', 'b', 'c', {'d': x}]}.x[3].d.y1[2], 2)  == 3 + 2 ")

	// node, err := parser.ParseString("{ 'a': 1, 'b': 2, 'c': 3, 'd': 4 }")

	// node, err := parser.ParseString("10.53 + ['testing']")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	// result, err := node.Eval(types.Context{
	// 	"x": types.Object{
	// 		"y": types.Object{
	// 			"z": types.Number(5),
	// 		},
	// 		"y1": types.Array{
	// 			types.Number(1),
	// 			types.Number(2),
	// 			types.Number(3),
	// 		},
	// 	},
	// })

	ctx, _ := types.JSONToContext(`
		{
			"y": {
				"z": 5
			},
			"y1": [1, 2, 3],
		}`)

	result, err := node.Eval(ctx)

	fmt.Println("result", result, err)

	// parsed, err := parser.ParseReader("", strings.NewReader("test.0.abc"))
	// if err != nil {
	// 	t.Errorf("Error: %v", err)
	// 	return
	// }

	// fmt.Println(parsed)

}
