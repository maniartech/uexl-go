//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/maniartech/uexl/code"
	"github.com/maniartech/uexl/compiler"
	"github.com/maniartech/uexl/parser"
	"github.com/maniartech/uexl/parser/errors"
	"github.com/maniartech/uexl/types"
	"github.com/maniartech/uexl/vm"
)

// evalError carries optional position info from the parser.
type evalError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}

type evalResponse struct {
	Ok               bool        `json:"ok"`
	Result           any         `json:"result,omitempty"`
	Errors           []evalError `json:"errors,omitempty"`
	ContextTime      int64       `json:"contextTime"`     // ns: JSON parse
	CompilationTime  int64       `json:"compilationTime"` // ns: parse + compile expression
	ExecutionTime    int64       `json:"executionTime"`   // ns: VM run
	Bytecode         string      `json:"bytecode,omitempty"`
	CompiledBytecode string      `json:"compiledBytecode,omitempty"`
}

type benchmarkResponse struct {
	Ok                  bool        `json:"ok"`
	Errors              []evalError `json:"errors,omitempty"`
	Iterations          int         `json:"iterations"`
	WarmupIterations    int         `json:"warmupIterations"`
	DurationNs          int64       `json:"durationNs"`
	ExecutionsPerSecond float64     `json:"executionsPerSecond"`
	AverageExecutionNs  float64     `json:"averageExecutionNs"`
}

func jsEvalUExL(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return respond(errResp("expression is required", "", 0, 0))
	}

	expr := args[0].String()

	// ── Timing note for future developers ──────────────────────────────────────
	// We use Go's standard time.Since() here, but in WASM that does NOT give
	// native OS-level nanosecond precision. When compiled for GOOS=js/GOARCH=wasm,
	// the Go runtime has no direct OS clock access — it runs inside the browser's
	// JS sandbox. The entire syscall layer, including the clock, is implemented by
	// wasm_exec.js, which delegates to JavaScript's performance.now().
	//
	// Call chain:
	//   time.Now() → Go runtime syscall → wasm_exec.js → performance.now()
	//
	// Since ~2018, all major browsers intentionally reduce performance.now()
	// resolution to ~100 µs (sometimes coarser) as a Spectre side-channel
	// mitigation. This causes readings to be quantized to ~100 µs boundaries,
	// which is why the playground shows values like 0ns, 99.8µs, 199.9µs
	// instead of precise values like 23µs or 66µs.
	//
	// Higher precision (~5 µs) is restored when the page is cross-origin isolated
	// (COOP: same-origin + COEP: require-corp headers). GitHub Pages cannot set
	// custom HTTP headers without a Service Worker shim, and we have chosen not
	// to add that complexity for a playground tool. The coarse timings are still
	// useful as a relative scale indicator (fast/slow) even if not cycle-accurate.
	// ─────────────────────────────────────────────────────────────────────────────

	// Phase 1: parse context JSON
	var contextVars map[string]any
	t0 := time.Now()
	if len(args) >= 2 {
		if ctxJSON := args[1].String(); ctxJSON != "" && ctxJSON != "{}" {
			if err := json.Unmarshal([]byte(ctxJSON), &contextVars); err != nil {
				return respond(errResp(fmt.Sprintf("invalid context JSON: %s", err.Error()), "", 0, 0))
			}
		}
	}
	contextTime := time.Since(t0).Nanoseconds()

	// Phase 2: parse + compile expression
	t1 := time.Now()
	node, parseErr := parser.ParseString(expr)
	if parseErr != nil {
		return respond(parseErrResp(parseErr))
	}
	comp := compiler.New()
	if compErr := comp.Compile(node); compErr != nil {
		return respond(parseErrResp(compErr))
	}
	compilationTime := time.Since(t1).Nanoseconds()

	compiledBytecode, marshalErr := json.Marshal(comp.ByteCode())
	if marshalErr != nil {
		return respond(errResp(fmt.Sprintf("failed to serialize compiled bytecode: %s", marshalErr.Error()), "internal-error", 0, 0))
	}

	// Phase 3: execute
	t2 := time.Now()
	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	result, runErr := machine.Run(comp.ByteCode(), contextVars)
	if runErr != nil {
		return respond(errResp(runErr.Error(), "runtime-error", 0, 0))
	}
	executionTime := time.Since(t2).Nanoseconds()

	return respond(evalResponse{
		Ok:               true,
		Result:           result,
		ContextTime:      contextTime,
		CompilationTime:  compilationTime,
		ExecutionTime:    executionTime,
		Bytecode:         disassemble(comp.ByteCode()),
		CompiledBytecode: string(compiledBytecode),
	})
}

func jsCompile(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return respond(errResp("expression is required", "", 0, 0))
	}

	expr := args[0].String()

	t1 := time.Now()
	node, parseErr := parser.ParseString(expr)
	if parseErr != nil {
		return respond(parseErrResp(parseErr))
	}
	comp := compiler.New()
	if compErr := comp.Compile(node); compErr != nil {
		return respond(parseErrResp(compErr))
	}
	compilationTime := time.Since(t1).Nanoseconds()

	compiledBytecode, marshalErr := json.Marshal(comp.ByteCode())
	if marshalErr != nil {
		return respond(errResp(fmt.Sprintf("failed to serialize compiled bytecode: %s", marshalErr.Error()), "internal-error", 0, 0))
	}

	return respond(evalResponse{
		Ok:               true,
		CompilationTime:  compilationTime,
		Bytecode:         disassemble(comp.ByteCode()),
		CompiledBytecode: string(compiledBytecode),
	})
}

// jsExecuteBytecode is an alternative entry point that accepts pre-compiled bytecode and context, and executes it directly without parsing/compilation overhead. This is not currently
// exposed to the playground UI, but can be used for testing or as a lower-level API.
func jsExecuteBytecode(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return respond(errResp("bytecode is required", "", 0, 0))
	}

	// Phase 1: parse bytecode JSON
	var bc compiler.ByteCode
	t0 := time.Now()
	if err := json.Unmarshal([]byte(args[0].String()), &bc); err != nil {
		return respond(errResp(fmt.Sprintf("invalid bytecode JSON: %s", err.Error()), "", 0, 0))
	}

	var contextVars map[string]any
	if len(args) >= 2 {
		if ctxJSON := args[1].String(); ctxJSON != "" && ctxJSON != "{}" {
			if err := json.Unmarshal([]byte(ctxJSON), &contextVars); err != nil {
				return respond(errResp(fmt.Sprintf("invalid context JSON: %s", err.Error()), "", 0, 0))
			}
		}
	}
	contextTime := time.Since(t0).Nanoseconds()

	// Phase 2: execute
	t1 := time.Now()
	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})
	result, runErr := machine.Run(&bc, contextVars)
	if runErr != nil {
		return respond(errResp(runErr.Error(), "runtime-error", 0, 0))
	}
	executionTime := time.Since(t1).Nanoseconds()

	return respond(evalResponse{
		Ok:            true,
		Result:        result,
		ContextTime:   contextTime,
		ExecutionTime: executionTime,
	})
}

func jsBenchmarkBytecode(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return respondBenchmark(benchmarkErrResp("bytecode is required", ""))
	}

	// Phase 1: parse bytecode JSON
	var bc compiler.ByteCode
	if err := json.Unmarshal([]byte(args[0].String()), &bc); err != nil {
		return respondBenchmark(benchmarkErrResp(fmt.Sprintf("invalid bytecode JSON: %s", err.Error()), ""))
	}

	warmupIterations := 100
	durationMs := 1500
	if len(args) >= 3 {
		warmupIterations = maxInt(0, args[2].Int())
	}
	if len(args) >= 4 {
		durationMs = maxInt(250, args[3].Int())
	}

	var baseContext map[string]any
	if len(args) >= 2 {
		if ctxJSON := args[1].String(); ctxJSON != "" && ctxJSON != "{}" {
			if err := json.Unmarshal([]byte(ctxJSON), &baseContext); err != nil {
				return respondBenchmark(benchmarkErrResp(fmt.Sprintf("invalid context JSON: %s", err.Error()), ""))
			}
		}
	}

	machine := vm.New(vm.LibContext{
		Functions:    vm.Builtins,
		PipeHandlers: vm.DefaultPipeHandlers,
	})

	// Warmup: excluded from measurements.
	for i := 0; i < warmupIterations; i++ {
		if _, err := machine.Run(&bc, mutateContext(baseContext, i)); err != nil {
			return respondBenchmark(benchmarkErrResp(err.Error(), "runtime-error"))
		}
	}

	start := time.Now()
	iterations := 0
	for time.Since(start) < time.Duration(durationMs)*time.Millisecond {
		if _, err := machine.Run(&bc, mutateContext(baseContext, warmupIterations+iterations)); err != nil {
			return respondBenchmark(benchmarkErrResp(err.Error(), "runtime-error"))
		}
		iterations++
	}
	elapsed := time.Since(start)
	if iterations == 0 {
		return respondBenchmark(benchmarkErrResp("benchmark completed with zero iterations", "runtime-error"))
	}

	return respondBenchmark(benchmarkResponse{
		Ok:                  true,
		Iterations:          iterations,
		WarmupIterations:    warmupIterations,
		DurationNs:          elapsed.Nanoseconds(),
		ExecutionsPerSecond: float64(iterations) / elapsed.Seconds(),
		AverageExecutionNs:  float64(elapsed.Nanoseconds()) / float64(iterations),
	})
}

// disassemble formats the bytecode as a human-readable string,
// with pipe predicate blocks inlined (indented) under their OpPipe line.
func disassemble(bc *compiler.ByteCode) string {
	var sb strings.Builder

	// Context variables
	if len(bc.ContextVars) > 0 {
		sb.WriteString("=== Context Vars ===\n")
		for i, v := range bc.ContextVars {
			sb.WriteString(fmt.Sprintf("%04d  %s\n", i, v))
		}
		sb.WriteString("\n")
	}

	// Instructions, with pipe blocks inlined
	sb.WriteString("=== Instructions ===\n")
	writeInstructions(&sb, bc.Instructions, bc.Constants, "")

	return sb.String()
}

// writeInstructions writes a disassembled instruction stream to sb.
// prefix is prepended to every line (used for indenting pipe blocks).
func writeInstructions(sb *strings.Builder, ins code.Instructions, constants []types.Value, prefix string) {
	i := 0
	for i < len(ins) {
		op := code.Opcode(ins[i])
		def, err := code.Lookup(ins[i])
		if err != nil {
			sb.WriteString(fmt.Sprintf("%sERROR: %s\n", prefix, err))
			i++
			continue
		}
		operands, read := code.ReadOperands(def, ins[i+1:])

		line := fmt.Sprintf("%s%04d %s %v\n", prefix, i, op.String(), operands)
		sb.WriteString(line)

		// For OpPipe, inline the predicate block indented beneath
		if op == code.OpPipe && len(operands) == 3 {
			pipeTypeIdx := operands[0]
			blockIdx := operands[2]

			// Resolve pipe name from constants
			pipeName := "pipe"
			if pipeTypeIdx >= 0 && pipeTypeIdx < len(constants) {
				if s, ok := constants[pipeTypeIdx].AnyVal.(string); ok {
					pipeName = s
				} else if constants[pipeTypeIdx].IsString() {
					pipeName = constants[pipeTypeIdx].StrVal
				}
			}

			// Resolve the InstructionBlock
			if blockIdx >= 0 && blockIdx < len(constants) {
				if blk, ok := constants[blockIdx].AnyVal.(*compiler.InstructionBlock); ok && blk != nil && len(blk.Instructions) > 0 {
					sb.WriteString(fmt.Sprintf("%s  ; %s predicate:\n", prefix, pipeName))
					writeInstructions(sb, blk.Instructions, constants, prefix+"  ")
				}
			}
		}

		i += 1 + read
	}
}

// parseErrResp converts parser ErrorList or a single ParserError into a response.
func parseErrResp(err error) evalResponse {
	switch e := err.(type) {
	case errors.ErrorList:
		errs := make([]evalError, 0, len(e))
		for _, pe := range e {
			errs = append(errs, evalError{
				Message: pe.Message,
				Code:    string(pe.Code),
				Line:    pe.Line,
				Column:  pe.Column,
			})
		}
		return evalResponse{Ok: false, Errors: errs}
	case *errors.ParserError:
		return evalResponse{Ok: false, Errors: []evalError{{
			Message: e.Message,
			Code:    string(e.Code),
			Line:    e.Line,
			Column:  e.Column,
		}}}
	case errors.ParserError:
		return evalResponse{Ok: false, Errors: []evalError{{
			Message: e.Message,
			Code:    string(e.Code),
			Line:    e.Line,
			Column:  e.Column,
		}}}
	default:
		return errResp(err.Error(), "", 0, 0)
	}
}

func errResp(msg, code string, line, col int) evalResponse {
	return evalResponse{Ok: false, Errors: []evalError{{Message: msg, Code: code, Line: line, Column: col}}}
}

func respondBenchmark(r benchmarkResponse) string {
	b, err := json.Marshal(r)
	if err != nil {
		return `{"ok":false,"errors":[{"message":"internal marshal error"}]}`
	}
	return string(b)
}

func benchmarkErrResp(msg, code string) benchmarkResponse {
	return benchmarkResponse{Ok: false, Errors: []evalError{{Message: msg, Code: code}}}
}

func mutateContext(base map[string]any, iteration int) map[string]any {
	if len(base) == 0 {
		return map[string]any{"__bench_iteration": iteration}
	}
	next := make(map[string]any, len(base)+1)
	for k, v := range base {
		switch t := v.(type) {
		case float64:
			next[k] = t + float64(iteration)
		case float32:
			next[k] = float64(t) + float64(iteration)
		case int:
			next[k] = t + iteration
		case int8:
			next[k] = int(t) + iteration
		case int16:
			next[k] = int(t) + iteration
		case int32:
			next[k] = int(t) + iteration
		case int64:
			next[k] = int(t) + iteration
		case uint:
			next[k] = int(t) + iteration
		case uint8:
			next[k] = int(t) + iteration
		case uint16:
			next[k] = int(t) + iteration
		case uint32:
			next[k] = int(t) + iteration
		case uint64:
			next[k] = int(t) + iteration
		case bool:
			if iteration%2 == 0 {
				next[k] = t
			} else {
				next[k] = !t
			}
		case string:
			next[k] = fmt.Sprintf("%s_%d", t, iteration)
		default:
			next[k] = v
		}
	}
	next["__bench_iteration"] = iteration
	return next
}

func maxInt(minValue, value int) int {
	if value < minValue {
		return minValue
	}
	return value
}

func respond(r evalResponse) string {
	b, err := json.Marshal(r)
	if err != nil {
		return `{"ok":false,"errors":[{"message":"internal marshal error"}]}`
	}
	return string(b)
}

func main() {
	js.Global().Set("evalUExL", js.FuncOf(jsEvalUExL))
	js.Global().Set("compileUExL", js.FuncOf(jsCompile))
	js.Global().Set("executeBytecode", js.FuncOf(jsExecuteBytecode))
	// Block forever — keeps the WASM instance alive.
	select {}
}
