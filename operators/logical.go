package operators

import "github.com/maniartech/uexl_go/core"

func LogicalAnd(a, b any) (any, error) {
	return core.IsTruthy(a) && core.IsTruthy(b), nil
}

func LogicalOr(a, b any) (any, error) {

	if core.IsTruthy(a) {
		return a, nil
	}

	if core.IsTruthy(b) {
		return b, nil
	}

	return false, nil
}

func LogicalXor(a, b any) (any, error) {
	return core.IsTruthy(a) != core.IsTruthy(b), nil
}

func LogicalNand(a, b any) (any, error) {
	return !(core.IsTruthy(a) && core.IsTruthy(b)), nil
}

// Urinary operators
func LogicalNot(a any) (any, error) {
	return !core.IsTruthy(a), nil
}

func init() {
	Registry.Register("&&", LogicalAnd)
	Registry.Register("||", LogicalOr)
}
