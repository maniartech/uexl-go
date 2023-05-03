package types

type Value interface {

	// Type returns the type of the value
	Type() string

	// IsTruthy returns true if the value is truethy
	IsTruthy() bool

	// IsZero returns true if the value is zero
	IsZero() bool

	// Equals returns true if the value is equal to the other value
	Equals(other Value) bool

	// Compare returns -1 if the value is less than the other value, 0 if the value is equal to the other value, 1 if the value is greater than the other value
	Compare(other Value) (int, error)

	// String returns the string representation of the value
	String() string
}
