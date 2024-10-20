package restrict

import (
	"fmt"
	"github.com/el-mike/restrict/internal/utils"
	"strings"
)

// PermissionErrors - an alias type for a slice of PermissionError, with extra helper methods.
type PermissionErrors []*PermissionError

// First - returns the first PermissionError encountered when performing authorization.
func (ae PermissionErrors) First() *PermissionError {
	if ae == nil || len(ae) == 0 {
		return nil
	}

	return ae[0]
}

// GetByRoleName - returns PermissionError structs specific to given Role.
func (ae PermissionErrors) GetByRoleName(roleName string) PermissionErrors {
	if ae == nil {
		return nil
	}

	result := PermissionErrors{}

	for _, e := range ae {
		if e.RoleName == roleName {
			result = append(result, e)
		}
	}

	return result
}

// GetByAction - returns PermissionError structs specific to given Action.
func (ae PermissionErrors) GetByAction(action string) PermissionErrors {
	if ae == nil {
		return nil
	}

	result := PermissionErrors{}

	for _, e := range ae {
		if e.Action == action {
			result = append(result, e)
		}
	}

	return nil
}

// GetFailedActions - returns all Actions for which access was denied.
func (ae PermissionErrors) GetFailedActions() []string {
	actions := []string{}

	for _, e := range ae {
		if utils.StringSliceContains(actions, e.Action) {
			actions = append(actions, e.Action)
		}
	}

	return actions
}

// AccessDeniedError - thrown when AccessRequest could not be satisfied due to
// insufficient privileges.
type AccessDeniedError struct {
	Request *AccessRequest
	Errors  PermissionErrors
}

// newAccessDeniedError - returns new AccessDeniedError instance.
func newAccessDeniedError(request *AccessRequest, errors PermissionErrors) *AccessDeniedError {
	return &AccessDeniedError{
		Request: request,
		Errors:  errors,
	}
}

// Error - error interface implementation.
func (e *AccessDeniedError) Error() string {
	failedActions := e.Errors.GetFailedActions()

	return fmt.Sprintf("Access denied for Actions: %s on Resource: %s", strings.Join(failedActions, ", "), e.Request.Resource.GetResourceName())
}

// PermissionError - thrown when Permission is not granted for a given Action.
type PermissionError struct {
	Action         string
	RoleName       string
	ResourceName   string
	ConditionError error
}

// newPermissionError - returns new PermissionError instance.
func newPermissionError(action, roleName, resourceName string, conditionError error) *PermissionError {
	return &PermissionError{
		Action:         action,
		RoleName:       roleName,
		ResourceName:   resourceName,
		ConditionError: conditionError,
	}
}

// Error - error interface implementation.
func (e *PermissionError) Error() string {
	return fmt.Sprintf("Permission for Action: %v is not granted for Resource: %v", e.Action, e.ResourceName)
}

// FailedCondition - helper function for retrieving underlying failed Condition.
func (e *PermissionError) FailedCondition() Condition {
	if e.ConditionError != nil {
		if conditionErr, ok := e.ConditionError.(*ConditionNotSatisfiedError); ok {
			return conditionErr.condition
		}
	}

	return nil
}

// ConditionNotSatisfiedError - thrown when given Condition for given AccessRequest.
type ConditionNotSatisfiedError struct {
	condition Condition
	request   *AccessRequest
	reason    error
}

// NewConditionNotSatisfiedError - returns new ConditionNotSatisfiedError instance.
func NewConditionNotSatisfiedError(condition Condition, request *AccessRequest, reason error) *ConditionNotSatisfiedError {
	return &ConditionNotSatisfiedError{
		condition: condition,
		request:   request,
		reason:    reason,
	}
}

// Error - error interface implementation.
func (e *ConditionNotSatisfiedError) Error() string {
	return fmt.Sprintf("Condition: \"%v\" was not satisfied! %s", e.condition.Type(), e.reason.Error())
}

// Reason - returns underlying reason (an error) of failing Condition.
func (e *ConditionNotSatisfiedError) Reason() error {
	return e.reason
}

// FailedCondition - returns failed Condition.
func (e *ConditionNotSatisfiedError) FailedCondition() Condition {
	return e.condition
}

// FailedRequest - returns failed AccessRequest.
func (e *ConditionNotSatisfiedError) FailedRequest() *AccessRequest {
	return e.request
}
