package vm

import "github.com/maniartech/uexl_go/types"

// Re-export Value and related types from types package for backward compatibility
type Value = types.Value

// Re-export type constants
const (
	TypeFloat  = types.TypeFloat
	TypeString = types.TypeString
	TypeBool   = types.TypeBool
	TypeAny    = types.TypeAny
	TypeNull   = types.TypeNull
)

// Re-export constructors
var (
	newFloatValue  = types.NewFloatValue
	newStringValue = types.NewStringValue
	newBoolValue   = types.NewBoolValue
	newNullValue   = types.NewNullValue
	newAnyValue    = types.NewAnyValue
)
