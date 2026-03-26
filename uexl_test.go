package uexl_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/maniartech/uexl"
	"github.com/stretchr/testify/assert"
)

// ── helpers ──────────────────────────────────────────────────────────────────

var bg = context.Background()

// simple functions for test envs
func addFn(args ...any) (any, error) {
	a, _ := args[0].(float64)
	b, _ := args[1].(float64)
	return a + b, nil
}

func constFn(args ...any) (any, error) { return 42.0, nil }

func errorFn(args ...any) (any, error) { return nil, fmt.Errorf("intentional error") }

// testLib implements uexl.Lib for testing WithLib.
type testLib struct {
	functions uexl.Functions
	pipes     uexl.PipeHandlers
	globals   map[string]any
}

func (l testLib) Apply(cfg *uexl.EnvConfig) {
	if l.functions != nil {
		cfg.AddFunctions(l.functions)
	}
	if l.pipes != nil {
		cfg.AddPipeHandlers(l.pipes)
	}
	if l.globals != nil {
		cfg.AddGlobals(l.globals)
	}
}

// panicOnApplyLib calls fn inside Apply — used to test EnvConfig nil panics.
type panicOnApplyLib struct {
	fn func(*uexl.EnvConfig)
}

func (l panicOnApplyLib) Apply(cfg *uexl.EnvConfig) { l.fn(cfg) }

// ── package-level Eval ───────────────────────────────────────────────────────

func TestEval_basic(t *testing.T) {
	result, err := uexl.Eval("price * qty", map[string]any{"price": 5.0, "qty": 3.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 15.0, result)
}

func TestEval_noVars(t *testing.T) {
	result, err := uexl.Eval("1 + 2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 3.0, result)
}

func TestEval_parseError(t *testing.T) {
	_, err := uexl.Eval("1 +", nil)
	assert.Error(t, err)
}

func TestEval_errorIsParseErrors(t *testing.T) {
	_, err := uexl.Eval("1 +", nil)
	if err == nil {
		t.Fatal("expected parse error")
	}
	var pe *uexl.ParseErrors
	assert.True(t, errors.As(err, &pe), "expected *ParseErrors, got %T", err)
	if pe != nil {
		assert.Greater(t, len(pe.Errors), 0)
	}
}

func TestEval_stringExpr(t *testing.T) {
	result, err := uexl.Eval(`"hello"`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, "hello", result)
}

func TestEval_boolExpr(t *testing.T) {
	result, err := uexl.Eval("true", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, true, result)
}

func TestEval_nullExpr(t *testing.T) {
	result, err := uexl.Eval("null", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Nil(t, result)
}

// ── Default and NewEnv ───────────────────────────────────────────────────────

func TestDefault_hasBuiltins(t *testing.T) {
	env := uexl.Default()
	result, err := env.Eval(bg, "len('hi')", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 2.0, result)
}

func TestDefault_singleton(t *testing.T) {
	a := uexl.Default()
	b := uexl.Default()
	assert.Same(t, a, b, "Default() must return the same pointer")
}

func TestDefault_hasDefaultPipes(t *testing.T) {
	result, err := uexl.Eval("[1, 2, 3] |map: $item * 2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, []any{2.0, 4.0, 6.0}, result)
}

func TestNewEnv_blankSlate(t *testing.T) {
	env := uexl.NewEnv()

	// Compiling a function call against a blank env should fail.
	_, err := env.Compile("len('hi')")
	assert.Error(t, err, "blank env should reject unknown function 'len'")
	assert.Contains(t, err.Error(), "unknown function")
}

func TestNewEnv_noGlobals(t *testing.T) {
	env := uexl.NewEnv()
	info := env.Info()
	assert.Empty(t, info.Functions)
	assert.Empty(t, info.PipeHandlers)
	assert.Empty(t, info.Globals)
}

func TestDefaultWith_extendsDefault(t *testing.T) {
	env := uexl.DefaultWith(uexl.WithFunctions(uexl.Functions{"myConst": constFn}))

	// Has stdlib
	r1, err := env.Eval(bg, "len('abc')", nil)
	if err != nil {
		t.Fatalf("unexpected len error: %v", err)
	}
	assert.Equal(t, 3.0, r1)

	// Has custom fn
	r2, err := env.Eval(bg, "myConst()", nil)
	if err != nil {
		t.Fatalf("unexpected myConst error: %v", err)
	}
	assert.Equal(t, 42.0, r2)
}

func TestDefaultWith_doesNotMutateDefault(t *testing.T) {
	before := uexl.Default().Info()
	_ = uexl.DefaultWith(uexl.WithFunctions(uexl.Functions{"extra": constFn}))
	after := uexl.Default().Info()

	assert.Equal(t, before.Functions, after.Functions, "Default() info must not change after DefaultWith")
	assert.False(t, uexl.Default().HasFunction("extra"))
}

// ── Env.Compile ──────────────────────────────────────────────────────────────

func TestEnv_Compile_andEval(t *testing.T) {
	env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"add": addFn}))
	ce, err := env.Compile("add(x, y)")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}

	r1, err := ce.Eval(bg, map[string]any{"x": 1.0, "y": 2.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 3.0, r1)

	r2, err := ce.Eval(bg, map[string]any{"x": 10.0, "y": 20.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 30.0, r2)
}

func TestEnv_Compile_parseError(t *testing.T) {
	_, err := uexl.Default().Compile("1 +")
	assert.Error(t, err)
}

func TestEnv_Compile_unknownFunction(t *testing.T) {
	env := uexl.NewEnv()
	_, err := env.Compile("discount(price)")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown function")
	assert.Contains(t, err.Error(), "discount")
}

func TestEnv_Compile_unknownFunctionInPipePredicate(t *testing.T) {
	// Register 'map' pipe so pipe compiles, but not 'secret' function in predicate.
	env := uexl.NewEnv(
		uexl.WithPipeHandlers(uexl.PipeHandlers{
			"map": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
		}),
	)
	_, err := env.Compile("[1, 2] |map: secret($item)")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown function")
	assert.Contains(t, err.Error(), "secret")
}

func TestEnv_Compile_knownFunction(t *testing.T) {
	_, err := uexl.Default().Compile("len('hello')")
	assert.NoError(t, err)
}

// ── MustCompile ──────────────────────────────────────────────────────────────

func TestMustCompile_packageLevel(t *testing.T) {
	ce := uexl.MustCompile("1 + 2")
	assert.NotNil(t, ce)
}

func TestMustCompile_panicsOnBadExpr(t *testing.T) {
	assert.Panics(t, func() { uexl.MustCompile("1 +") })
}

func TestMustCompile_panicsWithPrefix(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		assert.Contains(t, fmt.Sprintf("%v", r), "uexl: MustCompile:")
	}()
	uexl.MustCompile("1 +")
}

func TestEnv_MustCompile_valid(t *testing.T) {
	ce := uexl.Default().MustCompile("1 + 2")
	assert.NotNil(t, ce)
}

func TestEnv_MustCompile_panicsOnBadExpr(t *testing.T) {
	assert.Panics(t, func() { uexl.Default().MustCompile("1 +") })
}

func TestEnv_MustCompile_panicsWithPrefix(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		assert.Contains(t, fmt.Sprintf("%v", r), "uexl: Env.MustCompile:")
	}()
	uexl.Default().MustCompile("1 +")
}

// ── CompiledExpr ─────────────────────────────────────────────────────────────

func TestCompiledExpr_Variables_basic(t *testing.T) {
	ce, err := uexl.Default().Compile("price * qty")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	vars := ce.Variables()
	// Sorted alphabetically.
	assert.Equal(t, []string{"price", "qty"}, vars)
}

func TestCompiledExpr_Variables_noVars(t *testing.T) {
	ce, err := uexl.Default().Compile("1 + 2")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	vars := ce.Variables()
	assert.NotNil(t, vars)
	assert.Empty(t, vars)
}

func TestCompiledExpr_Variables_sorted(t *testing.T) {
	ce, err := uexl.Default().Compile("z + a + m")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	vars := ce.Variables()
	assert.Equal(t, []string{"a", "m", "z"}, vars)
}

func TestCompiledExpr_Variables_independent(t *testing.T) {
	ce, err := uexl.Default().Compile("price * qty")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	v1 := ce.Variables()
	v1[0] = "MUTATED"
	v2 := ce.Variables()
	assert.Equal(t, "price", v2[0], "mutating returned slice must not affect next call")
}

func TestCompiledExpr_Env(t *testing.T) {
	env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"add": addFn}))
	ce, err := env.Compile("add(1, 2)")
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	assert.Same(t, env, ce.Env())
}

func TestCompiledExpr_Eval_cancelledContext(t *testing.T) {
	ce := uexl.MustCompile("1 + 2")
	ctx, cancel := context.WithCancel(bg)
	cancel() // already cancelled
	_, err := ce.Eval(ctx, nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestCompiledExpr_Eval_nilVars(t *testing.T) {
	ce := uexl.MustCompile("42")
	result, err := ce.Eval(bg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 42.0, result)
}

// ── Env.Extend ───────────────────────────────────────────────────────────────

func TestEnv_Extend_inherits(t *testing.T) {
	parent := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"constFn": constFn}))
	child := parent.Extend()

	r, err := child.Eval(bg, "constFn()", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 42.0, r)
}

func TestEnv_Extend_override(t *testing.T) {
	overridden := func(args ...any) (any, error) { return 99.0, nil }
	parent := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"fn": constFn}))
	child := parent.Extend(uexl.WithFunctions(uexl.Functions{"fn": overridden}))

	rc, err := child.Eval(bg, "fn()", nil)
	if err != nil {
		t.Fatalf("child eval error: %v", err)
	}
	assert.Equal(t, 99.0, rc)

	rp, err := parent.Eval(bg, "fn()", nil)
	if err != nil {
		t.Fatalf("parent eval error: %v", err)
	}
	assert.Equal(t, 42.0, rp, "parent must retain original function")
}

func TestEnv_Extend_parentUnchanged(t *testing.T) {
	parent := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"fn": constFn}))
	_ = parent.Extend(uexl.WithFunctions(uexl.Functions{"extra": addFn}))

	assert.False(t, parent.HasFunction("extra"), "parent must not gain child's functions")
}

func TestEnv_Extend_additionalPipes(t *testing.T) {
	doublePipe := func(ctx uexl.PipeContext, input any) (any, error) {
		arr, ok := input.([]any)
		if !ok {
			return nil, fmt.Errorf("double pipe expects array")
		}
		out := make([]any, len(arr))
		for i, v := range arr {
			out[i] = v.(float64) * 2
		}
		return out, nil
	}
	parent := uexl.NewEnv()
	child := parent.Extend(uexl.WithPipeHandlers(uexl.PipeHandlers{"double": doublePipe}))

	assert.True(t, child.HasPipe("double"))
	assert.False(t, parent.HasPipe("double"))
}

func TestEnv_Extend_multiLevel(t *testing.T) {
	l1 := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"f1": constFn}))
	l2 := l1.Extend(uexl.WithFunctions(uexl.Functions{"f2": addFn}))
	l3 := l2.Extend(uexl.WithFunctions(uexl.Functions{"f3": constFn}))

	assert.True(t, l3.HasFunction("f1"))
	assert.True(t, l3.HasFunction("f2"))
	assert.True(t, l3.HasFunction("f3"))
	assert.False(t, l1.HasFunction("f2"))
	assert.False(t, l1.HasFunction("f3"))
}

func TestEnv_Extend_inheritsGlobals(t *testing.T) {
	parent := uexl.NewEnv(uexl.WithGlobals(map[string]any{"rate": 0.1}))
	child := parent.Extend()

	result, err := child.Eval(bg, "rate", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 0.1, result)
}

// ── WithLib ───────────────────────────────────────────────────────────────────

func TestWithLib_appliesFunctionsAndPipes(t *testing.T) {
	lib := testLib{
		functions: uexl.Functions{"myConst": constFn},
		pipes:     uexl.PipeHandlers{"id": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil }},
		globals:   map[string]any{"gKey": "gVal"},
	}
	env := uexl.NewEnv(uexl.WithLib(lib))

	assert.True(t, env.HasFunction("myConst"))
	assert.True(t, env.HasPipe("id"))
	assert.True(t, env.HasGlobal("gKey"))
}

func TestWithLib_calledDuringConstruction(t *testing.T) {
	called := false
	lib := panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) { called = true }}
	uexl.NewEnv(uexl.WithLib(lib))
	assert.True(t, called)
}

func TestWithLib_overriddenByLaterOption(t *testing.T) {
	overriddenFn := func(args ...any) (any, error) { return 99.0, nil }
	lib := testLib{functions: uexl.Functions{"fn": constFn}}
	env := uexl.NewEnv(
		uexl.WithLib(lib),
		uexl.WithFunctions(uexl.Functions{"fn": overriddenFn}),
	)

	r, err := env.Eval(bg, "fn()", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 99.0, r, "later WithFunctions should override lib")
}

func TestWithLib_nil_panics(t *testing.T) {
	assert.Panics(t, func() { uexl.WithLib(nil) })
}

// ── WithGlobals ───────────────────────────────────────────────────────────────

func TestWithGlobals_usedWhenNoVar(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"rate": 0.2}))
	result, err := env.Eval(bg, "rate", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 0.2, result)
}

func TestWithGlobals_shadowedByVars(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"x": 1.0}))
	result, err := env.Eval(bg, "x", map[string]any{"x": 99.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, 99.0, result, "eval var must shadow global")
}

func TestWithGlobals_nilValueShadows(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"x": 42.0}))
	result, err := env.Eval(bg, "x ?? 'fallback'", map[string]any{"x": nil})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, "fallback", result, "nil eval var must shadow global and trigger nullish coalescing")
}

func TestWithGlobals_inheritedByExtend(t *testing.T) {
	parent := uexl.NewEnv(uexl.WithGlobals(map[string]any{"g": "global"}))
	child := parent.Extend()
	result, err := child.Eval(bg, "g", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, "global", result)
}

// ── Introspection ─────────────────────────────────────────────────────────────

func TestEnv_HasFunction_true(t *testing.T) {
	env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"fn": constFn}))
	assert.True(t, env.HasFunction("fn"))
}

func TestEnv_HasFunction_false(t *testing.T) {
	assert.False(t, uexl.NewEnv().HasFunction("nonexistent"))
}

func TestEnv_HasFunction_emptyString(t *testing.T) {
	assert.False(t, uexl.Default().HasFunction(""))
}

func TestEnv_HasFunction_inheritedFromParent(t *testing.T) {
	parent := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"fn": constFn}))
	child := parent.Extend()
	assert.True(t, child.HasFunction("fn"))
}

func TestEnv_HasPipe_true(t *testing.T) {
	env := uexl.NewEnv(uexl.WithPipeHandlers(uexl.PipeHandlers{
		"myPipe": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
	}))
	assert.True(t, env.HasPipe("myPipe"))
}

func TestEnv_HasPipe_false(t *testing.T) {
	assert.False(t, uexl.NewEnv().HasPipe("nonexistent"))
}

func TestEnv_HasPipe_emptyString(t *testing.T) {
	assert.False(t, uexl.Default().HasPipe(""))
}

func TestEnv_HasGlobal_true(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"key": "val"}))
	assert.True(t, env.HasGlobal("key"))
}

func TestEnv_HasGlobal_false(t *testing.T) {
	assert.False(t, uexl.NewEnv().HasGlobal("nonexistent"))
}

func TestEnv_HasGlobal_emptyString(t *testing.T) {
	assert.False(t, uexl.Default().HasGlobal(""))
}

func TestEnvInfo_sorted(t *testing.T) {
	env := uexl.NewEnv(
		uexl.WithFunctions(uexl.Functions{
			"z": constFn, "a": constFn, "m": constFn,
		}),
		uexl.WithPipeHandlers(uexl.PipeHandlers{
			"zp": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
			"ap": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
		}),
		uexl.WithGlobals(map[string]any{"zg": 1, "ag": 2}),
	)
	info := env.Info()
	assert.Equal(t, []string{"a", "m", "z"}, info.Functions)
	assert.Equal(t, []string{"ap", "zp"}, info.PipeHandlers)
	assert.Equal(t, []string{"ag", "zg"}, info.Globals)
}

func TestEnvInfo_stable(t *testing.T) {
	env := uexl.Default()
	i1 := env.Info()
	i2 := env.Info()
	assert.Equal(t, i1.Functions, i2.Functions)
	assert.Equal(t, i1.PipeHandlers, i2.PipeHandlers)
}

func TestEnvInfo_independent(t *testing.T) {
	env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"fn": constFn}))
	info := env.Info()
	original := make([]string, len(info.Functions))
	copy(original, info.Functions)

	// Mutate the slice.
	info.Functions[0] = "CORRUPTED"

	// New call must return unmodified data.
	info2 := env.Info()
	assert.Equal(t, original, info2.Functions)
}

func TestEnvInfo_emptyEnv(t *testing.T) {
	env := uexl.NewEnv()
	info := env.Info()
	assert.NotNil(t, info.Functions)
	assert.NotNil(t, info.PipeHandlers)
	assert.NotNil(t, info.Globals)
	assert.Empty(t, info.Functions)
	assert.Empty(t, info.PipeHandlers)
	assert.Empty(t, info.Globals)
}

func TestEnvInfo_String_format(t *testing.T) {
	env := uexl.NewEnv(
		uexl.WithFunctions(uexl.Functions{
			"bar": constFn,
			"foo": constFn,
		}),
		uexl.WithPipeHandlers(uexl.PipeHandlers{
			"myPipe": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
		}),
		uexl.WithGlobals(map[string]any{"version": "1.0"}),
	)
	s := env.Info().String()
	expected := "Env:\n  Functions (2): bar, foo\n  PipeHandlers (1): myPipe\n  Globals (1): version\n"
	assert.Equal(t, expected, s)
}

func TestEnvInfo_String_Stringer(t *testing.T) {
	info := uexl.EnvInfo{
		Functions:    []string{"a", "b"},
		PipeHandlers: []string{"p"},
		Globals:      []string{},
	}
	s := fmt.Sprintf("%v", info) // exercises Stringer interface
	assert.True(t, strings.HasPrefix(s, "Env:"))
	assert.Contains(t, s, "Functions (2): a, b")
	assert.Contains(t, s, "PipeHandlers (1): p")
	assert.Contains(t, s, "Globals (0):")
}

// ── Validate ──────────────────────────────────────────────────────────────────

func TestValidate_packageLevel_valid(t *testing.T) {
	assert.NoError(t, uexl.Validate("1 + 2"))
}

func TestValidate_packageLevel_parseError(t *testing.T) {
	assert.Error(t, uexl.Validate("1 +"))
}

func TestValidate_packageLevel_unknownFunction(t *testing.T) {
	// Default env does not have "noSuchFn".
	err := uexl.Validate("noSuchFn()")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown function")
}

func TestEnv_Validate_valid(t *testing.T) {
	assert.NoError(t, uexl.Default().Validate("len('hi')"))
}

func TestEnv_Validate_parseError(t *testing.T) {
	err := uexl.Default().Validate("1 +")
	assert.Error(t, err)
}

func TestEnv_Validate_noArtifact(t *testing.T) {
	// Validate signature is error only — verified at compile time by the API.
	err := uexl.Default().Validate("1 + 2")
	assert.NoError(t, err)
}

// ── Nil-guard panics ──────────────────────────────────────────────────────────

func TestWithFunctions_nil_panics(t *testing.T) {
	assert.PanicsWithValue(t, "uexl: WithFunctions: fns must not be nil",
		func() { uexl.WithFunctions(nil) })
}

func TestWithPipeHandlers_nil_panics(t *testing.T) {
	assert.PanicsWithValue(t, "uexl: WithPipeHandlers: pipes must not be nil",
		func() { uexl.WithPipeHandlers(nil) })
}

func TestWithGlobals_nil_panics(t *testing.T) {
	assert.PanicsWithValue(t, "uexl: WithGlobals: vars must not be nil",
		func() { uexl.WithGlobals(nil) })
}

func TestEnvConfig_AddFunctions_nil_panics(t *testing.T) {
	assert.Panics(t, func() {
		uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
			cfg.AddFunctions(nil)
		}}))
	})
}

func TestEnvConfig_AddFunctions_nil_panicsWithMessage(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, "uexl: EnvConfig.AddFunctions: fns must not be nil", fmt.Sprintf("%v", r))
	}()
	uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddFunctions(nil)
	}}))
}

func TestEnvConfig_AddPipeHandlers_nil_panics(t *testing.T) {
	assert.Panics(t, func() {
		uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
			cfg.AddPipeHandlers(nil)
		}}))
	})
}

func TestEnvConfig_AddPipeHandlers_nil_panicsWithMessage(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, "uexl: EnvConfig.AddPipeHandlers: pipes must not be nil", fmt.Sprintf("%v", r))
	}()
	uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddPipeHandlers(nil)
	}}))
}

func TestEnvConfig_AddGlobals_nil_panics(t *testing.T) {
	assert.Panics(t, func() {
		uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
			cfg.AddGlobals(nil)
		}}))
	})
}

func TestEnvConfig_AddGlobals_nil_panicsWithMessage(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, "uexl: EnvConfig.AddGlobals: vars must not be nil", fmt.Sprintf("%v", r))
	}()
	uexl.NewEnv(uexl.WithLib(panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddGlobals(nil)
	}}))
}

// ── Concurrency ───────────────────────────────────────────────────────────────

func TestCompiledExpr_concurrentEval(t *testing.T) {
	ce := uexl.MustCompile("price * qty")
	var wg sync.WaitGroup
	const goroutines = 50
	results := make([]any, goroutines)
	errs := make([]error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			p := float64(n)
			results[n], errs[n] = ce.Eval(bg, map[string]any{"price": p, "qty": 2.0})
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d error: %v", i, err)
		} else {
			assert.Equal(t, float64(i)*2.0, results[i], "goroutine %d", i)
		}
	}
}

func TestEnv_concurrentEval(t *testing.T) {
	env := uexl.Default()
	var wg sync.WaitGroup
	const goroutines = 50
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			result, err := env.Eval(bg, "x + 1", map[string]any{"x": float64(n)})
			if err != nil {
				t.Errorf("goroutine %d: %v", n, err)
				return
			}
			assert.Equal(t, float64(n)+1, result)
		}(i)
	}
	wg.Wait()
}

func TestEnv_concurrentInfo(t *testing.T) {
	env := uexl.Default()
	var wg sync.WaitGroup
	const goroutines = 50
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info := env.Info()
			assert.NotEmpty(t, info.Functions)
		}()
	}
	wg.Wait()
}

func TestEnv_concurrentExtend(t *testing.T) {
	parent := uexl.Default()
	var wg sync.WaitGroup
	const goroutines = 50
	envs := make([]*uexl.Env, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fnName := fmt.Sprintf("fn%d", n)
			envs[n] = parent.Extend(uexl.WithFunctions(uexl.Functions{fnName: constFn}))
		}(i)
	}
	wg.Wait()
	// Verify parent is unchanged and each child has its own function.
	for i, env := range envs {
		if env == nil {
			t.Errorf("env %d is nil", i)
			continue
		}
		assert.True(t, env.HasFunction(fmt.Sprintf("fn%d", i)))
	}
	// Parent must not have gained any child functions.
	for i := 0; i < goroutines; i++ {
		assert.False(t, parent.HasFunction(fmt.Sprintf("fn%d", i)))
	}
}

func TestEnv_concurrentValidate(t *testing.T) {
	env := uexl.Default()
	var wg sync.WaitGroup
	const goroutines = 50
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			var err error
			if n%2 == 0 {
				err = env.Validate("1 + 2")
				if err != nil {
					t.Errorf("goroutine %d: unexpected validate error: %v", n, err)
				}
			} else {
				err = env.Validate("1 +")
				if err == nil {
					t.Errorf("goroutine %d: expected validate error", n)
				}
			}
		}(i)
	}
	wg.Wait()
}

// ── Error type assertions ─────────────────────────────────────────────────────

func TestParseErrors_type(t *testing.T) {
	_, err := uexl.Eval("1 +", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *uexl.ParseErrors
	assert.True(t, errors.As(err, &pe))
}

func TestParserError_singleError(t *testing.T) {
	_, err := uexl.Eval("1 +", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var pe *uexl.ParseErrors
	if errors.As(err, &pe) && pe != nil {
		if assert.Greater(t, len(pe.Errors), 0) {
			single := pe.Errors[0]
			// ParserError should have Line and Column populated.
			assert.Greater(t, single.Line, 0)
		}
	}
}

// ── Options: WithFunctions merging ────────────────────────────────────────────

func TestWithFunctions_merging(t *testing.T) {
	fn1 := func(args ...any) (any, error) { return "fn1", nil }
	fn2 := func(args ...any) (any, error) { return "fn2", nil }
	env := uexl.NewEnv(
		uexl.WithFunctions(uexl.Functions{"a": fn1}),
		uexl.WithFunctions(uexl.Functions{"b": fn2}),
	)
	assert.True(t, env.HasFunction("a"))
	assert.True(t, env.HasFunction("b"))
}

func TestWithFunctions_laterWins(t *testing.T) {
	fn1 := func(args ...any) (any, error) { return "first", nil }
	fn2 := func(args ...any) (any, error) { return "second", nil }
	env := uexl.NewEnv(
		uexl.WithFunctions(uexl.Functions{"fn": fn1}),
		uexl.WithFunctions(uexl.Functions{"fn": fn2}),
	)
	r, err := env.Eval(bg, "fn()", nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	assert.Equal(t, "second", r)
}

// ── Options: WithPipeHandlers ─────────────────────────────────────────────────

func TestWithPipeHandlers_merging(t *testing.T) {
	pipe1 := func(ctx uexl.PipeContext, input any) (any, error) { return input, nil }
	pipe2 := func(ctx uexl.PipeContext, input any) (any, error) { return input, nil }
	env := uexl.NewEnv(
		uexl.WithPipeHandlers(uexl.PipeHandlers{"p1": pipe1}),
		uexl.WithPipeHandlers(uexl.PipeHandlers{"p2": pipe2}),
	)
	assert.True(t, env.HasPipe("p1"))
	assert.True(t, env.HasPipe("p2"))
}

// ── mergeVars edge cases ──────────────────────────────────────────────────────

func TestMergeVars_bothEmpty(t *testing.T) {
	// Empty globals + empty vars → expression uses no vars, just constant.
	env := uexl.NewEnv()
	result, err := env.Eval(bg, "7", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 7.0, result)
}

func TestMergeVars_globalsOnly(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"g": 5.0}))
	result, err := env.Eval(bg, "g", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 5.0, result)
}

func TestMergeVars_varsOnly(t *testing.T) {
	env := uexl.NewEnv()
	result, err := env.Eval(bg, "v", map[string]any{"v": 9.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 9.0, result)
}

func TestMergeVars_bothNonEmpty(t *testing.T) {
	env := uexl.NewEnv(uexl.WithGlobals(map[string]any{"g": 1.0, "x": 100.0}))
	result, err := env.Eval(bg, "g + x", map[string]any{"x": 2.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// x is shadowed by eval var; g comes from globals.
	assert.Equal(t, 3.0, result)
}

// ── Pipe integration ──────────────────────────────────────────────────────────

func TestDefaultPipes_map(t *testing.T) {
	result, err := uexl.Eval("[1, 2, 3] |map: $item * 10", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, []any{10.0, 20.0, 30.0}, result)
}

func TestDefaultPipes_filter(t *testing.T) {
	result, err := uexl.Eval("[1, 2, 3, 4] |filter: $item > 2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, []any{3.0, 4.0}, result)
}

func TestDefaultPipes_reduce(t *testing.T) {
	result, err := uexl.Eval("[1, 2, 3, 4] |reduce: ($acc || 0) + $item", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 10.0, result)
}

// ── Env.Eval runtime errors ───────────────────────────────────────────────────

func TestEnv_Eval_runtimeError(t *testing.T) {
	env := uexl.NewEnv(uexl.WithFunctions(uexl.Functions{"errFn": errorFn}))
	_, err := env.Eval(bg, "errFn()", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "intentional error")
}

// ── EnvConfig.AddFunctions (through WithLib) ──────────────────────────────────

func TestEnvConfig_AddFunctions_valid(t *testing.T) {
	lib := panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddFunctions(uexl.Functions{"libFn": constFn})
	}}
	env := uexl.NewEnv(uexl.WithLib(lib))
	assert.True(t, env.HasFunction("libFn"))
}

func TestEnvConfig_AddPipeHandlers_valid(t *testing.T) {
	lib := panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddPipeHandlers(uexl.PipeHandlers{
			"myp": func(ctx uexl.PipeContext, input any) (any, error) { return input, nil },
		})
	}}
	env := uexl.NewEnv(uexl.WithLib(lib))
	assert.True(t, env.HasPipe("myp"))
}

func TestEnvConfig_AddGlobals_valid(t *testing.T) {
	lib := panicOnApplyLib{fn: func(cfg *uexl.EnvConfig) {
		cfg.AddGlobals(map[string]any{"libGlobal": true})
	}}
	env := uexl.NewEnv(uexl.WithLib(lib))
	assert.True(t, env.HasGlobal("libGlobal"))
}
