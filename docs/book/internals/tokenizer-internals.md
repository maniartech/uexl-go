# Tokenizer internals

This page documents the implementation details of the uEXL tokenizer so contributors can understand its design, performance characteristics, and extension points. It complements the public syntax docs, focusing on how tokens are produced.

## High-level goals
- Correctness for UTF-8 input and escape sequences
- High throughput with low allocations
- Clear error reporting with precise positions
- Stable token contracts consumed by the parser

## Core data structures
- Token
  - Type: parser/constants.TokenType
  - Value: TokenValue (tagged union for strongly typed payloads)
    - Kind: TVKNumber | TVKString | TVKBoolean | TVKNull | TVKIdentifier | TVKOperator
    - Num, Str, Bool: concrete value fields
  - Token: original lexeme (as appeared in source)
  - Line, Column: 1-based source position
  - IsSingleQuoted: only set for string tokens
- Token helper methods
  - AsFloat() (float64, bool)
  - AsString() (string, bool)
  - AsBool() (bool, bool)

These methods provide a compatibility and ergonomics layer for consumers that need typed values without switching on TokenValue.Kind.

## Tokenizer engine
### Type and rationale
- Type: Hand-written, single-pass, streaming lexer with UTF-8–aware cursor
- Why this design?
  - Performance: Tight control over branching, ASCII fast paths, and allocation patterns significantly outperforms table- or regex-driven lexers for short expressions.
  - Correctness: Explicit UTF-8 decode and cached rune/size avoid accidental byte-wise bugs on multibyte code points.
  - Simplicity: The token surface is small and stable; a custom lexer remains readable and easy to extend for domain-specific tokens (pipes, optional chaining forms).
  - Error quality: Direct access to positions enables precise, friendly errors without indirection.
  - Portability: No lexer generators or runtime deps; pure Go, easy to embed.

- Fields
  - input: source string
  - pos, line, column: cursor and positions
  - curRune, curSize: cached decode of current rune (avoids repeated utf8.DecodeRuneInString)
  - strBuf: reusable byte buffer for unescaping strings (minimizes allocations)
- Construction
  - NewTokenizer(input string) *Tokenizer: initializes state and decodes first rune

### Cursor operations
- current(): returns cached rune at pos
- peek(): decodes the next rune using curSize to advance correctly across multibyte UTF-8 sequences
- advance(): moves by curSize, updates line/column (line++ and column reset on '\n'), and refreshes cache via setCur()
- skipWhitespace(): ASCII fast path for ' ', '\n', '\t', '\r', with fallback to unicode.IsSpace for non-ASCII

### Token production
- NextToken(): dispatches based on current rune
  - Numbers → readNumber
  - Identifiers/keywords (and $ identifiers) → readIdentifierOrKeyword
  - Strings (regular and raw) → readString
  - Punctuation: parens/brackets/braces/comma/dot/colon
  - Question-prefixed tokens: readQuestionOrNullish for ?, ?., ?[
  - Pipes and bitwise or: readPipeOrBitwiseOr
  - Operators: readOperator
  - Invalid characters yield a structured parser error with location

### Numbers: readNumber
- ASCII fast-path for integer and fractional digits (avoids rune decoding)
- Optional fractional part only if '.' followed by digit (so '.' alone is a dot token for member access)
- Optional scientific notation (e or E with optional sign) validated before committing
- Fast manual parse for common sizes (<= 18 digits int/frac) using precomputed pow10[]; otherwise fallback to strconv.ParseFloat
- Emits TokenNumber with TokenValue{Kind: TVKNumber, Num: float64}

### Identifiers and keywords: readIdentifierOrKeyword
- Accepts leading letter/underscore/$ with ASCII fast path for subsequent chars; supports unicode letters via unicode.IsLetter
- Recognizes keywords: true/false → TokenBoolean; null → TokenNull; as → TokenAs
- Others emit TokenIdentifier with TokenValue{Kind: TVKIdentifier, Str: original}

### Strings: readString
- Supports raw strings with r"..." / r'...'
  - Doubled quote handling inside raw strings ("" → ", '' → ')
  - No backslash escapes processed
- Regular strings
  - Double-quoted: handles common escapes and \u/\U unicode forms via unescapeDoubleQuoted; falls back to strconv.Unquote on rare/invalid sequences
  - Single-quoted: common escapes handled by unescapeStringFast
- Tracks IsSingleQuoted for downstream consumers
- Emits TokenString with unescaped Value.Str and original Token for provenance

### Operators and special tokens
- readOperator covers: '??', '&&', '++', '**', '--' (context-aware split for unary sequences), '==', '!=', '<=', '>=', '<<', '>>', and single-char ops
- readQuestionOrNullish recognizes '?', '?.', '?[' into dedicated token types
- readPipeOrBitwiseOr recognizes
  - '||' as logical-or operator
  - '|name:' as a named pipe (ASCII letters only), or '|:' as default pipe
  - '|' (single) as bitwise-or operator

### Error handling
- Invalid character produces an error (not a sentinel token) with ErrInvalidCharacter
- Unterminated quote and invalid string/number forms return typed parser errors with line/column
- Parser integrates these via NewParserWithValidation or via Parser.advance() error-to-token conversion for error-tolerant parsing

## Performance techniques
- Cached rune (curRune) and size (curSize) eliminate redundant utf8.DecodeRuneInString calls
- ASCII fast paths for digits, letters, and whitespace
- Manual number parsing for common sizes avoids strconv allocations
- Reusable strBuf for string unescaping
- Pre-allocated charStrings[128] to avoid tiny string allocations for single-character punctuation

## Utilities
- PreloadTokens() []Token for debugging token streams
- PrintTokens() helper

## Extension points
- New keywords or operators can be added in readIdentifierOrKeyword and readOperator respectively
- Additional pipe name rules can be introduced inside readPipeOrBitwiseOr
- For new escape sequences, extend unescapeDoubleQuoted / unescapeStringFast

## Invariants and contracts
- Line/Column are 1-based and advance correctly across newlines
- peek() always returns the rune immediately following current rune, UTF-8 correct
- Token.Value’s Kind must match Type semantics (e.g., TokenNumber ⇒ TVKNumber)
- Token.Token preserves original lexeme for error messages and AST provenance
