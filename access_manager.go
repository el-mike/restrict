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

	roleName := request.Subject.GetRole()
	resourceName := request.Resource.GetResourceName()

	if roleName == "" || resourceName == "" {
		return newRequestMalformedError(request, fmt.Errorf("Missing roleName or resourceName"))
	}

	return am.authorize(request, roleName, resourceName, []string{})
}

// authorize - helper function for decoupling role and resource names retrieval from recursive search.
func (am *AccessManager) authorize(request *AccessRequest, roleName, resourceName string, checkedRoles []string) error {
	role, err := am.policyManager.GetRole(roleName)
	if err != nil {
		return err
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
			return newRequestMalformedError(request, fmt.Errorf("Action cannot be empty"))
		}

		validationError := am.validateAction(grants, action, request)

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
					return newRoleInheritanceCycleError(checkedRoles)
				}

				if err := am.authorize(parentRequest, parent, resourceName, checkedRoles); err != nil {
					// If the returned error is not an access error (meaning something not related
					// to Policy validation happened), we want to return the actual error to the caller.
					// Otherwise, returned error should not override the original one.
					if !am.isAccessError(err) {
						return err
					}
				} else {
					// If .authorize call with parent Role has returned nil,
					// that means the request is satisfied.
					validationError = nil
				}
			}
		}

		// If request has not been granted, abort the loop and return an error.
		if validationError != nil {
			if am.isAccessError(validationError) {
				return newAccessDeniedError(request, action, validationError)
			} else {
				return validationError
			}
		}
	}

	return nil
}

// isAccessError - returns true if error is caused by reasons related to Policy validation
// (Conditions checks, missing Permissions etc.), false otherwise (malformed Policy, reflection errors, etc.)
func (am *AccessManager) isAccessError(err error) bool {
	switch err.(type) {
	case *ConditionNotSatisfiedError,
		*PermissionNotGrantedError,
		*AccessDeniedError:
		return true

	default:
		return false
	}
}

// hasAction - checks if grants list contains a Permission for given action.
func (am *AccessManager) validateAction(permissions []*Permission, action string, request *AccessRequest) error {
	var conditionsCheckError error

	for _, permission := range permissions {
		if permission.Action == action {
			if permission.Conditions == nil || len(permission.Conditions) == 0 || request.SkipConditions {
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

	if conditionsCheckError != nil {
		return conditionsCheckError
	}

	return newPermissionNotGrantedError(action, request.Resource.GetResourceName())
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
