package uexl

import "fmt"

// AsFloat64 converts v to float64.
// Accepts float64 directly, or widens int, int64, and float32.
// Returns an error if v cannot be represented as float64, including nil.
func AsFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("uexl: AsFloat64: cannot convert %T to float64", v)
	}
}

// AsBool converts v to bool.
// No truthy coercion — only bool values are accepted.
// Returns an error for any other type, including nil.
func AsBool(v any) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("uexl: AsBool: cannot convert %T to bool", v)
	}
	return b, nil
}

// AsString converts v to string.
// No fmt.Sprint fallback — only string values are accepted.
// Returns an error for any other type, including nil.
func AsString(v any) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("uexl: AsString: cannot convert %T to string", v)
	}
	return s, nil
}

// AsSlice converts v to []any.
// No element conversion — only []any values are accepted.
// Returns an error for any other type, including nil.
func AsSlice(v any) ([]any, error) {
	s, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("uexl: AsSlice: cannot convert %T to []any", v)
	}
	return s, nil
}

// AsMap converts v to map[string]any.
// No key/value conversion — only map[string]any values are accepted.
// Returns an error for any other type, including nil.
func AsMap(v any) (map[string]any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("uexl: AsMap: cannot convert %T to map[string]any", v)
	}
	return m, nil
}
