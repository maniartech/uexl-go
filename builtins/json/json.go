// Package json is the UExL JSON family: parseJson (string -> value) and toJson (value -> string).
// Datetime/duration values serialize to their ISO 8601 forms. Attach via uexl.WithJSON().
package json

import (
	encjson "encoding/json"
	"fmt"
	"time"

	"github.com/maniartech/uexl/builtins/fn"
	"github.com/maniartech/uexl/types"
)

// Builtins maps function names to their implementations.
var Builtins = map[string]fn.Func{
	"parseJson": builtinParseJson,
	"toJson":    builtinToJson,
}

func builtinParseJson(args ...any) (any, error) {
	if err := fn.Arity("parseJson", args, 1, 1); err != nil {
		return nil, err
	}
	s, err := fn.Str("parseJson", args, 0)
	if err != nil {
		return nil, err
	}
	var out any
	if err := encjson.Unmarshal([]byte(s), &out); err != nil {
		return nil, fmt.Errorf("parseJson: %w", err)
	}
	return out, nil // encoding/json yields float64, string, bool, nil, []any, map[string]any
}

func builtinToJson(args ...any) (any, error) {
	if err := fn.Arity("toJson", args, 1, 2); err != nil {
		return nil, err
	}
	pretty := false
	if len(args) == 2 {
		b, err := fn.Bool("toJson", args, 1)
		if err != nil {
			return nil, err
		}
		pretty = b
	}
	ready := jsonReady(args[0])
	var (
		bytes []byte
		err   error
	)
	if pretty {
		bytes, err = encjson.MarshalIndent(ready, "", "  ")
	} else {
		bytes, err = encjson.Marshal(ready)
	}
	if err != nil {
		return nil, fmt.Errorf("toJson: %w", err)
	}
	return string(bytes), nil
}

// jsonReady deep-converts UExL temporal values to their ISO 8601 string forms so they serialize
// portably (a datetime becomes an RFC 3339 instant, a duration an ISO 8601 duration).
func jsonReady(v any) any {
	switch x := v.(type) {
	case types.DateTime:
		return time.UnixMilli(x.Millis).UTC().Format(time.RFC3339)
	case types.Duration:
		return types.FormatISODuration(x.Millis)
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = jsonReady(e)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			out[k] = jsonReady(e)
		}
		return out
	default:
		return v
	}
}
