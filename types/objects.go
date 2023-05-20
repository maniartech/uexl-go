package types

import (
	"encoding/json"
	"errors"
)

// Object represents an object in the uexl language.
type Object map[string]Value

// Type returns the type of the object.
func (o Object) Type() string {
	return "object"
}

// IsTruthy returns true if the object is not nil or empty.
func (o Object) IsTruthy() bool {
	return len(o) > 0
}

// IsZero returns true if the object is nil or empty.
func (o Object) IsZero() bool {
	return len(o) == 0
}

// Equals returns true if the object values are equal. It performs a deep comparison.
func (o Object) Equals(other Value) bool {
	if otherObject, ok := other.(Object); ok {
		if len(o) != len(otherObject) {
			return false
		}

		for k, v := range o {
			if otherV, ok := otherObject[k]; ok {
				if !otherV.Equals(v) {
					return false
				}
			} else {
				return false
			}
		}

		// All keys and values are equal
		return true
	}

	// Not an Object
	return false
}

// Compare returns 0 if the object values are equal, -1 if the object is nil or empty, and 1 if the object is not nil or empty.
func (o Object) Compare(other Value) (int, error) {
	return 0, errors.New("comparision not supported for objects")
}

// String returns the string representation of the object.
func (o Object) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// Add merges the other object into this object. If a key exists in both objects, the value from the other object is used.
func (o Object) Add(other Value) (Value, error) {
	if otherV, ok := other.(Object); ok {
		for k, v := range otherV {
			o[k] = v
		}
		return o, nil
	} else if otherV, ok := other.(Array); ok {
		// prepend the current object to the array
		otherV = append(Array{o}, otherV...)
		return otherV, nil
	}

	return nil, errors.New("invalid type for addition")
}

// Minus removes the keys from this object that are in the other object.
func (o Object) Minus(other Object) Object {
	for k := range other {
		delete(o, k)
	}
	return o
}

// Dot returns the value of the key in the object. If the key does not exist, nil is returned.
func (o Object) Dot(key string) (Value, error) {
	if v, ok := o[key]; ok {
		return v, nil
	}
	return nil, errors.New("key not found: " + key)
}
