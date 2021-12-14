package restrict

import (
	"fmt"

	"github.com/el-mike/restrict/internal/utils"
)

// ValueDescriptor - describes a value that will be tested in its parent Condition.
type ValueDescriptor struct {
	// Source - source of the value, one of the predefined enum type (ValueSource).
	Source ValueSource `json:"source,omitempty" yaml:"source,omitempty"`
	// Field - field on the given ValueSource that should hold the value.
	Field string `json:"field,omitempty" yaml:"field,omitempty"`
	// Value - explicit value taken when using ValueSource.Explicit as value source.
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

// GetValue - returns real value represented by given ValueDescriptor.
func (vd *ValueDescriptor) GetValue(request *AccessRequest) (interface{}, error) {
	if vd == nil {
		return nil, newValueDescriptorMalformedError(vd, fmt.Errorf("ValueDescriptor cannot be nil"))
	}

	if vd.Source == Explicit {
		return vd.Value, nil
	}

	if vd.Field == "" {
		return nil, newValueDescriptorMalformedError(vd, fmt.Errorf("Field cannot be empty for Source: \"%s\"", vd.Source.String()))
	}

	var source interface{}

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
		return nil, newValueDescriptorMalformedError(vd, fmt.Errorf("Source could not be find"))
	}

	if utils.IsMap(source) {
		return utils.GetMapValue(source, vd.Field), nil
	}

	if utils.HasField(source, vd.Field) {
		return utils.GetStructFieldValue(source, vd.Field), nil
	}

	return nil, newValueDescriptorMalformedError(vd, fmt.Errorf("Field \"%s\" does not exist on Source: \"%s\"", vd.Field, vd.Source.String()))
}
