package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/maniartech/uexl_go/parser/constants"
	"github.com/maniartech/uexl_go/parser/errors"
)

type Token struct {
	Type           constants.TokenType
	Value          any    // Stores native parsed values (float64, string without quotes, etc.)
	Token          string // Stores the original token string
	Line           int
	Column         int
	IsSingleQuoted bool // Only set for constants.TokenString
}

// Pre-allocated common single character strings to avoid allocations
var charStrings [128]string

func init() {
	for i := 0; i < 128; i++ {
		charStrings[i] = string(rune(i))
	}
}

type Tokenizer struct {
	input  string
	pos    int
	line   int
	column int
	// cache of current rune to avoid repeated utf8.DecodeRuneInString calls
	curRune rune
	curSize int
	// reusable buffer for string unescaping to avoid allocations
	strBuf []byte
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s) at %d:%d", t.Type, t.Token, t.Line, t.Column)
}

func NewTokenizer(input string) *Tokenizer {
	tz := &Tokenizer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		strBuf: make([]byte, 0, 256), // preallocate buffer for string processing
	}
	tz.setCur()
	return tz
}

func (t *Tokenizer) NextToken() (Token, error) {
	t.skipWhitespace()

	if t.pos >= len(t.input) {
		return Token{Type: constants.TokenEOF, Line: t.line, Column: t.column}, nil
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
		return t.singleCharToken(constants.TokenLeftParen)
	case ch == ')':
		return t.singleCharToken(constants.TokenRightParen)
	case ch == '[':
		return t.singleCharToken(constants.TokenLeftBracket)
	case ch == ']':
		return t.singleCharToken(constants.TokenRightBracket)
	case ch == '{':
		return t.singleCharToken(constants.TokenLeftBrace)
	case ch == '}':
		return t.singleCharToken(constants.TokenRightBrace)
	case ch == ',':
		return t.singleCharToken(constants.TokenComma)
	case ch == '.':
		return t.singleCharToken(constants.TokenDot)
	case ch == ':':
		return t.singleCharToken(constants.TokenColon)
	case ch == '?':
		return t.readQuestionOrNullish()
	case ch == '|':
		return t.readPipeOrBitwiseOr()
	default:
		// Check if it's a valid operator character
		if isOperatorChar(ch) {
			return t.readOperator()
		}
		// Return actual error instead of error token
		t.advance() // consume the invalid character
		return Token{}, errors.NewParserError(
			errors.ErrInvalidCharacter,
			t.line,
			t.column,
			fmt.Sprintf("invalid character: '%c'", ch),
		)
	}
}

func (t *Tokenizer) readNumber() (Token, error) {
	start := t.pos
	startColumn := t.column

	// Read integer part (mandatory when entering here)
	// ASCII fast path: digits 0-9
	for t.pos < len(t.input) {
		c := t.input[t.pos]
		if c >= '0' && c <= '9' {
			t.pos++
			t.column++
			continue
		}
		break
	}

	intEnd := t.pos
	hasDot := false

	// Read optional fractional part: only if '.' is followed by a digit
	if t.pos < len(t.input) && t.input[t.pos] == '.' {
		// Peek next rune to ensure it's a digit; otherwise, treat '.' as a separate token (e.g., member access)
		if t.pos+1 < len(t.input) {
			next := t.input[t.pos+1]
			if next >= '0' && next <= '9' {
				// consume '.'
				t.pos++
				t.column++
				// consume fractional digits
				hasDot = true
				for t.pos < len(t.input) {
					c := t.input[t.pos]
					if c >= '0' && c <= '9' {
						t.pos++
						t.column++
						continue
					}
					break
				}
			}
		}
	}

	// Check for scientific notation - only if followed by proper exponent
	hasExp := false
	if t.pos < len(t.input) && (t.input[t.pos] == 'e' || t.input[t.pos] == 'E') {
		// Look ahead to see if this is a valid exponent
		savedPos := t.pos
		savedColumn := t.column
		t.pos++ // consume 'e' or 'E'
		t.column++

		// Optional sign
		if t.pos < len(t.input) && (t.input[t.pos] == '+' || t.input[t.pos] == '-') {
			t.pos++
			t.column++
		}

		// Must have at least one digit after e/E or optional sign
		if t.pos >= len(t.input) || !(t.input[t.pos] >= '0' && t.input[t.pos] <= '9') {
			// Not a valid scientific notation, backtrack
			t.pos = savedPos
			t.column = savedColumn
			t.setCur()
		} else {
			// Valid scientific notation, consume all digits
			for t.pos < len(t.input) {
				c := t.input[t.pos]
				if c >= '0' && c <= '9' {
					t.pos++
					t.column++
					continue
				}
				break
			}
			hasExp = true
		}
	}

	originalToken := t.input[start:t.pos]
	// Fast path: no exponent, parse simple int or decimal manually to avoid allocations
	if !hasExp {
		s := originalToken
		// Determine boundaries for int and frac parts
		intPart := s[:intEnd-start]
		var fracPart string
		if hasDot {
			fracPart = s[len(intPart)+1:]
		}

		// Limit digits to avoid overflow; fallback to strconv for very long parts
		if len(intPart) <= 18 && len(fracPart) <= 18 {
			var u uint64
			for i := 0; i < len(intPart); i++ {
				c := intPart[i]
				if c < '0' || c > '9' {
					// Shouldn't happen; fallback
					goto stdparse
				}
				u = u*10 + uint64(c-'0')
			}
			fv := float64(u)
			if len(fracPart) > 0 {
				var fu uint64
				for i := 0; i < len(fracPart); i++ {
					c := fracPart[i]
					if c < '0' || c > '9' {
						goto stdparse
					}
					fu = fu*10 + uint64(c-'0')
				}
				fv += float64(fu) / pow10[len(fracPart)]
			}
			t.setCur()
			return Token{Type: constants.TokenNumber, Value: fv, Token: originalToken, Line: t.line, Column: startColumn}, nil
		}
	}

stdparse:
	value, err := strconv.ParseFloat(originalToken, 64)
	if err != nil {
		// Invalid number format - return error
		errMsg := errors.GetErrorMessage(errors.ErrInvalidNumber)
		return Token{}, errors.NewParserError(errors.ErrInvalidNumber, t.line, startColumn, fmt.Sprintf("%s: '%s'", errMsg, originalToken))
	}
	t.setCur()
	return Token{Type: constants.TokenNumber, Value: value, Token: originalToken, Line: t.line, Column: startColumn}, nil
}

func (t *Tokenizer) readIdentifierOrKeyword() (Token, error) {
	start := t.pos
	startColumn := t.column

	// Allow the first character to be a letter, underscore, or dollar sign
	r := t.current()
	if isLetter(r) || r == '_' || r == '$' {
		// ASCII fast path for subsequent chars
		if t.pos < len(t.input) && t.input[t.pos] < 0x80 {
			// consume first ASCII char
			t.pos++
			t.column++
			for t.pos < len(t.input) {
				c := t.input[t.pos]
				if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$' {
					t.pos++
					t.column++
					continue
				}
				break
			}
			t.setCur()
		} else {
			t.advance()
		}
	}

	// Continue reading alphanumeric characters, underscores, and dollar signs
	// NO dot handling - dots will be separate tokens
	for t.pos < len(t.input) {
		ch := t.current()
		if isLetter(ch) || isDigit(ch) || ch == '_' || ch == '$' {
			// ASCII fast path
			if t.pos < len(t.input) && t.input[t.pos] < 0x80 {
				t.pos++
				t.column++
				continue
			}
			t.advance()
		} else {
			break
		}
	}

	originalToken := t.input[start:t.pos]
	switch originalToken {
	case "true", "false":
		return Token{Type: constants.TokenBoolean, Value: originalToken == "true", Token: originalToken, Line: t.line, Column: startColumn}, nil
	case "null":
		return Token{Type: constants.TokenNull, Value: nil, Token: originalToken, Line: t.line, Column: startColumn}, nil
	case "as":
		return Token{Type: constants.TokenAs, Value: originalToken, Token: originalToken, Line: t.line, Column: startColumn}, nil
	default:
		return Token{Type: constants.TokenIdentifier, Value: originalToken, Token: originalToken, Line: t.line, Column: startColumn}, nil
	}
}
func (t *Tokenizer) readString() (Token, error) {
	start := t.pos
	startColumn := t.column

	// Check for raw string prefix
	rawString := false
	if t.pos < len(t.input) && t.input[t.pos] == 'r' {
		rawString = true
		// Advance past 'r'
		t.pos++
		t.column++
		t.setCur()
	}

	// Get the quote character and advance past it
	quote := t.current()
	t.pos++
	t.column++
	if quote == '\n' {
		t.line++
		t.column = 1
	}
	t.setCur()

	// Read until closing quote
	if rawString {
		// Raw string: handle doubled quotes for escaping the delimiter
		for t.pos < len(t.input) {
			if t.input[t.pos] == byte(quote) {
				// Check if it's a doubled quote (escaped delimiter)
				if t.pos+1 < len(t.input) && t.input[t.pos+1] == byte(quote) {
					// Skip both quotes (this is an escaped delimiter)
					t.pos += 2
					t.column += 2
				} else {
					// Single quote - this is the end of the string
					break
				}
			} else {
				if t.input[t.pos] == '\n' {
					t.line++
					t.column = 1
				} else {
					t.column++
				}
				t.pos++
			}
		}
	} else {
		// Regular string: handle backslash escapes
		for t.pos < len(t.input) && t.input[t.pos] != byte(quote) {
			if t.input[t.pos] == '\\' {
				t.pos++ // skip escape character
				t.column++
				if t.pos < len(t.input) {
					if t.input[t.pos] == '\n' {
						t.line++
						t.column = 1
					} else {
						t.column++
					}
					t.pos++
				}
			} else {
				if t.input[t.pos] == '\n' {
					t.line++
					t.column = 1
				} else {
					t.column++
				}
				t.pos++
			}
		}
	}

	// Check if we found the closing quote
	if t.pos >= len(t.input) {
		// Unterminated string error
		errMsg := errors.GetErrorMessage(errors.ErrUnterminatedQuote)
		return Token{}, errors.NewParserError(errors.ErrUnterminatedQuote, t.line, startColumn, errMsg)
	}

	if t.pos < len(t.input) {
		t.pos++ // consume closing quote
		t.column++
		if t.input[t.pos-1] == '\n' {
			t.line++
			t.column = 1
		}
	}
	t.setCur()

	originalToken := t.input[start:t.pos]
	var value string
	isSingleQuoted := false
	if rawString {
		// For raw strings, extract content between quotes and handle doubled quotes
		content := originalToken[2 : len(originalToken)-1] // Remove 'r' prefix and quotes
		if quote == '"' {
			// Fast path: check if any doubled quotes exist first
			if !containsDoubledQuote(content, '"') {
				value = content
			} else {
				value = t.unescapeRawString(content, '"')
			}
		} else {
			isSingleQuoted = true
			// Fast path: check if any doubled quotes exist first
			if !containsDoubledQuote(content, '\'') {
				value = content
			} else {
				value = t.unescapeRawString(content, '\'')
			}
		}
	} else if quote == '"' {
		content := originalToken[1 : len(originalToken)-1] // Remove quotes
		if !containsEscape(content) {
			value = content
		} else {
			quoted := originalToken
			unescaped, err := strconv.Unquote(quoted)
			if err != nil {
				errMsg := errors.GetErrorMessage(errors.ErrInvalidString)
				return Token{}, errors.NewParserError(errors.ErrInvalidString, t.line, startColumn, errMsg+": '"+originalToken+"'")
			}
			value = unescaped
		}
	} else if quote == '\'' {
		isSingleQuoted = true
		content := originalToken[1 : len(originalToken)-1]
		if !containsEscape(content) {
			value = content
		} else {
			value = t.unescapeStringFast(content)
		}
	}
	return Token{Type: constants.TokenString, Value: value, Token: originalToken, Line: t.line, Column: startColumn, IsSingleQuoted: isSingleQuoted}, nil
}

func (t *Tokenizer) readPipeOrBitwiseOr() (Token, error) {
	startColumn := t.column
	t.advance() // consume first '|'
	if t.current() == '|' {
		t.advance() // consume second '|'
		operator := "||"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	if t.current() == ':' {
		t.advance() // consume ':'
		pipeValue := ":"
		return Token{Type: constants.TokenPipe, Value: pipeValue, Token: pipeValue, Line: t.line, Column: startColumn}, nil
	}

	// Manual scan for optional [A-Za-z]+ followed by ':'
	// We only accept ASCII letters for pipe names as before.
	i := t.pos
	for i < len(t.input) {
		r, size := utf8.DecodeRuneInString(t.input[i:])
		if r == utf8.RuneError && size == 1 {
			break
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			i += size
			continue
		}
		break
	}
	// If at least one letter and next char is ':' then it's a named pipe
	if i > t.pos && i < len(t.input) && t.input[i] == ':' {
		pipeName := t.input[t.pos:i]
		// Advance over the name and ':' without allocating
		for t.pos < i {
			t.advance()
		}
		t.advance() // consume ':'
		return Token{Type: constants.TokenPipe, Value: pipeName, Token: pipeName, Line: t.line, Column: startColumn}, nil
	}

	operator := "|"
	return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
}

func (t *Tokenizer) readOperator() (Token, error) {
	// This function does not handle operators starting with '|' because that is
	// handled by the readPipeOrBitwiseOr function.

	start := t.pos
	startColumn := t.column

	// Handle nullish coalescing operator '??'
	if t.current() == '?' && t.peek() == '?' {
		t.advance()
		t.advance()
		operator := "??"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle && operator
	if t.current() == '&' && t.peek() == '&' {
		t.advance()
		t.advance()
		operator := "&&"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle ++ operator
	if t.current() == '+' && t.peek() == '+' {
		t.advance()
		t.advance()
		operator := "++"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle ** operator (power)
	if t.current() == '*' && t.peek() == '*' {
		t.advance()
		t.advance()
		operator := "**"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle -- operator, but be smart about unary vs postfix contexts
	// For expressions like "--10", "---5", we want to treat as separate "-" tokens
	if t.current() == '-' && t.peek() == '-' {
		// Look ahead to see what comes after the second '-'
		nextChar := rune(0)
		if t.pos+2 < len(t.input) {
			nextChar, _ = utf8.DecodeRuneInString(t.input[t.pos+2:])
		}

		// If followed by a digit, letter (identifier), or another operator (like another -),
		// treat as separate minus tokens for unary contexts
		if isDigit(nextChar) || isLetter(nextChar) || isOperatorChar(nextChar) {
			// Return single '-' token
			t.advance()
			operator := "-"
			return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
		} else {
			// Treat as decrement operator (for cases like "i--" in postfix contexts)
			t.advance()
			t.advance()
			operator := "--"
			return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
		}
	}

	// Handle == operator
	if t.current() == '=' && t.peek() == '=' {
		t.advance()
		t.advance()
		operator := "=="
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle != operator
	if t.current() == '!' && t.peek() == '=' {
		t.advance()
		t.advance()
		operator := "!="
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle consecutive ! operators - always treat as separate tokens for multiple negation
	// For expressions like "!!true" or "!!!false", we want separate "!" tokens
	if t.current() == '!' && t.peek() == '!' {
		// Return single '!' token - let the parser handle multiple consecutive unary operators
		t.advance()
		operator := "!"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle <= operator
	if t.current() == '<' && t.peek() == '=' {
		t.advance()
		t.advance()
		operator := "<="
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle >= operator
	if t.current() == '>' && t.peek() == '=' {
		t.advance()
		t.advance()
		operator := ">="
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle << operator
	if t.current() == '<' && t.peek() == '<' {
		t.advance()
		t.advance()
		operator := "<<"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle >> operator
	if t.current() == '>' && t.peek() == '>' {
		t.advance()
		t.advance()
		operator := ">>"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	// Handle single-character operators
	t.advance()
	operator := t.input[start:t.pos]
	return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
}

// readQuestionOrNullish handles tokens that start with '?': either '?' or '??'
func (t *Tokenizer) readQuestionOrNullish() (Token, error) {
	startColumn := t.column
	// current is '?'
	t.advance()

	// Priority: '??' > '?.' > '?[' > '?'
	if t.current() == '?' {
		// it's '??'
		t.advance()
		operator := "??"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
	}

	if t.current() == '.' {
		// '?.'
		t.advance()
		tokenStr := "?."
		return Token{Type: constants.TokenQuestionDot, Value: tokenStr, Token: tokenStr, Line: t.line, Column: startColumn}, nil
	}

	if t.current() == '[' {
		// '?['
		t.advance()
		tokenStr := "?["
		return Token{Type: constants.TokenQuestionLeftBracket, Value: tokenStr, Token: tokenStr, Line: t.line, Column: startColumn}, nil
	}

	// single '?' -> used by conditional operator parsing
	operator := "?"
	return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
}

func (t *Tokenizer) singleCharToken(tokenType constants.TokenType) (Token, error) {
	ch := t.current()
	var charValue string
	if ch < 128 {
		charValue = charStrings[ch]
	} else {
		charValue = string(ch)
	}
	token := Token{Type: tokenType, Value: charValue, Token: charValue, Line: t.line, Column: t.column}
	t.advance()
	return token, nil
}

func (t *Tokenizer) current() rune {
	return t.curRune
}

// peek the next character without advancing the position
func (t *Tokenizer) peekNext() rune {
	if t.pos+1 >= len(t.input) {
		return 0
	}
	b := t.input[t.pos+1]
	if b < 0x80 {
		return rune(b)
	}
	r, _ := utf8.DecodeRuneInString(t.input[t.pos+1:])
	return r
}

func (t *Tokenizer) advance() {
	if t.pos >= len(t.input) {
		return
	}
	r := t.curRune
	size := t.curSize
	if r == '\n' {
		t.line++
		t.column = 1
	} else {
		t.column++
	}
	t.pos += size
	t.setCur()
}

func (t *Tokenizer) skipWhitespace() {
	// Fast path for common ASCII whitespace
	for t.pos < len(t.input) {
		if t.pos >= len(t.input) {
			return
		}
		c := t.input[t.pos]
		if c == ' ' || c == '\n' || c == '\t' || c == '\r' {
			t.advance()
			continue
		}
		// Fallback for non-ASCII spaces
		if t.current() != 0 && unicode.IsSpace(t.current()) {
			t.advance()
			continue
		}
		break
	}
}

func isDigit(r rune) bool {
	if r <= 127 {
		return r >= '0' && r <= '9'
	}
	return unicode.IsDigit(r)
}

func isLetter(r rune) bool {
	if r <= 127 {
		// ASCII letter or underscore
		rr := r | 32 // fold case for [A-Za-z]
		return (rr >= 'a' && rr <= 'z') || r == '_'
	}
	return unicode.IsLetter(r) || r == '_'
}

func isOperatorChar(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '%', '<', '>', '=', '!', '&', '|', '^', '?':
		return true
	default:
		return false
	}
}

func (t *Tokenizer) peek() rune {
	if t.pos+1 >= len(t.input) {
		return 0
	}
	b := t.input[t.pos+1]
	if b < 0x80 {
		return rune(b)
	}
	r, _ := utf8.DecodeRuneInString(t.input[t.pos+1:])
	return r
}

// setCur decodes and caches the rune at current position.
func (t *Tokenizer) setCur() {
	if t.pos >= len(t.input) {
		t.curRune = 0
		t.curSize = 0
		return
	}
	// ASCII fast path
	if t.input[t.pos] < 0x80 {
		t.curRune = rune(t.input[t.pos])
		t.curSize = 1
		return
	}
	r, size := utf8.DecodeRuneInString(t.input[t.pos:])
	t.curRune = r
	t.curSize = size
}

// containsEscape checks if a string contains backslash escapes without allocating
func containsEscape(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			return true
		}
	}
	return false
}

// containsDoubledQuote checks if a string contains doubled quotes without allocating
func containsDoubledQuote(s string, quote byte) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == quote && s[i+1] == quote {
			return true
		}
	}
	return false
}

// unescapeRawString handles doubled quotes in raw strings using the internal buffer
func (t *Tokenizer) unescapeRawString(content string, quote byte) string {
	t.strBuf = t.strBuf[:0] // reset buffer
	for i := 0; i < len(content); {
		if content[i] == quote && i+1 < len(content) && content[i+1] == quote {
			t.strBuf = append(t.strBuf, quote)
			i += 2
		} else {
			t.strBuf = append(t.strBuf, content[i])
			i++
		}
	}
	return string(t.strBuf)
}

// unescapeStringFast handles common escape sequences using the internal buffer
func (t *Tokenizer) unescapeStringFast(s string) string {
	// Try to use strconv.Unquote by converting single-quoted to double-quoted
	if s != "" {
		if unquoted, err := strconv.Unquote("\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""); err == nil {
			return unquoted
		}
	}

	// Fallback: manual unescaping using internal buffer
	t.strBuf = t.strBuf[:0] // reset buffer
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '\'':
				t.strBuf = append(t.strBuf, '\'')
				i += 2
			case '"':
				t.strBuf = append(t.strBuf, '"')
				i += 2
			case '\\':
				t.strBuf = append(t.strBuf, '\\')
				i += 2
			case 'n':
				t.strBuf = append(t.strBuf, '\n')
				i += 2
			case 't':
				t.strBuf = append(t.strBuf, '\t')
				i += 2
			case 'r':
				t.strBuf = append(t.strBuf, '\r')
				i += 2
			case 'b':
				t.strBuf = append(t.strBuf, '\b')
				i += 2
			case 'f':
				t.strBuf = append(t.strBuf, '\f')
				i += 2
			case 'a':
				t.strBuf = append(t.strBuf, '\a')
				i += 2
			case 'v':
				t.strBuf = append(t.strBuf, '\v')
				i += 2
			case '/':
				t.strBuf = append(t.strBuf, '/')
				i += 2
			default:
				// Unknown escape sequence, keep the backslash and next char
				t.strBuf = append(t.strBuf, s[i])
				t.strBuf = append(t.strBuf, s[i+1])
				i += 2
			}
		} else {
			t.strBuf = append(t.strBuf, s[i])
			i++
		}
	}
	return string(t.strBuf)
}

// Precomputed powers of 10 for fast decimal parsing
var pow10 = [...]float64{
	1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000,
	10000000000, 100000000000, 1000000000000, 10000000000000, 100000000000000,
	1000000000000000, 10000000000000000, 100000000000000000, 1000000000000000000,
}

// PreloadTokens preloads all tokens in the input string.
// This is a helper function for debugging and testing. It preloads all tokens in the input string
// and returns them as a slice. After this function is called, the tokenizer will be at the end of the input.
func (t *Tokenizer) PreloadTokens() []Token {
	tokens := []Token{}

	// Loop until the end of the input
	for {
		token, err := t.NextToken()
		if err != nil {
			// On error, create an error token for debugging purposes
			errorToken := Token{
				Type:   constants.TokenError,
				Value:  err.Error(),
				Token:  err.Error(),
				Line:   t.line,
				Column: t.column,
			}
			tokens = append(tokens, errorToken)
			break
		}
		tokens = append(tokens, token)
		if token.Type == constants.TokenEOF {
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
