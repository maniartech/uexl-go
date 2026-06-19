// Package strings is the UExL string family beyond the core built-ins (len/substr/contains/str/join and
// the unicode views): case, trim, pad, replace, split, prefix/suffix tests, indexOf, repeat. Attach via
// uexl.WithStrings().
package strings

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/maniartech/uexl/builtins/fn"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"upper":      str1("upper", strings.ToUpper),
	"lower":      str1("lower", strings.ToLower),
	"trim":       trimFunc("trim", strings.TrimSpace, strings.Trim),
	"trimStart":  trimFunc("trimStart", func(s string) string { return strings.TrimLeft(s, " \t\n\r") }, strings.TrimLeft),
	"trimEnd":    trimFunc("trimEnd", func(s string) string { return strings.TrimRight(s, " \t\n\r") }, strings.TrimRight),
	"replace":    builtinReplace,
	"split":      builtinSplit,
	"startsWith": str2Bool("startsWith", strings.HasPrefix),
	"endsWith":   str2Bool("endsWith", strings.HasSuffix),
	"indexOf":    builtinIndexOf,
	"repeat":     builtinRepeat,
	"padStart":   padFunc("padStart", true),
	"padEnd":     padFunc("padEnd", false),
}

func str1(name string, f func(string) string) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 1, 1); err != nil {
			return nil, err
		}
		s, err := fn.Str(name, args, 0)
		if err != nil {
			return nil, err
		}
		return f(s), nil
	}
}

func str2Bool(name string, f func(s, t string) bool) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 2, 2); err != nil {
			return nil, err
		}
		s, e1 := fn.Str(name, args, 0)
		t, e2 := fn.Str(name, args, 1)
		if e1 != nil {
			return nil, e1
		}
		if e2 != nil {
			return nil, e2
		}
		return f(s, t), nil
	}
}

// trimFunc trims whitespace with one arg, or a custom cutset with two.
func trimFunc(name string, ws func(string) string, cut func(string, string) string) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 1, 2); err != nil {
			return nil, err
		}
		s, err := fn.Str(name, args, 0)
		if err != nil {
			return nil, err
		}
		if len(args) == 2 {
			set, err := fn.Str(name, args, 1)
			if err != nil {
				return nil, err
			}
			return cut(s, set), nil
		}
		return ws(s), nil
	}
}

func builtinReplace(args ...any) (any, error) {
	if err := fn.Arity("replace", args, 3, 3); err != nil {
		return nil, err
	}
	s, e1 := fn.Str("replace", args, 0)
	old, e2 := fn.Str("replace", args, 1)
	nw, e3 := fn.Str("replace", args, 2)
	for _, e := range []error{e1, e2, e3} {
		if e != nil {
			return nil, e
		}
	}
	return strings.ReplaceAll(s, old, nw), nil
}

func builtinSplit(args ...any) (any, error) {
	if err := fn.Arity("split", args, 2, 2); err != nil {
		return nil, err
	}
	s, e1 := fn.Str("split", args, 0)
	sep, e2 := fn.Str("split", args, 1)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}
	parts := strings.Split(s, sep)
	out := make([]any, len(parts))
	for i, p := range parts {
		out[i] = p
	}
	return out, nil
}

func builtinIndexOf(args ...any) (any, error) {
	if err := fn.Arity("indexOf", args, 2, 2); err != nil {
		return nil, err
	}
	s, e1 := fn.Str("indexOf", args, 0)
	sub, e2 := fn.Str("indexOf", args, 1)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}
	return float64(strings.Index(s, sub)), nil // byte index; -1 if absent
}

func builtinRepeat(args ...any) (any, error) {
	if err := fn.Arity("repeat", args, 2, 2); err != nil {
		return nil, err
	}
	s, e1 := fn.Str("repeat", args, 0)
	n, e2 := fn.Int("repeat", args, 1)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}
	if n < 0 {
		return nil, fmt.Errorf("repeat: count must be non-negative, got %d", n)
	}
	return strings.Repeat(s, n), nil
}

func padFunc(name string, atStart bool) fn.Func {
	return func(args ...any) (any, error) {
		if err := fn.Arity(name, args, 3, 3); err != nil {
			return nil, err
		}
		s, e1 := fn.Str(name, args, 0)
		length, e2 := fn.Int(name, args, 1)
		padStr, e3 := fn.Str(name, args, 2)
		for _, e := range []error{e1, e2, e3} {
			if e != nil {
				return nil, e
			}
		}
		cur := utf8.RuneCountInString(s)
		if cur >= length || padStr == "" {
			return s, nil
		}
		padRunes := []rune(padStr)
		var b strings.Builder
		for i := 0; i < length-cur; i++ {
			b.WriteRune(padRunes[i%len(padRunes)])
		}
		if atStart {
			return b.String() + s, nil
		}
		return s + b.String(), nil
	}
}
