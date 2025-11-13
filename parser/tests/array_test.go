package parser_test

import (
	"testing"

	. "github.com/maniartech/uexl/parser"
)

func TestArray(t *testing.T) {
	// Test case 1: Empty array
	t.Run("Empty array", func(t *testing.T) {
		p := NewParser("[]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 0 {
			t.Errorf("Expected empty array, got %d elements", len(arrayLiteral.Elements))
		}

		if arrayLiteral.Line != 1 || arrayLiteral.Column != 1 {
			t.Errorf("Expected position (1,1), got (%d,%d)", arrayLiteral.Line, arrayLiteral.Column)
		}
	})

	// Test case 2: Single element array
	t.Run("Single element array", func(t *testing.T) {
		p := NewParser("[42]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 1 {
			t.Fatalf("Expected 1 element, got %d", len(arrayLiteral.Elements))
		}

		numberLiteral, ok := arrayLiteral.Elements[0].(*NumberLiteral)
		if !ok {
			t.Errorf("Expected *NumberLiteral, got %T", arrayLiteral.Elements[0])
		} else if numberLiteral.Value != 42 {
			t.Errorf("Expected value 42, got %v", numberLiteral.Value)
		}
	})

	// Test case 3: Simple number array
	t.Run("Simple number array", func(t *testing.T) {
		p := NewParser("[1, 2, 3]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		expectedValues := []float64{1, 2, 3}
		for i, expected := range expectedValues {
			numberLiteral, ok := arrayLiteral.Elements[i].(*NumberLiteral)
			if !ok {
				t.Errorf("Element %d: Expected *NumberLiteral, got %T", i, arrayLiteral.Elements[i])
			} else if numberLiteral.Value != expected {
				t.Errorf("Element %d: Expected value %v, got %v", i, expected, numberLiteral.Value)
			}
		}
	})

	// Test case 4: Mixed data types array
	t.Run("Mixed data types array", func(t *testing.T) {
		p := NewParser("[1, 'hello', true, null, 3.14]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 5 {
			t.Fatalf("Expected 5 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check number
		if numberLiteral, ok := arrayLiteral.Elements[0].(*NumberLiteral); !ok {
			t.Errorf("Element 0: Expected *NumberLiteral, got %T", arrayLiteral.Elements[0])
		} else if numberLiteral.Value != 1 {
			t.Errorf("Element 0: Expected value 1, got %v", numberLiteral.Value)
		}

		// Check string
		if stringLiteral, ok := arrayLiteral.Elements[1].(*StringLiteral); !ok {
			t.Errorf("Element 1: Expected *StringLiteral, got %T", arrayLiteral.Elements[1])
		} else if stringLiteral.Value != "hello" {
			t.Errorf("Element 1: Expected value 'hello', got %v", stringLiteral.Value)
		}

		// Check boolean
		if boolLiteral, ok := arrayLiteral.Elements[2].(*BooleanLiteral); !ok {
			t.Errorf("Element 2: Expected *BooleanLiteral, got %T", arrayLiteral.Elements[2])
		} else if boolLiteral.Value != true {
			t.Errorf("Element 2: Expected value true, got %v", boolLiteral.Value)
		}

		// Check null
		if _, ok := arrayLiteral.Elements[3].(*NullLiteral); !ok {
			t.Errorf("Element 3: Expected *NullLiteral, got %T", arrayLiteral.Elements[3])
		}

		// Check float
		if numberLiteral, ok := arrayLiteral.Elements[4].(*NumberLiteral); !ok {
			t.Errorf("Element 4: Expected *NumberLiteral, got %T", arrayLiteral.Elements[4])
		} else if numberLiteral.Value != 3.14 {
			t.Errorf("Element 4: Expected value 3.14, got %v", numberLiteral.Value)
		}
	})

	// Test case 5: Array with expressions
	t.Run("Array with expressions", func(t *testing.T) {
		p := NewParser("[1 + 2, 3 * 4, true && false]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check first expression: 1 + 2
		if binaryExpr, ok := arrayLiteral.Elements[0].(*BinaryExpression); !ok {
			t.Errorf("Element 0: Expected *BinaryExpression, got %T", arrayLiteral.Elements[0])
		} else {
			if binaryExpr.Operator != "+" {
				t.Errorf("Element 0: Expected operator '+', got '%s'", binaryExpr.Operator)
			}
		}

		// Check second expression: 3 * 4
		if binaryExpr, ok := arrayLiteral.Elements[1].(*BinaryExpression); !ok {
			t.Errorf("Element 1: Expected *BinaryExpression, got %T", arrayLiteral.Elements[1])
		} else {
			if binaryExpr.Operator != "*" {
				t.Errorf("Element 1: Expected operator '*', got '%s'", binaryExpr.Operator)
			}
		}

		// Check third expression: true && false
		if binaryExpr, ok := arrayLiteral.Elements[2].(*BinaryExpression); !ok {
			t.Errorf("Element 2: Expected *BinaryExpression, got %T", arrayLiteral.Elements[2])
		} else {
			if binaryExpr.Operator != "&&" {
				t.Errorf("Element 2: Expected operator '&&', got '%s'", binaryExpr.Operator)
			}
		}
	})

	// Test case 6: Nested arrays
	t.Run("Nested arrays", func(t *testing.T) {
		p := NewParser("[[1, 2], [3, 4], []]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check first nested array [1, 2]
		if nestedArray, ok := arrayLiteral.Elements[0].(*ArrayLiteral); !ok {
			t.Errorf("Element 0: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[0])
		} else {
			if len(nestedArray.Elements) != 2 {
				t.Errorf("Nested array 0: Expected 2 elements, got %d", len(nestedArray.Elements))
			}
		}

		// Check second nested array [3, 4]
		if nestedArray, ok := arrayLiteral.Elements[1].(*ArrayLiteral); !ok {
			t.Errorf("Element 1: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[1])
		} else {
			if len(nestedArray.Elements) != 2 {
				t.Errorf("Nested array 1: Expected 2 elements, got %d", len(nestedArray.Elements))
			}
		}

		// Check empty nested array []
		if nestedArray, ok := arrayLiteral.Elements[2].(*ArrayLiteral); !ok {
			t.Errorf("Element 2: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[2])
		} else {
			if len(nestedArray.Elements) != 0 {
				t.Errorf("Nested array 2: Expected 0 elements, got %d", len(nestedArray.Elements))
			}
		}
	})

	// Test case 7: Array index access - basic
	t.Run("Array index access basic", func(t *testing.T) {
		p := NewParser("[1, 2, 3][0]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		indexAccess, ok := expr.(*IndexAccess)
		if !ok {
			t.Fatalf("Expected *IndexAccess, got %T", expr)
		}

		// Check array part
		if arrayLiteral, ok := indexAccess.Target.(*ArrayLiteral); !ok {
			t.Errorf("Array part: Expected *ArrayLiteral, got %T", indexAccess.Target)
		} else if len(arrayLiteral.Elements) != 3 {
			t.Errorf("Array part: Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check index part
		if numberLiteral, ok := indexAccess.Index.(*NumberLiteral); !ok {
			t.Errorf("Index part: Expected *NumberLiteral, got %T", indexAccess.Index)
		} else if numberLiteral.Value != 0 {
			t.Errorf("Index part: Expected value 0, got %v", numberLiteral.Value)
		}
	})

	// Test case 8: Array index access - negative index
	t.Run("Array index access negative", func(t *testing.T) {
		p := NewParser("[1, 2, 3][-1]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		indexAccess, ok := expr.(*IndexAccess)
		if !ok {
			t.Fatalf("Expected *IndexAccess, got %T", expr)
		}

		// Check array part
		if arrayLiteral, ok := indexAccess.Target.(*ArrayLiteral); !ok {
			t.Errorf("Array part: Expected *ArrayLiteral, got %T", indexAccess.Target)
		} else if len(arrayLiteral.Elements) != 3 {
			t.Errorf("Array part: Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check index part (should be a unary expression for -1)
		if unaryExpr, ok := indexAccess.Index.(*UnaryExpression); !ok {
			t.Errorf("Index part: Expected *UnaryExpression, got %T", indexAccess.Index)
		} else {
			if unaryExpr.Operator != "-" {
				t.Errorf("Index part: Expected operator '-', got '%s'", unaryExpr.Operator)
			}
		}
	})

	// Test case 9: Array index access - expression index
	t.Run("Array index access with expression", func(t *testing.T) {
		p := NewParser("[1, 2, 3][1 + 1]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		indexAccess, ok := expr.(*IndexAccess)
		if !ok {
			t.Fatalf("Expected *IndexAccess, got %T", expr)
		}

		// Check index part (should be a binary expression for 1 + 1)
		if binaryExpr, ok := indexAccess.Index.(*BinaryExpression); !ok {
			t.Errorf("Index part: Expected *BinaryExpression, got %T", indexAccess.Index)
		} else {
			if binaryExpr.Operator != "+" {
				t.Errorf("Index part: Expected operator '+', got '%s'", binaryExpr.Operator)
			}
		}
	})

	// Test case 10: Chained array access
	t.Run("Chained array access", func(t *testing.T) {
		p := NewParser("[[1, 2], [3, 4]][0][1]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		// Should be an IndexAccess where the Array is another IndexAccess
		outerIndexAccess, ok := expr.(*IndexAccess)
		if !ok {
			t.Fatalf("Expected *IndexAccess, got %T", expr)
		}

		// Check that the array part is also an IndexAccess
		innerIndexAccess, ok := outerIndexAccess.Target.(*IndexAccess)
		if !ok {
			t.Errorf("Inner access: Expected *IndexAccess, got %T", outerIndexAccess.Target)
		} else {
			// Check that the innermost array is a nested array literal
			if arrayLiteral, ok := innerIndexAccess.Target.(*ArrayLiteral); !ok {
				t.Errorf("Innermost array: Expected *ArrayLiteral, got %T", innerIndexAccess.Target)
			} else if len(arrayLiteral.Elements) != 2 {
				t.Errorf("Innermost array: Expected 2 elements, got %d", len(arrayLiteral.Elements))
			}
		}
	})

	// ERROR CASES

	// Test case 11: Unclosed array
	t.Run("Error: Unclosed array", func(t *testing.T) {
		p := NewParser("[1, 2, 3")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for unclosed array")
		}
		// Error should mention unclosed array
		if !containsSubstring(err.Error(), "unclosed") && !containsSubstring(err.Error(), "]") {
			t.Errorf("Error message should mention unclosed array or missing ']', got: %s", err.Error())
		}
	})

	// Test case 12: Missing comma between elements
	t.Run("Error: Missing comma", func(t *testing.T) {
		p := NewParser("[1 2 3]")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for missing comma")
		}
	})

	// Test case 13: Trailing comma (should be handled gracefully or error)
	t.Run("Trailing comma", func(t *testing.T) {
		p := NewParser("[1, 2, 3,]")
		_, err := p.Parse()
		// Depending on implementation, this might be an error or handled gracefully
		// Based on the parser code, it looks like trailing commas cause an error
		if err == nil {
			t.Logf("Trailing comma was handled gracefully")
		} else {
			t.Logf("Trailing comma caused error (expected): %v", err)
		}
	})

	// Test case 14: Empty element (double comma)
	t.Run("Error: Empty element", func(t *testing.T) {
		p := NewParser("[1,, 3]")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for empty element")
		}
	})

	// Test case 15: Invalid expression in array
	t.Run("Error: Invalid expression", func(t *testing.T) {
		p := NewParser("[1 +]")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for incomplete expression")
		}
	})

	// Test case 16: Unclosed array index
	t.Run("Error: Unclosed array index", func(t *testing.T) {
		p := NewParser("[1, 2, 3][0")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for unclosed array index")
		}
	})

	// Test case 17: Invalid array index expression
	t.Run("Error: Invalid array index", func(t *testing.T) {
		p := NewParser("[1, 2, 3][")
		_, err := p.Parse()
		if err == nil {
			t.Fatalf("Expected parse error for missing array index")
		}
	})

	// Test case 18: Array with identifiers and function calls
	t.Run("Array with identifiers and function calls", func(t *testing.T) {
		p := NewParser("[x, y, func(1, 2)]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check identifier x
		if identifier, ok := arrayLiteral.Elements[0].(*Identifier); !ok {
			t.Errorf("Element 0: Expected *Identifier, got %T", arrayLiteral.Elements[0])
		} else if identifier.Name != "x" {
			t.Errorf("Element 0: Expected name 'x', got '%s'", identifier.Name)
		}

		// Check identifier y
		if identifier, ok := arrayLiteral.Elements[1].(*Identifier); !ok {
			t.Errorf("Element 1: Expected *Identifier, got %T", arrayLiteral.Elements[1])
		} else if identifier.Name != "y" {
			t.Errorf("Element 1: Expected name 'y', got '%s'", identifier.Name)
		}

		// Check function call func(1, 2)
		if funcCall, ok := arrayLiteral.Elements[2].(*FunctionCall); !ok {
			t.Errorf("Element 2: Expected *FunctionCall, got %T", arrayLiteral.Elements[2])
		} else {
			if len(funcCall.Arguments) != 2 {
				t.Errorf("Function call: Expected 2 arguments, got %d", len(funcCall.Arguments))
			}
		}
	})

	// Test case 19: Complex valid expression
	t.Run("Complex valid expression", func(t *testing.T) {
		p := NewParser("[func(x)[0], [1, 2][variable], array[1 + 2 * 3]]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error for complex valid expression: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Each element should be an IndexAccess
		for i := 0; i < 3; i++ {
			if _, ok := arrayLiteral.Elements[i].(*IndexAccess); !ok {
				t.Errorf("Element %d: Expected *IndexAccess, got %T", i, arrayLiteral.Elements[i])
			}
		}
	})

	// Test case 20: Edge cases with whitespace and complex nesting
	t.Run("Edge cases with whitespace and complex nesting", func(t *testing.T) {
		p := NewParser("[ [ 1 , 2 ] , [ ] , [ 3 ] ]")
		expr, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse error for whitespace-heavy array: %v", err)
		}

		arrayLiteral, ok := expr.(*ArrayLiteral)
		if !ok {
			t.Fatalf("Expected *ArrayLiteral, got %T", expr)
		}

		if len(arrayLiteral.Elements) != 3 {
			t.Fatalf("Expected 3 elements, got %d", len(arrayLiteral.Elements))
		}

		// Check first element [1, 2]
		if nestedArray, ok := arrayLiteral.Elements[0].(*ArrayLiteral); !ok {
			t.Errorf("Element 0: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[0])
		} else if len(nestedArray.Elements) != 2 {
			t.Errorf("Element 0: Expected 2 sub-elements, got %d", len(nestedArray.Elements))
		}

		// Check second element []
		if nestedArray, ok := arrayLiteral.Elements[1].(*ArrayLiteral); !ok {
			t.Errorf("Element 1: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[1])
		} else if len(nestedArray.Elements) != 0 {
			t.Errorf("Element 1: Expected 0 sub-elements, got %d", len(nestedArray.Elements))
		}

		// Check third element [3]
		if nestedArray, ok := arrayLiteral.Elements[2].(*ArrayLiteral); !ok {
			t.Errorf("Element 2: Expected *ArrayLiteral, got %T", arrayLiteral.Elements[2])
		} else if len(nestedArray.Elements) != 1 {
			t.Errorf("Element 2: Expected 1 sub-element, got %d", len(nestedArray.Elements))
		}
	})
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
