package restrict

import (
	"fmt"
	"reflect"
)

const (
	// EmptyConditionType - EmptyCondition's type identifier.
	EmptyConditionType = "EMPTY"
	// NotEmptyConditionType - NotEmptyCondition's type identifier.
	NotEmptyConditionType = "NOT_EMPTY"
)

// baseEmptyCondition - describes fields needed by Empty/NotEmpty Conditions.
type baseEmptyCondition struct {
	// ID - Condition's id, useful when there is a need to identify failing Condition.
	ID string `json:"name,omitempty" yaml:"name,omitempty"`
	// Value - ValueDescriptor for the value being checked.
	Value *ValueDescriptor `json:"value" yaml:"value"`
}

// EmptyCondition - Condition for testing whether given value is empty.
type EmptyCondition baseEmptyCondition

// Type - returns Condition's type.
func (c *EmptyCondition) Type() string {
	return EmptyConditionType
}

// Check - returns true if value is empty (zero-like), false otherwise.
func (c *EmptyCondition) Check(request *AccessRequest) error {
	value, err := c.Value.GetValue(request)
	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	empty := reflect.ValueOf(value).IsZero()

	if !empty {
		return NewConditionNotSatisfiedError(c, request, fmt.Errorf("value \"%v\" is not empty", value))
	}

	return nil
}

type NotEmptyCondition baseEmptyCondition

// Type - returns Condition's type.
func (c *NotEmptyCondition) Type() string {
	return NotEmptyConditionType
}

// Check - returns true if value is not empty (zero-like), false otherwise.
func (c *NotEmptyCondition) Check(request *AccessRequest) error {
	value, err := c.Value.GetValue(request)
	if err != nil {
		return err
	}

	if value == nil {
		return NewConditionNotSatisfiedError(c, request, fmt.Errorf("value \"%v\" is empty", value))
	}

	empty := reflect.ValueOf(value).IsZero()

	if empty {
		return NewConditionNotSatisfiedError(c, request, fmt.Errorf("value \"%v\" is empty", value))
	}

	return nil
}
