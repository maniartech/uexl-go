package parser

import (
	"io"
	"io/ioutil"

	"github.com/maniartech/uexl_go/ast"
	"github.com/maniartech/uexl_go/types"
)

// ParseString parses the given expression and returns the AST Node.
// It allows you to parse the expression without having to create a file.
// For example:
//
//	node, err := ParseString("1 + 2")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(node.Eval(nil))
func ParseString(expr string) (ast.Node, error) {
	parser := NewParser(expr)
	result, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	// Convert the new parser's Expression type to ast.Node
	return convertExpressionToAST(result), nil
}

// ParseReaderNew parses input from an io.Reader and returns an AST Node.
// This function uses the new parser implementation.
func ParseReaderNew(filename string, r io.Reader) (ast.Node, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data))
}

// convertExpressionToAST converts the new parser's Expression types to the old AST types
// This provides compatibility between the new parser and existing AST evaluation code
func convertExpressionToAST(expr Expression) ast.Node {
	if expr == nil {
		return nil
	}

	switch e := expr.(type) {
	case *NumberLiteral:
		// NewNumberNode returns (*NumberNode, error), so we need to handle the error
		node, err := ast.NewNumberNode(e.Value, 0, e.Line, e.Column)
		if err != nil {
			// Fallback to a default value on error
			fallback, _ := ast.NewNumberNode("0", 0, e.Line, e.Column)
			return fallback
		}
		return node
	case *StringLiteral:
		// Create the AST node directly with the processed value and original token
		node := &ast.StringNode{
			BaseNode: &ast.BaseNode{
				Type:   ast.NodeTypeString,
				Line:   e.Line,
				Column: e.Column,
				Token:  e.Token,
			},
			Value:          types.String(e.Value),
			IsSingleQuoted: e.IsSingleQuoted,
		}
		return node
	case *BooleanLiteral:
		// NewBooleanNode expects string token, not bool value
		token := "false"
		if e.Value {
			token = "true"
		}
		node, err := ast.NewBooleanNode(token, 0, e.Line, e.Column)
		if err != nil {
			fallback, _ := ast.NewBooleanNode("false", 0, e.Line, e.Column)
			return fallback
		}
		return node
	case *NullLiteral:
		return ast.NewNullNode("null", 0, e.Line, e.Column)
	case *Identifier:
		node, err := ast.NewIdentifierNode(e.Name, 0, e.Line, e.Column)
		if err != nil {
			fallback, _ := ast.NewIdentifierNode("unknown", 0, e.Line, e.Column)
			return fallback
		}
		return node
	case *ArrayLiteral:
		var items []ast.Node
		for _, elem := range e.Elements {
			items = append(items, convertExpressionToAST(elem))
		}
		// Generate a proper token string for the array
		token := generateArrayToken(e)
		return ast.NewArrayNode(token, items, 0, e.Line, e.Column)
	case *ObjectLiteral:
		var items []ast.ObjectItem
		for key, value := range e.Properties {
			items = append(items, ast.ObjectItem{
				Key:   types.String(key), // Convert string to types.String
				Value: convertExpressionToAST(value),
			})
		}
		return ast.NewObjectNode("", items, 0, e.Line, e.Column)
	case *BinaryExpression:
		left := convertExpressionToAST(e.Left)
		right := convertExpressionToAST(e.Right)
		opType := ast.GetOperatorType(e.Operator)
		// Generate a token string for the expression
		token := generateExpressionToken(e)
		return ast.NewExpressionNode(token, e.Operator, opType, left, right, 0, e.Line, e.Column)
	case *FunctionCall:
		var args []ast.Node
		for _, arg := range e.Arguments {
			args = append(args, convertExpressionToAST(arg))
		}
		// Use NewFunctionNode instead of NewFunctionCallNode
		if ident, ok := e.Function.(*Identifier); ok {
			return ast.NewFunctionNode("", ident.Name, args, 0, e.Line, e.Column)
		}
		// Fallback: create a generic function node
		return ast.NewFunctionNode("", "unknown", args, 0, e.Line, e.Column)
	case *MemberAccess:
		object := convertExpressionToAST(e.Object)
		return ast.NewDotExpressionNode("", object, e.Property, 0, e.Line, e.Column)
	case *IndexAccess:
		// For now, convert IndexAccess to a DotExpression with string representation of index
		array := convertExpressionToAST(e.Array)
		// Convert index to string representation
		indexStr := "0" // default
		if numLit, ok := e.Index.(*NumberLiteral); ok {
			indexStr = numLit.Value
		}
		return ast.NewDotExpressionNode("", array, indexStr, 0, e.Line, e.Column)
	case *PipeExpression:
		var expressions []ast.Node
		for _, expr := range e.Expressions {
			expressions = append(expressions, convertExpressionToAST(expr))
		}
		// NewPipeNode signature: (token string, nodes []Node, offset, line, col int)
		return ast.NewPipeNode("", expressions, 0, e.Line, e.Column)
	default:
		// Fallback - create a generic node
		fallback, _ := ast.NewIdentifierNode("unknown", 0, 1, 1)
		return fallback
	}
}

// generateExpressionToken creates a string representation of a binary expression
func generateExpressionToken(expr *BinaryExpression) string {
	leftToken := getTokenFromExpression(expr.Left)
	rightToken := getTokenFromExpression(expr.Right)
	return leftToken + " " + expr.Operator + " " + rightToken
}

// generateArrayToken creates a string representation of an array
func generateArrayToken(arr *ArrayLiteral) string {
	if len(arr.Elements) == 0 {
		return "[]"
	}

	var tokens []string
	for _, elem := range arr.Elements {
		tokens = append(tokens, getTokenFromExpression(elem))
	}

	result := "["
	for i, token := range tokens {
		if i > 0 {
			result += ", "
		}
		result += token
	}
	result += "]"
	return result
}

// getTokenFromExpression extracts a token string from any expression
func getTokenFromExpression(expr Expression) string {
	switch e := expr.(type) {
	case *NumberLiteral:
		return e.Value
	case *StringLiteral:
		if e.IsRaw {
			return "r'" + e.Value + "'"
		}
		return "'" + e.Value + "'"
	case *BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *NullLiteral:
		return "null"
	case *Identifier:
		return e.Name
	case *BinaryExpression:
		return generateExpressionToken(e)
	case *ArrayLiteral:
		return generateArrayToken(e)
	default:
		return "unknown"
	}
}
