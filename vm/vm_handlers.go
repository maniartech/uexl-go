package vm

import (
	"fmt"
	"math"
	"strings"

	"github.com/maniartech/uexl/code"
)

func (vm *VM) getContextValue(name string) (any, error) {
	if vm.contextVarsValues == nil {
		return nil, fmt.Errorf("context variables not set")
	}
	value, exists := vm.contextVarsValues[name]
	if !exists {
		return nil, fmt.Errorf("context variable %q not found", name)
	}
	return value, nil
}

// executeBinaryExpressionValues evaluates binary expressions using Value types directly (zero-alloc)
// This is the optimized dispatch path — mirrors executeComparisonOperationValues
func (vm *VM) executeBinaryExpressionValues(operator code.Opcode, left, right Value) error {
	// Fast path: same types (most common)
	if left.Typ == right.Typ {
		switch left.Typ {
		case TypeFloat:
			return vm.executeNumberArithmetic(operator, left.FloatVal, right.FloatVal)
		case TypeString:
			if operator == code.OpAdd {
				return vm.executeStringAddition(left.StrVal, right.StrVal)
			}
			return vm.executeStringBinaryOperation(operator, left.StrVal, right.StrVal)
		case TypeBool:
			return vm.executeBooleanBinaryOperation(operator, left.BoolVal, right.BoolVal)
		}
	}

	// Mixed types or TypeAny — fall back to any-based dispatch
	return vm.executeBinaryExpression(operator, left.ToAny(), right.ToAny())
}

// executeBinaryExpression evaluates the binary expression using any types (legacy fallback for mixed/complex types)
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
		return vm.executeNumberArithmetic(operator, l, r)
	case string:
		r, ok := right.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", right)
		}
		// Type-specific dispatch for string operations
		if operator == code.OpAdd {
			return vm.executeStringAddition(leftVal, r)
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

// executeBinaryArithmeticOperation evaluates the binary arithmetic expression with optimized fast paths
// executeNumberArithmetic handles arithmetic operations with type-specific parameters
// This eliminates interface conversion overhead by accepting float64 directly
func (vm *VM) executeNumberArithmetic(operator code.Opcode, left, right float64) error {
	// Fast path for common operations without NaN/Inf checks
	switch operator {
	case code.OpAdd:
		return vm.pushFloat64(left + right)
	case code.OpSub:
		return vm.pushFloat64(left - right)
	case code.OpMul:
		return vm.pushFloat64(left * right)
	case code.OpDiv:
		// Zero check optimization - most divisions are non-zero
		if right == 0 {
			return fmt.Errorf("division by zero")
		}
		return vm.pushFloat64(left / right)
	}

	// Expensive checks only for operations that need them
	isNanOrInf := math.IsNaN(left) || math.IsInf(left, 0) || math.IsNaN(right) || math.IsInf(right, 0)
	isBitwiseOp := operator == code.OpBitwiseAnd || operator == code.OpBitwiseOr || operator == code.OpBitwiseXor || operator == code.OpShiftLeft || operator == code.OpShiftRight

	if isBitwiseOp && isNanOrInf {
		return fmt.Errorf("bitwise requires finite integers")
	}

	switch operator {
	case code.OpPow:
		// Optimized power operation with special case handling
		if left == 1 && isNanOrInf {
			return vm.pushFloat64(math.NaN())
		}
		return vm.pushFloat64(math.Pow(left, right))
	case code.OpMod:
		return vm.pushFloat64(math.Mod(left, right))
	// Bitwise operations with fast integer path
	case code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight:
		// Fast path: check if values are already integers
		if left == math.Trunc(left) && right == math.Trunc(right) {
			l := int64(left)
			r := int64(right)
			switch operator {
			case code.OpBitwiseAnd:
				return vm.pushFloat64(float64(l & r))
			case code.OpBitwiseOr:
				return vm.pushFloat64(float64(l | r))
			case code.OpBitwiseXor:
				return vm.pushFloat64(float64(l ^ r))
			case code.OpShiftLeft:
				if r < 0 || r >= 64 {
					return fmt.Errorf("shift count %d out of range [0, 63]", r)
				}
				return vm.pushFloat64(float64(l << uint(r)))
			case code.OpShiftRight:
				if r < 0 || r >= 64 {
					return fmt.Errorf("shift count %d out of range [0, 63]", r)
				}
				return vm.pushFloat64(float64(l >> uint(r)))
			}
		}
		return fmt.Errorf("bitwise operations require integerish operands (no decimals), got %v and %v", left, right)
	default:
		return fmt.Errorf("unknown arithmetic operator: %v", operator)
	}
}

// executeStringAddition handles string concatenation with type-specific parameters
// This eliminates interface conversion overhead by accepting string directly
func (vm *VM) executeStringAddition(left, right string) error {
	// Direct string concatenation without interface boxing
	return vm.pushString(left + right)
}

func (vm *VM) executeStringBinaryOperation(operator code.Opcode, left, right any) error {
	switch operator {
	case code.OpAdd:
		l, lok := left.(string)
		r, rok := right.(string)
		if !lok || !rok {
			return fmt.Errorf("string addition requires string operands, got %T and %T", left, right)
		}

		// For now, use simple concatenation.
		// TODO: Future optimization could use strings.Builder for concatenation chains
		result := l + r
		return vm.Push(result)
	default:
		return fmt.Errorf("unsupported string operation: %s", operator.String())
	}
}

func (vm *VM) executeBooleanBinaryOperation(operator code.Opcode, left, right bool) error {
	switch operator {
	case code.OpLogicalAnd:
		return vm.pushBool(left && right)
	case code.OpLogicalOr:
		return vm.pushBool(left || right)
	default:
		return fmt.Errorf("unsupported boolean operation: %s", operator.String())
	}
}

// executeNumberComparisonOperation optimized for maximum performance
func (vm *VM) executeNumberComparisonOperation(operator code.Opcode, left, right float64) error {
	// Optimized switch with inline operations - uses pushBool for zero allocation
	switch operator {
	case code.OpEqual:
		return vm.pushBool(left == right)
	case code.OpNotEqual:
		return vm.pushBool(left != right)
	case code.OpGreaterThan:
		return vm.pushBool(left > right)
	case code.OpGreaterThanOrEqual:
		return vm.pushBool(left >= right)
	default:
		return fmt.Errorf("unknown comparison operator: %v", operator)
	}
}

func (vm *VM) executeStringComparisonOperation(operator code.Opcode, left, right string) error {
	switch operator {
	case code.OpEqual:
		return vm.pushBool(left == right)
	case code.OpNotEqual:
		return vm.pushBool(left != right)
	default:
		return fmt.Errorf("unknown string comparison operator: %v", operator)
	}
}

func (vm *VM) executeBooleanComparisonOperation(operator code.Opcode, left, right bool) error {
	switch operator {
	case code.OpEqual:
		return vm.pushBool(left == right)
	case code.OpNotEqual:
		return vm.pushBool(left != right)
	default:
		return fmt.Errorf("unknown boolean comparison operator: %v", operator)
	}
}

// executeUnaryExpressionValue dispatches unary operations using Value types directly (zero-alloc)
func (vm *VM) executeUnaryExpressionValue(operator code.Opcode, operand Value) error {
	switch operator {
	case code.OpMinus:
		if operand.Typ == TypeFloat {
			return vm.pushFloat64(-operand.FloatVal)
		}
		return vm.executeUnaryMinusOperation(operand.ToAny())
	case code.OpBang:
		if operand.Typ == TypeBool {
			return vm.pushBool(!operand.BoolVal)
		}
		// For non-bool, use truthiness check on Value directly
		return vm.pushBool(!isTruthyValue(operand))
	case code.OpBitwiseNot:
		if operand.Typ == TypeFloat {
			v := operand.FloatVal
			if v != float64(int64(v)) {
				return fmt.Errorf("bitwise operations require integerish operands (no decimals), got %v", v)
			}
			return vm.pushFloat64(float64(^int64(v)))
		}
		return vm.executeUnaryBitwiseNotOperation(operand.ToAny())
	default:
		return fmt.Errorf("unknown unary operator: %v", operator)
	}
}

func (vm *VM) executeUnaryExpression(operator code.Opcode, operand any) error {
	switch operator {
	case code.OpMinus:
		return vm.executeUnaryMinusOperation(operand)
	case code.OpBang:
		return vm.executeUnaryBangOperation(operand)
	case code.OpBitwiseNot:
		return vm.executeUnaryBitwiseNotOperation(operand)
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

func (vm *VM) executeUnaryBitwiseNotOperation(operand any) error {
	// Extract float64 value
	var value float64
	switch v := operand.(type) {
	case float64:
		value = v
	case int:
		value = float64(v)
	default:
		return fmt.Errorf("bitwise NOT requires numeric operand, got %T", operand)
	}

	// Validate integerish (no decimals)
	if value != float64(int64(value)) {
		return fmt.Errorf("bitwise operations require integerish operands (no decimals), got %v", value)
	}

	// Perform bitwise NOT and push result (zero allocations)
	return vm.pushFloat64(float64(^int64(value)))
}

// executeComparisonOperationValues - zero-alloc comparison using Value types
func (vm *VM) executeComparisonOperationValues(operator code.Opcode, left, right Value) error {
	// Fast path: same types (most common)
	if left.Typ == right.Typ {
		switch left.Typ {
		case TypeFloat:
			return vm.executeNumberComparisonOperation(operator, left.FloatVal, right.FloatVal)
		case TypeString:
			return vm.executeStringComparisonOperation(operator, left.StrVal, right.StrVal)
		case TypeBool:
			return vm.executeBooleanComparisonOperation(operator, left.BoolVal, right.BoolVal)
		case TypeNull:
			// null == null is true, null != null is false
			switch operator {
			case code.OpEqual:
				return vm.pushBoolValue(true)
			case code.OpNotEqual:
				return vm.pushBoolValue(false)
			default:
				return fmt.Errorf("cannot compare null values with %v", operator)
			}
		}
	}

	// Mixed types or TypeAny - fall back to any comparison
	return vm.executeComparisonOperation(operator, left.ToAny(), right.ToAny())
}

func (vm *VM) executeComparisonOperation(operator code.Opcode, left, right any) error {
	switch l := left.(type) {
	case float64:
		r, ok := right.(float64)
		if !ok {
			return fmt.Errorf("number comparison requires float64 operands, got %T and %T", left, right)
		}
		return vm.executeNumberComparisonOperation(operator, l, r)
	case string:
		r, ok := right.(string)
		if !ok {
			return fmt.Errorf("string comparison requires string operands, got %T and %T", left, right)
		}
		return vm.executeStringComparisonOperation(operator, l, r)
	case bool:
		r, ok := right.(bool)
		if !ok {
			return fmt.Errorf("boolean comparison requires bool operands, got %T and %T", left, right)
		}
		return vm.executeBooleanComparisonOperation(operator, l, r)
	default:
		return fmt.Errorf("unsupported comparison for type: %T", left)
	}
}

func (vm *VM) buildArray(length int) []any {
	startIndex := vm.sp - length
	elements := make([]any, length)
	for i := 0; i < length; i++ {
		elements[i] = vm.stack[startIndex+i].ToAny()
	}
	vm.sp = startIndex
	return elements
}

func (vm *VM) buildObject(startIndex, endIndex int) (map[string]any, error) {
	object := make(map[string]any)
	for i := startIndex; i < endIndex; i += 2 {
		keyVal := vm.stack[i].ToAny()
		key, ok := keyVal.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got %T", keyVal)
		}
		object[key] = vm.stack[i+1].ToAny()
	}
	vm.sp = startIndex
	return object, nil
}

func (vm *VM) executeIndex(operand, index any) error {
	switch arr := operand.(type) {
	case []any, string:
		return vm.executeIndexValue(arr, index)
	case map[string]any:
		return vm.executeMapIndexAccess(arr, index)
	case nil:
		return fmt.Errorf("cannot index nil")
	}
	return fmt.Errorf("indexing not supported for %T", operand)
}

func (vm *VM) executeMemberAccess(container, index any) error {
	switch v := container.(type) {
	case map[string]any:
		return vm.executeMapIndexAccess(v, index)
	case []any, string:
		return vm.executeIndexValue(v, index)
	case nil:
		return fmt.Errorf("cannot access member of nil")
	}
	return fmt.Errorf("member access not supported for %T", container)
}

func (vm *VM) executeIndexValue(target any, index any) error {
	var idx int
	switch v := index.(type) {
	case float64:
		idx = int(v)
	case int:
		idx = v
	default:
		return fmt.Errorf("array index must be int, got %T", index)
	}

	switch v := target.(type) {
	case []any:
		if idx < 0 || idx >= len(v) {
			return fmt.Errorf("array index out of bounds: %d", idx)
		}
		return vm.Push(v[idx])
	case string:
		if idx < 0 || idx >= len(v) {
			return fmt.Errorf("string index out of bounds: %d", idx)
		}
		return vm.Push(string(v[idx]))
	default:
		return fmt.Errorf("unsupported target type for indexing: %T", target)
	}
}

func (vm *VM) executeMapIndexAccess(container, index any) error {
	key, ok := index.(string)
	if !ok {
		return fmt.Errorf("object key must be string, got %T", index)
	}
	if container == nil {
		return fmt.Errorf("cannot access property of nil")
	}
	value, exists := container.(map[string]any)[key]
	if !exists {
		return fmt.Errorf("key %q not found in object", key)
	}
	return vm.Push(value)
}
func (vm *VM) callFunction(funcIndex, numArgs uint16) error {

	if int(funcIndex) < 0 || int(funcIndex) >= len(vm.constants) {
		return fmt.Errorf("function index out of bounds: %d", funcIndex)
	}

	funcNameVal := vm.constants[funcIndex]
	functionName, ok := funcNameVal.AsString()
	if !ok {
		return fmt.Errorf("function name at constant index %d is not a string, got %v", funcIndex, funcNameVal.ToAny())
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

// isTruthyValue - zero-alloc truthiness check for Value types
func isTruthyValue(val Value) bool {
	switch val.Typ {
	case TypeBool:
		return val.BoolVal
	case TypeFloat:
		return val.FloatVal != 0
	case TypeString:
		return val.StrVal != ""
	case TypeNull:
		return false
	case TypeAny:
		return isTruthy(val.AnyVal)
	default:
		return false
	}
}

func isNullish(val any) bool {
	return val == nil
}

// executeStringConcat efficiently concatenates multiple strings from the stack
func (vm *VM) executeStringConcat(count int) error {
	if count < 2 {
		return fmt.Errorf("string concatenation requires at least 2 operands")
	}

	// Fast path for 2-string concatenation (most common case)
	if count == 2 {
		right := vm.popValue()
		left := vm.popValue()

		if left.Typ == TypeString && right.Typ == TypeString {
			return vm.pushString(left.StrVal + right.StrVal)
		}

		// One or both need conversion
		var leftStr, rightStr string
		if left.Typ == TypeString {
			leftStr = left.StrVal
		} else {
			leftStr = fmt.Sprintf("%v", left.ToAny())
		}
		if right.Typ == TypeString {
			rightStr = right.StrVal
		} else {
			rightStr = fmt.Sprintf("%v", right.ToAny())
		}
		return vm.pushString(leftStr + rightStr)
	}

	// General case for 3+ strings
	// Calculate total length and convert to strings in one pass
	var totalLen int
	stringValues := make([]string, count)

	for i := count - 1; i >= 0; i-- {
		val := vm.popValue()
		if val.Typ == TypeString {
			stringValues[i] = val.StrVal
			totalLen += len(val.StrVal)
		} else {
			// Convert non-string to string
			str := fmt.Sprintf("%v", val.ToAny())
			stringValues[i] = str
			totalLen += len(str)
		}
	}

	// Use strings.Builder for efficient concatenation
	var builder strings.Builder
	builder.Grow(totalLen) // Pre-allocate exact capacity
	for _, str := range stringValues {
		builder.WriteString(str)
	}

	return vm.pushString(builder.String())
}

// executeStringPatternMatch optimizes string comparison patterns like:
// variable == "prefix" + dynamic + "suffix"
// Stack layout: [target, prefix, middle, suffix]
// Uses zero-allocation pattern matching without string concatenation
func (vm *VM) executeStringPatternMatch(prefixIdx, suffixIdx int) error {
	const requiredStackElements = 4 // target, prefix, middle, suffix

	// Validate stack has sufficient elements
	if len(vm.stack) < requiredStackElements {
		return fmt.Errorf("insufficient stack elements for string pattern match: need %d, have %d",
			requiredStackElements, len(vm.stack))
	}

	// Stack has: [target, prefix, middle, suffix] from bottom to top
	// Pop in reverse order
	suffixVal := vm.popValue()
	middleVal := vm.popValue()
	prefixVal := vm.popValue()
	targetVal := vm.popValue()

	// Convert all values to strings (fast path for strings)
	var targetStr, prefixStr, middleStr, suffixStr string

	if targetVal.Typ == TypeString {
		targetStr = targetVal.StrVal
	} else {
		targetStr = fmt.Sprintf("%v", targetVal.ToAny())
	}

	if prefixVal.Typ == TypeString {
		prefixStr = prefixVal.StrVal
	} else {
		prefixStr = fmt.Sprintf("%v", prefixVal.ToAny())
	}

	if middleVal.Typ == TypeString {
		middleStr = middleVal.StrVal
	} else {
		middleStr = fmt.Sprintf("%v", middleVal.ToAny())
	}

	if suffixVal.Typ == TypeString {
		suffixStr = suffixVal.StrVal
	} else {
		suffixStr = fmt.Sprintf("%v", suffixVal.ToAny())
	}

	// Zero-allocation pattern matching
	// Check total length first
	expectedLen := len(prefixStr) + len(middleStr) + len(suffixStr)
	if len(targetStr) != expectedLen {
		return vm.pushBool(false)
	}

	// Check prefix match (length already validated above)
	if len(prefixStr) > 0 && targetStr[:len(prefixStr)] != prefixStr {
		return vm.pushBool(false)
	}

	// Check suffix match (length already validated above)
	if len(suffixStr) > 0 && targetStr[len(targetStr)-len(suffixStr):] != suffixStr {
		return vm.pushBool(false)
	}

	// Check middle match
	middleStart := len(prefixStr)
	middleEnd := len(targetStr) - len(suffixStr)
	if middleStart <= middleEnd {
		actualMiddle := targetStr[middleStart:middleEnd]
		result := actualMiddle == middleStr
		return vm.pushBool(result)
	}

	return vm.pushBool(false)
}
