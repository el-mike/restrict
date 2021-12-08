package restrict

import (
	"fmt"
	"reflect"
)

// EqualConditionType - EqualCondition's identifier.
const EqualConditionType string = "EQUAL"
const NotEqualConditionType string = "NOT_EQUAL"

// EqualCondition - Condition for testing whether given value is equal
// to some other value.
type EqualCondition struct {
	// Name - Condition's name, useful when there is a need to identify failing Condition.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Left - ValueDescriptor for left operand of equality check.
	Left *ValueDescriptor `json:"value,omitempty" yaml:"value,omitempty"`
	// Right - ValueDescriptor for right operand of equality check.
	Right *ValueDescriptor `json:"equals,omitempty" yaml:"equals,omitempty"`
}

// Name - returns Condition's name.
func (c *EqualCondition) Type() string {
	return EqualConditionType
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *EqualCondition) Check(request *AccessRequest) error {
	left, right, err := unpackDescriptors(c.Left, c.Right, request)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(left, right) {
		return newConditionNotSatisfiedError(c, request, fmt.Errorf("Values \"%v\" and \"%v\" are not equal", left, right))
	}

	return nil
}

// EqualCondition - Condition for testing whether given value is not equal
// to some other value.
type NotEqualCondition struct {
	// Name - Condition's name, useful when there is a need to identify failing Condition.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Left - ValueDescriptor for left operand of equality check.
	Left *ValueDescriptor `json:"value,omitempty" yaml:"value,omitempty"`
	// Right - ValueDescriptor for right operand of equality check.
	Right *ValueDescriptor `json:"equals,omitempty" yaml:"equals,omitempty"`
}

// Name - returns Condition's name.
func (c *NotEqualCondition) Type() string {
	return NotEqualConditionType
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *NotEqualCondition) Check(request *AccessRequest) error {
	left, right, err := unpackDescriptors(c.Left, c.Right, request)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(left, right) {
		return newConditionNotSatisfiedError(c, request, fmt.Errorf("Values \"%v\" and \"%v\" are equal", left, right))
	}

	return nil
}

// unpackDescriptors - helper function for unpacking ValueDescriptors' values.
func unpackDescriptors(left, right *ValueDescriptor, request *AccessRequest) (interface{}, interface{}, error) {
	leftValue, err := left.GetValue(request)
	if err != nil {
		return nil, nil, err
	}

	rightValue, err := right.GetValue(request)
	if err != nil {
		return nil, nil, err
	}

	return leftValue, rightValue, nil
}
