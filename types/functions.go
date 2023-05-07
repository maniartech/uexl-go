package types

// Function represents a function in the uexl language.
type Function func(...interface{}) (interface{}, error)
