package vm

import "github.com/maniartech/uexl/types"

// Re-export Value and related types from types package for backward compatibility
type Value = types.Value

// DateTime and Duration are the boundary forms of temporal values (re-exported from types).
type (
	DateTime = types.DateTime
	Duration = types.Duration
)

// Re-export type constants
const (
	TypeFloat    = types.TypeFloat
	TypeString   = types.TypeString
	TypeBool     = types.TypeBool
	TypeAny      = types.TypeAny
	TypeNull     = types.TypeNull
	TypeDateTime = types.TypeDateTime
	TypeDuration = types.TypeDuration
)

// Re-export constructors
var (
	newFloatValue    = types.NewFloatValue
	newStringValue   = types.NewStringValue
	newBoolValue     = types.NewBoolValue
	newNullValue     = types.NewNullValue
	newAnyValue      = types.NewAnyValue
	newDateTimeValue = types.NewDateTimeValue
	newDurationValue = types.NewDurationValue
)
