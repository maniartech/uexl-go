package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
)

func TestStrings(t *testing.T) {

	// testStrs is list of random strings to test the parser.
	// It contains all the possible combinations of single and double quotes.
	// It also contains strings with escaped characters, strings with
	// escaped quotes and unicode characters. It contains around 50
	// strings to test the parser.
	testStrs := []string{
		"This \u00A1sentence has a Unicode character!",
		"This sentence has \"escaped quotes\"",
		"This sentence has an \\escaped character",
		"\u00A1This sentence has a Unicode character and \"escaped quotes\"",
		"This sentence has an \\escaped character and \"escaped quotes\"",
		"This sentence has an \\escaped character and \u00A1a Unicode character!",
		"\u00A1This sentence has a Unicode character and an \\escaped character",
		"This sentence has \"escaped quotes\" and \u00A1a Unicode character!",
		"\u00A1This sentence has a Unicode character and \"escaped quotes\" and an \\escaped character",
		"This sentence has an \\escaped character and \"escaped quotes\" and \u00A1a Unicode character!",
		"This sentence has an \\escaped character and \u00A1a Unicode character! and \"escaped quotes\"",
		"\u00A1This sentence has a Unicode character and an \\escaped character and \"escaped quotes\"",
		"This sentence has \"escaped quotes\" and an \\escaped character and \u00A1a Unicode character!",
		"\u00A1This sentence has a Unicode character and \"escaped quotes\" and an \\escaped character and \u00A1a Unicode character!",
		"This sentence has an \\escaped character and \"escaped quotes\" and \u00A1a Unicode character! and an \\escaped character",
	}

	for _, s := range testStrs {
		// sdq := "\"" + s + "\""

		// Test case for single quoted string
		ssq := "'" + s + "'"
		parsed, err := ParseReader("", strings.NewReader(ssq))
		if err != nil {
			t.Errorf("Error: %v", err)
			continue
		}

		node := parsed.(*ast.StringNode)
		token := node.Token
		value := node.Value

		fmt.Println("Token: ", token)
		fmt.Println("Value: ", value)
	}
}

// func TestStrings(t *testing.T) {
// Test case for 'hello'
// got, _ := ParseReader("", strings.NewReader("'hello'"))
// gotNode := got.(ast.StringNode)
// token := gotNode.Token
// value := gotNode.Value
// if token != "'hello'" {
// 	t.Errorf("Token: Expected 'hello', got %v", token)
// }
// if value != "hello" {
// 	t.Errorf("Value: Expected hello, got %v", token)
// }

// // Test case for 'world'
// got, _ = ParseReader("", strings.NewReader("\"world\""))
// gotNode = got.(ast.StringNode)
// token = gotNode.Token
// value = gotNode.Value
// if token != "world" {
// 	t.Errorf("Token: Expected world, got %v", token)
// }
// if value != "world" {
// 	t.Errorf("Value: Expected world, got %v", token)
// }

// // Test case for 'Hello' + 'World'
// got, _ = ParseReader("", strings.NewReader("'Hello' + 'World'"))
// gotExpNode := got.(ast.ExpressionNode)
// expToken := gotExpNode.Token
// if expToken != "'Hello' + 'World'" {
// 	t.Errorf("Value: Expected 'Hello' + 'World', got %v", expToken)
// }
// leftVal := gotExpNode.Left.(ast.StringNode).Value
// rightVal := gotExpNode.Right.(ast.StringNode).Value
// if leftVal != "Hello" {
// 	t.Errorf("Expected Hello, got %v", leftVal)
// }
// if rightVal != "World" {
// 	t.Errorf("Expected World, got %v", rightVal)
// }

// Test case for 'Hello' + ('Hello' + 'World')
// got, _ = ParseReader("", strings.NewReader("'Hello' + ('There' + 'World')"))
// gotExpNode = got.(ast.ExpressionNode)
// expToken = gotExpNode.Token
// if expToken != "'Hello' + ('There' + 'World')" {
// 	t.Errorf("Expected 'Hello' + ('There' + 'World'), got %v", expToken)
// }
// leftVal = gotExpNode.Left.(ast.StringNode).Value
// leftValToken := gotExpNode.Left.(ast.StringNode).Token
// rightValToken := gotExpNode.Right.(ast.ExpressionNode).Token
// rightLeftVal := gotExpNode.Right.(ast.ExpressionNode).Left.(ast.StringNode).Value
// rightRightVal := gotExpNode.Right.(ast.ExpressionNode).Right.(ast.StringNode).Value
// if leftValToken != "'Hello'" {
// 	t.Errorf("Expected 'Hello', got %v", leftValToken)
// }
// if leftVal != "Hello" {
// 	t.Errorf("Expected Hello, got %v", leftVal)
// }
// if rightValToken != "'There' + 'World'" {
// 	t.Errorf("Expected 'There' + 'World', got %v", rightValToken)
// }
// if rightLeftVal != "There" {
// 	t.Errorf("Expected There, got %v", rightLeftVal)
// }
// if rightRightVal != "World" {
// 	t.Errorf("Expected World, got %v", rightRightVal)
// }
// }
