//go:build ignore

// kitchensink_example.go — Complete walkthrough of the UExL public API.
//
// This file is a living reference implementation for the spec in spec.md.
// It is excluded from normal compilation (build tag "ignore") so it can be
// maintained alongside the spec before the public API is fully implemented.
//
// Run once the API is implemented:
//
//	go run designdocs/public-api/kitchensink_example.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/maniartech/uexl"
)

// ─────────────────────────────────────────────────────────────────────────────
// § User Lib — Finance domain bundle
//
// A Lib groups related functions, pipe handlers, and globals into one
// composable unit. Library authors ship this; app developers call WithLib.
// ─────────────────────────────────────────────────────────────────────────────

// FinanceLib bundles domain-specific finance extensions.
// It implements the uexl.Lib interface with a single Apply method.
type FinanceLib struct {
	// DefaultTaxRate is the fallback rate, baked in at construction.
	DefaultTaxRate float64
}

func (fl FinanceLib) Apply(cfg *uexl.EnvConfig) {
	// ── Functions ────────────────────────────────────────────────────────────

	cfg.AddFunctions(uexl.Functions{
		// discount(price, pct) → price reduced by pct percent
		// Usage: discount(100, 20) → 80
		"discount": func(args ...any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("discount: expected 2 args, got %d", len(args))
			}
			price, ok1 := args[0].(float64)
			pct, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("discount: both args must be numeric")
			}
			return price * (1 - pct/100), nil
		},

		// tax(amount, rate) → amount with rate% tax added
		// Usage: tax(100, 8.5) → 108.5
		"tax": func(args ...any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("tax: expected 2 args, got %d", len(args))
			}
			amount, ok1 := args[0].(float64)
			rate, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("tax: both args must be numeric")
			}
			return amount * (1 + rate/100), nil
		},

		// round(v, places) → v rounded to N decimal places
		// Usage: round(3.14159, 2) → 3.14
		"round": func(args ...any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("round: expected 2 args, got %d", len(args))
			}
			v, ok1 := args[0].(float64)
			places, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("round: both args must be numeric")
			}
			shift := 1.0
			for i := 0; i < int(places); i++ {
				shift *= 10
			}
			return float64(int(v*shift+0.5)) / shift, nil
		},
	})

	// ── Pipe Handlers ────────────────────────────────────────────────────────

	cfg.AddPipeHandlers(uexl.PipeHandlers{
		// |top: expr — evaluates expr for each item, keeps the N highest values.
		// Usage: prices |top: $item > threshold
		// (Here simplified to a map-like predicate returning a score.)
		"score": func(ctx uexl.PipeContext, input any) (any, error) {
			items, ok := input.([]any)
			if !ok {
				return nil, fmt.Errorf("score pipe: expected array, got %T", input)
			}
			type scored struct {
				val   any
				score float64
			}
			results := make([]any, 0, len(items))
			for i, item := range items {
				if ctx.Context().Err() != nil {
					return nil, ctx.Context().Err()
				}
				val, err := ctx.EvalItem(item, i) // runs predicate with $item, $index
				if err != nil {
					return nil, err
				}
				f, ok := val.(float64)
				if !ok {
					continue // skip non-numeric scores
				}
				if f > 0 {
					results = append(results, item)
				}
			}
			return results, nil
		},
	})

	// ── Globals ───────────────────────────────────────────────────────────────
	// Globals are available in every expression evaluated against this env
	// without the caller needing to pass them as per-call vars.

	cfg.AddGlobals(map[string]any{
		"TAX_RATE":         fl.DefaultTaxRate, // e.g. 8.5 for 8.5%
		"FREE_SHIP_THRESH": 50.0,              // orders above this get free shipping
		"SHIP_COST":        5.99,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// § Package-level compiled expressions
//
// MustCompile is safe for top-level var declarations where the expression is
// known-correct at write time. It panics immediately on startup if invalid,
// rather than silently failing later.
// ─────────────────────────────────────────────────────────────────────────────

// These are compiled once at program startup and reused across all goroutines.
var (
	// Subtotal: sum of unit prices × quantities across line items.
	subtotalRule = uexl.MustCompile(
		"items |reduce: $acc + $item.price * $item.qty",
	)

	// Shipping: free above threshold, otherwise flat fee.
	shippingRule = uexl.MustCompile(
		"subtotal >= FREE_SHIP_THRESH ? 0 : SHIP_COST",
	)

	// Adult filter: keep only users aged 18+.
	adultFilter = uexl.MustCompile(
		"users |filter: $item.age >= 18 |map: $item.name",
	)
)

// ─────────────────────────────────────────────────────────────────────────────
// main — walks through every section of the public API
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	section("1. Zero-setup: Package-level Eval")
	demoQuickstart()

	section("2. Validate — syntax & compile-time check with no artifact")
	demoValidate()

	section("3. Error handling — parse errors, compile errors, runtime errors")
	demoErrors()

	section("4. DefaultWith — stdlib + custom functions")
	demoDefaultWith()

	section("5. WithGlobals — org-wide constants in env")
	demoGlobals()

	section("6. WithLib — reusable domain bundle (FinanceLib)")
	demoFinanceLib()

	section("7. Extend — multi-level environment chain")
	demoExtend()

	section("8. Compile once, Eval many — the hot path")
	demoCompileAndEval()

	section("9. Variables() — compile-time variable introspection")
	demoVariables()

	section("10. Introspection — HasFunction, HasPipe, HasGlobal, Info")
	demoIntrospection()

	section("11. Complex pipe expressions — map, filter, reduce, sort, groupBy")
	demoPipes()

	section("12. Result coercion helpers — AsFloat64, AsBool, AsString, AsSlice")
	demoResultHelpers()

	section("13. Concurrency — 100 goroutines, one *CompiledExpr")
	demoConcurrency()

	section("14. Context cancellation — timeout on long-running expression")
	demoCancellation()
}

// ─────────────────────────────────────────────────────────────────────────────
// § 1. Zero-setup quickstart
// ─────────────────────────────────────────────────────────────────────────────

func demoQuickstart() {
	// uexl.Eval is the absolute minimum — no env setup, no imports beyond uexl.
	// Internally it calls Default().Eval(context.Background(), expr, vars).

	result, err := uexl.Eval("price * qty * (1 - discount / 100)", map[string]any{
		"price":    149.99,
		"qty":      3.0,
		"discount": 10.0, // 10% off
	})
	must(err)

	total, err := uexl.AsFloat64(result)
	must(err)
	fmt.Printf("  order total (10%% off): %.2f\n", total) // 404.97
}

// ─────────────────────────────────────────────────────────────────────────────
// § 2. Validate
// ─────────────────────────────────────────────────────────────────────────────

func demoValidate() {
	// Validate checks syntax AND function existence — no *CompiledExpr allocated.
	// Ideal for REST validation endpoints and CI lint steps.

	valid := []string{
		"price * qty",
		"len('hello') > 3",
		"[1, 2, 3] |map: $item * 2",
	}
	for _, expr := range valid {
		if err := uexl.Validate(expr); err != nil {
			log.Fatalf("  unexpected validation error: %v", err)
		}
		fmt.Printf("  ✓  %s\n", expr)
	}

	invalid := []string{
		"price *",          // syntax error: missing right operand
		"unknownFn(price)", // compile error: function not registered
	}
	for _, expr := range invalid {
		err := uexl.Validate(expr)
		if err == nil {
			log.Fatalf("  expected validation error for: %s", expr)
		}
		fmt.Printf("  ✗  %-28s → %v\n", expr, err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// § 3. Error handling
// ─────────────────────────────────────────────────────────────────────────────

func demoErrors() {
	// ── Parse error ───────────────────────────────────────────────────────────
	_, err := uexl.Eval("1 +", nil)
	if err != nil {
		var single uexl.ParserError
		var multi uexl.ParseErrors
		switch {
		case errors.As(err, &single):
			fmt.Printf("  parse error at %d:%d — %s\n",
				single.Line, single.Column, single.Message)
		case errors.As(err, &multi):
			fmt.Printf("  %d parse errors; first: %s\n",
				len(multi.Errors), multi.Errors[0].Message)
		}
	}

	// ── Compile error: unknown function ───────────────────────────────────────
	_, err = uexl.Eval("npv(rate, cashflows)", map[string]any{
		"rate": 0.1, "cashflows": []any{100.0, 200.0},
	})
	if err != nil {
		fmt.Printf("  compile error: %v\n", err)
	}

	// ── Runtime error: type mismatch ──────────────────────────────────────────
	env := uexl.DefaultWith(uexl.WithFunctions(uexl.Functions{
		"strictAdd": func(args ...any) (any, error) {
			a, ok1 := args[0].(float64)
			b, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("strictAdd: expected two numbers")
			}
			return a + b, nil
		},
	}))
	_, err = env.Eval(context.Background(), "strictAdd(x, y)", map[string]any{
		"x": "hello", "y": 1.0, // x is wrong type
	})
	if err != nil {
		fmt.Printf("  runtime error: %v\n", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// § 4. DefaultWith — stdlib + custom functions
// ─────────────────────────────────────────────────────────────────────────────

func demoDefaultWith() {
	// DefaultWith = Default().Extend(opts...) — keeps all stdlib built-ins
	// and adds your custom functions on top.

	env := uexl.DefaultWith(uexl.WithFunctions(uexl.Functions{
		// clamp(v, lo, hi) → v clamped to [lo, hi]
		"clamp": func(args ...any) (any, error) {
			if len(args) != 3 {
				return nil, fmt.Errorf("clamp: expected 3 args")
			}
			v, _ := args[0].(float64)
			lo, _ := args[1].(float64)
			hi, _ := args[2].(float64)
			if v < lo {
				return lo, nil
			}
			if v > hi {
				return hi, nil
			}
			return v, nil
		},
	}))

	// Mix stdlib (len) with custom (clamp) in the same expression.
	result, err := env.Eval(context.Background(),
		"clamp(len(tags), 0, 10)",
		map[string]any{
			"tags": []any{"go", "uexl", "expressions", "pipes", "finance", "rules",
				"dynamic", "eval", "bytecode", "vm", "compiler", "parser"},
		},
	)
	must(err)
	fmt.Printf("  clamp(12 tags, 0, 10) = %.0f\n", result) // 10
}

// ─────────────────────────────────────────────────────────────────────────────
// § 5. WithGlobals — org-wide constants
// ─────────────────────────────────────────────────────────────────────────────

func demoGlobals() {
	// Globals are available in every expression without per-call vars.
	// Per-call vars shadow globals with the same name (see §4 of spec).

	env := uexl.DefaultWith(uexl.WithGlobals(map[string]any{
		"PI":      3.141592653589793,
		"E":       2.718281828459045,
		"MAX_AGE": 120.0,
	}))

	// PI is a global — no per-call vars needed.
	area, err := env.Eval(context.Background(),
		"PI * radius * radius",
		map[string]any{"radius": 5.0},
	)
	must(err)
	fmt.Printf("  circle area (r=5): %.4f\n", area) // 78.5398

	// Per-call var "PI" shadows the global — useful for overrides.
	areaOverride, err := env.Eval(context.Background(),
		"PI * radius * radius",
		map[string]any{"radius": 5.0, "PI": 3.14}, // rough PI for quick math
	)
	must(err)
	fmt.Printf("  circle area (rough PI): %.2f\n", areaOverride) // 78.50
}

// ─────────────────────────────────────────────────────────────────────────────
// § 6. WithLib — reusable domain bundle
// ─────────────────────────────────────────────────────────────────────────────

func demoFinanceLib() {
	// WithLib registers a Lib in one call — FinanceLib brings discount, tax,
	// round functions, the |score: pipe, and TAX_RATE / FREE_SHIP_THRESH globals.

	finEnv := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	// discount + tax combo using both custom functions and a stdlib global.
	result, err := finEnv.Eval(context.Background(),
		"round(tax(discount(price, discPct), TAX_RATE), 2)",
		map[string]any{
			"price":   199.99,
			"discPct": 15.0, // 15% off → 169.99
			// tax applied: 169.99 * 1.085 = 184.44
		},
	)
	must(err)
	fmt.Printf("  final price after 15%% discount + %.1f%% tax: %.2f\n",
		8.5, result) // 184.44

	// Shipping: free when subtotal >= FREE_SHIP_THRESH global (50.0)
	ship, err := finEnv.Eval(context.Background(),
		"subtotal >= FREE_SHIP_THRESH ? 0 : SHIP_COST",
		map[string]any{"subtotal": 62.50},
	)
	must(err)
	fmt.Printf("  shipping cost (subtotal=62.50): %.2f\n", ship) // 0

	ship2, err := finEnv.Eval(context.Background(),
		"subtotal >= FREE_SHIP_THRESH ? 0 : SHIP_COST",
		map[string]any{"subtotal": 30.00},
	)
	must(err)
	fmt.Printf("  shipping cost (subtotal=30.00): %.2f\n", ship2) // 5.99
}

// ─────────────────────────────────────────────────────────────────────────────
// § 7. Extend — multi-level environment chain
// ─────────────────────────────────────────────────────────────────────────────

func demoExtend() {
	// Level 0: standard library only
	base := uexl.Default()

	// Level 1: add finance domain lib
	finEnv := base.Extend(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	// Level 2: per-tenant override — different tax rate and a tenant-specific fn
	tenantEnv := finEnv.Extend(
		uexl.WithGlobals(map[string]any{
			"TAX_RATE": 5.0, // tenant in a lower-tax jurisdiction
		}),
		uexl.WithFunctions(uexl.Functions{
			// loyaltyBonus(price) → subtract 3% loyalty credit
			"loyaltyBonus": func(args ...any) (any, error) {
				if len(args) != 1 {
					return nil, fmt.Errorf("loyaltyBonus: expected 1 arg")
				}
				price, _ := args[0].(float64)
				return price * 0.97, nil
			},
		}),
	)

	// base: tax at 8.5%
	r1, _ := finEnv.Eval(context.Background(), "round(tax(100, TAX_RATE), 2)", nil)
	fmt.Printf("  finEnv  (8.5%% tax): %.2f\n", r1) // 108.50

	// tenant: tax at 5%, plus loyalty discount, using its own extended function
	r2, _ := tenantEnv.Eval(context.Background(),
		"round(loyaltyBonus(tax(100, TAX_RATE)), 2)", nil)
	fmt.Printf("  tenantEnv (5%% tax + loyalty): %.2f\n", r2) // 101.85

	// parent is unaffected — TAX_RATE is still 8.5 here
	r3, _ := finEnv.Eval(context.Background(), "TAX_RATE", nil)
	fmt.Printf("  finEnv TAX_RATE unchanged: %.1f\n", r3) // 8.5

	// child can use parent's functions (discount is inherited from finEnv)
	r4, _ := tenantEnv.Eval(context.Background(),
		"round(discount(200, 10), 2)", nil)
	fmt.Printf("  tenantEnv inherits discount(): %.2f\n", r4) // 180.00
}

// ─────────────────────────────────────────────────────────────────────────────
// § 8. Compile once, Eval many — the hot path
// ─────────────────────────────────────────────────────────────────────────────

func demoCompileAndEval() {
	finEnv := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	// Compile parses + compiles + validates function names.
	// The returned *CompiledExpr is immutable and goroutine-safe.
	orderTotal, err := finEnv.Compile(
		"round(tax(items |reduce: $acc + $item.price * $item.qty, TAX_RATE), 2)",
	)
	must(err)

	orders := []map[string]any{
		{"items": []any{
			map[string]any{"price": 29.99, "qty": 2.0},
			map[string]any{"price": 9.99, "qty": 5.0},
		}},
		{"items": []any{
			map[string]any{"price": 149.99, "qty": 1.0},
		}},
		{"items": []any{
			map[string]any{"price": 4.99, "qty": 10.0},
			map[string]any{"price": 19.99, "qty": 3.0},
		}},
	}

	fmt.Println("  Order totals (with 8.5% tax):")
	for i, vars := range orders {
		result, err := orderTotal.Eval(context.Background(), vars)
		must(err)
		total, _ := uexl.AsFloat64(result)
		fmt.Printf("    order %d: $%.2f\n", i+1, total)
	}

	// Package-level MustCompile uses Default() env — already compiled at init time.
	// subtotalRule, shippingRule, adultFilter are defined at top of this file.
	fmt.Printf("  subtotalRule variables: %v\n", subtotalRule.Variables())
	fmt.Printf("  shippingRule variables: %v\n", shippingRule.Variables())
}

// ─────────────────────────────────────────────────────────────────────────────
// § 9. Variables() — what does this expression depend on?
// ─────────────────────────────────────────────────────────────────────────────

func demoVariables() {
	finEnv := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	exprs := []string{
		"price * qty - discount(price, discPct)",
		"[1, 2, 3] |map: $item * multiplier",
		"len(name) > 0 && isActive",
		"1 + 2", // no variables
	}

	for _, expr := range exprs {
		ce, err := finEnv.Compile(expr)
		must(err)
		fmt.Printf("  %-50s → vars: %v\n", expr, ce.Variables())
	}

	// Practical use: preflight check before eval
	rule, _ := finEnv.Compile("gross - round(discount(gross, discPct), 2)")
	record := map[string]any{"gross": 500.0} // missing discPct

	missing := []string{}
	for _, v := range rule.Variables() {
		if _, ok := record[v]; !ok {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		fmt.Printf("  preflight: missing fields: %s\n", strings.Join(missing, ", "))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// § 10. Introspection
// ─────────────────────────────────────────────────────────────────────────────

func demoIntrospection() {
	finEnv := uexl.DefaultWith(
		uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}),
		uexl.WithFunctions(uexl.Functions{
			"clamp": func(args ...any) (any, error) { return nil, nil },
		}),
	)

	// HasFunction / HasPipe / HasGlobal — O(1) lookup, goroutine-safe
	fmt.Printf("  HasFunction(\"discount\"): %v\n", finEnv.HasFunction("discount"))
	fmt.Printf("  HasFunction(\"unknown\"):  %v\n", finEnv.HasFunction("unknown"))
	fmt.Printf("  HasPipe(\"map\"):          %v\n", finEnv.HasPipe("map"))
	fmt.Printf("  HasPipe(\"score\"):        %v\n", finEnv.HasPipe("score"))
	fmt.Printf("  HasGlobal(\"TAX_RATE\"):   %v\n", finEnv.HasGlobal("TAX_RATE"))
	fmt.Printf("  HasGlobal(\"PI\"):         %v\n", finEnv.HasGlobal("PI"))

	// Info() — sorted snapshot for logging, diagnostics, and tab-completion
	info := finEnv.Info()
	fmt.Printf("\n  Registered functions (%d):\n    %s\n",
		len(info.Functions), strings.Join(info.Functions, ", "))
	fmt.Printf("  Registered pipes (%d):\n    %s\n",
		len(info.PipeHandlers), strings.Join(info.PipeHandlers, ", "))
	fmt.Printf("  Registered globals (%d):\n    %s\n",
		len(info.Globals), strings.Join(info.Globals, ", "))

	// Env() on a CompiledExpr: back-reference to the compiling environment
	rule, _ := finEnv.Compile("discount(price, discPct)")
	envInfo := rule.Env().Info()
	fmt.Printf("\n  rule compiled against env with %d functions\n",
		len(envInfo.Functions))
}

// ─────────────────────────────────────────────────────────────────────────────
// § 11. Complex pipe expressions
// ─────────────────────────────────────────────────────────────────────────────

func demoPipes() {
	env := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	products := []any{
		map[string]any{"name": "Laptop", "price": 999.99, "category": "Electronics", "stock": 5.0},
		map[string]any{"name": "Headphones", "price": 149.99, "category": "Electronics", "stock": 0.0},
		map[string]any{"name": "Desk", "price": 349.99, "category": "Furniture", "stock": 12.0},
		map[string]any{"name": "Chair", "price": 199.99, "category": "Furniture", "stock": 3.0},
		map[string]any{"name": "Monitor", "price": 449.99, "category": "Electronics", "stock": 8.0},
		map[string]any{"name": "Keyboard", "price": 79.99, "category": "Electronics", "stock": 20.0},
		map[string]any{"name": "Lamp", "price": 39.99, "category": "Furniture", "stock": 15.0},
	}

	vars := map[string]any{"products": products}

	// ── |map: ─────────────────────────────────────────────────────────────────
	// Project each product to its name.
	names, _ := env.Eval(context.Background(),
		"products |map: $item.name", vars)
	fmt.Printf("  names: %v\n", names)

	// ── |filter: ──────────────────────────────────────────────────────────────
	// Keep only in-stock products.
	inStock, _ := env.Eval(context.Background(),
		"products |filter: $item.stock > 0 |map: $item.name", vars)
	fmt.Printf("  in stock: %v\n", inStock)

	// ── |reduce: ──────────────────────────────────────────────────────────────
	// Sum all product prices.
	totalValue, _ := env.Eval(context.Background(),
		"products |reduce: $acc + $item.price", vars)
	fmt.Printf("  total catalog value: %.2f\n", totalValue)

	// ── |sort: ────────────────────────────────────────────────────────────────
	// Sort products by price ascending and return names.
	sortedByPrice, _ := env.Eval(context.Background(),
		"products |sort: $item.price |map: $item.name", vars)
	fmt.Printf("  sorted by price: %v\n", sortedByPrice)

	// ── |groupBy: ─────────────────────────────────────────────────────────────
	// Group products by category.
	byCategory, _ := env.Eval(context.Background(),
		"products |groupBy: $item.category", vars)
	if groups, ok := byCategory.(map[string]any); ok {
		for cat, items := range groups {
			if arr, ok := items.([]any); ok {
				fmt.Printf("  %s (%d items)\n", cat, len(arr))
			}
		}
	}

	// ── Chained multi-pipe expression ─────────────────────────────────────────
	// "Electronics in stock, sorted by price, show name: price"
	summary, _ := env.Eval(context.Background(),
		`products
			|filter: $item.category == "Electronics" && $item.stock > 0
			|sort:   $item.price
			|map:    $item.name`,
		vars,
	)
	fmt.Printf("  electronics (in stock, by price): %v\n", summary)

	// ── Pipe + custom function ─────────────────────────────────────────────────
	// Apply finance discount to each product's price.
	discountedPrices, _ := env.Eval(context.Background(),
		"products |filter: $item.stock > 0 |map: round(discount($item.price, 10), 2)",
		vars,
	)
	fmt.Printf("  10%% discounted prices (in stock): %v\n", discountedPrices)

	// ── |some: and |every: ────────────────────────────────────────────────────
	anyInStock, _ := env.Eval(context.Background(),
		"products |some: $item.stock > 0", vars)
	allInStock, _ := env.Eval(context.Background(),
		"products |every: $item.stock > 0", vars)
	fmt.Printf("  any in stock: %v, all in stock: %v\n", anyInStock, allInStock)

	// ── |find: ────────────────────────────────────────────────────────────────
	cheapest, _ := env.Eval(context.Background(),
		`products |find: $item.price < 50`, vars)
	if item, ok := cheapest.(map[string]any); ok {
		fmt.Printf("  first item under $50: %s (%.2f)\n",
			item["name"], item["price"])
	}

	// ── |unique: ──────────────────────────────────────────────────────────────
	categories, _ := env.Eval(context.Background(),
		"products |map: $item.category |unique: $item", vars)
	fmt.Printf("  unique categories: %v\n", categories)

	// ── |chunk: ───────────────────────────────────────────────────────────────
	// Split products into pages of 3 for pagination.
	pages, _ := env.Eval(context.Background(),
		"products |chunk: 3", vars)
	if chunks, ok := pages.([]any); ok {
		fmt.Printf("  %d pages of ≤3 products\n", len(chunks))
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// § 12. Result coercion helpers
// ─────────────────────────────────────────────────────────────────────────────

func demoResultHelpers() {
	env := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))
	vars := map[string]any{"x": 42.0, "name": "Alice", "active": true}

	// AsFloat64 — safe numeric extraction with numeric widening
	r, _ := env.Eval(context.Background(), "x * 2.5", vars)
	f, err := uexl.AsFloat64(r)
	must(err)
	fmt.Printf("  AsFloat64: %.1f\n", f) // 105.0

	// AsBool — strict, no truthy coercion (0 is NOT false here, only bool is)
	r2, _ := env.Eval(context.Background(), "active && x > 40", vars)
	b, err := uexl.AsBool(r2)
	must(err)
	fmt.Printf("  AsBool:    %v\n", b) // true

	// AsBool rejects non-booleans (explicit semantics)
	r3, _ := env.Eval(context.Background(), "x", vars) // returns float64
	_, err = uexl.AsBool(r3)
	fmt.Printf("  AsBool(42): error (expected) → %v\n", err)

	// AsString
	r4, _ := env.Eval(context.Background(), "name", vars)
	s, err := uexl.AsString(r4)
	must(err)
	fmt.Printf("  AsString:  %q\n", s) // "Alice"

	// AsSlice
	r5, _ := env.Eval(context.Background(),
		"[1, 2, 3] |map: $item * $item", nil)
	arr, err := uexl.AsSlice(r5)
	must(err)
	fmt.Printf("  AsSlice:   %v (len=%d)\n", arr, len(arr)) // [1 4 9]

	// AsMap
	r6, _ := env.Eval(context.Background(),
		`{"user": name, "score": x}`, vars)
	m, err := uexl.AsMap(r6)
	must(err)
	fmt.Printf("  AsMap:     user=%v score=%v\n", m["user"], m["score"])
}

// ─────────────────────────────────────────────────────────────────────────────
// § 13. Concurrency — share one *CompiledExpr across goroutines
// ─────────────────────────────────────────────────────────────────────────────

func demoConcurrency() {
	finEnv := uexl.DefaultWith(uexl.WithLib(FinanceLib{DefaultTaxRate: 8.5}))

	// Compile once — *CompiledExpr is immutable, goroutine-safe.
	priceRule, err := finEnv.Compile("round(tax(discount(price, discPct), TAX_RATE), 2)")
	must(err)

	const goroutines = 100
	results := make([]float64, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			// Each goroutine supplies different per-call vars.
			price := 100.0 + float64(idx)
			disc := float64(idx % 20) // vary discount 0–19%

			result, err := priceRule.Eval(context.Background(), map[string]any{
				"price":   price,
				"discPct": disc,
			})
			if err != nil {
				log.Printf("goroutine %d error: %v", idx, err)
				return
			}
			f, _ := uexl.AsFloat64(result)
			results[idx] = f
		}(i)
	}

	wg.Wait()

	// Spot-check a few results
	fmt.Printf("  goroutine  0: price=100, disc= 0%% → %.2f\n", results[0])
	fmt.Printf("  goroutine 10: price=110, disc=10%% → %.2f\n", results[10])
	fmt.Printf("  goroutine 50: price=150, disc=10%% → %.2f\n", results[50])
	fmt.Printf("  all %d goroutines completed with no data races\n", goroutines)
}

// ─────────────────────────────────────────────────────────────────────────────
// § 14. Context cancellation
// ─────────────────────────────────────────────────────────────────────────────

func demoCancellation() {
	env := uexl.Default()

	// Simulate a very large dataset — the VM checks ctx.Done() between opcodes.
	bigData := make([]any, 100_000)
	for i := range bigData {
		bigData[i] = float64(i)
	}

	// Short deadline: 1 ms — will cancel before the reduce finishes.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := env.Eval(ctx,
		"data |reduce: $acc + $item", // sum 100k items
		map[string]any{"data": bigData},
	)

	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("  ✓ expression cancelled by context deadline as expected")
	} else if err != nil {
		fmt.Printf("  cancelled with error: %v\n", err)
	} else {
		// Completed before deadline — machine was fast enough; not an error.
		fmt.Println("  completed before deadline (fast machine)")
	}

	// Generous deadline — expression completes normally.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	result, err := env.Eval(ctx2,
		"[1, 2, 3, 4, 5] |reduce: $acc + $item",
		nil,
	)
	must(err)
	fmt.Printf("  small reduce with 5s deadline: %.0f\n", result) // 15
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func section(title string) {
	bar := strings.Repeat("─", 60)
	fmt.Printf("\n%s\n  %s\n%s\n", bar, title, bar)
}

func must(err error) {
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}
