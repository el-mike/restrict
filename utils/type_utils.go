package utils

import "reflect"

// IsSameType - returns true if both arguments are the same type,
// false otherwise.
func IsSameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
