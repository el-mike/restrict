package restrict

import (
	"github.com/el-Mike/restrict/utils"
)

// ValueDescriptor - describes a value that will be tested in its parent Condition.
type ValueDescriptor struct {
	Source ValueSource `json:"source"`
	Field  string      `json:"field"`
	Value  interface{} `json:"value"`
}

// GetValue - returns real value represented by given ValueDescriptor.
func (vd *ValueDescriptor) GetValue(request *AccessRequest) (interface{}, error) {
	if vd.Source == Explicit {
		return vd.Value, nil
	}

	if vd.Field == "" {
		return nil, NewValueDescriptorMalformedError(vd)
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
		return nil, NewValueDescriptorMalformedError(vd)
	}

	if utils.IsMap(source) {
		return utils.GetMapValue(source, vd.Field), nil
	}

	if utils.HasField(source, vd.Field) {
		return utils.GetStructFieldValue(source, vd.Field), nil
	}

	return nil, NewValueDescriptorMalformedError(vd)
}
