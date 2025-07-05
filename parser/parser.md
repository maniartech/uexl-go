# UExl-Go Parser Documentation

## Overview

The UExl-Go parser is a production-ready, recursive descent parser for the UExl (Unified Expression Language) that follows industry-standard practices for error handling, tokenization, and AST construction. This document provides comprehensive information for developers working with or extending the parser.

## Architecture Classification

### Parser Type: **Recursive Descent Parser**
- **Parsing Strategy**: Top-down recursive descent with operator precedence
- **Error Recovery**: Fail-fast with comprehensive error reporting
- **Memory Model**: Direct AST construction (no intermediate representations)
- **Threading Model**: Single-threaded, non-concurrent (thread-safe for read-only operations)

### Industry Standards Compliance
- **Go Standard Library Patterns**: Follows `go/parser` and `go/scanner` error handling conventions
- **Error Propagation**: Direct error returns with `(result, error)` pattern
- **API Design**: Consistent with Go standard library parser interfaces
- **Error Messages**: Rich context with position information and error codes

## Core Components

### 1. Tokenizer (`tokenizer.go`)

**Classification**: Finite State Machine Tokenizer
**Pattern**: Industry-standard lexical analysis with direct error propagation

```go
// Industry-standard tokenizer interface
func (t *Tokenizer) NextToken() (Token, error)
```

**Key Features**:
- **Position Tracking**: Line and column information for all tokens
- **Error Handling**: Direct error returns (no error tokens)
- **Type Safety**: Strongly typed tokens with native value parsing
- **Unicode Support**: Full UTF-8 support for identifiers and strings

**Supported Token Types**:
```go
// Literals
TokenNumber      // 42, 3.14, 1.5e-10
TokenString      // "hello", 'world', r"raw"
TokenBoolean     // true, false
TokenNull        // null

// Identifiers and Keywords
TokenIdentifier  // variable, $1, user_name
TokenAs          // as (for pipe aliases)

// Operators
TokenOperator    // +, -, *, /, %, ==, !=, <, >, <=, >=, &&, ||, !, &, |, ^, <<, >>

// Structural
TokenLeftParen   // (
TokenRightParen  // )
TokenLeftBracket // [
TokenRightBracket// ]
TokenLeftBrace   // {
TokenRightBrace  // }
TokenComma       // ,
TokenDot         // .
TokenColon       // :
TokenPipe        // |, |:, |map:, |filter:, etc.

// Special
TokenEOF         // End of input
TokenError       // Error state (used internally)
```

**Example Usage**:
```go
tokenizer := NewTokenizer("user.name + 42")
for {
    token, err := tokenizer.NextToken()
    if err != nil {
        log.Fatal(err)
    }
    if token.Type == constants.TokenEOF {
        break
    }
    fmt.Printf("%s: %v\n", token.Type, token.Value)
}
```

### 2. Parser (`parser.go`)

**Classification**: Recursive Descent Parser with Precedence Climbing
**Pattern**: Pratt Parser approach for expression parsing

**Operator Precedence** (highest to lowest):
1. **Member Access**: `.` (object.property, array[index])
2. **Unary Operators**: `-`, `!` (right-associative)
3. **Multiplicative**: `*`, `/`, `%` (left-associative)
4. **Additive**: `+`, `-` (left-associative)
5. **Bitwise Shift**: `<<`, `>>` (left-associative)
6. **Comparison**: `<`, `>`, `<=`, `>=` (left-associative)
7. **Equality**: `==`, `!=` (left-associative)
8. **Bitwise AND**: `&` (left-associative)
9. **Bitwise XOR**: `^` (left-associative)
10. **Bitwise OR**: `|` (left-associative)
11. **Logical AND**: `&&` (left-associative)
12. **Logical OR**: `||` (left-associative)
13. **Pipe Operations**: `|:`, `|map:`, `|filter:`, etc. (left-associative)

**Example Precedence Resolution**:
```go
// Input: "a + b * c.d"
// Parsed as: a + (b * (c.d))
&BinaryExpression{
    Left: &Identifier{Name: "a"},
    Operator: "+",
    Right: &BinaryExpression{
        Left: &Identifier{Name: "b"},
        Operator: "*",
        Right: &MemberAccess{
            Object: &Identifier{Name: "c"},
            Property: "d"
        }
    }
}
```

### 3. AST Types (`types.go`)

**Classification**: Strongly-typed AST with position tracking
**Pattern**: Visitor pattern compatible with interface-based polymorphism

**Core Interfaces**:
```go
type Node interface {
    Position() (line, column int)
}

type Expression interface {
    Node
    expressionNode() // Type marker method
}
```

**Expression Types**:

#### Literals
```go
type NumberLiteral struct {
    Value  string // Original string representation
    Line   int
    Column int
}

type StringLiteral struct {
    Value          string // Processed value (unquoted, unescaped)
    Token          string // Original token with quotes
    IsRaw          bool   // Raw string indicator
    IsSingleQuoted bool   // Quote type tracking
    Line           int
    Column         int
}

type BooleanLiteral struct {
    Value  bool
    Line   int
    Column int
}

type NullLiteral struct {
    Line   int
    Column int
}
```

#### Complex Expressions
```go
type BinaryExpression struct {
    Left     Expression
    Operator string
    Right    Expression
    Line     int
    Column   int
}

type UnaryExpression struct {
    Operator string
    Operand  Expression
    Line     int
    Column   int
}

type MemberAccess struct {
    Object   Expression
    Property string
    Line     int
    Column   int
}

type FunctionCall struct {
    Function  Expression
    Arguments []Expression
    Line      int
    Column    int
}
```

#### Collections
```go
type ArrayLiteral struct {
    Elements []Expression
    Line     int
    Column   int
}

type ObjectLiteral struct {
    Properties map[string]Expression
    Line       int
    Column     int
}
```

#### Pipe Expressions
```go
type PipeExpression struct {
    Expressions []Expression
    PipeTypes   []string
    Aliases     []string
    Line        int
    Column      int
}
```

## Best Practices Implemented

### 1. Error Handling Best Practices

**Industry Standard Compliance**:
- **No Panic Usage**: All error conditions return proper Go errors
- **Direct Error Propagation**: Tokenizer returns `(Token, error)`
- **Rich Error Context**: All errors include position and contextual information
- **Error Categorization**: Structured error codes for different error types

**Example Error Usage**:
```go
parser := NewParser("1 + ")
result, err := parser.Parse()
if err != nil {
    if parseErrors, ok := err.(*errors.ParseErrors); ok {
        for _, e := range parseErrors.Errors {
            fmt.Printf("Error %s at %d:%d: %s\n",
                e.Code, e.Line, e.Column, e.Message)
        }
    }
}
```

### 2. Memory Management Best Practices

**Efficient Memory Usage**:
- **Minimal Allocations**: Reuse token structures where possible
- **Position Tracking**: Lightweight position information
- **String Interning**: Efficient operator and keyword handling

**Example**:
```go
// Efficient token creation
func (t *Tokenizer) singleCharToken(tokenType constants.TokenType) (Token, error) {
    token := Token{
        Type:   tokenType,
        Token:  string(t.current()),
        Line:   t.line,
        Column: t.column,
    }
    t.advance()
    return token, nil
}
```

### 3. Type Safety Best Practices

**Strong Typing**:
- **Token Value Types**: Native Go types in token values
- **Expression Polymorphism**: Interface-based expression types
- **Position Information**: Consistent position tracking

**Example**:
```go
// Type-safe token values
func (t *Tokenizer) readNumber() (Token, error) {
    // ... parsing logic ...

    // Store both string and parsed value
    return Token{
        Type:  constants.TokenNumber,
        Value: value, // float64 or int64
        Token: numberStr,
        Line:  startLine,
        Column: startColumn,
    }, nil
}
```

### 4. Parser State Management

**Context-Aware Parsing**:
- **Sub-expression Tracking**: Prevents invalid pipe usage in sub-expressions
- **Parenthesis Context**: Allows pipes within parentheses
- **State Restoration**: Proper cleanup of parser state

**Example**:
```go
func (p *Parser) parseGroupedExpression() Expression {
    p.advance() // consume '('

    // Save and set context
    wasInParenthesis := p.inParenthesis
    wasSubExpressionActive := p.subExpressionActive
    p.inParenthesis = true
    p.subExpressionActive = true

    expr := p.parseExpression()

    // Restore context
    p.inParenthesis = wasInParenthesis
    p.subExpressionActive = wasSubExpressionActive

    // ... rest of parsing
    return expr
}
```

## Usage Examples

### Basic Expression Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/maniartech/uexl_go/parser"
)

func main() {
    // Simple arithmetic
    result, err := parser.ParseString("2 + 3 * 4")
    if err != nil {
        log.Fatal(err)
    }

    // Object member access
    result, err = parser.ParseString("user.profile.name")
    if err != nil {
        log.Fatal(err)
    }

    // Function calls
    result, err = parser.ParseString("max(a, b) + length(arr)")
    if err != nil {
        log.Fatal(err)
    }

    // Complex expressions with pipes
    result, err = parser.ParseString("[1, 2, 3] |map: $1 * 2 |filter: $1 > 3")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Advanced Usage with Error Handling

```go
func parseWithDetailedErrors(input string) {
    parser, err := parser.NewParserWithValidation(input)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }

    result, err := parser.Parse()
    if err != nil {
        if parseErrors, ok := err.(*errors.ParseErrors); ok {
            fmt.Printf("Found %d parsing errors:\n", len(parseErrors.Errors))
            for _, e := range parseErrors.Errors {
                fmt.Printf("  %s at line %d, column %d: %s\n",
                    e.Code, e.Line, e.Column, e.Message)
            }
        } else {
            fmt.Printf("Parse error: %v\n", err)
        }
        return
    }

    fmt.Printf("Successfully parsed: %T\n", result)
}
```

### Custom Tokenization

```go
func exploreTokens(input string) {
    tokenizer := parser.NewTokenizer(input)

    fmt.Printf("Tokenizing: %s\n", input)
    for {
        token, err := tokenizer.NextToken()
        if err != nil {
            fmt.Printf("Tokenization error: %v\n", err)
            break
        }

        if token.Type == constants.TokenEOF {
            break
        }

        fmt.Printf("  %s: '%s' (value: %v) at %d:%d\n",
            token.Type, token.Token, token.Value, token.Line, token.Column)
    }
}
```

## Language Features Supported

### 1. Data Types
- **Numbers**: Integers, floats, scientific notation (`42`, `3.14`, `1.5e-10`)
- **Strings**: Single/double quoted, raw strings (`"hello"`, `'world'`, `r"raw\string"`)
- **Booleans**: `true`, `false`
- **Null**: `null`
- **Arrays**: `[1, 2, 3]`, `[user.name, "literal", 42]`
- **Objects**: `{"key": value, "nested": {"inner": true}}`

### 2. Operators
- **Arithmetic**: `+`, `-`, `*`, `/`, `%`
- **Comparison**: `<`, `>`, `<=`, `>=`, `==`, `!=`
- **Logical**: `&&`, `||`, `!`
- **Bitwise**: `&`, `|`, `^`, `<<`, `>>`
- **Unary**: `-` (negation), `!` (logical not)

### 3. Member Access
- **Dot Notation**: `object.property`, `user.profile.name`
- **Array Indexing**: `array[0]`, `data[user.id]`
- **Chained Access**: `users[0].profile.address.street`

### 4. Function Calls
- **Simple Calls**: `func()`, `max(a, b)`
- **Method Calls**: `object.method(arg1, arg2)`
- **Nested Calls**: `outer(inner(value))`

### 5. Pipe Expressions
- **Basic Pipes**: `value |: transform`
- **Typed Pipes**: `array |map: $1 * 2`
- **Pipe Chains**: `data |filter: $1 > 0 |map: $1 * 2 |reduce: $1 + $2`
- **Aliases**: `value as $temp |: process($temp)`

## Extension Points

### 1. Adding New Token Types

```go
// In constants/tokens.go
const (
    // ... existing tokens ...
    TokenCustom TokenType = iota + 100
)

// In tokenizer.go
func (t *Tokenizer) readCustomToken() (Token, error) {
    // Custom tokenization logic
    return Token{
        Type:  constants.TokenCustom,
        Value: customValue,
        Token: tokenString,
        Line:  t.line,
        Column: t.column,
    }, nil
}
```

### 2. Adding New Expression Types

```go
// Define new expression type
type CustomExpression struct {
    // Custom fields
    Line   int
    Column int
}

func (ce *CustomExpression) expressionNode()      {}
func (ce *CustomExpression) Position() (int, int) { return ce.Line, ce.Column }

// Add parsing method to parser
func (p *Parser) parseCustomExpression() Expression {
    // Custom parsing logic
    return &CustomExpression{
        Line:   p.current.Line,
        Column: p.current.Column,
    }
}
```

### 3. Custom Error Types

```go
// In errors/errors.go
const (
    ErrCustomError ErrorCode = "custom-error"
)

// Usage in parser
func (p *Parser) customParsingMethod() {
    if someCondition {
        p.addError(errors.ErrCustomError, "Custom error message")
        return
    }
}
```

## Performance Characteristics

### Time Complexity
- **Tokenization**: O(n) where n is input length
- **Parsing**: O(n) for most expressions, O(n²) worst case for deeply nested structures
- **Memory**: O(n) for AST storage, O(log n) for parser stack

### Memory Usage
- **Token Storage**: ~40 bytes per token (typical)
- **AST Nodes**: ~80-120 bytes per node (depending on type)
- **Parser State**: ~200 bytes (constant)

### Optimization Features
- **Single-pass Parsing**: No backtracking required
- **Minimal Allocations**: Efficient memory usage patterns
- **Position Caching**: Optimized position tracking

## Testing and Validation

The parser includes comprehensive test coverage:

### Test Categories
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Full parsing pipeline testing
3. **Error Tests**: Error condition validation
4. **Edge Case Tests**: Boundary condition testing
5. **Performance Tests**: Benchmark validation

### Running Tests
```bash
# Run all parser tests
go test ./parser/tests/... -v

# Run specific test categories
go test ./parser/tests/ -run TestParser -v
go test ./parser/tests/ -run TestError -v
go test ./parser/tests/ -run TestTokenizer -v

# Run with coverage
go test ./parser/tests/... -cover -v
```

## Conclusion

The UExl-Go parser represents a production-ready implementation that follows industry best practices for:

- **Error Handling**: Direct error propagation with rich context
- **Type Safety**: Strong typing throughout the parsing pipeline
- **Performance**: Efficient single-pass parsing with minimal allocations
- **Extensibility**: Well-defined extension points for customization
- **Maintainability**: Clear separation of concerns and comprehensive testing

The parser is suitable for production use in expression evaluation, data transformation, and rule engine applications.

