package utils

import "reflect"

// IsSameType - returns true if both arguments are the same type,
// false otherwise.
func IsSameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// IsStruct - returns true if argument is a struct, false otherwise.
func IsStruct(a interface{}) bool {
	return reflect.ValueOf(a).Kind() == reflect.Struct
}
