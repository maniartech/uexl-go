package vm

import (
	"fmt"
	"math"

	"github.com/maniartech/uexl_go/code"
)

// executeBinaryExpression evaluates the binary expression by popping the top two elements from the stack, applying the operator, and pushing the result back onto the stack.
func (vm *VM) executeBinaryExpression(operator code.Opcode, left, right any) error {
	switch leftVal := left.(type) {
	case float64, int:
		var l, r float64
		// Convert left operand to float64
		switch v := left.(type) {
		case float64:
			l = v
		case int:
			l = float64(v)
		default:
			return fmt.Errorf("expected number, got %T", left)
		}
		// Convert right operand to float64
		switch v := right.(type) {
		case float64:
			r = v
		case int:
			r = float64(v)
		default:
			return fmt.Errorf("expected number, got %T", right)
		}
		return vm.executeBinaryArithmeticOperation(operator, l, r)
	case string:
		r, ok := right.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", right)
		}
		return vm.executeStringBinaryOperation(operator, leftVal, r)
	case bool:
		r, ok := right.(bool)
		if !ok {
			return fmt.Errorf("expected bool, got %T", right)
		}
		return vm.executeBooleanBinaryOperation(operator, leftVal, r)
	default:
		return fmt.Errorf("unsupported binary operation for type: %T", left)
	}
}

// executeBinaryArithmeticOperation evaluates the binary arithmetic expression by popping the top two elements from the stack, applying the operator, and pushing the result back onto the stack.
func (vm *VM) executeBinaryArithmeticOperation(operator code.Opcode, left, right any) error {
	leftValue := left.(float64)
	rightValue := right.(float64)

	switch operator {
	case code.OpAdd:
		vm.Push(leftValue + rightValue)
	case code.OpSub:
		vm.Push(leftValue - rightValue)
	case code.OpMul:
		vm.Push(leftValue * rightValue)
	case code.OpDiv:
		if rightValue == 0 {
			return fmt.Errorf("division by zero")
		}
		vm.Push(leftValue / rightValue)
	case code.OpPow:
		vm.Push(math.Pow(leftValue, rightValue))
	case code.OpMod:
		vm.Push(math.Mod(leftValue, rightValue))
	// Bitwise operations
	case code.OpBitwiseAnd:
		vm.Push(float64(int(leftValue) & int(rightValue)))
	case code.OpBitwiseOr:
		vm.Push(float64(int(leftValue) | int(rightValue)))
	case code.OpBitwiseXor:
		vm.Push(float64(int(leftValue) ^ int(rightValue)))
	case code.OpShiftLeft:
		vm.Push(float64(int(leftValue) << int(rightValue)))
	case code.OpShiftRight:
		vm.Push(float64(int(leftValue) >> int(rightValue)))
	default:
		return fmt.Errorf("unknown operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeStringBinaryOperation(operator code.Opcode, left, right any) error {
	switch operator {
	case code.OpAdd:
		l, lok := left.(string)
		r, rok := right.(string)
		if !lok || !rok {
			return fmt.Errorf("string addition requires string operands, got %T and %T", left, right)
		}
		return vm.Push(l + r)
	default:
		return fmt.Errorf("unsupported string operation: %s", operator.String())
	}
}

func (vm *VM) executeBooleanBinaryOperation(operator code.Opcode, left, right bool) error {
	switch operator {
	case code.OpLogicalAnd:
		vm.Push(left && right)
	case code.OpLogicalOr:
		vm.Push(left || right)
	default:
		return fmt.Errorf("unsupported boolean operation: %s", operator.String())
	}
	return nil
}

func (vm *VM) executeNumberComparisonOperation(operator code.Opcode, left, right any) error {
	leftValue, lok := left.(float64)
	rightValue, rok := right.(float64)
	if !lok || !rok {
		return fmt.Errorf("number comparison requires float64 operands, got %T and %T", left, right)
	}
	switch operator {
	case code.OpEqual:
		vm.Push(leftValue == rightValue)
	case code.OpNotEqual:
		vm.Push(leftValue != rightValue)
	case code.OpGreaterThan:
		vm.Push(leftValue > rightValue)
	case code.OpGreaterThanOrEqual:
		vm.Push(leftValue >= rightValue)
	default:
		return fmt.Errorf("unknown comparison operator: %v", operator)
	}
	return nil
}
func (vm *VM) executeStringComparisonOperation(operator code.Opcode, left, right any) error {
	leftValue, lok := left.(string)
	rightValue, rok := right.(string)
	if !lok || !rok {
		return fmt.Errorf("string comparison requires string operands, got %T and %T", left, right)
	}
	switch operator {
	case code.OpEqual:
		vm.Push(leftValue == rightValue)
	case code.OpNotEqual:
		vm.Push(leftValue != rightValue)
	default:
		return fmt.Errorf("unknown string comparison operator: %v", operator)
	}
	return nil
}
func (vm *VM) executeBooleanComparisonOperation(operator code.Opcode, left, right any) error {
	leftValue, lok := left.(bool)
	rightValue, rok := right.(bool)
	if !lok || !rok {
		return fmt.Errorf("boolean comparison requires bool operands, got %T and %T", left, right)
	}
	switch operator {
	case code.OpEqual:
		vm.Push(leftValue == rightValue)
	case code.OpNotEqual:
		vm.Push(leftValue != rightValue)
	default:
		return fmt.Errorf("unknown boolean comparison operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeUnaryExpression(operator code.Opcode, operand any) error {
	switch operator {
	case code.OpMinus:
		return vm.executeUnaryMinusOperation(operand)
	case code.OpBang:
		return vm.executeUnaryBangOperation(operand)
	default:
		return fmt.Errorf("unknown operand type: %T", operand)
	}
}

func (vm *VM) executeUnaryMinusOperation(operand any) error {
	switch v := operand.(type) {
	case float64:
		vm.Push(-v)
	case int:
		vm.Push(float64(-v))
	default:
		return fmt.Errorf("unknown operand type: %T", operand)
	}
	return nil
}

func (vm *VM) executeUnaryBangOperation(operand any) error {
	switch v := operand.(type) {
	case bool:
		vm.Push(!v)
	default:
		// Unary Logical Not converts anything falsy to false
		vm.Push(!isTruthy(operand))
	}
	return nil
}

func (vm *VM) executeComparisonOperation(operator code.Opcode, left, right any) error {
	switch left.(type) {
	case float64:
		return vm.executeNumberComparisonOperation(operator, left, right)
	case string:
		return vm.executeStringComparisonOperation(operator, left, right)
	case bool:
		return vm.executeBooleanComparisonOperation(operator, left, right)
	default:
		return fmt.Errorf("unsupported comparison for type: %T", left)
	}
}

func (vm *VM) buildArray(length int) []any {
	startIndex := vm.sp - length
	elements := make([]any, length)
	for i := 0; i < length; i++ {
		elements[i] = vm.stack[startIndex+i]
	}
	vm.sp = startIndex
	return elements
}

func (vm *VM) buildObject(startIndex, endIndex int) (map[string]any, error) {
	object := make(map[string]any)
	for i := startIndex; i < endIndex; i += 2 {
		key, ok := vm.stack[i].(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got %T", vm.stack[i])
		}
		object[key] = vm.stack[i+1]
	}
	vm.sp = startIndex
	return object, nil
}

func (vm *VM) executeIndex(operand, index any, nullish bool) error {
	switch arr := operand.(type) {
	case []any, string:
		return vm.executeIndexValue(arr, index, nullish)
	case map[string]any:
		return vm.executeMapIndexAccess(arr, index, nullish)
	case nil:
		if nullish {
			return vm.Push(nil)
		}
	
	}
	return fmt.Errorf("indexing not supported for %T", operand)
}

func (vm *VM) executeMemberAccess(container, index any, nullish bool) error {
	switch v := container.(type) {
	case map[string]any:
		return vm.executeMapIndexAccess(v, index, nullish)
	case []any, string:
		return vm.executeIndexValue(v, index, nullish)
	case nil:
		if nullish {
			return vm.Push(nil)
		}
	}
	return fmt.Errorf("member access not supported for %T", container)
}

// Suggestion for merging array and string index functions
func (vm *VM) executeIndexValue(target any, index any, nullish bool) error {
	idx := 0
	switch v := index.(type) {
	case float64:
		idx = int(v)
	case int:
		idx = v
	default:
		return fmt.Errorf("array index must be int, got %T", index)
	}

	// Determine the type of the target and perform the appropriate indexing
	switch v := target.(type) {
	case []any:
		if idx < 0 || idx >= len(v) {
			if nullish {
				return vm.Push(nil)
			}
			return fmt.Errorf("array index out of bounds: %d", idx)

		}
		return vm.Push(v[idx])
	case string:
		if idx < 0 || idx >= len(v) {
			if nullish {
				return vm.Push(nil)
			}
			return fmt.Errorf("string index out of bounds: %d", idx)
		}
		return vm.Push(string(v[idx]))
	default:
		if nullish {
			return vm.Push(nil)
		}
		return fmt.Errorf("unsupported target type for indexing: %T", target)
	}
}

func (vm *VM) executeMapIndexAccess(container, index any, nullish bool) error {
	key, ok := index.(string)
	if !ok {
		return fmt.Errorf("object key must be string, got %T", index)
	}
	if container == nil {
		if nullish {
			return vm.Push(nil)
		}
		return fmt.Errorf("cannot access property of nil")
	}
	value, exists := container.(map[string]any)[key]
	if !exists {
		if nullish {
			return vm.Push(nil)
		}
		return fmt.Errorf("key %q not found in object", key)
	}
	return vm.Push(value)
}

func (vm *VM) callFunction(funcIndex, numArgs uint16) error {

	if int(funcIndex) < 0 || int(funcIndex) >= len(vm.constants) {
		return fmt.Errorf("function index out of bounds: %d", funcIndex)
	}

	functionName, ok := vm.constants[funcIndex].(string)
	if !ok {
		return fmt.Errorf("function name at constant index %d is not a string, got %T", funcIndex, vm.constants[funcIndex])
	}
	function, exists := vm.functionContext[functionName]
	if !exists {
		return fmt.Errorf("function %s not found in context", functionName)
	}
	args := make([]any, numArgs)
	for i := 0; i < int(numArgs); i++ {
		if vm.sp == 0 {
			return fmt.Errorf("not enough arguments on stack for function %s", functionName)
		}
		args[int(numArgs)-1-i] = vm.Pop()
	}
	functionResult, err := function(args...)
	if err != nil {
		return fmt.Errorf("error calling function %s: %w", functionName, err)
	}
	if functionResult == nil {
		return nil
	}
	return vm.Push(functionResult)
}

func isTruthy(val any) bool {
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != ""
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	default:
		return val != nil
	}
}

func normalizeFalsyToFalse(val any) any {
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		if v == 0 {
			return false
		}
	case int:
		if v == 0 {
			return false
		}
	case string:
		if v == "" {
			return false
		}
	case []any:
		if len(v) == 0 {
			return false
		}
	case map[string]any:
		if len(v) == 0 {
			return false
		}
	}
	return val
}
