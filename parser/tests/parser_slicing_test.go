package parser_test

import (
	"testing"

	"github.com/maniartech/uexl_go/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_Slicing(t *testing.T) {
	tests := []struct {
		input    string
		expected parser.Expression
	}{
		{
			input: "arr[1:4]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1},
				End:    &parser.NumberLiteral{Value: 4},
			},
		},
		{
			input: "arr[:4]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End:    &parser.NumberLiteral{Value: 4},
			},
		},
		{
			input: "arr[1:]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1},
				End:    nil,
			},
		},
		{
			input: "arr[:]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End:    nil,
			},
		},
		{
			input: "arr[1:4:2]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1},
				End:    &parser.NumberLiteral{Value: 4},
				Step:   &parser.NumberLiteral{Value: 2},
			},
		},
		{
			input: "arr[::2]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End:    nil,
				Step:   &parser.NumberLiteral{Value: 2},
			},
		},
		{
			input: "arr[1::2]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  &parser.NumberLiteral{Value: 1},
				End:    nil,
				Step:   &parser.NumberLiteral{Value: 2},
			},
		},
		{
			input: "arr[:4:2]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End:    &parser.NumberLiteral{Value: 4},
				Step:   &parser.NumberLiteral{Value: 2},
			},
		},
		{
			input: "arr[-1:]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start: &parser.UnaryExpression{
					Operator: "-",
					Operand:  &parser.NumberLiteral{Value: 1},
				},
				End: nil,
			},
		},
		{
			input: "arr[:-1]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End: &parser.UnaryExpression{
					Operator: "-",
					Operand:  &parser.NumberLiteral{Value: 1},
				},
			},
		},
		{
			input: "arr[-2:-1]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start: &parser.UnaryExpression{
					Operator: "-",
					Operand:  &parser.NumberLiteral{Value: 2},
				},
				End: &parser.UnaryExpression{
					Operator: "-",
					Operand:  &parser.NumberLiteral{Value: 1},
				},
			},
		},
		{
			input: "arr[::-1]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start:  nil,
				End:    nil,
				Step: &parser.UnaryExpression{
					Operator: "-",
					Operand:  &parser.NumberLiteral{Value: 1},
				},
			},
		},
		{
			input: "arr[a():b(1):c*2]",
			expected: &parser.SliceExpression{
				Target: &parser.Identifier{Name: "arr"},
				Start: &parser.FunctionCall{
					Function:  &parser.Identifier{Name: "a"},
					Arguments: []parser.Expression{},
				},
				End: &parser.FunctionCall{
					Function: &parser.Identifier{Name: "b"},
					Arguments: []parser.Expression{
						&parser.NumberLiteral{Value: 1},
					},
				},
				Step: &parser.BinaryExpression{
					Left:     &parser.Identifier{Name: "c"},
					Operator: "*",
					Right:    &parser.NumberLiteral{Value: 2},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.NewParser(tt.input)
			expr, err := p.Parse()

			assert.NoError(t, err)
			assert.NotNil(t, expr)

			// Function to compare expressions, ignoring line and column numbers
			var compareExpr func(expected, actual parser.Expression)
			compareExpr = func(expected, actual parser.Expression) {
				if expected == nil {
					assert.Nil(t, actual)
					return
				}
				assert.NotNil(t, actual)
				assert.Equal(t, expected.Type(), actual.Type())

				switch exp := expected.(type) {
				case *parser.SliceExpression:
					act := actual.(*parser.SliceExpression)
					compareExpr(exp.Target, act.Target)
					compareExpr(exp.Start, act.Start)
					compareExpr(exp.End, act.End)
					compareExpr(exp.Step, act.Step)
				case *parser.Identifier:
					act := actual.(*parser.Identifier)
					assert.Equal(t, exp.Name, act.Name)
				case *parser.NumberLiteral:
					act := actual.(*parser.NumberLiteral)
					assert.Equal(t, exp.Value, act.Value)
				case *parser.UnaryExpression:
					act := actual.(*parser.UnaryExpression)
					assert.Equal(t, exp.Operator, act.Operator)
					compareExpr(exp.Operand, act.Operand)
				case *parser.BinaryExpression:
					act := actual.(*parser.BinaryExpression)
					assert.Equal(t, exp.Operator, act.Operator)
					compareExpr(exp.Left, act.Left)
					compareExpr(exp.Right, act.Right)
				case *parser.FunctionCall:
					act := actual.(*parser.FunctionCall)
					compareExpr(exp.Function, act.Function)
					assert.Equal(t, len(exp.Arguments), len(act.Arguments))
					for i := range exp.Arguments {
						compareExpr(exp.Arguments[i], act.Arguments[i])
					}
				default:
					t.Fatalf("unhandled expression type: %T", exp)
				}
			}

			compareExpr(tt.expected, expr)
		})
	}
}
