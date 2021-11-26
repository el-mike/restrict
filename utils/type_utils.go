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

// IsMap - returns true if argument is a Map, false otherwise.
func IsMap(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
}

// GetMapValue - returns a value under given key in passed map.
func GetMapValue(mapValue interface{}, keyValue interface{}) interface{} {
	rMapValue := reflect.ValueOf(mapValue)

	if rMapValue.Kind() == reflect.Ptr {
		rMapValue = rMapValue.Elem()
	}

	return rMapValue.MapIndex(reflect.ValueOf(keyValue)).Interface()
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

// GetStructFieldValue - returns a value under given field in passed struct.
func GetStructFieldValue(structValue interface{}, fieldName string) interface{} {
	rStructValue := reflect.ValueOf(structValue)

	if rStructValue.Kind() == reflect.Ptr {
		rStructValue = rStructValue.Elem()
	}

	return rStructValue.FieldByName(fieldName).Interface()
}
