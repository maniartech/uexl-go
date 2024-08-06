package parser

import (
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
)

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
	column := t.column
	for t.pos < len(t.input) && (isLetter(t.current()) || isDigit(t.current()) || t.current() == '_') {
		t.advance()
	}
	value := t.input[start:t.pos]
	switch value {
	case "true", "false":
		return Token{Type: TokenBoolean, Value: value, Line: t.line, Column: column}
	case "null":
		return Token{Type: TokenNull, Value: value, Line: t.line, Column: column}
	default:
		return Token{Type: TokenIdentifier, Value: value, Line: t.line, Column: column}
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

func (t *Tokenizer) readPipeOrBitwiseOr() Token {
	// start := t.pos
	t.advance()
	if t.current() == ':' {
		t.advance()
		return Token{Type: TokenPipe, Value: "|:", Line: t.line, Column: t.column - 2}
	}
	return Token{Type: TokenOperator, Value: "|", Line: t.line, Column: t.column - 1}
}

func (t *Tokenizer) readOperator() Token {
	start := t.pos
	column := t.column

	// Check for two-character operators
	if t.pos+1 < len(t.input) {
		twoCharOp := t.input[t.pos : t.pos+2]
		if twoCharOp == "&&" || twoCharOp == "||" || twoCharOp == "==" || twoCharOp == "!=" || twoCharOp == "<=" || twoCharOp == ">=" || twoCharOp == "<<" || twoCharOp == ">>" || twoCharOp == "|:" {
			t.advance()
			t.advance()
			if twoCharOp == "|:" {
				return Token{Type: TokenPipe, Value: twoCharOp, Line: t.line, Column: column}
			}
			return Token{Type: TokenOperator, Value: twoCharOp, Line: t.line, Column: column}
		}
	}

	// Single-character operators
	t.advance()
	return Token{Type: TokenOperator, Value: t.input[start:t.pos], Line: t.line, Column: column}
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
