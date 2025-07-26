package vm

import (
	"fmt"
	"math"

	"github.com/maniartech/uexl_go/code"
	"github.com/maniartech/uexl_go/compiler"
	"github.com/maniartech/uexl_go/parser"
)

const StackSize = 1024

var True = parser.BooleanLiteral{Value: true}
var False = parser.BooleanLiteral{Value: false}
var Null = parser.NullLiteral{}

type VM struct {
	constants    []parser.Node
	contextVars  []parser.Node
	instructions code.Instructions
	stack        []parser.Node
	sp           int
}

func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		contextVars:  bytecode.ContextVars,
		instructions: bytecode.Instructions,
		stack:        make([]parser.Node, StackSize),
		sp:           0,
	}
}
func (vm *VM) Push(node parser.Node) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = node
	vm.sp++
	return nil
}

func (vm *VM) Pop() parser.Node {
	if vm.sp == 0 {
		return nil
	}
	vm.sp--
	node := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil // Clear the reference
	return node
}

func (vm *VM) LastPoppedStackElem() parser.Node {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Top() parser.Node {
	// Check if the stack is empty
	if vm.sp == 0 {
		return nil
	}
	// Return the top element without removing it
	return vm.stack[vm.sp-1]
}

// BinaryExpression evaluates the binary expression by popping the top two elements from the stack, applying the operator, and pushing the result back onto the stack.
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
	default:
		return fmt.Errorf("unknown operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeNumberComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	leftValue := left.(*parser.NumberLiteral).Value
	rightValue := right.(*parser.NumberLiteral).Value
	switch operator {
	case code.OpEqual:
		if leftValue == rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	case code.OpNotEqual:
		if leftValue != rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	case code.OpGreaterThan:
		if leftValue > rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	case code.OpGreaterThanOrEqual:
		if leftValue >= rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
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
		if leftValue == rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	case code.OpNotEqual:
		if leftValue != rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	default:
		return fmt.Errorf("unknown string comparison operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeBooleanComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	leftValue := left.(*parser.BooleanLiteral).Value
	rightValue := right.(*parser.BooleanLiteral).Value
	switch operator {
	case code.OpEqual:
		if leftValue == rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	case code.OpNotEqual:
		if leftValue != rightValue {
			vm.Push(&True)
		} else {
			vm.Push(&False)
		}
	default:
		return fmt.Errorf("unknown boolean comparison operator: %v", operator)
	}
	return nil
}

func (vm *VM) executeComparisonOperation(operator code.Opcode, left, right parser.Node) error {
	// Check if both left and right are NumberLiteral nodes
	switch left := left.(type) {
	case *parser.NumberLiteral:
		if right, ok := right.(*parser.NumberLiteral); ok {
			return vm.executeNumberComparisonOperation(operator, left, right)
		}
	case *parser.StringLiteral:
		if right, ok := right.(*parser.StringLiteral); ok {
			return vm.executeStringComparisonOperation(operator, left, right)
		}
	case *parser.BooleanLiteral:
		if right, ok := right.(*parser.BooleanLiteral); ok {
			return vm.executeBooleanComparisonOperation(operator, left, right)
		}
	default:
		return fmt.Errorf("unsupported comparison operation for types: %T and %T", left, right)
	}
	return nil
}

func (vm *VM) executeUnaryArithmeticOperation(operator code.Opcode, operand parser.Node) error {
	operandValue := operand.(*parser.NumberLiteral).Value
	switch operator {
	case code.OpMinus:
		vm.Push(&parser.NumberLiteral{Value: -operandValue})
	default:
		return fmt.Errorf("unknown unary operator: %v", operator)
	}
	return nil
}

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod:
			right := vm.Pop()
			left := vm.Pop()
			vm.executeBinaryArithmeticOperation(opcode, left, right)
			ip += 1
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			right := vm.Pop()
			left := vm.Pop()
			vm.executeComparisonOperation(opcode, left, right)
		case code.OpMinus:
			operand := vm.Pop()
			vm.executeUnaryArithmeticOperation(opcode, operand)
			ip += 1
		default:
			return fmt.Errorf("unknown opcode: %v at ip=%d", opcode, ip)
		}
	}
	return nil
}
