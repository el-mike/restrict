package utils

import "reflect"

// IsSameType - returns true if both arguments are the same type, false otherwise.
func IsSameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

// IsStruct - returns true if argument is a struct, false otherwise.
func IsStruct(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Struct
}

// HasField - returns true if given field exists on passed struct, false otherwise.
func HasField(value interface{}, fieldName string) bool {
	rValue := reflect.ValueOf(value)

	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	if rValue.Kind() != reflect.Struct {
		return false
	}

	return rValue.FieldByName(fieldName).IsValid()
}
