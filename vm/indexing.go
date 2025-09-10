package vm

import (
	"fmt"
	"reflect"
	"strconv"
)

func (vm *VM) executeIndexExpression(left, index any, optional bool) error {
	if left == nil {
		if optional {
			return vm.Push(nil)
		}
		return fmt.Errorf("cannot index a null value")
	}

	switch typedLeft := left.(type) {
	case []any:
		return vm.executeArrayIndex(typedLeft, index)
	case map[string]any:
		return vm.executeObjectKey(typedLeft, index)
	case string:
		return vm.executeStringIndex(typedLeft, index)
	default:
		return fmt.Errorf("invalid type for index: %s", reflect.TypeOf(left).String())
	}
}

func (vm *VM) executeArrayIndex(array []any, index any) error {
	idxVal, ok := index.(float64)
	if !ok {
		return fmt.Errorf("array index must be a number, got %s", reflect.TypeOf(index).String())
	}

	intIdx := int(idxVal)
	if float64(intIdx) != idxVal {
		return fmt.Errorf("array index must be an integer, got %f", idxVal)
	}

	max := len(array)
	if intIdx < 0 {
		intIdx = max + intIdx
	}

	if intIdx < 0 || intIdx >= max {
		return fmt.Errorf("array index out of bounds: %d", intIdx)
	}

	return vm.Push(array[intIdx])
}

func (vm *VM) executeObjectKey(obj map[string]any, key any) error {
	var keyStr string
	switch v := key.(type) {
	case string:
		keyStr = v
	case float64:
		keyStr = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		keyStr = strconv.FormatBool(v)
	default:
		return fmt.Errorf("object key must be a string, number, or boolean, got %s", reflect.TypeOf(key).String())
	}

	val, exists := obj[keyStr]
	if !exists {
		return fmt.Errorf("key not found in object: %s", keyStr)
	}

	return vm.Push(val)
}

func (vm *VM) executeStringIndex(str string, index any) error {
	idxVal, ok := index.(float64)
	if !ok {
		return fmt.Errorf("string index must be a number, got %s", reflect.TypeOf(index).String())
	}

	intIdx := int(idxVal)
	if float64(intIdx) != idxVal {
		return fmt.Errorf("string index must be an integer, got %f", idxVal)
	}

	runes := []rune(str)
	max := len(runes)
	if intIdx < 0 {
		intIdx = max + intIdx
	}

	if intIdx < 0 || intIdx >= max {
		return fmt.Errorf("string index out of bounds: %d", intIdx)
	}

	return vm.Push(string(runes[intIdx]))
}
