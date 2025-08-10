package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

func (vm *VM) Run() error {
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod,
			code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight,
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
		case code.OpTrue:
			err := vm.Push(&parser.BooleanLiteral{Value: true})
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpFalse:
			err := vm.Push(&parser.BooleanLiteral{Value: false})
			if err != nil {
				return err
			}
			frame.ip += 1
		case code.OpArray:
			length := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			array := vm.buildArray(int(length))
			err := vm.Push(&parser.ArrayLiteral{Elements: array})
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
			err = vm.Push(&parser.ObjectLiteral{Properties: object})
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
		case code.OpIdentifier:
			identIndex := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			ident := vm.systemVars[identIndex].(*parser.Identifier).Name
			val, ok := vm.getPipeVar(ident)
			if !ok {
				return fmt.Errorf("undefined pipe variable: %s", ident)
			}
			if err := vm.Push(val); err != nil {
				return err
			}
			frame.ip += 3

		case code.OpPipe:
			pipeTypeIdx := code.ReadUint16(frame.instructions[frame.ip+1 : frame.ip+3])
			aliasIdx := code.ReadUint16(frame.instructions[frame.ip+3 : frame.ip+5])
			pipeType := vm.constants[pipeTypeIdx].(*parser.StringLiteral).Value
			alias := vm.constants[aliasIdx].(*parser.StringLiteral).Value

			// Pop lambda/body and input value from the stack
			lambda := vm.Pop()
			input := vm.Pop()

			handler, ok := vm.pipeHandlers[pipeType]
			if !ok {
				return fmt.Errorf("unknown pipe type: %s", pipeType)
			}
			result, err := handler(input, lambda, alias, vm)
			if err != nil {
				return err
			}
			if err := vm.Push(result); err != nil {
				return err
			}
			frame.ip += 5
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, frame.ip)
		}

	}
	return nil
}
