package vm

import (
	"fmt"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/parser"
)

func (vm *VM) Run() error {
	ip := 0
	ins := vm.instructions
	for ip < len(ins) {
		opcode := code.Opcode(ins[ip])
		switch opcode {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1 : ip+3])
			err := vm.Push(vm.constants[constIndex])
			if err != nil {
				return err
			}
			ip += 3
		case code.OpContextVar:
			varIndex := code.ReadUint16(ins[ip+1 : ip+3])
			err := vm.Push(vm.contextVars[varIndex])
			if err != nil {
				return err
			}
			ip += 3
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod,
			code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight,
			code.OpLogicalAnd, code.OpLogicalOr:
			right := vm.Pop()
			left := vm.Pop()
			err := vm.executeBinaryExpression(opcode, left, right)
			if err != nil {
				return err
			}
			ip += 1
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			right := vm.Pop()
			left := vm.Pop()
			err := vm.executeComparisonOperation(opcode, left, right)
			if err != nil {
				return err
			}
			ip += 1
		case code.OpMinus, code.OpBang:
			operand := vm.Pop()
			err := vm.executeUnaryExpression(opcode, operand)
			if err != nil {
				return err
			}
			ip += 1
		case code.OpTrue:
			err := vm.Push(&parser.BooleanLiteral{Value: true})
			if err != nil {
				return err
			}
			ip += 1
		case code.OpFalse:
			err := vm.Push(&parser.BooleanLiteral{Value: false})
			if err != nil {
				return err
			}
			ip += 1
		case code.OpArray:
			length := code.ReadUint16(ins[ip+1 : ip+3])
			array := vm.buildArray(int(length))
			err := vm.Push(&parser.ArrayLiteral{Elements: array})
			if err != nil {
				return err
			}
			ip += 3
		case code.OpObject:
			length := code.ReadUint16(ins[ip+1 : ip+3]) // length is already number of stack elements
			object, err := vm.buildObject(vm.sp-int(length), vm.sp)
			if err != nil {
				return err
			}
			err = vm.Push(&parser.ObjectLiteral{Properties: object})
			if err != nil {
				return err
			}
			ip += 3
		case code.OpIndex:
			index := vm.Pop()
			array := vm.Pop()
			err := vm.executeIndex(array, index)
			if err != nil {
				return err
			}
			ip += 1
		case code.OpCallFunction:
			funcIndex := code.ReadUint16(ins[ip+1 : ip+3])
			numArgs := code.ReadUint16(ins[ip+3 : ip+5])
			err := vm.callFunction(funcIndex, numArgs)
			if err != nil {
				return err
			}
			ip += 5
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, ip)
		}

	}
	return nil
}
