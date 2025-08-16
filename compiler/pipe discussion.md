# Pipes, Identifiers, Scopes, Frames – Design Summary

## 1. Variable Categories

- contextVars: External/environment variables (e.g. foo, bar) compiled via OpContextVar.
- systemVars: Pipe-local symbolic identifiers (currently names starting with $ like $item, $index) compiled via OpIdentifier.
- (Planned) Unified runtime scopes: Replace dedicated pipeScopes with a generic scope stack usable for pipes, future functions, blocks.

## 2. Original Problem

Test `[1,2] |map: $item * 2` failed: `undefined pipe variable: $item`.
Cause: Compiler emitted lambda code inline BEFORE OpPipe. VM executed OpIdentifier before pipe handler could bind $item.

## 3. Solutions Considered

| Approach                  | Description                                                                     | Pros                       | Cons                                                         |
| ------------------------- | ------------------------------------------------------------------------------- | -------------------------- | ------------------------------------------------------------ |
| Inline lambda (current)   | Emit lambda instructions, then OpPipe                                           | Simple emit                | Wrong execution order for $item                              |
| compileBlock + OpPipe ref | Compile lambda in its own scope; store block in constants; OpPipe references it | Correct ordering; reusable | Requires minor constant representation for instruction block |
| Jump / skip inline        | Emit lambda then jump over; OpPipe later replays                                | Avoids new constant type   | More complex IP management                                   |
| AST eval at runtime       | Store AST node, interpret per element                                           | No compiler change         | Slower; duplicate logic                                      |

Chosen direction: Separate compiled lambda per pipe step (future‑proof, matches “Writing A Compiler In Go” style).

## 4. Frames vs Pipe Scopes

- Frames: Execution context (instructions, ip, basePointer). Needed once you introduce re-entrant or nested executable units (future functions, pipe lambdas as blocks).
- pipeScopes (narrow): Only for pipe vars.
- Recommendation (industry standard): Keep Frames + a unified scopes []map[string]parser.Node. Drop specialized pipeScopes to avoid duplication.

## 5. Current State (Code Observations)

- VM struct still includes: pipeScopes and frames.
- OpIdentifier: Reads systemVars[identIndex] → name → looks up in getPipeVar (pipeScopes).
- OpPipe: Pops lambda & input, then calls handler (assumes lambda result already computed – not correct for map semantics).

## 6. Required Adjustments (If adopting block approach)

Compiler:

1. When seeing a PipeExpression:
   - Compile left/input expression (first segment).
   - For each subsequent pipe segment:
     - enterScope()
     - Compile pipe segment expression (lambda/body) into isolated instructions.
     - Capture those instructions (compileBlock).
     - exitScope()
     - Add a CompiledBlock node (implements parser.Node) to constants.
     - Emit OpPipe with operands: pipeTypeIdx, aliasIdx, lambdaBlockIdx.
2. Identifiers starting with $ still emit OpIdentifier with index into SystemVars.

VM:

1. OpPipe handler:

   - Read pipeType, alias, lambdaBlockIdx.
   - Retrieve CompiledBlock (instructions).
   - For map/filter style:
     - Expect input array on stack (or passed).
     - For each element:
       - Push new scope map (top of scopes).
       - Bind $item, $index, user alias (if provided).
       - Push a frame with lambda instructions (ip=0, basePointer=sp).
       - Run until frame completes (end of instructions).
       - Pop frame, pop scope.
       - Collect result (map) or test predicate (filter).
   - Push final transformed array (or value).

2. OpIdentifier:

   - identIndex → systemVars → name
   - Resolve: search scopes from top down; if not found error.

3. Scope helpers:
   - pushScope(): scopes = append(scopes, map[string]parser.Node{})
   - popScope(): scopes = scopes[:len(scopes)-1]
   - setVar(name, val): scopes[len-1][name] = val
   - getVar(name): walk backwards.

## 7. Minimal Transition Plan

Phase 1 (quick fix):

- Keep current pipeScopes.
- Delay lambda execution: Modify compiler to push (not execute) lambda representation (e.g. AST node) then let handler execute/evaluate it manually (temporary).

Phase 2 (proper):

- Implement compileBlock & CompiledBlock node.
- Replace pipeScopes with unified scopes.
- Implement per-element execution via frames.

Phase 3 (extension ready):

- Add future OpFunction / closures reusing same mechanism.

## 8. Testing Focus

- Update existing TestPipeFunction to expect new constant layout (pipe type, alias, block).
- Add tests:
  - `[1,2,3] |map: $item + 1`
  - Chained: `[1,2] |map: $item * 2 |map: $item + 1`
  - Alias: `[1,2] |map(x): $x * 3` (if alias syntax added)
  - Error: `$item` outside pipe (should fail).

## 9. Error Handling Guidelines

- OpIdentifier: return descriptive error if variable not found in scopes.
- OpPipe: error if input type mismatch (e.g., map expects array).
- Compile: error if pipe lambda missing expression.

## 10. Summary

To avoid bloat yet stay extensible:

- Keep frames (already present).
- Replace specialized pipeScopes with a unified scopes stack.
- Compile each pipe lambda in an isolated scope and store as a constant (CompiledBlock).
- Execute lambda per element using a new frame + pushed scope.
  This mirrors well-established VM architectures and positions your expression engine for future features without
