package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
)

func (vm *VM) setBaseInstructions(bytecode *compiler.ByteCode) {
	vm.constants = bytecode.Constants
	vm.contextVars = bytecode.ContextVars
	vm.systemVars = bytecode.SystemVars

	mainFrame := NewFrame(bytecode.Instructions, 0)
	frames := make([]*Frame, MaxFrames)
	pipeScopes := make([]map[string]any, 0)
	stack := make([]any, StackSize)
	aliasVars := make(map[string]any)

	frames[0] = mainFrame
	vm.frames = frames
	vm.framesIdx = 1
	vm.pipeScopes = pipeScopes
	vm.stack = stack
	vm.aliasVars = aliasVars
}

func (vm *VM) run() error {
	frame := vm.currentFrame()
	for frame.ip < len(frame.instructions) {
		opcode := code.Opcode(frame.instructions[frame.ip])
		switch opcode {
		case code.OpConstant:
			constIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			err := vm.Push(vm.constants[constIndex])
			if err != nil {
				return err
			}
			frame.ip += 3
		case code.OpContextVar:
			varIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			err := vm.Push(vm.contextVars[varIndex])
			if err != nil {
				return err
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod, code.OpPow, code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight,
			code.OpLogicalAnd, code.OpLogicalOr:
			right := vm.Pop()
			left := vm.Pop()
			err := vm.executeBinaryExpression(opcode, left, right)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			right := vm.Pop()
			left := vm.Pop()
			err := vm.executeComparisonOperation(opcode, left, right)
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpMinus, code.OpBang:
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
			value := vm.Pop()
			if isTruthy(value) {
				// keep the value as the result of the chain
				if err := vm.Push(value); err != nil {
					return err
				}
				frame.ip = int(pos)
			} else {
				frame.ip += 3
			}
		case code.OpJumpIfFalsy:
			pos := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			value := vm.Pop()
			if !isTruthy(value) {
				// Push the falsy value as the result of the chain
				if err := vm.Push(value); err != nil {
					return err
				}
				frame.ip = int(pos)
			} else {
				frame.ip += 3
			}
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
			index := vm.Pop()
			array := vm.Pop()
			err := vm.executeIndex(array, index)
			if err != nil {
				return err
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

			pipeType := vm.constants[pipeTypeIdx].(string)
			alias := vm.systemVars[aliasIdx].(string)
			block := vm.constants[blockIdx] // Should be *compiler.InstructionBlock or nil

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
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, frame.ip)
		}

	}
	return nil
}

func (vm *VM) Run(bytecode *compiler.ByteCode) (any, error) {
	vm.setBaseInstructions(bytecode)
	err := vm.run()
	if err != nil {
		return nil, err
	}
	return vm.LastPoppedStackElem(), nil
}
