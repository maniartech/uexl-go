package parser

import (
	"fmt"
	"strings"

	"github.com/maniartech/uexl_go/parser/errors"
)

// Constants for common pipe types
const (
	DefaultPipeType = "pipe"
)

// Parser represents a parser for the expression language
type Parser struct {
	tokenizer           *Tokenizer
	current             Token
	errors              []errors.ParserError // Changed from []string
	pos                 int
	subExpressionActive bool
	inParenthesis       bool
}

// IndexAccess represents array index access expressions
type IndexAccess struct {
	Array  Expression
	Index  Expression
	Line   int
	Column int
}

func (ia *IndexAccess) expressionNode()      {}
func (ia *IndexAccess) Position() (int, int) { return ia.Line, ia.Column }

// NewParser creates a new parser instance with the given input
func NewParser(input string) *Parser {
	p := &Parser{
		tokenizer: NewTokenizer(input),
	}
	p.advance()
	return p
}

// Parse parses the input and returns an Expression or an error
func (p *Parser) Parse() (Expression, error) {
	expr := p.parseExpression()

	if len(p.errors) > 0 {
		return nil, &errors.ParseErrors{Errors: p.errors}
	}

	if expr == nil {
		return nil, errors.NewParserError(errors.ErrEmptyExpression, p.current.Line, p.current.Column, "expression is nil")
	}

	if p.current.Type != TokenEOF {
		return nil, errors.NewParserErrorWithToken(errors.ErrUnexpectedToken, p.current.Line, p.current.Column, "unexpected token at end", p.current.Token)
	}

	return expr, nil
}

// parseExpression is the entry point for parsing expressions
func (p *Parser) parseExpression() Expression {
	return p.parsePipeExpression()
}

// parsePipeExpression parses pipe expressions with proper error handling
func (p *Parser) parsePipeExpression() Expression {
	// Only disallow pipe expressions in sub-expressions that aren't in parentheses
	if p.subExpressionActive && !p.inParenthesis {
		p.addError(errors.ErrPipeInSubExpression, errors.GetErrorMessage(errors.ErrPipeInSubExpression))
		return nil
	}

	// Handle leading pipe at the start of the expression
	if p.pos == 1 && p.current.Type == TokenPipe {
		return p.handleLeadingPipe()
	}

	firstExpression := p.parseLogicalOr()
	if firstExpression == nil {
		return nil
	}

	if p.subExpressionActive && !p.inParenthesis {
		return firstExpression
	}

	aliases := []string{}
	expressions := []Expression{firstExpression}
	pipeTypes := []string{DefaultPipeType}

	startLine, startColumn := expressions[0].Position()

	alias, e := p.parsePipeAlias()
	if e != nil {
		// Handle specific error types
		if parserErr, ok := e.(errors.ParserError); ok {
			p.errors = append(p.errors, parserErr)
		} else {
			p.addError(errors.ErrInvalidAlias, e.Error())
		}
		return nil
	}
	aliases = append(aliases, alias)

	for p.current.Type == TokenPipe {
		if !p.processPipeSegment(&expressions, &pipeTypes, &aliases) {
			return nil
		}
	}

	if len(expressions) == 1 {
		return expressions[0]
	}

	// Defensive check for consistency
	if len(expressions) > len(pipeTypes) {
		p.addError(errors.ErrMissingPipeType, errors.GetErrorMessage(errors.ErrMissingPipeType))
		return nil
	}

	return &PipeExpression{
		Expressions: expressions,
		PipeTypes:   pipeTypes,
		Aliases:     aliases,
		Line:        startLine,
		Column:      startColumn,
	}
}

func (p *Parser) parsePipeAlias() (string, error) {
	if p.current.Type == TokenAs {
		// If we're in a sub-expression (not top-level pipe), error
		if p.subExpressionActive {
			return "", errors.NewParserError(errors.ErrAliasInSubExpr, p.current.Line, p.current.Column, errors.GetErrorMessage(errors.ErrAliasInSubExpr))
		}
		p.advance() // consume 'as'

		if p.current.Type != TokenIdentifier || !strings.HasPrefix(p.current.Token, "$") {
			return "", errors.NewParserError(errors.ErrMissingDollarSign, p.current.Line, p.current.Column, errors.GetErrorMessage(errors.ErrMissingDollarSign))
		}

		alias := p.current.Token
		p.advance()

		return alias, nil
	}
	return "", nil
}

func (p *Parser) parseLogicalOr() Expression {
	return p.parseBinaryOp(p.parseLogicalAnd, "||")
}

func (p *Parser) parseLogicalAnd() Expression {
	return p.parseBinaryOp(p.parseBitwiseOr, "&&")
}

func (p *Parser) parseBitwiseOr() Expression {
	return p.parseBinaryOp(p.parseBitwiseXor, "|")
}

func (p *Parser) parseBitwiseXor() Expression {
	return p.parseBinaryOp(p.parseBitwiseAnd, "^")
}

func (p *Parser) parseBitwiseAnd() Expression {
	return p.parseBinaryOp(p.parseEquality, "&")
}

func (p *Parser) parseEquality() Expression {
	return p.parseBinaryOp(p.parseComparison, "==", "!=")
}

func (p *Parser) parseComparison() Expression {
	return p.parseBinaryOp(p.parseBitwiseShift, "<", ">", "<=", ">=")
}

func (p *Parser) parseBitwiseShift() Expression {
	return p.parseBinaryOp(p.parseAdditive, "<<", ">>")
}

func (p *Parser) parseAdditive() Expression {
	return p.parseBinaryOp(p.parseMultiplicative, "+", "-")
}

func (p *Parser) parseMultiplicative() Expression {
	return p.parseBinaryOp(p.parseUnary, "*", "/", "%")
}

func (p *Parser) parseUnary() Expression {
	if p.current.Type == TokenOperator {
		// Check if the operator is "-" or "!"
		if opValue, ok := p.current.Value.(string); ok && (opValue == "-" || opValue == "!") {
			op := p.current
			p.advance()
			expr := p.parseUnary()
			// Use the string value from the type assertion
			return &UnaryExpression{Operator: opValue, Operand: expr, Line: op.Line, Column: op.Column}
		}
	}
	return p.parseMemberAccess()
}

func (p *Parser) parseMemberAccess() Expression {
	expr := p.parsePrimary()

	for {
		// Handle array index access
		if p.current.Type == TokenLeftBracket {
			bracket := p.current
			p.advance() // consume '['

			// Save previous state
			wasInParenthesis := p.inParenthesis
			wasSubExpressionActive := p.subExpressionActive

			// Allow expressions within array index
			p.inParenthesis = true
			p.subExpressionActive = true

			// Parse the index expression
			indexExpr := p.parseExpression()

			if p.current.Type != TokenRightBracket {
				p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']' after array index", "]")
			} else {
				p.advance() // consume ']'
			}

			// Restore previous state
			p.inParenthesis = wasInParenthesis
			p.subExpressionActive = wasSubExpressionActive

			// Create an index access expression
			expr = &IndexAccess{
				Array:  expr,
				Index:  indexExpr,
				Line:   bracket.Line,
				Column: bracket.Column,
			}
			continue // check for more member access operations
		}

		// Handle dot access
		if p.current.Type == TokenDot {
			dot := p.current
			p.advance()

			// Check for end of input or unexpected token after dot
			if p.current.Type != TokenIdentifier {
				p.addErrorWithExpected(errors.ErrExpectedIdentifier, "expected identifier after '.'", "identifier")
				return expr // Return what we have so far since this is an error
			}

			property, ok := p.current.Value.(string)
			if !ok {
				property = p.current.Token
			}
			p.advance()

			expr = &MemberAccess{
				Object:   expr,
				Property: property,
				Line:     dot.Line,
				Column:   dot.Column,
			}
			continue // check for more member access operations
		}

		// Handle function call after member or index access
		if p.current.Type == TokenLeftParen {
			expr = p.parseFunctionCall(expr)
			continue
		}

		// No more member access operations
		break
	}

	return expr
}

func (p *Parser) parsePrimary() Expression {
	switch p.current.Type {
	case TokenNumber:
		return p.parseNumber()
	case TokenString:
		return p.parseString()
	case TokenBoolean:
		return p.parseBoolean()
	case TokenNull:
		return p.parseNull()
	case TokenIdentifier:
		return p.parseIdentifierOrFunctionCall()
	case TokenLeftParen:
		return p.parseGroupedExpression()
	case TokenLeftBracket:
		return p.parseArray()
	case TokenLeftBrace:
		return p.parseObject()
	default:
		p.addErrorWithToken(errors.ErrUnexpectedToken, "unexpected token")
		p.advance()
		return nil
	}
}

func (p *Parser) parseIdentifierOrFunctionCall() Expression {
	identifier := &Identifier{
		Name:   p.current.Token,
		Line:   p.current.Line,
		Column: p.current.Column,
	}
	p.advance()

	if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(identifier)
	}

	return identifier
}

func (p *Parser) parseFunctionCall(function Expression) Expression {
	openParen := p.current
	p.advance() // consume '('

	// Save previous state
	wasInParenthesis := p.inParenthesis
	wasSubExpressionActive := p.subExpressionActive

	// Set flags for function arguments
	p.inParenthesis = true // Allow parenthesized expressions in function args
	p.subExpressionActive = true

	args := []Expression{}

	if p.current.Type != TokenRightParen {
		for {
			args = append(args, p.parseExpression())
			if p.current.Type != TokenComma {
				break
			}
			p.advance() // consume ','
		}
	}

	if p.current.Type != TokenRightParen {
		p.addErrorWithExpected(errors.ErrUnclosedFunction, "expected ')' after function arguments", ")")
		return nil
	}
	p.advance() // consume ')'

	// Restore previous state
	p.inParenthesis = wasInParenthesis
	p.subExpressionActive = wasSubExpressionActive

	return &FunctionCall{
		Function:  function,
		Arguments: args,
		Line:      openParen.Line,
		Column:    openParen.Column,
	}
}

func (p *Parser) parseNumber() Expression {
	token := p.current
	p.advance()
	return &NumberLiteral{Value: token.Token, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseString() Expression {
	token := p.current
	p.advance()
	// Value already has quotes removed
	value, ok := token.Value.(string)
	if !ok {
		// Fallback to removing quotes manually if type assertion fails
		value = strings.Trim(token.Token, "'\"")
	}

	// Check if this was a raw string by looking at the original token
	isRaw := strings.HasPrefix(token.Token, "r'") || strings.HasPrefix(token.Token, "r\"")
	isSingleQuoted := false
	if isRaw {
		// For raw strings, check the second character
		if len(token.Token) > 1 && token.Token[1] == '\'' {
			isSingleQuoted = true
		}
	} else {
		// For regular strings, check the first character
		if len(token.Token) > 0 && token.Token[0] == '\'' {
			isSingleQuoted = true
		}
	}

	return &StringLiteral{
		Value:          value,
		Token:          token.Token,
		IsRaw:          isRaw,
		IsSingleQuoted: isSingleQuoted,
		Line:           token.Line,
		Column:         token.Column,
	}
}

func (p *Parser) parseBoolean() Expression {
	token := p.current
	p.advance()
	value, ok := token.Value.(bool)
	if !ok {
		value = token.Token == "true"
	}
	return &BooleanLiteral{Value: value, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseNull() Expression {
	token := p.current
	p.advance()
	return &NullLiteral{Line: token.Line, Column: token.Column}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.advance() // consume '('

	// Set both flags to track that we're in a parenthesized expression
	wasInParenthesis := p.inParenthesis
	p.inParenthesis = true
	wasSubExpressionActive := p.subExpressionActive
	p.subExpressionActive = true

	expr := p.parseExpression()
	if p.current.Type != TokenRightParen {
		p.addErrorWithExpected(errors.ErrExpectedToken, "expected ')'", ")")
	} else {
		p.advance() // consume ')'
	}

	// Restore previous state
	p.inParenthesis = wasInParenthesis
	p.subExpressionActive = wasSubExpressionActive

	return expr
}

func (p *Parser) parseArray() Expression {
	token := p.current
	p.advance() // consume '['

	// Save previous state
	wasInParenthesis := p.inParenthesis
	wasSubExpressionActive := p.subExpressionActive

	// Set flags for array elements
	p.inParenthesis = true
	p.subExpressionActive = true

	elements := []Expression{}

	if p.current.Type != TokenRightBracket {
		for {
			elements = append(elements, p.parseExpression())
			if p.current.Type != TokenComma {
				break
			}
			p.advance() // consume ','

			// Check for trailing comma
			if p.current.Type == TokenRightBracket {
				p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']'", "]")
				break
			}
		}
	}

	if p.current.Type != TokenRightBracket {
		p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']'", "]")
	} else {
		p.advance() // consume ']'
	}

	// Restore previous state
	p.inParenthesis = wasInParenthesis
	p.subExpressionActive = wasSubExpressionActive

	return &ArrayLiteral{Elements: elements, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseObject() Expression {
	token := p.current
	p.advance() // consume '{'

	// Save previous state
	wasInParenthesis := p.inParenthesis
	wasSubExpressionActive := p.subExpressionActive

	// Set flags for object properties
	p.inParenthesis = true
	p.subExpressionActive = true

	properties := make(map[string]Expression)
	for p.current.Type != TokenRightBrace {
		if p.current.Type != TokenString {
			p.addErrorWithExpected(errors.ErrInvalidObjectKey, "expected string key", "string")
			break
		}

		// Use the Value field which should contain the unquoted string
		var key string
		if strValue, ok := p.current.Value.(string); ok && strValue != "" {
			key = strValue
		} else {
			// Fallback: remove quotes manually if needed
			key = strings.Trim(p.current.Token, "'\"")
		}

		// Advance past the string token
		p.advance()

		if p.current.Type != TokenColon {
			p.addErrorWithExpected(errors.ErrExpectedToken, "expected ':'", ":")
			break
		}
		p.advance()

		value := p.parseExpression()
		if value == nil {
			p.addError(errors.ErrInvalidObjectValue, fmt.Sprintf("invalid value for key '%s'", key))
			break
		}

		properties[key] = value

		if p.current.Type != TokenComma {
			break
		}
		p.advance() // consume ','

		// Check for trailing comma
		if p.current.Type == TokenRightBrace {
			break
		}
	}

	if p.current.Type != TokenRightBrace {
		p.addErrorWithExpected(errors.ErrUnclosedObject, "expected '}'", "}")
	} else {
		p.advance() // consume '}'
	}

	// Restore previous state
	p.inParenthesis = wasInParenthesis
	p.subExpressionActive = wasSubExpressionActive

	return &ObjectLiteral{Properties: properties, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseBinaryOp(parseFunc func() Expression, operators ...string) Expression {
	left := parseFunc()

	for p.current.Type == TokenOperator {
		// Type assertion for operator value
		opValue, ok := p.current.Value.(string)
		if !ok || !contains(operators, opValue) {
			break
		}
		op := p.current
		p.advance()
		right := parseFunc()
		left = &BinaryExpression{Left: left, Operator: opValue, Right: right, Line: op.Line, Column: op.Column}
	}

	return left
}

// advance moves to the next token in the input
func (p *Parser) advance() {
	p.current = p.tokenizer.NextToken()
	p.pos++

	// Check if the tokenizer returned an error token
	if p.current.Type == TokenError {
		if errorCode, ok := p.current.Value.(errors.ErrorCode); ok {
			p.addError(errorCode, p.current.Token)
		} else {
			p.addError(errors.ErrInvalidToken, p.current.Token)
		}
	}
}

// addError adds an error message with current position information
func (p *Parser) addError(code errors.ErrorCode, message string) {
	p.errors = append(p.errors, errors.NewParserError(code, p.current.Line, p.current.Column, message))
}

// addErrorWithToken adds an error with token information
func (p *Parser) addErrorWithToken(code errors.ErrorCode, message string) {
	p.errors = append(p.errors, errors.NewParserErrorWithToken(code, p.current.Line, p.current.Column, message, p.current.Token))
}

// addErrorWithExpected adds an error with expected token information
func (p *Parser) addErrorWithExpected(code errors.ErrorCode, message, expected string) {
	p.errors = append(p.errors, errors.NewParserErrorWithExpected(code, p.current.Line, p.current.Column, message, p.current.Token, expected))
}

func (p *Parser) handleLeadingPipe() Expression {
	p.advance() // consume the pipe
	// Check if it's followed by 'as' (empty pipe with alias)
	if p.current.Type == TokenAs {
		p.addErrorWithToken(errors.ErrEmptyPipeWithAlias, errors.GetErrorMessage(errors.ErrEmptyPipeWithAlias))
		p.advance() // consume 'as'
		if p.current.Type == TokenIdentifier {
			p.advance() // consume alias identifier
		}
	} else {
		p.addError(errors.ErrEmptyPipe, errors.GetErrorMessage(errors.ErrEmptyPipe))
	}
	p.consumeRemainingTokens()
	return nil
}

// processPipeSegment processes a single pipe segment and returns false if parsing should stop
func (p *Parser) processPipeSegment(expressions *[]Expression, pipeTypes *[]string, aliases *[]string) bool {
	op := p.current
	p.advance()

	pipeType := p.determinePipeType(op)
	*pipeTypes = append(*pipeTypes, pipeType)

	// Check for empty pipe with alias immediately after consuming pipe
	if p.current.Type == TokenAs {
		p.addErrorWithToken(errors.ErrEmptyPipeWithAlias, errors.GetErrorMessage(errors.ErrEmptyPipeWithAlias))
		p.advance() // consume 'as'
		if p.current.Type == TokenIdentifier {
			p.advance() // consume alias identifier
		}
		p.consumeRemainingTokens()
		return false
	}

	nextExpr := p.parseLogicalOr()
	if nextExpr == nil {
		p.addError(errors.ErrEmptyPipe, errors.GetErrorMessage(errors.ErrEmptyPipe))
		p.consumeRemainingTokens()
		return false
	}

	*expressions = append(*expressions, nextExpr)
	alias, e := p.parsePipeAlias()
	if e != nil {
		// Handle specific error types
		if parserErr, ok := e.(errors.ParserError); ok {
			p.errors = append(p.errors, parserErr)
		} else {
			p.addError(errors.ErrInvalidAlias, e.Error())
		}
		return false
	}
	*aliases = append(*aliases, alias)
	return true
}

// determinePipeType extracts the pipe type from the pipe token
func (p *Parser) determinePipeType(op Token) string {
	if op.Value != nil {
		if strValue, ok := op.Value.(string); ok && strValue != "" {
			// If the value is just ":", treat as default pipe
			// This allows syntax like |: to be interpreted as a normal pipe.
			if strValue == ":" {
				return DefaultPipeType
			}
			return strValue
		}
	}
	return DefaultPipeType
}

// consumeRemainingTokens consumes all tokens until EOF to prevent further errors
func (p *Parser) consumeRemainingTokens() {
	for p.current.Type != TokenEOF {
		p.advance()
	}
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
