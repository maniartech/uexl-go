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
			vm.Push(vm.constants[constIndex])
			ip += 3
		case code.OpContextVar:
			varIndex := code.ReadUint16(ins[ip+1 : ip+3])
			vm.Push(vm.contextVars[varIndex])
			ip += 3
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod,
			code.OpBitwiseAnd, code.OpBitwiseOr, code.OpBitwiseXor, code.OpShiftLeft, code.OpShiftRight,
			code.OpLogicalAnd, code.OpLogicalOr:
			right := vm.Pop()
			left := vm.Pop()
			vm.executeBinaryExpression(opcode, left, right)
			ip += 1
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			right := vm.Pop()
			left := vm.Pop()
			vm.executeComparisonOperation(opcode, left, right)
			ip += 1
		case code.OpMinus, code.OpBang:
			operand := vm.Pop()
			vm.executeUnaryExpression(opcode, operand)
			ip += 1
		case code.OpTrue:
			vm.Push(&parser.BooleanLiteral{Value: true})
			ip += 1
		case code.OpFalse:
			vm.Push(&parser.BooleanLiteral{Value: false})
			ip += 1
		case code.OpArray:
			length := code.ReadUint16(ins[ip+1 : ip+3])
			array := vm.buildArray(int(length))
			vm.Push(&parser.ArrayLiteral{Elements: array})
			ip += 3
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, ip)
		}

	}
	return nil
}
