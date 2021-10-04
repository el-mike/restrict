package restrict

import (
	"reflect"

	"github.com/el-Mike/restrict/utils"
)

const IS_EQUAL_CONDITION_NAME string = "IS_EQUAL"

type IsEqualCondition struct {
	Value      interface{} `json:"equals,omitempty" yaml:"equals,omitempty"`
	ContextKey string      `json:"contextKey,omitempty" yaml:"contextKey,omitempty"`
}

// Name - returns Condition's name.
func (c *IsEqualCondition) Name() string {
	return IS_EQUAL_CONDITION_NAME
}

// Check - returns true if passed value is the same as Value set for Condition,
// false otherwise.
func (c *IsEqualCondition) Check(value interface{}, request *AccessRequest) bool {
	if c.ContextKey != "" {
		return reflect.DeepEqual(value, request.Context[c.ContextKey])
	}

	if !utils.IsSameType(value, c.Value) {
		return false
	}

	return reflect.DeepEqual(value, c.Value)

}
