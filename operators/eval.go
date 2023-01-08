package operators

func Eval(op string, a, b any) (any, error) {
	return Registry.Get(op)(a, b)
}
