// Package restrict provides an authorization library, with a hybrid of RBAC and ABAC models.
package restrict

import "fmt"

// AccessManager - manages all of the defined Permissions and Roles,
// provides an interface to perform authorization checks and add/remove
// Permissions and Roles.
type AccessManager struct {
	// Instance of PolicyProvider, responsible for providing PolicyDefinition.
	policyManager PolicyProvider
}

// NewAccessManager - returns new AccessManager instance.
func NewAccessManager(policyManager PolicyProvider) *AccessManager {
	return &AccessManager{
		policyManager: policyManager,
	}
}

// Authorize - checks if given AccessRequest can be satisfied given currently loaded policy.
// Returns error if access is not granted or any other problem occurred, nil otherwise.
func (am *AccessManager) Authorize(request *AccessRequest) error {
	if request.Subject == nil || request.Resource == nil {
		return NewRequestMalformedError(request, fmt.Errorf("Subject or Resource not defined"))
	}

	roleName := request.Subject.GetRole()
	resourceName := request.Resource.GetResourceName()

	if roleName == "" || resourceName == "" {
		return NewRequestMalformedError(request, fmt.Errorf("Missing roleName or resourceName"))
	}

	return am.authorize(request, roleName, resourceName, []string{})
}

// authorize - helper function for decoupling role and resource names retrieval from
// recursive search.
func (am *AccessManager) authorize(request *AccessRequest, roleName, resourceName string, checkedRoles []string) error {
	role, err := am.policyManager.GetRole(roleName)
	if err != nil {
		return err
	}

	if role.Grants == nil {
		return NewNoAvailablePermissionsError(role.ID)
	}

	grants := role.Grants[resourceName]
	parents := role.Parents

	// If given role has no permissions granted, and no parents to
	// fall back on, return an error.
	if len(grants) == 0 && len(parents) == 0 {
		return NewNoAvailablePermissionsError(role.ID)
	}

	for _, action := range request.Actions {
		if action == "" {
			return NewRequestMalformedError(request, fmt.Errorf("Action cannot be empty"))
		}

		authorizeError := am.validateAction(grants, action, request)

		// If access is not granted for given action on current Role, check if
		// any parent Role can satisfy the request.
		if authorizeError != nil && len(parents) > 0 {
			checkedRoles = append(checkedRoles, roleName)

			for _, parent := range parents {
				parentRequest := &AccessRequest{
					Resource: request.Resource,
					Actions:  []string{action},
					Context:  request.Context,
				}

				// If parent has already been checked, we want to return an error - otherwise
				// this function will fall into infinite loop.
				if am.isChecked(parent, checkedRoles) {
					return NewRoleInheritanceCycleError(checkedRoles)
				}

				if err := am.authorize(parentRequest, parent, resourceName, checkedRoles); err != nil {
					// If the returned error is not an access error (meaning something not related
					// to Policy validation happened), we want to return the actual error to the caller.
					// Otherwise, returned error should not override the original one.
					if !am.isAccessError(err) {
						return err
					}
				} else {
					// If .Authorize call with parent Role has returned nil,
					// that means the request is satisfied.
					authorizeError = nil
				}
			}
		}

		// If request has not been granted, abort the loop and return an error.
		if authorizeError != nil {
			return NewAccessDeniedError(request, action, authorizeError)
		}
	}

	return nil
}

// isAccessError - returns true if error is caused by reasons related to Policy validation
// (Conditions checks, missing Permissions etc.), false otherwise (malformed Policy, reflection errors, etc.)
func (am *AccessManager) isAccessError(err error) bool {
	switch err.(type) {
	case *ConditionNotSatisfiedError,
		*NoAvailablePermissionsError,
		*PermissionNotGrantedError,
		*AccessDeniedError:
		return true

	default:
		return false
	}
}

// hasAction - checks if grants list contains given action.
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

	return NewPermissionNotGrantedError(action, request.Resource.GetResourceName())
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

// isChecked - returns true if role is in parents slice, false otherwise.
func (am *AccessManager) isChecked(role string, parents []string) bool {
	for _, parent := range parents {
		if role == parent {
			return true
		}
	}

	return false
}
