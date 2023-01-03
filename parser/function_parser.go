package parser

import "github.com/maniartech/uexl_go/ast"

// parseFunction parses the function node
func parseFunction(token string, fn, args interface{}, offset, line, col int) (ast.Node, error) {
	argsSlice := toIfaceSlice(args)

	functionName := toIfaceSlice(fn)[0].(string)

	if len(argsSlice) == 0 {
		return ast.NewFunctionNode(token, functionName, []ast.Node{}, offset, line, col), nil
	}

	params := []interface{}{argsSlice[0]}
	restArgs := toIfaceSlice(argsSlice[1])

	for _, v := range restArgs {
		params = append(params, toIfaceSlice(v)[2])
	}

	return ast.NewFunctionNode(token, functionName, toNodesSlice(params), offset, line, col), nil
}
