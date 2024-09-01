package parser

import (
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
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
	default:
		return "Unknown"
	}
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Tokenizer struct {
	input  string
	pos    int
	line   int
	column int
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
	case isLetter(ch):
		return t.readIdentifierOrKeyword()
	case ch == '"' || ch == '\'':
		return t.readString()
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
	case ch == '$':
		return t.readDollarIdentifier()
	case ch == '|':
		return t.readPipeOrBitwiseOr()
	default:
		return t.readOperator()
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
	return Token{Type: TokenNumber, Value: t.input[start:t.pos], Line: t.line, Column: t.column - (t.pos - start)}
}

func (t *Tokenizer) readIdentifierOrKeyword() Token {
	start := t.pos
	for t.pos < len(t.input) && (isLetter(t.current()) || isDigit(t.current()) || t.current() == '_') {
		t.advance()
	}
	value := t.input[start:t.pos]
	switch value {
	case "true", "false":
		return Token{Type: TokenBoolean, Value: value, Line: t.line, Column: t.column - (t.pos - start)}
	case "null":
		return Token{Type: TokenNull, Value: value, Line: t.line, Column: t.column - (t.pos - start)}
	default:
		return Token{Type: TokenIdentifier, Value: value, Line: t.line, Column: t.column - (t.pos - start)}
	}
}

func (t *Tokenizer) readString() Token {
	quote := t.current()
	start := t.pos
	t.advance() // consume opening quote
	for t.pos < len(t.input) && t.current() != quote {
		if t.current() == '\\' {
			t.advance() // skip escape character
		}
		t.advance()
	}
	if t.pos < len(t.input) {
		t.advance() // consume closing quote
	}
	return Token{Type: TokenString, Value: t.input[start:t.pos], Line: t.line, Column: t.column - (t.pos - start)}
}

func (t *Tokenizer) readDollarIdentifier() Token {
	start := t.pos
	t.advance() // consume '$'
	for t.pos < len(t.input) && isDigit(t.current()) {
		t.advance()
	}
	return Token{Type: TokenIdentifier, Value: t.input[start:t.pos], Line: t.line, Column: t.column - (t.pos - start)}
}

var pipePattern = regexp.MustCompile(`(?m)^(?P<pipe>[a-z]+)?:`)

func (t *Tokenizer) readPipeOrBitwiseOr() Token {
	t.advance() // consume first '|'
	if t.current() == '|' {
		t.advance() // consume second '|'
		return Token{Type: TokenOperator, Value: "||", Line: t.line, Column: t.column - 2}
	}

	// Fetch the next 10 characters or the rest of the input if less than 10 characters are available
	nextChars := t.input[t.pos:int(math.Min(float64(t.pos+10), float64(len(t.input))))]

	pipeMatch := pipePattern.FindStringSubmatch(nextChars)
	if len(pipeMatch) > 0 {
		pipeName := pipeMatch[1]
		for range pipeName {
			t.advance()
		}
		t.advance() // consume ':'
		return Token{Type: TokenPipe, Value: pipeName, Line: t.line, Column: t.column - len(pipeMatch[0]) - 1}
	}

	return Token{Type: TokenOperator, Value: "|", Line: t.line, Column: t.column - 1}
}

func (t *Tokenizer) readOperator() Token {
	start := t.pos
	switch t.current() {
	case '&':
		if t.peek() == '&' {
			t.advance()
			t.advance()
			return Token{Type: TokenOperator, Value: "&&", Line: t.line, Column: t.column - 2}
		}
	case '|':
		if t.peek() == '|' {
			t.advance()
			t.advance()
			return Token{Type: TokenOperator, Value: "||", Line: t.line, Column: t.column - 2}
		}
	}
	// Handle single-character operators
	for t.pos < len(t.input) && isOperatorChar(t.current()) {
		t.advance()
	}
	return Token{Type: TokenOperator, Value: t.input[start:t.pos], Line: t.line, Column: t.column - (t.pos - start)}
}

func (t *Tokenizer) singleCharToken(tokenType TokenType) Token {
	token := Token{Type: tokenType, Value: string(t.current()), Line: t.line, Column: t.column}
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
