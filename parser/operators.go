package parser

import "github.com/maniartech/uexl_go/parser/constants"

// startOperator returns true if the given character is the start of an operator.
func StartOperator(ch rune) bool {
	return constants.StartOperator(ch)
}
