package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokenizer *Tokenizer
	current   Token
	errors    []string
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
	left := p.parseLogicalOr()

	for p.current.Type == TokenPipe {
		op := p.current
		p.advance()

		var pipeType string
		if p.current.Type == TokenIdentifier {
			pipeType = p.current.Value
			p.advance()
		}

		if p.current.Type != TokenColon {
			p.addError("expected ':' after pipe operator")
			return left
		}
		p.advance()

		right := p.parseExpression()
		left = &PipeExpression{Left: left, Right: right, PipeType: pipeType, Line: op.Line, Column: op.Column}
	}

	return left
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
	if p.current.Type == TokenOperator && (p.current.Value == "-" || p.current.Value == "!") {
		op := p.current
		p.advance()
		expr := p.parseUnary()
		return &UnaryExpression{Operator: op.Value, Operand: expr, Line: op.Line, Column: op.Column}
	}
	return p.parseMemberAccess()
}

func (p *Parser) parseMemberAccess() Expression {
	expr := p.parsePrimary()

	for p.current.Type == TokenDot {
		dot := p.current
		p.advance()
		if p.current.Type != TokenIdentifier {
			p.addError("expected identifier after '.'")
			return expr
		}
		property := p.current.Value
		p.advance()
		expr = &MemberAccess{Object: expr, Property: property, Line: dot.Line, Column: dot.Column}
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
	identifier := &Identifier{Name: p.current.Value, Line: p.current.Line, Column: p.current.Column}
	p.advance()

	if p.current.Type == TokenLeftParen {
		return p.parseFunctionCall(identifier)
	}

	return identifier
}

func (p *Parser) parseFunctionCall(function Expression) Expression {
	openParen := p.current
	p.advance() // consume '('
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

	return &FunctionCall{Function: function, Arguments: args, Line: openParen.Line, Column: openParen.Column}
}

func (p *Parser) parseNumber() Expression {
	token := p.current
	p.advance()
	return &NumberLiteral{Value: token.Value, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseString() Expression {
	token := p.current
	p.advance()
	// Remove surrounding quotes
	value := strings.Trim(token.Value, "'\"")
	return &StringLiteral{Value: value, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseBoolean() Expression {
	token := p.current
	p.advance()
	value := token.Value == "true"
	return &BooleanLiteral{Value: value, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseNull() Expression {
	token := p.current
	p.advance()
	return &NullLiteral{Line: token.Line, Column: token.Column}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.advance() // consume '('
	expr := p.parseExpression()
	if p.current.Type != TokenRightParen {
		p.addError("expected ')'")
	} else {
		p.advance() // consume ')'
	}
	return expr
}

func (p *Parser) parseArray() Expression {
	token := p.current
	p.advance() // consume '['
	elements := []Expression{}
	for p.current.Type != TokenRightBracket {
		elements = append(elements, p.parseExpression())
		if p.current.Type == TokenComma {
			p.advance()
		} else {
			break
		}
	}
	if p.current.Type != TokenRightBracket {
		p.addError("expected ']'")
	} else {
		p.advance() // consume ']'
	}
	return &ArrayLiteral{Elements: elements, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseObject() Expression {
	token := p.current
	p.advance() // consume '{'
	properties := make(map[string]Expression)
	for p.current.Type != TokenRightBrace {
		if p.current.Type != TokenString {
			p.addError("expected string key")
			break
		}
		key := p.current.Value
		p.advance()
		if p.current.Type != TokenColon {
			p.addError("expected ':'")
			break
		}
		p.advance()
		value := p.parseExpression()
		properties[key] = value
		if p.current.Type == TokenComma {
			p.advance()
		} else {
			break
		}
	}
	if p.current.Type != TokenRightBrace {
		p.addError("expected '}'")
	} else {
		p.advance() // consume '}'
	}
	return &ObjectLiteral{Properties: properties, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseBinaryOp(parseFunc func() Expression, operators ...string) Expression {
	left := parseFunc()

	for p.current.Type == TokenOperator && contains(operators, p.current.Value) {
		op := p.current
		p.advance()
		right := parseFunc()
		left = &BinaryExpression{Left: left, Operator: op.Value, Right: right, Line: op.Line, Column: op.Column}
	}

	return left
}

func (p *Parser) advance() {
	p.current = p.tokenizer.NextToken()
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
