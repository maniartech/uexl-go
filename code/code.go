package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions represents a sequence of bytecode instructions.
type Opcode byte

const (
	OpConstant Opcode = iota
	OpContextVar
	OpIdentifier
	OpStore
	OpPop
	OpAdd
	OpSub
	OpMul
	OpPow
	OpDiv
	OpMod
	OpMinus
	OpBang
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual
	OpLogicalAnd
	OpLogicalOr
	OpBitwiseAnd
	OpBitwiseOr
	OpBitwiseXor
	OpBitwiseNot
	OpShiftLeft
	OpShiftRight
	OpJump
	OpJumpIfTruthy
	OpJumpIfFalsy
	OpArray
	OpIndex
	OpObject
	OpCallFunction
	OpPipe
)

func (op Opcode) String() string {
	def, ok := definations[op]
	if !ok {
		return fmt.Sprintf("UNKNOWN(%d)", op)
	}
	return def.Name
}

type Definition struct {
	Name          string
	OperandWidths []int
}

var definations = map[Opcode]*Definition{
	OpConstant:           {"OpConstant", []int{2}},
	OpContextVar:         {"OpContextVar", []int{2}},
	OpIdentifier:         {"OpIdentifier", []int{2}},
	OpStore:              {"OpStore", []int{2}},
	OpPop:                {"OpPop", []int{}},
	OpAdd:                {"OpAdd", []int{}},
	OpSub:                {"OpSub", []int{}},
	OpMul:                {"OpMul", []int{}},
	OpPow:                {"OpPow", []int{}},
	OpDiv:                {"OpDiv", []int{}},
	OpMod:                {"OpMod", []int{}},
	OpTrue:               {"OpTrue", []int{}},
	OpFalse:              {"OpFalse", []int{}},
	OpEqual:              {"OpEqual", []int{}},
	OpNotEqual:           {"OpNotEqual", []int{}},
	OpGreaterThan:        {"OpGreaterThan", []int{}},
	OpGreaterThanOrEqual: {"OpGreaterThanOrEqual", []int{}},
	OpMinus:              {"OpMinus", []int{}},
	OpBang:               {"OpBang", []int{}},
	OpLogicalAnd:         {"OpLogicalAnd", []int{}},
	OpLogicalOr:          {"OpLogicalOr", []int{}},
	OpBitwiseAnd:         {"OpBitwiseAnd", []int{}},
	OpBitwiseOr:          {"OpBitwiseOr", []int{}},
	OpBitwiseXor:         {"OpBitwiseXor", []int{}},
	OpBitwiseNot:         {"OpBitwiseNot", []int{}},
	OpShiftLeft:          {"OpShiftLeft", []int{}},
	OpShiftRight:         {"OpShiftRight", []int{}},
	OpJump:               {"OpJump", []int{2}},
	OpJumpIfTruthy:       {"OpJumpIfTruthy", []int{2}},
	OpJumpIfFalsy:        {"OpJumpIfFalsy", []int{2}},
	OpArray:              {"OpArray", []int{2}},
	OpObject:             {"OpHash", []int{2}},
	OpIndex:              {"OpIndex", []int{}},
	OpCallFunction:       {"OpCallFunction", []int{2, 2}},
	OpPipe:               {"OpPipe", []int{2, 2, 2}}, // pipeTypeIdx, aliasIdx, blockIdx
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definations[Opcode(op)]
	if !ok {
		// TODO: handle error properly
		return nil, fmt.Errorf("unknown opcode: %d", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definations[op]
	if !ok {
		// should never happen, but if it does, we panic
		panic(fmt.Sprintf("unknown opcode: %d", op))
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}
	return instruction
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		op := Opcode(ins[i])
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			i++
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s %v\n", i, op.String(), operands)
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}
func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }
