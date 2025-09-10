package vm

import (
	"fmt"
	"reflect"
)

func (vm *VM) executeSliceExpression(target, start, end, step any, optional bool) error {
	if target == nil {
		if optional {
			return vm.Push(nil)
		}
		return fmt.Errorf("cannot slice a null value")
	}

	switch typedTarget := target.(type) {
	case []any:
		return vm.sliceArray(typedTarget, start, end, step)
	case string:
		return vm.sliceString(typedTarget, start, end, step)
	default:
		return fmt.Errorf("invalid type for slice: %s", reflect.TypeOf(target).String())
	}
}

func (vm *VM) sliceArray(arr []any, start, end, step any) error {
	st, err := vm.parseSliceStep(step)
	if err != nil {
		return err
	}

	// Set defaults based on step direction
	var defaultStart, defaultEnd int
	if st > 0 {
		defaultStart = 0
		defaultEnd = len(arr)
	} else {
		defaultStart = len(arr) - 1
		defaultEnd = -1
	}

	s, err := vm.parseSliceIndex(start, defaultStart)
	if err != nil {
		return err
	}

	e, err := vm.parseSliceIndex(end, defaultEnd)
	if err != nil {
		return err
	}

	s = vm.adjustSliceIndex(s, len(arr))
	if e != -1 {
		e = vm.adjustSliceIndex(e, len(arr))
	}

	var result []any
	if st > 0 {
		if s >= e {
			return vm.Push([]any{})
		}
		for i := s; i < e; i += st {
			result = append(result, arr[i])
		}
	} else {
		if s <= e {
			return vm.Push([]any{})
		}
		for i := s; i > e; i += st {
			result = append(result, arr[i])
		}
	}

	return vm.Push(result)
}

func (vm *VM) sliceString(str string, start, end, step any) error {
	runes := []rune(str)
	st, err := vm.parseSliceStep(step)
	if err != nil {
		return err
	}

	// Set defaults based on step direction
	var defaultStart, defaultEnd int
	if st > 0 {
		defaultStart = 0
		defaultEnd = len(runes)
	} else {
		defaultStart = len(runes) - 1
		defaultEnd = -1
	}

	s, err := vm.parseSliceIndex(start, defaultStart)
	if err != nil {
		return err
	}

	e, err := vm.parseSliceIndex(end, defaultEnd)
	if err != nil {
		return err
	}

	s = vm.adjustSliceIndex(s, len(runes))
	if e != -1 {
		e = vm.adjustSliceIndex(e, len(runes))
	}

	var result []rune
	if st > 0 {
		if s >= e {
			return vm.Push("")
		}
		for i := s; i < e; i += st {
			result = append(result, runes[i])
		}
	} else {
		if s <= e {
			return vm.Push("")
		}
		for i := s; i > e; i += st {
			result = append(result, runes[i])
		}
	}

	return vm.Push(string(result))
}

func (vm *VM) parseSliceIndex(val any, defaultVal int) (int, error) {
	if val == nil {
		return defaultVal, nil
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("slice index must be a number, got %s", reflect.TypeOf(val).String())
	}

	intVal := int(floatVal)
	if float64(intVal) != floatVal {
		return 0, fmt.Errorf("slice index must be an integer, got %g", floatVal)
	}

	return intVal, nil
}

func (vm *VM) parseSliceStep(val any) (int, error) {
	if val == nil {
		return 1, nil
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("slice step must be a number, got %s", reflect.TypeOf(val).String())
	}

	intVal := int(floatVal)
	if float64(intVal) != floatVal {
		return 0, fmt.Errorf("slice step must be an integer, got %g", floatVal)
	}

	if intVal == 0 {
		return 0, fmt.Errorf("slice step cannot be zero")
	}

	return intVal, nil
}

func (vm *VM) adjustSliceIndex(index, length int) int {
	if index < 0 {
		index += length
	}
	if index < 0 {
		return 0
	}
	if index > length {
		return length
	}
	return index
}
