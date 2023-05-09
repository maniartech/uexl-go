package types

// Adder interface represents types that can perform addition.
type Adder interface {
	// Add performs addition with another Value and returns the result and an error.
	Add(other Value) (Value, error)
}

// Subtractor interface represents types that can perform subtraction.
type Subtractor interface {
	// Subtract performs subtraction with another Value and returns the result and an error.
	Subtract(other Value) (Value, error)
}

// Multiplier interface represents types that can perform multiplication.
type Multiplier interface {
	// Multiply performs multiplication with another Value and returns the result and an error.
	Multiply(other Value) (Value, error)
}

// Divider interface represents types that can perform division.
type Divider interface {
	// Divide performs division with another Value and returns the result and an error.
	Divide(other Value) (Value, error)
}

// Modulus interface represents types that can perform modulo operation.
type Modulus interface {
	// Mod performs modulo operation with another Value and returns the result and an error.
	Mod(other Value) (Value, error)
}

// Comparer interface represents types that can perform comparison.
type Comparer interface {
	// Compare performs comparison with another Value and returns an integer and an error.
	// The integer result is:
	// - negative if the current instance is less than the other,
	// - zero if they are equal,
	// - positive if the current instance is greater than the other.
	Compare(other Value) (int, error)
}

// And interface represents types that can perform logical AND operation.
type And interface {
	// And performs logical AND operation with another Value and returns the result and an error.
	And(other Value) (Value, error)
}

// Or interface represents types that can perform logical OR operation.
type Or interface {
	// Or performs logical OR operation with another Value and returns the result and an error.
	Or(other Value) (Value, error)
}

// Not interface represents types that can perform logical NOT operation.
type Not interface {
	// Not performs logical NOT operation and returns the result and an error.
	Not() (Value, error)
}

// Dot interface represents types that can perform dot operation.
type Dot interface {
	// Dot performs dot operation with another Value and returns the result and an error.
	Dot(other Value) (Value, error)
}
