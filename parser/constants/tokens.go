package constants

// TokenType represents the type of a token in the UExl language
type TokenType int

// Token type constants - these represent the different types of tokens
// that can be recognized by the tokenizer
const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdentifier
	TokenOperator
	// Optional chaining composite tokens
	TokenQuestionDot         // '?.'
	TokenQuestionLeftBracket // '?['
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
	TokenDateTime // d"..." datetime literal; canonical epoch ms in TokenValue.Int
	TokenDuration // duration literal (e.g. 5m, 1.5h, 30ms); canonical ms in TokenValue.Int
	TokenError    // Special token type for tokenizer errors
)

// String returns the string representation of a token type
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
	case TokenQuestionDot:
		return "QuestionDot"
	case TokenQuestionLeftBracket:
		return "QuestionLeftBracket"
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
	case TokenDateTime:
		return "DateTime"
	case TokenDuration:
		return "Duration"
	case TokenError:
		return "Error"
	default:
		return "Unknown"
	}
}
