package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/maniartech/uexl_go/ast"
	. "github.com/maniartech/uexl_go/parser"
)

func TestArray(t *testing.T) {
	// Test case 1: Simple array [1, 2, 3]
	t.Run("Simple number array", func(t *testing.T) {
		got, err := ParseReaderNew("", strings.NewReader("[1, 2, 3]"))
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		gotNode, ok := got.(*ast.ArrayNode)
		if !ok {
			t.Fatalf("Expected *ast.ArrayNode, got %T", got)
		}

		if len(gotNode.Items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(gotNode.Items))
		}

		expectedValues := []float64{1, 2, 3}
		for i, expectedValue := range expectedValues {
			numberNode, ok := gotNode.Items[i].(*ast.NumberNode)
			if !ok {
				t.Errorf("Expected *ast.NumberNode at index %d, got %T", i, gotNode.Items[i])
				continue
			}
			if float64(numberNode.Value) != expectedValue {
				t.Errorf("Expected %v at index %d, got %v", expectedValue, i, numberNode.Value)
			}
		}
	})

	// Test case 2: Array with expressions [1 + 2, 3 + 4, true || false]
	t.Run("Array with expressions", func(t *testing.T) {
		got, err := ParseReaderNew("", strings.NewReader("[1 + 2, 3 + 4, true || false]"))
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		gotNode, ok := got.(*ast.ArrayNode)
		if !ok {
			t.Fatalf("Expected *ast.ArrayNode, got %T", got)
		}

		if len(gotNode.Items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(gotNode.Items))
		}

		// Test first expression: 1 + 2
		expNode1, ok := gotNode.Items[0].(*ast.ExpressionNode)
		if !ok {
			t.Errorf("Expected *ast.ExpressionNode at index 0, got %T", gotNode.Items[0])
		} else {
			if expNode1.Token != "1 + 2" {
				t.Errorf("Token: Expected '1 + 2', got %v", expNode1.Token)
			}

			leftNode, ok := expNode1.Left.(*ast.NumberNode)
			if !ok {
				t.Errorf("Expected *ast.NumberNode for left operand, got %T", expNode1.Left)
			} else {
				if leftNode.Token != "1" || float64(leftNode.Value) != 1 {
					t.Errorf("Left operand: Expected token '1' and value 1, got token '%s' and value %v", leftNode.Token, leftNode.Value)
				}
			}

			rightNode, ok := expNode1.Right.(*ast.NumberNode)
			if !ok {
				t.Errorf("Expected *ast.NumberNode for right operand, got %T", expNode1.Right)
			} else {
				if rightNode.Token != "2" || float64(rightNode.Value) != 2 {
					t.Errorf("Right operand: Expected token '2' and value 2, got token '%s' and value %v", rightNode.Token, rightNode.Value)
				}
			}
		}

		// Test second expression: 3 + 4
		expNode2, ok := gotNode.Items[1].(*ast.ExpressionNode)
		if !ok {
			t.Errorf("Expected *ast.ExpressionNode at index 1, got %T", gotNode.Items[1])
		} else {
			if expNode2.Token != "3 + 4" {
				t.Errorf("Token: Expected '3 + 4', got %v", expNode2.Token)
			}
		}

		// Test third expression: true || false
		expNode3, ok := gotNode.Items[2].(*ast.ExpressionNode)
		if !ok {
			t.Errorf("Expected *ast.ExpressionNode at index 2, got %T", gotNode.Items[2])
		} else {
			if expNode3.Token != "true || false" {
				t.Errorf("Token: Expected 'true || false', got %v", expNode3.Token)
			}

			leftBool, ok := expNode3.Left.(*ast.BooleanNode)
			if !ok {
				t.Errorf("Expected *ast.BooleanNode for left operand, got %T", expNode3.Left)
			} else {
				if leftBool.Token != "true" || !leftBool.Value {
					t.Errorf("Left boolean: Expected token 'true' and value true, got token '%s' and value %v", leftBool.Token, leftBool.Value)
				}
			}

			rightBool, ok := expNode3.Right.(*ast.BooleanNode)
			if !ok {
				t.Errorf("Expected *ast.BooleanNode for right operand, got %T", expNode3.Right)
			} else {
				if rightBool.Token != "false" || rightBool.Value {
					t.Errorf("Right boolean: Expected token 'false' and value false, got token '%s' and value %v", rightBool.Token, rightBool.Value)
				}
			}
		}
	})

	// Test case 3: Nested array with null [null, ['A', 'B']]
	t.Run("Array with null and nested array", func(t *testing.T) {
		got, err := ParseReaderNew("", strings.NewReader("[null, ['A', 'B']]"))
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		gotNode, ok := got.(*ast.ArrayNode)
		if !ok {
			t.Fatalf("Expected *ast.ArrayNode, got %T", got)
		}

		if len(gotNode.Items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(gotNode.Items))
		}

		// Test null element
		nullNode, ok := gotNode.Items[0].(*ast.NullNode)
		if !ok {
			t.Errorf("Expected *ast.NullNode at index 0, got %T", gotNode.Items[0])
		} else {
			if nullNode.Token != "null" {
				t.Errorf("Token: Expected 'null', got %v", nullNode.Token)
			}
		}

		// Test nested array ['A', 'B']
		nestedArray, ok := gotNode.Items[1].(*ast.ArrayNode)
		if !ok {
			t.Errorf("Expected *ast.ArrayNode at index 1, got %T", gotNode.Items[1])
		} else {
			if nestedArray.Token != "['A', 'B']" {
				t.Errorf("Token: Expected \"['A', 'B']\", got %v", nestedArray.Token)
			}

			if len(nestedArray.Items) != 2 {
				t.Errorf("Expected 2 items in nested array, got %d", len(nestedArray.Items))
			} else {
				// Test 'A'
				stringNodeA, ok := nestedArray.Items[0].(*ast.StringNode)
				if !ok {
					t.Errorf("Expected *ast.StringNode at nested index 0, got %T", nestedArray.Items[0])
				} else {
					if stringNodeA.Token != "'A'" || string(stringNodeA.Value) != "A" {
						t.Errorf("String A: Expected token \"'A'\" and value \"A\", got token '%s' and value '%s'", stringNodeA.Token, stringNodeA.Value)
					}
				}

				// Test 'B'
				stringNodeB, ok := nestedArray.Items[1].(*ast.StringNode)
				if !ok {
					t.Errorf("Expected *ast.StringNode at nested index 1, got %T", nestedArray.Items[1])
				} else {
					if stringNodeB.Token != "'B'" || string(stringNodeB.Value) != "B" {
						t.Errorf("String B: Expected token \"'B'\" and value \"B\", got token '%s' and value '%s'", stringNodeB.Token, stringNodeB.Value)
					}
				}
			}
		}
	})

	// Test case 4: Two-dimensional array [[1, 2, 3], [4, 5, 6]]
	t.Run("Two-dimensional array", func(t *testing.T) {
		got, err := ParseReaderNew("", strings.NewReader("[[1, 2, 3], [4, 5, 6]]"))
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		gotNode, ok := got.(*ast.ArrayNode)
		if !ok {
			t.Fatalf("Expected *ast.ArrayNode, got %T", got)
		}

		if len(gotNode.Items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(gotNode.Items))
		}

		// Test first array [1, 2, 3]
		firstArr, ok := gotNode.Items[0].(*ast.ArrayNode)
		if !ok {
			t.Errorf("Expected *ast.ArrayNode at index 0, got %T", gotNode.Items[0])
		} else {
			if firstArr.Token != "[1, 2, 3]" {
				t.Errorf("Token: Expected '[1, 2, 3]', got %v", firstArr.Token)
			}

			expectedFirstValues := []float64{1, 2, 3}
			if len(firstArr.Items) != 3 {
				t.Errorf("Expected 3 items in first array, got %d", len(firstArr.Items))
			} else {
				for i, expectedValue := range expectedFirstValues {
					numberNode, ok := firstArr.Items[i].(*ast.NumberNode)
					if !ok {
						t.Errorf("Expected *ast.NumberNode at first array index %d, got %T", i, firstArr.Items[i])
						continue
					}
					if numberNode.Token != fmt.Sprintf("%.0f", expectedValue) {
						t.Errorf("Token: Expected '%.0f', got %v", expectedValue, numberNode.Token)
					}
					if float64(numberNode.Value) != expectedValue {
						t.Errorf("Value: Expected %v, got %v", expectedValue, numberNode.Value)
					}
				}
			}
		}

		// Test second array [4, 5, 6]
		secondArr, ok := gotNode.Items[1].(*ast.ArrayNode)
		if !ok {
			t.Errorf("Expected *ast.ArrayNode at index 1, got %T", gotNode.Items[1])
		} else {
			if secondArr.Token != "[4, 5, 6]" {
				t.Errorf("Token: Expected '[4, 5, 6]', got %v", secondArr.Token)
			}

			expectedSecondValues := []float64{4, 5, 6}
			if len(secondArr.Items) != 3 {
				t.Errorf("Expected 3 items in second array, got %d", len(secondArr.Items))
			} else {
				for i, expectedValue := range expectedSecondValues {
					numberNode, ok := secondArr.Items[i].(*ast.NumberNode)
					if !ok {
						t.Errorf("Expected *ast.NumberNode at second array index %d, got %T", i, secondArr.Items[i])
						continue
					}
					if numberNode.Token != fmt.Sprintf("%.0f", expectedValue) {
						t.Errorf("Token: Expected '%.0f', got %v", expectedValue, numberNode.Token)
					}
					if float64(numberNode.Value) != expectedValue {
						t.Errorf("Value: Expected %v, got %v", expectedValue, numberNode.Value)
					}
				}
			}
		}
	})

	// Test case 5: Empty array
	t.Run("Empty array", func(t *testing.T) {
		got, err := ParseReaderNew("", strings.NewReader("[]"))
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		gotNode, ok := got.(*ast.ArrayNode)
		if !ok {
			t.Fatalf("Expected *ast.ArrayNode, got %T", got)
		}

		if len(gotNode.Items) != 0 {
			t.Errorf("Expected empty array, got %d items", len(gotNode.Items))
		}

		if gotNode.Token != "[]" {
			t.Errorf("Token: Expected '[]', got %v", gotNode.Token)
		}
	})
}
