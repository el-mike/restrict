// Package restrict provides an authorization library, with a hybrid of RBAC and ABAC models.
package restrict

import (
	"fmt"
	"github.com/el-mike/restrict/internal/utils"
)

// AccessManager - an entity responsible for checking the authorization. It uses underlying
// PolicyProvider to test an AccessRequest against currently used PolicyDefinition.
type AccessManager struct {
	// PolicyProvider instance, responsible for providing PolicyDefinition.
	policyManager PolicyProvider
}

// NewAccessManager - returns new AccessManager instance.
func NewAccessManager(policyManager PolicyProvider) *AccessManager {
	return &AccessManager{
		policyManager: policyManager,
	}
}

// Authorize - checks if given AccessRequest can be satisfied given currently loaded policy.
// Returns an error if access is not granted or any other problem occurred, nil otherwise.
func (am *AccessManager) Authorize(request *AccessRequest) error {
	if request.Subject == nil || request.Resource == nil {
		return newRequestMalformedError(request, fmt.Errorf("Subject or Resource not defined"))
	}

	roles := request.Subject.GetRoles()
	resourceName := request.Resource.GetResourceName()

	if len(roles) == 0 || resourceName == "" {
		return newRequestMalformedError(request, fmt.Errorf("missing roles or resourceName"))
	}

	allPermissionErrors := PermissionErrors{}

	for _, roleName := range roles {
		permissionErrors, err := am.authorize(request, roleName, resourceName, []string{})

		// If error is not authorization-specific, we return immediately.
		if err != nil {
			return err
		}

		// If AccessRequest is satisfied by a Role, we return immediately.
		if permissionErrors == nil {
			return nil
		}

		// Otherwise, we save it to PermissionErrors, so we can return it to the caller
		// if no Role satisfies the AccessRequest.
		allPermissionErrors = append(allPermissionErrors, permissionErrors...)
	}

	return newAccessDeniedError(request, allPermissionErrors)
}

// authorize - helper function for decoupling role and resource names retrieval from recursive search.
func (am *AccessManager) authorize(request *AccessRequest, roleName, resourceName string, checkedRoles []string) (PermissionErrors, error) {
	role, err := am.policyManager.GetRole(roleName)
	if err != nil {
		return nil, err
	}

	var grants Permissions

	if role.Grants == nil {
		grants = Permissions{}
	} else {
		grants = role.Grants[resourceName]
	}

	parents := role.Parents
	allPermissionErrors := PermissionErrors{}

	for _, action := range request.Actions {
		if action == "" {
			return nil, newRequestMalformedError(request, fmt.Errorf("action cannot be empty"))
		}

		permissionErrors, err := am.validateAction(grants, action, roleName, request)
		// If non-policy related error happened, we return it directly.
		if err != nil {
			return nil, err
		}

		// If access is not granted for given action on current Role, check if
		// any parent Role can satisfy the request.
		if permissionErrors != nil && len(parents) > 0 {
			checkedRoles = append(checkedRoles, roleName)

			for _, parent := range parents {
				parentRequest := &AccessRequest{
					Subject:  request.Subject,
					Resource: request.Resource,
					Actions:  []string{action},
					Context:  request.Context,
				}

				// If parent has already been checked, we want to return an error - otherwise
				// this function will fall into infinite loop.
				if utils.StringSliceContains(checkedRoles, parent) {
					return nil, newRoleInheritanceCycleError(checkedRoles)
				}

				parentPermissionErrors, err := am.authorize(parentRequest, parent, resourceName, checkedRoles)
				if err != nil {
					return nil, err
				}

				// If .authorize call with parent Role has returned nil,
				// that means the request is satisfied.
				// Otherwise, returned errors should not override the original ones.
				if parentPermissionErrors == nil {
					permissionErrors = nil
				}
			}
		}

		// If request has not been granted, abort the loop and return an error,
		// skipping rest of the Actions in the request.
		if permissionErrors != nil {
			allPermissionErrors = append(allPermissionErrors, permissionErrors...)

			// If CompleteValidation is false, we want to return early.
			// Otherwise, loop will continue for the rest of the actions.
			if !request.CompleteValidation {
				return allPermissionErrors, nil
			}
		}
	}

	if len(allPermissionErrors) > 0 {
		return allPermissionErrors, nil
	}

	return nil, nil
}

// validateAction - checks whether a Permission is granted for Action.
func (am *AccessManager) validateAction(permissions []*Permission, action, roleName string, request *AccessRequest) (PermissionErrors, error) {
	permissionErrors := PermissionErrors{}

	for _, permission := range permissions {
		if permission.Action == action {
			// If a Permission with given Action is found, and has no Conditions, access should be granted.
			if len(permission.Conditions) == 0 || request.SkipConditions {
				return nil, nil
			}

			conditionErrors, err := am.checkConditions(permission, request)
			// If non-policy related error happened, we return it directly.
			if err != nil {
				return nil, err
			}

			// If error is nil, Conditions have been satisfied.
			if conditionErrors == nil {
				return nil, nil
			}

			// Otherwise, we add new PermissionError to result slice.
			permissionError := newPermissionError(action, roleName, request.Resource.GetResourceName(), conditionErrors)
			permissionErrors = append(permissionErrors, permissionError)
		}
	}

	// If there are no permissionErrors at this point, this means there was no Permission
	// for given Action (even one with failing Conditions).
	// In such case, a PermissionError with no ConditionErrors is added to reflect that.
	if len(permissionErrors) == 0 {
		permissionError := newPermissionError(action, roleName, request.Resource.GetResourceName(), nil)
		permissionErrors = append(permissionErrors, permissionError)
	}

	return permissionErrors, nil
}

// checkConditions - returns nil if all conditions specified for given actions
// are satisfied, error otherwise.
func (am *AccessManager) checkConditions(permission *Permission, request *AccessRequest) (ConditionErrors, error) {
	if permission.Conditions == nil {
		return nil, nil
	}

	conditionErrors := ConditionErrors{}

	for _, condition := range permission.Conditions {
		if err := condition.Check(request); err != nil {
			// If error returned is ConditionNotSatisfiedError, we add it to the result slice.
			// Otherwise, we want to abort immediately and return it directly.
			if conditionError, ok := err.(*ConditionNotSatisfiedError); ok {
				conditionErrors = append(conditionErrors, conditionError)

				// If CompleteValidation is not enabled, we return first encountered error.
				if !request.CompleteValidation {
					break
				}
			} else {
				return nil, err
			}
		}
	}

	if len(conditionErrors) > 0 {
		return conditionErrors, nil
	}

	return nil, nil
}

// IsAccessError - returns true if error is caused by reasons related to Policy validation
// (Conditions checks, missing Permissions etc.), false otherwise (malformed Policy, reflection errors, etc.)
func (am *AccessManager) IsAccessError(err error) bool {
	switch err.(type) {
	case *ConditionNotSatisfiedError,
		*PermissionError,
		*AccessDeniedError:
		return true

	default:
		return false
	}
}
