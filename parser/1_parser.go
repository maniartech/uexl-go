package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokenizer           *Tokenizer
	current             Token
	errors              []string
	pos                 int
	subExpressionActive bool
	inParenthesis       bool // Add flag to track if we're inside parentheses
}

func NewParser(input string) *Parser {
	p := &Parser{
		tokenizer: NewTokenizer(input),
	}
	p.advance()
	return p
}

func (p *Parser) Parse() (Expression, error) {
	expr := p.parseExpression()

	if len(p.errors) != 0 {
		return nil, fmt.Errorf("parsing errors: %v", p.errors)
	}

	if expr == nil {
		return nil, fmt.Errorf("expression is nil")
	}

	if p.current.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token at end: %v", p.current)
	}
	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %v", p.errors)
	}
	return expr, nil
}

func (p *Parser) parseExpression() Expression {
	return p.parsePipeExpression()
}

func (p *Parser) parsePipeExpression() Expression {
	// Only disallow pipe expressions in sub-expressions that aren't in parentheses
	if p.subExpressionActive && !p.inParenthesis {
		p.addError("pipe expressions cannot be sub-expressions")
		return nil
	}

	firstExpression := p.parseLogicalOr()

	if firstExpression == nil {
		return nil
	}

	// Don't process pipes if inside a subexpression that's not a parenthesized expr
	if p.subExpressionActive && !p.inParenthesis {
		return firstExpression
	}

	expressions := []Expression{firstExpression}
	pipeTypes := []string{"pipe"}

	startLine, startColumn := expressions[0].Position()

	for p.current.Type == TokenPipe {
		op := p.current
		p.advance()

		pipeType := "pipe" // default pipe type
		if op.Value != nil {
			// Convert op.Value to string if possible
			if strValue, ok := op.Value.(string); ok && strValue != "" {
				pipeType = strValue
			}
		}
		pipeTypes = append(pipeTypes, pipeType)

		expressions = append(expressions, p.parseLogicalOr())
	}

	if len(expressions) == 1 {
		return expressions[0]
	}

	if len(expressions) > len(pipeTypes) {
		// insert default pipe type at the beginning of the pipeTypes slice
		pipeTypes = append([]string{"pipe"}, pipeTypes...)
	}

	return &PipeExpression{
		Expressions: expressions,
		PipeTypes:   pipeTypes,
		Line:        startLine,
		Column:      startColumn,
	}
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
				p.addError("expected ']' after array index")
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
				p.addError("expected identifier after '.'")
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
		p.addError(fmt.Sprintf("unexpected token: %v", p.current))
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
		p.addError("expected ')' after function arguments")
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
	return &StringLiteral{Value: value, Line: token.Line, Column: token.Column}
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
		p.addError("expected ')'")
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
				p.addError("expected ']'")
				break
			}
		}
	}

	if p.current.Type != TokenRightBracket {
		p.addError("expected ']'")
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
			p.addError("expected string key")
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
			p.addError("expected ':'")
			break
		}
		p.advance()

		value := p.parseExpression()
		if value == nil {
			p.addError(fmt.Sprintf("invalid value for key '%s'", key))
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
		p.addError("expected '}'")
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

func (p *Parser) advance() {
	p.current = p.tokenizer.NextToken()
	if p.pos == 0 {
		if p.current.Type == TokenPipe {
			p.addError("unexpected token: " + p.current.Token)
		}
	}
	p.pos++
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: %s", p.current.Line, p.current.Column, msg))
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Add a new expression type for array index access
type IndexAccess struct {
	Array  Expression
	Index  Expression
	Line   int
	Column int
}

func (ia *IndexAccess) expressionNode()      {}
func (ia *IndexAccess) Position() (int, int) { return ia.Line, ia.Column }
