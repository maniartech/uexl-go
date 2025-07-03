package parser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/maniartech/uexl_go/parser/errors"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdentifier
	TokenOperator
	TokenLeftParen
	TokenRightParen
	TokenLeftBracket
	TokenRightBracket
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenColon
	TokenPipe
	TokenString
	TokenBoolean
	TokenNull
	TokenDollar
	TokenAs
	TokenError // Special token type for tokenizer errors
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "Number"
	case TokenIdentifier:
		return "Identifier"
	case TokenOperator:
		return "Operator"
	case TokenLeftParen:
		return "LeftParen"
	case TokenRightParen:
		return "RightParen"
	case TokenLeftBracket:
		return "LeftBracket"
	case TokenRightBracket:
		return "RightBracket"
	case TokenLeftBrace:
		return "LeftBrace"
	case TokenRightBrace:
		return "RightBrace"
	case TokenComma:
		return "Comma"
	case TokenDot:
		return "Dot"
	case TokenColon:
		return "Colon"
	case TokenPipe:
		return "Pipe"
	case TokenString:
		return "String"
	case TokenBoolean:
		return "Boolean"
	case TokenNull:
		return "Null"
	case TokenDollar:
		return "Dollar"
	case TokenAs:
		return "As"
	case TokenError:
		return "Error"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type   TokenType
	Value  any    // Stores native parsed values (float64, string without quotes, etc.)
	Token  string // Stores the original token string
	Line   int
	Column int
}

type Tokenizer struct {
	input  string
	pos    int
	line   int
	column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s) at %d:%d", t.Type, t.Token, t.Line, t.Column)
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

func (t *Tokenizer) NextToken() Token {
	t.skipWhitespace()

	if t.pos >= len(t.input) {
		return Token{Type: TokenEOF, Line: t.line, Column: t.column}
	}

	switch ch := t.current(); {
	case isDigit(ch):
		return t.readNumber()
	case ch == '$':
		return t.readIdentifierOrKeyword()
	case ch == '"' || ch == '\'':
		return t.readString()
	case ch == 'r':
		if t.peekNext() == '"' || t.peekNext() == '\'' {
			return t.readString()
		}
		fallthrough
	case isLetter(ch):
		return t.readIdentifierOrKeyword()
	case ch == '(':
		return t.singleCharToken(TokenLeftParen)
	case ch == ')':
		return t.singleCharToken(TokenRightParen)
	case ch == '[':
		return t.singleCharToken(TokenLeftBracket)
	case ch == ']':
		return t.singleCharToken(TokenRightBracket)
	case ch == '{':
		return t.singleCharToken(TokenLeftBrace)
	case ch == '}':
		return t.singleCharToken(TokenRightBrace)
	case ch == ',':
		return t.singleCharToken(TokenComma)
	case ch == '.':
		return t.singleCharToken(TokenDot)
	case ch == ':':
		return t.singleCharToken(TokenColon)
	case ch == '|':
		return t.readPipeOrBitwiseOr()
	default:
		// Check if it's a valid operator character
		if isOperatorChar(ch) {
			return t.readOperator()
		}
		// Invalid character - create error token
		errMsg := errors.GetErrorMessage(errors.ErrInvalidCharacter)
		token := Token{
			Type:   TokenError,
			Value:  errors.ErrInvalidCharacter,
			Token:  fmt.Sprintf("%s: '%c'", errMsg, ch),
			Line:   t.line,
			Column: t.column,
		}
		t.advance() // consume the invalid character
		return token
	}
}

func (t *Tokenizer) readNumber() Token {
	start := t.pos
	for t.pos < len(t.input) && (isDigit(t.current()) || t.current() == '.') {
		t.advance()
	}
	if t.pos < len(t.input) && (t.current() == 'e' || t.current() == 'E') {
		t.advance()
		if t.pos < len(t.input) && (t.current() == '+' || t.current() == '-') {
			t.advance()
		}
		for t.pos < len(t.input) && isDigit(t.current()) {
			t.advance()
		}
	}
	originalToken := t.input[start:t.pos]
	value, err := strconv.ParseFloat(originalToken, 64)
	if err != nil {
		// Invalid number format - create error token
		errMsg := errors.GetErrorMessage(errors.ErrInvalidNumber)
		return Token{
			Type:   TokenError,
			Value:  errors.ErrInvalidNumber,
			Token:  fmt.Sprintf("%s: '%s'", errMsg, originalToken),
			Line:   t.line,
			Column: t.column - (t.pos - start),
		}
	}
	return Token{Type: TokenNumber, Value: value, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}
}

func (t *Tokenizer) readIdentifierOrKeyword() Token {
	start := t.pos
	hasDot := false

	// Allow the first character to be a letter, underscore, or dollar sign
	if isLetter(t.current()) || t.current() == '_' || t.current() == '$' {
		t.advance()
	}

	for t.pos < len(t.input) && (isLetter(t.current()) || isDigit(t.current()) || t.current() == '_' || t.current() == '.') {
		if t.current() == '.' {
			if hasDot {
				// Error: consecutive dots - create an error token
				errMsg := errors.GetErrorMessage(errors.ErrConsecutiveDots)
				return Token{
					Type:   TokenError,
					Value:  errors.ErrConsecutiveDots,
					Token:  errMsg,
					Line:   t.line,
					Column: t.column,
				}
			}
			hasDot = true
		} else {
			hasDot = false
		}
		t.advance()
	}

	originalToken := t.input[start:t.pos]
	switch originalToken {
	case "true", "false":
		return Token{Type: TokenBoolean, Value: originalToken == "true", Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}
	case "null":
		return Token{Type: TokenNull, Value: nil, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}
	case "as":
		return Token{Type: TokenAs, Value: originalToken, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}
	default:
		return Token{Type: TokenIdentifier, Value: originalToken, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}
	}
}
func (t *Tokenizer) readString() Token {
	start := t.pos
	startColumn := t.column

	// Check for raw string prefix
	rawString := false
	if t.input[t.pos] == 'r' {
		rawString = true
		// Advance past 'r'
		t.advance()
	}

	// Get the quote character and advance past it
	quote := t.current()
	t.advance() // consume opening quote

	// Read until closing quote
	for t.pos < len(t.input) && t.current() != quote {
		if !rawString && t.current() == '\\' {
			t.advance() // skip escape character
		}
		t.advance()
	}

	// Check if we found the closing quote
	if t.pos >= len(t.input) {
		// Unterminated string error
		errMsg := errors.GetErrorMessage(errors.ErrUnterminatedQuote)
		return Token{
			Type:   TokenError,
			Value:  errors.ErrUnterminatedQuote,
			Token:  errMsg,
			Line:   t.line,
			Column: startColumn,
		}
	}

	if t.pos < len(t.input) {
		t.advance() // consume closing quote
	}

	var value string
	originalToken := t.input[start:t.pos]
	if rawString {
		// Remove 'r' prefix with quotes from the original token
		value = originalToken[2 : len(originalToken)-1]
	} else {
		value = originalToken[1 : len(originalToken)-1]

		// Unescape the string
		value = strings.ReplaceAll(value, "\\\\", "\\")
		value = strings.ReplaceAll(value, "\\n", "\n")
		value = strings.ReplaceAll(value, "\\t", "\t")
		value = strings.ReplaceAll(value, "\\\"", "\"")
		value = strings.ReplaceAll(value, "\\'", "'")
	}

	return Token{Type: TokenString, Value: value, Token: originalToken, Line: t.line, Column: startColumn}
}

// Ref: https://regex101.com/r/w6qtHq/1
var pipePattern = regexp.MustCompile(`(?m)^(?P<pipe>[a-z]+)?:`)

func (t *Tokenizer) readPipeOrBitwiseOr() Token {
	start := t.pos
	t.advance() // consume first '|'
	if t.current() == '|' {
		t.advance() // consume second '|'
		operator := "||"
		return Token{Type: TokenOperator, Value: operator, Token: operator, Line: t.line, Column: t.column - (t.pos - start)}
	}

	if t.current() == ':' {
		t.advance() // consume ':'
		pipeValue := ":"
		return Token{Type: TokenPipe, Value: pipeValue, Token: pipeValue, Line: t.line, Column: t.column - (t.pos - start)}
	}

	// Fetch the next 10 characters or the rest of the input if less than 10 characters are available
	nextChars := t.input[t.pos:int(
		math.Min(float64(t.pos+10), float64(len(t.input))),
	)]

	pipeMatch := pipePattern.FindStringSubmatch(nextChars)
	if len(pipeMatch) > 0 {
		pipeName := pipeMatch[1]
		for range pipeName {
			t.advance()
		}
		t.advance() // consume ':'
		return Token{Type: TokenPipe, Value: pipeName, Token: pipeName, Line: t.line, Column: t.column - (t.pos - start)}
	}

	operator := "|"
	return Token{Type: TokenOperator, Value: operator, Token: operator, Line: t.line, Column: t.column - (t.pos - start)}
}

func (t *Tokenizer) readOperator() Token {
	// This function does not handle operators starting with '|' because that is
	// handled by the readPipeOrBitwiseOr function.

	start := t.pos
	startColumn := t.column

	// Handle && operator
	if t.current() == '&' && t.peek() == '&' {
		t.advance()
		t.advance()
		operator := "&&"
		return Token{Type: TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}
	}

	// Handle single-character operators
	for t.pos < len(t.input) && isOperatorChar(t.current()) {
		t.advance()
	}
	operator := t.input[start:t.pos]
	return Token{Type: TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}
}

func (t *Tokenizer) singleCharToken(tokenType TokenType) Token {
	charValue := string(t.current())
	token := Token{Type: tokenType, Value: charValue, Token: charValue, Line: t.line, Column: t.column}
	t.advance()
	return token
}

func (t *Tokenizer) current() rune {
	if t.pos >= len(t.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(t.input[t.pos:])
	return r
}

// peek the next character without advancing the position
func (t *Tokenizer) peekNext() rune {
	if t.pos+1 >= len(t.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(t.input[t.pos+1:])
	return r
}

func (t *Tokenizer) advance() {
	if t.pos >= len(t.input) {
		return
	}
	if t.input[t.pos] == '\n' {
		t.line++
		t.column = 1
	} else {
		t.column++
	}
	t.pos += utf8.RuneLen(t.current())
}

func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.input) && unicode.IsSpace(t.current()) {
		t.advance()
	}
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isOperatorChar(r rune) bool {
	return strings.ContainsRune("+-*/%<>=!&|^", r)
}

func (t *Tokenizer) peek() rune {
	if t.pos+1 >= len(t.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(t.input[t.pos+1:])
	return r
}

// PreloadTokens preloads all tokens in the input string.
// This is a helper function for debugging and testing. It preloads all tokens in the input string
// and returns them as a slice. After this function is called, the tokenizer will be at the end of the input.
func (t *Tokenizer) PreloadTokens() []Token {
	tokens := []Token{}

	// Loop until the end of the input
	for {
		token := t.NextToken()
		tokens = append(tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}
	return tokens
}

func (t *Tokenizer) PrintTokens() {
	for i, t := range t.PreloadTokens() {
		fmt.Printf("%d: %s\n", i, t)
	}
}
