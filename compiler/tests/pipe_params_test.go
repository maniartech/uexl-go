package compiler_test

import (
	"testing"

	"github.com/maniartech/uexl/code"
	"github.com/maniartech/uexl/compiler"
	"github.com/stretchr/testify/assert"
)

// compileExpr is a helper that parses and compiles an expression, returning the bytecode.
func compileExpr(t *testing.T, input string) *compiler.ByteCode {
	t.Helper()
	node := parse(input)
	comp := compiler.New()
	err := comp.Compile(node)
	if err != nil {
		t.Fatalf("compile error for %s: %v", input, err)
	}
	return comp.ByteCode()
}

// findOpPipeInstructions locates the OpPipe instruction at the given byte offset in the bytecode
// and returns its four operands: [pipeTypeIdx, aliasIdx, blockIdx, argsIdx].
func findOpPipeOperands(t *testing.T, ins code.Instructions) [4]int {
	t.Helper()
	for i := 0; i < len(ins); {
		op := code.Opcode(ins[i])
		def, err := code.Lookup(ins[i])
		if err != nil {
			t.Fatalf("unknown opcode at %d: %v", i, err)
		}
		if op == code.OpPipe {
			if len(def.OperandWidths) != 4 {
				t.Fatalf("OpPipe must have 4 operands, got %d", len(def.OperandWidths))
			}
			ops, _ := code.ReadOperands(def, ins[i+1:])
			return [4]int{ops[0], ops[1], ops[2], ops[3]}
		}
		read := 0
		for _, w := range def.OperandWidths {
			read += w
		}
		i += 1 + read
	}
	t.Fatal("no OpPipe instruction found in bytecode")
	return [4]int{}
}

// TestPipeParams_Compiler_NoArgs verifies that pipes without args always emit argsIdx == 0xFFFF.
func TestPipeParams_Compiler_NoArgs(t *testing.T) {
	inputs := []string{
		"[1,2,3] |window: $window",
		"[1,2] |chunk: $chunk",
		"[1,2] |map: $item * 2",
		"[1,2] |filter: $item > 1",
	}
	for _, input := range inputs {
		input := input
		t.Run(input, func(t *testing.T) {
			bc := compileExpr(t, input)
			ops := findOpPipeOperands(t, bc.Instructions)
			assert.Equal(t, 0xFFFF, ops[3],
				"no-args pipe must have argsIdx == 0xFFFF for: %s", input)
		})
	}
}

// TestPipeParams_Compiler_WithArgs verifies that OpPipe's 4th operand points to the args constant.
func TestPipeParams_Compiler_WithArgs(t *testing.T) {
	tests := []struct {
		input    string
		wantArgs []any
	}{
		{`[1,2,3] |window(3): $window`, []any{float64(3)}},
		{`[1,2,3,4,5] |chunk(4): $chunk`, []any{float64(4)}},
		{`arr |myPipe(3, "desc", true): $item`, []any{float64(3), "desc", true}},
		{`arr |myPipe(null): $item`, []any{nil}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			bc := compileExpr(t, tt.input)
			ops := findOpPipeOperands(t, bc.Instructions)
			argsIdx := ops[3]
			assert.NotEqual(t, 0xFFFF, argsIdx,
				"pipe with args must not have argsIdx == 0xFFFF for: %s", tt.input)
			if argsIdx >= len(bc.Constants) {
				t.Fatalf("argsIdx %d is out of constants range %d", argsIdx, len(bc.Constants))
			}
			gotArgs, ok := bc.Constants[argsIdx].ToAny().([]any)
			assert.True(t, ok, "constants[argsIdx] should be []any, got %T", bc.Constants[argsIdx].ToAny())
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

// TestPipeParams_Compiler_ChainedPipes verifies independent argsIdx for chained pipes.
func TestPipeParams_Compiler_ChainedPipes(t *testing.T) {
	t.Run("first no-args second has args", func(t *testing.T) {
		bc := compileExpr(t, "arr |filter: $item > 0 |window(3): $window")
		// Find all OpPipe instructions
		var opPipes [][4]int
		ins := bc.Instructions
		for i := 0; i < len(ins); {
			op := code.Opcode(ins[i])
			def, err := code.Lookup(ins[i])
			if err != nil {
				t.Fatalf("unknown opcode: %v", err)
			}
			if op == code.OpPipe {
				ops, _ := code.ReadOperands(def, ins[i+1:])
				opPipes = append(opPipes, [4]int{ops[0], ops[1], ops[2], ops[3]})
			}
			read := 0
			for _, w := range def.OperandWidths {
				read += w
			}
			i += 1 + read
		}
		if len(opPipes) != 2 {
			t.Fatalf("expected 2 OpPipe instructions, got %d", len(opPipes))
		}
		assert.Equal(t, 0xFFFF, opPipes[0][3], "filter pipe argsIdx must be 0xFFFF")
		assert.NotEqual(t, 0xFFFF, opPipes[1][3], "window pipe argsIdx must not be 0xFFFF")
		windowArgsIdx := opPipes[1][3]
		gotArgs, ok := bc.Constants[windowArgsIdx].ToAny().([]any)
		assert.True(t, ok)
		assert.Equal(t, []any{float64(3)}, gotArgs)
	})

	t.Run("both chained pipes have args", func(t *testing.T) {
		bc := compileExpr(t, "arr |window(3): $window |chunk(2): $chunk")
		var opPipes [][4]int
		ins := bc.Instructions
		for i := 0; i < len(ins); {
			op := code.Opcode(ins[i])
			def, err := code.Lookup(ins[i])
			if err != nil {
				t.Fatalf("unknown opcode: %v", err)
			}
			if op == code.OpPipe {
				ops, _ := code.ReadOperands(def, ins[i+1:])
				opPipes = append(opPipes, [4]int{ops[0], ops[1], ops[2], ops[3]})
			}
			read := 0
			for _, w := range def.OperandWidths {
				read += w
			}
			i += 1 + read
		}
		if len(opPipes) != 2 {
			t.Fatalf("expected 2 OpPipe instructions, got %d", len(opPipes))
		}
		// window args
		windowArgsIdx := opPipes[0][3]
		assert.NotEqual(t, 0xFFFF, windowArgsIdx)
		windowArgs, ok := bc.Constants[windowArgsIdx].ToAny().([]any)
		assert.True(t, ok)
		assert.Equal(t, []any{float64(3)}, windowArgs)
		// chunk args
		chunkArgsIdx := opPipes[1][3]
		assert.NotEqual(t, 0xFFFF, chunkArgsIdx)
		chunkArgs, ok2 := bc.Constants[chunkArgsIdx].ToAny().([]any)
		assert.True(t, ok2)
		assert.Equal(t, []any{float64(2)}, chunkArgs)
	})
}

// TestPipeParams_Compiler_EmptyArgs verifies that |pipe(): predicate emits argsIdx == 0xFFFF
// (empty parens treated as nil args — same as no parens).
func TestPipeParams_Compiler_EmptyArgs(t *testing.T) {
	bc := compileExpr(t, "arr |window(): $window")
	ops := findOpPipeOperands(t, bc.Instructions)
	assert.Equal(t, 0xFFFF, ops[3], "empty parens should produce argsIdx == 0xFFFF")
}

// TestPipeParams_Compiler_Bytecode_NoArgs checks the exact bytecode for a no-args pipe.
// This is a regression check: OpPipe must always be a 9-byte instruction.
func TestPipeParams_Compiler_Bytecode_NoArgs(t *testing.T) {
	// constants: 1.0(0), 2.0(1), "map"(2), 2.0(3 — predicate's literal), InstructionBlock(4)
	cases := []compilerTestCase{
		{`[1,2] |map: $item * 2`, []any{1.0, 2.0, "map", 2.0, nil}, []code.Instructions{
			code.Make(code.OpConstant, 0), // 1.0
			code.Make(code.OpConstant, 1), // 2.0
			code.Make(code.OpArray, 2),
			code.Make(code.OpPipe, 2, 0, 4, 0xFFFF), // argsIdx == 0xFFFF
		}},
	}
	runCompilerTestCases(t, cases)
}
