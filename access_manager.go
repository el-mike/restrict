// Package restrict provides an authorization library, with a hybrid of RBAC and ABAC models.
package restrict

import "fmt"

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
		return newRequestMalformedError(request, fmt.Errorf("Missing roles or resourceName"))
	}

	accessErrors := PermissionErrors{}

	for _, roleName := range roles {
		permissionDeniedError, err := am.authorize(request, roleName, resourceName, []string{})

		if err != nil {
			return err
		}

		// If AccessRequest is satisfied by a Role, we return immediately.
		if permissionDeniedError == nil {
			return nil
		}

		// Otherwise, we save it to PermissionErrors, so we can return it to the caller
		// if no Role satisfies the AccessRequest.
		accessErrors = append(accessErrors, permissionDeniedError)
	}

	return newAccessDeniedError(request, accessErrors)
}

// authorize - helper function for decoupling role and resource names retrieval from recursive search.
func (am *AccessManager) authorize(request *AccessRequest, roleName, resourceName string, checkedRoles []string) (*PermissionError, error) {
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

	for _, action := range request.Actions {
		if action == "" {
			return nil, newRequestMalformedError(request, fmt.Errorf("Action cannot be empty"))
		}

		validationError := am.validateAction(grants, action, roleName, request)

		// If access is not granted for given action on current Role, check if
		// any parent Role can satisfy the request.
		if validationError != nil && len(parents) > 0 {
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
				if am.isRoleChecked(parent, checkedRoles) {
					return nil, newRoleInheritanceCycleError(checkedRoles)
				}

				permissionDeniedError, err := am.authorize(parentRequest, parent, resourceName, checkedRoles)

				// If the returned error is not an access error (meaning something not related
				// to Policy validation happened), we want to return the actual error to the caller.
				if err != nil {
					return nil, err
				}

				// If .authorize call with parent Role has returned nil,
				// that means the request is satisfied.
				// Otherwise, returned error should not override the original one.
				if permissionDeniedError == nil {
					validationError = nil
				}
			}
		}

		// If request has not been granted, abort the loop and return an error,
		// skipping rest of the Actions in the request.
		if validationError != nil {
			if permissionDeniedError, ok := validationError.(*PermissionError); ok {
				return permissionDeniedError, nil
			}

			// If the error is not PermissionError, error was not specific to Policy validation.
			return nil, validationError
		}
	}

	return nil, nil
}

// validateAction - checks whether a Permission is granted for Action.
func (am *AccessManager) validateAction(permissions []*Permission, action, roleName string, request *AccessRequest) error {
	var conditionsCheckError error

	for _, permission := range permissions {
		if permission.Action == action {
			if len(permission.Conditions) == 0 || request.SkipConditions {
				return nil
			}

			conditionsCheckError = am.checkConditions(permission, request)

			// If error is nil, Conditions have been satisfied. Otherwise, we don't want to
			// break the loop, as other Permissions can have the same action.
			if conditionsCheckError == nil {
				return nil
			}
		}
	}

	// If condition error is not related to Policy validation, return it directly.
	// Otherwise, we wrap it in a PermissionError.
	if conditionsCheckError != nil && !am.isAccessError(conditionsCheckError) {
		return conditionsCheckError
	}

	return newPermissionError(action, roleName, request.Resource.GetResourceName(), conditionsCheckError)
}

// checkConditions - returns nil if all conditions specified for given actions
// are satisfied, error otherwise.
func (am *AccessManager) checkConditions(permission *Permission, request *AccessRequest) error {
	if permission.Conditions == nil {
		return nil
	}

	for _, condition := range permission.Conditions {
		if err := condition.Check(request); err != nil {
			return err
		}
	}

	return nil
}

// isRoleChecked - returns true if role is in parents slice, false otherwise.
func (am *AccessManager) isRoleChecked(role string, parents []string) bool {
	for _, parent := range parents {
		if role == parent {
			return true
		}
	}

	return false
}

// isAccessError - returns true if error is caused by reasons related to Policy validation
// (Conditions checks, missing Permissions etc.), false otherwise (malformed Policy, reflection errors, etc.)
func (am *AccessManager) isAccessError(err error) bool {
	switch err.(type) {
	case *ConditionNotSatisfiedError,
		*PermissionError,
		*AccessDeniedError:
		return true

	default:
		return false
	}
}
