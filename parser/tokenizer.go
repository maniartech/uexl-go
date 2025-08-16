package parser

import (
	"fmt"
	"math"
	"regexp"
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

	// Read integer part and decimal part
	for t.pos < len(t.input) && (isDigit(t.current()) || t.current() == '.') {
		t.advance()
	}

	// Check for scientific notation - only if followed by proper exponent
	if t.pos < len(t.input) && (t.current() == 'e' || t.current() == 'E') {
		// Look ahead to see if this is a valid exponent
		savedPos := t.pos
		t.advance() // consume 'e' or 'E'

		// Optional sign
		if t.pos < len(t.input) && (t.current() == '+' || t.current() == '-') {
			t.advance()
		}

		// Must have at least one digit after e/E or optional sign
		if t.pos >= len(t.input) || !isDigit(t.current()) {
			// Not a valid scientific notation, backtrack
			t.pos = savedPos
		} else {
			// Valid scientific notation, consume all digits
			for t.pos < len(t.input) && isDigit(t.current()) {
				t.advance()
			}
		}
	}

	originalToken := t.input[start:t.pos]
	value, err := strconv.ParseFloat(originalToken, 64)
	if err != nil {
		// Invalid number format - return error
		errMsg := errors.GetErrorMessage(errors.ErrInvalidNumber)
		return Token{}, errors.NewParserError(errors.ErrInvalidNumber, t.line, t.column-(t.pos-start), fmt.Sprintf("%s: '%s'", errMsg, originalToken))
	}
	return Token{Type: constants.TokenNumber, Value: value, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}, nil
}

func (t *Tokenizer) readIdentifierOrKeyword() (Token, error) {
	start := t.pos

	// Allow the first character to be a letter, underscore, or dollar sign
	if isLetter(t.current()) || t.current() == '_' || t.current() == '$' {
		t.advance()
	}

	// Continue reading alphanumeric characters, underscores, and dollar signs
	// NO dot handling - dots will be separate tokens
	for t.pos < len(t.input) {
		ch := t.current()

		if isLetter(ch) || isDigit(ch) || ch == '_' || ch == '$' {
			t.advance()
		} else {
			break
		}
	}

	originalToken := t.input[start:t.pos]
	switch originalToken {
	case "true", "false":
		return Token{Type: constants.TokenBoolean, Value: originalToken == "true", Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}, nil
	case "null":
		return Token{Type: constants.TokenNull, Value: nil, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}, nil
	case "as":
		return Token{Type: constants.TokenAs, Value: originalToken, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}, nil
	default:
		return Token{Type: constants.TokenIdentifier, Value: originalToken, Token: originalToken, Line: t.line, Column: t.column - (t.pos - start)}, nil
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
		t.advance()
	}

	// Get the quote character and advance past it
	quote := t.current()
	t.advance() // consume opening quote

	// Read until closing quote
	if rawString {
		// Raw string: handle doubled quotes for escaping the delimiter
		for t.pos < len(t.input) {
			if t.current() == quote {
				// Check if it's a doubled quote (escaped delimiter)
				if t.pos+1 < len(t.input) && rune(t.input[t.pos+1]) == quote {
					// Skip both quotes (this is an escaped delimiter)
					t.advance()
					t.advance()
				} else {
					// Single quote - this is the end of the string
					break
				}
			} else {
				t.advance()
			}
		}
	} else {
		// Regular string: handle backslash escapes
		for t.pos < len(t.input) && t.current() != quote {
			if t.current() == '\\' {
				t.advance() // skip escape character
			}
			t.advance()
		}
	}

	// Check if we found the closing quote
	if t.pos >= len(t.input) {
		// Unterminated string error
		errMsg := errors.GetErrorMessage(errors.ErrUnterminatedQuote)
		return Token{}, errors.NewParserError(errors.ErrUnterminatedQuote, t.line, startColumn, errMsg)
	}

	if t.pos < len(t.input) {
		t.advance() // consume closing quote
	}

	originalToken := t.input[start:t.pos]
	var value string
	isSingleQuoted := false
	if rawString {
		// For raw strings, extract content between quotes and handle doubled quotes
		content := originalToken[2 : len(originalToken)-1] // Remove 'r' prefix and quotes
		// Replace doubled quotes with single quotes
		if quote == '"' {
			value = strings.ReplaceAll(content, `""`, `"`)
		} else {
			value = strings.ReplaceAll(content, `''`, `'`)
			isSingleQuoted = true
		}
	} else if quote == '"' {
		quoted := originalToken
		unescaped, err := strconv.Unquote(quoted)
		if err != nil {
			errMsg := errors.GetErrorMessage(errors.ErrInvalidString)
			return Token{}, errors.NewParserError(errors.ErrInvalidString, t.line, startColumn, errMsg+": '"+originalToken+"'")
		}
		value = unescaped
	} else if quote == '\'' {
		isSingleQuoted = true
		content := originalToken[1 : len(originalToken)-1]
		value = t.unescapeString(content)
	}
	return Token{Type: constants.TokenString, Value: value, Token: originalToken, Line: t.line, Column: startColumn, IsSingleQuoted: isSingleQuoted}, nil
}

// Ref: https://regex101.com/r/w6qtHq/1
var pipePattern = regexp.MustCompile(`(?m)^(?P<pipe>[a-zA-Z]+)?:`)

func (t *Tokenizer) readPipeOrBitwiseOr() (Token, error) {
	start := t.pos
	t.advance() // consume first '|'
	if t.current() == '|' {
		t.advance() // consume second '|'
		operator := "||"
		return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: t.column - (t.pos - start)}, nil
	}

	if t.current() == ':' {
		t.advance() // consume ':'
		pipeValue := ":"
		return Token{Type: constants.TokenPipe, Value: pipeValue, Token: pipeValue, Line: t.line, Column: t.column - (t.pos - start)}, nil
	}

	// Fetch the next 10 characters or the rest of the input if less than 10 characters are available
	nextChars := t.input[t.pos:int(
		math.Min(float64(t.pos+10), float64(len(t.input))),
	)]

	pipeMatch := pipePattern.FindStringSubmatch(nextChars)
	if len(pipeMatch) > 1 && pipeMatch[1] != "" {
		pipeName := pipeMatch[1]
		for range pipeName {
			t.advance()
		}
		t.advance() // consume ':'
		return Token{Type: constants.TokenPipe, Value: pipeName, Token: pipeName, Line: t.line, Column: t.column - (t.pos - start)}, nil
	}

	operator := "|"
	return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: t.column - (t.pos - start)}, nil
}

func (t *Tokenizer) readOperator() (Token, error) {
	// This function does not handle operators starting with '|' because that is
	// handled by the readPipeOrBitwiseOr function.

	start := t.pos
	startColumn := t.column

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

	// Handle -- operator, but only in postfix contexts or when followed by non-digit
	// For expressions like "--10", we want to treat as two separate "-" tokens
	if t.current() == '-' && t.peek() == '-' {
		// Look ahead to see what comes after the second '-'
		nextChar := rune(0)
		if t.pos+2 < len(t.input) {
			nextChar, _ = utf8.DecodeRuneInString(t.input[t.pos+2:])
		}

		// If followed by a digit or letter (identifier), treat as two separate minus tokens
		if isDigit(nextChar) || isLetter(nextChar) {
			// Return single '-' token
			t.advance()
			operator := "-"
			return Token{Type: constants.TokenOperator, Value: operator, Token: operator, Line: t.line, Column: startColumn}, nil
		} else {
			// Treat as decrement operator
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

func (t *Tokenizer) singleCharToken(tokenType constants.TokenType) (Token, error) {
	charValue := string(t.current())
	token := Token{Type: tokenType, Value: charValue, Token: charValue, Line: t.line, Column: t.column}
	t.advance()
	return token, nil
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
	if t.pos < len(t.input) && t.input[t.pos] == '\n' {
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

// unescapeString manually handles common escape sequences
// This is a fallback when strconv.Unquote fails
func (t *Tokenizer) unescapeString(s string) string {
	// Try to use strconv.Unquote by converting single-quoted to double-quoted
	// This allows us to use Go's robust unescaping for most cases
	if s != "" {
		if unquoted, err := strconv.Unquote("\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""); err == nil {
			return unquoted
		}
	}

	// Fallback: manual unescaping for Python-style single-quoted strings
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '\'':
				result.WriteByte('\'')
				i += 2
			case '"':
				result.WriteByte('"')
				i += 2
			case '\\':
				result.WriteByte('\\')
				i += 2
			case 'n':
				result.WriteByte('\n')
				i += 2
			case 't':
				result.WriteByte('\t')
				i += 2
			case 'r':
				result.WriteByte('\r')
				i += 2
			case 'b':
				result.WriteByte('\b')
				i += 2
			case 'f':
				result.WriteByte('\f')
				i += 2
			case 'a':
				result.WriteByte('\a')
				i += 2
			case 'v':
				result.WriteByte('\v')
				i += 2
			case '/':
				result.WriteByte('/')
				i += 2
			case 'u':
				// Unicode escape sequence \uXXXX
				if i+5 < len(s) {
					hexStr := s[i+2 : i+6]
					if codePoint, err := strconv.ParseUint(hexStr, 16, 16); err == nil {
						result.WriteRune(rune(codePoint))
						i += 6
						continue
					}
				}
				// If unicode parsing fails, keep the escape sequence as-is
				result.WriteByte(s[i])
				i++
			case 'U':
				// Unicode escape sequence \UXXXXXXXX
				if i+9 < len(s) {
					hexStr := s[i+2 : i+10]
					if codePoint, err := strconv.ParseUint(hexStr, 16, 32); err == nil {
						result.WriteRune(rune(codePoint))
						i += 10
						continue
					}
				}
				result.WriteByte(s[i])
				i++
			default:
				// Unknown escape sequence, keep the backslash and next char
				result.WriteByte(s[i])
				result.WriteByte(s[i+1])
				i += 2
			}
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}
