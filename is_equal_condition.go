package restrict

import (
	"fmt"
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
