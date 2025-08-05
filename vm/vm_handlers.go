package vm

import (
	"fmt"
	"math"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

// executeBinaryExpression evaluates the binary expression by popping the top two elements from the stack, applying the operator, and pushing the result back onto the stack.
func (vm *VM) executeBinaryExpression(operator code.Opcode, left, right parser.Node) error {
	if left.Type() != right.Type() {
		return fmt.Errorf("type mismatch: cannot apply %s to %T and %T",
			operator.String(), left, right)
	}
	switch left.(type) {
	case *parser.NumberLiteral:
		return vm.executeBinaryArithmeticOperation(operator, left, right)
	case *parser.StringLiteral:
		return vm.executeStringBinaryOperation(operator, left, right)
	case *parser.BooleanLiteral:
		return vm.executeBooleanBinaryOperation(operator, left, right)
	default:
		return fmt.Errorf("unsupported binary operation for type: %T", left)
	}
}

// executeBinaryArithmeticOperation evaluates the binary arithmetic expression by popping the top two elements from the stack, applying the operator, and pushing the result back onto the stack.
func (vm *VM) executeBinaryArithmeticOperation(operator code.Opcode, left, right parser.Node) error {
	leftValue := left.(*parser.NumberLiteral).Value
	rightValue := right.(*parser.NumberLiteral).Value

	switch operator {
	case code.OpAdd:
		vm.Push(&parser.NumberLiteral{Value: leftValue + rightValue})
	case code.OpSub:
		vm.Push(&parser.NumberLiteral{Value: leftValue - rightValue})
	case code.OpMul:
		vm.Push(&parser.NumberLiteral{Value: leftValue * rightValue})
	case code.OpDiv:
		if rightValue == 0 {
			return fmt.Errorf("division by zero")
		}
		vm.Push(&parser.NumberLiteral{Value: leftValue / rightValue})
	case code.OpMod:
		vm.Push(&parser.NumberLiteral{Value: math.Mod(leftValue, rightValue)})
	// Bitwise operations
	case code.OpBitwiseAnd:
		vm.Push(&parser.NumberLiteral{Value: float64(int(leftValue) & int(rightValue))})
	case code.OpBitwiseOr:
		vm.Push(&parser.NumberLiteral{Value: float64(int(leftValue) | int(rightValue))})
	case code.OpBitwiseXor:
		vm.Push(&parser.NumberLiteral{Value: float64(int(leftValue) ^ int(rightValue))})
	case code.OpShiftLeft:
		vm.Push(&parser.NumberLiteral{Value: float64(int(leftValue) << int(rightValue))})
	case code.OpShiftRight:
		vm.Push(&parser.NumberLiteral{Value: float64(int(leftValue) >> int(rightValue))})
	default:
		return fmt.Errorf("unknown operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeStringBinaryOperation(operator code.Opcode, left, right parser.Node) error {
	switch operator {
	case code.OpAdd:
		return vm.Push(&parser.StringLiteral{
			Value: left.(*parser.StringLiteral).Value + right.(*parser.StringLiteral).Value})
	default:
		return fmt.Errorf("unsupported string operation: %s", operator.String())
	}
}

func (vm *VM) executeBooleanBinaryOperation(operator code.Opcode, left, right parser.Node) error {
	switch operator {
	case code.OpLogicalAnd:
		leftValue := left.(*parser.BooleanLiteral).Value
		rightValue := right.(*parser.BooleanLiteral).Value
		vm.Push(&parser.BooleanLiteral{Value: leftValue && rightValue})
	case code.OpLogicalOr:
		leftValue := left.(*parser.BooleanLiteral).Value
		rightValue := right.(*parser.BooleanLiteral).Value
		vm.Push(&parser.BooleanLiteral{Value: leftValue || rightValue})
	default:
		return fmt.Errorf("unsupported boolean operation: %s", operator.String())
	}
	return nil
}

func (vm *VM) executeNumberComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	leftValue := left.(*parser.NumberLiteral).Value
	rightValue := right.(*parser.NumberLiteral).Value
	switch operator {
	case code.OpEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue == rightValue})
	case code.OpNotEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue != rightValue})
	case code.OpGreaterThan:
		vm.Push(&parser.BooleanLiteral{Value: leftValue > rightValue})
	case code.OpGreaterThanOrEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue >= rightValue})
	default:
		return fmt.Errorf("unknown comparison operator: %v", operator)
	}
	return nil
}
func (vm *VM) executeStringComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	leftValue := left.(*parser.StringLiteral).Value
	rightValue := right.(*parser.StringLiteral).Value
	switch operator {
	case code.OpEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue == rightValue})
	case code.OpNotEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue != rightValue})
	default:
		return fmt.Errorf("unknown string comparison operator: %v", operator)
	}
	return nil
}
func (vm *VM) executeBooleanComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	// Check if both left and right are NumberLiteral nodes
	leftValue := left.(*parser.BooleanLiteral).Value
	rightValue := right.(*parser.BooleanLiteral).Value
	switch operator {
	case code.OpEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue == rightValue})
	case code.OpNotEqual:
		vm.Push(&parser.BooleanLiteral{Value: leftValue != rightValue})
	default:
		return fmt.Errorf("unknown boolean comparison operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeUnaryExpression(operator code.Opcode, operand parser.Node) error {
	switch operand := operand.(type) {
	case *parser.NumberLiteral:
		return vm.executeUnaryNumericOperation(operator, operand)
	case *parser.StringLiteral:
		return fmt.Errorf("unary operations not supported for strings")
	case *parser.BooleanLiteral:
		return vm.executeUnaryBooleanOperation(operator, operand)
	default:
		return fmt.Errorf("unknown operand type: %T", operand)
	}
}

func (vm *VM) executeUnaryNumericOperation(operator code.Opcode, operand parser.Node) error {
	operandValue := operand.(*parser.NumberLiteral).Value
	switch operator {
	case code.OpMinus:
		vm.Push(&parser.NumberLiteral{Value: -operandValue})
	default:
		return fmt.Errorf("unknown unary operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeUnaryBooleanOperation(operator code.Opcode, operand parser.Node) error {
	operandValue := operand.(*parser.BooleanLiteral).Value
	switch operator {
	case code.OpBang:
		vm.Push(&parser.BooleanLiteral{Value: !operandValue})
	default:
		return fmt.Errorf("unknown unary operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	if left.Type() != right.Type() {
		return fmt.Errorf("type mismatch: cannot compare %T with %T", left, right)
	}
	switch left.(type) {
	case *parser.NumberLiteral:
		return vm.executeNumberComparisonOperation(operator, left, right)
	case *parser.StringLiteral:
		return vm.executeStringComparisonOperation(operator, left, right)
	case *parser.BooleanLiteral:
		return vm.executeBooleanComparisonOperation(operator, left, right)
	default:
		return fmt.Errorf("unsupported comparison for type: %T", left)
	}
}

func (vm *VM) buildArray(length int) []parser.Expression {
	// Calculate the start index on the stack
	startIndex := vm.sp - length

	elements := make([]parser.Expression, length)
	for i := range length {
		elem, ok := vm.stack[startIndex+i].(parser.Expression)
		if !ok {
			panic(fmt.Sprintf("expected parser.Expression on stack, got %T", vm.stack[startIndex+i]))
		}
		elements[i] = elem
	}

	// Update the stack pointer to remove the elements
	vm.sp = startIndex

	return elements
}
