package ntypes

import "time"

// Bool - returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// BoolValue - returns the value of the bool pointer passed in or false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// Int8 - returns a pointer to the int8 value passed in.
func Int8(v int8) *int8 {
	return &v
}

// Int8Value - returns the value of the int8 pointer passed in or 0 if the pointer is nil.
func Int8Value(v *int8) int8 {
	if v != nil {
		return *v
	}
	return 0
}

// Int16 - returns a pointer to the int16 value passed in.
func Int16(v int16) *int16 {
	return &v
}

// Int16Value - returns the value of the int16 pointer passed in or 0 if the pointer is nil.
func Int16Value(v *int16) int16 {
	if v != nil {
		return *v
	}
	return 0
}

// Int32 - returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

// Int32Value - returns the value of the int32 pointer passed in or 0 if the pointer is nil.
func Int32Value(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

// Rune - returns a pointer to the rune value passed in.
func Rune(v rune) *rune {
	return &v
}

// RuneValue - returns the value of the rune pointer passed in or 0 if the pointer is nil.
func RuneValue(v *rune) rune {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 - returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

// Int64Value - returns the value of the int64 pointer passed in or 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Int - returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// Uint8 - returns a pointer to the uint8 value passed in.
func Uint8(v uint8) *uint8 {
	return &v
}

// Uint8Value - returns the value of the uint8 pointer passed in or 0 if the pointer is nil.
func Uint8Value(v *uint8) uint8 {
	if v != nil {
		return *v
	}
	return 0
}

// Byte - returns a pointer to the byte value passed in.
func Byte(v byte) *byte {
	return &v
}

// ByteValue - returns the value of the byte pointer passed in or 0 if the pointer is nil.
func ByteValue(v *byte) byte {
	if v != nil {
		return *v
	}
	return 0
}

// Uint16 - returns a pointer to the uint16 value passed in.
func Uint16(v uint16) *uint16 {
	return &v
}

// Uint16Value - returns the value of the uint16 pointer passed in or 0 if the pointer is nil.
func Uint16Value(v *uint16) uint16 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint32 - returns a pointer to the uint32 value passed in.
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint32Value - returns the value of the uint32 pointer passed in or 0 if the pointer is nil.
func Uint32Value(v *uint32) uint32 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint64 - returns a pointer to the uint64 value passed in.
func Uint64(v uint64) *uint64 {
	return &v
}

// Uint64Value - returns the value of the uint64 pointer passed in or 0 if the pointer is nil.
func Uint64Value(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint - returns a pointer to the uint value passed in.
func Uint(v uint) *uint {
	return &v
}

// UintValue - returns the value of the uint pointer passed in or 0 if the pointer is nil.
func UintValue(v *uint) uint {
	if v != nil {
		return *v
	}
	return 0
}

// Float32 - returns a pointer to the float32 value passed in.
func Float32(v float32) *float32 {
	return &v
}

// Float32Value - returns the value of the float32 pointer passed in or 0 if the pointer is nil.
func Float32Value(v *float32) float32 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64 - returns a pointer to the float64 value passed in.
func Float64(v float64) *float64 {
	return &v
}

// Float64Value - returns the value of the float64 pointer passed in or 0 if the pointer is nil.
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// String - returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue - returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Time - returns a pointer to the time.Time value passed in.
func Time(v time.Time) *time.Time {
	return &v
}

// TimeValue - returns the value of the time.Time pointer passed in or time.Time{} if the pointer is nil.
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}
