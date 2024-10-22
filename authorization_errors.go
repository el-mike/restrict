package restrict

import (
	"fmt"
	"github.com/el-mike/restrict/internal/utils"
	"strings"
)

// PermissionErrors - an alias type for a slice of PermissionError, with extra helper methods.
type PermissionErrors []*PermissionError

// First - returns the first PermissionError encountered when performing authorization.
// Especially helpful when AccessRequest was set to fail early.
func (ae PermissionErrors) First() *PermissionError {
	if len(ae) == 0 {
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

	return result
}

// GetFailedActions - returns all Actions for which access was denied.
func (ae PermissionErrors) GetFailedActions() []string {
	actions := []string{}

	for _, e := range ae {
		if !utils.StringSliceContains(actions, e.Action) {
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
	preparedActions := []string{}

	for _, action := range e.Errors.GetFailedActions() {
		preparedActions = append(preparedActions, fmt.Sprintf("\"%s\"", action))
	}

	actionsNoun := "Action"
	if len(preparedActions) > 1 {
		actionsNoun = "Actions"
	}

	return fmt.Sprintf("access denied for %s: %s on Resource: \"%s\"", actionsNoun, strings.Join(preparedActions, ", "), e.Request.Resource.GetResourceName())
}

// ConditionErrors - an alias type for a slice of ConditionNotSatisfiedError.
type ConditionErrors []*ConditionNotSatisfiedError

// First - returns the first ConditionErrors encountered when validating given Action.
// Especially helpful when AccessRequest was set to fail early.
func (ce ConditionErrors) First() *ConditionNotSatisfiedError {
	if len(ce) == 0 {
		return nil
	}

	return ce[0]
}

// PermissionError - thrown when Permission is not granted for a given Action.
type PermissionError struct {
	Action          string
	RoleName        string
	ResourceName    string
	ConditionErrors ConditionErrors
}

// newPermissionError - returns new PermissionError instance.
func newPermissionError(action, roleName, resourceName string, conditionErrors ConditionErrors) *PermissionError {
	return &PermissionError{
		Action:          action,
		RoleName:        roleName,
		ResourceName:    resourceName,
		ConditionErrors: conditionErrors,
	}
}

// Error - error interface implementation.
func (e *PermissionError) Error() string {
	if len(e.ConditionErrors) > 0 {
		return fmt.Sprintf("Permission for Action: \"%v\" is not granted for Resource: \"%v\" due to failed Conditions", e.Action, e.ResourceName)
	}

	return fmt.Sprintf("Permission for Action: \"%v\" is not granted for Resource: \"%v\"", e.Action, e.ResourceName)
}

// ConditionNotSatisfiedError - thrown when given Condition for given AccessRequest.
type ConditionNotSatisfiedError struct {
	Condition Condition
	Request   *AccessRequest
	Reason    error
}

// NewConditionNotSatisfiedError - returns new ConditionNotSatisfiedError instance.
func NewConditionNotSatisfiedError(condition Condition, request *AccessRequest, reason error) *ConditionNotSatisfiedError {
	return &ConditionNotSatisfiedError{
		Condition: condition,
		Request:   request,
		Reason:    reason,
	}
}

// Error - error interface implementation.
func (e *ConditionNotSatisfiedError) Error() string {
	return fmt.Sprintf("Condition: \"%v\" was not satisfied, reason: %s", e.Condition.Type(), e.Reason.Error())
}
