package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/maniartech/uexl_go/parser/errors"
)

// Constants for common pipe types
const (
	DefaultPipeType = constants.DefaultPipeType
)

// Parser represents a parser for the expression language
type Parser struct {
	tokenizer           *Tokenizer
	current             Token
	errors              []errors.ParserError // Changed from []string
	pos                 int
	subExpressionActive bool
	inParenthesis       bool
	options             Options
}

// NewParser creates a new parser instance with the given input
// This function maintains backward compatibility and does not return errors
func NewParser(input string) *Parser {
	return NewParserWithOptions(input, DefaultOptions())
}

// NewParserWithOptions creates a new parser instance with the given input and options
func NewParserWithOptions(input string, opt Options) *Parser {
	p := &Parser{
		tokenizer: NewTokenizerWithOptions(input, opt),
		options:   opt,
	}
	p.advance()
	return p
}

// NewParserWithValidation creates a new parser instance with the given input
// and returns an error if the input is invalid or initialization fails
func NewParserWithValidation(input string) (*Parser, error) {
	if input == "" {
		return nil, errors.NewParserError(errors.ErrEmptyExpression, 1, 1, constants.MsgEmptyExpression)
	}

	p := &Parser{
		tokenizer: NewTokenizerWithOptions(input, DefaultOptions()),
	}

	// Get the first token and check for immediate errors
	p.advance()

	// If the first token is an error, return it immediately
	if p.current.Type == constants.TokenError {
		// Error tokens store error message in Str field
		return nil, errors.NewParserError(errors.ErrInvalidToken, p.current.Line, p.current.Column, p.current.Value.Str)
	}

	return p, nil
}

// Parse parses the input and returns an Expression or an error
func (p *Parser) Parse() (Expression, error) {
	expr := p.parseExpression()

	if len(p.errors) > 0 {
		return nil, &errors.ParseErrors{Errors: p.errors}
	}

	if expr == nil {
		return nil, errors.NewParserError(errors.ErrEmptyExpression, p.current.Line, p.current.Column, constants.MsgEmptyExpression)
	}

	if p.current.Type != constants.TokenEOF {
		return nil, errors.NewParserErrorWithToken(errors.ErrUnexpectedToken, p.current.Line, p.current.Column, constants.MsgUnexpectedToken, p.current.Token)
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
	if p.pos == 1 && p.current.Type == constants.TokenPipe {
		return p.handleLeadingPipe()
	}

	firstExpression := p.parseConditional()
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

	for p.current.Type == constants.TokenPipe {
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

	programNode := &ProgramNode{
		PipeExpressions: make([]PipeExpression, 0, len(expressions)),
		Line:            startLine,
		Column:          startColumn,
	}

	for i, expr := range expressions {
		if expr == nil {
			p.addError(errors.ErrInvalidExpression, fmt.Sprintf("nil expression at index %d", i))
			continue
		}
		programNode.PipeExpressions = append(programNode.PipeExpressions, PipeExpression{
			Expression: expr,
			PipeType:   pipeTypes[i],
			Alias:      aliases[i],
			Index:      i,
			Line:       startLine,
			Column:     startColumn,
		})
	}

	return programNode
}

// parseConditional parses the ternary conditional operator and integrates with nullish coalescing/logical or
// Precedence: logical-or / nullish (same level) > conditional (?:) > pipe
// Associativity: conditional is right-associative
func (p *Parser) parseConditional() Expression {
	condition := p.parseLogicalOr()
	if condition == nil {
		return nil
	}

	// Check for '?'
	if p.current.Type == constants.TokenOperator {
		if p.current.Value.Kind == TVKOperator && p.current.Value.Str == "?" {
			qmark := p.current
			p.advance() // consume '?'

			// Parse consequent using conditional to support nesting (right-assoc)
			consequent := p.parseConditional()
			if consequent == nil {
				p.addErrorWithExpected(errors.ErrExpectedToken, "expected expression after '?'", ": or expression")
			}

			if p.current.Type != constants.TokenColon {
				p.addErrorWithExpected(errors.ErrExpectedToken, "expected ':' in conditional expression", ":")
			} else {
				p.advance() // consume ':'
			}

			alternate := p.parseConditional()
			if alternate == nil {
				p.addError(errors.ErrInvalidExpression, "invalid expression after ':' in conditional")
			}

			return &ConditionalExpression{
				Condition:  condition,
				Consequent: consequent,
				Alternate:  alternate,
				Line:       qmark.Line,
				Column:     qmark.Column,
			}
		}
	}

	return condition
}

func (p *Parser) parsePipeAlias() (string, error) {
	if p.current.Type == constants.TokenAs {
		// If we're in a sub-expression (not top-level pipe), error
		if p.subExpressionActive {
			return "", errors.NewParserError(errors.ErrAliasInSubExpr, p.current.Line, p.current.Column, errors.GetErrorMessage(errors.ErrAliasInSubExpr))
		}
		p.advance() // consume 'as'

		if p.current.Type != constants.TokenIdentifier || !strings.HasPrefix(p.current.Token, constants.SymbolDollar) {
			return "", errors.NewParserError(errors.ErrMissingDollarSign, p.current.Line, p.current.Column, errors.GetErrorMessage(errors.ErrMissingDollarSign))
		}

		alias := p.current.Token
		p.advance()

		return alias, nil
	}
	return "", nil
}

func (p *Parser) parseLogicalOr() Expression {
	// Logical OR handles only ||; nullish (??) is at a tighter precedence level
	return p.parseBinaryOp(p.parseLogicalAnd, constants.SymbolLogicalOr)
}

func (p *Parser) parseLogicalAnd() Expression {
	return p.parseBinaryOp(p.parseBitwiseOr, constants.SymbolLogicalAnd)
}

func (p *Parser) parseBitwiseOr() Expression {
	return p.parseBinaryOp(p.parseBitwiseXor, constants.SymbolBitwiseOr)
}

func (p *Parser) parseBitwiseXor() Expression {
	return p.parseBinaryOp(p.parseBitwiseAnd, constants.SymbolBitwiseXor)
}

func (p *Parser) parseBitwiseAnd() Expression {
	return p.parseBinaryOp(p.parseEquality, constants.SymbolBitwiseAnd)
}

func (p *Parser) parseEquality() Expression {
	return p.parseBinaryOp(p.parseComparison, constants.SymbolEqual, constants.SymbolNotEqual)
}

func (p *Parser) parseComparison() Expression {
	// Comparison is looser than nullish; parse nullish first
	return p.parseBinaryOp(p.parseNullish, "<", ">", "<=", ">=")
}

func (p *Parser) parseBitwiseShift() Expression {
	return p.parseBinaryOp(p.parseAdditive, "<<", ">>")
}

// parseNullish parses the nullish coalescing operator (??)
// Precedence: looser than arithmetic/shift, tighter than comparison/equality and all logical ops
// Associativity: left-to-right
func (p *Parser) parseNullish() Expression {
	return p.parseBinaryOp(p.parseBitwiseShift, "??")
}

func (p *Parser) parseAdditive() Expression {
	return p.parseBinaryOp(p.parseMultiplicative, "+", "-")
}

func (p *Parser) parseMultiplicative() Expression {
	return p.parseBinaryOp(p.parsePower, "*", "/", "%")
}

func (p *Parser) parsePower() Expression {
	// Power operator is right-associative and has higher precedence than unary
	// To mirror JavaScript precedence, if the left side is a chain of unary operators
	// (e.g., -, !), and we see '**', we must lift those unary operators to wrap the
	// entire exponent expression:  -2**3  =>  -(2**3),  !!2**3  =>  !!(2**3)

	left := p.parseUnary()

	if p.current.Type == constants.TokenOperator {
		if p.current.Value.Kind == TVKOperator && p.current.Value.Str == "**" {
			op := p.current
			p.advance()

			// Peel leading unary operators from the left side if present
			ops, base := peelLeadingUnary(left)

			// Right-associative parse for the right operand
			right := p.parsePower()

			powerExpr := &BinaryExpression{Left: base, Operator: op.Value.Str, Right: right, Line: op.Line, Column: op.Column}

			// Re-apply peeled unary operators around the power expression (outermost first)
			if len(ops) > 0 {
				expr := Expression(powerExpr)
				for _, uop := range ops {
					expr = &UnaryExpression{Operator: uop, Operand: expr, Line: op.Line, Column: op.Column}
				}
				return expr
			}

			return powerExpr
		}
	}

	return left
}

// peelLeadingUnary extracts a sequence of leading unary operators ('-' and '!')
// from the given expression. It returns the operators in outer-to-inner order
// and the innermost non-unary base expression.
func peelLeadingUnary(expr Expression) ([]string, Expression) {
	ops := []string{}

	for {
		if ue, ok := expr.(*UnaryExpression); ok {
			if ue.Operator == "-" || ue.Operator == "!" {
				ops = append(ops, ue.Operator)
				expr = ue.Operand
				continue
			}
		}
		break
	}
	return ops, expr
}

func (p *Parser) parseUnary() Expression {
	if p.current.Type == constants.TokenOperator {
		// Check if the operator is "-" or "!"
		if p.current.Value.Kind == TVKOperator && (p.current.Value.Str == "-" || p.current.Value.Str == "!") {
			op := p.current
			opStr := op.Value.Str
			p.advance()
			expr := p.parseUnary()
			return &UnaryExpression{Operator: opStr, Operand: expr, Line: op.Line, Column: op.Column}
		}
	}
	return p.parseMemberAccess()
}

func (p *Parser) parseMemberAccess() Expression {
	expr := p.parsePrimary()

	for {
		// Handle array index/slice access
		if p.current.Type == constants.TokenLeftBracket || p.current.Type == constants.TokenQuestionLeftBracket {
			expr = p.parseIndexOrSliceExpression(expr)
			continue // check for more member access operations
		}

		// Handle dot access
		if p.current.Type == constants.TokenDot || p.current.Type == constants.TokenQuestionDot {
			dot := p.current
			p.advance()

			// Disambiguate after '.'
			// 1) .<identifier> => MemberAccess (object/property)
			// 2) .<number>     => MemberAccess (property access, runtime converts if needed)
			// 3) .(expr)       => IndexAccess with arbitrary expression
			// Special-case: '?.[' optional element access (JS-style). Treat as optional index.
			if dot.Type == constants.TokenQuestionDot && p.current.Type == constants.TokenLeftBracket {
				bracket := p.current
				p.advance() // consume '['

				// Save and set flags for index expression
				wasInParenthesis := p.inParenthesis
				wasSubExpressionActive := p.subExpressionActive
				p.inParenthesis = true
				p.subExpressionActive = true

				indexExpr := p.parseExpression()

				if p.current.Type != constants.TokenRightBracket {
					p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']' after array index", "]")
				} else {
					p.advance() // consume ']'
				}

				// Restore flags
				p.inParenthesis = wasInParenthesis
				p.subExpressionActive = wasSubExpressionActive

				expr = &IndexAccess{
					Target:   expr,
					Index:    indexExpr,
					Optional: true,
					Line:     bracket.Line,
					Column:   bracket.Column,
				}
				continue
			}

			switch p.current.Type {
			case constants.TokenIdentifier:
				property := p.current.Value.Str
				if property == "" {
					property = p.current.Token
				}
				p.advance()
				expr = &MemberAccess{
					Target:   expr,
					Property: PropS(property),
					Optional: dot.Type == constants.TokenQuestionDot,
					Line:     dot.Line,
					Column:   dot.Column,
				}
				continue

			case constants.TokenNumber:
				// Treat as property access: obj.<number> → MemberAccess
				// If the numeric token contains a decimal point (e.g., "1.5"),
				// split it into multiple member accesses: obj.1.5 => ((obj).1).5
				tok := p.current
				p.advance()

				tokenStr := tok.Token
				if strings.Contains(tokenStr, ".") {
					parts := strings.Split(tokenStr, ".")
					for _, part := range parts {
						if part == "" {
							// Skip empty segments defensively; shouldn't occur for valid numbers
							continue
						}
						// Parse each segment as integer
						// We avoid using the float64 value to preserve exact segment semantics
						var segVal int
						// Support leading zeros (e.g., "05")
						// Atoi handles them fine
						if iv, err := strconv.Atoi(part); err == nil {
							segVal = iv
						} else {
							// Fallback: attempt via float then cast
							if tok.Value.Kind == TVKNumber {
								segVal = int(tok.Value.Num)
							} else {
								// As a last resort, keep as raw string
								expr = &MemberAccess{
									Target:   expr,
									Property: PropS(part),
									Optional: dot.Type == constants.TokenQuestionDot,
									Line:     dot.Line,
									Column:   dot.Column,
								}
								continue
							}
						}
						expr = &MemberAccess{
							Target:   expr,
							Property: PropI(segVal),
							Optional: dot.Type == constants.TokenQuestionDot,
							Line:     dot.Line,
							Column:   dot.Column,
						}
					}
				} else {
					// Simple integer-like number
					var prop Property
					if tok.Value.Kind == TVKNumber {
						prop = PropI(int(tok.Value.Num))
					} else {
						prop = PropS(tok.Token)
					}
					expr = &MemberAccess{
						Target:   expr,
						Property: prop,
						Optional: dot.Type == constants.TokenQuestionDot,
						Line:     dot.Line,
						Column:   dot.Column,
					}
				}
				continue

			case constants.TokenLeftParen:
				// Treat .(expr) as index access using grouped expression
				// Save previous state similar to bracket indexing to allow full expressions
				wasInParenthesis := p.inParenthesis
				wasSubExpressionActive := p.subExpressionActive
				p.inParenthesis = true
				p.subExpressionActive = true

				indexExpr := p.parseGroupedExpression()

				// Restore previous state
				p.inParenthesis = wasInParenthesis
				p.subExpressionActive = wasSubExpressionActive

				expr = &IndexAccess{
					Target:   expr,
					Index:    indexExpr,
					Optional: false,
					Line:     dot.Line,
					Column:   dot.Column,
				}
				continue

			case constants.TokenQuestionLeftBracket:
				// This handles cases like `obj.?[index]`
				expr = p.parseIndexOrSliceExpression(expr)
				continue

			default:
				// Unexpected token after '.'
				// For optional dot access, only identifier is valid per spec
				if dot.Type == constants.TokenQuestionDot {
					p.addErrorWithExpected(errors.ErrExpectedIdentifier, "expected identifier after '?.'", "identifier")
				} else {
					p.addErrorWithExpected(errors.ErrExpectedIdentifier, "expected identifier, number, or '(...)' after '.'", "identifier|number|(")
				}
				return expr
			}
		}

		// Handle function call after member or index access
		if p.current.Type == constants.TokenLeftParen {
			switch expr.(type) {
			case *Identifier, *FunctionCall:
				expr = p.parseFunctionCall(expr)
				continue // allow chaining after function call
			default:
				p.addErrorWithToken(errors.ErrInvalidSyntax, "function calls are only allowed after identifiers or function calls, not after member access or index access")
				return expr
			}
		}

		// No more member access operations
		break
	}

	return expr
}

func (p *Parser) parseIndexOrSliceExpression(target Expression) Expression {
	bracket := p.current
	p.advance() // consume '[' or '?['

	// Save previous state
	wasInParenthesis := p.inParenthesis
	wasSubExpressionActive := p.subExpressionActive

	// Allow expressions within array index/slice
	p.inParenthesis = true
	p.subExpressionActive = true

	// Peek ahead to see if this is a slice. A slice must contain at least one ':'.
	// We can't just check p.current, because it could be `[1:2]`.
	// A simple way is to check if there is a colon before the closing bracket.
	// This requires a lookahead in the tokenizer. Since we don't have that,
	// we will parse based on what we see.

	var start, end, step Expression

	// Check for the end of the slice/index right away
	if p.current.Type == constants.TokenRightBracket {
		p.addErrorWithExpected(errors.ErrUnclosedArray, "expected expression or ':' in slice/index", "expression or ':'")
		p.advance() // consume ']' to avoid infinite loops
		return nil
	}

	// It's a slice if the first token is a colon, e.g., `[:end]`
	if p.current.Type == constants.TokenColon {
		p.advance() // consume first ':'
		// `start` is nil
		if p.current.Type != constants.TokenRightBracket && p.current.Type != constants.TokenColon {
			end = p.parseExpression()
		}

		if p.current.Type == constants.TokenColon {
			p.advance() // consume second ':'
			if p.current.Type != constants.TokenRightBracket {
				step = p.parseExpression()
			}
		}
	} else {
		// It could be an index `[index]` or a slice `[start:...]`
		start = p.parseExpression()

		if p.current.Type == constants.TokenColon {
			// It's a slice
			p.advance() // consume first ':'

			if p.current.Type != constants.TokenRightBracket && p.current.Type != constants.TokenColon {
				end = p.parseExpression()
			}

			if p.current.Type == constants.TokenColon {
				p.advance() // consume second ':'
				if p.current.Type != constants.TokenRightBracket {
					step = p.parseExpression()
				}
			}
		} else {
			// It's a simple index access
			if p.current.Type != constants.TokenRightBracket {
				p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']' after array index", "]")
			} else {
				p.advance() // consume ']'
			}

			// Restore previous state
			p.inParenthesis = wasInParenthesis
			p.subExpressionActive = wasSubExpressionActive

			return &IndexAccess{
				Target:   target,
				Index:    start,
				Optional: bracket.Type == constants.TokenQuestionLeftBracket,
				Line:     bracket.Line,
				Column:   bracket.Column,
			}
		}
	}

	if p.current.Type != constants.TokenRightBracket {
		p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']' after slice", "]")
	} else {
		p.advance() // consume ']'
	}

	// Restore previous state
	p.inParenthesis = wasInParenthesis
	p.subExpressionActive = wasSubExpressionActive

	return &SliceExpression{
		Target:   target,
		Start:    start,
		End:      end,
		Step:     step,
		Optional: bracket.Type == constants.TokenQuestionLeftBracket,
		Line:     bracket.Line,
		Column:   bracket.Column,
	}
}

func (p *Parser) parsePrimary() Expression {
	switch p.current.Type {
	case constants.TokenNumber:
		return p.parseNumber()
	case constants.TokenString:
		return p.parseString()
	case constants.TokenBoolean:
		return p.parseBoolean()
	case constants.TokenNull:
		return p.parseNull()
	case constants.TokenIdentifier:
		return p.parseIdentifierOrFunctionCall()
	case constants.TokenLeftParen:
		return p.parseGroupedExpression()
	case constants.TokenLeftBracket:
		return p.parseArray()
	case constants.TokenLeftBrace:
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

	// Don't handle function calls here - let parseMemberAccess handle them
	// so we can detect chaining
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

	if p.current.Type != constants.TokenRightParen {
		for {
			args = append(args, p.parseExpression())
			if p.current.Type != constants.TokenComma {
				break
			}
			p.advance() // consume ','
		}
	}

	if p.current.Type != constants.TokenRightParen {
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
	return &NumberLiteral{Value: token.Value.Num, Line: token.Line, Column: token.Column}
}

func (p *Parser) parseString() Expression {
	token := p.current
	p.advance()
	// Value already has quotes removed
	value := token.Value.Str
	if value == "" {
		// Fallback to removing quotes manually if empty
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
	value := token.Value.Bool
	if token.Value.Kind != TVKBoolean {
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
	if p.current.Type != constants.TokenRightParen {
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

	if p.current.Type != constants.TokenRightBracket {
		for {
			elements = append(elements, p.parseExpression())
			if p.current.Type != constants.TokenComma {
				break
			}
			p.advance() // consume ','

			// Check for trailing comma
			if p.current.Type == constants.TokenRightBracket {
				p.addErrorWithExpected(errors.ErrUnclosedArray, "expected ']'", "]")
				break
			}
		}
	}

	if p.current.Type != constants.TokenRightBracket {
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
	for p.current.Type != constants.TokenRightBrace {
		if p.current.Type != constants.TokenString {
			p.addErrorWithExpected(errors.ErrInvalidObjectKey, "expected string key", "string")
			break
		}

		// Use the Value field which should contain the unquoted string
		var key string
		if p.current.Value.Kind == TVKString && p.current.Value.Str != "" {
			key = p.current.Value.Str
		} else {
			// Fallback: remove quotes manually if needed
			key = strings.Trim(p.current.Token, "'\"")
		}

		// Advance past the string token
		p.advance()

		if p.current.Type != constants.TokenColon {
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

		if p.current.Type != constants.TokenComma {
			break
		}
		p.advance() // consume ','

		// Check for trailing comma
		if p.current.Type == constants.TokenRightBrace {
			break
		}
	}

	if p.current.Type != constants.TokenRightBrace {
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

	for p.current.Type == constants.TokenOperator {
		// Check operator value
		if p.current.Value.Kind != TVKOperator || !contains(operators, p.current.Value.Str) {
			break
		}
		op := p.current
		opValue := op.Value.Str
		p.advance()
		right := parseFunc()
		left = &BinaryExpression{Left: left, Operator: opValue, Right: right, Line: op.Line, Column: op.Column}
	}

	return left
}

// advance moves to the next token in the input
func (p *Parser) advance() {
	token, err := p.tokenizer.NextToken()
	if err != nil {
		// Convert error to ParserError if it's not already
		if parserErr, ok := err.(errors.ParserError); ok {
			p.errors = append(p.errors, parserErr)
		} else {
			p.errors = append(p.errors, errors.NewParserError(
				errors.ErrInvalidToken,
				p.tokenizer.line,
				p.tokenizer.column,
				err.Error(),
			))
		}
		// Set current to an error token so parsing can continue
		p.current = Token{
			Type:   constants.TokenError,
			Value:  TokenValue{Kind: TVKString, Str: err.Error()},
			Token:  err.Error(),
			Line:   p.tokenizer.line,
			Column: p.tokenizer.column,
		}
		return
	}

	p.current = token
	p.pos++
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
	if p.current.Type == constants.TokenAs {
		p.addErrorWithToken(errors.ErrEmptyPipeWithAlias, errors.GetErrorMessage(errors.ErrEmptyPipeWithAlias))
		p.advance() // consume 'as'
		if p.current.Type == constants.TokenIdentifier {
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
	if p.current.Type == constants.TokenAs {
		p.addErrorWithToken(errors.ErrEmptyPipeWithAlias, errors.GetErrorMessage(errors.ErrEmptyPipeWithAlias))
		p.advance() // consume 'as'
		if p.current.Type == constants.TokenIdentifier {
			p.advance() // consume alias identifier
		}
		p.consumeRemainingTokens()
		return false
	}

	nextExpr := p.parseConditional()
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
	if op.Value.Kind == TVKString && op.Value.Str != "" {
		strValue := op.Value.Str
		// If the value is just ":", treat as default pipe
		// This allows syntax like |: to be interpreted as a normal pipe.
		if strValue == ":" {
			return DefaultPipeType
		}
		return strValue
	}
	return DefaultPipeType
}

// consumeRemainingTokens consumes all tokens until EOF to prevent further errors
func (p *Parser) consumeRemainingTokens() {
	for p.current.Type != constants.TokenEOF {
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
