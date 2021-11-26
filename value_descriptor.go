package restrict

import (
	"reflect"

	"github.com/el-Mike/restrict/utils"
)

// ValueSource - enum type for source of value for given ValueDescriptor.
type ValueSource int

const (
	// SubjectField - value that comes from Subject's field.
	SubjectField ValueSource = iota
	// ResourceField - value taht comes from Resource's field.
	ResourceField
	// ContextField - value that comes from Context's field.
	ContextField
	// Explicit - value set explicitly in PolicyDefinition.
	Explicit
)

// ValueDescriptor - describes a value that will be tested in its parent Condition.
type ValueDescriptor struct {
	Source ValueSource `json:"source"`
	Field  string      `json:"field"`
	Value  interface{} `json:"value"`
}

// GetValue - returns real value represented by given ValueDescriptor.
func (vd *ValueDescriptor) GetValue(request *AccessRequest) interface{} {
	if vd.Source == Explicit {
		return vd.Value
	}

	var source interface{} = nil

	if vd.Source == SubjectField {
		source = request.Subject
	}

	if vd.Source == ResourceField {
		source = request.Resource
	}

	if vd.Source == ContextField {
		source = request.Context
	}

	if source == nil {
		return nil
	}

	if reflect.ValueOf(source).Kind() == reflect.Map {
		return reflect.ValueOf(source).MapIndex(reflect.ValueOf(vd.Field)).Interface()
	}

	if utils.HasField(source, vd.Field) {
		return reflect.ValueOf(source).Elem().FieldByName(vd.Field).Interface()
	}

	return nil
}
