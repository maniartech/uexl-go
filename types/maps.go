package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Map map[string]any

type Context map[string]Value

func (ctx Context) ShallowCopy() Context {
	newCtx := make(Context)
	for k, v := range ctx {
		newCtx[k] = v
	}
	return newCtx
}

// ValueAtPath is a helper function to get a value from a context using a path.
// For example, if you have a map like this:
//
//	map[string]interface{}{
//		"foo": map[string]interface{}{
//			"bar": "baz",
//		},
//	}
//
// You can get the value of "baz" using the path "foo.bar".
// If the path is not found, nil is returned.
// If the member is a slice, you can get the value of a specific index using the path "foo.bar.0".
// If the member is a slice of maps, you can get the value of a specific index using the path "foo.bar.0.baz".
// For example, if you have a map like this:
//
//	map[string]interface{}{
//		"foo": map[string]interface{}{
//			"bar": []interface{}{
//				map[string]interface{}{
//					"baz": "qux",
//				},
//			},
//		},
//	}
//
// You can get the value of "qux" using the path "foo.bar.0.baz".
func (m Context) ValueAtPath(path string) (Value, error) {
	keys := strings.Split(path, ".")

	if len(keys) == 0 {
		return nil, fmt.Errorf("cannot get value from %v", m)
	} else if len(keys) == 1 {
		return m[keys[0]], nil
	}

	ret, err := getValue(m, keys)
	if err != nil {
		return nil, err
	}

	val, ok := ret.(Value)
	if !ok {
		return nil, fmt.Errorf("cannot get value from %v", m)
	}

	return val, nil
}

func getValue(m interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return m, nil
	}

	switch v := m.(type) {
	// `Context` is a type alias for a map of strings to `Value` types. It is used to represent a
	// context in which a template is executed, where the keys are variable names and the values are
	// their corresponding values. The `ValueAtPath` method is a helper function to get a value from
	// the context using a path.
	case Context:
		return getValue(v[keys[0]], keys[1:])
	case Object:
		return getValue(v[keys[0]], keys[1:])
	case Map:
		return getValue(v[keys[0]], keys[1:])
	case map[string]interface{}:
		return getValue(v[keys[0]], keys[1:])

	case []interface{}:
		index, _ := strconv.Atoi(keys[0])
		return getValue(v[index], keys[1:])
	case Array:
		index, _ := strconv.Atoi(keys[0])
		return getValue(v[index], keys[1:])
	default:
		return nil, fmt.Errorf("cannot get value from %v", m)
	}
}
