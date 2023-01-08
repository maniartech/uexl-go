package operators

type OperatorFn func(a, b any) (any, error)

type registry map[string]OperatorFn

var Registry registry = registry{}

func (r registry) Register(name string, f OperatorFn) {
	r[name] = f
}

func (r registry) Unregister(name string) {
	delete(r, name)
}

func (r registry) Get(name string) OperatorFn {
	if f, ok := r[name]; ok {
		return f
	}
	return nil
}
