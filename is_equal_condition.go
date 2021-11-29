package restrict

import (
	"fmt"
	"reflect"
)

// IsEqualConditionType - IsEqualCondition's identifier.
const IsEqualConditionType string = "IS_EQUAL"

// IsEqualCondition - Condition for testing whether given value is equal
// to some other value.
type IsEqualCondition struct {
	Name   string           `json:"name,omitempty" yaml:"name,omitempty"`
	Value  *ValueDescriptor `json:"value,omitempty" yaml:"value,omitempty"`
	Equals *ValueDescriptor `json:"equals,omitempty" yaml:"equals,omitempty"`
}

// Name - returns Condition's name.
func (c *IsEqualCondition) Type() string {
	return IsEqualConditionType
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *IsEqualCondition) Check(request *AccessRequest) error {
	value, err := c.Value.GetValue(request)
	if err != nil {
		return err
	}

	equals, err := c.Equals.GetValue(request)
	if err != nil {
		return err
	}

	equal := reflect.DeepEqual(value, equals)

	if !equal {
		return NewConditionNotSatisfiedError(c, request, fmt.Errorf("Values \"%v\" and \"%v\" are not equal", value, equals))
	}

	return nil
}
