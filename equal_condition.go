package restrict

import (
	"fmt"
	"reflect"
)

const (
	// EqualConditionType - EqualCondition's type identifier.
	EqualConditionType = "EQUAL"
	//NotEqualConditionType - NotEqualCondition's type identifier.
	NotEqualConditionType = "NOT_EQUAL"
)

// baseEqualCondition - describes fields needed by Equal/NotEqual Conditions.
type baseEqualCondition struct {
	// ID - Condition's id, useful when there is a need to identify failing Condition.
	ID string `json:"name,omitempty" yaml:"name,omitempty"`
	// Left - ValueDescriptor for left operand of equality check.
	Left *ValueDescriptor `json:"left" yaml:"left"`
	// Right - ValueDescriptor for right operand of equality check.
	Right *ValueDescriptor `json:"right" yaml:"right"`
}

// EqualCondition - checks whether given value (Left) is equal to some other value (Right).
type EqualCondition baseEqualCondition

// Type - returns Condition's type.
func (c *EqualCondition) Type() string {
	return EqualConditionType
}

// Check - returns true if values are equal, false otherwise.
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

// EqualCondition - checks whether given value (Left) is not equal to some other value (Right).
type NotEqualCondition baseEqualCondition

// Type - returns Condition's type.
func (c *NotEqualCondition) Type() string {
	return NotEqualConditionType
}

// Check - returns true if values are not equal, false otherwise.
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
