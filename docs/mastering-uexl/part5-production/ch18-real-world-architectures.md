# Chapter 18: Real-World Architectures

> "Embedding an expression engine is not about replacing your Go code — it is about finding the exact boundary where business rules need to change faster than your deployment cycle."

---

## 18.1 The Architecture Decision

Before embedding UExL, answer this question: *Which parts of your system change more often than you ship?*

Business rules — discount thresholds, eligibility criteria, routing logic, content conditions — are the canonical answer. Code changes require a deployment; expression changes stored in a database do not.

UExL is a good fit when:
- Non-developer stakeholders need to author or modify rules
- Rule changes need zero-downtime deployment
- Rules are data, not algorithms
- You need tenant-specific or role-specific rule variants

UExL is **not** a good fit when:
- Logic is deeply recursive or requires state across multiple expressions
- You need database queries or I/O inside the rule
- Rules are complex enough to need their own test suites and abstractions

---

## 18.2 Pattern 1: The Rule Engine

A rule engine evaluates named conditions against a context, returning a boolean (match/no match) or a computed value.

### Structure

```go
type Rule struct {
    Name       string
    Expression string
    Priority   int
}

type RuleEngine struct {
    env     *uexl.Env
    rules   []compiledRule
}

type compiledRule struct {
    Rule
    compiled *uexl.CompiledExpr
}

func NewRuleEngine(env *uexl.Env, rules []Rule) (*RuleEngine, error) {
    compiled := make([]compiledRule, 0, len(rules))
    for _, r := range rules {
        ce, err := env.Compile(r.Expression)
        if err != nil {
            return nil, fmt.Errorf("rule %q invalid: %w", r.Name, err)
        }
        compiled = append(compiled, compiledRule{Rule: r, compiled: ce})
    }
    // Sort by priority (higher first)
    sort.Slice(compiled, func(i, j int) bool {
        return compiled[i].Priority > compiled[j].Priority
    })
    return &RuleEngine{env: env, rules: compiled}, nil
}
```

### First-match evaluation

```go
// Match returns the first rule that evaluates to a truthy value.
func (re *RuleEngine) Match(ctx context.Context, vars map[string]any) (*Rule, error) {
    for _, r := range re.rules {
        result, err := r.compiled.Eval(ctx, vars)
        if err != nil {
            return nil, fmt.Errorf("rule %q failed: %w", r.Name, err)
        }
        if isTruthy(result) {
            return &r.Rule, nil
        }
    }
    return nil, nil // no match
}

func isTruthy(v any) bool {
    switch v := v.(type) {
    case bool:
        return v
    case float64:
        return v != 0
    case string:
        return v != ""
    case nil:
        return false
    default:
        return true
    }
}
```

### Usage — ShopLogic discount eligibility

```go
rules := []Rule{
    {Name: "platinum-early-access", Priority: 100,
        Expression: "customer.tier == 'platinum' && product.stock < 20"},
    {Name: "loyal-customer",        Priority: 50,
        Expression: "customer.totalSpent > 2000 && customer.memberSince < '2023-01-01'"},
    {Name: "sale-eligible",         Priority: 10,
        Expression: "customer.active"},
}

engine, _ := NewRuleEngine(shopEnv, rules)

matched, err := engine.Match(ctx, EvalContext(product, customer))
if matched != nil {
    applyDiscount(matched.Name, product)
}
```

---

## 18.3 Pattern 2: Multi-Tenant Expression Isolation

In multi-tenant SaaS, each tenant may have their own set of rules and their own set of approved functions. Use `Env.Extend()` to create per-tenant environments that share the base function library but add tenant-specific functions and restrict scope.

```go
// Base env — shared function library, no globals
baseEnv := uexl.NewEnv(
    uexl.WithFunctions(shoplib.Functions),
    uexl.WithPipeHandlers(uexl.DefaultPipeHandlers),
)

// Per-tenant extension — add tenant globals and any tenant-specific functions
func tenantEnv(tenantID string, cfg TenantConfig) *uexl.Env {
    return baseEnv.Extend(
        uexl.WithGlobals(map[string]any{
            "TAX_RATE":     cfg.TaxRate,
            "CURRENCY":     cfg.Currency,
            "MAX_DISCOUNT": cfg.MaxDiscount,
        }),
    )
}

// Tenant expressions are compiled against their specific env
compiled, err := tenantEnv(tenantID, cfg).Compile(tenantRule.Expression)
```

Each tenant's expressions can only call functions registered in their env. Tenant A cannot call a function you registered only for Tenant B. This is structural isolation — no runtime checks needed.

---

## 18.4 Pattern 3: Access Control Policies

Expression-based access control lets you write authorization rules in a readable, auditable format instead of scattered conditional logic.

```go
type Policy struct {
    Resource string
    Action   string
    Rule     string // expression returning bool
}

var policies = []Policy{
    {Resource: "order", Action: "view",   Rule: "user.id == order.customerId || user.role == 'admin'"},
    {Resource: "order", Action: "cancel", Rule: "user.role == 'admin' || (user.id == order.customerId && order.status == 'pending')"},
    {Resource: "product", Action: "edit", Rule: "user.role == 'admin' || user.role == 'manager'"},
}

type ACL struct {
    rules map[string]*uexl.CompiledExpr
}

func NewACL(env *uexl.Env, policies []Policy) (*ACL, error) {
    rules := make(map[string]*uexl.CompiledExpr, len(policies))
    for _, p := range policies {
        key := p.Resource + ":" + p.Action
        compiled, err := env.Compile(p.Rule)
        if err != nil {
            return nil, fmt.Errorf("ACL policy %q invalid: %w", key, err)
        }
        rules[key] = compiled
    }
    return &ACL{rules: rules}, nil
}

func (a *ACL) Allow(ctx context.Context, resource, action string, vars map[string]any) (bool, error) {
    key := resource + ":" + action
    rule, ok := a.rules[key]
    if !ok {
        return false, nil // deny by default
    }
    result, err := rule.Eval(ctx, vars)
    if err != nil {
        return false, err
    }
    b, ok := result.(bool)
    return b && ok, nil
}
```

### Usage

```go
allowed, err := acl.Allow(ctx, "order", "cancel", map[string]any{
    "user":  map[string]any{"id": "u1", "role": "customer"},
    "order": map[string]any{"customerId": "u1", "status": "pending"},
})
```

---

## 18.5 Pattern 4: Dynamic Pricing and Configuration

Expressions as pricing rules let non-engineers adjust pricing without a code deployment.

```go
type PricingRule struct {
    Name     string
    Formula  string // expression computing float price
    Priority int
}

// Separate rules by priority — the first rule with a truthy 'applicable' check wins
// Or: use formula that returns null when not applicable, fall through with ??

// Chained rule example using ?? to pass through
// actual_price = platinum_price ?? gold_price ?? standard_price
```

### Chained nullish pricing

```go
// Store as three separate named expressions
var (
    platinumPrice = mustCompile(shopEnv,
        "customer.tier == 'platinum' ? product.basePrice * 0.75 : null")
    goldPrice     = mustCompile(shopEnv,
        "customer.tier == 'gold' ? product.basePrice * 0.88 : null")
    standardPrice = mustCompile(shopEnv,
        "product.basePrice")
)

func computePrice(ctx context.Context, vars map[string]any) (float64, error) {
    for _, rule := range []*uexl.CompiledExpr{platinumPrice, goldPrice, standardPrice} {
        result, err := rule.Eval(ctx, vars)
        if err != nil {
            return 0, err
        }
        if result != nil {
            if price, ok := result.(float64); ok {
                return price, nil
            }
        }
    }
    return 0, errors.New("no pricing rule matched")
}
```

---

## 18.6 Pattern 5: CMS Content Visibility

Content management systems often have "show this block to segment X" conditions. Expression-driven visibility rules let content editors control targeting without developer help.

```go
type ContentBlock struct {
    ID         string
    Visibility string // expression returning bool
    Content    string
}

func filterVisibleBlocks(ctx context.Context, blocks []ContentBlock, env *uexl.Env, vars map[string]any) ([]ContentBlock, error) {
    var visible []ContentBlock
    for _, block := range blocks {
        // Empty visibility = always visible
        if block.Visibility == "" {
            visible = append(visible, block)
            continue
        }
        compiled, err := env.Compile(block.Visibility)
        if err != nil {
            // Log and skip — don't break entire page
            slog.Warn("invalid visibility expression", "block", block.ID, "error", err)
            continue
        }
        result, err := compiled.Eval(ctx, vars)
        if err != nil || !isTruthy(result) {
            continue
        }
        visible = append(visible, block)
    }
    return visible, nil
}
```

> **Production note:** For CMS use cases with many blocks, compile visibility expressions at save time and cache the `*CompiledExpr`, not the expression string. The pattern above re-compiles on every page render, which is only appropriate for prototypes.

---

## 18.7 Hot Reloading Rules

When rules are stored in a database and administrators update them without downtime, use a versioned rule store with atomic replacement:

```go
type RuleStore struct {
    mu      sync.RWMutex
    engine  *RuleEngine
}

func (rs *RuleStore) Reload(env *uexl.Env, rules []Rule) error {
    engine, err := NewRuleEngine(env, rules)
    if err != nil {
        return err  // keep old engine on compile failure
    }
    rs.mu.Lock()
    rs.engine = engine
    rs.mu.Unlock()
    return nil
}

func (rs *RuleStore) Match(ctx context.Context, vars map[string]any) (*Rule, error) {
    rs.mu.RLock()
    engine := rs.engine
    rs.mu.RUnlock()
    return engine.Match(ctx, vars)
}
```

`RuleEngine` itself is immutable once created. Swap the pointer under a write lock; all in-flight `Match` calls hold the read lock and complete against the old engine. New calls pick up the new engine.

---

## 18.8 Summary

- Rule engines, access control, pricing, and CMS targeting are the canonical UExL use cases.
- Compile expressions at startup or rule-load time, not per request.
- Use `Env.Extend()` to isolate per-tenant function sets and globals.
- Deny by default in access control — return `false` when no matching policy exists.
- Use nullish chaining (`??`) in pricing rule chains to fall through cleanly.
- Hot-reload rules atomically with a `sync.RWMutex` on the compiled rule engine pointer.
- For CMS: cache compiled expressions with the rule, not the string.

---

## Exercises

**18.1 — Recall.** Why is using `env.Extend()` for per-tenant environments better than creating an entirely new env per tenant on every request?

**18.2 — Apply.** Implement an A/B test router using UExL. The router has `n` variants, each with an expression that returns `true` when the request context qualifies for that variant. The router returns the name of the first matching variant, or `"control"` if none match.

**18.3 — Extend.** ShopLogic needs a rule that activates a "seasonal promo" only if: the customer has been a member for at least one year, their tier is not `'standard'`, and they have not already used the promo this month (stored as `customer.lastPromoMonth`). Design the expression, the context fields it requires, and the Go loading pattern for this rule. How would you update the rule mid-season without redeploying?
