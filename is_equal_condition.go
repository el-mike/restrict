package restrict

import (
	"reflect"
)

// IsEqualConditionName - IsEqualCondition's identifier.
const IsEqualConditionName string = "IS_EQUAL"

// IsEqualCondition - Condition for testing whether given value is equal
// to some other value.
type IsEqualCondition struct {
	Value  *ValueDescriptor `json:"value,omitempty" yaml:"value,omitempty"`
	Equals *ValueDescriptor `json:"equals,omitempty" yaml:"equals,omitempty"`
}

// Name - returns Condition's name.
func (c *IsEqualCondition) Name() string {
	return IsEqualConditionName
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *IsEqualCondition) Check(request *AccessRequest) bool {
	value := c.Value.GetValue(request)
	equals := c.Equals.GetValue(request)

	return reflect.DeepEqual(value, equals)
}
