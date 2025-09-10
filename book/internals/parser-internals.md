# Parser internals

This page explains how the uEXL parser turns tokens into an AST, including precedence rules, features like nullish coalescing and optional chaining, and the typed property model for member/index access. It’s meant for contributors extending grammar or behavior.

## Design goals
- Clear, predictable precedence and associativity
- Helpful, structured errors with token context
- No hidden global state; explicit Parser options
- Compatibility with the compiler/VM pipeline

## Construction and options
- NewParser(input string) → *Parser: convenience with DefaultOptions()
- NewParserWithOptions(input string, opt Options) → *Parser
- NewParserWithValidation(input string) (*Parser, error): fails fast on invalid first token or empty input
- Options
  - EnableNullish, EnableOptionalChaining, EnablePipes (feature toggles)
  - MaxDepth (0 = unlimited)

## Parser type and rationale
- Type: Hand-written recursive descent with precedence climbing (and right-associative handling for '**')
- Why this design?
  - Readability and control: Grammar is compact; hand-written descent maps directly to the language spec and keeps operator precedence explicit in code.
  - Extensibility: Adding domain-specific constructs (pipes, optional/nullable member access, slicing) is straightforward without grammar rewrites or generator tooling.
  - Error messages: Fine-grained control over recovery and context-aware diagnostics (e.g., pipe placement, expected tokens) produces actionable errors.
  - Performance: For expression-sized inputs, a tight descent parser has minimal overhead vs. table-driven parsers; no reflective dispatch or generator runtime.
  - Portability: No parser generator dependency; pure Go, easy to evolve in lockstep with tokenizer.

## Core loop
- Parser holds
  - tokenizer: *Tokenizer
  - current: Token (lookahead)
  - errors: []errors.ParserError
  - pos: token index (for pipe alias logic)
  - subExpressionActive, inParenthesis: context flags for disambiguation and error policies
- advance(): pulls next token; converts tokenizer errors into ParserErrors while keeping parsing resumable
- Parse(): entrypoint; builds a full Expression or returns aggregated ParseErrors

## Grammar overview (selected)
- Expression → PipeExpression
- PipeExpression → Conditional ('|' PipeSegment)* [with alias support via `as $id`]
- Conditional → LogicalOr ('?' Conditional ':' Conditional)?
- Nullish → BitwiseShift ('??' BitwiseShift)*
- Bitwise/Equality/Comparison/Additive/Multiplicative → standard left-associative groups
- Power '**' → right-associative with special handling for leading unary operators
- Unary → ('-' | '!') Unary | MemberAccess
- MemberAccess → Primary ('.' (Identifier | Number | '(' Expression ')') | '[' Expr ']' | '?[' Expr ']' | '?.' Identifier | '?.[' Expr ']')*
- Primary → Number | String | Boolean | Null | Identifier | Array | Object | '(' Expression ')'

## AST nodes (highlights)
- BinaryExpression, UnaryExpression, ConditionalExpression
- Literals: NumberLiteral, StringLiteral (with IsRaw, IsSingleQuoted), BooleanLiteral, NullLiteral
- Identifier, ArrayLiteral, ObjectLiteral
- FunctionCall
- MemberAccess { Target, Property: Property, Optional }
- IndexAccess { Target, Index, Optional }
- SliceExpression { Target, Start, End, Step, Optional }
- PipeExpression and ProgramNode for top-level pipe chains

### Property typing
- Property is a tagged union: PropString | PropInt
  - PropS(s string) / PropI(i int) constructors
  - IsString()/IsInt() inspectors
- Member access after '.' supports
  - .identifier → PropS(identifier)
  - .number → PropI(int(number))
  - .(expr) → becomes IndexAccess with arbitrary expression
  - For numeric tokens containing dots (e.g., .1.5), the parser splits into segments: .1 .5

### Optional chaining and nullish
- Token forms handled by tokenizer: '?.', '?[' and '??'
- Optional chaining is carried by Optional flags on MemberAccess/IndexAccess/SliceExpression
- Nullish '??' participates at the specified precedence level

### Arrays, objects, and calls
- Arrays: '[' elements ']' with commas
- Objects: '{' string-key ':' expr (',' ... ) '}'
- FunctionCall: postfixed to identifiers or nested function calls only (not after member/index access directly)

## Pipes
- The parser accepts either '|' (default pipe) or '|name:' for named pipes; each segment may end with optional alias `as $var`.
- Top-level only: pipes are disallowed within sub-expressions except inside parentheses; errors ErrPipeInSubExpression/ErrEmptyPipe/ErrEmptyPipeWithAlias guide users
- ProgramNode contains a sequence of PipeExpression entries with PipeType, Alias, and Index

## Error handling strategy
- Errors collected in p.errors with exact line/column and the offending token when relevant
- Helper methods: addError, addErrorWithToken, addErrorWithExpected
- After a critical pipe error (e.g., empty segment), parser consumes remaining tokens to avoid cascades

## Sub-expression flags
- subExpressionActive and inParenthesis track context for
  - Restricting pipes within nested contexts
  - Allowing rich expressions in brackets, grouped expressions, and function args

## Power operator and unary peeling
- To mirror JS precedence, '**' is right-associative; leading '-' and '!' operators are peeled and reapplied around the exponent expression so `-2**3` parses as `-(2**3)`

## Contracts with the tokenizer
- Token.Type and Token.Value.Kind must be consistent; parser checks operator strings via Token.Value.Kind == TVKOperator
- Identifiers/numbers/dots are distinct tokens (e.g., '.' never gets merged into a number unless preceded by digits and followed by a digit at lex time)

## Debugging tips
- parser/tests/* contains unit tests covering precedence, chaining, slices, nullish/optional, and pipes
- Use tokenizer.PrintTokens() or PreloadTokens() for token stream inspection while debugging parser behavior

## Future extensions
- More literal forms (template strings, hex/bin numbers) would be added by augmenting tokenizer and parsePrimary
- Additional operators or precedence tweaks should update parseBinaryOp wiring and tests
