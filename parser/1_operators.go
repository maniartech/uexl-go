package parser

// Operator enumeration
type Operator int

const (
	// Arithmetic operators
	OperatorPlus     Operator = iota // +
	OperatorMinus                    // -
	OperatorMultiply                 // *
	OperatorDivide                   // /
	OperatorModulo                   // %

	// Bitwise operators
	OperatorBitwiseAnd // &
	OperatorBitwiseOr  // |
	OperatorBitwiseXor // ^
	OperatorBitwiseNot // ~

	// Shift operators
	OperatorLeftShift  // <<
	OperatorRightShift // >>

	// Assignment operators
	OperatorAssign           // =
	OperatorAddAssign        // +=
	OperatorSubtractAssign   // -=
	OperatorMultiplyAssign   // *=
	OperatorDivideAssign     // /=
	OperatorModuloAssign     // %=
	OperatorBitwiseAndAssign // &=
	OperatorBitwiseOrAssign  // |=
	OperatorBitwiseXorAssign // ^=

	// Increment and decrement operators
	OperatorIncrement // ++
	OperatorDecrement // --

	// Comparison operators
	OperatorEqual          // ==
	OperatorNotEqual       // !=
	OperatorGreaterThan    // >
	OperatorLessThan       // <
	OperatorGreaterOrEqual // >=
	OperatorLessOrEqual    // <=
	// Logical operators
	OperatorLogicalAnd // &&
	OperatorLogicalOr  // ||

	// Conditional operators
	OperatorConditional // ? :

	// Pipe operators
	OperatorPipe      // |:
	OperatorNamedPipe // |map: |filter:
)

// startOperator returns true if the given character is the start of an operator.
func StartOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' ||
		ch == '&' || ch == '|' || ch == '^' || ch == '~' ||
		ch == '=' || ch == '>' || ch == '<' || ch == '!' || ch == '?'
}
