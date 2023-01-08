package functions

type Function func([]any) (any, error)

type registry map[string]Function

var Registry = registry{}

func (r registry) Register(name string, f Function) {
	r[name] = f
}

func (r registry) Unregister(name string) {
	delete(r, name)
}

func (r registry) Get(name string) Function {
	if f, ok := r[name]; ok {
		return f
	}
	return nil
}
