package vm

import (
	"fmt"
	"reflect"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
)

func (vm *VM) setBaseInstructions(bytecode *compiler.ByteCode, contextVarsValues map[string]any) {
	vm.constants = bytecode.Constants
	vm.contextVars = bytecode.ContextVars
	vm.systemVars = bytecode.SystemVars

	// Pre-resolve context variables into a lookup slice for O(1) access
	// Optimization: Only rebuild cache if context values pointer changed or cache size mismatches
	// Compare map pointers using reflect to avoid rebuilding cache when same map is reused
	var lastPtr, newPtr uintptr
	if vm.lastContextValues != nil {
		lastPtr = reflect.ValueOf(vm.lastContextValues).Pointer()
	}
	if contextVarsValues != nil {
		newPtr = reflect.ValueOf(contextVarsValues).Pointer()
	}
	contextValuesChanged := lastPtr != newPtr || len(vm.contextVarCache) != len(vm.contextVars)
	vm.contextVarsValues = contextVarsValues
	vm.lastContextValues = contextVarsValues

	if len(vm.contextVars) > 0 && contextValuesChanged {
		if vm.contextVarCache == nil || cap(vm.contextVarCache) < len(vm.contextVars) {
			vm.contextVarCache = make([]Value, len(vm.contextVars))
		} else {
			vm.contextVarCache = vm.contextVarCache[:len(vm.contextVars)]
		}

		for i, varName := range vm.contextVars {
			if value, exists := contextVarsValues[varName]; exists {
				// Store actual value as Value (can be nil, which is valid)
				vm.contextVarCache[i] = newAnyValue(value)
			} else {
				// Store sentinel to indicate missing variable
				vm.contextVarCache[i] = newAnyValue(contextVarNotProvided)
			}
		}
	}

	// Reset execution state
	vm.sp = 0
	vm.framesIdx = 1

	// Reuse existing frame instead of allocating
	if vm.frames[0] == nil {
		vm.frames[0] = NewFrame(bytecode.Instructions, 0)
	} else {
		vm.frames[0].instructions = bytecode.Instructions
		vm.frames[0].ip = 0
		vm.frames[0].basePointer = 0
	}

	// Clear pipe scopes (preserve capacity)
	vm.pipeScopes = vm.pipeScopes[:0]

	// Clear alias vars only if non-empty (avoid iteration cost)
	if len(vm.aliasVars) > 0 {
		// For small maps, clearing is faster than allocating new map
		// For large maps or frequent clears, Go runtime optimizes clear() (Go 1.21+)
		clear(vm.aliasVars)
	}
}

func (vm *VM) run() error {
	frame := vm.currentFrame()
	for frame.ip < len(frame.instructions) {
		opcode := code.Opcode(frame.instructions[frame.ip])
		switch opcode {
		case code.OpConstant:
			constIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			// Push Value directly from constants - zero allocation!
			err := vm.pushValue(vm.constants[constIndex])
			if err != nil {
				return err
			}
			frame.ip += 3
		case code.OpContextVar:
			varIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			// Fast path: use pre-resolved cache (O(1) array access)
			if int(varIndex) < len(vm.contextVarCache) {
				value := vm.contextVarCache[varIndex]
				// Check if sentinel value (stored as AnyVal)
				if value.IsAny() {
					if _, isMissing := value.AnyVal.(contextVarMissing); isMissing {
						// Variable was not provided in context
						return fmt.Errorf("context variable %q not found", vm.contextVars[varIndex])
					}
				}
				// Push the Value directly (zero-alloc!)
				if err := vm.pushValue(value); err != nil {
					return err
				}
			} else {
				// Fallback to map lookup (should not happen in normal execution)
				value, err := vm.getContextValue(vm.contextVars[varIndex])
				if err != nil {
					return err
				}
				if err := vm.Push(value); err != nil {
					return err
				}
			}
			frame.ip += 3
		case code.OpStore:
			varIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			value := vm.Pop()
			aliasName := vm.systemVars[varIndex].(string)
			if len(vm.pipeScopes) > 0 {
				vm.pipeScopes[0][aliasName] = value // Set in outermost scope
			} else {
				// creating a new pipe scope if none exists
				vm.pipeScopes = append(vm.pipeScopes, make(map[string]any))
				vm.pipeScopes[0][aliasName] = value
			}
			frame.ip += 3
		case code.OpIdentifier:
			identIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			ident := vm.systemVars[identIndex].(string)
			val, ok := vm.getPipeVar(ident)
			if !ok {
				return fmt.Errorf("undefined pipe variable: %s", ident)
			}
			if err := vm.Push(val); err != nil {
				return err
			}
			frame.ip += 3
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod, code.OpPow, code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight, code.OpLogicalAnd, code.OpLogicalOr:
			right := vm.Pop()
			left := vm.Pop()
			err := vm.executeBinaryExpression(opcode, left, right)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			right, left := vm.pop2Values()
			err := vm.executeComparisonOperationValues(opcode, left, right)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpMinus, code.OpBang, code.OpBitwiseNot:
			operand := vm.Pop()
			err := vm.executeUnaryExpression(opcode, operand)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpJump:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			frame.ip = int(pos)
		case code.OpJumpIfTruthy:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			value := vm.popValue()
			if isTruthyValue(value) {
				// keep the value as the result of the chain
				if err := vm.pushValue(value); err != nil {
					return err
				}
				frame.ip = int(pos)
			} else {
				frame.ip += 3
			}
		case code.OpJumpIfFalsy:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			value := vm.popValue()
			if !isTruthyValue(value) {
				// Push the falsy value as the result of the chain
				if err := vm.pushValue(value); err != nil {
					return err
				}
				frame.ip = int(pos)
			} else {
				frame.ip += 3
			}
		case code.OpJumpIfNotNullish:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			value := vm.Pop()
			if !isNullish(value) {
				if err := vm.Push(value); err != nil {
					return err
				}
				frame.ip = int(pos)
			} else {
				frame.ip += 3
			}
		case code.OpJumpIfNullish:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			if vm.sp > 0 && vm.stack[vm.sp-1].IsNull() {
				frame.ip = int(pos)
				continue
			}
			frame.ip += 3
		case code.OpPop:
			vm.Pop()
			frame.ip += 1
		case code.OpTrue:
			err := vm.Push(true)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpFalse:
			err := vm.Push(false)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpNull:
			err := vm.Push(nil)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpArray:
			length := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			array := vm.buildArray(int(length))
			err := vm.Push(array)
			if err != nil {
				return err
			}
			frame.ip += 3
		case code.OpObject:
			length := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3]) // length is already number of stack elements
			object, err := vm.buildObject(vm.sp-int(length), vm.sp)
			if err != nil {
				return err
			}
			err = vm.Push(object)
			if err != nil {
				return err
			}
			frame.ip += 3
		case code.OpIndex:
			// Bounds check to prevent index out of range panics
			if frame.ip+1 >= len(frame.instructions) {
				return fmt.Errorf("instruction pointer out of bounds: OpIndex requires 1-byte operand")
			}
			optional := frame.instructions[frame.ip+1] == 1
			index := vm.Pop()
			left := vm.Pop()
			if err := vm.executeIndexExpression(left, index, optional); err != nil {
				// TODO: Only skip errors related to nullish base when in safe mode
				// e.g. if left is not indexable type, should still raise error even in safe mode
				// Currently, all errors are skipped in safe mode
				// This requires error type checking in executeIndexExpression and executeMemberAccess
				// which is not implemented yet.
				if vm.safeMode {
					if perr := vm.Push(nil); perr != nil {
						return perr
					}
				} else {
					return err
				}
			}
			frame.ip += 2
		case code.OpSlice:
			// Bounds check to prevent index out of range panics
			if frame.ip+1 >= len(frame.instructions) {
				return fmt.Errorf("instruction pointer out of bounds: OpSlice requires 1-byte operand")
			}
			optional := frame.instructions[frame.ip+1] == 1
			step := vm.Pop()
			end := vm.Pop()
			start := vm.Pop()
			target := vm.Pop()

			if err := vm.executeSliceExpression(target, start, end, step, optional); err != nil {
				return err
			}
			frame.ip += 2
		case code.OpMemberAccess:
			prop := vm.Pop()
			target := vm.Pop()
			if err := vm.executeMemberAccess(target, prop); err != nil {
				// TODO: Only skip errors related to nullish base when in safe mode
				// e.g. if left is not indexable type, should still raise error even in safe mode
				// Currently, all errors are skipped in safe mode
				// This requires error type checking in executeIndexExpression and executeMemberAccess
				// which is not implemented yet.
				if vm.safeMode {
					if perr := vm.Push(nil); perr != nil {
						return perr
					}
				} else {
					return err
				}
			}
			frame.ip += 1
		case code.OpCallFunction:
			funcIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			numArgs := code.ReadUint16(frame.instructions[frame.ip+3 : frame.ip+5])
			err := vm.callFunction(funcIndex, numArgs)
			if err != nil {
				return err
			}
			frame.ip += 5
		case code.OpPipe:
			pipeTypeIdx := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			aliasIdx := code.ReadUint16(frame.instructions[frame.ip+3 : frame.ip+5])
			blockIdx := code.ReadUint16(frame.instructions[frame.ip+5 : frame.ip+7])

			pipeTypeVal := vm.constants[pipeTypeIdx]
			pipeType, _ := pipeTypeVal.AsString() // Extract string from Value
			alias := vm.systemVars[aliasIdx].(string)
			block := vm.constants[blockIdx].ToAny() // Convert Value to any for compatibility

			input := vm.Pop()

			handler, ok := vm.pipeHandlers[pipeType]
			if !ok {
				return fmt.Errorf("unknown pipe type: %s", pipeType)
			}
			result, err := handler(input, block, alias, vm)
			if err != nil {
				return err
			}
			if err := vm.Push(result); err != nil {
				return err
			}
			frame.ip += 7
		case code.OpSafeModeOn:
			vm.safeMode = true
			frame.ip += 1
		case code.OpSafeModeOff:
			vm.safeMode = false
			frame.ip += 1
		case code.OpStringConcat:
			count := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			err := vm.executeStringConcat(int(count))
			if err != nil {
				return err
			}
			frame.ip += 3
		case code.OpStringPatternMatch:
			prefixIdx := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			suffixIdx := code.ReadUint16(frame.instructions[frame.ip+3 : frame.ip+5])
			err := vm.executeStringPatternMatch(int(prefixIdx), int(suffixIdx))
			if err != nil {
				return err
			}
			frame.ip += 5
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, frame.ip)
		}

	}
	return nil
}

func (vm *VM) Run(bytecode *compiler.ByteCode, contextValues map[string]any) (any, error) {
	vm.setBaseInstructions(bytecode, contextValues)
	err := vm.run()
	if err != nil {
		return nil, err
	}
	return vm.LastPoppedStackElem(), nil
}
