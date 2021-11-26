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
func (vd *ValueDescriptor) GetValue(request *AccessRequest) interface{} {
	if vd.Source == Explicit {
		return vd.Value
	}

	if vd.Field == "" {
		return nil
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
		return nil
	}

	if utils.IsMap(source) {
		return utils.GetMapValue(source, vd.Field)
	}

	if utils.HasField(source, vd.Field) {
		return utils.GetStructFieldValue(source, vd.Field)
	}

	return nil
}
