package restrict

import (
	"reflect"

	"github.com/el-Mike/restrict/utils"
)

// IsEqualConditionName - IsEqualCondition's identifier.
const IsEqualConditionName string = "IS_EQUAL"

// IsEqualCondition - Condition for testing whether given value is equal
// to some other value.
type IsEqualCondition struct {
	Value      interface{} `json:"equals,omitempty" yaml:"equals,omitempty"`
	ContextKey string      `json:"contextKey,omitempty" yaml:"contextKey,omitempty"`
}

// Name - returns Condition's name.
func (c *IsEqualCondition) Name() string {
	return IsEqualConditionName
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *IsEqualCondition) Check(value interface{}, request *AccessRequest) bool {
	if c.ContextKey != "" {
		return reflect.DeepEqual(value, request.Context[c.ContextKey])
	}

	if !utils.IsSameType(value, c.Value) {
		return false
	}

	return reflect.DeepEqual(value, c.Value)

}
